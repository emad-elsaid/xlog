package rss

import (
	"encoding/xml"
	"flag"
	"fmt"
	"html/template"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/emad-elsaid/xlog"
)

var domain string
var description string
var limit int

func init() {
	flag.StringVar(&domain, "rss.domain", "", "RSS domain name to be used for RSS feed. without HTTPS://")
	flag.StringVar(&description, "rss.description", "", "RSS feed description")
	flag.IntVar(&limit, "rss.limit", 30, "Limit the number of items in the RSS feed to this amount")

	xlog.RegisterExtension(RSS{})
}

type RSS struct{}

func (RSS) Name() string { return "rss" }
func (RSS) Init() {
	xlog.RegisterWidget(xlog.WidgetHead, 0, metaTag)
	xlog.RegisterBuildPage("/+/feed.rss", false)
	xlog.RegisterLink(links)
	xlog.Get(`/+/feed.rss`, feed)
}

type rssLink struct{}

func (rssLink) Icon() string { return "fa-solid fa-rss" }
func (rssLink) Name() string { return "RSS" }
func (rssLink) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"href": "/+/feed.rss",
	}
}

func links(p xlog.Page) []xlog.Command {
	return []xlog.Command{rssLink{}}
}

func metaTag(p xlog.Page) template.HTML {
	tag := `<link href="/+/feed.rss" rel="alternate" title="%s" type="application/rss+xml">`
	return template.HTML(fmt.Sprintf(tag, template.JSEscapeString(xlog.Config.Sitename)))
}

type rss struct {
	Version string  `xml:"version,attr"`
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Language    string `xml:"language"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	PubDate     time.Time `xml:"pubDate"`
	GUID        string    `xml:"guid"`
	Link        string    `xml:"link"`
}

func feed(r xlog.Request) xlog.Output {
	f := rss{
		Version: "2.0",
		Channel: Channel{
			Title: xlog.Config.Sitename,
			Link: (&url.URL{
				Scheme: "https",
				Host:   domain,
				Path:   "/+/feed.rss",
			}).String(),
			Description: description,
			Language:    "en-US",
			Items:       []Item{},
		},
	}

	pages := xlog.Pages(r.Context())
	slices.SortFunc(pages, func(a, b xlog.Page) int {
		if modtime := b.ModTime().Compare(a.ModTime()); modtime != 0 {
			return modtime
		}

		return strings.Compare(a.Name(), b.Name())
	})

	if len(pages) > limit {
		pages = pages[0:limit]
	}

	for _, p := range pages {
		f.Channel.Items = append(f.Channel.Items, Item{
			Title:       p.Name(),
			Description: string(p.Render()),
			PubDate:     p.ModTime(),
			GUID:        p.Name(),
			Link: (&url.URL{
				Scheme: "https",
				Host:   domain,
				Path:   "/" + p.Name(),
			}).String(),
		})
	}

	buff, err := xml.MarshalIndent(f, "", "    ")
	if err != nil {
		return xlog.InternalServerError(err)
	}

	return xlog.PlainText(xml.Header + string(buff))
}
