# Open questions

**Nothing on this page is decided. Nothing elsewhere may cite this page as authority.**

1. ~~Machine clients~~ **DECIDED 2026-08-31 (Arnold): operator-declared machine
   clients in v0** — `client_credentials`, declaration = approval, each issuance
   attaches to/creates a grant, revocable (`contract-grants.md` G10, D10a–D10c). HITL for
   machines: future.
2. ~~Two routes → one upstream~~ **DECIDED 2026-08-31 (Arnold): allowed** — different
   audiences/catalogs/policies/credential refs onto one backend (`contract.md §5`).
3. ~~Store engine~~ **DECIDED 2026-08-31 (Arnold): pure-Go SQLite behind a store
   port** — a Go interface with a conformance suite every adapter must pass, so another
   database can be adopted later; memory adapter passes the same suite, `rehearse`-only
   (`contract.md §11, §13`).
4. ~~Rate limits~~ **DECIDED 2026-08-31 via B8** — per-IP per 60 s: register 30
   (inherited), authorize 60, token 120; keyed after B6 proxy unwrapping.
5. **CIMD live fetch: forbidden entirely, or opt-in behind an explicit key?** v0 serves
   vendored documents either way; the question is whether the option exists at all.
6. ~~In-cluster TLS/mTLS to upstreams~~ **DECIDED 2026-08-31 (Arnold): recipe
   obligation only in v0** — the product cannot enforce network position; §14 states
   it plainly; `validate --deep` warns on plain-http upstreams beyond loopback/private
   addresses. Signed gateway header / mTLS stay `future.md` candidates.
7. ~~Token TTLs~~ **DECIDED 2026-08-31** — `contract-boundaries.md` B8 (access 10 min
   = revocation lag; refresh = grant expiry; consent 5 min; code 2 min).
8. ~~Console-pairing~~ **DECIDED 2026-08-31 (Arnold): yes, loopback-only** — the
   one-machine tutorial identity, fail-closed on any non-loopback host/issuer/listen
   address, fixed `console-operator` subject, loud warning (`contract.md §13`, B1).
9. **Name and namespace check before first publish**: `atesaki` on GitHub, Go module
   path, container registry, and the search-query family — per the naming rule, the
   name locks at first publish. The local dir `atesaki-core` is a placeholder; rename
   freely.
10. **The corpus-expansion ask to mcp-sso**: which sections first, and who writes them —
    the fixture agent's next block of work should be sequenced against Atesaki's needs
    (§07, §09, §10, §11 before Go AS work can start proving parity).
11. **EMA / ID-JAG**: watch item. When Entra ships it, does Atesaki's AS accept the
    grant type? (Today: Okta-only in the wild; not a v0 decision.)

## Added 2026-08-30 after the freeze-readiness review (verdict: not freeze-ready)

12. ~~Freeze model~~ **DECIDED 2026-08-30 (Arnold):** no extra tiers. Acceptance tests
    may exist and run red before implementation. The one enforced rule: the
    implementation never changes contract or tests in its PRs without Arnold
    explicitly aware — CI hash-guard, dedicated contract-change PRs
    (`quality-bar.md`, "Who may change the contract"). §19's receipt semantics stay
    as-is for the shared mcp-sso corpus. **The CI mechanism was superseded on
    2026-09-01 by #52 before it landed; the owner-awareness rule remains.**
13. ~~Parity profiles~~ **DECIDED 2026-08-30 (Arnold): split + list.** Every shared
    fixture is labeled *portable* (every implementation passes) or *host* (TypeScript
    reference only). Atesaki keeps a short public list of intentional differences —
    `docs/deltas.md`, starting with per-route discovery/challenge. Corpus stays in
    mcp-sso; Atesaki pins an exact corpus version + manifest hash. Topology adopted:
    `atesaki-core` is the product repo (rename to `atesaki` at publish), planning
    folder archived, no third repo until a third implementation exists.
14. ~~Contracts that must exist before freeze~~ **DRAFTED 2026-08-30 as
    `docs/contract-boundaries.md`** (B1–B8). The `[R]` fail-closed defaults (symlink
    refusal, owner/mode `0600`, `0700` directories, ASCII-only hosts, pinned
    `ES256`/`RS256`, no `cache-control` relay) were **confirmed by Arnold 2026-08-31**;
    numbers are in B8.
34. ~~Caps and skew numbers~~ **DECIDED 2026-08-31 (Arnold: "ok")** — one table,
    `contract-boundaries.md` B8. All `[R]` fail-closed defaults in that page confirmed
    the same way. Any change from here is a contract change.

## Added 2026-08-31 after review round 3 — `[P]` items (proposed, awaiting Arnold)

