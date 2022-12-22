package github

import (
	"flag"
	"fmt"
	"html/template"
	"net/url"

	. "github.com/emad-elsaid/xlog"
)

var repo string
var branch string

func init() {
	flag.StringVar(&repo, "github.repo", "", "Github repository to use for 'edit on Github' quick action")
	flag.StringVar(&branch, "github.branch", "master", "Github repository branch to use for 'edit on Github' quick action")
	RegisterQuickCommand(quickCommands)
}

func quickCommands(p Page) []Command {
	if len(repo) == 0 {
		return nil
	}

	return []Command{editOnGithub{page: p}}
}

type editOnGithub struct {
	page Page
}

func (e editOnGithub) Icon() string {
	return "fa-solid fa-pen"
}
func (e editOnGithub) Name() string {
	return "Edit on Github"
}
func (e editOnGithub) Link() string {
	return fmt.Sprintf("%s/edit/%s/%s", repo, branch, url.PathEscape(e.page.FileName()))
}
func (e editOnGithub) OnClick() template.JS {
	return ""
}
func (e editOnGithub) Widget() template.HTML {
	return ""
}
