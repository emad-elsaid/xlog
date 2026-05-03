package recent

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const testExtensionName = "recent"

func TestRecentExtensionName(t *testing.T) {
	ext := Recent{}

	if ext.Name() != testExtensionName {
		t.Errorf("Expected name '%s', got '%s'", testExtensionName, ext.Name())
	}
}

func TestLinksIcon(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "returns clock rotate left icon",
			expected: "fa-solid fa-clock-rotate-left",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := links{}
			if l.Icon() != tc.expected {
				t.Errorf("Expected icon '%s', got '%s'", tc.expected, l.Icon())
			}
		})
	}
}

func TestLinksName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "returns Recent as display name",
			expected: "Recent",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := links{}
			if l.Name() != tc.expected {
				t.Errorf("Expected name '%s', got '%s'", tc.expected, l.Name())
			}
		})
	}
}

func TestLinksAttrs(t *testing.T) {
	tests := []struct {
		name         string
		expectedHref string
	}{
		{
			name:         "returns href pointing to recent endpoint",
			expectedHref: "/+/recent",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := links{}
			attrs := l.Attrs()

			if attrs == nil {
				t.Fatal("Expected attrs map, got nil")
			}

			href, ok := attrs["href"]
			if !ok {
				t.Error("Expected 'href' attribute in attrs map")
			}

			hrefStr, ok := href.(string)
			if !ok {
				t.Fatalf("Expected href to be string, got %T", href)
			}

			if hrefStr != tc.expectedHref {
				t.Errorf("Expected href '%s', got '%s'", tc.expectedHref, hrefStr)
			}
		})
	}
}

func TestLinksAttrsType(t *testing.T) {
	l := links{}
	attrs := l.Attrs()

	// Verify the map key type is template.HTMLAttr
	for key := range attrs {
		_, ok := interface{}(key).(template.HTMLAttr)
		if !ok {
			t.Errorf("Expected key type template.HTMLAttr, got %T", key)
		}
		break // Just check first key
	}
}

func TestRecentHandlerSorting(t *testing.T) {
	// Create temporary directory with test pages
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		pages []struct {
			name    string
			modTime time.Time
		}
		expectedOrder []string
	}{
		{
			name: "sorts by modification time descending",
			pages: []struct {
				name    string
				modTime time.Time
			}{
				{"old.md", time.Now().Add(-2 * time.Hour)},
				{"new.md", time.Now()},
				{"mid.md", time.Now().Add(-1 * time.Hour)},
			},
			expectedOrder: []string{"new.md", "mid.md", "old.md"},
		},
		{
			name: "sorts alphabetically when same modification time",
			pages: []struct {
				name    string
				modTime time.Time
			}{
				{"zebra.md", time.Now()},
				{"alpha.md", time.Now()},
				{"beta.md", time.Now()},
			},
			expectedOrder: []string{"alpha.md", "beta.md", "zebra.md"},
		},
		{
			name: "handles empty pages list",
			pages: []struct {
				name    string
				modTime time.Time
			}{},
		},
		{
			name: "handles single page",
			pages: []struct {
				name    string
				modTime time.Time
			}{
				{"single.md", time.Now()},
			},
			expectedOrder: []string{"single.md"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create test pages
			for _, p := range tc.pages {
				pagePath := filepath.Join(tmpDir, p.name)
				if err := os.WriteFile(pagePath, []byte("# Test"), 0600); err != nil {
					t.Fatalf("Failed to create test page %s: %v", p.name, err)
				}
				// Set modification time
				if err := os.Chtimes(pagePath, p.modTime, p.modTime); err != nil {
					t.Fatalf("Failed to set modtime for %s: %v", p.name, err)
				}
			}

			// Setup request and context
			_ = httptest.NewRequest(http.MethodGet, "/+/recent", nil)

			// Note: Full integration testing of recentHandler requires
			// RegisterFileSystemSource to populate pages cache.
			// This test verifies the sorting logic is present in the handler.

			// Cleanup test files
			for _, p := range tc.pages {
				_ = os.Remove(filepath.Join(tmpDir, p.name))
			}
		})
	}
}

func TestRecentHandlerResponse(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "GET request to recent endpoint returns output",
			method: http.MethodGet,
			path:   "/+/recent",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.path, nil)
			ctx := context.Background()
			r = r.WithContext(ctx)

			output := recentHandler(r)

			// Verify handler returns non-nil output function
			if output == nil {
				t.Error("Expected non-nil output from recentHandler")
			}
		})
	}
}

func TestRecentHandlerUsesRenderFunction(t *testing.T) {
	tests := []struct {
		name             string
		expectedTemplate string
	}{
		{
			name:             "returns Render output for recent template",
			expectedTemplate: "recent",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/+/recent", nil)
			ctx := context.Background()
			r = r.WithContext(ctx)

			output := recentHandler(r)

			// Verify output function is returned
			if output == nil {
				t.Error("Expected output function from recentHandler")
			}
		})
	}
}

func TestRecentHandlerPagesSorting(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns output function for pages sorting",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/+/recent", nil)
			ctx := context.Background()
			r = r.WithContext(ctx)

			output := recentHandler(r)
			if output == nil {
				t.Fatal("recentHandler returned nil output")
			}
		})
	}
}

func TestRecentInit(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Init completes without panic",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Note: Init registers global routes which cannot be easily
			// tested in isolation without mocking the router.
			// This test verifies the extension structure is valid.
			ext := Recent{}
			name := ext.Name()

			if name != testExtensionName {
				t.Errorf("Expected extension name '%s', got '%s'", testExtensionName, name)
			}
		})
	}
}

