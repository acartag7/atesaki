MODEL: claude-opus-class or fable   EFFORT: xhigh   TOOL: Claude Code in ~/project/atesaki-core
WHY: house rule — new surface = full threat model before v0.1; every enumerated edge
gets a negative test proving the rejection. docs/threat-model.md is a seed. Contract-
change PR.

Read first, fully: docs/threat-model.md · docs/contract.md · docs/contract-grants.md ·
docs/contract-boundaries.md · docs/deltas.md · fixtures/ (from 03) ·
~/project/mcp-sso/docs/threat-model.md (inherited surface).

DELIVERABLES
1. Complete `docs/threat-model.md`: assets, trust boundaries (client ↔ Atesaki ↔ IdP ↔
   upstream ↔ operator CLI ↔ filesystem ↔ store), attacker profiles with capabilities,
   and for every attacker × surface: attack, held-by rule (file:section), and the
   fixture id that proves the refusal. Cover at least: OWASP top-10 as applicable;
   path traversal; TOCTOU around claims/rotation/revocation; open redirects; tenant/
   route isolation; crypto misuse (alg pinning, kid, key rotation); log injection;
   CSRF on consent/approve; ReDoS on any regex fed untrusted input; timing-safe
   comparison of secrets/hashes; proxy/host-header trust; cache poisoning; caps on
   every untrusted input; enforcement-plane outage (decider, store, JWKS); approval
   binding and approve-then-swap; replay/idempotency for every state-changing
   endpoint; the audit trail as a target; confused-deputy via the relay.
2. `docs/negative-matrix.md`: the table attacker × surface × fixture id. A row with
   no fixture is a finding, listed at the top — never hidden.
3. Known-accepted residuals: each named, owner-tagged (`[O]` with a decisions.md row
   or it is not accepted), with the lever that reduces it.
4. Every "never/always/cannot" you add must trace to enforcing rule text and a fixture.

HARD RULES: read-only against contract pages other than threat-model.md and the new
matrix; a rule the threat model needs but the contract lacks is a **contract gap**,
recorded in the PR — do not add the rule yourself. No code.

DONE WHEN: no attacker row without a held-by rule; every held-by rule has a fixture or
is in the findings list; lint green.

REPORT: uncovered rows; contract gaps; residuals proposed for Arnold's acceptance.
