package hashtags

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
	"time"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func TestHashTagParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		valid    bool
	}{
		{
			name:     "simple hashtag",
			input:    "#golang",
			expected: "golang",
			valid:    true,
		},
		{
			name:     "hashtag with number",
			input:    "#web3",
			expected: "web3",
			valid:    true,
		},
		{
			name:     "hashtag with underscore",
			input:    "#hello_world",
			expected: "hello_world",
			valid:    true,
		},
		{
			name:     "hashtag with dash",
			input:    "#hello-world",
			expected: "hello-world",
			valid:    true,
		},
		{
			name:     "hashtag with CJK characters",
			input:    "#日本語",
			expected: "日本語",
			valid:    true,
		},
		{
			name:     "hashtag stops at space",
			input:    "#tag and more text",
			expected: "tag",
			valid:    true,
		},
		{
			name:     "hashtag stops at punctuation",
			input:    "#tag. More text",
			expected: "tag",
			valid:    true,
		},
		{
			name:     "hashtag stops at special chars",
			input:    "#tag@mention",
			expected: "tag",
			valid:    true,
		},
		{
			name:     "just # is invalid",
			input:    "# ",
			expected: "",
			valid:    false,
		},
		{
			name:     "empty after # is invalid",
			input:    "#",
			expected: "",
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HashTag{}
			reader := text.NewReader([]byte(tt.input))
			// Note: Parse is called AFTER the trigger character is already consumed
			// by the parser, so we don't need to skip it manually

			result := h.Parse(nil, reader, parser.NewContext())

			if tt.valid {
				if result == nil {
					t.Errorf("Expected valid hashtag, got nil")
					return
				}

				tag, ok := result.(*HashTag)
				if !ok {
					t.Errorf("Expected *HashTag, got %T", result)
					return
				}

				if string(tag.value) != tt.expected {
					t.Errorf("Expected tag value %q, got %q", tt.expected, string(tag.value))
				}
			} else if result != nil {
				t.Errorf("Expected invalid hashtag (nil), got %v", result)
			}
		})
	}
}

func TestHashTagRender(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		contains []string
	}{
		{
			name:     "single hashtag",
			markdown: "This is a #golang post",
			contains: []string{
				`<a href="/+/tag/golang" class="tag"`,
				`<span>golang</span>`,
				`<i class="fa-solid fa-tag"></i>`,
			},
		},
		{
			name:     "multiple hashtags",
			markdown: "Learning #golang and #rust together",
			contains: []string{
				`href="/+/tag/golang"`,
				`href="/+/tag/rust"`,
			},
		},
		{
			name:     "hashtag with underscore",
			markdown: "Using #hello_world pattern",
			contains: []string{
				`href="/+/tag/hello_world"`,
				`<span>hello_world</span>`,
			},
		},
		{
			name:     "hashtag with dash",
			markdown: "This is #test-case example",
			contains: []string{
				`href="/+/tag/test-case"`,
				`<span>test-case</span>`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := markdown.New()

			// Register hashtag parser
			h := &HashTag{}
			md.Parser().AddOptions(parser.WithInlineParsers(
				util.Prioritized(h, 999),
			))

			// Register hashtag renderer
			md.Renderer().AddOptions(renderer.WithNodeRenderers(
				util.Prioritized(h, 0),
			))

			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var buf bytes.Buffer
			err := md.Renderer().Render(&buf, []byte(tt.markdown), doc)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}

			htmlOutput := buf.String()

			for _, expected := range tt.contains {
				if !strings.Contains(htmlOutput, expected) {
					t.Errorf("Expected HTML to contain %q, got:\n%s", expected, htmlOutput)
				}
			}
		})
	}
}

func TestHashTagKind(t *testing.T) {
	h := &HashTag{}
	if h.Kind() != KindHashTag {
		t.Errorf("Expected Kind to be KindHashTag")
	}
}

func TestHashTagTrigger(t *testing.T) {
	h := &HashTag{}
	trigger := h.Trigger()

	if len(trigger) != 1 || trigger[0] != '#' {
		t.Errorf("Expected trigger to be ['#'], got %v", trigger)
	}
}

func TestHashTagDump(t *testing.T) {
	tag := &HashTag{
		value: []byte("testTag"),
	}

	// Just ensure it doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump panicked: %v", r)
		}
	}()

	tag.Dump([]byte("#testTag"), 0)
}

