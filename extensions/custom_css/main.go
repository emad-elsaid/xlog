package custom_css

import (
	"flag"
	"html/template"

	. "github.com/emad-elsaid/xlog"
)

var custom_css_file string

func init() {
	flag.StringVar(&custom_css_file, "custom_css", "", "Custom CSS file path")
	RegisterExtension(CustomCSS{})
}

type CustomCSS struct{}

func (CustomCSS) Name() string { return "custom-css" }
func (CustomCSS) Init() {
	if custom_css_file != "" {
		RegisterWidget(HEAD_WIDGET, 1, custom_css_tag)
	}
}

func custom_css_tag(_ Page) template.HTML {
	return template.HTML(`<link rel="stylesheet" href="` + custom_css_file + `">`)
}
