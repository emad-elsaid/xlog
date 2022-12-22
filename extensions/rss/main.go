package rss

import (
	"encoding/xml"
	"flag"
	"fmt"
	"html/template"
	"net/url"
	"sort"
	"time"

	. "github.com/emad-elsaid/xlog"
)

var domain string
var description string
var limit int

func init() {
	flag.StringVar(&domain, "rss.domain", "", "RSS domain name to be used for RSS feed. without HTTPS://")
	flag.StringVar(&description, "rss.description", "", "RSS feed description")
	flag.IntVar(&limit, "rss.limit", 30, "Limit the number of items in the RSS feed to this amount")

	RegisterWidget(HEAD_WIDGET, 0, metaTag)
	RegisterBuildPage("/+/feed.rss", false)
	Get(`/\+/feed.rss`, feed)
}

func metaTag(p Page, r Request) template.HTML {
	tag := `<link href="/+/feed.rss" rel="alternate" title="%s" type="application/rss+xml">`
	return template.HTML(fmt.Sprintf(tag, template.JSEscapeString(SITENAME)))
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

func feed(w Response, r Request) Output {
	f := rss{
		Version: "2.0",
		Channel: Channel{
			Title: SITENAME,
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

	pages := []Page{}

	EachPage(r.Context(), func(p Page) {
		pages = append(pages, p)
	})

	sort.Slice(pages, func(i, j int) bool {
		return pages[i].ModTime().After(pages[j].ModTime())
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
		return InternalServerError(err)
	}

	return PlainText(xml.Header + string(buff))
}
