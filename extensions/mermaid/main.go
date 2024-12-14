package mermaid

import (
	"fmt"
	"html/template"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
	shortcode "github.com/emad-elsaid/xlog/extensions/shortcode"
)

func init() {
	RegisterExtension(Mermaid{})
}

type Mermaid struct{}

func (Mermaid) Name() string { return "mermaid" }
func (Mermaid) Init() {
	shortcode.RegisterShortCode("mermaid", shortcode.ShortCode{Render: renderer})
}

//go:embed script.html
var script string

const pre = `<pre class="mermaid" style="background: transparent;text-align:center;">%s</pre>`

func renderer(md Markdown) template.HTML {
	html := fmt.Sprintf(pre, md)
	return template.HTML(html + script)

}
