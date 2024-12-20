package shortcode

import (
	"bytes"
	"fmt"
	"html/template"

	. "github.com/emad-elsaid/xlog"
)

type ShortCode struct {
	Render  func(Markdown) template.HTML
	Default string
}

func render(i Markdown) string {
	var b bytes.Buffer
	MarkdownConverter().Convert([]byte(i), &b)
	return b.String()
}

func container(cls string, content Markdown) template.HTML {
	tpl := `<article class="message %s"><div class="message-body">%s</div></article>`
	return template.HTML(fmt.Sprintf(tpl, cls, render(content)))
}

var shortcodes = map[string]ShortCode{
	"info":    {Render: func(c Markdown) template.HTML { return container("is-info", c) }},
	"success": {Render: func(c Markdown) template.HTML { return container("is-success", c) }},
	"warning": {Render: func(c Markdown) template.HTML { return container("is-warning", c) }},
	"alert":   {Render: func(c Markdown) template.HTML { return container("is-danger", c) }},
}

func RegisterShortCode(name string, shortcode ShortCode) {
	shortcodes[name] = shortcode
}
