package images

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func init() {
	RegisterExtension(Images{})
}

type Images struct{}

func (Images) Name() string { return "images" }
func (Images) Init() {
	MarkdownConverter().Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(columnizeImagesParagraph{}, 0),
		),
	)
	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&imagesColumnsRenderer{}, 0),
	))
}
