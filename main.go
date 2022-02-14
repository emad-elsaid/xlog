package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

func main() {
	cwd, _ := os.Getwd()
	source := flag.String("source", cwd, "Directory that will act as a storage")
	flag.Parse()

	absSource, err := filepath.Abs(*source)
	if err != nil {
		log.Fatal(err)
	}

	os.Chdir(absSource)

	GET("/", func(w Response, r Request) Output {
		return Redirect("/index")
	})

	GET("/{page}", func(w Response, r Request) Output {
		vars := VARS(r)
		page := NewPage(vars["page"])

		if !page.Exists() {
			return Redirect("/" + page.Name() + "/edit")
		}

		return Render("layout", "view", Locals{
			"content": template.HTML(page.Render()),
		})
	})

	POST("/{page}", func(w Response, r Request) Output {
		vars := VARS(r)
		page := NewPage(vars["page"])
		content := r.FormValue("content")

		if content != "" {
			page.Write(content)
			return Redirect("/" + page.Name())
		} else if page.Exists() {
			page.Delete()
		}

		return Redirect("/")
	})

	GET("/{page}/edit", func(w Response, r Request) Output {
		vars := VARS(r)
		page := NewPage(vars["page"])

		return Render("edit_layout", "edit", Locals{
			"action":  page.Name(),
			"content": page.Content(),
			"csrf":    CSRF(r),
		})
	})

	Helpers()
	Start()
}
