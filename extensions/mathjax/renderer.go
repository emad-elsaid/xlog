package mathjax

import (
	"bytes"
	"embed"
	"html/template"

	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

//go:embed js
var js embed.FS

const script = `
<script>
MathJax = {
  tex: {
    displayMath: [['$$', '$$'], ['\\[', '\\]']],
    inlineMath: [['$', '$'], ['\\(', '\\)']]
  },
  svg: {fontCache: 'global'}
};
</script>
<script type="text/javascript" src="/js/tex-chtml-full.js" async></script>`

func init() {
	RegisterStaticDir(js)
	RegisterWidget(HEAD_WIDGET, 1, headScript)
	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&InlineMathRenderer{startDelim: `\(`, endDelim: `\)`}, 0),
		util.Prioritized(&MathBlockRenderer{startDelim: `\[`, endDelim: `\]`}, 0),
	))
}

func headScript(_ Page) template.HTML {
	return template.HTML(script)
}

type InlineMathRenderer struct {
	startDelim string
	endDelim   string
}

func (r *InlineMathRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindInlineMath, r.renderInlineMath)
}

func (r *InlineMathRenderer) renderInlineMath(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(`<span>` + r.startDelim)
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			segment := c.(*ast.Text).Segment
			value := segment.Value(source)
			if bytes.HasSuffix(value, []byte("\n")) {
				w.Write(value[:len(value)-1])
				if c != n.LastChild() {
					w.Write([]byte(" "))
				}
			} else {
				w.Write(value)
			}
		}
		return ast.WalkSkipChildren, nil
	}
	w.WriteString(r.endDelim + `</span>`)
	return ast.WalkContinue, nil
}

type MathBlockRenderer struct {
	startDelim string
	endDelim   string
}

func (r *MathBlockRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindMathBlock, r.renderMathBlock)
}

func (r *MathBlockRenderer) renderMathBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*MathBlock)
	if entering {
		_, _ = w.WriteString(`<p>` + r.startDelim)
		l := n.Lines().Len()
		for i := 0; i < l; i++ {
			line := n.Lines().At(i)
			w.Write(line.Value(source))
		}
	} else {
		_, _ = w.WriteString(r.endDelim + `</p>` + "\n")
	}
	return ast.WalkContinue, nil
}