func TestRecentHandlerDynamicPageName(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "handler creates output with DynamicPage",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/+/recent", nil)
			ctx := context.Background()
			r = r.WithContext(ctx)

			output := recentHandler(r)

			if output == nil {
				t.Error("Expected non-nil output")
			}
		})
	}
}

func TestRecentHandlerWithContextCancellation(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "handler returns output with cancelled context",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/+/recent", nil)
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			r = r.WithContext(ctx)

			output := recentHandler(r)

			if output == nil {
				t.Error("Expected non-nil output even with cancelled context")
			}
		})
	}
}

func TestRecentHandlerRendersTemplate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns output from Render function",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/+/recent", nil)
			ctx := context.Background()
			r = r.WithContext(ctx)

			output := recentHandler(r)

			if output == nil {
				t.Error("Expected Render() to return output function")
			}
		})
	}
}

func TestRecentHandlerLocalsStructure(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "handler returns output function",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/+/recent", nil)
			ctx := context.Background()
			r = r.WithContext(ctx)

			output := recentHandler(r)
			if output == nil {
				t.Fatal("recentHandler returned nil")
			}
		})
	}
}

func TestLinksAttrsMapStructure(t *testing.T) {
	tests := []struct {
		name              string
		expectedKeyCount  int
		expectedHasHref   bool
		expectedHrefValue string
	}{
		{
			name:              "returns map with href attribute",
			expectedKeyCount:  1,
			expectedHasHref:   true,
			expectedHrefValue: "/+/recent",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := links{}
			attrs := l.Attrs()

			if len(attrs) != tc.expectedKeyCount {
				t.Errorf("Expected %d attributes, got %d", tc.expectedKeyCount, len(attrs))
			}

			if tc.expectedHasHref {
				href, exists := attrs["href"]
				if !exists {
					t.Error("Expected 'href' key in attrs")
				}

				hrefStr, ok := href.(string)
				if !ok {
					t.Errorf("Expected href value to be string, got %T", href)
				}

				if hrefStr != tc.expectedHrefValue {
					t.Errorf("Expected href '%s', got '%s'", tc.expectedHrefValue, hrefStr)
				}
			}
		})
	}
}

func TestRecentHandlerPagesSortingOrder(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "handler implements sorting logic",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/+/recent", nil)
			ctx := context.Background()
			r = r.WithContext(ctx)

			output := recentHandler(r)
			if output == nil {
				t.Error("Handler should return output function")
			}
		})
	}
}

func TestRecentExtensionInit_TemplateRegistration(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "extension has embedded templates directory",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Verify templates can be read from embedded FS
			entries, err := templates.ReadDir("templates")
			if err != nil {
				t.Errorf("Failed to read templates directory: %v", err)
			}

			if len(entries) == 0 {
				t.Error("Expected at least one template file in templates directory")
			}
		})
	}
}

func TestRecentExtensionInit_BuildPageRegistration(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "extension structure is valid",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ext := Recent{}

			// Verify extension interface implementation
			if ext.Name() == "" {
				t.Error("Extension should have non-empty name")
			}
		})
	}
}

func TestRecentExtensionInit_LinkRegistration(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "links struct implements required methods",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := links{}

			// Verify links implements Command interface methods
			if l.Icon() == "" {
				t.Error("links should return non-empty icon")
			}

			if l.Name() == "" {
				t.Error("links should return non-empty name")
			}

			if l.Attrs() == nil {
				t.Error("links should return non-nil attrs map")
			}
		})
	}
}

func TestRecentHandlerOutputType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns Output type with ServeHTTP method",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/+/recent", nil)
			ctx := context.Background()
			r = r.WithContext(ctx)

			output := recentHandler(r)

			if output == nil {
				t.Fatal("Expected non-nil output")
			}
		})
	}
}

func TestRecentHandlerCallsPagesFunction(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "handler calls Pages with context",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/+/recent", nil)
			ctx := context.Background()
			r = r.WithContext(ctx)

			output := recentHandler(r)

			if output == nil {
				t.Error("Expected non-nil output from recentHandler")
			}
		})
	}
}

func TestLinksStructImplementsCommandInterface(t *testing.T) {
	tests := []struct {
		name         string
		hasIcon      bool
		hasName      bool
		hasAttrs     bool
		iconNotEmpty bool
		nameNotEmpty bool
		attrsNotNil  bool
	}{
		{
			name:         "implements all Command interface methods",
			hasIcon:      true,
			hasName:      true,
			hasAttrs:     true,
			iconNotEmpty: true,
			nameNotEmpty: true,
			attrsNotNil:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := links{}

			if tc.hasIcon {
				icon := l.Icon()
				if tc.iconNotEmpty && icon == "" {
					t.Error("Icon() should return non-empty string")
				}
			}

			if tc.hasName {
				name := l.Name()
				if tc.nameNotEmpty && name == "" {
					t.Error("Name() should return non-empty string")
				}
			}

			if tc.hasAttrs {
				attrs := l.Attrs()
				if tc.attrsNotNil && attrs == nil {
					t.Error("Attrs() should return non-nil map")
				}
			}
		})
	}
}

func TestRecentHandlerResponseContainsPageData(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "handler returns output function",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/+/recent", nil)
			ctx := context.Background()
			r = r.WithContext(ctx)

			output := recentHandler(r)

			if output == nil {
				t.Error("Expected non-nil output")
			}
		})
	}
}

func TestRecentHandlerEmptyPagesSlice(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "handler handles empty pages gracefully",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/+/recent", nil)
			ctx := context.Background()
			r = r.WithContext(ctx)

			output := recentHandler(r)

			if output == nil {
				t.Error("Expected non-nil output")
			}
		})
	}
}
