package autolink_pages

import (
	"html/template"
	"sort"
	"strings"
	"testing"
	"time"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/text"
)

// TestNormalizedPageSorting tests that pages are sorted by name length (descending)
func TestNormalizedPageSorting(t *testing.T) {
	pages := []*NormalizedPage{
		{normalizedName: "a"},
		{normalizedName: "very-long-name"},
		{normalizedName: "medium"},
		{normalizedName: "short"},
	}

	sort.Sort(fileInfoByNameLength(pages))

	// After sorting, longest should be first
	if pages[0].normalizedName != "very-long-name" {
		t.Errorf("Expected 'very-long-name' first, got '%s'", pages[0].normalizedName)
	}

	// Verify descending order
	for i := 0; i < len(pages)-1; i++ {
		if len(pages[i].normalizedName) < len(pages[i+1].normalizedName) {
			t.Errorf("Pages not sorted by length: %s (%d) should come before %s (%d)",
				pages[i].normalizedName, len(pages[i].normalizedName),
				pages[i+1].normalizedName, len(pages[i+1].normalizedName))
		}
	}
}

// TestFileInfoByNameLength tests the sort interface implementation
func TestFileInfoByNameLength(t *testing.T) {
	pages := []*NormalizedPage{
		{normalizedName: "short"},
		{normalizedName: "very-long-name"},
		{normalizedName: "medium-one"},
	}

	list := fileInfoByNameLength(pages)

	// Test Len
	if list.Len() != 3 {
		t.Errorf("Expected 3 pages, got %d", list.Len())
	}

	// Test Less (longer names should be "less" to sort first)
	if !list.Less(1, 0) { // "very-long-name" (14) should be less than "short" (5)
		t.Error("Longer names should sort before shorter names")
	}

	if list.Less(0, 1) { // "short" (5) should not be less than "very-long-name" (14)
		t.Error("Shorter names should not sort before longer names")
	}

	// Test Swap
	original0 := list[0]
	original1 := list[1]
	list.Swap(0, 1)
	if list[0] != original1 || list[1] != original0 {
		t.Error("Swap did not exchange elements correctly")
	}
}

// TestPageLinkNode tests the PageLink AST node
func TestPageLinkNode(t *testing.T) {
	// Create a mock page (we just need something with a Name)
	mockPage := &mockPage{name: "test-page.md"}

	link := &PageLink{
		page: mockPage,
	}

	// Test Kind
	if link.Kind() != KindPageLink {
		t.Errorf("Expected Kind to be KindPageLink, got %v", link.Kind())
	}

	// Test Dump (should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump panicked: %v", r)
		}
	}()
	link.Dump([]byte("test source"), 0)
}

// TestPageLinkParser_Trigger tests the parser triggers
func TestPageLinkParser_Trigger(t *testing.T) {
	parser := &pageLinkParser{}
	triggers := parser.Trigger()

	expectedTriggers := []byte{' ', '*', '_', '~', '('}
	if len(triggers) != len(expectedTriggers) {
		t.Fatalf("Expected %d triggers, got %d", len(expectedTriggers), len(triggers))
	}

	for i, expected := range expectedTriggers {
		if triggers[i] != expected {
			t.Errorf("Trigger %d: expected '%c', got '%c'", i, expected, triggers[i])
		}
	}
}

