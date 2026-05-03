package date

import (
	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func init() {
	xlog.RegisterExtension(Date{})
}

type Date struct{}

func (Date) Name() string { return "date" }
func (Date) Init() {
	xlog.RegisterTemplate(templates, "templates")
	xlog.RegisterLink(links)
	xlog.RegisterBuildPage(`/+/calendar`, true)

	xlog.Get(`/+/date/{date}`, dateHandler)
	xlog.Get(`/+/calendar`, calendarHandler)

	xlog.MarkdownConverter().Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&dateParser{}, 999),
	))
	xlog.MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&dateRenderer{}, 0),
	))
}
