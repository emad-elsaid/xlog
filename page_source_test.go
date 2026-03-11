package xlog

import (
	"context"
	"html/template"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockPageSource implements PageSource for testing
type mockPageSource struct {
	pages map[string]Page
	eachFunc func(context.Context, func(Page))
}

func (m *mockPageSource) Page(name string) Page {
	if m.pages == nil {
		return nil
	}
	return m.pages[name]
}

func (m *mockPageSource) Each(ctx context.Context, fn func(Page)) {
	if m.eachFunc != nil {
		m.eachFunc(ctx, fn)
	}
}

// mockPage implements Page for testing
type mockPage struct {
	name   string
	exists bool
}

func (m *mockPage) Name() string                   { return m.name }
func (m *mockPage) Exists() bool                   { return m.exists }
func (m *mockPage) FileName() string               { return m.name + ".md" }
func (m *mockPage) Content() Markdown              { return "" }
func (m *mockPage) Render() template.HTML          { return "" }
func (m *mockPage) Delete() bool                   { return false }
func (m *mockPage) Write(Markdown) bool            { return false }
func (m *mockPage) ModTime() time.Time             { return time.Time{} }
func (m *mockPage) AST() ([]byte, ast.Node)        { return nil, nil }

func TestNewPage(t *testing.T) {
	// Save original sources and restore after test
	originalSources := sources
	defer func() { sources = originalSources }()

	tests := []struct {
		name           string
		pageName       string
		setupSources   func()
		expectedExists bool
		expectedName   string
	}{
		{
			name:     "page exists in first source",
			pageName: "test-page",
			setupSources: func() {
				sources = []PageSource{
					&mockPageSource{
						pages: map[string]Page{
							"test-page": &mockPage{name: "test-page", exists: true},
						},
					},
				}
			},
			expectedExists: true,
			expectedName:   "test-page",
		},
		{
			name:     "page exists in second source",
			pageName: "second-page",
			setupSources: func() {
				sources = []PageSource{
					&mockPageSource{
						pages: map[string]Page{
							"other-page": &mockPage{name: "other-page", exists: true},
						},
					},
					&mockPageSource{
						pages: map[string]Page{
							"second-page": &mockPage{name: "second-page", exists: true},
						},
					},
				}
			},
			expectedExists: true,
			expectedName:   "second-page",
		},
		{
			name:     "page doesn't exist in any source",
			pageName: "missing-page",
			setupSources: func() {
				sources = []PageSource{
					&mockPageSource{
						pages: map[string]Page{
							"other-page": &mockPage{name: "other-page", exists: true},
						},
					},
				}
			},
			expectedExists: false,
		},
		{
			name:     "page exists but not marked as existing",
			pageName: "non-existing-page",
			setupSources: func() {
				sources = []PageSource{
					&mockPageSource{
						pages: map[string]Page{
							"non-existing-page": &mockPage{name: "non-existing-page", exists: false},
						},
					},
				}
			},
			expectedExists: false,
		},
		{
			name:     "empty sources list",
			pageName: "any-page",
			setupSources: func() {
				sources = []PageSource{}
			},
			expectedExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupSources()
			
			page := NewPage(tt.pageName)
			
			if tt.expectedExists {
				require.NotNil(t, page, "Expected page to exist")
				assert.True(t, page.Exists(), "Expected page.Exists() to return true")
				assert.Equal(t, tt.expectedName, page.Name(), "Page name mismatch")
			} else {
				if page != nil {
					assert.False(t, page.Exists(), "Expected page.Exists() to return false")
				}
			}
		})
	}
}

func TestRegisterPageSource(t *testing.T) {
	// Save original sources and restore after test
	originalSources := sources
	defer func() { sources = originalSources }()

	// Reset sources to a known state
	sources = []PageSource{
		&mockPageSource{
			pages: map[string]Page{
				"existing": &mockPage{name: "existing", exists: true},
			},
		},
	}
	initialLen := len(sources)

	// Register a new page source
	newSource := &mockPageSource{
		pages: map[string]Page{
			"new-page": &mockPage{name: "new-page", exists: true},
		},
	}
	RegisterPageSource(newSource)

	// Verify the source was added
	assert.Equal(t, initialLen+1, len(sources), "Source count should increase by 1")

	// Verify new source is prepended (checked first)
	page := NewPage("new-page")
	require.NotNil(t, page, "Should find page from newly registered source")
	assert.True(t, page.Exists())
	assert.Equal(t, "new-page", page.Name())

	// Verify old sources still work
	page = NewPage("existing")
	require.NotNil(t, page, "Should still find page from original source")
	assert.True(t, page.Exists())
	assert.Equal(t, "existing", page.Name())
}

func TestRegisterPageSourcePriority(t *testing.T) {
	// Save original sources and restore after test
	originalSources := sources
	defer func() { sources = originalSources }()

	// Create two sources with the same page name
	sources = []PageSource{
		&mockPageSource{
			pages: map[string]Page{
				"conflict": &mockPage{name: "conflict-old", exists: true},
			},
		},
	}

	// Register a new source with the same page name
	newSource := &mockPageSource{
		pages: map[string]Page{
			"conflict": &mockPage{name: "conflict-new", exists: true},
		},
	}
	RegisterPageSource(newSource)

	// The newly registered source should take priority
	page := NewPage("conflict")
	require.NotNil(t, page)
	assert.Equal(t, "conflict-new", page.Name(), "Newly registered source should take priority")
}

func TestPageSourceIntegration(t *testing.T) {
	// This test verifies that markdownFS can be used as a PageSource
	// We use the current directory which should have some .md files
	
	// Save original sources and restore after test
	originalSources := sources
	defer func() { sources = originalSources }()

	// Create a markdownFS for current directory
	sources = []PageSource{
		newMarkdownFS("."),
	}

	// Test that NewPage returns a non-nil page even for non-existent pages
	// This is how the actual implementation works
	page := NewPage("nonexistent-test-page-12345")
	require.NotNil(t, page, "NewPage should always return a page object")
	assert.Equal(t, "nonexistent-test-page-12345", page.Name())
	
	// Since the page doesn't exist, Exists() should return false
	assert.False(t, page.Exists(), "Page should not exist")
}
