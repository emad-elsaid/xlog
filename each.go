package xlog

import (
	"context"
	"reflect"
	"regexp"
	"runtime"

	"golang.org/x/sync/errgroup"
)

// a List of directories that should be ignored by directory walking function.
// for example the versioning extension can register `.versions` directory to be
// ignored
var ignoredPaths = []*regexp.Regexp{
	regexp.MustCompile(`^\.`), // Ignore any hidden directory
}

// IgnorePath Register a pattern to be ignored when walking directories.
func IgnorePath(r *regexp.Regexp) {
	ignoredPaths = append(ignoredPaths, r)
}

// IsIgnoredPath checks if a file path should be ignored according to the list
// of ignored paths. page source implementations can use it to ignore files from
// their sources
func IsIgnoredPath(d string) bool {
	for _, v := range ignoredPaths {
		if v.MatchString(d) {
			return true
		}
	}

	return false
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

var concurrency = runtime.NumCPU() * 4

// MapPage Similar to EachPage but iterates concurrently and accumulates
// returns in a slice
func MapPage[T any](ctx context.Context, f func(Page) T) []T {
	if pages == nil {
		populatePagesCache(ctx)
	}

	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(concurrency)

	output := make([]T, 0, len(pages))
	ch := make(chan T, concurrency)
	done := make(chan bool)

	go func() {
		for v := range ch {
			output = append(output, v)
		}
		done <- true
	}()

	for _, p := range pages {
		select {
		case <-ctx.Done():
			break
		default:
			grp.Go(func() (err error) {
				val := f(p)
				if isNil(val) {
					return
				}

				ch <- val

				return
			})
		}
	}

	grp.Wait()
	close(ch)
	<-done

	return output
}

// From https://stackoverflow.com/a/77341451/22401486
func isNil[T any](t T) bool {
	v := reflect.ValueOf(t)
	kind := v.Kind()
	// Must be one of these types to be nillable
	return !v.IsValid() || (kind == reflect.Ptr ||
		kind == reflect.Interface ||
		kind == reflect.Slice ||
		kind == reflect.Map ||
		kind == reflect.Chan ||
		kind == reflect.Func) &&
		v.IsNil()
}

func clearPagesCache(p Page) (err error) {
	pages = nil
	return nil
}

func populatePagesCache(ctx context.Context) {
	pages = make([]Page, 0, 1000)
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
