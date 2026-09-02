MODEL: claude-opus-class or fable   EFFORT: high   TOOL: Claude Code in ~/project/atesaki-core MILESTONE: M2 (docs/roadmap.md), before packet 06 (the grants CLI ships in slice 2). PRECONDITION: Arnold has APPROVED the #24 proposal in docs/open-questions.md (or amended it, implement his words, not the proposal). Pre-freeze; before packet 07. WHY: the last open design cell with real content. The strict serial plan blocks on it.

Read first: docs/open-questions.md #24 (the proposal and Arnold's ruling) · docs/contract-grants.md G7 (approvers), G13 (verbs) · docs/contract-boundaries.md B1 (Route.spec.grant.approvers), B7 · docs/decisions.md · docs/threat-model.md.

DELIVERABLE: one contract-change PR:
1. Write the decided authority contract into contract-grants.md G13 (and B1's approvers row): who may run each verb, the approver identity model (per the ruling, e.g. {osUser, subject?} entries and the self-approval rule with its named residual), list/pending output visibility (purpose shown or redacted, his call), unknown-id behavior, the audit fields each verb writes, and, if ruled, the `--approver <label>` field recorded as `claimed_approver`, evidence, never authority, with the **numeric effective uid** (the authority value: B2's `0600` store is openable by one uid, so several listed OS users cannot satisfy B2 in v0; configured `osUser` names resolve to uids at boot and both are recorded) and an invocation correlation id on every verb's durable event; G13 names the model honestly as **single-operator authority** (every `kubectl exec` is one OS user; Atesaki authenticates no person) and the platform's exec audit trail is the human-identity source, a §14 recipe obligation. No sentence may read as if Atesaki authenticated a person.
2. Ledger row in docs/decisions.md with the receipt; strike #24 in open-questions.
3. The fixtures the new rules need, named by intent in the PR (unauthorized OS user refused per verb; self-approval refused where checkable; the named residual documented, not tested away), proven in packet 03 phase 2 as **suite receipts** with an injected effective-identity port (the fixture profile has no command carrier), not here: the profile does not exist when this packet runs.
4. Mirror threat-model rows (approver spoofing, store-file permission boundary).

HARD RULES: contract-change PR only; nothing implemented; nothing beyond what Arnold ruled; lint green; every "never" added carries its wrong-build test.

DONE WHEN: PR open, lint green, #24 struck with the receipt, packet 06's precondition satisfiable.
