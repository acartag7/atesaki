# Grants — dispensing, policy, and lifecycle (GrantV0)

**Status: DRAFT.** Owns everything about grants; `contract.md §3` points here. Tags:
`[O:date]` owner decision (receipt in `docs/decisions.md`) · `[S:mcp-sso §N]` inherited
· `[D ← rule]` derived from a named rule · `[P]` **proposed, awaiting Arnold — not
decided** · `[#]` a number in `contract-boundaries.md` B8. Timestamps: RFC 3339 UTC,
exactly 3 ms digits. Field names: `snake_case`.

## G1. What a grant is `[O:2026-08-30]`

Atesaki dispenses **grants**, never standing access. Every credential traces to one
grant: **who** (subject + client), **for what** (one route = one audience), **why**
(purpose — evidence, never authority `[O:ROADMAP ch.2]`), **for how long** (expiry;
no "until revoked"). A grant is a **bounded lineage**: an interactive grant owns
exactly one refresh family; a **machine grant owns no family** (tokens only — the
explicit machine exception `[D ← O:2026-08-31 machine clients]`); a new authorization
is a new grant; authority never carries over (D2). Revocation is **bounded**: refresh
stops now; issued access tokens live to `exp` ≤ access TTL — docs say "access ends
within *N* of revocation", never "immediately". All ids are library-generated and
unguessable.

**Purpose** is stored on `grant_request` and `grant` rows (whether `grants list`
shows or redacts it is #24's call), written to exactly one durable event (`grant_issued` /
`grant_machine_issued`), and never appears on flow lines (G12).

**`resource` is required** on every authorize and machine-token request and must equal
one route's audience byte-exact; omitted or unknown → `invalid_target`. A multi-route
gateway has no default resource to fall back to `[D ← D1 multi-route; RFC 8707]` (D13).

## G2. Records

Each record carries `state`; fields marked `?` are absent unless the state named in
G6 sets them.

`grant_request` — `request_id`, `state`, `subject`, `client_id`, `resource`,
`requested_scopes_raw[]` (pre-ceiling, the hash input), `scopes[]`
(post-ceiling), `requested_duration_s`, `purpose`, `requested_hash`, `policy_version`,
`created_at`, `resolved_at?`, `terminal_at?`, `preapproval_id?`, `grant_id?`.

`preapproval` — `preapproval_id`, `request_id`, `state`, `created_at`,
`awaiting_expires_at`, `approver?`, `decided_at?`, `approved_scopes[]?`,
`approved_duration_s?`, `approved_hash?`, `claim_expires_at?`, `claimed_at?`,
`claimed_by_request_id?`, `invalidated_at?`, `terminal_at?`.

`grant` — `grant_id`, `state`, `kind` (`interactive` | `machine`), `subject`,
`client_id`, `resource`, `scopes[]`, `approved_duration_s`, `purpose`, `created_at`,
`request_id?`, `preapproval_id?`, `approved_hash?`, `activated_at?`,
`grant_expires_at?`, `revoked_at?`, `terminal_at?`, `family_id?`, `declaration_digest?`.

`authorization_code` — inherited shape `[S:mcp-sso §7.3]` **plus** `grant_id`,
`request_id` (D6).

`grant_event` — `event_id`, `seq`, `at`, `reason` (B7 durable set), `grant_id?`,
`request_id?`, `preapproval_id?`, `subject?`, `client_id?`, `resource?`, allowlisted
fields per reason. Durable outbox; JSONL fan-out best-effort (G12).

`machine_tombstone` — `client_id`, `resource`, `declaration_digest`, `revoked_at`,
`revoked_by`.

