package xlog

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/emad-elsaid/memoize"
	"github.com/emad-elsaid/memoize/cache/adapters/hashicorp"
	"github.com/fsnotify/fsnotify"
	lru "github.com/hashicorp/golang-lru/v2"
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
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				slog.Error("Can't watch files for change", "error", err)
			}
			defer watcher.Close()

			err = watcher.Add(m.path)
			if err != nil {
				slog.Error("Can't add Markdown FS path", "error", err)
			}

			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						slog.Error("Filesystem watcher returned not OK from Events")
						return
					}

					if event.Has(fsnotify.Write) || event.Has(fsnotify.Remove) {
						if !strings.HasSuffix(event.Name, ".md") {
							continue
						}

						name := strings.TrimSuffix(event.Name, ".md")
						name, _ = filepath.Rel(m.path, name)
						cp := m._page(name)
						Trigger(Changed, cp)

						m.cache.Remove(name)
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						slog.Error("Filesystem watcher returned not OK from Errors")
						return
					}

					slog.Error("Filesystem watcher error", "error", err)
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
