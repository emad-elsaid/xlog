package xlog

import "html/template"

// Command defines a structure used for actions and links
type Command interface {
	// Icon returns the Fontawesome icon class name for the Command
	Icon() string
	// Name of the command. to be displayed in the list
	Name() string
	// Attrs a map of attributes to their values
	Attrs() map[template.HTMLAttr]any
}

// RegisterCommand registers a new command
func (app *App) RegisterCommand(c func(Page) []Command) {
	app.commands = append(app.commands, c)
}

// Commands returns the list of commands for a page
func (app *App) Commands(p Page) []Command {
	cmds := []Command{}
	for _, c := range app.commands {
		cmds = append(cmds, c(p)...)
	}
	return cmds
}

// RegisterQuickCommand registers a new quick command
func (app *App) RegisterQuickCommand(c func(Page) []Command) {
	app.quickCommands = append(app.quickCommands, c)
}

// RegisterLink registers a new link
func (app *App) RegisterLink(l func(Page) []Command) {
	app.links = append(app.links, l)
}

// QuickCommands returns the list of quick commands for a page
func (app *App) QuickCommands(p Page) []Command {
	cmds := []Command{}
	for _, c := range app.quickCommands {
		cmds = append(cmds, c(p)...)
	}
	return cmds
}

// Links returns a list of links for a Page
func (app *App) Links(p Page) []Command {
	cmds := []Command{}
	for _, l := range app.links {
		cmds = append(cmds, l(p)...)
	}
	return cmds
}
