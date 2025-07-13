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
func (Manifest) Init(app *App) {
	app.Get("/manifest.json", manifest)
	app.RegisterBuildPage("/manifest.json", false)
	app.RegisterWidget(WidgetHead, 1, head)
	app.RegisterTemplate(templates, "templates")
}

func manifest(r Request) Output {
	app := GetApp()
	return app.Cache(app.Render("manifest", Locals{"sitename": app.GetConfig().Sitename}))
}

func head(Page) template.HTML {
	return template.HTML(`<link rel="manifest" href="/manifest.json">`)
}
