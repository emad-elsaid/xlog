package search

import (
	"context"
	"embed"
	"html/template"
	"regexp"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

const MIN_SEARCH_KEYWORD = 3

//go:embed templates
var templates embed.FS

func init() {
	Get(`/\+/search`, searchHandler)
	PrependWidget(SIDEBAR_WIDGET, sidebar)
	Template(templates, "templates")
}

func sidebar(_ *Page, _ Request) template.HTML {
	if READONLY {
		return ""
	}

	return Partial("search-widget", nil)
}

func searchHandler(w Response, r Request) Output {
	return Render("search-datalist", Locals{
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

	EachPage(ctx, func(p *Page) {
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
