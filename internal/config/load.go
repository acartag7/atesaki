package config

import (
	"fmt"
	"os"
)

// Load reads and validates the configuration file. It touches nothing but the
// file, the environment variables it references, and the metadata of the files
// it references. Any refusal means no Config.
func Load(path string) (*Config, Refusals) {
	c := &collector{}
	fi, err := os.Stat(path)
	if err != nil {
		c.addf("config", "", "B2.file", "%s: %v", path, unwrapPathError(err))
		return nil, c.rs
	}
	if !fi.Mode().IsRegular() {
		c.addf("config", "", "B2.file", "%s: not a regular file", path)
		return nil, c.rs
	}
	if fi.Size() > configMaxBytes {
		c.addf("config", "", "B5.config-size", "%s: larger than %d bytes", path, configMaxBytes)
		return nil, c.rs
	}
	data, err := os.ReadFile(path)
	if err != nil {
		c.addf("config", "", "B2.file", "%s: %v", path, unwrapPathError(err))
		return nil, c.rs
	}
	return parse(data, true)
}

// Parse validates configuration bytes, including reference resolution.
func Parse(data []byte) (*Config, Refusals) { return parse(data, true) }

// ParseStructural validates everything except reference resolution (B2 env
// and file checks). Used where the environment is not the one that will run.
func ParseStructural(data []byte) (*Config, Refusals) { return parse(data, false) }

func parse(data []byte, resolveRefs bool) (*Config, Refusals) {
	c := &collector{}
	if len(data) > configMaxBytes {
		c.addf("config", "", "B5.config-size", "larger than %d bytes", configMaxBytes)
		return nil, c.rs
	}
	docs := decodeStream(data, c)
	if c.failed() {
		return nil, c.rs
	}
	cfg := parseDocuments(docs, c)
	if c.failed() {
		return nil, c.rs
	}
	crossValidate(cfg, c)
	if c.failed() {
		return nil, c.rs
	}
	if resolveRefs {
		checkRefs(cfg, c)
		if c.failed() {
			return nil, c.rs
		}
	}
	return cfg, nil
}

// Summary is the one-line success report for `atesaki validate`.
func (cfg *Config) Summary() string {
	return fmt.Sprintf("ok: Gateway/%s (%s), %d route(s)", cfg.Gateway.Name, cfg.Gateway.Identity.Provider, len(cfg.Routes))
}
