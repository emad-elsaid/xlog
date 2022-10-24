// This script is used to compile assets to local directory
package main

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const DEST = "public"

var CSS_DEST = path.Join(DEST, "style.css")
var JS_DEST = path.Join(DEST, "script.js")

//go:embed custom.css
var CUSTOM_CSS []byte

var CSS_URLS = []string{
	"https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.css",
	"https://cdn.jsdelivr.net/npm/bulma-prefers-dark@0.1.0-beta.1/css/bulma-prefers-dark.css",
	"https://cdn.jsdelivr.net/npm/codemirror@5.65.4/lib/codemirror.css",
	"https://codemirror.net/5/addon/hint/show-hint.css",
}

type zipURL = string
type zipPath = string

var CSS_ZIP = map[zipURL]map[zipPath]string{
	"https://use.fontawesome.com/releases/v6.1.1/fontawesome-free-6.1.1-web.zip": {
		"fontawesome-free-6.1.1-web/css/all.min.css": CSS_DEST,
		"fontawesome-free-6.1.1-web/webfonts/":       path.Join(DEST, "webfonts"),
	},
}

var JS_URLS = []string{
	"https://cdn.jsdelivr.net/npm/codemirror@5.65.4/lib/codemirror.min.js",
	"https://cdn.jsdelivr.net/npm/codemirror@5.65.4/addon/mode/overlay.js",
	"https://cdn.jsdelivr.net/npm/codemirror@5.65.4/mode/markdown/markdown.js",
	"https://cdn.jsdelivr.net/npm/codemirror@5.65.4/mode/xml/xml.js",
	"https://cdn.jsdelivr.net/npm/codemirror@5.65.4/mode/gfm/gfm.js",
	"https://cdn.jsdelivr.net/npm/codemirror@5.65.4/mode/javascript/javascript.js",
	"https://cdn.jsdelivr.net/npm/codemirror@5.65.4/mode/go/go.js",
	"https://codemirror.net/5/addon/hint/show-hint.js",
}

func main() {
	// ensure DEST exists
	if _, err := os.Stat(DEST); err != nil {
		err := os.Mkdir(DEST, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := urlsToFile(CSS_URLS, CSS_DEST)
	if err != nil {
		log.Fatal(err)
	}

	err = urlsToFile(JS_URLS, JS_DEST)
	if err != nil {
		log.Fatal(err)
	}

	for url, files := range CSS_ZIP {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		var buf bytes.Buffer
		io.Copy(&buf, resp.Body)
		resp.Body.Close()

		z, err := zip.NewReader(bytes.NewReader(buf.Bytes()), resp.ContentLength)
		if err != nil {
			log.Fatal(err)
		}

		for _, zf := range z.File {
			for f, d := range files {
				if !strings.HasPrefix(zf.Name, f) {
					continue
				}

				dpath := path.Join(d, zf.Name[len(f):])
				log.Println("Extracting to", dpath)

				if _, err := os.Stat(path.Dir(dpath)); err != nil {
					log.Println("checking dir: ", path.Dir(dpath))
					os.Mkdir(path.Dir(dpath), 0744)
				}

				if zf.FileInfo().IsDir() {
					os.Mkdir(dpath, 0744)
					continue
				}

				var flags int
				if strings.HasSuffix(zf.Name, ".css") {
					flags = os.O_APPEND | os.O_WRONLY | os.O_CREATE
				} else {
					flags = os.O_CREATE | os.O_WRONLY
				}

				dest, err := os.OpenFile(dpath, flags, 0744)
				if err != nil {
					log.Fatal("Opening the destination file ", err)
				}

				b, err := zf.Open()
				if err != nil {
					log.Fatal(err)
				}

				content, err := io.ReadAll(b)
				b.Close()
				if err != nil {
					log.Fatal("reading all failed: ", zf, err)
				}

				// this strips all ../ references to reference files in assets instead of parent dir
				// fontawesome does that which forced me to have the css in a separate file under assets/fontawesome for a while
				// I wanted to have just one css and one js in the root of `assets/` that the solution for it
				if strings.HasSuffix(zf.Name, ".css") {
					replaced_content := strings.NewReplacer("../", "").Replace(string(content))
					content = []byte(replaced_content)
				}

				dest.Write(content)
				dest.Close()
			}
		}

		f, err := os.OpenFile(CSS_DEST, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0744)
		if err != nil {
			log.Fatal(err)
		}
		f.Write(CUSTOM_CSS)
		f.Close()

		if err = mergeLines(CSS_DEST); err != nil {
			log.Fatal(err)
		}

	}
}

func urlsToFile(urls []string, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, v := range urls {
		log.Printf("Downloading: %s", v)

		resp, err := http.Get(v)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, resp.Body)
		resp.Body.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func mergeLines(f string) error {
	c, err := os.ReadFile(f)
	if err != nil {
		return err
	}

	c = bytes.ReplaceAll(c, []byte("\n"), []byte(""))
	return os.WriteFile(f, c, 0644)
}
