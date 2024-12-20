package images

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
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
