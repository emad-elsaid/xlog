package xlog

import (
	"context"
	"regexp"
	"runtime"

	"golang.org/x/sync/errgroup"
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

func Pages(ctx context.Context) []Page {
	if pages == nil {
		populatePagesCache(ctx)
	}

	return pages[:]
}

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

// EachPageCon Similar to EachPage but iterates concurrently
func EachPageCon(ctx context.Context, f func(Page)) {
	if pages == nil {
		populatePagesCache(ctx)
	}

	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(runtime.NumCPU() * 4)

	currentPages := pages
	for _, p := range currentPages {
		select {
		case <-ctx.Done():
			break
		default:
			grp.Go(func() (err error) { f(p); return })
		}
	}

	grp.Wait()
}

// TODO check if changing the cache based on the page would make xlog faster
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
