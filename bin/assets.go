// This script is used to compile assets to local directory
package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var CSS_URLS = []string{
	"https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.css",
	"https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.css",
}

var JS_URLS = []string{
	"https://cdn.jsdelivr.net/simplemde/latest/simplemde.min.js",
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
