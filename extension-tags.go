package main

import (
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
		"results": search(r.Context(), tag),
	})
}
