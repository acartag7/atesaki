#!/usr/bin/env python3
"""Generate and enforce Atesaki's contract manifest.

Run without arguments to rewrite CONTRACT-MANIFEST.json. Use --check to compare the
committed manifest with the current tree. The PR and mutation-test modes implement
the contract-change visibility gate used by CI.
"""
from __future__ import annotations

import argparse
import hashlib
import json
import os
import pathlib
import re
import shutil
import stat
import subprocess
import sys
import tempfile
from collections.abc import Iterable
from typing import Any


ROOT = pathlib.Path(__file__).resolve().parents[1]
MANIFEST_NAME = "CONTRACT-MANIFEST.json"
CONTRACT_LABEL = "contract-change"

EXACT_GUARDED_PATHS = (
    ".github/workflows/contract-guard.yml",
    "docs/contract-boundaries.md",
    "docs/contract-grants.md",
    "docs/contract.md",
    "docs/decisions.md",
    "docs/deltas.md",
    "docs/negative-matrix.md",
    "docs/quality-bar.md",
    "docs/threat-model.md",
    "tools/contract-lint.py",
    "tools/contract-manifest.py",
    "tools/schema-check.py",
)
GUARDED_DIRECTORIES = ("fixtures", "schema")
GUARDED_DECLARATIONS = tuple(
    sorted((*EXACT_GUARDED_PATHS, *(f"{path}/**" for path in GUARDED_DIRECTORIES)))
)
CONTRACT_CHANGE_PREFIXES = ("docs/", "fixtures/", "prompts/", "schema/", "tools/")


class ManifestError(Exception):
    """A deterministic, user-facing manifest or gate failure."""


def display_path(path: str) -> str:
    """Return a bounded, single-line representation for CI logs."""
    if len(path) > 4096:
        path = path[:4096] + "..."
    return json.dumps(path, ensure_ascii=True)


def relative_path(path: pathlib.Path, root: pathlib.Path) -> str:
    return path.relative_to(root).as_posix()


def require_regular_file(path: pathlib.Path, root: pathlib.Path) -> None:
    mode = path.lstat().st_mode
    if not stat.S_ISREG(mode):
        raise ManifestError(f"{display_path(relative_path(path, root))}: guarded path is not a regular file")


def guarded_files(root: pathlib.Path) -> list[pathlib.Path]:
    files: list[pathlib.Path] = []
    for name in EXACT_GUARDED_PATHS:
        path = root / name
        if path.exists() or path.is_symlink():
            require_regular_file(path, root)
            files.append(path)

    for name in GUARDED_DIRECTORIES:
        directory = root / name
        if not directory.exists() and not directory.is_symlink():
            continue
        if directory.is_symlink() or not directory.is_dir():
            raise ManifestError(f"{display_path(name)}: guarded directory is not a real directory")
        for current, directory_names, file_names in os.walk(directory, followlinks=False):
            current_path = pathlib.Path(current)
            for child_name in directory_names:
                child = current_path / child_name
                if child.is_symlink():
                    raise ManifestError(
                        f"{display_path(relative_path(child, root))}: symlink inside guarded directory"
                    )
            for child_name in file_names:
                child = current_path / child_name
                require_regular_file(child, root)
                files.append(child)
    return sorted(files, key=lambda path: relative_path(path, root))


