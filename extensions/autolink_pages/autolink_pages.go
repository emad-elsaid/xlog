package autolink_pages

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"sort"
	"strings"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

//go:embed templates
var templates embed.FS

type fileInfoByNameLength []Page

func (a fileInfoByNameLength) Len() int           { return len(a) }
func (a fileInfoByNameLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fileInfoByNameLength) Less(i, j int) bool { return len(a[i].Name()) > len(a[j].Name()) }

func init() {
	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&extension{}, -1),
	))
	MarkDownRenderer.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&extension{}, 999),
	))
	Listen(AfterWrite, UpdatePagesList)
	Listen(AfterDelete, UpdatePagesList)

	RegisterWidget(AFTER_VIEW_WIDGET, backlinksSection)
	RegisterAutocomplete(autocomplete(0))

	RegisterTemplate(templates, "templates")
}

type extension struct{}

func (h *extension) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindPageLink, render)
}

func (_ *extension) Trigger() []byte {
	// ' ' indicates any white spaces and a line head
	return []byte{' ', '*', '_', '~', '('}
}

func (s *extension) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
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

	var found Page
	var m int
	var url string

	for _, p := range autolinkPages {
		if len(line) < len(p.Name()) {
			continue
		}

		// Found a page
		if strings.EqualFold(string(line[0:len(p.Name())]), p.Name()) {
			found = p
			url = p.Name()
			m = len(p.Name())
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

var autolinkPages []Page

func UpdatePagesList(_ Page) (err error) {
	ps := []Page{}
	EachPage(context.Background(), func(p Page) {
		ps = append(ps, p)
	})
	sort.Sort(fileInfoByNameLength(ps))
	autolinkPages = ps
	return
}

var KindPageLink = ast.NewNodeKind("PageLink")

type PageLink struct {
	ast.BaseInline
	page  Page
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

func render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
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

func countTodos(p Page) (total int, done int) {
	tasks := FindAllInAST[*east.TaskCheckBox](p.AST(), east.KindTaskCheckBox)
	for _, v := range tasks {
		total++
		if v.IsChecked {
			done++
		}
	}

	return
}

func backlinksSection(p Page, r Request) template.HTML {
	if p.Name() == INDEX {
		return ""
	}

	pages := []Page{}

	EachPage(context.Background(), func(a Page) {
		// a page shouldn't mention itself
		if a.Name() == p.Name() {
			return
		}

		if containLinkTo(a.AST(), p) {
			pages = append(pages, a)
		}
	})

	return Partial("backlinks", Locals{"pages": pages})
}

func containLinkTo(n ast.Node, p Page) bool {
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

type autocomplete int

func (a autocomplete) StartChar() string {
	return "@"
}

func (a autocomplete) Suggestions() []*Suggestion {
	suggestions := []*Suggestion{}

	EachPage(context.Background(), func(p Page) {
		suggestions = append(suggestions, &Suggestion{
			Text:        p.Name(),
			DisplayText: p.Name(),
		})
	})

	return suggestions
}
