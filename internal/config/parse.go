package config

import (
	"fmt"
	"regexp"
	"strings"
)

var headerNameRe = regexp.MustCompile(`^[A-Za-z0-9-]+$`)

// parseDocuments applies the B1 envelope and the per-kind field tables. It
// records every structural refusal it can find; cross-object rules run later.
func parseDocuments(docs []any, c *collector) *Config {
	cfg := &Config{}
	gateways := 0
	for i, d := range docs {
		res := fmt.Sprintf("document %d", i+1)
		o, ok := newObj(c, res, "", d)
		if !ok {
			continue
		}
		apiVersion, _ := o.str("apiVersion", true)
		kind, _ := o.str("kind", true)
		name := ""
		if meta, mok := o.object("metadata", true); mok {
			name, _ = meta.str("name", true)
			meta.finish()
		}
		if apiVersion != "" && apiVersion != APIVersion {
			c.addf(res, "apiVersion", "B1.api-version", "unknown apiVersion %q; only %s is accepted", apiVersion, APIVersion)
		}
		if kind != "" && kind != "Gateway" && kind != "Route" {
			c.addf(res, "kind", "B1.kind", "unknown kind %q; only Gateway and Route are accepted", kind)
		}
		if name != "" && !isIdentifier(name) {
			c.add(res, "metadata.name", "B3.identifier", "must match [a-z][a-z0-9-]{0,62}")
		}
		if (kind == "Gateway" || kind == "Route") && name != "" {
			res = kind + "/" + name
			o.resource = res
		}
		specV, sok := o.raw("spec")
		o.finish()
		if !sok {
			c.add(res, "spec", "B1.required", "required field is missing")
			continue
		}
		switch kind {
		case "Gateway":
			gateways++
			if gateways > 1 {
				c.add(res, "kind", "B1.envelope", "exactly one Gateway document is accepted")
				continue
			}
			cfg.Gateway = parseGateway(c, res, name, specV)
		case "Route":
			if r := parseRoute(c, res, name, specV); r != nil {
				cfg.Routes = append(cfg.Routes, r)
			}
		}
	}
	if gateways == 0 {
		c.add("config", "", "B1.envelope", "no Gateway document")
	}
	if len(cfg.Routes) == 0 {
		c.add("config", "", "B1.envelope", "at least one Route document is required")
	}
	return cfg
}

func refField(c *collector, res, field, raw string) (Ref, bool) {
	r, err := parseRef(raw)
	if err != nil {
		c.add(res, field, "B2.reference", err.Error())
		return Ref{}, false
	}
	return r, true
}

