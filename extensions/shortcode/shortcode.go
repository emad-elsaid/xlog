package shortcode

import (
	"bytes"
	"fmt"
	"html/template"

	. "github.com/emad-elsaid/xlog"
)

type ShortCodeFunc func(Markdown) template.HTML

func render(i Markdown) string {
	var b bytes.Buffer
	MarkDownRenderer.Convert([]byte(i), &b)
	return b.String()
}

func container(cls string, content Markdown) template.HTML {
	tpl := `<article class="message %s"><div class="message-body">%s</div></article>`
	return template.HTML(fmt.Sprintf(tpl, cls, render(content)))
}

var shortcodes = map[string]ShortCodeFunc{
	"info":    func(c Markdown) template.HTML { return container("is-info", c) },
	"success": func(c Markdown) template.HTML { return container("is-success", c) },
	"warning": func(c Markdown) template.HTML { return container("is-warning", c) },
	"alert":   func(c Markdown) template.HTML { return container("is-danger", c) },
}

func RegisterShortCode(name string, shortcode ShortCodeFunc) {
	shortcodes[name] = shortcode
}
