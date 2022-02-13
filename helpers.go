package main

import "html/template"

func Helpers() {
	helpers["partial"] = func(v string, data interface{}) template.HTML {
		return template.HTML(partial(v, data))
	}
}
