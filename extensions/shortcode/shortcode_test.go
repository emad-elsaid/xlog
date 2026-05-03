package shortcode

import (
	"html/template"
	"strings"
	"testing"

	"github.com/emad-elsaid/xlog"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		input    xlog.Markdown
		contains []string
	}{
		{
			name:     "plain text markdown",
			input:    xlog.Markdown("Hello World"),
			contains: []string{"Hello World"},
		},
		{
			name:     "markdown with bold",
			input:    xlog.Markdown("**bold text**"),
			contains: []string{"<strong>", "bold text", "</strong>"},
		},
		{
			name:     "markdown with italic",
			input:    xlog.Markdown("*italic text*"),
			contains: []string{"<em>", "italic text", "</em>"},
		},
		{
			name:     "markdown with link",
			input:    xlog.Markdown("[link](https://example.com)"),
			contains: []string{"<a", "href=\"https://example.com\"", "link", "</a>"},
		},
		{
			name:     "markdown with code",
			input:    xlog.Markdown("`code`"),
			contains: []string{"<code>", "code", "</code>"},
		},
		{
			name:     "multiline markdown",
			input:    xlog.Markdown("line 1\nline 2\nline 3"),
			contains: []string{"line 1", "line 2", "line 3"},
		},
		{
			name:     "empty markdown",
			input:    xlog.Markdown(""),
			contains: []string{},
		},
		{
			name:     "markdown with heading",
			input:    xlog.Markdown("# Heading"),
			contains: []string{"<h1", "Heading", "</h1>"},
		},
		{
			name:     "markdown with list",
			input:    xlog.Markdown("- item 1\n- item 2"),
			contains: []string{"<ul>", "<li>", "item 1", "item 2", "</li>", "</ul>"},
		},
		{
			name:     "special characters",
			input:    xlog.Markdown("< > & \" '"),
			contains: []string{"&lt;", "&gt;", "&amp;"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := render(tc.input)

			for _, expected := range tc.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("render(%q) result missing %q\nGot: %s", tc.input, expected, result)
				}
			}
		})
	}
}

func TestContainer(t *testing.T) {
	tests := []struct {
		name       string
		class      string
		content    xlog.Markdown
		wantClass  string
		wantHTML   []string
		notContain []string
	}{
		{
			name:      "info container with plain text",
			class:     "is-info",
			content:   xlog.Markdown("Information message"),
			wantClass: "is-info",
			wantHTML:  []string{"<article", "class=\"message is-info\"", "<div", "class=\"message-body\"", "Information message", "</div>", "</article>"},
		},
		{
			name:      "success container with markdown",
			class:     "is-success",
			content:   xlog.Markdown("**Success!**"),
			wantClass: "is-success",
			wantHTML:  []string{"<article", "class=\"message is-success\"", "<strong>", "Success!", "</strong>"},
		},
		{
			name:      "warning container with multiple lines",
			class:     "is-warning",
			content:   xlog.Markdown("Warning line 1\nWarning line 2"),
			wantClass: "is-warning",
			wantHTML:  []string{"class=\"message is-warning\"", "Warning line 1", "Warning line 2"},
		},
		{
			name:      "danger container with link",
			class:     "is-danger",
			content:   xlog.Markdown("[Alert](https://alert.com)"),
			wantClass: "is-danger",
			wantHTML:  []string{"class=\"message is-danger\"", "<a", "href=\"https://alert.com\"", "Alert"},
		},
		{
			name:       "empty content",
			class:      "is-info",
			content:    xlog.Markdown(""),
			wantClass:  "is-info",
			wantHTML:   []string{"<article", "class=\"message is-info\"", "<div", "class=\"message-body\""},
			notContain: []string{},
		},
		{
			name:      "custom class",
			class:     "custom-class",
			content:   xlog.Markdown("Custom"),
			wantClass: "custom-class",
			wantHTML:  []string{"class=\"message custom-class\"", "Custom"},
		},
		{
			name:      "content with special characters",
			class:     "is-info",
			content:   xlog.Markdown("< > & \" '"),
			wantClass: "is-info",
			wantHTML:  []string{"&lt;", "&gt;", "&amp;"},
		},
		{
			name:      "content with code",
			class:     "is-warning",
			content:   xlog.Markdown("`code block`"),
			wantClass: "is-warning",
			wantHTML:  []string{"class=\"message is-warning\"", "<code>", "code block", "</code>"},
		},
		{
			name:      "content with heading",
			class:     "is-success",
			content:   xlog.Markdown("## Subheading"),
			wantClass: "is-success",
			wantHTML:  []string{"class=\"message is-success\"", "<h2", "Subheading", "</h2>"},
		},
		{
			name:      "multiline with list",
			class:     "is-danger",
			content:   xlog.Markdown("Items:\n- One\n- Two"),
			wantClass: "is-danger",
			wantHTML:  []string{"class=\"message is-danger\"", "Items:", "<ul>", "<li>", "One", "Two"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := container(tc.class, tc.content)
			resultStr := string(result)

			// Verify it returns template.HTML type
			if _, ok := interface{}(result).(template.HTML); !ok {
				t.Errorf("container() should return template.HTML, got %T", result)
			}

			// Verify class is present
			if !strings.Contains(resultStr, tc.wantClass) {
				t.Errorf("container() result missing class %q\nGot: %s", tc.wantClass, resultStr)
			}

			// Verify all expected HTML fragments
			for _, expected := range tc.wantHTML {
				if !strings.Contains(resultStr, expected) {
					t.Errorf("container() result missing %q\nGot: %s", expected, resultStr)
				}
			}

			// Verify forbidden content is not present
			for _, forbidden := range tc.notContain {
				if strings.Contains(resultStr, forbidden) {
					t.Errorf("container() result should not contain %q\nGot: %s", forbidden, resultStr)
				}
			}
		})
	}
}

