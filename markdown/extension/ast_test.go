package extension

import (
	"bytes"
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
	"github.com/emad-elsaid/xlog/markdown/testutil"
	"github.com/emad-elsaid/xlog/markdown/text"
)

// getNodeText extracts text from a node without using deprecated Text() method.
func getNodeText(node ast.Node, source []byte) []byte {
	// Handle specific node types with their own text extraction methods
	switch n := node.(type) {
	case *ast.Text:
		return n.Value(source)
	case *ast.String:
		return n.Value
	case *ast.AutoLink:
		return n.Label(source)
	case *ast.RawHTML:
		return n.Segments.Value(source)
	}

	// Try to use Lines() method for block nodes only (check Type() instead of interface)
	if node.Type() == ast.TypeBlock {
		if linesNode, ok := node.(interface{ Lines() *text.Segments }); ok {
			lines := linesNode.Lines()
			if lines != nil && lines.Len() > 0 {
				return lines.Value(source)
			}
		}
	}

	// For nodes with children (like CodeSpan, Emphasis, Link, Blockquote, List, etc.)
	// recursively collect text
	var buf bytes.Buffer
	for c := node.FirstChild(); c != nil; c = c.NextSibling() {
		buf.Write(getNodeText(c, source))
		// Add newline for soft line breaks
		if sb, ok := c.(interface{ SoftLineBreak() bool }); ok && sb.SoftLineBreak() {
			buf.WriteByte('\n')
		}
	}
	return buf.Bytes()
}

func TestASTBlockNodeText(t *testing.T) {
	var cases = []struct {
		Name   string
		Source string
		T1     string
		T2     string
		C      bool
	}{
		{
			Name: "DefinitionList",
			Source: `c1
:   c2
    c3

a

c4
:   c5
    c6`,
			T1: `c1c2
c3`,
			T2: `c4c5
c6`,
		},
		{
			Name: "Table",
			Source: `| h1 | h2 |
| -- | -- |
| c1 | c2 |

a


| h3 | h4 |
| -- | -- |
| c3 | c4 |`,

			T1: `h1h2c1c2`,
			T2: `h3h4c3c4`,
		},
	}

	for _, cs := range cases {
		t.Run(cs.Name, func(t *testing.T) {
			s := []byte(cs.Source)
			md := markdown.New(
				markdown.WithRendererOptions(
					html.WithUnsafe(),
				),
				markdown.WithExtensions(
					DefinitionList,
					Table,
				),
			)
			n := md.Parser().Parse(text.NewReader(s))
			c1 := n.FirstChild()
			c2 := c1.NextSibling().NextSibling()
			if cs.C {
				c1 = c1.FirstChild()
				c2 = c2.FirstChild()
			}
			if !bytes.Equal(getNodeText(c1, s), []byte(cs.T1)) {

				t.Errorf("%s unmatch:\n%s", cs.Name, testutil.DiffPretty(getNodeText(c1, s), []byte(cs.T1)))

			}
			if !bytes.Equal(getNodeText(c2, s), []byte(cs.T2)) {

				t.Errorf("%s(EOF) unmatch: %s", cs.Name, testutil.DiffPretty(getNodeText(c2, s), []byte(cs.T2)))

			}
		})
	}

}

func TestASTInlineNodeText(t *testing.T) {
	var cases = []struct {
		Name   string
		Source string
		T1     string
	}{
		{
			Name:   "Strikethrough",
			Source: `~c1 *c2*~`,
			T1:     `c1 c2`,
		},
	}

	for _, cs := range cases {
		t.Run(cs.Name, func(t *testing.T) {
			s := []byte(cs.Source)
			md := markdown.New(
				markdown.WithRendererOptions(
					html.WithUnsafe(),
				),
				markdown.WithExtensions(
					Strikethrough,
				),
			)
			n := md.Parser().Parse(text.NewReader(s))
			c1 := n.FirstChild().FirstChild()
			if !bytes.Equal(getNodeText(c1, s), []byte(cs.T1)) {

				t.Errorf("%s unmatch:\n%s", cs.Name, testutil.DiffPretty(getNodeText(c1, s), []byte(cs.T1)))

			}
		})
	}

}
