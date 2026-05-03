package extension

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	gast "github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
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

func TestFootnoteConfig_SetOption(t *testing.T) {
	tests := []struct {
		name     string
		option   renderer.OptionName
		value    any
		validate func(*testing.T, *FootnoteConfig)
	}{
		{
			name:   "optFootnoteIDPrefix",
			option: renderer.OptionName("FootnoteIDPrefix"),
			value:  []byte("test-prefix-"),
			validate: func(t *testing.T, c *FootnoteConfig) {
				if string(c.IDPrefix) != "test-prefix-" {
					t.Errorf("expected IDPrefix to be 'test-prefix-', got %q", c.IDPrefix)
				}
			},
		},
		{
			name:   "optFootnoteIDPrefixFunction",
			option: renderer.OptionName("FootnoteIDPrefixFunction"),
			value: func(n gast.Node) []byte {
				return []byte("dynamic-")
			},
			validate: func(t *testing.T, c *FootnoteConfig) {
				if c.IDPrefixFunction == nil {
					t.Error("expected IDPrefixFunction to be set")
				} else {
					result := c.IDPrefixFunction(nil)
					if string(result) != "dynamic-" {
						t.Errorf("expected IDPrefixFunction to return 'dynamic-', got %q", result)
					}
				}
			},
		},
		{
			name:   "optFootnoteLinkTitle",
			option: renderer.OptionName("FootnoteLinkTitle"),
			value:  []byte("Link to footnote ^^"),
			validate: func(t *testing.T, c *FootnoteConfig) {
				if string(c.LinkTitle) != "Link to footnote ^^" {
					t.Errorf("expected LinkTitle to be 'Link to footnote ^^', got %q", c.LinkTitle)
				}
			},
		},
		{
			name:   "optFootnoteBacklinkTitle",
			option: renderer.OptionName("FootnoteBacklinkTitle"),
			value:  []byte("Return to content"),
			validate: func(t *testing.T, c *FootnoteConfig) {
				if string(c.BacklinkTitle) != "Return to content" {
					t.Errorf("expected BacklinkTitle to be 'Return to content', got %q", c.BacklinkTitle)
				}
			},
		},
		{
			name:   "optFootnoteLinkClass",
			option: renderer.OptionName("FootnoteLinkClass"),
			value:  []byte("custom-footnote-link"),
			validate: func(t *testing.T, c *FootnoteConfig) {
				if string(c.LinkClass) != "custom-footnote-link" {
					t.Errorf("expected LinkClass to be 'custom-footnote-link', got %q", c.LinkClass)
				}
			},
		},
		{
			name:   "optFootnoteBacklinkClass",
			option: renderer.OptionName("FootnoteBacklinkClass"),
			value:  []byte("custom-backlink"),
			validate: func(t *testing.T, c *FootnoteConfig) {
				if string(c.BacklinkClass) != "custom-backlink" {
					t.Errorf("expected BacklinkClass to be 'custom-backlink', got %q", c.BacklinkClass)
				}
			},
		},
		{
			name:   "optFootnoteBacklinkHTML",
			option: renderer.OptionName("FootnoteBacklinkHTML"),
			value:  []byte("&uarr;"),
			validate: func(t *testing.T, c *FootnoteConfig) {
				if string(c.BacklinkHTML) != "&uarr;" {
					t.Errorf("expected BacklinkHTML to be '&uarr;', got %q", c.BacklinkHTML)
				}
			},
		},
		{
			name:   "unknown option delegates to Config",
			option: renderer.OptionName("UnknownOption"),
			value:  "test-value",
			validate: func(t *testing.T, c *FootnoteConfig) {
				// Verify it was passed to underlying Config by checking no panic occurred.
				// The html.Config should handle unknown options gracefully.
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config := NewFootnoteConfig()
			config.SetOption(tc.option, tc.value)
			tc.validate(t, &config)
		})
	}
}

func TestFootnoteBlockParser_CanAcceptIndentedLine(t *testing.T) {
	parser := NewFootnoteBlockParser()
	if parser.(*footnoteBlockParser).CanAcceptIndentedLine() {
		t.Error("CanAcceptIndentedLine should return false")
	}
}
