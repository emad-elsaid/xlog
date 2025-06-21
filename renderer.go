package xlog

import (
	"sync"

	chroma_html "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/emoji"
	"github.com/emad-elsaid/xlog/markdown/extension"
	"github.com/emad-elsaid/xlog/markdown/highlighting"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
)

// The instance of markdown renderer. this is what takes the page content and
// converts it to HTML. it defines what features to use from goldmark and what
// options to turn on
var MarkdownConverter = sync.OnceValue(func() markdown.Markdown {
	return markdown.New(
		markdown.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithCustomStyle(styles.Get(Config.CodeStyle)),
				highlighting.WithFormatOptions(
					chroma_html.WithLineNumbers(true),
				),
			),
			extension.Typographer,
			emoji.Emoji,
		),

		markdown.WithParserOptions(
			parser.WithAutoHeadingID(),
		),

		markdown.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)
})

// FindInAST takes an AST node and walks the tree depth first
// searching for a node of a specific type can be used to find first image,
// link, paragraph...etc
func FindInAST[t ast.Node](n ast.Node) (found t, ok bool) {
	if n == nil {
		return
	}

	ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if found, ok = n.(t); ok {
			return ast.WalkStop, nil
		}

		return ast.WalkContinue, nil
	})

	return
}

// Extract all nodes of a specific type from the AST
func FindAllInAST[t ast.Node](n ast.Node) (a []t) {
	if n == nil {
		return
	}

	ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		if casted, ok := n.(t); ok {
			a = append(a, casted)
		}
		return ast.WalkContinue, nil
	})

	return
}
