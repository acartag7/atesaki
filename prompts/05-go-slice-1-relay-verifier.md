MODEL: grok-4.5   EFFORT: high   FALLBACK: gpt-5.6-terra   TOOL: Codex CLI or pi, in a fresh clone of ~/project/atesaki-core
PRECONDITION: contract tagged `contract-v0-freeze`; Atesaki fixtures hash-locked; gate
green; **mcp-sso's §8 verifier fixtures frozen with receipts** — the honest cross-lane
input: this slice's verifier is inherited §07/§08 behavior, proven against the shared
corpus, never re-authored here.
WHY: the first things that run — the Atesaki fixture RUNNER, then relay + verifier +
the operator verbs proven against it. Beyond the frozen §8 set it depends on NO
mcp-sso fixture — the fastest path to a binary you can point at a real MCP.

Read first, fully, at the freeze tag: docs/contract.md (all) · docs/contract-
boundaries.md (all) · docs/contract-grants.md G1 (`resource` rule), G12 (the flow-event classes this slice
emits), and G9 (verifier
uses `grant_id`/`exp`) only · docs/deltas.md D1, D4 · schema/** · fixtures/** for:
nevers 1–5 and 7 (never 6 is slice 2's — redirect identity does not exist here),
relay rules, ladder rungs, B3/B5/B6/B7 · prompts/README.md conventions.

SCOPE — build exactly this, nothing from slices 2–3:
- **Deliverable 0 — the Atesaki fixture runner (Go)**: loads a fixture in the Atesaki
  profile (packet 03, deliverable 0); composes the binary from `given.config` through
  its ports (clock; seeded randomness mirroring mcp-sso's HMAC-SHA256 counter stream;
  keys; recorded outbound HTTP only); sends `when.request`; compares exactly — status,
  headers (RE2 where stated), body, `then.events` by reason and class, `then.state`
  over the logical records. No expectations of its own; a skipped locked fixture
  fails. It also runs the frozen mcp-sso portable §8 fixtures (numeric-clause subset,
  same spine) — that run is how the verifier below is proven.
- `cmd/atesaki`: verbs `validate`, `validate --deep`, `routes` (verifier
  + relay path against a mock upstream and a locally minted token). `rehearse` is
  **not in this slice** — the contract's rehearse runs Atesaki's full AS flow against a
  mock IdP (slices 2–3); a stub that mocks Atesaki itself proves nothing. `serve` =
  relay + verifier only; **the OAuth routes are not registered in this slice** — a
  request to them is an unrouted path, `404` per B3: nothing reachable that the
  contract does not define, no invented status, no contract gap. PRM/AS metadata
  served per B3.
- Config: YAML stream → typed structs via the schema (02); **everything validates
  before any side effect**; refusals name resource/field/rule (B1); references (B2)
  incl. the every-open filesystem invariants; console loopback rule.
- Canonicalization exactly per B3 (host grammar, paths, reserved paths, collision,
  single-decode inbound matching, absolute-form authority tuple equality).
- Egress layer (§7): every outbound byte through named profiles; proxy from env by
  default; CA bundle per destination; failures name the hop.
- Relay (§6): header allowlists both ways; strip upstream `WWW-Authenticate`; map
  upstream 401/403 → `502 upstream_auth_failed`; SSE unbuffered; client disconnect
  cancels upstream; query-only passthrough; top-level array refused.
- Verifier: ES256 tokens, `iss` exact, `aud` = route audience byte-exact, `exp`,
  `scope` ⊇ `requireScope`, duplicate `Authorization` occurrences fail closed, the
  per-route challenge (D1, B7). Keys: the signing key ref (B1) — public part used here.
- Forwarded-header trust (B6) and Host/Origin gate before body parsing.
- Flow audit lines (G12 class F) to the JSONL sink; caps (B5) enforced before parse.

HARD RULES: implement only what the locked fixtures and the cited sections require;
no special-casing fixture strings; no invented library APIs (grep the module source);
stdlib first — every dependency justified in the PR with version and publish date and
age per the target repo's dependency policy file (read it; assume no number); **never touch docs/ or fixtures/** (the gate fails you) — a gap is a PR
comment, then continue; fail closed on every ambiguity; no `value || default` on a
security selector; `O_NOFOLLOW` opens; explicit directory modes.

VERIFY: the Atesaki runner green on every locked slice-1 fixture, zero skips, and
green on the frozen mcp-sso §8 portable set; point `serve` at a real local MCP server (any stdio→HTTP MCP
you can run) with a static-header credential and show a tool call succeed through the
relay with a locally minted token and be refused with a token for another route.
Name the real server you used.

DONE WHEN: all slice-1 fixtures green, zero skips; `validate`/`routes`/`serve` work on the three valid example configs; gate + lint green; PR opened.

REPORT: fixtures passed/failed by id; the real MCP you relayed; every contract gap;
every place the contract was ambiguous and what you did NOT decide.
