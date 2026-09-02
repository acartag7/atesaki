MODEL: gpt-5.6-sol   EFFORT: high   TOOL: Codex CLI in ~/project/atesaki-core (after M6)
MILESTONE: M7 (docs/roadmap.md). PRECONDITION: `rehearse` merged (packet 15); the
client-matrix staleness window ruled (B8 row, packet 14 item 9).
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
2. The recipe embeds the `idp-request` output the binary already produces (packet
   06) — it designs no IdP ask of its own.
3. `Dockerfile` (distroless/static; non-root; read-only root) and `deploy/kustomize/`
   example: Deployment, Service, Ingress with the `atesaki routes` path list,
   NetworkPolicy derivable from documented egress, secrets as refs. Kubernetes facts
   the contract states (§14, #56) and the kit must embody, pinned to the cluster
   version it was tested on: secrets, CA bundles, and CIMD documents arrive as
   `env:` from Secret keys (today's default volume mounts are root-owned symlinks and
   B2 refuses them by design); one supported volume type named (a ReadWriteOnce
   PersistentVolumeClaim) mounted over a path the image does not own, with `fsGroup`
   and `fsGroupChangePolicy` so the non-root process can create the store
   subdirectory (`0700`) on first boot — **first boot and restart proven on a real
   cluster**, never assumed; the ingress listed in `trustedProxies` (otherwise every
   user shares one rate-limit bucket) and applying the same request-framing rule as
   Go's parser; **no path rewrite** at the ingress (audiences are byte-exact);
   `livez`/`readyz` exactly as #61 rules them; TLS terminates at the ingress and the
   issuer still derives from `externalBaseUrl`; image pinned by digest;
   `readOnlyRootFilesystem: true`.
4. The client matrix consumes the `rehearse` profiles packet 15 shipped (Codex CLI
   is per-install CIMD, per the live evidence — not DCR) and the live proofs from
   packets 06 and 07; this packet designs no client flow and edits no fixture.
5. Every "never/always/cannot" in the recipe traces to a rule and a fixture; a
   sentence that does not is removed or becomes a contract gap.
6. Operations section states plainly, with the enforcing rule: access ends within
   one access TTL of revocation (G1); a retried refresh after a lost response ends
   the grant and the agent signs in again (A10′); a group removed after activation
   does not revoke the grant — refresh rechecks nothing until expiry or revocation,
   the levers are a shorter `maxDuration` and `grants revoke`; what a store-file loss
   means (key rotation, every token dies); audit rotation preserves every line
   already written (reopen, never truncate) while flow-event loss stays the accepted,
   counted residual (G12) — no losslessness claim; how active streams end at
   shutdown (#61 as ruled); the exec audit trail of the platform is the
   accountability source for CLI approvals in a container (#24 as ruled).

HARD RULES: docs and deploy artifacts only; no changes to contract pages; nothing
employer-internal; every command in the recipe actually run once against the real
binary (name the run).

DONE WHEN: a fresh reader goes from zero to a real sign-in following only the recipe +
`idp-request` output + `rehearse` for the console mode and one real IdP mode; **every
other mode section carries either a tested-on date or the literal banner
"UNVERIFIED — written from the contract, not from a run"** — no silently unproven
mode ships.

REPORT: the two dry runs; every step that required knowledge not in the recipe.