> Receipts below that cite `contract.md §3.1–3.3` refer to material that now lives in
> `contract-grants.md` (G-numbers). History kept; current location stated here.

35. **`[P]` rulings needed** (each written into the contract *as a proposal*, tagged
    `[P]`, so no reviewer mistakes it for decided):
    a. ~~#26~~ **DECIDED 2026-08-31 (Arnold): the agent sends `purpose` and
       `requested_duration` in the authorize request** (G4, D5 carriers).
    b. ~~Decider outage~~ **DECIDED 2026-08-31 (Arnold): `temporarily_unavailable`**
       (G7, B7, mirrored in contract.md §11 and future.md).
    c. ~~expiry start~~ **DECIDED 2026-08-31 (Arnold): at activation** (code exchange),
       not at approval (G6 A9, G9).
    d. ~~claim-time freshness~~ **DECIDED 2026-08-31 (Arnold, bundle)** — A6.
    e. ~~origin-only base URL~~ **DECIDED 2026-08-31 (Arnold, bundle)** — B1/B3.
    f. ~~rung-4 profile~~ **DECIDED 2026-08-31 (Arnold): generic signed-JWT profile** — B4.
    g. ~~replay cache~~ **DECIDED 2026-08-31 (Arnold, bundle): none in v0**, named residual.
    h. → superseded by #42 (decided: deny rules apply + boot contradiction check).
    i. ~~B8 numbers~~ **DECIDED 2026-08-31 (Arnold, bundle)** — all B8 rows confirmed;
       total-pending-cap outcome `temporarily_unavailable` confirmed (G11).
36. **The machine-readable config schema** (JSON Schema, generated from B1, validated by
    a mutation suite) is the freeze artifact for configuration; B1 is the reference
    until it exists. Work, not a question. Same for the portable logical record schemas
    (G2 / D8).

## Added 2026-08-31 after review round 4

**Process finding (house rule: >3 review rounds = spec gap).** Rounds 2–4 kept
*discovering* in the grants machinery because it was written as prose. Fix applied:
`contract-grants.md` G6 is now an **operation table** (operation × preconditions ×
mutations × event × output × public error) so completeness is checked mechanically;
future additions to grants extend the table, never the prose.

37. ~~machine revocation~~ **DECIDED 2026-08-31 (Arnold): sticky tombstone** bound to
    the declaration digest; issuance refuses until the operator changes the declaration
    (G10, A13, D10c).
38. ~~hash scope~~ **DECIDED 2026-08-31 (Arnold): what-to-whom only** — `redirect_uri`
    and `code_challenge` excluded; final consent owns the recipient binding (G3).
39. ~~deny > allow~~ **DECIDED 2026-08-31 (Arnold, bundle)** — G7.
40. ~~denial consumes JTI~~ **DECIDED 2026-08-31 (Arnold): yes, spent** (A7, D11).
41. ~~console loopback http~~ **DECIDED 2026-08-31 (Arnold, bundle)** — B3.
42. ~~machines + rules~~ **DECIDED 2026-08-31 (Arnold): `deny` rules apply to machine
    issuance; `allow`/`escalate` ignored; `validate`/`serve` refuse a declaration wholly
    denied by its route's rules** (G7, G10, A12).
