package xlog

import "html/template"

type Command interface {
	Name() string
	OnClick() template.JS
	Widget(*Page) template.HTML
}

var commands = []Command{}

// RegisterCommand registers a command
func RegisterCommand(c Command) {
	commands = append(commands, c)
}
