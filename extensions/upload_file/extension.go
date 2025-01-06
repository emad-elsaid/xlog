package upload_file

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"slices"
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
	RegisterExtension(UploadFile{})
}

type UploadFile struct{}

func (UploadFile) Name() string { return "upload-file" }
func (UploadFile) Init() {
	if Config.Readonly {
		return
	}

	RequireHTMX()
	RegisterCommand(func(p Page) []Command {
		if !p.Exists() {
			return nil
		}

		return []Command{
			Upload{p: p},
			Screenshot{p: p},
			RecordScreen{p: p},
			RecordCamera{p: p},
			RecordAudio{p: p},
		}
	})

	Post("/+/upload-file/form", UploadForm)
	Post("/+/upload-file/screenshot-form", ScreenshotForm)
	Post("/+/upload-file/record-screen-form", RecordScreenForm)
	Post("/+/upload-file/record-camera-form", RecordCameraForm)
	Post("/+/upload-file/record-audio-form", RecordAudioForm)

	Post(`/+/upload-file`, uploadFileHandler)
	RegisterTemplate(templates, "templates")
}

func uploadFileHandler(r Request) Output {
	r.ParseMultipartForm(MAX_FILE_UPLOAD)

	fileName := r.FormValue("page")

	page := NewPage(fileName)
	if page == nil || (fileName != "" && !page.Exists()) {
		return NotFound("page not found")
	}

	var output string
	f, h, _ := r.FormFile("file")
	if f != nil && h != nil {
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

		if slices.Contains(IMAGES_EXTENSIONS, ext) {
			output = fmt.Sprintf("![](/%s)", p)
		} else if slices.Contains(VIDEOS_EXTENSIONS, ext) {
			output = fmt.Sprintf("<video controls src=\"/%s\"></video>", p)
		} else if slices.Contains(AUDIO_EXTENSIONS, ext) {
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

func filterChars(str string, exclude string) string {
	pattern := regexp.MustCompile("[" + regexp.QuoteMeta(exclude) + "]")

	return pattern.ReplaceAllString(str, "")
}
