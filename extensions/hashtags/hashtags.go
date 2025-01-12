package hashtags

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"slices"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/extensions/shortcode"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

//go:embed templates
var templates embed.FS

func init() {
	RegisterExtension(Hashtags{})
}

type Hashtags struct{}

func (Hashtags) Name() string { return "hashtags" }
func (Hashtags) Init() {
	Get(`/+/tags`, tagsHandler)
	Get(`/+/tag/{tag}`, tagHandler)
	RegisterWidget(WidgetAfterView, 1, relatedPages)
	RegisterBuildPage("/+/tags", true)
	RegisterLink(links)
	RegisterTemplate(templates, "templates")
	shortcode.RegisterShortCode("hashtag-pages", shortcode.ShortCode{Render: hashtagPages})

	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&HashTag{}, 0),
	))
	MarkdownConverter().Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&HashTag{}, 999),
	))
}

func links(Page) []Command {
	return []Command{link{}}
}

type link struct{}

func (l link) Icon() string { return "fa-solid fa-tags" }
func (l link) Name() string { return "Hashtags" }
func (l link) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"href": "/+/tags",
	}
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
	l, _ := block.PeekLine()
	if len(l) < 1 {
		return nil
	}

	var line string = string(l)

	var i int
	for ui, c := range line {
		if ui == 0 {
			i += utf8.RuneLen(c)
			continue
		}

		if !(unicode.In(c, unicode.Letter, unicode.Number, unicode.Dash) || c == '_') || unicode.IsSpace(c) {
			break
		}

		i += utf8.RuneLen(c)
	}
	if i > len(line) || i == 1 {
		return nil
	}
	block.Advance(i)
	tag := line[1:i]
	return &HashTag{value: []byte(tag)}
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
	fmt.Fprintf(writer, `<a href="/+/tag/%s" class="tag">%s</a>`, tag.value, tag.value)
	RegisterBuildPage(fmt.Sprintf("/+/tag/%s", tag.value), true)
	RegisterBuildPage(fmt.Sprintf("/+/tag/%s", strings.ToLower(string(tag.value))), true)
	return ast.WalkContinue, nil
}

func tagsHandler(r Request) Output {
	tags := map[string][]Page{}
	var lck sync.Mutex

	EachPage(r.Context(), func(a Page) {
		set := map[string]bool{}
		_, tree := a.AST()
		hashes := FindAllInAST[*HashTag](tree)
		for _, v := range hashes {
			val := strings.ToLower(string(v.value))

			// don't use same tag twice for same page
			if _, ok := set[val]; ok {
				continue
			}

			set[val] = true

			lck.Lock()
			if ps, ok := tags[val]; ok {
				tags[val] = append(ps, a)
			} else {
				tags[val] = []Page{a}
			}
			lck.Unlock()
		}
	})

	return Render("tags", Locals{
		"page": DynamicPage{NameVal: "Hashtags"},
		"tags": tags,
	})
}

func tagHandler(r Request) Output {
	tag := r.PathValue("tag")

	return Render("tag", Locals{
		"page":  DynamicPage{NameVal: "#" + tag},
		"pages": tagPages(r.Context(), tag),
	})
}

func tagPages(ctx context.Context, keyword string) []Page {
	return MapPage(ctx, func(p Page) Page {
		if p.Name() == Config.Index {
			return nil
		}

		_, tree := p.AST()
		tags := FindAllInAST[*HashTag](tree)
		for _, t := range tags {
			if strings.EqualFold(string(t.value), keyword) {
				return p
			}
		}

		return nil
	})
}

func relatedPages(p Page) template.HTML {
	if p.Name() == Config.Index {
		return ""
	}

	_, tree := p.AST()
	found_hashtags := FindAllInAST[*HashTag](tree)
	hashtags := map[string]bool{}
	for _, v := range found_hashtags {
		hashtags[strings.ToLower(string(v.value))] = true
	}

	pages := MapPage(context.Background(), func(rp Page) Page {
		if rp.Name() == p.Name() {
			return nil
		}

		_, tree := rp.AST()
		page_hashtags := FindAllInAST[*HashTag](tree)
		for _, h := range page_hashtags {
			if _, ok := hashtags[strings.ToLower(string(h.value))]; ok {
				return rp
			}
		}

		return nil
	})

	return Partial("related-hashtags-pages", Locals{
		"pages": pages,
	})
}

func hashtagPages(hashtag Markdown) template.HTML {
	hashtag_value := strings.Trim(string(hashtag), "# \n")
	pages := tagPages(context.Background(), hashtag_value)

	slices.SortFunc(pages, func(a, b Page) int {
		if modtime := b.ModTime().Compare(a.ModTime()); modtime != 0 {
			return modtime
		}

		return strings.Compare(a.Name(), b.Name())
	})

	output := Partial("hashtag-pages", Locals{"pages": pages})
	return template.HTML(output)
}
