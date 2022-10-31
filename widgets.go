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
	widgetFunc func(Page, Request) template.HTML
)

// List of widgets spaces that extensions can use to register a widgetFunc to
// inject content into.
const (
	SIDEBAR_WIDGET    widgetSpace = iota // widgets rengered in the sidebar
	AFTER_VIEW_WIDGET                    // widgets rendered after the content of the view page
	ACTION_WIDGET                        // widgets rendered in the actions row of the view page
	HEAD_WIDGET                          // widgets rendered in page <head> tag
)

// A map to keep track of list of widget functions registered in each widget space
var widgets = map[widgetSpace]*plist[widgetFunc]{}

// RegisterWidget Register a function to a widget space. functions registered
// will be executed in order of priority lower to higher when rendering view or
// edit page. the return values of these widgetfuncs will pass down to the
// template and injected in reserved places.
func RegisterWidget(s widgetSpace, priority float32, f widgetFunc) {
	widgets[s] = widgets[s].insert(priority, f)
}

// This is used by view and edit routes to render all widgetfuncs registered for
// specific widget space.
func RenderWidget(s widgetSpace, p Page, r Request) (o template.HTML) {
	widgets[s].each(func(f widgetFunc) {
		o += f(p, r)
	})
	return
}
