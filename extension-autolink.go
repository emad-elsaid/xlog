package main

import (
	"fmt"
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
}

func renderLinkToPages(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	if n.Parent().Kind() == ast.KindLink {
		fmt.Fprintf(writer, string(n.Text(source)))
		return ast.WalkContinue, nil
	}

	pages := []*Page{}
	WalkPages(context.Background(), func(p *Page) {
		pages = append(pages, p)
	})

	sort.Sort(fileInfoByNameLength(pages))
	text := string(n.Text(source))

	for _, p := range pages {
		reg := regexp.MustCompile(`(?imU)(^|\W)(` + regexp.QuoteMeta(p.Name) + `)(\W|$)`)
		text = reg.ReplaceAllString(text, `$1<a href="`+p.Name+`">$2</a>$3`)
	}

	fmt.Fprintf(writer, text)
	return ast.WalkContinue, nil
}

type fileInfoByNameLength []*Page

func (a fileInfoByNameLength) Len() int           { return len(a) }
func (a fileInfoByNameLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fileInfoByNameLength) Less(i, j int) bool { return len(a[i].Name) > len(a[j].Name) }
