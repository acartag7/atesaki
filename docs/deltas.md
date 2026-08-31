# Intentional differences from the mcp-sso reference

**Rule (`[O:2026-08-30]`):** Atesaki passes every frozen *portable* fixture in the
shared corpus with zero skips — portable fixtures are never skipped and never edited by
Atesaki. A behavior Atesaki intentionally changes is listed here first; mcp-sso labels
the fixtures pinning the reference-only behavior *host*; Atesaki writes its own. A
behavioral difference not on this list is a bug in Atesaki.

Tags: `[O:date]` owner decision · `[D ← rule]` derived consequence (reversing it means
changing the named rule) · `[P]` proposed, awaiting Arnold.

| # | Difference | Reference | Atesaki | Fixtures |
| --- | --- | --- | --- | --- |
| D1 `[O:08-30]` | Discovery / challenge location | Origin-root PRM; one bridge = one resource | Per-route path-inserted PRM (B3); each route's challenge points at its own | mcp-sso `08-…/8.4-duplicate-authorization-fails-closed` is **host**; Atesaki adds route-scoped challenge fixtures |
| D2 `[O:08-30]` | Scope accumulation | Stored-DCR re-auth unions active refresh records (§11, §9.3 step 5) | None — bounded lineage (G1); unreachable under the stateless-DCR pin, listed so no mode reintroduces it | §09/§11 accumulation → host; Atesaki: second grant inherits nothing |
| D3 `[D ← G4/A1]` | Consent-token claims | client, resource, scopes, ceiling (§7.1) | + `purpose`, `approved_duration_s`, `request_id`; `resource` remains the single route authority | §07.1 exact-claims → host |
| D4 `[D ← G9]` | Access-token claims | `sub`, `client_id`, `scope`, `iss`, `aud`, `iat`, `exp` (§7.2) | + `grant_id`; machine tokens additionally `grant_kind=machine` | §07.2 exact-claims → host |
| D5 `[O:08-30; carriers O:08-31]` | Authorize flow | Steps 1–4 → sign consent, one pass (§9.3) | Policy step after step 4, before 6 (G4); `purpose`/`requested_duration` are authorize parameters carried through the singleton guard, the §17.11 signed params and synthetic callback, and console pairing's round trip; escalation ends the pass; re-run claims (A6) | §09.3 single-pass → host; Atesaki: escalate, dedupe, claim win/lose/freshness, outage, per-tuple cap full, total cap full |
| D6 `[D ← A8/A9]` | Authorization code | Bound to client, redirect, resource, PKCE (§7.3); exchange **burns code → signs response → stores family** (§9.4); binding failure consumes the code, wrong `resource` does not | + bound to `grant_id`, `request_id`; exchange checks grant `issued` ∧ not revoked; **consumption on failure inherited unchanged** (A9′); **sign-before-commit** `[O:2026-08-31 G8]`, family and activation committed atomically — no burn-without-family state | §07.3/§09.4 exact order → host; Atesaki: code for revoked grant refuses (consumed); wrong-resource not consumed; commit failure leaves nothing consumed |
| D7 `[D ← A10/A10′]` | Refresh family | Rotation **mutates then compensates** failures by revoking the family (§7.4); replay or client mismatch revokes the **family** | Carries `grant_id`, `grant_expires_at`; rotation checks grant status inside the tx; successor expiry ≤ grant; replay/mismatch revokes **family and grant** (bounded lineage); sign-before-commit `[O:2026-08-31 G8]`, no compensation path | §07.4 rotation order / family-only theft response → host; Atesaki: rotation at/after grant expiry refuses; replay revokes the grant |
| D8 `[D ← G2]` | Store contents | codes, consent JTIs, families, revocations, instance (§12) | + `grant_request`, `preapproval`, `grant`, `grant_event` (outbox), `machine_tombstone` — exact logical schema is corpus work (#36) | §12 exact row sets → host; Atesaki: own logical rows |
| D9 `[D ← A11]` | Revocation | RFC 7009 revokes the family | `grants revoke` **and** RFC 7009 on any token of the family both revoke the **grant** and its family in one tx; access tokens live to `exp` (bounded, `[O:08-30]`) | §09 revoke family-only → host |
| D10a `[O:08-31]` | Machine registration & credentials | Stored-DCR machine clients, generated `mcc_…` ids, out-of-band lifecycle (§17.2) | Declared in config (`machineClients[]`), B3 identifier grammar, secret by reference; only the token-endpoint grammar and client authentication are inherited | §17.2 registration/lifecycle → host |
| D10b `[D ← A12]` | Machine token claims | `sub`/`client_id`/`gty` triad over an `mcc_` id (§17.2) | `sub` = `client_id` = declared B3-grammar id; `gty=client_credentials` kept (verifier classification inherited); + `grant_id`, `grant_kind=machine`; scopes = requested ∩ declared | §17.2 exact claims → host |
| D10c `[D ← A12/A13]` | Machine lifecycle | Stateless issuance; no grant object | One `active` machine grant per (client, resource); issuance attaches or creates; revocation writes a **tombstone** bound to the declaration digest `[O:2026-08-31]` | Atesaki-only fixtures |
| D11 `[O:08-31]` | Consent denial | Denial does **not** consume the consent JTI (§9.3) | Denial consumes it; the request is terminal `abandoned` (A7) | §09.3 denial-reusable → host |
| D13 `[D ← D1 multi-route; RFC 8707]` | `resource` parameter | Omitted `resource` defaults to the single configured resource (§9.3 step 3) | **Required** on authorize and machine-token requests; must equal a route audience; omitted/unknown → `invalid_target` (G1) | §09.3 default-resource → host; Atesaki: omitted resource refused |
| D12 `[D ← G8]` | Failure atomicity | Post-consumption clock failure leaves the JTI consumed by design (§9.3) | All fallible steps precede the transaction; JTI consumption commits with the grant; failure before commit consumes nothing | §09.3 consumed-on-late-failure → host |

That is the whole list. A growing list is a design smell to raise, not a convenience.
