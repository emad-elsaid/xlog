package recent

import (
	"embed"
	"html/template"
	"slices"
	"strings"

	_ "embed"

	"github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	xlog.RegisterExtension(Recent{})
}

type Recent struct{}

func (Recent) Name() string { return "recent" }
func (Recent) Init() {
	xlog.Get(`/+/recent`, recentHandler)
	xlog.RegisterBuildPage("/+/recent", true)
	xlog.RegisterTemplate(templates, "templates")
	xlog.RegisterLink(func(xlog.Page) []xlog.Command { return []xlog.Command{links{}} })
}

func recentHandler(r xlog.Request) xlog.Output {
	rp := xlog.Pages(r.Context())
	slices.SortFunc(rp, func(a, b xlog.Page) int {
		if modtime := b.ModTime().Compare(a.ModTime()); modtime != 0 {
			return modtime
		}

		return strings.Compare(a.Name(), b.Name())
	})

	return xlog.Render("recent", xlog.Locals{
		"page":  xlog.DynamicPage{NameVal: "Recent"},
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