def hash_file(path: pathlib.Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as source:
        for block in iter(lambda: source.read(1024 * 1024), b""):
            digest.update(block)
    return digest.hexdigest()


def manifest_data(root: pathlib.Path) -> dict[str, Any]:
    return {
        "version": 1,
        "algorithm": "sha256",
        "guarded_paths": list(GUARDED_DECLARATIONS),
        "files": {
            relative_path(path, root): hash_file(path)
            for path in guarded_files(root)
        },
    }


def manifest_bytes(root: pathlib.Path) -> bytes:
    return (json.dumps(manifest_data(root), indent=2, sort_keys=True) + "\n").encode()


def write_manifest(root: pathlib.Path) -> None:
    destination = root / MANIFEST_NAME
    payload = manifest_bytes(root)
    descriptor, temporary_name = tempfile.mkstemp(prefix=f".{MANIFEST_NAME}.", dir=root)
    temporary = pathlib.Path(temporary_name)
    try:
        os.fchmod(descriptor, 0o644)
        with os.fdopen(descriptor, "wb") as target:
            target.write(payload)
            target.flush()
            os.fsync(target.fileno())
        os.replace(temporary, destination)
    finally:
        if temporary.exists():
            temporary.unlink()
    print(f"wrote {MANIFEST_NAME} ({len(manifest_data(root)['files'])} files)")


def reject_duplicate_keys(pairs: list[tuple[str, Any]]) -> dict[str, Any]:
    result: dict[str, Any] = {}
    for key, value in pairs:
        if key in result:
            raise ManifestError(f"{MANIFEST_NAME}: duplicate JSON key {key!r}")
        result[key] = value
    return result


def load_manifest(root: pathlib.Path) -> tuple[dict[str, Any], bytes]:
    path = root / MANIFEST_NAME
    try:
        payload = path.read_bytes()
    except FileNotFoundError as error:
        raise ManifestError(f"{MANIFEST_NAME}: file is missing") from error
    try:
        data = json.loads(payload, object_pairs_hook=reject_duplicate_keys)
    except (json.JSONDecodeError, UnicodeDecodeError) as error:
        raise ManifestError(f"{MANIFEST_NAME}: invalid JSON: {error}") from error
    if not isinstance(data, dict):
        raise ManifestError(f"{MANIFEST_NAME}: top level must be an object")
    return data, payload


def check_manifest(root: pathlib.Path) -> list[str]:
    try:
        actual, actual_bytes = load_manifest(root)
        expected = manifest_data(root)
    except (ManifestError, OSError) as error:
        return [str(error)]

    errors: list[str] = []
    if set(actual) != set(expected):
        errors.append(f"{MANIFEST_NAME}: expected exactly the keys {sorted(expected)}")
    if actual.get("version") != expected["version"]:
        errors.append(f"{MANIFEST_NAME}: version mismatch")
    if actual.get("algorithm") != expected["algorithm"]:
        errors.append(f"{MANIFEST_NAME}: algorithm mismatch")
    if actual.get("guarded_paths") != expected["guarded_paths"]:
        errors.append(f"{MANIFEST_NAME}: guarded path declaration mismatch")

    actual_files = actual.get("files")
    if not isinstance(actual_files, dict):
        errors.append(f"{MANIFEST_NAME}: files must be an object")
    else:
        for name in sorted(set(actual_files) | set(expected["files"])):
            if name not in actual_files:
                errors.append(f"{display_path(name)}: missing from {MANIFEST_NAME}")
            elif name not in expected["files"]:
                errors.append(f"{display_path(name)}: guarded file is missing from the tree")
            elif not isinstance(actual_files[name], str) or not re.fullmatch(r"[0-9a-f]{64}", actual_files[name]):
                errors.append(f"{display_path(name)}: invalid SHA-256 in {MANIFEST_NAME}")
            elif actual_files[name] != expected["files"][name]:
                errors.append(f"{display_path(name)}: SHA-256 mismatch")

    if not errors and actual_bytes != manifest_bytes(root):
        errors.append(f"{MANIFEST_NAME}: content is not the canonical regeneration")
    return errors


def is_guarded(path: str) -> bool:
    return path in EXACT_GUARDED_PATHS or any(path.startswith(f"{name}/") for name in GUARDED_DIRECTORIES)


def allowed_in_contract_change(path: str) -> bool:
    return path == MANIFEST_NAME or is_guarded(path) or path.startswith(CONTRACT_CHANGE_PREFIXES)


def git(root: pathlib.Path, *arguments: str) -> str:
    completed = subprocess.run(
        ("git", *arguments), cwd=root, check=False, capture_output=True, text=True
    )
    if completed.returncode:
        detail = completed.stderr.strip() or completed.stdout.strip() or "git failed"
        raise ManifestError(detail)
    return completed.stdout


def changed_paths(root: pathlib.Path, base: str, head: str) -> list[str]:
    completed = subprocess.run(
        ("git", "diff", "--name-only", "-z", "--no-renames", base, head, "--"),
        cwd=root,
        check=False,
        capture_output=True,
    )
    if completed.returncode:
        detail = completed.stderr.decode(errors="replace").strip() or "git diff failed"
        raise ManifestError(detail)
    return sorted({item.decode(errors="surrogateescape") for item in completed.stdout.split(b"\0") if item})


def event_has_contract_label(event_path: pathlib.Path) -> bool:
    try:
        event = json.loads(event_path.read_text())
        labels = event["pull_request"]["labels"]
    except (OSError, json.JSONDecodeError, KeyError, TypeError) as error:
        raise ManifestError(f"GitHub event: cannot read pull_request.labels: {error}") from error
    if not isinstance(labels, list):
        raise ManifestError("GitHub event: pull_request.labels must be an array")
    return any(isinstance(label, dict) and label.get("name") == CONTRACT_LABEL for label in labels)


def check_pr(root: pathlib.Path, base: str, head: str, contract_change: bool) -> list[str]:
    try:
        changed = changed_paths(root, base, head)
    except ManifestError as error:
        return [f"PR diff: {error}"]
    guarded = [path for path in changed if is_guarded(path)]
    if not guarded:
        return []
    if not contract_change:
        return [
            f"{display_path(path)}: guarded path changed without the {CONTRACT_LABEL!r} label"
            for path in guarded
        ]
    return [
        f"{display_path(path)}: a {CONTRACT_LABEL!r} PR may change only guarded paths, "
        f"{MANIFEST_NAME}, or files under docs/, fixtures/, prompts/, schema/, and tools/"
        for path in changed
        if not allowed_in_contract_change(path)
    ]


def tracked_and_untracked_files(root: pathlib.Path) -> Iterable[str]:
    completed = subprocess.run(
        ("git", "ls-files", "--cached", "--others", "--exclude-standard", "-z"),
        cwd=root,
        check=True,
        capture_output=True,
    )
    return (item.decode(errors="surrogateescape") for item in completed.stdout.split(b"\0") if item)


def copy_test_tree(source: pathlib.Path, destination: pathlib.Path) -> None:
    for name in tracked_and_untracked_files(source):
        original = source / name
        target = destination / name
        target.parent.mkdir(parents=True, exist_ok=True)
        if original.is_symlink():
            target.symlink_to(os.readlink(original))
        else:
            shutil.copy2(original, target)


def commit(root: pathlib.Path, message: str) -> str:
    git(root, "add", "-A")
    git(root, "commit", "--quiet", "-m", message)
    return git(root, "rev-parse", "HEAD").strip()


def run_tool(root: pathlib.Path, *arguments: str) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        (sys.executable, "tools/contract-manifest.py", *arguments),
        cwd=root,
        check=False,
        capture_output=True,
        text=True,
    )


