# Atesaki core contract

**Status: DRAFT. Freeze happens only when Arnold says the word.**

Provenance tags: `[O:date]` decision on record · `[E:Fn/Rn]` evidence item from
`atesaki/evidence/migration-playbook-mcp-setup-pain-2026-08-26.md` · `[K:§X]` starter-kit
conformance group (`seed/CONFORMANCE.md`, provenance in `seed/PROVENANCE.md`) ·
`[S:mcp-sso §N]` inherited mcp-sso contract clause (citations to be pinned to an
immutable mcp-sso SHA at freeze).

## 1. Roles

Atesaki plays three roles in one process:

1. **Authorization server** for the people signing in: it sends them to the company IdP,
   shows consent, and mints its own tokens. `[S:mcp-sso §09]`
2. **Resource-server front** for every route: it verifies its own tokens — signature,
   audience, scope — on every request. `[S:mcp-sso §08]`
3. **Relay** to the upstream MCP server, injecting the upstream credential outbound only.

Two walls that define the product:

- The IdP's tokens never reach a client. The client only ever holds a token Atesaki
  minted. `[S:mcp-sso §01]`
- The client's token never reaches an upstream. There is no passthrough mode, and no
  config key can create one. `[O:2026-08-16, invariant 3]`

## 2. Configuration

Two resource kinds, Kubernetes-shaped so a local file, a ConfigMap, and a future CRD are
the same YAML `[O:2026-08-16]`:

- `Gateway` — external base URL, identity (IdP block + ladder rung), client registration
  (CIMD documents, DCR, loopback redirects), egress profiles, store, audit.
- `Route` — path, upstream URL + egress profile + credential, required scope. The
  route's audience is derived: `externalBaseUrl + path + mcpEndpoint`.

Rules:

- **Unknown is refusal.** An unrecognized kind, apiVersion, or field rejects the whole
  config load. An engine that skips what it doesn't understand fails open by omission.
- **Secrets are references** (`env:`, `file:` in v0). A literal secret value anywhere in
  config is a boot refusal, not a warning. `[K:§J]`
- **Everything validates before any side effect.** No listener, no state file, no
  outbound call until every boot assertion passes. A missing secret reference is refused
  by the exact key name. `[E:F1/R6]` `[K:§K]`
- Two routes with the same audience are a boot refusal. `[K:§D]`
- An allowlist the active mode requires (hosts and origins always; redirects in
  every hosted mode) refuses to boot when empty — empty never means "allow all".
  mcp-sso §05 deliberately permits empty lists in specific modes; where a mode is
  inherited, its §05 mode-conditional semantics apply verbatim rather than this
  blanket sentence.

## 3. Tokens and identity

- Per-route audience is the security wall: a token minted for `/splunk` is refused at
  `/playbook` with the correct challenge, before any relay. `[K:§D]`
- Scopes come from the route's catalog. Group membership sets a **ceiling** on what a
  user can be granted (`groupsToScopes`); a scope the ceiling excludes is never granted,
  regardless of what the client requests. `[E:F8/R7]` `[S:mcp-sso §11]`
- Identity comes from the IdP's id_token. Atesaki never requests IdP-API scopes it did
  not configure and never synthesizes `api://…` scopes. `[E:F4/R2]`
- Redirect URI is **required, explicit config with no derived default**; boot fails if
  it does not sit on the gateway's own configured origin+path. `[E:F5/R1]`
- Token verification behavior — signature, audience, expiry, duplicate `Authorization`
  headers failing closed — is inherited clause-for-clause. `[S:mcp-sso §07, §08]`
  Parity is proven by the shared corpus, not asserted.

### 3.1 Grants → `contract-grants.md`

Everything about dispensing lives in `contract-grants.md`: what a grant is (G1), the
records (G2), hashes and digests (G3), the flow and its parameter carriers (G4 =
delta D5), states (G5), the **operation table** (G6), policy (G7), signing/commit
discipline (G8), expiry propagation (G9), machine grants (G10), pending-state bounds
(G11), audit durability classes (G12), verbs (G13), public outcomes (G14). One rule, one
place; nothing there is restated here.

