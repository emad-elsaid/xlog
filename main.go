package main

import (
	"bytes"
	"flag"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func main() {
	cwd, _ := os.Getwd()
	source := flag.String("source", cwd, "Directory that will act as a storage")
	flag.Parse()

	absSource, err := filepath.Abs(*source)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(absSource)
	if err != nil {
		log.Fatal(err)
	}

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
		tools := template.HTML("")
		for _, v := range TOOLS_WIDGETS {
			tools += v(&page, r)
		}
		sidebar := template.HTML("")
		for _, v := range SIDEBAR_WIDGETS {
			sidebar += v(&page, r)
		}

		return Render("view", Locals{
			"edit":    "/" + page.Name + "/edit",
			"title":   page.Name,
			"updated": page.ModTime().Format("2006-01-02 15:04"),
			"content": template.HTML(html),
			"tools":   tools,
			"sidebar": sidebar,
		})
	})

	POST("/{page}", func(w Response, r Request) Output {
		vars := VARS(r)
		page := NewPage(vars["page"])
		content := r.FormValue("content")

		page.Write(content)
		return Redirect("/" + page.Name)
	})

	GET("/{page}/edit", func(w Response, r Request) Output {
		vars := VARS(r)
		page := NewPage(vars["page"])

		return Render("edit", Locals{
			"title":   page.Name,
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

// WIDGETS ===================================================

var NAVBAR_START_WIDGETS = []func() template.HTML{}

func NAVBAR_START(f func() template.HTML) {
	NAVBAR_START_WIDGETS = append(NAVBAR_START_WIDGETS, f)
}

var TOOLS_WIDGETS = []func(*Page, Request) template.HTML{}

func TOOL(f func(*Page, Request) template.HTML) {
	TOOLS_WIDGETS = append(TOOLS_WIDGETS, f)
}

var SIDEBAR_WIDGETS = []func(*Page, Request) template.HTML{}

func SIDEBAR(f func(*Page, Request) template.HTML) {
	SIDEBAR_WIDGETS = append(SIDEBAR_WIDGETS, f)
}