func TestContainer_StructureValidity(t *testing.T) {
	// Verify the HTML structure is always valid
	result := container("is-info", xlog.Markdown("test"))
	resultStr := string(result)

	// Must contain opening and closing tags in correct order
	articleStart := strings.Index(resultStr, "<article")
	articleEnd := strings.Index(resultStr, "</article>")
	divStart := strings.Index(resultStr, "<div")
	divEnd := strings.Index(resultStr, "</div>")

	if articleStart == -1 {
		t.Error("container() missing opening <article> tag")
	}
	if articleEnd == -1 {
		t.Error("container() missing closing </article> tag")
	}
	if divStart == -1 {
		t.Error("container() missing opening <div> tag")
	}
	if divEnd == -1 {
		t.Error("container() missing closing </div> tag")
	}

	// Verify nesting order
	if articleStart >= divStart || divStart >= divEnd || divEnd >= articleEnd {
		t.Errorf("container() has invalid tag nesting order\nGot: %s", resultStr)
	}
}

func TestRegisterShortCode(t *testing.T) {
	tests := []struct {
		name          string
		shortcodeName string
		shortcode     ShortCode
	}{
		{
			name:          "register new shortcode",
			shortcodeName: "testcode",
			shortcode: ShortCode{
				Render:  func(m xlog.Markdown) template.HTML { return template.HTML("test") },
				Default: "default",
			},
		},
		{
			name:          "register with empty default",
			shortcodeName: "emptytest",
			shortcode: ShortCode{
				Render:  func(m xlog.Markdown) template.HTML { return template.HTML("") },
				Default: "",
			},
		},
		{
			name:          "register overwriting existing",
			shortcodeName: "testoverwrite",
			shortcode: ShortCode{
				Render:  func(m xlog.Markdown) template.HTML { return template.HTML("new") },
				Default: "new default",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Should not panic
			RegisterShortCode(tc.shortcodeName, tc.shortcode)

			// Verify registration by checking internal map
			shortcodesMutex.RLock()
			registered, exists := shortcodes[tc.shortcodeName]
			shortcodesMutex.RUnlock()

			if !exists {
				t.Errorf("RegisterShortCode(%q) did not register shortcode", tc.shortcodeName)
			}

			if registered.Default != tc.shortcode.Default {
				t.Errorf("RegisterShortCode(%q) Default = %q, want %q",
					tc.shortcodeName, registered.Default, tc.shortcode.Default)
			}
		})
	}
}

func TestRegisterShortCode_Concurrency(t *testing.T) {
	// Test thread-safety of RegisterShortCode
	const goroutines = 100
	done := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			name := "concurrent"
			RegisterShortCode(name, ShortCode{
				Render:  func(m xlog.Markdown) template.HTML { return template.HTML("test") },
				Default: "",
			})
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < goroutines; i++ {
		<-done
	}

	// Should not panic and shortcode should be registered
	shortcodesMutex.RLock()
	_, exists := shortcodes["concurrent"]
	shortcodesMutex.RUnlock()

	if !exists {
		t.Error("Concurrent RegisterShortCode failed to register")
	}
}

func TestBuiltinShortcodes(t *testing.T) {
	// Verify built-in shortcodes are registered
	builtins := []string{"info", "success", "warning", "alert"}

	for _, name := range builtins {
		t.Run(name, func(t *testing.T) {
			shortcodesMutex.RLock()
			sc, exists := shortcodes[name]
			shortcodesMutex.RUnlock()

			if !exists {
				t.Errorf("Built-in shortcode %q is not registered", name)
			}

			if sc.Render == nil {
				t.Errorf("Built-in shortcode %q has nil Render function", name)
			}

			// Test rendering
			result := sc.Render(xlog.Markdown("test content"))
			resultStr := string(result)

			// All built-ins should produce article containers
			if !strings.Contains(resultStr, "<article") {
				t.Errorf("Built-in shortcode %q does not produce article container", name)
			}

			if !strings.Contains(resultStr, "message-body") {
				t.Errorf("Built-in shortcode %q missing message-body class", name)
			}

			if !strings.Contains(resultStr, "test content") {
				t.Errorf("Built-in shortcode %q did not render content", name)
			}
		})
	}
}

func TestBuiltinShortcodes_ClassMapping(t *testing.T) {
	// Verify correct CSS class for each built-in
	tests := []struct {
		name      string
		wantClass string
	}{
		{"info", "is-info"},
		{"success", "is-success"},
		{"warning", "is-warning"},
		{"alert", "is-danger"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			shortcodesMutex.RLock()
			sc := shortcodes[tc.name]
			shortcodesMutex.RUnlock()

			result := sc.Render(xlog.Markdown("content"))
			resultStr := string(result)

			if !strings.Contains(resultStr, tc.wantClass) {
				t.Errorf("Shortcode %q should have class %q\nGot: %s", tc.name, tc.wantClass, resultStr)
			}
		})
	}
}
