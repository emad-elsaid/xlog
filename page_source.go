package xlog

import (
	"context"
	"errors"
	"io/fs"
	"path"
	"path/filepath"
)

type PageSource interface {
	// Page takes a page name and return a Page struct
	Page(string) Page
	// Each iterates over all pages in the source
	Each(context.Context, func(Page))
}

var sources = []PageSource{
	&markdownCWDFS{},
}

func NewPage(name string) (p Page) {
	for i := range sources {
		p = sources[i].Page(name)
		if p != nil && p.Exists() {
			return
		}
	}

	return
}

func RegisterPageSource(p PageSource) {
	sources = append([]PageSource{p}, sources...)
}

// MarkdownCWDFS a current directory markdown pages
type markdownCWDFS struct{}

// NewPage Creates an instance of Page with name. if no name is passed it's assumed INDEX
func (m *markdownCWDFS) Page(name string) Page {
	if name == "" {
		name = INDEX
	}

	return &page{
		name: name,
	}
}

func (m *markdownCWDFS) Each(ctx context.Context, f func(Page)) {
	filepath.WalkDir(".", func(name string, d fs.DirEntry, err error) error {
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
