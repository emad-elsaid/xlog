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
		attrHref:   link,
		attrHxPost: link,
	}
}

// ScreenshotForm renders the screenshot capture interface for a specified page.
// It accepts a page name via the "page" form parameter and returns an HTML interface
// that allows users to capture and save screenshots directly to the page.
func ScreenshotForm(r xlog.Request) xlog.Output {
	name := r.FormValue("page")

	return xlog.Render("screenshot", map[string]any{
		attrAction: "/+/upload-file?page=" + url.QueryEscape(name),
		attrCSRF:   xlog.CSRF(r),
	})
}
