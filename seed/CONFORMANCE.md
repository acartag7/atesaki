> **OBSOLETE: evidence only, not normative.** Imported from the starter kit (`PROVENANCE.md`). Several rows contradict the current contract: discovery is RFC 9728 path insertion under `/.well-known/…/<route>/mcp` (B3), not `/<route>/.well-known/…`; a malicious upstream reflecting its own credential in a relayed body is an accepted residual, not a promise (`contract.md §6`); multi-replica stores are out of v0 (`§13`); readiness excludes upstream reachability (open question #61); `atesaki check` is the predecessor of `validate --deep`. The contract pages win everywhere they differ.

# atesaki conformance checklist

Test cases for an MCP OAuth gateway, grouped by concern. Each item names the **failure mode it guards against**. Most were observed live operating an mcp-sso-based gateway into a real enterprise cluster, not invented.

Legend: `⇒` = expected result. Codes are HTTP unless noted.

---

## A. OAuth discovery & metadata (per route)

- [ ] `GET /<route>/.well-known/oauth-protected-resource/<mcpPath>` ⇒ `200`, `resource` == the route's audience, `authorization_servers` == issuer. _Guards: multi-path discovery. RFC 9728 path-insertion must resolve per route, not just at host root._
- [ ] `GET /.well-known/oauth-authorization-server` ⇒ `200`, issuer + all endpoints absolute under `externalBaseUrl`.
- [ ] Metadata origin derives from `externalBaseUrl`, **not** the request `Host`. _Guards: a spoofed Host rewriting discovery URLs._
- [ ] `GET /oauth/jwks` ⇒ active `kid` present; rotated-out `kid` still verifies until its tokens expire.

## B. Fail-closed auth surface

- [ ] `POST /<route>/mcp` no token ⇒ `401` + `WWW-Authenticate: Bearer resource_metadata="…"` pointing at THIS route's metadata.
- [ ] Malformed / wrong-signature / expired token ⇒ `401`.
- [ ] Valid token, wrong issuer or wrong `aud` ⇒ `401`. _Guards: token replay across gateways/paths (see D)._
- [ ] Valid token missing `required` scope ⇒ `403`.
- [ ] DCR with non-allowlisted `redirect_uri` ⇒ `400 invalid_redirect_uri`.
- [ ] DCR with `redirect_uri: "*"` or wildcard ⇒ rejected at config-load AND at register.

## C. DNS-rebinding (Host / Origin)

- [ ] Foreign `Origin` on `/<route>/mcp` ⇒ `403`, checked in the pre-handler before the body is read.
- [ ] Unlisted `Host` ⇒ `403`.
- [ ] Absent `Origin` (native client) ⇒ allowed; browser cross-site `Origin` ⇒ blocked unless allowlisted.
- [ ] Empty allowlist ⇒ fails **closed** (deny), never open.

## D. Multi-path isolation: the new surface

- [ ] Token minted for `/splunk` (aud=`…/splunk/mcp`) replayed at `/conduktor` ⇒ `401`. _Guards: the core multi-tenant invariant. If this passes-through, multi-path is a lateral-movement hole._
- [ ] Each route advertises only its own scope catalog; consent for route A does not grant route B.
- [ ] Two routes may share one IdP but produce distinct audiences.
- [ ] Config-load rejects duplicate `resource` values across routes.
- [ ] Route `path` collision or a route path shadowing `/oauth`/`/.well-known` ⇒ config error.

## E. Upstream relay & the 403 masquerade

- [ ] Upstream returns `401`/`403` (gateway's own credential insufficient) ⇒ client sees a **non-`401`** error (e.g. `502 upstream_auth_failed`), and the upstream's `WWW-Authenticate` is **stripped**. _Guards: the re-auth loop. A relayed `403` makes OAuth clients re-login forever, re-auth mints a valid token the upstream refuses identically._
- [ ] Client token failure (gateway-side) still ⇒ `401` + challenge. The two auth failures are distinguishable to the client.
- [ ] Upstream `5xx`/unreachable ⇒ `502` with a JSON-RPC-safe error, cause logged.
- [ ] Response headers relayed by allowlist; hop-by-hop + upstream auth dropped; `Mcp-Session-Id`/`Mcp-Protocol-Version` preserved.
- [ ] Injected credential never appears in any client-visible response header or body.

## F. Egress: proxy + CA per destination

- [ ] IdP token/JWKS calls traverse the route's configured proxy (not env-implicit). _Guards: Node fetch ignoring `HTTPS_PROXY`, the call that hangs first in a proxied network._
- [ ] Upstream on a no-proxy internal net is reached **directly** even when a proxy profile exists globally.
- [ ] Upstream TLS validates against the route's `caBundleRef` without polluting global trust. _Guards: internal-CA upstreams (e.g. a corporate device CA) that a stock Node image doesn't trust._
- [ ] `insecureSkipVerify: true` emits a loud warning and is refused when `NODE_ENV=production`.

## G. Transport correctness (Streamable HTTP)

- [ ] `POST /<route>/mcp` with JSON body ⇒ body forwarded byte-for-byte (no re-serialization).
- [ ] `GET /<route>/mcp` (SSE) ⇒ stream relayed unbuffered; events arrive incrementally.
- [ ] Long-lived SSE is **not** severed by an idle timer when `timeouts.idle: none`.
- [ ] Client disconnect mid-stream aborts the upstream fetch. _Guards: cancellation. Bind the abort to the RESPONSE socket close, not the request stream, a buffered POST closes its request stream in ~1ms and would abort every call before it starts._
- [ ] Normal buffered `POST` completes without spurious abort (regression test for the above).
- [ ] `DELETE /<route>/mcp` session-end forwarded; `Mcp-Session-Id` round-trips both directions.

## H. State & multi-replica

- [ ] Access tokens verify with **no** store lookup (stateless hot path).
- [ ] Auth code / consent JTI is single-use; replay ⇒ rejected.
- [ ] With ≥2 replicas + shared store, an auth ceremony started on replica A completes on replica B.
- [ ] `store: memory` with replicas>1 ⇒ loud startup warning.
- [ ] SSE session affinity honored (a session's stream stays on its replica/upstream).

## I. Observability

- [ ] No `catch {}` on any request path; every caught error logs its full `cause` chain with the syscall code.
- [ ] `readyz` returns non-200 when: any upstream unreachable, IdP JWKS unreachable, signing key unloaded, or store unwritable. _Guards: a shallow `200` health check masking a fully broken upstream._
- [ ] `livez` and `readyz` are distinct; liveness stays green during a transient upstream outage.
- [ ] Audit-sink write failure ⇒ request still succeeds AND a loud error/metric is emitted (never silent). _Guards: a root-owned/again-unwritable audit file silently recording nothing._
- [ ] `atesaki check <route>` self-tests IdP reachability, upstream reachability + TLS, signing keys, redirect URIs, store writability.

## J. Secret non-exposure

- [ ] No secret material accepted inline in config (schema enforces `*Ref`).
- [ ] Secrets never passed as process argv. _Guards: tokens visible in `ps` (observed with an `mcp-remote --header "Authorization: Bearer …"` invocation)._
- [ ] Log redaction covers every sink, including error messages and structured fields.

## K. Deployment ergonomics

- [ ] Default bind `0.0.0.0` in a container; reachable via pod IP behind an ip-target LB.
- [ ] Runs under `readOnlyRootFilesystem: true` (no writes to cwd; state dir is mountable).
- [ ] Health path is configurable to match the orchestrator (avoids `/healthz`-vs-`/health` target-group flap → 502).
- [ ] Whole config validated at boot against the schema; any missing secret/ref fails **closed** at startup with an actionable message, never at request time.
- [ ] Shipped chart documents required egress ports (IdP-via-proxy, proxy port, each upstream port) so the NetworkPolicy is derivable.
