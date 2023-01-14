package autolink_pages

import (
	"fmt"

	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func init() {
	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&pageLinkRenderer{}, -1),
	))
}

type pageLinkRenderer struct{}

func (h *pageLinkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindPageLink, render)
}

func render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*PageLink)
	if !entering {
		return ast.WalkContinue, nil
	}

	w.WriteString(`<a href="`)
	url := []byte(n.url)
	label := n.value.Text(source)

	w.Write(util.EscapeHTML(util.URLEscape(url, false)))
	w.WriteString(`">`)

	if total, done := countTodos(n.page); total > 0 {
		isDone := ""
		if total == done {
			isDone = "is-success"
		}
		fmt.Fprintf(w, `<span class="tag is-rounded %s">%d/%d</span> `, isDone, done, total)
	}

	w.Write(util.EscapeHTML(label))
	w.WriteString(`</a>`)
	return ast.WalkContinue, nil
}
