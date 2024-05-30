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

func NewPage(name string) (p Page) {
	for i := range sources {
		p = sources[i].Page(name)
		if p != nil && p.Exists() {
			return
		}
	}

	return
}

func RegisterPageSource(p PageSource) {
	sources = append([]PageSource{p}, sources...)
}
