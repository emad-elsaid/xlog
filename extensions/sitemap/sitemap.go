package sitemap

import (
	"flag"
	"fmt"
	"net/url"

	. "github.com/emad-elsaid/xlog"
)

var SITEMAP_DOMAIN string

func init() {
	flag.StringVar(&SITEMAP_DOMAIN, "sitemap.domain", "", "domain name without protocol or trailing / to use for sitemap loc")
	Get(`/sitemap\.xml`, handler)
	RegisterBuildPage("/sitemap.xml", false)
}

func handler(w Response, r Request) Output {
	output := `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`

	EachPage(r.Context(), func(p Page) {
		output += fmt.Sprintf("<url><loc>https://%s/%s</loc></url>", SITEMAP_DOMAIN, url.PathEscape(p.Name()))
	})

	output += `</urlset>`

	return PlainText(output)
}
