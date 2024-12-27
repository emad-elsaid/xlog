package file_operations

import (
	"html/template"
	"log/slog"
	"net/url"

	. "github.com/emad-elsaid/xlog"
)

type PageDelete struct {
	page Page
}

func (PageDelete) Icon() string { return "fa-solid fa-trash" }
func (PageDelete) Name() string { return "Delete" }
func (f PageDelete) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"href":       "/+/file/delete?page=" + url.QueryEscape(f.page.Name()),
		"hx-delete":  "/+/file/delete?page=" + url.QueryEscape(f.page.Name()),
		"hx-confirm": "Are you sure?",
	}
}

func (f PageDelete) Handler(r Request) Output {
	name := r.FormValue("page")
	page := NewPage(name)
	if page == nil || !page.Exists() {
		slog.Error("Can't delete page", "page", page, "name", name)
	} else {
		page.Delete()
	}

	return func(w Response, r Request) {
		w.Header().Add("HX-Redirect", "/")
	}
}
