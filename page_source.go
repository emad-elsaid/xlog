package xlog

import (
	"context"
)

type PageSource interface {
	// Page takes a page name and return a Page struct
	Page(string) Page
	// Each iterates over all pages in the source
	Each(context.Context, func(Page))
}

var sources = []PageSource{
	newMarkdownFS("."),
}

// NewPage creates a Page instance for the given page name by querying all registered page sources.
// It returns the first existing page found across all sources, or nil if no page exists.
// Page sources are checked in registration order (most recently registered first).
func NewPage(name string) (p Page) {
	for i := range sources {
		p = sources[i].Page(name)
		if p != nil && p.Exists() {
			return
		}
	}

	return
}

// RegisterPageSource registers a new page source at the beginning of the sources list.
// New sources take precedence over previously registered sources when resolving pages.
// This allows extensions to override or supplement the default markdown file system.
func RegisterPageSource(p PageSource) {
	sources = append([]PageSource{p}, sources...)
}
