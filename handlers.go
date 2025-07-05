package xlog

import (
	"context"
	"html/template"
	"os"

	"github.com/gorilla/csrf"
)

// Define the catch all HTTP routes, parse CLI flags and take actions like
// building the static pages and exit, or start the HTTP server
func Start(ctx context.Context) {
	app := GetApp()
	app.Start(ctx)
}

// Redirect to `/index` to render the index page.
func rootHandler(r Request) Output {
	return Redirect("/" + Config.Index)
}

// Shows a page. the page name is the path itself. if the page doesn't exist it
// redirect to edit page otherwise will render it to HTML
func getPageHandler(r Request) Output {
	page := NewPage(r.PathValue("page"))

	if page == nil {
		return NoContent()
	}

	if !page.Exists() {
		// if it's a directory get back to home page
		if s, err := os.Stat(page.Name()); err == nil && s.IsDir() {
			return Redirect(Config.Index)
		}

		// if it's a static file serve it
		if output, err := staticHandler(r); err == nil {
			return output
		}

		// if it's readonly mode quit now
		if Config.Readonly {
			return NotFound("can't find page")
		}

		// Allow extensions to handle this page if it's not readonly mode like
		// opening an editor or something
		Trigger(PageNotFound, page)

		page = DynamicPage{
			NameVal: page.Name(),
			RenderFn: func() template.HTML {
				str := "Page doesn't exist"

				return template.HTML(str)
			},
		}
	}

	return Render("page", Locals{
		"page": page,
		"csrf": csrf.Token(r),
	})
}
