package manifest

import (
	"embed"
	"html/template"

	"github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	xlog.RegisterExtension(Manifest{})
}

type Manifest struct{}

func (Manifest) Name() string { return "manifest" }
func (Manifest) Init() {
	xlog.Get("/manifest.json", manifest)
	xlog.RegisterBuildPage("/manifest.json", false)
	xlog.RegisterWidget(xlog.WidgetHead, 1, head)
	xlog.RegisterTemplate(templates, "templates")
}

func manifest(r xlog.Request) xlog.Output {
	return xlog.Cache(xlog.Render("manifest", xlog.Locals{"sitename": xlog.Config.Sitename}))
}

func head(xlog.Page) template.HTML {
	return template.HTML(`<link rel="manifest" href="/manifest.json">`)
}
