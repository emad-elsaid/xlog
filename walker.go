package xlog

import (
	"context"
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"regexp"
	"sync"
)

func init() {
	Listen(AfterWrite, clearWalkPagesCache)
	Listen(AfterDelete, clearWalkPagesCache)
}

// a List of directories that should be ignored by directory walking function.
// for example the versioning extension can register `.versions` directory to be
// ignored
var ignoredDirs = []*regexp.Regexp{
	regexp.MustCompile(`\..+`),
	regexp.MustCompile(PUBLIC_PATH),
}

// Register a pattern to be ignored when walking directories.
func IgnoreDir(r *regexp.Regexp) {
	ignoredDirs = append(ignoredDirs, r)
}

var walkPagesCache []*Page
var walkPagesCacheMutex sync.RWMutex

// WalkPages iterates on all available pages. many extensions
// uses it to get all pages and maybe parse them and extract needed information
func WalkPages(ctx context.Context, f func(*Page)) {
	if walkPagesCache == nil {
		populateWalkPagesCache(ctx)
	}

	walkPagesCacheMutex.RLock()
	defer walkPagesCacheMutex.RUnlock()

	for _, p := range walkPagesCache {
		select {
		case <-ctx.Done():
			return
		default:
			f(p)
		}
	}
}

func clearWalkPagesCache(_ *Page) (err error) {
	walkPagesCacheMutex.Lock()
	defer walkPagesCacheMutex.Unlock()

	walkPagesCache = nil
	return nil
}

func populateWalkPagesCache(ctx context.Context) {
	walkPagesCacheMutex.Lock()
	defer walkPagesCacheMutex.Unlock()

	walkPagesCache = []*Page{}

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
			return errors.New("Context stopped")

		default:
			ext := path.Ext(name)
			basename := name[:len(name)-len(ext)]

			if ext == ".md" {
				walkPagesCache = append(walkPagesCache, &Page{Name: basename})
			}

		}

		return nil
	})
}