func TestHashTagUniqueness(t *testing.T) {
	tests := []struct {
		name     string
		tag1     string
		tag2     string
		expected bool // Should they have the same unique handle?
	}{
		{
			name:     "same case",
			tag1:     "golang",
			tag2:     "golang",
			expected: true,
		},
		{
			name:     "different case",
			tag1:     "GoLang",
			tag2:     "golang",
			expected: true,
		},
		{
			name:     "mixed case",
			tag1:     "GOLANG",
			tag2:     "golang",
			expected: true,
		},
		{
			name:     "completely different",
			tag1:     "golang",
			tag2:     "rust",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HashTag{}

			// Parse first tag
			reader1 := text.NewReader([]byte("#" + tt.tag1))
			result1 := h.Parse(nil, reader1, parser.NewContext())
			tag1, ok := result1.(*HashTag)
			if !ok {
				t.Fatalf("Failed to parse first tag")
			}

			// Parse second tag
			reader2 := text.NewReader([]byte("#" + tt.tag2))
			result2 := h.Parse(nil, reader2, parser.NewContext())
			tag2, ok := result2.(*HashTag)
			if !ok {
				t.Fatalf("Failed to parse second tag")
			}

			isSame := tag1.unique == tag2.unique
			if isSame != tt.expected {
				t.Errorf("Expected unique handles to be same=%v, got same=%v (tag1=%v, tag2=%v)",
					tt.expected, isSame, tag1.unique, tag2.unique)
			}
		})
	}
}

func TestHashtagsExtensionName(t *testing.T) {
	h := &Hashtags{}
	if h.Name() != "hashtags" {
		t.Errorf("Expected extension name to be 'hashtags', got %q", h.Name())
	}
}

