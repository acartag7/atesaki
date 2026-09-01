# Decisions ledger — the receipts behind every `[O]` tag

Every `[O:date]` in the contract set points at a row here. "Receipt" says *how* Arnold
decided: **picker** = an explicit choice in the design session's decision prompt on
that date; **words** = Arnold's own sentence, quoted or paraphrased; **PR** = recorded
as an owner comment on the named PR. A tag without a row here is a defect.

| Date | Decision | Receipt |
| --- | --- | --- |
| 2026-08-16 | Invariants: no passthrough; secrets as references; routes independently addressable; upstream OAuth client is core (later); no `passthrough` auth mode | ROADMAP.md (planning folder), "Answers — DECIDED 2026-08-16" |
| 2026-08-16 | Config is Kubernetes-resource-shaped from day one (Gateway/Route; local file = ConfigMap = future CRD) | ROADMAP.md, "Config is Kubernetes-resource-shaped" |
| 2026-08-16 | Accepted residual: a transparent relay does not sanitize upstream response bodies (a malicious upstream can reflect its own credential) | ROADMAP.md invariant 2, proxy-experiment proof |
| 2026-08-30 | The hosted front door is core Atesaki (reverses local-first) | picker |
| 2026-08-30 | New repo, Go; mcp-sso's rigor via contract + corpus, not code | picker + words ("golang, but please get the rigor but not the complications of mcp-sso") |
| 2026-08-30 | Credentials are dispensed, not configured: who / why / how long | words ("core part of all this should be that we dispense credentials for agents … and ask why they need it for and how long") |
| 2026-08-30 | Ladder of constrained modes, not forbidden shapes | words ("if you have no way to configure people will just do it") |
| 2026-08-30 | Proxy-aware from the start | words ("Also proxy aware from the start") |
| 2026-08-30 | Contract-change visibility gate instead of freeze tiers | words ("just make sure implementation doesn't try to write on the contract, in the PRs, without explicitly me being aware") |
| 2026-08-30 | Portable/host fixture split + public deltas list; topology (atesaki-core product repo, corpus in mcp-sso, planning archived) | picker |
| 2026-08-30 | Bounded revocation | picker |
| 2026-08-30 | Bounded lineage: mandatory expiry, one grant = one family, no cross-grant accumulation | picker |
| 2026-08-30 | Durable embedded store; memory rehearse-only | picker |
| 2026-08-30 | Stateless DCR only in v0 | picker |
| 2026-08-30 | Rung 4 in v0, signed assertions only | picker |
| 2026-08-30 | A policy/risk step decides: auto-approve / escalate to HITL / deny | words ("we can have kind of a policy/risk assessment some ops auto approved other escalate and wait for hitl") |
| 2026-08-30 | Policy per grant at request time; escalation completes by approve-then-re-run | picker |
| 2026-08-31 | The six mcp-sso §19 Phase-0 gaps: ~~anchors~~ (superseded by the §19 simplification); boot+suite evidence kinds; logical store vocabulary; capture-and-validate tokens; Fastify host; typed-vs-generic throws | words ("agree") → PR #338 comment |
| 2026-08-31 | Rule language: escalate by default, allow-list auto-approve, AND-only vocabulary, 24 h claim window | words ("seems reasonable") |
| 2026-08-31 | Numbers (B8) and `[R]` fail-closed defaults | words ("ok") |
| 2026-08-31 | Store = pure-Go SQLite behind a store port + conformance suite | words ("ok, make it interface/port to later on use another db type") |
| 2026-08-31 | Console pairing in v0, loopback-only | picker |
| 2026-08-31 | Backend reachability = recipe obligation only | picker |
| 2026-08-31 | Machine clients, operator-declared | picker |
| 2026-08-31 | Two routes may share one upstream | picker |
| 2026-08-31 | Agent sends `purpose` / `requested_duration` in the authorize request | picker |
| 2026-08-31 | Decider outage → `temporarily_unavailable` | picker |
| 2026-08-31 | ~~§19 anchor keys: bijective base-26 `[a-z]+`~~ **SUPERSEDED the same day** — see the §19 simplification row | picker → PR #338 comment |
| 2026-08-31 | Machine revocation is sticky (tombstone bound to declaration digest) | picker |
| 2026-08-31 | Grant clock starts at activation (code exchange) | picker |
| 2026-08-31 | Approval hash = what-to-whom only (no transport context) | picker |
| 2026-08-31 | Consent denial consumes the consent token | picker |
| 2026-08-31 | Machine clients: `deny` rules apply + boot contradiction check | picker (after explanation; Arnold: "explain me this better because i may go to 2") |
| 2026-08-31 | Sign before commit (G8); replaces the reference's burn→sign→store order | picker |
| 2026-08-31 | Audit: two durability classes — durable grant events, best-effort flow events (G12) | picker |
| 2026-08-31 | Rung 4 = generic signed-JWT assertion profile (B4) | picker |
| 2026-08-31 | Bundle: deny overrides allow; console may use loopback http; sweeper 60 s / pairing 10 min / approval window 24 h / request-target 8 KiB; claim-time freshness re-check; externalBaseUrl origin-only; no assertion replay cache in v0 | picker ("Confirm all six") |
| 2026-08-31 | Remaining B8 numbers (pending caps 3 / 10 000, retention 30 d, hops 10, JWKS 64 KiB / 32 keys, stale 24 h, assertion 8 KiB, groups 128 × 256 B) and total-pending-cap → `temporarily_unavailable` | picker ("Confirm all") |
| 2026-08-31 | ~~mcp-sso §19: prose is the only rule-carrier; fenced-block CI ban~~ **SUPERSEDED the same day** — see the §19 simplification row | picker → PR #338 comment |
| 2026-08-31 | **§19 simplification:** coverage counted by existing clause number, not by sentence; anchors, marker-word selector, fenced-block ban, and the "fail closed" spelling question cancelled; runner next | words ("are you putting rules based on words? This is plain stupid! I said that I don't want this complications") + picker → PR #338 and #339 comments |
| 2026-08-31 | Portable §8.4 fixture via the existing RE2 matcher accepting both documented metadata locations; host 8.4 unchanged, **frozen only once the runner passes it** — a receipt cannot precede a runner; no new assertion model | picker → PR #338 comment |
| 2026-09-01 | Contract mismatches trigger a concrete proposal and owner discussion; the owner chooses the current focused PR or a smaller linked PR; no automated contract gate | words ("this should be embedded everywhere"; "do not stop"; update the PR or split it) + Atesaki PR #1 withdrawal |
