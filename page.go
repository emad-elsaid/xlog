package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Page struct {
	Name string
}

func NewPage(name string) Page {
	if name == "" {
		name = "index"
	}

	return Page{
		Name: name,
	}
}

func (p *Page) FileName() string {
	return p.Name + ".md"
}

func (p *Page) Exists() bool {
	_, err := os.Stat(p.FileName())
	return err == nil
}

func (p *Page) Render() (html string, refs []string) {
	content := p.Content()
	content = preProcess(content)
	html = renderMarkdown(content)
	html, refs, _ = postProcess(html)
	return
}

func (p *Page) Content() string {
	dat, err := ioutil.ReadFile(p.FileName())
	if err != nil {
		fmt.Printf("Can't open `%s`, err: %s\n", p.Name, err)
		return ""
	}
	return string(dat)
}

func (p *Page) Delete() bool {
	if p.Exists() {
		err := os.Remove(p.FileName())
		if err != nil {
			fmt.Printf("Can't delete `%s`, err: %s\n", p.Name, err)
			return false
		}
	}
	return true
}

func (p *Page) Write(content string) bool {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	err := ioutil.WriteFile(p.FileName(), []byte(content), 0644)
	if err != nil {
		fmt.Printf("Can't write `%s`, err: %s\n", p.Name, err)
		return false
	}
	return true
}
