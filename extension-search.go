package main

import (
	"context"
	"html/template"
	"io/ioutil"
	"path"
	"regexp"
)

func init() {
	NAVBAR_START(searchNavbarStartWidget)
	GET("/+/search", searchHandler)
}

func searchNavbarStartWidget() template.HTML {
	return template.HTML(partial("extension/search-navbar", nil))
}

func searchHandler(w Response, r Request) Output {
	return Render("extension/search-datalist", Locals{
		"results": search(r.Context(), r.FormValue("q")),
	})
}

type searchResult struct {
	Page string
	Line string
}

func search(ctx context.Context, keyword string) []searchResult {
	results := []searchResult{}
	if len(keyword) < MIN_SEARCH_KEYWORD {
		return results
	}

	files, _ := ioutil.ReadDir(".")
	reg := regexp.MustCompile(`(?imU)^(.*` + regexp.QuoteMeta(keyword) + `.*)$`)

	for _, file := range files {
		select {
		case <-ctx.Done():
			break
		default:

			name := file.Name()
			ext := path.Ext(name)
			basename := name[:len(name)-len(ext)]

			if !file.IsDir() && ext == ".md" {

				match := reg.FindString(file.Name())
				if len(match) > 0 {
					results = append(results, searchResult{
						Page: basename,
						Line: "Matches the file name",
					})
					continue
				}

				f, err := ioutil.ReadFile(file.Name())
				if err != nil {
					continue
				}

				match = reg.FindString(string(f))
				if len(match) > 0 {
					results = append(results, searchResult{
						Page: basename,
						Line: match,
					})
				}
			}

		}
	}

	return results
}
