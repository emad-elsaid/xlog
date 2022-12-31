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
// recursively calling itself in search for a node of a specific kind
// can be used to find first image, link, paragraph...etc
func FindInAST[t ast.Node](n ast.Node, kind ast.NodeKind) (found t, ok bool) {
	if n.Kind() == kind {
		if found, ok := n.(t); ok {
			return found, true
		}
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if a, ok := FindInAST[t](c, kind); ok {
			return a, true
		}
	}

	return found, false
}

// Extract all nodes of a specific type from the AST
func FindAllInAST[t ast.Node](n ast.Node, kind ast.NodeKind) (a []t) {
	if n.Kind() == kind {
		typed, _ := n.(t)
		a = []t{typed}
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		a = append(a, FindAllInAST[t](c, kind)...)
	}

	return a
}
