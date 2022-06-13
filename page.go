package main

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
	"time"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

type (
	Page struct {
		Name string
	}

	PageEvent        int
	PageEventHandler func(*Page) error
	PageEventsMap    map[PageEvent][]PageEventHandler
)

const (
	BeforeWrite PageEvent = iota
	AfterWrite
	AfterDelete
)

var ignoredDirs = []*regexp.Regexp{}

func IGNORE_DIR(r *regexp.Regexp) {
	ignoredDirs = append(ignoredDirs, r)
}

var PageEvents = PageEventsMap{}
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

func NewPage(name string) Page {
	if name == "" {
		name = "index"
	}

	return Page{
		Name: name,
	}
}

func (p *Page) FileName() string {
	return filepath.FromSlash(p.Name) + ".md"
}

func (p *Page) Exists() bool {
	_, err := os.Stat(p.FileName())
	return err == nil
}

func (p *Page) Render() string {
	content := p.Content()
	content = preProcess(content)

	var buf bytes.Buffer
	if err := MarkDownRenderer.Convert([]byte(content), &buf); err != nil {
		return err.Error()
	}

	return buf.String()
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

func (p *Page) ModTime() time.Time {
	s, err := os.Stat(p.FileName())
	if err != nil {
		return time.Time{}
	}

	return s.ModTime()
}

func (p *Page) RTL() bool {
	return regexp.MustCompile(`\p{Arabic}`).MatchString(p.Content())
}

func (p *Page) AST() ast.Node {
	return MarkDownRenderer.Parser().Parse(text.NewReader([]byte(p.Content())))
}

func WalkPages(ctx context.Context, f func(*Page)) {
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
				f(&Page{
					Name: basename,
				})
			}

		}

		return nil
	})
}

func (c PageEventsMap) Listen(e PageEvent, h PageEventHandler) {
	if _, ok := c[e]; !ok {
		c[e] = []PageEventHandler{}
	}

	c[e] = append(c[e], h)
}

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

type (
	preProcessor func(string) string
)

var (
	preProcessors = []preProcessor{}
)

func PREPROCESSOR(f preProcessor) { preProcessors = append(preProcessors, f) }

func preProcess(content string) string {
	for _, v := range preProcessors {
		content = v(content)
	}

	return content
}
