package main

import (
	"context"
	"html/template"
	"sort"
	"time"
)

func init() {
	WIDGET(SIDEBAR_WIDGET, recent)
	GET(`/\+/recent`, recentHandler)
	EXTENSION_PAGE("/+/recent")
	HELPER("ago", func(t time.Time) string {
		return ago(time.Now().Sub(t))
	})
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

	return Render("extension/recent", Locals{
		"title":   "Recent",
		"pages":   rp,
		"sidebar": renderWidget(SIDEBAR_WIDGET, nil, r),
	})
}

func recent(p *Page, r Request) template.HTML {
	return template.HTML(partial("extension/recent-sidebar", nil))
}

type recentPages []*Page

func (a recentPages) Len() int           { return len(a) }
func (a recentPages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a recentPages) Less(i, j int) bool { return a[i].ModTime().After(a[j].ModTime()) }
