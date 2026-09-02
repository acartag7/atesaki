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

The contract pages are the authority. Where this page or a packet paraphrases a rule
to explain a gotcha, the contract's sentence wins and the paraphrase is not a rule.

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
| mcp-sso lane | the parity runner is landing as serial PRs on `main` (fifteen of PRs 374–392 merged, 388 and 391 open, 384 and 389 closed; `main` at `1f6911b`); draft PR 340 is superseded; only the two draft 8.4 fixtures exist, no `MANIFEST.json` |
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
   mcp-sso's guarded CIMD fetcher (§17.1.5: resolve, validate the address, dial that
   address directly, no redirects, bounded work) — the parity corpus pins it. *The
   catch:* that defense assumes a direct dial; through a corporate proxy the proxy
   resolves the name and the validated-IP dial is gone, so "inherit through the egress
   profile" contradicts the clause it inherits. *Recommend:* opt-in live fetch behind
   `clients.cimd.liveFetch: {egressProfile, allowedOrigins[]}` — exact `https`
   origins, never patterns; a document URL whose origin is not listed is refused
   before any network call. With an operator-chosen destination the DNS-pinning
   defense is not what stops SSRF (the allowlist is), so the fetch may use the named
   profile; when that profile is direct (`none`), the inherited validated-IP dial
   applies as written. The remaining inherited caps (document size, timeouts,
   redirects refused, content type, cache) are cited by clause. A new `deltas.md` row.
   Vendored documents remain the default and the only mode `rehearse` needs.
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
   *What happens:* with today's default volume behavior, Secret and ConfigMap mounts
   are symlinks into a root-owned directory and `subPath` mounts are root-owned
   regular files (an alpha feature in Kubernetes 1.37 adds ownership fields; off by
   default). B2 refuses both, by design. `env:` references work — for secrets and for
   `caBundleRef`, which is already a B2 reference. `knownCimd[]` is the one field
   with no reference form: it takes bare paths, so on the platform the config is
   shaped for it needs image-baked or init-copied files. *If open:* packet 08
   discovers this on the day it writes the recipe. *Recommend:* `knownCimd[]` entries
   become B2 references (`env:` or `file:`), one small contract change now; the
   recipe states, pinned to the Kubernetes versions it was tested on, that secrets,
   CA bundles, and CIMD documents arrive as `env:` from Secret keys, that `file:` is
   for hosts where the runtime user owns a `0600` regular file, and that the store
   path is a subdirectory Atesaki creates under the volume (mount roots fail the
   parent-directory rule).
7. **#57 — policy rules cannot name "any Codex install".** `[decide]`
   *What happens:* G7's `clientIn` matches exact client ids; Codex's id is a
   per-install URL. *If open:* a route rule "auto-approve read scopes for Codex" is
   unwritable; every user needs a rule edit. *Recommend:* add `clientOriginIn` (the
   origin of a CIMD client id URL; exact origins, never patterns) beside `clientIn`,
   AND-only like the rest; DCR clients have no origin and never match it.
8. **B8 configurability** (flagged 2026-08-31) — accept fixed numbers, no `limits:`
   block, until a deployment produces evidence. `[decide]`
