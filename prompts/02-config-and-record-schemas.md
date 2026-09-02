MODEL: gpt-5.6-sol   EFFORT: xhigh   TOOL: Codex CLI in ~/project/atesaki-core
MILESTONE: M1 residuals (docs/roadmap.md). RESCOPED 2026-09-02 under open question
#54 — the earlier JSON-Schema-plus-Python version of this packet is in git history
(`00ed981`) and is not dispatched. PRECONDITION: #54 ruled as proposed (phases 1 and
3); #56 landed in the contract (phase 2). If #54 is ruled the other way, STOP and
report; the old packet is the fallback.
WHY: the Go validator merged (PR 5, PR 6) with a 71-case refusal suite that names the
rule per case — that suite is the mutation suite. What is still missing is the
mechanical proof that B1 and the parser agree field by field, and the G2 record types
the fixture profile (packet 03) needs for `given.state`/`then.state`. Two validators
for one input would be the parser-differential class; one source (Go) with generated
artifacts is not.

RUN MODE: serial phases, one PR per run. On each fresh run fetch `origin/main`,
inspect which phases merged, implement only the first incomplete phase, open one
self-explanatory PR, verify it, stop. Never stack phases.

Read first, fully: docs/contract-boundaries.md B1 (every row, nested fields inside
type cells included), B2, B4 (the `assertion` object's fields) · docs/contract-grants.md
G2, G5, G6 (the columns that set fields), G12 · docs/deltas.md D8 · internal/config/**
(reader.go's accessors are where the registry hooks in; types.go; parse.go) ·
docs/open-questions.md #54, #56 · seed/atesaki.schema.json (evidence only; stale shapes).

PHASE 1 — `test(config): B1 to parser drift check`
- The parser records every field path it accepts as it reads: dotted paths, `[]` for
  list elements, `<name>` for map keys (`spec.egress.profiles.<name>.caBundleRef`,
  `spec.machineClients[].routes[].scopes[]`), plus the type class the accessor used
  (`string`, `integer`, `boolean`, `list`, `mapping`, `ref`, `url`, `path`,
  `duration`, `union:<tag>`) and requiredness as observed: required, optional, or
  per-variant (`entra`/`oidc`/`header`/`console`; credential `type`; `keys` arm).
  The registry is produced by parsing the three valid examples plus one synthetic
  document per union arm so every arm is visited; a field never visited is a test
  failure ("no positive case for arm X").
- The test parses B1's two tables and B4's `assertion` shape from
  `docs/contract-boundaries.md` (markdown rows `| field | type | rule |`; nested
  fields written inside the type cell as `{a, b?, c}` — state the parsing rules in the
  test file's header comment) into the same path/type/requiredness triples.
- Compare both directions; print both diffs; empty or fail. Where the table and the
  parser disagree, STOP and list every disagreement as a contract gap in the PR —
  never pick a side in code.
- Done when both diffs are empty or every difference is an owner-acknowledged gap.

PHASE 2 — `feat(config): knownCimd entries are references`
- `clients.knownCimd[]` entries parse as B2 references: `env:NAME` holds the document
  text; `file:PATH` follows the B2 file invariants (already implemented for secrets).
  Refusals name the reference, never the content. Whether the content is a valid CIMD
  document stays deferred to slice 2's CIMD validator, as PR 5 stated.
- Refusal cases: bare path (no scheme), unknown scheme, missing env, env blank, file
  invariants — one file each under `testdata/refuse/`, `# expect:` naming the B2 rule.

PHASE 3 — `feat(records): G2 record types and generated schemas`
- `internal/records`: one Go type per G2 record — `grant_request`, `preapproval`,
  `grant`, the `authorization_code` delta fields, `grant_event`, `machine_tombstone`.
  Make illegal states unrepresentable: `state` is a closed enum; fields G6 sets only in
  a given state live in that state's variant (e.g. `grant`'s `active` variant owns
  `activated_at`, `grant_expires_at`, `family_id`), not as pointers on a flat struct.
  A projection function flattens a record to the logical `snake_case` form the runner
  compares in `then.state` (optional = omitted, never null). Timestamps: RFC 3339
  UTC, exactly 3 ms digits, one formatter.
- A generator writes `schema/records/<record>.schema.json` (JSON Schema 2020-12,
  `if/then` per state, `additionalProperties: false`) from the Go types; a golden test
  fails when the committed file differs from the generated one — the types are the
  source, the files are artifacts the fixture profile consumes.
- Mutation suite in Go, per record: one positive case per G5 state branch and per
  `grant_event` reason branch; one refusal per unconditional required field (missing)
  and per typed field (wrong type); one refusal per state-dependent required or
  forbidden field inside each branch (an approved `preapproval` alone carries several).
  Each case is validated against the generated schema by a small pinned JSON Schema
  library (name it, version, publish date, age) or by the projection's own decoder
  if no dependency is needed — state which and why.
- G2↔types drift test: parse G2's record paragraphs (field lists with `?` marking
  optional) and compare with the types' registered fields both ways; empty or fail.

HARD RULES: no config JSON Schema is written (#54); `seed/atesaki.schema.json` is
evidence, never copied; touches only `internal/config/**`, `internal/records/**`,
`schema/records/**`, and — only for gaps the owner asks you to fix — the B1 or G2
source section with a one-line reason per change; one phase per run.

DONE WHEN: phase 1 diffs empty; phase 2 refusals green; phase 3 golden and mutation
tests green and the G2 diff empty; `go test -race ./...` green; lint green.

REPORT EACH RUN: phase and PR; head SHA; the drift diffs as printed; every B1/G2
sentence the parser or types could not honor as written; positive-case counts per
union arm and per state branch.
