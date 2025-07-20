package todo

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func init() {
	RegisterExtension(TODO{})
}

type TODO struct{}

func (TODO) Name() string { return "todo" }
func (TODO) Init(app *App) {
	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&TaskCheckBoxHTMLRenderer{}, 0),
	))

	if !app.GetConfig().Readonly {
		app.RequireHTMX()
		app.Post(`/+/todo`, toggleHandler)
	}
}
