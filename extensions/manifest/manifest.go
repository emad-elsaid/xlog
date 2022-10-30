package manifest

import (
	"embed"
	"html/template"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	Get("/manifest.json", manifest)
	BuildPage("/manifest.json", false)
	Widget(HEAD_WIDGET, head)
	Template(templates, "templates")
}

func manifest(w Response, r Request) Output {
	return Render("manifest", Locals{"sitename": SITENAME})
}

func head(_ Page, _ Request) template.HTML {
	return template.HTML(`<link rel="manifest" href="/manifest.json">`)
}
