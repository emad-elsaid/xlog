package main

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/csrf"
)

const (
	STATIC_DIR_PATH  = "public"
	ASSETS_DIR_PATH  = "assets"
	VIEWS_EXTENSION  = ".html"
	CSRF_COOKIE_NAME = "xlog_csrf"
)

const (
	_        = iota
	KB int64 = 1 << (10 * iota)
	MB
	GB
	TB
	PB
)

//go:embed assets
var assets embed.FS

var (
	BIND_ADDRESS string
	router       = &Handler{}
	CSRF         = csrf.TemplateField

	dynamicSegmentWithPatternRegexp = regexp.MustCompile("{([^}]+):([^}]+)}")
	dynamicSegmentRegexp            = regexp.MustCompile("{([^}]+)}")
	middlewares                     = []func(http.Handler) http.Handler{
		methodOverrideHandler,
		csrf.Protect(
			[]byte(os.Getenv("SESSION_SECRET")),
			csrf.Path("/"),
			csrf.FieldName("csrf"),
			csrf.CookieName(CSRF_COOKIE_NAME),
		),
		RequestLoggerHandler,
	}
)

// Some aliases to make it easier
type (
	Response = http.ResponseWriter
	Request  = *http.Request
	Output   = http.HandlerFunc
	Locals   map[string]interface{} // passed to views/templates
)

func init() {
	flag.StringVar(&BIND_ADDRESS, "bind", "127.0.0.1:3000", "IP and port to bind the web server to")
	log.SetFlags(log.Ltime)

	GET("/"+ASSETS_DIR_PATH+"/.*", assetsHandler())
	GET("/"+STATIC_DIR_PATH+"/.*", staticHandler())
}

func server() *http.Server {
	compileViews()
	var handler http.Handler = router
	for _, v := range middlewares {
		handler = v(handler)
	}

	return &http.Server{
		Handler:      handler,
		Addr:         BIND_ADDRESS,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func START() {
	srv := server()
	log.Printf("Starting server: %s", BIND_ADDRESS)
	log.Fatal(srv.ListenAndServe())
}

// Mux/Handler ===========================================
type (
	RouteCheck func(Request) (Request, bool)
	Route      struct {
		checks []RouteCheck
		route  http.HandlerFunc
	}

	Handler struct {
		routes []Route
	}
)

func (h *Handler) ServeHTTP(w Response, r Request) {
ROUTES:
	for _, route := range h.routes {
		rn := r
		ok := false
		for _, check := range route.checks {
			if rn, ok = check(rn); !ok {
				continue ROUTES
			}
		}

		route.route(w, rn)
		return
	}
}

func checkMethod(method string) RouteCheck {
	return func(r Request) (Request, bool) { return r, r.Method == method }
}

const varsIndex int = iota + 1

func checkPath(path string) RouteCheck {
	path = dynamicSegmentWithPatternRegexp.ReplaceAllString(path, "(?P<$1>$2)")
	path = dynamicSegmentRegexp.ReplaceAllString(path, "(?P<$1>[^/]+)")
	path = "^" + path + "$"
	reg := regexp.MustCompile(path)
	groups := reg.SubexpNames()

	return func(r Request) (Request, bool) {
		if !reg.MatchString(r.URL.Path) {
			return r, false
		}

		values := reg.FindStringSubmatch(r.URL.Path)
		vars := map[string]string{}
		for i, g := range groups {
			vars[g] = values[i]
		}

		ctx := context.WithValue(r.Context(), varsIndex, vars)
		return r.WithContext(ctx), true
	}
}

func VARS(r Request) map[string]string {
	if rv := r.Context().Value(varsIndex); rv != nil {
		return rv.(map[string]string)
	}
	return map[string]string{}
}

// LOGGING ===============================================

const (
	DEBUG = "\033[97;42m"
	INFO  = "\033[97;42m"
)

func Log(level, label, text string, args ...interface{}) func() {
	start := time.Now()
	return func() {
		if len(args) > 0 {
			log.Printf("%s %s \033[0m (%s) %s %v", level, label, time.Now().Sub(start), text, args)
		} else {
			log.Printf("%s %s \033[0m (%s) %s", level, label, time.Now().Sub(start), text)
		}
	}
}

// ROUTES HELPERS ==========================================

type HandlerFunc func(http.ResponseWriter, *http.Request) http.HandlerFunc

func handlerFuncToHttpHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)(w, r)
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotFound)
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusBadRequest)
}

func Unauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusUnauthorized)
}

func InternalServerError(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Redirect(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func NoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func PlainText(text string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(text))
	}
}

func ROUTE(route http.HandlerFunc, checks ...RouteCheck) {
	router.routes = append(router.routes, Route{
		checks: checks,
		route:  route,
	})
}

func GET(path string, handler HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	ROUTE(
		applyMiddlewares(handlerFuncToHttpHandler(handler), middlewares...),
		checkMethod(http.MethodGet), checkPath(path),
	)
}

func POST(path string, handler HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	ROUTE(
		applyMiddlewares(handlerFuncToHttpHandler(handler), middlewares...),
		checkMethod(http.MethodPost), checkPath(path),
	)
}

func DELETE(path string, handler HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	ROUTE(
		applyMiddlewares(handlerFuncToHttpHandler(handler), middlewares...),
		checkMethod(http.MethodDelete), checkPath(path),
	)
}

// VIEWS ====================

//go:embed views
var views embed.FS
var templates *template.Template
var helpers = template.FuncMap{}

func compileViews() {
	templates = template.New("")
	fs.WalkDir(views, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(p, VIEWS_EXTENSION) && d.Type().IsRegular() {
			rel := strings.TrimPrefix(p, "views/")
			ext := path.Ext(rel)
			name := strings.TrimSuffix(rel, ext)
			defer Log(DEBUG, "View", name)()

			c, err := fs.ReadFile(views, p)
			if err != nil {
				return err
			}

			template.Must(templates.New(name).Funcs(helpers).Parse(string(c)))
		}

		return nil
	})
}

func partial(path string, data Locals) string {
	v := templates.Lookup(path)
	if v == nil {
		return fmt.Sprintf("view %s not found", path)
	}

	// set extra locals here
	if data == nil {
		data = Locals{}
	}

	data["SITENAME"] = SITENAME
	data["READONLY"] = READONLY

	w := bytes.NewBufferString("")
	err := v.Execute(w, data)
	if err != nil {
		return "rendering error " + path + " " + err.Error()
	}

	return w.String()
}

func Render(path string, data Locals) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, partial(path, data))
	}
}

// SERVER MIDDLEWARES ==============================

func applyMiddlewares(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// Derived from Gorilla middleware https://github.com/gorilla/handlers/blob/v1.5.1/handlers.go#L134
func methodOverrideHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			om := r.FormValue("_method")
			if om == "PUT" || om == "PATCH" || om == "DELETE" {
				r.Method = om
			}
		}
		h.ServeHTTP(w, r)
	})
}

func RequestLoggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Log(INFO, r.Method, r.URL.Path)()
		h.ServeHTTP(w, r)
	})
}

// HANDLERS ==============================

func assetsHandler() HandlerFunc {
	assetsServer := http.FileServer(http.FS(assets))
	return func(w Response, _ Request) Output {
		w.Header().Add("Cache-Control", "max-age=31536000")
		return assetsServer.ServeHTTP
	}
}

func staticHandler() HandlerFunc {
	dir := http.Dir(STATIC_DIR_PATH)
	server := http.FileServer(dir)
	staticHandler := http.StripPrefix("/"+STATIC_DIR_PATH, server)

	return func(w Response, r Request) Output {
		if strings.HasSuffix(r.URL.Path, "/") {
			return NotFound
		}

		w.Header().Add("Cache-Control", "max-age=31536000")
		return staticHandler.ServeHTTP
	}
}

// HELPERS FUNCTIONS ======================

func HELPER(name string, f interface{}) {
	if _, ok := helpers[name]; ok {
		log.Fatalf("Helper: %s has been defined already", name)
	}

	helpers[name] = f
}
