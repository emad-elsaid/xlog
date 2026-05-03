package markdown_test

import (
	"bytes"
	"testing"

	. "github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
	"github.com/emad-elsaid/xlog/markdown/testutil"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func TestAttributeAndAutoHeadingID(t *testing.T) {
	markdown := New(
		WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/options.txt", t, testutil.ParseCliCaseArg()...)
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		contains string
	}{
		{
			name:     "converts simple markdown to HTML",
			input:    []byte("# Hello World"),
			contains: "<h1>Hello World</h1>",
		},
		{
			name:     "converts paragraph",
			input:    []byte("This is a paragraph."),
			contains: "<p>This is a paragraph.</p>",
		},
		{
			name:     "converts emphasis",
			input:    []byte("This is *italic* text."),
			contains: "<em>italic</em>",
		},
		{
			name:     "converts strong",
			input:    []byte("This is **bold** text."),
			contains: "<strong>bold</strong>",
		},
		{
			name:     "converts code block",
			input:    []byte("```\ncode here\n```"),
			contains: "<pre><code>code here",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Convert(tt.input, &buf)

			if err != nil {
				t.Fatalf("Convert returned error: %v", err)
			}

			output := buf.String()
			if output == "" {
				t.Error("Convert produced empty output")
			}

			if tt.contains != "" && !bytes.Contains([]byte(output), []byte(tt.contains)) {
				t.Errorf("Expected output to contain %q, got:\n%s", tt.contains, output)
			}
		})
	}
}

func TestWithExtensions(t *testing.T) {
	// Create a mock extension that modifies the parser
	ext := &mockExtension{
		extendCalled: false,
	}

	md := New(WithExtensions(ext))

	if !ext.extendCalled {
		t.Error("Extension Extend method was not called")
	}

	// Verify the markdown instance works
	var buf bytes.Buffer
	err := md.Convert([]byte("# Test"), &buf)
	if err != nil {
		t.Errorf("Convert failed: %v", err)
	}
}

func TestWithParser(t *testing.T) {
	customParser := parser.NewParser(
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
	)

	md := New(WithParser(customParser))

	// Verify the parser was set
	if md.Parser() == nil {
		t.Error("Parser was not set")
	}

	// Verify it can convert
	var buf bytes.Buffer
	err := md.Convert([]byte("# Test"), &buf)
	if err != nil {
		t.Errorf("Convert with custom parser failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Convert produced no output")
	}
}

func TestWithRenderer(t *testing.T) {
	customRenderer := renderer.NewRenderer(
		renderer.WithNodeRenderers(
			util.Prioritized(html.NewRenderer(), 1000),
		),
	)

	md := New(WithRenderer(customRenderer))

	// Verify the renderer was set
	if md.Renderer() == nil {
		t.Error("Renderer was not set")
	}

	// Verify it can convert
	var buf bytes.Buffer
	err := md.Convert([]byte("**bold**"), &buf)
	if err != nil {
		t.Errorf("Convert with custom renderer failed: %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("<strong>bold</strong>")) {
		t.Errorf("Expected <strong> tag, got: %s", output)
	}
}

func TestSetParser(t *testing.T) {
	md := New()

	// Get original parser
	originalParser := md.Parser()
	if originalParser == nil {
		t.Fatal("Default parser should not be nil")
	}

	// Create and set a new parser
	newParser := parser.NewParser(
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
	)

	md.SetParser(newParser)

	// Verify the parser was changed
	currentParser := md.Parser()
	if currentParser == nil {
		t.Error("Parser should not be nil after SetParser")
	}

	// Verify it still works
	var buf bytes.Buffer
	err := md.Convert([]byte("# Heading"), &buf)
	if err != nil {
		t.Errorf("Convert after SetParser failed: %v", err)
	}
}

func TestSetRenderer(t *testing.T) {
	md := New()

	// Get original renderer
	originalRenderer := md.Renderer()
	if originalRenderer == nil {
		t.Fatal("Default renderer should not be nil")
	}

	// Create and set a new renderer
	newRenderer := renderer.NewRenderer(
		renderer.WithNodeRenderers(
			util.Prioritized(html.NewRenderer(), 1000),
		),
	)

	md.SetRenderer(newRenderer)

	// Verify the renderer was changed
	currentRenderer := md.Renderer()
	if currentRenderer == nil {
		t.Error("Renderer should not be nil after SetRenderer")
	}

	// Verify it still works
	var buf bytes.Buffer
	err := md.Convert([]byte("*italic*"), &buf)
	if err != nil {
		t.Errorf("Convert after SetRenderer failed: %v", err)
	}
}

func TestRenderer(t *testing.T) {
	md := New()

	r := md.Renderer()
	if r == nil {
		t.Error("Renderer() should not return nil")
	}

	// Verify the renderer is functional
	var buf bytes.Buffer
	err := md.Convert([]byte("Test content"), &buf)
	if err != nil {
		t.Errorf("Renderer failed to convert: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Renderer produced no output")
	}
}

func TestNewWithMultipleOptions(t *testing.T) {
	// Test combining multiple options
	customParser := parser.NewParser(
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
	)
	customRenderer := renderer.NewRenderer(
		renderer.WithNodeRenderers(util.Prioritized(html.NewRenderer(), 1000)),
	)
	ext := &mockExtension{}

	md := New(
		WithParser(customParser),
		WithRenderer(customRenderer),
		WithExtensions(ext),
	)

	if !ext.extendCalled {
		t.Error("Extension was not applied")
	}

	if md.Parser() == nil {
		t.Error("Parser was not set")
	}

	if md.Renderer() == nil {
		t.Error("Renderer was not set")
	}

	// Verify full functionality
	var buf bytes.Buffer
	err := md.Convert([]byte("# Combined Test\n\nWith **bold** text."), &buf)
	if err != nil {
		t.Errorf("Convert with combined options failed: %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("Combined Test")) {
		t.Error("Output missing expected content")
	}
}

func TestConvertWithParseOptions(t *testing.T) {
	var buf bytes.Buffer
	source := []byte("# Test")

	// Convert with parse options
	err := Convert(source, &buf)
	if err != nil {
		t.Fatalf("Convert with options failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Convert produced no output")
	}
}

// mockExtension is a test extension that tracks whether Extend was called.
type mockExtension struct {
	extendCalled bool
}

func (m *mockExtension) Extend(md Markdown) {
	m.extendCalled = true
	// Optionally modify the markdown instance
	md.Parser().AddOptions(parser.WithAutoHeadingID())
}
