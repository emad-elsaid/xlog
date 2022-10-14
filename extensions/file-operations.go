package extensions

import (
	"fmt"
	"html/template"
	"net/url"

	. "github.com/emad-elsaid/xlog"
)

func init() {
	WIDGET(TOOLS_WIDGET, fileOperationsDeleteWidget)
	WIDGET(TOOLS_WIDGET, fileOperationsRenameWidget)
	DELETE(`/\+/file/delete`, fileOperationsDeleteHandler)
	POST(`/\+/file/rename`, fileOperationsRenameHandler)
}

func fileOperationsDeleteWidget(p *Page, r Request) template.HTML {
	if !p.Exists() {
		return template.HTML("")
	}

	return template.HTML(
		Partial("extension/file-operations-delete", Locals{
			"csrf":   CSRF(r),
			"page":   p.Name,
			"action": "/+/file/delete?page=" + url.QueryEscape(p.Name),
		}),
	)
}

func fileOperationsRenameWidget(p *Page, r Request) template.HTML {
	if !p.Exists() {
		return template.HTML("")
	}

	return template.HTML(
		Partial("extension/file-operations-rename", Locals{
			"csrf":   CSRF(r),
			"page":   p.Name,
			"action": "/+/file/rename",
		}),
	)
}

func fileOperationsDeleteHandler(w Response, r Request) Output {
	if page := NewPage(r.FormValue("page")); page.Exists() {
		page.Delete()
	}

	return Redirect("/")
}

func fileOperationsRenameHandler(w Response, r Request) Output {
	old := NewPage(r.FormValue("old"))
	if !old.Exists() {
		return BadRequest
	}

	new := NewPage(r.FormValue("new"))
	new.Write(old.Content())

	old.Write(fmt.Sprintf("Renamed to: %s", new.Name))
	return NoContent
}
