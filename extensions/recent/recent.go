package recent

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"sort"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

//go:embed views
var views embed.FS

func init() {
	WIDGET(SIDEBAR_WIDGET, recent)
	GET(`/\+/recent`, recentHandler)
	EXTENSION_PAGE("/+/recent")

	fs, _ := fs.Sub(views, "views")
	VIEW(fs)
}

func recentHandler(_ Response, r Request) Output {
	rp := recentPages{}
	WalkPages(context.Background(), func(i *Page) {
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
	return template.HTML(Partial("recent-sidebar", nil))
}

type recentPages []*Page

func (a recentPages) Len() int           { return len(a) }
func (a recentPages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a recentPages) Less(i, j int) bool { return a[i].ModTime().After(a[j].ModTime()) }
