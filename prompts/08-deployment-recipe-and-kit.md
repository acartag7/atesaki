MODEL: gpt-5.6-sol   EFFORT: high   TOOL: Codex CLI in ~/project/atesaki-core (after slice 3)
WHY: contract.md §14 — the recipe is a shipped deliverable, per-mode, every obligation
stated; plus the container, the kustomize example, and the `idp-request` templates.

Read first, fully: docs/contract.md §4, §9, §13, §14 · docs/contract-boundaries.md
B1, B2, B6, B8 · docs/onboarding.md · docs/threat-model.md and docs/deltas.md (the field failure modes are already
encoded in the contract; the private evidence folder is **not** read by this packet —
that is what keeps packet 09's sanitization provable) · the shipped `cmd/atesaki` help output.

DELIVERABLES
1. `docs/recipe.md` — the §14 per-mode matrix, exactly: for each identity mode
   (dedicated Entra, shared Entra, generic OIDC, header signed-assertion, console):
   domain + TLS (who issues, URL chosen before IdP registration), IdP registration
   (client id, secret delivery, redirect URI registered BEFORE deploy, groups claim
   where used), secret provisioning (exact key names), network (trusted-proxy
   allowlist, backend reachability as the obligation the product cannot enforce,
   `NO_PROXY`, custom CA, NetworkPolicy limits), persistence (store + signing key,
   what is lost on restart), operations (audit permissions/rotation, rate limits,
   upgrade, stream shutdown), and the **client matrix** (registration model, callback
   behavior, version tested, date — staleness window: **proposed 90 days, awaiting the
   owner**. The B8 contract change precedes this deliverable only after owner approval.
   If unanswered, report this deliverable as blocked and continue with the others).
2. `idp-request` output templates per provider: the minimal ask, and the explicit
   "does NOT need" list (no Expose-an-API, no `api://` scopes, no per-client redirect
   churn), paste-ready for a ticket.
3. `Dockerfile` (distroless/static; non-root; read-only root; `0700` state dir) and
   `deploy/kustomize/` example: Deployment, Service, Ingress with the `atesaki routes`
   path list, NetworkPolicy derivable from documented egress, secrets as refs.
4. `rehearse` client profiles: Claude Code (CIMD), Codex CLI (DCR, loopback), a
   hosted client with a fixed callback — each a recorded flow the mock IdP satisfies.
5. Every "never/always/cannot" in the recipe traces to a rule and a fixture; a
   sentence that does not is removed or becomes a contract gap.

HARD RULES: docs and deploy artifacts only; no changes to contract pages; nothing
employer-internal; every command in the recipe actually run once against the real
binary (name the run).

DONE WHEN: a fresh reader goes from zero to a real sign-in following only the recipe +
`idp-request` output + `rehearse` for the console mode and one real IdP mode; **every
other mode section carries either a tested-on date or the literal banner
"UNVERIFIED — written from the contract, not from a run"** — no silently unproven
mode ships.

REPORT: the two dry runs; every step that required knowledge not in the recipe.
