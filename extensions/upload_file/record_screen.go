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
		attrHref:   link,
		attrHxPost: link,
	}
}

// RecordScreenForm renders the screen recording interface for a specified page.
// It accepts a page name via the "page" form parameter and returns an HTML interface
// that allows users to record their screen and save the recording to the page.
func RecordScreenForm(r xlog.Request) xlog.Output {
	name := r.FormValue("page")

	return xlog.Render("record-screen", map[string]any{
		attrAction: "/+/upload-file?page=" + url.QueryEscape(name),
		attrCSRF:   xlog.CSRF(r),
	})
}
