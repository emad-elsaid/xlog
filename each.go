package xlog

import (
	"context"
	"regexp"
)

// IgnorePath Register a pattern to be ignored when walking directories.
func IgnorePath(r *regexp.Regexp) {
	app := GetApp()
	app.IgnorePath(r)
}

// IsIgnoredPath checks if a file path should be ignored according to the list
// of ignored paths. page source implementations can use it to ignore files from
// their sources
func IsIgnoredPath(d string) bool {
	app := GetApp()
	return app.IsIgnoredPath(d)
}

func Pages(ctx context.Context) []Page {
	app := GetApp()
	return app.Pages(ctx)
}

// EachPage iterates on all available pages. many extensions
// uses it to get all pages and maybe parse them and extract needed information
func EachPage(ctx context.Context, f func(Page)) {
	app := GetApp()
	app.EachPage(ctx, f)
}

// MapPage Similar to EachPage but iterates concurrently and accumulates
// returns in a slice
func MapPage[T any](ctx context.Context, f func(Page) T) []T {
	app := GetApp()
	return MapPageGeneric(app, ctx, f)
}
