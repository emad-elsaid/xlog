package shortcode

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"

	"github.com/emad-elsaid/xlog"
)

type ShortCode struct {
	Render  func(xlog.Markdown) template.HTML
	Default string
}

func render(i xlog.Markdown) string {
	var b bytes.Buffer
	// Error ignored: Convert writes to buffer which doesn't fail in practice
	_ = xlog.MarkdownConverter().Convert([]byte(i), &b)
	return b.String()
}

func container(cls string, content xlog.Markdown) template.HTML {
	tpl := `<article class="message %s"><div class="message-body">%s</div></article>`
	// #nosec G203 -- cls is a predefined CSS class constant; content passes through MarkdownConverter which sanitizes
	return template.HTML(fmt.Sprintf(tpl, cls, render(content)))
}

var (
	shortcodes = map[string]ShortCode{
		"info":    {Render: func(c xlog.Markdown) template.HTML { return container("is-info", c) }},
		"success": {Render: func(c xlog.Markdown) template.HTML { return container("is-success", c) }},
		"warning": {Render: func(c xlog.Markdown) template.HTML { return container("is-warning", c) }},
		"alert":   {Render: func(c xlog.Markdown) template.HTML { return container("is-danger", c) }},
	}
	shortcodesMutex sync.RWMutex
)

func RegisterShortCode(name string, shortcode ShortCode) {
	shortcodesMutex.Lock()
	defer shortcodesMutex.Unlock()
	shortcodes[name] = shortcode
}
