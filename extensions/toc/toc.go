package toc

// TOC is the table of contents. It's the top-level object under which the
// rest of the table of contents resides.
type TOC struct {
	// Items holds the top-level headings under the table of contents.
	//
	// Items is empty if there are no headings in the document.
	Items Items
}

// Item is a single item in the table of contents.
type Item struct {
	// Title of this item in the table of contents.
	//
	// This may be blank for items that don't refer to a heading, and only
	// have sub-items.
	Title []byte

	// ID is the identifier for the heading that this item refers to. This
	// is the fragment portion of the link without the "#".
	//
	// This may be blank if the item doesn't have an id assigned to it, or
	// if it doesn't have a title.
	//
	// Enable AutoHeadingID in your parser if you expected these to be set
	// but they weren't.
	ID []byte

	// Items references children of this item.
	//
	// For a heading at level 3, Items, contains the headings at level 4
	// under that section.
	Items Items
}

// Items is a list of items in a table of contents.
type Items []*Item
