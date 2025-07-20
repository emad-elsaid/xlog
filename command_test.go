package xlog

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandsBehavior(t *testing.T) {
	app := newTestApp()
	app.commands = []func(Page) []Command{}

	page := &page{name: "test"}
	commands := app.Commands(page)
	require.Len(t, commands, 0, "Expected 0 commands")

	testCommand := func(p Page) []Command {
		return []Command{
			&testCommandImpl{
				icon:  "fa-test",
				name:  "Test Command",
				attrs: map[template.HTMLAttr]any{"href": "/test"},
			},
		}
	}

	app.RegisterCommand(testCommand)
	commands = app.Commands(page)
	require.Len(t, commands, 1, "Expected 1 command")
	require.Equal(t, "Test Command", commands[0].Name(), "Expected command name 'Test Command'")
}

func TestQuickCommandsBehavior(t *testing.T) {
	app := newTestApp()
	app.quickCommands = []func(Page) []Command{}

	page := &page{name: "test"}
	commands := app.QuickCommands(page)
	require.Len(t, commands, 0, "Expected 0 quick commands")

	testCommand := func(p Page) []Command {
		return []Command{
			&testCommandImpl{
				icon:  "fa-quick",
				name:  "Quick Test",
				attrs: map[template.HTMLAttr]any{"href": "/quick"},
			},
		}
	}

	app.RegisterQuickCommand(testCommand)
	commands = app.QuickCommands(page)
	require.Len(t, commands, 1, "Expected 1 quick command")
	require.Equal(t, "Quick Test", commands[0].Name(), "Expected quick command name 'Quick Test'")
}

func TestLinksBehavior(t *testing.T) {
	app := newTestApp()
	app.links = []func(Page) []Command{}

	page := &page{name: "test"}
	links := app.Links(page)
	require.Len(t, links, 0, "Expected 0 links")

	testLink := func(p Page) []Command {
		return []Command{
			&testCommandImpl{
				icon:  "fa-link",
				name:  "Test Link",
				attrs: map[template.HTMLAttr]any{"href": "/link"},
			},
		}
	}

	app.RegisterLink(testLink)
	links = app.Links(page)
	require.Len(t, links, 1, "Expected 1 link")
	require.Equal(t, "Test Link", links[0].Name(), "Expected link name 'Test Link'")
}

// Test command implementation for testing
type testCommandImpl struct {
	icon  string
	name  string
	attrs map[template.HTMLAttr]any
}

func (c *testCommandImpl) Icon() string {
	return c.icon
}

func (c *testCommandImpl) Name() string {
	return c.name
}

func (c *testCommandImpl) Attrs() map[template.HTMLAttr]any {
	return c.attrs
}
