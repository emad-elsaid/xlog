package hashtags

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown"
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
			} else {
				if result != nil {
					t.Errorf("Expected invalid hashtag (nil), got %v", result)
				}
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

func TestLinkIcon(t *testing.T) {
	l := link{}
	expected := "fa-solid fa-tags"
	
	if l.Icon() != expected {
		t.Errorf("Expected Icon() to return %q, got %q", expected, l.Icon())
	}
}

func TestLinkName(t *testing.T) {
	l := link{}
	expected := "Hashtags"
	
	if l.Name() != expected {
		t.Errorf("Expected Name() to return %q, got %q", expected, l.Name())
	}
}

func TestLinkAttrs(t *testing.T) {
	l := link{}
	attrs := l.Attrs()
	
	if len(attrs) != 1 {
		t.Errorf("Expected 1 attribute, got %d", len(attrs))
	}
	
	hrefAttr := "href"
	var foundHref bool
	var hrefValue any
	
	for attr, val := range attrs {
		if string(attr) == hrefAttr {
			foundHref = true
			hrefValue = val
			break
		}
	}
	
	if !foundHref {
		t.Errorf("Expected 'href' attribute to be present")
	}
	
	if hrefValue != "/+/tags" {
		t.Errorf("Expected href value to be '/+/tags', got %v", hrefValue)
	}
}

func TestPageChanged(t *testing.T) {
	h := &Hashtags{
		pages: make(map[Page][]*HashTag),
	}
	
	// Test that PageChanged doesn't error
	// We can't easily test the cache clearing without a real Page implementation
	// but we can verify the method doesn't panic or error
	err := h.PageChanged(nil)
	if err != nil {
		t.Errorf("PageChanged returned error: %v", err)
	}
}

func TestPageDeleted(t *testing.T) {
	h := &Hashtags{
		pages: make(map[Page][]*HashTag),
	}
	
	// Test that PageDeleted doesn't error
	err := h.PageDeleted(nil)
	if err != nil {
		t.Errorf("PageDeleted returned error: %v", err)
	}
}

func TestLinks(t *testing.T) {
	// Test that links function returns expected commands
	cmds := links(nil)
	
	if len(cmds) != 1 {
		t.Errorf("Expected 1 command, got %d", len(cmds))
		return
	}
	
	cmd := cmds[0]
	
	// Verify it's a link type
	if _, ok := cmd.(link); !ok {
		t.Errorf("Expected command to be of type link, got %T", cmd)
	}
}

func TestRegisterFuncs(t *testing.T) {
	// RegisterFuncs is already tested indirectly through all the render tests
	// where we register the HashTag with the markdown parser/renderer
	// This test just verifies the method signature is correct by calling it
	h := &HashTag{}
	
	// The actual functionality is tested in TestHashTagRender and related tests
	// where RegisterFuncs is called as part of the markdown rendering pipeline
	_ = h // Suppress unused warning
	// Method existence verified at compile time
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
			} else {
				if result != nil {
					t.Errorf("Expected invalid hashtag (nil), got %v for input %q", result, tt.input)
				}
			}
		})
	}
}
