package recent

import (
	"context"
	"embed"
	"sort"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	Get(`/\+/recent`, recentHandler)
	RegisterBuildPage("/+/recent", true)
	RegisterTemplate(templates, "templates")
	RegisterLink(func(_ Page) []Link { return []Link{links(0)} })
}

func recentHandler(_ Response, r Request) Output {
	rp := recentPages{}
	EachPage(context.Background(), func(i Page) {
		rp = append(rp, i)
	})

	sort.Sort(rp)

	if len(rp) > 100 {
		rp = rp[:100]
	}

	return Render("recent", Locals{
		"title": "Recent",
		"pages": rp,
	})
}

type recentPages []Page

func (a recentPages) Len() int           { return len(a) }
func (a recentPages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a recentPages) Less(i, j int) bool { return a[i].ModTime().After(a[j].ModTime()) }

type links int

func (l links) Icon() string { return "fa-solid fa-clock-rotate-left" }
func (l links) Name() string { return "Recent" }
func (l links) Link() string { return "/+/recent" }
