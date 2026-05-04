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
		attrHref:     link,
		attrHxPost:   link,
		attrHxTarget: targetBody,
		attrHxSwap:   swapEnd,
	}
}

// RecordAudioForm renders the audio recording interface for a specified page.
// It accepts a page name via the "page" form parameter and returns an HTML interface
// that allows users to record audio using their microphone and save it to the page.
func RecordAudioForm(r xlog.Request) xlog.Output {
	name := r.FormValue("page")

	return xlog.Render("record-audio", map[string]any{
		attrAction: "/+/upload-file?page=" + url.QueryEscape(name),
		attrCSRF:   xlog.CSRF(r),
	})
}
