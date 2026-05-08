package livecodes

import (
	"embed"
	"fmt"
	"html"
	"html/template"
	"strings"
	"sync/atomic"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
)

//go:embed js
var js embed.FS

const extensionName = "livecodes"

func init() {
	xlog.RegisterExtension(LiveCodes{})
}

type LiveCodes struct{}

func (LiveCodes) Name() string { return extensionName }
func (LiveCodes) Init() {
	xlog.RegisterStaticDir(js)
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

// transformLiveCodesBlocks transforms fenced code blocks with "livecodes" in info string.
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

			// Get the full info string (e.g., "jsx livecodes" or "python livecodes console=open")
			var infoString string
			if n.Info != nil {
				infoString = string(n.Info.Segment.Value(source))
			}

			// Check if "livecodes" appears in the info string
			if !strings.Contains(infoString, "livecodes") {
				continue
			}

			blocks = append(blocks, n)
		}

		return ast.WalkContinue, nil
	})

	for _, b := range blocks {
		var infoString string
		if b.Info != nil {
			infoString = string(b.Info.Segment.Value(source))
		}

		// Extract the language (first word before "livecodes")
		parts := strings.Fields(infoString)
		actualLang := ""
		if len(parts) > 0 && parts[0] != "livecodes" {
			actualLang = parts[0]
		}

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
	htmlStr := fmt.Sprintf(`<div class="livecodes-playground" data-lang="%s" id="%s">
<pre style="display:none;"><code>%s</code></pre>
<div class="livecodes-loading">Loading playground...</div>
</div>
`, html.EscapeString(node.language), html.EscapeString(playgroundID), html.EscapeString(content))

	// #nosec G203 - Content is escaped via html.EscapeString
	output := template.HTML(htmlStr + script)
	if _, err := w.Write([]byte(output)); err != nil {
		return ast.WalkStop, err
	}

	return ast.WalkContinue, nil
}
