package xlog

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"strings"
)

//go:embed templates
var defaultTemplates embed.FS

// RegisterTemplate registers a filesystem that contains templates
func (app *App) RegisterTemplate(t fs.FS, subDir string) {

	ts, _ := fs.Sub(t, subDir)
	app.templatesFSs = append(app.templatesFSs, ts)
}

// compileTemplates compiles all registered templates
func (app *App) compileTemplates() {
	const ext = ".html"

	// add default templates before everything else
	sub, _ := fs.Sub(defaultTemplates, "templates")
	app.templatesFSs = append([]fs.FS{sub}, app.templatesFSs...)
	// add theme directory after everything else to allow user to override any template
	if _, err := os.Stat("theme"); err == nil {
		app.templatesFSs = append(app.templatesFSs, os.DirFS("theme"))
	}

	app.templates = template.New("")
	for _, tfs := range app.templatesFSs {
		fs.WalkDir(tfs, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(p, ext) && d.Type().IsRegular() {
				ext := path.Ext(p)
				name := strings.TrimSuffix(p, ext)
				slog.Info("Template " + name)

				c, err := fs.ReadFile(tfs, p)
				if err != nil {
					return err
				}

				template.Must(app.templates.New(name).Funcs(app.helpers).Parse(string(c)))
			}

			return nil
		})
	}
}

// Partial executes a template by it's path name
func (app *App) Partial(path string, data Locals) template.HTML {
	v := app.templates.Lookup(path)
	if v == nil {
		return template.HTML(fmt.Sprintf("template %s not found", path))
	}

	if data == nil {
		data = Locals{}
	}

	data["config"] = app.config

	w := bytes.NewBufferString("")

	if err := v.Execute(w, data); err != nil {
		return template.HTML("rendering error " + path + " " + err.Error())
	}

	return template.HTML(w.String())
}
