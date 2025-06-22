package autolink_pages

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func init() {
	RegisterExtension(AutoLinkPages{})
}

type AutoLinkPages struct{}

func (AutoLinkPages) Name() string { return "autolink-pages" }
func (AutoLinkPages) Init() {
	if !Config.Readonly {
		Listen(PageChanged, UpdatePagesList)
		Listen(PageDeleted, UpdatePagesList)
	}

	RegisterWidget(WidgetAfterView, 1, backlinksSection)
	RegisterTemplate(templates, "templates")
	MarkdownConverter().Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&pageLinkParser{}, 999),
	))
	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&pageLinkRenderer{}, -1),
	))
}
