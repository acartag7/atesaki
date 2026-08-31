# Quality bar — how this repo works

**Status: DRAFT.**

## Order of work

Contract → threat model → tests → code. Acceptance tests are written from the contract
and may run red before an implementation exists — that is normal, not a defect.

## Who may change the contract `[O:2026-08-30]`

The rule, in Arnold's words: **the implementation never writes on the contract without
Arnold explicitly being aware.** Enforced, not advised (text rules are advisory; gates
are not):

- Contract docs and acceptance tests **will be** hash-guarded in CI once the repo is
  bootstrapped (first commit, workflow, named guarded paths, a mutation test proving an
  implementation-shaped PR fails — open question #30). **Until that gate exists, no
  implementation PR may be opened.** Once it exists: an implementation PR that touches
  a contract page or an acceptance-test file fails CI.
- A contract change rides only in a dedicated contract PR — no implementation code in
  it — so the diff Arnold reviews is the rule change and nothing else.
- A test weakened, an assert loosened, or a fixture edited inside an implementation PR
  is the failure mode this gate exists for.

A rule that cannot be turned into a test a wrong build fails is not a rule yet; it
lives in `open-questions.md`.

## One rule, one place

- A rule that already exists as an mcp-sso contract clause is **cited**
  (`mcp-sso §NN.C`), never reworded. A reworded copy drifts; a citation goes stale
  loudly when the source changes.
- A rule that is Atesaki's own lives in exactly one of three pages — `contract.md`
  (roles, config rules, ladder, routes, relay, egress, verbs, nevers),
  `contract-grants.md` (everything about grants), `contract-boundaries.md` (config
  reference, reference trust, canonicalization, assertions, caps, proxy trust, errors,
  numbers) — and nowhere else. Other pages point at it.

## Parity is the rigor

The authorization-server half of Atesaki is proven by mcp-sso's language-neutral
fixture corpus (mcp-sso `docs/contracts/19`): deterministic examples — pinned clock,
seeded randomness, corpus-only keys, recorded outbound HTTP — that the Go runner must
pass with **zero skips of frozen *portable* fixtures**. Fixtures labeled *host* bind
only the TypeScript reference (its config envelope, framework glue, filesystem
behavior) and are outside parity. Portable fixtures are **never skipped and never
edited** by Atesaki; where Atesaki intentionally differs, the difference is listed in
`docs/deltas.md` first, mcp-sso labels the reference-only fixtures *host*, and Atesaki
writes its own — a difference not on the list is a bug. A runner contains no
expectations of its own. `[O:2026-08-30]` Atesaki adds its own corpus sections for what mcp-sso does not cover (relay,
routes, ladder, egress); the same freeze discipline applies: draft until the
implementation passes it unchanged, frozen with a receipt, changed only with a contract
change logged in the freeze log.

## What "not the complications" means, concretely

- One binary. One YAML config. No framework adapters, no peer-dependency matrix, no
  scaffolding generator.
- Contract pages stay short enough to read in one sitting; depth goes into fixtures,
  not prose.
- Plain words in every operator-facing sentence. If the operator has to ask "what does
  that mean", the sentence failed.
- No option exists just because the library underneath had it. Every config key traces
  to an observed failure or a decision on record.

## Change protocol

A change is a **sequence of linked review units**, never one mixed diff: (1) the
contract PR — the sentence, the threat-model line if the attack surface moved, and
the fixture/test change, with a freeze-log entry if frozen; Arnold approves it; then
(2) the implementation PR that makes it pass, which touches no contract page and no
acceptance-test file (the visibility gate above rejects it otherwise; unit tests that
belong to the implementation are the implementation's own). Changing a frozen fixture without a
contract change is a bug report against an implementation, never a specification
change (mcp-sso §19.4). The gate itself does not exist until the repo has its first
commit, a CI workflow, and the guarded paths named — until then this is policy prose,
and no implementation PR may be opened.

## Claims discipline

No "never / always / cannot" ships in any doc unless it traces to an enforcing test or
code path. Docs are security surface. Release notes are hand-written and lead with what
the user can now do.
