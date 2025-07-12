package sitemap

import (
	"flag"
	"fmt"
	"net/url"
	"strings"

	. "github.com/emad-elsaid/xlog"
)

var SITEMAP_DOMAIN string

func init() {
	app := GetApp()
	flag.StringVar(&SITEMAP_DOMAIN, "sitemap.domain", "", "domain name without protocol or trailing / to use for sitemap loc")
	app.RegisterExtension(Sitemap{})
}

type Sitemap struct{}

func (Sitemap) Name() string { return "sitemap" }
func (Sitemap) Init() {
	app := GetApp()
	app.Get(`/sitemap.xml`, handler)
	app.RegisterBuildPage("/sitemap.xml", false)
}

func handler(r Request) Output {
	output := []string{`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`}
	app := GetApp()
	output = append(output, MapPage(app, r.Context(), func(p Page) string {
		return fmt.Sprintf("<url><loc>https://%s/%s</loc></url>", SITEMAP_DOMAIN, url.PathEscape(p.Name()))
	})...)

	output = append(output, `</urlset>`)

	return PlainText(strings.Join(output, "\n"))
}
