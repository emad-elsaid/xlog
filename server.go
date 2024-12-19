package xlog

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"
)

var (
	router = http.NewServeMux()
	// a function that renders CSRF hidden input field
	CSRF = csrf.TemplateField
)

type (
	// alias http.ResponseWriter for shorter handler declaration
	Response = http.ResponseWriter
	// alias *http.Request for shorter handler declaration
	Request = *http.Request
	// alias of http.HandlerFunc as output is expected from defined http handlers
	Output = http.HandlerFunc
	// map of string to any value used for template rendering
	Locals map[string]any // passed to templates
)

func defaultMiddlewares(readonly bool) (middlewares []func(http.Handler) http.Handler) {
	if !readonly {
		crsfOpts := []csrf.Option{
			csrf.Path("/"),
			csrf.FieldName("csrf"),
			csrf.CookieName(Config.CsrfCookieName),
			csrf.Secure(!Config.ServeInsecure),
		}

		sessionSecret := []byte(os.Getenv("SESSION_SECRET"))
		if len(sessionSecret) == 0 {
			sessionSecret = make([]byte, 128)
			rand.Read(sessionSecret)
		}

		middlewares = append(middlewares,
			methodOverrideHandler,
			csrf.Protect(sessionSecret, crsfOpts...))
	}

	middlewares = append(middlewares, requestLoggerHandler)

	return
}

func server() *http.Server {
	compileTemplates()
	var handler http.Handler = router
	for _, v := range defaultMiddlewares(Config.Readonly) {
		handler = v(handler)
	}

	return &http.Server{
		Handler:      handler,
		Addr:         Config.BindAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

// HandlerFunc is the type of an HTTP handler function + returns output function.
// it makes it easier to return the output directly instead of writing the output to w then return.
type HandlerFunc func(Response, Request) Output

func handlerFuncToHttpHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w Response, r Request) {
		handler(w, r)(w, r)
	}
}

// NotFound returns an output function that writes 404 NotFound to http response
func NotFound(msg string) Output {
	return func(w Response, r Request) {
		http.Error(w, "", http.StatusNotFound)
	}
}

// BadRequest returns an output function that writes BadRequest http response
func BadRequest(msg string) Output {
	return func(w Response, r Request) {
		http.Error(w, msg, http.StatusBadRequest)
	}
}

// Unauthorized returns an output function that writes Unauthorized http response
func Unauthorized(msg string) Output {
	return func(w Response, r Request) {
		http.Error(w, "", http.StatusUnauthorized)
	}
}

// InternalServerError returns an output function that writes InternalServerError http response
func InternalServerError(err error) Output {
	return func(w Response, r Request) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Redirect returns an output function that writes Found http response to provided URL
func Redirect(url string) Output {
	return func(w Response, r Request) {
		http.Redirect(w, r, url, http.StatusFound)
	}
}

// NoContent returns an output function that writes NoContent http status
func NoContent() Output {
	return func(w Response, r Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

// Noop is an output that doesn't do anything to the request. can be useful for a websocket upgrader
func Noop(w Response, r Request) {}

// PlainText returns an output function that writes text to response writer
func PlainText(text string) Output {
	return func(w Response, r Request) {
		w.Write([]byte(text))
	}
}

func JsonResponse(a any) Output {
	return func(w Response, r Request) {
		b, err := json.Marshal(a)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(b)
	}
}

// Get defines a new route that gets executed when the request matches path and
// method is http Get. the list of middlewares are executed in order
func Get(path string, handler HandlerFunc) {
	slog.Info("GET", "path", path, "func", callerName(handler))
	router.HandleFunc("GET "+path, handlerFuncToHttpHandler(handler))
}

// Post defines a new route that gets executed when the request matches path and
// method is http Post. the list of middlewares are executed in order
func Post(path string, handler HandlerFunc) {
	slog.Info("POST", "path", path, "func", callerName(handler))
	router.HandleFunc("POST "+path, handlerFuncToHttpHandler(handler))
}

// Delete defines a new route that gets executed when the request matches path and
// method is http Delete. the list of middlewares are executed in order
func Delete(path string, handler HandlerFunc) {
	slog.Info("DELETE", "path", path, "func", callerName(handler))
	router.HandleFunc("DELETE "+path, handlerFuncToHttpHandler(handler))
}

// Render returns an output function that renders partial with data and writes it as response
func Render(path string, data Locals) Output {
	return func(w Response, r Request) {
		fmt.Fprint(w, Partial(path, data))
	}
}

// Derived from Gorilla middleware https://github.com/gorilla/handlers/blob/v1.5.1/handlers.go#L134
func methodOverrideHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w Response, r Request) {
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
	return http.HandlerFunc(func(w Response, r Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		slog.Info(r.Method+" "+r.URL.Path, "time", time.Since(start))
	})
}

// Cache wraps Output and adds header to instruct the browser to cache the output
func Cache(out Output) Output {
	return func(w Response, r Request) {
		w.Header().Add("Cache-Control", "max-age=604800")
		out(w, r)
	}
}
