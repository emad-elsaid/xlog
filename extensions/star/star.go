package star

import (
	"embed"
	"fmt"
	"html/template"
	"net/url"
	"strings"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

const STARRED_PAGES = "starred"

//go:embed templates
var templates embed.FS

func init() {
	Widget(ACTION_WIDGET, starMeta)
	Widget(SIDEBAR_WIDGET, starredPages)
	Post(`/\+/star/{page:.*}`, starHandler)
	Delete(`/\+/star/{page:.*}`, unstarHandler)
	Template(templates, "templates")
}

func starredPages(p Page, r Request) template.HTML {
	pages := NewPage(STARRED_PAGES)
	content := strings.TrimSpace(pages.Content())
	if content == "" {
		return template.HTML("")
	}

	list := strings.Split(content, "\n")
	ps := make([]Page, 0, len(list))
	for _, v := range list {
		p := NewPage(v)
		ps = append(ps, p)
	}
	return Partial("starred", Locals{
		"pages": ps,
	})
}

func starMeta(p Page, r Request) template.HTML {
	if READONLY {
		return ""
	}

	starred := isStarred(p)

	return Partial("star-meta", Locals{
		"csrf":    CSRF(r),
		"starred": starred,
		"action":  fmt.Sprintf("/+/star/%s", url.PathEscape(p.Name())),
	})
}

func starHandler(w Response, r Request) Output {
	if READONLY {
		return Unauthorized("Readonly mode is active")
	}

	vars := Vars(r)
	page := NewPage(vars["page"])
	if !page.Exists() {
		return Redirect("/")
	}

	starred_pages := NewPage(STARRED_PAGES)
	starred_pages.Write(strings.TrimSpace(starred_pages.Content()) + "\n" + page.Name())
	return Redirect("/" + page.Name())
}

func unstarHandler(w Response, r Request) Output {
	if READONLY {
		return Unauthorized("Readonly mode is active")
	}

	vars := Vars(r)
	page := NewPage(vars["page"])
	if !page.Exists() {
		return Redirect("/")
	}

	starred_pages := NewPage(STARRED_PAGES)
	content := strings.Split(strings.TrimSpace(starred_pages.Content()), "\n")
	new_content := ""
	for _, v := range content {
		if v != page.Name() {
			new_content += "\n" + v
		}
	}
	starred_pages.Write(new_content)

	return Redirect("/" + page.Name())
}

func isStarred(p Page) bool {
	starred_page := NewPage(STARRED_PAGES)
	for _, k := range strings.Split(starred_page.Content(), "\n") {
		if strings.TrimSpace(k) == p.Name() {
			return true
		}
	}

	return false
}
