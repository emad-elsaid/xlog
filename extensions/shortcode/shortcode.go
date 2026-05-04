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

const (
	typeInfo    = "info"
	typeSuccess = "success"
	typeWarning = "warning"
	typeAlert   = "alert"
)

var (
	shortcodes = map[string]ShortCode{
		typeInfo:    {Render: func(c xlog.Markdown) template.HTML { return container("is-info", c) }},
		typeSuccess: {Render: func(c xlog.Markdown) template.HTML { return container("is-success", c) }},
		typeWarning: {Render: func(c xlog.Markdown) template.HTML { return container("is-warning", c) }},
		typeAlert:   {Render: func(c xlog.Markdown) template.HTML { return container("is-danger", c) }},
	}
	shortcodesMutex sync.RWMutex
)

// RegisterShortCode registers a new shortcode with the given name that can be
// used in markdown content via the {{name}} syntax.
func RegisterShortCode(name string, shortcode ShortCode) {
	shortcodesMutex.Lock()
	defer shortcodesMutex.Unlock()
	shortcodes[name] = shortcode
}
