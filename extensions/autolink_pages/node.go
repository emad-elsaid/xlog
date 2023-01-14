package autolink_pages

import (
	"fmt"

	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/ast"
)

var KindPageLink = ast.NewNodeKind("PageLink")

type PageLink struct {
	ast.BaseInline
	page  Page
	url   string
	value *ast.Text
}

func (_ *PageLink) Kind() ast.NodeKind {
	return KindPageLink
}

func (p *PageLink) Dump(source []byte, level int) {
	m := map[string]string{
		"value": fmt.Sprintf("%#v:%s", p.value, p.url),
	}
	ast.DumpHelper(p, source, level, m, nil)
}
