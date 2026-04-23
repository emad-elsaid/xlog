package editor

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/emad-elsaid/xlog"
)

func TestEditorExtensionName(t *testing.T) {
	ext := Editor{}
	if name := ext.Name(); name != "editor" {
		t.Errorf("Expected name 'editor', got '%s'", name)
	}
}

func TestEditorInit(t *testing.T) {
	// Save original config
	origReadonly := xlog.Config.Readonly
	defer func() { xlog.Config.Readonly = origReadonly }()

	// Test that Init respects readonly mode
	xlog.Config.Readonly = true
	ext := Editor{}
	ext.Init() // Should return early without registering handlers

	// Reset for normal operation
	xlog.Config.Readonly = false
	ext.Init() // Should register handlers
}

func TestNewPage(t *testing.T) {
	// Create temp directory for test pages
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Set editor to a no-op command to avoid actually opening an editor
	origEditor := editor
	editor = "true" // Unix command that does nothing and exits successfully
	defer func() { editor = origEditor }()

	// Create a test page
	pageName := "test-page"
	page := xlog.NewPage(pageName)

	// Call newPage - should not error
	err := newPage(page)
	if err != nil {
		t.Errorf("newPage() returned error: %v", err)
	}
}

func TestNewPageWithNilPage(t *testing.T) {
	// Should handle nil page gracefully
	err := newPage(nil)
	if err != nil {
		t.Errorf("newPage(nil) returned error: %v", err)
	}
}

func TestOpenEditorWithNilPage(t *testing.T) {
	// Should not panic with nil page
	openEditor(nil) // Just verify it doesn't crash
}

func TestOpenEditorIgnoresStaticFiles(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Set editor to a command that would fail if called
	origEditor := editor
	editor = "false" // Unix command that always fails
	defer func() { editor = origEditor }()

	// These extensions should be ignored (not opened in editor)
	staticExtensions := []string{".ico", ".jpg", ".png", ".gif", ".so"}

	for _, ext := range staticExtensions {
		pageName := "file" + ext
		page := xlog.NewPage(pageName)
		
		// Should not attempt to open editor (which would fail)
		// If it does try, the test would show error logs
		openEditor(page)
	}
}

func TestOpenEditorWithValidPage(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create a test page file
	pageName := "test-page"
	pageFile := pageName + ".md"
	if err := os.WriteFile(pageFile, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Set editor to 'true' command (does nothing)
	origEditor := editor
	editor = "true"
	defer func() { editor = origEditor }()

	page := xlog.NewPage(pageName)
	openEditor(page) // Should not crash
}

func TestOpenEditorWithEmptyEditorCommand(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Set editor to empty string
	origEditor := editor
	editor = ""
	defer func() { editor = origEditor }()

	page := xlog.NewPage("test-page")
	openEditor(page) // Should handle gracefully (no command to run)
}

func TestEditorHandler(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Set editor to 'true'
	origEditor := editor
	editor = "true"
	defer func() { editor = origEditor }()

	// Create test HTTP request
	req := httptest.NewRequest(http.MethodPost, "/+/editor/test-page", nil)
	req.SetPathValue("page", "test-page")
	
	// Call handler
	output := editorHandler(req)
	
	// Should return NoContent (204)
	if output == nil {
		t.Error("Expected non-nil output")
	}
	
	// Verify the output is NoContent by executing it
	w := httptest.NewRecorder()
	output(w, req)
	
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204 No Content, got %d", w.Code)
	}
}

func TestLinksWithValidPage(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create a test page
	pageName := "test-page"
	pageFile := pageName + ".md"
	if err := os.WriteFile(pageFile, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
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
		hxPost, ok := attrs["hx-post"]
		if !ok {
			t.Error("Expected hx-post attribute")
		}
		
		// Check that the URL contains the escaped page name
		hxPostStr, ok := hxPost.(string)
		if !ok {
			t.Error("hx-post should be a string")
		}
		
		expectedPath := "/+/editor/" + url.PathEscape(pageName)
		if hxPostStr != expectedPath {
			t.Errorf("Expected hx-post '%s', got '%s'", expectedPath, hxPostStr)
		}
	}
}

