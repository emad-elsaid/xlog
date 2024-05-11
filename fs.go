package xlog

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
)

// return file that exists in one of the FS structs.
// Prioritizing the end of the slice over earlier FSs.
type priorityFS []fs.FS

func (p priorityFS) Open(name string) (fs.File, error) {
	for i := len(p) - 1; i >= 0; i-- {
		cf := p[i]
		f, err := cf.Open(name)
		if err == nil {
			return f, err
		}
	}

	return nil, fs.ErrNotExist
}

//go:embed public
var assets embed.FS

var staticDirs = []fs.FS{assets}

// RegisterStaticDir adds a filesystem to the filesystems list scanned for files
// when serving static files. can be used to add a directory of CSS or JS files
// by extensions
func RegisterStaticDir(f fs.FS) {
	staticDirs = append(staticDirs, f)
}

func staticHandler(r Request) (Output, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	staticFSs := http.FS(
		priorityFS(
			append(staticDirs, os.DirFS(wd)),
		),
	)

	server := http.FileServer(staticFSs)

	cleanPath := path.Clean(r.URL.Path)

	if f, err := staticFSs.Open(cleanPath); err != nil {
		return nil, err
	} else {
		f.Close()
		return Cache(server.ServeHTTP), nil
	}
}
