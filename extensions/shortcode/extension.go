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
	RegisterAutocomplete(autocomplete{})
	MarkDownRenderer.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(&shortCodeParser{}, 0),
	))
	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&shortCodeRenderer{}, 0),
	))
	MarkDownRenderer.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(transformShortCodeBlocks{}, 0),
		),
	)
}
