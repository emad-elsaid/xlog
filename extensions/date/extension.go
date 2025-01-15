package date

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func init() {
	RegisterExtension(Date{})
}

type Date struct{}

func (Date) Name() string { return "date" }
func (Date) Init() {
	RegisterTemplate(templates, "templates")
	RegisterLink(links)
	RegisterBuildPage(`/+/calendar`, true)

	Get(`/+/date/{date}`, dateHandler)
	Get(`/+/calendar`, calendarHandler)

	MarkdownConverter().Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&dateParser{}, 999),
	))
	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&dateRenderer{}, 0),
	))
}
