MODEL: grok-4.5   EFFORT: high   FALLBACK: gpt-5.6-terra   TOOL: Codex CLI or pi, fresh clone
MILESTONE: M4 (docs/roadmap.md). PRECONDITION (#55): slice 1 merged (its runner
already executes the mcp-sso §8 portable subset); packet 03 phase 2 fixtures locked;
packet 14 landed (#53, #5, #56, #57, #58 as ruled); the mcp-sso corpus pinned at an
exact commit + `MANIFEST.json` hash, and the **exact frozen portable fixture-id set**
this build will pass listed in the PR before any code — acceptance is by ID, not
section name; every frozen portable id this slice does not yet pass is listed as
deferred with its reason (the lane may lag; nothing is skipped silently).
PINNED: contract.md @ <sha> · contract-grants.md @ <sha> · contract-boundaries.md @
<sha> · deltas.md @ <sha> · fixtures/MANIFEST.json @ <sha> · mcp-sso corpus <version>
@ <sha>.
WHY: the first real sign-in. Parity by fixture, not by memory — and the store, the
two-phase discipline, and the allow path land here so the rows with humans and
machines in them (M5) build on a proven base.

Read first, fully: ~/project/mcp-sso/docs/contracts/ §01–§12, §14, §17 (identity
ports, machine clients), §19 at the pinned SHA · mcp-sso/fixtures/** (portable only)
· docs/contract.md §2, §3, §4, §5, §7, §8, §9 (`idp-request`), §13 ·
docs/contract-grants.md G1–G9, G12, G14 (A1, A2, A3, A3′, A3″, A7, A8, A9, A9′, A10,
A10′, A10″, A11, A14 for these rows, E1–E3 are this slice's rows) ·
docs/contract-boundaries.md (all) · docs/deltas.md (every row — each changes what
you would otherwise inherit; read the row before implementing the clause it touches)
· schema/records/** · docs/negative-matrix.md rows tagged slice 2 · docs/roadmap.md
§M4 (gotcha rows 1, 2, 3, 6–10, 17, 21, 23, 27) · docs/decisions.md (why each rule is
the way it is — do not relitigate).

SCOPE — serial PRs, one behavior each, in this order:
1. `feat(store): port, memory adapter, conformance suite` — the Go interface with
   exactly the operations G6 needs (one transaction per A-row: authoritative reads,
   CAS predicates, mutations, durable events); the conformance suite as the contract
   (atomicity per row, CAS semantics, uniqueness, ordering-free comparison); memory
   passes; `serve` refuses memory (accepted by the runner and, later, `rehearse`).
2. `feat(store): sqlite adapter` — pure-Go driver chosen and pinned **first** (name
   it, version, publish date, age), and its open semantics proven by an integration
   test before any row lands: how it opens the main file, how it creates `-wal` and
   `-shm` (SQLite's Unix VFS gives them the database file's mode; observe what this
   driver does), whether it honors a no-follow open flag. B2 for the store is then:
   the directory is Atesaki's (`0700`, created by it, held open); the database file is
   created by Atesaki with `O_EXCL` and `0600` before the driver sees the path; after
   the driver opens, re-open the path with `O_NOFOLLOW`, `fstat`, and compare the
   inode with the file Atesaki created; sidecars re-checked under B2 on every open.
   A B2 sentence the driver cannot meet is a contract gap reported before this PR
   merges — never a weakened check. WAL; `busy_timeout`; `BEGIN IMMEDIATE` for every
   writing transaction (a deferred transaction that reads then writes fails with
   `SQLITE_BUSY_SNAPSHOT` under a concurrent writer, and the busy handler cannot
   retry it); `synchronous=FULL`; passes the same suite, also under contention (two
   goroutines on the same CAS).
3. `feat(as): metadata, stateless DCR, CIMD` — origin AS metadata (fields per
   §9.1/RFC 8414 as inherited; `scopes_supported` per the #53 ruling); stateless DCR
   per §9.2 with the redirect allowlist (§10); vendored CIMD from `knownCimd` refs;
   live fetch behind `clients.cimd.liveFetch` only if #5 allowed and exactly as
   ruled: the document origin must be on the exact allowlist before any network
   call; the inherited caps by clause; the validated-IP dial when the profile is
   direct; CIMD document validation.
4. `feat(as): authorize steps 1–4, ceiling, purpose and duration` — §9.3 steps 1–4
   exactly (`resource` required and exact, D13); effective scopes = requested ∩
   catalog ∩ group ceiling (#53) with `scope_ceiling_applied`; G-a parsing of
   `purpose`/`requested_duration` with the singleton guard across every carrier (D5:
   direct authorize, the §17.11 signed bridge params and synthetic callback, console
   pairing's round trip with the equal-or-`invalid_request` reconciliation); G-b
   `requested_hash` over raw scopes (G3, RFC 8785 — a pinned JCS implementation
   after supply-chain review, or a purpose-written canonical serializer for the fixed
   G3 shapes; either is tested against the RFC 8785 vectors; Go's `encoding/json` is
   not JCS and no "close enough" argument is accepted).
5. `feat(policy): built-in rules` — G7 vocabulary including `clientOriginIn` (#57);
   AND-only; deny overrides allow; escalate by default; purpose unreadable;
   `policy_version` digest; the boot contradiction check already in `validate` learns
   the new condition.
6. `feat(grants): A1, A2, A3 insert, A3′, A3″, lazy expiry` — consent token signed
   with `request_id`, `purpose`, `approved_duration_s` (D3); `allow` → consent page;
   `deny` → `access_denied`; escalate → insert `grant_request` + `preapproval` per
   A3 with dedupe and both caps, `approval_pending` + `request_id`; decider error →
   `temporarily_unavailable`. **This slice owns A14's lazy path** for every row it
   reads: the cap and dedupe read in A3, and the reads in A9/A10, first transition
   past-due rows they touch (G5) inside the same transaction with their durable
   events, or expired pending rows keep consuming the cap. M5 adds the sweeper. No
   approver exists yet: every pre-approval expires — contract-conformant, not an
   interim state.
7. `feat(grants): consent, exchange, rotation, revocation` — A7 (denial consumes the
   JTI, D11), A8, A9 (sign-before-commit, G8; family and activation committed
   atomically, D6; `grant_expires_at` from activation, G9), A9′ (consumed on binding
   failure and revoked grant; not consumed on wrong `resource`), A10, A10′ (replay or
   client mismatch revokes family AND grant, D7), A10″, A11 via RFC 7009 (D9), E1
   (`response_not_delivered`), E2, E3; access `exp = min(now + accessTTL,
   grant_expires_at)`; tokens carry `grant_id` (D4). Timing-safe comparison for
   client secrets (compare digests with `subtle.ConstantTimeCompare`).
8. `feat(identity): entra, oidc` — redirect flow, id_token verification (issuer,
   audience, nonce, expiry, signature under the JWKS fetched through
   `identity.egressProfile`), groups from the claim only; Entra **exactly as mcp-sso
   §17 says**: a groups overage marker is an identity refusal with its named audit
   reason (never a Graph call, never a silent empty ceiling); group identifiers are
   object ids (GUIDs) only — a display name as a `groupsToScopes` key is a boot
   refusal; the subject boundary §6.5; identity-failure pairs (rejection vs port
   throw) per §19.2; an `inherited` Atesaki fixture for every §17 clause with no
   frozen upstream fixture.
9. `feat(identity): header assertion, console pairing` — B4 in full: pinned `alg`
   per key, `kid` exact, JWKS through the profile with size/count caps, stale-interval
   refusal, bounded refetch on unknown `kid` (#58 as ruled), header exactly once,
   identity headers dropped everywhere but the identity leg; console pairing
   loopback-only before any state write, the pairing code printed and never audited,
   fixed `console-operator` subject, the loud tutorial warning.
10. `feat(cli): idp-request` — per provider (Entra dedicated, Entra shared, generic
    OIDC, header, console): the minimal ask (app type, the single redirect URI, the
    groups claim as object ids with group filtering or app-assigned groups so overage
    cannot occur, secret delivery) and the explicit does-not-need list
    (no Expose-an-API, no `api://` scopes, no per-client redirect churn); golden
    output per provider; contains no secret.
11. `test(e2e): real sign-in to tool call` — the named real input (below).

HARD RULES: as prompts/README.md. Every place Atesaki deviates from an mcp-sso clause
MUST correspond to a row in docs/deltas.md — a deviation without a row is a bug, stop
and report. Never edit a frozen fixture. Pin the corpus version and manifest hash in
the PR. A predicate that cannot be evaluated where G6 puts it (preflight vs in-tx) is
a contract gap — report, do not move it. Never weaken CAS or atomicity to pass a
fixture. Rate limits per client IP after B6 on register, authorize, token, approve,
revoke with the B8 budgets; limiter failure per #60 as ruled — never fail open on a
token-issuing path.

VERIFY: the pinned portable fixture-id set — zero skips; Atesaki phase-2 fixtures
green; the store conformance suite green on both adapters; crash tests at named
failpoints (between preflight and commit, between commit and response) for A8 and A9
with restart and exact state assertions; a real MCP client — Claude Code AND Codex
CLI, versions named — completing discovery, registration, sign-in through Entra or
generic OIDC, a consent page showing purpose and duration, a tool call on a route
with an `allow` rule, and `approval_pending` + `request_id` on a route without one;
Codex on two routes with different catalogs (the #53 case) succeeding; the console
rung on a laptop; RFC 7009 revocation observed (refresh refused, access dies at TTL).

DONE WHEN: parity status line published in the PR (corpus version passed, current,
frozen fixtures not yet passed with reasons); every inherited clause this slice
implements has a frozen upstream or an `inherited` Atesaki fixture, or is named
uncovered and excluded from the capability claim; real sign-in shown; packet-11
review clean.

REPORT: parity line; deferred fixtures and why; every contract gap; divergences from
mcp-sso you found that have no delta row (the most important findings); what each
client did with `approval_pending` and where the user finds the `request_id`.
