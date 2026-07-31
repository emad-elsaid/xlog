package xlog

import (
	"bytes"
	"html"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

// pageLinkPattern matches page link syntax [[page_name]] in markdown content.
const pageLinkPattern = `\[\[([^\]]+)\]\]`

// Markdown is used instead of string to make sure it's clear the string is markdown string.
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
}

func (p *page) Name() string {
	return p.name
}

// ValidPageName reports whether a page name is safe to resolve to a file
// inside the source directory. It rejects absolute paths, Windows drive
// letters, backslashes, and "." or ".." path segments, all of which could
// otherwise escape the source directory when the process has chdir'd into it.
// Page sources should reject names that fail this check before performing any
// filesystem operation.
func ValidPageName(name string) bool {
	if name == "" {
		return true
	}

	// Reject absolute paths on any OS.
	if filepath.IsAbs(filepath.FromSlash(name)) {
		return false
	}

	// Reject Windows drive-relative paths like "C:foo".
	if len(name) >= 2 && name[1] == ':' {
		return false
	}

	// Reject backslashes; page names use forward slashes as separators only.
	if strings.ContainsRune(name, '\\') {
		return false
	}

	// Reject "." and ".." path segments that could escape the source dir.
	for _, segment := range strings.Split(name, "/") {
		switch segment {
		case ".", "..":
			return false
		}
	}

	return true
}

func (p *page) FileName() string {
	if !ValidPageName(p.name) {
		return ""
	}

	return filepath.FromSlash(p.name) + ".md"
}

func (p *page) Exists() bool {
	_, err := os.Stat(p.FileName())
	return err == nil
}

func (p *page) Render() template.HTML {
	src, astNode := p.AST()

	var buf bytes.Buffer
	if err := MarkdownConverter().Renderer().Render(&buf, src, astNode); err != nil {
		slog.Error("Failed to render page", "page", p.Name(), "error", err)
		// #nosec G203 -- Error message is html.EscapeString-sanitized before conversion
		return template.HTML(html.EscapeString(err.Error()))
	}

	// #nosec G203 - buf.String() contains markdown-rendered HTML which is already escaped by the renderer
	return template.HTML(buf.String())
}

func (p *page) Content() Markdown {
	dat, err := os.ReadFile(p.FileName())
	if err != nil {
		return ""
	}
	return Markdown(dat)
}

func (p *page) Delete() bool {
	p.clearCache()

	if p.Exists() {
		err := os.Remove(p.FileName())
		if err != nil {
			slog.Error("Can't delete page", "page", p.Name(), "error", err)
			return false
		}
	}

	Trigger(PageDeleted, p)
	return true
}

func (p *page) Write(content Markdown) bool {
	name := p.FileName()
	if err := os.MkdirAll(filepath.Dir(name), 0700); err != nil {
		slog.Error("Can't create page directory", "page", p.Name(), "error", err)
		return false
	}

	content = Markdown(strings.ReplaceAll(string(content), "\r\n", "\n"))
	if err := os.WriteFile(name, []byte(content), 0600); err != nil {
		slog.Error("Can't write page", "page", p.Name(), "error", err)
		return false
	}

	p.clearCache()
	Trigger(PageChanged, p)
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
	p.l.Lock()
	defer p.l.Unlock()

	// Get current file modification time
	modtime := p.ModTime()

	// Regenerate cached content if:
	// 1. Cache is empty (content == nil)
	// 2. File has been modified since last cache (modtime changed)
	if p.content == nil || !modtime.Equal(p.lastUpdate) {
		c := p.Content()
		c = PreProcess(c)
		p.content = &c
		p.lastUpdate = modtime

		// Content changed, invalidate AST cache so it gets regenerated
		p.ast = nil
	}

	content := Markdown(*p.content)

	// Re-parse AST if cache is empty
	if p.ast == nil {
		p.ast = MarkdownConverter().Parser().Parse(text.NewReader([]byte(content)))
	}

	return []byte(content), p.ast
}

func (p *page) clearCache() {
	p.l.Lock()
	defer p.l.Unlock()

	p.content = nil
	p.ast = nil
	p.lastUpdate = time.Time{}
}

// DynamicPage implement Page interface and allow extensions to define a page to
// be passed to templates without having underlying file on desk.
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
