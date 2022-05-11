package main

import (
	"io/ioutil"
	"sort"
	"strings"
)

func Search(keyword string) []string {
	pages := []string{}
	files, _ := ioutil.ReadDir(".")
	sort.Sort(fileInfoByNameLength(files))

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			f, err := ioutil.ReadFile(file.Name())
			if err != nil {
				continue
			}

			basename := file.Name()[:len(file.Name())-3]
			if strings.Contains(string(f), keyword) {
				pages = append(pages, basename)
			}
		}
	}

	return pages
}