// TestPageLinkParser_Parse tests the parser with various scenarios
func TestPageLinkParser_Parse(t *testing.T) {
	// Setup mock pages
	autolinkPage_lck.Lock()
	autolinkPages = []*NormalizedPage{
		{
			page:           &mockPage{name: "long-page-name.md", filename: "long-page-name.md"},
			normalizedName: "long-page-name.md",
		},
		{
			page:           &mockPage{name: "test.md", filename: "test.md"},
			normalizedName: "test.md",
		},
	}
	autolinkPage_lck.Unlock()

	tests := []struct {
		name      string
		input     string
		expectNil bool
	}{
		{
			name:      "Match at start after space",
			input:     " long-page-name.md is great",
			expectNil: false,
		},
		{
			name:      "Match at start after asterisk",
			input:     "*test.md is here",
			expectNil: false,
		},
		{
			name:      "No match",
			input:     " nonexistent-page.md",
			expectNil: true,
		},
		{
			name:      "Match but followed by alphanumeric",
			input:     " test.mdx",
			expectNil: true,
		},
	}

	p := &pageLinkParser{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := text.NewReader([]byte(tt.input))
			pc := parser.NewContext()
			parent := ast.NewParagraph()

			result := p.Parse(parent, reader, pc)

			if tt.expectNil && result != nil {
				t.Errorf("Expected nil result, got %T", result)
			}
			if !tt.expectNil && result == nil {
				t.Error("Expected non-nil result, got nil")
			}
			if !tt.expectNil && result != nil {
				if result.Kind() != KindPageLink {
					t.Errorf("Expected PageLink node, got %v", result.Kind())
				}
			}
		})
	}
}








// TestNormalizedName tests case-insensitive matching
func TestNormalizedName(t *testing.T) {
	// Setup pages with different casings
	autolinkPage_lck.Lock()
	autolinkPages = []*NormalizedPage{
		{
			page:           &mockPage{name: "Test-Page.md", filename: "Test-Page.md"},
			normalizedName: strings.ToLower("Test-Page.md"),
		},
	}
	autolinkPage_lck.Unlock()

	// Test that we can match with different casings
	tests := []string{
		" Test-Page.md is here",
		" test-page.md is here",
		" TEST-PAGE.MD is here",
	}

	p := &pageLinkParser{}
	for _, input := range tests {
		reader := text.NewReader([]byte(input))
		pc := parser.NewContext()
		parent := ast.NewParagraph()

		result := p.Parse(parent, reader, pc)
		if result == nil {
			t.Errorf("Expected to match %q with case-insensitive search", input)
		}
	}
}

// TestContainLinkToFrom_RelativeWithContext tests context-aware relative link resolution
func TestContainLinkToFrom_RelativeWithContext(t *testing.T) {
	tests := []struct {
		name         string
		sourcePage   string
		targetPage   string
		linkDest     string
		shouldMatch  bool
		description  string
	}{
		{
			name:         "Same directory - simple filename",
			sourcePage:   "folder/source.md",
			targetPage:   "folder/target.md",
			linkDest:     "target.md",
			shouldMatch:  true,
			description:  "Relative link in same directory should match",
		},
		{
			name:         "Subdirectory - relative path",
			sourcePage:   "folder/source.md",
			targetPage:   "folder/sub/target.md",
			linkDest:     "sub/target.md",
			shouldMatch:  true,
			description:  "Relative link with subdirectory should match",
		},
		{
			name:         "Parent directory - relative path",
			sourcePage:   "folder/sub/source.md",
			targetPage:   "folder/target.md",
			linkDest:     "../target.md",
			shouldMatch:  true,
			description:  "Relative link to parent directory should match",
		},
		{
			name:         "Different folders - same basename",
			sourcePage:   "folder1/source.md",
			targetPage:   "folder2/target.md",
			linkDest:     "target.md",
			shouldMatch:  true, // Fallback to basename matching
			description:  "Basename fallback should match even in different folders",
		},
		{
			name:         "Different folders - relative path",
			sourcePage:   "folder1/source.md",
			targetPage:   "folder2/target.md",
			linkDest:     "../folder2/target.md",
			shouldMatch:  true,
			description:  "Explicit relative path to different folder should match",
		},
		{
			name:         "Root level - simple filename",
			sourcePage:   "source.md",
			targetPage:   "target.md",
			linkDest:     "target.md",
			shouldMatch:  true,
			description:  "Files at root level should match by name",
		},
		{
			name:         "Wrong target",
			sourcePage:   "folder/source.md",
			targetPage:   "folder/target.md",
			linkDest:     "other.md",
			shouldMatch:  false,
			description:  "Link to different file should not match",
		},
		{
			name:         "Absolute path mismatch",
			sourcePage:   "folder/source.md",
			targetPage:   "folder/target.md",
			linkDest:     "/other/target.md",
			shouldMatch:  false,
			description:  "Absolute path to different location should not match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create source and target pages
			sourcePage := &mockPage{name: tt.sourcePage, filename: tt.sourcePage}
			targetPage := &mockPage{name: tt.targetPage, filename: tt.targetPage}

			// Create a link node
			link := ast.NewLink()
			link.Destination = []byte(tt.linkDest)
			link.AppendChild(link, ast.NewString([]byte("link text")))

			// Create a paragraph containing the link
			para := ast.NewParagraph()
			para.AppendChild(para, link)

			// Test containLinkToFrom
			result := containLinkToFrom(para, sourcePage, targetPage)
			if result != tt.shouldMatch {
				t.Errorf("%s: expected %v, got %v. Source: %s, Target: %s, Link: %s",
					tt.description, tt.shouldMatch, result, tt.sourcePage, tt.targetPage, tt.linkDest)
			}
		})
	}
}

