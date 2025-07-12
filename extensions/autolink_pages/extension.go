package autolink_pages

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func init() {
	app := GetApp()
	app.RegisterExtension(AutoLinkPages{})
}

type AutoLinkPages struct{}

func (AutoLinkPages) Name() string { return "autolink-pages" }
func (AutoLinkPages) Init(app *App) {
	if !app.GetConfig().Readonly {
		app.Listen(PageChanged, UpdatePagesList)
		app.Listen(PageDeleted, UpdatePagesList)
	}

	app.RegisterWidget(WidgetAfterView, 1, backlinksSection)
	app.RegisterTemplate(templates, "templates")
	MarkdownConverter().Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&pageLinkParser{}, 999),
	))
	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&pageLinkRenderer{}, -1),
	))
}
