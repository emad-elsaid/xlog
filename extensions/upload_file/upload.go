package upload_file

import (
	"fmt"
	"html/template"
	"net/url"

	"github.com/emad-elsaid/xlog"
)

type Upload struct {
	p xlog.Page
}

func (Upload) Icon() string { return "fa-solid fa-file-arrow-up" }
func (Upload) Name() string { return "Upload File" }
func (u Upload) Attrs() map[template.HTMLAttr]any {
	link := fmt.Sprintf("/+/upload-file/form?page=%s", url.PathEscape(u.p.Name()))

	return map[template.HTMLAttr]any{
		attrHref:     link,
		attrHxPost:   link,
		attrHxTarget: targetBody,
		attrHxSwap:   swapEnd,
	}
}

func UploadForm(r xlog.Request) xlog.Output {
	name := r.FormValue("page")

	return xlog.Render("upload", map[string]any{
		attrAction: "/+/upload-file?page=" + url.QueryEscape(name),
		attrCSRF:   xlog.CSRF(r),
	})
}
