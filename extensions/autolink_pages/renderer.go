package autolink_pages

import (
	"fmt"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

type pageLinkRenderer struct{}

func (h *pageLinkRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindPageLink, render)
}

func render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		n := node.(*PageLink)
		url := n.page.Name()

		fmt.Fprintf(w,
			`<a href="/%s">`,
			util.EscapeHTML(util.URLEscape([]byte([]byte(url)), false)),
		)

		if total, done := countTodos(n.page); total > 0 {
			isDone := ""
			if total == done {
				isDone = "is-success"
			}
			fmt.Fprintf(w, `<span class="tag is-rounded %s">%d/%d</span> `, isDone, done, total)
		}
	} else {
		w.WriteString(`</a>`)
	}

	return ast.WalkContinue, nil
}
