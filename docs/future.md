# Not this version

Each item names the seam it plugs into (`contract.md §11`). Nothing here constrains v0 beyond keeping the seam's port shape honest.

- **Local daemon** (the former chapter 0: Go, keychain, loopback per upstream). A sibling deployment shape of the same relay core. Demoted from core by `[O:2026-08-30]`.
- **Edictum decider, per tool call.** Allow, block, or ask on each `tools/call` through the decider port, holding an in-flight MCP request for a human. Grant-time decisions are already v0 (`contract-grants.md` G7) with built-in rules only. Wiring an external decider (an Edictum Ruleset through the same seam) is this item; the per-call hold, its stream semantics, and client timeouts are the future part. Unreachable means block (fail closed); at grant time the public outcome is `temporarily_unavailable` `[O:2026-08-31]`.
- **Outbound OAuth client** (former chapter 5). Atesaki as the OAuth client to OAuth-protected upstreams, a credential-resolver strategy.
- **Token exchange (RFC 8693) and per-identity credentials.** Further resolver strategies, already sketched in the starter-kit reference YAML.
- **Purpose-bound broker** (former chapter 2). Leases bound to subject, client, and run. Waits on a trustworthy run-identity source (sandbox identity).
- **Kubernetes operator and CRDs.** The config is already resource-shaped; a CRD is the same YAML plus a controller. Only if demand shows up. The CNCF gateways own that layer today.
- **EMA and ID-JAG grant support.** When enterprise IdPs beyond Okta ship it.
- **Agent mailbox (Dengon).** A separate product, not an Atesaki chapter (`~/project/ideas-new/mcp-mailbox-bootstrap.md`).
