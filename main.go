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
			"content": content,
			"csrf":    CSRF(r),
		})
	})

	HELPER("navbarStart", func() template.HTML {
		o := template.HTML("")
		for _, v := range navbarWidgets {
			o += v()
		}
		return o
	})

	START()
}

// WIDGETS ===================================================

type widgetSpace int

const (
	TOOLS_WIDGET widgetSpace = iota
	SIDEBAR_WIDGET
	META_WIDGET
)

var (
	navbarWidgets = []func() template.HTML{}
	widgets       = map[widgetSpace][]func(*Page, Request) template.HTML{}
)

func NAVBAR_START(f func() template.HTML) { navbarWidgets = append(navbarWidgets, f) }

func WIDGET(s widgetSpace, f func(*Page, Request) template.HTML) {
	if _, ok := widgets[s]; !ok {
		widgets[s] = []func(*Page, Request) template.HTML{}
	}
	widgets[s] = append(widgets[s], f)
}

func renderWidget(s widgetSpace, p *Page, r Request) template.HTML {
	o := template.HTML("")
	ws, ok := widgets[s]
	if !ok {
		return o
	}

	for _, v := range ws {
		o += v(p, r)
	}
	return o
}
