package gpg

import (
	"testing"
)

// TestPGPName verifies the PGP extension name.
func TestPGPName(t *testing.T) {
	ext := PGP{}
	if got := ext.Name(); got != "pgp" {
		t.Errorf("Name() = %q, want %q", got, "pgp")
	}
}

// TestPGPInit verifies the Init method registers components correctly.
func TestPGPInit(t *testing.T) {
	// This test verifies that Init() can be called without panicking
	// and performs the necessary registrations. Full integration testing
	// of the registered components is done in separate tests.
	ext := PGP{}

	// Init should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Init() panicked: %v", r)
		}
	}()

	ext.Init()
}

// TestEXTConstant verifies the file extension constant.
func TestEXTConstant(t *testing.T) {
	expected := ".md.pgp"
	if EXT != expected {
		t.Errorf("EXT = %q, want %q", EXT, expected)
	}
}
