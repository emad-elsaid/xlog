package xlog

import "html/template"

// Command defines a structure used for 3 categories of lists:
// 1. Commands for Ctrl+K actions menu
// 2. Quick commands displayed in the default template at the top right of the page
// 3. Links displayed in the navigation bar
// The template decides where and how to display commands. it can choose to use them in a different way than the default template
type Command interface {
	// Icon returns the Fontawesome icon class name for the Command
	Icon() string
	// Name of the command. to be displayed in the list
	Name() string
	// Attrs a map of attributes to their values
	Attrs() map[template.HTMLAttr]any
}

// RegisterCommand registers a new command
func RegisterCommand(c func(Page) []Command) {
	app := GetApp()
	app.RegisterCommand(c)
}

// Commands return the list of commands for a page. when a page is displayed it
// executes all functions registered with RegisterCommand and collect all
// results in one slice. result can be passed to the view to render the commands
// list
func Commands(p Page) []Command {
	app := GetApp()
	return app.Commands(p)
}

func RegisterQuickCommand(c func(Page) []Command) {
	app := GetApp()
	app.RegisterQuickCommand(c)
}

// QuickCommands return the list of QuickCommands for a page. it executes all functions
// registered with RegisterQuickCommand and collect all results in one slice. result
// can be passed to the view to render the Quick commands list
func QuickCommands(p Page) []Command {
	app := GetApp()
	return app.QuickCommands(p)
}

// Register a new links function, should return a list of Links
func RegisterLink(l func(Page) []Command) {
	app := GetApp()
	app.RegisterLink(l)
}

// Links returns a list of links for a Page. it executes all functions
// registered with RegisterLink and collect them in one slice. Can be passed to
// the view to render in the footer for example.
func Links(p Page) []Command {
	app := GetApp()
	return app.Links(p)
}
