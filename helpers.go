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

var helpers = template.FuncMap{
	"ago":            ago,
	"properties":     Properties,
	"links":          Links,
	"widgets":        RenderWidget,
	"commands":       Commands,
	"quick_commands": QuickCommands,
	"isFontAwesome":  IsFontAwesome,
	"includeJS":      includeJS,
	"scripts":        scripts,
	"banner":         Banner,
	"emoji":          Emoji,
	"base":           path.Base,
	"dir":            dir,
	"raw":            raw,
}

var ErrHelperRegistered = errors.New("helper already registered")

// RegisterHelper registers a new helper function. all helpers are used when compiling
// templates. so registering helpers function must happen before the server
// starts as compiling templates happened right before starting the http server.
func RegisterHelper(name string, f any) error {
	if _, ok := helpers[name]; ok {
		return ErrHelperRegistered
	}

	helpers[name] = f

	return nil
}

// A function that takes time.duration and return a string representation of the
// duration in human readable way such as "3 seconds ago". "5 hours 30 minutes
// ago". The precision of this function is 2. which means it returns the largest
// unit of time possible and the next one after it. for example days + hours, or
// Hours + minutes or Minutes + seconds...etc.
func ago(t time.Time) string {
	if Config.Readonly {
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
			fmt.Fprintf(&o, "%d years ", years)
		case d >= month:
			months := d / month
			d -= months * month
			fmt.Fprintf(&o, "%d months ", months)
		case d >= week:
			weeks := d / week
			d -= weeks * week
			fmt.Fprintf(&o, "%d weeks ", weeks)
		case d >= day:
			days := d / day
			d -= days * day
			fmt.Fprintf(&o, "%d days ", days)
		case d >= time.Hour:
			hours := d / time.Hour
			d -= hours * time.Hour
			fmt.Fprintf(&o, "%d hours ", hours)
		case d >= time.Minute:
			minutes := d / time.Minute
			d -= minutes * time.Minute
			fmt.Fprintf(&o, "%d minutes ", minutes)
		case d >= time.Second:
			seconds := d / time.Second
			d -= seconds * time.Second
			fmt.Fprintf(&o, "%d seconds ", seconds)
		}
	}

	o.WriteString("ago")

	return o.String()
}

var js = []string{}

// RegisterJS adds a Javascript library URL/path to be included in the scripts used by the template.
func RegisterJS(f string) {
	if slices.Contains(js, f) {
		return
	}

	js = append(js, f)
}

// RequireHTMX registes HTML library, this helps include one version of HTMX.
func RequireHTMX() {
	RegisterJS("/public/htmx.min.js")
}

func includeJS(f string) template.HTML {
	RegisterJS(f)

	return ""
}

func scripts() template.HTML {
	var b strings.Builder
	for _, f := range js {
		fmt.Fprintf(&b, `<script src="%s" defer></script>`, f)
	}

	// #nosec G203 - js slice contains trusted paths from Config.Assets, safe to convert
	return template.HTML(b.String())
}

// IsFontAwesome checks if the given string is a FontAwesome icon class.
// It returns true if the string starts with "fa", which is the standard FontAwesome prefix.
func IsFontAwesome(i string) bool {
	return strings.HasPrefix(i, "fa")
}

// Banner extracts the banner image URL from a page.
// It looks for the first image in the first paragraph of the page's AST.
// Returns an empty string if no banner image is found or if the image destination is invalid.
// Relative image paths are converted to absolute paths based on the page's directory.
func Banner(p Page) string {
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

	if !path.IsAbs(dest) && !strings.HasPrefix(dest, "http") {
		d := path.Dir(p.FileName())
		dest = path.Join("/", d, dest)
	}

	return dest
}

// Emoji extracts the first emoji from a page's markdown content.
// It searches the page's AST for an Emoji node and returns its Unicode representation.
// Returns an empty string if no emoji is found.
func Emoji(p Page) string {
	_, tree := p.AST()
	if e, ok := FindInAST[*emojiAst.Emoji](tree); ok && e != nil {
		return string(e.Value.Unicode)
	}

	return ""
}

func dir(s string) string {
	v := path.Dir(s)

	if v == "." {
		return ""
	}

	return v
}

// raw a helper to output input string as safe HTML.
// WARNING: This intentionally bypasses HTML escaping. Only use with trusted content.
// #nosec G203 - intentional unsafe conversion for trusted template content
func raw(i string) template.HTML {
	return template.HTML(i)
}
