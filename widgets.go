package xlog

import "html/template"

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
	ACTION_WIDGET
	HEAD_WIDGET
)

// A map to keep track of list of widget functions registered in each widget space
var widgets = map[widgetSpace][]widgetFunc{}

// Register widget function to be rendered in a specific space before any other
// widget. functions registered by this function will have higher priority than
// the rest. this function is needed for example to register the search input
// before any other links in the sidebar
func PrependWidget(s widgetSpace, f func(*Page, Request) template.HTML) {
	if _, ok := widgets[s]; !ok {
		widgets[s] = []widgetFunc{}
	}
	widgets[s] = append([]widgetFunc{f}, widgets[s]...)
}

// Register a function to a widget space. functions registered will be executed
// in order when rendering view or edit page. the return values of these
// widgetfuncs will pass down to the template and injected in reserved places.
func Widget(s widgetSpace, f func(*Page, Request) template.HTML) {
	if _, ok := widgets[s]; !ok {
		widgets[s] = []widgetFunc{}
	}
	widgets[s] = append(widgets[s], f)
}

// This is used by view and edit routes to render all widgetfuncs registered for
// specific widget space.
func RenderWidget(s widgetSpace, p *Page, r Request) (o template.HTML) {
	for _, v := range widgets[s] {
		o += v(p, r)
	}
	return
}
