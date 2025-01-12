package activitypub

import (
	"flag"
	"fmt"
	"html/template"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	. "github.com/emad-elsaid/xlog"
)

var domain string
var username string
var summary string
var icon string
var image string

func init() {
	flag.StringVar(&domain, "activitypub.domain", "", "domain used for activitypub stream absolute URLs")
	flag.StringVar(&username, "activitypub.username", "", "username for activitypub actor")
	flag.StringVar(&summary, "activitypub.summary", "", "summary of the user for activitypub actor")
	flag.StringVar(&icon, "activitypub.icon", "/public/logo.png", "the path to the activitypub profile icon. mastodon use it as profile picture for example.")
	flag.StringVar(&image, "activitypub.image", "/public/logo.png", "the path to the activitypub profile image. mastodon use it as profile cover for example.")

	RegisterExtension(ActivityPub{})
}

type ActivityPub struct{}

func (ActivityPub) Name() string { return "activitypub" }
func (ActivityPub) Init() {
	Get(`/.well-known/webfinger`, webfinger)
	Get(`/+/activitypub/{user}/outbox/{page}`, outboxPage)
	Get(`/+/activitypub/{user}/outbox`, outbox)
	Get(`/+/activitypub/{user}`, profile)
	RegisterWidget(WidgetHead, 1, meta)
}

func meta(p Page) template.HTML {
	if domain == "" || username == "" {
		return ""
	}

	RegisterBuildPage("/.well-known/webfinger", false)
	RegisterBuildPage(fmt.Sprintf("/+/activitypub/@%s", username), true)
	RegisterBuildPage(fmt.Sprintf("/+/activitypub/@%s/outbox", username), true)

	o := fmt.Sprintf(`<link href='https://%s/+/activitypub/@%s' rel='alternate' type='application/activity+json'>`, domain, username)

	return template.HTML(o)
}

type webfingerResponse struct {
	Subject string              `json:"subject,omitempty"`
	Aliases []string            `json:"aliases,omitempty"`
	Links   []map[string]string `json:"links,omitempty"`
}

func webfinger(r Request) Output {
	if domain == "" || username == "" {
		return NoContent()
	}

	return JsonResponse(
		webfingerResponse{
			Subject: fmt.Sprintf("acct:%s@%s", username, domain),
			Aliases: []string{
				fmt.Sprintf("https://%s", domain),
				fmt.Sprintf("https://%s/+/activitypub/@%s", domain, username),
			},
			Links: []map[string]string{
				{
					"rel":  "http://webfinger.net/rel/profile-page",
					"type": "text/html",
					"href": fmt.Sprintf("https://%s", domain),
				},
				{
					"rel":  "self",
					"type": "application/activity+json",
					"href": fmt.Sprintf("https://%s/+/activitypub/@%s", domain, username),
				},
				// TODO we need to make sure this is actually needed
				{
					"rel":      "http://ostatus.org/schema/1.0/subscribe",
					"template": fmt.Sprintf("https://%s/authorize_interaction?uri={uri}", domain),
				},
			},
		},
	)
}

type profileResponse struct {
	Context           string            `json:"@context,omitempty"`
	ID                string            `json:"id,omitempty"`
	Type              string            `json:"type,omitempty"`
	PreferredUsername string            `json:"preferredUsername,omitempty"`
	Name              string            `json:"name,omitempty"`
	Summary           string            `json:"summary,omitempty"`
	URL               string            `json:"url,omitempty"`
	Inbox             string            `json:"inbox,omitempty"`
	Outbox            string            `json:"outbox,omitempty"`
	Endpoints         map[string]string `json:"endpoints,omitempty"`
	Icon              map[string]string `json:"icon,omitempty"`
	Image             map[string]string `json:"image,omitempty"`
}

func profile(r Request) Output {
	if domain == "" || username == "" {
		return NoContent()
	}

	if r.PathValue("user") != "@"+username {
		return NotFound("User not found")
	}

	return JsonResponse(
		profileResponse{
			Context:           "https://www.w3.org/ns/activitystreams",
			ID:                fmt.Sprintf("https://%s/+/activitypub/@%s", domain, username),
			Type:              "Person",
			PreferredUsername: username,
			Name:              username,
			Summary:           summary,
			URL:               fmt.Sprintf("https://%s", domain),
			Inbox:             fmt.Sprintf("https://%s/+/activitypub/@%s/inbox", domain, username),
			Outbox:            fmt.Sprintf("https://%s/+/activitypub/@%s/outbox", domain, username),
			Endpoints: map[string]string{
				"sharedInbox": fmt.Sprintf("https://%s/+/activitypub/@%s/inbox", domain, username),
			},
			Icon: map[string]string{
				"type":      "Image",
				"mediaType": "image/png",
				"url":       fmt.Sprintf("https://%s%s", domain, icon),
			},
			Image: map[string]string{
				"type":      "Image",
				"mediaType": "image/png",
				"url":       fmt.Sprintf("https://%s%s", domain, image),
			},
		},
	)

}

