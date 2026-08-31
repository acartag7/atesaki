MODEL: grok-4.5   EFFORT: high   FALLBACK: gpt-5.6-terra   TOOL: Codex CLI or pi, fresh clone
PRECONDITION: slices 1–2 merged; Atesaki grant fixtures (prompts/03 §4) locked; the
**#24 grants-authority contract packet is owner-approved and landed** — this packet
decides nothing about CLI authority.
WHY: the core promise — credentials are dispensed, not configured. Implement the
operation table literally.

Read first, fully: docs/contract-grants.md (every line; G6 is the specification) ·
docs/contract-boundaries.md B5, B7, B8 · schema/records/** · docs/deltas.md D2–D13 ·
fixtures/** for G6 rows and nevers 8–9 · docs/decisions.md (why each rule is the way
it is — do not relitigate).

SCOPE — build exactly this:
- **Store port** (Go interface) + **SQLite adapter** (pure-Go driver; name it, version,
  publish date) + the memory adapter from slice 2; a **store conformance suite** run
  against both adapters (§12-style tables: atomicity of each G6 row, CAS semantics,
  uniqueness, ordering-free comparison); `serve` refuses the memory adapter.
- The **operation table** A1–A15 including A3′, A3″, A6a, A6b, A9′, A10′, A10″,
  E1–E3: two-phase
  discipline (preflight outside the store; one transaction with only authoritative
  reads, CAS predicates, mutations, durable events); sign-before-commit; response after
  commit; lazy expiry + sweeper. One store transaction per A-row; E-rows have none.
- Hashes and digests exactly per G3 (RFC 8785 canonical JSON, domain-separation
  prefixes, sorted scope sets, purpose as hex).
- **Policy** (G7): built-in per-route rules, AND-only vocabulary, deny overrides allow,
  escalate by default, purpose never readable by rules, `policy_version` digest, boot
  contradiction check for machine declarations, outage → `temporarily_unavailable`.
- **Machine clients** (G10): `client_credentials`, declared clients, per-route digest,
  deny-only rules, sticky tombstone, one active grant per (client, route), token
  identity profile (D10b).
- **Audit** (G12): durable `grant_event` rows committed with each transaction; JSONL
  fan-out from the outbox best-effort with a loss counter; flow events direct.
- **CLI** `grants list|pending|approve|deny|revoke` per the landed #24 contract —
  nothing about authority is decided here.
- Onboarding flow end-to-end: agent sends `purpose`/`requested_duration`; consent page
  shows exactly the committed values; escalation → `approval_pending` + `request_id`;
  approver narrows; re-run claims once; bounded revocation observable.

HARD RULES: as prompts/README.md. A predicate that cannot be evaluated where G6 puts
it (preflight vs in-tx) is a contract gap — report, do not move it silently. Never
weaken CAS or atomicity to pass a fixture.

VERIFY: all grant fixtures green, zero skips (A3″ both-caps refusal and A15 purge
verified by fixture id); the two-runner claim race driven by a **deterministic
barrier** — both runners held at the CAS and released together, so the collision is
constructed, not hoped for (repetition alone proves nothing); never-8 and never-9 matrices green; the store
conformance suite green on both adapters; crashes injected at **named failpoints** — between preflight and commit, and between
commit and response — for A8 and A9, each followed by restart and exact state
assertions (nothing consumed/nothing issued, or committed-with-E1); a real end-to-end escalation with a real MCP
client: request → pending → `grants approve` → re-run → tool call → `grants revoke`
→ refresh refused, access dies at TTL.

DONE WHEN: above verified and named; gate green; PR opened.

REPORT: fixtures by id; crash-test and race-test outputs; every contract gap; every
operation-table cell you could not implement as written.
