MODEL: gpt-5.6-sol   EFFORT: high   TOOL: Codex CLI in ~/project/mcp-sso, FRESH session, do NOT resume the PR #338 thread STATUS NOTE: the runner already exists as draft PR #340, this packet applies ONLY if #340 is abandoned; otherwise review and land #340 instead. PRECONDITION: PR #338 merged (or checked out as the base). §19 is settled: clause-level coverage, portable/host profiles, boot/suite evidence kinds. It is NOT to be edited. WHY: fourteen hours of the previous session produced the protocol and two fixtures but no runner. This packet is code only. One deliverable. One allowed question.

READ, fully, then do not reopen: docs/contracts/19-parity-fixture-protocol.md · fixtures/schema/fixture.schema.json · fixtures/08-resource-server-verifier/*.json (the host 8.4 and the portable 8.4) · fixtures/README.md · src/ports/ (clock, randomness, store, identity, audit ports) · the Fastify adapter and the fastify-sqlite example composition.

THE ONE DELIVERABLE: the reference runner, TypeScript, in this repo:
- Loads a fixture; composes the reference implementation through the Fastify adapter from `given.config`; pins the clock port to `given.clock`; feeds `given.random.seed` through a seedable randomness port (add that port ONLY if none exists, narrowest change, production behavior identical, noted in the PR); loads keys from `fixtures/keys/`; allows outbound HTTP ONLY as recorded in `given.http`, an unrecorded outbound call fails the fixture; sends `when.request` as real HTTP; compares status, headers (RE2 `matches` where the fixture says so, exact otherwise, duplicate occurrences preserved), body, audit events, store effects, `then.outbound`.
- Contains no expectations of its own. A skipped `frozen` fixture is a failure.
- `pnpm test:fixtures` runs it; CI runs it. MANIFEST/CATALOGUE/FREEZE-LOG belong to the follow-up PR, not this one.
- Run both 8.4 fixtures unchanged and report pass/fail per fixture. **Freezing is NOT this packet**, status/receipt edits touch fixture files and belong to a separate follow-up PR (packet 10 phase D). On failure: STOP and report per §19.6, bug vs contract gap, with the diff, and do not edit the fixture.

HARD RULES
- **Do not edit anything under docs/contracts/ or fixtures/** (all of it, nested included). The CI guard fails you; if you believe a contract change is needed, that is the one allowed question.
- **One question maximum per run.** If you hit a second, stop and report both without choosing. Never ask about spellings, anchors, selectors, inventories, or naming.
- Do not post `@codex review` (quota-blocked). Do not open new PRs beyond this one.
- No new dependencies unless unavoidable; if one, name version and publish date and meet the repo's cooldown (mcp-sso runs a hard 15 days).
- Conventional commits, no AI trailers, branch + PR, never force-push main.

DONE WHEN: runner in CI; both 8.4 fixtures pass unchanged (or STOPped with the §19.6 report); zero doc or fixture edits. Freeze and MANIFEST/CATALOGUE come in the follow-up PR.

REPORT (short): commands run and results · pass/fail per fixture · the single question if any · nothing else.
