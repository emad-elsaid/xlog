package gpg

import (
	"html/template"
	"net/url"
	"path"

	"github.com/emad-elsaid/xlog"
)

const decryptableExt = ".pgp"

func commands(p xlog.Page) []xlog.Command {
	if !p.Exists() {
		return nil
	}

	if len(gpgId) == 0 {
		return nil
	}

	if path.Ext(p.FileName()) == decryptableExt {
		return []xlog.Command{
			&decryptCommand{page: p},
		}
	} else {
		return []xlog.Command{
			&encryptCommand{page: p},
		}
	}
}

type encryptCommand struct {
	page xlog.Page
}

func (e *encryptCommand) Icon() string { return "fa-solid fa-lock" }
func (e *encryptCommand) Name() string { return "Make private" }
func (e *encryptCommand) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"hx-post": "/+/gpg/encrypt/" + url.PathEscape(e.page.Name()),
	}
}

type decryptCommand struct {
	page xlog.Page
}

func (e *decryptCommand) Icon() string { return "fa-solid fa-lock-open has-text-danger" }
func (e *decryptCommand) Name() string { return "Make public" }
func (e *decryptCommand) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"hx-post": "/+/gpg/decrypt/" + url.PathEscape(e.page.Name()),
	}
}
