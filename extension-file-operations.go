package main

import (
	"html/template"
	"net/url"
)

func init() {
	WIDGET(TOOLS_WIDGET, fileOperationsWidget)
	DELETE(`/\+/file/delete`, fileOperationsDeleteHandler)
}

func fileOperationsWidget(p *Page, r Request) template.HTML {
	return template.HTML(
		partial("extension/file-operations", Locals{
			"csrf":   CSRF(r),
			"page":   p.Name,
			"action": "/+/file/delete?page=" + url.QueryEscape(p.Name),
		}),
	)
}

func fileOperationsDeleteHandler(w Response, r Request) Output {
	if page := NewPage(r.FormValue("page")); page.Exists() {
		page.Delete()
	}

	return Redirect("/")
}
