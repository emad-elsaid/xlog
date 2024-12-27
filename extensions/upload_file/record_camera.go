package upload_file

import (
	"fmt"
	"html/template"
	"net/url"

	"github.com/emad-elsaid/xlog"
)

type RecordCamera struct {
	p xlog.Page
}

func (RecordCamera) Icon() string { return "fa-solid fa-video" }
func (RecordCamera) Name() string { return "Record camera" }
func (s RecordCamera) Attrs() map[template.HTMLAttr]any {
	link := fmt.Sprintf("/+/upload-file/record-camera-form?page=%s", url.PathEscape(s.p.Name()))

	return map[template.HTMLAttr]any{
		"href":      link,
		"hx-post":   link,
		"hx-target": "body",
		"hx-swap":   "beforeend",
	}
}

func RecordCameraForm(r xlog.Request) xlog.Output {
	name := r.FormValue("page")

	return xlog.Render("record-camera", map[string]any{
		"action": "/+/upload-file?page=" + url.QueryEscape(name),
		"csrf":   xlog.CSRF(r),
	})
}
