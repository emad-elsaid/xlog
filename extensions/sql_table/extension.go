package sql_table

import (
	"embed"
	"flag"
	"fmt"
	"html/template"

	"github.com/emad-elsaid/types"
	"github.com/emad-elsaid/xlog"
	east "github.com/emad-elsaid/xlog/markdown/extension/ast"
)

//go:embed js
var js embed.FS

var sqlTableThreshold int

func init() {
	flag.IntVar(&sqlTableThreshold, "sql-table.threshold", 100, "If a table rows is more than this threshold it'll allow users to query it with SQL")
	xlog.RegisterExtension(Extension{})
}

type Extension struct{}

func (Extension) Name() string {
	return "sql_table"
}

func (Extension) Init() {
	xlog.RegisterWidget(xlog.WidgetAfterView, 1, script)
}

func script(p xlog.Page) template.HTML {
	if p == nil {
		return ""
	}

	_, a := p.AST()
	if a == nil {
		return ""
	}

	tables := xlog.FindAllInAST[*east.Table](a)
	if len(tables) == 0 {
		return ""
	}

	largeTableFound := types.Slice[*east.Table](tables).Any(func(t *east.Table) bool {
		return len(xlog.FindAllInAST[*east.TableRow](t)) >= sqlTableThreshold
	})
	if !largeTableFound {
		return ""
	}

	o, _ := js.ReadFile("js/sql_table.html")
	o = append(o, []byte(fmt.Sprintf("<script>const sqlTableThreshold = %d</script>", sqlTableThreshold))...)

	return template.HTML(o)
}
