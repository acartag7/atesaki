# Threat model (seed)

**Status: DRAFT seed.** The full model precedes v0.1 (house rule: new surface means a threat model before v0.1, and every enumerated edge gets a negative test).

## Assets

Upstream credentials · the token-signing key · minted tokens (access, refresh, codes, consent tokens) · grants, grant requests, pre-approvals, and their hashes · the policy rules and their version · the store file (grant authority survives restart there) · the audit event rows and JSONL · the IdP client secret · assertion verification keys (JWKS) and signed assertions (rung 4) · machine-client secrets.

## Attackers and what they try

| Who | Tries | Held by |
| --- | --- | --- |
| Signed-in user without the required group | Requests high scopes, calls tools anyway | group ceiling at consent + scope check per request (`contract.md §3`) |
| Agent (or user) with a plausible-sounding purpose | Lies in the purpose field | purpose is evidence, never authority; enforcement is expiry, ceiling, policy, revocation (G1, G7) |
| Holder of a stolen refresh token | Rotates forever, or replays a consumed one | family and grant both revoked on replay or client mismatch (A10′); family expiry ≤ grant expiry; rotation re-checks grant status inside the mutation (G6 A10, G9) |
| Holder of a stolen access token after revocation | Keeps calling tools | until that token's `exp`, at most one access TTL, the stated bounded promise (G1); the operator's lever is a shorter TTL |
| Holder of a route-A token | Uses it on route B | per-route audience wall (`contract.md §3`, never #5) |
| Requester after an approval | Re-runs with wider scopes or longer duration | claim matches the original `requested_hash`; the grant is created with the approved values; nothing wider can be presented (G6 A6) |
| Two concurrent retries | Both claim one pre-approval | one atomic compare-and-consume; exactly one wins (G6 A6/A6a) |
| Requester whose group was removed after approval | Claims a stale pre-approval | current policy and ceiling re-evaluated at claim; approval bypasses escalation, never a current deny or narrower ceiling (G6 A6b `[O:2026-08-31]`) |
| Approver | Approves own request, or widens | not-your-own rule; narrow-only; audited with who, when, and hashes (G7) |
| Anyone who can guess grant or request ids | Inspects or revokes others' grants | ids library-generated, unguessable; `grants` authority contract pending (#24) |
| Flooder | Fills the store with pending 24-hour records | per-tuple cap + total cap (both refuse `temporarily_unavailable`, A3″), identical requests dedupe, awaiting pre-approvals expire, retention purge (A15) (G11, B5) |
| Whoever controls the decider's availability | Takes it down hoping requests sail through | unreachable decider = `temporarily_unavailable` on the redirect channel, never allow, never silent escalate (G7 `[O:2026-08-31]`) |
| Malicious website in a user's browser | Drives a browser-resident agent into the gateway | Host/Origin allowlists pre-parse (`contract.md §6`) |
| Any client | Smuggles its own upstream auth or pins connections | header allowlist outbound (`§6`) |
| Client sending two authorities | Absolute-form target with a different host than `Host` | `400` unless byte-equal (B3); no silent discard |
| Client smuggling path separators | `%2F`/`%5C` in the raw path, or double-encoding surviving one decode | raw-byte separator scan before the single decode; residual `%` after it refuses; not routed, no upstream invoked (B3) |
| Client spoofing forwarded headers | Picks its own rate-limit bucket or host | forwarded headers honored only from `trustedProxies`; strict parse; rightward walk (B6) |
| Holder of a stolen signed assertion (rung 4) | Replays it | valid until its `exp`; an accepted residual for v0 (no replay cache, B4 `[O:2026-08-31]`), mitigated by short proxy assertion lifetimes |
| Whoever can serve or poison the JWKS | Feeds a rogue key | keys fetched via the configured egress profile only; `kid` exact; `alg` pinned in config; stale beyond the maximum interval refuses (B4) |
| Client supplying identity headers on non-identity paths | Impersonates via `x-forwarded-user` etc. | identity headers dropped before routing everywhere but the identity leg (B4) |
| Malicious or compromised upstream | Reflects its credential in a body; sends a challenge to loop clients | challenge stripped + status mapped to `502 upstream_auth_failed` (`§6`, never #3); body reflection is a named residual `[O:2026-08-16]` |
| Local attacker with filesystem access | Swaps the store or audit file for a symlink or a foreign file between restarts | every open re-checks `O_NOFOLLOW`, regular file, link count 1, owner, mode, parent directory (B2) |
| Anyone who can read the console pairing code | Becomes `console-operator` | console mode boots only on loopback; loud never-multi-user warning (`contract.md §13`) |
| Network position inside the cluster | Reaches the authless backend directly | the product cannot enforce this; recipe obligation (§14) + `validate --deep` warning; see Known-accepted |
| Operator mistake | Inline secret, empty allowlist, duplicate audience, wrong proxy, non-canonical URL, a machine declaration its own route's rules deny | boot refusals (`§2`, B1–B3, G7 contradiction check), deep readiness (`§9`) |
| Machine whose binding was revoked | Asks again after any restart or after the operator edits an unrelated route | tombstone bound to the per-(client, route) declaration digest; only a deliberate change to that binding clears it (G3, G10) |
| Client omitting `resource` | Hopes the gateway picks a route for it | `invalid_target`; no default resource in a multi-route gateway (G1, D13) |
| Whoever reads config, logs, or audit | Harvests credentials | references-only config; the credential appears in nothing Atesaki authors (`§2`, `§6`, `§10`) |

## Inherited surface

The authorization-server surface (authorize, token, register, consent, metadata) inherits mcp-sso's threat model and mitigations except where `contract-grants.md` changes the behavior. Those changes are the deltas D2–D13 in `deltas.md`, each with its own fixtures. Parity for the unchanged parts is proven by the portable corpus.

## Known-accepted (named, not hidden)

- A transparent relay does not sanitize upstream response bodies. `[O:2026-08-16]`
- Audit delivery can fail without failing requests. Durable evidence is the event row committed with the change (G6/G12 `[O:2026-08-31]`); delivery is best-effort and loud. `[O: starter kit]`
- The product cannot enforce that only the gateway reaches an authless backend. That is a recipe obligation with an advisory warning. `[O:2026-08-31]`
- A stolen rung-4 assertion is valid until its `exp`; no replay cache in v0 `[O:2026-08-31]`.
- Issued access tokens survive revocation for at most one access TTL. `[O:2026-08-30]`
- Store growth is rate-bound, not size-bound. Terminal rows purge after retention, but `grant_event` rows never purge in v0 and denied or abandoned rows accrue until retention, so sustained authorized traffic grows the store at whatever rate the per-IP budgets allow. Disk exhaustion is a named residual; the levers are budgets, retention, and monitoring.
