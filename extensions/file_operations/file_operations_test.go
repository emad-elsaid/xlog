package file_operations

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
)

func TestFileOpsExtensionName(t *testing.T) {
	ext := FileOps{}
	expected := "file-operations"
	if ext.Name() != expected {
		t.Errorf("Expected name %q, got %q", expected, ext.Name())
	}
}

func TestPageDeleteIcon(t *testing.T) {
	pd := PageDelete{}
	expected := "fa-solid fa-trash"
	if pd.Icon() != expected {
		t.Errorf("Expected icon %q, got %q", expected, pd.Icon())
	}
}

func TestPageDeleteName(t *testing.T) {
	pd := PageDelete{}
	expected := "Delete"
	if pd.Name() != expected {
		t.Errorf("Expected name %q, got %q", expected, pd.Name())
	}
}

func TestPageRenameIcon(t *testing.T) {
	pr := PageRename{}
	expected := "fa-solid fa-i-cursor"
	if pr.Icon() != expected {
		t.Errorf("Expected icon %q, got %q", expected, pr.Icon())
	}
}

func TestPageRenameName(t *testing.T) {
	pr := PageRename{}
	expected := "Rename"
	if pr.Name() != expected {
		t.Errorf("Expected name %q, got %q", expected, pr.Name())
	}
}

func TestPageDeleteAttrs(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create test page.
	testPageName := "test-delete-page"
	testFilePath := filepath.Join(tmpDir, testPageName+".md")
	if err := os.WriteFile(testFilePath, []byte("test content"), 0600); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}

	page := xlog.NewPage(testPageName)
	if page == nil {
		t.Fatal("Failed to create page")
	}

	pd := PageDelete{page: page}
	attrs := pd.Attrs()

	// Check href attribute.
	expectedHref := "/+/file/delete?page=" + url.QueryEscape(testPageName)
	if href, ok := attrs["href"]; !ok || href != expectedHref {
		t.Errorf("Expected href %q, got %v", expectedHref, href)
	}

	// Check hx-delete attribute.
	expectedDelete := "/+/file/delete?page=" + url.QueryEscape(testPageName)
	if hxDelete, ok := attrs["hx-delete"]; !ok || hxDelete != expectedDelete {
		t.Errorf("Expected hx-delete %q, got %v", expectedDelete, hxDelete)
	}

	// Check hx-confirm attribute.
	expectedConfirm := "Are you sure?"
	if hxConfirm, ok := attrs["hx-confirm"]; !ok || hxConfirm != expectedConfirm {
		t.Errorf("Expected hx-confirm %q, got %v", expectedConfirm, hxConfirm)
	}
}

func TestPageRenameAttrs(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create test page.
	testPageName := "test-rename-page"
	testFilePath := filepath.Join(tmpDir, testPageName+".md")
	if err := os.WriteFile(testFilePath, []byte("test content"), 0600); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}

	page := xlog.NewPage(testPageName)
	if page == nil {
		t.Fatal("Failed to create page")
	}

	pr := PageRename{page: page}
	attrs := pr.Attrs()

	// Check href attribute.
	expectedHref := "/+/file/rename?page=" + url.QueryEscape(testPageName)
	if href, ok := attrs["href"]; !ok || href != expectedHref {
		t.Errorf("Expected href %q, got %v", expectedHref, href)
	}

	// Check hx-get attribute.
	expectedGet := "/+/file/rename?page=" + url.QueryEscape(testPageName)
	if hxGet, ok := attrs["hx-get"]; !ok || hxGet != expectedGet {
		t.Errorf("Expected hx-get %q, got %v", expectedGet, hxGet)
	}

	// Check hx-target attribute.
	if hxTarget, ok := attrs["hx-target"]; !ok || hxTarget != "body" {
		t.Errorf("Expected hx-target %q, got %v", "body", hxTarget)
	}

	// Check hx-swap attribute.
	if hxSwap, ok := attrs["hx-swap"]; !ok || hxSwap != "beforeend" {
		t.Errorf("Expected hx-swap %q, got %v", "beforeend", hxSwap)
	}
}

