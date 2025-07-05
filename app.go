package xlog

import (
	"bytes"
	"context"
	"crypto/rand"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"

	"iter"
	"sort"

	"github.com/emad-elsaid/xlog/markdown/ast"
	emojiAst "github.com/emad-elsaid/xlog/markdown/emoji/ast"
	"github.com/gorilla/csrf"
	"gitlab.com/greyxor/slogor"
	"golang.org/x/sync/errgroup"
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
		app.Listen(PageChanged, func(p Page) error { app.clearPagesCache(p); return nil })
		app.Listen(PageDeleted, func(p Page) error { app.clearPagesCache(p); return nil })
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
		srv.Close()
	}()

	srv.ListenAndServe()
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

// RegisterExtension registers a new extension
func (app *App) RegisterExtension(e Extension) {
	app.extensions = append(app.extensions, e)
}

// RegisterWidget registers a widget function
func (app *App) RegisterWidget(s WidgetSpace, priority float32, f WidgetFunc) {

	pl, ok := app.widgets[s]
	if !ok {
		pl = new(priorityList[WidgetFunc])
		app.widgets[s] = pl
	}

	pl.Add(f, priority)
}

// RegisterCommand registers a new command
func (app *App) RegisterCommand(c func(Page) []Command) {
	app.commands = append(app.commands, c)
}

// RegisterQuickCommand registers a new quick command
func (app *App) RegisterQuickCommand(c func(Page) []Command) {
	app.quickCommands = append(app.quickCommands, c)
}

// RegisterLink registers a new link
func (app *App) RegisterLink(l func(Page) []Command) {
	app.links = append(app.links, l)
}

// RegisterPageSource registers a page source
func (app *App) RegisterPageSource(p PageSource) {
	app.sources = append([]PageSource{p}, app.sources...)
}

// RegisterPreprocessor registers a preprocessor function
func (app *App) RegisterPreprocessor(f Preprocessor) {
	app.preprocessors = append(app.preprocessors, f)
}

// RegisterProperty registers a function that returns a set of properties for a page
func (app *App) RegisterProperty(a func(Page) []Property) {
	app.propsSources = append(app.propsSources, a)
}

// RequireHTMX registers HTMX library
func (app *App) RequireHTMX() {
	app.includeJS("/public/htmx.min.js")
}

// PreProcess processes content through all registered preprocessors
func (app *App) PreProcess(content Markdown) Markdown {

	for _, v := range app.preprocessors {
		content = v(content)
	}

	return content
}

// initExtensions initializes all registered extensions
func (app *App) initExtensions() {
	if app.config.DisabledExtensions == "all" {
		slog.Info("extensions", "disabled", "all")
		return
	}

	disabled := strings.Split(app.config.DisabledExtensions, ",")
	disabledNames := []string{} // because the user can input wrong extension name
	enabledNames := []string{}
	for i := range app.extensions {
		if slices.Contains(disabled, app.extensions[i].Name()) {
			disabledNames = append(disabledNames, app.extensions[i].Name())
			continue
		}

		app.extensions[i].Init()
		enabledNames = append(enabledNames, app.extensions[i].Name())
	}

	slog.Info("extensions", "enabled", enabledNames, "disabled", disabled)
}

// RegisterTemplate registers a filesystem that contains templates
func (app *App) RegisterTemplate(t fs.FS, subDir string) {

	ts, _ := fs.Sub(t, subDir)
	app.templatesFSs = append(app.templatesFSs, ts)
}

// RegisterStaticDir adds a filesystem to the static files list
func (app *App) RegisterStaticDir(f fs.FS) {
	app.staticDirs = append(app.staticDirs, f)
}

// RegisterHelper registers a new helper function
func (app *App) RegisterHelper(name string, f any) error {

	if _, ok := app.helpers[name]; ok {
		return ErrHelperRegistered
	}

	app.helpers[name] = f
	return nil
}

// IgnorePath registers a pattern to be ignored when walking directories
func (app *App) IgnorePath(r *regexp.Regexp) {
	app.ignoredPaths = append(app.ignoredPaths, r)
}

// IsIgnoredPath checks if a file path should be ignored
func (app *App) IsIgnoredPath(d string) bool {

	for _, v := range app.ignoredPaths {
		if v.MatchString(d) {
			return true
		}
	}
	return false
}

// GetConfig returns the application configuration
func (app *App) GetConfig() *Configuration {
	return app.config
}

// GetRouter returns the application router
func (app *App) GetRouter() *http.ServeMux {
	return app.router
}

