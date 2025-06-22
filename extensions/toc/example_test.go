package toc

import (
	"os"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func Example() {
	src := []byte(`
# A section

Hello

# Another section

## A sub-section

### A sub-sub-section

Bye
`)

	markdown := markdown.New()

	// Request that IDs are automatically assigned to headers.
	markdown.Parser().AddOptions(parser.WithAutoHeadingID())
	// Alternatively, we can provide our own implementation of parser.IDs
	// and use,
	//
	//   pctx := parser.NewContext(parser.WithIDs(ids))
	//   doc := parser.Parse(text.NewReader(src), parser.WithContext(pctx))

	doc := markdown.Parser().Parse(text.NewReader(src))

	// Inspect the parsed Markdown document to find headers and build a
	// tree for the table of contents.
	tree, err := Inspect(doc, src)
	if err != nil {
		panic(err)
	}

	if len(tree.Items) == 0 {
		return
		// No table of contents because there are no headers.
	}

	// Render the tree as-is into a Markdown list.
	treeList := RenderList(tree)

	// Render the Markdown list into HTML.
	if err := markdown.Renderer().Render(os.Stdout, src, treeList); err != nil {
		panic(err)
	}

	// Output:
	// <ul>
	// <li>
	// <a href="#a-section">A section</a></li>
	// <li>
	// <a href="#another-section">Another section</a><ul>
	// <li>
	// <a href="#a-sub-section">A sub-section</a><ul>
	// <li>
	// <a href="#a-sub-sub-section">A sub-sub-section</a></li>
	// </ul>
	// </li>
	// </ul>
	// </li>
	// </ul>
}
