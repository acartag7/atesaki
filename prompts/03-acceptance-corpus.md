MODEL: gpt-5.6-sol   EFFORT: xhigh   TOOL: Codex CLI in ~/project/atesaki-core
FALLBACK: grok-4.5
MILESTONES: M2 (phases 0–1), M4 (phase 2), M5 (phase 3) — docs/roadmap.md.
WHY: the tests a WRONG build fails. Atesaki's own fixture corpus, in the shared §19
format, for everything mcp-sso's corpus does not cover. Fixtures for a slice are
written and hash-locked before that slice's code (per-slice freeze, #55); they run red
until the Go exists. Each focused fixture PR is owner-reviewed before merge.

RUN MODE: phases in order; a phase is many serial PRs (one behavior each). On each
run: fetch `origin/main`, find the first incomplete phase, do the next reviewable unit,
open one PR, stop. Later phases never start before the slice that consumes the
earlier phase has begun (phase 2 waits for M3 to start; phase 3 for M4).

Read first, fully: ~/project/mcp-sso/docs/contracts/19-parity-fixture-protocol.md and
fixtures/schema/fixture.schema.json (the spine: `given`/`when`/`then`, seeded
randomness stream, `{absent: true}`, RE2 matchers, `boot` kind, chains and captures,
`profile`) · docs/contract.md §4, §6, §9, §12 (nevers 1–9) · docs/contract-grants.md
G4, G5, G6 (every row), G14 · docs/contract-boundaries.md B2–B7 · docs/deltas.md
(every row, the ones packet 14 added included) · docs/threat-model.md ·
schema/records/** (packet 02 phase 3) · docs/roadmap.md §M2–§M5 test tables.

PHASE 0 — THE ATESAKI FIXTURE PROFILE (blocking prerequisite; the mcp-sso schema
admits only numeric §-clauses, mcp-sso's BridgeConfig, and mcp-sso record kinds).
`fixtures/schema/atesaki-fixture.schema.json`, derived from the mcp-sso schema with
the same spine, additionally accepting:
- clause ids in the G/B/D/never grammars (`G6.A12`, `B3`, `never-8`, `D1`) and the
  §-clauses of mcp-sso where a fixture pins an inherited rule;
- `given.config` as the Atesaki YAML stream, as a string — the Go runner validates it
  with the real parser (no config JSON Schema exists, #54); `given.env` for `env:`
  references (sentinels); `given.files` for B2 boot fixtures: paths, modes, owners
  (`self`/`other`), symlink and hard-link flags, sizes — a filesystem materialization
  the runner builds in a temp dir;
- `given.state`/`then.state` over the Atesaki logical records validated against
  `schema/records/*` plus the inherited mcp-sso records;
- `then.events` over the B7 reason set with its D/F class; durable events also appear
  in `then.state` as `grant_event` rows;
- `given.clients` for CIMD documents the fixture vendors; `given.mockIdp` scripts
  for the identity leg (id_token claims, JWKS) where the inherited `identity.checks`
  form does not fit a redirect flow — prefer the inherited form; add only what is
  needed, and say why.
The profile ships its own mutation suite: an mcp-sso-shaped config, an unknown clause
id, an unknown record field, an unknown reason code, a reason with the wrong class,
a chain with a gap — each REJECTED, named. A contract artifact of THIS repo; the
mcp-sso schema is never edited. `fixtures/MANIFEST.json` + `CATALOGUE.md` generation
with clause-level coverage: every numbered clause / never / operation row in
contract.md §4, §6, §12, contract-grants.md G1–G14 (each A/E row separately),
contract-boundaries.md B2–B7, and every deltas.md row maps to ≥1 fixture id or is
listed uncovered. Hashes of `locked` fixtures; the runner (packet 05) fails on a
mismatch or a skip.

PHASE 1 — SLICE-1 FIXTURES (M3), `profile: portable`, sentinels only, every outbound
call recorded:
1. Nevers 1, 3, 5, 7: never 1 with the upstream stub recording every received header
   and the client's bearer absent from all of them; never 3 asserting exactly `502`
   `upstream_auth_failed`, `WWW-Authenticate` absent, only B7-allowlisted response
   headers; never 5 with a REAL token minted under `given.keys` for route A presented
   at route B, expecting B's challenge; never 7 as the runner's own skip test.
2. Relay rules §6, each bullet: header allowlists both directions (a new header must
   be dropped); upstream 401/403 with a challenge; SSE streamed unbuffered with a
   client disconnect cancelling the upstream; buffered POST completing (the cancel
   binding regression); query-only passthrough; top-level JSON array refused;
   `Transfer-Encoding` plus `Content-Length` → `400`.
3. Boundaries: B3 host grammar accept/refuse pairs; inbound target — absolute-form
   with mismatched authority → `400`, encoded separators not routed, double-encoding
   not routed, dot segments not routed; B5 caps (`413`, `414`, `431`, `429` with
   `Retry-After`); B6 forwarded walk (trusted/untrusted peer, hop cap, malformed entry
   → `400`, duplicate `Host` → `400`, HTTP/2 `:authority` vs `Host` mismatch → `400`);
   B7 every row the relay side reaches, exact status and code.
4. Verifier and discovery: per-route PRM at the path-inserted location and the
   challenge pointing at it (D1); origin AS metadata documents; `iss`/`aud`/`exp`/
   `scope` refusals with the non-oracular shape; duplicate `Authorization` (the §8.4
   portable fixture stays mcp-sso's — do not duplicate it; cite its id).
5. `boot` kind: B2 file invariants via `given.files` (symlink, hard link, wrong owner,
   group-readable, group-writable parent, oversize); console loopback refusals; every
   B1 refusal already covered by `internal/config/testdata` is NOT re-fixtured — the
   coverage map points at that suite as `suite` evidence.
6. Every fixture `draft`; locked in one PR at the end of the phase with the owner's
   read as the receipt (#55).

PHASE 2 — SLICE-2 FIXTURES (M4): §4 rungs (each rung's boot refusals and acceptance;
rung 4: duplicate assertion header, unsigned header, wrong `kid`, stale JWKS beyond
the interval, identity headers stripped on non-identity paths, bounded refetch #58);
never 6 (dedicated rung whose IdP rejects the redirect: the IdP error surfaces, no
fallback); the scope-ceiling delta (#53: Codex-shaped union request on two routes);
`clientOriginIn` (#57); live CIMD fetch guards if #5 allowed; D1, D3, D4, D5's allow
branch, D6, D7, D11, D12, D13; operation rows A1, A2, A3 (insert), A3′, A3″, A7, A8,
A9, A9′ (consumed on binding failure, not on wrong `resource`), A10, A10′ (replay
revokes grant and family), A10″, A11 via RFC 7009, A14 lazy expiry for these rows,
E1–E3; identity-failure pairs per §19.2 (rejection vs port throw) for every identity
path; Entra groups overage → empty ceiling, no outbound call.

PHASE 3 — SLICE-3 FIXTURES (M5): A4, A5 (authority per packet 12: unauthorized OS
user refused per verb; self-approval refused where checkable), A6 (two-runner claim
race as a `suite` receipt with a barrier), A6a, A6b (freshness → invalidated), A12
(first-issuance race, reuse, digest mismatch, tombstone, deny rule, scope outside
declaration), A13, A14 sweeper, A15 purge idempotence; never 8 as the matrix written
in §12 (purpose shape × duration shape × boundaries × caps × races) and never 9 as
three independent mutants plus lineage, asserted on the real JWT `scope` claim; every
durable reason in B7 produced by at least one fixture; every G5 state reached.

HOSTILE-CONSTRUCTION RULES: build the person the title names; real foreign ids on
every id-taking action; exact refusal, never "any 4xx"; never catch the fixture's own
failure; fresh sentinel subjects per fixture; garbage is refusal, never a save. A
fixture and a fail-closed rule in conflict → the rule wins; record why.

HARD RULES: one self-explanatory fixture behavior per PR; contract pages unchanged
unless the owner accepts a proposal under `prompts/README.md`; no Go except the
schema/manifest tooling if it is Go; nothing may depend on the machine running it;
every fixture PR updates the rows it satisfies in `docs/negative-matrix.md` (packet 04).

DONE WHEN (per phase): schema-valid fixtures; coverage map lists every uncovered
clause explicitly; profile mutation suite green; the phase locked with its receipt.

REPORT: coverage counts per page and per phase; uncovered clauses; contract gaps
(rows you could not fixture as written and why) — the most valuable output.
