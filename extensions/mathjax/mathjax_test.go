package mathjax

import (
	"bytes"
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/text"
)

// Test inline math parser.
func TestInlineMathParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     string
		wantNode bool
	}{
		{
			name:     "single dollar inline math",
			input:    "$E = mc^2$",
			want:     "E = mc^2",
			wantNode: true,
		},
		{
			name:     "double dollar inline math",
			input:    "$$x^2 + y^2 = z^2$$",
			want:     "x^2 + y^2 = z^2",
			wantNode: true,
		},
		{
			name:     "inline math with spaces",
			input:    "$ x + y $",
			want:     "x + y",
			wantNode: true,
		},
		{
			name:     "empty inline math",
			input:    "$$",
			want:     "",
			wantNode: false,
		},
		{
			name:     "unclosed inline math returns text",
			input:    "$incomplete",
			want:     "",
			wantNode: false,
		},
		{
			name:     "inline math multiline",
			input:    "$x +\ny$",
			want:     "x +\ny",
			wantNode: true,
		},
		{
			name:     "inline math with special chars",
			input:    "$\\frac{a}{b}$",
			want:     "\\frac{a}{b}",
			wantNode: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := &inlineMathParser{}
			reader := text.NewReader([]byte(tc.input))
			pc := parser.NewContext()
			parent := ast.NewDocument()

			node := p.Parse(parent, reader, pc)

			if !tc.wantNode {
				if node != nil && node.Kind() == KindInlineMath {
					t.Errorf("Expected no InlineMath node, got one")
				}
				return
			}

			if node == nil {
				t.Fatal("Expected node, got nil")
			}

			if node.Kind() != KindInlineMath {
				t.Errorf("Expected InlineMath node, got %v", node.Kind())
				return
			}

			// Extract text from children
			var buf bytes.Buffer
			for c := node.FirstChild(); c != nil; c = c.NextSibling() {
				seg := c.(*ast.Text).Segment
				buf.Write(seg.Value([]byte(tc.input)))
			}

			got := buf.String()
			if got != tc.want {
				t.Errorf("Content mismatch:\ngot:  %q\nwant: %q", got, tc.want)
			}
		})
	}
}

func TestInlineMathParser_Trigger(t *testing.T) {
	p := &inlineMathParser{}
	trigger := p.Trigger()

	if len(trigger) != 1 || trigger[0] != '$' {
		t.Errorf("Expected trigger ['$'], got %v", trigger)
	}
}

// Test block math parser.
func TestMathJaxBlockParser_Open(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantNode  bool
		wantState parser.State
	}{
		{
			name:      "valid block math start",
			input:     "$$\n",
			wantNode:  true,
			wantState: parser.NoChildren,
		},
		{
			name:      "single dollar not block math",
			input:     "$\n",
			wantNode:  false,
			wantState: parser.NoChildren,
		},
		{
			name:      "no dollar sign",
			input:     "regular text\n",
			wantNode:  false,
			wantState: parser.NoChildren,
		},
		{
			name:      "empty line",
			input:     "",
			wantNode:  false,
			wantState: parser.NoChildren,
		},
		{
			name:      "triple dollar valid",
			input:     "$$$\n",
			wantNode:  true,
			wantState: parser.NoChildren,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := &mathJaxBlockParser{}
			reader := text.NewReader([]byte(tc.input))
			pc := parser.NewContext()
			parent := ast.NewDocument()
			pc.SetBlockOffset(-1)

			if tc.input != "" {
				pc.SetBlockOffset(0)
			}

			node, state := p.Open(parent, reader, pc)

			if tc.wantNode && node == nil {
				t.Fatal("Expected node, got nil")
			}

			if !tc.wantNode && node != nil {
				t.Errorf("Expected no node, got %v", node.Kind())
			}

			if state != tc.wantState {
				t.Errorf("State mismatch: got %v, want %v", state, tc.wantState)
			}
		})
	}
}

