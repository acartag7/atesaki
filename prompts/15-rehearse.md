MODEL: grok-4.5   EFFORT: high   FALLBACK: gpt-5.6-terra   TOOL: Codex CLI or pi, fresh clone
MILESTONE: M6 (docs/roadmap.md). PRECONDITION: slice 3 merged (packet 07); the
Atesaki fixture profile and runner in place (packets 03, 05).
WHY: onboarding step 4 — "before anything deploys, the whole flow runs on your
machine against a mock IdP, per client you care about". A rehearsal that mocks
Atesaki itself proves nothing; this one runs the real binary through its ports.

Read first, fully: docs/contract.md §9 (`rehearse` row), §4, §13 · docs/onboarding.md ·
docs/contract-boundaries.md B1 (console rules), B4 · fixtures/schema/atesaki-fixture.schema.json
and the runner package (packet 05) · docs/roadmap.md §M6.

SCOPE — build exactly this:
- `atesaki rehearse <config.yaml> [--profile <name>]`: composes the real binary
  through the runner's ports (clock, seeded randomness, recorded outbound HTTP, the
  memory store adapter — accepted here and only here), binds the listener and a mock
  IdP on loopback only, and drives each client profile as a fixture **chain** in the
  Atesaki profile: discovery → registration (CIMD or DCR per profile) → authorize →
  identity leg against the mock IdP → consent with `purpose`/`requested_duration` →
  token → one `/mcp` call against a recorded upstream. Per configured rung.
- Client profiles as recorded chains under `fixtures/rehearse/`: Claude Code (stable
  CIMD URL), Codex CLI (per-install CIMD document, loopback redirect), a DCR loopback
  client, a hosted client with a fixed callback. Each profile records the exact
  requests the real client sent in the live probe (`docs/deltas.md` D1 live check)
  — sentinels only, no real ids.
- Output: per profile and rung, the step reached and the exact refusal if any
  (resource, field, rule, or public error code); never a secret, never a token.
- A config that cannot rehearse (e.g., `identity.provider: header` with a JWKS URL)
  says why in one line and what the profile would need.

HARD RULES: loopback only for both listeners, refused otherwise before any state
write; never contacts the real IdP, a real upstream, or the network at all (the
runner's egress port fails any unrecorded call); the memory adapter is never
reachable from `serve`; no new behavior in the authorization server — a rehearsal
that needs one is a contract gap; prompts/README.md conventions.

VERIFY: every profile green on the three valid example configs; a deliberately
broken config (wrong redirect URI) fails at the named step; the onboarding page's
step 4 is literally true, run on a laptop and named.

DONE WHEN: above verified; packet-11 review clean; PR opened.

REPORT: profiles and rungs passed; the step each broken config stopped at; every
contract gap; anything a real client did in the live probe that the profile could
not encode.
