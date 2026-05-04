package xlog

import (
	"bytes"
	"embed"
	"fmt"
	"html"
	"html/template"
	"io/fs"
	"log/slog"
	"os"
	"path"
	"strings"
)

//go:embed templates
var defaultTemplates embed.FS
var templates *template.Template
var templatesFSs []fs.FS

// RegisterTemplate registers a filesystem that contains templates, specifying subDir as
// the subdirectory name that contains the templates. templates are registered
// such that the latest registered directory override older ones. template file
// extensions are signified by '.html' extension and the file path can
// be used as template name without this extension.
func RegisterTemplate(t fs.FS, subDir string) {
	ts, _ := fs.Sub(t, subDir)
	templatesFSs = append(templatesFSs, ts)
}

func compileTemplates() {
	const ext = ".html"

	// add default templates before everything else
	sub, _ := fs.Sub(defaultTemplates, "templates")
	templatesFSs = append([]fs.FS{sub}, templatesFSs...)
	// add theme directory after everything else to allow user to override any template
	if _, err := os.Stat("theme"); err == nil {
		templatesFSs = append(templatesFSs, os.DirFS("theme"))
	}

	templates = template.New("")
	for _, tfs := range templatesFSs {
		_ = fs.WalkDir(tfs, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("failed to walk template directory at %s: %w", p, err)
			}

			if strings.HasSuffix(p, ext) && d.Type().IsRegular() {
				ext := path.Ext(p)
				name := strings.TrimSuffix(p, ext)
				slog.Info("Template " + name)

				c, err := fs.ReadFile(tfs, p)
				if err != nil {
					return fmt.Errorf("failed to read template file %s: %w", p, err)
				}

				template.Must(templates.New(name).Funcs(helpers).Parse(string(c)))
			}

			return nil
		})
	}
}

// Partial executes a template by it's path name. it passes data to the
// template. returning the output of the template. in case of an error it will
// return the error string as the output.
func Partial(templatePath string, data Locals) template.HTML {
	v := templates.Lookup(templatePath)
	if v == nil {
		// Escape path to prevent potential XSS in error message
		// #nosec G203 - path is escaped with html.EscapeString, safe to convert
		return template.HTML(fmt.Sprintf("template %s not found", html.EscapeString(templatePath)))
	}

	if data == nil {
		data = Locals{}
	}

	data["config"] = Config

	w := bytes.NewBufferString("")

	if err := v.Execute(w, data); err != nil {
		// Escape error details to prevent potential XSS
		// #nosec G203 - both path and err are escaped with html.EscapeString, safe to convert
		return template.HTML("rendering error " + html.EscapeString(templatePath) + " " + html.EscapeString(err.Error()))
	}

	// #nosec G203 - w.String() contains output from template.Execute which auto-escapes, safe to convert
	return template.HTML(w.String())
}
