package activitypub

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	. "github.com/emad-elsaid/xlog"
)

var domain string
var username string
var summary string

func init() {
	flag.StringVar(&domain, "activitypub.domain", "", "domain used for activitypub.")
	flag.StringVar(&username, "activitypub.username", "xlog", "username for activitypub steam")
	flag.StringVar(&summary, "activitypub.summary", "", "summary of the user for activitypub")

	Get(`/\.well-known/webfinger`, webfinger)
	Get(`/\+/activitypub/@{user:.+}`, profile)
	Get(`/\+/activitypub/@{user:.+}/outbox`, outbox)
	Get(`/\+/activitypub/@{user:.+}/outbox/{page:[0-9]+}`, outboxPage)
}

type webFingerResponse struct {
	Subject string              `json:"subject,omitempty"`
	Aliases []string            `json:"aliases,omitempty"`
	Links   []map[string]string `json:"links,omitempty"`
}

func webfinger(w Response, r Request) Output {
	return JsonResponse(
		webFingerResponse{
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
	Icon              []string          `json:"icon,omitempty"`
	Image             []string          `json:"image,omitempty"`
}

func profile(w Response, r Request) Output {
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
			Icon:  []string{fmt.Sprintf("https://%s/public/logo.png", domain)},
			Image: []string{fmt.Sprintf("https://%s/public/logo.png", domain)},
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

func outbox(w Response, r Request) Output {
	count := 0
	EachPage(r.Context(), func(_ Page) { count += 1 })

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

func outboxPage(w Response, r Request) Output {
	vars := Vars(r)
	pageIndex, _ := strconv.ParseInt(vars["page"], 10, 64)
	var count int64
	var page Page
	EachPage(r.Context(), func(p Page) {
		count += 1
		if count == pageIndex {
			page = p
		}
	})

	return JsonResponse(
		outboxPageResponse{
			Context: "https://www.w3.org/ns/activitystreams",
			ID:      "https://{{.domain}}/+/activitypub/@{{.username}}/outbox/{{.page}}",
			Type:    "OrderedCollectionPage",
			Prev:    "https://{{.domain}}/+/activitypub/@{{.username}}/outbox/1",
			PartOf:  "https://{{.domain}}/+/activitypub/@{{.username}}/outbox",
			OrderedItems: []outboxPageItem{
				{
					ID:        "https://{{.domain}}/+/activitypub/@{{.username}}/statuses/109420100218327922/activity",
					Type:      "Create",
					Actor:     "https://{{.domain}}/+/activitypub/@{{.username}}",
					Published: page.ModTime(),
					To:        []string{"https://www.w3.org/ns/activitystreams#Public"},
					Object: outboxPageObject{
						ID:           "https://{{.domain}}/{{.name}}",
						Type:         "Note",
						Published:    page.ModTime(),
						URL:          "https://{{.domain}}/{{.name}}",
						AttributedTo: "https://{{.domain}}/+/activitypub/@{{.username}}",
						To:           []string{"https://www.w3.org/ns/activitystreams#Public"},
						Content:      page.Name() + "\n" + page.Content(),
					},
				},
			},
		},
	)
}