func parseGateway(c *collector, res, name string, specV any) *Gateway {
	s, ok := newObj(c, res, "spec", specV)
	if !ok {
		return nil
	}
	g := &Gateway{Name: name, Egress: map[string]EgressProfile{}}
	g.ExternalBaseURL, _ = s.str("externalBaseUrl", true)
	g.Listen, _ = s.str("listen", true)
	if raw, ok := s.str("signingKeyRef", true); ok {
		g.SigningKeyRef, _ = refField(c, res, "spec.signingKeyRef", raw)
	}

	if hosts, ok := s.strList("allowedHosts", true); ok {
		if len(hosts) == 0 {
			c.add(res, "spec.allowedHosts", "B1.empty-allowlist", "empty never means allow all")
		}
		for i, h := range hosts {
			if err := checkHostPort(h); err != nil {
				c.add(res, fmt.Sprintf("spec.allowedHosts[%d]", i), "B3.host", err.Error())
			}
		}
		g.AllowedHosts = hosts
	}
	if origins, ok := s.strList("allowedOrigins", true); ok {
		if len(origins) == 0 {
			c.add(res, "spec.allowedOrigins", "B1.empty-allowlist", "empty never means allow all")
		}
		for i, o := range origins {
			if _, err := checkURL(o, urlRules{schemes: []string{"https", "http"}, originOnly: true}); err != nil {
				c.add(res, fmt.Sprintf("spec.allowedOrigins[%d]", i), "B3.url", "origin: "+err.Error())
			}
		}
		g.AllowedOrigins = origins
	}
	g.TrustProxyHeaders = s.boolean("trustProxyHeaders", false)
	if proxies, ok := s.strList("trustedProxies", false); ok {
		for i, p := range proxies {
			if err := checkCIDROrIP(p); err != nil {
				c.add(res, fmt.Sprintf("spec.trustedProxies[%d]", i), "B6.trusted-proxy", err.Error())
			}
		}
		g.TrustedProxies = proxies
	}
	if g.TrustProxyHeaders && len(g.TrustedProxies) == 0 {
		c.add(res, "spec.trustProxyHeaders", "B6.trusted-proxies-required", "trustProxyHeaders: true requires a non-empty trustedProxies list")
	}

	g.Identity = parseIdentity(c, res, s)
	g.Clients = parseClients(c, res, s)

	if eg, ok := s.object("egress", false); ok {
		if profiles, ok := eg.object("profiles", false); ok {
			for pname := range profiles.m {
				pobj, pok := profiles.object(pname, true)
				field := "spec.egress.profiles." + pname
				if !isIdentifier(pname) {
					c.add(res, field, "B3.identifier", "profile name must match [a-z][a-z0-9-]{0,62}")
				}
				if !pok {
					continue
				}
				var ep EgressProfile
				if proxy, ok := pobj.str("proxy", true); ok {
					if proxy != "fromEnv" && proxy != "none" {
						if _, err := checkURL(proxy, urlRules{schemes: []string{"http", "https"}, allowPath: true}); err != nil {
							c.add(res, field+".proxy", "B3.url", "proxy must be fromEnv, none, or a canonical URL: "+err.Error())
						}
					}
					ep.Proxy = proxy
				}
				if raw, ok := pobj.str("caBundleRef", false); ok {
					if r, ok := refField(c, res, field+".caBundleRef", raw); ok {
						ep.CABundleRef = &r
					}
				}
				pobj.finish()
				g.Egress[pname] = ep
			}
			profiles.finish()
		}
		eg.finish()
	}

	if st, ok := s.object("store", true); ok {
		g.StorePath = absPath(c, res, st, "spec.store", "path")
		st.finish()
	}
	if au, ok := s.object("audit", true); ok {
		g.AuditPath = absPath(c, res, au, "spec.audit", "path")
		au.finish()
	}

	g.Health = Health{LivePath: defaultLivePath, ReadyPath: defaultReadyPath}
	if h, ok := s.object("health", false); ok {
		if p, ok := h.str("livePath", false); ok {
			g.Health.LivePath = pathField(c, res, "spec.health.livePath", p)
		}
		if p, ok := h.str("readyPath", false); ok {
			g.Health.ReadyPath = pathField(c, res, "spec.health.readyPath", p)
		}
		h.finish()
	}

	g.Tokens = Tokens{AccessTTL: defaultAccessTTL, ConsentTTL: defaultConsentTTL, CodeTTL: defaultCodeTTL}
	if t, ok := s.object("tokens", false); ok {
		if n, ok := t.positive("accessTtl", false); ok {
			g.Tokens.AccessTTL = n
		}
		if n, ok := t.positive("consentTtl", false); ok {
			g.Tokens.ConsentTTL = n
		}
		if n, ok := t.positive("codeTtl", false); ok {
			g.Tokens.CodeTTL = n
		}
		t.finish()
	}

	if list, ok := s.list("machineClients", false); ok {
		for i, v := range list {
			field := fmt.Sprintf("spec.machineClients[%d]", i)
			m, ok := newObj(c, res, field, v)
			if !ok {
				continue
			}
			var mc MachineClient
			if id, ok := m.str("id", true); ok {
				if !isIdentifier(id) {
					c.add(res, field+".id", "B3.identifier", "must match [a-z][a-z0-9-]{0,62}")
				}
				mc.ID = id
			}
			if raw, ok := m.str("secretRef", true); ok {
				mc.SecretRef, _ = refField(c, res, field+".secretRef", raw)
			}
			if p, ok := m.str("purpose", true); ok {
				if err := checkPurpose(p); err != nil {
					c.add(res, field+".purpose", "B5.purpose", err.Error())
				}
				mc.Purpose = p
			}
			if d, ok := m.positive("maxDuration", true); ok {
				if d > hardMaxDuration {
					c.addf(res, field+".maxDuration", "B8.max-duration", "exceeds the hard ceiling of %d seconds", hardMaxDuration)
				}
				mc.MaxDuration = d
			}
			if routes, ok := m.list("routes", true); ok {
				if len(routes) == 0 {
					c.add(res, field+".routes", "B1.required", "at least one route declaration is required")
				}
				for j, rv := range routes {
					rf := fmt.Sprintf("%s.routes[%d]", field, j)
					ro, ok := newObj(c, res, rf, rv)
					if !ok {
						continue
					}
					var mr MachineRoute
					if p, ok := ro.str("path", true); ok {
						mr.Path = pathField(c, res, rf+".path", p)
					}
					if scopes, ok := ro.strList("scopes", true); ok {
						if len(scopes) == 0 {
							c.add(res, rf+".scopes", "B1.required", "at least one scope is required")
						}
						for k, sc := range scopes {
							if err := checkScope(sc); err != nil {
								c.add(res, fmt.Sprintf("%s.scopes[%d]", rf, k), "B1.scope", err.Error())
							}
						}
						mr.Scopes = scopes
					}
					ro.finish()
					mc.Routes = append(mc.Routes, mr)
				}
			}
			m.finish()
			g.MachineClients = append(g.MachineClients, mc)
		}
	}
	s.finish()
	return g
}

