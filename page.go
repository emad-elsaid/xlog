package xlog

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Page struct {
	name string
}

func NewPage(name string) Page {
	if name == "" {
		name = "index"
	}

	return Page{
		name: name,
	}
}

func (p *Page) Name() string {
	return p.name
}

func (p *Page) FileName() string {
	return p.name + ".md"
}

func (p *Page) Exists() bool {
	_, err := os.Stat(p.FileName())
	return err == nil
}

func (p *Page) Render() string {
	html := renderMarkdown(p.Content())
	html = processShortCodes(html)
	return html
}

func (p *Page) Title() string {
	dat, err := ioutil.ReadFile(p.FileName())
	if err != nil {
		fmt.Printf("Can't open `%s`, err: %s\n", p.name, err)
		return ""
	}

	fileContent := string(dat)
	endOfLine := strings.Index(fileContent, "\n")
	title := fileContent[:endOfLine]
	return title
}

func (p *Page) Content() string {
	dat, err := ioutil.ReadFile(p.FileName())
	if err != nil {
		fmt.Printf("Can't open `%s`, err: %s\n", p.name, err)
		return ""
	}

	fileContent := string(dat)
	endOfLine := strings.Index(fileContent, "\n")

	fileContent = fileContent[endOfLine+1:]
	endOfLine = strings.Index(fileContent, "\n")

	content := fileContent[endOfLine+1:]
	return content
}

func (p *Page) Delete() bool {
	if p.Exists() {
		err := os.Remove(p.FileName())
		if err != nil {
			fmt.Printf("Can't delete `%s`, err: %s\n", p.name, err)
			return false
		}
	}
	return true
}

func (p *Page) Write(title, content string) bool {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	err := ioutil.WriteFile(p.FileName(), []byte(title+"\n=========\n"+content), 0644)
	if err != nil {
		fmt.Printf("Can't write `%s`, err: %s\n", p.name, err)
		return false
	}

	return true
}