type outboxResponse struct {
	Context    string `json:"@context,omitempty"`
	ID         string `json:"id,omitempty"`
	Type       string `json:"type,omitempty"`
	TotalItems int    `json:"totalItems,omitempty"`
	First      string `json:"first,omitempty"`
	Last       string `json:"last,omitempty"`
}

func outbox(r Request) Output {
	if domain == "" || username == "" {
		return NoContent()
	}

	if r.PathValue("user") != "@"+username {
		return NotFound("User not found")
	}

	count := 0
	EachPage(r.Context(), func(Page) {
		count += 1
		RegisterBuildPage(fmt.Sprintf("/+/activitypub/@%s/outbox/%d", username, count), false)
	})

	return JsonResponse(
		outboxResponse{
			Context:    "https://www.w3.org/ns/activitystreams",
			ID:         fmt.Sprintf("https://%s/+/activitypub/@%s/outbox", domain, username),
			Type:       "OrderedCollection",
			TotalItems: count,
			First:      fmt.Sprintf("https://%s/+/activitypub/@%s/outbox/1", domain, username),
			Last:       fmt.Sprintf("https://%s/+/activitypub/@%s/outbox/%d", domain, username, count),
		},
	)
}

type outboxPageResponse struct {
	Context      string           `json:"@context,omitempty"`
	ID           string           `json:"id,omitempty"`
	Type         string           `json:"type,omitempty"`
	Prev         string           `json:"prev,omitempty"`
	Next         string           `json:"next,omitempty"`
	PartOf       string           `json:"partOf,omitempty"`
	OrderedItems []outboxPageItem `json:"orderedItems"`
}

type outboxPageItem struct {
	ID        string    `json:"id,omitempty"`
	Type      string    `json:"type,omitempty"`
	Actor     string    `json:"actor"`
	Published time.Time `json:"published"`
	To        []string  `json:"to"`
	Object    outboxPageObject
}

type outboxPageObject struct {
	ID           string    `json:"id,omitempty"`
	Type         string    `json:"type,omitempty"`
	Published    time.Time `json:"published"`
	URL          string    `json:"url"`
	AttributedTo string    `json:"attributedTo"`
	To           []string  `json:"to"`
	Content      string
}

func outboxPage(r Request) Output {
	if domain == "" || username == "" {
		return NoContent()
	}

	if r.PathValue("user") != "@"+username {
		return NotFound("User not found")
	}

	pageIndex, _ := strconv.ParseInt(r.PathValue("page"), 10, 64)
	pageIndex--

	pages := Pages(r.Context())

	if int(pageIndex) >= len(pages) || pageIndex < 0 {
		return NotFound("page index is out of context")
	}

	var page Page
	slices.SortFunc(pages, func(a, b Page) int {
		if modtime := b.ModTime().Compare(a.ModTime()); modtime != 0 {
			return modtime
		}

		return strings.Compare(a.Name(), b.Name())
	})

	page = pages[pageIndex]

	var u url.URL
	u.Scheme = "https"
	u.Path = "/" + page.Name()
	u.Host = domain

	return JsonResponse(
		outboxPageResponse{
			Context: "https://www.w3.org/ns/activitystreams",
			ID:      fmt.Sprintf("https://%s/+/activitypub/@%s/outbox/%d", domain, username, pageIndex),
			Type:    "OrderedCollectionPage",
			Prev:    fmt.Sprintf("https://%s/+/activitypub/@%s/outbox/%d", domain, username, pageIndex-1),
			Next:    fmt.Sprintf("https://%s/+/activitypub/@%s/outbox/%d", domain, username, pageIndex+1),
			PartOf:  fmt.Sprintf("https://%s/+/activitypub/@%s/outbox", domain, username),
			OrderedItems: []outboxPageItem{
				{
					ID:        u.String(),
					Type:      "Create",
					Actor:     fmt.Sprintf("https://%s/+/activitypub/@%s", domain, username),
					Published: page.ModTime(),
					To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
					Object: outboxPageObject{
						ID:           u.String(),
						Type:         "Note",
						Published:    page.ModTime(),
						URL:          u.String(),
						AttributedTo: fmt.Sprintf("https://%s/+/activitypub/@%s", domain, username),
						To:           []string{"https://www.w3.org/ns/activitystreams#Public"},
						Content:      fmt.Sprintf("%s\n%s", page.Name(), u.String()),
					},
				},
			},
		},
	)
}
