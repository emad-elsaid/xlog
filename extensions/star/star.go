package star

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

const STARRED_PAGES = "starred"

func init() {
	RegisterExtension(Star{})
}

type Star struct{}

func (Star) Name() string { return "star" }
func (Star) Init() {
	RegisterLink(starredPages)
	IgnorePath(regexp.MustCompile(`^starred\.md$`))

	if !Config.Readonly {
		RequireHTMX()
		RegisterCommand(starAction)
		RegisterQuickCommand(starAction)
		Post(`/+/star/{page...}`, starHandler)
		Delete(`/+/star/{page...}`, unstarHandler)
	}
}

type starredPage struct {
	Page
}

func (s starredPage) Icon() string {
	if s.Emoji() == "" {
		return "fa-solid fa-star"
	} else {
		return s.Emoji()
	}
}

func (s starredPage) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"href": "/" + s.Name(),
	}
}

func starredPages(p Page) []Command {
	pages := NewPage(STARRED_PAGES)
	content := strings.TrimSpace(string(pages.Content()))
	if content == "" {
		return nil
	}

	list := strings.Split(content, "\n")
	ps := make([]Command, 0, len(list))
	for _, v := range list {
		p := starredPage{NewPage(v)}
		ps = append(ps, p)
	}

	return ps
}

type action struct {
	page    Page
	starred bool
}

func (l action) Icon() string {
	if l.starred {
		return "fa-solid fa-star"
	} else {
		return "fa-regular fa-star"
	}
}
func (l action) Name() string {
	if l.starred {
		return "Unstar"
	} else {
		return "Star"
	}
}
func (l action) Attrs() map[template.HTMLAttr]any {
	var method template.HTMLAttr = "hx-post"
	if l.starred {
		method = "hx-delete"
	}

	return map[template.HTMLAttr]any{
		method: fmt.Sprintf("/+/star/%s", url.PathEscape(l.page.Name())),
		"href": fmt.Sprintf("/+/star/%s", url.PathEscape(l.page.Name())),
	}
}

func starAction(p Page) []Command {
	starred := isStarred(p)
	return []Command{action{starred: starred, page: p}}
}

func starHandler(r Request) Output {
	page := NewPage(r.PathValue("page"))
	if !page.Exists() {
		return Redirect("/")
	}

	starred_pages := NewPage(STARRED_PAGES)
	new_content := strings.TrimSpace(string(starred_pages.Content())) + "\n" + page.Name()
	starred_pages.Write(Markdown(new_content))

	return func(w Response, r Request) {
		w.Header().Add("HX-Refresh", "true")
		w.WriteHeader(http.StatusNoContent)
	}
}

func unstarHandler(r Request) Output {
	page := NewPage(r.PathValue("page"))
	if !page.Exists() {
		return Redirect("/")
	}

	starred_pages := NewPage(STARRED_PAGES)
	content := strings.Split(strings.TrimSpace(string(starred_pages.Content())), "\n")
	new_content := ""
	for _, v := range content {
		if v != page.Name() {
			new_content += "\n" + v
		}
	}
	starred_pages.Write(Markdown(new_content))

	return func(w Response, r Request) {
		w.Header().Add("HX-Refresh", "true")
		w.WriteHeader(http.StatusNoContent)
	}
}

func isStarred(p Page) bool {
	starred_page := NewPage(STARRED_PAGES)
	for _, k := range strings.Split(string(starred_page.Content()), "\n") {
		if strings.TrimSpace(k) == p.Name() {
			return true
		}
	}

	return false
}
