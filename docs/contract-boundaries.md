# Trust-boundary contracts

**Status: DRAFT.** Companion to `contract.md` and `contract-grants.md`; owns the boundary
rules those pages only name. Tags as in `contract.md`, plus: **`[R]` = fail-closed
default confirmed by Arnold 2026-08-31** (changing one is a contract change); **`[P]` =
proposed by the design session 2026-08-31, awaiting Arnold — not decided**; `[#]` = a
number in B8. Parse, don't validate: every rule below produces a typed value once at the
boundary; interior code never re-checks it.

**Honesty note on B1:** the table below is the *reference* for the configuration
contract. The freeze artifact is a machine-readable schema generated from it and
validated against a mutation suite (open question #36). B1 is not claimed complete until
that artifact exists.

## B1. Configuration

**Envelope.** One YAML file, a stream of documents. Each document: `apiVersion:
atesaki/v1alpha1`, `kind: Gateway | Route`, `metadata.name` (identifier grammar, B3),
`spec`. Exactly one `Gateway`; one or more `Route`. Any other kind or apiVersion refuses
the whole load.

**`Gateway.spec`**

| Field | Type | Rule |
| --- | --- | --- |
| `externalBaseUrl` | URL | required; **origin only — empty path** `[O:2026-08-31]`; the issuer; B3 |
| `listen` | `host:port` | required; must be loopback when `identity.provider: console` |
| `signingKeyRef` | ref (B2) | required; the token-signing private key (ES256) |
| `allowedHosts[]`, `allowedOrigins[]` | [string] | required, non-empty (`contract.md §2`); matched against the effective host/origin (B6) |
| `trustProxyHeaders` | bool | default `false`; `true` requires non-empty `trustedProxies[]` else refuse |
| `trustedProxies[]` | [CIDR \| IP] | B6 |
| `identity` | tagged union on `provider` | `entra` \| `oidc` \| `header` \| `console`; a field belonging to an inactive variant refuses the load |
| `identity.registration` | `dedicated` \| `shared` | default `dedicated`; `shared` never implied (`contract.md §4`) |
| `identity.redirectUri` | URL | `entra`/`oidc`: required; must equal `externalBaseUrl` + `callbackPath` byte-exact |
| `identity.callbackPath` | path | default `/oauth/callback`; B3 path grammar; reserved-path collision refuses |
| `identity.clientId` | string | `entra`/`oidc`: required, non-blank |
| `identity.clientSecretRef` \| `identity.publicClient: true` | ref \| literal `true` | `entra`/`oidc`: exactly one of the two |
| `identity.tenantId` | string | `entra`: required |
| `identity.issuer` | URL | `oidc`: required |
| `identity.groupsClaim` | string | `entra`/`oidc` only, default `groups`; **refused** for `header` (its single groups selector is `assertion.groupsClaim`, B4) and `console` |
| `identity.assertion` | object | `header`: required (B4); refused for other variants |
| `identity.groupsToScopes` | map group → [scope] | optional; every scope ∈ union of route catalogs else refuse |
| `clients.knownCimd[]` | [path] | B2 file rules; each a valid CIMD document |
| `clients.allowLoopbackRedirects` | bool | default `true` |
| `clients.redirectAllowlist[]` | [URL \| origin] | required in every hosted mode (stateless DCR needs it, `contract.md §2`); exact origins/URLs, never `*` |
| `identity.egressProfile` | name | required for `entra`/`oidc` (every IdP call uses a named profile, `contract.md §7`); **refused** for `header` (its single JWKS egress selector is `assertion.keys.egressProfile`; a static `assertion.keys.jwksRef` needs none) and `console` |
| `identity.pairingCodeTtl` | duration | `console` only; default B8 |
| `egress.profiles.<name>` | `{proxy: fromEnv \| none \| URL, caBundleRef?: ref}` | name grammar B3 |
| `store.path` | path | required; B2 |
| `audit.path` | path | required (JSONL); B2 |
| `health.livePath`, `health.readyPath` | path | defaults `/healthz`, `/readyz`; reserved paths |
| `tokens.accessTtl`, `tokens.consentTtl`, `tokens.codeTtl` | duration | defaults B8; positive integer seconds |
| `machineClients[]` | [`{id, secretRef, purpose, maxDuration, routes[{path, scopes[]}]}`] | G10; `id` grammar B3; each route path must exist; scopes ⊆ that route's catalog; `purpose` B5 shape; `maxDuration` ≤ that route's `grant.maxDuration` ≤ hard ceiling |

**`Route.spec`**

| Field | Type | Rule |
| --- | --- | --- |
| `path` | path | required; B3 grammar; unique; no nesting (B3) |
| `mcpEndpoint` | path | default `/mcp` |
| `upstream.url` | URL | required; `http` \| `https`; B3; trusted target |
| `upstream.egressProfile` | name | required; must exist |
| `upstream.credential` | union on `type` | `{type: none}` \| `{type: static-header, header, scheme?, valueRef}`; unknown type refuses |
| `scopes.catalog[]`, `scopes.default[]` | [scope] | bounds `[S:mcp-sso §05]`: ≤128 entries, ≤256 bytes each; default ⊆ catalog; catalog non-empty |
| `authz.requireScope` | scope | required; ∈ catalog |
| `grant.maxDuration` | duration | optional; **default 8 h** (B8); ≤ hard ceiling |
| `grant.policy.rules[]` | [`{effect: allow \| deny, when: {scopesSubsetOf?, durationAtMost?, subjectInGroup?, clientIn?}}`] | G7 vocabulary; `when` non-empty; unknown key refuses the load; `deny` overrides `allow` `[O:2026-08-31]`; `policy_version` = SHA-256 of canonical `grant` block (G7) |
| `grant.approvers[]` | [group] | required unless every request is provably `allow`/`deny` — i.e. required when anything can escalate |

**Malformed input is refusal, never coercion:** wrong type, unknown field anywhere,
duplicate key, null for a required field, blank string after trimming, list-for-scalar,
scalar-for-list. Refusal names resource, field path, and rule. Config file size `[#]`;
YAML anchors, aliases, and custom tags refused.

## B2. Reference trust — `env:` and `file:`, and durable files

- `env:NAME` — `NAME` matches `[A-Z_][A-Z0-9_]*`. Unset **or empty** = missing = refuse.
  Whitespace-only value = refuse.
- **Filesystem invariants (all files Atesaki opens — secrets, store, audit, sidecars),
  on every open, not only creation** (restart is where the store becomes grant
  authority): open with `O_NOFOLLOW`; a
  symlink in the final component = refuse `[R]`; must be a regular file; link count
  exactly 1 `[R]`; owner = process user `[R]`; mode grants nothing to group/other
  (`0600` or stricter) `[R]`; the **parent directory** is owned by the process user and
  not group/other-writable `[R]`.
- **Secret-reference files only** (`file:` refs, `knownCimd`, key material): size `[#]`;
  exactly one trailing newline stripped; other leading/trailing whitespace = refuse; read
  once at boot, never re-read at request time. Store and audit files are not subject to
  these — they are opened for writing and grow.
- Directories Atesaki creates get mode `0700` `[R]` explicitly. Store and audit files are
  created with `O_EXCL`; **reopened** under the same rules as above. SQLite sidecar
  files (`-wal`, `-shm`) are held to the same owner/mode rules. An existing file that
  fails any rule = refuse to start, naming the rule — never "fix" permissions silently.
- A resolved secret is never logged, never echoed (errors name the *reference*), never
  placed in a token claim or audit record.

## B3. Canonicalization — URLs, paths, identifiers

One parser, one canonical form. Config values are **refused if not already
canonical** — never normalized: uppercase scheme or host, an explicit default port, a
trailing dot on a host, an empty path segment — each is a boot refusal naming the
spelling it wanted.

- **Portable host grammar** (the same for config URLs, `allowedHosts[]`, and the
  effective `Host`): a DNS name = RFC 1123 labels, lowercase `[a-z0-9-]`, no leading/
  trailing hyphen, ≤63 per label, ≤253 total, no trailing dot; **or** IPv4 = exactly four
  decimal octets, no leading zeros; **or** IPv6 = RFC 5952 canonical text (lowercase,
  maximally compressed) in brackets. Port = decimal 1–65535, no leading zeros, omitted
  when default. Any other spelling is refused, so Go and TypeScript cannot disagree on
  it (their default parsers are not the contract; this grammar is).
- **URLs.** `externalBaseUrl`: `https`, **origin only** `[O:2026-08-31]` (no path, query, fragment,
  userinfo) — **except** `identity.provider: console`, where `http://` on a loopback host
  is permitted `[O:2026-08-31]` (the five-minute path needs no local TLS); `upstream.url`: `http` \| `https`, path allowed, no query/fragment/userinfo;
  `redirectUri`, `issuer`, `jwksUrl`: `https`. Hosts lowercase, ASCII only (DNS name,
  IPv4, or bracketed IPv6) — Unicode/IDN hosts refused in v0 `[R]`.
- **Configured paths** (`Route.spec.path`, `mcpEndpoint`, `callbackPath`, health paths):
  begin with `/`; no trailing `/`; no empty, `.` or `..` segment; ASCII; no
  percent-encoding; segment grammar `[A-Za-z0-9._~-]+`.
- **Reserved paths** — refused as route paths or prefixes of them: `/oauth`,
  `/.well-known`, `callbackPath`, `health.*`, and every path Atesaki serves.
- **Route collision:** equal paths, or one a path-prefix of another (`/a`, `/a/b`),
  refuse. Two routes **may** share an upstream `[O:2026-08-31]`.
- **Audience** = `externalBaseUrl` + `Route.spec.path` + `mcpEndpoint`, byte-exact.
- **Metadata URLs (exact):** per route, RFC 9728 path-insertion —
  `/.well-known/oauth-protected-resource` + audience path (e.g.
  `/.well-known/oauth-protected-resource/playbook/mcp`); the authorization-server
  metadata lives at the origin: `/.well-known/oauth-authorization-server` and
  `/.well-known/openid-configuration` (one issuer). All are reserved; `atesaki routes`
  prints them; collision analysis includes them.
- **Inbound request-target:** parsed once. **Absolute-form** is accepted only when its
  scheme and authority byte-equal the effective authority tuple (B6); otherwise `400` — never silently
  discarded (two authorities in one request is a parser-differential surface). The **raw** path is
  first scanned for encoded separators (`%2F`, `%2f`, `%5C`, `%5c`) — any present ⇒
  **not routed** (the evidence is inspected before it is destroyed by decoding); the
  path is then percent-decoded exactly once; if the decoded path still contains `%`
  (double-encoding), an empty segment, or a dot segment, it is **not routed** (404). Matching is byte
  equality against canonical paths; no case folding, no trailing-slash tolerance.
- **Identifiers** (profile names, route names, machine-client ids): `[a-z][a-z0-9-]{0,62}`.

## B4. Rung 4 — the signed-assertion contract (`identity.provider: header`)

Generic signed-JWT profile (Cloudflare Access is one instance) `[O:2026-08-31]`:

- `assertion` = `{header, issuer, audience, alg, keys, subjectClaim?, groupsClaim?}`.
  `alg` ∈ {`ES256`, `RS256`} `[R]`, **pinned in config, never read from the token**.
  `keys` is a union: `{jwksUrl, egressProfile}` or `{jwksRef}` (B2 file). `kid` is
  required in the token and must match exactly one configured/fetched key; no `kid`, or
  a `kid` whose key's algorithm ≠ pinned `alg`, = refuse.
- **JWKS handling:** fetched through the named egress profile; document size `[#]`, key
  count `[#]`; refreshed on a schedule; a cached key set is used while stale up to a
  maximum stale interval `[#]`, after which verification refuses; no keys at boot =
  `serve` refuses to start. Rung 4 never degrades to "trust the header".
- **Header handling:** the configured header must occur **exactly once**; duplicate
  occurrences refuse (the §8.4 rule applied to identity). Missing assertion at the
  authorize identity leg = `access_denied` `401` on the **direct** channel — the
  resource-server bearer challenge belongs to `/mcp`, not here `[S:mcp-sso §9.3 step 1]`.
  On every path other than the identity leg, the configured header and the conventional
  identity headers (`cf-access-jwt-assertion`, `x-forwarded-user`, `x-forwarded-email`)
  are **dropped before routing**.
- **Claims:** `iss` byte-exact; `aud` contains the configured audience byte-exact; `exp`
  required and unexpired; `nbf` honored; `iat` not in the future beyond skew `[#]`;
  token size `[#]`. Subject = `subjectClaim` (default `sub`), which then crosses the
  inherited subject boundary `[S:mcp-sso §6.5]` — primitive, well-formed, no blank or
  replacement characters, bounded. Groups = `groupsClaim`, an array of strings, ≤128
  entries, ≤256 bytes each `[O:2026-08-31]`; a malformed groups claim yields **no groups** (empty
  ceiling → nothing can be granted), never a partial parse.
- **Replay:** a stolen assertion is valid until its `exp`; v0 has no assertion replay
  cache `[O:2026-08-31]` — a named accepted risk (`threat-model.md`), mitigated by the proxy's
  short assertion lifetimes.

## B5. Caps — every untrusted input is bounded

| Input | Cap | Failure |
| --- | --- | --- |
| Config file bytes; secret file bytes | `[#]` | boot refuse |
| Raw request-target (URL) bytes | `[#]` | `414`, before any parameter parsing |
| Request body, `/mcp` POST | `[#]` | `413`, body not read further |
| Request body, OAuth endpoints | inherited `[S:mcp-sso §05/§09 budgets]` | `413` |
| Header count / total header bytes | `[#]` | `431` |
| Purpose text | `[#]` UTF-8 bytes after trim; control characters refused | `invalid_request` |
| Open pending requests per (subject, client, route) | `[#]`; identical `requested_hash` dedupes (G11) | duplicate → `access_denied` + `approval_pending` with the existing id; non-duplicate at a full cap → `temporarily_unavailable` (A3″) |
| Total pending records | `[#]` | escalation refused `temporarily_unavailable` `[O:2026-08-31]` |
| Concurrent open streams per subject / per client | `[#]` | `429` + `Retry-After` |
| Non-stream upstream timeout | `[#]` | `502 upstream_unavailable` |
| SSE idle timeout | **none** (decided, `contract.md §6`) | client connection is the bound |
| Register / authorize / token rate per IP (after B6) | `[#]` | `429` |
| JWKS document / key count / assertion size (B4) | `[#]` | refuse |

A raw **byte cap** is enforced before parsing (body, headers, config, JWKS, assertion);
every structural cap (count, length, hops, entries) is enforced **before semantic use
or persistence** — a parsed value that exceeds its cap is refused and never acted on.

## B6. Forwarded-header and proxy trust

- `trustProxyHeaders: true` requires `trustedProxies[]`. Only `X-Forwarded-For`,
  `X-Forwarded-Proto`, and `X-Forwarded-Host` are recognized; RFC 7239 `Forwarded` is
  ignored in v0.
- **Effective authority is one tuple `(scheme, host, port)`:** HTTP/1.1 `Host` must
  occur exactly once (else `400`); HTTP/2 `:authority` is the host, and if a `Host`
  header is also present it must be byte-equal (else `400`); `X-Forwarded-Host` /
  `X-Forwarded-Proto` from a trusted peer replace host/scheme; otherwise the socket's
  facts hold. An absolute-form request-target must match the tuple in **scheme and
  authority** (B3), else `400`. Every competing spelling is rejected before routing.
- **Effective Origin:** the `Origin` header occurs at most once (else `400`); its
  value is `scheme://host[:port]` in the B3 host grammar with default port elided, or
  the literal `null`, treated as absent. Absent = allowed (non-browser MCP clients,
  `contract.md §6`); present = exact match against `allowedOrigins[]`.
- **Client IP algorithm:** if the socket peer is not in `trustedProxies`, the client IP
  is the socket peer and all forwarded headers are ignored. Otherwise: concatenate every
  `X-Forwarded-For` occurrence in order; split on commas; trim; **every** entry must
  parse as an IPv4/IPv6 literal, else `400`; at most `[#]` hops, else `400`; walk from
  the rightmost entry leftward skipping entries in `trustedProxies`; the first
  untrusted entry is the client IP; if every entry is trusted, the leftmost entry is the
  client IP. `X-Forwarded-Proto`/`-Host` are honored only from a trusted peer and must
  be single-valued and well-formed, else `400`.
- The **effective host** (after this step) is what Host/Origin validation and
  absolute-form comparison (B3) use. The issuer and all metadata URLs derive from
  `externalBaseUrl`, never from any header.
- Rate-limit identity = the client IP established here, never a bare header value.

## B7. Public error catalog, channels, and audit reason codes

Public responses are **non-oracular**: one shape for unknown, expired, revoked, and
not-yours. The channel rule is inherited: **direct** HTTP error until the redirect
destination is trusted (§9.3 steps 1–2), **redirect** after `[S:mcp-sso §9.3]`.

| Surface · channel | Public code | Status | When |
| --- | --- | --- | --- |
| authorize · direct | `access_denied` | 401 | no resolved subject (missing/invalid assertion, denied identity) |
| authorize · direct | `invalid_request` | 400 | ambiguous/malformed pre-validation parameters |
| authorize · direct | `invalid_client` / `invalid_redirect_uri` | 400 | inherited `[S:mcp-sso §9.3 step 2, §10]` |
| authorize · redirect | `unsupported_response_type` / `invalid_target` / `invalid_scope` | 302 | inherited `[S:mcp-sso §9.3 step 3]`; every redirect-channel row below is `302` |
| authorize · redirect | `invalid_request` | 302 | malformed purpose/duration (G4 G-a); duplicate or carrier-conflicting parameters |
| authorize · redirect | `access_denied` | — | user denial; policy `deny`; **escalation**: `error_description=approval_pending` + extension parameter `request_id` |
| authorize · redirect | `temporarily_unavailable` | 302 | decider unreachable `[O:2026-08-31]` |
| authorize · redirect | `temporarily_unavailable` | 302 | any pending cap hit — per-tuple non-duplicate or total (A3″, G11) `[D ← O:2026-08-31 total-cap outcome]` |
| authorize · direct/redirect | `invalid_target` | 400 / 302 | `resource` omitted or not a route audience (G1, D13) |
| register (POST) | `invalid_client_metadata` / `invalid_redirect_uri` / `invalid_request` | 400 | inherited `[S:mcp-sso §9.2]` |
| approve (POST) | `invalid_consent` / `invalid_origin` / `invalid_request` / `invalid_redirect_uri` / `invalid_grant` / `access_denied` | 400/403/302 | inherited `[S:mcp-sso §9.3 approve]`; redirect rows are `302` |
| revoke (POST) | — | 200 | RFC 7009: known or unknown token both answer 200; unknown is an admitted no-op, audited `unrecognized_token` `[S:mcp-sso §9, §13]` |
| revoke (POST) | `invalid_request` | 400 | duplicate/ambiguous form occurrences `[S:mcp-sso §9.4 occurrence gate]` |
| callback (redirect identity, bridge completion) | `server_error` | 302 | generic completion failure `[S:mcp-sso §17.11 bridge-completion row]`; the direct `500 internal_error` belongs to the claims-only completion Atesaki does not use |
| token | `invalid_request` / `invalid_client` / `unsupported_grant_type` / `invalid_scope` / `invalid_target` | 400/401 | inherited `[S:mcp-sso §9.4, §17.2]` |
| token | `invalid_grant` | 400 | unknown/expired/revoked code, family, or grant; grant not `issued`/`active`; one code for all |
| any OAuth | `internal_error` | 500 | generic (non-OAuth) throw from a port — never mapped to a client-auth error `[S:mcp-sso §14]` |
| any OAuth | `temporarily_unavailable` | 503 | store unreachable where the contract says fail closed |
| `/mcp` | `invalid_token` challenge | 401 | inherited shape `[S:mcp-sso §08]`, per-route metadata URL (D1) |
| `/mcp` | `insufficient_scope` challenge | 403 | inherited `[S:mcp-sso §08.3]` |
| `/mcp` | `forbidden` | 403 | Host/Origin not allowed (before body parse) |
| `/mcp` | `upstream_unavailable` / `upstream_auth_failed` | 502 | upstream unreachable / upstream refused the gateway's credential (`contract.md §6`) |
| any | — | 400 | absolute-form authority ≠ effective Host; malformed forwarded headers (B3, B6) |
| any | — | 404 | unrouted path (B3); empty body |
| any | — | 413 / 414 / 431 / 429 | caps (B5); `Retry-After` on 429 |

**Response headers relayed from an upstream** (exact allowlist; all else dropped):
`content-type`, `mcp-session-id`, `mcp-protocol-version`. `cache-control` is not
relayed `[R]`.

**Audit reason codes — the set for v0, closed at freeze** (a new code is a contract
change; a test rejects any reason constant absent from this list, and a coverage check
requires every G6 row to name only reasons from it). Class per G12: **D** durable
(committed with a state transition) or **F** flow (best-effort):
**D:** `request_allowed` · `request_denied_policy` · `request_escalated` ·
`request_unavailable` · `request_consented` · `request_abandoned` · `request_resolved` ·
`preapproval_approved` · `preapproval_denied` · `preapproval_claimed` ·
`preapproval_expired` · `preapproval_invalidated_stale` · `consent_denied` ·
`grant_issued` · `grant_activated` · `grant_expired` · `grant_revoked` ·
`grant_machine_issued` · `grant_machine_reused` · `grant_machine_revoked` ·
`token_refresh_rotated` · `token_refused_replay`.
**F:** `boot_refused` · `config_rejected` · `secret_ref_missing` · `identity_verified` ·
`identity_refused` · `assertion_duplicate` · `assertion_stale_keys` ·
`request_deduplicated` · `decider_unavailable` · `scope_ceiling_applied` (emitted by
G4 when the ceiling removes ≥1 requested scope) · `preapproval_claim_lost_race` ·
`response_not_delivered` · `retention_purged` · `unrecognized_token` · `token_refused_expired` ·
`token_refused_unknown` · `token_refused_binding` · `token_refused_tombstone` ·
`token_refused_deny_rule` · `token_refused_client_auth` · `token_refused_no_declaration` ·
`token_refused_scope` ·
`relay_forbidden_host` · `relay_forbidden_origin` · `relay_upstream_unavailable` ·
`relay_upstream_auth_failed` · `cap_exceeded` · `port_failure` · `audit_sink_failed`.

Free text never enters a **flow** event; purpose appears in exactly two durable events
(`grant_issued`, `grant_machine_issued`) and nowhere else in either stream.

## B8. Numbers

`[O:2026-08-31 — Arnold: "ok" to the defaults below]`; rows added after that ruling were confirmed the same day (bundle). `[D ← B1 unknown-is-refusal; interpretation flagged to Arnold — his "ok" decided the
numbers, not the key surface; his veto restores a `limits:` block]` Values that have a
B1 key (`tokens.*`, `identity.pairingCodeTtl`, `Route.spec.grant.maxDuration`) are
operator-configurable defaults; a route may raise
`maxDuration` only up to the hard ceiling. **Every other value is fixed in v0** — no
config key exists for it; a `limits:` block is future work only if evidence demands
it. Changing any value here is a contract change with a freeze-log entry.

| Marker | Value | Note |
| --- | --- | --- |
| access-token TTL | **10 min** | = the bounded-revocation lag |
| refresh-token lifetime | **= grant expiry** | no separate number (G8) |
| consent-token TTL | **5 min** | |
| authorization-code TTL | **2 min** | |
| purpose text | **512 UTF-8 bytes** | |
| `Route.spec.grant.maxDuration` default | **8 h** | hard ceiling **30 days** |
| pre-approval claim window | **24 h** | |
| config file / secret file | **1 MiB / 64 KiB** | |
| `/mcp` request body | **4 MiB** | |
| request headers | **100 / 64 KiB total** | |
| concurrent open streams | **16 per subject, 64 per client** | |
| non-stream upstream timeout | **60 s** | |
| per-IP budgets per 60 s | **register 30 · authorize 60 · token 120** | register inherited `[S:mcp-sso example]` |
| assertion clock skew | **60 s** | |
| open pending requests per (subject, client, route) `[O:2026-08-31]` | **3** | identical requests dedupe |
| approval window (awaiting pre-approval expires) `[O:2026-08-31]` | **24 h** | so a cap can never fill permanently |
| raw request-target bytes `[O:2026-08-31]` | **8 KiB** | `414` |
| total pending records (= `grant_request` rows in state `escalated`; the paired `preapproval` is not counted separately) `[O:2026-08-31]` | **10 000** | |
| terminal-record retention `[O:2026-08-31]` | **30 days**, then purged | |
| forwarded hops max `[O:2026-08-31]` | **10** | |
| JWKS document / key count `[O:2026-08-31]` | **64 KiB / 32 keys** | |
| JWKS maximum stale interval `[O:2026-08-31]` | **24 h** | |
| assertion token size `[O:2026-08-31]` | **8 KiB** | |
| groups claim `[O:2026-08-31]` | **≤128 entries, ≤256 bytes each** | mirrors scope bounds |
| expiry sweeper interval `[O:2026-08-31]` | **60 s** | lazy transition on read still applies |
| console pairing code TTL `[O:2026-08-31]` | **10 min** | loopback only |
