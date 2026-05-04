package custom_widget

import (
	"flag"
	"html/template"
	"os"

	"github.com/emad-elsaid/memoize"
	"github.com/emad-elsaid/xlog"
)

var head_file, before_view_file, after_view_file string

func init() {
	flag.StringVar(&head_file, "custom.head", "", "path to a file it's content will be included in every page <head> tag")
	flag.StringVar(&before_view_file, "custom.before_view", "", "path to a file it's content will be included in every page BEFORE the content of the page")
	flag.StringVar(&after_view_file, "custom.after_view", "", "path to a file it's content will be included in every page AFTER the content of the page")

	xlog.RegisterExtension(CustomWidget{})
}

type CustomWidget struct{}

func (CustomWidget) Name() string { return "custom-widget" }
func (CustomWidget) Init() {
	if head_file != "" {
		xlog.RegisterWidget(xlog.WidgetHead, 1, func(xlog.Page) template.HTML {
			return readFile(head_file)
		})
	}
	if before_view_file != "" {
		xlog.RegisterWidget(xlog.WidgetBeforeView, 1, func(xlog.Page) template.HTML {
			return readFile(before_view_file)
		})
	}
	if after_view_file != "" {
		xlog.RegisterWidget(xlog.WidgetAfterView, 1, func(xlog.Page) template.HTML {
			return readFile(after_view_file)
		})
	}
}

var readFile = memoize.New(func(f string) template.HTML {
	b, err := os.ReadFile(f) // #nosec G304 -- File path controlled by admin via command-line flags, not user input
	if err != nil {
		// Return error message as HTML comment instead of crashing server
		// #nosec G203 -- Error message is HTMLEscapeString-sanitized before template.HTML conversion
		return template.HTML("<!-- custom widget error: " + template.HTMLEscapeString(err.Error()) + " -->")
	}

	return template.HTML(b) // #nosec G203 -- Content is trusted admin HTML from admin-controlled file
})
