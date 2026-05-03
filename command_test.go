package xlog

import (
	"html/template"
	"testing"
)

// mockCommand implements the Command interface for testing.
type mockCommand struct {
	icon  string
	name  string
	attrs map[template.HTMLAttr]any
}

func (m mockCommand) Icon() string {
	return m.icon
}

func (m mockCommand) Name() string {
	return m.name
}

func (m mockCommand) Attrs() map[template.HTMLAttr]any {
	return m.attrs
}

func TestCommands_Empty(t *testing.T) {
	// Save original state
	originalCommands := commands
	defer func() { commands = originalCommands }()

	// Reset to empty
	commands = []func(Page) []Command{}

	page := &mockPage{name: "test", exists: true}
	result := Commands(page)

	if len(result) != 0 {
		t.Errorf("Expected 0 commands, got %d", len(result))
	}
}

func TestRegisterCommand_Single(t *testing.T) {
	// Save original state
	originalCommands := commands
	defer func() { commands = originalCommands }()

	// Reset to empty
	commands = []func(Page) []Command{}

	// Register a command
	RegisterCommand(func(p Page) []Command {
		return []Command{
			mockCommand{
				icon: "fa-test",
				name: "Test Command",
				attrs: map[template.HTMLAttr]any{
					"href": "/test",
				},
			},
		}
	})

	page := &mockPage{name: "test", exists: true}
	result := Commands(page)

	if len(result) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(result))
	}

	cmd := result[0]
	if cmd.Icon() != "fa-test" {
		t.Errorf("Expected icon 'fa-test', got '%s'", cmd.Icon())
	}

	if cmd.Name() != "Test Command" {
		t.Errorf("Expected name 'Test Command', got '%s'", cmd.Name())
	}
}

func TestRegisterCommand_Multiple(t *testing.T) {
	// Save original state
	originalCommands := commands
	defer func() { commands = originalCommands }()

	// Reset to empty
	commands = []func(Page) []Command{}

	// Register multiple commands
	RegisterCommand(func(p Page) []Command {
		return []Command{
			mockCommand{icon: "fa-one", name: "Command One"},
		}
	})

	RegisterCommand(func(p Page) []Command {
		return []Command{
			mockCommand{icon: "fa-two", name: "Command Two"},
			mockCommand{icon: "fa-three", name: "Command Three"},
		}
	})

	page := &mockPage{name: "test", exists: true}
	result := Commands(page)

	if len(result) != 3 {
		t.Fatalf("Expected 3 commands, got %d", len(result))
	}

	// Verify order is preserved
	if result[0].Name() != "Command One" {
		t.Errorf("Expected first command 'Command One', got '%s'", result[0].Name())
	}

	if result[1].Name() != "Command Two" {
		t.Errorf("Expected second command 'Command Two', got '%s'", result[1].Name())
	}

	if result[2].Name() != "Command Three" {
		t.Errorf("Expected third command 'Command Three', got '%s'", result[2].Name())
	}
}

func TestCommands_PageParameter(t *testing.T) {
	// Save original state
	originalCommands := commands
	defer func() { commands = originalCommands }()

	// Reset to empty
	commands = []func(Page) []Command{}

	// Register command that uses page name
	RegisterCommand(func(p Page) []Command {
		return []Command{
			mockCommand{
				name: "Edit " + p.Name(),
			},
		}
	})

	page := &mockPage{name: "my-page", exists: true}
	result := Commands(page)

	if len(result) != 1 {
		t.Fatalf("Expected 1 command, got %d", len(result))
	}

	if result[0].Name() != "Edit my-page" {
		t.Errorf("Expected 'Edit my-page', got '%s'", result[0].Name())
	}
}

func TestQuickCommands_Empty(t *testing.T) {
	// Save original state
	originalQuickCommands := quickCommands
	defer func() { quickCommands = originalQuickCommands }()

	// Reset to empty
	quickCommands = []func(Page) []Command{}

	page := &mockPage{name: "test", exists: true}
	result := QuickCommands(page)

	if len(result) != 0 {
		t.Errorf("Expected 0 quick commands, got %d", len(result))
	}
}

