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

// Start initializes and starts the application
func (app *App) Start(ctx context.Context) {
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
	if len(app.config.Build) > 0 {
		app.config.Readonly = true
	}

	if !app.config.Readonly {
		app.Listen(PageChanged, app.clearPagesCache)
		app.Listen(PageDeleted, app.clearPagesCache)
	}

	if err := os.Chdir(app.config.Source); err != nil {
		slog.Error("Failed to change dir to source", "error", err, "source", app.config.Source)
		os.Exit(1)
	}

	app.initExtensions()

	app.Get("/{$}", app.rootHandler)
	app.Get("/{page...}", app.getPageHandler)

	if len(app.config.Build) > 0 {
		if err := app.build(app.config.Build); err != nil {
			slog.Error("Failed to build static pages", "error", err)
			os.Exit(1)
		}

		return
	}

	srv := app.server()
	slog.Info("Starting server", "address", app.config.BindAddress)

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	srv.ListenAndServe()
}

// rootHandler redirects to the index page
func (app *App) rootHandler(r Request) Output {
	return app.Redirect("/" + app.config.Index)
}

// getPageHandler handles page requests
func (app *App) getPageHandler(r Request) Output {
	page := app.NewPage(r.PathValue("page"))

	if page == nil {
		return app.NoContent()
	}

	if !page.Exists() {
		// if it's a directory get back to home page
		if s, err := os.Stat(page.Name()); err == nil && s.IsDir() {
			return app.Redirect(app.config.Index)
		}

		// if it's a static file serve it
		if output, err := app.staticHandler(r); err == nil {
			return output
		}

		// if it's readonly mode quit now
		if app.config.Readonly {
			return app.NotFound("can't find page")
		}

		// Allow extensions to handle this page if it's not readonly mode like
		// opening an editor or something
		app.Trigger(PageNotFound, page)

		page = DynamicPage{
			NameVal: page.Name(),
			RenderFn: func() template.HTML {
				str := "Page doesn't exist"
				return template.HTML(str)
			},
		}
	}

	return app.Render("page", Locals{
		"page": page,
		"csrf": csrf.Token(r),
	})
}
