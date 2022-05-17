package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type postProcessor func(*goquery.Document)

var postProcessors = []postProcessor{}

func POSTPROCESSOR(f postProcessor) {
	postProcessors = append(postProcessors, f)
}

func postProcess(content string) (string, error) {
	r := strings.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}

	for _, v := range postProcessors {
		v(doc)
	}

	return doc.Html()
}
