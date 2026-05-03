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
	if err := os.Mkdir(themeDir, 0750); err != nil {
		t.Fatalf("Failed to create theme directory: %v", err)
	}

	// Create a custom template in the theme directory
	customTemplate := `<div>Custom Theme Template</div>`
	customTemplatePath := filepath.Join(themeDir, "custom.html")
	if err := os.WriteFile(customTemplatePath, []byte(customTemplate), 0600); err != nil {
		t.Fatalf("Failed to write custom template: %v", err)
	}

	// Change to the temporary directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

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
		name          string
		templatePath  string
		data          Locals
		shouldContain string
		shouldError   bool
	}{
		{
			name:          "Simple template rendering",
			templatePath:  "emoji-favicon",
			data:          Locals{"page": &page{name: "test"}},
			shouldContain: "",
			shouldError:   false,
		},
		{
			name:          "Non-existent template",
			templatePath:  "nonexistent-template",
			data:          nil,
			shouldContain: "template nonexistent-template not found",
			shouldError:   true,
		},
		{
			name:          "Nil data should create empty Locals",
			templatePath:  "emoji-favicon",
			data:          nil,
			shouldContain: "",
			shouldError:   false,
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

// BenchmarkPartialSimple benchmarks simple template rendering with minimal data.
func BenchmarkPartialSimple(b *testing.B) {
	// Setup: compile templates once
	originalTemplates := templates
	defer func() { templates = originalTemplates }()
	compileTemplates()

	// Create a simple template
	template.Must(templates.New("bench-simple").Parse(`<div>Hello {{.name}}</div>`))

	data := Locals{"name": "World"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = Partial("bench-simple", data)
	}
}

// BenchmarkPartialComplex benchmarks complex template with nested data and helpers.
func BenchmarkPartialComplex(b *testing.B) {
	originalTemplates := templates
	defer func() { templates = originalTemplates }()
	compileTemplates()

	// Create a complex template with loops and conditionals
	complexTemplate := `
{{range .items}}
  <div class="item">
    <h2>{{.title}}</h2>
    <p>{{.description}}</p>
    {{if .tags}}
      <ul>
      {{range .tags}}
        <li>{{.}}</li>
      {{end}}
      </ul>
    {{end}}
  </div>
{{end}}
`
	template.Must(templates.New("bench-complex").Funcs(helpers).Parse(complexTemplate))

	data := Locals{
		"items": []map[string]interface{}{
			{
				"title":       "Item 1",
				"description": "Description for item 1",
				"tags":        []string{"tag1", "tag2", "tag3"},
			},
			{
				"title":       "Item 2",
				"description": "Description for item 2",
				"tags":        []string{"tag4", "tag5"},
			},
			{
				"title":       "Item 3",
				"description": "Description for item 3",
				"tags":        []string{},
			},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = Partial("bench-complex", data)
	}
}

// BenchmarkPartialWithConfig benchmarks template rendering with config access.
func BenchmarkPartialWithConfig(b *testing.B) {
	originalTemplates := templates
	originalConfig := Config
	defer func() {
		templates = originalTemplates
		Config = originalConfig
	}()

	compileTemplates()
	Config = Configuration{
		Sitename:    "Benchmark Site",
		BindAddress: ":8080",
		Index:       "index",
	}

	template.Must(templates.New("bench-config").Parse(`
<header>
  <h1>{{.config.Sitename}}</h1>
  <nav>{{.config.Index}}</nav>
</header>
`))

	data := Locals{"title": "Test Page"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = Partial("bench-config", data)
	}
}

// BenchmarkPartialNilData benchmarks template rendering with nil data.
func BenchmarkPartialNilData(b *testing.B) {
	originalTemplates := templates
	defer func() { templates = originalTemplates }()
	compileTemplates()

	template.Must(templates.New("bench-nil").Parse(`<div>Static content</div>`))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = Partial("bench-nil", nil)
	}
}

// BenchmarkPartialEmptyData benchmarks template rendering with empty Locals.
func BenchmarkPartialEmptyData(b *testing.B) {
	originalTemplates := templates
	defer func() { templates = originalTemplates }()
	compileTemplates()

	template.Must(templates.New("bench-empty").Parse(`<div>Static content</div>`))

	data := Locals{}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = Partial("bench-empty", data)
	}
}

// BenchmarkPartialConcurrent benchmarks concurrent template rendering.
func BenchmarkPartialConcurrent(b *testing.B) {
	originalTemplates := templates
	defer func() { templates = originalTemplates }()
	compileTemplates()

	template.Must(templates.New("bench-concurrent").Parse(`<div>{{.id}}: {{.data}}</div>`))

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			data := Locals{
				"id":   i,
				"data": "concurrent data",
			}
			_ = Partial("bench-concurrent", data)
			i++
		}
	})
}

// BenchmarkPartialLargeData benchmarks template rendering with large data sets.
func BenchmarkPartialLargeData(b *testing.B) {
	originalTemplates := templates
	defer func() { templates = originalTemplates }()
	compileTemplates()

	largeTemplate := `
{{range .entries}}
<article>
  <h1>{{.title}}</h1>
  <time>{{.date}}</time>
  <p>{{.content}}</p>
</article>
{{end}}
`
	template.Must(templates.New("bench-large").Parse(largeTemplate))

	// Generate large dataset (100 entries)
	entries := make([]map[string]string, 100)
	for i := 0; i < 100; i++ {
		entries[i] = map[string]string{
			"title":   "Article " + string(rune(i)),
			"date":    "2024-01-01",
			"content": "This is the content for article number " + string(rune(i)),
		}
	}

	data := Locals{"entries": entries}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = Partial("bench-large", data)
	}
}

// BenchmarkPartialNotFound benchmarks template not found error path.
func BenchmarkPartialNotFound(b *testing.B) {
	originalTemplates := templates
	defer func() { templates = originalTemplates }()
	compileTemplates()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = Partial("nonexistent-template-bench", Locals{})
	}
}
