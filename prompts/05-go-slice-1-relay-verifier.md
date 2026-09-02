MODEL: grok-4.5   EFFORT: high   FALLBACK: gpt-5.6-terra   TOOL: Codex CLI or pi, in a fresh clone of ~/project/atesaki-core
MILESTONE: M3 (docs/roadmap.md). PRECONDITION (#55): the sections below are
SHA-pinned in this packet's header by the owner; packet 03 phase 1 fixtures merged
and read by the owner (§19 status `draft` until this runner passes them); **mcp-sso's
§8 verifier fixtures `frozen` in their files with the `receipt` object** — the one
cross-lane input: this slice's verifier is inherited §07/§08 behavior, proven against
the shared corpus, never re-authored here.
PINNED: contract.md @ <sha> · contract-boundaries.md @ <sha> · contract-grants.md @
<sha> (G1 `resource`, G9, G12 only) · deltas.md @ <sha> · fixtures/ @ <sha> ·
mcp-sso corpus @ <commit sha>. BLOCKING
RULINGS: #55, #59, #61, #63 (with values), and the JOSE library choice.
WHY: the first things that run — the Atesaki fixture RUNNER, then relay + verifier +
`validate --deep`, proven against it. Beyond the frozen §8 set it depends on NO
mcp-sso fixture — the fastest path to a binary you can point at a real MCP.

Read first, fully, at the pinned SHAs: docs/contract.md (all) · docs/contract-
boundaries.md (all) · docs/contract-grants.md G1 (`resource` rule), G9 (verifier uses
`grant_id`/`exp`), G12 (flow-event classes this slice emits) · docs/deltas.md D1, D4 ·
fixtures/** for nevers 1, 3, 5, 7, relay rules, B3/B5/B6/B7, the `boot` fixtures ·
docs/negative-matrix.md rows tagged slice 1 · docs/roadmap.md §M3 (the pipeline
order, the gotcha register rows 11–19, 24, 28) · prompts/README.md conventions ·
docs/open-questions.md #59 as ruled.

SCOPE — serial PRs, one invariant or protocol chain each, in this order (merge before the next):
1. `feat(runner): load and validate Atesaki fixtures` — runs every non-superseded
   fixture and reports by `status`; no manifest, lock, or hash to verify (#30, #50);
   loads a fixture in the Atesaki profile; validates
   `given.config` with the real parser; materializes `given.files` under `os.Root` in
   a fresh temp root per fixture, enforcing the profile's containment grammar
   (relative paths, no `..`, link targets inside the root, modes, count and byte
   caps; no ownership simulation) — a fixture that violates it is refused before any
   write; ports: clock,
   seeded randomness (the §19.2 HMAC-SHA256 counter stream, byte-exact), keys from
   `fixtures/keys`, recorded outbound HTTP only (an unrecorded call fails the fixture);
   chains and captures; exact comparison — status, headers (RE2 where stated),
   body, `then.events` by reason and class, `then.state` over the logical records
   (packet 02 phase 2 projection); absence assertions; a skipped fixture is a
   failure; zero expectations
   of its own. It also runs the frozen mcp-sso portable §8 fixtures (numeric-clause
   subset, same spine) — that run is how the verifier is proven. A test fails if any
   package outside the randomness port imports `crypto/rand`.
2. `feat(egress): profiles, proxy, CA per destination, references read once` —
   every `env:`/`file:` reference is resolved into a typed boot snapshot from the
   same descriptor `checkSecretFile` validated (B2 "read once at boot"); nothing
   reopens a path at request time. One `http.Transport` per profile; `fromEnv` = `http.ProxyFromEnvironment`; `none` = no proxy; URL = that
   proxy; `RootCAs` per profile from `caBundleRef`, never the global pool; TLS ≥ 1.2;
   failures name the hop (`proxy 403 at host:port`, `ECONNREFUSED`, certificate
   errors) — never a bare "fetch failed".
3. `feat(http): caps, authority, target parsing, host and origin gate` — in this
   order and no other: raw byte caps (`http.Server.MaxHeaderBytes` = B8 header total,
   `MaxBytesReader` on bodies, request-target cap → `414`) → one request-target parse
   (B3: absolute-form authority equality, raw separator scan, single decode,
   residual `%`/empty/dot segments → not routed) → effective authority (B6: `Host`
   exactly once, HTTP/2 `:authority` equality, forwarded headers only from
   `trustedProxies`, the rightward walk, hop cap, malformed → `400`) → Host/Origin gate
   (`Origin` at most once; absent allowed; present exact) → route match. Check header
   multiplicity with `len(r.Header[...])`, never `Header.Get` — except `Host`, which
   Go promotes to `Request.Host`: two HTTP/1.1 `Host` fields are refused `400` by the
   parser before the handler (fixtures pin the status; no audit line exists for it),
   and on HTTP/2 a `Host` field beside `:authority` stays in `Header` and must
   byte-equal `Request.Host`. `Transfer-Encoding: chunked` beside `Content-Length` is
   resolved by Go's parser (chunked wins, length dropped) — the fixture pins that
   observable; do not add a second parser. Header **bytes** via `MaxHeaderBytes`,
   which Go enforces with a 4 KiB read allowance and including the request line —
   measure the real threshold on the shipped Go version and pin that observable, do
   not assert an application byte count; header **count** (B8) counted in the
   handler → `431` (Go has no count limit). The pre-handler time envelope (#63 as
   ruled): `ReadHeaderTimeout`, body-read deadline, `IdleTimeout`, TLS handshake
   timeout, listener connection cap, unauthenticated per-IP budget on `/mcp` —
   proven with real-socket slow-header and slow-body tests, not `httptest`.
   Rate-limit identity = the B6 client IP; budgets per B8.
4. `feat(verify): ES256 verifier, per-route metadata and challenge` — one mature,
   pinned JOSE library (name it, version, publish date, age; a focused source and
   supply-chain review in the PR) behind a narrow verifier that admits exactly the
   contract's algorithms, key forms, and claims — compact JWS parsing (base64url
   framing, duplicate JSON members, `crit`, claim typing) is not hand-rolled; the
   token's `alg` must equal the configured algorithm and match the key's type — the
   allowlist never comes from the token (RFC 8725); `iss` exact; `aud` exact single
   string = route audience; `exp` **strict, no skew** (inherited §7.2; B8's 60 s is
   the rung-4 assertion skew, not an access-token grace); `scope` ⊇ `requireScope`;
   duplicate `Authorization` fails closed; the per-route challenge (D1, B7); every §7
   token clause implemented here has a frozen upstream fixture or an `inherited`
   Atesaki fixture; PRM at the path-inserted location;
   AS metadata documents at the origin (endpoints exist in slice 2; the documents are
   served now with the fields the contract fixes). Public key from `signingKeyRef`.
5. `feat(relay): allowlisted relay with streaming and cancel` — hand-built outbound
   request, NOT `httputil.ReverseProxy` (it adds `X-Forwarded-For` and forwards by
   blocklist); header allowlists both ways (§6, B7); upstream host and path from
   config, inbound contributes a query string at most; upstream 401/403 → `502
   upstream_auth_failed` with `WWW-Authenticate` stripped; unreachable → `502
   upstream_unavailable`; non-stream timeout 60 s; SSE flushed per event via
   `http.ResponseController`, no write deadline on a stream; client disconnect bound
   to the response context cancels the upstream; top-level JSON array refused;
   `DELETE` forwarded; `mcp-session-id` round-trips. The injected credential appears
   in nothing Atesaki authors (never 1).
6. `feat(serve): wire the relay, flow audit, validate --deep, health, shutdown` —
   boot order: validate in full → create the store directory `0700` if absent and
   check the path under B2 → open the audit sink (B2) → listen; nothing before
   validation passes. `livez`/`readyz` exactly per #61 as ruled (readiness per mode
   and capability: in this slice the store directory admissible under B2, the audit
   sink open, the key loaded — there is no database until slice 2; never upstream
   reachability); `SIGTERM` stops accepting, drains non-stream requests for the B8
   bound, cancels every stream's context, force-closes after the bound — Go's
   `Shutdown` alone waits forever on an open stream. Flow audit
   lines (G12 class F): JSON-encoded, one line, allowlisted fields, a formatter that
   never throws; `audit_sink_failed` counted and loud. Boot warning when listening on
   a non-loopback address with `trustProxyHeaders: false` (every user shares one
   rate-limit bucket). OAuth routes are NOT registered: a request to them is an
   unrouted path, `404` per B3. `validate --deep` per #59 as ruled: IdP metadata and
   JWKS through the identity profile, each upstream through its profile, the store
   path's writability by a create-and-remove probe in the store directory, the
   signing key parses; the §14 backend-reachability warning for plain `http://` on
   non-loopback, non-private addresses; never a state-changing call.
7. `test(e2e): real MCP behind the relay` — the named real input (below).

`rehearse` is NOT in this slice (it runs the full AS, M6). There is no `mint` verb in
v0 (§9's verb list is closed; a bounded one would be a contract change): the e2e
test mints its own token with the test signing key.

HARD RULES: implement only what the locked fixtures and the pinned sections require;
no special-casing fixture strings; no invented library APIs (grep the module source);
stdlib first — every dependency justified in the PR with version, publish date, and
age as human evidence (there is no age script; Dependabot cooldown is the gate); **never touch docs/ or fixtures silently** — follow the
review checkpoint in `prompts/README.md`; a gap is a PR comment, then continue; fail
closed on every ambiguity; no `value || default` on a security selector; `O_NOFOLLOW`
opens; explicit directory modes; `umask 077` at boot.

VERIFY: the Atesaki runner green on every slice-1 fixture, zero skips,
and green on the frozen mcp-sso §8 portable set; the exhaustion tests end at the #63
bounds; shutdown under live streams ends at the bound; `serve` in front of a real local MCP
server (any stdio→HTTP server you can run — name it and its version) with a
static-header credential: one tool call succeeds with a test-minted token; the same
token at a second route is refused with that route's challenge; the upstream's
recorded request headers contain no bearer token; `validate --deep` against that
setup reports every probe.

DONE WHEN: all slice-1 fixtures green, zero skips; §8 portable green; every
inherited clause this slice implements has a frozen upstream or an `inherited`
Atesaki fixture, or is named uncovered in the parity line and excluded from the
capability claim; `validate`/`routes`/`serve`/`validate --deep` work on the three
valid example configs; packet-11 review clean; the parity status line published.

REPORT: fixtures passed/failed by id; the real MCP you relayed; every contract gap;
every place the contract was ambiguous and what you did NOT decide.
