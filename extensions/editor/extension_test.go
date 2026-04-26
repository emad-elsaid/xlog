package editor

import (
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
		t.Errorf("Expected name 'editor', got '%s'", ext.Name())
	}
}

func TestOpenEditorWithNilPage(t *testing.T) {
	// Should not panic
	openEditor(nil)
}

func TestOpenEditorWithInvalidExtension(t *testing.T) {
	// Test files with extensions that should be ignored (.ico, .jpeg, etc.)
	testCases := []string{
		"test.ico",
		"test.jpeg",
		"test.so",
		"test.png",
	}
	
	for _, name := range testCases {
		page := xlog.NewPage(name)
		// Should not attempt to open editor (returns early)
		openEditor(page)
	}
}

func TestOpenEditorWithValidPage(t *testing.T) {
	// Save original editor value
	originalEditor := editor
	defer func() { editor = originalEditor }()
	
	// Set editor to a no-op command
	editor = "true"
	
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	
	// Create a test page
	pageName := "test-page.md"
	pagePath := filepath.Join(tmpDir, pageName)
	if err := os.WriteFile(pagePath, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}
	
	page := xlog.NewPage(pageName)
	// Should not panic
	openEditor(page)
}

func TestOpenEditorWithEmptyEditor(t *testing.T) {
	// Save original editor value
	originalEditor := editor
	defer func() { editor = originalEditor }()
	
	// Set editor to empty string
	editor = ""
	
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	
	// Create a test page
	pageName := "test-page.md"
	pagePath := filepath.Join(tmpDir, pageName)
	if err := os.WriteFile(pagePath, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}
	
	page := xlog.NewPage(pageName)
	// Should return early without error
	openEditor(page)
}

func TestEditorHandler(t *testing.T) {
	// Save original editor value
	originalEditor := editor
	defer func() { editor = originalEditor }()
	
	// Set editor to a no-op command
	editor = "true"
	
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	
	// Create a test page
	pageName := "test-page.md"
	pagePath := filepath.Join(tmpDir, pageName)
	if err := os.WriteFile(pagePath, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}
	
	// Create a test request
	req := httptest.NewRequest(http.MethodPost, "/+/editor/"+pageName, nil)
	req.SetPathValue("page", pageName)
	
	// Call handler
	output := editorHandler(req)
	
	// Execute the output function
	w := httptest.NewRecorder()
	output(w, req)
	
	// Should return NoContent (204)
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestLinksWithEmptyFileName(t *testing.T) {
	page := xlog.NewPage("")
	commands := links(page)
	
	// Empty page should still have a button since FileName() might not be empty
	// Just verify we don't panic
	if commands == nil {
		t.Log("Commands is nil for empty page name")
	}
}

func TestLinksWithValidPage(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	
	// Create a test page
	pageName := "test-page.md"
	pagePath := filepath.Join(tmpDir, pageName)
	if err := os.WriteFile(pagePath, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}
	
	page := xlog.NewPage(pageName)
	commands := links(page)
	
	if len(commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(commands))
	}
	
	if len(commands) > 0 {
		btn := commands[0].(editButton)
		if btn.Icon() != "fa-solid fa-pen" {
			t.Errorf("Expected icon 'fa-solid fa-pen', got '%s'", btn.Icon())
		}
		if btn.Name() != "Edit" {
			t.Errorf("Expected name 'Edit', got '%s'", btn.Name())
		}
		
		attrs := btn.Attrs()
		if _, ok := attrs["hx-post"]; !ok {
			t.Error("Expected hx-post attribute to be present")
		}
	}
}

func TestNewPageCallback(t *testing.T) {
	// Save original editor value
	originalEditor := editor
	defer func() { editor = originalEditor }()
	
	// Set editor to a no-op command
	editor = "true"
	
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	
	// Create a test page
	pageName := "test-page.md"
	pagePath := filepath.Join(tmpDir, pageName)
	if err := os.WriteFile(pagePath, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}
	
	page := xlog.NewPage(pageName)
	err := newPage(page)
	
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestEditButtonStructMethods(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	
	// Create a test page
	pageName := "test-page.md"
	pagePath := filepath.Join(tmpDir, pageName)
	if err := os.WriteFile(pagePath, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}
	
	page := xlog.NewPage(pageName)
	btn := editButton{page: page}
	
	if btn.Icon() != "fa-solid fa-pen" {
		t.Errorf("Expected icon 'fa-solid fa-pen', got '%s'", btn.Icon())
	}
	
	if btn.Name() != "Edit" {
		t.Errorf("Expected name 'Edit', got '%s'", btn.Name())
	}
	
	attrs := btn.Attrs()
	hxPost, ok := attrs["hx-post"]
	if !ok {
		t.Fatal("Expected hx-post attribute to be present")
	}
	
	expectedPath := "/+/editor/" + pageName
	if hxPost != expectedPath {
		t.Errorf("Expected hx-post '%s', got '%s'", expectedPath, hxPost)
	}
}
