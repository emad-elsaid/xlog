package main

import (
	"bytes"
	"regexp"
	"sort"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	"golang.org/x/net/context"
)

func init() {
	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&AutolinkPages{}, -1),
	))
}

type AutolinkPages struct{}

func (h *AutolinkPages) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindText, renderLinkToPages)
	reg.Register(ast.KindAutoLink, renderAutoLink)
}

func renderLinkToPages(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Text)
	segment := n.Segment
	if n.IsRaw() {
		w.Write(segment.Value(source))
	} else {

		if n.Parent().Kind() == ast.KindLink {
			w.Write(n.Text(source))
			return ast.WalkContinue, nil
		}

		pages := []*Page{}
		WalkPages(context.Background(), func(p *Page) {
			pages = append(pages, p)
		})

		sort.Sort(fileInfoByNameLength(pages))
		text := string(segment.Value(source))

		for _, p := range pages {
			reg := regexp.MustCompile(`(?imU)(^|\W)(` + regexp.QuoteMeta(p.Name) + `)(\W|$)`)
			text = reg.ReplaceAllString(text, `$1<a href="`+p.Name+`">$2</a>$3`)
		}

		w.Write([]byte(text))
	}

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
