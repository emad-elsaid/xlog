package xlog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sort"

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
	files, _ := ioutil.ReadDir(".")
	sort.Sort(fileInfoByNameLength(files))

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			basename := file.Name()[:len(file.Name())-3]
			linkPage(doc, basename)
		}
	}
}

func linkPage(doc *goquery.Document, basename string) {
	selector := fmt.Sprintf(":contains('%s')", basename)
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) != "#text" || s.ParentsFiltered("code,a,pre").Length() > 0 {
			return
		}

		text, _ := goquery.OuterHtml(s)
		a := fmt.Sprintf(`<a href="%s">%s</a>`, basename, basename)

		s.ReplaceWithHtml(strings.ReplaceAll(text, basename, a))
	})
}
