MODEL: gpt-5.6-sol   EFFORT: xhigh   TOOL: Codex CLI in ~/project/atesaki-core
FALLBACK: grok-4.5
WHY: the tests a WRONG build fails. Atesaki's own fixture corpus, in the shared §19
format, for everything mcp-sso's corpus does not cover. Contract-change PR. These are
hash-locked by the gate once Arnold approves; they may run red until Go exists.

Read first, fully: ~/project/mcp-sso/docs/contracts/19-parity-fixture-protocol.md and
fixtures/schema/fixture.schema.json (the format, incl. `profile`, seeded randomness,
absence assertions, RE2 matchers, `boot` and `suite` evidence kinds) · docs/contract.md
§6, §12 (nevers 1–9) · docs/contract-grants.md G5, G6 (every row), G14 ·
docs/contract-boundaries.md B3–B7 · docs/deltas.md · docs/threat-model.md.

DELIVERABLE 0 — THE ATESAKI FIXTURE-FORMAT PROFILE (blocking prerequisite: the shared
mcp-sso schema cannot encode Atesaki — it admits only numeric §-clauses, mcp-sso's
BridgeConfig, and mcp-sso record kinds). Author
`fixtures/schema/atesaki-fixture.schema.json` as a derived profile of the mcp-sso
schema (same given/when/then spine, seeded-randomness stream, RE2 matchers, boot/suite
kinds) additionally accepting: clause ids in the G/B/D/never grammars (`G6.A12`, `B3`,
`never-8`); `given.config` as the Atesaki YAML stream validated against
`schema/atesaki-config.schema.json` (packet 02); `given.state`/`then.state` over the
Atesaki logical records (`schema/records/*`); `then.events` over the B7 reason set with
its D/F class; a **filesystem materialization** input (`given.files`: paths, modes,
owners, symlink/hardlink flags) so B2's runtime rules are fixturable as `boot` kind.
The profile ships its own **mutation suite** — invalid clause ids, records, events,
and mcp-sso-shaped configs must be REJECTED, each named. A contract artifact of THIS
repo; the mcp-sso schema is never edited.

DELIVERABLES — `fixtures/` in this repo, `profile: portable`, one fixture per clause
instance, chains for flows, sentinels only, every outbound call recorded:
1. **Nevers 1–9** (contract.md §12): the matrices as written — never 8 across purpose
   shape × duration shape × boundaries × caps × races; never 9 as three independent
   scope mutants plus lineage, asserted on the real JWT `scope` claim; never 3 asserting
   exactly `502 upstream_auth_failed`, no `WWW-Authenticate`, allowlisted headers only;
   never 5 with a REAL token minted for route A presented at route B.
2. **Relay rules** (§6): each bullet; upstream stub answering 401/403 with a challenge;
   header allowlists both directions (a new header must be dropped); SSE streamed
   unbuffered; absolute-form target with mismatched authority → 400; top-level JSON
   array → refused.
3. **Ladder rungs** (§4): each rung's boot refusals and acceptance; rung 4 signed
   assertion — duplicate header, unsigned header, wrong `kid`, stale JWKS beyond the
   interval, identity headers stripped on non-identity paths.
4. **Operation table** (G6): one fixture per row and per named failure branch, as
   `boot`/HTTP/`suite` kinds as appropriate; A6 two-runner claim race; A6b freshness
   failure → invalidated; A12/A13/A14 machine issue/revoke/expire races including the
   per-route digest rule; A9′ consumption semantics (consumed on binding failure, not on
   wrong resource); A10′ replay revokes grant and family. Every durable event reason
   in B7 must be produced by at least one fixture; every G5 state reached.
5. **Boundaries**: B3 host grammar accept/refuse pairs; B6 forwarded-IP walk
   (trusted/untrusted peers, hop cap, malformed entry → 400); B7 every public error
   row reached at least once with its exact status and code.
6. `fixtures/MANIFEST.json` + `CATALOGUE.md` (generated) with **clause-level** coverage
   (the §19 simplification applies here too): every numbered clause / never / operation
   row in contract.md §6/§12, contract-grants.md G1–G14 (each A/E row separately,
   incl. A3″ and A15), contract-boundaries.md B2–B7, and **every deltas.md row D1–D13**
   maps to ≥1 fixture id or is listed as uncovered. Sentence quotes stay drift checks, never coverage units.

SLICE OWNERSHIP: never 6 (redirect identity) belongs to slice 2's fixture-ID set,
not slice 1's.

HOSTILE-CONSTRUCTION RULES (holes.md classes): build the person the title names; real
foreign ids on every id-taking action; exact refusal, never "any 4xx"; never catch the
fixture's own failure; fresh sentinel subjects per fixture; garbage is refusal, never a
save. A fixture and a fail-closed rule in conflict → the rule wins; record why.

HARD RULES: contract-change PR; contract pages unchanged (gaps go in the PR text, not
into the docs); no Go code; nothing may depend on the machine running it.

DONE WHEN: schema-valid fixtures; coverage map lists every uncovered clause explicitly;
lint green; gate green.

REPORT: coverage counts per page; uncovered clauses; contract gaps (rows you could not
fixture as written and why) — the most valuable output.
