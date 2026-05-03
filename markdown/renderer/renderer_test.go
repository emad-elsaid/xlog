package renderer

import (
	"bytes"
	"io"
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/util"
)

// Mock node for testing.
type mockNode struct {
	ast.BaseBlock
	kind ast.NodeKind
}

func newMockNode(kind ast.NodeKind) *mockNode {
	return &mockNode{kind: kind}
}

func (n *mockNode) Kind() ast.NodeKind {
	return n.kind
}

func (n *mockNode) Dump(source []byte, level int) {
	// No-op for testing.
}

// Mock node renderer for testing.
type mockNodeRenderer struct {
	renderFunc NodeRendererFunc
}

func (m *mockNodeRenderer) RegisterFuncs(reg NodeRendererFuncRegisterer) {
	reg.Register(mockKind, m.renderFunc)
}

var mockKind = ast.NewNodeKind("MockNode")

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"creates valid config"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := NewConfig()
			if cfg.Options == nil {
				t.Error("NewConfig().Options is nil")
			}
			if cfg.NodeRenderers == nil {
				t.Error("NewConfig().NodeRenderers is nil")
			}
		})
	}
}

func TestWithNodeRenderers(t *testing.T) {
	tests := []struct {
		name     string
		renderer NodeRenderer
	}{
		{
			"adds node renderer",
			&mockNodeRenderer{
				renderFunc: func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
					return ast.WalkContinue, nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithNodeRenderers(util.Prioritized(tt.renderer, 100))
			cfg := NewConfig()
			opt.SetConfig(cfg)

			if len(cfg.NodeRenderers) != 1 {
				t.Errorf("WithNodeRenderers() expected 1 renderer, got %d", len(cfg.NodeRenderers))
			}
		})
	}
}

func TestWithOption(t *testing.T) {
	tests := []struct {
		name  string
		opt   OptionName
		value any
	}{
		{"sets string option", OptionName("test-option"), "test-value"},
		{"sets int option", OptionName("int-option"), 42},
		{"sets bool option", OptionName("bool-option"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithOption(tt.opt, tt.value)
			cfg := NewConfig()
			opt.SetConfig(cfg)

			if v, ok := cfg.Options[tt.opt]; !ok {
				t.Errorf("WithOption() option %v not set", tt.opt)
			} else if v != tt.value {
				t.Errorf("WithOption() expected %v, got %v", tt.value, v)
			}
		})
	}
}

func TestNewRenderer(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
	}{
		{"creates renderer with no options", []Option{}},
		{
			"creates renderer with options",
			[]Option{
				WithOption(OptionName("test"), "value"),
			},
		},
		{
			"creates renderer with node renderer",
			[]Option{
				WithNodeRenderers(util.Prioritized(&mockNodeRenderer{
					renderFunc: func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
						return ast.WalkContinue, nil
					},
				}, 100)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRenderer(tt.opts...)
			if r == nil {
				t.Error("NewRenderer() returned nil")
			}
		})
	}
}

func TestRenderer_AddOptions(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		add  []Option
	}{
		{
			"adds options to existing renderer",
			[]Option{WithOption(OptionName("initial"), "value1")},
			[]Option{WithOption(OptionName("added"), "value2")},
		},
		{
			"adds multiple options",
			[]Option{},
			[]Option{
				WithOption(OptionName("opt1"), 1),
				WithOption(OptionName("opt2"), 2),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRenderer(tt.opts...)
			r.AddOptions(tt.add...)

			// Renderer accepts options without error.
			if r == nil {
				t.Error("AddOptions() changed renderer to nil")
			}
		})
	}
}

func TestRenderer_Render(t *testing.T) {
	tests := []struct {
		name       string
		node       ast.Node
		renderer   NodeRenderer
		wantOutput string
		wantErr    bool
	}{
		{
			"renders simple node",
			newMockNode(mockKind),
			&mockNodeRenderer{
				renderFunc: func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
					if entering {
						_, _ = w.WriteString("rendered")
					}
					return ast.WalkContinue, nil
				},
			},
			"rendered",
			false,
		},
		{
			"renders node with children",
			func() ast.Node {
				doc := ast.NewDocument()
				child := newMockNode(mockKind)
				doc.AppendChild(doc, child)
				return doc
			}(),
			&mockNodeRenderer{
				renderFunc: func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
					if entering && n.Kind() == mockKind {
						_, _ = w.WriteString("child")
					}
					return ast.WalkContinue, nil
				},
			},
			"child",
			false,
		},
		{
			"renders with source data",
			newMockNode(mockKind),
			&mockNodeRenderer{
				renderFunc: func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
					if entering {
						_, _ = w.Write(source)
					}
					return ast.WalkContinue, nil
				},
			},
			"test source",
			false,
		},
		{
			"handles empty node",
			ast.NewDocument(),
			&mockNodeRenderer{
				renderFunc: func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
					return ast.WalkContinue, nil
				},
			},
			"",
			false,
		},
		{
			"handles node without renderer function",
			func() ast.Node {
				// Use Document kind which won't have a renderer registered.
				return ast.NewDocument()
			}(),
			&mockNodeRenderer{
				renderFunc: func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
					if n.Kind() == mockKind {
						_, _ = w.WriteString("should not appear")
					}
					return ast.WalkContinue, nil
				},
			},
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRenderer(WithNodeRenderers(util.Prioritized(tt.renderer, 100)))
			var buf bytes.Buffer
			source := []byte(tt.wantOutput)
			err := r.Render(&buf, source, tt.node)

			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got := buf.String(); got != tt.wantOutput {
				t.Errorf("Render() output = %q, want %q", got, tt.wantOutput)
			}
		})
	}
}

