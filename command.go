package xlog

import "html/template"

// Command define a command that a user can invoke in view or edit page on a
// Page.
type Command interface {
	// Name of the command. to be displayed in the list
	Name() string
	// OnClick action. a Javascript code to invoke when the command is executed
	OnClick() template.JS
	// Widget a HTML snippet to embed in the page that include any needed
	// assets, HTML, JS the command needs
	Widget(Page) template.HTML
}

var commands = []Command{}

// RegisterCommand registers a new command
func RegisterCommand(c Command) {
	commands = append(commands, c)
}
