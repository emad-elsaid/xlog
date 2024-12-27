package upload_file

import (
	"fmt"
	"html/template"
	"net/url"

	"github.com/emad-elsaid/xlog"
)

type RecordAudio struct {
	p xlog.Page
}

func (RecordAudio) Icon() string { return "fa-solid fa-microphone" }
func (RecordAudio) Name() string { return "Record audio" }
func (s RecordAudio) Attrs() map[template.HTMLAttr]any {
	link := fmt.Sprintf("/+/upload-file/record-audio-form?page=%s", url.PathEscape(s.p.Name()))

	return map[template.HTMLAttr]any{
		"href":      link,
		"hx-post":   link,
		"hx-target": "body",
		"hx-swap":   "beforeend",
	}
}

func RecordAudioForm(r xlog.Request) xlog.Output {
	name := r.FormValue("page")

	return xlog.Render("record-audio", map[string]any{
		"action": "/+/upload-file?page=" + url.QueryEscape(name),
		"csrf":   xlog.CSRF(r),
	})
}
