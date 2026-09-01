MODEL: grok-4.5   EFFORT: high   FALLBACK: gpt-5.6-terra   TOOL: Codex CLI or pi, fresh clone
PRECONDITION: slice 1 merged (its runner already executes the mcp-sso §8 portable
subset); the mcp-sso corpus pinned at an exact commit + `MANIFEST.json` hash, and the
**exact frozen portable fixture-ID set** this build will pass listed in the PR before
any code — acceptance is by ID, not section name; the Atesaki delta fixtures for
D1/D3/D4 locked.
WHY: Atesaki's authorization-server half proves itself against the shared corpus —
parity by fixture, not by memory.

Read first, fully: ~/project/mcp-sso/docs/contracts/ §01–§11, §14, §17 (identity
ports), §19 at the SHA pinned in docs/contract.md · mcp-sso/fixtures/** (portable only)
· docs/contract.md §3, §4, §8, §13 · docs/contract-boundaries.md B4 · docs/deltas.md
(D1, D3, D4, D5, D6, D7, D11, D12, D13 change what you would otherwise inherit — read each
before implementing the clause it touches).

SCOPE — build exactly this:
- Extend the slice-1 Atesaki runner to the **full** frozen portable mcp-sso set (it
  already runs the §8 subset):
  loads a fixture, composes Atesaki through its ports (clock, seeded randomness, keys,
  recorded outbound HTTP — an unrecorded outbound call fails the fixture), sends the
  request, compares exactly. A runner contains no expectations of its own. A skipped
  frozen portable fixture is a failure.
- Authorization-server core: AS metadata (origin) + per-route PRM (D1); stateless DCR;
  CIMD from vendored documents only (no live fetch); authorize with the §9.3 steps
  1–4 exactly; then, in this slice, the policy step runs with **no rules configured,
  so every request escalates — which IS the contract** (escalate by default): insert
  `grant_request` + `preapproval` per A3, return `approval_pending` + `request_id`;
  with no approver CLI yet, every pre-approval expires via A14. Contract-conformant,
  fail-closed, and the A3/A3″/A14 fixtures pass for real — no invented interim state. Consent, approve, code exchange, and refresh rotation are **not in this
  slice** — they are grant-coupled (D3/D6/D7/D12, G8) and arrive with slice 3; the
  portable fixtures that exercise them are deferred, not skipped silently (list every
  deferred fixture id in the PR).
- Identity ports: Entra (id_token, redirect flow), generic OIDC, header-mode signed
  assertion per B4 in full, console pairing loopback-only. Subject boundary per
  mcp-sso §6.5.
- Store port (interface) with the **memory adapter** only in this slice, used only by
  the runner and tests; SQLite arrives in slice 3 with the conformance suite.

HARD RULES: as prompts/README.md. Additionally: every place Atesaki deviates from an
mcp-sso clause MUST correspond to a row in docs/deltas.md — a deviation without a row
is a bug, stop and report. Never edit a frozen fixture. Pin the corpus version and
manifest hash you ran against in the PR.

VERIFY: the pinned portable fixture-ID set — zero skips; Atesaki fixtures for §4
rungs, D1, never 6, and the A3/A3″/A14 escalation path; a real MCP client (Claude Code or Codex CLI) completing discovery, registration, and
the console identity leg to a real `approval_pending` + `request_id` — observed and
named with the client version. The full sign-in-to-tool-call proof
belongs to slice 3, when grants exist.

DONE WHEN: parity status line published in the PR (corpus version passed, current,
frozen fixtures not yet passed with reasons); real-client sign-in shown;
implementation checks green.

REPORT: parity line; deferred fixtures and why; every contract gap; divergences from
mcp-sso you found that have no delta row (these are the most important findings).
