package todo

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	east "github.com/emad-elsaid/xlog/markdown/extension/ast"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/text"
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

func TestToggleHandler(t *testing.T) {
	// Create temp directory for test pages
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create test page
	testPageName := "test-page"
	testContent := "# Test\n[ ] Task 1\n[x] Task 2\n"
	if err := os.WriteFile(testPageName+".md", []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}

	tests := []struct {
		name           string
		page           string
		pos            string
		checked        string
		expectedStatus int
	}{
		{
			name:           "Check unchecked task",
			page:           testPageName,
			pos:            "8",
			checked:        "true",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Uncheck checked task",
			page:           testPageName,
			pos:            "18",
			checked:        "",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Page not found",
			page:           "nonexistent",
			pos:            "0",
			checked:        "",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid position",
			page:           testPageName,
			pos:            "invalid",
			checked:        "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Position exceeds content length",
			page:           testPageName,
			pos:            "9999",
			checked:        "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset test page content for each test
			if err := os.WriteFile(testPageName+".md", []byte(testContent), 0644); err != nil {
				t.Fatalf("Failed to reset test page: %v", err)
			}

			// Create form data
			form := url.Values{}
			form.Add("page", tt.page)
			form.Add("pos", tt.pos)
			if tt.checked != "" {
				form.Add("checked", tt.checked)
			}

			// Create request
			req := httptest.NewRequest("POST", "/+/todo", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Execute handler - get the Output (http.HandlerFunc)
			output := toggleHandler(req)

			// Create a response recorder to capture the response
			w := httptest.NewRecorder()
			output(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestTaskCheckBoxHTMLRenderer_RegisterFuncs(t *testing.T) {
	r := &TaskCheckBoxHTMLRenderer{}

	// Create a mock registerer
	var registered bool
	var registeredKind ast.NodeKind

	mockRegisterer := &mockNodeRendererFuncRegisterer{
		registerFunc: func(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
			registered = true
			registeredKind = kind
		},
	}

	r.RegisterFuncs(mockRegisterer)

	if !registered {
		t.Error("Expected RegisterFuncs to call Register")
	}

	if registeredKind != east.KindTaskCheckBox {
		t.Errorf("Expected kind %v, got %v", east.KindTaskCheckBox, registeredKind)
	}
}

type mockNodeRendererFuncRegisterer struct {
	registerFunc func(ast.NodeKind, renderer.NodeRendererFunc)
}

func (m *mockNodeRendererFuncRegisterer) Register(kind ast.NodeKind, fn renderer.NodeRendererFunc) {
	if m.registerFunc != nil {
		m.registerFunc(kind, fn)
	}
}

func TestTaskCheckBoxHTMLRenderer_NotEntering(t *testing.T) {
	r := &TaskCheckBoxHTMLRenderer{}
	taskCheckBox := east.NewTaskCheckBox(false)

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	// Test with entering=false - should return early and write nothing
	status, err := r.renderTaskCheckBox(writer, []byte("test"), taskCheckBox, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if status != ast.WalkContinue {
		t.Errorf("Expected WalkContinue, got %v", status)
	}

	writer.Flush()

	// Should produce no output when not entering
	if buf.Len() > 0 {
		t.Errorf("Expected no output when entering=false, got: %s", buf.String())
	}
}

func TestTaskCheckBoxHTMLRenderer_Basic(t *testing.T) {
	r := &TaskCheckBoxHTMLRenderer{}

	tests := []struct {
		name      string
		isChecked bool
		readonly  bool
		expected  []string
	}{
		{
			name:      "Unchecked checkbox (editable)",
			isChecked: false,
			readonly:  false,
			expected:  []string{`<input name="checked" type="checkbox"`},
		},
		{
			name:      "Checked checkbox (editable)",
			isChecked: true,
			readonly:  false,
			expected:  []string{`checked=""`},
		},
		{
			name:      "Readonly checkbox",
			isChecked: false,
			readonly:  true,
			expected:  []string{`disabled=""`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save and restore config
			oldReadonly := Config.Readonly
			Config.Readonly = tt.readonly
			defer func() { Config.Readonly = oldReadonly }()

			// Create task checkbox node
			taskCheckBox := east.NewTaskCheckBox(tt.isChecked)

			// Create parent text block with lines
			parent := ast.NewTextBlock()
			// Add a segment to the lines - renderer expects at least one
			if !tt.readonly {
				seg := text.NewSegment(0, 10)
				parent.Lines().Append(seg)
			}
			taskCheckBox.SetParent(parent)

			// Render
			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)
			source := []byte("test content")

			status, err := r.renderTaskCheckBox(writer, source, taskCheckBox, true)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if status != ast.WalkContinue {
				t.Errorf("Expected WalkContinue, got %v", status)
			}

			writer.Flush()
			output := buf.String()

			for _, exp := range tt.expected {
				if !strings.Contains(output, exp) {
					t.Errorf("Expected output to contain %q, got: %s", exp, output)
				}
			}
		})
	}
}

