package markdown_test

import (
	"bytes"
	"testing"

	. "github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/ast"
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
			Name: "AtxHeading",
			Source: `# l1

a

# l2`,
			T1: `l1`,
			T2: `l2`,
		},
		{
			Name: "SetextHeading",
			Source: `l1
l2
===============

a

l3
l4
==============`,
			T1: `l1
l2`,
			T2: `l3
l4`,
		},
		{
			Name: "CodeBlock",
			Source: `    l1
    l2

a

    l3
	l4`,
			T1: `l1
l2
`,
			T2: `l3
l4
`,
		},
		{
			Name: "FencedCodeBlock",
			Source: "```" + `
l1
l2
` + "```" + `

a

` + "```" + `
l3
l4`,
			T1: `l1
l2
`,
			T2: `l3
l4
`,
		},
		{
			Name: "Blockquote",
			Source: `> l1
> l2

a

> l3
> l4`,
			T1: `l1
l2`,
			T2: `l3
l4`,
		},
		{
			Name: "List",
			Source: `- l1
  l2

a

- l3
  l4`,
			T1: `l1
l2`,
			T2: `l3
l4`,
			C: true,
		},
		{
			Name: "HTMLBlock",
			Source: `<div>
l1
l2
</div>

a

<div>
l3
l4`,
			T1: `<div>
l1
l2
</div>
`,
			T2: `<div>
l3
l4`,
		},
	}

	for _, cs := range cases {
		t.Run(cs.Name, func(t *testing.T) {
			s := []byte(cs.Source)
			md := New()
			n := md.Parser().Parse(text.NewReader(s))
			c1 := n.FirstChild()
			c2 := c1.NextSibling().NextSibling()
			if cs.C {
				c1 = c1.FirstChild()
				c2 = c2.FirstChild()
			}
			if !bytes.Equal(getNodeText(c1, s), []byte(cs.T1)) {

				t.Errorf("%s unmatch: %s", cs.Name, testutil.DiffPretty(getNodeText(c1, s), []byte(cs.T1)))

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
			Name:   "CodeSpan",
			Source: "`c1`",
			T1:     `c1`,
		},
		{
			Name:   "Emphasis",
			Source: `*c1 **c2***`,
			T1:     `c1 c2`,
		},
		{
			Name:   "Link",
			Source: `[label](url)`,
			T1:     `label`,
		},
		{
			Name:   "AutoLink",
			Source: `<http://url>`,
			T1:     `http://url`,
		},
		{
			Name:   "RawHTML",
			Source: `<span>c1</span>`,
			T1:     `<span>`,
		},
	}

	for _, cs := range cases {
		t.Run(cs.Name, func(t *testing.T) {
			s := []byte(cs.Source)
			md := New()
			n := md.Parser().Parse(text.NewReader(s))
			c1 := n.FirstChild().FirstChild()
			if !bytes.Equal(getNodeText(c1, s), []byte(cs.T1)) {
				t.Errorf("%s unmatch:\n%s", cs.Name, testutil.DiffPretty(getNodeText(c1, s), []byte(cs.T1)))
			}
		})
	}

}

func TestHasBlankPreviousLines(t *testing.T) {
	var cases = []struct {
		Name     string
		Source   string
		Node     func(n ast.Node) ast.Node
		Expected bool
	}{
		{
			Name: "nesting paragraphs in blockquotes",
			Source: `
> a
> 
> b
`,
			Node: func(n ast.Node) ast.Node {
				return n.FirstChild().FirstChild().NextSibling()
			},
			Expected: true,
		},
		{
			Name: "nesting HTML blocks in blockquotes",
			Source: `
> <!-- a -->
> 
> <!-- b -->
`,
			Node: func(n ast.Node) ast.Node {
				return n.FirstChild().FirstChild().NextSibling()
			},
			Expected: true,
		},
		{
			Name: "nesting HTML blocks in blockquotes",
			Source: `
> <!-- a -->
> <!-- b -->
`,
			Node: func(n ast.Node) ast.Node {
				return n.FirstChild().FirstChild().NextSibling()
			},
			Expected: false,
		},
		{
			Name: "nesting loose lists in blockquotes",
			Source: `
> - a
> 
> - b
`,
			Node: func(n ast.Node) ast.Node {
				return n.FirstChild().FirstChild().FirstChild().NextSibling()
			},
			Expected: true,
		},
		{
			Name: "nesting tight lists in blockquotes",
			Source: `
> - a
> - b
`,
			Node: func(n ast.Node) ast.Node {
				return n.FirstChild().FirstChild().FirstChild().NextSibling()
			},
			Expected: false,
		},
		{
			Name: "nesting paragraphs in lists",
			Source: `
- a

  b
`,
			Node: func(n ast.Node) ast.Node {
				return n.FirstChild().FirstChild().FirstChild().NextSibling()
			},
			Expected: true,
		},
	}
	md := New()
	for _, cs := range cases {
		t.Run(cs.Name, func(t *testing.T) {
			n := md.Parser().Parse(text.NewReader([]byte(cs.Source)))
			if cs.Node(n).HasBlankPreviousLines() != cs.Expected {
				t.Errorf("expected %v, got %v", cs.Expected, !cs.Expected)
			}
		})
	}
}
