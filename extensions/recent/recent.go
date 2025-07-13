package recent

import (
	"embed"
	"html/template"
	"slices"
	"strings"

	_ "embed"

	"github.com/emad-elsaid/xlog"
	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	xlog.RegisterExtension(Recent{})
}

type Recent struct{}

func (Recent) Name() string { return "recent" }
func (Recent) Init(app *xlog.App) {
	app.Get(`/+/recent`, recentHandler)
	app.RegisterBuildPage("/+/recent", true)
	app.RegisterTemplate(templates, "templates")
	app.RegisterLink(func(Page) []Command { return []Command{links{}} })
}

func recentHandler(r Request) Output {
	app := xlog.GetApp()
	rp := app.Pages(r.Context())
	slices.SortFunc(rp, func(a, b Page) int {
		if modtime := b.ModTime().Compare(a.ModTime()); modtime != 0 {
			return modtime
		}

		return strings.Compare(a.Name(), b.Name())
	})

	return app.Render("recent", Locals{
		"page":  DynamicPage{NameVal: "Recent"},
		"pages": rp,
	})
}

type links struct{}

func (l links) Icon() string { return "fa-solid fa-clock-rotate-left" }
func (l links) Name() string { return "Recent" }
func (l links) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"href": "/+/recent",
	}
}
