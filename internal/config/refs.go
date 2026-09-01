package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// checkRefs resolves every B2 reference far enough to refuse a missing one by
// its exact name, without reading or echoing any value. env: must be set and
// non-blank. file: must satisfy the filesystem invariants on open.
func checkRefs(cfg *Config, c *collector) {
	g := cfg.Gateway
	if g == nil {
		return
	}
	res := "Gateway/" + g.Name
	checkRef(c, res, "spec.signingKeyRef", g.SigningKeyRef)
	if g.Identity.ClientSecretRef != nil {
		checkRef(c, res, "spec.identity.clientSecretRef", *g.Identity.ClientSecretRef)
	}
	if a := g.Identity.Assertion; a != nil && a.Keys.JWKSRef != nil {
		checkRef(c, res, "spec.identity.assertion.keys.jwksRef", *a.Keys.JWKSRef)
	}
	for name, ep := range g.Egress {
		if ep.CABundleRef != nil {
			checkRef(c, res, "spec.egress.profiles."+name+".caBundleRef", *ep.CABundleRef)
		}
	}
	for i, mc := range g.MachineClients {
		checkRef(c, res, fmt.Sprintf("spec.machineClients[%d].secretRef", i), mc.SecretRef)
	}
	for i, p := range g.Clients.KnownCIMD {
		if err := checkSecretFile(p); err != nil {
			c.add(res, fmt.Sprintf("spec.clients.knownCimd[%d]", i), "B2.file", err.Error())
		}
	}
	for _, r := range cfg.Routes {
		if r.Upstream.Credential.ValueRef != nil {
			checkRef(c, "Route/"+r.Name, "spec.upstream.credential.valueRef", *r.Upstream.Credential.ValueRef)
		}
	}
}

func checkRef(c *collector, res, field string, r Ref) {
	switch r.Kind {
	case "env":
		v, ok := os.LookupEnv(r.Value)
		if !ok {
			c.addf(res, field, "B2.env-missing", "environment variable %s is not set", r.Value)
			return
		}
		if isBlank(v) {
			c.addf(res, field, "B2.env-missing", "environment variable %s is empty", r.Value)
		}
	case "file":
		if err := checkSecretFile(r.Value); err != nil {
			c.add(res, field, "B2.file", err.Error())
		}
	}
}

func isBlank(s string) bool {
	for _, r := range s {
		if r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			return false
		}
	}
	return true
}

// checkSecretFile applies the B2 invariants for a secret-reference file:
// opened with O_NOFOLLOW, a regular file, link count 1, owned by the process
// user, mode granting nothing to group or other, parent directory owned by the
// process user and not group/other-writable, size within the cap.
func checkSecretFile(path string) error {
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
	if err != nil {
		if errors.Is(err, syscall.ELOOP) {
			return fmt.Errorf("%s: symlink in the final path component is refused", path)
		}
		return fmt.Errorf("%s: cannot open: %v", path, unwrapPathError(err))
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("%s: cannot stat", path)
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("%s: not a regular file", path)
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("%s: cannot read ownership", path)
	}
	if uint64(st.Nlink) != 1 {
		return fmt.Errorf("%s: link count must be exactly 1", path)
	}
	if int(st.Uid) != os.Getuid() {
		return fmt.Errorf("%s: must be owned by the process user", path)
	}
	if fi.Mode().Perm()&0o077 != 0 {
		return fmt.Errorf("%s: mode must grant nothing to group or other (0600 or stricter)", path)
	}
	if fi.Size() > secretFileMax {
		return fmt.Errorf("%s: larger than %d bytes", path, secretFileMax)
	}
	dir := filepath.Dir(path)
	di, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("%s: cannot stat parent directory", path)
	}
	dst, ok := di.Sys().(*syscall.Stat_t)
	if !ok || int(dst.Uid) != os.Getuid() {
		return fmt.Errorf("%s: parent directory must be owned by the process user", path)
	}
	if di.Mode().Perm()&0o022 != 0 {
		return fmt.Errorf("%s: parent directory must not be group/other-writable", path)
	}
	return nil
}

func unwrapPathError(err error) error {
	var pe *os.PathError
	if errors.As(err, &pe) {
		return pe.Err
	}
	return err
}
