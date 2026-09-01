package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// B3 grammars. One parser, one canonical form: config values are refused when
// not already canonical, never normalized.

var (
	identifierRe = regexp.MustCompile(`^[a-z][a-z0-9-]{0,62}$`)
	dnsLabelRe   = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`)
	ipv4Re       = regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])(\.(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])){3}$`)
	pathSegRe    = regexp.MustCompile(`^[A-Za-z0-9._~-]+$`)
	envNameRe    = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)
	scopeRe      = regexp.MustCompile(`^[\x21\x23-\x5B\x5D-\x7E]+$`)
)

func isIdentifier(s string) bool { return identifierRe.MatchString(s) }

// checkHost applies the portable host grammar: lowercase RFC 1123 DNS name
// without trailing dot, canonical IPv4, or RFC 5952 IPv6 in brackets.
func checkHost(h string) error {
	if h == "" {
		return errors.New("empty host")
	}
	if strings.HasPrefix(h, "[") {
		if !strings.HasSuffix(h, "]") {
			return errors.New("unterminated IPv6 literal")
		}
		inner := h[1 : len(h)-1]
		ip := net.ParseIP(inner)
		if ip == nil || ip.To4() != nil {
			return errors.New("not an IPv6 address")
		}
		if ip.String() != inner {
			return fmt.Errorf("IPv6 must be written in RFC 5952 canonical form: %s", ip.String())
		}
		return nil
	}
	if strings.ContainsAny(h, ":") {
		return errors.New("IPv6 literals must be bracketed")
	}
	if ipv4Re.MatchString(h) {
		return nil
	}
	if looksNumeric(h) {
		return errors.New("IPv4 must be four decimal octets without leading zeros")
	}
	if len(h) > 253 {
		return errors.New("host longer than 253 characters")
	}
	if strings.HasSuffix(h, ".") {
		return errors.New("trailing dot on host")
	}
	for _, label := range strings.Split(h, ".") {
		if !dnsLabelRe.MatchString(label) {
			return fmt.Errorf("label %q is not lowercase RFC 1123", label)
		}
	}
	return nil
}

func looksNumeric(h string) bool {
	for _, r := range h {
		if r != '.' && (r < '0' || r > '9') {
			return false
		}
	}
	return true
}

func isLoopbackHost(h string) bool {
	if h == "localhost" {
		return true
	}
	inner := strings.TrimSuffix(strings.TrimPrefix(h, "["), "]")
	ip := net.ParseIP(inner)
	return ip != nil && ip.IsLoopback()
}

// checkPath applies the configured-path grammar: begins with '/', no trailing
// '/', no empty, '.' or '..' segment, ASCII, no percent-encoding.
func checkPath(p string) error {
	if p == "" || p[0] != '/' {
		return errors.New("must begin with /")
	}
	if len(p) > 1 && strings.HasSuffix(p, "/") {
		return errors.New("no trailing /")
	}
	segs := strings.Split(p[1:], "/")
	for _, s := range segs {
		switch {
		case s == "":
			return errors.New("empty path segment")
		case s == "." || s == "..":
			return errors.New("dot segments are refused")
		case !pathSegRe.MatchString(s):
			return fmt.Errorf("segment %q: only [A-Za-z0-9._~-] is accepted", s)
		}
	}
	return nil
}

type urlRules struct {
	schemes    []string
	originOnly bool // no path, query, fragment
	allowPath  bool // path allowed, still no query/fragment
}

// checkURL refuses any non-canonical spelling: uppercase scheme or host,
// explicit default port, userinfo, and whatever the rules forbid. It returns
// the parsed URL for derived values.
func checkURL(raw string, rules urlRules) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" || u.Host == "" || !strings.HasPrefix(raw, u.Scheme+"://") {
		return nil, errors.New("must be an absolute URL with a lowercase scheme")
	}
	okScheme := false
	for _, s := range rules.schemes {
		if u.Scheme == s {
			okScheme = true
		}
	}
	if !okScheme {
		return nil, fmt.Errorf("scheme must be one of %s", strings.Join(rules.schemes, ", "))
	}
	if u.User != nil {
		return nil, errors.New("userinfo is refused")
	}
	host, port := u.Hostname(), u.Port()
	if strings.Contains(u.Host, ":") && !strings.HasPrefix(u.Host, "[") && strings.Count(u.Host, ":") > 1 {
		return nil, errors.New("IPv6 literals must be bracketed")
	}
	if strings.HasPrefix(u.Host, "[") {
		host = "[" + host + "]"
	}
	if err := checkHost(host); err != nil {
		return nil, err
	}
	if port != "" {
		n, perr := strconv.Atoi(port)
		if perr != nil || n < 1 || n > 65535 || strconv.Itoa(n) != port {
			return nil, errors.New("port must be decimal 1-65535 without leading zeros")
		}
		if (u.Scheme == "https" && port == "443") || (u.Scheme == "http" && port == "80") {
			return nil, errors.New("explicit default port is refused")
		}
	}
	if u.RawQuery != "" || u.ForceQuery {
		return nil, errors.New("query is refused")
	}
	if u.Fragment != "" || strings.Contains(raw, "#") {
		return nil, errors.New("fragment is refused")
	}
	if rules.originOnly && u.Path != "" {
		return nil, errors.New("origin only: no path, not even a trailing /")
	}
	if u.Path != "" {
		if u.RawPath != "" || strings.Contains(u.Path, "%") {
			return nil, errors.New("percent-encoding in path is refused")
		}
		if err := checkPath(u.Path); err != nil {
			return nil, err
		}
	}
	return u, nil
}

