package xlog

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"os"
	"path"
)

// Define the catch all HTTP routes, parse CLI flags and take actions like
// building the static pages and exit, or start the HTTP server
func Start() {
	// Program Core routes. View, Edit routes and a route to write new content
	// to the page. + handling root path which just show `index` page.
	Get("/", rootHandler)
	Get("/edit/{page:.+}", getPageEditHandler)
	Get("/{page:.+}", getPageHandler)
	Post("/{page:.+}", postPageHandler)

	flag.Parse()

	if err := os.Chdir(SOURCE); err != nil {
		log.Fatal(err)
	}

	if len(BUILD) > 0 {
		READONLY = true

		if err := buildStaticSite(BUILD); err != nil {
			log.Printf("%s", err.Error())
		}
		os.Exit(0)
	}

	serve()
}

// Redirect to `/index` to render the index page.
func rootHandler(w Response, r Request) Output {
	return Redirect("/" + INDEX)
}

// Shows a page. the page name is the path itself. if the page doesn't exist it
// redirect to edit page otherwise will render it to HTML
func getPageHandler(w Response, r Request) Output {
	vars := Vars(r)
	page := NewPage(vars["page"])

	if !page.Exists() {
		if output, err := staticHandler(r); err == nil {
			return output
		}

		if READONLY {
			return NotFound("can't find page")
		}

		return Redirect("/edit/" + page.Name())
	}

	return Render("view", Locals{
		"title":   page.Emoji() + " " + page.Name(),
		"page":    page,
		"edit":    "/edit/" + page.Name(),
		"content": page.Render(),
		"csrf":    CSRF(r),
	})
}

// Edit page, gets the page from path
func getPageEditHandler(w Response, r Request) Output {
	if READONLY {
		return Unauthorized("Readonly mode is active")
	}

	vars := Vars(r)
	page := NewPage(vars["page"])

	var content Markdown
	if page.Exists() {
		content = page.Content()
	}

	return Render("edit", Locals{
		"title":        page.Emoji() + " " + page.Name(),
		"page":         page,
		"commands":     Commands(page),
		"content":      content,
		"autocomplete": autocompletes,
		"csrf":         CSRF(r),
	})
}

// Save new content of the page
func postPageHandler(w Response, r Request) Output {
	if READONLY {
		return Unauthorized("Readonly mode is active")
	}

	vars := Vars(r)
	page := NewPage(vars["page"])
	content := r.FormValue("content")
	page.Write(Markdown(content))

	return Redirect("/" + page.Name())
}

//go:embed public
var assets embed.FS

func staticHandler(r Request) (Output, error) {
	staticFSs := http.FS(priorityFS{
		assets,
		os.DirFS(SOURCE),
	})

	server := http.FileServer(staticFSs)

	cleanPath := path.Clean(r.URL.Path)

	if f, err := staticFSs.Open(cleanPath); err != nil {
		return nil, err
	} else {
		f.Close()
		return server.ServeHTTP, nil
	}
}
