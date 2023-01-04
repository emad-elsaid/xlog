package mermaid

import (
	"fmt"
	"html/template"

	. "github.com/emad-elsaid/xlog"
	shortcode "github.com/emad-elsaid/xlog/extensions/shortcode"
)

func init() {
	shortcode.ShortCode("mermaid", renderer)
}

const script = `
<script type="module">
  import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@9/dist/mermaid.esm.min.mjs';
  mermaid.initialize({
	startOnLoad: true,
    securityLevel: 'loose',
    theme: "neutral"
  });
</script>
`

func renderer(md Markdown) template.HTML {
	html := fmt.Sprintf(`<pre class="mermaid" style="background: transparent;text-align:center;">%s</pre>`, md)
	return template.HTML(html + script)

}