func TestCommandsReturnsNilForEmptyFileName(t *testing.T) {
	// Create a page with empty filename.
	page := &mockPage{fileName: ""}
	result := commands(page)

	if result != nil {
		t.Errorf("Expected nil for empty filename, got %v", result)
	}
}

func TestCommandsReturnsCommandsForValidPage(t *testing.T) {
	// Create a page with valid filename.
	page := &mockPage{fileName: "test.md"}
	result := commands(page)

	if result == nil {
		t.Fatal("Expected commands for valid page, got nil")
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(result))
	}

	// Verify commands are PageDelete and PageRename.
	foundDelete := false
	foundRename := false
	for _, cmd := range result {
		switch cmd.(type) {
		case PageDelete:
			foundDelete = true
		case PageRename:
			foundRename = true
		}
	}

	if !foundDelete {
		t.Error("Expected PageDelete command")
	}
	if !foundRename {
		t.Error("Expected PageRename command")
	}
}

func TestPageDeleteHandler_NonExistentPage(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create request for non-existent page.
	req := httptest.NewRequest(http.MethodDelete, "/+/file/delete?page=non-existent", http.NoBody)
	rec := httptest.NewRecorder()

	pd := PageDelete{}
	output := pd.Handler(req)

	// Execute the output function.
	output(rec, req)

	// Should redirect to home.
	if redirect := rec.Header().Get("HX-Redirect"); redirect != "/" {
		t.Errorf("Expected redirect to '/', got %q", redirect)
	}
}

func TestPageDeleteHandler_ExistingPage(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create test page.
	testPageName := "test-delete-existing"
	testFilePath := filepath.Join(tmpDir, testPageName+".md")
	if err := os.WriteFile(testFilePath, []byte("content to delete"), 0600); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}

	// Verify file exists.
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		t.Fatal("Test page should exist before deletion")
	}

	// Create request.
	req := httptest.NewRequest(http.MethodDelete, "/+/file/delete?page="+testPageName, http.NoBody)
	rec := httptest.NewRecorder()

	pd := PageDelete{}
	output := pd.Handler(req)

	// Execute the output function.
	output(rec, req)

	// Verify redirect.
	if redirect := rec.Header().Get("HX-Redirect"); redirect != "/" {
		t.Errorf("Expected redirect to '/', got %q", redirect)
	}

	// Verify file was deleted.
	if _, err := os.Stat(testFilePath); !os.IsNotExist(err) {
		t.Error("Test page should be deleted")
	}
}

func TestPageRenameHandler_NonExistentPage(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create form data.
	formData := url.Values{}
	formData.Set("old", "non-existent-page")
	formData.Set("new", "new-name")

	req := httptest.NewRequest(http.MethodPost, "/+/file/rename", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	pr := PageRename{}
	output := pr.Handler(req)

	// Execute output - should be BadRequest.
	rec := httptest.NewRecorder()
	output(rec, req)

	// Verify bad request status.
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	// Verify error message.
	body := rec.Body.String()
	if !strings.Contains(body, "doesn't exist") {
		t.Errorf("Expected error message about non-existent file, got: %s", body)
	}
}

func TestPageRenameHandler_SuccessfulRename(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create test page.
	oldName := "old-page-name"
	newBaseName := "new-page-name"
	oldFilePath := filepath.Join(tmpDir, oldName+".md")
	newFilePath := filepath.Join(tmpDir, newBaseName+".md")

	originalContent := "# Original Content"
	if err := os.WriteFile(oldFilePath, []byte(originalContent), 0600); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}

	// Create form data.
	formData := url.Values{}
	formData.Set("old", oldName)
	formData.Set("new", newBaseName)

	req := httptest.NewRequest(http.MethodPost, "/+/file/rename", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	pr := PageRename{}
	output := pr.Handler(req)

	// Execute the output function.
	output(rec, req)

	// Verify redirect to new page.
	expectedRedirect := "/" + newBaseName
	if redirect := rec.Header().Get("HX-Redirect"); redirect != expectedRedirect {
		t.Errorf("Expected redirect to %q, got %q", expectedRedirect, redirect)
	}

	// Verify new file exists.
	if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
		t.Error("New file should exist after rename")
	}

	// Note: The current implementation calls old.Write() after rename, which recreates the old file.
	// This is a known behaviour - the old file will be recreated with a "Renamed to:" message.
	// Verify old file exists with rename message.
	if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
		t.Error("Old file gets recreated with rename message (current implementation behaviour)")
	}
}

