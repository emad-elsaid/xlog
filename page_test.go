package xlog

import (
	"html/template"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestPageName(t *testing.T) {
	p := &page{name: "test-page"}
	if p.Name() != "test-page" {
		t.Errorf("Expected name 'test-page', got '%s'", p.Name())
	}
}

func TestPageFileName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"simple", "simple.md"},
		{"with/slash", filepath.FromSlash("with/slash.md")},
		{"nested/path/page", filepath.FromSlash("nested/path/page.md")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &page{name: tt.name}
			if p.FileName() != tt.expected {
				t.Errorf("Expected filename '%s', got '%s'", tt.expected, p.FileName())
			}
		})
	}
}

func TestPageExists(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	// Test non-existent page
	p := &page{name: "nonexistent"}
	if p.Exists() {
		t.Error("Expected page to not exist")
	}

	// Create a page file
	os.WriteFile("test.md", []byte("content"), 0644)
	p2 := &page{name: "test"}
	if !p2.Exists() {
		t.Error("Expected page to exist")
	}
}

func TestPageContent(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	content := "# Test Page\n\nThis is test content."
	os.WriteFile("test.md", []byte(content), 0644)

	p := &page{name: "test"}
	got := p.Content()
	if string(got) != content {
		t.Errorf("Expected content '%s', got '%s'", content, got)
	}

	// Test non-existent page returns empty content
	p2 := &page{name: "missing"}
	if p2.Content() != "" {
		t.Error("Expected empty content for non-existent page")
	}
}

func TestPageWrite(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	// Write new page
	p := &page{name: "test"}
	content := Markdown("# New Page\n\nContent here.")
	if !p.Write(content) {
		t.Error("Write failed")
	}

	// Verify file was created
	if !p.Exists() {
		t.Error("Page file not created")
	}

	// Verify content
	got := p.Content()
	if got != content {
		t.Errorf("Expected content '%s', got '%s'", content, got)
	}

	// Test write with nested path
	p2 := &page{name: "nested/path/page"}
	if !p2.Write(Markdown("nested content")) {
		t.Error("Write failed for nested path")
	}
	if !p2.Exists() {
		t.Error("Nested page file not created")
	}
}

func TestPageWriteNormalizesLineEndings(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	p := &page{name: "test"}
	content := Markdown("line1\r\nline2\r\nline3")
	p.Write(content)

	got := p.Content()
	expected := Markdown("line1\nline2\nline3")
	if got != expected {
		t.Errorf("Expected normalized content '%s', got '%s'", expected, got)
	}
}

func TestPageDelete(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	// Create a page
	p := &page{name: "test"}
	p.Write(Markdown("content"))

	if !p.Exists() {
		t.Error("Page should exist before delete")
	}

	// Delete it
	if !p.Delete() {
		t.Error("Delete failed")
	}

	if p.Exists() {
		t.Error("Page should not exist after delete")
	}

	// Delete non-existent page should still return true
	p2 := &page{name: "nonexistent"}
	if !p2.Delete() {
		t.Error("Delete of non-existent page should return true")
	}
}

func TestPageModTime(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	// Non-existent page returns zero time
	p := &page{name: "test"}
	if !p.ModTime().IsZero() {
		t.Error("Expected zero time for non-existent page")
	}

	// Create page and check mod time
	before := time.Now()
	p.Write(Markdown("content"))
	after := time.Now()

	modTime := p.ModTime()
	if modTime.Before(before) || modTime.After(after) {
		t.Errorf("ModTime %v not between %v and %v", modTime, before, after)
	}
}

func TestPageRender(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	// Simple markdown rendering
	p := &page{name: "test"}
	p.Write(Markdown("# Title\n\nParagraph text."))

	html := p.Render()
	htmlStr := string(html)

	// Should contain rendered HTML elements
	if len(htmlStr) == 0 {
		t.Error("Render returned empty HTML")
	}

	// Basic check for HTML tags (actual rendering depends on markdown converter)
	// We're just verifying it returns something and doesn't panic
}

func TestPageAST(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	p := &page{name: "test"}
	content := Markdown("# Heading\n\nSome text.")
	p.Write(content)

	source, tree := p.AST()
	if len(source) == 0 {
		t.Error("Expected non-empty source from AST")
	}
	if tree == nil {
		t.Error("Expected non-nil AST tree")
	}
}

func TestPageCaching(t *testing.T) {
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)

	p := &page{name: "test"}
	p.Write(Markdown("# Original"))

	// First AST call should cache
	_, ast1 := p.AST()
	if ast1 == nil {
		t.Error("Expected non-nil AST")
	}

	// Second call should return same cached AST
	_, ast2 := p.AST()
	if ast2 != ast1 {
		t.Error("Expected cached AST to be reused")
	}

	// After write, cache should be cleared
	time.Sleep(10 * time.Millisecond) // Ensure modtime changes
	p.Write(Markdown("# Modified"))
	_, ast3 := p.AST()
	if ast3 == ast1 {
		t.Error("Expected new AST after write, cache should be cleared")
	}
}

func TestPageClearCache(t *testing.T) {
	p := &page{name: "test"}
	content := Markdown("test")
	p.content = &content
	p.lastUpdate = time.Now()

	if p.content == nil {
		t.Error("Content should be set before clearCache")
	}

	p.clearCache()

	if p.content != nil {
		t.Error("Content should be nil after clearCache")
	}
	if p.ast != nil {
		t.Error("AST should be nil after clearCache")
	}
	if !p.lastUpdate.IsZero() {
		t.Error("lastUpdate should be zero after clearCache")
	}
}

func TestDynamicPageInterface(t *testing.T) {
	dp := DynamicPage{
		NameVal: "dynamic-test",
		RenderFn: func() template.HTML {
			return template.HTML("<p>Dynamic content</p>")
		},
	}

	// Test all interface methods
	if dp.Name() != "dynamic-test" {
		t.Errorf("Expected name 'dynamic-test', got '%s'", dp.Name())
	}

	if dp.FileName() != "" {
		t.Error("DynamicPage should return empty filename")
	}

	if dp.Exists() {
		t.Error("DynamicPage should not exist")
	}

	if dp.Content() != "" {
		t.Error("DynamicPage should return empty content")
	}

	if dp.Delete() {
		t.Error("DynamicPage Delete should return false")
	}

	if dp.Write(Markdown("test")) {
		t.Error("DynamicPage Write should return false")
	}

	if !dp.ModTime().IsZero() {
		t.Error("DynamicPage should return zero ModTime")
	}

	src, tree := dp.AST()
	if src != nil || tree != nil {
		t.Error("DynamicPage AST should return nil, nil")
	}

	html := dp.Render()
	if html != template.HTML("<p>Dynamic content</p>") {
		t.Errorf("Expected custom render output, got '%s'", html)
	}
}

func TestDynamicPageRenderWithoutFunction(t *testing.T) {
	dp := DynamicPage{
		NameVal:  "no-render",
		RenderFn: nil,
	}

	html := dp.Render()
	if html != "" {
		t.Error("DynamicPage with no RenderFn should return empty HTML")
	}
}

func TestMarkdownType(t *testing.T) {
	// Test that Markdown is a distinct type from string
	var md Markdown = "test content"
	if string(md) != "test content" {
		t.Error("Markdown should convert to string properly")
	}
}
