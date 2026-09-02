MODEL: claude-opus-class or fable   EFFORT: xhigh   TOOL: Claude Code in ~/project/atesaki-core
MILESTONE: M2 (docs/roadmap.md). Contract-change PRs only — one ruling per PR, lint
green each, nothing implemented. PRECONDITION per item: the owner's ruling on the
named open question, in his words. Implement his words, not the proposal.
WHY: these are the sentences the authorization-server and grants slices will
implement. Each was found by evidence (a live probe, the merged validator, the
Kubernetes file model) after the pages were written. Closing them now is what keeps
slice 2 from discovering them.

Read first, fully: docs/roadmap.md §2 · docs/open-questions.md #5, #53, #54, #55,
#56, #57, #58, #59 and the B8 configurability note in docs/contract-boundaries.md ·
docs/contract.md §2, §3, §8, §9, §14 · docs/contract-grants.md G4, G7 ·
docs/contract-boundaries.md B1, B2, B4, B5, B7, B8 · docs/deltas.md ·
docs/threat-model.md · docs/decisions.md · docs/quality-bar.md · the PR 5 description
("Interpretations to confirm") · tools/contract-lint.py (what it checks).

ITEMS — one PR each, in this order, each done when lint is green, the ledger row has
its receipt, the open question is struck with the receipt, and the draft fixture ids
the rule needs are listed in the PR (packet 03 writes them). Rulings the owner gives
in one sitting may share one PR when they touch different sentences.

