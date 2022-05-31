package main

import (
	"context"
	"html/template"
	"sort"
)

func init() {
	WIDGET(SIDEBAR_WIDGET, recentNotes)
}

func recentNotes(p *Page, r Request) template.HTML {
	rp := recentPages{}
	WalkPages(context.Background(), func(i *Page) {
		rp = append(rp, i)
	})

	sort.Sort(rp)

	if len(rp) > 10 {
		rp = rp[:10]
	}

	pages := []string{}

	for _, v := range rp {
		pages = append(pages, v.Name)
	}

	return template.HTML(partial("extension/recent-notes", Locals{
		"pages": pages,
	}))
}

type recentPages []*Page

func (a recentPages) Len() int           { return len(a) }
func (a recentPages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a recentPages) Less(i, j int) bool { return a[i].ModTime().After(a[j].ModTime()) }
