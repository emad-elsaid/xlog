package upload_file

import (
	"fmt"
	"html/template"
	"net/url"

	"github.com/emad-elsaid/xlog"
)

type Screenshot struct {
	p xlog.Page
}

func (Screenshot) Icon() string { return "fa-solid fa-camera" }
func (Screenshot) Name() string { return "Screenshot" }
func (s Screenshot) Attrs() map[template.HTMLAttr]any {
	link := fmt.Sprintf("/+/upload-file/screenshot-form?page=%s", url.PathEscape(s.p.Name()))

	return map[template.HTMLAttr]any{
		"href":    link,
		"hx-post": link,
	}
}

func ScreenshotForm(r xlog.Request) xlog.Output {
	name := r.FormValue("page")

	return xlog.Render("screenshot", map[string]any{
		"action": "/+/upload-file?page=" + url.QueryEscape(name),
		"csrf":   xlog.CSRF(r),
	})
}
