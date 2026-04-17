package opengraph

import (
	"html/template"
	"strings"
	"testing"
	"time"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestOpengraphExtensionName(t *testing.T) {
	ext := Opengraph{}
	if ext.Name() != "opengraph" {
		t.Errorf("Expected extension name to be 'opengraph', got '%s'", ext.Name())
	}
}

func TestRawText(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		limit    int
		expected string
	}{
		{
			name:     "simple text",
			markdown: "Hello world",
			limit:    100,
			expected: "Hello world",
		},
		{
			name:     "text with limit",
			markdown: "This is a longer text that should be truncated",
			limit:    20,
			expected: "This is a longer te",
		},
		{
			name:     "text with formatting",
			markdown: "This is **bold** and *italic* text",
			limit:    100,
			expected: "This is bold and italic text",
		},
		{
			name:     "multiple paragraphs",
			markdown: "First paragraph.\n\nSecond paragraph.",
			limit:    100,
			expected: "First paragraph . Second paragraph .",
		},
		{
			name:     "nil source",
			markdown: "",
			limit:    100,
			expected: "",
		},
		{
			name:     "text with links",
			markdown: "Check out [this link](https://example.com)",
			limit:    100,
			expected: "Check out this link",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tree ast.Node
			if tt.markdown != "" {
				mockPage := &mockTestPage{content: tt.markdown}
				_, tree = mockPage.AST()
			}

			result := rawText([]byte(tt.markdown), tree, tt.limit)
			if result != tt.expected {
				t.Errorf("rawText() = '%s', want '%s'", result, tt.expected)
			}
		})
	}
}

func TestOpengraphTags(t *testing.T) {
	// Setup test configuration
	originalSitename := Config.Sitename
	originalIndex := Config.Index
	defer func() {
		Config.Sitename = originalSitename
		Config.Index = originalIndex
	}()

	Config.Sitename = "Test Site"
	Config.Index = "index"
	domain = "example.com"
	twitterUsername = "@testuser"

	// Create a mock page
	mockPage := mockTestPage{
		name:    "Test Page",
		content: "# Test Page\n\nThis is a test page with some content.",
	}

	result := opengraphTags(mockPage)
	resultStr := string(result)

	// Check for required OpenGraph tags
	expectedTags := []string{
		`property="og:site_name" content="Test Site"`,
		`property="og:title" content="Test Page"`,
		`property="og:description"`,
		`property="og:type" content="website"`,
		`name="twitter:title" content="Test Page"`,
		`name="twitter:creator" content="@testuser"`,
		`name="twitter:site" content="@testuser"`,
		`name="description"`,
	}

	for _, tag := range expectedTags {
		if !strings.Contains(resultStr, tag) {
			t.Errorf("Expected tag not found: %s", tag)
		}
	}
}

func TestOpengraphTagsWithImage(t *testing.T) {
	Config.Sitename = "Test Site"
	domain = "example.com"
	twitterUsername = "@testuser"

	mockPage := mockTestPage{
		name:    "Page With Image",
		content: "# Page With Image\n\n![alt text](/image.png)\n\nSome content.",
	}

	result := opengraphTags(mockPage)
	resultStr := string(result)

	if !strings.Contains(resultStr, `property="og:image" content="https://example.com/image.png"`) {
		t.Error("Expected og:image tag with correct image URL")
	}

	if !strings.Contains(resultStr, `name="twitter:image" content="https://example.com/image.png"`) {
		t.Error("Expected twitter:image tag with correct image URL")
	}
}

func TestOpengraphTagsForIndexPage(t *testing.T) {
	Config.Sitename = "My Blog"
	Config.Index = "index"
	domain = "myblog.com"

	mockPage := mockTestPage{
		name:    "index",
		content: "# Welcome\n\nThis is the index page.",
	}

	result := opengraphTags(mockPage)
	resultStr := string(result)

	// For index page, title should be sitename
	if !strings.Contains(resultStr, `property="og:title" content="My Blog"`) {
		t.Error("Expected og:title to be sitename for index page")
	}
}

// Mock page implementation for testing
type mockTestPage struct {
	name    string
	content string
}

func (m mockTestPage) Name() string                      { return m.name }
func (m mockTestPage) FileName() string                  { return m.name + ".md" }
func (m mockTestPage) Exists() bool                      { return true }
func (m mockTestPage) Render() template.HTML             { return "" }
func (m mockTestPage) Content() Markdown                 { return Markdown(m.content) }
func (m mockTestPage) Delete() bool                      { return false }
func (m mockTestPage) Write(Markdown) bool               { return false }
func (m mockTestPage) ModTime() time.Time                { return time.Now() }
func (m mockTestPage) AST() ([]byte, ast.Node) {
	source := []byte(m.content)
	tree := MarkdownConverter().Parser().Parse(text.NewReader(source))
	return source, tree
}
