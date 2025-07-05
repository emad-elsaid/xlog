package xlog

import (
	"encoding/json"
	"net/http"
	"reflect"
	"runtime"
	"strings"

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

// Get defines a new route that gets executed when the request matches path and
// method is http Get.
func Get(path string, handler HandlerFunc) {
	app := GetApp()
	app.Get(path, handler)
}

// Post defines a new route that gets executed when the request matches path and
// method is http Post.
func Post(path string, handler HandlerFunc) {
	app := GetApp()
	app.Post(path, handler)
}

// Delete defines a new route that gets executed when the request matches path and
// method is http Delete.
func Delete(path string, handler HandlerFunc) {
	app := GetApp()
	app.Delete(path, handler)
}

// Render returns an output function that renders partial with data and writes it as response
func Render(path string, data Locals) Output {
	app := GetApp()
	return app.Render(path, data)
}

// NotFound returns an output function that writes 404 NotFound to http response
func NotFound(msg string) Output {
	app := GetApp()
	return app.NotFound(msg)
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
	app := GetApp()
	return app.Redirect(url)
}

// NoContent returns an output function that writes NoContent http status
func NoContent() Output {
	app := GetApp()
	return app.NoContent()
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
