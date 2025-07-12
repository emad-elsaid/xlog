package xlog

import (
	"context"
	"reflect"
	"regexp"

	"golang.org/x/sync/errgroup"
)

// IgnorePath registers a pattern to be ignored when walking directories
func (app *App) IgnorePath(r *regexp.Regexp) {
	app.ignoredPaths = append(app.ignoredPaths, r)
}

// IsIgnoredPath checks if a file path should be ignored
func (app *App) IsIgnoredPath(d string) bool {

	for _, v := range app.ignoredPaths {
		if v.MatchString(d) {
			return true
		}
	}
	return false
}

// Pages returns all pages
func (app *App) Pages(ctx context.Context) []Page {
	if app.pages == nil {
		app.populatePagesCache(ctx)
	}
	pages := app.pages
	return pages[:]
}

// EachPage iterates on all available pages
func (app *App) EachPage(ctx context.Context, f func(Page)) {
	if app.pages == nil {
		app.populatePagesCache(ctx)
	}
	currentPages := app.pages

	for _, p := range currentPages {
		select {
		case <-ctx.Done():
			return
		default:
			f(p)
		}
	}
}

// MapPage maps over all pages
func (app *App) MapPage(ctx context.Context, f func(Page) error) []error {
	sources := app.sources

	var errs []error
	for _, source := range sources {
		source.Each(ctx, func(p Page) {
			if err := f(p); err != nil {
				errs = append(errs, err)
			}
		})
	}
	return errs
}

// MapPage maps over all pages with generic return type
func MapPage[T any](app *App, ctx context.Context, f func(Page) T) []T {
	if app.pages == nil {
		app.populatePagesCache(ctx)
	}
	pages := app.pages
	concurrency := app.concurrency

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

// populatePagesCache populates the pages cache
func (app *App) populatePagesCache(ctx context.Context) {

	app.pages = make([]Page, 0, 1000)
	for _, s := range app.sources {
		select {
		case <-ctx.Done():
			return
		default:
			s.Each(ctx, func(p Page) {
				app.pages = append(app.pages, p)
			})
		}
	}
}

// isNil checks if a value is nil
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
