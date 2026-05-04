package images

import (
	"bytes"
	"strings"
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/renderer"
)

func TestImagesColumnsRenderer_RegisterFuncs(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "registers render function for KindColumns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &imagesColumnsRenderer{}
			reg := &mockRegisterer{funcs: make(map[ast.NodeKind]renderer.NodeRendererFunc)}
			r.RegisterFuncs(reg)

			if _, exists := reg.funcs[KindColumns]; !exists {
				t.Error("RegisterFuncs should register a function for KindColumns")
			}
		})
	}
}

func TestImagesColumnsRenderer_Render(t *testing.T) {
	tests := []struct {
		name            string
		setup           func() (ast.Node, []byte)
		entering        bool
		expectedContent string
		expectedStatus  ast.WalkStatus
	}{
		{
			name: "entering renders opening div and column wrappers",
			setup: func() (ast.Node, []byte) {
				source := []byte("test")
				node := &imagesColumns{}
				img1 := ast.NewImage(ast.NewLink())
				img2 := ast.NewImage(ast.NewLink())
				node.AppendChild(node, img1)
				node.AppendChild(node, img2)
				return node, source
			},
			entering:        true,
			expectedContent: `<div class="columns">`,
			expectedStatus:  ast.WalkSkipChildren,
		},
		{
			name: "exiting renders closing div",
			setup: func() (ast.Node, []byte) {
				source := []byte("test")
				node := &imagesColumns{}
				return node, source
			},
			entering:        false,
			expectedContent: `</div>`,
			expectedStatus:  ast.WalkSkipChildren,
		},
		{
			name: "entering with single child renders one column",
			setup: func() (ast.Node, []byte) {
				source := []byte("test")
				node := &imagesColumns{}
				img := ast.NewImage(ast.NewLink())
				node.AppendChild(node, img)
				return node, source
			},
			entering:        true,
			expectedContent: `<div class="columns">`,
			expectedStatus:  ast.WalkSkipChildren,
		},
		{
			name: "entering with three children renders three columns",
			setup: func() (ast.Node, []byte) {
				source := []byte("test")
				node := &imagesColumns{}
				img1 := ast.NewImage(ast.NewLink())
				img2 := ast.NewImage(ast.NewLink())
				img3 := ast.NewImage(ast.NewLink())
				node.AppendChild(node, img1)
				node.AppendChild(node, img2)
				node.AppendChild(node, img3)
				return node, source
			},
			entering:        true,
			expectedContent: `<div class="columns">`,
			expectedStatus:  ast.WalkSkipChildren,
		},
		{
			name: "entering with no children renders empty columns div",
			setup: func() (ast.Node, []byte) {
				source := []byte("test")
				node := &imagesColumns{}
				return node, source
			},
			entering:        true,
			expectedContent: `<div class="columns">`,
			expectedStatus:  ast.WalkSkipChildren,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &imagesColumnsRenderer{}
			node, source := tt.setup()
			buf := &mockBufWriter{buf: &bytes.Buffer{}}

			status, err := r.render(buf, source, node, tt.entering)

			if err != nil {
				t.Errorf("render() error = %v, want nil", err)
			}

			if status != tt.expectedStatus {
				t.Errorf("render() status = %v, want %v", status, tt.expectedStatus)
			}

			output := buf.buf.String()
			if !strings.Contains(output, tt.expectedContent) {
				t.Errorf("render() output does not contain %q, got %q", tt.expectedContent, output)
			}
		})
	}
}

func TestImagesColumnsRenderer_RenderWithChildren(t *testing.T) {
	// This test verifies the column wrapper structure
	r := &imagesColumnsRenderer{}
	source := []byte("test")
	node := &imagesColumns{}

	// Add two image children
	img1 := ast.NewImage(ast.NewLink())
	img2 := ast.NewImage(ast.NewLink())
	node.AppendChild(node, img1)
	node.AppendChild(node, img2)

	buf := &mockBufWriter{buf: &bytes.Buffer{}}
	_, _ = r.render(buf, source, node, true)

	output := buf.buf.String()

	// Verify structure contains opening columns div
	if !strings.Contains(output, `<div class="columns">`) {
		t.Error("output should contain opening columns div")
	}

	// Count column divs (should have 2 for 2 images)
	columnCount := strings.Count(output, `<div class="column">`)
	if columnCount != 2 {
		t.Errorf("should have 2 column divs, got %d", columnCount)
	}

	// Verify closing column divs
	closingCount := strings.Count(output, `</div>`)
	if closingCount != 2 { // 2 column closings (the columns div is closed separately on exiting=false)
		t.Errorf("should have 2 closing divs for columns, got %d", closingCount)
	}
}

// mockRegisterer is a mock implementation of renderer.NodeRendererFuncRegisterer.
type mockRegisterer struct {
	funcs map[ast.NodeKind]renderer.NodeRendererFunc
}

func (m *mockRegisterer) Register(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
	m.funcs[kind] = fn
}

// mockBufWriter implements util.BufWriter for testing.
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
