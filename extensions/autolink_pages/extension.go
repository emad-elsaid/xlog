package autolink_pages

import (
	"log/slog"
	"time"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func init() {
	RegisterExtension(AutoLinkPages{})
}

type AutoLinkPages struct{}

func (AutoLinkPages) Name() string { return "autolink-pages" }
func (AutoLinkPages) Init() {
	if !Config.Readonly {
		Listen(PageChanged, UpdatePagesList)
		Listen(PageDeleted, UpdatePagesList)
	}

	// Listen to BeforeCacheWarming event to pre-build trie
	// This ensures the trie is built once before concurrent AST parsing
	Listen(BeforeCacheWarming, func(Page) error {
		slog.Info("Pre-building autolink trie for fast lookups")
		start := time.Now()
		if err := UpdatePagesList(nil); err != nil {
			slog.Error("Failed to build autolink trie", "error", err)
			return err
		}
		slog.Info("Autolink trie built", "duration", time.Since(start))
		return nil
	})

	RegisterWidget(WidgetAfterView, 1, backlinksSection)
	RegisterTemplate(templates, "templates")
	MarkdownConverter().Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&pageLinkParser{}, 999),
	))
	MarkdownConverter().Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&pageLinkRenderer{}, -1),
	))
}
