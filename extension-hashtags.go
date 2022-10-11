package main

import (
	"context"
	"fmt"
	"html/template"
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
	WIDGET(SIDEBAR_WIDGET, hashtagsSidebar)
	WIDGET(AFTER_VIEW_WIDGET, relatedHashtagsPages)
	AUTOCOMPLETE(hashtagAutocomplete)

	GET(`/\+/tags`, tagsHandler)
	EXTENSION_PAGE("/+/tags")
	GET(`/\+/tag/{tag}`, tagHandler)
}

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
	fmt.Fprintf(writer, `<a href="/+/tag/%s" class="tag is-info is-light">#%s</a>`, tag.value, tag.value)
	EXTENSION_PAGE(fmt.Sprintf("/+/tag/%s", tag.value))
	return ast.WalkContinue, nil
}

func extractHashtags(n ast.Node) []*HashTag {
	a := []*HashTag{}

	if n.Kind() == KindHashTag {
		tag, _ := n.(*HashTag)
		a = []*HashTag{tag}
	}

	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		a = append(a, extractHashtags(c)...)
		if c == n.LastChild() {
			break
		}
	}

	return a
}

func tagsHandler(_ Response, r Request) Output {
	tags := map[string][]string{}
	WalkPages(context.Background(), func(a *Page) {
		set := map[string]bool{}
		hashes := extractHashtags(a.AST())
		for _, v := range hashes {
			val := strings.ToLower(string(v.value))

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

	return Render("extension/tags", Locals{
		"title":   "Hashtags",
		"tags":    tags,
		"sidebar": renderWidget(SIDEBAR_WIDGET, nil, r),
	})
}

func tagHandler(w Response, r Request) Output {
	vars := VARS(r)
	tag := vars["tag"]

	return Render("extension/tag", Locals{
		"title":   "#" + tag,
		"pages":   tagPages(r.Context(), tag),
		"sidebar": renderWidget(SIDEBAR_WIDGET, nil, r),
	})
}

func tagPages(ctx context.Context, keyword string) []*Page {
	results := []*Page{}

	WalkPages(ctx, func(p *Page) {
		tags := extractHashtags(p.AST())
		for _, t := range tags {
			if strings.EqualFold(string(t.value), keyword) {
				results = append(results, p)
				break
			}
		}
	})

	return results
}

func hashtagsSidebar(p *Page, r Request) template.HTML {
	return template.HTML(partial("extension/tags-sidebar", nil))
}

func relatedHashtagsPages(p *Page, r Request) template.HTML {
	if p.Name == "index" {
		return ""
	}

	found_hashtags := extractHashtags(p.AST())
	hashtags := map[string]bool{}
	for _, v := range found_hashtags {
		hashtags[string(v.value)] = true
	}

	pages := []*Page{}

	WalkPages(context.Background(), func(rp *Page) {
		if rp.Name == p.Name {
			return
		}

		page_hashtags := extractHashtags(rp.AST())
		for _, h := range page_hashtags {
			if _, ok := hashtags[string(h.value)]; ok {
				pages = append(pages, rp)
				return
			}
		}
	})

	return template.HTML(partial("extension/related-hashtags-pages", Locals{
		"pages": pages,
	}))
}

func hashtagAutocomplete() *Autocomplete {
	a := &Autocomplete{
		StartChar:   "#",
		Suggestions: []*Suggestion{},
	}

	set := map[string]bool{}
	WalkPages(context.Background(), func(a *Page) {
		hashes := extractHashtags(a.AST())
		for _, v := range hashes {
			set[strings.ToLower(string(v.value))] = true
		}
	})

	for t := range set {
		a.Suggestions = append(a.Suggestions, &Suggestion{
			Text:        "#" + t,
			DisplayText: t,
		})
	}

	return a
}
