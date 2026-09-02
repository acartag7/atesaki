MODEL: claude-opus-class or fable   EFFORT: high   TOOL: Claude Code in ~/project/atesaki-core (after M7 — the
live proof requires a working AS and the recipe; slice 1 alone cannot satisfy it)
MILESTONE: M8 (docs/roadmap.md). The name check (deliverable 1) runs on its own in
M0 — the repo is already public as `acartag7/atesaki` and the module path is in
`go.mod`, so a bad answer must arrive before more work is stamped with the name.
WHY: unshown work doesn't exist — and the name locks at first publish.

Read first: README.md · docs/decisions.md · seed/PROVENANCE.md (sanitization note) ·
docs/recipe.md if present · the `edictum-visibility` skill (positioning rules, verified-
claims traps) · the house rules on public copy (class-level only; no employer
customers, internal names, ROI or run counts).

DELIVERABLES
1. **Name check**, recorded in `docs/decisions.md` as an `[O]` row once Arnold picks:
   `atesaki` on GitHub (org + repo), Go module path, container registry namespace,
   npm/pypi collisions (informational), the search-query family and autocomplete;
   trademark screen. Report each as free / taken / congested.
2. **Sanitize**: `seed/` internal identifiers replaced with neutral examples; grep the
   whole tree for hostnames, tenant ids, AD group names, Vault paths, employer names;
   confirm `~/project/atesaki/evidence/` is NOT referenced by anything that ships (it
   never leaves the machine).
3. `SECURITY.md` and `LICENSE` exist since M0 — verify they still say the truth;
   `CHANGELOG.md` seeded; release notes template that leads with what the user can now
   do; the release workflow: binaries for linux/darwin with checksums, the container
   image pinned by digest, SBOM and provenance attestations (GitHub artifact
   attestations), `govulncheck` clean at the tag; the sanitization grep from
   deliverable 2 added to CI so it runs on every PR from here on.
4. README rewrite for strangers: lead with the pain killed ("front your MCP with your
   company login; the key never leaves the box; every grant says who, why, how long"),
   then the ten-minute path, then the honest scope (v0 slices shipped vs not), then the
   trust artifacts (contract set, deltas, decisions ledger, fixtures, threat model).
   Positioning line: agentgateway for companies whose IdP and network cooperate;
   Atesaki for the ones where they don't. No claim without a fixture or a receipt.
5. **Live verification before promotion**: a real sign-in and tool call through the
   published container with two real clients (name versions); the recipe followed
   verbatim from a clean machine.
6. Distribution loop plan (release → listings → proof → content): first three listings
   (awesome-mcp-gateways, mcp-auth compliance list, the MCP conformance registry when
   applicable) and the first post — drafted, not sent.

HARD RULES: nothing pushed public until Arnold says so; every public claim traced;
no AI mentions in commits/PRs/copy.

DONE WHEN: name decision recorded; tree sanitized (grep proof attached); live
verification named; README/SECURITY/CHANGELOG in a PR.

REPORT: name-check table; sanitization grep results; live-verification transcript
summary; the draft post and listings for Arnold's edit.
