// This script is used to compile assets to local directory
package main

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var CSS_URLS = []string{
	"https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.css",
	"https://cdn.jsdelivr.net/npm/codemirror@5.65.4/lib/codemirror.css",
}

var CSS_ZIP = map[string]map[string]string{
	"https://use.fontawesome.com/releases/v6.1.1/fontawesome-free-6.1.1-web.zip": {
		"fontawesome-free-6.1.1-web/css/all.min.css": "fontawesome/style.css",
		"fontawesome-free-6.1.1-web/webfonts/":       "webfonts",
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
}

const DEST = "assets"

var CSS_DEST = path.Join(DEST, "style.css")
var JS_DEST = path.Join(DEST, "script.js")

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

		buf := bytes.NewBuffer([]byte{})
		io.Copy(buf, resp.Body)
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

				dpath := path.Join(DEST, d, zf.Name[len(f):])
				log.Println("Extracting to", dpath)

				if _, err := os.Stat(path.Dir(dpath)); err != nil {
					log.Println("checking dir: ", path.Dir(dpath))
					os.Mkdir(path.Dir(dpath), 0700)
				}

				if zf.FileInfo().IsDir() {
					os.Mkdir(dpath, 0700)
					continue
				}

				dest, err := os.Create(dpath)
				if err != nil {
					log.Fatal("Opening the destination file ", err)
				}

				b, err := zf.Open()
				if err != nil {
					log.Fatal(err)
				}

				io.Copy(dest, b)

				dest.Close()
				b.Close()
			}
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
