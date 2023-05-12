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
	var rename PageRename
	var delete PageDelete

	RegisterCommand(commands)
	RegisterQuickCommand(commands)

	Post(`/\+/file/rename`, rename.Handler)
	Delete(`/\+/file/delete`, delete.Handler)
	RegisterTemplate(templates, "templates")
}

func commands(p Page) []Command {
	if READONLY {
		return []Command{}
	}

	return []Command{PageDelete{p}, PageRename{p}}
}

type PageRename struct {
	page Page
}

func (f PageRename) Icon() string {
	return "fa-solid fa-i-cursor"
}

func (f PageRename) Name() string {
	return "Rename"
}

func (f PageRename) OnClick() template.JS {
	return "renamePage(event)"
}

func (_ PageRename) Link() string { return "" }

func (f PageRename) Widget() template.HTML {
	if !f.page.Exists() {
		return ""
	}

	return Partial("file-operations-rename", Locals{
		"page":   f.page.Name(),
		"action": "/+/file/rename",
	})
}

func (f PageRename) Handler(w Response, r Request) Output {
	if READONLY {
		return Unauthorized("Readonly mode is active")
	}

	old := NewPage(r.FormValue("old"))
	if !old.Exists() {
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

func (f PageDelete) Icon() string {
	return "fa-solid fa-trash"
}

func (f PageDelete) Name() string {
	return "Delete"
}

func (_ PageDelete) Link() string { return "" }

func (f PageDelete) OnClick() template.JS {
	return "deletePage(event)"
}

func (f PageDelete) Widget() template.HTML {
	if !f.page.Exists() {
		return template.HTML("")
	}

	return Partial("file-operations-delete", Locals{
		"action": "/+/file/delete?page=" + url.QueryEscape(f.page.Name()),
	})
}

func (f PageDelete) Handler(w Response, r Request) Output {
	if READONLY {
		return Unauthorized("Readonly mode is active")
	}

	if page := NewPage(r.FormValue("page")); page.Exists() {
		page.Delete()
	}

	return NoContent()
}
