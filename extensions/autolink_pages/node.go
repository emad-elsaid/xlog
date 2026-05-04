package autolink_pages

import (
	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
)

var KindPageLink = ast.NewNodeKind("PageLink")

type PageLink struct {
	ast.BaseInline
	page xlog.Page
}

// Page returns the linked page.
func (p *PageLink) Page() xlog.Page {
	return p.page
}

func (*PageLink) Kind() ast.NodeKind {
	return KindPageLink
}

func (p *PageLink) Dump(source []byte, level int) {
	m := map[string]string{
		"value": p.page.Name(),
	}
	ast.DumpHelper(p, source, level, m, nil)
}
