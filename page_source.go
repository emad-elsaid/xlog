package xlog

import (
	"context"
	"sync"
)

type PageSource interface {
	// Page takes a page name and return a Page struct
	Page(string) Page
	// Each iterates over all pages in the source
	Each(context.Context, func(Page))
}

var (
	sources      = []PageSource{newMarkdownFS(".")}
	sourcesMutex sync.RWMutex
)

// NewPage creates a Page instance for the given page name by querying all registered page sources.
// It returns the first existing page found across all sources. If no existing page is found,
// it returns a non-existing page object (never nil) to prevent panics in handlers.
// Page sources are checked in registration order (most recently registered first).
func NewPage(name string) (p Page) {
	sourcesMutex.RLock()
	defer sourcesMutex.RUnlock()
	
	for i := range sources {
		p = sources[i].Page(name)
		if p != nil && p.Exists() {
			return
		}
	}

	// If no existing page found, return a non-nil page object.
	// This prevents panics when handlers use NewPage with user input.
	// Try to get a page from the first source if available.
	if len(sources) > 0 {
		p = sources[0].Page(name)
	}

	// Final fallback: create a minimal page object if all sources returned nil.
	// This should rarely happen as markdownFS always returns a page object.
	if p == nil {
		p = &page{name: name}
	}

	return
}

// RegisterPageSource registers a new page source at the beginning of the sources list.
// New sources take precedence over previously registered sources when resolving pages.
// This allows extensions to override or supplement the default markdown file system.
func RegisterPageSource(p PageSource) {
	sourcesMutex.Lock()
	defer sourcesMutex.Unlock()
	sources = append([]PageSource{p}, sources...)
}