func TestRegisterQuickCommand_Single(t *testing.T) {
	// Save original state
	originalQuickCommands := quickCommands
	defer func() { quickCommands = originalQuickCommands }()

	// Reset to empty
	quickCommands = []func(Page) []Command{}

	RegisterQuickCommand(func(p Page) []Command {
		return []Command{
			mockCommand{
				icon: "fa-quick",
				name: "Quick Action",
			},
		}
	})

	page := &mockPage{name: "test", exists: true}
	result := QuickCommands(page)

	if len(result) != 1 {
		t.Fatalf("Expected 1 quick command, got %d", len(result))
	}

	if result[0].Icon() != "fa-quick" {
		t.Errorf("Expected icon 'fa-quick', got '%s'", result[0].Icon())
	}
}

func TestQuickCommands_Multiple(t *testing.T) {
	// Save original state
	originalQuickCommands := quickCommands
	defer func() { quickCommands = originalQuickCommands }()

	// Reset to empty
	quickCommands = []func(Page) []Command{}

	RegisterQuickCommand(func(p Page) []Command {
		return []Command{
			mockCommand{name: "Quick One"},
			mockCommand{name: "Quick Two"},
		}
	})

	RegisterQuickCommand(func(p Page) []Command {
		return []Command{
			mockCommand{name: "Quick Three"},
		}
	})

	page := &mockPage{name: "test", exists: true}
	result := QuickCommands(page)

	if len(result) != 3 {
		t.Fatalf("Expected 3 quick commands, got %d", len(result))
	}
}

func TestLinks_Empty(t *testing.T) {
	// Save original state
	originalLinks := links
	defer func() { links = originalLinks }()

	// Reset to empty
	links = []func(Page) []Command{}

	page := &mockPage{name: "test", exists: true}
	result := Links(page)

	if len(result) != 0 {
		t.Errorf("Expected 0 links, got %d", len(result))
	}
}

func TestRegisterLink_Single(t *testing.T) {
	// Save original state
	originalLinks := links
	defer func() { links = originalLinks }()

	// Reset to empty
	links = []func(Page) []Command{}

	RegisterLink(func(p Page) []Command {
		return []Command{
			mockCommand{
				icon: "fa-link",
				name: "External Link",
				attrs: map[template.HTMLAttr]any{
					"href":   "https://example.com",
					"target": "_blank",
				},
			},
		}
	})

	page := &mockPage{name: "test", exists: true}
	result := Links(page)

	if len(result) != 1 {
		t.Fatalf("Expected 1 link, got %d", len(result))
	}

	if result[0].Name() != "External Link" {
		t.Errorf("Expected 'External Link', got '%s'", result[0].Name())
	}

	attrs := result[0].Attrs()
	if attrs["href"] != "https://example.com" {
		t.Errorf("Expected href 'https://example.com', got '%v'", attrs["href"])
	}
}

func TestLinks_Multiple(t *testing.T) {
	// Save original state
	originalLinks := links
	defer func() { links = originalLinks }()

	// Reset to empty
	links = []func(Page) []Command{}

	RegisterLink(func(p Page) []Command {
		return []Command{
			mockCommand{name: "Link One"},
		}
	})

	RegisterLink(func(p Page) []Command {
		return []Command{
			mockCommand{name: "Link Two"},
			mockCommand{name: "Link Three"},
		}
	})

	page := &mockPage{name: "test", exists: true}
	result := Links(page)

	if len(result) != 3 {
		t.Fatalf("Expected 3 links, got %d", len(result))
	}

	// Verify accumulation order
	if result[0].Name() != "Link One" {
		t.Errorf("Expected first link 'Link One', got '%s'", result[0].Name())
	}
}

func TestCommand_Attrs(t *testing.T) {
	// Test that attrs are properly stored and retrieved
	cmd := mockCommand{
		icon: "fa-test",
		name: "Test",
		attrs: map[template.HTMLAttr]any{
			"href":         "/path",
			"data-action":  "click",
			"data-confirm": "Are you sure?",
		},
	}

	attrs := cmd.Attrs()

	if len(attrs) != 3 {
		t.Errorf("Expected 3 attributes, got %d", len(attrs))
	}

	if attrs["href"] != "/path" {
		t.Errorf("Expected href '/path', got '%v'", attrs["href"])
	}

	if attrs["data-action"] != "click" {
		t.Errorf("Expected data-action 'click', got '%v'", attrs["data-action"])
	}
}
