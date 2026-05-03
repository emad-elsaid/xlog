package search

import (
	"context"
	"embed"
	"html/template"
	"regexp"

	_ "embed"

	"github.com/emad-elsaid/xlog"
)

const (
	MIN_SEARCH_KEYWORD = 3
	ExtensionName      = "search"
)

//go:embed templates
var templates embed.FS

func init() {
	xlog.RegisterExtension(Search{})
}

type Search struct{}

func (Search) Name() string { return ExtensionName }
func (Search) Init() {
	if xlog.Config.Readonly {
		return
	}

	xlog.RequireHTMX()
	xlog.Get(`/+/search`, searchFormHandler)
	xlog.Get(`/+/search-result`, searchResultHandler)
	xlog.RegisterWidget("search", 0, searchWidget)
	xlog.RegisterTemplate(templates, "templates")
}

func searchWidget(xlog.Page) template.HTML {
	return xlog.Partial("search", nil)
}

func searchFormHandler(r xlog.Request) xlog.Output {
	return xlog.Render("search-form", xlog.Locals{
		"page":    xlog.DynamicPage{NameVal: "Create"},
		"results": search(r.Context(), r.FormValue("q")),
	})
}

func searchResultHandler(r xlog.Request) xlog.Output {
	return xlog.Render("search-result", xlog.Locals{
		"results": search(r.Context(), r.FormValue("q")),
	})
}

type searchResult struct {
	Page xlog.Page
	Line string
}

func search(ctx context.Context, keyword string) []*searchResult {
	results := []*searchResult{}
	if len(keyword) < MIN_SEARCH_KEYWORD {
		return results
	}

	reg := regexp.MustCompile(`(?imU)^(.*` + regexp.QuoteMeta(keyword) + `.*)$`)

	return xlog.MapPage(ctx, func(p xlog.Page) *searchResult {
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
