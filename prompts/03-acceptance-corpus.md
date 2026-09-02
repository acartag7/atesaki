MODEL: gpt-5.6-sol   EFFORT: xhigh   TOOL: Codex CLI in ~/project/atesaki-core
FALLBACK: grok-4.5
MILESTONES: M2 (phases 0–1), M4 (phase 2), M5 (phase 3) — docs/roadmap.md.
WHY: the tests a WRONG build fails. Atesaki's own fixture corpus, in the shared §19
format, for everything mcp-sso's corpus does not cover. Fixtures for a slice are
written and **slice-locked** before that slice's code (the slice lock, #55): their
§19 status stays `draft` until the runner passes them; the lock is an owner-approved
`fixtures/LOCK-<slice>.json` listing fixture hashes and the pinned contract SHAs;
`frozen` needs a passing-runner receipt. They run red until the Go exists.

RUN MODE: phases in order; a phase is several serial PRs, **one per invariant or
protocol chain with all its fixtures** (a never, a G6 row family, a B-section's
grammar, a delta row) — never one fixture per PR. On each run: fetch `origin/main`,
find the first incomplete phase, do the next reviewable unit, open one PR, stop. Later phases never start before the slice that consumes the
earlier phase has begun (phase 2 waits for M3 to start; phase 3 for M4).

