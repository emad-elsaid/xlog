package xlog

import (
	"flag"
	"html/template"
	"log"
	"os"
	"time"
)

// Define the catch all HTTP routes, parse CLI flags and take actions like
// building the static pages and exit, or start the HTTP server
func Start() {
	// Program Core routes. View, Edit routes and a route to write new content
	// to the page. + handling root path which just show `index` page.
	GET("/", RootHandler)
	GET("/edit/{page:.*}", GetPageEditHandler)
	GET("/{page:.*}", GetPageHandler)
	POST("/{page:.*}", PostPageHandler)

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
func RootHandler(w Response, r Request) Output {
	return Redirect("/" + INDEX)
}

// Shows a page. the page name is the path itself. if the page doesn't exist it
// redirect to edit page otherwise will render it to HTML
func GetPageHandler(w Response, r Request) Output {
	vars := VARS(r)
	page := NewPage(vars["page"])

	if !page.Exists() {
		if READONLY {
			return NotFound
		} else {
			return Redirect("/edit/" + page.Name)
		}
	}

	return Render("view", Locals{
		"edit":      "/edit/" + page.Name,
		"title":     page.Emoji() + " " + page.Name,
		"updated":   ago(time.Now().Sub(page.ModTime())),
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
func GetPageEditHandler(w Response, r Request) Output {
	if READONLY {
		return NotFound
	}

	vars := VARS(r)
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
		"rtl":          page.RTL(),                           // is it Right-To-Left page?
		"tools":        RenderWidget(TOOLS_WIDGET, &page, r), // render all tools widgets
		"content":      content,
		"autocomplete": acs,
		"csrf":         CSRF(r),
	})
}

// Save new content of the page
func PostPageHandler(w Response, r Request) Output {
	if READONLY {
		return NotFound
	}

	vars := VARS(r)
	page := NewPage(vars["page"])
	content := r.FormValue("content")
	page.Write(content)

	return Redirect("/" + page.Name)
}
