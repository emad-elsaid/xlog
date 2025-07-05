package rtl

import (
	"html/template"

	. "github.com/emad-elsaid/xlog"
)

func init() {
	app := GetApp()
	app.RegisterExtension(RTL{})
}

type RTL struct{}

func (RTL) Name() string { return "rtl" }
func (RTL) Init() {
	app := GetApp()
	app.RegisterWidget(WidgetHead, 0, rtlWidget)
}

func rtlWidget(Page) template.HTML {
	return template.HTML(`<link rel="stylesheet" href="/public/rtl.css">`)
}
