package github

import (
	"flag"
	"fmt"
	"html/template"

	. "github.com/emad-elsaid/xlog"
)

var editUrl string
var repo string
var branch string

func init() {
	flag.StringVar(&editUrl, "github.url", "", "Repository url for 'edit on Github' quick action e.g https://github.com/emad-elsaid/xlog/edit/master/docs")
	flag.StringVar(&repo, "github.repo", "", "[Deprecated] Github repository to use for 'edit on Github' quick action e.g https://github.com/emad-elsaid/xlog")
	flag.StringVar(&branch, "github.branch", "master", "[Deprecated] Github repository branch to use for 'edit on Github' quick action")
	RegisterExtension(Github{})
}

type Github struct{}

func (Github) Name() string { return "github" }
func (Github) Init()        { RegisterQuickCommand(quickCommands) }

func quickCommands(p Page) []Command {
	if len(repo) == 0 && len(editUrl) == 0 {
		return nil
	}

	return []Command{editOnGithub{page: p}}
}

type editOnGithub struct {
	page Page
}

func (e editOnGithub) Icon() string { return "fa-brands fa-github" }
func (e editOnGithub) Name() string { return "Edit on Github" }
func (e editOnGithub) Link() string {
	if len(editUrl) > 0 {
		return fmt.Sprintf("%s/%s", editUrl, e.page.FileName())
	} else {
		return fmt.Sprintf("%s/edit/%s/%s", repo, branch, e.page.FileName())
	}
}
func (e editOnGithub) OnClick() template.JS  { return "" }
func (e editOnGithub) Widget() template.HTML { return "" }
