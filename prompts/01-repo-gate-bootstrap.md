MODEL: gpt-5.6-sol   EFFORT: high   TOOL: Codex CLI in ~/project/atesaki-core
PRECONDITION: Arnold has said "commit". The repo has zero commits until he does.
WHY: the contract-change visibility gate (docs/quality-bar.md, open question #30) is
policy prose until this lands. No implementation PR may open before it exists.

Read first, fully: docs/quality-bar.md ("Who may change the contract", "Change
protocol") · docs/open-questions.md #30 · tools/contract-lint.py · README.md.

DELIVERABLES
1. First commit: the current tree, exactly as is (no edits to any doc), message
   `chore: contract set v0 draft — nothing frozen`. Tag nothing yet.
2. `tools/contract-manifest.py`: writes/checks `CONTRACT-MANIFEST.json` — SHA-256 of
   every guarded path. Guarded paths, named explicitly: `docs/contract.md`,
   `docs/contract-grants.md`, `docs/contract-boundaries.md`, `docs/deltas.md`,
   `docs/decisions.md`, `docs/threat-model.md`, `docs/negative-matrix.md` (once 04
   lands), `docs/quality-bar.md`, `fixtures/**` (once 03 lands), `schema/**` (once 02
   lands), **and the gate itself**: `.github/workflows/contract-guard.yml`,
   `tools/contract-manifest.py`, `tools/contract-lint.py`, `tools/schema-check.py` — a
   gate that does not guard its own enforcement can be weakened by any implementation
   PR. **`CONTRACT-MANIFEST.json` itself is excluded from its own hash list** (a
   manifest that hashes itself never stabilizes); the workflow instead verifies the
   committed manifest equals a fresh regeneration.
   `--check` exits 1 on any mismatch and names the file.
3. `.github/workflows/contract-guard.yml`: on every PR — run `tools/contract-lint.py`
   and `tools/contract-manifest.py --check`. A PR may change guarded paths ONLY if it
   carries the label `contract-change` AND touches no file outside: the guarded set,
   `tools/`, `schema/`, `docs/`, `prompts/`, and `CONTRACT-MANIFEST.json` — so the
   contract packets 02 (adds `tools/schema-check.py`) and 04 (adds
   `docs/negative-matrix.md`) are admissible by construction. Any other PR that changes a guarded
   path fails with a message naming the file and the rule.
4. A **mutation test** in CI proving the gate works: a scripted PR-shaped diff that
   edits one line of `docs/contract.md` without the label must fail; the same edit with
   the label **plus its regenerated manifest** must pass. Both directions asserted.
5. The relock path: a contract PR regenerates `CONTRACT-MANIFEST.json`; the workflow
   verifies the regenerated manifest matches the tree. Arnold's approval of the PR is
   the owner relock — no separate ceremony.

HARD RULES: do not edit any contract page (the first commit is the tree as-is). Do
not add product code. Do not touch `seed/` beyond including it in the commit. Do not
push to `main` — branch + PR, Arnold merges.

DONE WHEN: CI green on the PR; the mutation test demonstrably fails the unlabeled
edit; `python3 tools/contract-lint.py` and `--check` both pass on the branch head.

REPORT: exact guarded set; the mutation test output for both directions; anything
in quality-bar.md that the implemented gate cannot enforce as written (contract gap).
