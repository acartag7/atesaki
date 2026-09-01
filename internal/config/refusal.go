// Package config loads the Atesaki YAML stream into typed values, refusing
// anything the contract does not accept. It implements docs/contract-boundaries.md
// B1 (configuration), B2 (reference trust), B3 (canonicalization), the B4
// assertion shape, and the config-side rows of B5/B8. Interior code never
// re-checks what these types guarantee.
package config

import (
	"fmt"
	"strings"
)

// Refusal names the resource, the field path, the contract rule, and a plain
// explanation. It never carries a resolved secret value.
type Refusal struct {
	Resource string
	Field    string
	Rule     string
	Detail   string
}

func (r Refusal) String() string {
	field := r.Field
	if field == "" {
		field = "-"
	}
	return fmt.Sprintf("%s %s %s: %s", r.Resource, field, r.Rule, r.Detail)
}

// Refusals is the complete list of reasons a load was refused. A load with any
// refusal produces no Config.
type Refusals []Refusal

func (rs Refusals) Error() string {
	lines := make([]string, 0, len(rs))
	for _, r := range rs {
		lines = append(lines, r.String())
	}
	return strings.Join(lines, "\n")
}

type collector struct {
	rs Refusals
}

func (c *collector) add(resource, field, rule, detail string) {
	c.rs = append(c.rs, Refusal{Resource: resource, Field: field, Rule: rule, Detail: detail})
}

func (c *collector) addf(resource, field, rule, format string, args ...any) {
	c.add(resource, field, rule, fmt.Sprintf(format, args...))
}

func (c *collector) failed() bool { return len(c.rs) > 0 }
