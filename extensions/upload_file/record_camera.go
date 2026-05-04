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
		attrHref:     link,
		attrHxPost:   link,
		attrHxTarget: targetBody,
		attrHxSwap:   swapEnd,
	}
}

// RecordCameraForm renders the camera recording interface for a specified page.
// It accepts a page name via the "page" form parameter and returns an HTML interface
// that allows users to record video using their camera and save it to the page.
func RecordCameraForm(r xlog.Request) xlog.Output {
	name := r.FormValue("page")

	return xlog.Render("record-camera", map[string]any{
		attrAction: "/+/upload-file?page=" + url.QueryEscape(name),
		attrCSRF:   xlog.CSRF(r),
	})
}
