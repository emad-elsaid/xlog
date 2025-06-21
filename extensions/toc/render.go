package toc

import "github.com/emad-elsaid/xlog/markdown/ast"

const _defaultMarker = '*'

// RenderList renders a table of contents as a nested list with a sane,
// default configuration for the ListRenderer.
//
// If the TOC is nil or empty, nil is returned.
// Do not call Goldmark's renderer if the returned node is nil.
func RenderList(toc *TOC) ast.Node {
	return new(ListRenderer).Render(toc)
}

// RenderOrderedList renders a table of contents as a nested, ordered
// list with a sane, default configuration for the ListRenderer.
//
// If the TOC is nil or empty, nil is returned.
// Do not call Goldmark's renderer if the returned node is nil.
func RenderOrderedList(toc *TOC) ast.Node {
	renderer := ListRenderer{Marker: '.'}
	return renderer.Render(toc)
}

// ListRenderer builds a nested list from a table of contents.
//
// For example,
//
//	# Foo
//	## Bar
//	## Baz
//	# Qux
//
//	// becomes
//
//	- Foo
//	  - Bar
//	  - Baz
//	- Qux
type ListRenderer struct {
	// Marker for elements of the list, e.g. '-', '*', etc.
	//
	// Defaults to '*'.
	Marker byte
}

// Render renders the table of contents into Markdown.
//
// If the TOC is nil or empty, nil is returned.
// Do not call Goldmark's renderer if the returned node is nil.
func (r *ListRenderer) Render(toc *TOC) ast.Node {
	if toc == nil {
		return nil
	}
	return r.renderItems(toc.Items)
}

func (r *ListRenderer) renderItems(items Items) ast.Node {
	if len(items) == 0 {
		return nil
	}

	mkr := r.Marker
	if mkr == 0 {
		mkr = _defaultMarker
	}

	list := ast.NewList(mkr)
	if list.IsOrdered() {
		list.Start = 1
	}
	for _, item := range items {
		list.AppendChild(list, r.renderItem(item))
	}
	return list
}

func (r *ListRenderer) renderItem(n *Item) ast.Node {
	item := ast.NewListItem(0)

	if t := n.Title; len(t) > 0 {
		title := ast.NewString(t)
		title.SetRaw(true)
		if len(n.ID) > 0 {
			link := ast.NewLink()
			link.Destination = append([]byte("#"), n.ID...)
			link.AppendChild(link, title)
			item.AppendChild(item, link)
		} else {
			item.AppendChild(item, title)
		}
	}

	if items := r.renderItems(n.Items); items != nil {
		item.AppendChild(item, items)
	}

	return item
}
