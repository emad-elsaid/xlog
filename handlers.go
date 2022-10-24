package xlog

import (
	"flag"
	"html/template"
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
	GET("/", rootHandler)
	GET("/edit/{page:.+}", getPageEditHandler)
	GET("/{page:.+}", getPageHandler)
	POST("/{page:.+}", postPageHandler)

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

		return Redirect("/edit/" + page.Name)
	}

	return Render("view", Locals{
		"edit":      "/edit/" + page.Name,
		"title":     page.Emoji() + " " + page.Name,
		"updated":   page.ModTime(),
		"content":   template.HTML(page.Render()),
		"tools":     RenderWidget(TOOLS_WIDGET, &page, r),      // all tools registered widgets
		"sidebar":   RenderWidget(SIDEBAR_WIDGET, &page, r),    // widgets registered for sidebar
		"action":    RenderWidget(ACTION_WIDGET, &page, r),     // widgets registered to be displayed under the page title
		"head":      RenderWidget(HEAD_WIDGET, &page, r),       // widgets registered to be displayed under the page title
		"afterView": RenderWidget(AFTER_VIEW_WIDGET, &page, r), // widgets registered to be displayed under the page content in the view page
	})
}

// Edit page, gets the page from path, if it doesn't exist it'll use the
// template.md content as default value
func getPageEditHandler(w Response, r Request) Output {
	if READONLY {
		return Unauthorized("Readonly mode is active")
	}

	vars := Vars(r)
	page := NewPage(vars["page"])

	var content string
	if page.Exists() {
		content = page.Content()
	} else if template := NewPage(TEMPLATE); template.Exists() {
		content = template.Content()
	}

	// Execute all Autocomplete functions and add them to a slice and pass it
	// down to the view
	acs := []*Autocomplete{}
	for _, v := range autocompletes {
		acs = append(acs, v())
	}

	return Render("edit", Locals{
		"title":        page.Name,
		"action":       page.Name,
		"tools":        RenderWidget(TOOLS_WIDGET, &page, r), // render all tools widgets
		"content":      content,
		"autocomplete": acs,
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
	page.Write(content)

	return Redirect("/" + page.Name)
}

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