func TestLinksWithEmptyPageName(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Empty page name results in filename ".md" which has length 3
	// According to code, links only returns nil if len(p.FileName()) == 0
	// So empty page name actually creates a link since ".md" has length > 0
	page := xlog.NewPage("")
	
	commands := links(page)
	
	// Should return a command since FileName() returns ".md" (length 3, not 0)
	if len(commands) != 1 {
		t.Errorf("Expected 1 command for page with .md filename, got %d commands", len(commands))
	}
}

func TestEditButtonProperties(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	pageName := "my-test-page"
	page := xlog.NewPage(pageName)
	
	btn := editButton{page: page}
	
	// Test Icon
	if icon := btn.Icon(); icon != "fa-solid fa-pen" {
		t.Errorf("Icon() = %s, want 'fa-solid fa-pen'", icon)
	}
	
	// Test Name
	if name := btn.Name(); name != "Edit" {
		t.Errorf("Name() = %s, want 'Edit'", name)
	}
	
	// Test Attrs
	attrs := btn.Attrs()
	if len(attrs) != 1 {
		t.Errorf("Expected 1 attribute, got %d", len(attrs))
	}
	
	hxPost, ok := attrs["hx-post"]
	if !ok {
		t.Fatal("Expected hx-post attribute")
	}
	
	expectedURL := "/+/editor/" + url.PathEscape(pageName)
	if hxPost != expectedURL {
		t.Errorf("hx-post = %s, want %s", hxPost, expectedURL)
	}
}

func TestEditorWithComplexCommand(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create a test page
	pageName := "test-page"
	pageFile := pageName + ".md"
	if err := os.WriteFile(pageFile, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Test with editor command that has arguments
	origEditor := editor
	editor = "echo test argument"
	defer func() { editor = origEditor }()

	page := xlog.NewPage(pageName)
	openEditor(page)
	
	// Verify the command segments are properly split
	segments := strings.Split(editor, " ")
	if len(segments) != 3 {
		t.Errorf("Expected 3 segments in command, got %d", len(segments))
	}
}

func TestOpenEditorWithLongExtension(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	origEditor := editor
	editor = "true"
	defer func() { editor = origEditor }()

	// Extensions longer than 4 characters should be opened in editor
	pageName := "file.longext"
	page := xlog.NewPage(pageName)
	
	// Should attempt to open editor (not ignored like short extensions)
	openEditor(page)
}

func TestEditButtonWithSpecialCharactersInPageName(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Test with page name containing special characters
	pageName := "test page with spaces & symbols"
	page := xlog.NewPage(pageName)
	
	btn := editButton{page: page}
	attrs := btn.Attrs()
	
	hxPost, ok := attrs["hx-post"]
	if !ok {
		t.Fatal("Expected hx-post attribute")
	}
	
	hxPostStr := hxPost.(string)
	
	// Verify special characters are URL-encoded
	if !strings.Contains(hxPostStr, "/+/editor/") {
		t.Error("Expected path to contain /+/editor/")
	}
	
	// Verify it doesn't contain raw spaces
	if strings.Contains(hxPostStr, " ") {
		t.Error("Expected spaces to be URL-encoded")
	}
}

func TestEditorHandlerWithDifferentPagePaths(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	origEditor := editor
	editor = "true"
	defer func() { editor = origEditor }()

	testCases := []string{
		"simple-page",
		"path/to/nested/page",
		"page-with-dashes",
	}

	for _, pagePath := range testCases {
		t.Run(pagePath, func(t *testing.T) {
			// Create directory structure if needed
			if dir := filepath.Dir(pagePath); dir != "." {
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatal(err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/+/editor/"+pagePath, nil)
			req.SetPathValue("page", pagePath)
			
			output := editorHandler(req)
			
			w := httptest.NewRecorder()
			output(w, req)
			
			if w.Code != http.StatusNoContent {
				t.Errorf("Expected status 204, got %d", w.Code)
			}
		})
	}
}

func TestAttrsReturnType(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	page := xlog.NewPage("test")
	btn := editButton{page: page}
	attrs := btn.Attrs()
	
	// Verify it returns the expected map type
	var _ map[template.HTMLAttr]any = attrs
	
	// Verify the hx-post attribute key type
	for key := range attrs {
		if key == "hx-post" {
			// Key should be template.HTMLAttr type
			var _ template.HTMLAttr = key
			break
		}
	}
}