func absPath(c *collector, res string, o *obj, prefix, name string) string {
	p, ok := o.str(name, true)
	if !ok {
		return ""
	}
	if !strings.HasPrefix(p, "/") {
		c.add(res, prefix+"."+name, "B2.path", "must be an absolute path")
	}
	return p
}

func pathField(c *collector, res, field, p string) string {
	if err := checkPath(p); err != nil {
		c.add(res, field, "B3.path", err.Error())
	}
	return p
}

func parseIdentity(c *collector, res string, s *obj) Identity {
	var id Identity
	o, ok := s.object("identity", true)
	if !ok {
		return id
	}
	f := func(name string) string { return "spec.identity." + name }
	provider, _ := o.str("provider", true)
	switch provider {
	case "entra", "oidc", "header", "console":
	case "":
	default:
		c.addf(res, f("provider"), "B1.identity-provider", "unknown provider %q; entra, oidc, header, or console", provider)
	}
	id.Provider = provider
	idp := provider == "entra" || provider == "oidc"

	if idp {
		id.Registration = "dedicated"
		if r, ok := o.str("registration", false); ok {
			if r != "dedicated" && r != "shared" {
				c.add(res, f("registration"), "B1.registration", "dedicated or shared")
			}
			id.Registration = r
		}
		if u, ok := o.str("redirectUri", true); ok {
			if _, err := checkURL(u, urlRules{schemes: []string{"https"}, allowPath: true}); err != nil {
				c.add(res, f("redirectUri"), "B3.url", err.Error())
			}
			id.RedirectURI = u
		}
		id.CallbackPath = defaultCallbackPath
		if p, ok := o.str("callbackPath", false); ok {
			id.CallbackPath = pathField(c, res, f("callbackPath"), p)
		}
		id.ClientID, _ = o.str("clientId", true)
		secret, hasSecret := o.str("clientSecretRef", false)
		pub, pubOK := o.literalTrue("publicClient")
		if hasSecret {
			if r, ok := refField(c, res, f("clientSecretRef"), secret); ok {
				id.ClientSecretRef = &r
			}
		}
		id.PublicClient = pub && pubOK
		if hasSecret == (pub && pubOK) {
			c.add(res, f("clientSecretRef"), "B1.client-auth", "exactly one of clientSecretRef or publicClient: true")
		}
		if id.Registration == "shared" && !hasSecret {
			c.add(res, f("registration"), "B1.shared-secret", "a shared app registration requires its own client secret")
		}
		if provider == "entra" {
			id.TenantID, _ = o.str("tenantId", true)
			o.forbid("issuer", "issuer belongs to the oidc variant")
		} else {
			o.forbid("tenantId", "tenantId belongs to the entra variant")
			if iss, ok := o.str("issuer", true); ok {
				if _, err := checkURL(iss, urlRules{schemes: []string{"https"}, allowPath: true}); err != nil {
					c.add(res, f("issuer"), "B3.url", err.Error())
				}
				id.Issuer = iss
			}
		}
		id.GroupsClaim = defaultGroupsClaim
		if g, ok := o.str("groupsClaim", false); ok {
			id.GroupsClaim = g
		}
		id.EgressProfile, _ = o.str("egressProfile", true)
		o.forbid("assertion", "assertion belongs to the header variant")
		o.forbid("pairingCodeTtl", "pairingCodeTtl belongs to the console variant")
	} else {
		for _, name := range []string{"registration", "redirectUri", "callbackPath", "clientId", "clientSecretRef", "publicClient", "tenantId", "issuer"} {
			o.forbid(name, name+" belongs to the entra/oidc variants")
		}
		o.forbid("groupsClaim", "refused for this variant; header mode selects groups with assertion.groupsClaim")
		o.forbid("egressProfile", "refused for this variant; header mode selects JWKS egress with assertion.keys.egressProfile")
		if provider == "header" {
			id.Assertion = parseAssertion(c, res, o)
			o.forbid("pairingCodeTtl", "pairingCodeTtl belongs to the console variant")
		}
		if provider == "console" {
			o.forbid("assertion", "assertion belongs to the header variant")
			id.PairingCodeTTL = defaultPairingCodeTTL
			if n, ok := o.positive("pairingCodeTtl", false); ok {
				id.PairingCodeTTL = n
			}
		}
	}

	if gts, ok := o.object("groupsToScopes", false); ok {
		id.GroupsToScopes = map[string][]string{}
		for group := range gts.m {
			scopes, ok := gts.strList(group, true)
			if strings.TrimSpace(group) == "" {
				c.add(res, f("groupsToScopes"), "B1.blank", "blank group name")
			}
			if !ok {
				continue
			}
			for i, sc := range scopes {
				if err := checkScope(sc); err != nil {
					c.add(res, fmt.Sprintf("%s.%s[%d]", f("groupsToScopes"), group, i), "B1.scope", err.Error())
				}
			}
			id.GroupsToScopes[group] = scopes
		}
		gts.finish()
	}
	o.finish()
	return id
}

