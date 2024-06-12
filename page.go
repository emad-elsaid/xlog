package xlog

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	emojiAst "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Markdown is used instead of string to make sure it's clear the string is markdown string
type Markdown string

// a Type that represent a page.
type Page interface {
	// Name returns page name without '.md' extension
	Name() string
	// returns the filename, makes sure it converts slashes to backslashes when
	// needed. this is safe to use when trying to access the file that represent the
	// page
	FileName() string
	Title() string
	// checks if the page underlying file exists on disk or not.
	Exists() bool
	// Renders the page content to HTML. it makes sure all preprocessors are called
	Render() template.HTML
	// Reads the underlying file and returns the content
	Content() Markdown
	// Deletes the file and makes sure it triggers the AfterDelete event
	Delete() bool
	// Overwrite page content with new content. making sure to trigger before and
	// after write events.
	Write(Markdown) bool
	// ModTime Return the last modification time of the underlying file
	ModTime() time.Time
	// Parses the page content and returns the Abstract Syntax Tree (AST).
	// extensions can use it to walk the tree and modify it or collect statistics or
	// parts of the page. for example the following "Emoji" function uses it to
	// extract the first emoji.
	AST() ([]byte, ast.Node)
	// Returns the first emoji of the page.
	Emoji() string
}

type page struct {
	name string

	l          sync.Mutex
	lastUpdate time.Time
	ast        ast.Node
	content    *Markdown
}

func (p *page) Name() string {
	return p.name
}

func (p *page) Title() string {
	log.Printf("getting name %s", p.name)
	_, ast := p.AST()
	log.Println("getting ast")
	mtitle, ok := ast.OwnerDocument().Meta()["title"].(string)
	log.Println("has meta?")
	if !ok {
		mtitle = p.Name()
	}
	return mtitle
}

func (p *page) FileName() string {
	return filepath.FromSlash(p.name) + ".md"
}

func (p *page) Exists() bool {
	_, err := os.Stat(p.FileName())
	return err == nil
}

func (p *page) Render() template.HTML {
	src, ast := p.AST()

	var buf bytes.Buffer
	if err := MarkDownRenderer.Renderer().Render(&buf, src, ast); err != nil {
		return template.HTML(err.Error())
	}

	return template.HTML(buf.String())
}

func (p *page) Content() Markdown {
	dat, err := os.ReadFile(p.FileName())
	if err != nil {
		return ""
	}
	return Markdown(dat)
}

func (p *page) preProcessedContent() Markdown {
	p.l.Lock()
	defer p.l.Unlock()

	modtime := p.ModTime()

	if p.content == nil || modtime.Equal(p.lastUpdate) {
		c := p.Content()
		c = PreProcess(c)
		p.content = &c
		p.lastUpdate = p.ModTime()
	}

	return Markdown(*p.content)
}

func (p *page) Delete() bool {
	defer Trigger(AfterDelete, p)

	p.clearCache()

	if p.Exists() {
		err := os.Remove(p.FileName())
		if err != nil {
			fmt.Printf("Can't delete `%s`, err: %s\n", p.Name(), err)
			return false
		}
	}
	return true
}

func (p *page) Write(content Markdown) bool {
	Trigger(BeforeWrite, p)
	defer Trigger(AfterWrite, p)

	p.clearCache()
	name := p.FileName()
	os.MkdirAll(filepath.Dir(name), 0700)

	content = Markdown(strings.ReplaceAll(string(content), "\r\n", "\n"))
	if err := os.WriteFile(name, []byte(content), 0644); err != nil {
		fmt.Printf("Can't write `%s`, err: %s\n", p.Name(), err)
		return false
	}
	return true
}

func (p *page) ModTime() time.Time {
	s, err := os.Stat(p.FileName())
	if err != nil {
		return time.Time{}
	}

	return s.ModTime()
}

func (p *page) AST() (source []byte, tree ast.Node) {
	lastModified := p.lastUpdate
	content := p.preProcessedContent()

	if p.ast == nil || p.lastUpdate != lastModified {
		p.ast = MarkDownRenderer.Parser().Parse(text.NewReader([]byte(content)))
	}

	return []byte(content), p.ast
}

func (p *page) Emoji() string {
	_, tree := p.AST()
	if e, ok := FindInAST[*emojiAst.Emoji](tree); ok {
		return string(e.Value.Unicode)
	}

	return ""
}

func (p *page) clearCache() {
	p.content = nil
	p.ast = nil
	p.lastUpdate = time.Time{}
}
