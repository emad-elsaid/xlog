package xlog

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func renderMarkdown(content string) string {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.DefinitionList, extension.Footnote, highlighting.Highlighting),
		goldmark.WithRendererOptions(html.WithHardWraps()),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return err.Error()
	}

	post := postProcess(buf.String())

	return post
}

func postProcess(content string) string {
	r := strings.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err.Error()
	}

	linkPages(doc)
	out, _ := doc.Html()
	return out
}

func linkPages(doc *goquery.Document) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			basename := (file.Name()[:len(file.Name())-3])
			selector := fmt.Sprintf(":contains('%s')", basename)

			doc.Find(selector).Each(func(i int, s *goquery.Selection) {
				if goquery.NodeName(s) != "#text" {
					return
				}
				if s.ParentsFiltered("code,a,pre").Length() > 0 {
					return
				}

				h, _ := goquery.OuterHtml(s)
				fmt.Println(selector, h)

				text, _ := goquery.OuterHtml(s)
				s.ReplaceWithHtml(strings.ReplaceAll(text, basename, fmt.Sprintf(`<a href="%s">%s</a>`, basename, basename)))
			})
		}
	}
}
