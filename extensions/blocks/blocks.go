package blocks

import (
	"embed"
	"html"
	"html/template"
	"io/fs"
	"strings"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/extensions/shortcode"
	"gopkg.in/yaml.v3"
)

//go:embed templates
var templates embed.FS

//go:embed public
var public embed.FS

func init() {
	xlog.RegisterExtension(Blocks{})
}

type Blocks struct{}

func (Blocks) Name() string { return "blocks" }
func (Blocks) Init() {
	RegisterShortCodes()
	xlog.RegisterTemplate(templates, "templates")
	xlog.RegisterStaticDir(public)
	registerBuildFiles()
	xlog.RegisterWidget(xlog.WidgetHead, 0, style)
}

// RegisterShortCodes walks the templates directory and registers each template
// as a shortcode that can be used in pages.
func RegisterShortCodes() {
	_ = fs.WalkDir(templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		name := strings.TrimPrefix(path, "templates/")
		name = strings.TrimSuffix(name, ".html")

		shortcode.RegisterShortCode(name, shortcode.ShortCode{Render: block(name)})

		return nil
	})
}

func registerBuildFiles() {
	_ = fs.WalkDir(public, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		xlog.RegisterBuildPage("/"+path, false)

		return nil
	})
}

func style(xlog.Page) template.HTML {
	return `<link rel="stylesheet" href="/public/blocks.css">`
}

func block(tpl string) func(xlog.Markdown) template.HTML {
	return func(in xlog.Markdown) template.HTML {
		b := map[string]any{}

		if err := yaml.Unmarshal([]byte(in), &b); err != nil {
			// Escape error message to prevent potential XSS if error contains user input.
			// #nosec G203 - Error string is HTML-escaped before conversion to template.HTML
			return template.HTML(html.EscapeString(err.Error()))
		}

		output := xlog.Partial(tpl, xlog.Locals(b))

		// #nosec G203 - xlog.Partial returns trusted output from html/template execution
		return template.HTML(output)
	}
}
