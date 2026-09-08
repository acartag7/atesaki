#!/usr/bin/env bash
# Review aid for public pushes. Pattern matching cannot identify every private
# name: the author must still inspect the complete public diff.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)"

scan() {
  local category="$1" pattern="$2"
  shift 2
  local status=0
  # Print filenames only: never echo a matched identifier into CI logs.
  git grep -Il -E -e "$pattern" -- "$@" || status=$?
  case "$status" in
    0) printf 'sanitization refused: %s\n' "$category" >&2; return 1 ;;
    1) return 0 ;;
    *) printf 'sanitization scan failed: %s\n' "$category" >&2; return "$status" ;;
  esac
}

scan home-path '(/Users/[^/[:space:]]+/|/home/[^/[:space:]]+/)' . ':!tools/sanitize.sh'
scan internal-host '([[:alnum:]-]+[.])+(corp|lan|local)([:/[:space:]"`]|$)' . ':!tools/sanitize.sh'
scan directory-name '(CN|OU|DC)=[[:alnum:]_-]+' . ':!tools/sanitize.sh'
scan vault-path '(vault://[^[:space:]]+|secret/data/[^[:space:]]+)' . ':!tools/sanitize.sh'
scan employer '(Swiss[ -]?Post|swisspost|[.]post[.]ch)' . ':!tools/sanitize.sh'
# UUIDs are valid tenant/group examples only when explicitly recognizable as
# placeholders. Report other UUIDs for review without printing their values.
python3 - <<'PYTHON'
import pathlib
import re
import subprocess
import sys

uuid = re.compile(rb"[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}")
placeholders = {b"00000000-0000-0000-0000-000000000000", b"11111111-1111-1111-1111-111111111111"}
# Existing evidence citations only. Other docs receive no exception.
evidence_notes = {
    "docs/deltas.md": b"`atesaki/evidence/prm-probe-2026-09-01/`",
    "docs/open-questions.md": b"`atesaki/evidence/prm-probe-2026-09-01/`",
    "docs/contract.md": b"`atesaki/evidence/migration-playbook-mcp-setup-pain-2026-08-26.md`",
}
paths = subprocess.check_output(["git", "ls-files", "-z"]).split(b"\0")
failed = False
for raw in filter(None, paths):
    path = pathlib.Path(raw.decode())
    if path.as_posix() == "tools/sanitize.sh":
        continue
    data = path.read_bytes()
    citation = evidence_notes.get(path.as_posix())
    evidence_data = data.replace(citation, b"") if citation else data
    if b"atesaki/evidence/" in evidence_data:
        print(f"sanitization refused: private evidence path in {path}", file=sys.stderr)
        failed = True
    if any(value not in placeholders for value in uuid.findall(data)):
        print(f"sanitization refused: tenant/group identifier in {path}", file=sys.stderr)
        failed = True
sys.exit(int(failed))
PYTHON
printf 'sanitization passed\n'
