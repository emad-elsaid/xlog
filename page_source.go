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

func NewPage(name string) (p Page) {
	app := GetApp()
	return app.NewPage(name)
}

func RegisterPageSource(p PageSource) {
	app := GetApp()
	app.RegisterPageSource(p)
}
