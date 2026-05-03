package xlog

import (
	"html/template"
	"strings"
	"testing"
)

func TestRegisterHelper_Success(t *testing.T) {
	// Save original helpers and restore
	originalHelpers := helpers
	defer func() { helpers = originalHelpers }()
	
	// Reset to known state
	helpers = template.FuncMap{}
	
	// Register a new helper
	testFunc := func() string { return "test" }
	err := RegisterHelper("testHelper", testFunc)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	// Verify it was registered
	if _, ok := helpers["testHelper"]; !ok {
		t.Error("Helper was not registered")
	}
}

func TestRegisterHelper_Duplicate(t *testing.T) {
	// Save original helpers and restore
	originalHelpers := helpers
	defer func() { helpers = originalHelpers }()
	
	// Reset to known state
	helpers = template.FuncMap{
		"existing": func() string { return "exists" },
	}
	
	// Try to register duplicate
	err := RegisterHelper("existing", func() string { return "new" })
	
	if err != ErrHelperRegistered {
		t.Errorf("Expected ErrHelperRegistered, got %v", err)
	}
	
	// Verify original was not replaced
	fn := helpers["existing"].(func() string)
	if fn() != "exists" {
		t.Error("Original helper was replaced")
	}
}

func TestRegisterJS_NewLibrary(t *testing.T) {
	// Save original js slice and restore
	originalJS := js
	defer func() { js = originalJS }()
	
	// Reset to empty
	js = []string{}
	
	// Register a library
	RegisterJS("/public/test.js")
	
	if len(js) != 1 {
		t.Fatalf("Expected 1 library, got %d", len(js))
	}
	
	if js[0] != "/public/test.js" {
		t.Errorf("Expected '/public/test.js', got '%s'", js[0])
	}
}

func TestRegisterJS_Duplicate(t *testing.T) {
	// Save original js slice and restore
	originalJS := js
	defer func() { js = originalJS }()
	
	// Reset to empty
	js = []string{}
	
	// Register same library twice
	RegisterJS("/public/test.js")
	RegisterJS("/public/test.js")
	
	if len(js) != 1 {
		t.Errorf("Expected 1 library (no duplicates), got %d", len(js))
	}
}

func TestRegisterJS_Multiple(t *testing.T) {
	// Save original js slice and restore
	originalJS := js
	defer func() { js = originalJS }()
	
	// Reset to empty
	js = []string{}
	
	// Register multiple libraries
	RegisterJS("/public/lib1.js")
	RegisterJS("/public/lib2.js")
	RegisterJS("/public/lib3.js")
	
	if len(js) != 3 {
		t.Fatalf("Expected 3 libraries, got %d", len(js))
	}
	
	expected := []string{
		"/public/lib1.js",
		"/public/lib2.js",
		"/public/lib3.js",
	}
	
	for i, lib := range expected {
		if js[i] != lib {
			t.Errorf("Expected js[%d] = '%s', got '%s'", i, lib, js[i])
		}
	}
}

func TestRequireHTMX(t *testing.T) {
	// Save original js slice and restore
	originalJS := js
	defer func() { js = originalJS }()
	
	// Reset to empty
	js = []string{}
	
	// Call RequireHTMX
	RequireHTMX()
	
	if len(js) != 1 {
		t.Fatalf("Expected 1 library (HTMX), got %d", len(js))
	}
	
	if js[0] != "/public/htmx.min.js" {
		t.Errorf("Expected HTMX path '/public/htmx.min.js', got '%s'", js[0])
	}
}

func TestRequireHTMX_Multiple(t *testing.T) {
	// Save original js slice and restore
	originalJS := js
	defer func() { js = originalJS }()
	
	// Reset to empty
	js = []string{}
	
	// Call RequireHTMX multiple times
	RequireHTMX()
	RequireHTMX()
	RequireHTMX()
	
	// Should only register once (uses RegisterJS which prevents duplicates)
	if len(js) != 1 {
		t.Errorf("Expected 1 library (no duplicates), got %d", len(js))
	}
}

func TestIncludeJS(t *testing.T) {
	// Save original js slice and restore
	originalJS := js
	defer func() { js = originalJS }()
	
	// Reset to empty
	js = []string{}
	
	// Call includeJS
	result := includeJS("/public/custom.js")
	
	// Should register the JS
	if len(js) != 1 {
		t.Errorf("Expected JS to be registered")
	}
	
	// Should return empty HTML (registration happens, no output)
	if result != "" {
		t.Errorf("Expected empty result, got '%s'", result)
	}
}

func TestScripts_Empty(t *testing.T) {
	// Save original js slice and restore
	originalJS := js
	defer func() { js = originalJS }()
	
	// Reset to empty
	js = []string{}
	
	result := scripts()
	
	if result != "" {
		t.Errorf("Expected empty scripts output, got '%s'", result)
	}
}

func TestScripts_Single(t *testing.T) {
	// Save original js slice and restore
	originalJS := js
	defer func() { js = originalJS }()
	
	// Reset with one script
	js = []string{"/public/test.js"}
	
	result := scripts()
	expected := `<script src="/public/test.js" defer></script>`
	
	if string(result) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestScripts_Multiple(t *testing.T) {
	// Save original js slice and restore
	originalJS := js
	defer func() { js = originalJS }()
	
	// Reset with multiple scripts
	js = []string{
		"/public/lib1.js",
		"/public/lib2.js",
		"/public/lib3.js",
	}
	
	result := string(scripts())
	
	// Should contain all three scripts
	if !strings.Contains(result, `<script src="/public/lib1.js" defer></script>`) {
		t.Error("Missing lib1.js in output")
	}
	if !strings.Contains(result, `<script src="/public/lib2.js" defer></script>`) {
		t.Error("Missing lib2.js in output")
	}
	if !strings.Contains(result, `<script src="/public/lib3.js" defer></script>`) {
		t.Error("Missing lib3.js in output")
	}
}

func TestIsFontAwesome_True(t *testing.T) {
	tests := []string{
		"fa",
		"fa-home",
		"fa-user",
		"fas fa-check",
		"far fa-circle",
	}
	
	for _, icon := range tests {
		if !IsFontAwesome(icon) {
			t.Errorf("Expected IsFontAwesome('%s') to be true", icon)
		}
	}
}

func TestIsFontAwesome_False(t *testing.T) {
	tests := []string{
		"",
		"icon-home",
		"material-icons",
		"bootstrap-icon",
		" fa-home", // Space prefix
	}
	
	for _, icon := range tests {
		if IsFontAwesome(icon) {
			t.Errorf("Expected IsFontAwesome('%s') to be false", icon)
		}
	}
}

func TestRaw(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected template.HTML
	}{
		{
			name:     "plain text",
			input:    "Hello World",
			expected: template.HTML("Hello World"),
		},
		{
			name:     "HTML tags",
			input:    "<strong>Bold</strong>",
			expected: template.HTML("<strong>Bold</strong>"),
		},
		{
			name:     "script tag",
			input:    "<script>alert('test')</script>",
			expected: template.HTML("<script>alert('test')</script>"),
		},
		{
			name:     "empty string",
			input:    "",
			expected: template.HTML(""),
		},
		{
			name:     "special characters",
			input:    "<div>&lt;&gt;&amp;</div>",
			expected: template.HTML("<div>&lt;&gt;&amp;</div>"),
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := raw(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
