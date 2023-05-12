package xlog

import (
	"context"
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
		populatePagesCache(ctx)
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

func populatePagesCache(ctx context.Context) {
	pages = []Page{}
	for _, s := range sources {
		select {
		case <-ctx.Done():
			return
		default:
			s.Each(ctx, func(p Page) {
				pages = append(pages, p)
			})
		}
	}
}