// compileTemplates compiles all registered templates
func (app *App) compileTemplates() {
	const ext = ".html"

	// add default templates before everything else
	sub, _ := fs.Sub(defaultTemplates, "templates")
	app.templatesFSs = append([]fs.FS{sub}, app.templatesFSs...)
	// add theme directory after everything else to allow user to override any template
	if _, err := os.Stat("theme"); err == nil {
		app.templatesFSs = append(app.templatesFSs, os.DirFS("theme"))
	}

	app.templates = template.New("")
	for _, tfs := range app.templatesFSs {
		fs.WalkDir(tfs, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(p, ext) && d.Type().IsRegular() {
				ext := path.Ext(p)
				name := strings.TrimSuffix(p, ext)
				slog.Info("Template " + name)

				c, err := fs.ReadFile(tfs, p)
				if err != nil {
					return err
				}

				template.Must(app.templates.New(name).Funcs(app.helpers).Parse(string(c)))
			}

			return nil
		})
	}
}

// Partial executes a template by it's path name
func (app *App) Partial(path string, data Locals) template.HTML {
	v := app.templates.Lookup(path)
	if v == nil {
		return template.HTML(fmt.Sprintf("template %s not found", path))
	}

	if data == nil {
		data = Locals{}
	}

	data["config"] = app.config

	w := bytes.NewBufferString("")

	if err := v.Execute(w, data); err != nil {
		return template.HTML("rendering error " + path + " " + err.Error())
	}

	return template.HTML(w.String())
}

// requestLoggerHandler logs HTTP requests
func (app *App) requestLoggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w Response, r Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		slog.Info(r.Method+" "+r.URL.Path, "time", time.Since(start))
	})
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

// clearPagesCache clears the pages cache
func (app *App) clearPagesCache(p Page) {
	app.pages = nil
}

// Helper methods that need to be implemented
func (app *App) ago(t time.Time) string {
	if app.config.Readonly {
		return t.Format("Monday 2 January 2006")
	}

	d := time.Since(t)

	const day = time.Hour * 24
	const week = day * 7
	const month = day * 30
	const year = day * 365
	const maxPrecision = 2

	var o strings.Builder

	if d.Seconds() < 1 {
		o.WriteString("Less than a second ")
	}

	for precision := 0; d.Seconds() > 1 && precision < maxPrecision; precision++ {
		switch {
		case d >= year:
			years := d / year
			d -= years * year
			o.WriteString(fmt.Sprintf("%d years ", years))
		case d >= month:
			months := d / month
			d -= months * month
			o.WriteString(fmt.Sprintf("%d months ", months))
		case d >= week:
			weeks := d / week
			d -= weeks * week
			o.WriteString(fmt.Sprintf("%d weeks ", weeks))
		case d >= day:
			days := d / day
			d -= days * day
			o.WriteString(fmt.Sprintf("%d days ", days))
		case d >= time.Hour:
			hours := d / time.Hour
			d -= hours * time.Hour
			o.WriteString(fmt.Sprintf("%d hours ", hours))
		case d >= time.Minute:
			minutes := d / time.Minute
			d -= minutes * time.Minute
			o.WriteString(fmt.Sprintf("%d minutes ", minutes))
		case d >= time.Second:
			seconds := d / time.Second
			d -= seconds * time.Second
			o.WriteString(fmt.Sprintf("%d seconds ", seconds))
		}
	}

	o.WriteString("ago")

	return o.String()
}

// Properties returns a list of properties for a page
func (app *App) Properties(p Page) map[string]Property {

	ps := map[string]Property{}
	for _, source := range app.propsSources {
		for _, pr := range source(p) {
			ps[pr.Name()] = pr
		}
	}
	return ps
}

// Links returns a list of links for a Page
func (app *App) Links(p Page) []Command {

	cmds := []Command{}
	for _, l := range app.links {
		cmds = append(cmds, l(p)...)
	}
	return cmds
}

// RenderWidget renders all widget functions registered for a specific widget space
func (app *App) RenderWidget(s WidgetSpace, p Page) template.HTML {

	w, ok := app.widgets[s]
	if !ok {
		return ""
	}

	var o template.HTML
	for f := range w.All() {
		o += f(p)
	}
	return o
}

// Commands returns the list of commands for a page
func (app *App) Commands(p Page) []Command {

	cmds := []Command{}
	for _, c := range app.commands {
		cmds = append(cmds, c(p)...)
	}
	return cmds
}

// QuickCommands returns the list of quick commands for a page
func (app *App) QuickCommands(p Page) []Command {

	cmds := []Command{}
	for _, c := range app.quickCommands {
		cmds = append(cmds, c(p)...)
	}
	return cmds
}

// IsFontAwesome checks if an icon is a FontAwesome icon
func (app *App) IsFontAwesome(i string) bool {
	return strings.HasPrefix(i, "fa")
}

// includeJS adds a JavaScript library URL/path
func (app *App) includeJS(f string) template.HTML {

	if !slices.Contains(app.js, f) {
		app.js = append(app.js, f)
	}
	return ""
}

