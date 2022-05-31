package main

import (
	"context"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func init() {
	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&HashTag{}, 0),
	))
	MarkDownRenderer.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&HashTag{}, 999),
	))
	SIDEBAR(hashtagsSidebar)

	GET(`/\+/tag/{tag}`, tagHandler)
}

var hashtagReg = regexp.MustCompile(`(?imU)#([[:alpha:]]\w+)(\W|$)`)

type HashTag struct {
	ast.BaseInline
	value []byte
}

func (h *HashTag) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindHashTag, renderHashtag)
}

func (h *HashTag) Trigger() []byte {
	return []byte{'#'}
}

func (h *HashTag) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < 1 {
		return nil
	}
	i := 1
	for ; i < len(line); i++ {
		c := line[i]
		if !(util.IsAlphaNumeric(c) || c == '_' || c == '-') {
			break
		}
	}
	if i > len(line) || i == 1 {
		return nil
	}
	block.Advance(i)
	tag := line[1:i]
	return &HashTag{value: tag}
}

func (h *HashTag) Dump(source []byte, level int) {
	m := map[string]string{
		"value": fmt.Sprintf("%#v", h.value),
	}
	ast.DumpHelper(h, source, level, m, nil)
}

var KindHashTag = ast.NewNodeKind("Hashtag")

func (h *HashTag) Kind() ast.NodeKind {
	return KindHashTag
}

func renderHashtag(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering || n.Kind() != KindHashTag {
		return ast.WalkContinue, nil
	}

	tag := n.(*HashTag)
	fmt.Fprintf(writer, `<a href="/+/tag/%s" class="tag is-info">#%s</a>`, tag.value, tag.value)
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
