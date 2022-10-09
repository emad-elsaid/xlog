package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Template file name. its content is used as default value for the editor when
// creating new pages.
const TEMPLATE_NAME = "template"

var (
	READONLY bool
	SITENAME string
)

func main() {
	// Uses current working directory as default value for source flag. If the
	// source flag is set by user the program changes working directory to is
	// and the rest of the program can use relative paths to access files
	cwd, _ := os.Getwd()
	source := flag.String("source", cwd, "Directory that will act as a storage")
	build := flag.String("build", "", "Build all pages as static site in this directory")
	flag.StringVar(&SITENAME, "sitename", "XLOG", "Site name is the name that appears on the header beside the logo and in the title tag")
	flag.BoolVar(&READONLY, "readonly", false, "Should xlog hide write operations, read-only means all write operations will be disabled")
	flag.Parse()

	if len(*build) > 0 {
		READONLY = true
	}

	absSource, err := filepath.Abs(*source)
	if err != nil {
		log.Fatal(err)
	}

	if err = os.Chdir(absSource); err != nil {
		log.Fatal(err)
	}

	// Program Core routes. View, Edit routes and a route to write new content
	// to the page. + handling root path which just show `index` page.
	GET("/", RootHandler)

	if !READONLY {
		GET("/edit/{page:.*}", GetPageEditHandler)
	}

	GET("/{page:.*}", GetPageHandler)

	if !READONLY {
		POST("/{page:.*}", PostPageHandler)
	}

	if len(*build) > 0 {
		if err := buildStaticSite(*build); err != nil {
			log.Printf("%s", err.Error())
		}
		os.Exit(0)
	}

	// Start the server
	START()
}

// Redirect to `/index` to render the index page.
func RootHandler(w Response, r Request) Output {
	return Redirect("/index")
}

// Shows a page. the page name is the path itself. if the page doesn't exist it
// redirect to edit page otherwise will render it to HTML
func GetPageHandler(w Response, r Request) Output {
	vars := VARS(r)
	page := NewPage(vars["page"])

	if !page.Exists() {
		if READONLY {
			return NotFound
		} else {
			return Redirect("/edit/" + page.Name)
		}
	}

	return Render("view", Locals{
		"edit":      "/edit/" + page.Name,
		"title":     page.Emoji() + " " + page.Name,
		"updated":   ago(time.Now().Sub(page.ModTime())),
		"content":   template.HTML(page.Render()),
		"tools":     renderWidget(TOOLS_WIDGET, &page, r),      // all tools registered widgets
		"sidebar":   renderWidget(SIDEBAR_WIDGET, &page, r),    // widgets registered for sidebar
		"meta":      renderWidget(META_WIDGET, &page, r),       // widgets registered to be displayed under the page title
		"afterView": renderWidget(AFTER_VIEW_WIDGET, &page, r), // widgets registered to be displayed under the page content in the view page
	})
}

// Edit page, gets the page from path, if it doesn't exist it'll use the
// template.md content as default value
func GetPageEditHandler(w Response, r Request) Output {
	vars := VARS(r)
	page := NewPage(vars["page"])

	var content string
	if page.Exists() {
		content = page.Content()
	} else if template := NewPage(TEMPLATE_NAME); template.Exists() {
		content = template.Content()
	}

	// Execute all Autocomplete functions and add them to a slice and pass it
	// down to the view
	acs := []*Autocomplete{}
	for _, v := range autocompletes {
		acs = append(acs, v())
	}

	return Render("edit", Locals{
		"title":        page.Name,
		"action":       page.Name,
		"rtl":          page.RTL(),                           // is it Right-To-Left page?
		"tools":        renderWidget(TOOLS_WIDGET, &page, r), // render all tools widgets
		"content":      content,
		"autocomplete": acs,
		"csrf":         CSRF(r),
	})
}

// Save new content of the page
func PostPageHandler(w Response, r Request) Output {
	vars := VARS(r)
	page := NewPage(vars["page"])
	content := r.FormValue("content")
	page.Write(content)

	return Redirect("/" + page.Name)
}

// WIDGETS ===================================================

type (
	// a type used to define list of widgets spaces. it's used to register
	// widgets to be injected in the view or edit pages
	widgetSpace int
	// a function that takes the current page and the HTTP request and returns
	// the widget. This can be used by extensions to define new widgets to be
	// rendered in view or edit pages. the extension should define this func
	// type and register it to be rendered in a specific widgetSpace such as
	// before or after the page.
	widgetFunc func(*Page, Request) template.HTML
)

