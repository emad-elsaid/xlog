package xlog

import (
	"context"
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"regexp"
)

func init() {
	Listen(AfterWrite, clearPagesCache)
	Listen(AfterDelete, clearPagesCache)
}

// a List of directories that should be ignored by directory walking function.
// for example the versioning extension can register `.versions` directory to be
// ignored
var ignoredDirs = []*regexp.Regexp{
	regexp.MustCompile(`\..+`), // Ignore any hidden directory
}

// IgnoreDirectory Register a pattern to be ignored when walking directories.
func IgnoreDirectory(r *regexp.Regexp) {
	ignoredDirs = append(ignoredDirs, r)
}

var pages []Page

// EachPage iterates on all available pages. many extensions
// uses it to get all pages and maybe parse them and extract needed information
func EachPage(ctx context.Context, f func(Page)) {
	if pages == nil {
		pages = populatePagesCache(ctx)
	}

	currentPages := pages
	for _, p := range currentPages {
		select {
		case <-ctx.Done():
			return
		default:
			f(p)
		}
	}
}

func clearPagesCache(_ Page) (err error) {
	pages = nil
	return nil
}

func populatePagesCache(ctx context.Context) []Page {
	pages := []Page{}

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
				pages = append(pages, NewPage(basename))
			}

		}

		return nil
	})

	return pages
}