func TestMathJaxBlockParser_Continue(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantState parser.State
	}{
		{
			name:      "math content continues",
			content:   "x^2 + y^2\n",
			wantState: parser.Continue | parser.NoChildren,
		},
		{
			name:      "closing delimiter",
			content:   "$$\n",
			wantState: parser.Close,
		},
		{
			name:      "empty line ends",
			content:   "",
			wantState: parser.NoChildren,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := &mathJaxBlockParser{}
			reader := text.NewReader([]byte(tc.content))
			pc := parser.NewContext()
			node := &MathBlock{}

			// Set required context
			pc.Set(mathBlockInfoKey, &mathBlockData{indent: 0})

			state := p.Continue(node, reader, pc)

			if state != tc.wantState {
				t.Errorf("State mismatch: got %v, want %v", state, tc.wantState)
			}
		})
	}
}

func TestMathJaxBlockParser_Close(t *testing.T) {
	p := &mathJaxBlockParser{}
	reader := text.NewReader([]byte(""))
	pc := parser.NewContext()
	node := &MathBlock{}

	pc.Set(mathBlockInfoKey, &mathBlockData{indent: 0})

	p.Close(node, reader, pc)

	// Verify context key cleared
	if pc.Get(mathBlockInfoKey) != nil {
		t.Error("Expected mathBlockInfoKey to be nil after Close")
	}
}

func TestMathJaxBlockParser_CanInterruptParagraph(t *testing.T) {
	p := &mathJaxBlockParser{}
	if !p.CanInterruptParagraph() {
		t.Error("Expected CanInterruptParagraph to return true")
	}
}

func TestMathJaxBlockParser_CanAcceptIndentedLine(t *testing.T) {
	p := &mathJaxBlockParser{}
	if p.CanAcceptIndentedLine() {
		t.Error("Expected CanAcceptIndentedLine to return false")
	}
}

func TestMathJaxBlockParser_Trigger(t *testing.T) {
	p := &mathJaxBlockParser{}
	trigger := p.Trigger()

	if trigger != nil {
		t.Errorf("Expected nil trigger, got %v", trigger)
	}
}

// Test nodes.
func TestInlineMath_Kind(t *testing.T) {
	node := &InlineMath{}
	if node.Kind() != KindInlineMath {
		t.Errorf("Expected KindInlineMath, got %v", node.Kind())
	}
}

func TestInlineMath_IsBlank(t *testing.T) {
	tests := []struct {
		name   string
		setup  func() (*InlineMath, []byte)
		expect bool
	}{
		{
			name: "blank content",
			setup: func() (*InlineMath, []byte) {
				source := []byte("   ")
				node := &InlineMath{}
				seg := text.NewSegment(0, 3)
				node.AppendChild(node, ast.NewTextSegment(seg))
				return node, source
			},
			expect: true,
		},
		{
			name: "non-blank content",
			setup: func() (*InlineMath, []byte) {
				source := []byte("x^2")
				node := &InlineMath{}
				seg := text.NewSegment(0, 3)
				node.AppendChild(node, ast.NewTextSegment(seg))
				return node, source
			},
			expect: false,
		},
		{
			name: "no children",
			setup: func() (*InlineMath, []byte) {
				node := &InlineMath{}
				return node, []byte("")
			},
			expect: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node, source := tc.setup()
			got := node.IsBlank(source)

			if got != tc.expect {
				t.Errorf("IsBlank() = %v, want %v", got, tc.expect)
			}
		})
	}
}

func TestInlineMath_Dump(t *testing.T) {
	node := &InlineMath{}
	source := []byte("test")
	// Should not panic
	node.Dump(source, 0)
}

func TestMathBlock_Kind(t *testing.T) {
	node := &MathBlock{}
	if node.Kind() != KindMathBlock {
		t.Errorf("Expected KindMathBlock, got %v", node.Kind())
	}
}

func TestMathBlock_IsRaw(t *testing.T) {
	node := &MathBlock{}
	if !node.IsRaw() {
		t.Error("Expected IsRaw to return true")
	}
}

