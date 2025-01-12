package recent

import (
	"embed"
	"html/template"
	"slices"
	"strings"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	RegisterExtension(Recent{})
}

type Recent struct{}

func (Recent) Name() string { return "recent" }
func (Recent) Init() {
	Get(`/+/recent`, recentHandler)
	RegisterBuildPage("/+/recent", true)
	RegisterTemplate(templates, "templates")
	RegisterLink(func(Page) []Command { return []Command{links{}} })
}

func recentHandler(r Request) Output {
	rp := Pages(r.Context())
	slices.SortFunc(rp, func(a, b Page) int {
		if modtime := b.ModTime().Compare(a.ModTime()); modtime != 0 {
			return modtime
		}

		return strings.Compare(a.Name(), b.Name())
	})

	return Render("recent", Locals{
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
