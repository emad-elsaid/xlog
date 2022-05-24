package main

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/context"
)

func init() {
	POSTPROCESSOR(linkPages)
}

func linkPages(doc *goquery.Document) {
	pages := []*Page{}
	WalkPages(context.Background(), func(p *Page) {
		pages = append(pages, p)
	})

	sort.Sort(fileInfoByNameLength(pages))

	for _, p := range pages {
		linkPage(doc, p.Name)
	}
}

func linkPage(doc *goquery.Document, basename string) {
	selector := fmt.Sprintf(":contains('%s')", basename)

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) != "#text" || s.ParentsFiltered("code,a,pre").Length() > 0 {
			return
		}

		text, _ := goquery.OuterHtml(s)
		reg := regexp.MustCompile(`(?imU)(^|\W)(` + regexp.QuoteMeta(basename) + `)(\W|$)`)

		s.ReplaceWithHtml(reg.ReplaceAllString(text, `$1<a href="`+basename+`">$2</a>$3`))
	})
}

type fileInfoByNameLength []*Page

func (a fileInfoByNameLength) Len() int           { return len(a) }
func (a fileInfoByNameLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fileInfoByNameLength) Less(i, j int) bool { return len(a[i].Name) > len(a[j].Name) }
