package file_operations

import (
	"embed"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	RegisterExtension(FileOps{})
}

type FileOps struct{}

func (FileOps) Name() string { return "file-operations" }
func (FileOps) Init() {
	if Config.Readonly {
		return
	}

	RequireHTMX()
	RegisterCommand(commands)
	RegisterQuickCommand(commands)
	RegisterTemplate(templates, "templates")
	Post(`/+/file/rename`, PageRename{}.Handler)
	Get(`/+/file/rename`, PageRename{}.Form)
	Delete(`/+/file/delete`, PageDelete{}.Handler)
}

func commands(p Page) []Command {
	if len(p.FileName()) == 0 {
		return nil
	}

	return []Command{PageDelete{p}, PageRename{p}}
}
