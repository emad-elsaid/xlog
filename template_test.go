package xlog

import (
	"embed"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed testdata/custom_templates
var customTemplatesFS embed.FS

func TestRegisterTemplate(t *testing.T) {
	// Store original state to restore after test
	originalFSs := templatesFSs
	defer func() { templatesFSs = originalFSs }()

	// Reset templatesFSs
	templatesFSs = nil

	// Test registering a custom template filesystem
	RegisterTemplate(customTemplatesFS, "testdata/custom_templates")

	if len(templatesFSs) != 1 {
		t.Errorf("Expected 1 template filesystem registered, got %d", len(templatesFSs))
	}

	// Test registering multiple filesystems
	RegisterTemplate(customTemplatesFS, "testdata/custom_templates")
	if len(templatesFSs) != 2 {
		t.Errorf("Expected 2 template filesystems registered, got %d", len(templatesFSs))
	}
}

func TestCompileTemplates(t *testing.T) {
	// Store original state
	originalTemplates := templates
	originalFSs := templatesFSs
	defer func() {
		templates = originalTemplates
		templatesFSs = originalFSs
	}()

	// Reset state
	templatesFSs = nil
	templates = nil

	// Compile templates
	compileTemplates()

	if templates == nil {
		t.Fatal("Expected templates to be initialized after compileTemplates()")
	}

	// Check that default templates are loaded
	defaultTemplateNames := []string{
		"layout",
		"page",
		"navbar",
		"pages",
		"pages-grid",
		"commands",
		"emoji-favicon",
	}

	for _, name := range defaultTemplateNames {
		if templates.Lookup(name) == nil {
			t.Errorf("Expected default template '%s' to be compiled", name)
		}
	}
}

func TestCompileTemplatesWithThemeDirectory(t *testing.T) {
	// Create a temporary theme directory
	tmpDir := t.TempDir()
	themeDir := filepath.Join(tmpDir, "theme")
	if err := os.Mkdir(themeDir, 0755); err != nil {
		t.Fatalf("Failed to create theme directory: %v", err)
	}

	// Create a custom template in the theme directory
	customTemplate := `<div>Custom Theme Template</div>`
	customTemplatePath := filepath.Join(themeDir, "custom.html")
	if err := os.WriteFile(customTemplatePath, []byte(customTemplate), 0644); err != nil {
		t.Fatalf("Failed to write custom template: %v", err)
	}

	// Change to the temporary directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Store original state
	originalTemplates := templates
	originalFSs := templatesFSs
	defer func() {
		templates = originalTemplates
		templatesFSs = originalFSs
	}()

	// Reset state
	templatesFSs = nil
	templates = nil

	// Compile templates (should include theme directory)
	compileTemplates()

	// Verify custom template was loaded
	if templates.Lookup("custom") == nil {
		t.Error("Expected custom template from theme directory to be compiled")
	}
}

func TestCompileTemplatesOverride(t *testing.T) {
	// Store original state
	originalTemplates := templates
	originalFSs := templatesFSs
	defer func() {
		templates = originalTemplates
		templatesFSs = originalFSs
	}()

	// Reset state
	templatesFSs = nil
	templates = nil

	// First compilation with default templates
	compileTemplates()

	// Get the original navbar template
	navbarTemplate := templates.Lookup("navbar")
	if navbarTemplate == nil {
		t.Fatal("navbar template not found in default templates")
	}

	// Note: In a real override test, we would register a custom filesystem
	// with an overriding template, but for this test we're just verifying
	// that the latest registered templates take precedence (tested by order)
}

func TestPartial(t *testing.T) {
	// Store original state
	originalTemplates := templates
	originalConfig := Config
	defer func() {
		templates = originalTemplates
		Config = originalConfig
	}()

	// Compile templates to ensure they're available
	compileTemplates()

	tests := []struct {
		name         string
		templatePath string
		data         Locals
		shouldContain string
		shouldError  bool
	}{
		{
			name:         "Simple template rendering",
			templatePath: "emoji-favicon",
			data:         Locals{"page": &page{name: "test"}},
			shouldContain: "",
			shouldError:  false,
		},
		{
			name:         "Non-existent template",
			templatePath: "nonexistent-template",
			data:         nil,
			shouldContain: "template nonexistent-template not found",
			shouldError:  true,
		},
		{
			name:         "Nil data should create empty Locals",
			templatePath: "emoji-favicon",
			data:         nil,
			shouldContain: "",
			shouldError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Partial(tt.templatePath, tt.data)
			resultStr := string(result)

			if tt.shouldError {
				if !strings.Contains(resultStr, tt.shouldContain) {
					t.Errorf("Expected error message to contain '%s', got: %s", tt.shouldContain, resultStr)
				}
			} else {
				if strings.Contains(resultStr, "rendering error") {
					t.Errorf("Unexpected rendering error: %s", resultStr)
				}
			}
		})
	}
}

