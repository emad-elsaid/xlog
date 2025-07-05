package recent

import (
	"context"
	"embed"
	"html/template"
	"slices"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	app := GetApp()
	app.RegisterExtension(Recent{})
}

type Recent struct{}

func (Recent) Name() string { return "recent" }
func (Recent) Init() {
	app := GetApp()
	app.RegisterWidget(WidgetAfterView, 1, recentWidget)
	app.RegisterTemplate(templates, "templates")
}

func recentWidget(p Page) template.HTML {
	if p == nil {
		return ""
	}

	app := GetApp()
	recentPages := getRecentPages(context.Background())
	if len(recentPages) == 0 {
		return ""
	}

	return app.Partial("recent", Locals{"pages": recentPages})
}

func getRecentPages(ctx context.Context) []Page {
	app := GetApp()
	var pages []Page
	app.EachPage(ctx, func(p Page) {
		if p.Exists() && !p.ModTime().IsZero() {
			pages = append(pages, p)
		}
	})

	slices.SortFunc(pages, func(a, b Page) int {
		return b.ModTime().Compare(a.ModTime())
	})

	if len(pages) > 10 {
		pages = pages[:10]
	}

	return pages
}