Exact portable logical schemas are corpus work (D8, #36).

## G3. Hashes and digests (mechanical)

`requested_hash` = SHA-256 over `atesaki-grant-request-v1\n` + RFC 8785 canonical JSON
of `{subject, client_id, resource, requested_scopes_raw (the client's scope list as
requested, sorted, deduplicated, **before** the ceiling), duration_s (integer),
purpose_hex}`; `approved_hash` = same over approved scopes and duration. Raw scopes are
the matching key so a later ceiling change cannot change identity — the ceiling is
enforcement (§9.3 step 4 and A6b), never identity; this is what makes the
group-removed-after-approval branch reachable (round-6 fix). **Not
hashed:** groups and `policy_version` (re-evaluated at claim, A6) and transport context
(`redirect_uri`, `code_challenge`) `[O:2026-08-31]` — the approver approves *what* to
*whom*; the OAuth recipient binding is owned by the final consent inside a flow whose
redirect and PKCE were validated (§9.3 steps 2–3).

`declaration_digest` is **per (client, route)** `[D ← O:tombstone]`: SHA-256 over
`atesaki-machine-declaration-v1\n` + canonical JSON of `{client_id, route_path, scopes
(sorted), purpose_hex, max_duration_s}`. Reordering routes or editing another route
never changes it, so a tombstone survives unrelated edits and clears only when *this*
binding is deliberately changed.

## G4. The flow and what carries the new parameters — D5, complete

`purpose` and `requested_duration` are **authorize-request parameters sent by the
client** `[O:2026-08-31]`. Carriers: the direct-authorize **singleton guard** (each at
most once per carrier; duplicate = `invalid_request`); the signed **§17.11 bridge-flow
params** and the **synthetic callback request** `[S:mcp-sso §17.11]`; **console
pairing**'s query + form round trip `[S:mcp-sso §17.5]` with this reconciliation rule
(mirroring the reference's `resource` rule): present in both carriers with equal
values = accepted; different values = `invalid_request`; singleton within each carrier.
Raw request-target bytes are capped before any parsing (B5).

Position in mcp-sso §9.3 (both identity paths deliver step 1's subject): steps 1–4 →
**G-a** parse `purpose` (B5 shape) and `requested_duration` (integer seconds, `1 ≤ d ≤
Route.spec.grant.maxDuration`; anything else `invalid_request`, redirect channel) →
**G-b** compute `requested_hash` (raw scopes); when step 4's ceiling removed ≥1
requested scope, emit flow `scope_ceiling_applied`; look up a claimable pre-approval
(A6) → **G-c** policy
(G7) → `allow`: sign consent (step 6) with `purpose`, `approved_duration_s`,
`request_id` (D3) → consent (A7/A8) → exchange (A9). Every request reaching G-c produces
a `grant_request` **except** an exact duplicate of an open escalated request, which is
answered from the existing record (A3 dedupe) and logged, with no new row.

## G5. States

```
grant_request: allowed | escalated | denied | unavailable        (created directly in its outcome state)
               allowed   → consented | abandoned
               escalated → resolved_approved | resolved_denied | resolved_expired | resolved_invalidated
preapproval:   awaiting  → approved | denied | expired
               approved  → claimed | expired | invalidated
grant (interactive): issued → active → expired | revoked ;  issued → expired | revoked
grant (machine):     active → expired | revoked
```
Terminal: `denied`, `unavailable`, `consented`, `abandoned`, `resolved_*`, `expired`,
`claimed`, `invalidated`, `revoked`. Nothing transitions backwards. A re-run is a
**new** `grant_request`. Expiry transitions fire lazily (any operation touching a
past-due row records the expiry first, inside its own transaction) or by the sweeper
(A14, `[#]`), whichever first — exactly one event per expiry.

## G6. Operation table — the contract

Two phases per operation `[O:2026-08-31 — see G8]`: **preflight** (outside the store: identity,
request validation, policy evaluation, approver-group check, client authentication,
signing inputs) and **one store transaction** containing only authoritative reads,
compare-and-set / uniqueness predicates, mutations, and durable events. Tokens are
signed **before** the transaction from pre-generated ids; the response is returned
**after** commit. Failure in preflight or before commit ⇒ nothing consumed, nothing
issued. **Event-only operations** (E-rows) write flow events with no store transaction.
A behavior not in this table does not exist.

| # | Operation | Preflight | In-tx predicates | Mutations | Durable events | Output / public error |
| --- | --- | --- | --- | --- | --- | --- |
| A1 | authorize → `allow` | §9.3 1–4 (incl. exact `resource`); G-a; policy `allow`; consent token signed with `request_id` | no claimable pre-approval for `requested_hash` (pending caps bind escalations only — A3/A3″) | insert `grant_request` `allowed` (`created_at`) | `request_allowed` | consent page |
| A2 | authorize → `deny` | same preflight as A1 up to policy; policy `deny` | — | insert `grant_request` `denied` (`resolved_at`) | `request_denied_policy` | 302 `access_denied` |
| A3 | authorize → `escalate` | same preflight as A1 up to policy; policy `escalate` | if an **open** (`escalated`, pre-approval `awaiting`) request with equal `requested_hash` exists for (subject, client, resource): dedupe; else per-tuple open count < cap **and total pending count < cap** | insert `grant_request` `escalated` + `preapproval` `awaiting` (`awaiting_expires_at` = now + approval window `[#]`); **or** no insert | `request_escalated` **or** flow `request_deduplicated` | 302 `access_denied` + `approval_pending` + `request_id` (existing id on dedupe) |
| A3″ | authorize → escalate, a pending cap full | same preflight as A1 up to policy; policy `escalate` | (per-tuple count = cap `[#]` with **no** equal-hash open request) **or** total pending = cap `[#]` | — (no insert) | flow `cap_exceeded` | 302 `temporarily_unavailable` `[D ← O:2026-08-31 total-cap outcome]` |
| A3′ | authorize, decider unreachable | same preflight as A1 up to policy; decider error | — | insert `grant_request` `unavailable` (`resolved_at`) | `request_unavailable`, flow `decider_unavailable` | 302 `temporarily_unavailable` `[O:2026-08-31]` |
| A4 | approver approves | approver authenticated (#24), ∈ route approvers, ≠ request subject; approved ⊆ requested; duration ≤ requested | `preapproval` `awaiting`, not past `awaiting_expires_at` | `preapproval` → `approved`: `approver`, `decided_at`, approved values, `approved_hash`, `claim_expires_at` = now + claim window `[#]` | `preapproval_approved` | CLI |
| A5 | approver denies | as A4 | `preapproval` `awaiting` | `preapproval` → `denied` (`approver`, `decided_at`); request → `resolved_denied` (`resolved_at`) | `preapproval_denied`, `request_resolved` | CLI |
| A6 | re-run claims | §9.3 1–4; G-a; **current** policy ≠ `deny` and approved scopes ⊆ **current** ceiling `[O:2026-08-31]`; consent token signed with the **approved** values and the new `request_id` | CAS: `preapproval` `approved` ∧ now < `claim_expires_at` ∧ its request's `requested_hash` equals ours ∧ same (subject, client, resource) | `preapproval` → `claimed` (`claimed_at`, `claimed_by_request_id`); original request → `resolved_approved` (`resolved_at`); insert new `grant_request` `allowed` | `preapproval_claimed`, `request_resolved`, `request_allowed` | consent page (approved values) |
| A6a | re-run, CAS lost | as A6 | CAS fails (already `claimed`/`expired`) | none in this tx; the request proceeds as **A3** in a new transaction (no open duplicate exists any more) | flow `preapproval_claim_lost_race`, then A3's | as A3 |
| A6b | re-run, freshness fails | as A6 but current policy `deny` or ceiling narrower | `preapproval` `approved` | `preapproval` → `invalidated` (`invalidated_at`); original request → `resolved_invalidated`; then **A3** (a fresh escalation; A3 dedupe cannot match an invalidated pre-approval) or **A2** if current policy is `deny` | `preapproval_invalidated_stale`, `request_resolved`, then A3's/A2's | as A3 / A2 |
| A7 | consent denied | consent token verifies; origin check | JTI unconsumed | consume JTI `[O:2026-08-31, D11]`; request → `abandoned` (`resolved_at`) | `consent_denied`, `request_abandoned` | 302 `access_denied` |
| A8 | consent approved | consent token verifies; origin check; approval `true` | JTI unconsumed | consume JTI; insert `grant` `issued` (`kind=interactive`, subject, client, resource, scopes, `approved_duration_s`, purpose, `request_id`, `preapproval_id?`, `approved_hash?`, `created_at`); insert `authorization_code` (+`grant_id`, `request_id`) (D6); request → `consented` (`resolved_at`, `grant_id`) | `grant_issued` (carries purpose), `request_consented` | 302 with code |
| A9 | code exchange | client auth; `resource` exact; tokens signed from pre-generated `family_id`/generation | code exists ∧ unconsumed ∧ PKCE + client/redirect/resource bind (inherited); `grant` `issued` ∧ not revoked; code TTL not elapsed | consume code; `grant` → `active` (`activated_at`, `grant_expires_at` = now + `approved_duration_s` `[O:2026-08-31]`, `family_id`); insert refresh family (`grant_id`, `grant_expires_at`) (D7) | `grant_activated` | tokens |
| A9′ | code exchange refused | as A9 | binding fails **or** grant not `issued` **or** TTL elapsed | code **consumed** on PKCE/client/redirect failure and on revoked/expired grant; **not** consumed on wrong `resource` `[S:mcp-sso §7.3]` (D6) | flow `token_refused_binding` / `token_refused_expired` / `token_refused_unknown` | `invalid_grant` / `invalid_target` |
| A10 | refresh rotation | client auth; successor tokens signed | family valid ∧ generation current ∧ `grant` `active` ∧ now < `grant_expires_at` | rotate generation; successor expiry ≤ `grant_expires_at` copied from `grant` | `token_refresh_rotated` | tokens |
| A10′ | rotation refused: replay or client mismatch | as A10 | consumed generation presented **or** client ≠ family client `[S:mcp-sso §7.4]` | **family revoked and `grant` → `revoked`** (`revoked_at`) — bounded lineage extends the reference's family-only theft response (D7) | `grant_revoked`, `token_refused_replay` | `invalid_grant` |
| A10″ | rotation refused: expired / unknown | as A10 | past `grant_expires_at` or unknown | lazy expiry if due (A14 semantics) | flow `token_refused_expired` / `token_refused_unknown` | `invalid_grant` |
| A11 | revoke grant (`grants revoke` **or** RFC 7009 on any token of the family, D9) | operator authority (#24) or valid revoke request | `grant` `issued` or `active` | `grant` → `revoked` (`revoked_at`); family revoked **if it exists** (issued grants have none yet) | `grant_revoked` | CLI / 200 |
| A12 | machine issuance (`client_credentials`) | client auth `[S:mcp-sso §9.4]`; `resource` exact; declaration exists for (client, route); requested ⊆ declared; **no matching explicit `deny` rule** `[O:2026-08-31]`; token signed | no `machine_tombstone` for (client, resource, current `declaration_digest`); reuse the `grant` `active` `kind=machine` for (client, resource) **only if** its `declaration_digest` equals current ∧ now < `grant_expires_at` (a past-due one is expired first, A14 semantics); **first-issuance race**: a losing insert on the one-active-per-(client, resource) uniqueness rolls back, discards its signed token, and deterministically retries once as the reuse path | else insert `grant` `active` (`kind=machine`, subject = client id, client, resource, scopes = declared ceiling, `approved_duration_s` = `maxDuration`, purpose = declared, `created_at` = `activated_at`, `grant_expires_at`, `declaration_digest`) | `grant_machine_issued` (on insert) / `grant_machine_reused`; refusals, in matrix order: flow `token_refused_client_auth` / `token_refused_no_declaration` / `token_refused_scope` / `token_refused_tombstone` / `token_refused_deny_rule` | access token (`exp` ≤ `grant_expires_at`, scopes = requested ∩ declared); refusal matrix, same order as the flow reasons: client auth fails → `invalid_client`; no declaration for (client, resource) → `invalid_target`; requested ⊄ declared → `invalid_scope`; tombstone → `invalid_grant`; deny rule → `invalid_grant` |
| A13 | revoke machine grant | operator authority | `grant` `kind=machine`, state `active` **or** `expired` (revocation is about the binding, not the row's state) | `grant` → `revoked` if `active`; insert `machine_tombstone` (client, resource, current `declaration_digest`) `[O:2026-08-31]` | `grant_revoked` (if active), `grant_machine_revoked` | CLI |
| A14 | expiry (lazy or sweeper `[#]`) | — | row past its bound, non-terminal | `preapproval` `awaiting` → `expired` + request → `resolved_expired`; `preapproval` `approved` → `expired` + request → `resolved_expired`; `grant_request` `allowed` → `abandoned` when its consent token expires; `grant` `issued` → `expired` when code TTL elapsed; `grant` `active` → `expired` | `preapproval_expired`, `request_resolved`, `request_abandoned`, `grant_expired` | — |
| A15 | retention purge (sweeper `[#]`) | — | terminal `grant_request` / `preapproval` / `grant` rows whose `terminal_at` (set by every terminal transition on all three records) is older than retention `[#]`; a `grant` only after its codes and family have expired; `grant_event` rows are never purged in v0 and may reference purged ids — events are history, a stated residual | batch-delete; idempotent — the delete's own predicate re-checks eligibility, so a race skips, never errors | flow `retention_purged` (count, oldest terminal timestamp) | — |
| E1 | response not delivered (post-commit) | — | — | none (rows already committed) | flow `response_not_delivered` | — |
| E2 | relay verification refusal (`/mcp`) | token verify | — | none | flow `token_refused_expired` / `token_refused_unknown`, relay line | 401 challenge |
| E3 | RFC 7009 unknown token | — | none | none | flow `unrecognized_token` | 200 |

## G7. Policy `[O:2026-08-30; vocabulary 2026-08-31]`

Input: subject (+ groups), client, route, scopes (post-ceiling), `requested_duration`,
purpose as opaque evidence. Output `allow` | `escalate` | `deny`.

- **Escalate by default; allow-list what may auto-approve.** AND-only `when` over
  scopes ⊆ set · duration ≤ ceiling · subject ∈ group · client ∈ list. **Explicit
  `deny` overrides `allow`** `[O:2026-08-31]`; a matching `allow` auto-approves; nothing matching
  escalates. No rule may reference purpose.
- **Machines** `[O:2026-08-31]`: rules apply **deny-only** — a matching `deny` refuses
  issuance (A12); `allow`/`escalate` are ignored because the declaration is the
  approval. **Boot contradiction check:** `validate`/`serve` refuse when a declaration
  is wholly denied by its route's rules ("declaration `X` on `/r` is denied by rule
  *n*"), so config cannot silently contradict itself.
- `policy_version` = SHA-256 of the RFC 8785 canonical JSON of `Route.spec.grant`.
- Unknown rule type/version/field refuses the whole config load. **v0 ships only the
  built-in decider**; an external decider (Edictum) is future work with no v0 config
  key — `contract.md §11` and `future.md` say the same.
- **Outage:** `temporarily_unavailable` on the redirect channel, event
  `decider_unavailable` `[O:2026-08-31]`.
- **Approvers**: route-configured group; narrow-only; never own request; audited.

## G8. Signing and commit discipline `[O:2026-08-31]`

Ids are generated in preflight; consent, access, and refresh tokens are **signed in
preflight** from those ids; the transaction commits rows and durable events; the
response is written after commit. Failure before commit ⇒ rollback, nothing consumed,
nothing issued. Failure to deliver an already-committed response ⇒ E1. This replaces
the reference's burn → sign → store order and its compensating family revocation
(D6/D7/D12) — a disclosed delta with its own fixtures.

## G9. Expiry propagation `[O:2026-08-31 activation start]`

Access `exp = min(now + accessTTL, grant_expires_at)`; refresh-family expiry ≤
`grant_expires_at`; successors copy from the `grant` row; access tokens carry
`grant_id` (D4). Before activation the only bound is the code TTL — never #8's
exchange-time boundary is "code TTL elapsed or grant revoked"; `grant_expires_at`
bounds rotation (A10″) and relay verification (E2).

## G10. Machine grants `[O:2026-08-31]`

Declared in config (`Gateway.spec.machineClients[]`, B1); ids in the B3 identifier
grammar. **Not** mcp-sso's stored-DCR `mcc_` model `[S:mcp-sso §17.2]` — inherited only:
`grant_type=client_credentials` grammar, client authentication, `resource`/`scope`
handling `[S:mcp-sso §9.4, §17.2 token clauses]` (D10a).

- **Token identity profile** (D10b): `sub` = `client_id` = the declared machine id;
  `gty=client_credentials` (inherited claim, drives the verifier's `machine`
  classification `[S:mcp-sso §17.2]`); `grant_kind=machine`; `grant_id`; `scope` =
  requested ∩ declared.
- **Lineage** (D10c): family-less; one `active` machine grant per (client, resource) at
  a time, reused only while its `declaration_digest` matches the current declaration
  (A12); `maxDuration` ≤ route `grant.maxDuration` ≤ hard ceiling `[#]`.
- **Rules deny-only + boot contradiction check** (G7) `[O:2026-08-31]`.
- **Sticky revocation** `[O:2026-08-31]`: A13 writes a tombstone bound to the per-route
  `declaration_digest` (G3); issuance refuses until the operator changes *that* binding.
  Expiry is not revocation.

## G11. Bounds on pending state

Open escalations per (subject, client, resource) `[#]`; identical `requested_hash`
dedupes (A3); awaiting pre-approvals expire after the approval window `[#]` (A14) so
the cap can never fill permanently; total pending `[#]`; terminal rows purged after
retention `[#]`. Total cap hit: `temporarily_unavailable` `[O:2026-08-31]`.

## G12. Audit durability classes `[O:2026-08-31]`

Two classes, one reason set (B7, each reason tagged): **durable events** are
`grant_event` rows committed in the operation's transaction (every G5 transition);
**flow events** — relay lines, token/identity/assertion refusals, caps, boot, port
failures, E-rows — are written best-effort to the JSONL sink with no store transaction.
Both fan into the same JSONL stream; loss is possible only for flow events and is
always announced (`audit_sink_failed` counter).

## G13. Verbs

`grants list` · `grants pending` · `grants approve <request-id> [--scopes …]
[--duration …]` · `grants deny <request-id>` · `grants revoke <grant-id>`. Authority
contract: #24.

## G14. Public outcomes

Escalation: 302 `error=access_denied&error_description=approval_pending&
request_id=<id>&iss=…&state=…` (dedupe returns the existing id). Any pending cap hit (A3″): `temporarily_unavailable`. Missing/unknown
`resource`: `invalid_target`. Exchange/rotation/tombstone/deny-rule refusals for
machines and exchanges: `invalid_grant`. Outage: `temporarily_unavailable`
`[O:2026-08-31]`. Everything else inherited (B7). Exact reasons only in audit.
