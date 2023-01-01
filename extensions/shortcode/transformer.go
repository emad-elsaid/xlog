package shortcode

import (
	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func init() {
	MarkDownRenderer.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(transformShortCodeBlocks(0), 0),
		),
	)
}

type transformShortCodeBlocks int

func (t transformShortCodeBlocks) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()
	blocks := []*ast.FencedCodeBlock{}

	ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		for c := node.FirstChild(); c != nil; c = c.NextSibling() {
			n, ok := c.(*ast.FencedCodeBlock)
			if !ok {
				continue
			}

			lang := string(n.Language(source))
			if _, ok := shortcodes[lang]; !ok {
				continue
			}

			blocks = append(blocks, n)
		}

		return ast.WalkContinue, nil
	})

	for _, b := range blocks {
		lang := string(b.Language(source))

		replacement := ShortCodeBlock{
			FencedCodeBlock: *b,
			fun:             shortcodes[lang],
		}

		parent := b.Parent()
		parent.ReplaceChild(parent, b, &replacement)
	}
}