// scripts returns the HTML for all registered JavaScript files
func (app *App) scripts() template.HTML {

	var b strings.Builder
	for _, f := range app.js {
		fmt.Fprintf(&b, `<script src="%s" defer></script>`, f)
	}
	return template.HTML(b.String())
}

// Banner returns the banner image for a page
func (app *App) Banner(p Page) string {
	_, a := p.AST()
	if a == nil {
		return ""
	}

	paragraph := a.FirstChild()
	if paragraph == nil || paragraph.Kind() != ast.KindParagraph {
		return ""
	}

	img := paragraph.FirstChild()
	if img == nil || img.Kind() != ast.KindImage {
		return ""
	}

	image, ok := img.(*ast.Image)
	if !ok {
		return ""
	}

	dest := string(image.Destination)
	if len(dest) == 0 || dest == "#" {
		return ""
	}

	if !(path.IsAbs(dest) || strings.HasPrefix(dest, "http")) {
		d := path.Dir(p.FileName())
		dest = path.Join("/", d, dest)
	}

	return dest
}

// Emoji returns the emoji for a page
func (app *App) Emoji(p Page) string {
	_, tree := p.AST()
	if e, ok := FindInAST[*emojiAst.Emoji](tree); ok && e != nil {
		return string(e.Value.Unicode)
	}
	return ""
}

// dir returns the directory name
func (app *App) dir(s string) string {
	v := path.Dir(s)
	if v == "." {
		return ""
	}
	return v
}

// raw returns safe HTML
func (app *App) raw(i string) template.HTML {
	return template.HTML(i)
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

// Trigger triggers an event
func (app *App) Trigger(e PageEvent, p Page) {
	handlers, ok := app.pageEvents[e]

	if !ok {
		return
	}

	for _, h := range handlers {
		if err := h(p); err != nil {
			slog.Error("Failed to execute handler for event", "event", e, "handler", h, "error", err)
		}
	}
}

// Listen registers an event handler
func (app *App) Listen(e PageEvent, h PageEventHandler) {

	if _, ok := app.pageEvents[e]; !ok {
		app.pageEvents[e] = []PageEventHandler{}
	}

	app.pageEvents[e] = append(app.pageEvents[e], h)
}

// staticHandler handles static file serving
func (app *App) staticHandler(r Request) (Output, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	staticFSs := http.FS(
		priorityFS(
			append(app.staticDirs, os.DirFS(wd)),
		),
	)

	server := http.FileServer(staticFSs)

	cleanPath := path.Clean(r.URL.Path)

	if f, err := staticFSs.Open(cleanPath); err != nil {
		return nil, err
	} else {
		f.Close()
		return app.Cache(server.ServeHTTP), nil
	}
}

// build builds the static site
func (app *App) build(buildDir string) error {
	srv := app.server()

	// building Index separately
	err := app.buildRoute(
		srv,
		"/"+app.config.Index,
		buildDir,
		path.Join(buildDir, "index.html"),
	)

	if err != nil {
		slog.Error("Index Page may not exist, make sure your Index Page exists", "index", app.config.Index, "error", err)
	}

	errs := app.MapPage(context.Background(), func(p Page) error {
		err := app.buildRoute(
			srv,
			"/"+p.Name(),
			path.Join(buildDir, p.Name()),
			path.Join(buildDir, p.Name(), "index.html"),
		)

		if err != nil {
			return fmt.Errorf("Failed to process page: %s, error: %w", p.Name(), err)
		}

		return nil
	})

	if err := errors.Join(errs...); err != nil {
		slog.Error(err.Error())
	}

	// If we render 404 page
	// Copy 404 page from dest/404/index.html to /dest/404.html
	if in, err := os.Open(path.Join(buildDir, app.config.NotFoundPage, "index.html")); err == nil {
		defer in.Close()
		out, err := os.Create(path.Join(buildDir, "404.html"))
		if err != nil {
			slog.Error("Failed to open dest/404.html", "error", err)
		}
		defer out.Close()
		io.Copy(out, in)
	}

	extensionPageEnclosed := app.extensionPageEnclosed
	extensionPage := app.extensionPage
	buildPerms := app.buildPerms

	for route := range extensionPageEnclosed {
		err := app.buildRoute(
			srv,
			route,
			path.Join(buildDir, route),
			path.Join(buildDir, route, "index.html"),
		)

		if err != nil {
			slog.Error("Failed to process extension page", "route", route, "error", err)
		}
	}

	for route := range extensionPage {
		err := app.buildRoute(
			srv,
			route,
			path.Join(buildDir, path.Dir(route)),
			path.Join(buildDir, route),
		)

		if err != nil {
			slog.Error("Failed to process extension page", "route", route, "error", err)
		}
	}

	assets := GetAssets()
	return fs.WalkDir(assets, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		destPath := path.Join(buildDir, p)

		if entry.IsDir() {
			if err := os.MkdirAll(destPath, buildPerms); err != nil {
				return err
			}
		} else if _, err := os.Stat(destPath); err == nil {
			slog.Warn("Asset file already exists", "path", destPath)
		} else {
			content, err := fs.ReadFile(assets, p)
			if err != nil {
				return err
			}

			if err := os.WriteFile(destPath, content, buildPerms); err != nil {
				return err
			}
		}

		return nil
	})
}