## 4. The ladder of constrained modes

A constraint the product refuses to model gets implemented by hand outside the audit
boundary — that is the evidence of waves 7–9. So every rung is explicit, named config;
**no rung is ever entered by silent fallback**; changing rungs is a config edit.

| Rung | Config | Guardrails |
| --- | --- | --- |
| Dedicated app (default) | `registration: dedicated` | none needed; `idp-request` emits the minimal ask |
| Shared app | `registration: shared` (explicit) | requires its own client secret; boot-warns the known failure modes (consent scope, assignment-required conflicts) |
| Borrowed callback path | `identity.callbackPath` (any path, e.g. `/signin-oidc`) | callback is served by Atesaki at that path — never forwarded through another application `[E:waves 7–9]` |
| No IdP change at all | `identity.provider: header` behind an existing SSO proxy | **Signed assertions only** `[O:2026-08-30]`: the fronting proxy's identity assertion must be a cryptographically signed token (e.g. Cloudflare Access JWT) verified by signature, issuer, audience, and expiry `[S:mcp-sso Path B / identity ports]`. A plain unsigned identity header is refused outright — there is no "trusted proxy by network position" mode and no key that creates one. Refuses to serve without a valid assertion. |

## 5. Routes, discovery, and the challenge

- Each route serves its own protected-resource metadata at the path-inserted well-known
  location, and each route's `WWW-Authenticate` challenge points at **its own** metadata
  — never at the origin root. `[K:§A]` This is **delta D1** from the mcp-sso reference
  (`docs/deltas.md`) — intentional, public, and covered by Atesaki-side fixtures rather
  than a silent skip of the reference's origin-root fixtures.
- Routes are independently addressable; there is no merged tool catalog.
  `[O:2026-08-16, invariant 1]`
- Two routes **may point at the same upstream** `[O:2026-08-31]` — different
  audiences, scope catalogs, policies, and credential references onto one backend
  (`/splunk-read` vs `/splunk-admin`); a token for one is dead at the other.
