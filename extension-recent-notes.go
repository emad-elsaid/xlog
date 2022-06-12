package main

import (
	"context"
	"html/template"
	"sort"
)

func init() {
	WIDGET(SIDEBAR_WIDGET, recentNotes)
	GET(`/\+/recent`, recentHandler)
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

	pages := []string{}

	for _, v := range rp {
		pages = append(pages, v.Name)
	}

	return Render("extension/recent", Locals{
		"title":   "Recent",
		"pages":   pages,
		"sidebar": renderWidget(SIDEBAR_WIDGET, nil, r),
	})
}

func recentNotes(p *Page, r Request) template.HTML {
	return template.HTML(partial("extension/recent-sidebar", nil))
}

type recentPages []*Page

func (a recentPages) Len() int           { return len(a) }
func (a recentPages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a recentPages) Less(i, j int) bool { return a[i].ModTime().After(a[j].ModTime()) }
