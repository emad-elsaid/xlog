package main

import (
	"context"
	"html/template"
	"regexp"
)

func init() {
	WIDGET(SIDEBAR_WIDGET, backlinksSidebar)
}

func backlinksSidebar(p *Page, r Request) template.HTML {
	pages := []string{}
	reg := regexp.MustCompile(`(?imU)(^|\W)(` + regexp.QuoteMeta(p.Name) + `)(\W|$)`)

	WalkPages(context.Background(), func(a *Page) {
		// a page shouldn't mention itself
		if a.Name == p.Name {
			return
		}

		if len(reg.FindString(a.Content())) > 0 {
			pages = append(pages, a.Name)
		}
	})

	return template.HTML(partial("extension/backlinks", Locals{"pages": pages}))
}
