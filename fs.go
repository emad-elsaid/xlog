package xlog

import (
	"io/fs"
	"net/http"
	"os"
	"path"
)

// priorityFS returns file that exists in one of the FS structs.
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

// RegisterStaticDir adds a filesystem to the static files list
func (app *App) RegisterStaticDir(f fs.FS) {
	app.staticDirs = append(app.staticDirs, f)
}

// staticHandler handles static file serving
func (app *App) staticHandler(r Request) (Output, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	staticFSs := http.FS(
		priorityFS(
			append(app.staticDirs, os.DirFS(wd)),
		),
	)

	server := http.FileServer(staticFSs)

	cleanPath := path.Clean(r.URL.Path)

	if f, err := staticFSs.Open(cleanPath); err != nil {
		return nil, err
	} else {
		f.Close()
		return app.Cache(server.ServeHTTP), nil
	}
}
