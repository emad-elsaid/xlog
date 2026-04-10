package heading

import (
	"bytes"
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func TestHeadingExtensionName(t *testing.T) {
	ext := Heading{}
	expected := "heading"
	if ext.Name() != expected {
		t.Errorf("Expected name %q, got %q", expected, ext.Name())
	}
}

func TestHeadingInit(t *testing.T) {
	// Test that Init doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Init() panicked: %v", r)
		}
	}()

	ext := Heading{}
	ext.Init()
}

func TestHeadingRenderer_RegisterFuncs(t *testing.T) {
	r := &headingRenderer{}

	var registeredKind ast.NodeKind
	var registeredFunc renderer.NodeRendererFunc

	mockReg := &mockNodeRendererFuncRegisterer{
		registerFunc: func(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
			registeredKind = kind
			registeredFunc = fn
		},
	}

	r.RegisterFuncs(mockReg)

	if registeredKind != ast.KindHeading {
		t.Errorf("Expected to register ast.KindHeading, got: %v", registeredKind)
	}

	if registeredFunc == nil {
		t.Error("Expected registered function to not be nil")
	}
}

func TestHeadingRenderer_RenderEntering(t *testing.T) {
	r := &headingRenderer{}
	buf := &mockBufWriter{buf: &bytes.Buffer{}}
	source := []byte("# Test Heading")

	tests := []struct {
		name     string
		level    int
		expected string
	}{
		{"H1", 1, "<h1>"},
		{"H2", 2, "<h2>"},
		{"H3", 3, "<h3>"},
		{"H4", 4, "<h4>"},
		{"H5", 5, "<h5>"},
		{"H6", 6, "<h6>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.buf.Reset()
			node := ast.NewHeading(tt.level)

			status, err := r.render(buf, source, node, true)

			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			if status != ast.WalkContinue {
				t.Errorf("Expected WalkContinue, got: %v", status)
			}

			html := buf.buf.String()
			if html != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, html)
			}
		})
	}
}

func TestHeadingRenderer_RenderExiting(t *testing.T) {
	r := &headingRenderer{}
	buf := &mockBufWriter{buf: &bytes.Buffer{}}
	source := []byte("# Test Heading")

	tests := []struct {
		name     string
		level    int
		expected string
	}{
		{"H1", 1, "</h1>\n"},
		{"H2", 2, "</h2>\n"},
		{"H3", 3, "</h3>\n"},
		{"H4", 4, "</h4>\n"},
		{"H5", 5, "</h5>\n"},
		{"H6", 6, "</h6>\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.buf.Reset()
			node := ast.NewHeading(tt.level)

			status, err := r.render(buf, source, node, false)

			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			if status != ast.WalkContinue {
				t.Errorf("Expected WalkContinue, got: %v", status)
			}

			html := buf.buf.String()
			if html != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, html)
			}
		})
	}
}

func TestHeadingRenderer_RenderExitingWithID(t *testing.T) {
	r := &headingRenderer{}
	buf := &mockBufWriter{buf: &bytes.Buffer{}}
	source := []byte("# Test Heading")

	node := ast.NewHeading(2)
	node.SetAttributeString("id", []byte("test-heading"))

	status, err := r.render(buf, source, node, false)

	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if status != ast.WalkContinue {
		t.Errorf("Expected WalkContinue, got: %v", status)
	}

	html := buf.buf.String()

	// Should contain the anchor link with proper classes
	expectedAnchor := ` <a class="show-on-parent-hover is-hidden has-text-grey" href="#test-heading">¶</a>`
	if !contains(html, expectedAnchor) {
		t.Errorf("Expected HTML to contain anchor %q, got: %q", expectedAnchor, html)
	}

	// Should also contain closing tag
	if !contains(html, "</h2>\n") {
		t.Errorf("Expected HTML to contain closing tag, got: %q", html)
	}
}

func TestHeadingRenderer_RenderExitingWithAttributes(t *testing.T) {
	r := &headingRenderer{}
	buf := &mockBufWriter{buf: &bytes.Buffer{}}
	source := []byte("# Test")

	node := ast.NewHeading(1)
	node.SetAttributeString("class", []byte("custom-class"))

	// Entering should render attributes
	status, err := r.render(buf, source, node, true)
	if err != nil {
		t.Fatalf("Render entering failed: %v", err)
	}
	if status != ast.WalkContinue {
		t.Errorf("Expected WalkContinue, got: %v", status)
	}

	html := buf.buf.String()
	if !contains(html, "<h1") {
		t.Errorf("Expected opening tag, got: %q", html)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Mock implementations
type mockNodeRendererFuncRegisterer struct {
	registerFunc func(ast.NodeKind, renderer.NodeRendererFunc)
}

func (m *mockNodeRendererFuncRegisterer) Register(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
	if m.registerFunc != nil {
		m.registerFunc(kind, fn)
	}
}

func (m *mockNodeRendererFuncRegisterer) Prioritized() util.PrioritizedSlice {
	return nil
}

type mockBufWriter struct {
	buf *bytes.Buffer
}

func (m *mockBufWriter) Write(p []byte) (int, error) {
	return m.buf.Write(p)
}

func (m *mockBufWriter) WriteByte(c byte) error {
	return m.buf.WriteByte(c)
}

func (m *mockBufWriter) WriteRune(r rune) (int, error) {
	return m.buf.WriteRune(r)
}

func (m *mockBufWriter) WriteString(s string) (int, error) {
	return m.buf.WriteString(s)
}

func (m *mockBufWriter) Buffered() int {
	return m.buf.Len()
}

func (m *mockBufWriter) Available() int {
	return m.buf.Cap() - m.buf.Len()
}

func (m *mockBufWriter) Flush() error {
	return nil
}
