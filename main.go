package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const MAX_FILE_UPLOAD = 50 * MB

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
			return Redirect("/" + page.Name() + "/edit")
		}

		html, refs := page.Render()
		refsIn := Search(page.name)

		return Render("view", Locals{
			"edit":         "/" + page.Name() + "/edit",
			"title":        page.Name(),
			"content":      template.HTML(html),
			"references":   refs,
			"referencedIn": refsIn,
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

				if strings.Contains(".jpg,.jpeg,.png,.gif", ext) {
					content += fmt.Sprintf("\n![](/%s)\n", p)
				} else {
					content += fmt.Sprintf("\n[%s](/%s)\n", p, p)
				}
			}

			page.Write(content)
			return Redirect("/" + page.Name())
		} else if page.Exists() {
			page.Delete()
		}

		return Redirect("/")
	})

	GET("/{page}/edit", func(w Response, r Request) Output {
		vars := VARS(r)
		page := NewPage(vars["page"])

		return Render("edit", Locals{
			"action":  page.Name(),
			"content": page.Content(),
			"csrf":    CSRF(r),
		})
	})

	Start()
}
