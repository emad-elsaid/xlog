package sitemap

import (
	"embed"
	"flag"
	"net/url"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS
var SITEMAP_DOMAIN string

func init() {
	flag.StringVar(&SITEMAP_DOMAIN, "sitemap.domain", "", "domain name without protocol or trailing / to use for sitemap loc")
	RegisterTemplate(templates, "templates")
	Get(`/sitemap\.xml`, handler)
	RegisterBuildPage("/sitemap.xml", false)
	RegisterHelper("urlescaper", url.PathEscape)
}

func handler(w Response, r Request) Output {
	pages := []Page{}
	EachPage(r.Context(), func(p Page) {
		pages = append(pages, p)
	})

	return Render("sitemap", Locals{
		"pages":  pages,
		"domain": SITEMAP_DOMAIN,
	})
}
