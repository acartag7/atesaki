MODEL: gpt-5.6-sol   EFFORT: xhigh   FALLBACK: gpt-5.6-terra   TOOL: Codex CLI in ~/project/mcp-sso, resuming the PR #338 thread
WHY: the shared corpus is the rigor Atesaki inherits. Owner decisions are recorded on
PR #338 (six Phase-0 gaps, 2026-08-31 06:20Z; §19 simplification superseding the anchor
and fenced-block decisions, 2026-08-31). CI at 96422b5 is green. Proceed in this order; each phase is its own PR.

Read first, fully: docs/contracts/19-parity-fixture-protocol.md · docs/contracts/18 ·
fixtures/README.md, fixtures/schema/fixture.schema.json, the draft 8.4 fixture · the
owner-decision comments on PR #338 (the latest, "§19 simplification", supersedes the
anchor and fenced-block ones) · §05, §07, §08, §09, §10, §11 ·
docs/verification-design.md · ~/project/atesaki-core/docs/deltas.md (the consumer's
declared divergences — label, never pin against, D1/D2/D3/D4/D5/D6/D7/D11/D12/D13).
STALE-STATE NOTE (2026-09-02, evening): #338 is MERGED at 09fb858. PR #340 is
CLOSED — the runner landed on `main` as serial PRs (#374–#400; `main` at 63ed987
"run the parity runner in CI"). Phases A–C are DONE on `main`: `pnpm test:parity`
runs in CI and passes the two draft 8.4 fixtures through Fastify, Express, and Hono
with zero skips (verified locally 2026-09-02). Phase D is NOT done: both 8.4
fixtures are still `draft`, no receipt exists, and there is no `MANIFEST.json`,
`CATALOGUE.md`, or hash-gate script — §19.9 requires the hash gate to exist before
the first freeze. Start at D. Side finding for mcp-sso: three live-run script tests
(`live-evidence-scripts`, `live-run-script`, the symlinked live-state cases) fail on
macOS with any temp dir and pass in Linux CI — a portability defect in the live
run-support scripts, unrelated to parity. What Atesaki needs, in order
(docs/roadmap.md §4): (1) the runner passes the portable 8.4 draft and it freezes
with a receipt; (2) `MANIFEST.json` + `CATALOGUE.md` + the CI hash gate; (3) the §08
verifier slice frozen — M3's cross-lane input; (4) §07/§09/§10/§11 fixtures labeled
portable/host against `deltas.md` including the scope-ceiling row packet 14 adds
(catalog refusal at §9.3 step 3 becomes host) — M4's input.

PHASE A — §19 simplification (contract PR): rewrite §19.4's coverage gate to **clause
level** — `MANIFEST.json` lists, per numbered clause in §05–§17, the fixtures that pin
it and the clauses with none; a fixture names its clause and quotes its sentence (the
existing schema); drift = stale quote. Remove every mention of per-sentence anchors,
the anchor grammar, the marker-word statement selector, and the fenced-block ban.
Withdraw PR #339. Owner decision recorded on PR #338/#339 (2026-08-31, "§19
simplification").

PHASE B — evidence kinds (contract PR): MANIFEST evidence kinds `fixture`, `boot`,
`suite` per the owner decision; the portable logical store vocabulary (§12 appendix or
§19); `given.identity.checks` `result` | `throw {kind: oauth|generic}`; capture-and-
validate token chains (`capture`, exact header/claims, signature verification under
the corpus key); Fastify as canonical host + Express/Hono adapter-drift runs. Close the
unresolved P1 review thread with the `boot`/`suite` decision.

PHASE C — the reference runner (TypeScript): per §19.1/19.2 — clock port, seedable
randomness port (add narrowly if missing; production unchanged), corpus keys, recorded
outbound HTTP only, real HTTP against the composed app, exact comparison including
absence assertions and RE2 matchers; MANIFEST/CATALOGUE/FREEZE-LOG; CI: hash mismatch
fails, skipped frozen fixture fails, host and portable counted separately.

PHASE D — freeze 8.4: run it unchanged; freeze with receipt, or STOP with the §19.6
report (bug vs contract gap). Never edit the fixture to match the code.

PHASE E — slices, labeling every fixture portable/host: slice 1 §08 verifier + §09.1
PRM/challenge (**only fixtures pinning the reference's exact origin-root challenge
envelope are host** — the portable 8.4 shows the pattern: a location-flexible matcher
stays portable; Atesaki D1); slice 2 §09 bridge
core with §07/§10/§11 — registration (DCR stateless + stored, CIMD), consent, PKCE,
token exchange (exact burn/sign/store order → host where Atesaki D6/D12 diverge),
refresh rotation (family-only theft response → host where D7 diverges), revocation,
denial not consuming the JTI (host, D11), default `resource` (host, D13), stored-DCR
accumulation (host, D2). Flows are chains, no sleeps. Official conformance suite
scores none of this — exclude nothing, note the anchored version.

HARD RULES: one fixture, one clause, sentence quoted; sentinels only; never weaken a
fail-closed control; ambiguity = contract gap → STOP, record, move on (§19.6);
negative fixtures build the named case (real tokens for the other resource, real
duplicate occurrences); exact error/challenge shapes; fresh sentinel subjects.
SCOPE FENCE: no corpus relocation, no Go, no product features, contract edits only in
the dedicated phase-A/B PRs.

DONE WHEN: A–D landed; slice 1 frozen; slice 2 drafted or better; MANIFEST names every
uncovered clause with its profile; zero skipped frozen fixtures.

REPORT: frozen ids + receipts; drafts and why; portable vs host counts; contract gaps
(most valuable); implementation bugs; the parity status line format Atesaki will
consume.
