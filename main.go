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

	GET("/", func(w Response, r Request) Output {
		return Redirect("/index")
	})

	GET("/{page}", func(w Response, r Request) Output {
		vars := VARS(r)
		page := NewPage(vars["page"])

		if !page.Exists() {
			return Redirect("/" + page.Name + "/edit")
		}

		return Render("view", Locals{
			"edit":    "/" + page.Name + "/edit",
			"title":   page.Name,
			"updated": page.ModTime().Format("2006-01-02 15:04"),
			"content": template.HTML(page.Render()),
			"navbar":  renderWidget(NAVBAR_WIDGET, &page, r),
			"tools":   renderWidget(TOOLS_WIDGET, &page, r),
			"sidebar": renderWidget(SIDEBAR_WIDGET, &page, r),
			"meta":    renderWidget(META_WIDGET, &page, r),
		})
	})

	POST("/{page}", func(w Response, r Request) Output {
		vars := VARS(r)
		page := NewPage(vars["page"])
		content := r.FormValue("content")

		page.Write(content)
		return Redirect("/" + page.Name)
	})

	GET("/{page}/edit", func(w Response, r Request) Output {
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
			"content": content,
			"csrf":    CSRF(r),
		})
	})

	START()
}

// WIDGETS ===================================================

type widgetSpace int
type widgetFunc func(*Page, Request) template.HTML

const (
	TOOLS_WIDGET widgetSpace = iota
	SIDEBAR_WIDGET
	META_WIDGET
	NAVBAR_WIDGET
)

var widgets = map[widgetSpace][]widgetFunc{}

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
