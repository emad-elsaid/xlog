package file_operations

import (
	"embed"
	"fmt"
	"html/template"
	"net/url"
	"os"
	"path"

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

	RegisterCommand(commands)
	RegisterQuickCommand(commands)
	RegisterTemplate(templates, "templates")
	Post(`/+/file/rename`, PageRename{}.Handler)
	Delete(`/+/file/delete`, PageDelete{}.Handler)
}

func commands(p Page) []Command {
	return []Command{PageDelete{p}, PageRename{p}}
}

type PageRename struct {
	page Page
}

func (PageRename) Icon() string         { return "fa-solid fa-i-cursor" }
func (PageRename) Name() string         { return "Rename" }
func (PageRename) OnClick() template.JS { return "renamePage(event)" }
func (PageRename) Link() string         { return "" }

func (f PageRename) Widget() template.HTML {
	if !f.page.Exists() {
		return ""
	}

	return Partial("file-operations-rename", Locals{
		"page":   f.page.Name(),
		"action": "/+/file/rename",
	})
}

func (f PageRename) Handler(r Request) Output {
	old := NewPage(r.FormValue("old"))
	if old == nil || !old.Exists() {
		return BadRequest("file doesn't exist")
	}

	ext := path.Ext(old.FileName())
	basename := r.FormValue("new")
	newName := basename + ext

	os.Rename(old.FileName(), newName)
	old.Write(Markdown(fmt.Sprintf("Renamed to: %s", basename)))

	return NoContent()
}

type PageDelete struct {
	page Page
}

func (PageDelete) Icon() string         { return "fa-solid fa-trash" }
func (PageDelete) Name() string         { return "Delete" }
func (PageDelete) Link() string         { return "" }
func (PageDelete) OnClick() template.JS { return "deletePage(event)" }

func (f PageDelete) Widget() template.HTML {
	if !f.page.Exists() {
		return template.HTML("")
	}

	return Partial("file-operations-delete", Locals{
		"action": "/+/file/delete?page=" + url.QueryEscape(f.page.Name()),
	})
}

func (f PageDelete) Handler(r Request) Output {
	if page := NewPage(r.FormValue("page")); page != nil && page.Exists() {
		page.Delete()
	}

	return NoContent()
}
