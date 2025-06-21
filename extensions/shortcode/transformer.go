package shortcode

import (
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/text"
)

type transformShortCodeBlocks struct{}

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

			if _, ok := shortcodes[string(n.Language(source))]; !ok {
				continue
			}

			blocks = append(blocks, n)
		}

		return ast.WalkContinue, nil
	})

	for _, b := range blocks {
		replacement := ShortCodeBlock{
			FencedCodeBlock: *b,
			fun:             shortcodes[string(b.Language(source))],
		}

		parent := b.Parent()
		parent.ReplaceChild(parent, b, &replacement)
	}
}
