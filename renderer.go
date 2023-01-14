package xlog

import (
	chroma_html "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// The instance of markdown renderer. this is what takes the page content and
// converts it to HTML. it defines what features to use from goldmark and what
// options to turn on
var MarkDownRenderer = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.DefinitionList,
		extension.Footnote,
		extension.Typographer,
		highlighting.NewHighlighting(
			highlighting.WithCustomStyle(styles.Dracula),
			highlighting.WithFormatOptions(
				chroma_html.WithLineNumbers(true),
			),
		),
		emoji.Emoji,
	),

	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithUnsafe(),
	),
)

// This is a function that takes an AST node and walks the tree depth first
// searching for a node of a specific type can be used to find first image,
// link, paragraph...etc
func FindInAST[t ast.Node](n ast.Node) (found t, ok bool) {
	ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if casted, success := n.(t); success {
			found = casted
			ok = true
			return ast.WalkStop, nil
		}

		return ast.WalkContinue, nil
	})

	return
}

// Extract all nodes of a specific type from the AST
func FindAllInAST[t ast.Node](n ast.Node) (a []t) {
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
