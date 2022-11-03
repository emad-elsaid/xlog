package hashtags

import (
	"context"
	"embed"
	"fmt"
	"html/template"
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

//go:embed templates
var templates embed.FS

func init() {
	Get(`/\+/tags`, tagsHandler)
	Get(`/\+/tag/{tag}`, tagHandler)
	RegisterWidget(AFTER_VIEW_WIDGET, 1, relatedPages)
	RegisterBuildPage("/+/tags", true)
	RegisterLink(func(_ Page) []Link { return []Link{link(0)} })
	RegisterAutocomplete(autocomplete(0))
	RegisterTemplate(templates, "templates")
	shortcode.ShortCode("hashtag-pages", hashtagPages)

	MarkDownRenderer.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&HashTag{}, 0),
	))
	MarkDownRenderer.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&HashTag{}, 999),
	))
}

type link int

func (l link) Icon() string { return "fa-solid fa-tags" }
func (l link) Name() string { return "Hashtags" }
func (l link) Link() string { return "/+/tags" }

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
	RegisterBuildPage(fmt.Sprintf("/+/tag/%s", tag.value), true)
	return ast.WalkContinue, nil
}

func tagsHandler(_ Response, r Request) Output {
	tags := map[string][]Page{}
	EachPage(context.Background(), func(a Page) {
		set := map[string]bool{}
		hashes := FindAllInAST[*HashTag](a.AST(), KindHashTag)
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
				tags[val] = []Page{a}
			}
		}
	})

	return Render("tags", Locals{
		"title": "Hashtags",
		"tags":  tags,
	})
}

func tagHandler(w Response, r Request) Output {
	vars := Vars(r)
	tag := vars["tag"]

	return Render("tag", Locals{
		"title": "#" + tag,
		"pages": tagPages(r.Context(), tag),
	})
}

func tagPages(ctx context.Context, keyword string) []Page {
	results := []Page{}

	EachPage(ctx, func(p Page) {
		if p.Name() == INDEX {
			return
		}

		tags := FindAllInAST[*HashTag](p.AST(), KindHashTag)
		for _, t := range tags {
			if strings.EqualFold(string(t.value), keyword) {
				results = append(results, p)
				break
			}
		}
	})

	return results
}

func relatedPages(p Page, r Request) template.HTML {
	if p.Name() == INDEX {
		return ""
	}

	found_hashtags := FindAllInAST[*HashTag](p.AST(), KindHashTag)
	hashtags := map[string]bool{}
	for _, v := range found_hashtags {
		hashtags[strings.ToLower(string(v.value))] = true
	}

	pages := []Page{}

	EachPage(context.Background(), func(rp Page) {
		if rp.Name() == p.Name() {
			return
		}

		page_hashtags := FindAllInAST[*HashTag](rp.AST(), KindHashTag)
		for _, h := range page_hashtags {
			if _, ok := hashtags[strings.ToLower(string(h.value))]; ok {
				pages = append(pages, rp)
				return
			}
		}
	})

	return Partial("related-hashtags-pages", Locals{
		"pages": pages,
	})
}

type autocomplete int

func (a autocomplete) StartChar() string {
	return "#"
}

func (a autocomplete) Suggestions() []*Suggestion {
	suggestions := []*Suggestion{}

	set := map[string]bool{}
	EachPage(context.Background(), func(a Page) {
		hashes := FindAllInAST[*HashTag](a.AST(), KindHashTag)
		for _, v := range hashes {
			set[strings.ToLower(string(v.value))] = true
		}
	})

	for t := range set {
		suggestions = append(suggestions, &Suggestion{
			Text:        "#" + t,
			DisplayText: t,
		})
	}

	return suggestions
}

func hashtagPages(hashtag string) string {
	hashtag = strings.Trim(hashtag, "# ")
	pages := tagPages(context.Background(), hashtag)

	output := string(Partial("hashtag-pages", Locals{"pages": pages}))
	output = strings.ReplaceAll(output, "\n", "")
	output = strings.TrimSpace(output)
	output += "\n"

	return output
}
