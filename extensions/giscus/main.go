package giscus

import (
	"flag"
	"fmt"
	"html/template"

	"github.com/emad-elsaid/xlog"
)

const tmpl = `
<script src="https://giscus.app/client.js"
        data-repo="%s"
        data-repo-id="%s"
        data-category="%s"
        data-category-id="%s"
        data-mapping="%s"
        data-strict="0"
        data-reactions-enabled="1"
        data-emit-metadata="0"
        data-input-position="bottom"
        data-theme="%s"
        data-lang="%s"
        crossorigin="anonymous"
        async>
</script>`

var (
	repo       string
	repoID     string
	category   string
	categoryID string
	mapping    string
	theme      string
	lang       string
)

func init() {
	flag.StringVar(&repo, "giscus-repo", "", "GitHub repository in format 'owner/repo' (e.g., 'emad-elsaid/xlog')")
	flag.StringVar(&repoID, "giscus-repo-id", "", "GitHub repository ID (get from giscus.app)")
	flag.StringVar(&category, "giscus-category", "", "GitHub Discussions category name")
	flag.StringVar(&categoryID, "giscus-category-id", "", "GitHub Discussions category ID (get from giscus.app)")
	flag.StringVar(&mapping, "giscus-mapping", "pathname", "Page-discussion mapping method (pathname, url, title, og:title)")
	flag.StringVar(&theme, "giscus-theme", "preferred_color_scheme", "Giscus theme (e.g., light, dark, preferred_color_scheme)")
	flag.StringVar(&lang, "giscus-lang", "en", "Language code for Giscus interface")
	xlog.RegisterExtension(Giscus{})
}

type Giscus struct{}

func (Giscus) Name() string { return "giscus" }
func (Giscus) Init() {
	xlog.RegisterWidget(xlog.WidgetAfterView, 2, widget)
}

func widget(p xlog.Page) template.HTML {
	// Require minimum configuration
	if repo == "" || repoID == "" || category == "" || categoryID == "" {
		return ""
	}

	script := fmt.Sprintf(
		tmpl,
		template.HTMLEscapeString(repo),
		template.HTMLEscapeString(repoID),
		template.HTMLEscapeString(category),
		template.HTMLEscapeString(categoryID),
		template.HTMLEscapeString(mapping),
		template.HTMLEscapeString(theme),
		template.HTMLEscapeString(lang),
	)
	// #nosec G203 -- Template string is constant, all values are HTML-escaped, flags are admin-controlled
	return template.HTML(script)
}
