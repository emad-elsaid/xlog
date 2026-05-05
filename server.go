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

	csrf "filippo.io/csrf/gorilla"
)

var (
	router = http.NewServeMux()
	// CSRF returns a placeholder token for templates.
	// Note: filippo.io/csrf/gorilla uses Fetch metadata headers, not tokens.
	// This is kept for template compatibility only.
	CSRF = csrf.Token //nolint:staticcheck // Deprecated but needed for template compatibility
)

type (
	// alias http.ResponseWriter for shorter handler declaration.
	Response = http.ResponseWriter
	// alias *http.Request for shorter handler declaration.
	Request = *http.Request
	// alias of http.HandlerFunc as output is expected from defined http handlers.
	Output = http.HandlerFunc
	// map of string to any value used for template rendering.
	Locals map[string]any // passed to templates
)

func defaultMiddlewares() (middlewares []func(http.Handler) http.Handler) {
	if !Config.Readonly {
		sessionSecret := []byte(os.Getenv("SESSION_SECRET"))
		if len(sessionSecret) == 0 {
			sessionSecret = make([]byte, 128)
			if _, err := rand.Read(sessionSecret); err != nil {
				slog.Error("Failed to generate session secret", "error", err)
				osExit(1)
			}
		}

		middlewares = append(middlewares,
			csrf.Protect(sessionSecret))
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

// NotFound returns an output function that writes 404 NotFound to http response.
func NotFound(msg string) Output {
	return func(w Response, r Request) {
		http.Error(w, msg, http.StatusNotFound)
	}
}

// BadRequest returns an output function that writes BadRequest http response.
func BadRequest(msg string) Output {
	return func(w Response, r Request) {
		http.Error(w, msg, http.StatusBadRequest)
	}
}

// Unauthorized returns an output function that writes 401 Unauthorized http response.
// Commonly used when authentication credentials are missing or invalid.
func Unauthorized(msg string) Output {
	return func(w Response, r Request) {
		http.Error(w, msg, http.StatusUnauthorized)
	}
}

// Forbidden returns an output function that writes 403 Forbidden http response.
// Used when the client is authenticated but lacks permission to access the resource.
func Forbidden(msg string) Output {
	return func(w Response, r Request) {
		http.Error(w, msg, http.StatusForbidden)
	}
}

// MethodNotAllowed returns an output function that writes 405 Method Not Allowed http response.
// Used when the HTTP method is not supported for the requested resource.
func MethodNotAllowed(msg string) Output {
	return func(w Response, r Request) {
		http.Error(w, msg, http.StatusMethodNotAllowed)
	}
}

// InternalServerError returns an output function that writes InternalServerError http response.
func InternalServerError(err error) Output {
	return func(w Response, r Request) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Redirect returns an output function that writes Found http response to provided URL.
func Redirect(url string) Output {
	return func(w Response, r Request) {
		http.Redirect(w, r, url, http.StatusFound)
	}
}

// NoContent returns an output function that writes NoContent http status.
func NoContent() Output {
	return func(w Response, r Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

// Created returns an output function that writes http status 201 to response
// with optional location header for the newly created resource.
func Created(location string) Output {
	return func(w Response, r Request) {
		if location != "" {
			w.Header().Set("Location", location)
		}
		w.WriteHeader(http.StatusCreated)
	}
}

// Accepted returns an output function that writes http status 202 to response
// indicating that the request has been accepted for processing but not completed.
func Accepted() Output {
	return func(w Response, r Request) {
		w.WriteHeader(http.StatusAccepted)
	}
}

// healthHandler is an HTTP handler that returns 200 OK for health checks.
// This endpoint is useful for load balancers, Kubernetes probes, and monitoring services.
func healthHandler(w Response, r Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK\n")); err != nil {
		slog.Error("failed to write health check response", "error", err)
	}
}

// PlainText returns an output function that writes text to response writer.
func PlainText(text string) Output {
	return func(w Response, r Request) {
		if _, err := w.Write([]byte(text)); err != nil {
			slog.Error("failed to write plain text response", "error", err)
		}
	}
}

// JsonResponse returns an output function that encodes data as JSON and writes it to response
// with appropriate content-type header.
func JsonResponse(a any) Output {
	return func(w Response, r Request) {
		b, err := json.Marshal(a)
		if err != nil {
			if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
				slog.Error("failed to write error response", "error", writeErr)
			}
			return
		}

		if _, err := w.Write(b); err != nil {
			slog.Error("failed to write JSON response", "error", err)
		}
	}
}

// registerRoute is a helper function that registers an HTTP route with the given method.
// It logs the registration and adds the handler to the router.
func registerRoute(method, path string, handler HandlerFunc) {
	slog.Info(method, "path", path, "func", funcStringer{handler})
	router.HandleFunc(method+" "+path, handlerFuncToHttpHandler(handler))
}

// Get defines a new route that gets executed when the request matches path and
// method is http Get.
func Get(path string, handler HandlerFunc) {
	registerRoute("GET", path, handler)
}

// Post defines a new route that gets executed when the request matches path and
// method is http Post.
func Post(path string, handler HandlerFunc) {
	registerRoute("POST", path, handler)
}

// Delete defines a new route that gets executed when the request matches path and
// method is http Delete.
func Delete(path string, handler HandlerFunc) {
	registerRoute("DELETE", path, handler)
}

// Put defines a new route that gets executed when the request matches path and
// method is http Put.
func Put(path string, handler HandlerFunc) {
	registerRoute("PUT", path, handler)
}

// Patch defines a new route that gets executed when the request matches path and
// method is http Patch.
func Patch(path string, handler HandlerFunc) {
	registerRoute("PATCH", path, handler)
}

// Head defines a new route that gets executed when the request matches path and
// method is http Head.
func Head(path string, handler HandlerFunc) {
	registerRoute("HEAD", path, handler)
}

// Options defines a new route that gets executed when the request matches path and
// method is http Options.
func Options(path string, handler HandlerFunc) {
	registerRoute("OPTIONS", path, handler)
}

// Render returns an output function that renders partial with data and writes it as response.
func Render(path string, data Locals) Output {
	return func(w Response, r Request) {
		if _, err := fmt.Fprint(w, Partial(path, data)); err != nil {
			slog.Error("failed to render template", "path", path, "error", err)
		}
	}
}

// sanitizeLogString removes newline and carriage return characters from log strings
// to prevent log injection attacks where malicious input could split log entries.
func sanitizeLogString(s string) string {
	return strings.NewReplacer("\n", " ", "\r", " ").Replace(s)
}

func requestLoggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w Response, r Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		logMsg := sanitizeLogString(r.Method + " " + r.URL.Path)
		slog.Info(logMsg, "time", time.Since(start)) // #nosec G706 -- Input sanitized via sanitizeLogString
	})
}

// Cache wraps Output and adds header to instruct the browser to cache the output.
func Cache(out Output) Output {
	return func(w Response, r Request) {
		w.Header().Add("Cache-Control", "max-age=604800")
		out(w, r)
	}
}

// NoCache wraps Output and adds headers to prevent browser and proxy caching.
// Useful for dynamic content, APIs, and pages that should always be fresh.
func NoCache(out Output) Output {
	return func(w Response, r Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
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
