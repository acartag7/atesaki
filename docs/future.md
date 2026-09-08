# Not this version

Each item names the seam it plugs into (`contract.md §11`). Nothing here constrains v0 beyond keeping the seam's port shape honest.

- **Local daemon** (the former chapter 0: Go, keychain, loopback per upstream). A sibling deployment shape of the same relay core. Demoted from core by `[O:2026-08-30]`.
- **Edictum decider, per tool call.** Allow, block, or ask on each `tools/call` through the decider port, holding an in-flight MCP request for a human. Grant-time decisions are already v0 (`contract-grants.md` G7) with built-in rules only. Wiring an external decider (an Edictum Ruleset through the same seam) is this item; the per-call hold, its stream semantics, and client timeouts are the future part. Unreachable means block (fail closed); at grant time the public outcome is `temporarily_unavailable` `[O:2026-08-31]`.
- **Outbound OAuth client** (former chapter 5). Atesaki as the OAuth client to OAuth-protected upstreams, a credential-resolver strategy.
- **Token exchange (RFC 8693) and per-identity credentials.** Further resolver strategies, already sketched in the starter-kit reference YAML.
- **Purpose-bound broker** (former chapter 2). Leases bound to subject, client, and run. Waits on a trustworthy run-identity source (sandbox identity).
- **Kubernetes operator and CRDs.** The config is already resource-shaped; a CRD is the same YAML plus a controller. Only if demand shows up. The CNCF gateways own that layer today.
- **EMA and ID-JAG grant support.** When enterprise IdPs beyond Okta ship it.
- **Agent mailbox (Dengon).** A separate product, not an Atesaki chapter (`~/project/ideas-new/mcp-mailbox-bootstrap.md`).

## Machine clients: deferred to v0.1

