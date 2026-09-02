# Atesaki v0 roadmap: feature by feature

**Status: DRAFT plan, 2026-09-02.** Sequence and scope live here; the contract pages stay the authority for every rule this page names. Every `[decide]` is a call only the owner makes; the plan is written under the recommendation so a session can start the moment the ruling lands. The December schedule in the planning folder (`atesaki-v0-roadmap-2026-09-01.html`) keeps its dates; this page replaces its milestone list.

## 0. How to read a milestone

Each milestone answers the same six questions, in this order:

| Heading | What it must say |
| --- | --- |
| **You can now** | the thing you can demonstrate, in the operator's or the user's words, never "component X exists" |
| **Security first** | the trust boundary the milestone opens, the fail-closed edges, and the gotchas that have bitten real deployments |
| **Tests** | what proves it: which fixtures, which suites, which real input, which negative-matrix rows |
| **Implement** | the serial PR list and the dispatch packet |
| **Gates** | what must be true to start; what must be true to call it done |
| **Decisions** | the `[decide]` items it consumes |

The contract pages are the authority. Where this page or a packet paraphrases a rule to explain a gotcha, the contract's sentence wins and the paraphrase is not a rule.

Rules that every milestone inherits (`prompts/README.md`, `docs/quality-bar.md`): one self-explanatory review unit per PR, where a unit is an invariant or a protocol chain with all its fixtures, never one test case · a slice's fixtures are merged and read by the owner before that slice's code starts · every implementation PR gets a packet-11 adversarial review before merge · never weaken a fail-closed control to pass a test · a contract mismatch is a proposal plus an owner checkpoint (#52), never a silent edit · verify by running a real input and name it.

## 1. Where we are (2026-09-02)

| Area | State |
| --- | --- |
| Contract set | drafted; ~50 rulings receipted in `decisions.md`; lint green; not frozen |
| Product code | config boundary and the two pure verbs merged (PR 5, PR 6); 3 valid examples, 71 refusal cases, build/vet/test green |
| Evidence | live discovery probe recorded (PR 4, open): D1 confirmed for Codex CLI 0.151.0 and Claude Code 2.1.257; #5 and #53 now carry evidence that forces an answer |
| mcp-sso lane | the parity runner is **done** and runs in CI (`main` at `63ed987`; verified locally: the two 8.4 fixtures pass through Fastify, Express, and Hono, zero skips); PR 340 closed; nothing is frozen yet, both fixtures are `draft` with no `receipt`, and the corpus is two fixtures; no manifest, catalogue, or hash gate is to be built (#30, #50) |
| Repo hygiene | public repo; no CI workflow, no branch protection, no LICENSE, SECURITY.md, or `.gitignore` |
| Packets | 02 conflicts with the merged Go validator (#54); 05–09 assume a single freeze that practice has already moved past (#55); nothing owns `rehearse`, `idp-request`, or the operator-side k8s facts |

## 2. Rulings the plan is waiting on

Ranked by what they block. Each item: what it is, what happens, what breaks if it stays open, where the same problem is already solved, and the recommendation the plan is written under. Items 1–10 block the authorization-server slice; 11–22 are process and packaging calls that can land any time before their packet.

1. **#62. Real clients cannot send `purpose` and `requested_duration`.** `[decide]` *What happens:* the 2026-09-01 probe recorded three Codex CLI 0.151.0 authorize requests; none carries either parameter, and neither Codex nor Claude Code has a way to add arbitrary authorize parameters, the MCP client flow is plain OAuth. Under G4 as ruled (2026-08-31, 35a) the absent parameters are `invalid_request`, so never 8 refuses the product's principal clients on day one. *Already solved:* the consent page is already a gateway-owned, signed stage the user sees (mcp-sso §9.3, approve POST with a signed consent token). *Recommend:* the consent page is the carrier, two fields, purpose and duration (defaulting to the route's maximum), POSTed with the consent; the policy step (G-c) runs at approve time on the submitted values. Executable state machine: the page after §9.3 steps 1–4 carries consent token C1 (stage `entry`); its POST consumes C1's JTI, validates, runs policy, deny terminates, allow issues, escalate records; when the submitted values hash to an approved pre-approval for the tuple, the claim (A6) happens in that transaction and the response is a locked confirmation page under a fresh token C2 (stage `confirm`) whose POST alone runs A7/A8. Two tokens, one discriminator. The authorize-parameter carrier is dropped and nothing prefills from authorize parameters in v0. This reverses 35a with evidence; packet 16 is to confirm both clients complete a page with two extra fields (not yet observed). Consequences: G4/G6 (A1–A3 move their policy predicate to approve time, D5 rewritten), D3 (the consent token no longer carries purpose and duration, the approve POST does, bound by the JTI), B7 rows (malformed purpose or duration is refused on the approve channel), threat model (purpose never travels in a URL, so it never reaches browser history, ingress logs, or referrers ; that finding is retired), and hostile-purpose fixtures (HTML, Unicode, control characters, size) with the inherited page controls (HTML escaping, CSP, `nosniff`, `Cache-Control: no-store`, referrer policy) cited by clause.
2. **#53, Codex requests the union of advertised scopes on every route.** `[decide]` *What happens:* Codex reads the origin AS metadata's `scopes_supported`, not the route's PRM, and sends the union at `/authorize`. Under the inherited §9.3 step 3 a scope outside the route catalog is `invalid_scope`. *If open:* every Codex login fails on any gateway whose routes have different catalogs, `/splunk-read` next to `/splunk-admin`, the shape the README sells. *Already solved:* the group ceiling (`groupsToScopes`) narrows silently and emits `scope_ceiling_applied`. *Recommend:* two stages, so two different refusals stay different: first the route catalog, requested scopes outside it are removed, `scope_ceiling_applied` fires, an empty result is `invalid_scope` (a new `deltas.md` row; the reference's catalog-refusal fixtures become *host*); then the inherited group ceiling exactly as mcp-sso §17.4 has it, an empty result is `access_denied`, an entitlement refusal, never a malformed-scope one. The narrowed `scope` is returned in the token response (RFC 6749 §5.1). Follow-up probe (packet 16): confirm Codex accepts a narrowed `scope`. Fallback if it does not: omit `scopes_supported` from the origin metadata (optional in RFC 8414) and re-probe.
3. **#5, CIMD documents: vendored only, or opt-in live fetch.** `[decide]` *What happens:* Codex mints one CIMD document per install per server entry. *If vendored-only:* onboarding needs one document per user per route collected before that user can sign in, which contradicts "no pre-provisioning". *Already solved:* mcp-sso's guarded CIMD fetcher (§17.1.5: resolve, validate the address, dial that address directly, no redirects, bounded work), the parity corpus pins it. *The catch:* that defense assumes a direct dial; through a corporate proxy the proxy resolves the name and the validated-IP dial is gone, so "inherit through the egress profile" contradicts the clause it inherits. *Recommend:* opt-in live fetch behind `clients.cimd.liveFetch: {egressProfile, allowedOrigins[]}`, exact `https` origins, never patterns; a document URL whose origin is not listed is refused before any network call. With an operator-chosen destination the DNS-pinning defense is not what stops SSRF (the allowlist is), so the fetch may use the named profile; when that profile is direct (`none`), the inherited validated-IP dial applies as written. The remaining inherited caps (document size, timeouts, redirects refused, content type, cache) are cited by clause, the contract specifies them; the shared corpus does not yet prove them (no CIMD fixture is frozen), so Atesaki's `inherited` fixtures carry them. In proxied mode the address-level defense is **lost** (the proxy resolves the name) and is an accepted residual that requires a trusted egress proxy; the safest v0 ships live fetch in **direct mode only** and defers proxied live fetch, in a proxied enterprise that means per-user vendoring for Codex until v0.1. A new `deltas.md` row either way. Vendored documents remain the default.
4. **#24. Who may run `atesaki grants`.** `[decide]` *If open:* packet 12, then the approvals in M4, cannot start. *Recommend:* the proposal as written, local CLI, the OS user's ability to open the store file is the authentication boundary, approvers as `{osUser, subject?}` entries, self-approval refused only where `subject` is present, named honestly as a **single-operator authority model**: the authority value is the numeric effective uid (B2's `0600` store is openable by one uid, so several listed OS users cannot satisfy B2 in v0); configured names resolve to uids at boot and both are recorded with an invocation correlation id; any supplied name is `claimed_approver` (evidence, never authority); the platform's exec audit trail is the human-identity source, a §14 recipe obligation. CLI authority is proven by suite receipts with an injected identity port, the fixture profile has no command carrier. No sentence may read as if Atesaki authenticated a person it did not.
5. **#60. What happens when the rate limiter itself fails.** `[decide]` *What happens:* the reference limiter fails **open** when it throws on authorize, approve, token, and revoke (fail-closed only for stored registration); Atesaki inherits that by silence, and B8 names budgets for register, authorize, and token only. *Already solved:* the decider-outage ruling. *Recommend:* limiter error = `temporarily_unavailable` on every OAuth path (a delta row); B8 gains approve = the authorize budget, revoke = the token budget.
6. **#63. The server has no pre-handler exhaustion envelope.** `[decide]` *What happens:* B5 bounds sizes and counts and B8 bounds authenticated streams, but nothing bounds time: no header-read timeout, no body-read deadline, no idle timeout, no connection cap, no pre-verification per-IP budget on `/mcp`. Slow headers, slow bodies, or anonymous connection churn exhaust the process before any route, identity, or subject limit runs. *Recommend (numbers for the "ok"):* header-read 10 s; body-read 60 s on non-stream requests; idle 120 s; connection cap 1 024 as a global semaphore; unauthenticated per-IP `/mcp` budget 120 per 60 s with bounded per-IP state (10 000 addresses, least-recently-seen evicted) plus a global anonymous budget so IP cardinality cannot exhaust the limiter; the TLS handshake bound lives in the pinned ingress recipe (Atesaki has no inbound TLS configuration and terminates TLS at the ingress). Blocks packet 05; real-socket slow-header and slow-body tests in M3.
7. **#64. Durable events are not guaranteed to reach the JSONL stream.** `[decide]` *What happens:* G12 says loss is possible only for flow events, but the JSONL fan-out of `grant_event` rows is best-effort with no cursor or retry, so a durable event can be missing from the advertised combined stream. *Already solved:* G2's `grant_event.seq`. *Recommend:* the store is the durable audit of record; the JSONL projector serializes sink writes, appends and `fsync`s, **then** advances a cursor kept in the store, a crash between the two duplicates, never loses; consumers deduplicate by `event_id` (the projector does not scan JSONL); flow events stay lossy and counted. One cursor, no dispatcher framework.
8. **#66. Hash bytes, A6b atomicity, A9′ error mapping.** `[decide]` *What happens:* `purpose_hex` has no byte, trim, or case rule; `approved_hash` names no exact member set; A6b does not say whether invalidation and the following A3/A2 are one transaction; A9′ lists two public errors without mapping each predicate. Two implementations would hash differently and refuse differently. *Recommend:* one byte-exact hash vector in G3; the approved object spelled out; A6b as one transaction; every A9′ predicate mapped to one error, with a delta for any inherited change.
9. **#67. What stays in v0.** `[decide]` The reviewer's judgment: with one person merging and December the target, defer **machine clients** (G10, A12, A13, D10a–D10c, tombstones) to v0.1, a separate capability with its own rows and fixtures whose removal touches nothing in the human loop, and keep **rung 4** (signed proxy assertions), because "no IdP change at all" is the positioning sentence and cutting it cuts the pitch. Not docs-only: the merged validator accepts `machineClients[]`, a valid example declares one, and record types, B7 reasons, G5 states, onboarding, and rehearsal carry machine shapes, deferral is one config-boundary PR (packet 02 phase 2a) that makes the field unknown for v0 and sweeps every sibling, **before** record schemas are generated. Rule it first. The reviewer also proposed deferring the periodic sweeper and retention purge (keeping lazy expiry) and proxied live CIMD. The plan below marks machine clients as M5-if-kept. Reverses 2026-08-31 (machine clients in v0) if taken.
10. **#61. Readiness and shutdown semantics.** `[decide]` *What happens:* B1 reserves `livePath`/`readyPath` and §14 asks the recipe to say how streams end at shutdown, but no sentence says what `readyz` checks or what `SIGTERM` does. *Recommend:* readiness per identity mode and per shipped capability, M3: store directory admissible under B2, audit sink open, signing key loaded (no database exists until M4); M4 adds database open with migration state current, and "identity JWKS fetched at boot" only for modes that fetch one (console and a static `jwksRef` never do), never upstream reachability, which is `validate --deep`'s job. `SIGTERM`: stop accepting, drain non-stream requests for a bounded time (B8), cancel every stream's context, force-close after the bound (Go's `Shutdown` alone waits forever on an open stream).
11. **#54. Packet 02 is superseded by the Go validator.** `[decide]` *Recommend:* drop the config JSON Schema; add a mechanical B1↔parser drift test in Go, parser→B1 strict at all times (no undocumented accepted field), B1→parser reported as pending gaps that may persist while the slice's fixtures are merged ahead of its code and must be empty at slice **completion**, closed by the slice's first, fixture-driven boundary PR; write the G2 records as Go types and generate `schema/records/*.schema.json` from them with a golden test. Reverses #36's config half only.
12. **#55. One freeze or a rolling one.** `[decide]` *Recommend:* per slice, and **no machinery** (#30, #50, #52, `quality-bar.md`: "no hash manifest or self-checking workflow decides this; review owns changes"). Before a slice's code starts, its fixtures are merged and the owner has read them, the PR approval is the record. A fixture carries only what §19 already puts in the file: `status` (`draft` until a runner passes it, then `frozen` with the `receipt` object naming implementation, version, commit, date) and the clause it pins. No lock file, no manifest, no catalogue, no hash gate: a fixture edited inside an implementation PR is a review-checkpoint finding (#52) that the diff shows. The runner runs every non-superseded fixture and reports by status; a skipped fixture is a failure (never 7). `contract-v0-freeze` is one git tag when the whole portable set is green.
13. **#56, B2 file rules on Kubernetes.** `[decide]` *Recommend:* `knownCimd[]` entries become B2 references; the recipe states, pinned to tested versions, that secrets, CA bundles (`caBundleRef` is already a reference), and CIMD documents arrive as `env:`; **and** B2 gains the one exception the code already has: the configuration file itself is read once, may be a symlink (a ConfigMap mount is one), and carries the size cap but no ownership or mode rule, it is the operator's input, not a secret.
14. **#57. Policy rules cannot name "any Codex install".** `[decide]` *Recommend:* add `clientOriginIn` (exact origins of CIMD client-id URLs) beside `clientIn`, AND-only; DCR clients never match it.
15. **#58. Bounded JWKS refetch on an unknown `kid`.** A number, so a ruling: at most one on-demand refetch per key set per 60 s. `[decide]`
16. **#59. What `validate --deep` sends to an upstream.** A `GET` with no credential, closed after the status line, reported as "transport path reachable", never "the backend works" (a proxy block page answers too). `[decide]`
17. **#65. The first upgrade.** `[decide]` No store schema version, migration, backup, or downgrade refusal exists, and B1 has one signing key. *Recommend:* a schema-version row and forward-only migrations inside one transaction, downgrade refused, SQLite online backup for the recipe, restore tested; key rotation in v0 is hard and honest: replace the key, restart, every credential dies. That needs a **credential epoch** (the key's fingerprint) on grants, refresh families, and codes, advanced by rotation and restore and refused on mismatch, because codes and refresh tokens are stored by hash, not signed, and would otherwise survive the key and a restored backup would resurrect revoked state. No key ring in v0.
18. **B8 configurability.** Accept fixed numbers, no `limits:` block. `[decide]`
19. **Client-matrix staleness window.** 90 days. `[decide]`
20. **Dependency floor.** 15 days like mcp-sso, majors 30. `[decide]`
21. **LICENSE.** mcp-sso is MIT; Apache-2.0 adds an explicit patent grant and notice requirements. Owner's choice; the plan carries either. `[decide]`
22. **Name check now** (open question #9). The repo is public and the module path is in `go.mod`; run the check in M0 as its own dispatch. `[decide]`

One wording fix rides with packet 14: B4 says `alg` is "never read from the token"; the executable rule is "the token's `alg` must equal the configured one and match the key's type, and the allowed set never comes from the token".

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
                          │              M4 slice 2: sign-in, store, the human loop (allow, escalate,
                          │                     │        approve, claim), grants CLI, idp-request
                          └──► 03 phase 3 ──────▼
                                         M5 slice 3: machine clients (if kept), sweeper, retention,
                                                     audit projection, upgrade path
                                                │
                                                ▼
                                         M6 rehearse (mock IdP, client profiles)
                                                │
                                                ▼
                                         M7 deployment kit + recipe ──► M8 publish
```

Serial inside the Atesaki lane. The two dashed inputs from mcp-sso are the only cross-lane dependencies and are named where they bite. The table below is the one authoritative dependency view; the milestone sections and `prompts/README.md` are derived from it.

| Packet | Produces | Consumes | Blocking rulings | Starts when | Done when |
| --- | --- | --- | --- | --- | --- |
| 16 (M2 step 0) | client compatibility spike: what Codex and Claude Code actually send and accept (authorize parameters, narrowed `scope`, `approval_pending`, per-install CIMD, consent POST) | probe server from PR 4 | none | now, two days, before any M2 ruling closes | evidence recorded; #62, #53, #5 rulings informed |
| 13 | CI, protection, license, SECURITY.md, cooldown, grammar fixes, STATE/ledger, name check | none | LICENSE; dependency floor; name check (#9) | now | CI required on `main`; red PR cannot merge |
| 02 | B1↔parser drift test; (phase 2a: `machineClients` removed if #67 defers); G2 record types + generated record schemas; `knownCimd` refs | merged config code | #54 (phases 1–2); #67 (phase 2a, before 2); #56 (phase 3) | after 13 | parser→B1 empty always; B1→parser pending list printed, empty at each slice completion; records golden; phase 3 after packet 14 item 3 and before packet 06 PR 3 |
| 14 | contract sentences for #62, #53, #5, #56, #57, PR-5 interpretations and header-name rule, #58, #59, #55, #60, #61, #63, #64, #65, #66, #67, B4 `alg` wording, B8 note, matrix window | rulings; packet 16 evidence | each item's ruling | any time a ruling lands | lint green; ledger receipts; fixture ids named for 03 |
| 12 | G13 authority text, B1 approvers row, audit fields, residuals | #24 ruling | #24 | after the ruling | lint green; fixture intent named for 03 phase 3 |
| 03 phase 0 | fixture profile + mutation suite | record schemas (02), packet 14 sentences | none | after 02 phase 2 and 14 | profile mutation suite green |
| 03 phase 1 | slice-1 fixtures, `draft`, merged and read by the owner | profile | none | after phase 0 | merged; uncovered clauses listed in the PR |
| 04 | threat model complete; negative matrix | 03 phase 1 ids | none | after 03 phase 1 | no attacker row without a rule; uncovered list published |
| 05 | runner, egress, pipeline, verifier, relay, serve, `--deep`, health/shutdown | merged phase-1 fixtures; mcp-sso §8 fixtures `frozen` in their files | #55, #59, #61, #63 (values), JOSE library | after 04 and the §8 fixtures | phase-1 + §8 green, zero skips; review clean; real MCP named |
| 03 phase 2 | slice-2 fixtures | profile; packet 14 | none | alongside 05; locked before 06 | locked |
| 06 | the slice-2 configuration fields (first PR, closes the B1→parser gaps), store port + SQLite with the server lock, AS, the whole human loop (allow, escalate, approve, claim, consent, exchange, rotation, revocation), the grants CLI, identity ports, `idp-request` | merged phase-2 fixtures; packet 12; packet 02 phase 3; mcp-sso §07/09/10/11 or listed deferrals | #62, #53, #5, #24, #57, #58, #60, #64, #66 | after 05, packet 12, 02 phase 3, and the phase-2 fixtures | parity line by clause; B1→parser empty; real sign-in to approval to tool call; review clean |
| 03 phase 3 | slice-3 fixtures (machine clients if kept, sweeper, retention, projection, upgrade) | profile | #67 | alongside 06; merged before 07 | merged and read |
| 07 | machine clients (if kept), sweeper, retention, JSONL projection with cursor, schema migration, backup | merged phase-3 fixtures | #67, #65 | after 06 and the phase-3 fixtures | every G6 row green; `contract-v0-freeze` |
| 15 | `rehearse` + client profiles | runner; full AS | none | after 07 | onboarding step 4 true |
| 08 | recipe, image, kustomize, client matrix | `idp-request` (06), profiles (15) | matrix window | after 15 | recipe run once end to end |
| 09 | README, CHANGELOG, release workflow, sanitization in CI, listings | everything | none | after 08 | live verification named |

### M0: Repo hardening

**You can now.** Every PR runs the suite before you see it; `main` cannot be force-pushed or merged red; a stranger can report a vulnerability; the tree says what it is licensed under.

**Security first.** This is the milestone that makes the later ones checkable. Go has no package-manager-level release-age gate and the module proxy's `.info` time is the **commit** time, not a publish time (a freshly tagged old commit would pass any age test), so there is no honest second layer to build: the cooldown is Dependabot `cooldown` (15 days, majors 30, item 20) plus GitHub's dependency-review check on every PR that touches `go.mod`, plus the PR template asking for the version's publish date and age as human evidence. `go.sum` is the integrity ledger, not a lockfile, `go.mod` plus minimal version selection is what pins the build list, verified with `go list -m all`. `govulncheck` runs in CI from a pinned version. Builds use `-mod=readonly`. Race detector on. Linux and macOS in the matrix (the B2 rules use `syscall.Stat_t` and `O_NOFOLLOW`; Windows is out of v0 and the README says so). Actions pinned by commit; the fixture corpus is supply-chain input from here on (its hashes are verified before anything is materialized, M3). Two grammar defects found in review are fixed here, with their mirrors: `checkHostPort` accepts a present but empty port (`gw.example.com:`), the unmirrored sibling of PR 6's URL-port fix, and the redirect-allowlist grammar accepts `http://` on any host where the inherited §10 allows it only on loopback; a static credential header name may be a transport or hop-by-hop field (`Host`, `Content-Length`, `Transfer-Encoding`, `Connection`) and must be refused, the sentence lands in packet 14 first, the code after.

**Tests.** Prove the gate bites: the CI PR carries one deliberately failing test in a temporary commit, shows red, and, with protection on and admin bypass disabled, the owner's own account cannot merge it; then drop the commit, show green. The sanitization grep (packet 09's categories: hostnames, tenant ids, group names, vault paths, employer names, private evidence paths) runs here first and before every public push from now on; M8 re-runs it, it does not introduce it.

**Implement** (packet 13):

| PR | Content |
| --- | --- |
| `ci: build, vet, gofmt, race tests, govulncheck, contract lint, sanitization grep` | matrix job plus one aggregate `ci` job that branch protection requires, with `if: always()` and an explicit check that every leg succeeded (a skipped required check counts as passing on GitHub); `gofmt` output tested, not printed; protection as a ruleset: no force-push, no deletion, required check, linear history, **bypass disabled for admins**; proven by a deliberately failed leg the owner cannot merge |
| `chore: license, security policy, gitignore` | LICENSE (owner picks; Apache-2.0 matches the rest of the family unless told otherwise), `SECURITY.md` with the disclosure channel, `.gitignore` for the binary and local scratch |
| `chore: dependency cooldown and review` | Dependabot `cooldown` at the ruled floor, dependency-review on PRs, PR template evidence line; no commit-age script |
| `fix(config): read the config file through one descriptor` | open → `fstat` for the regular-file check → read through `io.LimitReader(cap+1)` and refuse past the cap; the cap is enforced by the reader, so no size-then-read race exists (the secret-file path already opens once; the config path stats the name then reads the name, the TOCTOU sibling) |
| `refactor(config): one reserved-path source` | `main.go` and `validate.go` each carry the list today |
| `fix(config): host:port with an empty port; http redirects off loopback` | the two grammar defects whose rules already exist (B3; inherited §10's loopback-only `http`), each with a refusal case and its mirrors (IPv4, bracketed IPv6, origin, exact URL); the full inherited §10 entry grammar lands in M4 with its fixtures; the credential header-name refusal lands after packet 14 item 5 writes the sentence |
| `docs: name check` | packet 09's deliverable 1 run now as a report and a ledger row (open question #9) |
| `docs: record slices-before-freeze; refresh STATE; reorder packets` | ledger row with the PR 5 receipt; lane-1 rows corrected; `prompts/README.md` in this page's order |

Local, not a PR: delete the three merged branches; remove the dead `probe-a`/`probe-b` entries from the Codex global config (they error on every Codex start).

**Gates.** Start: none. Done: CI required on `main`; a red PR cannot merge; lint and tests green on the hardened tree.

**Decisions.** LICENSE (item 21); the dependency floor (item 20); the name check (item 22).

### M1: Config boundary (done) + residuals

**You can now.** Write `atesaki.yaml`, run `atesaki validate`, and be refused by the exact resource, field, and rule for 71 malformed shapes; `atesaki routes` prints the ingress path list. (PR 5, PR 6.)

**Security first.** Done: strict YAML (no anchors, aliases, tags, duplicate keys), unknown-is-refusal, references only, the B2 file invariants on `file:` refs, B3 grammars, B4 shape, B5 config caps, the G7 boot contradiction check. Residual: the six "interpretations to confirm" in PR 5's description are enforced by code but written nowhere in the contract, each becomes a B1 sentence or is reversed.

**Tests.** The refusal suite stays the mutation suite. Add the drift test (#54): a Go test parses the B1 tables in `contract-boundaries.md` (field path, type cell, required/optional per variant) and compares them, both directions, with the field registry the parser exposes. The registry records every accessor call, present or absent, so an optional field no example uses is still registered; a variant branch never executed shows up as a B1→parser miss. The parser→B1 direction fails at all times (an undocumented accepted field is a defect); the B1→parser direction is a pending list, printed on every run, allowed while a slice's fixtures are merged ahead of its code, and required empty at each slice's completion, closed by that slice's first, fixture-driven configuration PR. This is the executable form of "B1 is not claimed complete until the artifact exists". The record drift test compares G2's field lists **and** which G6 row sets each state-dependent field, so ownership by state is checked, not just names.

**Implement** (packet 02, rescoped):

| PR | Content |
| --- | --- |
| `docs(contract): write the six PR-5 interpretations into B1` | contract-change PR (packet 14 item 5); one line each; owner confirms or reverses |
| `test(config): B1 to parser drift check` | the parser registers every accepted path with its requiredness at accessor-call time; the test reads B1; both diffs printed; empty or fail |
| `feat(records): G2 record types and generated schemas` | Go types for `grant_request`, `preapproval`, `grant`, `authorization_code` delta, `grant_event`, `machine_tombstone` with state-dependent presence encoded as typed unions; `schema/records/*.schema.json` generated by a golden test; RFC 3339 3-ms timestamps; `snake_case`; needed by 03 phase 0 |
| `feat(config): knownCimd entries are references` | after packet 14 item 3 (#56); `env:`/`file:`; the B2 file rules apply to `file:`; does not block M2 |

**Gates.** Done: drift diffs empty; record schemas committed and reproducible.

**Decisions.** #54, #56.

### M2: Contract closure before any authorization-server code

**You can now.** Nothing new runs. Every rule the AS and grants slices will implement exists as a sentence, and every sentence that a wrong build could violate has a fixture id or is listed as uncovered by name.

**Security first.** This is where the OWASP pass happens on paper: injection (scope names, purpose, subject bytes, JSON encoding of audit fields), broken auth (assertion verification, client authentication, pairing code), authorization (audience wall, ceiling, approve-then-swap via hash binding), sensitive exposure (secrets by reference, non-oracular errors), misconfiguration (unknown-is-refusal, empty allowlists), SSRF (upstream from config only; CIMD/JWKS through egress profiles with address-range refusal), path traversal (B3 single decode, separator scan), deserialization (strict YAML, bounded JSON), TOCTOU (CAS inside the transaction), open redirects (exact allowlist), tenant isolation (per-route audience), crypto misuse (alg pinned per key, `kid` exact), log injection (allowlisted fields, JSON-encoded lines), CSRF on consent (origin check, signed consent token), ReDoS (RE2 only, and only in fixtures), timing (digest comparison in constant time), proxy trust (B6), cache poisoning (no `cache-control` relay, metadata from config), caps on every input (B5), enforcement- plane outage (decider, store, JWKS, and the limiter, #60), and the **build and fixture supply chain**: a hostile fixture revision is a file-write primitive unless the runner contains it; modules, actions, tools, and base images are pinned; attestations bind the reviewed inputs. The threat model also records the residual that a group removed after activation does not revoke the grant, refresh rechecks nothing until expiry or revocation; the levers are a shorter `maxDuration` and `grants revoke`.

**Tests.** The negative matrix is the deliverable: attacker × surface × fixture id, rows without a fixture at the top. The fixture profile ships its own mutation suite (an mcp-sso-shaped config, an unknown clause id, an unknown record field, an unknown reason code, each refused by name).

**Implement.**

| Step | Packet | Content |
| --- | --- | --- |
| 0 | 16 | the two-day client spike against the PR-4 probe server: what Codex and Claude Code send at authorize (no custom parameters, confirm), whether they accept a narrowed `scope`, how they surface `approval_pending` and `request_id`, whether a consent page with two extra fields completes, per-install CIMD stability; conclusions into open questions #62, #53, #5 |
| 1 | 14 | contract closure, one ruling per PR: #62 the consent-page carrier, #53 two-stage ceiling (new deltas row), #5 live fetch with its allowlist, #56 `knownCimd` refs and the config-file exception, #57 `clientOriginIn`, the PR-5 interpretations and the credential header-name rule, #58, #59, #55 per-slice fixtures without machinery, #60, #61, #63, #64, #65, #66, #67, B4 `alg` wording, B8 note, matrix window; lint green; fixtures named as drafts |
| 2 | 12 | grants authority (#24 ruling): G13 authority text, B1 approvers row, audit fields per verb, the container residual; its fixtures are named by intent and written in 03 phase 3 |
| 3 | 03 phase 0 | the Atesaki fixture profile (`fixtures/schema/atesaki-fixture.schema.json`): mcp-sso spine + G/B/D/never clause ids + `given.config` as the Atesaki stream (validated by the real parser in the runner) + records from `schema/records/*` + B7 reasons with class + `given.files` for B2 boot fixtures under a containment contract (relative paths only, link targets inside the root, no ownership simulation, count and byte caps); its mutation suite. No manifest, catalogue, or hash gate (#30, #50), fixture ids carry their clause and the PR lists what is still uncovered |
| 4 | 03 phase 1 | slice-1 fixtures: relay §6, nevers 1, 3, 5, 7, B3 host and target grammar, B5 caps, B6 forwarded walk, B7 rows the relay reaches; every fixture `draft`; the B1 refusal suite and the B8 boundaries named as covered by their Go tests |
| 5 | 04 | threat model completed + `negative-matrix.md`; rows for later slices point at planned fixture ids and are counted as uncovered until those land |

**Gates.** Start: the ruling each packet-14 item names; packet 12 needs #24; 03 phase 0 needs 02 phase 2. Done: lint green; profile mutation suite green; every slice-1 fixture schema-valid; matrix published with its uncovered list.

### M3: Slice 1: runner, relay, verifier, `validate --deep`

**You can now.** Run `atesaki serve` in front of a real MCP server that today needs a shared key. The key lives in the container. A request without a token gets the route's own challenge; a token for another route is refused with that route's challenge; a good token reaches the tool. `validate --deep` proves the IdP metadata, each upstream, the store path, and the signing key are reachable through the configured egress before anything is deployed. `readyz` answers only when the store directory is admissible, the audit sink is open, and the key is loaded (no database exists yet); `SIGTERM` drains (#61). There is no login yet: the demo token is minted by the end-to-end test with the test signing key, there is no `mint` verb in v0 (§9's verb list is closed; a bounded mint would be a contract change).

**Security first.**

- *Request pipeline order is the contract:* byte caps (B5) → one request-target parse (B3: absolute-form authority equality, raw separator scan, single decode) → effective authority (B6) → Host/Origin gate → route match → verifier → relay. A step that runs later than this list is a finding. Some of it happens in Go's parser before any handler runs, and the plan says so instead of pretending otherwise: two HTTP/1.1 `Host` fields are refused `400` by `net/http` (no audit line is possible); `Transfer-Encoding: chunked` beside `Content-Length` is resolved the RFC 9112 §6.3 way, chunked wins and the length is dropped before dispatch; header **bytes** are capped by `MaxHeaderBytes`, which Go enforces per protocol, HTTP/1.1 with a 4 KiB read allowance including the request line; HTTP/2 as decoded HPACK fields with 32 bytes of overhead each and a 320-byte allowance, so B8's 64 KiB is enforced as those two observables, each measured on the shipped Go version and pinned by its own fixture, never asserted as a portable raw byte count, and the ingress is required to be no weaker; the header **count** (B8) is counted in the handler because Go has no count limit. Fixtures pin those observables; the recipe requires the ingress to apply the same framing rule.
- *Time bounds before identity* (#63): header-read timeout, body-read deadline, idle timeout, TLS handshake timeout, a listener connection cap, and an unauthenticated per-IP budget on `/mcp`, proven with real-socket slow-header and slow-body tests, not with `httptest`.
- *References are read once, through the validated descriptor:* `checkSecretFile` today validates a file and closes it; `serve` resolves every `env:`/`file:` reference into a typed boot snapshot from the same descriptor it validated, and nothing reopens a path at request time (B2's "read once at boot").
- *Relay is hand-built, not `httputil.ReverseProxy`:* the reverse proxy adds `X-Forwarded-For` and forwards headers by blocklist; §6 requires allowlists both ways. Upstream host and path come from config; the inbound request contributes a query string at most.
- *Duplicate headers fail closed:* Go's `Header.Get` returns the first value; the verifier must check `len(Header["Authorization"])` and the origin step `len(Header["Origin"])`; `Host` is promoted to `Request.Host` (HTTP/1.1 duplicates never reach the handler; on HTTP/2 a `Host` field beside `:authority` stays in `Header` and must byte-equal `Request.Host`).
- *Streams:* flush per event (`http.ResponseController`), no write deadline on a stream, non-stream upstream timeout 60 s (B8), client disconnect bound to the response context so a buffered POST is not aborted the moment its body closes.
- *Egress:* one transport per profile; `fromEnv` = `http.ProxyFromEnvironment`, `none` = nil proxy, URL = fixed proxy; `RootCAs` per profile, never the global pool; TLS ≥ 1.2; a proxy CONNECT failure is reported as `proxy <code> at <host>:<port>`.
- *Verifier:* one mature, pinned JOSE library after a source and supply-chain review, wrapped by a narrow verifier that admits exactly the contract's algorithms, key forms, and claims, the house rule is "never hand-roll parsing over untrusted input when a proven library exists", and compact JWS parsing (base64url framing, duplicate JSON members, `crit`, claim typing) is that surface. The token's `alg` must equal the configured algorithm and match the key's type, and the allowed set never comes from the token; `aud` exact single string; `iss` exact; `exp` strict with **no skew** (inherited §7.2, B8's 60 s is the rung-4 assertion skew, not an access-token grace); `scope` ⊇ `requireScope`.
- *Rate-limit identity:* the client IP from B6. A deployment whose ingress is not in `trustedProxies` puts every user in one bucket, `serve` warns at boot when it listens on a non-loopback address with `trustProxyHeaders: false`.
- *Audit lines* are JSON-encoded, one line, allowlisted fields; a formatter for untrusted input never throws.

**Tests.**

| Kind | What |
| --- | --- |
| Atesaki fixtures (03 phase 1) | relay rules, nevers 1/3/5/7, B3/B5/B6/B7 rows, zero skips |
| mcp-sso portable | the frozen §8 verifier set by fixture id, run by the same runner; every §7 token clause the verifier implements has a frozen upstream fixture **or** an Atesaki fixture marked `inherited` that quotes the same sentence, no inherited clause is implemented untested because the lane lags |
| Unit | grammars, header parsing, forwarded walk, target parsing, table-driven, one row per B3/B6 sentence |
| Exhaustion | real sockets: slow headers, slow bodies, idle connections, anonymous stream churn, shutdown under live streams, each ends at the #63 bound |
| Real input | `serve` in front of a named local MCP server (any stdio→HTTP server you can run) with a static-header credential: one tool call succeeds with a test-minted token; the same token at a second route is refused with that route's challenge; the upstream's recorded request headers contain no bearer token (never 1) |
| Negative matrix | every slice-1 row flips from planned to a fixture id |
| Randomness lint | a test fails if any package outside the randomness port imports `crypto/rand` (fixture determinism) |

**Implement** (packet 05):

| PR | Content |
| --- | --- |
| `feat(runner): load and validate Atesaki fixtures` | runs every non-superseded fixture and reports by `status`; no manifest or lock file to verify (#30, #50); `given.files` built under `os.Root` (relative paths only, no `..`, link targets inside the root, modes as stated, no ownership simulation, the owner-mismatch rule stays a unit test with an injected stat; count and byte caps); profile validation, chain ordering, clock/randomness/keys/recorded-HTTP ports, exact comparison, absence assertions, RE2 matchers; skipped locked fixture = failure |
| `feat(egress): profiles, proxy, CA per destination` | the one outbound layer; hop-naming errors; references resolved once into the boot snapshot |
| `feat(http): caps, authority, target parsing, host and origin gate` | the pipeline order above; B6 walk |
| `feat(verify): ES256 verifier, per-route metadata and challenge` | pinned JOSE library behind the narrow verifier; PRM at the path-inserted location, AS metadata documents at the origin, challenge per route (D1) |
| `feat(relay): allowlisted relay with streaming and cancel` | §6 in full |
| `feat(serve): wire the relay, flow audit, validate --deep, health, shutdown` | boot order: validate → open store path (create dir `0700`) → open audit sink → listen; `livez`/`readyz` and `SIGTERM` drain per #61; `--deep` per #59 |
| `test(e2e): real MCP behind the relay` | the named real input |

**Gates.** Start: M2 done; mcp-sso §8 portable fixtures `frozen` in their files with the `receipt` object (the cross-lane input); slice-1 sections SHA-pinned in the packet (#55). Done: all slice-1 fixtures green, zero skips; §8 portable green; every inherited clause this slice implements has a frozen upstream or an `inherited` Atesaki fixture, or is named as uncovered in the parity line and excluded from the capability claim; packet-11 review clean; real input named in the PR.

**Decisions.** #55, #59, #61, #63, and the JOSE library choice.

### M4: Slice 2: sign-in, the store, the human loop, `idp-request`

**You can now.** Point Claude Code or Codex at `https://host/route/mcp`. It discovers the route, registers (CIMD or DCR), the user signs in with the company login (Entra, generic OIDC, a signed proxy assertion, or the loopback console), states a purpose and a duration on the consent page, and, where a route rule says `allow`, the agent gets tokens and calls a tool. Where nothing allows it, the flow ends with `approval_pending` and a request id; `atesaki grants pending` shows it, `grants approve <id>` narrows and approves, the user runs the flow again, sees the approved values, approves, and the tool call works. `grants deny`, `grants revoke`, and RFC 7009 end access within one access TTL. `atesaki idp-request` prints the ticket for the IdP team. This is the product promise for people; machines come next.

*Why this shape:* under the default policy everything escalates, so a slice without approvals ends every default flow in a dead end and proves nothing about the loop onboarding sells. The store, the two-phase discipline (G8), the conformance suite, and the whole interactive operation table land here on rows with humans in them; M5 adds the rows with machines and clocks in them.

**Security first.**

- *Parity is the rigor, by clause:* every inherited clause this slice implements has a frozen upstream portable fixture or an `inherited` Atesaki fixture quoting the same sentence; the parity line names clauses, never whole sections (§17 also contains device flow and stored machine registration that Atesaki does not ship); a mode ships only when every portable fixture relevant to it passes, otherwise the mode is deferred, not a failure inside its claim. Every deviation is a `deltas.md` row or a bug.
- *The consent page is the carrier* (#62 as ruled): purpose and duration arrive by POST with the signed consent, never in a URL; two consent tokens, C1 (`entry`) whose POST validates, runs policy, and claims an approved pre-approval when the hash matches, and C2 (`confirm`) on the locked page whose POST alone runs A7/A8; hostile purpose text (HTML, Unicode, control characters, size) is refused or escaped per the inherited page controls (HTML escaping, CSP, `nosniff`, `no-store`, referrer policy); each JTI binds one submission.
- *Scope ceiling in two stages* (#53): catalog narrowing (`invalid_scope` when empty), then the inherited group ceiling (`access_denied` when empty); never 9 holds by construction; the narrowed `scope` goes back in the token response.
- *Approve-then-swap:* the claim (A6) is a CAS on `requested_hash` plus tuple; the grant is created with the **approved** values; freshness re-check (A6b) against current policy and ceiling, in one transaction (#66 as ruled).
- *Two concurrent claims:* exactly one wins; the loser proceeds as A3 in a new transaction (A6a), proven with a deterministic barrier, never by repetition.
- *CLI authority* (#24 as ruled): the store file's permissions are the boundary; ids are exact; the effective uid and a correlation id are recorded; a supplied name is `claimed_approver`, evidence only.
- *Identity ports:* inherited, not reinterpreted. The id_token is verified (issuer, audience, nonce, expiry, signature under the fetched JWKS); groups come from the claim only. Entra, exactly as mcp-sso §17 says: a groups overage is an identity refusal with its named audit reason, never a Graph call and never a silent empty ceiling; group identifiers are object ids (GUIDs) only. `idp-request` asks the IdP team for group filtering or app-assigned groups so overage cannot happen.
- *Rung 4:* B4 in full; `alg` pinned per key; `kid` exact; identity headers stripped everywhere but the identity leg; JWKS through the egress profile, size and count caps, stale-interval refusal; #58 bounds the refetch on an unknown `kid`.
- *Console pairing:* loopback only, before any state write; the pairing code is printed, never audited.
- *CIMD live fetch* (#5, if allowed): only for a document URL whose origin is on the operator's exact allowlist; the inherited caps cited by clause; the validated-IP dial when the profile is direct.
- *Store:* the driver is chosen and pinned **first**, and its open semantics proven by an integration test before any G6 row lands: how it opens the main file, how it creates `-wal`/`-shm` (SQLite's Unix VFS gives sidecars the database file's mode, observed for the driver actually used), whether it honors no-follow. B2 for the store: the directory is Atesaki's (`0700`, created by it, held open); the database file is created by Atesaki with `O_EXCL` and `0600` before the driver sees the path; after the driver opens, the path is re-opened with `O_NOFOLLOW` and its inode compared with the file Atesaki created. `database/sql` with **one connection** (`SetMaxOpenConns(1)`) so per-connection pragmas (`busy_timeout`, foreign keys, `synchronous=FULL`) and `BEGIN IMMEDIATE` apply to every transaction; WAL; a schema-version row and forward-only migrations in one transaction, downgrade refused (#65); a **server-instance lock** held only by `serve` (a second `serve` refuses to start; the grants CLI and the backup command open the store transactionally and are not locked out); local filesystem only, WAL is not for network filesystems. A B2 sentence the driver cannot meet is a contract gap reported before the adapter merges, never a weakened check.
- *Two-phase discipline (G8):* ids and tokens are produced in preflight; the transaction holds only authoritative reads, CAS predicates, mutations, durable events; the response is written after commit; E1 on a lost response.
- *Durable events reach the JSONL stream* (#64): serialized sink writes; append and `fsync`, then advance the cursor kept in the store; a crash between the two duplicates, never loses; consumers deduplicate by `event_id`; flow events stay lossy and counted.
- *Canonical JSON for hashes (G3):* RFC 8785 via a pinned JCS implementation or a purpose-written canonical serializer for the fixed G3 shapes, tested against the RFC vectors and against the byte-exact vector #66 adds; Go's `encoding/json` is not JCS.
- *Lazy expiry is this slice's:* every cap and dedupe read in A3, and every A6, A9, A10 read, first transitions past-due rows it touches (G5) inside the same transaction, or expired pending rows keep consuming the cap. M5 adds the sweeper.
- *Limiter outage* (#60 as ruled): never fail open on a token-issuing path.
- *Timing-safe comparison* for client secrets and the pairing code.
- *The refresh race clients will hit:* a retried refresh after a lost response is a replay and ends family and grant (A10′). Stated in the recipe, and crash-tested here on both sides of the commit.

**Tests.**

| Kind | What |
| --- | --- |
| mcp-sso portable | the pinned corpus commit; the exact frozen portable fixture-id set this build passes, listed in the PR before code; zero skips; deferred ids listed with reasons and the modes they defer |
| Atesaki fixtures (03 phase 2) | rungs §4 (each rung's boot refusals and acceptance; rung 4 duplicate header, unsigned header, wrong `kid`, stale JWKS, header stripping, bounded refetch); never 6; the consent-page carrier (#62) with hostile-purpose cases; the two-stage ceiling (#53); `clientOriginIn` (#57); live-fetch allowlist (#5); the limiter-outage delta (#60); Entra overage refusal; D1, D3, D4, D5, D6, D7, D11, D12, D13; A1–A11 with every branch (A3′, A3″, A6a, A6b, A9′, A10′, A10″), A14's lazy path, E1–E3; the projector cursor (#64); never 8 and never 9 as the matrices in §12; an `inherited` fixture for every §7/§9/§10/§11/§17 clause implemented here without a frozen upstream fixture |
| Store conformance suite | one table per G6 row this slice implements: atomicity, CAS, uniqueness, ordering-free comparison; both adapters; under contention; `serve` refuses memory |
| Crash tests | the mutable-state × exit-path matrix for every response-returning row this slice ships (A6, A8, A9, A10), with named failpoints between preflight and commit and between commit and response; restart; exact state (nothing consumed / committed-with-E1); for A10, the client's retry after the lost response is constructed and its A10′ outcome asserted |
| Race | the two-runner claim with a barrier at the CAS |
| Real input | Claude Code and Codex CLI (versions named) through Entra or generic OIDC: sign-in, consent page with purpose and duration, tool call on an `allow` route; on a default route: `approval_pending` → `grants approve` (narrowed) → re-run → consent shows the approved values → tool call → `grants revoke` → refresh refused, access dies at TTL; Codex on two routes with different catalogs (the #53 case); the console rung on a laptop; how each client surfaces `approval_pending` and where the user finds `request_id` |
| `idp-request` | golden output per provider; contains no secret; lists what it does not need |

**Implement** (packet 06, rescoped; one PR per invariant or chain, not per test):

| PR | Content |
| --- | --- |
| `feat(config): the slice-2 configuration fields` | first PR of the slice, fixture-driven: `clients.cimd.liveFetch`, `clientOriginIn`, the approver objects, `knownCimd` references consumed, the credential header-name refusal, the full inherited §10 redirect-entry grammar, closes the B1→parser pending list |
| `feat(store): port, memory adapter, conformance suite` | the interface, the suite as the contract, memory passes |
| `feat(store): sqlite adapter` | pinned driver with its open/sidecar semantics proven by test first; one connection; pragmas; schema version and migrations; the server-instance lock (two `serve` refused; `serve` plus CLI and `serve` plus backup allowed, tested); the B2 enforcement above; passes the suite under contention |
| `feat(as): metadata, stateless DCR, CIMD` | origin AS metadata; DCR per §9.2; vendored CIMD; live fetch behind its key and allowlist |
| `feat(as): authorize steps 1–4 and the two-stage ceiling` | §9.3 1–4 exactly, `resource` required (D13), #53 |
| `feat(policy): built-in rules` | G7 vocabulary incl. `clientOriginIn`; escalate by default; `policy_version` |
| `feat(grants): consent page carrier, A1–A3 with lazy expiry` | purpose and duration on the page (#62); policy at approve time; escalation ends with `approval_pending`; A14's lazy path in every read |
| `feat(grants): approvals A4, A5 and the claim A6, A6a, A6b` | with the authority contract from packet 12 |
| `feat(cli): grants list, pending, approve, deny, revoke` | OS-user authentication, exact-id lookup, `claimed_approver`, audit fields |
| `feat(grants): consent, exchange, rotation, revocation` | A7–A11; G8 sign-before-commit; G9 expiry propagation; D4 claims; A10 crash matrix |
| `feat(audit): durable event projection with a cursor` | #64; loss counter for flow events |
| `feat(identity): entra, oidc` | redirect flow, id_token verification, groups as GUIDs, overage refusal; subject boundary §6.5 |
| `feat(identity): header assertion, console pairing` | B4 in full; loopback-only pairing |
| `feat(cli): idp-request` | per-provider templates, the does-not-need list |
| `test(e2e): real sign-in, approval, and tool call` | the named clients and versions |

**Gates.** Start: M3 done; packet 12 landed; packet 02 phase 3 merged; 03 phase 2 fixtures merged and read; the mcp-sso §07/§09/§10/§11 portable set frozen **or** the exact not-yet-frozen ids listed as deferred with the modes they defer. Done: parity line by clause; the B1→parser pending list empty; every inherited clause has a frozen upstream or an `inherited` fixture, or its mode is deferred; the real input above shown; packet-11 review clean.

**Decisions.** #62, #53, #5, #24, #57, #58, #60, #64, #66.

### M5: Slice 3: machines, clocks, and the upgrade path

**You can now.** Unattended agents, if #67 keeps them in v0, are declared as machine clients and get bounded, revocable, tombstone-guarded tokens; expiry fires on time without waiting for a request; terminal rows purge; the store has a schema version, a migration path, a backup command, and a documented hard key rotation.

**Security first.**

- *Machine issuance* (if kept): requested ⊆ declared; deny-only rules; tombstone on the per-route digest; one active grant per (client, resource) with the losing insert discarding its signed token and retrying once as reuse (A12).
- *Sweeper and lazy expiry:* exactly one event per expiry, inside the operation's own transaction; retention purge idempotent (A15).
- *Upgrade* (#65): forward-only migrations in one transaction, downgrade refused, backup via SQLite's online backup, restore tested; key rotation = replace, restart, every credential dies, true only because grants, refresh families, and codes carry the credential epoch (the key fingerprint) that rotation and restore advance and a mismatch refuses; a pre-rotation code or refresh token is proven unable to mint under the new key.

**Tests.**

| Kind | What |
| --- | --- |
| Atesaki fixtures (03 phase 3) | A12, A13, A14 sweeper, A15; the machine first-issuance race; the migration and downgrade-refusal `boot` fixtures; every durable reason produced by at least one fixture; every G5 state reached |
| Store conformance | the remaining rows on both adapters |
| Crash | failpoints around A12; a crash mid-migration leaves the old schema intact |
| Real input | a machine client via `client_credentials` on a route with a deny rule (if kept); a restore from backup on a real cluster; an upgrade of a real schema-v1 database created by the pinned M4 commit's binary (archived with its checksum as the upgrade fixture, no earlier release tag exists) |

**Implement** (packet 07, rescoped):

| PR | Content |
| --- | --- |
| `feat(grants): machine clients A12, A13` (if kept) | D10a–D10c, tombstones |
| `feat(grants): sweeper and retention A14, A15` | 60 s interval over every row kind (the lazy path exists since M4), idempotent purge |
| `feat(store): migrations, backup, restore` | #65 |
| `test(e2e): machine client and upgrade` | the named real input |

**Gates.** Start: M4 done; 03 phase 3 fixtures merged and read; #67 and #65 ruled. Done: parity line green on the whole portable set; every G6 row has a green fixture; `contract-v0-freeze` tag applied (#55).

**Decisions.** #67, #65.

### M6: `rehearse`

**You can now.** Before deploying, run the whole protocol on your laptop against a mock IdP with recorded client profiles, discovery → registration → authorize → callback → token → one `/mcp` call, per configured rung. It is a **protocol and configuration self-test**: it proves the gateway and the config agree with what the named clients did when they were recorded. It does not run Codex or Claude Code and cannot prove a changed client version, a browser, or the company's IdP registration; live-client receipts with a tested version and date (M8) are the only compatibility proof, and onboarding's step 4 sentence is corrected to say so.

**Security first.** The mock IdP and the rehearsal listener bind loopback only; the memory adapter is accepted here and only here; `rehearse` never contacts the real IdP or a real upstream (recorded exchanges only, the runner's egress port); client profiles are recorded flows in the fixture profile, not live clients; output names the profile, the rung, and the step that failed, never a secret.

**Tests.** `rehearse` is the runner in a trench coat: each client profile is a chain in the Atesaki fixture profile executed against the composed binary, and the matrix is explicit, an `allow`-policy chain reaches `/mcp`; an escalation chain ends green at the expected `approval_pending`, or invokes the grants CLI and re-runs to `/mcp`; a signed-header profile uses a recorded JWKS exchange or a static `jwksRef`; the console profile pairs on loopback. Golden output per profile and rung; a deliberately broken config fails at the named step; the only network is loopback.

**Implement** (packet 15): `feat(cli): rehearse with a mock IdP`, then one PR per client profile (Claude Code CIMD, Codex CLI CIMD-per-install, DCR loopback, a hosted client with a fixed callback).

**Gates.** Start: M5 done. Done: the onboarding page's step 4 is literally true.

### M7: Deployment kit and recipe

**You can now.** Deploy one container with one kustomize example; the ingress path list comes from `atesaki routes`; the recipe tells you, per identity mode, exactly what to ask the IdP team, which secret keys the binary reads, where state lives, what is lost on restart, and what the platform must enforce that the product cannot.

**Security first.** Distroless static image, non-root uid, read-only root filesystem, a persistent volume mounted over a path the image does not own, `ReadWriteOncePod` on a named CSI storage class where available, else `ReadWriteOnce`, with `replicas: 1`, `strategy: Recreate` (a rolling update overlaps two pods on one store), and the server-instance lock from M4 so a second `serve` refuses to start while the CLI and backup still work; `fsGroup` so the non-root process can create the store subdirectory (`0700`) on first boot, proven on first boot **and** restart on a real cluster; one ingress controller and version pinned (a generic `Ingress` guarantees neither framing nor timeouts nor path semantics); secrets and CIMD documents as `env:` from Secret keys (#56); CA bundles the same way; image pinned by digest; SBOM and provenance attestations on the release; NetworkPolicy egress derived from the documented ports; the ingress in `trustedProxies` (or every user shares one rate-limit bucket) and applying the same request-framing rule as Go's parser; **no path rewrite at the ingress** (audiences are byte-exact); `livez`/`readyz` as ruled in #61; audit rotation that preserves every line already written (reopen, never truncate) while flow-event loss stays the accepted, counted residual it is (G12); backend reachability stated as the obligation it is (§14).

**Tests.** Every command in the recipe run once against the real binary (named); kustomize builds; the container runs under `readOnlyRootFilesystem: true`; a fresh reader reaches a real sign-in from the recipe alone for console mode and one real IdP mode; every other mode carries a tested-on date or the literal "UNVERIFIED" banner.

**Implement** (packet 08): `docs: deployment recipe`, `build: container image`, `deploy: kustomize example`, `docs: client matrix`. The recipe consumes M4's `idp-request` output and M6's client profiles; it designs no client flow of its own.

**Gates.** Start: M6 done; the staleness window ruled. Done: as the packet says.

**Decisions.** Client-matrix window.

### M8: Publish

**You can now.** A stranger reads the README, follows the ten-minute path, and gets a real sign-in; the trust artifacts (contract set, deltas, decisions ledger, fixtures, threat model, negative matrix, parity line) are one click away; the name is checked and recorded; the tree contains nothing employer-internal.

**Security first.** Sanitization is provable: a grep for hostnames, tenant ids, group names, vault paths, employer names, and the private evidence folder runs in CI from here on. Release artifacts carry checksums and attestations. `SECURITY.md` names the channel and the response window. Every "never/always/cannot" in public copy traces to a fixture or a receipt.

**Tests.** Live verification before promotion: two real clients (versions named) through the published image following the recipe verbatim from a clean machine.

**Implement** (packet 09): as written, with the name check already done in M0.

## 4. Cross-cutting lanes

- **Packet 11: adversarial implementation review** on every implementation PR, fresh context, "report everything with confidence", fail-closed sweep, sibling sweep, test-diff suspicion, invented-API grep, concurrency and partial failure for any G6 row touched. A PR is done when the review re-runs clean.
- **Packet 10: the mcp-sso lane.** The runner is done and in CI (`main` `63ed987`; 4/4 on the two draft 8.4 fixtures across Fastify, Express, and Hono, zero skips, verified locally 2026-09-02); PR 340 is closed. What Atesaki still needs, in order, and with **no freeze machinery** (the owner's standing rule, #30, #50): (1) the two 8.4 fixtures flipped to `status: frozen` with the `receipt` object in the file, one PR, since the runner already passes them; (2) the §08 verifier clauses 8.1–8.3 and the remaining 8.4 input classes, plus the §07 token clauses the verifier delegates to, each fixture frozen in the PR where the reference passes it (M3's input, today the corpus is two fixtures); (3) §09/§10/§11 fixtures labeled portable/host against `deltas.md` including the new scope-ceiling row, then §17 identity (M4's input); (4) a one-line §19 edit removing `MANIFEST.json`, `CATALOGUE.md`, and the hash gate from the text, since the status and receipt in the file plus the freeze log are the whole record. One session a day on that lane through September is the cheapest schedule insurance there is.
- **Packet 16: the client compatibility spike** runs before any authorization-server design is finalized and again before publish: what the real clients send and accept is evidence, and evidence has already overturned one ruling (#62).
- **STATE.md** is refreshed in the PR that changes a lane's state, never a separate chore.

## 5. Gotcha register

Each row names the milestone that owns it and the proof that it is handled. A row without a proof is a finding.

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
| 22 | Clients show `approval_pending` differently; the user must find `request_id` | M2/M4 | packet 16 probe; M4 real input |
| 23 | Rung-4 unknown `kid` refetch can be driven by an attacker (#58) | M4 | bounded refetch fixture |
| 24 | Fixture determinism dies if any package reads `crypto/rand` directly | M3 | import lint test |
| 25 | Go has no package-manager cooldown, and the module proxy's time is commit time | M0 | Dependabot cooldown + dependency review + human publish-age evidence; no age script |
| 26 | Windows lacks `O_NOFOLLOW`/`Stat_t` semantics | M0 | README: Linux and macOS only in v0; CI matrix |
| 27 | Pairing code in logs | M4 | audit fixture: absent |
| 28 | `validate --deep` "real read" against an MCP upstream is undefined (#59) | M3 | contract sentence + fixture before the verb ships |
| 29 | A group removed after activation does not revoke the grant; refresh rechecks nothing | M2/M7 | threat-model residual; recipe: shorter `maxDuration`, `grants revoke` |
| 30 | A hostile fixture revision is a file-write, chmod, or link primitive on every machine that runs the corpus | M3 | `os.Root` containment; hostile-path fixtures refused; fixture changes reviewed in the PR diff |
| 31 | The inherited limiter fails open when it throws (#60) | M4 | delta row + fixture: `temporarily_unavailable` |
| 32 | A readiness probe that includes upstream reachability takes a multi-route gateway out of the balancer when one backend flaps (#61) | M3/M7 | `readyz` semantics fixed by contract; recipe probes only those |
| 33 | Two `Host` fields are refused by Go before the handler; no audit line exists for them | M3 | fixture pins `400`; the negative matrix row cites the parser |
| 34 | `alg` confusion: pinning means "must equal the configured one", not "ignore the header" | M3 | packet 14 wording; §7/§8 fixtures |
| 35 | Real clients send no custom authorize parameters; `purpose` in the authorize URL is dead on arrival (#62) | M2/M4 | packet 16 evidence; consent-page carrier fixtures |
| 36 | Purpose in a URL lands in browser history, ingress logs, referrers | M4 | POST carrier; fixture: purpose absent from every URL and flow line |
| 37 | An empty group ceiling is `access_denied` (inherited), not `invalid_scope` | M4 | two-stage ceiling fixtures |
| 38 | A required check that is *skipped* counts as passing on GitHub | M0 | `if: always()` plus explicit leg results; a failed leg proven unmergeable |
| 39 | `ReadWriteOnce` allows several pods on one node; a rolling update runs two pods on one SQLite file | M7 | `Recreate`, `replicas: 1`, RWOP where available, startup lock |
| 40 | `database/sql` opens more connections than one; pragmas are per connection | M4 | `SetMaxOpenConns(1)`; conformance suite under contention |
| 41 | Go's `Shutdown` waits forever on an open SSE handler | M3 | stream contexts cancelled at drain; force-close after the bound (#61) |
| 42 | No header-read or idle timeout by default; slowloris exhausts the process before any limit runs (#63) | M3 | real-socket exhaustion tests |
| 43 | The module proxy's `.info` time is commit time, not publish time | M0 | no commit-age gate; Dependabot cooldown plus review |
| 44 | A static credential header named `Host`, `Content-Length`, or `Transfer-Encoding` corrupts framing | M0/M3 | refusal case; forbidden-sibling table |
| 45 | `MaxHeaderBytes` is enforced with a 4 KiB read allowance and includes the request line | M3 | threshold measured and pinned, not asserted |
| 46 | Access tokens have no skew; B8's 60 s is for rung-4 assertions | M3 | strict `exp` fixture |
| 47 | `.info`-style "any HTTP status" reachability is satisfied by a proxy block page (#59) | M3 | `--deep` reports transport path only |
| 48 | The first upgrade has no migration or key-rotation story (#65) | M5/M7 | migration `boot` fixtures; recipe rotation section |
| 49 | `checkHostPort` accepts `host:` (empty port); the URL sibling was fixed in PR 6 | M0 | refusal cases for both grammars |
| 50 | Redirect-allowlist entries accept `http://` on non-loopback hosts (inherited §10 forbids) | M0 | refusal case; inherited §10 fixture |
| 51 | Codes and refresh tokens are stored by hash, not signed: replacing the key alone does not kill them, and a restored backup resurrects revoked state (#65) | M5 | credential epoch; pre-rotation code and refresh token refused under the new key |
| 52 | A re-run after approval has no purpose in the request; without a second token the claim is ambiguous across up to three pending requests (#62) | M4 | C1/C2 state machine fixtures |
| 53 | Appending JSONL and advancing a store cursor cannot be atomic (#64) | M4 | append, `fsync`, then cursor; crash fixture shows a duplicate, never a loss |
| 54 | A process lock that covers every Atesaki process locks out the approval CLI | M4 | lock held by `serve` only; `serve` plus CLI tested |
| 55 | An unbounded per-IP limiter map is itself a DoS surface | M3 | bounded per-IP state plus a global anonymous budget (#63) |
| 56 | GitHub admins bypass branch protection by default | M0 | ruleset with bypass disabled; the owner proven unable to merge red |
| 57 | Every PR before M8 is already public; sanitization at M8 cannot scrub history | M0 | sanitization grep in CI from M0 |
| 58 | HTTP/2 header limits are HPACK-field accounting, not bytes on the wire | M3 | protocol-specific fixtures |

## 6. Not in v0

`docs/future.md` is the list. Nothing above adds to it; the two seams this plan touches (external decider, outbound OAuth client) stay ports with no config key.
