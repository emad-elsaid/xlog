package mermaid

import (
	"embed"
	"html/template"

	. "github.com/emad-elsaid/xlog"
)

//go:embed script.html
var script string

// Prevent unused import removal
var _ = embed.FS{}

func init() {
	app := GetApp()
	app.RegisterExtension(Mermaid{})
}

type Mermaid struct{}

func (Mermaid) Name() string { return "mermaid" }
func (Mermaid) Init() {
	app := GetApp()
	app.RegisterWidget(WidgetHead, 0, mermaidWidget)
}

func mermaidWidget(Page) template.HTML {
	return template.HTML(script)
}