9. **Client-matrix staleness window** — 90 days. `[decide]`
10. **Name check now, not at publish** (open question #9). The repo is already
    public as `acartag7/atesaki` and the module path is in `go.mod`; the namespace is
    claimed in fact. Run packet 09's name check as its own small step during M0 so a
    bad answer arrives before more work is stamped with the name. `[decide]`
11. **#60 — what happens when the rate limiter itself fails.** `[decide]`
    *What happens:* the reference limiter fails **open** when it throws on authorize,
    approve, token, and revoke (fail-closed only for stored registration); Atesaki
    inherits that by silence, and B8 names budgets for register, authorize, and token
    only. *If open:* an in-process fault removes abuse controls from the token-issuing
    paths, and the implementer guesses the approve and revoke budgets. *Already
    solved:* the decider-outage ruling (unreachable = `temporarily_unavailable`).
    *Recommend:* limiter error = `temporarily_unavailable` on every OAuth path (a
    delta row), and B8 gains approve = the authorize budget, revoke = the token
    budget.
12. **#61 — readiness and shutdown semantics.** `[decide]`
    *What happens:* B1 reserves `livePath`/`readyPath` and §14 asks the recipe to say
    how streams end at shutdown, but no sentence says what `readyz` checks or what
    `SIGTERM` does. *If open:* the platform routes traffic before the store is open,
    or a probe that includes upstream reachability takes a multi-route gateway out of
    the load balancer whenever one backend flaps (the seed's own recipe did this).
    *Recommend:* `livez` = the process serves; `readyz` = store open, signing key
    loaded, and (from M4) the identity JWKS fetched at least once — never upstream
    reachability, which is `validate --deep`'s job; `SIGTERM` stops accepting, drains
    non-stream requests for a bounded time (a B8 number), ends streams, then exits.
    Owned by M3 (packet 05) once ruled.
13. **Dependency floor: 15 days like mcp-sso, or the 5-day house minimum.** Atesaki
    is the same kind of boundary as mcp-sso with a handful of dependencies; slow
    intake costs users nothing. The plan assumes 15 (majors 30). `[decide]`

Two gaps the plan records without asking for a ruling yet, because the implementer
will hit them and the checkpoint rule (#52) already says what to do: #58 (rung-4 JWKS
refetch on an unknown `kid`, bounded) and #59 (what `validate --deep` may send to an
MCP upstream as a "real read"). One wording fix rides with packet 14: B4 says `alg` is
"never read from the token"; the executable rule is "the token's `alg` must equal the
configured one and match the key's type, and the allowlist never comes from the
token" (RFC 7515 requires processing the header; RFC 8725 forbids trusting it).

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
cross-lane dependencies and are named where they bite. The table below is the one
authoritative dependency view; the milestone sections and `prompts/README.md` are
derived from it.

| Packet | Produces | Consumes | Blocking rulings | Starts when | Done when |
| --- | --- | --- | --- | --- | --- |
| 13 | CI, protection, license, SECURITY.md, cooldown, two config fixes, STATE/ledger | — | LICENSE; dependency floor (item 13); name check (#9) | now | CI required on `main`; red PR cannot merge |
| 02 | B1↔parser drift test; G2 record types + generated record schemas; `knownCimd` refs | merged config code | #54 (phases 1–2); #56 (phase 3) | after 13 | drift diffs empty; records golden; phase 3 after packet 14 item 3 |
| 14 | contract sentences for #53, #5, #56, #57, PR-5 interpretations, #58, #59, #55, #60, #61, B4 `alg` wording, B8 note, matrix window | rulings | each item's ruling | any time a ruling lands | lint green; ledger receipts; fixture ids named for 03 |
| 12 | G13 authority text, B1 approvers row, audit fields, residuals | #24 ruling | #24 | after the ruling | lint green; fixture intent named for 03 phase 3 |
| 03 phase 0 | fixture profile + mutation suite; manifest/catalogue tooling | record schemas (02), packet 14 sentences | — | after 02 phase 2 and 14 | profile mutation suite green |
| 03 phase 1 | slice-1 fixtures (draft, then locked) | profile | — | after phase 0 | locked with the owner's read |
| 04 | threat model complete; negative matrix | 03 phase 1 ids | — | after 03 phase 1 | no attacker row without a rule; uncovered list published |
| 05 | runner, egress, pipeline, verifier, relay, serve, `--deep`, health/shutdown | locked phase-1 fixtures; frozen mcp-sso §8 | #55, #59, #61 | after 04 and the §8 freeze | phase-1 + §8 green, zero skips; review clean; real MCP named |
| 03 phase 2 | slice-2 fixtures | profile; packet 14 | — | alongside 05; locked before 06 | locked |
| 06 | store port + SQLite, AS, allow path, identity ports, `idp-request` | locked phase-2 fixtures; mcp-sso §07/09/10/11 or listed deferrals | #53, #5, #57, #58, #60 | after 05 and phase 2 lock | parity line; real sign-in; review clean |
| 03 phase 3 | slice-3 fixtures incl. packet 12's | profile; packet 12 | — | alongside 06; locked before 07 | locked |
| 07 | approvals, claim, CLI, machine clients, sweeper, outbox | locked phase-3 fixtures; packet 12 | #24 | after 06 and phase 3 lock | every G6 row green; `contract-v0-freeze` |
| 15 | `rehearse` + client profiles | runner; full AS | — | after 07 | onboarding step 4 true |
| 08 | recipe, image, kustomize, client matrix | `idp-request` (06), profiles (15) | matrix window | after 15 | recipe run once end to end |
| 09 | README, CHANGELOG, release workflow, sanitization in CI, listings | everything | — | after 08 | live verification named |

### M0 — Repo hardening

**You can now:** every PR runs the suite before you see it; `main` cannot be
force-pushed or merged red; a stranger can report a vulnerability; the tree says
what it is licensed under.

**Security first.** This is the milestone that makes the later ones checkable. Go has
no package-manager-level release-age gate, so the dependency cooldown needs two
layers at the same value: Dependabot `cooldown` (15 days, majors 30 — item 13) **and**
a CI test that reads each `go.mod` requirement's publish time from the module proxy
(`@v/<version>.info`, module path escaped per the module reference: uppercase →
`!x`) and fails on anything younger unless an exceptions file names a GHSA/CVE and
the minimum fixing version. The proxy timestamp is a risk signal, not provenance;
`go.sum` is the integrity ledger, not a lockfile — `go.mod` plus minimal version
selection is what pins the build list, verified with `go list -m all`. `govulncheck`
runs in CI from a pinned version. Builds use `-mod=readonly`. Race detector on. Linux
and macOS in the matrix (the B2 rules use `syscall.Stat_t` and `O_NOFOLLOW`; Windows
is out of v0 and the README says so). Actions pinned by commit; the fixture corpus is
supply-chain input from here on (its hashes are verified before anything is
materialized, M3).

**Tests.** Prove the gate bites: the CI PR carries one deliberately failing test in a
temporary commit, shows red, drops the commit, shows green. The dependency-age test
runs against a fake `.info` response that is one day old and must fail.

**Implement** (packet 13):

| PR | Content |
| --- | --- |
| `ci: build, vet, gofmt, race tests, govulncheck, contract lint` | matrix job plus one aggregate `ci` job that branch protection requires (matrix legs have suffixed names); `gofmt` output tested, not printed; protection: no force-push, no deletion, required check, linear history |
| `chore: license, security policy, gitignore` | LICENSE (owner picks; Apache-2.0 matches the rest of the family unless told otherwise), `SECURITY.md` with the disclosure channel, `.gitignore` for the binary and local scratch |
| `chore: dependency cooldown in CI and Dependabot` | the age test, the exceptions file format, Dependabot `cooldown` at the same floors |
| `fix(config): read the config file through one descriptor` | open → `fstat` for the regular-file check → read through `io.LimitReader(cap+1)` and refuse past the cap; the cap is enforced by the reader, so no size-then-read race exists (the secret-file path already opens once; the config path stats the name then reads the name — the TOCTOU sibling) |
| `refactor(config): one reserved-path source` | `main.go` and `validate.go` each carry the list today |
| `docs: record slices-before-freeze; refresh STATE; reorder packets` | ledger row with the PR 5 receipt; lane-1 rows corrected; `prompts/README.md` in this page's order |

Local, not a PR: delete the three merged branches; remove the dead `probe-a`/`probe-b`
entries from the Codex global config (they error on every Codex start).

**Gates.** Start: none. Done: CI required on `main`; a red PR cannot merge; lint and
tests green on the hardened tree.

**Decisions:** LICENSE choice; the dependency floor (item 13); the name check (item 10, open question #9).

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
registry the parser exposes. The registry records every accessor call, present or
absent, so an optional field no example uses is still registered; a variant branch
never executed shows up as a B1→parser miss. A field in one place and not the other
fails. This is the executable form of "B1 is not claimed complete until the artifact
exists". The record drift test compares G2's field lists **and** which G6 row sets
each state-dependent field, so ownership by state is checked, not just names.

**Implement** (packet 02, rescoped):

| PR | Content |
| --- | --- |
| `docs(contract): write the six PR-5 interpretations into B1` | contract-change PR (packet 14 item 5); one line each; owner confirms or reverses |
| `test(config): B1 to parser drift check` | the parser registers every accepted path with its requiredness at accessor-call time; the test reads B1; both diffs printed; empty or fail |
| `feat(records): G2 record types and generated schemas` | Go types for `grant_request`, `preapproval`, `grant`, `authorization_code` delta, `grant_event`, `machine_tombstone` with state-dependent presence encoded as typed unions; `schema/records/*.schema.json` generated by a golden test; RFC 3339 3-ms timestamps; `snake_case`; needed by 03 phase 0 |
| `feat(config): knownCimd entries are references` | after packet 14 item 3 (#56); `env:`/`file:`; the B2 file rules apply to `file:`; does not block M2 |

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
`cache-control` relay, metadata from config), caps on every input (B5), enforcement-
plane outage (decider, store, JWKS, and — #60 — the limiter), and the **build and
fixture supply chain**: a hostile fixture revision is a file-write primitive unless
the runner contains it; modules, actions, tools, and base images are pinned;
attestations bind the reviewed inputs. The threat model also records the residual
that a group removed after activation does not revoke the grant — refresh rechecks
nothing until expiry or revocation; the levers are a shorter `maxDuration` and
`grants revoke`.

**Tests.** The negative matrix is the deliverable: attacker × surface × fixture id,
rows without a fixture at the top. The fixture profile ships its own mutation suite
(an mcp-sso-shaped config, an unknown clause id, an unknown record field, an unknown
reason code — each refused by name).

**Implement:**

| Step | Packet | Content |
| --- | --- | --- |
| 1 | 14 | contract closure: #53 scope ceiling (new deltas row), #5 live fetch key and its inherited guards, #56 `knownCimd` refs, #57 `clientOriginIn`, the token response's narrowed `scope`; one small PR per item, lint green, fixtures listed as drafts |
| 2 | 12 | grants authority (#24 ruling): G13 authority text, B1 approvers row, audit fields per verb, the container residual; its fixtures are named by intent and written in 03 phase 3 |
| 3 | 03 phase 0 | the Atesaki fixture profile (`fixtures/schema/atesaki-fixture.schema.json`): mcp-sso spine + G/B/D/never clause ids + `given.config` as the Atesaki stream (validated by the real parser in the runner) + records from `schema/records/*` + B7 reasons with class + `given.files` for B2 boot fixtures under a containment contract (relative paths only, link targets inside the root, no ownership simulation, count and byte caps); the manifest's clause inventory includes B1 (the config refusal suite as `suite` evidence) and B8 (each number exercised at its boundary); its mutation suite |
| 4 | 03 phase 1 | slice-1 fixtures: relay §6, nevers 1, 3, 5, 7, B3 host and target grammar, B5 caps, B6 forwarded walk, B7 rows the relay reaches; `MANIFEST.json`/`CATALOGUE.md` generated; every fixture `draft` |
| 5 | 04 | threat model completed + `negative-matrix.md`; rows for later slices point at planned fixture ids and are counted as uncovered until those land |

**Gates.** Start: the ruling each packet-14 item names; packet 12 needs #24; 03 phase
0 needs 02 phase 2. Done: lint green; profile mutation suite green; every slice-1
fixture schema-valid; matrix published with its uncovered list.

### M3 — Slice 1: runner, relay, verifier, `validate --deep`

**You can now:** run `atesaki serve` in front of a real MCP server that today needs a
shared key. The key lives in the container. A request without a token gets the
route's own challenge; a token for another route is refused with that route's
challenge; a good token reaches the tool. `validate --deep` proves the IdP metadata,
each upstream, the store path, and the signing key are reachable through the
configured egress before anything is deployed. `readyz` answers only when the store
is open and the key is loaded; `SIGTERM` drains (#61). There is no login yet: the
demo token is minted by the end-to-end test with the test signing key — there is no
`mint` verb in v0 (§9's verb list is closed; a bounded mint would be a contract
change).

**Security first.**

- *Request pipeline order is the contract:* byte caps (B5) → one request-target parse
  (B3: absolute-form authority equality, raw separator scan, single decode) →
  effective authority (B6) → Host/Origin gate → route match → verifier → relay. A
  step that runs later than this list is a finding. Some of it happens in Go's parser
  before any handler runs, and the plan says so instead of pretending otherwise: two
  HTTP/1.1 `Host` fields are refused `400` by `net/http` (no audit line is possible);
  `Transfer-Encoding: chunked` beside `Content-Length` is resolved the RFC 9112 §6.3
  way — chunked wins and the length is dropped before dispatch; header **bytes** are
  capped by `MaxHeaderBytes` while the header **count** (B8) is counted in the
  handler because Go has no count limit. Fixtures pin those observables; the recipe
  requires the ingress to apply the same framing rule.
- *Relay is hand-built, not `httputil.ReverseProxy`:* the reverse proxy adds
  `X-Forwarded-For` and forwards headers by blocklist; §6 requires allowlists both
  ways. Upstream host and path come from config; the inbound request contributes a
  query string at most.
- *Duplicate headers fail closed:* Go's `Header.Get` returns the first value; the
  verifier must check `len(Header["Authorization"])` and the origin step
  `len(Header["Origin"])`; `Host` is promoted to `Request.Host` (HTTP/1.1 duplicates
  never reach the handler; on HTTP/2 a `Host` field beside `:authority` stays in
  `Header` and must byte-equal `Request.Host`).
- *Streams:* flush per event (`http.ResponseController`), no write deadline on a
  stream, non-stream upstream timeout 60 s (B8), client disconnect bound to the
  response context so a buffered POST is not aborted the moment its body closes.
- *Egress:* one transport per profile; `fromEnv` = `http.ProxyFromEnvironment`, `none`
  = nil proxy, URL = fixed proxy; `RootCAs` per profile, never the global pool;
  TLS ≥ 1.2; a proxy CONNECT failure is reported as `proxy <code> at <host>:<port>`.
- *Verifier:* ES256 with the stdlib (`crypto/ecdsa`, raw `R||S` signatures, not
  ASN.1); the token's `alg` must equal the configured algorithm and match the key's
  type, and the allowlist never comes from the token; `aud` exact single string,
  `iss` exact, `exp` with B8 skew, `scope` ⊇ `requireScope`. No general JWT library:
  the two algorithms the contract allows are the only code paths that exist.
- *Rate-limit identity:* the client IP from B6. A deployment whose ingress is not in
  `trustedProxies` puts every user in one bucket — `serve` warns at boot when it
  listens on a non-loopback address with `trustProxyHeaders: false`.
- *Audit lines* are JSON-encoded, one line, allowlisted fields; a formatter for
  untrusted input never throws.

**Tests.**

| Kind | What |
| --- | --- |
| Atesaki fixtures (03 phase 1) | relay rules, nevers 1/3/5/7, B3/B5/B6/B7 rows — zero skips |
| mcp-sso portable | the frozen §8 verifier set by fixture id, run by the same runner; every §7 token clause the verifier implements has a frozen upstream fixture **or** an Atesaki fixture marked `inherited` that quotes the same sentence — no inherited clause is implemented untested because the lane lags |
| Unit | grammars, header parsing, forwarded walk, target parsing — table-driven, one row per B3/B6 sentence |
| Real input | `serve` in front of a named local MCP server (any stdio→HTTP server you can run) with a static-header credential: one tool call succeeds with a test-minted token; the same token at a second route is refused with that route's challenge; the upstream's recorded request headers contain no bearer token (never 1) |
| Negative matrix | every slice-1 row flips from planned to a fixture id |
| Randomness lint | a test fails if any package outside the randomness port imports `crypto/rand` (fixture determinism) |

**Implement** (packet 05):

| PR | Content |
| --- | --- |
| `feat(runner): load and validate Atesaki fixtures` | manifest hash check **before** anything is materialized; `given.files` built under `os.Root` (relative paths only, no `..`, link targets inside the root, modes as stated, no ownership simulation — the owner-mismatch rule stays a unit test with an injected stat; count and byte caps); profile validation, chain ordering, clock/randomness/keys/recorded-HTTP ports, exact comparison, absence assertions, RE2 matchers; skipped locked fixture = failure |
| `feat(egress): profiles, proxy, CA per destination` | the one outbound layer; hop-naming errors |
| `feat(http): caps, authority, target parsing, host and origin gate` | the pipeline order above; B6 walk |
| `feat(verify): ES256 verifier, per-route metadata and challenge` | PRM at the path-inserted location, AS metadata documents at the origin, challenge per route (D1) |
| `feat(relay): allowlisted relay with streaming and cancel` | §6 in full |
| `feat(serve): wire the relay, flow audit, validate --deep, health, shutdown` | boot order: validate → open store path (create dir `0700`) → open audit sink → listen; `livez`/`readyz` and `SIGTERM` drain per #61; `--deep` per #59 |
| `test(e2e): real MCP behind the relay` | the named real input |

**Gates.** Start: M2 done; mcp-sso §8 portable fixtures frozen with receipts (the
cross-lane input); slice-1 sections SHA-pinned in the packet (#55). Done: all slice-1
fixtures green, zero skips; §8 portable green; every inherited clause this slice
implements has a frozen upstream or an `inherited` Atesaki fixture, or is named as
uncovered in the parity line and excluded from the capability claim; packet-11
review clean; real input named in the PR.

**Decisions:** #55, #59, #61.

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
- *Identity ports:* inherited, not reinterpreted. The id_token is verified (issuer,
  audience, nonce, expiry, signature under the fetched JWKS); groups come from the
  claim only. Entra, exactly as mcp-sso §17 says: a groups **overage** (the claim
  replaced by an overage marker) is an identity refusal with its named audit reason,
  never a Graph call and never a silent empty ceiling; group identifiers are object
  ids (GUIDs) only — display names are not unique and are refused as `groupsToScopes`
  keys. `idp-request` asks the IdP team for group filtering or app-assigned groups so
  overage cannot happen for the app.
- *Rung 4:* B4 in full; `alg` pinned per key; `kid` exact; identity headers stripped
  everywhere but the identity leg; JWKS through the egress profile, size and count
  caps, stale-interval refusal; #58 bounds the refetch on an unknown `kid`.
- *Console pairing:* loopback only, before any state write; the pairing code is
  printed, never audited.
- *CIMD live fetch* (#5, if allowed): only for a document URL whose origin is on the
  operator's exact allowlist; the inherited caps (document size, timeouts, redirects
  refused, content type, cache) cited by clause; the validated-IP dial applies when
  the profile is direct, and the allowlist is the SSRF control when it is a proxy.
- *Store:* the driver is chosen and pinned **first**, and its open semantics proven
  by an integration test before any G6 row lands: how it opens the main file, how it
  creates `-wal`/`-shm` (SQLite's Unix VFS gives sidecars the database file's mode,
  so a `0600` database yields `0600` sidecars — the test observes it for the driver
  actually used), whether it honors no-follow on open. B2 for the store is then
  enforced as: the directory is Atesaki's (`0700`, created by it, held open), the
  database file is created by Atesaki with `O_EXCL` and `0600` before the driver ever
  sees the path, and after the driver opens, the path is re-opened with `O_NOFOLLOW`
  and its inode compared with the file Atesaki created — a swap inside a directory
  only this user can write is the same-user case B2 does not defend against, and the
  contract says so. WAL, `busy_timeout`, `BEGIN IMMEDIATE` for every writing
  transaction (a deferred transaction that reads then writes fails with
  `SQLITE_BUSY_SNAPSHOT` under a concurrent writer and the busy handler cannot retry
  it), `synchronous=FULL`. If the pinned driver cannot meet a B2 sentence, that is a
  contract gap reported before the adapter merges, never a weakened check.
- *Two-phase discipline (G8):* ids and tokens are produced in preflight; the
  transaction holds only authoritative reads, CAS predicates, mutations, durable
  events; the response is written after commit; E1 on a lost response.
- *Canonical JSON for hashes (G3):* RFC 8785 — Go's `encoding/json` is not JCS (it
  escapes U+2028/2029 even with HTML escaping off, and JCS orders keys by UTF-16 code
  units); use a pinned JCS implementation after supply-chain review, or a purpose-
  written canonical serializer for the fixed G3 shapes, and test either against the
  RFC 8785 vectors. No "the encoder is close enough" argument.
- *Lazy expiry is this slice's, not M5's:* every cap and dedupe read in A3, and every
  A9/A10 read, first transitions past-due rows it touches (G5) inside the same
  transaction, or expired pending rows keep consuming the cap — a cheap denial of
  authorization. M5 adds the sweeper.
- *Limiter outage* (#60 as ruled): never fail open on a token-issuing path.
- *Timing-safe comparison* for client secrets, machine secrets, and the pairing code
  (compare digests with `subtle.ConstantTimeCompare`).
- *Rate limits* per client IP after B6 on register, authorize, token.

**Tests.**

| Kind | What |
| --- | --- |
| mcp-sso portable | the pinned corpus version + `MANIFEST.json` hash; the exact frozen portable fixture-id set this build passes, listed in the PR before code; zero skips; deferred ids listed with reasons |
| Atesaki fixtures (03 phase 2) | rungs §4 (each rung's boot refusals and acceptance; rung 4 duplicate header, unsigned header, wrong `kid`, stale JWKS, header stripping), D1, D3, D4, D5's allow branch, D6, D7, D11, D12, D13, never 6, A1, A2, A3 insert, A3′, A3″, A7, A8, A9, A9′, A10, A10′, A10″, A11 (RFC 7009 path), A14 lazy expiry for these rows, E1–E3, the scope-ceiling delta, #57 rules, the limiter-outage delta, Entra overage refusal; plus an `inherited` fixture for every §17 identity clause this slice implements that has no frozen upstream fixture |
| Store conformance suite | one table per G6 row this slice implements: atomicity, CAS, uniqueness, ordering-free comparison; both adapters; `serve` refuses memory |
| Crash tests | named failpoints between preflight and commit, and between commit and response, for A8 and A9; restart; exact state (nothing consumed / committed-with-E1) |
| Real input | Claude Code and Codex CLI (versions named) through Entra or generic OIDC to a tool call; Codex on two routes with different catalogs (the #53 case); the console rung on a laptop |
| `idp-request` | golden output per provider; contains no secret; lists what it does not need |

**Implement** (packet 06, rescoped):

| PR | Content |
| --- | --- |
| `feat(store): port, memory adapter, conformance suite` | the interface, the suite as the contract, memory passes |
| `feat(store): sqlite adapter` | pinned driver with its open/sidecar semantics proven by test first; WAL, pragmas; the B2 enforcement above; passes the suite under contention |
| `feat(as): metadata, stateless DCR, CIMD` | origin AS metadata; DCR per §9.2; vendored CIMD; live fetch behind its key |
| `feat(as): authorize steps 1–4, ceiling, purpose and duration` | §9.3 1–4 exactly, G-a/G-b carriers (singleton guard, signed bridge params, console round trip), #53 |
| `feat(policy): built-in rules, allow and deny` | G7 vocabulary incl. `clientOriginIn`; escalate by default; `policy_version` |
| `feat(grants): A1, A2, A3 insert, A3′, A3″, lazy expiry` | consent signed with `request_id` (D3); escalation ends the pass with `approval_pending`; the A14 lazy transition for every row this slice reads |
| `feat(grants): consent, exchange, rotation, revocation` | A7, A8, A9, A9′, A10, A10′, A10″, A11 via RFC 7009; G8 sign-before-commit; G9 expiry propagation; D4 claims |
| `feat(identity): entra, oidc` | redirect flow, id_token verification, groups; subject boundary §6.5 |
| `feat(identity): header assertion, console pairing` | B4 in full; loopback-only pairing |
| `feat(cli): idp-request` | per-provider templates, the does-not-need list |
| `test(e2e): real sign-in to tool call` | the named clients and versions |

**Gates.** Start: M3 done; 03 phase 2 fixtures locked; the mcp-sso §07/§09/§10/§11
portable set frozen **or** the exact not-yet-frozen ids listed as deferred in the
packet (the lane may lag; the parity line says so, nothing is skipped silently).
Done: parity line published; every deferred id named; every inherited clause this
slice implements has a frozen upstream or an `inherited` Atesaki fixture, or is named
uncovered and excluded from the capability claim; real sign-in shown; packet-11
review clean.

**Decisions:** #53, #5, #57, #58, #60.

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
| `feat(grants): sweeper and retention A14, A15` | 60 s interval over every row kind (the lazy path exists since M4), idempotent purge |
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
a persistent volume (one supported type named, RWO) mounted over a path the image
does not own, with `fsGroup` so the non-root process can create the store
subdirectory (`0700`) on first boot — proven on first boot **and** restart on a real
cluster; secrets and CIMD documents as `env:` from Secret keys (#56); CA bundles the
same way; image pinned by digest; SBOM and provenance attestations on the release;
NetworkPolicy egress derived from the documented ports; the ingress in
`trustedProxies` (or every user shares one rate-limit bucket) and applying the same
request-framing rule as Go's parser; **no path rewrite at the ingress** (audiences are
byte-exact); `livez`/`readyz` as ruled in #61; audit rotation that preserves every
line already written (reopen, never truncate) while flow-event loss stays the
accepted, counted residual it is (G12); backend reachability stated as the
obligation it is (§14).

**Tests.** Every command in the recipe run once against the real binary (named);
kustomize builds; the container runs under `readOnlyRootFilesystem: true`; a fresh
reader reaches a real sign-in from the recipe alone for console mode and one real IdP
mode; every other mode carries a tested-on date or the literal "UNVERIFIED" banner.

**Implement** (packet 08): `docs: deployment recipe`, `build: container image`,
`deploy: kustomize example`, `docs: client matrix`. The recipe consumes M4's
`idp-request` output and M6's client profiles; it designs no client flow of its own.

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
- **Packet 10 — the mcp-sso lane.** The runner is on `main` as serial PRs (fifteen of
  374–392 merged, two open); PR 340 is history. What Atesaki needs, in order: (1) the runner passes the portable 8.4 draft
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
| 5 | With today's default volume behavior, Kubernetes Secret/ConfigMap mounts are root-owned symlinks; `file:` refs are refused by design (#56) | M1/M7 | recipe uses `env:`, pinned to tested versions; `knownCimd` accepts refs |
| 6 | Volume mount roots are root-owned; the store must live in a subdirectory Atesaki creates | M4/M7 | boot fixture with `given.files`; kustomize example |
| 7 | SQLite's Unix VFS gives `-wal`/`-shm` the database file's mode; a database created with a foreign mode leaks through its sidecars, and a pure-Go driver may differ | M4 | Atesaki creates the db `0600` before the driver opens it; driver integration test observes sidecar modes; B2 re-check on every open |
| 8 | A deferred SQLite transaction that reads then writes fails with `SQLITE_BUSY_SNAPSHOT` under a concurrent writer, and the busy handler cannot retry it | M4 | `BEGIN IMMEDIATE` for every writing transaction; conformance suite under contention |
| 9 | Entra groups overage replaces the claim with a marker | M4 | inherited §17 fixture: identity refusal with its reason, no outbound call; `idp-request` asks for group filtering |
| 10 | Entra group identifiers are object ids; display names are not unique | M4 | inherited §17 rule: GUIDs only; a name as a `groupsToScopes` key is a boot refusal |
| 11 | Go `Header.Get` hides duplicate `Authorization`/`Host`/`Origin` | M3 | §8.4 portable fixture; B6 duplicate-Host fixture |
| 12 | `httputil.ReverseProxy` adds forwarded headers and forwards by blocklist | M3 | hand-built relay; never-1 fixture records upstream headers |
| 13 | Go `url.Parse` accepts non-canonical spellings | M3 | B3 grammar unit tables (already in config) reused for inbound |
| 14 | `Transfer-Encoding: chunked` + `Content-Length` together: Go's parser drops the length and honors chunked before any handler | M3/M7 | fixture pins the observable; recipe requires the ingress to apply the same rule |
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
| 29 | A group removed after activation does not revoke the grant; refresh rechecks nothing | M2/M7 | threat-model residual; recipe: shorter `maxDuration`, `grants revoke` |
| 30 | A hostile fixture revision is a file-write, chmod, or link primitive on every machine that runs the corpus | M3 | `os.Root` containment; manifest hash before materialization; hostile-path fixtures refused |
| 31 | The inherited limiter fails open when it throws (#60) | M4 | delta row + fixture: `temporarily_unavailable` |
| 32 | A readiness probe that includes upstream reachability takes a multi-route gateway out of the balancer when one backend flaps (#61) | M3/M7 | `readyz` semantics fixed by contract; recipe probes only those |
| 33 | Two `Host` fields are refused by Go before the handler; no audit line exists for them | M3 | fixture pins `400`; the negative matrix row cites the parser |
| 34 | `alg` confusion: pinning means "must equal the configured one", not "ignore the header" | M3 | packet 14 wording; §7/§8 fixtures |

## 6. Not in v0

`docs/future.md` is the list. Nothing above adds to it; the two seams this plan touches
(external decider, outbound OAuth client) stay ports with no config key.
