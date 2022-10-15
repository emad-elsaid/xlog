package upload_file

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	. "github.com/emad-elsaid/xlog"
)

const MAX_FILE_UPLOAD = 1 * GB

var (
	IMAGES_EXTENSIONS = []string{".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp"}
	VIDEOS_EXTENSIONS = []string{".webm"}
	AUDIO_EXTENSIONS  = []string{".wave", ".ogg", ".opus", ".mp3"}
)

func init() {
	WIDGET(TOOLS_WIDGET, uploadFileWidget)
	POST(`/\+/upload-file`, uploadFileHandler)
}

func uploadFileWidget(p *Page, r Request) template.HTML {
	return template.HTML(
		Partial("extension/upload-file", Locals{
			"page":           p,
			"csrf":           CSRF(r),
			"action":         "/+/upload-file?page=" + url.QueryEscape(p.Name),
			"editModeAction": "/+/upload-file",
		}),
	)
}

func uploadFileHandler(w Response, r Request) Output {
	r.ParseMultipartForm(MAX_FILE_UPLOAD)

	fileName := r.FormValue("page")

	page := NewPage(fileName)
	if fileName != "" && !page.Exists() {
		return Redirect("/" + page.Name + "/edit")
	}

	var output string
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

		if containString(IMAGES_EXTENSIONS, ext) {
			output = fmt.Sprintf("![](/%s)", p)
		} else if containString(VIDEOS_EXTENSIONS, ext) {
			output = fmt.Sprintf("<video controls src=\"/%s\"></video>", p)
		} else if containString(AUDIO_EXTENSIONS, ext) {
			output = fmt.Sprintf("<audio controls src=\"/%s\"></audio>", p)
		} else {
			output = fmt.Sprintf("[%s](/%s)", mdName, p)
		}
	}

	if fileName != "" && page.Exists() {
		content := strings.TrimSpace(page.Content()) + "\n\n" + output + "\n"
		page.Write(content)
		return Redirect("/" + page.Name)
	}

	return PlainText(output)
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
