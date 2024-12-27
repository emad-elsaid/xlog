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

	// building Index separately
	err := buildRoute(
		srv,
		"/"+Config.Index,
		dest,
		path.Join(dest, "index.html"),
	)

	if err != nil {
		slog.Error("Index Page may not exist, make sure your Index Page exists", "index", Config.Index, "error", err)
	}

	errs := MapPage(context.Background(), func(p Page) error {
		err := buildRoute(
			srv,
			"/"+p.Name(),
			path.Join(dest, p.Name()),
			path.Join(dest, p.Name(), "index.html"),
		)

		if err != nil {
			return fmt.Errorf("Failed to process page: %s, error: %w", p.Name(), err)
		}

		return nil
	})

	if err := errors.Join(errs...); err != nil {
		slog.Error(err.Error())
	}

	// If we render 404 page
	// Copy 404 page from dest/404/index.html to /dest/404.html
	if in, err := os.Open(path.Join(dest, Config.NotFoundPage, "index.html")); err == nil {
		defer in.Close()
		out, err := os.Create(path.Join(dest, "404.html"))
		if err != nil {
			slog.Error("Failed to open dest/404.html", "error", err)
		}
		defer out.Close()
		io.Copy(out, in)
	}

	extension_page_enclosed.Range(func(route string, _ bool) bool {
		err := buildRoute(
			srv,
			route,
			path.Join(dest, route),
			path.Join(dest, route, "index.html"),
		)

		if err != nil {
			slog.Error("Failed to process extension page", "route", route, "error", err)
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
			slog.Error("Failed to process extension page", "route", route, "error", err)
		}

		return true
	})

	return fs.WalkDir(assets, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		destPath := path.Join(dest, p)

		if entry.IsDir() {
			if err := os.MkdirAll(destPath, build_perms); err != nil {
				return err
			}
		} else if _, err := os.Stat(destPath); err == nil {
			slog.Warn("Asset file already exists", "path", destPath)
		} else {
			content, err := fs.ReadFile(assets, p)
			if err != nil {
				return err
			}

			if err := os.WriteFile(destPath, content, build_perms); err != nil {
				return err
			}
		}

		return nil
	})
}

func buildRoute(srv *http.Server, route, dir, file string) error {
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		return err
	}

	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if err := os.MkdirAll(dir, build_perms); err != nil {
		return err
	}

	if rec.Result().StatusCode != http.StatusOK {
		return errors.New(rec.Result().Status)
	}

	body, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		return err
	}
	defer rec.Result().Body.Close()

	return os.WriteFile(file, body, build_perms)
}
