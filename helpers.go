package xlog

import (
	"errors"
	"fmt"
	"html/template"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/emad-elsaid/xlog/markdown/ast"
	emojiAst "github.com/emad-elsaid/xlog/markdown/emoji/ast"
)

var ErrHelperRegistered = errors.New("Helper already registered")

// RegisterHelper registers a new helper function. all helpers are used when compiling
// templates. so registering helpers function must happen before the server
// starts as compiling templates happened right before starting the http server.
func (app *App) RegisterHelper(name string, f any) error {
	if _, ok := app.helpers[name]; ok {
		return ErrHelperRegistered
	}

	app.helpers[name] = f
	return nil
}

// A function that takes time.duration and return a string representation of the
// duration in human readable way such as "3 seconds ago". "5 hours 30 minutes
// ago". The precision of this function is 2. which means it returns the largest
// unit of time possible and the next one after it. for example days + hours, or
// Hours + minutes or Minutes + seconds...etc
func (app *App) ago(t time.Time) string {
	if app.config.Readonly {
		return t.Format("Monday 2 January 2006")
	}

	d := time.Since(t)

	const day = time.Hour * 24
	const week = day * 7
	const month = day * 30
	const year = day * 365
	const maxPrecision = 2

	var o strings.Builder

	if d.Seconds() < 1 {
		o.WriteString("Less than a second ")
	}

	for precision := 0; d.Seconds() > 1 && precision < maxPrecision; precision++ {
		switch {
		case d >= year:
			years := d / year
			d -= years * year
			o.WriteString(fmt.Sprintf("%d years ", years))
		case d >= month:
			months := d / month
			d -= months * month
			o.WriteString(fmt.Sprintf("%d months ", months))
		case d >= week:
			weeks := d / week
			d -= weeks * week
			o.WriteString(fmt.Sprintf("%d weeks ", weeks))
		case d >= day:
			days := d / day
			d -= days * day
			o.WriteString(fmt.Sprintf("%d days ", days))
		case d >= time.Hour:
			hours := d / time.Hour
			d -= hours * time.Hour
			o.WriteString(fmt.Sprintf("%d hours ", hours))
		case d >= time.Minute:
			minutes := d / time.Minute
			d -= minutes * time.Minute
			o.WriteString(fmt.Sprintf("%d minutes ", minutes))
		case d >= time.Second:
			seconds := d / time.Second
			d -= seconds * time.Second
			o.WriteString(fmt.Sprintf("%d seconds ", seconds))
		}
	}

	o.WriteString("ago")

	return o.String()
}

// RegisterJS registers a JavaScript file to be included in the page
func (app *App) RegisterJS(f string) {
	app.includeJS(f)
}

// RequireHTMX registers HTMX library
func (app *App) RequireHTMX() {
	app.includeJS("/public/htmx.min.js")
}

// includeJS adds a JavaScript library URL/path
func (app *App) includeJS(f string) template.HTML {
	if !slices.Contains(app.js, f) {
		app.js = append(app.js, f)
	}
	return ""
}

// scripts returns the HTML for all registered JavaScript files
func (app *App) scripts() template.HTML {
	var b strings.Builder
	for _, f := range app.js {
		fmt.Fprintf(&b, `<script src="%s" defer></script>`, f)
	}
	return template.HTML(b.String())
}

// IsFontAwesome checks if an icon is a FontAwesome icon
func (app *App) IsFontAwesome(i string) bool {
	return strings.HasPrefix(i, "fa")
}

// Banner returns the banner image for a page
func (app *App) Banner(p Page) string {
	_, a := p.AST()
	if a == nil {
		return ""
	}

	paragraph := a.FirstChild()
	if paragraph == nil || paragraph.Kind() != ast.KindParagraph {
		return ""
	}

	img := paragraph.FirstChild()
	if img == nil || img.Kind() != ast.KindImage {
		return ""
	}

	image, ok := img.(*ast.Image)
	if !ok {
		return ""
	}

	dest := string(image.Destination)
	if len(dest) == 0 || dest == "#" {
		return ""
	}

	if !(path.IsAbs(dest) || strings.HasPrefix(dest, "http")) {
		d := path.Dir(p.FileName())
		dest = path.Join("/", d, dest)
	}

	return dest
}

// Emoji returns the emoji for a page
func (app *App) Emoji(p Page) string {
	_, tree := p.AST()
	if e, ok := FindInAST[*emojiAst.Emoji](tree); ok && e != nil {
		return string(e.Value.Unicode)
	}
	return ""
}

// dir returns the directory name
func (app *App) dir(s string) string {
	v := path.Dir(s)
	if v == "." {
		return ""
	}
	return v
}

// raw returns safe HTML
func (app *App) raw(i string) template.HTML {
	return template.HTML(i)
}