def run_mutation_test(root: pathlib.Path) -> list[str]:
    errors = check_manifest(root)
    if errors:
        return ["mutation setup: " + error for error in errors]
    with tempfile.TemporaryDirectory(prefix="atesaki-contract-guard-") as temporary_name:
        test_root = pathlib.Path(temporary_name) / "repo"
        test_root.mkdir()
        copy_test_tree(root, test_root)
        git(test_root, "init", "--quiet")
        git(test_root, "config", "user.name", "Contract Guard Test")
        git(test_root, "config", "user.email", "contract-guard@example.invalid")
        base = commit(test_root, "test: mutation base")

        contract = test_root / "docs" / "contract.md"
        original = contract.read_text()
        marker = "**Status: DRAFT. Freeze happens only when Arnold says the word.**"
        if marker not in original:
            return [f"mutation setup: {display_path('docs/contract.md')}: status marker not found"]
        contract.write_text(
            original.replace(
                marker,
                "**Status: DRAFT. Freeze happens only after the mutation test.**",
                1,
            )
        )
        unlabeled_head = commit(test_root, "test: mutate contract without label")
        event_path = test_root.parent / "event.json"
        event_path.write_text(json.dumps({"pull_request": {"labels": []}}))
        unlabeled = run_tool(
            test_root,
            "--check-pr",
            "--base",
            base,
            "--head",
            unlabeled_head,
            "--event",
            str(event_path),
        )
        if (
            unlabeled.returncode != 1
            or "docs/contract.md" not in unlabeled.stderr
            or CONTRACT_LABEL not in unlabeled.stderr
        ):
            return ["mutation unlabeled: expected exit 1 with the named contract-change label rule"]
        stale_manifest = run_tool(test_root, "--check")
        if stale_manifest.returncode != 1 or "docs/contract.md" not in stale_manifest.stderr:
            return ["mutation manifest: expected exit 1 naming docs/contract.md"]
        print("mutation unlabeled: expected failure")
        for line in unlabeled.stderr.strip().splitlines():
            print(" -", line)

        regenerated = run_tool(test_root)
        if regenerated.returncode:
            return ["mutation relock: " + (regenerated.stdout + regenerated.stderr).strip()]
        labeled_head = commit(test_root, "test: relock contract manifest")
        event_path.write_text(
            json.dumps({"pull_request": {"labels": [{"name": CONTRACT_LABEL}]}})
        )
        labeled = run_tool(
            test_root,
            "--check-pr",
            "--base",
            base,
            "--head",
            labeled_head,
            "--event",
            str(event_path),
        )
        manifest_check = run_tool(test_root, "--check")
        labeled_errors: list[str] = []
        if labeled.returncode:
            labeled_errors.append((labeled.stdout + labeled.stderr).strip())
        if manifest_check.returncode:
            labeled_errors.append((manifest_check.stdout + manifest_check.stderr).strip())
        lint = subprocess.run(
            (sys.executable, "tools/contract-lint.py"), cwd=test_root, check=False, capture_output=True, text=True
        )
        if lint.returncode:
            labeled_errors.append("tools/contract-lint.py: " + (lint.stdout + lint.stderr).strip())
        if labeled_errors:
            return ["mutation labeled: " + error for error in labeled_errors]

        readme = test_root / "README.md"
        readme.write_text(readme.read_text() + "\ncontract guard scope mutation\n")
        outside_head = commit(test_root, "test: add an out-of-scope change")
        outside = run_tool(
            test_root,
            "--check-pr",
            "--base",
            base,
            "--head",
            outside_head,
            "--event",
            str(event_path),
        )
        if outside.returncode != 1 or "README.md" not in outside.stderr:
            return ["mutation scope: expected exit 1 naming out-of-scope README.md"]
        print("mutation labeled with regenerated manifest: pass")
    return []


