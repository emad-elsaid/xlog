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
)

var extension_page = map[string]bool{}
var build_perms fs.FileMode = 0744

// EXTENSION_PAGE announces a path of a page as an extension page, enables the
// page to be exported when building static version of the knowledgebase
func EXTENSION_PAGE(p string) {
	extension_page[p] = true
}

func buildStaticSite(dest string) error {
	srv := server()

	// building Index separately
	err := buildRoute(srv, "/"+INDEX, dest, path.Join(dest, "index.html"))
	if err != nil {
		log.Printf("error while processing root path, err: %s", err.Error())
	}

	WalkPages(context.Background(), func(p *Page) {
		route := "/" + p.Name
		dir := path.Join(dest, p.Name)
		file := path.Join(dest, p.Name, "index.html")

		err := buildRoute(srv, route, dir, file)
		if err != nil {
			log.Printf("error while processing: %s, err: %s", p.Name, err.Error())
		}
	})

	for route := range extension_page {
		dir := path.Join(dest, route)
		file := path.Join(dest, route, "index.html")

		err := buildRoute(srv, route, dir, file)
		if err != nil {
			log.Printf("error while processing: %s, err: %s", route, err.Error())
			continue
		}
	}

	fs.WalkDir(assets, ".", func(p string, entry fs.DirEntry, err error) error {
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

	return nil
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
