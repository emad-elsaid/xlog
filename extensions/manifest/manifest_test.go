package manifest

import (
	"html/template"
	"strings"
	"testing"

	. "github.com/emad-elsaid/xlog"
)

func TestManifestExtensionName(t *testing.T) {
	ext := Manifest{}
	expected := "manifest"
	if ext.Name() != expected {
		t.Errorf("Expected extension name to be %q, got %q", expected, ext.Name())
	}
}

func TestManifestInit_DoesNotPanic(t *testing.T) {
	// Initialize extension (should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Init() panicked: %v", r)
		}
	}()

	ext := Manifest{}
	ext.Init()
}

func TestManifestHead_ReturnsLinkTag(t *testing.T) {
	// Call head widget with empty page
	page := Page(nil)
	result := head(page)

	// Check it returns link tag
	expected := `<link rel="manifest" href="/manifest.json">`
	if result != template.HTML(expected) {
		t.Errorf("Expected head widget to return %q, got %q", expected, result)
	}
}

func TestManifestHead_HTMLFormat(t *testing.T) {
	page := Page(nil)
	result := head(page)
	html := string(result)

	// Validate HTML structure
	if !strings.HasPrefix(html, "<link") {
		t.Error("Expected head widget to start with <link tag")
	}
	if !strings.Contains(html, `rel="manifest"`) {
		t.Error("Expected head widget to contain rel=\"manifest\"")
	}
	if !strings.Contains(html, `href="/manifest.json"`) {
		t.Error("Expected head widget to contain href=\"/manifest.json\"")
	}
	if !strings.HasSuffix(html, ">") {
		t.Error("Expected head widget to end with >")
	}
}

func TestManifestHead_NoPageDependency(t *testing.T) {
	// Test that head widget doesn't depend on page content
	page1 := Page(nil)
	page2 := Page(nil)

	result1 := head(page1)
	result2 := head(page2)

	if result1 != result2 {
		t.Error("Expected head widget to return same result regardless of page")
	}
}

func TestManifestTemplate_ContainsRequiredFields(t *testing.T) {
	// Read the embedded template directly
	tmplContent, err := templates.ReadFile("templates/manifest.html")
	if err != nil {
		t.Fatalf("Failed to read manifest template: %v", err)
	}

	content := string(tmplContent)

	// Check all required manifest fields are in template
	requiredFields := []string{
		`"name"`,
		`"short_name"`,
		`"start_url"`,
		`"display"`,
		`"icons"`,
		`"src"`,
		`"type"`,
		`"sizes"`,
	}

	for _, field := range requiredFields {
		if !strings.Contains(content, field) {
			t.Errorf("Expected manifest template to contain field %s", field)
		}
	}
}

func TestManifestTemplate_UsesConfigSitename(t *testing.T) {
	// Read the embedded template
	tmplContent, err := templates.ReadFile("templates/manifest.html")
	if err != nil {
		t.Fatalf("Failed to read manifest template: %v", err)
	}

	content := string(tmplContent)

	// Check template uses config.Sitename
	if !strings.Contains(content, "{{.config.Sitename}}") {
		t.Error("Expected manifest template to use {{.config.Sitename}}")
	}
}

func TestManifestTemplate_JSONStructure(t *testing.T) {
	// Read the embedded template
	tmplContent, err := templates.ReadFile("templates/manifest.html")
	if err != nil {
		t.Fatalf("Failed to read manifest template: %v", err)
	}

	content := string(tmplContent)

	// Basic JSON validation
	if !strings.HasPrefix(strings.TrimSpace(content), "{") {
		t.Error("Expected manifest template to start with {")
	}
	if !strings.HasSuffix(strings.TrimSpace(content), "}") {
		t.Error("Expected manifest template to end with }")
	}
}

func TestManifestTemplate_DisplayStandalone(t *testing.T) {
	tmplContent, err := templates.ReadFile("templates/manifest.html")
	if err != nil {
		t.Fatalf("Failed to read manifest template: %v", err)
	}

	content := string(tmplContent)

	// Check display mode is standalone
	if !strings.Contains(content, `"display": "standalone"`) {
		t.Error("Expected manifest template to set display mode to standalone")
	}
}

func TestManifestTemplate_StartURL(t *testing.T) {
	tmplContent, err := templates.ReadFile("templates/manifest.html")
	if err != nil {
		t.Fatalf("Failed to read manifest template: %v", err)
	}

	content := string(tmplContent)

	// Check start_url is root
	if !strings.Contains(content, `"start_url": "/"`) {
		t.Error("Expected manifest start_url to be /")
	}
}

func TestManifestTemplate_LogoIcon(t *testing.T) {
	tmplContent, err := templates.ReadFile("templates/manifest.html")
	if err != nil {
		t.Fatalf("Failed to read manifest template: %v", err)
	}

	content := string(tmplContent)

	// Check icon references logo.png
	if !strings.Contains(content, "public/logo.png") {
		t.Error("Expected manifest template to reference public/logo.png")
	}
}
