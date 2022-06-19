package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"
)

const TEMPLATE_NAME = "template"

func main() {
	cwd, _ := os.Getwd()
	source := flag.String("source", cwd, "Directory that will act as a storage")
	flag.Parse()

	absSource, err := filepath.Abs(*source)
	if err != nil {
		log.Fatal(err)
	}

	if err = os.Chdir(absSource); err != nil {
		log.Fatal(err)
	}

	GET("/", RootHandler)
	GET("/edit/{page:.*}", GetPageEditHandler)
	GET("/{page:.*}", GetPageHandler)
	POST("/{page:.*}", PostPageHandler)

	START()
}

func RootHandler(w Response, r Request) Output {
	return Redirect("/index")
}

func GetPageHandler(w Response, r Request) Output {
	vars := VARS(r)
	page := NewPage(vars["page"])

	if !page.Exists() {
		return Redirect("/edit/" + page.Name)
	}

	return Render("view", Locals{
		"edit":      "/edit/" + page.Name,
		"title":     page.Name,
		"updated":   ago(time.Now().Sub(page.ModTime())),
		"content":   template.HTML(page.Render()),
		"tools":     renderWidget(TOOLS_WIDGET, &page, r),
		"sidebar":   renderWidget(SIDEBAR_WIDGET, &page, r),
		"meta":      renderWidget(META_WIDGET, &page, r),
		"afterView": renderWidget(AFTER_VIEW_WIDGET, &page, r),
	})
}

func GetPageEditHandler(w Response, r Request) Output {
	vars := VARS(r)
	page := NewPage(vars["page"])

	var content string
	if page.Exists() {
		content = page.Content()
	} else if template := NewPage(TEMPLATE_NAME); template.Exists() {
		content = template.Content()
	}

	return Render("edit", Locals{
		"title":   page.Name,
		"action":  page.Name,
		"rtl":     page.RTL(),
		"tools":   renderWidget(TOOLS_WIDGET, &page, r),
		"content": content,
		"csrf":    CSRF(r),
	})
}

func PostPageHandler(w Response, r Request) Output {
	vars := VARS(r)
	page := NewPage(vars["page"])
	content := r.FormValue("content")
	page.Write(content)

	return Redirect("/" + page.Name)
}

// WIDGETS ===================================================

type (
	widgetSpace int
	widgetFunc  func(*Page, Request) template.HTML
)

const (
	TOOLS_WIDGET widgetSpace = iota
	SIDEBAR_WIDGET
	AFTER_VIEW_WIDGET
	META_WIDGET
)

var widgets = map[widgetSpace][]widgetFunc{}

func PREPEND_WIDGET(s widgetSpace, f func(*Page, Request) template.HTML) {
	if _, ok := widgets[s]; !ok {
		widgets[s] = []widgetFunc{}
	}
	widgets[s] = append([]widgetFunc{f}, widgets[s]...)
}

func WIDGET(s widgetSpace, f func(*Page, Request) template.HTML) {
	if _, ok := widgets[s]; !ok {
		widgets[s] = []widgetFunc{}
	}
	widgets[s] = append(widgets[s], f)
}

func renderWidget(s widgetSpace, p *Page, r Request) (o template.HTML) {
	for _, v := range widgets[s] {
		o += v(p, r)
	}
	return
}

// HELPERS

func ago(t time.Duration) (o string) {
	const day = time.Hour * 24
	const week = day * 7
	const month = day * 30
	const year = day * 365
	const maxPrecision = 2

	if t.Seconds() < 1 {
		return "seconds ago"
	}

	for precision := 0; t.Seconds() > 0 && precision < maxPrecision; precision++ {
		switch {
		case t >= year:
			years := t / year
			t -= years * year
			o += fmt.Sprintf("%d years ", years)
		case t >= month:
			months := t / month
			t -= months * month
			o += fmt.Sprintf("%d months ", months)
		case t >= week:
			weeks := t / week
			t -= weeks * week
			o += fmt.Sprintf("%d weeks ", weeks)
		case t >= day:
			days := t / day
			t -= days * day
			o += fmt.Sprintf("%d days ", days)
		case t >= time.Hour:
			hours := t / time.Hour
			t -= hours * time.Hour
			o += fmt.Sprintf("%d hours ", hours)
		case t >= time.Minute:
			minutes := t / time.Minute
			t -= minutes * time.Minute
			o += fmt.Sprintf("%d minutes ", minutes)
		case t >= time.Second:
			seconds := t / time.Second
			t -= seconds * time.Second
			o += fmt.Sprintf("%d seconds ", seconds)
		}
	}

	return o + "ago"
}
