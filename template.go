package xlog

import (
	"html/template"
	"io/fs"
)

// RegisterTemplate registers a filesystem that contains templates, specifying subDir as
// the subdirectory name that contains the templates. templates are registered
// such that the latest registered directory override older ones. template file
// extensions are signified by '.html' extension and the file path can
// be used as template name without this extension
func RegisterTemplate(t fs.FS, subDir string) {
	app := GetApp()
	app.RegisterTemplate(t, subDir)
}

// Partial executes a template by it's path name. it passes data to the
// template. returning the output of the template. in case of an error it will
// return the error string as the output
func Partial(path string, data Locals) template.HTML {
	app := GetApp()
	return app.Partial(path, data)
}