def print_errors(errors: Iterable[str]) -> int:
    errors = list(errors)
    if not errors:
        return 0
    for error in errors:
        print("contract guard:", error, file=sys.stderr)
    return 1


def parse_arguments() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    mode = parser.add_mutually_exclusive_group()
    mode.add_argument("--check", action="store_true", help="check the committed manifest")
    mode.add_argument("--check-pr", action="store_true", help="check a PR-shaped Git diff")
    mode.add_argument("--mutation-test", action="store_true", help="prove both contract-label directions")
    parser.add_argument("--base", help="base commit for --check-pr")
    parser.add_argument("--head", help="head commit for --check-pr")
    parser.add_argument("--event", type=pathlib.Path, help="GitHub event JSON containing PR labels")
    parser.add_argument("--contract-change", action="store_true", help=f"act as if {CONTRACT_LABEL!r} is present")
    return parser.parse_args()


def main() -> int:
    arguments = parse_arguments()
    try:
        if arguments.check:
            errors = check_manifest(ROOT)
            if not errors:
                print(f"{MANIFEST_NAME}: matches {len(manifest_data(ROOT)['files'])} guarded files")
            return print_errors(errors)
        if arguments.check_pr:
            if not arguments.base or not arguments.head:
                raise ManifestError("--check-pr requires --base and --head")
            has_label = arguments.contract_change
            if arguments.event:
                has_label = has_label or event_has_contract_label(arguments.event)
            errors = check_pr(ROOT, arguments.base, arguments.head, has_label)
            if not errors:
                print("contract PR policy: pass")
            return print_errors(errors)
        if arguments.mutation_test:
            return print_errors(run_mutation_test(ROOT))
        write_manifest(ROOT)
        return 0
    except (ManifestError, OSError, subprocess.SubprocessError) as error:
        return print_errors([str(error)])


if __name__ == "__main__":
    sys.exit(main())
