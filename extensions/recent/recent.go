package recent

import (
	"embed"
	"html/template"
	"slices"
	"strings"
	"sync"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

var cachedOutput Output
var cachedOutputLock sync.RWMutex

func init() {
	RegisterExtension(Recent{})
}

type Recent struct{}

func (Recent) Name() string { return "recent" }
func (Recent) Init() {
	Get(`/+/recent`, recentHandler)
	RegisterBuildPage("/+/recent", true)
	RegisterTemplate(templates, "templates")
	RegisterLink(func(Page) []Command { return []Command{links{}} })

	// Clear cache when pages change
	if !Config.Readonly {
		Listen(PageChanged, func(Page) error {
			cachedOutputLock.Lock()
			cachedOutput = nil
			cachedOutputLock.Unlock()
			return nil
		})
		Listen(PageDeleted, func(Page) error {
			cachedOutputLock.Lock()
			cachedOutput = nil
			cachedOutputLock.Unlock()
			return nil
		})
	}
}

func recentHandler(r Request) Output {
	// Check cache first
	cachedOutputLock.RLock()
	if cachedOutput != nil {
		output := cachedOutput
		cachedOutputLock.RUnlock()
		return output
	}
	cachedOutputLock.RUnlock()

	// Cache miss - render
	rp := Pages(r.Context())
	slices.SortFunc(rp, func(a, b Page) int {
		if modtime := b.ModTime().Compare(a.ModTime()); modtime != 0 {
			return modtime
		}

		return strings.Compare(a.Name(), b.Name())
	})

	// Pre-warm AST caches in parallel before rendering template
	// This parallelizes the expensive AST parsing that would otherwise
	// happen serially during template rendering
	var wg sync.WaitGroup
	sem := make(chan struct{}, 100) // Limit to 100 concurrent goroutines

	for _, p := range rp {
		wg.Add(1)
		go func(page Page) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire semaphore
			defer func() { <-sem }() // Release semaphore

			// Pre-parse AST in parallel
			// Banner(), Emoji(), Properties() all need AST
			_, _ = page.AST()
		}(p)
	}
	wg.Wait()

	output := Render("recent", Locals{
		"page":  DynamicPage{NameVal: "Recent"},
		"pages": rp,
	})

	// Store in cache
	cachedOutputLock.Lock()
	cachedOutput = output
	cachedOutputLock.Unlock()

	return output
}

type links struct{}

func (l links) Icon() string { return "fa-solid fa-clock-rotate-left" }
func (l links) Name() string { return "Recent" }
func (l links) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"href": "/+/recent",
	}
}
