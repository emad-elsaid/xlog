package xlog

import (
	"embed"
	"io/fs"
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

// GetAssets returns the embedded assets filesystem
func GetAssets() embed.FS {
	return assets
}

// RegisterStaticDir adds a filesystem to the filesystems list scanned for files
// when serving static files. can be used to add a directory of CSS or JS files
// by extensions
func RegisterStaticDir(f fs.FS) {
	staticDirs = append(staticDirs, f)
}

func staticHandler(r Request) (Output, error) {
	app := GetApp()
	return app.staticHandler(r)
}
