// Package toc provides support for building a Table of Contents from a
// goldmark Markdown document.
//
// The package operates in two stages: inspection and rendering. During
// inspection, the package analyzes an existing Markdown document, and builds
// a Table of Contents from it.
//
//	markdown := markdown.New(...)
//
//	parser := markdown.Parser()
//	doc := parser.Parse(text.NewReader(src))
//	tocTree, err := toc.Inspect(doc, src)
//
// During rendering, it converts the Table of Contents into a list of headings
// with nested items under each as a goldmark Markdown document. You may
// manipulate the TOC, removing items from it or simplifying it, before
// rendering.
//
//	if len(tocTree.Items) == 0 {
//		// No headings in the document.
//		return
//	}
//	tocList := toc.RenderList(tocTree)
//
// You can render that Markdown document using goldmark into whatever form you
// prefer.
//
//	renderer := markdown.Renderer()
//	renderer.Render(out, src, tocList)
//
// The following diagram summarizes the flow of information with goldmark-toc.
//
//	   src
//	+--------+                           +-------------------+
//	|        |   goldmark/Parser.Parse   |                   |
//	| []byte :---------------------------> goldmark/ast.Node |
//	|        |                           |                   |
//	+---.----+                           +-------.-----.-----+
//	    |                                        |     |
//	    '----------------.     .-----------------'     |
//	                      \   /                        |
//	                       \ /                         |
//	                        |                          |
//	                        | toc.Inspect              |
//	                        |                          |
//	                   +----v----+                     |
//	                   |         |                     |
//	                   | toc.TOC |                     |
//	                   |         |                     |
//	                   +----.----+                     |
//	                        |                          |
//	                        | toc/Renderer.Render      |
//	                        |                          |
//	              +---------v---------+                |
//	              |                   |                |
//	              | goldmark/ast.Node |                |
//	              |                   |                |
//	              +---------.---------+                |
//	                        |                          |
//	                        '-------.   .--------------'
//	                                 \ /
//	                                  |
//	         goldmark/Renderer.Render |
//	                                  |
//	                                  v
//	                              +------+
//	                              | HTML |
//	                              +------+
package toc
