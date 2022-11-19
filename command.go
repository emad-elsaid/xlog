package xlog

import "html/template"

// Command define a command that a user can invoke in view or edit page on a
// Page.
type Command interface {
	// Icon returns the Fontawesome icon class name for the Command
	Icon() string
	// Name of the command. to be displayed in the list
	Name() string
	// Link returns the link/url/path of the command if any
	Link() string
	// OnClick action. a Javascript code to invoke when the command is executed
	OnClick() template.JS
	// Widget a HTML snippet to embed in the page that include any needed
	// assets, HTML, JS the command needs
	Widget() template.HTML
}

var commands = []func(Page) []Command{defaultCommands}

// RegisterCommand registers a new command
func RegisterCommand(c func(Page) []Command) {
	commands = append(commands, c)
}

// Commands return the list of commands for a page. it executes all functions
// registered with RegisterCommand and collect all results in one slice. result
// can be passed to the view to render the commands list
func Commands(p Page) []Command {
	cmds := []Command{}
	for c := range commands {
		cmds = append(cmds, commands[c](p)...)
	}

	return cmds
}

var quickCommands = []func(Page) []Command{defaultCommands}

func RegisterQuickCommand(c func(Page) []Command) {
	quickCommands = append(quickCommands, c)
}

// QuickCommands return the list of QuickCommands for a page. it executes all functions
// registered with RegisterQuickCommand and collect all results in one slice. result
// can be passed to the view to render the Quick commands list
func QuickCommands(p Page) []Command {
	cmds := []Command{}
	for c := range quickCommands {
		cmds = append(cmds, quickCommands[c](p)...)
	}

	return cmds
}

type editQuickCommand struct{ page Page }

func (a editQuickCommand) Icon() string          { return "fa-solid fa-pen" }
func (a editQuickCommand) Name() string          { return "Edit" }
func (a editQuickCommand) Link() string          { return "/edit/" + a.page.Name() }
func (a editQuickCommand) OnClick() template.JS  { return "" }
func (a editQuickCommand) Widget() template.HTML { return "" }

func defaultCommands(p Page) []Command {
	if READONLY {
		return nil
	}

	return []Command{editQuickCommand{p}}
}
