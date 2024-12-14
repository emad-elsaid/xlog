package manifest

import (
	"embed"
	"html/template"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	RegisterExtension(Manifest{})
}

type Manifest struct{}

func (Manifest) Name() string { return "manifest" }
func (Manifest) Init() {
	Get("/manifest.json", manifest)
	RegisterBuildPage("/manifest.json", false)
	RegisterWidget(HEAD_WIDGET, 1, head)
	RegisterTemplate(templates, "templates")
}

func manifest(w Response, r Request) Output {
	return Cache(Render("manifest", Locals{"sitename": Config.Sitename}))
}

func head(_ Page) template.HTML {
	return template.HTML(`<link rel="manifest" href="/manifest.json">`)
}
