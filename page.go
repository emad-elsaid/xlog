package xlog

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

// Counter for file reads during development
var fileReadCount atomic.Int64

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
}

type page struct {
	name string

	l          sync.Mutex
	lastUpdate time.Time
	ast        ast.Node
	content    *Markdown
	banner     *string    // cached banner result
	modTime    *time.Time // cached modtime
}

func (p *page) Name() string {
	return p.name
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
	if err := MarkdownConverter().Renderer().Render(&buf, src, ast); err != nil {
		return template.HTML(err.Error())
	}

	return template.HTML(buf.String())
}

func (p *page) Content() Markdown {
	fileReadCount.Add(1) // Track file reads
	dat, err := os.ReadFile(p.FileName())
	if err != nil {
		return ""
	}
	return Markdown(dat)
}

// modTimeLocked returns modtime assuming lock is already held
func (p *page) modTimeLocked() time.Time {
	if p.modTime != nil {
		return *p.modTime
	}

	s, err := os.Stat(p.FileName())
	if err != nil {
		zero := time.Time{}
		p.modTime = &zero
		return zero
	}

	mt := s.ModTime()
	p.modTime = &mt
	return mt
}

func (p *page) preProcessedContent() Markdown {
	p.l.Lock()
	defer p.l.Unlock()

	modtime := p.modTimeLocked()

	if p.content == nil || !modtime.Equal(p.lastUpdate) {
		c := p.Content()
		c = PreProcess(c)
		p.content = &c
		p.lastUpdate = modtime
	}

	return Markdown(*p.content)
}

func (p *page) Delete() bool {
	defer Trigger(PageDeleted, p)

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
	defer Trigger(PageChanged, p)

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
	p.l.Lock()
	defer p.l.Unlock()
	return p.modTimeLocked()
}

func (p *page) AST() (source []byte, tree ast.Node) {
	lastModified := p.lastUpdate
	content := p.preProcessedContent()

	if p.ast == nil || p.lastUpdate != lastModified {
		p.ast = MarkdownConverter().Parser().Parse(text.NewReader([]byte(content)))
	}

	return []byte(content), p.ast
}

func (p *page) clearCache() {
	p.content = nil
	p.ast = nil
	p.banner = nil
	p.modTime = nil
	p.lastUpdate = time.Time{}
}

// loadContent pre-loads content into the cache without triggering file read counter
// Used during cache warming to batch file reads efficiently
func (p *page) loadContent(content Markdown, modtime time.Time) {
	p.l.Lock()
	defer p.l.Unlock()

	c := PreProcess(content)
	p.content = &c
	p.lastUpdate = modtime

	// Cache modtime too to avoid another stat call
	p.modTime = &modtime
}

// DynamicPage implement Page interface and allow extensions to define a page to
// be passed to templates without having underlying file on desk
type DynamicPage struct {
	NameVal  string
	RenderFn func() template.HTML
}

func (DynamicPage) FileName() string        { return "" }
func (DynamicPage) Exists() bool            { return false }
func (DynamicPage) Content() Markdown       { return "" }
func (DynamicPage) Delete() bool            { return false }
func (DynamicPage) Write(Markdown) bool     { return false }
func (DynamicPage) ModTime() time.Time      { return time.Time{} }
func (DynamicPage) AST() ([]byte, ast.Node) { return nil, nil }
func (d DynamicPage) Name() string          { return d.NameVal }
func (d DynamicPage) Render() template.HTML {
	if d.RenderFn != nil {
		return d.RenderFn()
	}

	return ""
}