43. ~~numbers~~ **DECIDED 2026-08-31 (Arnold, bundle)**: sweeper 60 s; pairing code
    10 min; approval window 24 h; request-target 8 KiB (B8). Remaining B8 rows and the
    total-pending-cap outcome confirmed in the same bundle (#35i).
44. ~~mcp-sso §19 anchor key space~~ **SUPERSEDED 2026-08-31** by the §19
    simplification (#50): no per-sentence anchors at all.
15. ~~mcp-sso §19 untracked~~ **RESOLVED (observed 2026-08-30, late):** §19 and
    `fixtures/` are committed, indexed in `contracts.md`, in the CI contract guards,
    with portable/host profiles and the 8.4 fixture labeled host citing D1. Remaining
    blocker is narrower: no runner, no `MANIFEST.json`, no `CATALOGUE.md` — nothing is
    frozen yet.
16. **mcp-sso citation pinning**: at freeze, every `[S:mcp-sso §N]` pins an immutable
    SHA and exact clause. (Two mis-cites already found and fixed: §02→§01; blanket §05.)
17. ~~Corpus determinism gaps~~ **RESOLVED in the current §19 text** (RE2 dialect,
    absence assertions, exact state/outbound expectations, seeded PRNG) — verify
    against the runner once it exists.

## Added 2026-08-30 with the purpose-bound-grants core (§3.1)

18. ~~Headless agents~~ **DECIDED 2026-08-31 with #1** — machine clients declared in
    config (G9).
19. ~~HITL~~ **merged into #32/#33 (decided)** — the escalate path *is* the second
    human.
20. ~~Purpose bounds~~ **DECIDED 2026-08-31** — free text, 512 UTF-8 bytes (B8);
    no rule may read it (`§3.3`).
21. ~~Duration numbers~~ **DECIDED 2026-08-31** — default `maxDuration` 8 h, hard
    ceiling 30 days (B8). ("Until revoked" → #22, decided no.)

## Added 2026-08-30 after review round 2 (verdict: not freeze-ready; §3.1 needs GrantV0)

22. ~~"until revoked"~~ **DECIDED 2026-08-30 (Arnold): no.** Every grant expires; no
    config key creates an open-ended grant (`contract.md §3.1`).
23. ~~revocation guarantee~~ **DECIDED 2026-08-30 (Arnold): bounded revocation.**
    Refresh family dies now; access tokens live to their `exp` ≤ access TTL. Docs say
    "access ends within N of revocation", never "immediately" (`§3.1`, `§3.2`).
24. **`atesaki grants` authority contract.** Who may run it (operator auth), output
    redaction, exact-id lookup only (no enumeration), unknown-id behavior, concurrent
    revoke-vs-rotate semantics. Work, once #23 is decided.
25. ~~Grant cardinality~~ **DECIDED 2026-08-30 (Arnold): bounded lineage.** One grant
    → one refresh family; a new authorization is a new grant; no scope carry-over
    (D2 confirmed).
26. **REOPENED by round 3 — who supplies purpose and duration, and where?** The
    policy step needs both *before* consent is signed, and mcp-sso's approve POST
    consumes an already-signed consent token — so the values must arrive earlier.
    **DECIDED 2026-08-31 (Arnold, see 35a) — was the proposal (`contract-grants.md` G4):** the client (the agent) sends them as
    extension parameters of the authorize request (`purpose`, `requested_duration`);
    the consent page shows them and offers approve/deny only; narrowing happens
    through the approver on escalation, never by the user editing the page. This adds
    no protocol stage. Rejected alternative: a user-entry page before consent (a new
    signed stage). Arnold decides.
27. ~~DCR mode~~ **DECIDED 2026-08-30 (Arnold): stateless DCR only in v0.** No stored
    mode, no config key; grants are revoked, not clients; D2's accumulation path is
    unreachable by construction (`contract.md §8`).
28. ~~Store durability~~ **DECIDED 2026-08-30 (Arnold): durable embedded store.**
    `serve` refuses in-memory; memory is `rehearse`-only; store-file loss ⇒ key
    rotation ⇒ all tokens die (`contract.md §3.2`). Engine choice narrowed in #3.
29. ~~Header mode in v0~~ **DECIDED 2026-08-30 (Arnold): yes, signed assertions
    only.** Unsigned identity headers refused; no trusted-proxy-by-network mode
    (`contract.md §4`). Remaining #14 work for this rung shrinks to: assertion
    verification contract (signature/issuer/audience/expiry/clock skew), stripping of
    inbound identity headers on every other path, and the rehearse scenario.
30. ~~Visibility gate~~ **ABANDONED 2026-09-01:** the proposed hash manifest,
    mutation harness, and pull-request workflow added maintenance work without owning
    the engineering decision. PR #1 closed without merge. Replaced by the review
    checkpoint in #52.
31. ~~GrantV0 contract~~ **DRAFTED as `contract.md §3.2`** (2026-08-30) from the
    #22/#23/#25 rulings. Still consumes open inputs #26 (consent authorship), #27 (DCR
    mode), #28 (store durability), #7, #20, #21. The never-8/never-9 acceptance matrices
    and Atesaki's grant fixtures still do not exist — written after those inputs close.
    D3 (consent token carries grant fields) and D4 (access token carries `grant_id`)
    are consequences of §3.2 recorded in `deltas.md` for confirmation at freeze.
32. ~~Policy step shape~~ **DECIDED 2026-08-30 (Arnold):** (a) **per grant, at request
    time** — per-tool-call `ask` stays Edictum's, future; (b) **approve, then re-run** —
    pre-approval bound to the exact request hash, narrow-only, claim window. Written as
    `contract.md §3.3`; D5 recorded.
33. ~~Rule language~~ **DECIDED 2026-08-31 (Arnold: "seems reasonable"):** escalate by
    default, allow-list what may auto-approve; vocabulary = scopes ⊆ set, duration ≤
    ceiling, subject ∈ group, client ∈ list, AND-only; explicit `deny` rules; purpose
    never readable by a rule; claim window 24 h (`contract.md §3.3`, B8). The exact
    YAML shape of a rule is B1 work at freeze, not a new decision.
19. → merged into #32/#33 (the second-human approval *is* the escalate path).
24. **Grants-CLI authority — PROPOSAL ready for Arnold (pre-07 contract packet):**
    v0 is **local CLI only**: authentication = the OS user's ability to open the store
    file (Unix permissions are the authn boundary; no network admin API — that is a
    whole new attack surface future.md can earn). All five verbs share it. Lookups are
    by exact id; `grants list`/`pending` show everything including purpose — the
    operator is inside the trust boundary by construction. Unknown id answers plainly
    ("not found" — non-oracular shaping is for the public surface, not the operator's
    own store). Concurrency is already owned by the operation table (A4/A5/A11/A13
    transactions). Every CLI action writes its durable event with the **OS username**
    as `approver`/`revoked_by`.
    **The wrinkle needing your call:** G7 says approvers are a route-configured
    *group* and an approver may never approve their own request — but groups are IdP
    concepts and the CLI runs as an OS user; the namespaces don't match. Proposed v0:
    `Route.spec.grant.approvers[]` becomes a list of `{osUser, subject?}` entries —
    approval requires the invoking OS user to be listed; when the entry carries
    `subject`, a request whose subject equals it is refused (the not-your-own rule);
    an entry without `subject` cannot be self-checked, which is a named residual.
    IdP-group approvers return with a web approval surface (future.md).
    Alternative: keep approvers as IdP groups and have the CLI take an authenticated
    subject token — real authn, but it needs the whole token flow the CLI was meant to
    bypass. Arnold decides; then the packet is a one-page contract-change PR.

## Added 2026-08-31 after review round 5

45. ~~audit durability~~ **DECIDED 2026-08-31 (Arnold): two classes** — durable grant
    events committed with the change; flow events best-effort and counted (G12, B7).
46. ~~signing order~~ **DECIDED 2026-08-31 (Arnold): sign before commit** (G8; D6/D7/D12
    disclosed).
47. **`[D]` recorded, confirm at freeze — `resource` is required** on authorize and
    machine-token requests (no multi-route default; `invalid_target`), D13.
48. **Contract coverage tooling (work, not a question):** a script in this repo that
    checks G5 edges ↔ exactly one G6 row, every B7 reason has a producer row and a
    class, every G2 field is set by some row, every `[O]` tag has a `decisions.md`
    row, every G-/B-/D-number reference resolves. Reviews should verify, not discover;
    this is what makes that true for the doc set itself.

49. ~~mcp-sso §19 code-fence rules~~ **SUPERSEDED 2026-08-31** by #50: no fenced-block
    ban, no marker-word scanning — here or in mcp-sso.
50. **DECIDED 2026-08-31 (Arnold): §19 simplification.** Coverage by existing clause
    number; a fixture names its clause and quotes its sentence (already the format);
    which clauses need fixtures is a human-reviewed list. Anchors, `[a-z]+` grammar,
    marker-word selector, fenced-block ban, and the "fail closed" spelling question are
    cancelled. Recorded on PR #338 and #339. **Process finding:** three owner-decision
    requests about accounting machinery reached Arnold before anyone asked whether the
    machinery was needed — the design session relayed instead of triaging. Rule from
    here: any request from a build/fixture session that adds *accounting* (selectors,
    anchors, bans) rather than *behavior* is first checked against "rigor without the
    complications" and comes to Arnold as "do we need this at all?", not as a choice
    between variants.
51. ~~portable 8.4~~ **DECIDED 2026-08-31 (Arnold):** existing regex matcher with the
    metadata path optional; host 8.4 unchanged, frozen only when the runner passes it; no schema/§19 change. The
    runner starts. (Triage applied: behavior need real, machinery not needed.)
52. **DECIDED 2026-09-01: implementation mismatches are review checkpoints,
    not automated gates.** The implementer states what does not fit and proposes a
    concrete change. The owner approves, rejects, or refines it. Keep the change in the
    current PR when the PR remains focused; otherwise split a smaller linked PR. Keep
    working on unaffected items while the decision is pending. Never add exception
    machinery to avoid the discussion.

## Process finding #2 (2026-08-31, Arnold: "9 rounds of review, something is wrong?")

He's right. Causes, recorded so the next product doesn't repeat them: (1) the artifact
grew ~4× during its own review (the dispensing core, the boundaries page, and the
packets were all added mid-loop) — review was made to carry design; (2) repair sweeps
themselves introduced defects until the lint existed — the lint should be built at
round 1, not round 5; (3) "iterate until a round finds nothing" cannot terminate on
prose — an adversarial reader of 1,500 normative lines always finds another cell.
**Rule from here:** round 9 is the last prose round; only its P0/P1s are swept; all
further truth-finding happens in terminating, executable form (lint, schema mutation
suites, fixtures, packet-11 code reviews) plus Arnold's own read at freeze.
