MODEL: gpt-5.6-sol   EFFORT: xhigh   TOOL: Codex CLI in ~/project/atesaki-core
WHY: B1 is a reference table until a machine-readable schema exists (open question
#36). The schema and its mutation suite are freeze artifacts. Contract-change PR.

RUN MODE: this packet lands through serial PRs. On each fresh run, fetch `origin/main`
and inspect which phases below have merged. Implement only the first incomplete phase.
Branch from current `origin/main`, open one self-explanatory PR, verify it, and stop.
After the owner merges it, start the next phase from the new `origin/main`. Never
stack phases. If a phase is still too large to review as one behavior, split it again
and report the remaining work.

PROVISIONAL PHASES
1. `feat(schema): add config schema validation`
   - `schema/atesaki-config.schema.json`
   - the three complete `schema/valid/` examples
   - the smallest `tools/schema-check.py` path that proves those examples pass
   - bidirectional B1 property-to-schema drift check; this phase cannot merge with
     an omitted or invented property
2. `test(schema): add config refusal mutations`
   - structural B1/B2/B3/B4 mutations, including malformed signed-assertion and
     `keys` union combinations
   - semantic B1/B2/B3 mutations, including machine scopes outside the route
     catalog and declarations wholly denied by their route rules
   - named-rule failure assertions in `tools/schema-check.py`, including the
     `validate` path for rules JSON Schema cannot express
   - split by boundary section if the refusal set is not one reviewable unit
3. `feat(schema): add logical record schemas`
   - the six `schema/records/*.schema.json` files
   - state-dependent required and absent fields
   - focused valid and refusal cases
   - `tools/schema-check.py` mapping from each record case to its logical schema,
     with the expected pass or named-rule refusal asserted
   - bidirectional G2 record-field-to-schema drift check; this phase cannot merge
     with an omitted or invented property
4. `test(schema): enforce contract-schema drift`
   - rerun the B1 property-to-schema comparison in both directions
   - rerun the G2 record-field-to-schema comparison in both directions
   - coverage check that every semantic rule listed for `validate` has an
     executable phase-2 mutation and named-rule assertion
   - final full-schema and mutation run

Read first, fully: docs/contract-boundaries.md (B1–B8, every field and refusal rule —
B4 in full: the signed-assertion union is a schema shape) · docs/contract.md §2, §4,
§13 · docs/contract-grants.md G2, G5, G6, G7, G10, G12 (record states and the
event/reason model are schema inputs) ·
docs/deltas.md D8 · docs/open-questions.md #36 · seed/atesaki.schema.json and
seed/PROVENANCE.md (evidence only — do NOT copy its stale shapes).

DELIVERABLES
1. `schema/atesaki-config.schema.json` — JSON Schema 2020-12 for the YAML stream:
   `Gateway` and `Route` documents, every B1 field with its type, required/optional
   per variant, the identity **tagged union** on `provider` (fields of inactive
   variants refused via `unevaluatedProperties: false`), `machineClients[]`,
   `grant.policy.rules[]` in the G7 vocabulary, references (`env:`/`file:`) as a
   pattern type, durations as positive integers, identifiers and paths in the B3
   grammars, host grammar from B3 as a pattern. `additionalProperties: false`
   everywhere — unknown is refusal.
2. `schema/mutations/` — a **refusal suite**: one YAML file per malformed-input rule
   in B1/B2/B3/B4 (wrong type, unknown field, duplicate key, null required, blank string,
   list-for-scalar, inline secret, non-canonical URL, nested routes, reserved-path
   route, empty catalog, malformed signed-assertion and `keys` union combinations,
   machine scopes outside the route catalog, declaration wholly denied by its route
   rules — that one is a semantic check, note it as out of schema scope and in
   `validate`'s scope). Each names the B-rule it exercises. Plus
   `schema/valid/` — at least three complete valid configs: console loopback,
   Entra dedicated with two routes sharing an upstream, header-mode signed assertion
   with a machine client.
3. `schema/records/*.schema.json` — the portable logical record schemas from G2:
   `grant_request`, `preapproval`, `grant`, `authorization_code` (delta fields only),
   `grant_event`, `machine_tombstone`; state-dependent required/absent fields via
   `if/then`; RFC 3339 UTC 3-ms-digit timestamps as a pattern; snake_case only.
4. `tools/schema-check.py` (stdlib + one pinned validator only if unavoidable — name
   the version and publish date): validates every `schema/valid/*` passes, every
   `schema/mutations/*` fails **for the named rule**, not merely fails, and every
   logical-record case is checked against its named schema for its expected result.
5. Drift checks: every B1 row must map to a config-schema property and vice versa;
   every G2 record field must map to its logical record schema and vice versa. Print
   both diffs. Where a contract table and schema disagree, STOP and list it as a
   contract gap — do not silently pick.

HARD RULES: one phase and one PR per run; touches only `schema/**`,
`tools/schema-check.py`, and — only for gaps the owner asks you to fix — the B1 or G2
source section itself with a one-line reason per change. No Go code. Do not start the
next phase before the current PR merges.

RUN DONE WHEN: the current phase is verified and its PR is open. PACKET DONE WHEN:
`tools/schema-check.py` is green across all merged phases; both drift diffs are empty
or every difference is listed as a gap; lint is green.

REPORT EACH RUN: phase and PR; exact head SHA; changed behavior; checks run; remaining
phases; B1↔config-schema and G2↔record-schema drift found so far; rules the schema
cannot express; open questions the field tables left unanswered.
