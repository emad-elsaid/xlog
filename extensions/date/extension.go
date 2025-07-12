package date

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func init() {
	app := GetApp()
	app.RegisterExtension(Date{})
}

type Date struct{}

func (Date) Name() string { return "date" }
func (Date) Init(app *App) {
	app.RegisterTemplate(templates, "templates")
	app.RegisterLink(links)
	app.RegisterBuildPage(`/+/calendar`, true)

	app.Get(`/+/date/{date}`, dateHandler)
	app.Get(`/+/calendar`, calendarHandler)

	MarkdownConverter().Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&dateParser{}, 999),
	))
	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&dateRenderer{}, 0),
	))
}