Owner ruling 2026-09-08 (#67): machine clients move to v0.1; signed proxy assertions stay in v0. Sweeper, retention purge and proxied CIMD are not deferred by this ruling. The following former machine contract is retained as deferred design, not v0 behavior. References within the retained excerpts describe their former locations; the excerpts below now own those deferred details. A v0.1 implementation needs its own contract and fixture review. The v0 configuration-boundary removal must precede the drift check and affected record generation; current Go acceptance is not changed by this docs PR.

### Deferred configuration and records

| Field | Type | Former rule |
| --- | --- | --- |
| `machineClients[]` | [`{id, secretRef, purpose, maxDuration, routes[{path, scopes[]}]}`] | G10; `id` grammar B3; each route path must exist; scopes ⊆ that route's catalog; `purpose` B5 shape; `maxDuration` ≤ that route's `grant.maxDuration` ≤ hard ceiling |

The deferred grant variant has kind=machine, subject=client_id, no refresh family, and a declaration_digest. Its machine_tombstone record carries client_id, resource, declaration_digest, revoked_at and revoked_by. Machine grants can be created active and transition to expired or revoked.

`declaration_digest` is per (client, route) `[D ← O:tombstone]`: SHA-256 over `atesaki-machine-declaration-v1\n` + canonical JSON of `{client_id, route_path, scopes (sorted), purpose_hex, max_duration_s}`. Reordering routes or editing another route never changes it, so a tombstone survives unrelated edits and clears only when *this* binding is deliberately changed.

### Deferred dispensing and revocation

Machine grants `[O:2026-08-31]`

Declared in config (`Gateway.spec.machineClients[]`, B1); ids in the B3 identifier grammar. Not mcp-sso's stored-DCR `mcc_` model `[S:mcp-sso §17.2]`. Inherited only: the `grant_type=client_credentials` grammar, client authentication, and `resource`/`scope` handling `[S:mcp-sso §9.4, §17.2 token clauses]` (D10a).

- **Token identity profile** (D10b). `sub` = `client_id` = the declared machine id; `gty=client_credentials` (inherited claim, drives the verifier's `machine` classification `[S:mcp-sso §17.2]`); `grant_kind=machine`; `grant_id`; `scope` = requested ∩ declared.
- **Lineage** (D10c). Family-less; one `active` machine grant per (client, resource) at a time, reused only while its `declaration_digest` matches the current declaration (A12); `maxDuration` ≤ route `grant.maxDuration` ≤ hard ceiling `[#]`.
- **Rules deny-only plus the boot contradiction check** (G7) `[O:2026-08-31]`.
- **Sticky revocation** `[O:2026-08-31]`. A13 writes a tombstone bound to the per-route `declaration_digest` (G3); issuance refuses until the operator changes *that* binding. Expiry is not revocation.

| # | Operation | Preflight | In-tx predicates | Mutations | Durable events | Output / public error |
| --- | --- | --- | --- | --- | --- | --- |
| A12 | machine issuance (`client_credentials`) | client auth `[S:mcp-sso §9.4]`; `resource` exact; declaration exists for (client, route); requested ⊆ declared; no matching explicit `deny` rule `[O:2026-08-31]`; token signed | no `machine_tombstone` for (client, resource, current `declaration_digest`); reuse the `grant` `active` `kind=machine` for (client, resource) only if its `declaration_digest` equals current ∧ now < `grant_expires_at` (a past-due one is expired first, A14 semantics); first-issuance race: a losing insert on the one-active-per-(client, resource) uniqueness rolls back, discards its signed token, and deterministically retries once as the reuse path | else insert `grant` `active` (`kind=machine`, subject = client id, client, resource, scopes = declared ceiling, `approved_duration_s` = `maxDuration`, purpose = declared, `created_at` = `activated_at`, `grant_expires_at`, `declaration_digest`) | `grant_machine_issued` (on insert) / `grant_machine_reused`; refusals, in matrix order: flow `token_refused_client_auth` / `token_refused_no_declaration` / `token_refused_scope` / `token_refused_tombstone` / `token_refused_deny_rule` | access token (`exp` ≤ `grant_expires_at`, scopes = requested ∩ declared); refusal matrix, same order as the flow reasons: client auth fails → `invalid_client`; no declaration for (client, resource) → `invalid_target`; requested ⊄ declared → `invalid_scope`; tombstone → `invalid_grant`; deny rule → `invalid_grant` |
| A13 | revoke machine grant | operator authority | `grant` `kind=machine`, state `active` or `expired` (revocation is about the binding, not the row's state) | `grant` → `revoked` if `active`; insert `machine_tombstone` (client, resource, current `declaration_digest`) `[O:2026-08-31]` | `grant_revoked` (if active), `grant_machine_revoked` | CLI |

The deferred machine-only audit set includes grant_machine_issued, grant_machine_reused, grant_machine_revoked, token_refused_tombstone, token_refused_deny_rule, token_refused_client_auth, token_refused_no_declaration and token_refused_scope. Machine rules are deny-only; the deferred configuration check refuses a declaration wholly denied by its route rules. A machine token's scope must be within both declared and requested scopes; its separate never-9 matrix belongs here. The operator's deferred onboarding declares routes, scopes, purpose and maximum duration for unattended agents.

### Deferred reference differences

| # | Difference | Reference | Deferred Atesaki design | Fixtures |
| --- | --- | --- | --- | --- |
| D10a `[O:08-31]` | Machine registration & credentials | Stored-DCR machine clients, generated `mcc_…` ids, out-of-band lifecycle (§17.2) | Declared in config (`machineClients[]`), B3 identifier grammar, secret by reference; only the token-endpoint grammar and client authentication are inherited | §17.2 registration/lifecycle → host |
| D10b `[D ← A12]` | Machine token claims | `sub`/`client_id`/`gty` triad over an `mcc_` id (§17.2) | `sub` = `client_id` = declared B3-grammar id; `gty=client_credentials` kept (verifier classification inherited); + `grant_id`, `grant_kind=machine`; scopes = requested ∩ declared | §17.2 exact claims → host |
| D10c `[D ← A12/A13]` | Machine lifecycle | Stateless issuance; no grant object | One `active` machine grant per (client, resource); issuance attaches or creates; revocation writes a tombstone bound to the declaration digest `[O:2026-08-31]` | Atesaki-only fixtures |
