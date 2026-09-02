# Atesaki v0 roadmap — feature by feature

**Status: DRAFT plan, 2026-09-02.** Sequence and scope live here; the contract pages
stay the authority for every rule this page names. Every `[decide]` is a call only the
owner makes; the plan is written under the recommendation so a session can start the
moment the ruling lands. The December schedule in the planning folder
(`atesaki-v0-roadmap-2026-09-01.html`) keeps its dates; this page replaces its
milestone list.

## 0. How to read a milestone

Each milestone answers the same six questions, in this order:

| Heading | What it must say |
| --- | --- |
| **You can now** | the thing you can demonstrate, in the operator's or the user's words — never "component X exists" |
| **Security first** | the trust boundary the milestone opens, the fail-closed edges, and the gotchas that have bitten real deployments |
| **Tests** | what proves it: which fixtures, which suites, which real input, which negative-matrix rows |
| **Implement** | the serial PR list and the dispatch packet |
| **Gates** | what must be true to start; what must be true to call it done |
| **Decisions** | the `[decide]` items it consumes |

Rules that every milestone inherits (`prompts/README.md`, `docs/quality-bar.md`):
one self-explanatory review unit per PR · fixtures for a slice are hash-locked before
that slice's code starts · every implementation PR gets a packet-11 adversarial
review before merge · never weaken a fail-closed control to pass a test · a contract
mismatch is a proposal plus an owner checkpoint (#52), never a silent edit · verify by
running a real input and name it.

## 1. Where we are (2026-09-02)

| Area | State |
| --- | --- |
| Contract set | drafted; ~50 rulings receipted in `decisions.md`; lint green; not frozen |
| Product code | config boundary and the two pure verbs merged (PR 5, PR 6); 3 valid examples, 71 refusal cases, build/vet/test green |
| Evidence | live discovery probe recorded (PR 4, open): D1 confirmed for Codex CLI 0.151.0 and Claude Code 2.1.257; #5 and #53 now carry evidence that forces an answer |
| mcp-sso lane | the parity runner is landing as serial PRs on `main` (nineteen merged 2026-09-01, two open); draft PR 340 is superseded; no fixture frozen yet, no `MANIFEST.json` |
| Repo hygiene | public repo; no CI workflow, no branch protection, no LICENSE, SECURITY.md, or `.gitignore` |
| Packets | 02 conflicts with the merged Go validator (#54); 05–09 assume a single freeze that practice has already moved past (#55); nothing owns `rehearse`, `idp-request`, or the operator-side k8s facts |

## 2. Rulings the plan is waiting on

Ranked by what they block. Each item: what it is, what happens, what breaks if it
stays open, where the same problem is already solved, and the recommendation the plan
is written under.

1. **#53 — Codex requests the union of advertised scopes on every route.** `[decide]`
   *What happens:* Codex reads the origin AS metadata's `scopes_supported`, not the
   route's PRM, and sends the union at `/authorize`. Under the inherited §9.3 step 3 a
   scope outside the route catalog is `invalid_scope`. *If open:* every Codex login
   fails on any gateway whose routes have different catalogs — `/splunk-read` next to
   `/splunk-admin`, the shape the README sells. *Already solved:* the group ceiling
   (`groupsToScopes`) narrows silently and emits `scope_ceiling_applied`. *Recommend:*
   make the route catalog part of the ceiling — requested scopes outside `catalog ∩
   group ceiling` are removed, `scope_ceiling_applied` fires, the narrowed `scope`
   is returned in the token response (RFC 6749 §5.1); an empty result is
   `invalid_scope`. It works regardless of which metadata document a client reads,
   and it reuses a mechanism that exists. It is a new `deltas.md` row (the reference's
   catalog-refusal fixtures become *host*). Follow-up probe: confirm Codex accepts a
   narrowed `scope`. Fallback if it does not: omit `scopes_supported` from the origin
   metadata (optional in RFC 8414) and re-probe.
2. **#5 — CIMD documents: vendored only, or opt-in live fetch.** `[decide]`
   *What happens:* Codex mints one CIMD document per install per server entry. *If
   vendored-only:* onboarding needs one document per user per route collected before
   that user can sign in, which contradicts "no pre-provisioning". *Already solved:*
   mcp-sso's guarded CIMD fetcher (bounded DNS, scripted transport, size caps) — the
   parity corpus pins it. *Recommend:* opt-in live fetch behind an explicit
   `clients.cimd.liveFetch: {egressProfile}` key, off by default; when on, the
   reference's guarded-fetch clauses are inherited verbatim through the named egress
   profile with a B5 document cap, `https` only, no redirects, private and loopback
   address ranges refused after resolution. Vendored documents remain the default and
   the only mode `rehearse` needs.
3. **#24 — who may run `atesaki grants`.** `[decide]`
   *If open:* packet 12, then M5, cannot start; approvals have no authority model.
   *Recommend:* the proposal as written (local CLI, the OS user's ability to open the
   store file is the authentication boundary, approvers as `{osUser, subject?}`
   entries, self-approval refused only where `subject` is present), with one
   amendment: in a container every `kubectl exec` is the same OS user, so the audit
   field `approver` would read `nonroot` for everyone. Add an optional `--approver
   <label>` recorded as *evidence, never authority*, and make the platform's exec
   audit trail a recipe obligation (§14). Named residual, not a hidden one.
4. **#54 — packet 02 is superseded by the Go validator.** `[decide]`
   *What happens:* the packet forbids Go and assumes no `validate` binary; both
   exist, and the Go refusal suite already is the mutation suite phase 2 describes.
   *If open:* the next session builds a second validator (JSON Schema plus Python) for
   the same input — the parser-differential class, and double maintenance. *Already
   solved:* `internal/config/testdata` names the rule per case. *Recommend:* drop the
   config JSON Schema; add a mechanical B1↔parser drift test in Go (field path, type,
   requiredness, both directions); write the G2 records as Go types and **generate**
   `schema/records/*.schema.json` from them with a golden test, because the fixture
   profile (packet 03) needs record schemas for `given.state`/`then.state`. The
   fixture runner validates `given.config` with the real parser, so no config schema
   is needed anywhere. Reverses #36's config half only.
5. **#55 — one freeze or a rolling one.** `[decide]`
   *What happens:* README already says "slices against the contract as it stands";
   `prompts/README.md`, `quality-bar.md`, and packets 05–07 still gate every line of
   Go on a single `contract-v0-freeze`. *If open:* the next implementer either
   refuses to write Go or invents a rule. *Recommend:* freeze **per slice**: when a
   slice starts, the sections it implements are SHA-pinned in its packet, their
   Atesaki fixtures hash-locked, and the mcp-sso citations for those sections pinned;
   the owner reads those pages, not all 1,400 lines at once. `contract-v0-freeze`
   becomes the tag applied when M5's parity line is green. Record the already-made
   "slices before freeze" decision (PR 5 merge as receipt) as a ledger row.
6. **#56 — B2 file rules on Kubernetes.** `[decide]`
   *What happens:* Secret and ConfigMap volume mounts are symlinks into a root-owned
   directory; `subPath` mounts are root-owned regular files. B2 refuses both, by
   design. `env:` references work. `knownCimd[]` and `caBundleRef` have no `env:`
   form, so on the platform the config is shaped for they need image-baked or
   init-copied files. *If open:* packet 08 discovers this on the day it writes the
   recipe. *Recommend:* `knownCimd[]` entries become B2 references (`env:` or
   `file:`), one small contract change now; the recipe states plainly that on
   Kubernetes secrets, CA bundles, and CIMD documents arrive as `env:` from Secret
   keys, `file:` is for hosts where the runtime user owns a `0600` file, and the
   store path is a subdirectory Atesaki creates under the volume (mount roots are
   root-owned and fail the parent-directory rule).
7. **#57 — policy rules cannot name "any Codex install".** `[decide]`
   *What happens:* G7's `clientIn` matches exact client ids; Codex's id is a
   per-install URL. *If open:* a route rule "auto-approve read scopes for Codex" is
   unwritable; every user needs a rule edit. *Recommend:* add `clientOriginIn` (the
   origin of a CIMD client id URL; exact origins, never patterns) beside `clientIn`,
   AND-only like the rest; DCR clients have no origin and never match it.
8. **B8 configurability** (flagged 2026-08-31) — accept fixed numbers, no `limits:`
   block, until a deployment produces evidence. `[decide]`
9. **Client-matrix staleness window** — 90 days. `[decide]`
10. **Name check now, not at publish.** The repo is already public as
    `acartag7/atesaki` and the module path is in `go.mod`; the namespace is claimed in
    fact. Run packet 09's name check as its own small step during M0 so a bad answer
    arrives before more work is stamped with the name. `[decide]`

Two gaps the plan records without asking for a ruling yet, because the implementer
will hit them and the checkpoint rule (#52) already says what to do: #58 (rung-4 JWKS
refetch on an unknown `kid`, bounded) and #59 (what `validate --deep` may send to an
MCP upstream as a "real read").

## 3. Milestones

```
M0 repo hardening ──► M1 config boundary (done + residuals)
                          │
                          ▼
                     M2 contract closure: rulings → packets 14, 12 → record types
                          → fixture profile + slice-1 fixtures (03 phase 0–1)
                          → threat model + negative matrix (04)
                          │
   mcp-sso §8 frozen ─────┼─────────────► M3 slice 1: runner + relay + verifier + validate --deep
                          │                     │
   mcp-sso §07/09/10/11 ──┼──► 03 phase 2 ──────▼
                          │              M4 slice 2: sign-in, store, the allow path, idp-request
                          │                     │
                          └──► 03 phase 3 ──────▼
                                         M5 slice 3: escalation, CLI, machine clients, sweeper, outbox
                                                │
                                                ▼
                                         M6 rehearse (mock IdP, client profiles)
                                                │
                                                ▼
                                         M7 deployment kit + recipe ──► M8 publish
```

Serial inside the Atesaki lane. The two dashed inputs from mcp-sso are the only
cross-lane dependencies and are named where they bite.

### M0 — Repo hardening

**You can now:** every PR runs the suite before you see it; `main` cannot be
force-pushed or merged red; a stranger can report a vulnerability; the tree says
what it is licensed under.

**Security first.** This is the milestone that makes the later ones checkable. Go has
no package-manager-level release-age gate, so the dependency cooldown needs two
layers: Renovate `minimumReleaseAge` (5 days, majors 30) **and** a CI test that reads
each `go.mod` requirement's publish time from the module proxy (`@v/<version>.info`)
and fails on anything younger unless an exceptions file names a GHSA/CVE and the
minimum fixing version. `govulncheck` runs in CI. Builds use `-mod=readonly`, the
lockfile is `go.sum`. Race detector on. Linux and macOS in the matrix (the B2 rules use
`syscall.Stat_t` and `O_NOFOLLOW`; Windows is out of v0 and the README says so).

**Tests.** Prove the gate bites: the CI PR carries one deliberately failing test in a
temporary commit, shows red, drops the commit, shows green. The dependency-age test
runs against a fake `.info` response that is one day old and must fail.

**Implement** (packet 13):

| PR | Content |
| --- | --- |
| `ci: build, vet, gofmt, race tests, govulncheck, contract lint` | workflow + required status check; branch protection: no force-push, no deletion, required check, linear history |
| `chore: license, security policy, gitignore` | LICENSE (owner picks; Apache-2.0 matches the rest of the family unless told otherwise), `SECURITY.md` with the disclosure channel, `.gitignore` for the binary and local scratch |
| `chore: dependency cooldown in CI and Renovate` | the age test, the exceptions file format, Renovate config |
| `fix(config): read the config file through one descriptor` | open → fstat → `io.LimitReader` (the secret-file path already does open-then-fstat; the config path stats the name then reads the name — the TOCTOU sibling) |
| `refactor(config): one reserved-path source` | `main.go` and `validate.go` each carry the list today |
| `docs: record slices-before-freeze; refresh STATE; reorder packets` | ledger row with the PR 5 receipt; lane-1 rows corrected; `prompts/README.md` in this page's order |

Local, not a PR: delete the three merged branches; remove the dead `probe-a`/`probe-b`
entries from the Codex global config (they error on every Codex start).

**Gates.** Start: none. Done: CI required on `main`; a red PR cannot merge; lint and
tests green on the hardened tree.

**Decisions:** LICENSE choice; #10 name check (item 10 above).

### M1 — Config boundary (done) + residuals

**You can now:** write `atesaki.yaml`, run `atesaki validate`, and be refused by the
exact resource, field, and rule for 71 malformed shapes; `atesaki routes` prints the
ingress path list. (PR 5, PR 6.)

**Security first.** Done: strict YAML (no anchors, aliases, tags, duplicate keys),
unknown-is-refusal, references only, the B2 file invariants on `file:` refs, B3
grammars, B4 shape, B5 config caps, the G7 boot contradiction check. Residual: the six
"interpretations to confirm" in PR 5's description are enforced by code but written
nowhere in the contract — each becomes a B1 sentence or is reversed.

**Tests.** The refusal suite stays the mutation suite. Add the drift test (#54): a Go
test parses the B1 tables in `contract-boundaries.md` (field path, type cell,
required/optional per variant) and compares them, both directions, with the field
registry the parser exposes; a field in one place and not the other fails. This is the
executable form of "B1 is not claimed complete until the artifact exists".

**Implement** (packet 02, rescoped):

| PR | Content |
| --- | --- |
| `docs(contract): write the six PR-5 interpretations into B1` | contract-change PR; one line each; owner confirms or reverses |
| `test(config): B1 to parser drift check` | the parser registers every accepted path with its requiredness; the test reads B1; both diffs printed; empty or fail |
| `feat(config): knownCimd entries are references` | after #56 lands; `env:`/`file:`; the B2 file rules apply to `file:` |
| `feat(records): G2 record types and generated schemas` | Go types for `grant_request`, `preapproval`, `grant`, `authorization_code` delta, `grant_event`, `machine_tombstone` with state-dependent presence encoded as typed unions; `schema/records/*.schema.json` generated by a golden test; RFC 3339 3-ms timestamps; `snake_case` |

**Gates.** Done: drift diffs empty; record schemas committed and reproducible.

**Decisions:** #54, #56.

### M2 — Contract closure before any authorization-server code

**You can now:** nothing new runs. Every rule the AS and grants slices will implement
exists as a sentence, and every sentence that a wrong build could violate has a
fixture id or is listed as uncovered by name.

**Security first.** This is where the OWASP pass happens on paper: injection (scope
names, purpose, subject bytes, JSON encoding of audit fields), broken auth (assertion
verification, client authentication, pairing code), authorization (audience wall,
ceiling, approve-then-swap via hash binding), sensitive exposure (secrets by reference,
non-oracular errors), misconfiguration (unknown-is-refusal, empty allowlists),
SSRF (upstream from config only; CIMD/JWKS through egress profiles with address-range
refusal), path traversal (B3 single decode, separator scan), deserialization (strict
YAML, bounded JSON), TOCTOU (CAS inside the transaction), open redirects (exact
allowlist), tenant isolation (per-route audience), crypto misuse (alg pinned per key,
`kid` exact), log injection (allowlisted fields, JSON-encoded lines), CSRF on consent
(origin check, signed consent token), ReDoS (RE2 only, and only in fixtures), timing
(digest comparison in constant time), proxy trust (B6), cache poisoning (no
`cache-control` relay, metadata from config), caps on every input (B5).

**Tests.** The negative matrix is the deliverable: attacker × surface × fixture id,
rows without a fixture at the top. The fixture profile ships its own mutation suite
(an mcp-sso-shaped config, an unknown clause id, an unknown record field, an unknown
reason code — each refused by name).

**Implement:**

| Step | Packet | Content |
| --- | --- | --- |
| 1 | 14 | contract closure: #53 scope ceiling (new deltas row), #5 live fetch key and its inherited guards, #56 `knownCimd` refs, #57 `clientOriginIn`, the token response's narrowed `scope`; one small PR per item, lint green, fixtures listed as drafts |
| 2 | 12 | grants authority (#24 ruling): G13 authority text, B1 approvers row, audit fields per verb, the container residual |
| 3 | 03 phase 0 | the Atesaki fixture profile (`fixtures/schema/atesaki-fixture.schema.json`): mcp-sso spine + G/B/D/never clause ids + `given.config` as the Atesaki stream (validated by the real parser in the runner) + records from `schema/records/*` + B7 reasons with class + `given.files` for B2 boot fixtures; its mutation suite |
| 4 | 03 phase 1 | slice-1 fixtures: relay §6, nevers 1, 3, 5, 7, B3 host and target grammar, B5 caps, B6 forwarded walk, B7 rows the relay reaches; `MANIFEST.json`/`CATALOGUE.md` generated; every fixture `draft` |
| 5 | 04 | threat model completed + `negative-matrix.md`; rows for later slices point at planned fixture ids and are counted as uncovered until those land |

**Gates.** Start: rulings 1–7. Done: lint green; profile mutation suite green; every
slice-1 fixture schema-valid; matrix published with its uncovered list.

### M3 — Slice 1: runner, relay, verifier, `validate --deep`

**You can now:** run `atesaki serve` in front of a real MCP server that today needs a
shared key. The key lives in the container. A request without a token gets the
route's own challenge; a token for another route is refused with that route's
challenge; a good token reaches the tool. `validate --deep` proves the IdP metadata,
each upstream, the store path, and the signing key are reachable through the
configured egress before anything is deployed. There is no login yet: the demo token
is minted by the end-to-end test with the test signing key — **no `mint` verb ships,
ever** (never 8 forbids a credential with no purpose and no expiry).

**Security first.**

- *Request pipeline order is the contract:* byte caps (B5) → one request-target parse
  (B3: absolute-form authority equality, raw separator scan, single decode) →
  effective authority (B6) → Host/Origin gate → route match → verifier → relay. A
  step that runs later than this list is a finding.
- *Relay is hand-built, not `httputil.ReverseProxy`:* the reverse proxy adds
  `X-Forwarded-For` and forwards headers by blocklist; §6 requires allowlists both
  ways. Upstream host and path come from config; the inbound request contributes a
  query string at most.
- *Duplicate headers fail closed:* Go's `Header.Get` returns the first value; the
  verifier must check `len(Header["Authorization"])`, the authority step
  `len(Header["Host"])`, the origin step `len(Header["Origin"])`.
- *Streams:* flush per event (`http.ResponseController`), no write deadline on a
  stream, non-stream upstream timeout 60 s (B8), client disconnect bound to the
  response context so a buffered POST is not aborted the moment its body closes.
- *Egress:* one transport per profile; `fromEnv` = `http.ProxyFromEnvironment`, `none`
  = nil proxy, URL = fixed proxy; `RootCAs` per profile, never the global pool;
  TLS ≥ 1.2; a proxy CONNECT failure is reported as `proxy <code> at <host>:<port>`.
- *Verifier:* ES256 with the stdlib (`crypto/ecdsa`, raw `R||S` signatures, not
  ASN.1), `alg` never read from the token, `aud` exact single string, `iss` exact,
  `exp` with B8 skew, `scope` ⊇ `requireScope`. No general JWT library: the two
  algorithms the contract allows are the only code paths that exist.
- *Smuggling:* a request carrying both `Transfer-Encoding` and `Content-Length` is
  refused `400` — pinned by a fixture, not assumed from the Go server's behavior.
- *Rate-limit identity:* the client IP from B6. A deployment whose ingress is not in
  `trustedProxies` puts every user in one bucket — `serve` warns at boot when it
  listens on a non-loopback address with `trustProxyHeaders: false`.
- *Audit lines* are JSON-encoded, one line, allowlisted fields; a formatter for
  untrusted input never throws.

**Tests.**

| Kind | What |
| --- | --- |
| Atesaki fixtures (03 phase 1) | relay rules, nevers 1/3/5/7, B3/B5/B6/B7 rows — zero skips |
| mcp-sso portable | the frozen §8 verifier set by fixture id, run by the same runner |
| Unit | grammars, header parsing, forwarded walk, target parsing — table-driven, one row per B3/B6 sentence |
| Real input | `serve` in front of a named local MCP server (any stdio→HTTP server you can run) with a static-header credential: one tool call succeeds with a test-minted token; the same token at a second route is refused with that route's challenge; the upstream's recorded request headers contain no bearer token (never 1) |
| Negative matrix | every slice-1 row flips from planned to a fixture id |
| Randomness lint | a test fails if any package outside the randomness port imports `crypto/rand` (fixture determinism) |

**Implement** (packet 05):

| PR | Content |
| --- | --- |
| `feat(runner): load and validate Atesaki fixtures` | profile validation, chain ordering, clock/randomness/keys/recorded-HTTP ports, exact comparison, absence assertions, RE2 matchers, hash check against `MANIFEST.json`; skipped locked fixture = failure |
| `feat(egress): profiles, proxy, CA per destination` | the one outbound layer; hop-naming errors |
| `feat(http): caps, authority, target parsing, host and origin gate` | the pipeline order above; B6 walk |
| `feat(verify): ES256 verifier, per-route metadata and challenge` | PRM at the path-inserted location, AS metadata documents at the origin, challenge per route (D1) |
| `feat(relay): allowlisted relay with streaming and cancel` | §6 in full |
| `feat(serve): wire the relay, flow audit, validate --deep` | boot order: validate → open store path (create dir `0700`) → open audit sink → listen; `--deep` per #59 |
| `test(e2e): real MCP behind the relay` | the named real input |

**Gates.** Start: M2 done; mcp-sso §8 portable fixtures frozen with receipts (the
cross-lane input); slice-1 sections SHA-pinned in the packet (#55). Done: all slice-1
fixtures green, zero skips; §8 portable green; packet-11 review clean; real input
named in the PR; parity line published.

**Decisions:** #55, #59.

### M4 — Slice 2: sign-in, the store, the allow path, `idp-request`

**You can now:** point Claude Code or Codex at `https://host/route/mcp`. It discovers
the route, registers (CIMD or DCR), the user signs in with the company login (Entra,
generic OIDC, a signed proxy assertion, or the loopback console), the agent states a
purpose and a duration, the consent page shows exactly those, and — where a route rule
says `allow` — the agent gets tokens and calls a tool. A request no rule allows ends
with `approval_pending` and a request id (nobody can approve yet; M5). `atesaki
idp-request` prints the ticket for the IdP team. RFC 7009 revokes a grant.

*Why this split:* the smallest deployable sign-in needs SQLite (`serve` refuses the
memory adapter) and the allow path (A1, A7, A8, A9, A10, A11). Building those here
puts the store port, the two-phase discipline (G8), and the conformance suite in place
one slice earlier, on the rows with the fewest branches. M5 adds the rows with humans
and machines in them.

**Security first.**

- *Parity is the rigor:* every §07/§09/§10/§11/§17 clause is inherited by fixture;
  every deviation is a `deltas.md` row or a bug. The dangerous corners — PKCE, `state`,
  consent-token signing and JTI consumption, code binding, refresh rotation and the
  theft response, redirect allowlisting — are the reference's, proven by the shared
  corpus.
- *Scope ceiling* (#53): effective scopes = requested ∩ catalog ∩ group ceiling;
  never 9 holds by construction; the narrowed `scope` goes back in the token response.
- *Identity ports:* the id_token is verified (issuer, audience, nonce, expiry,
  signature under the fetched JWKS); groups come from the claim only. Entra gotchas:
  a user in more than 200 groups gets no `groups` claim but an overage marker — that
  is **no groups**, an empty ceiling, never a Graph call; Entra emits group object
  ids, not names, unless the app is configured otherwise — `idp-request` says which
  and `groupsToScopes` keys must match what the token carries.
- *Rung 4:* B4 in full; `alg` pinned per key; `kid` exact; identity headers stripped
  everywhere but the identity leg; JWKS through the egress profile, size and count
  caps, stale-interval refusal; #58 bounds the refetch on an unknown `kid`.
- *Console pairing:* loopback only, before any state write; the pairing code is
  printed, never audited.
- *CIMD live fetch* (#5, if allowed): the reference's guarded fetcher — `https` only,
  no redirects, address-range refusal after DNS, document cap, cache TTL — through
  the named profile.
- *Store:* pure-Go SQLite (pinned, cooldown-aged), WAL, `busy_timeout`, `BEGIN
  IMMEDIATE` for every writing transaction (a deferred transaction that reads then
  writes fails with `SQLITE_BUSY_SNAPSHOT` under a concurrent writer and the busy
  handler cannot retry it), `synchronous=FULL`; `umask 077` at boot so the `-wal` and
  `-shm` sidecars SQLite creates from the database file's mode come out `0600`, then
  B2 re-checks them on every open; the store lives in a directory Atesaki creates
  with `0700` (#56).
- *Two-phase discipline (G8):* ids and tokens are produced in preflight; the
  transaction holds only authoritative reads, CAS predicates, mutations, durable
  events; the response is written after commit; E1 on a lost response.
- *Canonical JSON for hashes (G3):* RFC 8785 — Go's encoder HTML-escapes `<>&` and
  escapes U+2028/2029 differently from JCS; use a pinned JCS implementation or a
  restricted-grammar proof, and test against the JCS vectors either way.
- *Timing-safe comparison* for client secrets, machine secrets, and the pairing code
  (compare digests with `subtle.ConstantTimeCompare`).
- *Rate limits* per client IP after B6 on register, authorize, token.

**Tests.**

| Kind | What |
| --- | --- |
| mcp-sso portable | the pinned corpus version + `MANIFEST.json` hash; the exact frozen portable fixture-id set this build passes, listed in the PR before code; zero skips; deferred ids listed with reasons |
| Atesaki fixtures (03 phase 2) | rungs §4 (each rung's boot refusals and acceptance; rung 4 duplicate header, unsigned header, wrong `kid`, stale JWKS, header stripping), D1, D3, D4, D5's allow branch, D6, D7, D11, D12, D13, never 6, A1, A2, A3 insert, A3′, A3″, A7, A8, A9, A9′, A10, A10′, A10″, A11 (RFC 7009 path), A14 lazy expiry for these rows, E1–E3, the scope-ceiling delta, #57 rules |
| Store conformance suite | one table per G6 row this slice implements: atomicity, CAS, uniqueness, ordering-free comparison; both adapters; `serve` refuses memory |
| Crash tests | named failpoints between preflight and commit, and between commit and response, for A8 and A9; restart; exact state (nothing consumed / committed-with-E1) |
| Real input | Claude Code and Codex CLI (versions named) through Entra or generic OIDC to a tool call; Codex on two routes with different catalogs (the #53 case); the console rung on a laptop |
| `idp-request` | golden output per provider; contains no secret; lists what it does not need |

**Implement** (packet 06, rescoped):

| PR | Content |
| --- | --- |
| `feat(store): port, memory adapter, conformance suite` | the interface, the suite as the contract, memory passes |
| `feat(store): sqlite adapter` | pinned driver, WAL, pragmas, sidecar rules, passes the suite |
| `feat(as): metadata, stateless DCR, CIMD` | origin AS metadata; DCR per §9.2; vendored CIMD; live fetch behind its key |
| `feat(as): authorize steps 1–4, ceiling, purpose and duration` | §9.3 1–4 exactly, G-a/G-b carriers (singleton guard, signed bridge params, console round trip), #53 |
| `feat(policy): built-in rules, allow and deny` | G7 vocabulary incl. `clientOriginIn`; escalate by default; `policy_version` |
| `feat(grants): A1, A2, A3 insert, A3′, A3″` | consent signed with `request_id` (D3); escalation ends the pass with `approval_pending` |
| `feat(grants): consent, exchange, rotation, revocation` | A7, A8, A9, A9′, A10, A10′, A10″, A11 via RFC 7009; G8 sign-before-commit; G9 expiry propagation; D4 claims |
| `feat(identity): entra, oidc` | redirect flow, id_token verification, groups; subject boundary §6.5 |
| `feat(identity): header assertion, console pairing` | B4 in full; loopback-only pairing |
| `feat(cli): idp-request` | per-provider templates, the does-not-need list |
| `test(e2e): real sign-in to tool call` | the named clients and versions |

**Gates.** Start: M3 done; 03 phase 2 fixtures locked; the mcp-sso §07/§09/§10/§11
portable set frozen **or** the exact not-yet-frozen ids listed as deferred in the
packet (the lane may lag; the parity line says so, nothing is skipped silently).
Done: parity line published; every deferred id named; real sign-in shown; packet-11
review clean.

**Decisions:** #53, #5, #57, #58.

### M5 — Slice 3: dispensing complete

**You can now:** `atesaki grants pending` shows what waits; `grants approve <id>`
narrows and approves; the requester re-runs and the consent page shows the approved
values; `grants deny`, `grants revoke`; unattended agents are declared as machine
clients and get bounded, revocable, tombstone-guarded tokens; expiry fires on time;
terminal rows purge; every state transition is a durable event fanned out to JSONL
with a loss counter. This is the product promise, running.

**Security first.**

- *Approve-then-swap:* the claim (A6) is a CAS on `requested_hash` plus tuple; the
  grant is created with the **approved** values; freshness re-check (A6b) against
  current policy and ceiling.
- *Two concurrent claims:* exactly one wins; the loser proceeds as A3 in a new
  transaction (A6a) — proven with a deterministic barrier, never by repetition.
- *CLI authority* (#24 as ruled): the store file's permissions are the boundary;
  ids are exact, no enumeration on the public surface; `--approver` is evidence.
- *Machine issuance:* requested ⊆ declared; deny-only rules; tombstone on the
  per-route digest; one active grant per (client, resource) with the losing insert
  discarding its signed token and retrying once as reuse (A12).
- *Sweeper and lazy expiry:* exactly one event per expiry, inside the operation's own
  transaction; retention purge idempotent (A15).
- *Outbox:* durable rows committed with the change; JSONL fan-out best-effort with
  `audit_sink_failed` counted and loud; purpose appears in exactly two durable events
  and no flow line; free text never enters a flow event.
- *The refresh race clients will hit:* a retried refresh after a lost response is a
  replay and ends family and grant (A10′, inherited theft response plus bounded
  lineage). Stated in the recipe as "sign in again", not hidden.

**Tests.**

| Kind | What |
| --- | --- |
| Atesaki fixtures (03 phase 3) | A4, A5, A6, A6a, A6b, A12, A13, A14 sweeper, A15; never 8 and never 9 as the matrices written in §12; every durable reason produced by at least one fixture; every G5 state reached |
| Store conformance | the remaining rows on both adapters |
| Race | the two-runner claim with a barrier at the CAS; the machine first-issuance race |
| Crash | failpoints around A6 and A12 |
| Real input | a real escalation with a real client: request → pending → `grants approve` → re-run → tool call → `grants revoke` → refresh refused, access dies at TTL; a machine client via `client_credentials` on a route with a deny rule |
| Client matrix probe | how each client surfaces `approval_pending` + `request_id` (the user must be able to find the id) |

**Implement** (packet 07, rescoped):

| PR | Content |
| --- | --- |
| `feat(grants): approvals A4, A5 and the claim A6, A6a, A6b` | with the authority contract from packet 12 |
| `feat(cli): grants list, pending, approve, deny, revoke` | OS-user authentication, exact-id lookup, audit fields |
| `feat(grants): machine clients A12, A13` | D10a–D10c, tombstones |
| `feat(grants): sweeper and retention A14, A15` | 60 s interval, lazy path, idempotent purge |
| `feat(audit): durable outbox fan-out` | best-effort JSONL, loss counter, both classes in one stream |
| `test(e2e): escalation end to end` | the named real input |

**Gates.** Start: M4 done; 03 phase 3 locked; packet 12 landed. Done: parity line
green on the whole portable set; every G6 row has a green fixture; `contract-v0-freeze`
tag applied (#55).

**Decisions:** #24.

### M6 — `rehearse`

**You can now:** before deploying, run the whole flow on your laptop against a mock
IdP for each client you care about — discovery → registration → authorize → callback →
token → one `/mcp` call — per configured rung. It proves the gateway, the clients, and
the config agree. It cannot prove the company's IdP registration; onboarding step 6
does that.

**Security first.** The mock IdP and the rehearsal listener bind loopback only; the
memory adapter is accepted here and only here; `rehearse` never contacts the real IdP
or a real upstream (recorded exchanges only, the runner's egress port); client
profiles are recorded flows in the fixture profile, not live clients; output names
the profile, the rung, and the step that failed, never a secret.

**Tests.** `rehearse` is the runner in a trench coat: each client profile is a chain
in the Atesaki fixture profile executed against the composed binary. Golden output
per profile; a deliberately broken config fails at the named step.

**Implement** (packet 15): `feat(cli): rehearse with a mock IdP`, then one PR per
client profile (Claude Code CIMD, Codex CLI CIMD-per-install, DCR loopback, a hosted
client with a fixed callback).

**Gates.** Start: M5 done. Done: the onboarding page's step 4 is literally true.

### M7 — Deployment kit and recipe

**You can now:** deploy one container with one kustomize example; the ingress path
list comes from `atesaki routes`; the recipe tells you, per identity mode, exactly
what to ask the IdP team, which secret keys the binary reads, where state lives, what
is lost on restart, and what the platform must enforce that the product cannot.

**Security first.** Distroless static image, non-root uid, read-only root filesystem,
`/data` volume with the store in a subdirectory Atesaki creates (`0700`); secrets as
`env:` from Secret keys (#56); CA bundles the same way; image pinned by digest;
SBOM and provenance attestations on the release; NetworkPolicy egress derived from
the documented ports; the ingress in `trustedProxies` (or every user shares one
rate-limit bucket); **no path rewrite at the ingress** (audiences are byte-exact);
`livez`/`readyz` distinct; audit file rotation without losing lines; backend
reachability stated as the obligation it is (§14).

**Tests.** Every command in the recipe run once against the real binary (named);
kustomize builds; the container runs under `readOnlyRootFilesystem: true`; a fresh
reader reaches a real sign-in from the recipe alone for console mode and one real IdP
mode; every other mode carries a tested-on date or the literal "UNVERIFIED" banner.

**Implement** (packet 08): `docs: deployment recipe`, `build: container image`,
`deploy: kustomize example`, `docs: client matrix`.

**Gates.** Start: M6 done; the staleness window ruled. Done: as the packet says.

**Decisions:** client-matrix window.

### M8 — Publish

**You can now:** a stranger reads the README, follows the ten-minute path, and gets a
real sign-in; the trust artifacts (contract set, deltas, decisions ledger, fixtures,
threat model, negative matrix, parity line) are one click away; the name is checked
and recorded; the tree contains nothing employer-internal.

**Security first.** Sanitization is provable: a grep for hostnames, tenant ids, group
names, vault paths, employer names, and the private evidence folder runs in CI from
here on. Release artifacts carry checksums and attestations. `SECURITY.md` names the
channel and the response window. Every "never/always/cannot" in public copy traces to
a fixture or a receipt.

**Tests.** Live verification before promotion: two real clients (versions named)
through the published image following the recipe verbatim from a clean machine.

**Implement** (packet 09): as written, with the name check already done in M0.

## 4. Cross-cutting lanes

- **Packet 11 — adversarial implementation review** on every implementation PR,
  fresh context, "report everything with confidence", fail-closed sweep, sibling
  sweep, test-diff suspicion, invented-API grep, concurrency and partial failure for
  any G6 row touched. A PR is done when the review re-runs clean.
- **Packet 10 — the mcp-sso lane.** The runner is on `main` as serial PRs; PR 340 is
  history. What Atesaki needs, in order: (1) the runner passes the portable 8.4 draft
  and it freezes with a receipt; (2) `MANIFEST.json` + `CATALOGUE.md` + the CI hash
  gate; (3) the §08 verifier slice frozen (M3's input); (4) §07/§09/§10/§11 fixtures
  labeled portable/host against `deltas.md` including the new scope-ceiling row
  (M4's input). One session a day on that lane through September is the cheapest
  schedule insurance there is.
- **STATE.md** is refreshed in the PR that changes a lane's state — never a separate
  chore.

## 5. Gotcha register

Each row names the milestone that owns it and the proof that it is handled. A row
without a proof is a finding.

| # | Gotcha | Owner | Proof |
| --- | --- | --- | --- |
| 1 | Codex sends the union of AS-metadata scopes on every route (#53) | M2/M4 | ceiling fixture; real Codex on two routes |
| 2 | Codex mints one CIMD document per install per server (#5) | M2/M4 | live-fetch fixtures; onboarding without pre-provisioning shown |
| 3 | Codex client ids are per-install URLs; rules cannot name "Codex" (#57) | M2/M4 | `clientOriginIn` fixture |
| 4 | claude.ai connector rejects multi-segment paths (anthropics/claude-ai-mcp#738) | M7 | client-matrix row marked unsupported with the upstream issue |
| 5 | Kubernetes Secret/ConfigMap mounts are root-owned symlinks; `file:` refs are refused by design (#56) | M1/M7 | recipe uses `env:`; `knownCimd` accepts refs |
| 6 | Volume mount roots are root-owned; the store must live in a subdirectory Atesaki creates | M4/M7 | boot fixture with `given.files`; kustomize example |
| 7 | SQLite creates `-wal`/`-shm` sidecars with the database file's mode masked by the umask; a permissive umask or a foreign db mode leaks | M4 | `umask 077` at boot; B2 re-check on every open; boot fixture |
| 8 | A deferred SQLite transaction that reads then writes fails with `SQLITE_BUSY_SNAPSHOT` under a concurrent writer, and the busy handler cannot retry it | M4 | `BEGIN IMMEDIATE` for every writing transaction; conformance suite under contention |
| 9 | Entra groups overage (>200 groups) omits the claim | M4 | identity fixture: overage marker → empty ceiling, no outbound call |
| 10 | Entra emits group object ids, not names | M4 | `idp-request` states it; recipe example uses ids |
| 11 | Go `Header.Get` hides duplicate `Authorization`/`Host`/`Origin` | M3 | §8.4 portable fixture; B6 duplicate-Host fixture |
| 12 | `httputil.ReverseProxy` adds forwarded headers and forwards by blocklist | M3 | hand-built relay; never-1 fixture records upstream headers |
| 13 | Go `url.Parse` accepts non-canonical spellings | M3 | B3 grammar unit tables (already in config) reused for inbound |
| 14 | `Transfer-Encoding` + `Content-Length` together | M3 | fixture: `400` |
| 15 | ECDSA JWS signatures are raw `R\|\|S`, not ASN.1 | M3 | §8 portable fixtures verify real tokens |
| 16 | `aud` may be a string or an array in the wild; Atesaki mints and verifies a single string | M3 | inherited §07 fixture |
| 17 | Go's JSON encoder is not RFC 8785 (HTML escaping, U+2028/2029) | M4 | JCS vector test; fixture pins `requested_hash` bytes |
| 18 | Buffered POST closes its request body immediately; cancel must bind to the response | M3 | fixture: buffered POST completes; streaming client disconnect cancels upstream |
| 19 | Ingress not in `trustedProxies` puts every user in one rate-limit bucket | M3/M7 | boot warning; recipe network section |
| 20 | Ingress path rewrite breaks byte-exact audiences | M7 | recipe: no rewrite; `routes` output used verbatim |
| 21 | A retried refresh after a lost response is a replay that ends the grant (A10′) | M4/M7 | fixture; recipe "sign in again" sentence |
| 22 | Clients show `approval_pending` differently; the user must find `request_id` | M5 | client-matrix probe per client |
| 23 | Rung-4 unknown `kid` refetch can be driven by an attacker (#58) | M4 | bounded refetch fixture |
| 24 | Fixture determinism dies if any package reads `crypto/rand` directly | M3 | import lint test |
| 25 | Go has no package-manager cooldown | M0 | CI age test + Renovate |
| 26 | Windows lacks `O_NOFOLLOW`/`Stat_t` semantics | M0 | README: Linux and macOS only in v0; CI matrix |
| 27 | Pairing code in logs | M4 | audit fixture: absent |
| 28 | `validate --deep` "real read" against an MCP upstream is undefined (#59) | M3 | contract sentence + fixture before the verb ships |

## 6. Not in v0

`docs/future.md` is the list. Nothing above adds to it; the two seams this plan touches
(external decider, outbound OAuth client) stay ports with no config key.
