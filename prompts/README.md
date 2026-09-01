# Prompts — dispatch packets for every remaining piece of work

Each active packet is pasted into a **fresh session** of the named tool. Retired packets
are marked and never dispatched. This design repo writes rules, tests, and prompts; it
never runs the builder itself. Fill every `<…>` before pasting. Model lines follow the
house template and are the owner's to swap — review seats are never downgraded to a
cheaper tier.

## Order and gates

```
 BEFORE FREEZE   00 RETIRED — prose review closed after process finding #2
                 01 RETIRED — repository bootstrap landed; contract gate abandoned
                 02 config + record schemas         (serial phases; one PR per run)
                 03 Atesaki corpus: fixture-format profile FIRST, then the fixtures
                 04 threat model + negative matrix
                 packet 12: the #24 grants-authority contract (needs Arnold's ruling on the proposal) — before 07
                          │
                    ARNOLD READS EVERY PAGE, SAYS "FREEZE"  →  tag contract-v0-freeze
                          │  strictly serial from here — no parallel dispatch:
 AFTER FREEZE    05 Go slice 1 — Atesaki fixture RUNNER + relay + verifier
                    (needs the frozen mcp-sso §8 verifier fixtures — the one cross-lane input)
                 06 Go slice 2 — AS parity (identity flows; grant-dependent endpoints
                    REFUSE until slice 3 — never a fail-open stub)
                 07 Go slice 3 — grants, policy, store, CLI, audit
                 08 deployment recipe + kit (after 07)
                 09 publish readiness (after 08 — not after slice 1)

 PARALLEL (mcp-sso)         10 corpus phases; #340 IS the runner — rebase packets to its head
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
  source before using it.
- Conventional commits (`type(scope): why`). No AI co-author trailers. No "Made with"
  lines. Feature branch + PR. Never force-push `main`.
- **Report shape:** what shipped · what did not and why · contract gaps found (the most
  valuable output) · bugs found in the contract or upstream · exact commands run and
  their results.
