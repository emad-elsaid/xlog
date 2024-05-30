package xlog

import (
	"context"
	"errors"
	"io/fs"
	"path"
	"path/filepath"

	"github.com/emad-elsaid/memoize"
)

func newMarkdownFS(path string) *markdownFS {
	m := markdownFS{path: path}

	m._page = memoize.New(func(name string) Page {
		if name == "" {
			name = INDEX
		}

		return &page{
			name: name,
		}
	})

	return &m
}

// MarkdownCWDFS a current directory markdown pages
type markdownFS struct {
	path  string
	_page func(string) Page
}

// Page Creates an instance of Page with name. if no name is passed it's assumed INDEX
func (m *markdownFS) Page(name string) Page {
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
