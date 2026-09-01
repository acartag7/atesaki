# Quality bar — how this repo works

**Status: DRAFT.**

## Order of work

Contract → threat model → tests → code. Acceptance tests are written from the contract
and may run red before an implementation exists — that is normal, not a defect.

## Who may change the contract `[O:2026-08-30]` `[O:2026-09-01]`

The implementation never changes the contract or its acceptance fixtures silently.
The review process enforces that rule:

- When the implementation does not fit, the implementer states the exact conflict and
  proposes a concrete change. The implementer continues with unaffected work.
- The owner decides whether the focused PR may include the change or whether the
  change belongs in a smaller linked PR. The implementer applies that decision and
  finishes the work.
- The implementer never invents an exception, weakens a test, or edits a fixture
  silently.

No hash manifest or self-checking pull-request workflow decides this. Automation checks
product behavior. Review owns changes to the contract.

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

Start with one reviewable unit per PR. If implementation exposes a contract mismatch,
state the problem and propose the change before editing the contract or its acceptance
fixtures. The owner then chooses whether to keep the change in the current focused PR or
split it into a linked contract PR.

A frozen fixture still changes only with the contract and freeze-log receipt required
by mcp-sso §19.4. Unit tests that belong to the implementation stay with the
implementation.

## Claims discipline

No "never / always / cannot" ships in any doc unless it traces to an enforcing test or
code path. Docs are security surface. Release notes are hand-written and lead with what
the user can now do.
