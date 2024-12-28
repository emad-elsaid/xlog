package custom_widget

import (
	"flag"
	"html/template"
	"os"

	"github.com/emad-elsaid/memoize"
	. "github.com/emad-elsaid/xlog"
)

var head_file, before_view_file, after_view_file string

func init() {
	flag.StringVar(&head_file, "custom.head", "", "path to a file it's content will be included in every page <head> tag")
	flag.StringVar(&before_view_file, "custom.before_view", "", "path to a file it's content will be included in every page BEFORE the content of the page")
	flag.StringVar(&after_view_file, "custom.after_view", "", "path to a file it's content will be included in every page AFTER the content of the page")

	RegisterExtension(CustomWidget{})
}

type CustomWidget struct{}

func (CustomWidget) Name() string { return "custom-widget" }
func (CustomWidget) Init() {
	if head_file != "" {
		RegisterWidget(WidgetHead, 1, func(Page) template.HTML {
			return readFile(head_file)
		})
	}
	if before_view_file != "" {
		RegisterWidget(WidgetBeforeView, 1, func(Page) template.HTML {
			return readFile(before_view_file)
		})
	}
	if after_view_file != "" {
		RegisterWidget(WidgetAfterView, 1, func(Page) template.HTML {
			return readFile(after_view_file)
		})
	}
}

var readFile = memoize.New(func(f string) template.HTML {
	b, err := os.ReadFile(f)
	if err != nil {
		panic(err)
	}

	return template.HTML(b)
})
