package mermaid

import (
	"fmt"
	"html/template"

	_ "embed"

	"github.com/emad-elsaid/xlog"
	shortcode "github.com/emad-elsaid/xlog/extensions/shortcode"
)

func init() {
	xlog.RegisterExtension(Mermaid{})
}

type Mermaid struct{}

func (Mermaid) Name() string { return "mermaid" }
func (Mermaid) Init() {
	shortcode.RegisterShortCode("mermaid", shortcode.ShortCode{Render: renderer})
}

//go:embed script.html
var script string

const pre = `<pre class="mermaid" style="background: transparent;text-align:center;">%s</pre>`

func renderer(md xlog.Markdown) template.HTML {
	htmlContent := fmt.Sprintf(pre, md)
	// #nosec G203 - Mermaid diagram syntax requires unescaped content. The content is rendered
	// by mermaid.js in the browser which sanitizes and prevents XSS. Escaping would break diagrams.
	return template.HTML(htmlContent + script)
}
