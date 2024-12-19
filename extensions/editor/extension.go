package editor

import (
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/emad-elsaid/xlog"
)

var editor string

func init() {
	flag.StringVar(&editor, "editor", os.Getenv("EDITOR"), "command to use to open pages for editing")

	xlog.RegisterExtension(Editor{})
}

type Editor struct{}

func (Editor) Name() string { return "editor" }
func (Editor) Init() {
	if xlog.Config.Readonly {
		return
	}

	xlog.RegisterQuickCommand(links)
	xlog.Post(`/+/editor/{page...}`, editorHandler)
	xlog.Listen(xlog.PageNotFound, newPage)
}

func newPage(p xlog.Page) error {
	openEditor(p)

	return nil
}

var ignoredPages = []string{"favicon.ico"}

func openEditor(page xlog.Page) {
	if slices.Contains(ignoredPages, page.Name()) {
		return
	}

	segments := strings.Split(editor, " ")
	if len(segments) == 0 {
		return
	}

	name := segments[0]
	args := append(segments[1:], page.FileName())
	cmd := exec.Command(name, args...)

	if err := cmd.Start(); err != nil {
		slog.Error("Error start command", "command", cmd.String(), "error", err)
	}
}

func editorHandler(r xlog.Request) xlog.Output {
	page := xlog.NewPage(r.PathValue("page"))
	slog.Info("Editing page", "name", page)

	openEditor(page)

	return xlog.NoContent()
}

func links(p xlog.Page) []xlog.Command {
	return []xlog.Command{
		editButton{page: p},
	}
}

type editButton struct {
	page xlog.Page
}

func (editButton) Icon() string {
	return "fa-solid fa-pen"
}
func (editButton) Name() string {
	return "Edit"
}

func (editButton) Link() string { return "" }
func (e editButton) OnClick() template.JS {
	action := fmt.Sprintf("/+/editor/%s", url.PathEscape(e.page.Name()))
	script := `
     const data = new FormData()
     data.append('csrf', document.querySelector('input[name=csrf]').value);
     fetch("%s", {method: 'POST', body: data});
`
	return template.JS(fmt.Sprintf(script, action))
}
func (editButton) Widget() template.HTML { return "" }