Read first, fully: ~/project/mcp-sso/docs/contracts/19-parity-fixture-protocol.md and
fixtures/schema/fixture.schema.json (the spine: `given`/`when`/`then`, seeded
randomness stream, `{absent: true}`, RE2 matchers, `boot` kind, chains and captures,
`profile`) · docs/contract.md §4, §6, §9, §12 (nevers 1–9) · docs/contract-grants.md
G4, G5, G6 (every row), G14 · docs/contract-boundaries.md B2–B7 · docs/deltas.md
(every row, the ones packet 14 added included) · docs/threat-model.md ·
schema/records/** (packet 02 phase 2) · docs/roadmap.md §M2–§M5 test tables.

PHASE 0 — THE ATESAKI FIXTURE PROFILE (blocking prerequisite; the mcp-sso schema
admits only numeric §-clauses, mcp-sso's BridgeConfig, and mcp-sso record kinds).
`fixtures/schema/atesaki-fixture.schema.json`, derived from the mcp-sso schema with
the same spine, additionally accepting:
- clause ids in the G/B/D/never grammars (`G6.A12`, `B3`, `never-8`, `D1`) and the
  §-clauses of mcp-sso where a fixture pins an inherited rule;
- `given.config` as the Atesaki YAML stream, as a string — the Go runner validates it
  with the real parser (no config JSON Schema exists, #54); `given.env` for `env:`
  references (sentinels); `given.files` for B2 boot fixtures under a **containment
  contract** written into the schema: paths are relative, one segment grammar, no
  `.`/`..`/absolute/empty segments; symlink and hard-link targets must resolve inside
  the materialization root; modes as octal strings; sizes bounded (per-file and total
  byte caps, file-count cap); **no ownership simulation** (`owner: other` is not
  expressible — the owner-mismatch rule is proven by a Go unit test with an injected
  stat, recorded as `suite` evidence); the runner materializes under `os.Root` and
  refuses any fixture that violates the grammar before touching the disk (the corpus
  is supply-chain input; a hostile fixture must not become a file-write primitive);
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
a chain with a gap, a `given.files` path with `..`, an absolute path, a link target
outside the root — each REJECTED, named. A contract artifact of THIS repo; the
mcp-sso schema is never edited. `fixtures/MANIFEST.json` + `CATALOGUE.md` generation
with clause-level coverage: every numbered clause / never / operation row in
contract.md §4, §6, §9 (verbs), §12, contract-grants.md G1–G14 (each A/E row
separately), contract-boundaries.md **B1–B8** (B1 rows point at the config refusal
suite as `suite` evidence; every B8 number has a fixture that exercises its exact
boundary), and every deltas.md row maps to ≥1 fixture id or is listed uncovered. A
fixture may carry `inherited: <mcp-sso clause>` when it pins an inherited sentence
that has no frozen upstream portable fixture yet; it is superseded when the upstream
one freezes, and until then it is what makes "inherited" mean "tested". Hashes of
`locked` fixtures; the runner (packet 05) verifies them before materializing anything
and fails on a mismatch or a skip.

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
   `Transfer-Encoding: chunked` beside `Content-Length` → the observable Go's parser
   produces (chunked honored, length dropped) — pinned, not wished away.
3. Boundaries: B3 host grammar accept/refuse pairs; inbound target — absolute-form
   with mismatched authority → `400`, encoded separators not routed, double-encoding
   not routed, dot segments not routed; B5 caps (`413`, `414`, `431`, `429` with
   `Retry-After`); B6 forwarded walk (trusted/untrusted peer, hop cap, malformed entry
   → `400`, duplicate HTTP/1.1 `Host` → `400` — refused by Go's parser before any
   handler, so the fixture asserts the status and **no** audit event, and the
   negative-matrix row cites the parser; HTTP/2 `Host` beside `:authority` with a
   different value → `400`); header count over B8 → `431`; B7 every row the relay
   side reaches, exact status and code.
4. Verifier and discovery: per-route PRM at the path-inserted location and the
   challenge pointing at it (D1); origin AS metadata documents; `iss`/`aud`/`exp`/
   `scope` refusals with the non-oracular shape; `alg` in the token ≠ the configured
   algorithm, or a key of the wrong type → refused; duplicate `Authorization` (the
   §8.4 portable fixture stays mcp-sso's — do not duplicate it; cite its id); an
   `inherited` fixture for every §7 token clause the verifier implements that has no
   frozen upstream fixture; `livez`/`readyz` and the drain per #61 as ruled.
5. `boot` kind: B2 file invariants via `given.files` (symlink, hard link,
   group-readable, group-writable parent, oversize — never wrong owner, which the
   profile cannot express and the injected-stat Go suite receipt covers); console
   loopback refusals; every
   B1 refusal already covered by `internal/config/testdata` is NOT re-fixtured — the
   coverage map points at that suite as `suite` evidence.
6. Every fixture `draft`; the slice lock written in one PR at the end of the phase
   with the owner's read as the receipt (#55); `frozen` only after packet 05's runner
   passes them.

PHASE 2 — SLICE-2 FIXTURES (M4, the whole human loop): §4 rungs (each rung's boot
refusals and acceptance; rung 4: duplicate assertion header, unsigned header, wrong
`kid`, stale JWKS beyond the interval, identity headers stripped on non-identity
paths, bounded refetch #58); never 6 (dedicated rung whose IdP rejects the redirect:
the IdP error surfaces, no fallback); the slice-2 configuration fields (`boot`
fixtures the old parser must fail: `clients.cimd.liveFetch`, `clientOriginIn`, the
approver objects, `knownCimd` references, forbidden credential header names, the
inherited §10 redirect-entry grammar); the consent-page carrier (#62): purpose and
duration POSTed with the consent, absent from every URL and flow line, hostile
purpose (HTML, Unicode, control characters, over cap) refused or escaped per the
inherited page controls, policy evaluated on the submitted values; the two-stage
ceiling (#53: catalog empty → `invalid_scope`; group ceiling empty → `access_denied`
inherited; the Codex-shaped union request on two routes);
`clientOriginIn` (#57); live CIMD fetch if #5 allowed: origin not on the allowlist
refused before any network call, the inherited caps cited by clause; the limiter-
outage delta (#60); D1, D3, D4, D5's allow branch, D6, D7, D11, D12, D13; operation
rows A1, A2, A3 (insert), A3′, A3″, A7, A8, A9, A9′ (consumed on binding failure, not
on wrong `resource`), A10, A10′ (replay revokes grant and family), A10″, A11 via RFC
7009, A4, A5, A6, A6a (two-runner barrier as a `suite` receipt), A6b (one
transaction, #66), the packet-12 authority rules as **suite receipts** with an
injected effective-identity port (the profile has no command carrier: unauthorized
uid refused per verb; self-approval refused where checkable; `claimed_approver`
never authority), the #62 state machine (C1 `entry` POST → deny / allow / escalate
/ claim; C2 `confirm` POST → A7/A8; a C1 replay refused; a C2 presented at the
entry stage refused), A14's lazy transitions inside A3's cap/dedupe
read and inside A6/A9/A10, E1–E3; the projector cursor (#64: a durable event reaches
JSONL after a sink failure and a restart); never 8 and never 9 as the matrices
written in §12 (purpose shape × duration shape × boundaries × caps × races; three
scope mutants plus lineage on the real JWT `scope` claim); the A10 crash pair (commit
then lost response; the client's retry ends the grant);
identity-failure pairs per §19.2 (rejection vs port throw) for every identity path;
Entra groups overage → the inherited identity refusal with its reason, no outbound
call; a display name as a `groupsToScopes` key → boot refusal; an `inherited` fixture
for every §17 identity clause this slice implements that has no frozen upstream
fixture.

PHASE 3 — SLICE-3 FIXTURES (M5): A12 (first-issuance race, reuse, digest mismatch,
tombstone, deny rule, scope outside declaration) and A13 if #67 keeps machine
clients; A14 sweeper; A15 purge idempotence; the migration and downgrade-refusal
`boot` fixtures (#65); every durable reason in B7 produced by at least one fixture;
every G5 state reached.

HOSTILE-CONSTRUCTION RULES: build the person the title names; real foreign ids on
every id-taking action; exact refusal, never "any 4xx"; never catch the fixture's own
failure; fresh sentinel subjects per fixture; garbage is refusal, never a save. A
fixture and a fail-closed rule in conflict → the rule wins; record why.

HARD RULES: one invariant or protocol chain per PR with all its fixtures; contract pages unchanged
unless the owner accepts a proposal under `prompts/README.md`; no Go except the
schema/manifest tooling if it is Go; nothing may depend on the machine running it;
every fixture PR updates the rows it satisfies in `docs/negative-matrix.md` (packet 04).

DONE WHEN (per phase): schema-valid fixtures; coverage map lists every uncovered
clause explicitly; profile mutation suite green; the phase locked with its receipt.

REPORT: coverage counts per page and per phase; uncovered clauses; contract gaps
(rows you could not fixture as written and why) — the most valuable output.
