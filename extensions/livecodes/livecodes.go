package livecodes

import (
	"fmt"
	"html/template"
	"strings"
	"sync/atomic"

	_ "embed"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
)

const extensionName = "livecodes"

func init() {
	xlog.RegisterExtension(LiveCodes{})
}

type LiveCodes struct{}

func (LiveCodes) Name() string { return extensionName }
func (LiveCodes) Init() {
	LiveCodes{}.Extend(xlog.MarkdownConverter())
}

func (LiveCodes) Extend(md markdown.Markdown) {
	md.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(transformLiveCodesBlocks{}, 0),
		),
	)
	md.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&liveCodesRenderer{}, 0),
		),
	)
}

//go:embed script.html
var script string

// KindLiveCodesBlock is a NodeKind for LiveCodes blocks.
var KindLiveCodesBlock = ast.NewNodeKind("LiveCodesBlock")

// LiveCodesBlock represents a code block that should be rendered as a LiveCodes playground.
type LiveCodesBlock struct {
	ast.FencedCodeBlock
	language string
}

func (l *LiveCodesBlock) Kind() ast.NodeKind {
	return KindLiveCodesBlock
}

// transformLiveCodesBlocks transforms fenced code blocks with "live" prefix to LiveCodes blocks.
type transformLiveCodesBlocks struct{}

func (t transformLiveCodesBlocks) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()
	blocks := []*ast.FencedCodeBlock{}

	// Error ignored: Walk errors are not expected in this transformation
	_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		for c := node.FirstChild(); c != nil; c = c.NextSibling() {
			n, ok := c.(*ast.FencedCodeBlock)
			if !ok {
				continue
			}

			lang := string(n.Language(source))
			// Support "live-<language>" syntax (e.g., live-js, live-python)
			if !strings.HasPrefix(lang, "live-") && !strings.HasPrefix(lang, "live.") {
				continue
			}

			blocks = append(blocks, n)
		}

		return ast.WalkContinue, nil
	})

	for _, b := range blocks {
		lang := string(b.Language(source))
		// Extract actual language (e.g., "live-js" -> "js", "live.html" -> "html")
		actualLang := strings.TrimPrefix(lang, "live-")
		actualLang = strings.TrimPrefix(actualLang, "live.")

		replacement := LiveCodesBlock{
			FencedCodeBlock: *b,
			language:        actualLang,
		}

		parent := b.Parent()
		parent.ReplaceChild(parent, b, &replacement)
	}
}

// liveCodesRenderer renders LiveCodes blocks.
type liveCodesRenderer struct{}

var playgroundCounter uint64

func (r *liveCodesRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindLiveCodesBlock, r.render)
}

func (r *liveCodesRenderer) render(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	node, ok := n.(*LiveCodesBlock)
	if !ok {
		return ast.WalkContinue, nil
	}

	lines := node.Lines()
	content := ""
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		content += string(line.Value(source))
	}

	// Generate a unique ID for this playground
	id := atomic.AddUint64(&playgroundCounter, 1)
	playgroundID := fmt.Sprintf("livecodes-%d", id)

	// Create the playground container with the code embedded
	html := fmt.Sprintf(`<div class="livecodes-playground" data-lang="%s" id="%s">
<pre style="display:none;"><code>%s</code></pre>
<div class="livecodes-loading">Loading playground...</div>
</div>
`, htmlEscape(node.language), htmlEscape(playgroundID), htmlEscape(content))

	// #nosec G203 - Content is escaped via htmlEscape function
	output := template.HTML(html + script)
	if _, err := w.Write([]byte(output)); err != nil {
		return ast.WalkStop, err
	}

	return ast.WalkContinue, nil
}

// htmlEscape escapes HTML special characters.
func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&#39;",
	)
	return replacer.Replace(s)
}
