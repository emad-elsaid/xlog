package xlog

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/csrf"
)

const xCSRF_COOKIE_NAME = "xlog_csrf"

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
			csrf.CookieName(xCSRF_COOKIE_NAME),
		),
		requestLoggerHandler,
	}
)

type (
	// alias http.ResponseWriter for shorter handler declaration
	Response = http.ResponseWriter
	// alias *http.Request for shorter handler declaration
	Request = *http.Request
	// alias of http.HandlerFunc as output is expected from defined http handlers
	Output = http.HandlerFunc
	// map of string to any value used for template rendering
	Locals map[string]interface{} // passed to templates
)

func init() {
	flag.StringVar(&bind_address, "bind", "127.0.0.1:3000", "IP and port to bind the web server to")
}

func serve() {
	srv := server()
	log.Printf("Starting server: %s", bind_address)
	log.Fatal(srv.ListenAndServe())
}

func server() *http.Server {
	compileTemplates()
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

type (
	// Checks a request for conditions, may modify request returning the new
	// request and if conditions are met.
	//
	// Can be used to check request method, path or other attributes against
	// expected values.
	RouteCheck func(Request) (Request, bool)
	// A route has handler function and set of RouteChecks if all checks are
	// true, the last request will be passed to the handler function
	Route struct {
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

type contextIndex int

const varsIndex contextIndex = iota + 1

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

// Vars returns a map of key = dynamic segment in the path, value = the value of
// the segment registering a path like `"/edit/{page:.+}"` calling `Vars in the
// handler will return a map with one key `page` and the value is that part of
// the path in r
func Vars(r Request) map[string]string {
	if rv := r.Context().Value(varsIndex); rv != nil {
		return rv.(map[string]string)
	}
	return map[string]string{}
}

// HandlerFunc is the type of an HTTP handler function + returns output function.
// it makes it easier to return the output directly instead of writing the output to w then return.
type HandlerFunc func(http.ResponseWriter, *http.Request) Output

func handlerFuncToHttpHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)(w, r)
	}
}

// NotFound returns an output function that writes 404 NotFound to http response
func NotFound(msg string) Output {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusNotFound)
	}
}

// BadRequest returns an output function that writes BadRequest http response
func BadRequest(msg string) Output {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, msg, http.StatusBadRequest)
	}
}

// Unauthorized returns an output function that writes Unauthorized http response
func Unauthorized(msg string) Output {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusUnauthorized)
	}
}

// InternalServerError returns an output function that writes InternalServerError http response
func InternalServerError(err error) Output {
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

// NoContent returns an output function that writes NoContent http status
func NoContent() Output {
	return noContent
}

func noContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// PlainText returns an output function that writes text to response writer
func PlainText(text string) Output {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(text))
	}
}

// Match Adds a new HTTP handler function to the list of routes with a list of checks functions.
// the list of checks are executed when a request comes in if all of them returned true the handler function gets executed.
func Match(route http.HandlerFunc, checks ...RouteCheck) Route {
	r := Route{
		checks: checks,
		route:  route,
	}
	router.routes = append(router.routes, r)

	return r
}

// Get defines a new route that gets executed when the request matches path and
// method is http Get. the list of middlewares are executed in order
func Get(path string, handler HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) Route {
	return Match(
		applyMiddlewares(handlerFuncToHttpHandler(handler), middlewares...),
		checkMethod(http.MethodGet), checkPath(path),
	)
}

// Post defines a new route that gets executed when the request matches path and
// method is http Post. the list of middlewares are executed in order
func Post(path string, handler HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) Route {
	return Match(
		applyMiddlewares(handlerFuncToHttpHandler(handler), middlewares...),
		checkMethod(http.MethodPost), checkPath(path),
	)
}

// Delete defines a new route that gets executed when the request matches path and
// method is http Delete. the list of middlewares are executed in order
func Delete(path string, handler HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) Route {
	return Match(
		applyMiddlewares(handlerFuncToHttpHandler(handler), middlewares...),
		checkMethod(http.MethodDelete), checkPath(path),
	)
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
		defer echo(xINFO, r.Method, r.URL.Path)()
		h.ServeHTTP(w, r)
	})
}
