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
	setupLogger()

	if !Config.Readonly {
		Listen(PageChanged, clearPagesCache)
		Listen(PageChanged, clearPagesCache)
		Listen(PageDeleted, clearPagesCache)
	}

	initExtensions()

	Get("/{$}", rootHandler)
	Get("/{page...}", getPageHandler)

	if err := os.Chdir(Config.Source); err != nil {
		slog.Error("Failed to change dir to source", "error", err, "source", Config.Source)
		os.Exit(1)
	}

	if len(Config.Build) > 0 {
		Config.Readonly = true

		if err := buildStaticSite(Config.Build); err != nil {
			slog.Error("Failed to build static pages", "error", err)
			os.Exit(1)
		}

		return
	}

	srv := server()
	slog.Info("Starting server", "address", Config.BindAddress)

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
		if s, err := os.Stat(page.Name()); err == nil && s.IsDir() {
			return Redirect(Config.Index)
		}
		if output, err := staticHandler(r); err == nil {
			return output
		}

		if Config.Readonly {
			return NotFound("can't find page")
		}

		// Allow extensions to handle the event
		Trigger(PageNotFound, page)

		return NotFound("Page does not exist")
	}

	return Render("view", Locals{
		"title":   page.Emoji() + " " + page.Name(),
		"page":    page,
		"content": page.Render(),
		"csrf":    CSRF(r),
	})
}
