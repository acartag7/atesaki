# atesaki — design starter kit

A multi-path MCP OAuth gateway: one host, one Authorization Server, N routes.
Each route is an MCP resource with its own upstream, its own injected credential,
and — the load-bearing invariant — **its own audience**. The access token's `aud`
is the security boundary between paths.

These three files are a starting-point spec, distilled from operating an
mcp-sso-based gateway into a real enterprise cluster. Most requirements here map
to a failure that was actually observed, not a hypothetical.

## Files

| File | What it is |
|------|------------|
| `atesaki.reference.yaml` | Annotated reference config — the route model by example. Every field is commented with the failure mode it prevents. Two sample routes (`/splunk`, `/conduktor`) share one IdP but keep distinct audiences and scope catalogs. |
| `atesaki.schema.json` | JSON Schema (draft 2020-12) to seed `atesaki validate`. Strict where strictness fails closed: https-only issuer/resource, secret-refs only (inline secrets forbidden), non-`*` redirect allowlist, `credential` as a `oneOf` on `type`. |
| `CONFORMANCE.md` | 52 test cases in 11 groups (A–K), each annotated with the failure mode it guards against. |

## The model in one breath

```
externalBaseUrl ─ one Authorization Server (issuer) ─┬─ route /splunk    aud=…/splunk/mcp    → upstream A
                                                     └─ route /conduktor aud=…/conduktor/mcp → upstream B
```

- **Discovery is per route** (RFC 9728 path-insertion); a route's `WWW-Authenticate`
  points at its own protected-resource metadata.
- **Tokens are per-audience** (RFC 8707 resource indicators); a `/splunk` token is
  rejected at `/conduktor`. That test is `CONFORMANCE.md` §D — the one that turns
  multi-path from a feature into a vuln if it regresses.
- **Egress and CA trust are per-destination.** The IdP is usually reached through a
  forward proxy; internal upstreams direct. Node's `fetch` won't pick a proxy from
  the environment, so it's explicit config.
- **Credential injection is a strategy**: `static-header` | `token-exchange` (RFC
  8693, so the backend sees *who*) | `per-identity`.

## Validate the reference against the schema

```bash
# any draft 2020-12 validator; example with python
pip install jsonschema pyyaml
python - <<'PY'
import json, yaml, jsonschema
schema = json.load(open("atesaki.schema.json"))
cfg    = yaml.safe_load(open("atesaki.reference.yaml"))
jsonschema.Draft202012Validator.check_schema(schema)
errs = list(jsonschema.Draft202012Validator(schema).iter_errors(cfg))
print("clean" if not errs else f"{len(errs)} errors")
PY
```

## Two decisions left open for you

1. **Secret-ref scheme.** Modeled as `scheme:location` (`k8s:secret/KEY`,
   `file:/path`, `env:VAR`, `vault:…`), with the scheme allowlisted in the schema
   pattern. Swap the enum for your resolver set.
2. **Shared vs. per-route AS.** Modeled as one shared AS (single issuer, per-path
   audiences). If you go per-route instead, the `authorizationServer` block moves
   inside `route` and §D's cross-path replay test changes shape.
