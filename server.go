package xlog

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/csrf"
)

var (
	// a function that returns the CSRF token
	CSRF = csrf.Token
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

// HandlerFunc is the type of an HTTP handler function + returns output function.
// it makes it easier to return the output directly instead of writing the output to w then return.
type HandlerFunc func(Request) Output

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
	slog.Info("GET", "path", path, "func", funcStringer{handler})
	app.router.HandleFunc("GET "+path, app.handlerFuncToHttpHandler(handler))
}

// Post registers a POST route
func (app *App) Post(path string, handler HandlerFunc) {
	slog.Info("POST", "path", path, "func", funcStringer{handler})
	app.router.HandleFunc("POST "+path, app.handlerFuncToHttpHandler(handler))
}

// Delete registers a DELETE route
func (app *App) Delete(path string, handler HandlerFunc) {
	slog.Info("DELETE", "path", path, "func", funcStringer{handler})
	app.router.HandleFunc("DELETE "+path, app.handlerFuncToHttpHandler(handler))
}

func (app *App) handlerFuncToHttpHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w Response, r Request) {
		handler(r)(w, r)
	}
}

// NotFound returns an output function that writes 404 NotFound to http response
func (app *App) NotFound(msg string) Output {
	return func(w Response, r Request) {
		http.Error(w, msg, http.StatusNotFound)
	}
}

// BadRequest returns an output function that writes BadRequest http response
func (app *App) BadRequest(msg string) Output {
	return func(w Response, r Request) {
		http.Error(w, msg, http.StatusBadRequest)
	}
}

// InternalServerError returns an output function that writes InternalServerError http response
func (app *App) InternalServerError(err error) Output {
	return func(w Response, r Request) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Redirect returns an output function that writes Found http response to provided URL
func (app *App) Redirect(url string) Output {
	return func(w Response, r Request) {
		http.Redirect(w, r, url, http.StatusFound)
	}
}

// NoContent returns an output function that writes NoContent http status
func (app *App) NoContent() Output {
	return func(w Response, r Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

// PlainText returns an output function that writes text to response writer
func (app *App) PlainText(text string) Output {
	return func(w Response, r Request) {
		w.Write([]byte(text))
	}
}

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

// Render returns an output function that renders partial with data and writes it as response
func (app *App) Render(path string, data Locals) Output {
	return func(w Response, r Request) {
		fmt.Fprint(w, app.Partial(path, data))
	}
}

// requestLoggerHandler logs HTTP requests
func (app *App) requestLoggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w Response, r Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		slog.Info(r.Method+" "+r.URL.Path, "time", time.Since(start))
	})
}

// Cache wraps Output and adds header to instruct the browser to cache the output
func (app *App) Cache(out Output) Output {
	return func(w Response, r Request) {
		w.Header().Add("Cache-Control", "max-age=604800")
		out(w, r)
	}
}

type funcStringer struct{ any }

func (f funcStringer) String() string {
	const xlogPrefix = "emad-elsaid/xlog/"
	const ghPrefix = "github.com/"
	name := runtime.FuncForPC(reflect.ValueOf(f.any).Pointer()).Name()
	name = strings.TrimPrefix(name, ghPrefix)
	name = strings.TrimPrefix(name, xlogPrefix)
	return name
}
