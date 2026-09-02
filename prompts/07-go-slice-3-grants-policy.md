MODEL: grok-4.5   EFFORT: high   FALLBACK: gpt-5.6-terra   TOOL: Codex CLI or pi, fresh clone
MILESTONE: M5 (docs/roadmap.md). PRECONDITION (#55): slices 1–2 merged; packet 03
phase 3 fixtures locked; the **#24 grants-authority contract (packet 12) is
owner-approved and landed** — this packet decides nothing about CLI authority.
PINNED: contract-grants.md @ <sha> · contract-boundaries.md @ <sha> · deltas.md @
<sha> · fixtures/MANIFEST.json @ <sha>.
WHY: the core promise — credentials are dispensed, not configured — completed:
approvals, machine clients, expiry, retention, durable audit fan-out. Implement the
remaining operation-table rows literally on the store and discipline slice 2 proved.

Read first, fully: docs/contract-grants.md (every line; G6 is the specification; this
slice's rows: A4, A5, A6, A6a, A6b, A12, A13, A14 sweeper, A15) · docs/contract-boundaries.md
B5, B7, B8 · schema/records/** · docs/deltas.md D2, D10a, D10b, D10c · fixtures/**
phase 3 · docs/negative-matrix.md rows tagged slice 3 · docs/roadmap.md §M5 (gotcha
rows 21, 22) · docs/decisions.md (do not relitigate) · the landed packet-12 text (G13,
B1 approvers row).

SCOPE — serial PRs, one behavior each, in this order:
1. `feat(grants): approvals A4, A5` — approver authenticated per packet 12; ∈ route
   approvers; ≠ request subject where checkable (the named residual stays named);
   approved ⊆ requested; duration ≤ requested; `approved_hash`; `claim_expires_at`;
   durable events.
2. `feat(grants): the claim A6, A6a, A6b` — CAS on `preapproval` `approved` ∧ now <
   `claim_expires_at` ∧ equal `requested_hash` ∧ same tuple; consent signed with the
   approved values and the new `request_id`; CAS lost → A3 in a new transaction;
   freshness (current policy ≠ `deny` ∧ approved ⊆ current ceiling) fails → invalidated
   then A3 or A2. The two-runner race is proven with a **deterministic barrier** —
   both runners held at the CAS and released together, so the collision is
   constructed, never hoped for.
3. `feat(cli): grants list, pending, approve, deny, revoke` — per packet 12: the OS
   user's ability to open the store file is the boundary; exact-id lookup; unknown id
   answers plainly; `--approver` label as evidence if ruled; every verb writes its
   durable event with the fields packet 12 names; output never contains a token or a
   secret; ids only on argv, never secrets.
4. `feat(grants): machine clients A12, A13` — `client_credentials` per the inherited
   token grammar and client authentication (timing-safe); declaration lookup per
   (client, route); requested ⊆ declared; deny-only rules; tombstone check on the
   current `declaration_digest` (G3); reuse of the active grant only on digest match
   and before expiry; the first-issuance race: a losing insert on the one-active-per-
   (client, resource) uniqueness rolls back, discards its signed token, retries once
   as reuse; D10b claims; A13 tombstone on revoke (state `active` or `expired`).
5. `feat(grants): sweeper and retention A14, A15` — 60 s interval (B8) over every
   row kind; the lazy path exists since slice 2 and is extended to every operation
   this slice adds; exactly one event per expiry;
   retention purge after 30 days (B8), idempotent, a `grant` only after its codes
   and family expired; `grant_event` rows never purged; flow `retention_purged` with
   count and oldest terminal timestamp.
6. `feat(audit): durable outbox fan-out` — `grant_event` rows committed in the
   operation's transaction (already, since slice 2); the fan-out to the JSONL sink
   best-effort with `audit_sink_failed` counted and loud; both classes in one stream;
   purpose in exactly `grant_issued` and `grant_machine_issued`, nowhere else; free
   text never in a flow event.
7. `test(e2e): escalation end to end` — the named real input (below).

HARD RULES: as prompts/README.md. A predicate that cannot be evaluated where G6 puts
it is a contract gap — report, do not move it silently. Never weaken CAS or
atomicity to pass a fixture. Nothing about authority beyond packet 12's words.

VERIFY: all phase-3 fixtures green, zero skips (A3″ both-caps refusal and A15 purge
by fixture id); the barrier race for A6 and the machine first-issuance race; never 8
and never 9 matrices green; the store conformance suite green on both adapters for
every row; crashes at named failpoints around A6 and A12 with restart and exact state
assertions; a real end-to-end escalation with a real MCP client (version named):
request → pending → `grants approve` (narrowed) → re-run → consent page shows the
approved values → tool call → `grants revoke` → refresh refused, access dies at TTL;
a machine client via `client_credentials` on a route with a deny rule, refused, then
allowed on another route, then tombstoned by `grants revoke` and refused after
restart; the client-matrix probe: how each real client surfaces `approval_pending`
and where the user finds `request_id`.

DONE WHEN: above verified and named; the parity line green on the whole frozen
portable set; every G6 row has a green fixture; packet-11 review clean; the owner
applies `contract-v0-freeze` (#55).

REPORT: fixtures by id; crash-test and race-test outputs; every contract gap; every
operation-table cell you could not implement as written; the client-matrix probe
results.
