package xlog

import "html/template"

// Command defines a structure used for 3 categories of lists:
// 1. Commands for Ctrl+K actions menu
// 2. Quick commands displayed in the default template at the top right of the page
// 3. Links displayed in the navigation bar
// The template decides where and how to display commands. it can choose to use them in a different way than the default template.
type Command interface {
	// Icon returns the Fontawesome icon class name for the Command.
	Icon() string
	// Name of the command. to be displayed in the list.
	Name() string
	// Attrs a map of attributes to their values.
	Attrs() map[template.HTMLAttr]any
}

var commands = []func(Page) []Command{}

// RegisterCommand registers a new command function. When a page is displayed,
// all registered command functions are executed to collect commands for the
// Ctrl+K actions menu.
func RegisterCommand(c func(Page) []Command) {
	commands = append(commands, c)
}

// Commands returns the list of commands for a page. When a page is displayed,
// it executes all functions registered with RegisterCommand and collects all
// results in one slice. The result can be passed to the view to render the
// commands list for the Ctrl+K actions menu.
func Commands(p Page) []Command {
	cmds := []Command{}
	for c := range commands {
		cmds = append(cmds, commands[c](p)...)
	}

	return cmds
}

var quickCommands = []func(Page) []Command{}

// RegisterQuickCommand registers a new quick command function. Quick commands are
// displayed prominently at the top right of the page in the default template.
func RegisterQuickCommand(c func(Page) []Command) {
	quickCommands = append(quickCommands, c)
}

// QuickCommands returns the list of QuickCommands for a page. It executes all
// functions registered with RegisterQuickCommand and collects all results in
// one slice. The result can be passed to the view to render the quick commands
// list, which are displayed prominently at the top right of the page in the
// default template.
func QuickCommands(p Page) []Command {
	cmds := []Command{}
	for c := range quickCommands {
		cmds = append(cmds, quickCommands[c](p)...)
	}

	return cmds
}

var links = []func(Page) []Command{}

// RegisterLink registers a new links function. The function should return a
// list of links that will be displayed in navigation areas such as the footer.
func RegisterLink(l func(Page) []Command) {
	links = append(links, l)
}

// Links returns a list of links for a Page. It executes all functions
// registered with RegisterLink and collects them in one slice. The result
// can be passed to the view to render in navigation areas such as the footer.
func Links(p Page) []Command {
	lnks := []Command{}
	for l := range links {
		lnks = append(lnks, links[l](p)...)
	}
	return lnks
}
