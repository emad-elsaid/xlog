package main

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

const MAX_FILE_UPLOAD = 50 * MB

var IMAGES_EXTENSIONS = []string{".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp"}
var VIDEOS_EXTENSIONS = []string{".webm"}

func init() {
	WIDGET(TOOLS_WIDGET, uploadFileWidget)
	POST(`/\+/upload-file/{page}`, uploadFileHandler)
}

func uploadFileWidget(p *Page, r Request) template.HTML {
	return template.HTML(
		partial("extension/upload-file", Locals{
			"page":   p,
			"csrf":   CSRF(r),
			"action": "/+/upload-file/" + p.Name,
		}),
	)
}

func uploadFileHandler(w Response, r Request) Output {
	r.ParseMultipartForm(MAX_FILE_UPLOAD)

	vars := VARS(r)
	page := NewPage(vars["page"])

	if !page.Exists() {
		return Redirect("/" + page.Name + "/edit")
	}

	content := page.Content()
	f, h, _ := r.FormFile("file")
	if f != nil {
		defer f.Close()
		c, _ := io.ReadAll(f)
		ext := strings.ToLower(path.Ext(h.Filename))
		name := fmt.Sprintf("%x%s", sha256.Sum256(c), ext)
		p := path.Join(STATIC_DIR_PATH, name)
		mdName := filterChars(h.Filename, "[]")

		os.Mkdir(STATIC_DIR_PATH, 0700)
		out, err := os.Create(p)
		if err != nil {
			return InternalServerError(err)
		}

		f.Seek(io.SeekStart, 0)
		_, err = io.Copy(out, f)
		if err != nil {
			return InternalServerError(err)
		}

		content = strings.TrimSpace(content)

		if containString(IMAGES_EXTENSIONS, ext) {
			content += fmt.Sprintf("\n\n![](/%s)\n", p)
		} else if containString(VIDEOS_EXTENSIONS, ext) {
			content += fmt.Sprintf("\n\n<video controls src=\"%s\"></video>\n", p)
		} else {
			content += fmt.Sprintf("\n\n[%s](/%s)\n", mdName, p)
		}
	}

	page.Write(content)

	return Redirect("/" + page.Name)
}

func containString(slice []string, str string) bool {
	for k := range slice {
		if slice[k] == str {
			return true
		}
	}

	return false
}

func filterChars(str string, exclude string) string {
	pattern := regexp.MustCompile("[" + regexp.QuoteMeta(exclude) + "]")

	return pattern.ReplaceAllString(str, "")
}
