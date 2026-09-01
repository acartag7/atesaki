package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

// Tags the core YAML schema resolves plain scalars and collections to. Anything
// else (anchors, aliases, merge keys, timestamps, binary, custom tags) is a
// refusal under B1: "YAML anchors, aliases, and custom tags refused".
var allowedTags = map[string]bool{
	"!!str": true, "!!int": true, "!!bool": true, "!!float": true,
	"!!null": true, "!!map": true, "!!seq": true,
}

// decodeStream parses the YAML stream into one generic tree per document
// (map[string]any / []any / scalars) after refusing every non-core construct
// and every duplicate mapping key. Refusals are collected, never coerced.
func decodeStream(data []byte, c *collector) []any {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	var docs []any
	for i := 1; ; i++ {
		resource := fmt.Sprintf("document %d", i)
		var n yaml.Node
		err := dec.Decode(&n)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			c.addf(resource, "", "B1.yaml-syntax", "YAML parse error: %v", err)
			break
		}
		root := &n
		if root.Kind == yaml.DocumentNode {
			if len(root.Content) == 0 {
				c.add(resource, "", "B1.empty-document", "document is empty")
				continue
			}
			root = root.Content[0]
		}
		before := len(c.rs)
		checkNode(root, resource, "", c)
		if len(c.rs) != before {
			continue
		}
		var tree any
		if err := root.Decode(&tree); err != nil {
			c.addf(resource, "", "B1.yaml-syntax", "cannot decode document: %v", err)
			continue
		}
		if tree == nil {
			c.add(resource, "", "B1.empty-document", "document is null")
			continue
		}
		docs = append(docs, tree)
	}
	return docs
}

func checkNode(n *yaml.Node, resource, path string, c *collector) {
	if n.Kind == yaml.AliasNode {
		c.add(resource, path, "B1.yaml-alias", "YAML aliases are refused")
		return
	}
	if n.Anchor != "" {
		c.add(resource, path, "B1.yaml-anchor", "YAML anchors are refused")
	}
	if n.Tag != "" && !allowedTags[n.Tag] {
		c.addf(resource, path, "B1.yaml-tag", "tag %s is refused; only plain scalars, maps, and lists are accepted", n.Tag)
	}
	switch n.Kind {
	case yaml.MappingNode:
		seen := map[string]bool{}
		for i := 0; i+1 < len(n.Content); i += 2 {
			k, v := n.Content[i], n.Content[i+1]
			if k.Kind != yaml.ScalarNode || k.Tag != "!!str" {
				c.add(resource, path, "B1.wrong-type", "mapping keys must be plain strings")
				continue
			}
			child := k.Value
			if path != "" {
				child = path + "." + k.Value
			}
			if seen[k.Value] {
				c.add(resource, child, "B1.duplicate-key", "duplicate key")
			}
			seen[k.Value] = true
			checkNode(v, resource, child, c)
		}
	case yaml.SequenceNode:
		for i, v := range n.Content {
			checkNode(v, resource, fmt.Sprintf("%s[%d]", path, i), c)
		}
	}
}
