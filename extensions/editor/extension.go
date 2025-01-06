package editor

import (
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/emad-elsaid/xlog"
)

var editor string

func init() {
	flag.StringVar(&editor, "editor", os.Getenv("EDITOR"), "command to use to open pages for editing")

	xlog.RequireHTMX()
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

func openEditor(page xlog.Page) {
	if page == nil {
		return
	}

	// if it's like a .ico, .jpeg, .so...etc ignore it, it's not a page we
	// should create, maybe just a static file that's missing
	if ext := len(filepath.Ext(page.Name())); ext > 0 && ext <= 4 {
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
	if len(p.FileName()) == 0 {
		return nil
	}

	return []xlog.Command{editButton{page: p}}
}

type editButton struct {
	page xlog.Page
}

func (editButton) Icon() string { return "fa-solid fa-pen" }
func (editButton) Name() string { return "Edit" }
func (e editButton) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"hx-post": fmt.Sprintf("/+/editor/%s", url.PathEscape(e.page.Name())),
	}
}
