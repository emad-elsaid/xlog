package github

import (
	"flag"
	"fmt"
	"html/template"

	. "github.com/emad-elsaid/xlog"
)

var editUrl string

func init() {
	flag.StringVar(&editUrl, "github.url", "", "Repository url for 'edit on Github' quick action e.g https://github.com/emad-elsaid/xlog/edit/master/docs")
	RegisterExtension(Github{})
}

type Github struct{}

func (Github) Name() string { return "github" }
func (Github) Init() {
	if len(editUrl) == 0 {
		return
	}

	RegisterQuickCommand(quickCommands)
}

func quickCommands(p Page) []Command {
	if len(p.FileName()) == 0 {
		return nil
	}

	return []Command{editOnGithub{page: p}}
}

type editOnGithub struct {
	page Page
}

func (e editOnGithub) Icon() string { return "fa-brands fa-github" }
func (e editOnGithub) Name() string { return "Edit on Github" }
func (e editOnGithub) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"href": fmt.Sprintf("%s/%s", editUrl, e.page.FileName()),
	}
}
