package blocks

import (
	"embed"
	"testing"

	"github.com/emad-elsaid/xlog"
)

func TestBlocksExtensionName(t *testing.T) {
	ext := Blocks{}
	expected := "blocks"

	if ext.Name() != expected {
		t.Errorf("Expected name %q, got %q", expected, ext.Name())
	}
}

func TestTemplatesEmbedded(t *testing.T) {
	// Verify embedded templates exist
	entries, err := templates.ReadDir("templates")
	if err != nil {
		t.Fatalf("Failed to read templates directory: %v", err)
	}

	if len(entries) == 0 {
		t.Error("Expected embedded templates, got none")
	}

	// Verify expected template files exist
	expectedTemplates := []string{"book.html", "github-user.html", "hero.html", "person.html"}
	for _, expectedFile := range expectedTemplates {
		found := false
		for _, entry := range entries {
			if entry.Name() == expectedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected template %q not found", expectedFile)
		}
	}
}

func TestPublicEmbedded(t *testing.T) {
	// Verify embedded public files exist
	entries, err := public.ReadDir("public")
	if err != nil {
		t.Fatalf("Failed to read public directory: %v", err)
	}

	if len(entries) == 0 {
		t.Error("Expected embedded public files, got none")
	}

	// Verify CSS file exists
	_, err = public.ReadFile("public/blocks.css")
	if err != nil {
		t.Errorf("Expected blocks.css to exist: %v", err)
	}
}

func TestBlockRenderInvalidYAML(t *testing.T) {
	// Test that invalid YAML returns error HTML
	invalidYAML := "this is not: valid: yaml: at: all:"
	renderFn := block("hero")
	result := renderFn(xlog.Markdown(invalidYAML))

	// Should contain error message for invalid YAML
	resultStr := string(result)
	if resultStr == "" {
		t.Error("Expected error message for invalid YAML, got empty string")
	}
	// The error should mention yaml or unmarshal
	if !containsAny(resultStr, []string{"yaml", "unmarshal", "cannot"}) {
		t.Errorf("Expected YAML error message, got: %q", resultStr)
	}
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

func TestStyleWidget(t *testing.T) {
	// style function doesn't use the Page parameter, so nil is fine
	result := style(nil)
	expected := `<link rel="stylesheet" href="/public/blocks.css">`

	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, string(result))
	}
}

func TestRegisterShortCodesFindsTemplates(t *testing.T) {
	// This tests that RegisterShortCodes can walk the embedded templates
	// We can't easily test the side effect without mocking xlog registration,
	// but we can verify templates are present
	entries, err := templates.ReadDir("templates")
	if err != nil {
		t.Fatalf("RegisterShortCodes would fail: %v", err)
	}

	templateCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			templateCount++
		}
	}

	if templateCount == 0 {
		t.Error("RegisterShortCodes would register no templates")
	}

	// Verify we have the expected template count
	if templateCount != 4 {
		t.Errorf("Expected 4 templates, found %d", templateCount)
	}
}

func TestEmbedFSNotEmpty(t *testing.T) {
	// Verify templates embed.FS is not empty
	if templates == (embed.FS{}) {
		t.Error("templates embed.FS should not be empty")
	}

	// Verify public embed.FS is not empty
	if public == (embed.FS{}) {
		t.Error("public embed.FS should not be empty")
	}
}
