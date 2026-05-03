package images

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

type imagesColumnsRenderer struct{}

func (s *imagesColumnsRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindColumns, s.render)
}

func (s *imagesColumnsRenderer) render(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(`<div class="columns">`)

		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			_, _ = w.WriteString(`<div class="column">`)
			_ = MarkdownConverter().Renderer().Render(w, source, c)
			_, _ = w.WriteString(`</div>`)
		}

	} else {
		_, _ = w.WriteString(`</div>`)
	}

	return ast.WalkSkipChildren, nil
}
