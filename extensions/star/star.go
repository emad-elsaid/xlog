package star

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"

	_ "embed"

	"github.com/emad-elsaid/xlog"
)

const (
	STARRED_PAGES = "starred"
	actionUnstar  = "Unstar"
	actionStar    = "Star"
)

func init() {
	xlog.RegisterExtension(Star{})
}

type Star struct{}

func (Star) Name() string { return "star" }
func (Star) Init() {
	xlog.RegisterLink(starredPages)
	xlog.IgnorePath(regexp.MustCompile(`^starred\.md$`))

	if !xlog.Config.Readonly {
		xlog.RequireHTMX()
		xlog.RegisterCommand(starAction)
		xlog.RegisterQuickCommand(starAction)
		xlog.Post(`/+/star/{page...}`, starHandler)
		xlog.Delete(`/+/star/{page...}`, unstarHandler)
	}
}

type starredPage struct {
	Page xlog.Page
}

func (s starredPage) Icon() string {
	if e := xlog.Emoji(s.Page); e == "" {
		return "fa-solid fa-star"
	} else {
		return e
	}
}

func (s starredPage) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"href": "/" + s.Page.Name(),
	}
}

func (s starredPage) Name() string {
	return path.Base(s.Page.Name())
}

func starredPages(p xlog.Page) []xlog.Command {
	pages := xlog.NewPage(STARRED_PAGES)
	if pages == nil {
		return nil
	}

	content := strings.TrimSpace(string(pages.Content()))
	if content == "" {
		return nil
	}

	list := strings.Split(content, "\n")
	ps := make([]xlog.Command, 0, len(list))
	for _, v := range list {
		p := starredPage{Page: xlog.NewPage(v)}
		ps = append(ps, p)
	}

	return ps
}

type action struct {
	page    xlog.Page
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
		return actionUnstar
	} else {
		return actionStar
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

func starAction(p xlog.Page) []xlog.Command {
	if !p.Exists() {
		return nil
	}

	starred := isStarred(p)
	return []xlog.Command{action{starred: starred, page: p}}
}

func starHandler(r xlog.Request) xlog.Output {
	page := xlog.NewPage(r.PathValue("page"))

	if page == nil || !page.Exists() {
		return xlog.Redirect("/")
	}

	starred_pages := xlog.NewPage(STARRED_PAGES)
	if starred_pages == nil {
		return xlog.Redirect("/")
	}

	new_content := strings.TrimSpace(string(starred_pages.Content())) + "\n" + page.Name()
	starred_pages.Write(xlog.Markdown(new_content))

	return func(w xlog.Response, r xlog.Request) {
		w.Header().Add("HX-Refresh", "true")
		w.WriteHeader(http.StatusNoContent)
	}
}

func unstarHandler(r xlog.Request) xlog.Output {
	page := xlog.NewPage(r.PathValue("page"))
	if page == nil || !page.Exists() {
		return xlog.Redirect("/")
	}

	starred_pages := xlog.NewPage(STARRED_PAGES)
	if starred_pages == nil {
		return xlog.Redirect("/")
	}

	content := strings.Split(strings.TrimSpace(string(starred_pages.Content())), "\n")
	new_content := ""
	for _, v := range content {
		if v != page.Name() {
			new_content += "\n" + v
		}
	}
	starred_pages.Write(xlog.Markdown(new_content))

	return func(w xlog.Response, r xlog.Request) {
		w.Header().Add("HX-Refresh", "true")
		w.WriteHeader(http.StatusNoContent)
	}
}

func isStarred(p xlog.Page) bool {
	starred_page := xlog.NewPage(STARRED_PAGES)
	if starred_page == nil {
		return false
	}

	for _, k := range strings.Split(string(starred_page.Content()), "\n") {
		if strings.TrimSpace(k) == p.Name() {
			return true
		}
	}

	return false
}
