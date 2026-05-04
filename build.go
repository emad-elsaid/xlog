// Package xlog provides static site generation capabilities for digital gardening.
// The build functionality converts markdown pages into a static HTML site.
package xlog

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path"

	"github.com/emad-elsaid/types"
)

var extension_page = types.Map[string, bool]{}
var extension_page_enclosed = types.Map[string, bool]{}
var build_perms fs.FileMode = 0744

// RegisterBuildPage registers a path of a page to export when building static version
// of the knowledgebase. encloseInDir will write the output to p/index.html
// instead instead of writing to p directly. that can help have paths with no
// .html extension to be served with the exact name.
func RegisterBuildPage(p string, encloseInDir bool) {
	if encloseInDir {
		extension_page_enclosed.Store(p, true)
	} else {
		extension_page.Store(p, true)
	}
}

func build(dest string) error {
	srv := server()

	buildIndexPage(srv, dest)
	buildAllPages(srv, dest)
	copy404Page(dest)
	buildExtensionPages(srv, dest)
	return copyAssets(dest)
}

func buildIndexPage(srv *http.Server, dest string) {
	err := buildRoute(
		srv,
		"/"+Config.Index,
		dest,
		path.Join(dest, "index.html"),
	)

	if err != nil {
		slog.Error("Index Page may not exist, make sure your Index Page exists", "index", Config.Index, "error", err)
	}
}

func buildAllPages(srv *http.Server, dest string) {
	errs := MapPage(context.Background(), func(p Page) error {
		err := buildRoute(
			srv,
			"/"+p.Name(),
			path.Join(dest, p.Name()),
			path.Join(dest, p.Name(), "index.html"),
		)

		if err != nil {
			return fmt.Errorf("failed to process page: %s, error: %w", p.Name(), err)
		}

		return nil
	})

	if err := errors.Join(errs...); err != nil {
		slog.Error(err.Error())
	}
}

func copy404Page(dest string) {
	// If we render 404 page
	// Copy 404 page from dest/404/index.html to /dest/404.html
	// #nosec G304 -- File path constructed from Config.NotFoundPage (controlled config) and dest (build output dir)
	in, err := os.Open(path.Join(dest, Config.NotFoundPage, "index.html"))
	if err != nil {
		return
	}
	defer func() {
		if closeErr := in.Close(); closeErr != nil {
			slog.Error("Failed to close input 404 file", "error", closeErr)
		}
	}()
	// #nosec G304 -- File path is controlled build output destination, not user input
	out, err := os.Create(path.Join(dest, "404.html"))
	if err != nil {
		slog.Error("Failed to open dest/404.html", "error", err)
		return
	}
	defer func() {
		if err := out.Close(); err != nil {
			slog.Error("Failed to close output 404 file", "error", err)
		}
	}()
	if _, err := io.Copy(out, in); err != nil {
		slog.Error("Failed to copy 404 file", "error", err)
	}
}

func buildExtensionPages(srv *http.Server, dest string) {
	extension_page_enclosed.Range(func(route string, _ bool) bool {
		err := buildRoute(
			srv,
			route,
			path.Join(dest, route),
			path.Join(dest, route, "index.html"),
		)

		if err != nil {
			slog.Error("Failed to build extension page",
				"route", route,
				"error", err,
				"hint", "check if the route handler is registered and returns valid HTML")
		}

		return true
	})

	extension_page.Range(func(route string, _ bool) bool {
		err := buildRoute(
			srv,
			route,
			path.Join(dest, path.Dir(route)),
			path.Join(dest, route),
		)

		if err != nil {
			slog.Error("Failed to build extension page",
				"route", route,
				"error", err,
				"hint", "check if the route handler is registered and returns valid content")
		}

		return true
	})
}

func copyAssets(dest string) error {
	return fs.WalkDir(assets, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk asset path %q: %w", p, err)
		}

		destPath := path.Join(dest, p)

		if entry.IsDir() {
			if mkdirErr := os.MkdirAll(destPath, build_perms); mkdirErr != nil {
				return fmt.Errorf("failed to create asset directory %q: %w", destPath, mkdirErr)
			}
			return nil
		}

		if _, statErr := os.Stat(destPath); statErr == nil {
			slog.Warn("Asset file already exists, skipping", "path", destPath)
			return nil
		}

		content, err := fs.ReadFile(assets, p)
		if err != nil {
			return fmt.Errorf("failed to read embedded asset %q: %w", p, err)
		}

		if err := os.WriteFile(destPath, content, build_perms); err != nil {
			return fmt.Errorf("failed to write asset file %q: %w", destPath, err)
		}
		return nil
	})
}

func buildRoute(srv *http.Server, route, dir, file string) error {
	// #nosec G704 -- Route is internal, constructed from page names in build process, not SSRF
	req, err := http.NewRequest(http.MethodGet, route, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request for route %q: %w", route, err)
	}

	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	// #nosec G703 -- Dir is controlled build output directory, no path traversal risk
	if mkdirErr := os.MkdirAll(dir, build_perms); mkdirErr != nil {
		return fmt.Errorf("failed to create directory %q: %w", dir, mkdirErr)
	}

	if rec.Result().StatusCode != http.StatusOK {
		return fmt.Errorf("route %q returned HTTP %s (expected 200 OK)", route, rec.Result().Status)
	}

	body, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		return fmt.Errorf("failed to read response body for route %q: %w", route, err)
	}
	defer func() {
		if err := rec.Result().Body.Close(); err != nil {
			slog.Error("Failed to close response body", "error", err)
		}
	}()

	// #nosec G703 -- File is controlled build output destination, no path traversal risk
	if err := os.WriteFile(file, body, build_perms); err != nil {
		return fmt.Errorf("failed to write output file %q: %w", file, err)
	}
	return nil
}
