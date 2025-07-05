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
)

// build builds the static site
func (app *App) build(buildDir string) error {
	srv := app.server()

	// building Index separately
	err := app.buildRoute(
		srv,
		"/"+app.config.Index,
		buildDir,
		path.Join(buildDir, "index.html"),
	)

	if err != nil {
		slog.Error("Index Page may not exist, make sure your Index Page exists", "index", app.config.Index, "error", err)
	}

	errs := app.MapPage(context.Background(), func(p Page) error {
		err := app.buildRoute(
			srv,
			"/"+p.Name(),
			path.Join(buildDir, p.Name()),
			path.Join(buildDir, p.Name(), "index.html"),
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
	if in, err := os.Open(path.Join(buildDir, app.config.NotFoundPage, "index.html")); err == nil {
		defer in.Close()
		out, err := os.Create(path.Join(buildDir, "404.html"))
		if err != nil {
			slog.Error("Failed to open dest/404.html", "error", err)
		}
		defer out.Close()
		io.Copy(out, in)
	}

	extensionPageEnclosed := app.extensionPageEnclosed
	extensionPage := app.extensionPage
	buildPerms := app.buildPerms

	for route := range extensionPageEnclosed {
		err := app.buildRoute(
			srv,
			route,
			path.Join(buildDir, route),
			path.Join(buildDir, route, "index.html"),
		)

		if err != nil {
			slog.Error("Failed to process extension page", "route", route, "error", err)
		}
	}

	for route := range extensionPage {
		err := app.buildRoute(
			srv,
			route,
			path.Join(buildDir, path.Dir(route)),
			path.Join(buildDir, route),
		)

		if err != nil {
			slog.Error("Failed to process extension page", "route", route, "error", err)
		}
	}

	return fs.WalkDir(assets, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		destPath := path.Join(buildDir, p)

		if entry.IsDir() {
			if err := os.MkdirAll(destPath, buildPerms); err != nil {
				return err
			}
		} else if _, err := os.Stat(destPath); err == nil {
			slog.Warn("Asset file already exists", "path", destPath)
		} else {
			content, err := fs.ReadFile(assets, p)
			if err != nil {
				return err
			}

			if err := os.WriteFile(destPath, content, buildPerms); err != nil {
				return err
			}
		}

		return nil
	})
}

// buildRoute builds a single route
func (app *App) buildRoute(srv *http.Server, route, dir, file string) error {
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		return err
	}

	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if err := os.MkdirAll(dir, app.buildPerms); err != nil {
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

	return os.WriteFile(file, body, app.buildPerms)
}

// RegisterBuildPage registers a build page
func (app *App) RegisterBuildPage(p string, encloseInDir bool) {

	if encloseInDir {
		app.extensionPageEnclosed[p] = true
	} else {
		app.extensionPage[p] = true
	}
}
