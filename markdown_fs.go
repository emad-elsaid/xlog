package xlog

import (
	"context"
	"errors"
	"io/fs"
	"log"
	"log/slog"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/emad-elsaid/memoize"
	"github.com/emad-elsaid/memoize/cache/adapters/hashicorp"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/rjeczalik/notify"
)

func newMarkdownFS(p string) *markdownFS {
	cache, err := lru.New[string, Page](1000)
	if err != nil {
		slog.Error("Can't create cache for pages", "error", err)
		panic("Can't continue without cache instance")
	}

	m := markdownFS{
		cache: cache,
		path:  p,
	}

	m._page = memoize.NewWithCache(
		hashicorp.LRU(cache),
		func(name string) Page {
			if name == "" {
				name = INDEX
			}

			return &page{
				name: name,
			}
		},
	)

	m.watch = sync.OnceFunc(func() {
		go func() {
			events := make(chan notify.EventInfo, 1)

			absPath, err := filepath.Abs(m.path)
			if err != nil {
				log.Fatal(err)
			}

			if err := notify.Watch(m.path+"/...", events, notify.All); err != nil {
				slog.Error("Can't watch files for change", "error", err)
			}
			defer notify.Stop(events)

			for {
				switch ei := <-events; ei.Event() {
				case notify.Write, notify.Remove, notify.Rename:
					relPath, err := filepath.Rel(absPath, ei.Path())
					if err != nil {
						slog.Error("Can't resolve relative path", "error", err)
						continue
					}

					if !strings.HasSuffix(relPath, ".md") {
						continue
					}

					name := strings.TrimSuffix(relPath, ".md")
					cp := m._page(name)
					Trigger(Changed, cp)

					m.cache.Remove(name)
				}
			}
		}()
	})

	return &m
}

// MarkdownFS a current directory markdown pages
type markdownFS struct {
	path  string
	cache *lru.Cache[string, Page]
	_page func(string) Page
	watch func()
}

// Page Creates an instance of Page with name. if no name is passed it's assumed INDEX
func (m *markdownFS) Page(name string) Page {
	m.watch()

	return m._page(name)
}

func (m *markdownFS) Each(ctx context.Context, f func(Page)) {
	filepath.WalkDir(m.path, func(name string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			for _, v := range ignoredDirs {
				if v.MatchString(name) {
					return fs.SkipDir
				}
			}

			return nil
		}

		select {

		case <-ctx.Done():
			return errors.New("context stopped")

		default:
			ext := path.Ext(name)
			basename := name[:len(name)-len(ext)]

			if ext == ".md" {
				f(m.Page(basename))
			}

		}

		return nil
	})
}
