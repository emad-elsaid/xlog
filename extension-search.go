package main

import (
	"context"
	"html/template"
	"regexp"
)

const MIN_SEARCH_KEYWORD = 3

func init() {
	NAVBAR_START(searchNavbarStartWidget)
	GET("/+/search", searchHandler)
}

func searchNavbarStartWidget() template.HTML {
	return template.HTML(partial("extension/search-navbar", nil))
}

func searchHandler(w Response, r Request) Output {
	return Render("extension/search-datalist", Locals{
		"results": search(r.Context(), r.FormValue("q")),
	})
}

type searchResult struct {
	Page string
	Line string
}

func search(ctx context.Context, keyword string) []searchResult {
	results := []searchResult{}
	if len(keyword) < MIN_SEARCH_KEYWORD {
		return results
	}

	reg := regexp.MustCompile(`(?imU)^(.*` + regexp.QuoteMeta(keyword) + `.*)$`)

	WalkPages(ctx, func(p *Page) {
		match := reg.FindString(p.Name)
		if len(match) > 0 {
			results = append(results, searchResult{
				Page: p.Name,
				Line: "Matches the file name",
			})
			return
		}

		match = reg.FindString(p.Content())
		if len(match) > 0 {
			results = append(results, searchResult{
				Page: p.Name,
				Line: match,
			})
		}
	})

	return results
}
