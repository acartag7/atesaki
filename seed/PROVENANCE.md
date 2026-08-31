# Seed provenance

These four files are the **starter kit** distilled 2026-07/08 from operating the
personal `splunk-gateway` spike (mcp-sso 0.2.0 fronting a Splunk MCP). They are
**evidence and seed material, not contract**: `docs/contract.md` cites their
conformance groups as `[K:§X]`.

Imported 2026-08-30 from the archived spike (`splunk-gateway.tar.gz`,
`atesaki/` directory) because the spike repo is not checked out on this machine
and a contract must not cite a path that may not exist.

**Known stale against the v0 contract (evidence, not v0 — do not read as the pin):**
`CONFORMANCE.md §E` line "never appears in any client-visible response header or body"
(contract narrows to Atesaki-authored output; upstream body reflection is an accepted
residual); `§H` multi-replica requirements and `NODE_ENV` checks (v0 is single-node Go);
`§K` mentions an `atesaki check` verb that does not exist (the verbs are `validate`,
`validate --deep`, `idp-request`, `rehearse`, `routes`, `grants`, `serve`);
`atesaki.reference.yaml` state store `redis` and audit `stdout` (v0 pins a durable
embedded SQLite store and JSONL); credential types `token-exchange` / `per-identity`
(future.md); its `policy` block is the **per-tool-call** `authorization.toolPolicy`
idea (future.md) — grant-time policy is v0 but shaped as `Route.spec.grant.policy`
(B1), not as the seed spells it; `atesaki.schema.json` admits Redis/Postgres/OTLP, has
no `none` credential, no grant/purpose/duration/machine-client fields, and
`^https://`-only URL validation (the canonicalization contract is now B3). `CONFORMANCE.md §A` spells the per-route metadata URL as `/<route>/.well-known/…`
(the contract uses RFC 9728 path insertion `/.well-known/oauth-protected-resource/<route>/mcp`, B3);
`§I` assumes a JWKS rotation contract the seed never declared (B4 now does); `§K` deep
`readyz` semantics predate `validate --deep`. The seed
models `server`/`routes`; the contract models `Gateway`/`Route` resources. Where the
seed and the contract pages disagree, the contract wins and the seed is history.

**Sanitize before any public push:** `atesaki.reference.yaml` still carries
internal-infrastructure identifiers (an internal CA bundle filename and an
internal host-naming pattern in examples). Replace with neutral examples before
this repo gets a public remote.
