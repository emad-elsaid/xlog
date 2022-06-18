package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

const TEMPLATE_NAME = "template"

func main() {
	cwd, _ := os.Getwd()
	source := flag.String("source", cwd, "Directory that will act as a storage")
	flag.Parse()

	absSource, err := filepath.Abs(*source)
	if err != nil {
		log.Fatal(err)
	}

	if err = os.Chdir(absSource); err != nil {
		log.Fatal(err)
	}

	GET("/", RootHandler)
	GET("/edit/{page:.*}", GetPageEditHandler)
	GET("/{page:.*}", GetPageHandler)
	POST("/{page:.*}", PostPageHandler)

	START()
}

func RootHandler(w Response, r Request) Output {
	return Redirect("/index")
}

func GetPageHandler(w Response, r Request) Output {
	vars := VARS(r)
	page := NewPage(vars["page"])

	if !page.Exists() {
		return Redirect("/edit/" + page.Name)
	}

	return Render("view", Locals{
		"edit":      "/edit/" + page.Name,
		"title":     page.Name,
		"updated":   page.ModTime().Format("2006-01-02 15:04"),
		"content":   template.HTML(page.Render()),
		"tools":     renderWidget(TOOLS_WIDGET, &page, r),
		"sidebar":   renderWidget(SIDEBAR_WIDGET, &page, r),
		"meta":      renderWidget(META_WIDGET, &page, r),
		"afterView": renderWidget(AFTER_VIEW_WIDGET, &page, r),
	})
}

func GetPageEditHandler(w Response, r Request) Output {
	vars := VARS(r)
	page := NewPage(vars["page"])

	var content string
	if page.Exists() {
		content = page.Content()
	} else if template := NewPage(TEMPLATE_NAME); template.Exists() {
		content = template.Content()
	}

	return Render("edit", Locals{
		"title":   page.Name,
		"action":  page.Name,
		"rtl":     page.RTL(),
		"tools":   renderWidget(TOOLS_WIDGET, &page, r),
		"content": content,
		"csrf":    CSRF(r),
	})
}

func PostPageHandler(w Response, r Request) Output {
	vars := VARS(r)
	page := NewPage(vars["page"])
	content := r.FormValue("content")
	page.Write(content)

	return Redirect("/" + page.Name)
}

// WIDGETS ===================================================

type (
	widgetSpace int
	widgetFunc  func(*Page, Request) template.HTML
)

const (
	TOOLS_WIDGET widgetSpace = iota
	SIDEBAR_WIDGET
	AFTER_VIEW_WIDGET
	META_WIDGET
)

var widgets = map[widgetSpace][]widgetFunc{}

func PREPEND_WIDGET(s widgetSpace, f func(*Page, Request) template.HTML) {
	if _, ok := widgets[s]; !ok {
		widgets[s] = []widgetFunc{}
	}
	widgets[s] = append([]widgetFunc{f}, widgets[s]...)
}

func WIDGET(s widgetSpace, f func(*Page, Request) template.HTML) {
	if _, ok := widgets[s]; !ok {
		widgets[s] = []widgetFunc{}
	}
	widgets[s] = append(widgets[s], f)
}

func renderWidget(s widgetSpace, p *Page, r Request) (o template.HTML) {
	for _, v := range widgets[s] {
		o += v(p, r)
	}
	return
}
