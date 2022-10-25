package xlog

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"strings"
)

const TEMPLATE_EXTENSION = ".html"

//go:embed templates
var defaultTemplates embed.FS
var templates *template.Template
var helpers = template.FuncMap{}
var templatesFSs []fs.FS

// Template registers a filesystem that contains templates, templates are
// registered such that the latest directory override older ones. template file
// extensions are signified by TEMPLATE_EXTENSION constant and the file path can
// be used as template name without this extension
func Template(t fs.FS) {
	templatesFSs = append(templatesFSs, t)
}

func compileTemplates() {
	// add default templates before everything else
	sub, _ := fs.Sub(defaultTemplates, "templates")
	templatesFSs = append([]fs.FS{sub}, templatesFSs...)

	templates = template.New("")
	for _, tfs := range templatesFSs {
		fs.WalkDir(tfs, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if strings.HasSuffix(p, TEMPLATE_EXTENSION) && d.Type().IsRegular() {
				ext := path.Ext(p)
				name := strings.TrimSuffix(p, ext)
				defer Log(DEBUG, "View", name)()

				c, err := fs.ReadFile(tfs, p)
				if err != nil {
					return err
				}

				template.Must(templates.New(name).Funcs(helpers).Parse(string(c)))
			}

			return nil
		})
	}
}

func Partial(path string, data Locals) string {
	v := templates.Lookup(path)
	if v == nil {
		return fmt.Sprintf("template %s not found", path)
	}

	// set extra locals here
	if data == nil {
		data = Locals{}
	}

	data["SITENAME"] = SITENAME
	data["READONLY"] = READONLY
	data["SIDEBAR"] = SIDEBAR

	w := bytes.NewBufferString("")
	err := v.Execute(w, data)
	if err != nil {
		return "rendering error " + path + " " + err.Error()
	}

	return w.String()
}
