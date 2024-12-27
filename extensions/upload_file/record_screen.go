package upload_file

import (
	"fmt"
	"html/template"
	"net/url"

	"github.com/emad-elsaid/xlog"
)

type RecordScreen struct {
	p xlog.Page
}

func (RecordScreen) Icon() string { return "fa-solid fa-desktop" }
func (RecordScreen) Name() string { return "Record screen" }
func (s RecordScreen) Attrs() map[template.HTMLAttr]any {
	link := fmt.Sprintf("/+/upload-file/record-screen-form?page=%s", url.PathEscape(s.p.Name()))

	return map[template.HTMLAttr]any{
		"href":    link,
		"hx-post": link,
	}
}

func RecordScreenForm(r xlog.Request) xlog.Output {
	name := r.FormValue("page")

	return xlog.Render("record-screen", map[string]any{
		"action": "/+/upload-file?page=" + url.QueryEscape(name),
		"csrf":   xlog.CSRF(r),
	})
}
