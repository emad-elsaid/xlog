package sitemap

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"time"

	. "github.com/emad-elsaid/xlog"
)

var domain string

func init() {
	flag.StringVar(&domain, "sitemap.domain", "", "domain name to be used for sitemap URLs")
	app := GetApp()
	app.RegisterExtension(Sitemap{})
}

type Sitemap struct{}

func (Sitemap) Name() string { return "sitemap" }
func (Sitemap) Init() {
	app := GetApp()
	app.RegisterBuildPage("/sitemap.xml", false)
	app.Get("/sitemap.xml", sitemapHandler)
}

func sitemapHandler(r Request) Output {
	app := GetApp()
	return app.Cache(func(w Response, r Request) {
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprint(w, generateSitemap(app, r.Context()))
	})
}

func generateSitemap(app *App, ctx context.Context) string {
	urlset := URLSet{
		XML: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	app.EachPage(ctx, func(p Page) {
		if p.Exists() {
			urlset.URLs = append(urlset.URLs, URL{
				Loc:        fmt.Sprintf("https://%s/%s", domain, p.Name()),
				LastMod:    p.ModTime().Format(time.RFC3339),
				ChangeFreq: "weekly",
				Priority:   "0.5",
			})
		}
	})

	output, _ := xml.MarshalIndent(urlset, "", "  ")
	return xml.Header + string(output)
}

type URLSet struct {
	XML  string `xml:"xmlns,attr"`
	URLs []URL  `xml:"url"`
}

type URL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod"`
	ChangeFreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}
