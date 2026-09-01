# Atesaki

**Status: DRAFT contract. Nothing is frozen.** Product code is being built slice by
slice against the contract as it stands; each PR names the sections it implements.

```
go build ./cmd/atesaki
atesaki validate atesaki.yaml   # pure validation: reads config, touches nothing else
atesaki routes   atesaki.yaml   # the route and well-known path list, as JSON
```

A team has a useful internal service an AI agent could use, but the only way in is a
shared key nobody wants to hand out. The operator runs one container in front of it and
writes one YAML file. From then on: a person points their agent at a URL, signs in with
their normal company login, approves what the agent may do, and works. The key never
leaves the container. The person who couldn't get the security team to change anything
gets a generated one-page request saying exactly what to ask for — and a dress-rehearsal
command that proves the login works before anything is deployed.

Positioning: *agentgateway is the gateway for companies whose IdP and network cooperate;
Atesaki is the front door for the ones where they don't* — change-controlled IdPs,
proxied egress, internal CAs, no staging environment.

## Decisions on record (2026-08-30, Arnold)

1. **The hosted front door is core Atesaki.** This reverses the 2026-08-16 roadmap's
   local-first order; the local daemon becomes a later chapter (`docs/future.md`).
2. **New repo, Go.** mcp-sso's rigor without its complications: Atesaki takes mcp-sso's
   *contract clauses and parity corpus*, not its code. The Go authorization-server half
   must pass every frozen fixture **labeled portable** with zero skips (mcp-sso
   `docs/contracts/19`); fixtures labeled *host* bind only the TypeScript reference.
   Intentional behavioral differences are few, public, and listed in `docs/deltas.md`
   — nothing differs silently. Atesaki pins an exact corpus version + manifest hash
   as a test-only input. `[O:2026-08-30]`
3. **Credentials are dispensed, not configured.** Core, not a feature: every grant
   records who, why (stated purpose — recorded, bounded, never trusted), and for how
   long (expiry bounds the whole token family; revocable; `atesaki grants`). See
   `docs/contract-grants.md`. The run-bound lease (the full broker) stays future.
   Every later decision (2026-08-31: policy step, machine clients, store, signing order,
   audit classes, rung 4, numbers) is in `docs/decisions.md` with its receipt.
4. Evidence base: `~/project/atesaki/evidence/migration-playbook-mcp-setup-pain-2026-08-26.md`
   (F1–F9, R1–R11), the splunk-gateway spike, and the starter kit (imported at
   `seed/` — reference YAML, schema, 52-case conformance checklist).

## Documents

| File | Job |
| --- | --- |
| `docs/quality-bar.md` | how we work; when a rule may change |
| `docs/contract.md` | roles, config rules, tokens, ladder, routes, relay, egress, verbs, the nevers |
| `docs/contract-grants.md` | GrantV0: dispensing, the records and state machines, hashes, the operation table, policy, machine grants |
| `docs/contract-boundaries.md` | typed config schema, reference trust, canonicalization, signed-assertion rung, caps, proxy trust, error catalog |
| `docs/deltas.md` | the public list of intentional differences from the mcp-sso reference |
| `docs/onboarding.md` | how an operator and a user join |
| `docs/threat-model.md` | what a signed-in attacker can try (seed) |
| `docs/open-questions.md` | only what is still open — nothing here is decided |
| `docs/future.md` | not this version |
| `docs/decisions.md` | the ledger — every `[O]` tag has a row saying how Arnold decided it |
| `tools/contract-lint.py` | coverage lint: dangling refs, reasons without producers, states without operations, decisions without receipts |
| `prompts/` | dispatch packets for every remaining piece of work, in execution order, with the gates between them (`prompts/README.md`) |

There is no `http-api.md`: Atesaki has no website. Its HTTP surface (OAuth endpoints,
well-known documents, `/​<route>/mcp`) is part of `docs/contract.md`.

## Relationship to mcp-sso

- mcp-sso (TypeScript) is the **reference implementation** of the authorization-server
  contract set. Atesaki's Go AS is an **implementation** under mcp-sso §19: it proves
  itself against the shared fixture corpus, clause by clause.
- Where a rule already exists as an mcp-sso contract clause, this repo **cites it**
  (`mcp-sso §NN.C`) instead of rewording it. One rule lives in one place.
- Standing ask to mcp-sso (work item, not a question): grow the corpus beyond §08 —
  §07 crypto/token, §09 AS-lite bridge, §10 redirect policy, §11 scopes — because those
  are the clauses the Go implementation must prove before it can claim parity.

## Sequence (after freeze — not before)

contract → threat model → acceptance tests (corpus-driven) → freeze → implementation
prompts handed to build sessions (`prompts/`, in order, gates stated). This repo's
design chat writes rules, tests, and prompts only.
