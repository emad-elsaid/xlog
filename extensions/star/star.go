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
	RegisterExtension(Star{})
}

type Star struct{}

func (Star) Name() string { return "star" }
func (Star) Init() {
	RegisterLink(starredPages)
	RegisterTemplate(templates, "templates")

	if !Config.Readonly {
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

func (s starredPage) Link() string {
	return "/" + s.Name()
}

func starredPages(p Page) []Link {
	pages := NewPage(STARRED_PAGES)
	content := strings.TrimSpace(string(pages.Content()))
	if content == "" {
		return nil
	}

	list := strings.Split(content, "\n")
	ps := make([]Link, 0, len(list))
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
func (action) Link() string         { return "" }
func (action) OnClick() template.JS { return "star(event)" }
func (l action) Widget() template.HTML {
	starred := isStarred(l.page)

	return Partial("star", Locals{
		"starred": starred,
		"action":  fmt.Sprintf("/+/star/%s", url.PathEscape(l.page.Name())),
	})
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
	return NoContent()
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

	return NoContent()
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
