package xlog

import (
	"html/template"
	"iter"
	"sort"
)

type (
	// WidgetSpace used to represent a widgets spaces. it's used to register
	// widgets to be injected in the view or edit pages
	WidgetSpace string
	// WidgetFunc a function that takes the current page and returns the widget.
	// This can be used by extensions to define new widgets to be rendered in
	// view or edit pages. the extension should define this func type and
	// register it to be rendered in a specific widgetSpace such as before or
	// after the page.
	WidgetFunc func(Page) template.HTML
)

// List of widgets spaces that extensions can use to register a WidgetFunc to
// inject content into.
var (
	WidgetAfterView  WidgetSpace = "after_view"  // widgets rendered after the content of the view page
	WidgetBeforeView WidgetSpace = "before_view" // widgets rendered before the content of the view page
	WidgetHead       WidgetSpace = "head"        // widgets rendered in page <head> tag
)

// A map to keep track of list of widget functions registered in each widget space
var widgets = map[WidgetSpace]*priorityList[WidgetFunc]{}

// RegisterWidget Register a function to a widget space. functions registered
// will be executed in order of priority lower to higher when rendering view or
// edit page. the return values of these widgetfuncs will pass down to the
// template and injected in reserved places.
func RegisterWidget(s WidgetSpace, priority float32, f WidgetFunc) {
	pl, ok := widgets[s]
	if !ok {
		pl = new(priorityList[WidgetFunc])
		widgets[s] = pl
	}

	pl.Add(f, priority)
}

// This is used by view and edit routes to render all widgetfuncs registered for
// specific widget space.
func RenderWidget(s WidgetSpace, p Page) (o template.HTML) {
	w, ok := widgets[s]
	if !ok {
		return
	}

	for f := range w.All() {
		o += f(p)
	}
	return
}

type priorityItem[T any] struct {
	Item     T
	Priority float32
}

type priorityList[T any] struct {
	items []priorityItem[T]
}

func (pl *priorityList[T]) Add(item T, priority float32) {
	pl.items = append(pl.items, priorityItem[T]{Item: item, Priority: priority})
	pl.sortByPriority()
}

func (pl *priorityList[T]) sortByPriority() {
	sort.Slice(pl.items, func(i, j int) bool {
		return pl.items[i].Priority < pl.items[j].Priority
	})
}

// An iterator over all items
func (pl *priorityList[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range pl.items {
			if !yield(v.Item) {
				return
			}
		}
	}
}
