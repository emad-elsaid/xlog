package autolink_pages

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func init() {
	RegisterExtension(AutoLinkPages{})
}

type AutoLinkPages struct{}

func (AutoLinkPages) Name() string { return "autolink-pages" }
func (AutoLinkPages) Init() {
	RegisterAutocomplete(autocomplete{})

	if !Config.Readonly {
		Listen(PageChanged, UpdatePagesList)
		Listen(PageDeleted, UpdatePagesList)
	}

	RegisterWidget(WidgetAfterView, 1, backlinksSection)
	RegisterTemplate(templates, "templates")
	MarkDownRenderer.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&pageLinkParser{}, 999),
	))
	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&pageLinkRenderer{}, -1),
	))
}