func parseAssertion(c *collector, res string, id *obj) *Assertion {
	a, ok := id.object("assertion", true)
	if !ok {
		return nil
	}
	f := func(name string) string { return "spec.identity.assertion." + name }
	var as Assertion
	if h, ok := a.str("header", true); ok {
		if !headerNameRe.MatchString(h) {
			c.add(res, f("header"), "B4.header-name", "header name must match [A-Za-z0-9-]+")
		}
		as.Header = h
	}
	as.Issuer, _ = a.str("issuer", true)
	as.Audience, _ = a.str("audience", true)
	if alg, ok := a.str("alg", true); ok {
		if alg != "ES256" && alg != "RS256" {
			c.add(res, f("alg"), "B4.alg", "alg is pinned to ES256 or RS256")
		}
		as.Alg = alg
	}
	if keys, ok := a.object("keys", true); ok {
		jwksURL, hasURL := keys.str("jwksUrl", false)
		egress, hasEgress := keys.str("egressProfile", false)
		jwksRef, hasRef := keys.str("jwksRef", false)
		if hasURL {
			if _, err := checkURL(jwksURL, urlRules{schemes: []string{"https"}, allowPath: true}); err != nil {
				c.add(res, f("keys.jwksUrl"), "B3.url", err.Error())
			}
			as.Keys.JWKSURL = jwksURL
		}
		as.Keys.EgressProfile = egress
		if hasRef {
			if r, ok := refField(c, res, f("keys.jwksRef"), jwksRef); ok {
				if r.Kind != "file" {
					c.add(res, f("keys.jwksRef"), "B4.keys-union", "jwksRef must be a file: reference")
				}
				as.Keys.JWKSRef = &r
			}
		}
		fetched := hasURL && hasEgress && !hasRef
		static := hasRef && !hasURL && !hasEgress
		if !fetched && !static {
			c.add(res, f("keys"), "B4.keys-union", "keys is exactly {jwksUrl, egressProfile} or {jwksRef}")
		}
		keys.finish()
	}
	as.SubjectClaim = "sub"
	if s, ok := a.str("subjectClaim", false); ok {
		as.SubjectClaim = s
	}
	as.GroupsClaim, _ = a.str("groupsClaim", false)
	a.finish()
	return &as
}

