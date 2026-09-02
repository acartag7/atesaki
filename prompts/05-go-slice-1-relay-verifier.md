MODEL: grok-4.5   EFFORT: high   FALLBACK: gpt-5.6-terra   TOOL: Codex CLI or pi, in a fresh clone of ~/project/atesaki-core
MILESTONE: M3 (docs/roadmap.md). PRECONDITION (per-slice freeze, #55): the sections
below are SHA-pinned in this packet's header by the owner; packet 03 phase 1 fixtures
hash-locked; **mcp-sso's §8 verifier fixtures frozen with receipts** — the one
cross-lane input: this slice's verifier is inherited §07/§08 behavior, proven against
the shared corpus, never re-authored here.
PINNED: contract.md @ <sha> · contract-boundaries.md @ <sha> · contract-grants.md @
<sha> (G1 `resource`, G9, G12 only) · deltas.md @ <sha> · fixtures/MANIFEST.json @
<sha> · mcp-sso corpus <version> @ <sha>.
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

SCOPE — serial PRs, one behavior each, in this order (merge before the next):
1. `feat(runner): load and validate Atesaki fixtures` — loads a fixture in the
   Atesaki profile; validates `given.config` with the real parser; materializes
   `given.files` in a temp dir with the stated modes/owners/links; ports: clock,
   seeded randomness (the §19.2 HMAC-SHA256 counter stream, byte-exact), keys from
   `fixtures/keys`, recorded outbound HTTP only (an unrecorded call fails the fixture);
   chains and captures; exact comparison — status, headers (RE2 where stated),
   body, `then.events` by reason and class, `then.state` over the logical records
   (packet 02 phase 3 projection); absence assertions; hash check against
   `fixtures/MANIFEST.json`; a skipped locked fixture is a failure; zero expectations
   of its own. It also runs the frozen mcp-sso portable §8 fixtures (numeric-clause
   subset, same spine) — that run is how the verifier is proven. A test fails if any
   package outside the randomness port imports `crypto/rand`.
2. `feat(egress): profiles, proxy, CA per destination` — one `http.Transport` per
   profile; `fromEnv` = `http.ProxyFromEnvironment`; `none` = no proxy; URL = that
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
   multiplicity with `len(r.Header[...])`, never `Header.Get`. `Transfer-Encoding` +
   `Content-Length` together → `400`, pinned by fixture. Rate-limit identity = the B6
   client IP; budgets per B8.
4. `feat(verify): ES256 verifier, per-route metadata and challenge` — standard
   library only (`crypto/ecdsa`, raw `R||S` signatures, `crypto/rsa` unused here);
   `alg` never read from the token; `iss` exact; `aud` exact single string = route
   audience; `exp` with B8 skew; `scope` ⊇ `requireScope`; duplicate `Authorization`
   fails closed; the per-route challenge (D1, B7); PRM at the path-inserted location;
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
6. `feat(serve): wire the relay, flow audit, validate --deep` — boot order: validate
   in full → create the store directory `0700` if absent and check the path under B2
   → open the audit sink (B2) → listen; nothing before validation passes. Flow audit
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

`rehearse` is NOT in this slice (it runs the full AS, M6). No `mint` verb ships, ever:
the e2e test mints its own token with the test signing key.

HARD RULES: implement only what the locked fixtures and the pinned sections require;
no special-casing fixture strings; no invented library APIs (grep the module source);
stdlib first — every dependency justified in the PR with version, publish date, and
age against `tools/depage.py`; **never touch docs/ or fixtures silently** — follow the
review checkpoint in `prompts/README.md`; a gap is a PR comment, then continue; fail
closed on every ambiguity; no `value || default` on a security selector; `O_NOFOLLOW`
opens; explicit directory modes; `umask 077` at boot.

VERIFY: the Atesaki runner green on every locked slice-1 fixture, zero skips, and
green on the frozen mcp-sso §8 portable set; `serve` in front of a real local MCP
server (any stdio→HTTP server you can run — name it and its version) with a
static-header credential: one tool call succeeds with a test-minted token; the same
token at a second route is refused with that route's challenge; the upstream's
recorded request headers contain no bearer token; `validate --deep` against that
setup reports every probe.

DONE WHEN: all slice-1 fixtures green, zero skips; §8 portable green;
`validate`/`routes`/`serve`/`validate --deep` work on the three valid example configs;
packet-11 review clean; the parity status line published in the PR.

REPORT: fixtures passed/failed by id; the real MCP you relayed; every contract gap;
every place the contract was ambiguous and what you did NOT decide.
