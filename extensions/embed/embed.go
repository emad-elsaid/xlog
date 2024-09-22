package embed

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/extensions/shortcode"
)

func init() {
	shortcode.RegisterShortCode("embed", embedShortcode)
}

func embedShortcode(in xlog.Markdown) template.HTML {
	p := xlog.NewPage(strings.TrimSpace(string(in)))
	if p == nil || !p.Exists() {
		return template.HTML(fmt.Sprintf("Page: %s doesn't exist", in))
	}

	return p.Render()
}
