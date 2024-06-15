package recent

import (
	"embed"
	"sort"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	Get(`/+/recent`, recentHandler)
	RegisterBuildPage("/+/recent", true)
	RegisterTemplate(templates, "templates")
	RegisterLink(func(_ Page) []Link { return []Link{links(0)} })
}

func recentHandler(_ Response, r Request) Output {
	var rp recentPages = Pages(r.Context())
	sort.Sort(rp)

	return Render("recent", Locals{
		"title": "Recent",
		"pages": rp,
	})
}

type recentPages []Page

func (a recentPages) Len() int           { return len(a) }
func (a recentPages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a recentPages) Less(i, j int) bool { return a[i].ModTime(false).After(a[j].ModTime(false)) }

type links int

func (l links) Icon() string { return "fa-solid fa-clock-rotate-left" }
func (l links) Name() string { return "Recent" }
func (l links) Link() string { return "/+/recent" }
