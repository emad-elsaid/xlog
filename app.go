package xlog

import (
	"crypto/rand"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"path"
	"regexp"
	"runtime"
	"time"

	"github.com/gorilla/csrf"
)

// Property represent a piece of information about the current page such as last
// update time, number of versions, number of words, reading time...etc
type Property interface {
	// Icon returns the fontawesome icon class name or emoji
	Icon() string
	// Name returns the name of the property
	Name() string
	// Value returns the value of the property
	Value() any
}

// PageEvent represents different events that can occur with a page
type PageEvent int

// PageEventHandler is a function that handles a page event
type PageEventHandler func(Page) error

// List of page events
const (
	PageChanged PageEvent = iota
	PageDeleted
	PageNotFound // user requested a page that's not found
)

// Extension represents a plugin that can be registered with the application
type Extension interface {
	Name() string
	Init()
}

// Command defines a structure used for actions and links
type Command interface {
	// Icon returns the Fontawesome icon class name for the Command
	Icon() string
	// Name of the command. to be displayed in the list
	Name() string
	// Attrs a map of attributes to their values
	Attrs() map[template.HTMLAttr]any
}

// Preprocessor is a function that takes the whole page content and returns a
// modified version of the content
type Preprocessor func(Markdown) Markdown

// WidgetSpace used to represent a widgets spaces. it's used to register
// widgets to be injected in the view or edit pages
type WidgetSpace string

// WidgetFunc a function that takes the current page and returns the widget.
// This can be used by extensions to define new widgets to be rendered in
// view or edit pages. the extension should define this func type and
// register it to be rendered in a specific widgetSpace such as before or
// after the page.
type WidgetFunc func(Page) template.HTML

// List of widgets spaces that extensions can use to register a WidgetFunc to
// inject content into.
var (
	WidgetAfterView  WidgetSpace = "after_view"  // widgets rendered after the content of the view page
	WidgetBeforeView WidgetSpace = "before_view" // widgets rendered before the content of the view page
	WidgetHead       WidgetSpace = "head"        // widgets rendered in page <head> tag
)

//go:embed public
var assets embed.FS

// GetAssets returns the embedded assets filesystem
func GetAssets() embed.FS {
	return assets
}

// priorityFS returns file that exists in one of the FS structs.
// Prioritizing the end of the slice over earlier FSs.
type priorityFS []fs.FS

func (p priorityFS) Open(name string) (fs.File, error) {
	for i := len(p) - 1; i >= 0; i-- {
		cf := p[i]
		f, err := cf.Open(name)
		if err == nil {
			return f, err
		}
	}

	return nil, fs.ErrNotExist
}

//go:embed templates
var defaultTemplates embed.FS

// App represents the main application instance that encapsulates all global state
type App struct {
	// Configuration
	config *Configuration

	// HTTP server components
	router        *http.ServeMux
	csrfProtect   func(http.Handler) http.Handler
	requestLogger func(http.Handler) http.Handler

	// Templates and static files
	templates    *template.Template
	templatesFSs []fs.FS
	staticDirs   []fs.FS

	// Extensions and plugins
	extensions []Extension

	// Widgets
	widgets map[WidgetSpace]*priorityList[WidgetFunc]

	// Commands
	commands      []func(Page) []Command
	quickCommands []func(Page) []Command
	links         []func(Page) []Command

	// Events
	pageEvents map[PageEvent][]PageEventHandler

	// Pages and caching
	pages        []Page
	ignoredPaths []*regexp.Regexp
	concurrency  int

	// Properties
	propsSources []func(Page) []Property

	// Page sources
	sources []PageSource

	// Preprocessors
	preprocessors []Preprocessor

	// Helpers and JavaScript
	helpers template.FuncMap
	js      []string

	// Build-related
	extensionPage         map[string]bool
	extensionPageEnclosed map[string]bool
	buildPerms            fs.FileMode
}

// Global instance of the application
var globalApp *App

// GetApp returns the global application instance
func GetApp() *App {
	return globalApp
}

// init creates the global application instance
func init() {
	globalApp = &App{
		config:                &Config,
		router:                http.NewServeMux(),
		widgets:               make(map[WidgetSpace]*priorityList[WidgetFunc]),
		pageEvents:            make(map[PageEvent][]PageEventHandler),
		ignoredPaths:          []*regexp.Regexp{regexp.MustCompile(`^\.`)}, // Ignore any hidden directory
		concurrency:           runtime.NumCPU() * 4,
		propsSources:          []func(Page) []Property{DefaultProps},
		sources:               []PageSource{newMarkdownFS(".")},
		preprocessors:         []Preprocessor{},
		helpers:               template.FuncMap{},
		js:                    []string{},
		extensionPage:         make(map[string]bool),
		extensionPageEnclosed: make(map[string]bool),
		buildPerms:            0744,
		staticDirs:            []fs.FS{GetAssets()},
	}

	// Initialize default helpers
	globalApp.initDefaultHelpers()
}

