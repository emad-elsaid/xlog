// Package blocks provides customizable content blocks through YAML-based shortcodes.
//
// This extension allows you to create rich, styled content blocks in your markdown
// pages using YAML syntax combined with HTML templates. Each template in the
// templates/ directory automatically becomes available as a shortcode.
//
// # Usage
//
// In your markdown, use the shortcode syntax with YAML data:
//
//	{{% blockname %}}
//	title: Example Block
//	content: This is my content
//	highlight: true
//	{{% /blockname %}}
//
// The extension will parse the YAML data and render it using the corresponding
// template (templates/blockname.html). Template data is accessible as variables.
//
// # Creating Custom Blocks
//
// 1. Add a template file: templates/myblock.html
// 2. Use Go template syntax: {{.title}}, {{.content}}, etc.
// 3. The block becomes immediately available as {{% myblock %}}
//
// # Built-in Styling
//
// The extension includes blocks.css automatically in all pages via the head widget.
// Custom block templates can reference CSS classes defined in public/blocks.css.
//
// # Error Handling
//
// If YAML parsing fails, the error message is safely HTML-escaped and displayed
// inline, making debugging straightforward without breaking page rendering.
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