0. `docs(contract): the consent page carries purpose and duration` (#62, as ruled)
   - contract-grants.md G4 rewritten: the two fields live on the consent page and
     arrive in the approve POST with the signed consent, bound by its JTI; the
     policy step (G-c) runs at approve time on the submitted values; the
     authorize-parameter carrier and its singleton guard are removed (a future
     client that can send them gets prefilled fields, nothing more); a re-run after
     approval shows the approved values and offers approve or deny only.
   - G6: A1/A2/A3 preflight and public outcomes move to the approve operation
     (the authorize step ends at the consent page for every request that passes
     §9.3 steps 1–4 and the ceiling); A6's claim happens at approve time too; the
     rows are rewritten, not patched — the table stays the specification.
   - deltas.md D3 and D5 rewritten; B7: malformed purpose or duration is
     `invalid_request` on the approve channel; the escalation and deny outcomes
     become approve-channel redirects.
   - threat-model.md: purpose never travels in a URL (browser history, ingress logs,
     referrers) — a held-by row, not a residual; hostile purpose (HTML, Unicode,
     control characters, size) refused or escaped under the inherited page controls
     (HTML escaping, CSP, `nosniff`, `Cache-Control: no-store`, referrer policy),
     cited by mcp-sso clause.
   - onboarding.md user step 3 rewritten (the person states the purpose and
     duration on the page the gateway shows).
   - decisions.md: the reversal of 35a with the probe as receipt.
1. `docs(contract): the route catalog is part of the scope ceiling` (#53, as ruled)
   - contract.md §3: two stages — first the route catalog (a requested scope outside
     it is removed, `scope_ceiling_applied`, an empty result → `invalid_scope`), then
     the inherited group ceiling unchanged (mcp-sso §17.4: an empty result →
     `access_denied`). Two refusals stay two refusals.
   - contract-grants.md G4 G-b: `scope_ceiling_applied` fires when either stage
     removed ≥1 scope; `requested_scopes_raw` unchanged (the hash input).
   - The token response carries the narrowed `scope` (RFC 6749 §5.1); the consent page
     shows the narrowed set.
   - deltas.md: a new row (reference: catalog refusal at §9.3 step 3 → host; Atesaki:
     narrowing) with its fixtures column; the group-ceiling outcome is not a delta.
   - threat-model.md: the row "client requests the union" with the held-by rule.
   - Note for packet 10: which mcp-sso fixtures become host.
   - If the ruling is the fallback (omit `scopes_supported`), write that instead and
     record the re-probe as a packet-06 verification step.
2. `docs(contract): opt-in live CIMD fetch` (#5, as ruled)
   - B1: `clients.cimd.liveFetch: {egressProfile, allowedOrigins[]}` (absent =
     vendored only; origins exact `https`, B3 host grammar, never patterns); B5: a
     CIMD document cap row; B8: the number (needs the owner's "ok"); contract.md §8
     rewritten; the inherited guarded-fetch clauses cited by mcp-sso §17.1.5, not
     reworded, with the one stated difference as a `deltas.md` row: the reference
     dials a validated address directly and forbids a proxy; Atesaki refuses any
     origin not on the allowlist before a network call, then uses the named profile
     (validated-IP dial when direct; the allowlist is the SSRF control when proxied);
     threat-model row for a client-supplied document URL.
   - If ruled vendored-only: onboarding.md loses "no pre-provisioning" for CIMD
     clients and says what the operator collects, per client.
3. `docs(contract): knownCimd entries are references, and the config-file exception` (#56)
   - B1 row: `[ref]` (B2); B2: `env:` for a document is a JSON string. `caBundleRef`
     is already a reference — do not touch it. B2 gains the exception the merged
     code embodies: the configuration file is read once, may be a symlink (a
     ConfigMap mount is one), carries the B5 size cap, and has no ownership or mode
     rule — it is the operator's input, not a secret. Recipe obligations for Kubernetes go
     into contract.md §14 as one paragraph (secrets, CA bundles, CIMD documents as
     `env:`; `file:` only where the runtime user owns a `0600` regular file; store
     path is a subdirectory Atesaki creates), stated as today's default volume
     behavior pinned to tested versions, not as a law of the platform.
4. `docs(contract): clientOriginIn in the rule vocabulary` (#57)
   - G7 and B1 `grant.policy.rules[].when`: `clientOriginIn: [origin]` in the B3 host
     grammar with scheme `https`; matches a CIMD client id whose origin equals an
     entry; DCR clients never match; AND-only; `policy_version` covers it.
   - The G7 boot contradiction check learns the new condition (a machine client
     never matches `clientOriginIn`).
5. `docs(contract): the six PR-5 interpretations and the header-name rule` — each
   as one B1 sentence, or reversed if the owner says so: `identity.registration`
   refused for `header` and `console`; padded strings refused, not trimmed;
   `grant.approvers` required unless some rule matches every request;
   `identity.redirectUri` must be `https`; duplicate Route names refused;
   `upstream.credential` required with `{type: none}` explicit; **and**
   `upstream.credential.header` refuses transport, framing, and hop-by-hop names
   (`Host`, `Content-Length`, `Transfer-Encoding`, `Connection`, `Trailer`,
   `Upgrade`, `TE`, `Keep-Alive`, `Proxy-*`) — the code lands in packet 13 PR 6.
6. `docs(contract): bounded JWKS refetch on an unknown kid` (#58, as ruled; number
   to B8).
7. `docs(contract): what validate --deep sends to an upstream` (#59, as ruled) —
   the verb reports "transport path reachable", never "the backend works".
8. `docs: per-slice freeze` (#55) — quality-bar.md "Order of work" and the packet
   preconditions say: sections SHA-pinned in the slice packet, fixtures hash-locked,
   owner has read those pages; `contract-v0-freeze` applied when the whole portable
   set is green. Ledger row. README status line updated.
9. `docs(contract): B8 configurability` — resolve the flagged B8 note per the ruling.
10. `docs(contract): the client-matrix staleness window` — §14 gains the number (B8
    row).
11. `docs(contract): limiter outage and the approve/revoke budgets` (#60, as ruled)
    — a `deltas.md` row against mcp-sso §09's fail-open limiter; B7 channel row; B8
    budgets for approve and revoke; threat-model row.
12. `docs(contract): readiness and shutdown` (#61, as ruled) — contract.md §9
    (`serve`) gains `livez`/`readyz` semantics and the `SIGTERM` sequence; B8 the
    drain bound; §14's "how active streams end" points at it.
13. `docs(contract): the alg sentence in B4` — "never read from the token" becomes
    "the token's `alg` must equal the configured algorithm and match the key's type;
    the allowed set never comes from the token" (RFC 7515 requires processing the
    header; RFC 8725 forbids trusting it). Mechanical wording; no ruling needed, but
    the ledger notes it.
14. `docs(contract): the pre-handler exhaustion envelope` (#63, as ruled) — B8
    rows for header-read timeout, body-read deadline, idle timeout, TLS handshake
    timeout, listener connection cap, unauthenticated per-IP `/mcp` budget; B5 gains
    the time dimension; threat-model rows for slow headers, slow bodies, anonymous
    stream churn, shutdown under live streams.
15. `docs(contract): durable events reach the JSONL stream at least once` (#64, as
    ruled) — G12 and G2 (`grant_event` outbox): a persisted projector cursor over
    `seq`, resume on restart, duplicates deduplicated by `event_id`; flow events
    lossy and counted; the "loss is possible only for flow events" sentence made
    true by construction.
16. `docs(contract): schema version, migration, backup, hard key rotation` (#65, as
    ruled) — contract.md §13/§14 persistence and upgrade sentences; B2 for the
    backup file; a `deltas.md` row is not needed (Atesaki-only).
17. `docs(contract): hash vector, A6b atomicity, A9′ mapping` (#66, as ruled) — G3
    gains the byte-exact vector and the approved object; A6b is one transaction;
    A9′ predicates each map to one public error, with a `deltas.md` row if the
    inherited mapping differs.
18. `docs(contract): v0 scope` (#67, as ruled) — if machine clients defer: G10,
    A12, A13, D10a–D10c, and the B1 `machineClients[]` row move to `future.md` with
    a pointer, the boot contradiction check's machine branch is retired, and
    `validate` refuses `machineClients[]` as unknown until v0.1; §13's scope pin and
    the README say so. If kept: nothing changes.

HARD RULES: one ruling per PR; nothing beyond the owner's words; every "never"
added carries the test a wrong build fails, named as a fixture id for packet 03;
no fixture written here; no Go; lint green before every push; the ledger receipt is
the owner's sentence, quoted.

REPORT PER PR: the ruling as received; the sentences added, by page and section; the
draft fixture ids listed; any sibling sentence the change touched (every page that
restated the rule — there should be none, quality-bar "one rule, one place").
