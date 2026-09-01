package config

import (
	"fmt"
	"strings"
)

// crossValidate applies the rules that need more than one object: the
// gateway URL against the identity variant, reserved paths, route collisions,
// audiences, profile references, allowlists, group ceilings, machine
// declarations, the G7 contradiction check, and approver requiredness.
func crossValidate(cfg *Config, c *collector) {
	g := cfg.Gateway
	if g == nil {
		return
	}
	res := "Gateway/" + g.Name
	console := g.Identity.Provider == "console"

	// externalBaseUrl: https origin; console may use http on loopback and must be loopback.
	schemes := []string{"https"}
	if console {
		schemes = []string{"https", "http"}
	}
	if g.ExternalBaseURL != "" {
		u, err := checkURL(g.ExternalBaseURL, urlRules{schemes: schemes, originOnly: true})
		if err != nil {
			c.add(res, "spec.externalBaseUrl", "B3.url", err.Error())
		} else if console && !isLoopbackHost(bracket(u.Hostname())) {
			c.add(res, "spec.externalBaseUrl", "B1.console-loopback", "console pairing is loopback-only; refusing a non-loopback externalBaseUrl")
		}
	}
	if g.Listen != "" {
		host, port, err := checkListen(g.Listen)
		if err != nil {
			c.add(res, "spec.listen", "B3.host", err.Error())
		} else {
			g.ListenHost, g.ListenPort = host, port
			if console && !isLoopbackHost(host) {
				c.add(res, "spec.listen", "B1.console-loopback", "console pairing is loopback-only; listen must be a loopback address")
			}
		}
	}

	// Identity cross rules.
	id := &g.Identity
	if id.Provider == "entra" || id.Provider == "oidc" {
		if id.RedirectURI != "" && g.ExternalBaseURL != "" && id.RedirectURI != g.ExternalBaseURL+id.CallbackPath {
			c.addf(res, "spec.identity.redirectUri", "B1.redirect-uri", "must equal externalBaseUrl + callbackPath byte-exact (%s)", g.ExternalBaseURL+id.CallbackPath)
		}
		if id.EgressProfile != "" {
			if _, ok := g.Egress[id.EgressProfile]; !ok {
				c.addf(res, "spec.identity.egressProfile", "B1.egress-profile-unknown", "no egress profile named %q", id.EgressProfile)
			}
		}
	}
	if id.Provider == "header" && id.Assertion != nil && id.Assertion.Keys.EgressProfile != "" {
		if _, ok := g.Egress[id.Assertion.Keys.EgressProfile]; !ok {
			c.addf(res, "spec.identity.assertion.keys.egressProfile", "B1.egress-profile-unknown", "no egress profile named %q", id.Assertion.Keys.EgressProfile)
		}
	}
	if !console && len(g.Clients.RedirectAllowlist) == 0 {
		c.add(res, "spec.clients.redirectAllowlist", "B1.redirect-allowlist", "required in every hosted mode; empty never means allow all")
	}

	// Reserved paths (B3).
	reserved := []string{"/oauth", "/.well-known", g.Health.LivePath, g.Health.ReadyPath}
	if id.CallbackPath != "" {
		reserved = append(reserved, id.CallbackPath)
	}
	if g.Health.LivePath == g.Health.ReadyPath {
		c.add(res, "spec.health.readyPath", "B3.reserved-path", "livePath and readyPath must differ")
	}

	// Routes.
	unionCatalog := map[string]bool{}
	byPath := map[string]*Route{}
	audiences := map[string]string{}
	names := map[string]bool{}
	for i, r := range cfg.Routes {
		rres := "Route/" + r.Name
		if r.Name != "" {
			if names[r.Name] {
				c.add(rres, "metadata.name", "B1.duplicate-key", "duplicate Route name")
			}
			names[r.Name] = true
		}
		if r.Path == "" {
			continue
		}
		for _, sc := range r.Catalog {
			unionCatalog[sc] = true
		}
		for _, p := range reserved {
			if p == "" {
				continue
			}
			if under(r.Path, p) || under(p, r.Path) {
				c.addf(rres, "spec.path", "B3.reserved-path", "collides with reserved path %s", p)
			}
		}
		for j := 0; j < i; j++ {
			o := cfg.Routes[j]
			if o.Path == "" {
				continue
			}
			if under(r.Path, o.Path) || under(o.Path, r.Path) {
				c.addf(rres, "spec.path", "B3.route-collision", "collides with Route/%s path %s", o.Name, o.Path)
			}
		}
		byPath[r.Path] = r
		if g.ExternalBaseURL != "" {
			r.Audience = g.ExternalBaseURL + r.MCPPath
			if other, dup := audiences[r.Audience]; dup {
				c.addf(rres, "spec.mcpEndpoint", "B1.audience-duplicate", "same audience as Route/%s: %s", other, r.Audience)
			}
			audiences[r.Audience] = r.Name
		}
		if r.Upstream.EgressProfile != "" {
			if _, ok := g.Egress[r.Upstream.EgressProfile]; !ok {
				c.addf(rres, "spec.upstream.egressProfile", "B1.egress-profile-unknown", "no egress profile named %q", r.Upstream.EgressProfile)
			}
		}
		if len(r.Approvers) == 0 && canEscalate(r) {
			c.add(rres, "spec.grant.approvers", "B1.approvers-required", "some request can escalate under these rules, so approvers are required")
		}
	}

	// Group ceiling scopes must exist in some route catalog.
	for group, scopes := range id.GroupsToScopes {
		for i, sc := range scopes {
			if !unionCatalog[sc] {
				c.addf(res, fmt.Sprintf("spec.identity.groupsToScopes.%s[%d]", group, i), "B1.groups-to-scopes", "scope %q is in no route catalog", sc)
			}
		}
	}

	// Machine clients (G10) and the G7 contradiction check.
	seenIDs := map[string]bool{}
	for i, mc := range g.MachineClients {
		field := fmt.Sprintf("spec.machineClients[%d]", i)
		if mc.ID != "" {
			if seenIDs[mc.ID] {
				c.add(res, field+".id", "B1.duplicate-key", "duplicate machine client id")
			}
			seenIDs[mc.ID] = true
		}
		for j, mr := range mc.Routes {
			rf := fmt.Sprintf("%s.routes[%d]", field, j)
			r, ok := byPath[mr.Path]
			if !ok {
				c.addf(res, rf+".path", "B1.machine-route-unknown", "no Route with path %s", mr.Path)
				continue
			}
			for k, sc := range mr.Scopes {
				if !contains(r.Catalog, sc) {
					c.addf(res, fmt.Sprintf("%s.scopes[%d]", rf, k), "B1.machine-scope", "scope %q is not in the catalog of Route/%s", sc, r.Name)
				}
			}
			if mc.MaxDuration > r.MaxDuration {
				c.addf(res, field+".maxDuration", "B1.machine-duration", "exceeds grant.maxDuration of Route/%s (%d)", r.Name, r.MaxDuration)
			}
			if n, denied := whollyDenied(r, mc, mr); denied {
				c.addf(res, rf, "G7.declaration-denied", "declaration %q on %s is denied by rule %d of Route/%s", mc.ID, mr.Path, n, r.Name)
			}
		}
	}
}

