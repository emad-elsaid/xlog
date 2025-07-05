package html

import (
	"flag"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
	"github.com/emad-elsaid/xlog/markdown/util"
)

var allowHTML bool

func init() {
	flag.BoolVar(&allowHTML, "html", false, "Allow HTML in markdown")
	app := xlog.GetApp()
	app.RegisterExtension(HTML{})
}

type HTML struct{}

func (HTML) Name() string { return "html" }
func (HTML) Init() {
	if allowHTML {
		xlog.MarkdownConverter().Parser().AddOptions(parser.WithBlockParsers(
			util.Prioritized(parser.NewHTMLBlockParser(), 0),
		))
		xlog.MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
			util.Prioritized(html.NewRenderer(), 0),
		))
	}
}
