package shortcode

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func init() {
	RegisterExtension(ShortCodeEx{})
}

type ShortCodeEx struct{}

func (ShortCodeEx) Name() string { return "shortcode" }
func (ShortCodeEx) Init() {
	MarkdownConverter().Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(&shortCodeParser{}, 0),
	))
	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&shortCodeRenderer{}, 0),
	))
	MarkdownConverter().Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(transformShortCodeBlocks{}, 0),
		),
	)
}