// TestContainLinkToFrom_AbsoluteLink tests absolute link handling with context
func TestContainLinkToFrom_AbsoluteLink(t *testing.T) {
	sourcePage := &mockPage{name: "folder/source.md", filename: "folder/source.md"}
	targetPage := &mockPage{name: "other/target.md", filename: "other/target.md"}

	// Create an absolute link
	link := ast.NewLink()
	link.Destination = []byte("/other/target.md")
	link.AppendChild(link, ast.NewString([]byte("link text")))

	para := ast.NewParagraph()
	para.AppendChild(para, link)

	// Absolute links should match regardless of source location
	if !containLinkToFrom(para, sourcePage, targetPage) {
		t.Error("Absolute link should match target page")
	}
}

// TestContainLinkToFrom_PageLink tests PageLink handling with context
func TestContainLinkToFrom_PageLink(t *testing.T) {
	sourcePage := &mockPage{name: "folder/source.md", filename: "folder/source.md"}
	targetPage := &mockPage{name: "folder/target.md", filename: "folder/target.md"}

	// Create a PageLink node
	pageLink := &PageLink{
		page: targetPage,
	}
	pageLink.AppendChild(pageLink, ast.NewString([]byte("target.md")))

	para := ast.NewParagraph()
	para.AppendChild(para, pageLink)

	// PageLink should match when filenames are the same
	if !containLinkToFrom(para, sourcePage, targetPage) {
		t.Error("PageLink should match target page")
	}
}

// TestContainLinkToFrom_ComplexPath tests path normalization
func TestContainLinkToFrom_ComplexPath(t *testing.T) {
	sourcePage := &mockPage{name: "a/b/c/source.md", filename: "a/b/c/source.md"}
	targetPage := &mockPage{name: "a/b/target.md", filename: "a/b/target.md"}

	// Link with unnecessary path segments (../c/../target.md should resolve to ../target.md)
	link := ast.NewLink()
	link.Destination = []byte("../c/../target.md")
	link.AppendChild(link, ast.NewString([]byte("link text")))

	para := ast.NewParagraph()
	para.AppendChild(para, link)

	// Path should be normalized and match
	if !containLinkToFrom(para, sourcePage, targetPage) {
		t.Error("Complex relative path should be normalized and match")
	}
}


// mockPage is a minimal Page implementation for testing
type mockPage struct {
	name     string
	filename string
}

func (m *mockPage) Name() string     { return m.name }
func (m *mockPage) FileName() string { return m.filename }
func (m *mockPage) Exists() bool     { return true }
func (m *mockPage) Render() template.HTML {
	return template.HTML("<h1>Mock Page</h1>")
}
func (m *mockPage) Content() Markdown {
	return Markdown("# Mock Page\nContent")
}
func (m *mockPage) Delete() bool        { return false }
func (m *mockPage) Write(Markdown) bool { return false }
func (m *mockPage) ModTime() time.Time  { return time.Now() }
func (m *mockPage) AST() ([]byte, ast.Node) {
	return []byte("# Mock Page\nContent"), ast.NewDocument()
}
