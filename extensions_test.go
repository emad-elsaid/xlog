package xlog

import (
	"testing"
)

type mockExtension struct {
	name        string
	initialized bool
}

func (m *mockExtension) Name() string {
	return m.name
}

func (m *mockExtension) Init() {
	m.initialized = true
}

func TestRegisterExtension(t *testing.T) {
	// Save original extensions and restore after test
	originalExtensions := extensions
	defer func() { extensions = originalExtensions }()

	extensions = []Extension{}

	mock := &mockExtension{name: "test-extension"}
	RegisterExtension(mock)

	if len(extensions) != 1 {
		t.Errorf("Expected 1 extension, got %d", len(extensions))
	}

	if extensions[0].Name() != "test-extension" {
		t.Errorf("Expected extension name 'test-extension', got '%s'", extensions[0].Name())
	}
}

func TestInitExtensions_AllDisabled(t *testing.T) {
	// Save original state
	originalExtensions := extensions
	originalConfig := Config
	defer func() {
		extensions = originalExtensions
		Config = originalConfig
	}()

	// Setup
	mock1 := &mockExtension{name: "ext1"}
	mock2 := &mockExtension{name: "ext2"}
	extensions = []Extension{mock1, mock2}
	Config.DisabledExtensions = "all"

	// Execute
	initExtensions()

	// Verify none were initialized
	if mock1.initialized {
		t.Error("Expected ext1 to not be initialized when all extensions disabled")
	}
	if mock2.initialized {
		t.Error("Expected ext2 to not be initialized when all extensions disabled")
	}
}

func TestInitExtensions_SpecificDisabled(t *testing.T) {
	// Save original state
	originalExtensions := extensions
	originalConfig := Config
	defer func() {
		extensions = originalExtensions
		Config = originalConfig
	}()

	// Setup
	mock1 := &mockExtension{name: "ext1"}
	mock2 := &mockExtension{name: "ext2"}
	mock3 := &mockExtension{name: "ext3"}
	extensions = []Extension{mock1, mock2, mock3}
	Config.DisabledExtensions = "ext1,ext3"

	// Execute
	initExtensions()

	// Verify
	if mock1.initialized {
		t.Error("Expected ext1 to not be initialized (disabled)")
	}
	if !mock2.initialized {
		t.Error("Expected ext2 to be initialized (enabled)")
	}
	if mock3.initialized {
		t.Error("Expected ext3 to not be initialized (disabled)")
	}
}

func TestInitExtensions_AllEnabled(t *testing.T) {
	// Save original state
	originalExtensions := extensions
	originalConfig := Config
	defer func() {
		extensions = originalExtensions
		Config = originalConfig
	}()

	// Setup
	mock1 := &mockExtension{name: "ext1"}
	mock2 := &mockExtension{name: "ext2"}
	extensions = []Extension{mock1, mock2}
	Config.DisabledExtensions = ""

	// Execute
	initExtensions()

	// Verify all initialized
	if !mock1.initialized {
		t.Error("Expected ext1 to be initialized")
	}
	if !mock2.initialized {
		t.Error("Expected ext2 to be initialized")
	}
}

func TestInitExtensions_InvalidExtensionName(t *testing.T) {
	// Save original state
	originalExtensions := extensions
	originalConfig := Config
	defer func() {
		extensions = originalExtensions
		Config = originalConfig
	}()

	// Setup - disable non-existent extension
	mock1 := &mockExtension{name: "ext1"}
	extensions = []Extension{mock1}
	Config.DisabledExtensions = "nonexistent,ext1"

	// Execute
	initExtensions()

	// Verify ext1 was disabled correctly despite invalid name in list
	if mock1.initialized {
		t.Error("Expected ext1 to not be initialized (disabled)")
	}
}
