package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Every env reference used by the valid examples. The refusal case
// env-missing references ATESAKI_TEST_UNSET_SECRET, which must stay unset.
var exampleEnv = []string{
	"ATESAKI_SIGNING_KEY", "ENTRA_CLIENT_SECRET", "CORP_CA_BUNDLE",
	"SPLUNK_READ_TOKEN", "SPLUNK_ADMIN_TOKEN", "NIGHTLY_REPORT_SECRET", "GRAFANA_TOKEN",
}

func setExampleEnv(t *testing.T) {
	t.Helper()
	for _, k := range exampleEnv {
		t.Setenv(k, "set-for-test")
	}
	if os.Getenv("ATESAKI_TEST_UNSET_SECRET") != "" {
		t.Fatal("ATESAKI_TEST_UNSET_SECRET must not be set while testing")
	}
}

func TestValidExamples(t *testing.T) {
	setExampleEnv(t)
	files, err := filepath.Glob("testdata/valid/*.yaml")
	if err != nil || len(files) == 0 {
		t.Fatalf("no valid examples found: %v", err)
	}
	for _, f := range files {
		t.Run(filepath.Base(f), func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			cfg, refusals := Parse(data)
			if refusals != nil {
				t.Fatalf("expected acceptance, got:\n%s", refusals.Error())
			}
			if cfg.Gateway == nil || len(cfg.Routes) == 0 {
				t.Fatal("accepted config without gateway or routes")
			}
			for _, r := range cfg.Routes {
				if r.Audience != cfg.Gateway.ExternalBaseURL+r.Path+r.MCPEndpoint {
					t.Errorf("route %s audience %q not derived byte-exact", r.Name, r.Audience)
				}
				if r.MetadataURL != "/.well-known/oauth-protected-resource"+r.MCPPath {
					t.Errorf("route %s metadata path %q", r.Name, r.MetadataURL)
				}
				if r.MaxDuration < 1 || r.MaxDuration > hardMaxDuration {
					t.Errorf("route %s maxDuration %d outside bounds", r.Name, r.MaxDuration)
				}
			}
			if cfg.Gateway.Tokens.AccessTTL != defaultAccessTTL {
				t.Errorf("access TTL default not applied: %d", cfg.Gateway.Tokens.AccessTTL)
			}
		})
	}
}

// Each refusal case names the rule it exercises on its first line. The case
// passes only when the config is refused AND that rule is among the reasons.
func TestRefusalSuite(t *testing.T) {
	setExampleEnv(t)
	files, err := filepath.Glob("testdata/refuse/*.yaml")
	if err != nil || len(files) == 0 {
		t.Fatalf("no refusal cases found: %v", err)
	}
	for _, f := range files {
		t.Run(filepath.Base(f), func(t *testing.T) {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			expect := expectedRule(t, data)
			cfg, refusals := Parse(data)
			if cfg != nil || refusals == nil {
				t.Fatalf("expected refusal for rule %s, config was accepted", expect)
			}
			for _, r := range refusals {
				if r.Rule == expect {
					return
				}
			}
			t.Fatalf("refused, but not for %s:\n%s", expect, refusals.Error())
		})
	}
}

func expectedRule(t *testing.T, data []byte) string {
	t.Helper()
	sc := bufio.NewScanner(strings.NewReader(string(data)))
	if !sc.Scan() {
		t.Fatal("empty case file")
	}
	line := sc.Text()
	const prefix = "# expect: "
	if !strings.HasPrefix(line, prefix) {
		t.Fatalf("first line must be %q<rule>, got %q", prefix, line)
	}
	return strings.TrimSpace(strings.TrimPrefix(line, prefix))
}

