package main

import (
	"context"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

func init() {
	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&HashTag{}, 0),
	))
	SIDEBAR(hashtagsSidebar)

	GET("/+/tag/{tag}", tagHandler)
}

var hashtagReg = regexp.MustCompile(`(?imU)#([[:alpha:]]\w+)(\W|$)`)

type HashTag struct{}

func (h *HashTag) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindText, h.renderHashtag)
}

func (h *HashTag) renderHashtag(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	text := string(n.Text(source))
	text = hashtagReg.ReplaceAllString(text, `<a href="/+/tag/$1" class="tag is-info">#$1</a>$2`)

	fmt.Fprintf(writer, text)
	return ast.WalkContinue, nil
}

func tagHandler(w Response, r Request) Output {
	vars := VARS(r)
	tag := "#" + vars["tag"]

	return Render("extension/tag", Locals{
		"title":   tag,
		"results": tagPages(r.Context(), tag),
	})
}

type tagResult struct {
	Page string
	Line string
}

func tagPages(ctx context.Context, keyword string) []tagResult {
	results := []tagResult{}
	reg := regexp.MustCompile(`(?imU)^(.*` + regexp.QuoteMeta(keyword) + `.*)$`)

	WalkPages(ctx, func(p *Page) {
		match := reg.FindString(p.Content())
		if len(match) > 0 {
			results = append(results, tagResult{
				Page: p.Name,
				Line: match,
			})
		}
	})

	return results
}

func hashtagsSidebar(p *Page, r Request) template.HTML {
	tags := map[string][]string{}
	WalkPages(context.Background(), func(a *Page) {
		set := map[string]bool{}
		hashes := hashtagReg.FindAllStringSubmatch(a.Content(), -1)
		for _, v := range hashes {
			val := strings.ToLower(v[1])

			// don't use same tag twice for same page
			if _, ok := set[val]; ok {
				continue
			}

			set[val] = true
			if ps, ok := tags[val]; ok {
				tags[val] = append(ps, a.Name)
			} else {
				tags[val] = []string{a.Name}
			}
		}
	})

	return template.HTML(partial("extension/tags-sidebar", Locals{
		"tags": tags,
	}))
}
