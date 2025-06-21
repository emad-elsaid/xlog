package mathjax

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func init() {
	RegisterExtension(Mathjax{})
}

type Mathjax struct{}

func (Mathjax) Name() string { return "mathjax" }
func (Mathjax) Init() {
	RegisterStaticDir(js)
	registerBuildFiles()
	MarkdownConverter().Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&inlineMathParser{}, 999),
		),
		parser.WithBlockParsers(
			util.Prioritized(&mathJaxBlockParser{}, 999),
		),
	)
	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&InlineMathRenderer{startDelim: `\(`, endDelim: `\)`}, 0),
		util.Prioritized(&MathBlockRenderer{startDelim: `\[`, endDelim: `\]`}, 0),
	))
}
