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

// RegisterPageSource registers a page source
func (app *App) RegisterPageSource(p PageSource) {
	app.sources = append([]PageSource{p}, app.sources...)
}
