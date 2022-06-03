package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"golang.org/x/net/context"
)

func init() {
	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&AutolinkPages{}, -1),
	))
	MarkDownRenderer.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&AutolinkPages{}, 999),
	))
	PageEvents.Listen(AfterWrite, UpdatePagesList)
	PageEvents.Listen(AfterDelete, UpdatePagesList)
}

type AutolinkPages struct{}

func (h *AutolinkPages) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindPageLink, renderPageLink)
	reg.Register(ast.KindAutoLink, renderAutoLink)
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

	var m int
	var url string

	for _, p := range autolinkPages {
		if len(line) < len(p.Name) {
			continue
		}

		// Found a page
		if strings.EqualFold(string(line[0:len(p.Name)]), p.Name) {
			url = p.Name
			m = len(p.Name)
			break
		}
	}

	if m == 0 || // can't find a page match
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
		url:   url,
		value: n,
	}
	return link
}

var KindPageLink = ast.NewNodeKind("PageLink")

type PageLink struct {
	ast.BaseInline
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
	_, _ = w.WriteString(`<a href="`)
	url := []byte(n.url)
	label := n.value.Text(source)
	_, _ = w.Write(util.EscapeHTML(util.URLEscape(url, false)))
	_, _ = w.WriteString(`">`)
	_, _ = w.Write(util.EscapeHTML(label))
	_, _ = w.WriteString(`</a>`)
	return ast.WalkContinue, nil
}

func renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*ast.AutoLink)
	if !entering {
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString(`<a href="`)
	url := n.URL(source)
	label := n.Label(source)
	limit := 30
	if len(label) > limit {
		label = []byte(string(label[0:limit]) + "…")
	}

	if n.AutoLinkType == ast.AutoLinkEmail && !bytes.HasPrefix(bytes.ToLower(url), []byte("mailto:")) {
		_, _ = w.WriteString("mailto:")
	}
	_, _ = w.Write(util.EscapeHTML(util.URLEscape(url, false)))
	if n.Attributes() != nil {
		_ = w.WriteByte('"')
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString(`">`)
	}
	_, _ = w.Write(util.EscapeHTML(label))
	_, _ = w.WriteString(`</a>`)
	return ast.WalkContinue, nil
}

type fileInfoByNameLength []*Page

func (a fileInfoByNameLength) Len() int           { return len(a) }
func (a fileInfoByNameLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fileInfoByNameLength) Less(i, j int) bool { return len(a[i].Name) > len(a[j].Name) }
