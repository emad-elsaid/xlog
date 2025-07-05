package link_preview

import (
	"embed"
	"html/template"
	"net/http"
	"strings"
	"time"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
)

//go:embed templates
var templates embed.FS

func init() {
	app := GetApp()
	app.RegisterExtension(LinkPreview{})
}

type LinkPreview struct{}

func (LinkPreview) Name() string { return "link-preview" }
func (LinkPreview) Init() {
	app := GetApp()
	app.RegisterTemplate(templates, "templates")
	app.RegisterWidget(WidgetAfterView, 1, linkPreviewWidget)
}

func linkPreviewWidget(p Page) template.HTML {
	if p == nil {
		return ""
	}

	_, tree := p.AST()
	if tree == nil {
		return ""
	}

	links := FindAllInAST[*ast.Link](tree)
	if len(links) == 0 {
		return ""
	}

	var previews []template.HTML
	for _, link := range links {
		url := string(link.Destination)
		if strings.HasPrefix(url, "http") {
			preview := getLinkPreview(url)
			if preview != "" {
				previews = append(previews, preview)
			}
		}
	}

	if len(previews) == 0 {
		return ""
	}

	app := GetApp()
	return app.Partial("link-preview", Locals{"previews": previews})
}

func getLinkPreview(url string) template.HTML {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Simple preview - in a real implementation you'd parse the HTML
	// and extract title, description, image, etc.
	return template.HTML(`<div class="link-preview"><a href="` + url + `">` + url + `</a></div>`)
}