func TestMathBlock_Dump(t *testing.T) {
	node := &MathBlock{}
	source := []byte("test")
	// Should not panic
	node.Dump(source, 0)
}

// Test renderers.
func TestInlineMathRenderer_RegisterFuncs(t *testing.T) {
	r := &InlineMathRenderer{startDelim: `\(`, endDelim: `\)`}

	// Mock registerer
	registered := false
	mockReg := &mockNodeRendererFuncRegisterer{
		registerFunc: func(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
			if kind == KindInlineMath {
				registered = true
			}
		},
	}

	r.RegisterFuncs(mockReg)

	if !registered {
		t.Error("Expected InlineMath renderer to be registered")
	}
}

func TestInlineMathRenderer_RenderInlineMath(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		entering bool
		wantHTML string
	}{
		{
			name:     "entering renders opening",
			content:  "x^2",
			entering: true,
			wantHTML: `<span>\(x^2`,
		},
		{
			name:     "exiting renders closing",
			content:  "x^2",
			entering: false,
			wantHTML: `\)</span>` + script,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := &InlineMathRenderer{startDelim: `\(`, endDelim: `\)`}
			source := []byte(tc.content)
			node := &InlineMath{}
			seg := text.NewSegment(0, len(tc.content))
			node.AppendChild(node, ast.NewTextSegment(seg))

			buf := &mockBufWriter{buf: &bytes.Buffer{}}
			status, err := r.renderInlineMath(buf, source, node, tc.entering)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tc.entering && status != ast.WalkSkipChildren {
				t.Errorf("Expected WalkSkipChildren when entering, got %v", status)
			}

			if !tc.entering && status != ast.WalkContinue {
				t.Errorf("Expected WalkContinue when exiting, got %v", status)
			}

			got := buf.buf.String()
			if got != tc.wantHTML {
				t.Errorf("HTML mismatch:\ngot:  %q\nwant: %q", got, tc.wantHTML)
			}
		})
	}
}

func TestMathBlockRenderer_RegisterFuncs(t *testing.T) {
	r := &MathBlockRenderer{startDelim: `\[`, endDelim: `\]`}

	registered := false
	mockReg := &mockNodeRendererFuncRegisterer{
		registerFunc: func(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
			if kind == KindMathBlock {
				registered = true
			}
		},
	}

	r.RegisterFuncs(mockReg)

	if !registered {
		t.Error("Expected MathBlock renderer to be registered")
	}
}

func TestMathBlockRenderer_RenderMathBlock(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		entering bool
		wantHTML string
	}{
		{
			name:     "entering renders opening",
			content:  "x^2 + y^2 = z^2",
			entering: true,
			wantHTML: `<p>\[x^2 + y^2 = z^2`,
		},
		{
			name:     "exiting renders closing",
			content:  "",
			entering: false,
			wantHTML: `\]</p>` + "\n" + script,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := &MathBlockRenderer{startDelim: `\[`, endDelim: `\]`}
			source := []byte(tc.content)
			node := &MathBlock{}

			if tc.entering && tc.content != "" {
				seg := text.NewSegment(0, len(tc.content))
				node.Lines().Append(seg)
			}

			buf := &mockBufWriter{buf: &bytes.Buffer{}}
			status, err := r.renderMathBlock(buf, source, node, tc.entering)

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if status != ast.WalkContinue {
				t.Errorf("Expected WalkContinue, got %v", status)
			}

			got := buf.buf.String()
			if got != tc.wantHTML {
				t.Errorf("HTML mismatch:\ngot:  %q\nwant: %q", got, tc.wantHTML)
			}
		})
	}
}

// Integration test: full markdown parsing.
func TestMathjaxIntegration(t *testing.T) {
	t.Skip("Integration test - requires full markdown setup")
}

// Mock for testing.
type mockNodeRendererFuncRegisterer struct {
	registerFunc func(ast.NodeKind, renderer.NodeRendererFunc)
}

func (m *mockNodeRendererFuncRegisterer) Register(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
	if m.registerFunc != nil {
		m.registerFunc(kind, fn)
	}
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
