package main

import (
	"context"
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type fileInfoByNameLength []*Page

func (a fileInfoByNameLength) Len() int           { return len(a) }
func (a fileInfoByNameLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fileInfoByNameLength) Less(i, j int) bool { return len(a[i].Name) > len(a[j].Name) }

func init() {
	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&AutolinkPages{}, -1),
	))
	MarkDownRenderer.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&AutolinkPages{}, 999),
	))
	PageEvents.Listen(AfterWrite, UpdatePagesList)
	PageEvents.Listen(AfterDelete, UpdatePagesList)

	WIDGET(AFTER_VIEW_WIDGET, backlinksSidebar)
}

type AutolinkPages struct{}

func (h *AutolinkPages) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindPageLink, renderPageLink)
}

func (_ *AutolinkPages) Trigger() []byte {
	// ' ' indicates any white spaces and a line head
	return []byte{' ', '*', '_', '~', '('}
}

var autolinkPages []*Page

func UpdatePagesList(_ *Page) (err error) {
	ps := []*Page{}
	WalkPages(context.Background(), func(p *Page) {
		ps = append(ps, p)
	})
	sort.Sort(fileInfoByNameLength(ps))
	autolinkPages = ps
	return
}

func (s *AutolinkPages) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	if pc.IsInLinkLabel() {
		return nil
	}

	if autolinkPages == nil {
		UpdatePagesList(nil)
	}

	line, segment := block.PeekLine()
	consumes := 0
	start := segment.Start
	c := line[0]
	// advance if current position is not a line head.
	if c == ' ' || c == '*' || c == '_' || c == '~' || c == '(' {
		consumes++
		start++
		line = line[1:]
	}

	var found *Page
	var m int
	var url string

	for _, p := range autolinkPages {
		if len(line) < len(p.Name) {
			continue
		}

		// Found a page
		if strings.EqualFold(string(line[0:len(p.Name)]), p.Name) {
			found = p
			url = p.Name
			m = len(p.Name)
			break
		}
	}

	if found == nil ||
		(len(line) > m && util.IsAlphaNumeric(line[m])) { // next character is word character
		block.Advance(consumes)
		return nil
	}

	if consumes != 0 {
		s := segment.WithStop(segment.Start + 1)
		ast.MergeOrAppendTextSegment(parent, s)
	}
	consumes += m
	block.Advance(consumes)
	n := ast.NewTextSegment(text.NewSegment(start, start+m))
	link := &PageLink{
		page:  found,
		url:   "/" + url,
		value: n,
	}
	return link
}

var KindPageLink = ast.NewNodeKind("PageLink")

type PageLink struct {
	ast.BaseInline
	page  *Page
	url   string
	value *ast.Text
}

func (_ *PageLink) Kind() ast.NodeKind {
	return KindPageLink
}

func (p *PageLink) Dump(source []byte, level int) {
	m := map[string]string{
		"value": fmt.Sprintf("%#v:%s", p.value, p.url),
	}
	ast.DumpHelper(p, source, level, m, nil)
}

func renderPageLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*PageLink)
	if !entering {
		return ast.WalkContinue, nil
	}

	w.WriteString(`<a href="`)
	url := []byte(n.url)
	label := n.value.Text(source)

	w.Write(util.EscapeHTML(util.URLEscape(url, false)))
	w.WriteString(`">`)

	if total, done := countTodos(n.page); total > 0 {
		isDone := ""
		if total == done {
			isDone = "is-success"
		}
		fmt.Fprintf(w, `<span class="tag is-rounded %s">%d/%d</span> `, isDone, done, total)
	}

	w.Write(util.EscapeHTML(label))
	w.WriteString(`</a>`)
	return ast.WalkContinue, nil
}

func countTodos(p *Page) (total int, done int) {
	tasks := extractTodos(p.AST())
	for _, v := range tasks {
		total++
		if v.IsChecked {
			done++
		}
	}

	return
}

func extractTodos(n ast.Node) []*east.TaskCheckBox {
	a := []*east.TaskCheckBox{}

	if n.Kind() == east.KindTaskCheckBox {
		t, _ := n.(*east.TaskCheckBox)
		a = []*east.TaskCheckBox{t}
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		a = append(a, extractTodos(c)...)
		if c == n.LastChild() {
			break
		}
	}

	return a
}

func backlinksSidebar(p *Page, r Request) template.HTML {
	pages := []*Page{}

	WalkPages(context.Background(), func(a *Page) {
		// a page shouldn't mention itself
		if a.Name == p.Name {
			return
		}

		if containLinkTo(a.AST(), p) {
			pages = append(pages, a)
		}
	})

	return template.HTML(partial("extension/backlinks", Locals{"pages": pages}))
}

func containLinkTo(n ast.Node, p *Page) bool {
	if n.Kind() == KindPageLink {
		t, _ := n.(*PageLink)
		if t.page.FileName() == p.FileName() {
			return true
		}
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		if containLinkTo(c, p) {
			return true
		}

		if c == n.LastChild() {
			break
		}
	}

	return false
}