// buildRoute builds a single route
func (app *App) buildRoute(srv *http.Server, route, dir, file string) error {
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		return err
	}

	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if err := os.MkdirAll(dir, app.buildPerms); err != nil {
		return err
	}

	if rec.Result().StatusCode != http.StatusOK {
		return errors.New(rec.Result().Status)
	}

	body, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		return err
	}
	defer rec.Result().Body.Close()

	return os.WriteFile(file, body, app.buildPerms)
}

// MapPage maps over all pages
func (app *App) MapPage(ctx context.Context, f func(Page) error) []error {
	sources := app.sources

	var errs []error
	for _, source := range sources {
		source.Each(ctx, func(p Page) {
			if err := f(p); err != nil {
				errs = append(errs, err)
			}
		})
	}
	return errs
}

// RegisterBuildPage registers a build page
func (app *App) RegisterBuildPage(p string, encloseInDir bool) {

	if encloseInDir {
		app.extensionPageEnclosed[p] = true
	} else {
		app.extensionPage[p] = true
	}
}

// Pages returns all pages
func (app *App) Pages(ctx context.Context) []Page {
	if app.pages == nil {
		app.populatePagesCache(ctx)
	}
	pages := app.pages
	return pages[:]
}

// EachPage iterates on all available pages
func (app *App) EachPage(ctx context.Context, f func(Page)) {
	if app.pages == nil {
		app.populatePagesCache(ctx)
	}
	currentPages := app.pages

	for _, p := range currentPages {
		select {
		case <-ctx.Done():
			return
		default:
			f(p)
		}
	}
}

// MapPageGeneric maps over all pages with generic return type
func MapPageGeneric[T any](app *App, ctx context.Context, f func(Page) T) []T {
	if app.pages == nil {
		app.populatePagesCache(ctx)
	}
	pages := app.pages
	concurrency := app.concurrency

	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(concurrency)

	output := make([]T, 0, len(pages))
	ch := make(chan T, concurrency)
	done := make(chan bool)

	go func() {
		for v := range ch {
			output = append(output, v)
		}
		done <- true
	}()

	for _, p := range pages {
		select {
		case <-ctx.Done():
			break
		default:
			grp.Go(func() (err error) {
				val := f(p)
				if isNil(val) {
					return
				}

				ch <- val

				return
			})
		}
	}

	grp.Wait()
	close(ch)
	<-done

	return output
}

// populatePagesCache populates the pages cache
func (app *App) populatePagesCache(ctx context.Context) {

	app.pages = make([]Page, 0, 1000)
	for _, s := range app.sources {
		select {
		case <-ctx.Done():
			return
		default:
			s.Each(ctx, func(p Page) {
				app.pages = append(app.pages, p)
			})
		}
	}
}

// isNil checks if a value is nil
func isNil[T any](t T) bool {
	v := reflect.ValueOf(t)
	kind := v.Kind()
	// Must be one of these types to be nillable
	return !v.IsValid() || (kind == reflect.Ptr ||
		kind == reflect.Interface ||
		kind == reflect.Slice ||
		kind == reflect.Map ||
		kind == reflect.Chan ||
		kind == reflect.Func) &&
		v.IsNil()
}

type priorityItem[T any] struct {
	Item     T
	Priority float32
}

type priorityList[T any] struct {
	items []priorityItem[T]
}

func (pl *priorityList[T]) Add(item T, priority float32) {
	pl.items = append(pl.items, priorityItem[T]{Item: item, Priority: priority})
	pl.sortByPriority()
}

func (pl *priorityList[T]) sortByPriority() {
	sort.Slice(pl.items, func(i, j int) bool {
		return pl.items[i].Priority < pl.items[j].Priority
	})
}

// An iterator over all items
func (pl *priorityList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range pl.items {
			if !yield(v.Item) {
				return
			}
		}
	}
}

type lastUpdateProp struct{ page Page }

func (a lastUpdateProp) Icon() string { return "fa-solid fa-clock" }
func (a lastUpdateProp) Name() string { return "modified" }
func (a lastUpdateProp) Value() any {
	app := GetApp()
	return app.ago(a.page.ModTime())
}

func DefaultProps(p Page) []Property {
	if p.ModTime().IsZero() {
		return nil
	}

	return []Property{
		lastUpdateProp{p},
	}
}
