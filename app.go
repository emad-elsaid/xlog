package xlog

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"regexp"
	"runtime"
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

// RegisterPageSource registers a page source
func (app *App) RegisterPageSource(p PageSource) {
	app.sources = append([]PageSource{p}, app.sources...)
}

// GetConfig returns the application configuration
func (app *App) GetConfig() *Configuration {
	return app.config
}

// GetRouter returns the application router
func (app *App) GetRouter() *http.ServeMux {
	return app.router
}

// clearPagesCache clears the pages cache
func (app *App) clearPagesCache(p Page) {
	app.pages = nil
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
