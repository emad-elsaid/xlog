package xlog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	emojiAst "github.com/yuin/goldmark-emoji/ast"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

// a Type that represent a page.
type Page struct {
	Name string // page name without '.md' extension
	ast  ast.Node
}

// The instance of markdown renderer. this is what takes the page content and
// converts it to HTML. it defines what features to use from goldmark and what
// options to turn on
var MarkDownRenderer = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.DefinitionList,
		extension.Footnote,
		highlighting.Highlighting,
		emoji.Emoji,
	),

	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithUnsafe(),
	),
)

// Create an instance of Page with name. if no name is passed it's assumed "index"
func NewPage(name string) Page {
	if name == "" {
		name = INDEX
	}

	return Page{
		Name: name,
	}
}

// returns the filename, makes sure it converts slashes to backslashes when
// needed. this is safe to use when trying to access the file that represent the
// page
func (p *Page) FileName() string {
	return filepath.FromSlash(p.Name) + ".md"
}

// checks if the page underlying file exists on disk or not.
func (p *Page) Exists() bool {
	_, err := os.Stat(p.FileName())
	return err == nil
}

// Renders the page content to HTML. it makes sure all preprocessors are called
func (p *Page) Render() string {
	content := p.Content()
	content = preProcess(content)

	var buf bytes.Buffer
	if err := MarkDownRenderer.Convert([]byte(content), &buf); err != nil {
		return err.Error()
	}

	return buf.String()
}

// Reads the underlying file and returns the content
func (p *Page) Content() string {
	dat, err := os.ReadFile(p.FileName())
	if err != nil {
		return ""
	}
	return string(dat)
}

// Deletes the file and makes sure it triggers the AfterDelete event
func (p *Page) Delete() bool {
	defer Trigger(AfterDelete, p)

	if p.Exists() {
		err := os.Remove(p.FileName())
		if err != nil {
			fmt.Printf("Can't delete `%s`, err: %s\n", p.Name, err)
			return false
		}
	}
	return true
}

// Overwrite page content with new content. making sure to trigger before and
// after write events.
func (p *Page) Write(content string) bool {
	Trigger(BeforeWrite, p)
	defer Trigger(AfterWrite, p)

	name := p.FileName()
	os.MkdirAll(filepath.Dir(name), 0700)

	content = strings.ReplaceAll(content, "\r\n", "\n")
	if err := os.WriteFile(name, []byte(content), 0644); err != nil {
		fmt.Printf("Can't write `%s`, err: %s\n", p.Name, err)
		return false
	}
	return true
}

// Return the last modification time of the underlying file
func (p *Page) ModTime() time.Time {
	s, err := os.Stat(p.FileName())
	if err != nil {
		return time.Time{}
	}

	return s.ModTime()
}

// Parses the page content and returns the Abstract Syntax Tree (AST).
// extensions can use it to walk the tree and modify it or collect statistics or
// parts of the page. for example the following "Emoji" function uses it to
// extract the first emoji.
func (p *Page) AST() ast.Node {
	if p.ast == nil {
		p.ast = MarkDownRenderer.Parser().Parse(text.NewReader([]byte(p.Content())))
	}

	return p.ast
}

// Returns the first emoji of the page.
func (p *Page) Emoji() string {
	if e, ok := FindInAST[*emojiAst.Emoji](p.AST(), emojiAst.KindEmoji); ok {
		return string(e.Value.Unicode)
	}

	return ""
}

// This is a function that takes an AST node and walks the tree depth first
// recursively calling itself in search for a node of a specific kind
// can be used to find first image, link, paragraph...etc
func FindInAST[t ast.Node](n ast.Node, kind ast.NodeKind) (found t, ok bool) {
	if n.Kind() == kind {
		if found, ok := n.(t); ok {
			return found, true
		}
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if a, ok := FindInAST[t](c, kind); ok {
			return a, true
		}
	}

	return found, false
}

// Extract all nodes of a specific type from the AST
func FindAllInAST[t ast.Node](n ast.Node, kind ast.NodeKind) (a []t) {
	if n.Kind() == kind {
		typed, _ := n.(t)
		a = []t{typed}
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		a = append(a, FindAllInAST[t](c, kind)...)
	}

	return a
}
