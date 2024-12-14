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
	RegisterExtension(Search{})
}

type Search struct{}

func (Search) Name() string { return "search" }
func (Search) Init() {
	if Config.Readonly {
		return
	}

	Get(`/+/search`, searchHandler)
	RegisterWidget("search", 0, searchWidget)
	RegisterTemplate(templates, "templates")
}

func searchWidget(_ Page) template.HTML {
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

	EachPage(ctx, func(p Page) {
		match := reg.FindString(p.Name())
		if len(match) > 0 {
			results = append(results, searchResult{
				Page: p.Name(),
				Line: "Matches the file name",
			})
			return
		}

		match = reg.FindString(string(p.Content()))
		if len(match) > 0 {
			results = append(results, searchResult{
				Page: p.Name(),
				Line: match,
			})
		}
	})

	return results
}
