package config

import (
	"fmt"
	"strings"
)

// obj reads one mapping of the generic tree with the B1 malformed-input rules:
// wrong type, unknown field, null for a required field, blank string after
// trimming, list-for-scalar, scalar-for-list. Every accessor marks the key as
// seen; finish() refuses whatever was never read.
type obj struct {
	c        *collector
	resource string
	path     string
	m        map[string]any
	seen     map[string]bool
}

func newObj(c *collector, resource, path string, v any) (*obj, bool) {
	m, ok := v.(map[string]any)
	if !ok {
		c.addf(resource, path, "B1.wrong-type", "expected a mapping, got %s", typeName(v))
		return nil, false
	}
	return &obj{c: c, resource: resource, path: path, m: m, seen: map[string]bool{}}, true
}

func (o *obj) field(name string) string {
	if o.path == "" {
		return name
	}
	return o.path + "." + name
}

func (o *obj) has(name string) bool {
	_, ok := o.m[name]
	return ok
}

func (o *obj) raw(name string) (any, bool) {
	o.seen[name] = true
	v, ok := o.m[name]
	return v, ok
}

func (o *obj) missing(name string) {
	o.c.add(o.resource, o.field(name), "B1.required", "required field is missing")
}

// str returns a non-blank string. Absent optional fields return ("", false)
// without a refusal.
func (o *obj) str(name string, required bool) (string, bool) {
	v, ok := o.raw(name)
	if !ok {
		if required {
			o.missing(name)
		}
		return "", false
	}
	if v == nil {
		if required {
			o.c.add(o.resource, o.field(name), "B1.null-required", "null is not a value for a required field")
		} else {
			o.c.add(o.resource, o.field(name), "B1.wrong-type", "null is not a string")
		}
		return "", false
	}
	s, isStr := v.(string)
	if !isStr {
		o.c.addf(o.resource, o.field(name), "B1.wrong-type", "expected a string, got %s", typeName(v))
		return "", false
	}
	if strings.TrimSpace(s) == "" {
		o.c.add(o.resource, o.field(name), "B1.blank", "blank string")
		return "", false
	}
	if s != strings.TrimSpace(s) {
		o.c.add(o.resource, o.field(name), "B1.blank", "leading or trailing whitespace")
		return "", false
	}
	return s, true
}

// integer returns an integer. YAML floats, bools, and strings are refused.
func (o *obj) integer(name string, required bool) (int64, bool) {
	v, ok := o.raw(name)
	if !ok {
		if required {
			o.missing(name)
		}
		return 0, false
	}
	if v == nil {
		o.c.add(o.resource, o.field(name), "B1.null-required", "null is not an integer")
		return 0, false
	}
	switch n := v.(type) {
	case int:
		return int64(n), true
	case int64:
		return n, true
	case uint64:
		if n > 1<<62 {
			o.c.add(o.resource, o.field(name), "B1.wrong-type", "integer out of range")
			return 0, false
		}
		return int64(n), true
	}
	o.c.addf(o.resource, o.field(name), "B1.wrong-type", "expected an integer, got %s", typeName(v))
	return 0, false
}

// positive returns an integer >= 1.
func (o *obj) positive(name string, required bool) (int64, bool) {
	n, ok := o.integer(name, required)
	if !ok {
		return 0, false
	}
	if n < 1 {
		o.c.add(o.resource, o.field(name), "B1.wrong-type", "expected a positive integer")
		return 0, false
	}
	return n, true
}

func (o *obj) boolean(name string, def bool) bool {
	v, ok := o.raw(name)
	if !ok {
		return def
	}
	b, isBool := v.(bool)
	if !isBool {
		o.c.addf(o.resource, o.field(name), "B1.wrong-type", "expected true or false, got %s", typeName(v))
		return def
	}
	return b
}

// literalTrue accepts only the literal `true` (B1: `identity.publicClient: true`).
func (o *obj) literalTrue(name string) (present bool, ok bool) {
	v, has := o.raw(name)
	if !has {
		return false, true
	}
	if b, isBool := v.(bool); isBool && b {
		return true, true
	}
	o.c.add(o.resource, o.field(name), "B1.wrong-type", "only the literal true is accepted")
	return true, false
}

func (o *obj) list(name string, required bool) ([]any, bool) {
	v, ok := o.raw(name)
	if !ok {
		if required {
			o.missing(name)
		}
		return nil, false
	}
	if v == nil {
		o.c.add(o.resource, o.field(name), "B1.null-required", "null is not a list")
		return nil, false
	}
	l, isList := v.([]any)
	if !isList {
		o.c.addf(o.resource, o.field(name), "B1.wrong-type", "expected a list, got %s", typeName(v))
		return nil, false
	}
	return l, true
}

// strList returns a list of non-blank strings.
func (o *obj) strList(name string, required bool) ([]string, bool) {
	l, ok := o.list(name, required)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(l))
	good := true
	for i, v := range l {
		f := fmt.Sprintf("%s[%d]", o.field(name), i)
		s, isStr := v.(string)
		if !isStr {
			o.c.addf(o.resource, f, "B1.wrong-type", "expected a string, got %s", typeName(v))
			good = false
			continue
		}
		if strings.TrimSpace(s) == "" || s != strings.TrimSpace(s) {
			o.c.add(o.resource, f, "B1.blank", "blank or padded string")
			good = false
			continue
		}
		out = append(out, s)
	}
	return out, good
}

func (o *obj) object(name string, required bool) (*obj, bool) {
	v, ok := o.raw(name)
	if !ok {
		if required {
			o.missing(name)
		}
		return nil, false
	}
	if v == nil {
		o.c.add(o.resource, o.field(name), "B1.null-required", "null is not a mapping")
		return nil, false
	}
	return newObj(o.c, o.resource, o.field(name), v)
}

// forbid refuses a field that belongs to an inactive variant.
func (o *obj) forbid(name, why string) {
	if o.has(name) {
		o.seen[name] = true
		o.c.add(o.resource, o.field(name), "B1.inactive-variant", why)
	}
}

func (o *obj) finish() {
	for k := range o.m {
		if !o.seen[k] {
			o.c.add(o.resource, o.field(k), "B1.unknown-field", "unknown field")
		}
	}
}

func typeName(v any) string {
	switch v.(type) {
	case nil:
		return "null"
	case string:
		return "string"
	case bool:
		return "boolean"
	case int, int64, uint64:
		return "integer"
	case float64:
		return "float"
	case []any:
		return "list"
	case map[string]any:
		return "mapping"
	}
	return fmt.Sprintf("%T", v)
}
