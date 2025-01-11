package xlog

import (
	"context"
	"flag"
	"html/template"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/gorilla/csrf"
	"gitlab.com/greyxor/slogor"
)

// Define the catch all HTTP routes, parse CLI flags and take actions like
// building the static pages and exit, or start the HTTP server
func Start(ctx context.Context) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	flag.Parse()

	// Setup logger
	level := slogor.SetLevel(slog.LevelDebug)
	timeFmt := slogor.SetTimeFormat(time.TimeOnly)
	handler := slogor.NewHandler(os.Stderr, level, timeFmt)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// if a static site is going to be built then lets also turn on read only
	// mode
	if len(Config.Build) > 0 {
		Config.Readonly = true
	}

	if !Config.Readonly {
		Listen(PageChanged, clearPagesCache)
		Listen(PageDeleted, clearPagesCache)
	}

	if err := os.Chdir(Config.Source); err != nil {
		slog.Error("Failed to change dir to source", "error", err, "source", Config.Source)
		os.Exit(1)
	}

	initExtensions()

	Get("/{$}", rootHandler)
	Get("/{page...}", getPageHandler)

	if len(Config.Build) > 0 {
		if err := build(Config.Build); err != nil {
			slog.Error("Failed to build static pages", "error", err)
			os.Exit(1)
		}

		return
	}

	srv := server()
	slog.Info("Starting server", "address", Config.BindAddress)

	go func() {
		<-ctx.Done()
		srv.Close()
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
