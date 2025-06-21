package toc

import (
	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/util"
)

// Extender extends a Goldmark Markdown parser and renderer to always include
// a table of contents in the output.
//
// To use this, install it into your Goldmark Markdown object.
//
//	md := markdown.New(
//	  // ...
//	  markdown.WithParserOptions(parser.WithAutoHeadingID()),
//	  markdown.WithExtensions(
//	    // ...
//	    &toc.Extender{
//	    },
//	  ),
//	)
//
// This will install the default Transformer. For more control, install the
// Transformer directly on the Markdown Parser.
//
// NOTE: Unless you've supplied your own parser.IDs implementation, you'll
// need to enable the WithAutoHeadingID option on the parser to generate IDs
// and links for headings.
type Extender struct {
	// Title is the title of the table of contents section.
	// Defaults to "Table of Contents" if unspecified.
	Title string

	// TitleDepth is the heading depth for the Title.
	// Defaults to 1 (<h1>) if unspecified.
	TitleDepth int

	// MinDepth is the minimum depth of the table of contents.
	// Headings with a level lower than the specified depth will be ignored.
	// See the documentation for MinDepth for more information.
	//
	// Defaults to 0 (no limit) if unspecified.
	MinDepth int

	// MaxDepth is the maximum depth of the table of contents.
	// Headings with a level greater than the specified depth will be ignored.
	// See the documentation for MaxDepth for more information.
	//
	// Defaults to 0 (no limit) if unspecified.
	MaxDepth int

	// ListID is the id for the list of TOC items rendered in the HTML.
	//
	// See the documentation for Transformer.ListID for more information.
	ListID string

	// TitleID is the id for the Title heading rendered in the HTML.
	//
	// See the documentation for Transformer.TitleID for more information.
	TitleID string

	// Compact controls whether empty items should be removed
	// from the table of contents.
	//
	// See the documentation for Compact for more information.
	Compact bool
}

// Extend adds support for rendering a table of contents to the provided
// Markdown parser/renderer.
func (e *Extender) Extend(md markdown.Markdown) {
	md.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{
				Title:      e.Title,
				TitleDepth: e.TitleDepth,
				MinDepth:   e.MinDepth,
				MaxDepth:   e.MaxDepth,
				ListID:     e.ListID,
				TitleID:    e.TitleID,
				Compact:    e.Compact,
			}, 100),
		),
	)
}
