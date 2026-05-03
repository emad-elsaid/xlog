package autolink_pages

import (
	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func init() {
	xlog.RegisterExtension(AutoLinkPages{})
}

type AutoLinkPages struct{}

func (AutoLinkPages) Name() string { return "autolink-pages" }
func (AutoLinkPages) Init() {
	if !xlog.Config.Readonly {
		xlog.Listen(xlog.PageChanged, UpdatePagesList)
		xlog.Listen(xlog.PageDeleted, UpdatePagesList)
	}

	xlog.RegisterWidget(xlog.WidgetAfterView, 1, backlinksSection)
	xlog.RegisterTemplate(templates, "templates")
	xlog.MarkdownConverter().Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&pageLinkParser{}, 999),
	))
	xlog.MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&pageLinkRenderer{}, -1),
	))
}
