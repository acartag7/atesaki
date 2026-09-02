MODEL: claude-opus-class or fable   EFFORT: high   TOOL: Claude Code in ~/project/atesaki-core (with the probe server from PR 4's evidence folder)
MILESTONE: M2 step 0 (docs/roadmap.md). Two days, no PR to product code; the output
is evidence in the private folder plus one-line conclusions in `docs/open-questions.md`
(#62, #53, #5) and a short PR like PR 4.
WHY: evidence has already overturned one ruling (#62: real clients send no custom
authorize parameters) and forced two others (#53, #5). Before the authorization-
server contract closes, what Codex CLI and Claude Code actually send and accept is
measured, not assumed — and measured again before publish.

Read first, fully: `~/project/atesaki/evidence/prm-probe-2026-09-01/` (the probe
server, run notes, `probe.log`) · docs/deltas.md D1 live check · docs/open-questions.md
#5, #53, #62 · docs/contract-grants.md G4, G14 · docs/roadmap.md §2 items 1–3 ·
STATE.md (the leftover local-scope probe entries to reuse or remove).

PROBES — each against the loopback probe server, each recorded with client name and
exact version, request log kept in the evidence folder, sentinels only:
1. **Authorize parameters.** Confirm for Codex CLI and observe for Claude Code (the
   interactive `/mcp` menu) that the authorize request carries only the standard
   parameters. Note any client setting that could add parameters (none is expected).
2. **Narrowed `scope`.** Return a token response whose `scope` is a strict subset of
   the request; does each client accept the token and call the tool? Does it retry,
   refuse, or re-authorize? (#53 fallback decision hangs on this.)
3. **Consent page with extra fields.** Serve a consent page carrying two extra form
   fields (purpose, duration) that POST with the approval; does the flow complete in
   the client's browser hand-off for both clients? (#62.)
4. **`approval_pending`.** Answer the authorize redirect with `error=access_denied`,
   `error_description=approval_pending`, and a `request_id` extension parameter;
   record exactly what each client shows the user and whether `request_id` is
   visible anywhere. (G14; the user must be able to find the id.)
5. **Per-install CIMD.** Repeat the PR-4 observation for Codex on a second machine
   or a fresh install; confirm the document id differs per install and is stable per
   server entry; record Claude Code's stable document URL. (#5.)
6. **Multi-segment paths.** Re-check Claude Code and Codex with a two-segment route
   path; record whether the claude.ai connector can be tested (quick tunnels failed
   last time — say if they still do).

HARD RULES: nothing employer-internal; no real tenant, hostname, or credential; the
probe server binds loopback only; remove the local-scope client entries you add;
conclusions go to the open questions as dated evidence notes, never as decisions.

DONE WHEN: six probes recorded with versions and dates; the evidence folder updated;
the three open questions carry their evidence lines; PR open.

REPORT: per probe, per client: what was sent, what was shown, what the user could
find; anything that changes a roadmap recommendation.