func bracket(host string) string {
	if strings.Contains(host, ":") {
		return "[" + host + "]"
	}
	return host
}

// under reports whether p equals prefix or sits below it at a segment boundary.
func under(p, prefix string) bool {
	return p == prefix || strings.HasPrefix(p, prefix+"/")
}

// canEscalate is true unless some rule matches every possible request on the
// route: no group or client condition, scopes condition covering the whole
// catalog, duration condition at or above the route maximum. With escalate-by-
// default (G7) such a rule is the only way a request is provably allow/deny.
func canEscalate(r *Route) bool {
	for _, rule := range r.Rules {
		if rule.SubjectInGroup != "" || len(rule.ClientIn) > 0 {
			continue
		}
		if rule.ScopesSubsetOf != nil && !subset(r.Catalog, rule.ScopesSubsetOf) {
			continue
		}
		if rule.DurationAtMost != nil && *rule.DurationAtMost < r.MaxDuration {
			continue
		}
		return false
	}
	return true
}

// whollyDenied reports the first deny rule that matches every request a
// machine declaration can make (G7 boot contradiction check). Machines have
// no groups, so a subjectInGroup condition never matches; scopes are any
// subset of the declaration; duration is the declaration's maxDuration.
func whollyDenied(r *Route, mc MachineClient, mr MachineRoute) (int, bool) {
	for i, rule := range r.Rules {
		if rule.Effect != "deny" || rule.SubjectInGroup != "" {
			continue
		}
		if rule.ScopesSubsetOf != nil && !subset(mr.Scopes, rule.ScopesSubsetOf) {
			continue
		}
		if rule.DurationAtMost != nil && mc.MaxDuration > *rule.DurationAtMost {
			continue
		}
		if len(rule.ClientIn) > 0 && !contains(rule.ClientIn, mc.ID) {
			continue
		}
		return i, true
	}
	return 0, false
}

func subset(a, of []string) bool {
	for _, x := range a {
		if !contains(of, x) {
			return false
		}
	}
	return true
}
