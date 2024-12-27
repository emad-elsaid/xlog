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
	router = http.NewServeMux()
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

func defaultMiddlewares() (middlewares []func(http.Handler) http.Handler) {
	if !Config.Readonly {
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
			csrf.Protect(sessionSecret, crsfOpts...))
	}

	middlewares = append(middlewares, requestLoggerHandler)

	return
}

func server() *http.Server {
	compileTemplates()
	var handler http.Handler = router
	for _, v := range defaultMiddlewares() {
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
type HandlerFunc func(Request) Output

func handlerFuncToHttpHandler(handler HandlerFunc) http.HandlerFunc {
	return func(w Response, r Request) {
		handler(r)(w, r)
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
// method is http Get.
func Get(path string, handler HandlerFunc) {
	slog.Info("GET", "path", path, "func", funcStringer{handler})
	router.HandleFunc("GET "+path, handlerFuncToHttpHandler(handler))
}

// Post defines a new route that gets executed when the request matches path and
// method is http Post.
func Post(path string, handler HandlerFunc) {
	slog.Info("POST", "path", path, "func", funcStringer{handler})
	router.HandleFunc("POST "+path, handlerFuncToHttpHandler(handler))
}

// Delete defines a new route that gets executed when the request matches path and
// method is http Delete.
func Delete(path string, handler HandlerFunc) {
	slog.Info("DELETE", "path", path, "func", funcStringer{handler})
	router.HandleFunc("DELETE "+path, handlerFuncToHttpHandler(handler))
}

// Render returns an output function that renders partial with data and writes it as response
func Render(path string, data Locals) Output {
	return func(w Response, r Request) {
		fmt.Fprint(w, Partial(path, data))
	}
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

type funcStringer struct{ any }

func (f funcStringer) String() string {
	const xlogPrefix = "emad-elsaid/xlog/"
	const ghPrefix = "github.com/"
	name := runtime.FuncForPC(reflect.ValueOf(f.any).Pointer()).Name()
	name = strings.TrimPrefix(name, ghPrefix)
	name = strings.TrimPrefix(name, xlogPrefix)
	return name
}
