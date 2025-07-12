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
	app := GetApp()
	app.RegisterExtension(Search{})
}

type Search struct{}

func (Search) Name() string { return "search" }
func (Search) Init() {
	app := GetApp()
	if app.GetConfig().Readonly {
		return
	}

	app.RequireHTMX()
	app.Get(`/+/search`, searchFormHandler)
	app.Get(`/+/search-result`, searchResultHandler)
	app.RegisterWidget("search", 0, searchWidget)
	app.RegisterTemplate(templates, "templates")
}

func searchWidget(Page) template.HTML {
	app := GetApp()
	return app.Partial("search", nil)
}

func searchFormHandler(r Request) Output {
	app := GetApp()
	return app.Render("search-form", Locals{
		"page":    DynamicPage{NameVal: "Create"},
		"results": search(r.Context(), r.FormValue("q")),
	})
}

func searchResultHandler(r Request) Output {
	app := GetApp()
	return app.Render("search-result", Locals{
		"results": search(r.Context(), r.FormValue("q")),
	})
}

type searchResult struct {
	Page Page
	Line string
}

func search(ctx context.Context, keyword string) []*searchResult {
	results := []*searchResult{}
	if len(keyword) < MIN_SEARCH_KEYWORD {
		return results
	}

	reg := regexp.MustCompile(`(?imU)^(.*` + regexp.QuoteMeta(keyword) + `.*)$`)

	app := GetApp()
	return MapPage(app, ctx, func(p Page) *searchResult {
		match := reg.FindString(p.Name())
		if len(match) > 0 {
			return &searchResult{
				Page: p,
				Line: "Matches the file name",
			}
		}

		match = reg.FindString(string(p.Content()))
		if len(match) > 0 {
			return &searchResult{
				Page: p,
				Line: match,
			}
		}

		return nil
	})
}