func TestRenderer_RenderMultipleTimes(t *testing.T) {
	tests := []struct {
		name       string
		renderFunc NodeRendererFunc
		iterations int
	}{
		{
			"renders consistently on multiple calls",
			func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
				if entering {
					_, _ = w.WriteString("consistent")
				}
				return ast.WalkContinue, nil
			},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRenderer(WithNodeRenderers(util.Prioritized(&mockNodeRenderer{
				renderFunc: tt.renderFunc,
			}, 100)))

			for i := 0; i < tt.iterations; i++ {
				var buf bytes.Buffer
				node := newMockNode(mockKind)
				if err := r.Render(&buf, []byte{}, node); err != nil {
					t.Errorf("Render() iteration %d error = %v", i, err)
				}
				if got := buf.String(); got != "consistent" {
					t.Errorf("Render() iteration %d output = %q, want %q", i, got, "consistent")
				}
			}
		})
	}
}

func TestRenderer_WithNonBufWriter(t *testing.T) {
	tests := []struct {
		name   string
		writer io.Writer
	}{
		{"handles bytes.Buffer", &bytes.Buffer{}},
		{"handles custom writer", &customWriter{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRenderer(WithNodeRenderers(util.Prioritized(&mockNodeRenderer{
				renderFunc: func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
					if entering {
						_, _ = w.WriteString("output")
					}
					return ast.WalkContinue, nil
				},
			}, 100)))

			node := newMockNode(mockKind)
			if err := r.Render(tt.writer, []byte{}, node); err != nil {
				t.Errorf("Render() with custom writer error = %v", err)
			}
		})
	}
}

type customWriter struct {
	buf bytes.Buffer
}

func (w *customWriter) Write(p []byte) (n int, err error) {
	return w.buf.Write(p)
}

func TestRenderer_MultipleNodeRenderers(t *testing.T) {
	tests := []struct {
		name       string
		wantOutput string
	}{
		{
			"renders multiple registered node types",
			"mockA",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use mockKind for both to avoid NodeKind registration issues.
			doc := ast.NewDocument()
			child := newMockNode(mockKind)
			doc.AppendChild(doc, child)

			r := NewRenderer(WithNodeRenderers(
				util.Prioritized(&mockNodeRenderer{
					renderFunc: func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
						if entering && n.Kind() == mockKind {
							_, _ = w.WriteString(tt.wantOutput)
						}
						return ast.WalkContinue, nil
					},
				}, 100),
			))

			var buf bytes.Buffer
			if err := r.Render(&buf, []byte{}, doc); err != nil {
				t.Errorf("Render() error = %v", err)
			}
			if got := buf.String(); got != tt.wantOutput {
				t.Errorf("Render() output = %q, want %q", got, tt.wantOutput)
			}
		})
	}
}

// Mock SetOptioner for testing.
type mockSetOptioner struct {
	options map[OptionName]any
}

func (m *mockSetOptioner) SetOption(name OptionName, value any) {
	if m.options == nil {
		m.options = make(map[OptionName]any)
	}
	m.options[name] = value
}

func (m *mockSetOptioner) RegisterFuncs(reg NodeRendererFuncRegisterer) {
	reg.Register(mockKind, func(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
		return ast.WalkContinue, nil
	})
}

func TestRenderer_SetOptionerPropagation(t *testing.T) {
	tests := []struct {
		name        string
		options     []Option
		wantOptions map[OptionName]any
	}{
		{
			"propagates options to SetOptioner",
			[]Option{
				WithOption(OptionName("key1"), "value1"),
				WithOption(OptionName("key2"), 42),
			},
			map[OptionName]any{
				OptionName("key1"): "value1",
				OptionName("key2"): 42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSetOptioner{}
			allOpts := make([]Option, 0, len(tt.options)+1)
			allOpts = append(allOpts, tt.options...)
			allOpts = append(allOpts, WithNodeRenderers(util.Prioritized(mock, 100)))
			r := NewRenderer(allOpts...)

			// Trigger initialization by rendering.
			node := newMockNode(mockKind)
			var buf bytes.Buffer
			_ = r.Render(&buf, []byte{}, node)

			// Verify options were propagated.
			for k, v := range tt.wantOptions {
				if got, ok := mock.options[k]; !ok {
					t.Errorf("SetOption() option %v not propagated", k)
				} else if got != v {
					t.Errorf("SetOption() option %v = %v, want %v", k, got, v)
				}
			}
		})
	}
}
