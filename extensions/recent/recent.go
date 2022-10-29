package recent

import (
	"context"
	"embed"
	"html/template"
	"sort"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	Widget(SIDEBAR_WIDGET, recent)
	Get(`/\+/recent`, recentHandler)
	BuildPage("/+/recent", true)
	Template(templates, "templates")
}

func recentHandler(_ Response, r Request) Output {
	rp := recentPages{}
	EachPage(context.Background(), func(i *Page) {
		rp = append(rp, i)
	})

	sort.Sort(rp)

	if len(rp) > 100 {
		rp = rp[:100]
	}

	return Render("recent", Locals{
		"title":   "Recent",
		"pages":   rp,
		"sidebar": RenderWidget(SIDEBAR_WIDGET, nil, r),
	})
}

func recent(p *Page, r Request) template.HTML {
	return Partial("recent-sidebar", nil)
}

type recentPages []*Page

func (a recentPages) Len() int           { return len(a) }
func (a recentPages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a recentPages) Less(i, j int) bool { return a[i].ModTime().After(a[j].ModTime()) }
