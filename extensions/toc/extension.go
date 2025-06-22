package toc

import (
	"embed"
	"html/template"

	"github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	xlog.RegisterExtension(Extension{})
}

type Extension struct{}

func (Extension) Name() string { return "toc" }
func (Extension) Init() {
	xlog.RegisterWidget(xlog.WidgetBeforeView, 0, tocWidget)
	xlog.RegisterTemplate(templates, "templates")
}

func tocWidget(p xlog.Page) template.HTML {
	if p == nil {
		return ""
	}

	src, doc := p.AST()
	if src == nil || doc == nil {
		return ""
	}

	tree, err := Inspect(doc, src, MaxDepth(1))
	if err != nil {
		return ""
	}

	if len(tree.Items) == 0 {
		return ""
	}

	return xlog.Partial("toc", xlog.Locals{"tree": tree})
}
