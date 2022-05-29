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

	err = os.Chdir(absSource)
	if err != nil {
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

		html := page.Render()
		tools := template.HTML("")
		for _, v := range TOOLS_WIDGETS {
			tools += v(&page, r)
		}
		sidebar := template.HTML("")
		for _, v := range SIDEBAR_WIDGETS {
			sidebar += v(&page, r)
		}

		return Render("view", Locals{
			"edit":    "/" + page.Name + "/edit",
			"title":   page.Name,
			"updated": page.ModTime().Format("2006-01-02 15:04"),
			"content": template.HTML(html),
			"tools":   tools,
			"sidebar": sidebar,
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
		template := NewPage(TEMPLATE_NAME)

		var content string
		if page.Exists() {
			content = page.Content()
		} else if template.Exists() {
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
		for _, v := range NAVBAR_START_WIDGETS {
			o += v()
		}
		return o
	})

	Start()
}

// WIDGETS ===================================================

var (
	NAVBAR_START_WIDGETS = []func() template.HTML{}
	TOOLS_WIDGETS        = []func(*Page, Request) template.HTML{}
	SIDEBAR_WIDGETS      = []func(*Page, Request) template.HTML{}
)

func NAVBAR_START(f func() template.HTML)          { NAVBAR_START_WIDGETS = append(NAVBAR_START_WIDGETS, f) }
func TOOL(f func(*Page, Request) template.HTML)    { TOOLS_WIDGETS = append(TOOLS_WIDGETS, f) }
func SIDEBAR(f func(*Page, Request) template.HTML) { SIDEBAR_WIDGETS = append(SIDEBAR_WIDGETS, f) }
