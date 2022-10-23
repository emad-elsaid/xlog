package xlog

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

//go:embed assets
var assets embed.FS

var (
	bind_address string
	router       = &mux{}
	// a function that renders CSRF hidden input field
	CSRF = csrf.TemplateField

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
		requestLoggerHandler,
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
	flag.StringVar(&bind_address, "bind", "127.0.0.1:3000", "IP and port to bind the web server to")
}

func server() *http.Server {
	compileViews()
	var handler http.Handler = router
	for _, v := range middlewares {
		handler = v(handler)
	}

	return &http.Server{
		Handler:      handler,
		Addr:         bind_address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func serve() {
	srv := server()
	log.Printf("Starting server: %s", bind_address)
	log.Fatal(srv.ListenAndServe())
}

type (
	RouteCheck func(Request) (Request, bool)
	Route      struct {
		checks []RouteCheck
		route  http.HandlerFunc
	}

	mux struct {
		routes []Route
	}
)

func (h *mux) ServeHTTP(w Response, r Request) {
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

func Vars(r Request) map[string]string {
	if rv := r.Context().Value(varsIndex); rv != nil {
		return rv.(map[string]string)
	}
	return map[string]string{}
}

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

// HandlerFunc is the type of an HTTP handler function + returns output function.
// it makes it easier to return the output directly instead of writing the output to w then return.
type HandlerFunc func(http.ResponseWriter, *http.Request) Output

func handlerFuncToHttpHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)(w, r)
	}
}

// NotFound an output function that writes 404 NotFound to http response
func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNotFound)
}

// BadRequest an output function that writes BadRequest http response
func BadRequest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusBadRequest)
}

// Unauthorized an output function that writes Unauthorized http response
func Unauthorized(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusUnauthorized)
}

// InternalServerError an output function that writes InternalServerError http response
func InternalServerError(err error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Redirect returns an output function that writes Found http response to provided URL
func Redirect(url string) Output {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, url, http.StatusFound)
	}
}

// NoContent an output function that writes NoContent http response
func NoContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// PlainText returns an output function that writes text to response writer
func PlainText(text string) Output {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(text))
	}
}

// ROUTE Adds a new HTTP handler function to the list of routes with a list of checks functions.
// the list of checks are executed when a request comes in if all of them returned true the handler function gets executed.
func ROUTE(route http.HandlerFunc, checks ...RouteCheck) Route {
	r := Route{
		checks: checks,
		route:  route,
	}
	router.routes = append(router.routes, r)

	return r
}

// GET defines a new route that gets executed when the request matches path and
// method is http GET. the list of middlewares are executed in order
func GET(path string, handler HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) Route {
	return ROUTE(
		applyMiddlewares(handlerFuncToHttpHandler(handler), middlewares...),
		checkMethod(http.MethodGet), checkPath(path),
	)
}

// POST defines a new route that gets executed when the request matches path and
// method is http POST. the list of middlewares are executed in order
func POST(path string, handler HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) Route {
	return ROUTE(
		applyMiddlewares(handlerFuncToHttpHandler(handler), middlewares...),
		checkMethod(http.MethodPost), checkPath(path),
	)
}

// DELETE defines a new route that gets executed when the request matches path and
// method is http DELETE. the list of middlewares are executed in order
func DELETE(path string, handler HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) Route {
	return ROUTE(
		applyMiddlewares(handlerFuncToHttpHandler(handler), middlewares...),
		checkMethod(http.MethodDelete), checkPath(path),
	)
}

//go:embed views
var defaultViews embed.FS
var templates *template.Template
var helpers = template.FuncMap{}
var views []fs.FS

// VIEW registers a filesystem as a view, views are registered such that the
// latest view directory override older ones. views file extensions are
// signified by VIEWS_EXTENSION constant and the file path can be used as
// template name without this extension
func VIEW(view fs.FS) {
	views = append(views, view)
}

func compileViews() {
	// add default views before everything else
	sub, _ := fs.Sub(defaultViews, "views")
	views = append([]fs.FS{sub}, views...)

	templates = template.New("")
	for _, view := range views {
		fs.WalkDir(view, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(p, VIEWS_EXTENSION) && d.Type().IsRegular() {
				ext := path.Ext(p)
				name := strings.TrimSuffix(p, ext)
				defer Log(DEBUG, "View", name)()

				c, err := fs.ReadFile(view, p)
				if err != nil {
					return err
				}

				template.Must(templates.New(name).Funcs(helpers).Parse(string(c)))
			}

			return nil
		})
	}
}

func Partial(path string, data Locals) string {
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
	data["SIDEBAR"] = SIDEBAR

	w := bytes.NewBufferString("")
	err := v.Execute(w, data)
	if err != nil {
		return "rendering error " + path + " " + err.Error()
	}

	return w.String()
}

// Render returns an output function that renders partial with data and writes it as response
func Render(path string, data Locals) Output {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, Partial(path, data))
	}
}

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

func requestLoggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Log(INFO, r.Method, r.URL.Path)()
		h.ServeHTTP(w, r)
	})
}

// HELPER registers a new helper function. all helpers are used when compiling
// view/templates so registering helpers function must happen before the server
// starts as compiling views happend right before starting the http server.
func HELPER(name string, f interface{}) {
	if _, ok := helpers[name]; ok {
		log.Fatalf("Helper: %s already registered", name)
	}

	helpers[name] = f
}
