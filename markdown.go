package xlog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

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

	post, err := postProcess(buf.String())
	if err != nil {
		return err.Error()
	}

	return post
}

func postProcess(content string) (string, error) {
	r := strings.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}

	linkPages(doc)
	return doc.Html()
}

func linkPages(doc *goquery.Document) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			basename := file.Name()[:len(file.Name())-3]
			selector := fmt.Sprintf(":contains('%s')", basename)

			doc.Find(selector).Each(func(i int, s *goquery.Selection) {
				if goquery.NodeName(s) != "#text" {
					return
				}
				if s.ParentsFiltered("code,a,pre").Length() > 0 {
					return
				}

				text, _ := goquery.OuterHtml(s)
				s.ReplaceWithHtml(strings.ReplaceAll(text, basename, fmt.Sprintf(`<a href="%s">%s</a>`, basename, basename)))
			})
		}
	}
}
