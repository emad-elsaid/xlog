package xlog

import (
	"html/template"
	"sort"
)

type (
	// a type used to represent a widgets spaces. it's used to register
	// widgets to be injected in the view or edit pages
	WidgetSpace string
	// a function that takes the current page and returns the widget. This can
	// be used by extensions to define new widgets to be rendered in view or
	// edit pages. the extension should define this func type and register it to
	// be rendered in a specific widgetSpace such as before or after the page.
	WidgetFunc func(Page) template.HTML
)

// List of widgets spaces that extensions can use to register a WidgetFunc to
// inject content into.
const (
	SIDEBAR_WIDGET    = "sidebar"    // widgets rendered in the sidebar
	AFTER_VIEW_WIDGET = "after_view" // widgets rendered after the content of the view page
	HEAD_WIDGET       = "head"       // widgets rendered in page <head> tag
)

// A map to keep track of list of widget functions registered in each widget space
var widgets = map[WidgetSpace]byPriority[WidgetFunc]{}

// RegisterWidget Register a function to a widget space. functions registered
// will be executed in order of priority lower to higher when rendering view or
// edit page. the return values of these widgetfuncs will pass down to the
// template and injected in reserved places.
func RegisterWidget(s WidgetSpace, priority float32, f WidgetFunc) {
	widgets[s] = append(widgets[s], priorityItem[WidgetFunc]{
		priority: priority,
		value:    f,
	})
	sort.Sort(widgets[s])
}

// This is used by view and edit routes to render all widgetfuncs registered for
// specific widget space.
func RenderWidget(s WidgetSpace, p Page) (o template.HTML) {
	for _, f := range widgets[s] {
		o += f.value(p)
	}
	return
}
