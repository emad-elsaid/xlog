package autolink_pages

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/ast"
)

var KindPageLink = ast.NewNodeKind("PageLink")

type PageLink struct {
	ast.BaseInline
	page Page
}

func (_ *PageLink) Kind() ast.NodeKind {
	return KindPageLink
}

func (p *PageLink) Dump(source []byte, level int) {
	m := map[string]string{
		"value": p.page.Title(),
	}
	ast.DumpHelper(p, source, level, m, nil)
}
