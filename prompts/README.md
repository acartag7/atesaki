# Prompts — dispatch packets for every remaining piece of work

Each active packet is pasted into a **fresh session** of the named tool. Retired packets
are marked and never dispatched. This design repo writes rules, tests, and prompts; it
never runs the builder itself. Fill every `<…>` before pasting. Model lines follow the
house template and are the owner's to swap — review seats are never downgraded to a
cheaper tier. The milestone each packet serves, with the capability it must produce,
the security surface it opens, and the tests that prove it, is in `docs/roadmap.md`.

## Order and gates

```
 RETIRED         00 prose review · 01 repository bootstrap (contract gate abandoned)

 M0              13 repo hardening — CI required on main, license, SECURITY.md,
                    dependency cooldown, two config nits, STATE refresh
 M1              02 config drift check + record types (RESCOPED, waits on #54)
 M2              14 contract closure — #53 scope ceiling, #5 live fetch, #56 knownCimd
                    refs, #57 clientOriginIn, PR-5 interpretations, #58, #59, #55,
                    #60 limiter outage, #61 readiness, B4 alg wording, B8 note, matrix window
                 12 grants authority (waits on the #24 ruling; fixtures written in 03 phase 3)
                 03 phase 0 (fixture profile) then phase 1 (slice-1 fixtures)
                 04 threat model + negative matrix
                          │
                 PER-SLICE FREEZE (#55): the slice's sections SHA-pinned in its packet,
                 its fixtures hash-locked, the owner has read those pages
                          │  strictly serial from here — no parallel dispatch:
 M3              05 Go slice 1 — runner + relay + verifier + validate --deep
                    (needs the frozen mcp-sso §8 verifier fixtures — cross-lane input;
                    03 phase 2 is written alongside and locked before 06)
 M4              06 Go slice 2 — sign-in, store port + SQLite, the allow
                    path (A1, A7–A11), idp-request
                    (needs mcp-sso §07/§09/§10/§11 frozen, or every unfrozen id listed
                    as deferred in the PR — never skipped silently)
                    (03 phase 3 written alongside, locked before 07)
 M5              07 Go slice 3 — escalation, approvals CLI, machine
                    clients, sweeper, outbox  →  tag contract-v0-freeze
 M6              15 rehearse — mock IdP, client profiles
 M7              08 deployment recipe + kit
 M8              09 publish readiness (name check already done in M0)

 PARALLEL (mcp-sso)         10 corpus phases — the runner is on main as serial PRs
 EVERY IMPLEMENTATION PR    11 adversarial implementation review
```

## Conventions every packet shares

- **Read-first lists are exact.** Read the whole file, not the first screen.
- **One self-explanatory review unit per PR.** If a packet or milestone needs a large
  change, decouple it into small serial PRs. Merge each PR before branching the next
  from current `main`. Each PR explains and verifies its own behavior.
- **The implementation never changes the contract silently.** When something does not
  fit, state the exact conflict and propose a concrete change. Continue with unaffected
  work. The owner decides whether the focused PR may include the change or whether to
  split a smaller linked PR.
- **Ambiguity is a contract gap, not your call.** Record the exact question in the PR
  under "contract gaps" and continue with the next item. Stop only the affected work.
  Never invent a rule.
- **Never weaken a fail-closed control to make a test pass.**
- **Verify by running** — the test suite, the binary, a real input. "Compiles" is not
  done. Name the real input you ran.
- **No invented APIs.** Grep every external symbol against the actual dependency
  source before using it. Every new dependency: exact version, publish date, age
  against the cooldown, and why the standard library was not enough.
- **Fixtures before code, per slice.** A slice's Atesaki fixtures are hash-locked and
  its sections pinned before its first implementation PR. A skipped locked fixture is
  a failure. Fixture phases for the next slice are written alongside the current
  slice's code and locked before the next slice starts.
- **Packets cite, they do not legislate.** A packet names the row, the artifact to
  produce, and the proof to run. Where a packet paraphrases a rule, the contract's
  sentence wins; a paraphrase that disagrees with it is a bug in the packet.
- Conventional commits (`type(scope): why`). No AI co-author trailers. No "Made with"
  lines. Feature branch + PR. Never force-push `main`.
- **Report shape:** what shipped · what did not and why · contract gaps found (the most
  valuable output) · bugs found in the contract or upstream · exact commands run and
  their results.
