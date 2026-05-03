package shortcode

import (
	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func init() {
	xlog.RegisterExtension(ShortCodeEx{})
}

type ShortCodeEx struct{}

func (s ShortCodeEx) Name() string { return "shortcode" }
func (s ShortCodeEx) Init() {
	s.Extend(xlog.MarkdownConverter())
}

func (s ShortCodeEx) Extend(md markdown.Markdown) {
	md.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(&shortCodeParser{}, 0),
	))
	md.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&shortCodeRenderer{}, 0),
	))
	md.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(transformShortCodeBlocks{}, 0),
		),
	)
}
