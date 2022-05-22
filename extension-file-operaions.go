package main

import "html/template"

func init() {
	TOOL(fileOperationsWidget)
	DELETE("/+/file/delete/{page}", fileOperationsDeleteHandler)
}

func fileOperationsWidget(p *Page, r Request) template.HTML {
	return template.HTML(
		partial("extension/file-operations", Locals{
			"csrf":   CSRF(r),
			"page":   p.Name,
			"action": "/+/file/delete/" + p.Name,
		}),
	)
}

func fileOperationsDeleteHandler(w Response, r Request) Output {
	vars := VARS(r)
	page := NewPage(vars["page"])

	if page.Exists() {
		page.Delete()
	}

	return Redirect("/")
}
