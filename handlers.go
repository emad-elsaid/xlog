package xlog

import (
	"context"
	"flag"
	"html/template"
	"log/slog"
	"os"
	"runtime"
	"time"

	csrf "filippo.io/csrf/gorilla"
	"gitlab.com/greyxor/slogor"
)

const keyPage = "page"

// osExit is a variable that allows testing of os.Exit calls.
var osExit = os.Exit

// printVersion outputs the version information and runtime details.
func printVersion() {
	slog.Info("xlog version", "version", Version, "go", runtime.Version(), "os", runtime.GOOS, "arch", runtime.GOARCH)
}

// listPages outputs all markdown pages, one per line.
// Format: filename (without .md extension) for easy piping to other commands.
func listPages() {
	allPages := Pages(context.Background())
	for _, p := range allPages {
		slog.Info(p.Name())
	}
}

// Define the catch all HTTP routes, parse CLI flags and take actions like
// building the static pages and exit, or start the HTTP server.
func Start(ctx context.Context) {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	// Only parse flags if not already parsed (avoids conflict with testing framework)
	if !flag.Parsed() {
		flag.Parse()
	}

	// Handle version flag first - most lightweight operation
	if Config.ShowVersion {
		printVersion()
		osExit(0)
	}

	// Handle completion flag early before any logging or extension initialization
	if Config.Completion != "" {
		handleCompletion(Config.Completion)
		osExit(0)
	}

	// Handle doctor flag early before extension initialization but after basic setup
	if Config.RunDoctor {
		// Setup minimal logger for doctor output
		level := slogor.SetLevel(slog.LevelInfo)
		timeFmt := slogor.SetTimeFormat(time.TimeOnly)
		handler := slogor.NewHandler(os.Stderr, level, timeFmt)
		logger := slog.New(handler)
		slog.SetDefault(logger)

		Doctor()
		osExit(0)
	}

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
		osExit(1)
	}

	// Handle list flag after chdir to source directory
	if Config.ListPages {
		// Setup minimal logger for list output
		level := slogor.SetLevel(slog.LevelInfo)
		timeFmt := slogor.SetTimeFormat(time.TimeOnly)
		handler := slogor.NewHandler(os.Stderr, level, timeFmt)
		logger := slog.New(handler)
		slog.SetDefault(logger)

		listPages()
		osExit(0)
	}

	// Handle stats flag after chdir to source directory
	if Config.ShowStats {
		Stats(ctx)
		osExit(0)
	}

	initExtensions()

	// Register health check endpoint for load balancers and monitoring
	router.HandleFunc("GET /+/health", healthHandler)

	Get("/{$}", rootHandler)
	Get("/{page...}", getPageHandler)

	if len(Config.Build) > 0 {
		if err := build(Config.Build); err != nil {
			slog.Error("Failed to build static pages", "error", err)
			osExit(1)
		}

		return
	}

	srv := server()
	slog.Info("Starting server", "address", Config.BindAddress)

	go func() {
		<-ctx.Done()
		if err := srv.Close(); err != nil {
			slog.Error("Failed to close server", "error", err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("Server error", "error", err)
	}
}

// Redirect to `/index` to render the index page.
func rootHandler(r Request) Output {
	return Redirect("/" + Config.Index)
}

// Shows a page. the page name is the path itself. if the page doesn't exist it
// redirect to edit page otherwise will render it to HTML.
func getPageHandler(r Request) Output {
	page := NewPage(r.PathValue("page"))

	if page == nil {
		return NoContent()
	}

	if !page.Exists() {
		// if it's a directory get back to home page
		// #nosec G703 -- page.Name() is sanitized by PageSource abstraction, no path traversal risk
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

				// #nosec G203 - str is a constant string literal, safe to convert
				return template.HTML(str)
			},
		}
	}

	return Render("page", Locals{
		keyPage: page,
		"csrf":  csrf.Token(r), //nolint:staticcheck // Deprecated but needed for template compatibility
	})
}