- Path-based URLs, full stop (`https://host/playbook/mcp`). Hosted-client bugs against
  multi-segment paths are their bugs (anthropics/claude-ai-mcp#738). `[O:ROADMAP ch.1]`

## 6. Relay rules

All promoted from observed failures — none is advisory:

- Upstream answers 401/403 (the **gateway's** credential is bad) ⇒ the client sees a
  non-401 (`502 upstream_auth_failed`) and the upstream's `WWW-Authenticate` is
  stripped. A relayed challenge sends OAuth clients into a permanent re-login loop.
  `[K:§E]` — the splunk-gateway spike violated this; the rule exists so no composition
  can.
- Client-side auth failure still answers `401` + Atesaki's own challenge. The two
  failures are distinguishable. `[K:§E]`
- Request headers cross to the upstream by **allowlist** (MCP framing only:
  `mcp-protocol-version`, `accept`, `content-type`, `mcp-session-id`, `last-event-id`).
  Response headers relay by allowlist. A new header defaults to dropped. `[K:§E]`
- The upstream host and path come from trusted config only; the inbound request
  contributes at most a query string. `[K:§E, SSRF]`
- SSE streams through unbuffered; client disconnect cancels the upstream call; no idle
  timer severs a long-lived stream. `[K:§G]`
- Host and Origin are validated against allowlists before any body parsing, on every
  method the route serves. `[K:§C]` `[S:mcp-sso §08]`
- A top-level JSON array on `/mcp` is refused (current protocol removed batching).
- The injected credential appears in nothing **Atesaki authors**: no header it sets,
  no log, no audit record, no error, no response it generates. `[K:§J]`
  `[O:invariant 2]` A malicious or broken upstream reflecting its own credential in
  a relayed response body is outside a transparent relay's control — that is a named
  accepted risk (`threat-model.md`), not covered by this promise. `[O:2026-08-16]`

## 7. Egress — proxy-aware from byte one

- Every outbound byte (IdP metadata, JWKS, token endpoint, upstream calls) goes through
  one egress layer that honors `HTTPS_PROXY`/`NO_PROXY` **by default** — the inverse of
  the runtime defaults that produced `[E:F3]` and wave 9.
- Per-destination profiles: proxy or direct, CA bundle per destination — the IdP is
  usually reached through the corporate proxy, internal upstreams directly. `[K:§F]`
- A transport failure names the hop and the cause (`proxy 403 at <host>:<port>`,
  `ECONNREFUSED`, certificate errors) — never a bare "fetch failed". `[K:§I]`

## 8. Client registration

- CIMD documents are served from **vendored files** (`clients.knownCimd`) with no
  runtime fetch in v0. `[E:F3/R3]` (Whether live fetch exists at all, even opt-in, is
  an open question — not decided here.)
- DCR and CIMD are both on by default; loopback redirects with ephemeral ports are
  allowed by explicit default `[E:F9/R8]`. Redirect-allowlist semantics inherit
  `[S:mcp-sso §10]`.
- **DCR is stateless in v0** `[O:2026-08-30]`: the client id carries its own
  registration, nothing is stored, registrations survive restarts by construction, and
  the cross-grant scope-accumulation path (`deltas.md` D2) is structurally unreachable.
  Consequence, stated plainly: there is no per-client revocation — operators revoke
  **grants**, not clients. Stored DCR is not in v0 and has no config key.

## 9. The verbs

| Verb | Contract |
| --- | --- |
| `atesaki validate` | Pure validation: reads config, touches nothing else — no network, no store, no files written; refusals name the exact key. `[K:§K]` |
| `atesaki validate --deep` | Active diagnostics, distinct from validation: probes IdP metadata/JWKS, each upstream, store writability, signing key — through the configured egress profiles. Performs real reads and one store write-probe; never a state-changing call against the IdP or an upstream tool. Emits the **backend-reachability warning** (§14) for any plain-`http://` upstream on a non-loopback, non-private address. `[K:§K]` |
| `atesaki idp-request` | Emits the complete, reviewable IdP change request for the configured rung: app type, the single redirect URI, claims configuration, secret delivery — and the list of settings it does **not** need. Output is paste-ready for a ticket. `[E:F5, "I can't do this"]` |
| `atesaki rehearse` | Runs the full flow locally against a mock IdP — discovery → registration (CIMD and DCR) → authorize → callback → token → one `/mcp` call — per configured client profile and rung, before anything deploys. `[E:F6/R9]` |
| `atesaki routes` | Prints the machine-readable route + well-known path list, ready for ingress generation. `[E:F2/R5]` |
| `atesaki grants` | `list` · `pending` · `approve <request-id>` (narrow only) · `deny <request-id>` · `revoke <grant-id>` — `contract-grants.md` G13; authority contract: open question #24. |
| `atesaki serve` | The gateway. Boots only if `validate` would pass. |

## 10. Audit

Two durability classes, one reason set (`contract-grants.md` G12 `[O:2026-08-31]`): **durable
events** — every grant state transition, committed as a row in the operation's
transaction and fanned out — and **flow events** — relay lines per `/mcp` request,
token/identity/assertion refusals, caps, boot, port failures — written best-effort to
the sink. Never credential contents; purpose only in the issuance event. `[O:ROADMAP invariant 5]` Emission is guaranteed;
delivery is not: a failing sink does not fail the request, so lines can be lost — and
every loss is announced by a loud error signal and a counter. Silence about loss is the
forbidden outcome, not loss itself. `[O: starter kit, audit.failOpen]` (Rehearsal
`rehearse` proves the mock flow; only a live sign-in against the real IdP proves the
registration — the onboarding page says which step proves what.)

## 11. Extension seams (named now, built later)

Four ports, from ADR-001 D5: **inbound verifier** (the AS/verify half), **decider**
(**in v0 for grant requests** — `contract-grants.md` G7, **built-in per-route rules only**; an external decider such as Edictum is future work with no v0 config key;
unreachable = `temporarily_unavailable`, never allow `[O:2026-08-31]`; per-tool-call allow/block/ask through the same seam is future),
**credential resolver** (static-header in v0; token-exchange, per-identity, outbound
OAuth later), **event sink** (audit). Plus the **store port** (§13): SQLite and memory
adapters in v0, a conformance suite as the contract, other databases later
`[O:2026-08-31]`. `future.md` holds what plugs into them.

## 12. The nevers, each with the test a wrong build fails

1. **Never passthrough.** Test: a route configured every legal way; the upstream
   records its received headers; the client's bearer token appears in none of them.
2. **Never a literal secret in config.** Test: a config with an inline secret value
   boots; the only acceptable outcome is refusal naming the key — acceptance of the
   config is failure.
3. **Never relay an upstream challenge.** Test: upstream stub answers 403 + `WWW-Authenticate: Basic realm="x"`;
   the client must receive exactly `502` with body code `upstream_auth_failed`, **no**
   `WWW-Authenticate` header, and only B7-allowlisted response headers; any other
   status, the header's presence, or an extra header is failure.
4. **Never serve with an empty allowlist the active mode requires.** Hosts and origins
   are required in every hosted mode; redirect allowlists follow mcp-sso §05's
   mode-conditional rule (an empty list is a refusal only in the modes §05 says so).
   Test: for each mode, boot with the empty list that mode requires; refusal is the
   only pass; and for a mode where §05 permits empty, boot must succeed — both
   directions are asserted, so a blanket rule fails the test too.
5. **Never mint across the wall.** Test: obtain a real token on route A, present it on
   route B; anything but 401-with-B's-challenge is failure. (Real token, not a fake id
   — a made-up token only proves signature checking.)
6. **Never auto-degrade rungs.** Test: dedicated-rung config whose IdP app rejects the
   redirect; the gateway must surface the IdP error, not fall back to another rung;
   any successful login is failure.
7. **Never skip a frozen fixture.** Test: the corpus runner reports a skip; CI fails.
8. **Never dispense without a purpose and an expiry.** An acceptance **matrix**, not
   one test: purpose absent / empty / whitespace-only / control characters / over
   cap / wrong type; duration absent / zero / negative / fractional / overflow / over
   the route maximum — each refused at the request step (a completed grant is
   failure). Boundary: code exchange refuses when the code TTL has elapsed or the grant
   is revoked (expiry starts at activation, `contract-grants.md` G9 `[O:2026-08-31]`; the exchange also refuses a code whose PKCE/client/redirect binding fails — consuming it — and refuses without consuming on wrong `resource`, D6); at exactly
   `grant_expires_at`, rotation refuses and relay verification rejects the token. Caps: every access token's `exp` ≤ `grant_expires_at`; every refresh
   successor's expiry ≤ `grant_expires_at` — asserted on the real minted tokens, so the
   test never has to assume a token that outlives its grant exists. Races: a rotation
   whose commit happens after expiry or revocation is refused inside the atomic
   mutation. The public refusal is the non-oracular `invalid_grant`; the exact reason
   lives in audit.
9. **Never stamp a scope nothing granted.** Interactive: a minted token's scopes are
   always a subset of (route catalog ∩ the subject's group ceiling ∩ the signed consent
   decision). Machine: a subset of (declared scopes ∩ requested scopes), with no
   consent or ceiling axis — a separate matrix (mutants: scope outside the declaration;
   scope outside the request). Three independent mutants, each asserted on the real JWT `scope` claim:
   a scope absent from the catalog; a scope inside the catalog but outside the
   ceiling; a scope inside both but not in what the user approved. Plus lineage: a
   second grant by the same subject+client must not inherit the first grant's scopes
   (deltas.md D2). (The playbook anti-pattern — a scope stamped into a
   JWT that no authority granted.)

## 12a. Trust-boundary surfaces → `contract-boundaries.md`

The configuration reference table and malformed-input rules (B1 — the machine-readable schema is the freeze artifact, #36), `env:`/`file:` reference
trust (B2), URL/path canonicalization and route collision (B3), the rung-4 signed-
assertion contract (B4), caps (B5), forwarded-header and proxy trust (B6), and the
public error catalog with audit reason codes (B7) live in `contract-boundaries.md`.
Those pages are one contract; nothing there is restated here.

## 13. V0 scope pin `[O:2026-08-30, restoring the SHAPE decisions]`

**In v0:** identity — Entra, generic OIDC, header-mode with **signed assertions only**
(rung 4, `§4`), and **console pairing, loopback-only** `[O:2026-08-31]` — the
one-machine tutorial identity: the server prints a one-time code, the person pastes it
in a browser, and becomes the fixed `console-operator` subject; boot **refuses** any
non-loopback `externalBaseUrl`, issuer, or listen address before writing state; in
console mode — and only there — `externalBaseUrl` may be `http://` on loopback `[O:2026-08-31]`; boots
with a loud "tutorial, never multi-user" warning `[S:mcp-sso console-pairing]`;
registration — CIMD
(vendored) + DCR + loopback redirects; credentials — `none` and `static-header`;
grants — `contract-grants.md` G1–G14 in full, including **operator-declared machine
clients** (G10) `[O:2026-08-31]`;
registration — **stateless DCR only**; store — a **store port** `[O:2026-08-31]` (a Go
interface plus a conformance suite every adapter must pass, so another database can be
adopted later without touching the core); v0 adapters: **pure-Go SQLite** (default, one
file) and **memory** (passes the same suite; accepted by `rehearse`, refused by
`serve`); egress profiles; all six verbs; audit JSONL; **ship
artifacts** — one container image and one worked kustomize example (ingress from
`atesaki routes`, NetworkPolicy derivable from documented egress ports `[K:§K]`).
**Not in v0:** everything in `future.md`, human-in-the-loop approval for machine
clients, multi-replica stores.

## 14. Deployment recipe obligations `[E:R10]`

The recipe (a doc shipped with v0, not code) is a **per-mode matrix** — header mode
needs no IdP secret; generic OIDC may have no groups claim — and MUST state, per
irreducible part: **domain + TLS** — one hostname per gateway, who issues the cert,
and that the public URL is also the OAuth issuer, so it is chosen before the IdP app
is registered; **IdP registration** — client id, secret delivery path, the redirect
URI registered *before* deploy, groups-claim emission where the mode uses it (what
`idp-request` generates); **secret provisioning** — the exact key names the binary
reads, and how rotation/reload works; **network** — trusted-proxy allowlist for
forwarded headers; **backend reachability** `[O:2026-08-31]`: "nothing but the gateway
reaches the authless backend" is a **recipe obligation the product cannot enforce** —
the recipe requires NetworkPolicy or equivalent network position per platform and
says so plainly; Atesaki's only lever is advisory: `validate --deep` **warns** when an
upstream is plain `http://` on an address that is neither loopback nor
private/cluster-local; `NO_PROXY` and custom-CA installation; what NetworkPolicy can
and cannot enforce; **persistence** — where grant state and
the signing key live and what is lost on restart; **operations** — audit file
permissions and rotation, rate limits, upgrade behavior, and how active streams end at
shutdown; **client matrix** — per client: registration model, callback behavior,
version last tested, date — an entry older than the stated window fails the release
gate rather than "kept current" by intention. Every "never/always" in the recipe traces
to enforcing code or a test, like any other doc.
