package toc

import (
	"bytes"
	"fmt"
	"io"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/util"
)

// InspectOption customizes the behavior of Inspect.
type InspectOption interface {
	apply(*inspectOptions)
}

type inspectOptions struct {
	minDepth int
	maxDepth int
	compact  bool
}

// MinDepth limits the depth of the table of contents.
// Headings with a level lower than the specified depth will be ignored.
//
// For example, given the following:
//
//	# Foo
//	## Bar
//	### Baz
//	# Quux
//	## Qux
//
// MinDepth(3) will result in the following:
//
//	TOC{Items: ...}
//	 |
//	 +--- &Item{Title: "Baz", ID: "baz"}
//
// Whereas, MinDepth(2) will result in the following:
//
//	TOC{Items: ...}
//	 |
//	 +--- &Item{Title: "Bar", ID: "bar", Items: ...}
//	 |     |
//	 |     +--- &Item{Title: "Baz", ID: "baz"}
//	 |
//	 +--- &Item{Title: "Qux", ID: "qux"}
//
// A value of 0 or less will result in no limit.
//
// The default is no limit.
func MinDepth(depth int) InspectOption {
	return minDepthOption(depth)
}

type minDepthOption int

func (d minDepthOption) apply(opts *inspectOptions) {
	opts.minDepth = int(d)
}

func (d minDepthOption) String() string {
	return fmt.Sprintf("MinDepth(%d)", int(d))
}

// MaxDepth limits the depth of the table of contents.
// Headings with a level greater than the specified depth will be ignored.
//
// For example, given the following:
//
//	# Foo
//	## Bar
//	### Baz
//	# Quux
//	## Qux
//
// MaxDepth(1) will result in the following:
//
//	TOC{Items: ...}
//	 |
//	 +--- &Item{Title: "Foo", ID: "foo"}
//	 |
//	 +--- &Item{Title: "Quux", ID: "quux", Items: ...}
//
// Whereas, MaxDepth(2) will result in the following:
//
//	TOC{Items: ...}
//	 |
//	 +--- &Item{Title: "Foo", ID: "foo", Items: ...}
//	 |     |
//	 |     +--- &Item{Title: "Bar", ID: "bar"}
//	 |
//	 +--- &Item{Title: "Quux", ID: "quux", Items: ...}
//	       |
//	       +--- &Item{Title: "Qux", ID: "qux"}
//
// A value of 0 or less will result in no limit.
//
// The default is no limit.
func MaxDepth(depth int) InspectOption {
	return maxDepthOption(depth)
}

type maxDepthOption int

func (d maxDepthOption) apply(opts *inspectOptions) {
	opts.maxDepth = int(d)
}

func (d maxDepthOption) String() string {
	return fmt.Sprintf("MaxDepth(%d)", int(d))
}

// Compact instructs Inspect to remove empty items from the table of contents.
// Children of removed items will be promoted to the parent item.
//
// For example, given the following:
//
//	# A
//	### B
//	#### C
//	# D
//	#### E
//
// Compact(false), which is the default, will result in the following:
//
//	TOC{Items: ...}
//	 |
//	 +--- &Item{Title: "A", ...}
//	 |     |
//	 |     +--- &Item{Title: "", ...}
//	 |           |
//	 |           +--- &Item{Title: "B", ...}
//	 |                 |
//	 |                 +--- &Item{Title: "C"}
//	 |
//	 +--- &Item{Title: "D", ...}
//	       |
//	       +--- &Item{Title: "", ...}
//	             |
//	             +--- &Item{Title: "", ...}
//	                   |
//	                   +--- &Item{Title: "E", ...}
//
// Whereas, Compact(true) will result in the following:
//
//	TOC{Items: ...}
//	 |
//	 +--- &Item{Title: "A", ...}
//	 |     |
//	 |     +--- &Item{Title: "B", ...}
//	 |           |
//	 |           +--- &Item{Title: "C"}
//	 |
//	 +--- &Item{Title: "D", ...}
//	       |
//	       +--- &Item{Title: "E", ...}
//
// Notice that the empty items have been removed
// and the generated TOC is more compact.
func Compact(compact bool) InspectOption {
	return compactOption(compact)
}

type compactOption bool

func (c compactOption) apply(opts *inspectOptions) {
	opts.compact = bool(c)
}