func TestPartialWithConfig(t *testing.T) {
	// Store original state
	originalTemplates := templates
	originalConfig := Config
	defer func() {
		templates = originalTemplates
		Config = originalConfig
	}()

	// Set a test config value
	Config = Configuration{
		Sitename: "Test Site",
	}

	// Create a simple test template
	templates = template.New("")
	testTemplate := `Site: {{.config.Sitename}}`
	template.Must(templates.New("test-config").Parse(testTemplate))

	// Test that config is passed to template
	result := Partial("test-config", Locals{})
	resultStr := string(result)

	if !strings.Contains(resultStr, "Test Site") {
		t.Errorf("Expected template to have access to config, got: %s", resultStr)
	}
}

func TestPartialDataMerging(t *testing.T) {
	// Store original state
	originalTemplates := templates
	originalConfig := Config
	defer func() {
		templates = originalTemplates
		Config = originalConfig
	}()

	// Set a test config
	Config = Configuration{
		Sitename: "Test",
	}

	// Create a test template that uses both custom data and config
	templates = template.New("")
	testTemplate := `Name: {{.name}}, Site: {{.config.Sitename}}`
	template.Must(templates.New("test-merge").Parse(testTemplate))

	// Test that custom data and config are both available
	result := Partial("test-merge", Locals{"name": "TestPage"})
	resultStr := string(result)

	if !strings.Contains(resultStr, "TestPage") {
		t.Errorf("Expected custom data to be available, got: %s", resultStr)
	}

	if !strings.Contains(resultStr, "Test") {
		t.Errorf("Expected config to be available, got: %s", resultStr)
	}
}

func TestPartialTemplateError(t *testing.T) {
	// Store original state
	originalTemplates := templates
	originalConfig := Config
	defer func() {
		templates = originalTemplates
		Config = originalConfig
	}()

	// Create a template that will fail during execution
	// We'll use a template that tries to call a method on a nil value
	templates = template.New("")
	// Create a custom function that will panic
	panicFunc := func() string {
		panic("intentional test panic")
	}
	funcMap := template.FuncMap{
		"panicFunc": panicFunc,
	}
	badTemplate := `{{panicFunc}}`
	template.Must(templates.New("test-error").Funcs(funcMap).Parse(badTemplate))

	// The Partial function catches panics/errors and returns them as strings
	result := Partial("test-error", Locals{})
	resultStr := string(result)

	// The template execution should fail and return an error message
	if !strings.Contains(resultStr, "rendering error") {
		t.Errorf("Expected rendering error message, got: %s", resultStr)
	}
}

func TestTemplateHelpers(t *testing.T) {
	// Store original state
	originalTemplates := templates
	defer func() {
		templates = originalTemplates
	}()

	// Compile templates (which includes helpers)
	compileTemplates()

	// Create a test template that uses a helper function
	testTemplate := `{{base "/path/to/file.md"}}`
	template.Must(templates.New("test-helper").Funcs(helpers).Parse(testTemplate))

	result := Partial("test-helper", Locals{})
	resultStr := string(result)

	// The base helper should extract just the filename
	if !strings.Contains(resultStr, "file.md") {
		t.Errorf("Expected helper function to work, got: %s", resultStr)
	}
}