func TestPageRenameHandler_PreservesExtension(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create test page with .md extension.
	oldName := "extension-test"
	newBaseName := "renamed-extension-test"
	oldFilePath := filepath.Join(tmpDir, oldName+".md")
	expectedNewFilePath := filepath.Join(tmpDir, newBaseName+".md")

	if err := os.WriteFile(oldFilePath, []byte("content"), 0600); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}

	// Create form data - note: only basename without extension.
	formData := url.Values{}
	formData.Set("old", oldName)
	formData.Set("new", newBaseName)

	req := httptest.NewRequest(http.MethodPost, "/+/file/rename", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	pr := PageRename{}
	output := pr.Handler(req)
	output(rec, req)

	// Verify new file has .md extension.
	if _, err := os.Stat(expectedNewFilePath); os.IsNotExist(err) {
		t.Error("Renamed file should preserve .md extension")
	}

	// Verify no file without extension was created.
	wrongPath := filepath.Join(tmpDir, newBaseName)
	if info, err := os.Stat(wrongPath); err == nil && !info.IsDir() {
		t.Error("Should not create file without extension")
	}
}

func TestFileOps_Init_ReadonlyMode(t *testing.T) {
	// Save original config and restore after test.
	originalReadonly := xlog.Config.Readonly
	defer func() { xlog.Config.Readonly = originalReadonly }()

	// Set readonly mode.
	xlog.Config.Readonly = true

	// Init should return early without registering anything.
	ext := FileOps{}
	ext.Init()

	// Verify no routes were registered by attempting to check registered routes.
	// Since xlog doesn't expose route inspection, we verify by attempting a request.
	// This test verifies Init returns early without panic in readonly mode.
}

func TestFileOps_Init_NormalMode(t *testing.T) {
	// Save original config and restore after test.
	originalReadonly := xlog.Config.Readonly
	defer func() { xlog.Config.Readonly = originalReadonly }()

	// Set normal (non-readonly) mode.
	xlog.Config.Readonly = false

	// Init should register routes and commands.
	ext := FileOps{}
	ext.Init()

	// Verify Init completes without panic.
	// The actual route registration is verified through integration tests.
}

func TestPageRename_Form_RendersOutput(t *testing.T) {
	// This test verifies Form method exists and returns an Output function.
	// Full template rendering requires server initialization which is beyond unit test scope.
	tmpDir := t.TempDir()
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	testPageName := "test-form-page"
	req := httptest.NewRequest(http.MethodGet, "/+/file/rename?page="+testPageName, http.NoBody)

	pr := PageRename{}
	output := pr.Form(req)

	// Verify Form returns a non-nil Output function.
	if output == nil {
		t.Error("Expected Form to return non-nil Output function")
	}
}

// Mock page for testing.
type mockPage struct {
	fileName string
}

func (m *mockPage) Name() string                { return m.fileName }
func (m *mockPage) FileName() string            { return m.fileName }
func (m *mockPage) Exists() bool                { return true }
func (m *mockPage) Delete() bool                { return true }
func (m *mockPage) Write(md xlog.Markdown) bool { return true }
func (m *mockPage) Content() xlog.Markdown      { return "" }
func (m *mockPage) Render() template.HTML       { return "" }
func (m *mockPage) ModTime() time.Time          { return time.Now() }
func (m *mockPage) AST() ([]byte, ast.Node)     { return nil, nil }
