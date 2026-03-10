package xlog

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
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

// warmASTCache pre-warms AST cache in background for better first-request performance
func warmASTCache() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Cache warming panicked", "error", r)
			}
		}()

		slog.Info("Starting AST cache warming in background")
		start := time.Now()

		ctx := context.Background()
		pages := Pages(ctx)
		total := len(pages)

		// Trigger BeforeCacheWarming event - extensions can prepare their caches
		// For example: autolink_pages extension builds its trie here
		Trigger(BeforeCacheWarming, nil)

		// Use aggressive concurrency - we're I/O bound (file reads)
		concurrency := runtime.NumCPU() * 10
		if concurrency > 400 {
			concurrency = 400
		}
		slog.Info("Cache warming configuration",
			"pages", total,
			"concurrency", concurrency,
			"cpus", runtime.NumCPU())

		// PHASE 1: Pre-load all file contents into memory
		// This batches file I/O to reduce syscall overhead
		slog.Info("Phase 1: Pre-loading file contents into memory")
		preloadStart := time.Now()

		type pageContent struct {
			page    Page
			content Markdown
			modtime time.Time
		}

		contentsCh := make(chan pageContent, concurrency*2)
		pageInputCh := make(chan Page, concurrency*2)
		var preloadWg sync.WaitGroup

		// Workers for reading AND preprocessing files
		for i := 0; i < concurrency; i++ {
			preloadWg.Add(1)
			go func() {
				defer preloadWg.Done()
				for p := range pageInputCh {
					stat, err := os.Stat(p.FileName())
					if err != nil {
						slog.Error("Failed to stat file", "page", p.Name(), "error", err)
						continue
					}

					content, err := os.ReadFile(p.FileName())
					if err != nil {
						slog.Error("Failed to pre-load file", "page", p.Name(), "error", err)
						continue
					}

					contentsCh <- pageContent{page: p, content: Markdown(content), modtime: stat.ModTime()}
				}
			}()
		}

		// Second set of workers for preprocessing (CPU intensive)
		var preprocessWg sync.WaitGroup
		for i := 0; i < concurrency; i++ {
			preprocessWg.Add(1)
			go func() {
				defer preprocessWg.Done()
				for pc := range contentsCh {
					// Type assert to access internal method
					if p, ok := pc.page.(*page); ok {
						p.loadContent(pc.content, pc.modtime)
					}
				}
			}()
		}

		// Feed pages to file reading workers
		go func() {
			for _, p := range pages {
				pageInputCh <- p
			}
			close(pageInputCh)
		}()

		// Wait for file readers, then close preprocessing channel
		preloadWg.Wait()
		close(contentsCh)
		preprocessWg.Wait()

		preloadElapsed := time.Since(preloadStart)
		slog.Info("Phase 1 complete: Files pre-loaded",
			"duration", preloadElapsed,
			"pages", total,
			"pages_per_sec", fmt.Sprintf("%.0f", float64(total)/preloadElapsed.Seconds()))

		// PHASE 2: Parse AST (now with cached content, no file I/O)
		slog.Info("Phase 2: Parsing AST from cached content")
		parseStart := time.Now()

		// Worker pool pattern
		pagesCh := make(chan Page, concurrency*2)
		var wg sync.WaitGroup

		// Atomic counter for lock-free progress tracking
		var completed atomic.Int64
		progressTicker := time.NewTicker(2 * time.Second)
		defer progressTicker.Stop()

		// Progress reporter goroutine
		go func() {
			for range progressTicker.C {
				count := completed.Load()
				if count > 0 && count <= int64(total) {
					slog.Info("Phase 2 progress: AST parsing",
						"completed", count,
						"total", total,
						"percent", fmt.Sprintf("%.1f%%", float64(count)/float64(total)*100))
				}
				if count >= int64(total) {
					return
				}
			}
		}()

		// Start worker goroutines
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for page := range pagesCh {
					func() {
						defer func() {
							if r := recover(); r != nil {
								slog.Error("Page AST parsing panicked", "page", page.Name(), "error", r)
							}
						}()

						// Parse AST (content already cached, no file I/O)
						_, _ = page.AST()

						// Increment completed count (lock-free)
						completed.Add(1)
					}()
				}
			}()
		}

		// Feed pages to workers
		for _, p := range pages {
			pagesCh <- p
		}
		close(pagesCh)

		// Wait for all workers to complete
		wg.Wait()

		parseElapsed := time.Since(parseStart)
		slog.Info("Phase 2 complete: AST parsing finished",
			"duration", parseElapsed,
			"pages", total,
			"pages_per_sec", fmt.Sprintf("%.0f", float64(total)/parseElapsed.Seconds()))

		elapsed := time.Since(start)
		slog.Info("Cache warming completed",
			"pages", total,
			"duration", elapsed,
			"pages_per_sec", fmt.Sprintf("%.0f", float64(total)/elapsed.Seconds()))

		// Log file read statistics
		readCount := fileReadCount.Load()
		slog.Info("File read statistics",
			"total_reads", readCount,
			"pages", total,
			"reads_per_page", fmt.Sprintf("%.2f", float64(readCount)/float64(total)))
	}()
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
