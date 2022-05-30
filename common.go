package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	_ "embed"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const (
	APP_NAME         = "xlog"
	STATIC_DIR_PATH  = "public"
	ASSETS_DIR_PATH  = "assets"
	VIEWS_EXTENSION  = ".html"
	CSRF_COOKIE_NAME = APP_NAME + "_csrf"
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
	AUTORELOAD   bool
	STARTUP_TIME = time.Now()
	router       = mux.NewRouter()
	VARS         = mux.Vars
	CSRF         = csrf.TemplateField
	middlewares  = []func(http.Handler) http.Handler{
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
type Response = http.ResponseWriter
type Request = *http.Request
type Output = http.HandlerFunc
type Locals map[string]interface{} // passed to views/templates

func init() {
	flag.StringVar(&BIND_ADDRESS, "bind", "127.0.0.1:3000", "IP and port to bind the web server to")
	flag.BoolVar(&AUTORELOAD, "autoload", false, "reload the page when the server restarts")
	log.SetFlags(log.Ltime)
	HELPER("autoload", func() bool { return AUTORELOAD })
}

func Start() {
	compileViews()
	router.PathPrefix("/" + STATIC_DIR_PATH).Handler(staticWithoutDirectoryListingHandler())
	router.PathPrefix("/" + ASSETS_DIR_PATH).Handler(http.FileServer(http.FS(assets)))

	if AUTORELOAD {
		router.PathPrefix("/autoreload/token").HandlerFunc(func(w Response, _ Request) { fmt.Fprintf(w, "%d", STARTUP_TIME.Unix()) })
	}

	var handler http.Handler = router
	for _, v := range middlewares {
		handler = v(handler)
	}

	http.Handle("/", handler)

	srv := &http.Server{
		Handler:      handler,
		Addr:         BIND_ADDRESS,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting server: %s", BIND_ADDRESS)
	log.Fatal(srv.ListenAndServe())
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

func GET(path string, handler HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	router.Methods("GET").Path(path).HandlerFunc(applyMiddlewares(handlerFuncToHttpHandler(handler), middlewares...))
}

func POST(path string, handler HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	router.Methods("POST").Path(path).HandlerFunc(applyMiddlewares(handlerFuncToHttpHandler(handler), middlewares...))
}

func DELETE(path string, handler HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	router.Methods("DELETE").Path(path).HandlerFunc(applyMiddlewares(handlerFuncToHttpHandler(handler), middlewares...))
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
			rel := strings.TrimPrefix(p, "views"+string(os.PathSeparator))
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

// HANDLERS MIDDLEWARES =============================

// First middleware gets executed first
func applyMiddlewares(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// SERVER MIDDLEWARES ==============================

func staticWithoutDirectoryListingHandler() http.Handler {
	dir := http.Dir(STATIC_DIR_PATH)
	server := http.FileServer(dir)
	handler := http.StripPrefix("/public", server)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})
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

// HELPERS FUNCTIONS ======================

func HELPER(name string, f interface{}) {
	if _, ok := helpers[name]; ok {
		log.Fatalf("Helper: %s has been defined already", name)
	}

	helpers[name] = f
}

func atoi32(s string) int32 {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int32(i)
}

func atoi64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// VALIDATION ============================

type ValidationErrors map[string][]error

func (v ValidationErrors) Add(field string, err error) {
	v[field] = append(v[field], err)
}

func ValidateStringPresent(val, key, label string, ve ValidationErrors) {
	if len(strings.TrimSpace(val)) == 0 {
		ve.Add(key, fmt.Errorf("%s can't be empty", label))
	}
}

func ValidateStringLength(val, key, label string, ve ValidationErrors, min, max int) {
	l := len(strings.TrimSpace(val))
	if l < min || l > max {
		ve.Add(key, fmt.Errorf("%s has to be between %d and %d characters, length is %d", label, min, max, l))
	}
}

func ValidateStringNumeric(val, key, label string, ve ValidationErrors) {
	for _, c := range val {
		if !strings.ContainsRune("0123456789", c) {
			ve.Add(key, fmt.Errorf("%s has to consist of numbers", label))
			return
		}
	}
}

func ValidateISBN13(val, key, label string, ve ValidationErrors) {
	if len(val) != 13 {
		ve.Add(key, fmt.Errorf("%s has to be 13 digits", label))
		return
	}

	sum := 0
	for i, s := range val {
		digit, _ := strconv.Atoi(string(s))
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}

	if sum%10 != 0 {
		ve.Add(key, fmt.Errorf("%s is not a valid ISBN13 number", label))
	}
}

func ValidateImage(val io.Reader, key, label string, ve ValidationErrors, maxw, maxh int) {
	if val == nil {
		return
	}

	image, _, err := image.Decode(val)
	if err != nil {
		ve.Add(key, fmt.Errorf("%s has an unsupported format supported formats are JPG, GIF, PNG", label))
		return
	}

	sz := image.Bounds().Size()
	if sz.X > maxw {
		ve.Add(key, fmt.Errorf("%s width should be less than %d px", label, maxw))
	}
	if sz.Y > maxh {
		ve.Add(key, fmt.Errorf("%s height should be less than %d px", label, maxh))
	}
}

func ValidateInt32Min(val int32, key, label string, ve ValidationErrors, min int32) {
	if val < min {
		ve.Add(key, fmt.Errorf("%s shouldn't be less than %d", label, min))
	}
}
