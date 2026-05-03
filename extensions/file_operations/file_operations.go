package file_operations

import (
	"embed"

	_ "embed"

	"github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	xlog.RegisterExtension(FileOps{})
}

type FileOps struct{}

func (FileOps) Name() string { return "file-operations" }
func (FileOps) Init() {
	if xlog.Config.Readonly {
		return
	}

	xlog.RequireHTMX()
	xlog.RegisterCommand(commands)
	xlog.RegisterQuickCommand(commands)
	xlog.RegisterTemplate(templates, "templates")
	xlog.Post(`/+/file/rename`, PageRename{}.Handler)
	xlog.Get(`/+/file/rename`, PageRename{}.Form)
	xlog.Delete(`/+/file/delete`, PageDelete{}.Handler)
}

func commands(p xlog.Page) []xlog.Command {
	if len(p.FileName()) == 0 {
		return nil
	}

	return []xlog.Command{PageDelete{p}, PageRename{p}}
}
