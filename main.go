package main

import (
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

const MAX_FILE_UPLOAD = 50 * MB
const MIN_SEARCH_KEYWORD = 3

var IMAGES_EXTENSIONS = []string{".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp"}

func main() {
	cwd, _ := os.Getwd()
	source := flag.String("source", cwd, "Directory that will act as a storage")
	flag.Parse()

	absSource, err := filepath.Abs(*source)
	if err != nil {
		log.Fatal(err)
	}

	os.Chdir(absSource)

	GET("/", func(w Response, r Request) Output {
		return Redirect("/index")
	})

	GET("/{page}", func(w Response, r Request) Output {
		vars := VARS(r)
		page := NewPage(vars["page"])

		if !page.Exists() {
			return Redirect("/" + page.Name + "/edit")
		}

		html := page.Render()

		return Render("view", Locals{
			"edit":    "/" + page.Name + "/edit",
			"title":   page.Name,
			"content": template.HTML(html),
		})
	})

	POST("/{page}", func(w Response, r Request) Output {
		r.ParseMultipartForm(MAX_FILE_UPLOAD)

		vars := VARS(r)
		page := NewPage(vars["page"])
		content := r.FormValue("content")

		if content != "" {
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
					content += fmt.Sprintf("\n![](/%s)\n", p)
				} else {
					content += fmt.Sprintf("\n[%s](/%s)\n", mdName, p)
				}
			}

			page.Write(content)
			return Redirect("/" + page.Name)
		} else if page.Exists() {
			page.Delete()
		}

		return Redirect("/")
	})

	GET("/{page}/edit", func(w Response, r Request) Output {
		vars := VARS(r)
		page := NewPage(vars["page"])

		return Render("edit", Locals{
			"action":  page.Name,
			"content": page.Content(),
			"csrf":    CSRF(r),
		})
	})

	HELPER("navbarStart", func() template.HTML {
		o := template.HTML("")
		for _, v := range NAVBAR_START_WIDGETS {
			o += v()
		}
		return o
	})

	Start()
}

func renderMarkdown(content string) string {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			highlighting.Highlighting,
			emoji.Emoji,
		),

		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return err.Error()
	}

	return buf.String()
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

var NAVBAR_START_WIDGETS = []func() template.HTML{}

func NAVBAR_START(f func() template.HTML) {
	NAVBAR_START_WIDGETS = append(NAVBAR_START_WIDGETS, f)
}
