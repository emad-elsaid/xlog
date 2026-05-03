package todo

import (
	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func init() {
	xlog.RegisterExtension(TODO{})
}

type TODO struct{}

func (TODO) Name() string { return "todo" }
func (TODO) Init() {
	xlog.MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&TaskCheckBoxHTMLRenderer{}, 0),
	))

	if !xlog.Config.Readonly {
		xlog.RequireHTMX()
		xlog.Post(`/+/todo`, toggleHandler)
	}
}
