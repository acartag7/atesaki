# Not this version

Each item names the seam it plugs into (`contract.md §11`). Nothing here constrains v0
beyond keeping the seam's port shape honest.

- **Local daemon** (the former chapter 0, Go, keychain, loopback-per-upstream): a
  sibling deployment shape of the same relay core. Demoted from core by
  `[O:2026-08-30]`.
- **Edictum decider, per tool call**: allow/block/**ask** on each `tools/call` through
  the decider port — holding an in-flight MCP request for a human. Grant-time decisions
  are already v0 (`contract-grants.md` G7) with **built-in rules only**; wiring an
  external decider (an Edictum Ruleset through the same seam) is this item; what's future is the per-call hold, its stream semantics, and
  client timeouts. Unreachable = block (fail closed); at grant time the public outcome
  is `temporarily_unavailable` `[O:2026-08-31]`.
- **Outbound OAuth client** (former chapter 5): Atesaki as the OAuth client to
  OAuth-protected upstreams — a credential-resolver strategy.
- **Token exchange (RFC 8693) and per-identity credentials**: further resolver
  strategies, already sketched in the starter-kit reference YAML.
- **Purpose-bound broker** (former chapter 2): leases bound to subject+client+run;
  waits on a trustworthy run-identity source (sandbox identity).
- **Kubernetes operator / CRDs**: the config is already resource-shaped; a CRD is the
  same YAML plus a controller. Only if demand shows up — the CNCF gateways own that
  layer today.
- **EMA / ID-JAG grant support**: when enterprise IdPs beyond Okta ship it.
- **Agent mailbox (Dengon)**: a separate product, not an Atesaki chapter
  (`~/project/ideas-new/mcp-mailbox-bootstrap.md`).
