package blocks

import (
	"embed"
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

func RegisterShortCodes() {
	fs.WalkDir(templates, "templates", func(path string, d fs.DirEntry, err error) error {
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
	fs.WalkDir(public, ".", func(path string, d fs.DirEntry, err error) error {
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
			return template.HTML(err.Error())
		}

		output := xlog.Partial(tpl, xlog.Locals(b))

		return template.HTML(output)
	}
}
