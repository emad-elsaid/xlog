package xlog

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"runtime"
)

// Define the catch all HTTP routes, parse CLI flags and take actions like
// building the static pages and exit, or start the HTTP server
func Start(ctx context.Context) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	flag.Parse()

	// Program Core routes. View, Edit routes and a route to write new content
	// to the page. + handling root path which just show `index` page.
	Get("/{$}", rootHandler)
	Get("/edit/{page...}", getPageEditHandler)
	Get("/{page...}", getPageHandler)
	Post("/{page...}", postPageHandler)

	if err := os.Chdir(SOURCE); err != nil {
		slog.Error("Failed to change dir to source", "error", err, "source", SOURCE)
		os.Exit(1)
	}

	if len(BUILD) > 0 {
		READONLY = true

		if err := buildStaticSite(BUILD); err != nil {
			slog.Error("Failed to build static pages", "error", err)
			os.Exit(1)
		}

		return
	}

	srv := server()
	slog.Info("Starting server", "address", bindAddress)

	go func() {
		select {
		case <-ctx.Done():
			srv.Close()
			return
		}
	}()

	srv.ListenAndServe()
}

// Redirect to `/index` to render the index page.
func rootHandler(w Response, r Request) Output {
	return Redirect("/" + INDEX)
}

// Shows a page. the page name is the path itself. if the page doesn't exist it
// redirect to edit page otherwise will render it to HTML
func getPageHandler(w Response, r Request) Output {
	page := NewPage(r.PathValue("page"))

	if !page.Exists() {
		if s, err := os.Stat(page.Name()); err == nil && s.IsDir() {
			return Redirect(INDEX)
		}
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

	page := NewPage(r.PathValue("page"))

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

	page := NewPage(r.PathValue("page"))
	content := r.FormValue("content")
	page.Write(Markdown(content))

	return Redirect("/" + page.Name())
}
