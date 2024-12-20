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
	Get(`/+/date/{date}`, dateHandler)
	MarkdownConverter().Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&dateParser{}, 999),
	))
	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&dateRenderer{}, 0),
	))
}
