package todo

import (
	"fmt"
	"html/template"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	east "github.com/emad-elsaid/xlog/markdown/extension/ast"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

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

	_, _ = w.WriteString(`<input name="checked" type="checkbox" `)

	if n.IsChecked {
		_, _ = w.WriteString(`checked="" `)
	}

	if Config.Readonly {
		_, _ = w.WriteString(`disabled="" `)
	} else if p.Kind() == ast.KindTextBlock {
		if l := p.Lines(); l != nil {

			vals := fmt.Sprintf(`{"page": decodeURI(document.location.pathname.substr(1)), "pos": %d}
`,
				l.At(0).Start,
			)

			_, _ = fmt.Fprintf(w, `hx-post="/+/todo" hx-vals="js:%s"`, template.HTMLEscapeString(vals))
		}
	}

	_, _ = w.WriteString("> ")
	return ast.WalkContinue, nil
}
