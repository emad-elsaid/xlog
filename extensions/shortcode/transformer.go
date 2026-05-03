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

			shortcodesMutex.RLock()
			_, exists := shortcodes[string(n.Language(source))]
			shortcodesMutex.RUnlock()

			if !exists {
				continue
			}

			blocks = append(blocks, n)
		}

		return ast.WalkContinue, nil
	})

	for _, b := range blocks {
		shortcodesMutex.RLock()
		fn := shortcodes[string(b.Language(source))]
		shortcodesMutex.RUnlock()

		replacement := ShortCodeBlock{
			FencedCodeBlock: *b,
			fun:             fn,
		}

		parent := b.Parent()
		parent.ReplaceChild(parent, b, &replacement)
	}
}
