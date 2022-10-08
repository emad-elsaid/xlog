package main

import (
	"context"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
)

func buildStaticSite(dest string) error {
	srv := server()

	WalkPages(context.Background(), func(p *Page) {
		req, err := http.NewRequest(http.MethodGet, "/"+p.Name, nil)
		if err != nil {
			log.Printf("error while processing: %s, err: %s", p.Name, err.Error())
			return
		}

		rec := httptest.NewRecorder()
		srv.Handler.ServeHTTP(rec, req)

		dir := path.Join(dest, p.Name)
		file := path.Join(dest, p.Name, "index.html")
		if p.Name == "index" {
			dir = dest
			file = path.Join(dest, "index.html")
		}

		if err := os.MkdirAll(dir, 0700); err != nil {
			log.Printf("error while processing: %s, err: %s", p.Name, err.Error())
			return
		}

		if rec.Result().StatusCode != http.StatusOK {
			log.Printf("error while processing: %s, err: %s", p.Name, rec.Result().Status)
			return
		}

		body, err := io.ReadAll(rec.Result().Body)
		if err != nil {
			log.Printf("error while processing: %s, err: %s", p.Name, err.Error())
			return
		}
		defer rec.Result().Body.Close()

		err = os.WriteFile(file, body, 0700)
		if err != nil {
			log.Printf("error while processing: %s, err: %s", p.Name, err.Error())
			return
		}
	})

	fs.WalkDir(assets, ".", func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		destPath := path.Join(dest, p)

		if entry.IsDir() {
			if err := os.MkdirAll(destPath, 0700); err != nil {
				return err
			}
		} else {
			content, err := fs.ReadFile(assets, p)
			if err != nil {
				return err
			}

			if err := os.WriteFile(destPath, content, 0700); err != nil {
				return err
			}
		}

		return nil
	})

	return nil
}
