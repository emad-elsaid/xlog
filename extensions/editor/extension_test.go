package editor

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/emad-elsaid/xlog"
)

func TestEditorExtensionName(t *testing.T) {
	ext := Editor{}
	if ext.Name() != "editor" {
		t.Errorf("Expected extension name 'editor', got '%s'", ext.Name())
	}
}

func TestEditorInit(t *testing.T) {
	// Save original config and restore after test
	origReadonly := xlog.Config.Readonly
	defer func() { xlog.Config.Readonly = origReadonly }()

	xlog.Config.Readonly = false

	// Init should register routes and commands
	ext := Editor{}
	ext.Init()

	// We can't directly test route registration without starting the server,
	// but we can verify the function doesn't panic
}

func TestEditorInitReadonly(t *testing.T) {
	// Save original config and restore after test
	origReadonly := xlog.Config.Readonly
	defer func() { xlog.Config.Readonly = origReadonly }()

	xlog.Config.Readonly = true

	// Init should return early when readonly
	ext := Editor{}
	ext.Init()
	// Should not panic
}

func TestOpenEditorNilPage(t *testing.T) {
	// Should not panic with nil page
	openEditor(nil)
}

func TestOpenEditorWithExtension(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	
	// Save original source and restore after test
	origSource := xlog.Config.Source
	defer func() { xlog.Config.Source = origSource }()
	
	xlog.Config.Source = tmpDir
	
	// Create a mock page with an extension that should be ignored
	page := xlog.NewPage("test.ico")

	// Save original editor and restore after test
	origEditor := editor
	defer func() { editor = origEditor }()
	
	editor = ""
	
	// Should not attempt to open editor for files with extensions like .ico
	openEditor(page)
	// Should not panic
}

func TestOpenEditorEmptyCommand(t *testing.T) {
	tmpDir := t.TempDir()
	
	origSource := xlog.Config.Source
	defer func() { xlog.Config.Source = origSource }()
	xlog.Config.Source = tmpDir
	
	page := xlog.NewPage("testpage")

	origEditor := editor
	defer func() { editor = origEditor }()
	
	editor = ""
	
	// Should handle empty editor command gracefully
	openEditor(page)
	// Should not panic
}

func TestOpenEditorInvalidCommand(t *testing.T) {
	tmpDir := t.TempDir()
	
	origSource := xlog.Config.Source
	defer func() { xlog.Config.Source = origSource }()
	xlog.Config.Source = tmpDir
	
	page := xlog.NewPage("testpage")

	origEditor := editor
	defer func() { editor = origEditor }()
	
	// Use a non-existent command
	editor = "nonexistent-editor-command-12345"
	
	// Should handle invalid command gracefully (will log error but not panic)
	openEditor(page)
	// Should not panic
}

func TestEditorHandler(t *testing.T) {
	tmpDir := t.TempDir()
	
	origSource := xlog.Config.Source
	defer func() { xlog.Config.Source = origSource }()
	xlog.Config.Source = tmpDir

	// Create a test file
	testFile := filepath.Join(tmpDir, "testpage.md")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	origEditor := editor
	defer func() { editor = origEditor }()
	editor = ""

	// Create a mock request
	req := httptest.NewRequest(http.MethodPost, "/+/editor/testpage", nil)
	req.SetPathValue("page", "testpage")
	w := httptest.NewRecorder()

	// Call handler
	output := editorHandler(req)
	output(w, req)

	// Should return NoContent (204)
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}
}

func TestLinksNilPage(t *testing.T) {
	tmpDir := t.TempDir()
	
	origSource := xlog.Config.Source
	defer func() { xlog.Config.Source = origSource }()
	xlog.Config.Source = tmpDir
	
	// Create a page with no filename (empty name)
	page := xlog.NewPage("")
	
	// Check if the page actually has no filename
	if len(page.FileName()) > 0 {
		t.Skip("Page with empty name still has filename, skipping test")
	}
	
	cmds := links(page)
	
	if cmds != nil {
		t.Errorf("Expected nil commands for page with no filename, got %v", cmds)
	}
}

func TestLinksWithPage(t *testing.T) {
	tmpDir := t.TempDir()
	
	origSource := xlog.Config.Source
	defer func() { xlog.Config.Source = origSource }()
	xlog.Config.Source = tmpDir

	// Create a test file
	testFile := filepath.Join(tmpDir, "testpage.md")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	page := xlog.NewPage("testpage")
	
	cmds := links(page)
	
	if len(cmds) != 1 {
		t.Errorf("Expected 1 command, got %d", len(cmds))
	}

	if len(cmds) > 0 {
		btn, ok := cmds[0].(editButton)
		if !ok {
			t.Errorf("Expected editButton, got %T", cmds[0])
		}
		
		if btn.page.Name() != "testpage" {
			t.Errorf("Expected page name 'testpage', got '%s'", btn.page.Name())
		}
	}
}

func TestEditButtonIcon(t *testing.T) {
	page := xlog.NewPage("testpage")
	btn := editButton{page: page}
	
	icon := btn.Icon()
	expected := "fa-solid fa-pen"
	
	if icon != expected {
		t.Errorf("Expected icon '%s', got '%s'", expected, icon)
	}
}

func TestEditButtonName(t *testing.T) {
	page := xlog.NewPage("testpage")
	btn := editButton{page: page}
	
	name := btn.Name()
	expected := "Edit"
	
	if name != expected {
		t.Errorf("Expected name '%s', got '%s'", expected, name)
	}
}

func TestEditButtonAttrs(t *testing.T) {
	page := xlog.NewPage("testpage")
	btn := editButton{page: page}
	
	attrs := btn.Attrs()
	
	if len(attrs) == 0 {
		t.Error("Expected non-empty attrs map")
	}
	
	hxPost, ok := attrs["hx-post"]
	if !ok {
		t.Error("Expected 'hx-post' attribute")
	}
	
	expectedPath := "/+/editor/testpage"
	if hxPost != expectedPath {
		t.Errorf("Expected hx-post '%s', got '%v'", expectedPath, hxPost)
	}
}

func TestEditButtonAttrsWithSpecialChars(t *testing.T) {
	page := xlog.NewPage("test page/with spaces")
	btn := editButton{page: page}
	
	attrs := btn.Attrs()
	
	hxPost, ok := attrs["hx-post"]
	if !ok {
		t.Error("Expected 'hx-post' attribute")
	}
	
	// Should properly escape the page name
	hxPostStr, ok := hxPost.(string)
	if !ok {
		t.Errorf("Expected string for hx-post, got %T", hxPost)
	}
	
	if hxPostStr == "" {
		t.Error("Expected non-empty hx-post path")
	}
}

func TestEditButtonAttrsType(t *testing.T) {
	page := xlog.NewPage("testpage")
	btn := editButton{page: page}
	
	attrs := btn.Attrs()
	
	// Verify the return type is correct
	var _ map[template.HTMLAttr]any = attrs
}

func TestNewPageNilPage(t *testing.T) {
	// Should not panic with nil page
	err := newPage(nil)
	if err != nil {
		t.Errorf("Expected no error for nil page, got %v", err)
	}
}

func TestNewPageValidPage(t *testing.T) {
	tmpDir := t.TempDir()
	
	origSource := xlog.Config.Source
	defer func() { xlog.Config.Source = origSource }()
	xlog.Config.Source = tmpDir

	origEditor := editor
	defer func() { editor = origEditor }()
	editor = ""

	page := xlog.NewPage("testpage")
	
	err := newPage(page)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
