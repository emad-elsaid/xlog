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
	RegisterExtension(Recent{})
}

type Recent struct{}

func (Recent) Name() string { return "recent" }
func (Recent) Init() {
	Get(`/+/recent`, recentHandler)
	RegisterBuildPage("/+/recent", true)
	RegisterTemplate(templates, "templates")
	RegisterLink(func(_ Page) []Link { return []Link{links{}} })
}

func recentHandler(r Request) Output {
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
func (a recentPages) Less(i, j int) bool { return a[i].ModTime().After(a[j].ModTime()) }

type links struct{}

func (l links) Icon() string { return "fa-solid fa-clock-rotate-left" }
func (l links) Name() string { return "Recent" }
func (l links) Link() string { return "/+/recent" }