// Ref is a B2 secret reference: env:NAME or file:/path. The referenced value
// is never stored here.
type Ref struct {
	Kind  string // "env" | "file"
	Value string // the NAME or the path, safe to print
}

func (r Ref) String() string { return r.Kind + ":" + r.Value }

func parseRef(s string) (Ref, error) {
	switch {
	case strings.HasPrefix(s, "env:"):
		name := s[4:]
		if !envNameRe.MatchString(name) {
			return Ref{}, errors.New("env reference name must match [A-Z_][A-Z0-9_]*")
		}
		return Ref{Kind: "env", Value: name}, nil
	case strings.HasPrefix(s, "file:"):
		p := s[5:]
		if p == "" || p[0] != '/' {
			return Ref{}, errors.New("file reference must be an absolute path")
		}
		return Ref{Kind: "file", Value: p}, nil
	}
	return Ref{}, errors.New("secrets are references: env:NAME or file:/path; a literal value is refused")
}

func checkScope(s string) error {
	if !scopeRe.MatchString(s) {
		return errors.New("scope must be RFC 6749 scope-token characters")
	}
	if len(s) > 256 {
		return errors.New("scope longer than 256 bytes")
	}
	return nil
}

// checkPurpose applies the B5 purpose shape: non-blank, at most 512 UTF-8
// bytes after trimming, no control characters.
func checkPurpose(s string) error {
	t := strings.TrimSpace(s)
	if t == "" {
		return errors.New("blank purpose")
	}
	if !utf8.ValidString(t) {
		return errors.New("purpose is not valid UTF-8")
	}
	if len(t) > purposeMaxBytes {
		return fmt.Errorf("purpose longer than %d bytes", purposeMaxBytes)
	}
	for _, r := range t {
		if unicode.IsControl(r) {
			return errors.New("control characters in purpose")
		}
	}
	return nil
}

func checkCIDROrIP(s string) error {
	if strings.Contains(s, "/") {
		_, ipnet, err := net.ParseCIDR(s)
		if err != nil {
			return errors.New("not a CIDR")
		}
		if ipnet.String() != s {
			return fmt.Errorf("CIDR must be canonical: %s", ipnet.String())
		}
		return nil
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return errors.New("not an IP address")
	}
	if ip.String() != s {
		return fmt.Errorf("IP must be canonical: %s", ip.String())
	}
	return nil
}

// checkHostPort accepts host or host:port in the B3 grammar, for allowedHosts
// entries matched against the effective authority.
func checkHostPort(s string) error {
	host, port := s, ""
	if strings.HasPrefix(s, "[") {
		end := strings.Index(s, "]")
		if end < 0 {
			return errors.New("unterminated IPv6 literal")
		}
		host = s[:end+1]
		rest := s[end+1:]
		if rest != "" {
			if !strings.HasPrefix(rest, ":") {
				return errors.New("must be host or host:port")
			}
			port = rest[1:]
		}
	} else if i := strings.LastIndex(s, ":"); i >= 0 {
		host, port = s[:i], s[i+1:]
	}
	if err := checkHost(host); err != nil {
		return err
	}
	if port != "" {
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 || strconv.Itoa(n) != port {
			return errors.New("port must be decimal 1-65535 without leading zeros")
		}
	}
	return nil
}

// checkListen parses host:port with the B3 host grammar.
func checkListen(s string) (host string, port int, err error) {
	h, p, err := net.SplitHostPort(s)
	if err != nil {
		return "", 0, errors.New("must be host:port")
	}
	if strings.Contains(h, ":") {
		h = "[" + h + "]"
	}
	if err := checkHost(h); err != nil {
		return "", 0, err
	}
	n, err := strconv.Atoi(p)
	if err != nil || n < 1 || n > 65535 || strconv.Itoa(n) != p {
		return "", 0, errors.New("port must be decimal 1-65535 without leading zeros")
	}
	return h, n, nil
}

// checkOriginOrURL accepts an exact origin or exact URL for the redirect
// allowlist; wildcards are refused.
func checkOriginOrURL(s string) error {
	if strings.Contains(s, "*") {
		return errors.New("wildcards are refused")
	}
	_, err := checkURL(s, urlRules{schemes: []string{"https", "http"}, allowPath: true})
	return err
}
