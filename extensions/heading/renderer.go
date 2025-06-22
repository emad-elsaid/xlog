package heading

import (
	"fmt"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func init() {
	RegisterExtension(Heading{})
}

type Heading struct{}

func (Heading) Name() string { return "heading" }
func (Heading) Init() {
	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&headingRenderer{}, 0),
	))
}

type headingRenderer struct{}

func (s *headingRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindHeading, s.render)
}

func (s *headingRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.Heading)
	if entering {
		_, _ = w.WriteString("<h")
		_ = w.WriteByte("0123456"[n.Level])
		if n.Attributes() != nil {
			html.RenderAttributes(w, node, html.HeadingAttributeFilter)
		}
		_ = w.WriteByte('>')
	} else {

		if id, ok := node.AttributeString("id"); ok {
			w.WriteString(fmt.Sprintf(` <a class="show-on-parent-hover is-hidden has-text-grey" href="#%s">¶</a>`, id))
		}

		_, _ = w.WriteString("</h")
		_ = w.WriteByte("0123456"[n.Level])
		_, _ = w.WriteString(">\n")
	}
	return ast.WalkContinue, nil
}
