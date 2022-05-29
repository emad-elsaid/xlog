package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type (
	preProcessor  func(string) string
	postProcessor func(*goquery.Document)
)

var (
	preProcessors  = []preProcessor{}
	postProcessors = []postProcessor{}
)

func PREPROCESSOR(f preProcessor)   { preProcessors = append(preProcessors, f) }
func POSTPROCESSOR(f postProcessor) { postProcessors = append(postProcessors, f) }

func preProcess(content string) string {
	for _, v := range preProcessors {
		content = v(content)
	}

	return content
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
