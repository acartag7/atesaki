MODEL: gpt-5.6-sol   EFFORT: xhigh   TOOL: Codex CLI in ~/project/atesaki-core
WHY: B1 is a reference table until a machine-readable schema exists (open question
#36). The schema and its mutation suite are freeze artifacts. Contract-change PR.

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
   in B1/B2/B3 (wrong type, unknown field, duplicate key, null required, blank string,
   list-for-scalar, inline secret, non-canonical URL, nested routes, reserved-path
   route, empty catalog, machine scopes outside the route catalog, declaration wholly
   denied by its route rules — that one is a semantic check, note it as out of schema
   scope and in `validate`'s scope). Each names the B-rule it exercises. Plus
   `schema/valid/` — at least three complete valid configs: console loopback,
   Entra dedicated with two routes sharing an upstream, header-mode signed assertion
   with a machine client.
3. `schema/records/*.schema.json` — the portable logical record schemas from G2:
   `grant_request`, `preapproval`, `grant`, `authorization_code` (delta fields only),
   `grant_event`, `machine_tombstone`; state-dependent required/absent fields via
   `if/then`; RFC 3339 UTC 3-ms-digit timestamps as a pattern; snake_case only.
4. `tools/schema-check.py` (stdlib + one pinned validator only if unavoidable — name
   the version and publish date): validates every `schema/valid/*` passes and every
   `schema/mutations/*` fails **for the named rule**, not merely fails.
5. Drift check: every B1 row must map to a schema property and vice versa; print the
   diff. Where the table and the schema disagree, STOP and list it as a contract gap —
   do not silently pick.

HARD RULES: one focused schema PR; touches only `schema/**`,
`tools/schema-check.py`, and — only for gaps the owner asks you to fix — B1 itself with
a one-line reason per change. No Go code.

DONE WHEN: `tools/schema-check.py` green; drift diff empty or every diff listed as a
gap; lint green.

REPORT: B1↔schema drift list; rules the schema cannot express (belong to `validate`);
open questions the field table left unanswered.
