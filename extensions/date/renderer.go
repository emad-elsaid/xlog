package date

import (
	"fmt"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

type dateRenderer struct{}

func (s *dateRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindDate, s.render)
}

func (s *dateRenderer) render(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	node, ok := n.(*DateNode)
	if !ok {
		return ast.WalkContinue, nil
	}

	path := fmt.Sprintf(`/+/date/%s`, node.time.Format("2-1-2006"))
	RegisterBuildPage(path, true)

	fmt.Fprintf(w, ` <a href="%s" class="tag"><span class="icon"><i class="fa-regular fa-clock"></i></span><span>%s<span></a> `, path, node.time.Format("2 January 2006"))

	return ast.WalkContinue, nil
}
