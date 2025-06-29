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
	"unique"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/extensions/shortcode"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
)

//go:embed templates
var templates embed.FS

func init() {
	h := Hashtags{
		pages: make(map[Page][]*HashTag),
	}

	RegisterExtension(&h)
}

type Hashtags struct {
	pages map[Page][]*HashTag
	mu    sync.Mutex
}

func (*Hashtags) Name() string { return "hashtags" }
func (h *Hashtags) Init() {
	Get(`/+/tags`, h.tagsHandler)
	Get(`/+/tag/{tag}`, h.tagHandler)
	RegisterWidget(WidgetAfterView, 1, h.relatedPages)
	RegisterBuildPage("/+/tags", true)
	RegisterLink(links)
	RegisterTemplate(templates, "templates")
	shortcode.RegisterShortCode("hashtag-pages", shortcode.ShortCode{Render: h.hashtagPages})
	shortcode.RegisterShortCode("hashtag-pages-grid", shortcode.ShortCode{Render: h.hashtagPagesGrid})

	Listen(PageChanged, h.PageChanged)
	Listen(PageDeleted, h.PageDeleted)

	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&HashTag{}, 0),
	))
	MarkdownConverter().Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&HashTag{}, 999),
	))
}

func (h *Hashtags) PageChanged(p Page) error {
	delete(h.pages, p)

	return nil
}

func (h *Hashtags) PageDeleted(p Page) error {
	return h.PageChanged(p)
}

func (h *Hashtags) hashtagsFor(p Page) []*HashTag {
	h.mu.Lock()
	defer h.mu.Unlock()

	if tags, ok := h.pages[p]; ok {
		return tags
	}

	_, tree := p.AST()
	tags := FindAllInAST[*HashTag](tree)
	h.pages[p] = tags

	return tags
}

func (h *Hashtags) tagsHandler(r Request) Output {
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

func (h *Hashtags) tagHandler(r Request) Output {
	tag := r.PathValue("tag")

	return Render("tag", Locals{
		"page":  DynamicPage{NameVal: "#" + tag},
		"pages": h.tagPages(r.Context(), tag),
	})
}

func (h *Hashtags) tagPages(ctx context.Context, hashtag string) []Page {
	uniqHandle := unique.Make(strings.ToLower(hashtag))

	return MapPage(ctx, func(p Page) Page {
		if p.Name() == Config.Index {
			return nil
		}

		tags := h.hashtagsFor(p)
		for _, t := range tags {
			if uniqHandle == t.unique {
				return p
			}
		}

		return nil
	})
}

func (h *Hashtags) relatedPages(p Page) template.HTML {
	if p.Name() == Config.Index {
		return ""
	}

	_, tree := p.AST()
	found_hashtags := FindAllInAST[*HashTag](tree)
	hashtags := map[unique.Handle[string]]bool{}
	for _, v := range found_hashtags {
		hashtags[v.unique] = true
	}

	pages := MapPage(context.Background(), func(rp Page) Page {
		if rp.Name() == p.Name() {
			return nil
		}

		_, tree := rp.AST()
		page_hashtags := FindAllInAST[*HashTag](tree)
		for _, h := range page_hashtags {
			if _, ok := hashtags[h.unique]; ok {
				return rp
			}
		}

		return nil
	})

	return Partial("related-hashtags-pages", Locals{
		"pages": pages,
	})
}

func (h *Hashtags) hashtagPages(hashtag Markdown) template.HTML {
	hashtag_value := strings.Trim(string(hashtag), "# \n")
	pages := h.tagPages(context.Background(), hashtag_value)

	slices.SortFunc(pages, func(a, b Page) int {
		if modtime := b.ModTime().Compare(a.ModTime()); modtime != 0 {
			return modtime
		}

		return strings.Compare(a.Name(), b.Name())
	})

	output := Partial("hashtag-pages", Locals{"pages": pages})
	return template.HTML(output)
}

func (h *Hashtags) hashtagPagesGrid(hashtag Markdown) template.HTML {
	hashtag_value := strings.Trim(string(hashtag), "# \n")
	pages := h.tagPages(context.Background(), hashtag_value)

	slices.SortFunc(pages, func(a, b Page) int {
		if modtime := b.ModTime().Compare(a.ModTime()); modtime != 0 {
			return modtime
		}

		return strings.Compare(a.Name(), b.Name())
	})

	output := Partial("hashtag-pages-grid", Locals{"pages": pages})
	return template.HTML(output)
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
	value  []byte
	unique unique.Handle[string]
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
	return &HashTag{
		value:  []byte(tag),
		unique: unique.Make(strings.ToLower(tag)),
	}
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
	fmt.Fprintf(writer, `<a href="/+/tag/%s" class="tag"><span class="icon"><i class="fa-solid fa-tag"></i></span><span>%s</span></a>`, tag.value, tag.value)
	RegisterBuildPage(fmt.Sprintf("/+/tag/%s", tag.value), true)
	RegisterBuildPage(fmt.Sprintf("/+/tag/%s", strings.ToLower(string(tag.value))), true)
	return ast.WalkContinue, nil
}
