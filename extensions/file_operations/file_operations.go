package file_operations

import (
	"embed"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	app := GetApp()
	app.RegisterExtension(FileOps{})
}

type FileOps struct{}

func (FileOps) Name() string { return "file-operations" }
func (FileOps) Init(app *App) {
	if app.GetConfig().Readonly {
		return
	}

	app.RequireHTMX()
	app.RegisterCommand(commands)
	app.RegisterQuickCommand(commands)
	app.RegisterTemplate(templates, "templates")
	app.Post(`/+/file/rename`, PageRename{}.Handler)
	app.Get(`/+/file/rename`, PageRename{}.Form)
	app.Delete(`/+/file/delete`, PageDelete{}.Handler)
}

func commands(p Page) []Command {
	if len(p.FileName()) == 0 {
		return nil
	}

	return []Command{PageDelete{p}, PageRename{p}}
}