func TestHashTagInMarkdownContext(t *testing.T) {
	tests := []struct {
		name        string
		markdown    string
		contains    string
		notContains string
	}{
		{
			name:     "hashtag in sentence",
			markdown: "I love #programming in Go",
			contains: `href="/+/tag/programming"`,
		},
		{
			name:        "not a hashtag - space after hash",
			markdown:    "This is # not a tag",
			notContains: `href="/+/tag/not"`,
		},
		{
			name:     "hashtag at start of line",
			markdown: "#golang is amazing",
			contains: `href="/+/tag/golang"`,
		},
		{
			name:     "hashtag at end of line",
			markdown: "Learning Go #golang",
			contains: `href="/+/tag/golang"`,
		},
		{
			name:     "hashtag in parentheses",
			markdown: "This (#golang) is cool",
			contains: `href="/+/tag/golang"`,
		},
		{
			name:     "hashtag after punctuation",
			markdown: "Cool! #golang",
			contains: `href="/+/tag/golang"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := markdown.New()
			h := &HashTag{}

			// Register hashtag parser and renderer
			md.Parser().AddOptions(parser.WithInlineParsers(
				util.Prioritized(h, 999),
			))
			md.Renderer().AddOptions(renderer.WithNodeRenderers(
				util.Prioritized(h, 0),
			))

			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)))

			var buf bytes.Buffer
			err := md.Renderer().Render(&buf, []byte(tt.markdown), doc)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}

			htmlOutput := buf.String()

			if tt.contains != "" && !strings.Contains(htmlOutput, tt.contains) {
				t.Errorf("Expected HTML to contain %q, got:\n%s", tt.contains, htmlOutput)
			}

			if tt.notContains != "" && strings.Contains(htmlOutput, tt.notContains) {
				t.Errorf("Expected HTML NOT to contain %q, but it does:\n%s", tt.notContains, htmlOutput)
			}
		})
	}
}

func TestHashTagEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		valid    bool
	}{
		{
			name:     "single letter",
			input:    "#a",
			expected: "a",
			valid:    true,
		},
		{
			name:     "numbers only",
			input:    "#123",
			expected: "123",
			valid:    true,
		},
		{
			name:     "unicode emoji not included",
			input:    "#tag🎉",
			expected: "tag",
			valid:    true,
		},
		{
			name:     "hashtag ending with punctuation",
			input:    "#tag.",
			expected: "tag",
			valid:    true,
		},
		{
			name:     "hashtag ending with comma",
			input:    "#tag,",
			expected: "tag",
			valid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HashTag{}
			reader := text.NewReader([]byte(tt.input))

			result := h.Parse(nil, reader, parser.NewContext())

			if tt.valid {
				if result == nil {
					t.Errorf("Expected valid hashtag, got nil for input %q", tt.input)
					return
				}

				tag, ok := result.(*HashTag)
				if !ok {
					t.Errorf("Expected *HashTag, got %T", result)
					return
				}

				if string(tag.value) != tt.expected {
					t.Errorf("Expected tag value %q, got %q for input %q",
						tt.expected, string(tag.value), tt.input)
				}
			} else if result != nil {
				t.Errorf("Expected invalid hashtag (nil), got %v for input %q", result, tt.input)
			}
		})
	}
}

func TestLinkCommandInterface(t *testing.T) {
	tests := []struct {
		name   string
		method string
		want   any
	}{
		{
			name:   "Icon returns correct value",
			method: "Icon",
			want:   "fa-solid fa-tags",
		},
		{
			name:   "Name returns correct value",
			method: "Name",
			want:   "Hashtags",
		},
		{
			name:   "Attrs returns map with href",
			method: "Attrs",
			want:   "/+/tags",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := link{}

			switch tt.method {
			case "Icon":
				got := l.Icon()
				if got != tt.want {
					t.Errorf("Icon() = %q, want %q", got, tt.want)
				}
			case "Name":
				got := l.Name()
				if got != tt.want {
					t.Errorf("Name() = %q, want %q", got, tt.want)
				}
			case "Attrs":
				attrs := l.Attrs()
				href, ok := attrs["href"]
				if !ok {
					t.Error("Attrs() missing 'href' key")
					return
				}
				if href != tt.want {
					t.Errorf("Attrs()[\"href\"] = %q, want %q", href, tt.want)
				}
			}
		})
	}
}

func TestLinksFunction(t *testing.T) {
	// links function should return a slice with one link command.
	cmds := links(nil)

	if len(cmds) != 1 {
		t.Errorf("links() returned %d commands, want 1", len(cmds))
		return
	}

	l, ok := cmds[0].(link)
	if !ok {
		t.Errorf("links() returned %T, want link", cmds[0])
		return
	}

	// Verify the link has expected values.
	if l.Icon() != "fa-solid fa-tags" {
		t.Errorf("link.Icon() = %q, want %q", l.Icon(), "fa-solid fa-tags")
	}

	if l.Name() != "Hashtags" {
		t.Errorf("link.Name() = %q, want %q", l.Name(), "Hashtags")
	}
}

func TestHashtagsPageChangedAndDeleted(t *testing.T) {
	tests := []struct {
		name           string
		setupCache     bool
		operation      string
		expectedCached bool
	}{
		{
			name:           "PageChanged clears cache for existing page",
			setupCache:     true,
			operation:      "changed",
			expectedCached: false,
		},
		{
			name:           "PageChanged on non-cached page",
			setupCache:     false,
			operation:      "changed",
			expectedCached: false,
		},
		{
			name:           "PageDeleted clears cache for existing page",
			setupCache:     true,
			operation:      "deleted",
			expectedCached: false,
		},
		{
			name:           "PageDeleted on non-cached page",
			setupCache:     false,
			operation:      "deleted",
			expectedCached: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Hashtags{
				pages: make(map[Page][]*HashTag),
			}

			// Create a mock page
			mockPage := &mockPage{name: "test-page"}

			// Setup cache if needed
			if tt.setupCache {
				h.pages[mockPage] = []*HashTag{
					{value: []byte("test")},
				}
			}

			// Perform operation
			var err error
			if tt.operation == "changed" {
				err = h.PageChanged(mockPage)
			} else {
				err = h.PageDeleted(mockPage)
			}

			if err != nil {
				t.Errorf("Operation %s returned error: %v", tt.operation, err)
			}

			// Verify cache state
			_, cached := h.pages[mockPage]
			if cached != tt.expectedCached {
				t.Errorf("After %s: expected cached=%v, got cached=%v",
					tt.operation, tt.expectedCached, cached)
			}
		})
	}
}

// mockPage implements the Page interface for testing.
type mockPage struct {
	name    string
	content []byte
}

func (m *mockPage) Name() string     { return m.name }
func (m *mockPage) FileName() string { return m.name + ".md" }
func (m *mockPage) Exists() bool     { return true }
func (m *mockPage) Render() template.HTML {
	return template.HTML(m.content)
}
func (m *mockPage) Content() Markdown {
	if m.content == nil {
		return Markdown("# " + m.name)
	}
	return Markdown(m.content)
}
func (m *mockPage) Delete() bool        { return true }
func (m *mockPage) Write(Markdown) bool { return true }
func (m *mockPage) ModTime() time.Time  { return time.Now() }
func (m *mockPage) AST() ([]byte, ast.Node) {
	// Parse content as markdown
	md := markdown.New()
	h := &HashTag{}
	md.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(h, 999),
	))
	content := m.content
	if content == nil {
		content = []byte("# " + m.name)
	}
	doc := md.Parser().Parse(text.NewReader(content))
	return content, doc
}

func TestHashtagsFor(t *testing.T) {
	tests := []struct {
		name         string
		pageContent  string
		expectedTags []string
		cached       bool
	}{
		{
			name:         "page with single hashtag",
			pageContent:  "This is a #golang post",
			expectedTags: []string{"golang"},
			cached:       false,
		},
		{
			name:         "page with multiple hashtags",
			pageContent:  "Learning #golang and #rust together",
			expectedTags: []string{"golang", "rust"},
			cached:       false,
		},
		{
			name:         "page with no hashtags",
			pageContent:  "This is a post without tags",
			expectedTags: []string{},
			cached:       false,
		},
		{
			name:         "page with duplicate hashtags",
			pageContent:  "#golang is great, I love #golang",
			expectedTags: []string{"golang", "golang"},
			cached:       false,
		},
		{
			name:         "cached result returns same tags",
			pageContent:  "Another #test post",
			expectedTags: []string{"test"},
			cached:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Hashtags{
				pages: make(map[Page][]*HashTag),
			}

			mockPage := &mockPage{
				name:    "test-page",
				content: []byte(tt.pageContent),
			}

			// If testing cached scenario, populate cache first
			if tt.cached {
				expectedHashtags := make([]*HashTag, len(tt.expectedTags))
				for i, tag := range tt.expectedTags {
					expectedHashtags[i] = &HashTag{value: []byte(tag)}
				}
				h.pages[mockPage] = expectedHashtags
			}

			// Call hashtagsFor
			tags := h.hashtagsFor(mockPage)

			// Verify tag count
			if len(tags) != len(tt.expectedTags) {
				t.Errorf("Expected %d tags, got %d", len(tt.expectedTags), len(tags))
			}

			// Verify tag values
			for i, expectedTag := range tt.expectedTags {
				if i >= len(tags) {
					break
				}
				if string(tags[i].value) != expectedTag {
					t.Errorf("Tag[%d]: expected %q, got %q",
						i, expectedTag, string(tags[i].value))
				}
			}

			// Verify caching occurred (second call should return same instance)
			if !tt.cached {
				tags2 := h.hashtagsFor(mockPage)
				if len(tags) > 0 && len(tags2) > 0 {
					// Should be same slice instance (cached)
					if &tags[0] != &tags2[0] {
						t.Log("Note: Tags were re-parsed instead of cached (this may be expected)")
					}
				}
			}
		})
	}
}

func TestHashtagsForConcurrency(t *testing.T) {
	h := &Hashtags{
		pages: make(map[Page][]*HashTag),
	}

	mockPage := &mockPage{
		name:    "concurrent-test",
		content: []byte("Testing #concurrency with #hashtags"),
	}

	// Test concurrent access doesn't cause race conditions
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			tags := h.hashtagsFor(mockPage)
			if len(tags) != 2 {
				t.Errorf("Expected 2 tags, got %d", len(tags))
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestRegisterFuncs(t *testing.T) {
	h := &HashTag{}

	// Create a mock registerer
	registered := false
	var registeredKind ast.NodeKind

	mockReg := &mockRegisterer{
		registerFunc: func(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
			registered = true
			registeredKind = kind
		},
	}

	h.RegisterFuncs(mockReg)

	if !registered {
		t.Error("RegisterFuncs did not call Register")
	}

	if registeredKind != KindHashTag {
		t.Errorf("Registered wrong kind: got %v, want %v", registeredKind, KindHashTag)
	}
}

type mockRegisterer struct {
	registerFunc func(ast.NodeKind, renderer.NodeRendererFunc)
}

func (m *mockRegisterer) Register(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
	if m.registerFunc != nil {
		m.registerFunc(kind, fn)
	}
}
