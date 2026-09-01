// Command atesaki is the gateway binary. This build carries the pure
// configuration verbs: validate and routes (docs/contract.md §9).
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/acartag7/atesaki/internal/config"
)

const usage = `usage:
  atesaki validate <config.yaml>   validate the configuration; touches nothing else
  atesaki routes   <config.yaml>   print the route and well-known path list as JSON
`

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, usage)
		return 2
	}
	switch args[0] {
	case "validate":
		return runValidate(args[1:])
	case "routes":
		return runRoutes(args[1:])
	}
	fmt.Fprintf(os.Stderr, "unknown verb %q\n%s", args[0], usage)
	return 2
}

func loadArg(args []string) (*config.Config, int) {
	for _, a := range args {
		if a == "--deep" {
			fmt.Fprintln(os.Stderr, "validate --deep is not in this build")
			return nil, 2
		}
	}
	if len(args) != 1 {
		fmt.Fprint(os.Stderr, usage)
		return nil, 2
	}
	cfg, refusals := config.Load(args[0])
	if refusals != nil {
		for _, r := range refusals {
			fmt.Fprintln(os.Stderr, r.String())
		}
		fmt.Fprintf(os.Stderr, "refused: %d problem(s)\n", len(refusals))
		return nil, 1
	}
	return cfg, 0
}

func runValidate(args []string) int {
	cfg, code := loadArg(args)
	if cfg == nil {
		return code
	}
	fmt.Println(cfg.Summary())
	return 0
}

type routeOut struct {
	Name                      string `json:"name"`
	Path                      string `json:"path"`
	MCPPath                   string `json:"mcpPath"`
	Audience                  string `json:"audience"`
	ProtectedResourceMetadata string `json:"protectedResourceMetadata"`
}

type routesOut struct {
	Issuer                      string     `json:"issuer"`
	Routes                      []routeOut `json:"routes"`
	AuthorizationServerMetadata []string   `json:"authorizationServerMetadata"`
	Reserved                    []string   `json:"reserved"`
}

func runRoutes(args []string) int {
	cfg, code := loadArg(args)
	if cfg == nil {
		return code
	}
	g := cfg.Gateway
	out := routesOut{
		Issuer:                      g.ExternalBaseURL,
		AuthorizationServerMetadata: []string{"/.well-known/oauth-authorization-server", "/.well-known/openid-configuration"},
		Reserved:                    []string{"/oauth", "/.well-known", g.Health.LivePath, g.Health.ReadyPath},
	}
	if g.Identity.CallbackPath != "" {
		out.Reserved = append(out.Reserved, g.Identity.CallbackPath)
	}
	for _, r := range cfg.Routes {
		out.Routes = append(out.Routes, routeOut{
			Name: r.Name, Path: r.Path, MCPPath: r.MCPPath, Audience: r.Audience,
			ProtectedResourceMetadata: r.MetadataURL,
		})
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
