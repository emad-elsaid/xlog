package xlog

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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

func init() {
	PageEvents.Listen(AfterWrite, clearWalkPagesCache)
	PageEvents.Listen(AfterDelete, clearWalkPagesCache)
}

type (
	// a Type that represent a page.
	Page struct {
		Name string // page name without '.md' extension
		ast  ast.Node
	}

	// a type used to define events to be used when the page is manipulated for
	// example modified, renamed, deleted...etc.
	PageEvent int
	// a function that handles a page event. this should be implemented by an
	// extension and then registered. it will get executed when the event is
	// triggered
	PageEventHandler func(*Page) error
	// a map of all handlers functions registered for each page event.
	PageEventsMap map[PageEvent][]PageEventHandler
)

// List of page events. extensions can use these events to register a function
// to be executed when this event is triggered. extensions that require to be
// notified when the page is created or overwritten or deleted should register
// an event handler for the interesting events.
const (
	BeforeWrite PageEvent = iota
	AfterWrite
	AfterDelete
)

// a List of directories that should be ignored by directory walking function.
// for example the versioning extension can register `.versions` directory to be
// ignored
var ignoredDirs = []*regexp.Regexp{}

// Register a pattern to be ignored when walking directories.
func IGNORE_DIR(r *regexp.Regexp) {
	ignoredDirs = append(ignoredDirs, r)
}

// a map to keep all page events and respective list of event handlers
var PageEvents = PageEventsMap{}

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
	dat, err := ioutil.ReadFile(p.FileName())
	if err != nil {
		return ""
	}
	return string(dat)
}

// Deletes the file and makes sure it triggers the AfterDelete event
func (p *Page) Delete() bool {
	defer PageEvents.Trigger(AfterDelete, p)

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
	PageEvents.Trigger(BeforeWrite, p)
	defer PageEvents.Trigger(AfterWrite, p)

	name := p.FileName()
	os.MkdirAll(filepath.Dir(name), 0700)

	content = strings.ReplaceAll(content, "\r\n", "\n")
	if err := ioutil.WriteFile(name, []byte(content), 0644); err != nil {
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
	if e, ok := ExtractFirstFromAST[*emojiAst.Emoji](p.AST(), emojiAst.KindEmoji); ok {
		return string(e.Value.Unicode)
	}

	return ""
}

// This is a function that takes an AST node and walks the tree depth first
// recursively calling itself in search for a node of a specific kind
// can be used to find first image, link, paragraph...etc
func ExtractFirstFromAST[t ast.Node](n ast.Node, kind ast.NodeKind) (found t, ok bool) {
	if n.Kind() == kind {
		if found, ok := n.(t); ok {
			return found, true
		}
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if a, ok := ExtractFirstFromAST[t](c, kind); ok {
			return a, true
		}
	}

	return found, false
}

// Extract all nodes of a specific type from the AST
func ExtractAllFromAST[t ast.Node](n ast.Node, kind ast.NodeKind) (a []t) {
	if n.Kind() == kind {
		typed, _ := n.(t)
		a = []t{typed}
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		a = append(a, ExtractAllFromAST[t](c, kind)...)
	}

	return a
}

var walkPagesCache []*Page
var walkPagesCacheMutex sync.RWMutex

// this function is useful to iterate on all available pages. many extensions
// uses it to get all pages and maybe parse them and extract needed information
func WalkPages(ctx context.Context, f func(*Page)) {
	if walkPagesCache == nil {
		populateWalkPagesCache(ctx)
	}

	walkPagesCacheMutex.RLock()
	defer walkPagesCacheMutex.RUnlock()

	for _, p := range walkPagesCache {
		select {
		case <-ctx.Done():
			return
		default:
			f(p)
		}
	}
}

func clearWalkPagesCache(_ *Page) (err error) {
	walkPagesCacheMutex.Lock()
	defer walkPagesCacheMutex.Unlock()

	walkPagesCache = nil
	return nil
}

func populateWalkPagesCache(ctx context.Context) {
	walkPagesCacheMutex.Lock()
	defer walkPagesCacheMutex.Unlock()

	walkPagesCache = []*Page{}

	filepath.WalkDir(".", func(name string, d fs.DirEntry, err error) error {
		if d.IsDir() && name == STATIC_DIR_PATH {
			return fs.SkipDir
		}

		if d.IsDir() {
			for _, v := range ignoredDirs {
				if v.MatchString(name) {
					return fs.SkipDir
				}
			}

			return nil
		}

		select {

		case <-ctx.Done():
			return errors.New("Context stopped")

		default:
			ext := path.Ext(name)
			basename := name[:len(name)-len(ext)]

			if ext == ".md" {
				walkPagesCache = append(walkPagesCache, &Page{Name: basename})
			}

		}

		return nil
	})
}

// Register an event handler to be executed when PageEvent is triggered.
// extensions can use this to register hooks under specific page events.
// extensions that keeps a cached version of the pages list for example needs to
// register handlers to update its cache
func (c PageEventsMap) Listen(e PageEvent, h PageEventHandler) {
	if _, ok := c[e]; !ok {
		c[e] = []PageEventHandler{}
	}

	c[e] = append(c[e], h)
}

// Trigger event handlers for a specific page event. page methods use this function to trigger all registered handlers when the page is edited or deleted for example.
func (c PageEventsMap) Trigger(e PageEvent, p *Page) {
	if _, ok := c[e]; !ok {
		return
	}

	for _, h := range c[e] {
		if err := h(p); err != nil {
			log.Printf("Executing Event %#v handler %#v failed with error: %s\n", e, h, err)
		}
	}
}

// PREPROCESSORS =======================

// A PreProcessor is a function that takes the whole page content and returns a
// modified version of the content. extensions should define this type and
// register is so that when page is rendered it will execute all of them in
// order like a pipeline each function output is passed as an input to the next.
// at the end the last preprocessor output is then rendered to HTML
type (
	PreProcessor func(string) string
)

// List of registered preprocessor functions
var (
	preProcessors = []PreProcessor{}
)

// Register a PREPROCESSOR function. extensions should use this function to
// register a preprocessor.
func PREPROCESSOR(f PreProcessor) { preProcessors = append(preProcessors, f) }

// This function take the page content and pass it through all registered
// preprocessors and return the last preprocessor output to the caller
func preProcess(content string) string {
	for _, v := range preProcessors {
		content = v(content)
	}

	return content
}
