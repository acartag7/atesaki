package config

import (
	"runtime"
	"testing"
)

func TestCIGateProbe(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Fatal("deliberate macOS matrix failure for the owner merge refusal proof")
	}
}
