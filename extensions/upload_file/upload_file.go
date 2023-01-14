package upload_file

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	_ "embed"

	. "github.com/emad-elsaid/xlog"
)

const gb = 1 << (10 * 3)
const MAX_FILE_UPLOAD = 1 * gb
const PUBLIC_PATH = "public"

//go:embed templates
var templates embed.FS

var (
	IMAGES_EXTENSIONS = []string{".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp"}
	VIDEOS_EXTENSIONS = []string{".webm"}
	AUDIO_EXTENSIONS  = []string{".wave", ".ogg", ".opus", ".mp3"}
)

func init() {
	RegisterCommand(func(p Page) []Command {
		if READONLY {
			return []Command{}
		}

		return []Command{
			command{
				page:    p,
				icon:    "fa-solid fa-file-arrow-up",
				name:    "Upload File",
				onClick: "upload(event)",
				main:    true,
			},
			command{
				page:    p,
				icon:    "fa-solid fa-camera",
				name:    "Screenshot",
				onClick: "screenshot(event)",
			},
			command{
				page:    p,
				icon:    "fa-solid fa-desktop",
				name:    "Record Screen",
				onClick: "record(event)",
			},
			command{
				page:    p,
				icon:    "fa-solid fa-video",
				name:    "Record Camera",
				onClick: "recordCamera(event)",
			},
			command{
				page:    p,
				icon:    "fa-solid fa-microphone",
				name:    "Record Audio",
				onClick: "recordAudio(event)",
			},
			command{
				page:    p,
				icon:    "fa-solid fa-pen",
				name:    "Toggle Vim",
				onClick: "",
			},
		}
	})

	Post(`/\+/upload-file`, uploadFileHandler)
	RegisterTemplate(templates, "templates")
}

type command struct {
	page    Page
	icon    string
	name    string
	onClick template.JS
	main    bool
}

func (u command) Icon() string {
	return u.icon
}

func (u command) Name() string {
	return u.name
}

func (u command) Link() string { return "" }

func (u command) OnClick() template.JS {
	return u.onClick
}

func (u command) Widget() template.HTML {
	if !u.main {
		return ""
	}

	return Partial("upload-file", Locals{
		"page":           u.page,
		"action":         "/+/upload-file?page=" + url.QueryEscape(u.page.Name()),
		"editModeAction": "/+/upload-file",
	})
}

func uploadFileHandler(w Response, r Request) Output {
	r.ParseMultipartForm(MAX_FILE_UPLOAD)

	fileName := r.FormValue("page")

	page := NewPage(fileName)
	if fileName != "" && !page.Exists() {
		return Redirect("/" + page.Name() + "/edit")
	}

	var output string
	f, h, _ := r.FormFile("file")
	if f != nil {
		defer f.Close()
		c, _ := io.ReadAll(f)
		ext := strings.ToLower(path.Ext(h.Filename))
		name := fmt.Sprintf("%x%s", sha256.Sum256(c), ext)
		p := path.Join(PUBLIC_PATH, name)
		mdName := filterChars(h.Filename, "[]")

		os.Mkdir(PUBLIC_PATH, 0700)
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
		content := strings.TrimSpace(string(page.Content())) + "\n\n" + output + "\n"
		page.Write(Markdown(content))
		return Redirect("/" + page.Name())
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
