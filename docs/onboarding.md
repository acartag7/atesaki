# Onboarding

**Status: DRAFT.**

## The operator (the person who runs the front door)

0. *Five-minute first run, no IdP:* `identity.provider: console` on loopback. The
   server prints a one-time code; paste it in the browser; you are the
   `console-operator`. It refuses to run on any non-loopback address — it is a
   tutorial, not a deployment. Then continue below for real users.
1. Write `atesaki.yaml`: one `Gateway`, one `Route` pointing at the MCP server that
   today needs a shared key. Secrets as `env:`/`file:` references.
2. `atesaki validate` — fix what it names. Nothing has touched the network yet.
3. `atesaki idp-request` — paste the output into the IdP team's ticket queue. It asks
   for the minimum: one app, one redirect URI, groups on the id_token, one secret.
   It also lists what it does *not* need, so the reviewer can say yes quickly.
4. Before anything deploys: `atesaki rehearse` — the whole flow runs on your machine
   against a **mock** IdP, per client you care about. It proves the gateway, the
   clients, and your config agree. It cannot prove the company IdP registration.
5. Deploy the container. `atesaki routes` gives the ingress the exact path list.
6. **The proof step:** one real sign-in against the real IdP, from one real client.
   Only this proves the ticket landed as asked. Do it before announcing the route.
7. Add backend N+1 by adding a `Route` block. No new hostname, no new IdP ticket.
8. Decide per route what may auto-approve and what is refused outright; everything
   else waits for a human. `atesaki grants pending` shows requests waiting;
   `grants approve <id>` (you may shorten the duration or drop scopes, never add) or
   `grants deny <id>`. The requester then runs their flow again: they see the values
   *you* approved on their consent page, approve them, and it completes.
9. Unattended agents (cron, CI, a server-side agent) don't get a browser flow: declare
   them as machine clients in the config — which routes, which scopes, why, and for how
   long at most. Their access is a grant like any other: listed, bounded, revocable.

If the IdP ticket is impossible: pick a rung (`docs/contract.md §4`) explicitly. The
degraded modes are supported and named — the product never picks one for you.

## The user (the person whose agent connects)

1. Add the route URL (`https://host/playbook/mcp`) to the agent.
2. Sign in with the normal company login when the browser opens.
3. Your agent states **why** it needs access and **for how long** (the route caps the
   duration). You approve or refuse what it asked — the page shows exactly the client,
   the route, the scopes (already capped by your groups), the stated purpose, and the
   duration; the clock starts when the agent exchanges its approval for tokens (not at the first tool call, and not when you approved). Nothing shown can
   change after you approve.
4. Sometimes the answer is "this needs approval": the flow ends, you're told a
   request id, and an approver decides — approve (possibly with less than you asked
   for), or refuse. Once approved, run the same request again: the consent page shows
   what was actually approved, you approve that, and it completes without a second
   wait. Some requests are refused outright by the route's rules; the answer then is
   simply no.
5. Work. Access ends at the expiry you chose — or within a few minutes of the
   operator revoking the grant (the exact maximum is stated by the deployment). After
   that the agent asks again, with a fresh purpose. Re-login behavior after a gateway
   restart is likewise stated by the deployment, not discovered.

## What onboarding never asks

- No keys pasted into any client config.
- No per-user IdP tickets.
- No pre-provisioning. A person with no grant simply has no token yet; the first
  request to a route answers with the standard challenge that starts sign-in — nothing
  about them has to exist in Atesaki beforehand.