func parseClients(c *collector, res string, s *obj) Clients {
	cl := Clients{AllowLoopbackRedirects: true}
	o, ok := s.object("clients", false)
	if !ok {
		return cl
	}
	if paths, ok := o.strList("knownCimd", false); ok {
		for i, p := range paths {
			if !strings.HasPrefix(p, "/") {
				c.add(res, fmt.Sprintf("spec.clients.knownCimd[%d]", i), "B2.path", "must be an absolute path")
			}
		}
		cl.KnownCIMD = paths
	}
	cl.AllowLoopbackRedirects = o.boolean("allowLoopbackRedirects", true)
	if list, ok := o.strList("redirectAllowlist", false); ok {
		for i, e := range list {
			if err := checkOriginOrURL(e); err != nil {
				c.add(res, fmt.Sprintf("spec.clients.redirectAllowlist[%d]", i), "B3.url", err.Error())
			}
		}
		cl.RedirectAllowlist = list
	}
	o.finish()
	return cl
}

func parseRoute(c *collector, res, name string, specV any) *Route {
	s, ok := newObj(c, res, "spec", specV)
	if !ok {
		return nil
	}
	r := &Route{Name: name, MCPEndpoint: defaultMCPEndpoint, MaxDuration: defaultMaxDuration}
	if p, ok := s.str("path", true); ok {
		r.Path = pathField(c, res, "spec.path", p)
	}
	if p, ok := s.str("mcpEndpoint", false); ok {
		r.MCPEndpoint = pathField(c, res, "spec.mcpEndpoint", p)
	}

	if up, ok := s.object("upstream", true); ok {
		if u, ok := up.str("url", true); ok {
			if _, err := checkURL(u, urlRules{schemes: []string{"http", "https"}, allowPath: true}); err != nil {
				c.add(res, "spec.upstream.url", "B3.url", err.Error())
			}
			r.Upstream.URL = u
		}
		r.Upstream.EgressProfile, _ = up.str("egressProfile", true)
		if cr, ok := up.object("credential", true); ok {
			t, _ := cr.str("type", true)
			r.Upstream.Credential.Type = t
			switch t {
			case "none":
				for _, n := range []string{"header", "scheme", "valueRef"} {
					cr.forbid(n, n+" belongs to the static-header credential")
				}
			case "static-header":
				if h, ok := cr.str("header", true); ok {
					if !headerNameRe.MatchString(h) {
						c.add(res, "spec.upstream.credential.header", "B1.header-name", "header name must match [A-Za-z0-9-]+")
					}
					r.Upstream.Credential.Header = h
				}
				r.Upstream.Credential.Scheme, _ = cr.str("scheme", false)
				if raw, ok := cr.str("valueRef", true); ok {
					if ref, ok := refField(c, res, "spec.upstream.credential.valueRef", raw); ok {
						r.Upstream.Credential.ValueRef = &ref
					}
				}
			case "":
			default:
				c.addf(res, "spec.upstream.credential.type", "B1.credential-type", "unknown credential type %q; none or static-header", t)
			}
			cr.finish()
		}
		up.finish()
	}

	if sc, ok := s.object("scopes", true); ok {
		if catalog, ok := sc.strList("catalog", true); ok {
			if len(catalog) == 0 {
				c.add(res, "spec.scopes.catalog", "B1.required", "catalog must not be empty")
			}
			if len(catalog) > scopeMaxEntries {
				c.addf(res, "spec.scopes.catalog", "B1.scope", "more than %d scopes", scopeMaxEntries)
			}
			seen := map[string]bool{}
			for i, s := range catalog {
				if err := checkScope(s); err != nil {
					c.add(res, fmt.Sprintf("spec.scopes.catalog[%d]", i), "B1.scope", err.Error())
				}
				if seen[s] {
					c.add(res, fmt.Sprintf("spec.scopes.catalog[%d]", i), "B1.scope", "duplicate scope")
				}
				seen[s] = true
			}
			r.Catalog = catalog
		}
		if def, ok := sc.strList("default", false); ok {
			for i, s := range def {
				if !contains(r.Catalog, s) {
					c.add(res, fmt.Sprintf("spec.scopes.default[%d]", i), "B1.scope", "default scope is not in the catalog")
				}
			}
			r.DefaultScopes = def
		}
		sc.finish()
	}

	if az, ok := s.object("authz", true); ok {
		if rs, ok := az.str("requireScope", true); ok {
			if !contains(r.Catalog, rs) {
				c.add(res, "spec.authz.requireScope", "B1.scope", "requireScope is not in the catalog")
			}
			r.RequireScope = rs
		}
		az.finish()
	}

	if gr, ok := s.object("grant", false); ok {
		if d, ok := gr.positive("maxDuration", false); ok {
			if d > hardMaxDuration {
				c.addf(res, "spec.grant.maxDuration", "B8.max-duration", "exceeds the hard ceiling of %d seconds", hardMaxDuration)
			}
			r.MaxDuration = d
		}
		if pol, ok := gr.object("policy", false); ok {
			if rules, ok := pol.list("rules", false); ok {
				for i, rv := range rules {
					field := fmt.Sprintf("spec.grant.policy.rules[%d]", i)
					ro, ok := newObj(c, res, field, rv)
					if !ok {
						continue
					}
					var rule PolicyRule
					if e, ok := ro.str("effect", true); ok {
						if e != "allow" && e != "deny" {
							c.add(res, field+".effect", "G7.effect", "allow or deny")
						}
						rule.Effect = e
					}
					if w, ok := ro.object("when", true); ok {
						if len(w.m) == 0 {
							c.add(res, field+".when", "G7.when-empty", "when must name at least one condition")
						}
						if sc, ok := w.strList("scopesSubsetOf", false); ok {
							for k, s := range sc {
								if err := checkScope(s); err != nil {
									c.add(res, fmt.Sprintf("%s.when.scopesSubsetOf[%d]", field, k), "B1.scope", err.Error())
								}
							}
							rule.ScopesSubsetOf = sc
						}
						if d, ok := w.positive("durationAtMost", false); ok {
							rule.DurationAtMost = &d
						}
						rule.SubjectInGroup, _ = w.str("subjectInGroup", false)
						if cl, ok := w.strList("clientIn", false); ok {
							if len(cl) == 0 {
								c.add(res, field+".when.clientIn", "B1.required", "clientIn must not be empty")
							}
							rule.ClientIn = cl
						}
						w.finish()
					}
					ro.finish()
					r.Rules = append(r.Rules, rule)
				}
			}
			pol.finish()
		}
		if ap, ok := gr.strList("approvers", false); ok {
			r.Approvers = ap
		}
		gr.finish()
	}
	s.finish()

	r.MCPPath = r.Path + r.MCPEndpoint
	r.MetadataURL = prmPrefix + r.MCPPath
	return r
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}
