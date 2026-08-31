#!/usr/bin/env python3
"""Generate and check Atesaki's contract manifest."""
from __future__ import annotations

import argparse
import hashlib
import json
import pathlib
import subprocess
import sys
from typing import Any


ROOT = pathlib.Path(__file__).resolve().parents[1]
MANIFEST = ROOT / "CONTRACT-MANIFEST.json"
CONTRACT_LABEL = "contract-change"
GUARDED_FILES = (
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
GUARDED_PATHS = tuple(
    sorted((*GUARDED_FILES, *(f"{path}/**" for path in GUARDED_DIRECTORIES)))
)


def contract_files() -> list[pathlib.Path]:
    files: list[pathlib.Path] = []
    for name in GUARDED_FILES:
        path = ROOT / name
        if not path.exists() and not path.is_symlink():
            continue
        if path.is_symlink() or not path.is_file():
            raise ValueError(f"{name}: guarded path must be a regular file")
        files.append(path)

    for name in GUARDED_DIRECTORIES:
        directory = ROOT / name
        if not directory.exists() and not directory.is_symlink():
            continue
        if directory.is_symlink() or not directory.is_dir():
            raise ValueError(f"{name}: guarded directory must be a real directory")
        for path in directory.rglob("*"):
            if path.is_symlink():
                raise ValueError(f"{path.relative_to(ROOT)}: symlink in guarded directory")
            if path.is_file():
                files.append(path)
    return sorted(files)


def file_hash(path: pathlib.Path) -> str:
    digest = hashlib.sha256()
    with path.open("rb") as source:
        for block in iter(lambda: source.read(1024 * 1024), b""):
            digest.update(block)
    return digest.hexdigest()


def manifest_data() -> dict[str, Any]:
    return {
        "version": 1,
        "algorithm": "sha256",
        "guarded_paths": list(GUARDED_PATHS),
        "files": {
            path.relative_to(ROOT).as_posix(): file_hash(path)
            for path in contract_files()
        },
    }


def manifest_text(data: dict[str, Any]) -> str:
    return json.dumps(data, indent=2, sort_keys=True) + "\n"


def write_manifest() -> int:
    MANIFEST.write_text(manifest_text(manifest_data()))
    print(f"wrote {MANIFEST.name}")
    return 0


def check_manifest() -> int:
    expected = manifest_data()
    try:
        text = MANIFEST.read_text()
        actual = json.loads(text)
    except (OSError, json.JSONDecodeError) as error:
        raise ValueError(f"{MANIFEST.name}: {error}") from error
    if actual == expected and text == manifest_text(expected):
        print(f"{MANIFEST.name}: matches the guarded files")
        return 0

    actual_files = actual.get("files", {}) if isinstance(actual, dict) else {}
    if not isinstance(actual_files, dict):
        actual_files = {}
    changed = sorted(
        name
        for name in set(actual_files) | set(expected["files"])
        if actual_files.get(name) != expected["files"].get(name)
    )
    detail = ", ".join(changed) if changed else "manifest structure"
    raise ValueError(
        f"{MANIFEST.name} is stale for {detail}; run python3 tools/contract-manifest.py"
    )


def guarded(path: str) -> bool:
    return path in GUARDED_FILES or any(
        path.startswith(f"{directory}/") for directory in GUARDED_DIRECTORIES
    )


def pull_request_has_label(event_path: pathlib.Path) -> bool:
    try:
        labels = json.loads(event_path.read_text())["pull_request"]["labels"]
    except (OSError, json.JSONDecodeError, KeyError, TypeError) as error:
        raise ValueError(f"cannot read pull request labels: {error}") from error
    if not isinstance(labels, list):
        raise ValueError("pull request labels must be an array")
    return any(
        isinstance(label, dict) and label.get("name") == CONTRACT_LABEL
        for label in labels
    )


def check_pull_request(base: str, head: str, event_path: pathlib.Path) -> int:
    diff = subprocess.run(
        ("git", "diff", "--name-only", "-z", "--no-renames", base, head, "--"),
        cwd=ROOT,
        check=False,
        capture_output=True,
    )
    if diff.returncode:
        raise ValueError(diff.stderr.decode(errors="replace").strip() or "git diff failed")
    changed = sorted(
        path.decode(errors="surrogateescape")
        for path in diff.stdout.split(b"\0")
        if path and guarded(path.decode(errors="surrogateescape"))
    )
    if changed and not pull_request_has_label(event_path):
        raise ValueError(
            f"guarded files changed without the {CONTRACT_LABEL!r} label: "
            + ", ".join(changed)
        )
    print("contract-change visibility: pass")
    return 0


def parse_arguments() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    mode = parser.add_mutually_exclusive_group()
    mode.add_argument("--check", action="store_true")
    mode.add_argument("--check-pr", action="store_true")
    parser.add_argument("--base")
    parser.add_argument("--head")
    parser.add_argument("--event", type=pathlib.Path)
    return parser.parse_args()


def main() -> int:
    arguments = parse_arguments()
    try:
        if arguments.check:
            return check_manifest()
        if arguments.check_pr:
            if not arguments.base or not arguments.head or not arguments.event:
                raise ValueError("--check-pr requires --base, --head, and --event")
            return check_pull_request(arguments.base, arguments.head, arguments.event)
        return write_manifest()
    except (OSError, ValueError) as error:
        print(f"contract guard: {error}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
