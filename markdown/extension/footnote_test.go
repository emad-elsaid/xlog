package extension

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	gast "github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
	"github.com/emad-elsaid/xlog/markdown/testutil"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func TestFootnote(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			Footnote,
		),
	)
	testutil.DoTestCaseFile(md, "_test/footnote.txt", t, testutil.ParseCliCaseArg()...)
}

type footnoteID struct {
}

func (a *footnoteID) Transform(node *gast.Document, reader text.Reader, pc parser.Context) {
	node.Meta()["footnote-prefix"] = "article12-"
}

func TestFootnoteOptions(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			NewFootnote(
				WithFootnoteIDPrefix("article12-"),
				WithFootnoteLinkClass("link-class"),
				WithFootnoteBacklinkClass("backlink-class"),
				WithFootnoteLinkTitle("link-title-%%-^^"),
				WithFootnoteBacklinkTitle("backlink-title"),
				WithFootnoteBacklinkHTML("^"),
			),
		),
	)

	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          1,
			Description: "Footnote with options",
			Markdown: `That's some text with a footnote.[^1]

Same footnote.[^1]

Another one.[^2]

[^1]: And that's the footnote.
[^2]: Another footnote.
`,
			Expected: `<p>That's some text with a footnote.<sup id="article12-fnref:1"><a href="#article12-fn:1" class="link-class" title="link-title-2-1" role="doc-noteref">1</a></sup></p>
<p>Same footnote.<sup id="article12-fnref1:1"><a href="#article12-fn:1" class="link-class" title="link-title-2-1" role="doc-noteref">1</a></sup></p>
<p>Another one.<sup id="article12-fnref:2"><a href="#article12-fn:2" class="link-class" title="link-title-1-2" role="doc-noteref">2</a></sup></p>
<div class="footnotes" role="doc-endnotes">
<hr>
<ol>
<li id="article12-fn:1">
<p>And that's the footnote.&#160;<a href="#article12-fnref:1" class="backlink-class" title="backlink-title" role="doc-backlink">^</a>&#160;<a href="#article12-fnref1:1" class="backlink-class" title="backlink-title" role="doc-backlink">^</a></p>
</li>
<li id="article12-fn:2">
<p>Another footnote.&#160;<a href="#article12-fnref:2" class="backlink-class" title="backlink-title" role="doc-backlink">^</a></p>
</li>
</ol>
</div>`,
		},
		t,
	)

	md = markdown.New(
		markdown.WithParserOptions(
			parser.WithASTTransformers(
				util.Prioritized(&footnoteID{}, 100),
			),
		),
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			NewFootnote(
				WithFootnoteIDPrefixFunction(func(n gast.Node) []byte {
					v, ok := n.OwnerDocument().Meta()["footnote-prefix"]
					if ok {
						return util.StringToReadOnlyBytes(v.(string))
					}
					return nil
				}),
				WithFootnoteLinkClass([]byte("link-class")),
				WithFootnoteBacklinkClass([]byte("backlink-class")),
				WithFootnoteLinkTitle([]byte("link-title-%%-^^")),
				WithFootnoteBacklinkTitle([]byte("backlink-title")),
				WithFootnoteBacklinkHTML([]byte("^")),
			),
		),
	)

	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          2,
			Description: "Footnote with an id prefix function",
			Markdown: `That's some text with a footnote.[^1]

Same footnote.[^1]

Another one.[^2]

[^1]: And that's the footnote.
[^2]: Another footnote.
`,
			Expected: `<p>That's some text with a footnote.<sup id="article12-fnref:1"><a href="#article12-fn:1" class="link-class" title="link-title-2-1" role="doc-noteref">1</a></sup></p>
<p>Same footnote.<sup id="article12-fnref1:1"><a href="#article12-fn:1" class="link-class" title="link-title-2-1" role="doc-noteref">1</a></sup></p>
<p>Another one.<sup id="article12-fnref:2"><a href="#article12-fn:2" class="link-class" title="link-title-1-2" role="doc-noteref">2</a></sup></p>
<div class="footnotes" role="doc-endnotes">
<hr>
<ol>
<li id="article12-fn:1">
<p>And that's the footnote.&#160;<a href="#article12-fnref:1" class="backlink-class" title="backlink-title" role="doc-backlink">^</a>&#160;<a href="#article12-fnref1:1" class="backlink-class" title="backlink-title" role="doc-backlink">^</a></p>
</li>
<li id="article12-fn:2">
<p>Another footnote.&#160;<a href="#article12-fnref:2" class="backlink-class" title="backlink-title" role="doc-backlink">^</a></p>
</li>
</ol>
</div>`,
		},
		t,
	)
}
