package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func postProcess(content string) (string, []string, error) {
	r := strings.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", nil, err
	}

	references := linkPages(doc)
	html, err := doc.Html()
	return html, references, err
}

func linkPages(doc *goquery.Document) []string {
	files, _ := ioutil.ReadDir(".")
	sort.Sort(fileInfoByNameLength(files))
	references := []string{}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			basename := file.Name()[:len(file.Name())-3]
			if linkPage(doc, basename) {
				references = append(references, basename)
			}
		}
	}

	return references
}

func linkPage(doc *goquery.Document, basename string) bool {
	selector := fmt.Sprintf(":contains('%s')", basename)
	found := false

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) != "#text" || s.ParentsFiltered("code,a,pre").Length() > 0 {
			return
		}

		found = true
		text, _ := goquery.OuterHtml(s)
		reg := regexp.MustCompile(`(?imU)(?:^|\s)(` + regexp.QuoteMeta(basename) + `)(?:\s|$)`)

		s.ReplaceWithHtml(reg.ReplaceAllString(text, ` <a href="$1">$1</a> `))
	})

	return found
}

type fileInfoByNameLength []os.FileInfo

func (a fileInfoByNameLength) Len() int           { return len(a) }
func (a fileInfoByNameLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fileInfoByNameLength) Less(i, j int) bool { return len(a[i].Name()) > len(a[j].Name()) }
