package sql_table

import (
	"embed"
	"html/template"

	"github.com/emad-elsaid/xlog"
)

//go:embed js
var js embed.FS

func init() {
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
	o, _ := js.ReadFile("js/sql_table.html")
	return template.HTML(o)
}
