package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

type fileInfoByNameLength []os.FileInfo

func (a fileInfoByNameLength) Len() int           { return len(a) }
func (a fileInfoByNameLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fileInfoByNameLength) Less(i, j int) bool { return len(a[i].Name()) > len(a[j].Name()) }
