package todo

import (
	"testing"
)

func TestTaskListRegexp(t *testing.T) {
	tests := []struct {
		input       string
		shouldMatch bool
		replaceWith string
		expected    string
	}{
		{"[ ] Task", true, "[x] ", "[x] Task"},
		{"[x] Task", true, "[ ] ", "[ ] Task"},
		{"[X] Task", true, "[ ] ", "[ ] Task"},
		{"[ ]Task", true, "[x] ", "[x] Task"},
		{"[  ] Task", false, "", "[  ] Task"},
		{"[] Task", false, "", "[] Task"},
		{"Task without checkbox", false, "", "Task without checkbox"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			matches := taskListRegexp.MatchString(tt.input)
			if matches != tt.shouldMatch {
				t.Errorf("Expected match=%v for input %q, got %v", tt.shouldMatch, tt.input, matches)
			}

			if tt.shouldMatch {
				result := taskListRegexp.ReplaceAllString(tt.input, tt.replaceWith)
				if result != tt.expected {
					t.Errorf("Expected %q, got %q", tt.expected, result)
				}
			}
		})
	}
}

func TestTODOExtensionName(t *testing.T) {
	ext := TODO{}

	if ext.Name() != "todo" {
		t.Errorf("Expected name 'todo', got '%s'", ext.Name())
	}
}

func TestToggleLogic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		checked  bool
		expected string
	}{
		{
			name:     "Uncheck box with space",
			input:    "[ ] Task",
			checked:  false,
			expected: "[ ] Task",
		},
		{
			name:     "Check box with space",
			input:    "[ ] Task",
			checked:  true,
			expected: "[x] Task",
		},
		{
			name:     "Uncheck box with lowercase x",
			input:    "[x] Task",
			checked:  false,
			expected: "[ ] Task",
		},
		{
			name:     "Uncheck box with uppercase X",
			input:    "[X] Task",
			checked:  false,
			expected: "[ ] Task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replacement := "[ ] "
			if tt.checked {
				replacement = "[x] "
			}

			result := taskListRegexp.ReplaceAllString(tt.input, replacement)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
