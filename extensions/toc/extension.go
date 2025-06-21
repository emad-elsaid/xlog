package toc

import (
	"embed"
	"html/template"

	"github.com/emad-elsaid/xlog"
	gtoc "github.com/emad-elsaid/xlog/markdown-toc"
)

//go:embed templates
var templates embed.FS

func init() {
	xlog.RegisterExtension(Extension{})
}

type Extension struct{}

func (Extension) Name() string { return "toc" }
func (Extension) Init() {
	xlog.RegisterWidget(xlog.WidgetBeforeView, 0, TOC)
	xlog.RegisterTemplate(templates, "templates")
}

func TOC(p xlog.Page) template.HTML {
	if p == nil {
		return ""
	}

	src, doc := p.AST()
	if src == nil || doc == nil {
		return ""
	}

	tree, err := gtoc.Inspect(doc, src, gtoc.MaxDepth(1))
	if err != nil {
		return ""
	}

	if len(tree.Items) == 0 {
		return ""
	}

	return xlog.Partial("toc", xlog.Locals{"tree": tree})
}