// initDefaultHelpers initializes the default helper functions
func (app *App) initDefaultHelpers() {
	app.helpers = template.FuncMap{
		"ago":            app.ago,
		"properties":     app.Properties,
		"links":          app.Links,
		"widgets":        app.RenderWidget,
		"commands":       app.Commands,
		"quick_commands": app.QuickCommands,
		"isFontAwesome":  app.IsFontAwesome,
		"includeJS":      app.includeJS,
		"scripts":        app.scripts,
		"banner":         app.Banner,
		"emoji":          app.Emoji,
		"base":           path.Base,
		"dir":            app.dir,
		"raw":            app.raw,
	}
}

// server creates and configures the HTTP server
func (app *App) server() *http.Server {
	app.compileTemplates()
	var handler http.Handler = app.router
	for _, v := range app.defaultMiddlewares() {
		handler = v(handler)
	}

	return &http.Server{
		Handler:      handler,
		Addr:         app.config.BindAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

// defaultMiddlewares returns the default middleware stack
func (app *App) defaultMiddlewares() []func(http.Handler) http.Handler {
	var middlewares []func(http.Handler) http.Handler

	if !app.config.Readonly {
		crsfOpts := []csrf.Option{
			csrf.Path("/"),
			csrf.FieldName("csrf"),
			csrf.CookieName(app.config.CsrfCookieName),
			csrf.Secure(!app.config.ServeInsecure),
		}

		sessionSecret := []byte(os.Getenv("SESSION_SECRET"))
		if len(sessionSecret) == 0 {
			sessionSecret = make([]byte, 128)
			rand.Read(sessionSecret)
		}

		app.csrfProtect = csrf.Protect(sessionSecret, crsfOpts...)
		middlewares = append(middlewares, app.csrfProtect)
	}

	middlewares = append(middlewares, app.requestLoggerHandler)

	return middlewares
}

// Get registers a GET route
func (app *App) Get(path string, handler HandlerFunc) {
	app.router.HandleFunc("GET "+path, app.handlerFuncToHttpHandler(handler))
}

// Post registers a POST route
func (app *App) Post(path string, handler HandlerFunc) {
	app.router.HandleFunc("POST "+path, app.handlerFuncToHttpHandler(handler))
}

// Delete registers a DELETE route
func (app *App) Delete(path string, handler HandlerFunc) {
	app.router.HandleFunc("DELETE "+path, app.handlerFuncToHttpHandler(handler))
}

// HandlerFunc is the type of an HTTP handler function + returns output function.
// Request is an alias of *http.Request for shorter handler declaration
// Response is an alias http.ResponseWriter for shorter handler declaration
// Output is an alias of http.HandlerFunc as output is expected from defined http handlers
// Locals is a map of string to any value used for template rendering

func (app *App) handlerFuncToHttpHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w Response, r Request) {
		handler(r)(w, r)
	}
}

// RegisterPageSource registers a page source
func (app *App) RegisterPageSource(p PageSource) {
	app.sources = append([]PageSource{p}, app.sources...)
}

// RequireHTMX registers HTMX library
func (app *App) RequireHTMX() {
	app.includeJS("/public/htmx.min.js")
}

// GetConfig returns the application configuration
func (app *App) GetConfig() *Configuration {
	return app.config
}

// GetRouter returns the application router
func (app *App) GetRouter() *http.ServeMux {
	return app.router
}

// requestLoggerHandler logs HTTP requests
func (app *App) requestLoggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w Response, r Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		slog.Info(r.Method+" "+r.URL.Path, "time", time.Since(start))
	})
}

// clearPagesCache clears the pages cache
func (app *App) clearPagesCache(p Page) {
	app.pages = nil
}

// Redirect returns a redirect response
func (app *App) Redirect(url string) Output {
	return func(w Response, r Request) {
		http.Redirect(w, r, url, http.StatusFound)
	}
}

// NoContent returns a no content response
func (app *App) NoContent() Output {
	return func(w Response, r Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

// NotFound returns a not found response
func (app *App) NotFound(msg string) Output {
	return func(w Response, r Request) {
		http.Error(w, msg, http.StatusNotFound)
	}
}

// Render returns a rendered template response
func (app *App) Render(path string, data Locals) Output {
	return func(w Response, r Request) {
		fmt.Fprint(w, app.Partial(path, data))
	}
}

// BadRequest returns a bad request response
func (app *App) BadRequest(msg string) Output {
	return func(w Response, r Request) {
		http.Error(w, msg, http.StatusBadRequest)
	}
}

// InternalServerError returns an internal server error response
func (app *App) InternalServerError(err error) Output {
	return func(w Response, r Request) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// PlainText returns a plain text response
func (app *App) PlainText(text string) Output {
	return func(w Response, r Request) {
		w.Write([]byte(text))
	}
}

// JsonResponse returns a JSON response
func (app *App) JsonResponse(a any) Output {
	return func(w Response, r Request) {
		b, err := json.Marshal(a)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(b)
	}
}

// Cache wraps Output and adds header to instruct the browser to cache the output
func (app *App) Cache(out Output) Output {
	return func(w Response, r Request) {
		w.Header().Add("Cache-Control", "max-age=604800")
		out(w, r)
	}
}

// NewPage creates a new page
func (app *App) NewPage(name string) Page {

	for i := range app.sources {
		p := app.sources[i].Page(name)
		if p != nil && p.Exists() {
			return p
		}
	}

	return &page{name: name}
}
