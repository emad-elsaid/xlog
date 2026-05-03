package shortcode_test

import (
	"html/template"
	"testing"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/extensions/shortcode"
	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransformShortCodeBlocks_RegisteredShortcode(t *testing.T) {
	// Register a test shortcode
	shortcode.RegisterShortCode("testcode", shortcode.ShortCode{
		Render:  func(m xlog.Markdown) template.HTML { return template.HTML("rendered") },
		Default: "",
	})

	input := "```testcode\ncontent\n```"
	doc := parseMarkdown(t, input)

	// Verify transformation occurred
	var foundShortCodeBlock bool
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if _, ok := n.(*shortcode.ShortCodeBlock); ok {
				foundShortCodeBlock = true
			}
		}
		return ast.WalkContinue, nil
	})

	assert.True(t, foundShortCodeBlock, "Expected ShortCodeBlock node after transformation")
}

func TestTransformShortCodeBlocks_UnregisteredShortcode(t *testing.T) {
	input := "```unregistered\ncontent\n```"
	doc := parseMarkdown(t, input)

	// Verify no transformation occurred
	var foundFencedCodeBlock bool
	var foundShortCodeBlock bool
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if _, ok := n.(*ast.FencedCodeBlock); ok {
				foundFencedCodeBlock = true
			}
			if _, ok := n.(*shortcode.ShortCodeBlock); ok {
				foundShortCodeBlock = true
			}
		}
		return ast.WalkContinue, nil
	})

	assert.True(t, foundFencedCodeBlock, "Expected FencedCodeBlock to remain unchanged")
	assert.False(t, foundShortCodeBlock, "Expected no ShortCodeBlock for unregistered shortcode")
}

func TestTransformShortCodeBlocks_MultipleShortcodes(t *testing.T) {
	// Register multiple test shortcodes
	shortcode.RegisterShortCode("first", shortcode.ShortCode{
		Render:  func(m xlog.Markdown) template.HTML { return template.HTML("first") },
		Default: "",
	})
	shortcode.RegisterShortCode("second", shortcode.ShortCode{
		Render:  func(m xlog.Markdown) template.HTML { return template.HTML("second") },
		Default: "",
	})

	input := "```first\ncontent1\n```\n\n```second\ncontent2\n```"
	doc := parseMarkdown(t, input)

	// Count ShortCodeBlock nodes
	var count int
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if _, ok := n.(*shortcode.ShortCodeBlock); ok {
				count++
			}
		}
		return ast.WalkContinue, nil
	})

	assert.Equal(t, 2, count, "Expected 2 ShortCodeBlock nodes")
}

func TestTransformShortCodeBlocks_MixedContent(t *testing.T) {
	shortcode.RegisterShortCode("mixed", shortcode.ShortCode{
		Render:  func(m xlog.Markdown) template.HTML { return template.HTML("mixed") },
		Default: "",
	})

	input := "```go\ncode\n```\n\n```mixed\ncontent\n```\n\n```python\nmore code\n```"
	doc := parseMarkdown(t, input)

	// Count different block types
	var shortcodeBlocks, fencedBlocks int
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			switch n.(type) {
			case *shortcode.ShortCodeBlock:
				shortcodeBlocks++
			case *ast.FencedCodeBlock:
				fencedBlocks++
			}
		}
		return ast.WalkContinue, nil
	})

	assert.Equal(t, 1, shortcodeBlocks, "Expected 1 ShortCodeBlock")
	assert.Equal(t, 2, fencedBlocks, "Expected 2 FencedCodeBlock nodes")
}

func TestTransformShortCodeBlocks_EmptyLanguage(t *testing.T) {
	input := "```\ncontent\n```"
	doc := parseMarkdown(t, input)

	// Verify no transformation for blocks without language
	var foundShortCodeBlock bool
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if _, ok := n.(*shortcode.ShortCodeBlock); ok {
				foundShortCodeBlock = true
			}
		}
		return ast.WalkContinue, nil
	})

	assert.False(t, foundShortCodeBlock, "Expected no transformation for empty language")
}

func TestTransformShortCodeBlocks_ConcurrentRegistration(t *testing.T) {
	// Test thread-safety during transformation with concurrent registration
	done := make(chan bool)

	// Goroutine that registers shortcodes
	go func() {
		for i := 0; i < 100; i++ {
			shortcode.RegisterShortCode("concurrent", shortcode.ShortCode{
				Render:  func(m xlog.Markdown) template.HTML { return template.HTML("test") },
				Default: "",
			})
		}
		done <- true
	}()

	// Main goroutine that transforms
	for i := 0; i < 100; i++ {
		input := "```concurrent\ntest\n```"
		_ = parseMarkdown(t, input)
	}

	<-done
}

func TestTransformShortCodeBlocks_NestedStructure(t *testing.T) {
	shortcode.RegisterShortCode("outer", shortcode.ShortCode{
		Render:  func(m xlog.Markdown) template.HTML { return template.HTML("outer") },
		Default: "",
	})

	// Fenced code blocks cannot be nested in markdown, but test within other blocks
	input := "> Quote\n>\n> ```outer\n> content\n> ```"
	doc := parseMarkdown(t, input)

	// Verify transformation works in nested context
	var foundShortCodeBlock bool
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if _, ok := n.(*shortcode.ShortCodeBlock); ok {
				foundShortCodeBlock = true
			}
		}
		return ast.WalkContinue, nil
	})

	assert.True(t, foundShortCodeBlock, "Expected ShortCodeBlock in nested structure")
}

func TestTransformShortCodeBlocks_PreservesContent(t *testing.T) {
	shortcode.RegisterShortCode("preserve", shortcode.ShortCode{
		Render:  func(m xlog.Markdown) template.HTML { return template.HTML("test") },
		Default: "",
	})

	input := "```preserve\noriginal content\nline 2\n```"
	doc := parseMarkdown(t, input)

	// Verify content is preserved in ShortCodeBlock
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if scb, ok := n.(*shortcode.ShortCodeBlock); ok {
				// Access underlying FencedCodeBlock content
				lines := scb.Lines()
				require.NotNil(t, lines, "Expected content to be preserved")
			}
		}
		return ast.WalkContinue, nil
	})
}

// parseMarkdown is a helper that parses markdown with shortcode extension enabled.
func parseMarkdown(t *testing.T, input string) *ast.Document {
	t.Helper()

	md := markdown.New(
		markdown.WithExtensions(&shortcode.ShortCodeEx{}),
	)

	reader := text.NewReader([]byte(input))
	doc := md.Parser().Parse(reader)
	require.NotNil(t, doc, "Expected non-nil document")

	docPtr, ok := doc.(*ast.Document)
	require.True(t, ok, "Expected document to be *ast.Document")
	return docPtr
}
