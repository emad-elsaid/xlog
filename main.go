package main

import (
	"bytes"
	"crypto/sha256"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

const MAX_FILE_UPLOAD = 50 * MB

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

		html, refs := page.Render()

		return Render("view", Locals{
			"edit":         "/" + page.Name + "/edit",
			"title":        page.Name,
			"content":      template.HTML(html),
			"references":   refs,
			"referencedIn": Search(page.Name),
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
				p := path.Join("public", name)
				mdName := filterChars(h.Filename, "[]")

				os.Mkdir("public", 0700)
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

	Start()
}

func Search(keyword string) []string {
	pages := []string{}
	files, _ := ioutil.ReadDir(".")
	sort.Sort(fileInfoByNameLength(files))

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			f, err := ioutil.ReadFile(file.Name())
			if err != nil {
				continue
			}

			basename := file.Name()[:len(file.Name())-3]
			if strings.Contains(string(f), keyword) {
				pages = append(pages, basename)
			}
		}
	}

	return pages
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
