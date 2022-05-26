package main

import (
	"context"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	POSTPROCESSOR(linkHashtags)
	SIDEBAR(hashtagsSidebar)

	GET("/+/tag/{tag}", tagHandler)
}

var hashtagReg = regexp.MustCompile(`(?imU)#([[:alpha:]]\w+)(\W|$)`)

func linkHashtags(doc *goquery.Document) {
	selector := fmt.Sprintf(":contains('#')")

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) != "#text" || s.ParentsFiltered("code,a,pre").Length() > 0 {
			return
		}

		text, _ := goquery.OuterHtml(s)
		s.ReplaceWithHtml(hashtagReg.ReplaceAllString(text, `<a href="/+/tag/$1" class="tag is-info">#$1</a>$2`))
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

func hashtagsSidebar(p *Page, r Request) template.HTML {
	tags := map[string][]string{}
	WalkPages(context.Background(), func(a *Page) {
		set := map[string]bool{}
		hashes := hashtagReg.FindAllStringSubmatch(a.Content(), -1)
		for _, v := range hashes {
			val := strings.ToLower(v[1])

			// don't use same tag twice for same page
			if _, ok := set[val]; ok {
				continue
			}

			set[val] = true
			if ps, ok := tags[val]; ok {
				tags[val] = append(ps, a.Name)
			} else {
				tags[val] = []string{a.Name}
			}
		}
	})

	return template.HTML(partial("extension/tags-sidebar", Locals{
		"tags": tags,
	}))
}
