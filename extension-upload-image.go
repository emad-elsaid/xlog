package main

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"strings"
)

const MAX_FILE_UPLOAD = 50 * MB

var IMAGES_EXTENSIONS = []string{".jpg", ".jpeg", ".png", ".gif", ".svg", ".webp"}

func init() {
	SIDEBAR(uploadImageWidget)
	GET("/+/upload-image/{page}", uploadImagePageHandler)

	POST("/+/upload-image/{page}", uploadImageHandler)
}

func uploadImageWidget(p *Page) template.HTML {
	return template.HTML(
		partial("extension/upload-image-sidebar", Locals{
			"page": p,
		}),
	)
}

func uploadImagePageHandler(w Response, r Request) Output {
	vars := VARS(r)
	page := NewPage(vars["page"])

	if !page.Exists() {
		return Redirect("/" + page.Name + "/edit")
	}

	return Render("extension/upload-image", Locals{
		"title":  "Upload image",
		"page":   page,
		"action": "/+/upload-image/" + page.Name,
		"csrf":   CSRF(r),
	})
}

func uploadImageHandler(w Response, r Request) Output {
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

		if containString(IMAGES_EXTENSIONS, ext) {
			content += fmt.Sprintf("\n![](/%s)\n", p)
		} else {
			content += fmt.Sprintf("\n[%s](/%s)\n", mdName, p)
		}
	}

	page.Write(content)

	return Redirect("/" + page.Name)
}