// List of widgets spaces that extensions can use to register a widgetFunc to
// inject content into.
const (
	TOOLS_WIDGET widgetSpace = iota
	SIDEBAR_WIDGET
	AFTER_VIEW_WIDGET
	META_WIDGET
)

// A map to keep track of list of widget functions registered in each widget space
var widgets = map[widgetSpace][]widgetFunc{}

// Register widget function to be rendered in a specific space before any other
// widget. functions registered by this function will have higher priority than
// the rest. this function is needed for example to register the search input
// before any other links in the sidebar
func PREPEND_WIDGET(s widgetSpace, f func(*Page, Request) template.HTML) {
	if _, ok := widgets[s]; !ok {
		widgets[s] = []widgetFunc{}
	}
	widgets[s] = append([]widgetFunc{f}, widgets[s]...)
}

// Register a function to a widget space. functions registered will be executed
// in order when rendering view or edit page. the return values of these
// widgetfuncs will pass down to the template and injected in reserved places.
func WIDGET(s widgetSpace, f func(*Page, Request) template.HTML) {
	if _, ok := widgets[s]; !ok {
		widgets[s] = []widgetFunc{}
	}
	widgets[s] = append(widgets[s], f)
}

// This is used by view and edit routes to render all widgetfuncs registered for
// specific widget space.
func renderWidget(s widgetSpace, p *Page, r Request) (o template.HTML) {
	for _, v := range widgets[s] {
		o += v(p, r)
	}
	return
}

// HELPERS

// A function that takes time.duration and return a string representation of the
// duration in human readable way such as "3 seconds ago". "5 hours 30 minutes
// ago". The precision of this function is 2. which means it returns the largest
// unit of time possible and the next one after it. for example days + hours, or
// Hours + minutes or Minutes + seconds...etc
func ago(t time.Duration) string {
	const day = time.Hour * 24
	const week = day * 7
	const month = day * 30
	const year = day * 365
	const maxPrecision = 2

	var o strings.Builder

	if t.Seconds() < 1 {
		o.WriteString("Less than a second ")
	}

	for precision := 0; t.Seconds() > 1 && precision < maxPrecision; precision++ {
		switch {
		case t >= year:
			years := t / year
			t -= years * year
			o.WriteString(fmt.Sprintf("%d years ", years))
		case t >= month:
			months := t / month
			t -= months * month
			o.WriteString(fmt.Sprintf("%d months ", months))
		case t >= week:
			weeks := t / week
			t -= weeks * week
			o.WriteString(fmt.Sprintf("%d weeks ", weeks))
		case t >= day:
			days := t / day
			t -= days * day
			o.WriteString(fmt.Sprintf("%d days ", days))
		case t >= time.Hour:
			hours := t / time.Hour
			t -= hours * time.Hour
			o.WriteString(fmt.Sprintf("%d hours ", hours))
		case t >= time.Minute:
			minutes := t / time.Minute
			t -= minutes * time.Minute
			o.WriteString(fmt.Sprintf("%d minutes ", minutes))
		case t >= time.Second:
			seconds := t / time.Second
			t -= seconds * time.Second
			o.WriteString(fmt.Sprintf("%d seconds ", seconds))
		}
	}

	o.WriteString("ago")

	return o.String()
}

// AUTOCOMPLETE ================================================

// Autocomplete defines what character triggeres the autocomplete feature and
// what is the list to display in this case.
type Autocomplete struct {
	StartChar   string
	Suggestions []*Suggestion
}

// Suggestions represent an item in the list of autocomplete menu in the edit page
type Suggestion struct {
	Text        string // The text that gets injected in the editor if this option is choosen
	DisplayText string // The display text for this item in the menu. this can be more cosmetic.
}

// This is a function that returns an auto completer instance. this function
// should be defined by extensions and registered to be executed when rendering
// the edit page
type Autocompleter func() *Autocomplete

// Holds a list of registered autocompleter functions
var autocompletes = []Autocompleter{}

// this function registers an autocompleter function. it should be used by an
// extension to register a new autocompleter function. these functions are going
// to be executed when rendering the edit page.
func AUTOCOMPLETE(a Autocompleter) {
	autocompletes = append(autocompletes, a)
}
