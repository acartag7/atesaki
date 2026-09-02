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
the rule needs are listed in the PR (packet 03 writes them):

1. `docs(contract): the route catalog is part of the scope ceiling` (#53, as ruled)
   - contract.md §3: effective scopes = requested ∩ route catalog ∩ group ceiling; a
     requested scope outside is removed, never refused; empty result → `invalid_scope`.
   - contract-grants.md G4 G-b: `scope_ceiling_applied` fires when the catalog or the
     group ceiling removed ≥1 scope; `requested_scopes_raw` unchanged (the hash input).
   - The token response carries the narrowed `scope` (RFC 6749 §5.1); the consent page
     shows the narrowed set.
   - deltas.md: a new row (reference: catalog refusal at §9.3 step 3 → host; Atesaki:
     narrowing) with its fixtures column.
   - threat-model.md: the row "client requests the union" with the held-by rule.
   - Note for packet 10: which mcp-sso fixtures become host.
   - If the ruling is the fallback (omit `scopes_supported`), write that instead and
     record the re-probe as a packet-06 verification step.
2. `docs(contract): opt-in live CIMD fetch` (#5, as ruled)
   - B1: `clients.cimd.liveFetch: {egressProfile}` (absent = vendored only); B5: a
     CIMD document cap row; B8: the number (needs the owner's "ok"); contract.md §8
     rewritten; the inherited guarded-fetch clauses cited by mcp-sso section, not
     reworded; threat-model row (SSRF via a client-supplied document URL: `https`
     only, no redirects, private and loopback ranges refused after resolution,
     through the named profile).
   - If ruled vendored-only: onboarding.md loses "no pre-provisioning" for CIMD
     clients and says what the operator collects, per client.
3. `docs(contract): knownCimd entries are references` (#56)
   - B1 row: `[ref]` (B2); B2: `env:` for a document is a JSON string; recipe
     obligations for Kubernetes go into contract.md §14 as one paragraph (secrets, CA
     bundles, CIMD documents as `env:`; `file:` only where the runtime user owns a
     `0600` regular file; store path is a subdirectory Atesaki creates).
4. `docs(contract): clientOriginIn in the rule vocabulary` (#57)
   - G7 and B1 `grant.policy.rules[].when`: `clientOriginIn: [origin]` in the B3 host
     grammar with scheme `https`; matches a CIMD client id whose origin equals an
     entry; DCR clients never match; AND-only; `policy_version` covers it.
   - The G7 boot contradiction check learns the new condition (a machine client
     never matches `clientOriginIn`).
5. `docs(contract): the six PR-5 interpretations` — each as one B1 sentence, or
   reversed if the owner says so: `identity.registration` refused for `header` and
   `console`; padded strings refused, not trimmed; `grant.approvers` required unless
   some rule matches every request; `identity.redirectUri` must be `https`; duplicate
   Route names refused; `upstream.credential` required with `{type: none}` explicit.
6. `docs(contract): bounded JWKS refetch on an unknown kid` (#58, number to B8).
7. `docs(contract): what validate --deep sends to an upstream` (#59).
8. `docs: per-slice freeze` (#55) — quality-bar.md "Order of work" and the packet
   preconditions say: sections SHA-pinned in the slice packet, fixtures hash-locked,
   owner has read those pages; `contract-v0-freeze` applied when the whole portable
   set is green. Ledger row. README status line updated.
9. `docs(contract): B8 configurability and the client-matrix window` — resolve the
   flagged B8 note per the ruling; §14 gains the staleness window number (B8 row).

HARD RULES: one ruling per PR; nothing beyond the owner's words; every "never"
added carries the test a wrong build fails, named as a fixture id for packet 03;
no fixture written here; no Go; lint green before every push; the ledger receipt is
the owner's sentence, quoted.

REPORT PER PR: the ruling as received; the sentences added, by page and section; the
draft fixture ids listed; any sibling sentence the change touched (every page that
restated the rule — there should be none, quality-bar "one rule, one place").
