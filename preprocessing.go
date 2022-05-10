package main

import "regexp"

var imgUrlReg = regexp.MustCompile(`(?imU)^(https\:\/\/[^ ]+(svg|jpg|jpeg|gif|png|webp))$`)

func preProcess(content string) string {
	content = imgUrlReg.ReplaceAllString(content, `<img src="$1"/>`)
	return content
}
