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

	"github.com/emad-elsaid/xlog"
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
	xlog.RegisterExtension(UploadFile{})
}

type UploadFile struct{}

func (UploadFile) Name() string { return "upload-file" }
func (UploadFile) Init() {
	if xlog.Config.Readonly {
		return
	}

	xlog.RequireHTMX()
	xlog.RegisterCommand(func(p xlog.Page) []xlog.Command {
		if !p.Exists() {
			return nil
		}

		return []xlog.Command{
			Upload{p: p},
			Screenshot{p: p},
			RecordScreen{p: p},
			RecordCamera{p: p},
			RecordAudio{p: p},
		}
	})

	xlog.Post("/+/upload-file/form", UploadForm)
	xlog.Post("/+/upload-file/screenshot-form", ScreenshotForm)
	xlog.Post("/+/upload-file/record-screen-form", RecordScreenForm)
	xlog.Post("/+/upload-file/record-camera-form", RecordCameraForm)
	xlog.Post("/+/upload-file/record-audio-form", RecordAudioForm)

	xlog.Post(`/+/upload-file`, uploadFileHandler)
	xlog.RegisterTemplate(templates, "templates")
}

func uploadFileHandler(r xlog.Request) xlog.Output {
	if err := r.ParseMultipartForm(MAX_FILE_UPLOAD); err != nil {
		return xlog.BadRequest(err.Error())
	}

	fileName := r.FormValue("page")

	page := xlog.NewPage(fileName)
	if page == nil || (fileName != "" && !page.Exists()) {
		return xlog.NotFound("page not found")
	}

	var output string
	f, h, _ := r.FormFile("file")
	if f != nil && h != nil {
		defer func() { _ = f.Close() }()
		c, _ := io.ReadAll(f)
		ext := strings.ToLower(path.Ext(h.Filename))
		name := fmt.Sprintf("%x%s", sha256.Sum256(c), ext)
		p := path.Join(PUBLIC_PATH, name)
		mdName := filterChars(h.Filename, "[]")

		if err := os.Mkdir(PUBLIC_PATH, 0700); err != nil && !os.IsExist(err) {
			return xlog.InternalServerError(err)
		}
		out, err := os.Create(p) // #nosec G304 -- Path is Join(PUBLIC_PATH, sha256+ext); sha256 prevents traversal (never starts with ..)
		if err != nil {
			return xlog.InternalServerError(err)
		}

		if _, err := f.Seek(0, io.SeekStart); err != nil {
			return xlog.InternalServerError(err)
		}
		_, err = io.Copy(out, f)
		if err != nil {
			return xlog.InternalServerError(err)
		}

		switch {
		case slices.Contains(IMAGES_EXTENSIONS, ext):
			output = fmt.Sprintf("![](/%s)", p)
		case slices.Contains(VIDEOS_EXTENSIONS, ext):
			output = fmt.Sprintf("<video controls src=\"/%s\"></video>", p)
		case slices.Contains(AUDIO_EXTENSIONS, ext):
			output = fmt.Sprintf("<audio controls src=\"/%s\"></audio>", p)
		default:
			output = fmt.Sprintf("[%s](/%s)", mdName, p)
		}
	}

	if fileName != "" && page.Exists() {
		content := strings.TrimSpace(string(page.Content())) + "\n\n" + output + "\n"
		page.Write(xlog.Markdown(content))
		return xlog.Redirect("/" + page.Name())
	}

	return xlog.PlainText(output)
}

func filterChars(str string, exclude string) string {
	pattern := regexp.MustCompile("[" + regexp.QuoteMeta(exclude) + "]")

	return pattern.ReplaceAllString(str, "")
}
