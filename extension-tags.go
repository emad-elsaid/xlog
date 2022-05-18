package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	POSTPROCESSOR(linkHashtags)

	GET("/+/tag/{tag}", tagHandler)
}

func linkHashtags(doc *goquery.Document) {
	selector := fmt.Sprintf(":contains('#')")
	reg := regexp.MustCompile(`(?imU)#(\w+)(\W|$)`)

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) != "#text" || s.ParentsFiltered("code,a,pre").Length() > 0 {
			return
		}

		text, _ := goquery.OuterHtml(s)
		s.ReplaceWithHtml(reg.ReplaceAllString(text, `<a href="/+/tag/$1" class="tag is-info">#$1</a>$2`))
	})
}

func tagHandler(w Response, r Request) Output {
	vars := VARS(r)
	tag := "#" + vars["tag"]

	return Render("extension/tag", Locals{
		"title":   tag,
		"results": tagPages(r.Context(), tag),
	})
}

type tagResult struct {
	Page string
	Line string
}

func tagPages(ctx context.Context, keyword string) []tagResult {
	results := []tagResult{}
	reg := regexp.MustCompile(`(?imU)^(.*` + regexp.QuoteMeta(keyword) + `.*)$`)

	WalkPages(ctx, func(p *Page) {
		match := reg.FindString(p.Content())
		if len(match) > 0 {
			results = append(results, tagResult{
				Page: p.Name,
				Line: match,
			})
		}
	})

	return results
}
