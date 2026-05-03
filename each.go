package xlog

import (
	"context"
	"reflect"
	"regexp"
	"runtime"
	"sync"

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

var (
	pages      []Page
	pagesMutex sync.RWMutex
)

func Pages(ctx context.Context) []Page {
	pagesMutex.RLock()
	cached := pages
	pagesMutex.RUnlock()

	if cached == nil {
		populatePagesCache(ctx)
		pagesMutex.RLock()
		cached = pages
		pagesMutex.RUnlock()
	}

	return cached[:]
}

// EachPage iterates on all available pages. many extensions
// uses it to get all pages and maybe parse them and extract needed information
func EachPage(ctx context.Context, f func(Page)) {
	pagesMutex.RLock()
	cached := pages
	pagesMutex.RUnlock()

	if cached == nil {
		populatePagesCache(ctx)
		pagesMutex.RLock()
		cached = pages
		pagesMutex.RUnlock()
	}

	for _, p := range cached {
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
	pagesMutex.RLock()
	cached := pages
	pagesMutex.RUnlock()

	if cached == nil {
		populatePagesCache(ctx)
		pagesMutex.RLock()
		cached = pages
		pagesMutex.RUnlock()
	}

	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(concurrency)

	output := make([]T, 0, len(cached))
	ch := make(chan T, concurrency)
	done := make(chan bool)

	go func() {
		for v := range ch {
			output = append(output, v)
		}
		done <- true
	}()

	for _, p := range cached {
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
	pagesMutex.Lock()
	pages = nil
	pagesMutex.Unlock()
	return nil
}

func populatePagesCache(ctx context.Context) {
	pagesMutex.Lock()
	defer pagesMutex.Unlock()

	// Double-check after acquiring lock
	if pages != nil {
		return
	}

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
