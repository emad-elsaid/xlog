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

//go:embed public
var assets embed.FS

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
		staticDirs:            []fs.FS{assets},
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
