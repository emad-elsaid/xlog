package custom_widget

import (
	"flag"
	"html/template"
	"os"

	. "github.com/emad-elsaid/xlog"
)

var head_file, before_view_file, after_view_file string

func init() {
	flag.StringVar(&head_file, "custom_head", "", "path to a file it's content will be included in every page <head> tag")
	flag.StringVar(&before_view_file, "custom_before_view", "", "path to a file it's content will be included in every page BEFORE the content of the page")
	flag.StringVar(&after_view_file, "custom_after_view", "", "path to a file it's content will be included in every page AFTER the content of the page")

	RegisterExtension(CustomWidget{})
}

type CustomWidget struct{}

func (CustomWidget) Name() string { return "custom-widget" }
func (CustomWidget) Init() {
	RegisterWidget(HEAD_WIDGET, 1, custom_head)
	RegisterWidget(BEFORE_VIEW_WIDGET, 1, custom_before_view)
	RegisterWidget(AFTER_VIEW_WIDGET, 1, custom_after_view)
}

var head_content []byte

func custom_head(_ Page) template.HTML {
	if head_file == "" {
		return ""
	}

	if head_content == nil {
		var err error
		head_content, err = os.ReadFile(head_file)
		if err != nil {
			panic(err)
		}
	}

	return template.HTML(head_content)
}

var before_view_content []byte

func custom_before_view(_ Page) template.HTML {
	if before_view_file == "" {
		return ""
	}

	if before_view_content == nil {
		var err error
		before_view_content, err = os.ReadFile(before_view_file)
		if err != nil {
			panic(err)
		}
	}

	return template.HTML(before_view_content)
}

var after_view_content []byte

func custom_after_view(_ Page) template.HTML {
	if after_view_file == "" {
		return ""
	}

	if after_view_content == nil {
		var err error
		after_view_content, err = os.ReadFile(after_view_file)
		if err != nil {
			panic(err)
		}
	}

	return template.HTML(after_view_content)
}
