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

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/extensions/shortcode"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
)

//go:embed templates
var templates embed.FS

const (
	keyHashtags  = "hashtags"
	nameHashtags = "Hashtags"
	keyPages     = "pages"
	iconTags     = "fa-solid fa-tags"
	pathTags     = "/+/tags"
)

func init() {
	h := Hashtags{
		pages: make(map[xlog.Page][]*HashTag),
	}

	xlog.RegisterExtension(&h)
}

// Hashtags extension provides hashtag parsing and navigation functionality.
// It maintains a thread-safe cache of parsed hashtags per page.
type Hashtags struct {
	pages map[xlog.Page][]*HashTag
	mu    sync.Mutex
}

// Name returns the unique identifier for the hashtags extension.
func (*Hashtags) Name() string { return keyHashtags }

// Init registers routes, widgets, templates, and shortcodes for the hashtags extension.
func (h *Hashtags) Init() {
	xlog.Get(pathTags, h.tagsHandler)
	xlog.Get(`/+/tag/{tag}`, h.tagHandler)
	xlog.RegisterWidget(xlog.WidgetAfterView, 1, h.relatedPages)
	xlog.RegisterBuildPage(pathTags, true)
	xlog.RegisterLink(links)
	xlog.RegisterTemplate(templates, "templates")
	shortcode.RegisterShortCode("hashtag-pages", shortcode.ShortCode{Render: h.hashtagPages})
	shortcode.RegisterShortCode("hashtag-pages-grid", shortcode.ShortCode{Render: h.hashtagPagesGrid})

	xlog.Listen(xlog.PageChanged, h.PageChanged)
	xlog.Listen(xlog.PageDeleted, h.PageDeleted)

	xlog.MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&HashTag{}, 0),
	))
	xlog.MarkdownConverter().Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&HashTag{}, 999),
	))
}

// PageChanged clears the hashtag cache for a page when its content changes.
func (h *Hashtags) PageChanged(p xlog.Page) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.pages, p)

	return nil
}

// PageDeleted clears the hashtag cache for a deleted page.
func (h *Hashtags) PageDeleted(p xlog.Page) error {
	return h.PageChanged(p)
}

func (h *Hashtags) hashtagsFor(p xlog.Page) []*HashTag {
	h.mu.Lock()
	defer h.mu.Unlock()

	if tags, ok := h.pages[p]; ok {
		return tags
	}

	_, tree := p.AST()
	tags := xlog.FindAllInAST[*HashTag](tree)
	h.pages[p] = tags

	return tags
}

func (h *Hashtags) tagsHandler(r xlog.Request) xlog.Output {
	tags := map[string][]xlog.Page{}
	var lck sync.Mutex

	xlog.EachPage(r.Context(), func(a xlog.Page) {
		set := map[string]bool{}
		_, tree := a.AST()
		hashes := xlog.FindAllInAST[*HashTag](tree)
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
				tags[val] = []xlog.Page{a}
			}
			lck.Unlock()
		}
	})

	return xlog.Render("tags", xlog.Locals{
		"page": xlog.DynamicPage{NameVal: nameHashtags},
		"tags": tags,
	})
}

func (h *Hashtags) tagHandler(r xlog.Request) xlog.Output {
	tag := r.PathValue("tag")

	return xlog.Render("tag", xlog.Locals{
		"page":   xlog.DynamicPage{NameVal: "#" + tag},
		keyPages: h.tagPages(r.Context(), tag),
	})
}

