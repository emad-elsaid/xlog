package xlog

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"log"
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

func buildStaticSite(dest string) error {
	srv := server()

	// building Index separately
	err := buildRoute(
		srv,
		"/"+INDEX,
		dest,
		path.Join(dest, "index.html"),
	)

	if err != nil {
		log.Printf("Index Page may not exist, make sure your Index Page exists, err: %s", err.Error())
	}

	EachPageCon(context.Background(), func(p Page) {
		err := buildRoute(
			srv,
			"/"+p.Name(),
			path.Join(dest, p.Name()),
			path.Join(dest, p.Name(), "index.html"),
		)

		if err != nil {
			log.Printf("error while processing: %s, err: %s", p.Name(), err.Error())
		}
	})

	// If we render 404 page
	// Copy 404 page from dest/404/index.html to /dest/404.html
	if in, err := os.Open(path.Join(dest, NOT_FOUND_PAGE, "index.html")); err == nil {
		defer in.Close()
		out, err := os.Create(path.Join(dest, "404.html"))
		if err != nil {
			log.Printf("error while opening dest/404.html, err: %s", err.Error())
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
			log.Printf("error while processing: %s, err: %s", route, err.Error())
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
			log.Printf("error while processing: %s, err: %s", route, err.Error())
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
			log.Printf("Asset %s already exists", destPath)
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