func TestRefusalsNeverEchoSecretValues(t *testing.T) {
	setExampleEnv(t)
	data, err := os.ReadFile("testdata/refuse/credential-inline-value.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, refusals := Parse(data)
	if refusals == nil {
		t.Fatal("expected refusal")
	}
	if strings.Contains(refusals.Error(), "abc123") {
		t.Fatalf("refusal echoed the inline value:\n%s", refusals.Error())
	}
}

func TestSecretFileInvariants(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "secret")
	if err := os.WriteFile(good, []byte("s\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := checkSecretFile(good); err != nil {
		t.Fatalf("0600 regular file must be accepted: %v", err)
	}

	loose := filepath.Join(dir, "loose")
	if err := os.WriteFile(loose, []byte("s\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := checkSecretFile(loose); err == nil || !strings.Contains(err.Error(), "mode") {
		t.Fatalf("group/other-readable file must be refused by mode, got %v", err)
	}

	link := filepath.Join(dir, "link")
	if err := os.Symlink(good, link); err != nil {
		t.Fatal(err)
	}
	if err := checkSecretFile(link); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("symlink must be refused, got %v", err)
	}

	if err := checkSecretFile(dir); err == nil || !strings.Contains(err.Error(), "regular") {
		t.Fatalf("directory must be refused, got %v", err)
	}
	if err := checkSecretFile(filepath.Join(dir, "missing")); err == nil {
		t.Fatal("missing file must be refused")
	}

	big := filepath.Join(dir, "big")
	if err := os.WriteFile(big, make([]byte, secretFileMax+1), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := checkSecretFile(big); err == nil || !strings.Contains(err.Error(), "larger") {
		t.Fatalf("oversized file must be refused, got %v", err)
	}

	hard := filepath.Join(dir, "hard")
	if err := os.Link(good, hard); err == nil {
		if err := checkSecretFile(hard); err == nil || !strings.Contains(err.Error(), "link count") {
			t.Fatalf("hard-linked file must be refused, got %v", err)
		}
	}
}

func TestLoadFile(t *testing.T) {
	setExampleEnv(t)
	cfg, refusals := Load("testdata/valid/console-loopback.yaml")
	if refusals != nil {
		t.Fatalf("expected acceptance: %s", refusals.Error())
	}
	if cfg.Gateway.Identity.Provider != "console" || cfg.Gateway.Identity.PairingCodeTTL != defaultPairingCodeTTL {
		t.Fatalf("console defaults not applied: %+v", cfg.Gateway.Identity)
	}
	if _, refusals := Load("testdata/does-not-exist.yaml"); refusals == nil || refusals[0].Rule != "B2.file" {
		t.Fatalf("missing config must refuse with B2.file, got %v", refusals)
	}
	big := filepath.Join(t.TempDir(), "big.yaml")
	if err := os.WriteFile(big, make([]byte, configMaxBytes+1), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, refusals := Load(big); refusals == nil || refusals[0].Rule != "B5.config-size" {
		t.Fatalf("oversized config must refuse with B5.config-size, got %v", refusals)
	}
}

func TestHostGrammar(t *testing.T) {
	accept := []string{"gw.example.com", "127.0.0.1", "[::1]", "[2001:db8::1]", "a", "0.0.0.0"}
	refuse := []string{"", "GW.example.com", "gw.example.com.", "-a.example.com", "127.000.0.1", "1.2.3", "::1", "[2001:DB8::1]", "[2001:db8:0:0:0:0:0:1]", "exämple.com", "a b"}
	for _, h := range accept {
		if err := checkHost(h); err != nil {
			t.Errorf("%q must be accepted: %v", h, err)
		}
	}
	for _, h := range refuse {
		if err := checkHost(h); err == nil {
			t.Errorf("%q must be refused", h)
		}
	}
}

func TestPathGrammar(t *testing.T) {
	accept := []string{"/a", "/a/b-c", "/x.y_z~1"}
	refuse := []string{"", "a", "/", "/a/", "/a//b", "/a/./b", "/a/../b", "/a%2Fb", "/a b", "/ä"}
	for _, p := range accept {
		if err := checkPath(p); err != nil {
			t.Errorf("%q must be accepted: %v", p, err)
		}
	}
	for _, p := range refuse {
		if err := checkPath(p); err == nil {
			t.Errorf("%q must be refused", p)
		}
	}
}

func TestURLGrammar(t *testing.T) {
	origin := urlRules{schemes: []string{"https"}, originOnly: true}
	for _, u := range []string{"https://gw.example.com", "https://gw.example.com:8443", "https://[2001:db8::1]:8443"} {
		if _, err := checkURL(u, origin); err != nil {
			t.Errorf("%q must be accepted as origin: %v", u, err)
		}
	}
	for _, u := range []string{"HTTPS://gw.example.com", "https://gw.example.com/", "https://gw.example.com:443", "https://gw.example.com?x", "https://gw.example.com#f", "https://u@gw.example.com", "http://gw.example.com", "https://GW.example.com", "gw.example.com", "https://gw.example.com:08443"} {
		if _, err := checkURL(u, origin); err == nil {
			t.Errorf("%q must be refused as origin", u)
		}
	}
	withPath := urlRules{schemes: []string{"http", "https"}, allowPath: true}
	if _, err := checkURL("http://127.0.0.1:9000/mcp", withPath); err != nil {
		t.Errorf("upstream URL with path must be accepted: %v", err)
	}
	for _, u := range []string{"http://127.0.0.1:9000/mcp/", "http://127.0.0.1:9000/a//b", "http://127.0.0.1:9000/a%2Fb"} {
		if _, err := checkURL(u, withPath); err == nil {
			t.Errorf("%q must be refused: non-canonical path", u)
		}
	}
}

func TestRefGrammar(t *testing.T) {
	for _, s := range []string{"env:NAME", "env:A_B9", "file:/etc/x"} {
		if _, err := parseRef(s); err != nil {
			t.Errorf("%q must parse: %v", s, err)
		}
	}
	for _, s := range []string{"", "secret", "env:", "env:lower", "env:9X", "file:", "file:relative", "vault:x", "ENV:X"} {
		if _, err := parseRef(s); err == nil {
			t.Errorf("%q must be refused", s)
		}
	}
}

func TestListenGrammar(t *testing.T) {
	for _, s := range []string{"127.0.0.1:8080", "0.0.0.0:8443", "[::]:8443", "[::1]:1"} {
		if _, _, err := checkListen(s); err != nil {
			t.Errorf("%q must be accepted: %v", s, err)
		}
	}
	for _, s := range []string{"8080", ":8080", "127.0.0.1", "127.0.0.1:0", "127.0.0.1:08080", "localhost:70000", "::1:8080"} {
		if _, _, err := checkListen(s); err == nil {
			t.Errorf("%q must be refused", s)
		}
	}
}