func (h *Hashtags) tagPages(ctx context.Context, hashtag string) []xlog.Page {
	uniqHandle := unique.Make(strings.ToLower(hashtag))

	return xlog.MapPage(ctx, func(p xlog.Page) xlog.Page {
		if p.Name() == xlog.Config.Index {
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

func (h *Hashtags) relatedPages(p xlog.Page) template.HTML {
	if p.Name() == xlog.Config.Index {
		return ""
	}

	_, tree := p.AST()
	found_hashtags := xlog.FindAllInAST[*HashTag](tree)
	hashtags := map[unique.Handle[string]]bool{}
	for _, v := range found_hashtags {
		hashtags[v.unique] = true
	}

	pages := xlog.MapPage(context.Background(), func(rp xlog.Page) xlog.Page {
		if rp.Name() == p.Name() {
			return nil
		}

		_, tree := rp.AST()
		page_hashtags := xlog.FindAllInAST[*HashTag](tree)
		for _, h := range page_hashtags {
			if _, ok := hashtags[h.unique]; ok {
				return rp
			}
		}

		return nil
	})

	return xlog.Partial("related-hashtags-pages", xlog.Locals{
		keyPages: pages,
	})
}

func (h *Hashtags) hashtagPages(hashtag xlog.Markdown) template.HTML {
	hashtag_value := strings.Trim(string(hashtag), "# \n")
	pages := h.tagPages(context.Background(), hashtag_value)

	slices.SortFunc(pages, func(a, b xlog.Page) int {
		if modtime := b.ModTime().Compare(a.ModTime()); modtime != 0 {
			return modtime
		}

		return strings.Compare(a.Name(), b.Name())
	})

	output := xlog.Partial("hashtag-pages", xlog.Locals{keyPages: pages})
	// #nosec G203 - xlog.Partial executes html/template which auto-escapes, safe to convert
	return template.HTML(output)
}

func (h *Hashtags) hashtagPagesGrid(hashtag xlog.Markdown) template.HTML {
	hashtag_value := strings.Trim(string(hashtag), "# \n")
	pages := h.tagPages(context.Background(), hashtag_value)

	slices.SortFunc(pages, func(a, b xlog.Page) int {
		if modtime := b.ModTime().Compare(a.ModTime()); modtime != 0 {
			return modtime
		}

		return strings.Compare(a.Name(), b.Name())
	})

	output := xlog.Partial("hashtag-pages-grid", xlog.Locals{keyPages: pages})
	// #nosec G203 - xlog.Partial executes html/template which auto-escapes, safe to convert
	return template.HTML(output)
}

func links(xlog.Page) []xlog.Command {
	return []xlog.Command{link{}}
}

type link struct{}

// Icon returns the FontAwesome icon class for the hashtags link.
func (l link) Icon() string { return iconTags }

// Name returns the display name for the hashtags navigation link.
func (l link) Name() string { return nameHashtags }

// Attrs returns HTML attributes for rendering the hashtags link.
func (l link) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"href": pathTags,
	}
}

// HashTag represents a hashtag node in the markdown AST.
// It stores both the raw value and a case-insensitive unique handle for matching.
type HashTag struct {
	ast.BaseInline
	value  []byte
	unique unique.Handle[string]
}

// RegisterFuncs registers the render function for hashtag nodes.
func (h *HashTag) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindHashTag, renderHashtag)
}

// Trigger returns the character that triggers hashtag parsing ('#').
func (h *HashTag) Trigger() []byte {
	return []byte{'#'}
}

// Parse extracts a hashtag from the markdown text stream.
// It accepts letters, numbers, dashes, and underscores, stopping at whitespace or punctuation.
// Hashtags are case-insensitive and support Unicode characters.
func (h *HashTag) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	l, _ := block.PeekLine()
	if len(l) < 1 {
		return nil
	}

	line := string(l)

	var i int
	for ui, c := range line {
		if ui == 0 {
			i += utf8.RuneLen(c)
			continue
		}

		if (!unicode.In(c, unicode.Letter, unicode.Number, unicode.Dash) && c != '_') || unicode.IsSpace(c) {
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

// Dump outputs debug information about the hashtag node for AST introspection.
func (h *HashTag) Dump(source []byte, level int) {
	m := map[string]string{
		"value": fmt.Sprintf("%#v", h.value),
	}
	ast.DumpHelper(h, source, level, m, nil)
}

// KindHashTag uniquely identifies hashtag nodes in the AST.
var KindHashTag = ast.NewNodeKind("Hashtag")

// Kind returns the unique node kind identifier for hashtag nodes.
func (h *HashTag) Kind() ast.NodeKind {
	return KindHashTag
}

func renderHashtag(writer util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering || n.Kind() != KindHashTag {
		return ast.WalkContinue, nil
	}

	tag := n.(*HashTag)
	if _, err := fmt.Fprintf(writer, `<a href="/+/tag/%s" class="tag"><span class="icon"><i class="fa-solid fa-tag"></i></span><span>%s</span></a>`, tag.value, tag.value); err != nil {
		return ast.WalkStop, err
	}
	xlog.RegisterBuildPage(fmt.Sprintf("/+/tag/%s", tag.value), true)
	xlog.RegisterBuildPage(fmt.Sprintf("/+/tag/%s", strings.ToLower(string(tag.value))), true)
	return ast.WalkContinue, nil
}