func (c compactOption) String() string {
	return fmt.Sprintf("Compact(%v)", bool(c))
}

// Inspect builds a table of contents by inspecting the provided document.
//
// The table of contents is represents as a tree where each item represents a
// heading or a heading level with zero or more children.
// The returned TOC will be empty if there are no headings in the document.
//
// For example,
//
//	# Section 1
//	## Subsection 1.1
//	## Subsection 1.2
//	# Section 2
//	## Subsection 2.1
//	# Section 3
//
// Will result in the following items.
//
//	TOC{Items: ...}
//	 |
//	 +--- &Item{Title: "Section 1", ID: "section-1", Items: ...}
//	 |     |
//	 |     +--- &Item{Title: "Subsection 1.1", ID: "subsection-1-1"}
//	 |     |
//	 |     +--- &Item{Title: "Subsection 1.2", ID: "subsection-1-2"}
//	 |
//	 +--- &Item{Title: "Section 2", ID: "section-2", Items: ...}
//	 |     |
//	 |     +--- &Item{Title: "Subsection 2.1", ID: "subsection-2-1"}
//	 |
//	 +--- &Item{Title: "Section 3", ID: "section-3"}
//
// You may analyze or manipulate the table of contents before rendering it.
func Inspect(n ast.Node, src []byte, options ...InspectOption) (*TOC, error) {
	var opts inspectOptions
	for _, opt := range options {
		opt.apply(&opts)
	}

	// Appends an empty subitem to the given node
	// and returns a reference to it.
	appendChild := func(n *Item) *Item {
		child := new(Item)
		n.Items = append(n.Items, child)
		return child
	}

	// Returns the last subitem of the given node,
	// creating it if necessary.
	lastChild := func(n *Item) *Item {
		if len(n.Items) > 0 {
			return n.Items[len(n.Items)-1]
		}
		return appendChild(n)
	}

	var root Item

	stack := []*Item{&root} // inv: len(stack) >= 1
	err := ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		heading, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}
		if opts.minDepth > 0 && heading.Level < opts.minDepth {
			return ast.WalkSkipChildren, nil
		}

		if opts.maxDepth > 0 && heading.Level > opts.maxDepth {
			return ast.WalkSkipChildren, nil
		}

		// The heading is deeper than the current depth.
		// Append empty items to match the heading's level.
		for len(stack) < heading.Level {
			parent := stack[len(stack)-1]
			stack = append(stack, lastChild(parent))
		}

		// The heading is shallower than the current depth.
		// Move back up the stack until we reach the heading's level.
		if len(stack) > heading.Level {
			stack = stack[:heading.Level]
		}

		parent := stack[len(stack)-1]
		target := lastChild(parent)
		if len(target.Title) > 0 || len(target.Items) > 0 {
			target = appendChild(parent)
		}

		target.Title = util.UnescapePunctuations(nodeText(src, heading))
		if id, ok := n.AttributeString("id"); ok {
			target.ID, _ = id.([]byte)
		}

		return ast.WalkSkipChildren, nil
	})

	if opts.compact {
		compactItems(&root.Items)
	}

	return &TOC{Items: root.Items}, err
}

// compactItems removes items with no titles
// from the given list of items.
//
// Children of removed items will be promoted to the parent item.
func compactItems(items *Items) {
	for i := 0; i < len(*items); i++ {
		item := (*items)[i]
		if len(item.Title) > 0 {
			compactItems(&item.Items)
			continue
		}

		children := item.Items
		newItems := make(Items, 0, len(*items)-1+len(children))
		newItems = append(newItems, (*items)[:i]...)
		newItems = append(newItems, children...)
		newItems = append(newItems, (*items)[i+1:]...)
		*items = newItems
		i-- // start with first child
	}
}

func nodeText(src []byte, n ast.Node) []byte {
	var buf bytes.Buffer
	writeNodeText(src, &buf, n)
	return buf.Bytes()
}

func writeNodeText(src []byte, dst io.Writer, n ast.Node) {
	switch n := n.(type) {
	case *ast.Text:
		_, _ = dst.Write(n.Segment.Value(src))
	case *ast.String:
		_, _ = dst.Write(n.Value)
	default:
		for c := n.FirstChild(); c != nil; c = c.NextSibling() {
			writeNodeText(src, dst, c)
		}
	}
}
