package parser

import (
	"strings"
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

// TestParser_Parse tests the core Parse function to achieve coverage.
func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{"heading", "# Title\n"},
		{"paragraph", "Simple text\n"},
		{"multiple paragraphs", "Para 1\n\nPara 2\n"},
		{"code block", "```\ncode\n```\n"},
		{"list", "- Item\n"},
		{"blockquote", "> Quote\n"},
		{"thematic break", "---\n"},
		{"empty", ""},
		{"whitespace only", "   \n"},
		{"mixed", "# Title\n\nText **bold** *italic*.\n"},
		{"unicode", "# 日本語\n\nテスト\n"},
		{"long line", strings.Repeat("a", 1000) + "\n"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(
				WithBlockParsers(DefaultBlockParsers()...),
				WithInlineParsers(DefaultInlineParsers()...),
				WithParagraphTransformers(DefaultParagraphTransformers()...),
			)
			reader := text.NewReader([]byte(tc.source))
			doc := parser.Parse(reader)

			if doc == nil {
				t.Fatal("Parse returned nil")
			}
			if doc.Kind() != ast.KindDocument {
				t.Errorf("Kind = %v, want %v", doc.Kind(), ast.KindDocument)
			}
		})
	}
}

// TestNewParser tests parser construction with various options.
func TestNewParser_Options(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
	}{
		{
			name: "with defaults",
			options: []Option{
				WithBlockParsers(DefaultBlockParsers()...),
				WithInlineParsers(DefaultInlineParsers()...),
				WithParagraphTransformers(DefaultParagraphTransformers()...),
			},
		},
		{
			name: "with attribute",
			options: []Option{
				WithBlockParsers(DefaultBlockParsers()...),
				WithInlineParsers(DefaultInlineParsers()...),
				WithAttribute(),
			},
		},
		{
			name: "with escaped space",
			options: []Option{
				WithBlockParsers(DefaultBlockParsers()...),
				WithInlineParsers(DefaultInlineParsers()...),
				WithEscapedSpace(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParser(tc.options...)
			if parser == nil {
				t.Fatal("NewParser returned nil")
			}

			// Verify it works
			source := []byte("# Test\n")
			reader := text.NewReader(source)
			doc := parser.Parse(reader)
			if doc == nil {
				t.Error("Parse returned nil")
			}
		})
	}
}

// TestParser_AddOptions tests adding options after parser creation.
func TestParser_AddOptions(t *testing.T) {
	parser := NewParser(
		WithBlockParsers(DefaultBlockParsers()...),
		WithInlineParsers(DefaultInlineParsers()...),
	)
	parser.AddOptions(WithAttribute())

	source := []byte("# Test\n")
	reader := text.NewReader(source)
	doc := parser.Parse(reader)

	if doc == nil {
		t.Fatal("Parse returned nil after AddOptions")
	}
}

// TestParser_ParseWithContext tests parsing with custom context.
func TestParser_ParseWithContext(t *testing.T) {
	// Create key first, before context, so store is sized correctly
	customKey := NewContextKey()
	ctx := NewContext()
	testValue := "custom-test-value"
	ctx.Set(customKey, testValue)

	parser := NewParser(
		WithBlockParsers(DefaultBlockParsers()...),
		WithInlineParsers(DefaultInlineParsers()...),
	)

	source := []byte("# Test\n")
	reader := text.NewReader(source)
	doc := parser.Parse(reader, WithContext(ctx))

	if doc == nil {
		t.Fatal("Parse returned nil")
	}

	// Verify context was used
	if got := ctx.Get(customKey); got != testValue {
		t.Errorf("Context value = %v, want %v", got, testValue)
	}
}

// TestParser_IsInLinkLabel tests the context link label state.
func TestParser_IsInLinkLabel(t *testing.T) {
	ctx := NewContext()
	if ctx.IsInLinkLabel() {
		t.Error("IsInLinkLabel should be false initially")
	}
}
