package file_operations

import (
	"embed"
	"fmt"
	"html/template"
	"net/url"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	var rename PageRename
	var delete PageDelete

	RegisterCommand(rename)
	RegisterCommand(delete)
	Post(`/\+/file/rename`, rename.Handler)
	Delete(`/\+/file/delete`, delete.Handler)
	RegisterTemplate(templates, "templates")
}

type PageRename int

func (f PageRename) Name() string {
	return "Rename Page"
}

func (f PageRename) OnClick() template.JS {
	return "renamePage()"
}

func (f PageRename) Widget(p Page) template.HTML {
	if !p.Exists() {
		return ""
	}

	return Partial("file-operations-rename", Locals{
		"page":   p.Name(),
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

	new := NewPage(r.FormValue("new"))
	new.Write(old.Content())

	old.Write(fmt.Sprintf("Renamed to: %s", new.Name()))
	return NoContent()
}

type PageDelete int

func (f PageDelete) Name() string {
	return "Delete Page"
}

func (f PageDelete) OnClick() template.JS {
	return "deletePage()"
}

func (f PageDelete) Widget(p Page) template.HTML {
	if !p.Exists() {
		return template.HTML("")
	}

	return Partial("file-operations-delete", Locals{
		"action": "/+/file/delete?page=" + url.QueryEscape(p.Name()),
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
