package todo

import (
	"fmt"

	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func init() {
	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&TaskCheckBoxHTMLRenderer{}, 0),
	))
}

type TaskCheckBoxHTMLRenderer struct{}

func (r *TaskCheckBoxHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(east.KindTaskCheckBox, r.renderTaskCheckBox)
}

func (r *TaskCheckBoxHTMLRenderer) renderTaskCheckBox(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*east.TaskCheckBox)
	p := n.Parent()

	w.WriteString(`<input type="checkbox" `)

	if n.IsChecked {
		w.WriteString(`checked="" `)
	}

	if READONLY {
		w.WriteString(`disabled="" `)
	} else if p.Kind() == ast.KindTextBlock {
		if l := p.Lines(); l != nil {
			w.WriteString(fmt.Sprintf(`data-pos="%d" `, l.At(0).Start))
		}
	}

	w.WriteString("> ")
	return ast.WalkContinue, nil
}
