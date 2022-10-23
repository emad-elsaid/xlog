package hashtags

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"strings"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/extensions/shortcode"

	_ "embed"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

//go:embed views
var views embed.FS

func init() {
	WIDGET(SIDEBAR_WIDGET, sidebar)
	WIDGET(AFTER_VIEW_WIDGET, relatedPages)

	GET(`/\+/tags`, tagsHandler)
	GET(`/\+/tag/{tag}`, tagHandler)

	EXTENSION_PAGE("/+/tags")

	AUTOCOMPLETE(autocompleter)
	shortcode.SHORTCODE("hashtag-pages", hashtagPages)

	fs, _ := fs.Sub(views, "views")
	VIEW(fs)

	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&HashTag{}, 0),
	))
	MarkDownRenderer.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&HashTag{}, 999),
	))
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

func tagsHandler(_ Response, r Request) Output {
	tags := map[string][]*Page{}
	WalkPages(context.Background(), func(a *Page) {
		set := map[string]bool{}
		hashes := ExtractAllFromAST[*HashTag](a.AST(), KindHashTag)
		for _, v := range hashes {
			val := strings.ToLower(string(v.value))

			// don't use same tag twice for same page
			if _, ok := set[val]; ok {
				continue
			}

			set[val] = true
			if ps, ok := tags[val]; ok {
				tags[val] = append(ps, a)
			} else {
				tags[val] = []*Page{a}
			}
		}
	})

	return Render("tags", Locals{
		"title":   "Hashtags",
		"tags":    tags,
		"sidebar": RenderWidget(SIDEBAR_WIDGET, nil, r),
	})
}

func tagHandler(w Response, r Request) Output {
	vars := Vars(r)
	tag := vars["tag"]

	return Render("tag", Locals{
		"title":   "#" + tag,
		"pages":   tagPages(r.Context(), tag),
		"sidebar": RenderWidget(SIDEBAR_WIDGET, nil, r),
	})
}

func tagPages(ctx context.Context, keyword string) []*Page {
	results := []*Page{}

	WalkPages(ctx, func(p *Page) {
		if p.Name == INDEX {
			return
		}

		tags := ExtractAllFromAST[*HashTag](p.AST(), KindHashTag)
		for _, t := range tags {
			if strings.EqualFold(string(t.value), keyword) {
				results = append(results, p)
				break
			}
		}
	})

	return results
}

func sidebar(p *Page, r Request) template.HTML {
	return template.HTML(Partial("tags-sidebar", nil))
}

func relatedPages(p *Page, r Request) template.HTML {
	if p.Name == INDEX {
		return ""
	}

	found_hashtags := ExtractAllFromAST[*HashTag](p.AST(), KindHashTag)
	hashtags := map[string]bool{}
	for _, v := range found_hashtags {
		hashtags[strings.ToLower(string(v.value))] = true
	}

	pages := []*Page{}

	WalkPages(context.Background(), func(rp *Page) {
		if rp.Name == p.Name {
			return
		}

		page_hashtags := ExtractAllFromAST[*HashTag](rp.AST(), KindHashTag)
		for _, h := range page_hashtags {
			if _, ok := hashtags[strings.ToLower(string(h.value))]; ok {
				pages = append(pages, rp)
				return
			}
		}
	})

	return template.HTML(Partial("related-hashtags-pages", Locals{
		"pages": pages,
	}))
}

func autocompleter() *Autocomplete {
	a := &Autocomplete{
		StartChar:   "#",
		Suggestions: []*Suggestion{},
	}

	set := map[string]bool{}
	WalkPages(context.Background(), func(a *Page) {
		hashes := ExtractAllFromAST[*HashTag](a.AST(), KindHashTag)
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

func hashtagPages(hashtag string) string {
	hashtag = strings.Trim(hashtag, "# ")
	pages := tagPages(context.Background(), hashtag)

	output := Partial("hashtag-pages", Locals{"pages": pages})
	output = strings.ReplaceAll(output, "\n", "")
	output = strings.TrimSpace(output)
	return output + "\n"
}
