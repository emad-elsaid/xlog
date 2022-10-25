package manifest

import (
	"embed"
	"html/template"
	"io/fs"

	. "github.com/emad-elsaid/xlog"
)

//go:embed views
var views embed.FS

func init() {
	f, _ := fs.Sub(views, "views")
	Template(f)
	Get("/manifest.json", manifest)
	ExtensionPage("/manifest.json")
	Widget(HEAD_WIDGET, head)
}

func manifest(w Response, r Request) Output {
	return Render("manifest", Locals{"sitename": SITENAME})
}

func head(_ *Page, _ Request) template.HTML {
	return template.HTML(`<link rel="manifest" href="/manifest.json">`)
}
