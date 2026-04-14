package html

import (
	"context"
	"html/template"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
)

func TestHTMLExtensionName(t *testing.T) {
	ext := HTML{}
	if ext.Name() != "html" {
		t.Errorf("Expected extension name 'html', got '%s'", ext.Name())
	}
}

func TestHTMLSource_Page(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Create test HTML files
	os.WriteFile("test.html", []byte("<h1>Test HTML</h1>"), 0644)
	os.WriteFile("test2.htm", []byte("<h1>Test HTM</h1>"), 0644)
	os.WriteFile("test3.xhtml", []byte("<h1>Test XHTML</h1>"), 0644)

	source := &htmlSource{}

	tests := []struct {
		name     string
		pageName string
		wantNil  bool
	}{
		{"html extension", "test", false},
		{"htm extension", "test2", false},
		{"xhtml extension", "test3", false},
		{"non-existent page", "nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := source.Page(tt.pageName)
			if tt.wantNil {
				if page != nil {
					t.Errorf("Expected nil page, got %v", page)
				}
			} else {
				if page == nil {
					t.Errorf("Expected non-nil page, got nil")
				} else if page.Name() != tt.pageName {
					t.Errorf("Expected page name '%s', got '%s'", tt.pageName, page.Name())
				}
			}
		})
	}
}

func TestHTMLSource_Each(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Create test HTML files in various directories
	os.WriteFile("root.html", []byte("<h1>Root</h1>"), 0644)
	os.MkdirAll("subdir", 0755)
	os.WriteFile("subdir/nested.htm", []byte("<h1>Nested</h1>"), 0644)
	os.WriteFile("another.xhtml", []byte("<h1>Another</h1>"), 0644)
	os.WriteFile("not-html.txt", []byte("Not HTML"), 0644) // Should be ignored

	source := &htmlSource{}
	found := make(map[string]bool)

	ctx := context.Background()
	source.Each(ctx, func(p xlog.Page) {
		found[p.Name()] = true
	})

	expected := []string{"root", "subdir/nested", "another"}
	for _, name := range expected {
		if !found[name] {
			t.Errorf("Expected to find page '%s'", name)
		}
	}

	if found["not-html"] {
		t.Errorf("Should not have found non-HTML file")
	}
}

func TestHTMLSource_Each_ContextCancellation(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	// Create multiple files
	for i := 0; i < 10; i++ {
		os.WriteFile(filepath.Join(tmpDir, "file"+string(rune('0'+i))+".html"), []byte("<h1>Test</h1>"), 0644)
	}

	source := &htmlSource{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	count := 0
	source.Each(ctx, func(p xlog.Page) {
		count++
	})

	// Should stop early due to context cancellation
	// Note: depending on timing, it might process 0-1 files before seeing cancellation
	if count > 5 {
		t.Errorf("Expected context cancellation to stop iteration early, got %d files", count)
	}
}

func TestPage_Name(t *testing.T) {
	p := &page{name: "test/page", ext: ".html"}
	if p.Name() != "test/page" {
		t.Errorf("Expected name 'test/page', got '%s'", p.Name())
	}
}

func TestPage_FileName(t *testing.T) {
	p := &page{name: "test/page", ext: ".html"}
	expected := filepath.FromSlash("test/page") + ".html"
	if p.FileName() != expected {
		t.Errorf("Expected filename '%s', got '%s'", expected, p.FileName())
	}
}

func TestPage_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.WriteFile("exists.html", []byte("<h1>Exists</h1>"), 0644)

	tests := []struct {
		name   string
		page   *page
		exists bool
	}{
		{"existing page", &page{name: "exists", ext: ".html"}, true},
		{"non-existing page", &page{name: "nothere", ext: ".html"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.page.Exists() != tt.exists {
				t.Errorf("Expected Exists() = %v, got %v", tt.exists, tt.page.Exists())
			}
		})
	}
}

func TestPage_Content(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	content := "<h1>Test Content</h1>"
	os.WriteFile("test.html", []byte(content), 0644)

	p := &page{name: "test", ext: ".html"}
	got := string(p.Content())
	if got != content {
		t.Errorf("Expected content '%s', got '%s'", content, got)
	}
}

func TestPage_Content_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	p := &page{name: "nonexistent", ext: ".html"}
	content := p.Content()
	if string(content) != "" {
		t.Errorf("Expected empty content for non-existent file, got '%s'", content)
	}
}

func TestPage_Render(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	htmlContent := "<h1>Test Content</h1>"
	os.WriteFile("test.html", []byte(htmlContent), 0644)

	p := &page{name: "test", ext: ".html"}
	rendered := p.Render()
	expected := template.HTML(htmlContent)
	if rendered != expected {
		t.Errorf("Expected rendered content '%s', got '%s'", expected, rendered)
	}
}

func TestPage_ModTime(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	before := time.Now().Add(-1 * time.Second)
	os.WriteFile("test.html", []byte("<h1>Test</h1>"), 0644)
	after := time.Now().Add(1 * time.Second)

	p := &page{name: "test", ext: ".html"}
	modTime := p.ModTime()

	if modTime.Before(before) || modTime.After(after) {
		t.Errorf("ModTime %v should be between %v and %v", modTime, before, after)
	}
}

func TestPage_ModTime_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	p := &page{name: "nonexistent", ext: ".html"}
	modTime := p.ModTime()
	if !modTime.IsZero() {
		t.Errorf("Expected zero time for non-existent file, got %v", modTime)
	}
}

func TestPage_Write(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	p := &page{name: "test", ext: ".html"}
	content := xlog.Markdown("<h1>New Content</h1>")

	if !p.Write(content) {
		t.Error("Write() should return true on success")
	}

	if !p.Exists() {
		t.Error("File should exist after Write()")
	}

	got := p.Content()
	if string(got) != string(content) {
		t.Errorf("Expected content '%s', got '%s'", content, got)
	}
}

func TestPage_Write_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	p := &page{name: "subdir/nested/test", ext: ".html"}
	content := xlog.Markdown("<h1>Nested Content</h1>")

	if !p.Write(content) {
		t.Error("Write() should return true on success")
	}

	if !p.Exists() {
		t.Error("File should exist after Write() in nested directory")
	}
}

func TestPage_Write_NormalizesLineEndings(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	p := &page{name: "test", ext: ".html"}
	content := xlog.Markdown("<h1>Line 1</h1>\r\n<h2>Line 2</h2>")

	p.Write(content)

	got := p.Content()
	expected := "<h1>Line 1</h1>\n<h2>Line 2</h2>"
	if string(got) != expected {
		t.Errorf("Expected normalized content '%s', got '%s'", expected, got)
	}
}

func TestPage_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	os.WriteFile("test.html", []byte("<h1>Test</h1>"), 0644)

	p := &page{name: "test", ext: ".html"}
	if !p.Exists() {
		t.Error("File should exist before Delete()")
	}

	if !p.Delete() {
		t.Error("Delete() should return true on success")
	}

	if p.Exists() {
		t.Error("File should not exist after Delete()")
	}
}

func TestPage_Delete_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	p := &page{name: "nonexistent", ext: ".html"}
	if !p.Delete() {
		t.Error("Delete() should return true even for non-existent file")
	}
}

func TestPage_AST(t *testing.T) {
	p := &page{name: "test", ext: ".html"}
	data, node := p.AST()

	if len(data) != 0 {
		t.Errorf("Expected empty byte slice, got %v", data)
	}

	if node == nil {
		t.Error("Expected non-nil AST node")
	}
}

func TestSupportedExtensions(t *testing.T) {
	expected := []string{".htm", ".html", ".xhtml"}
	if len(SUPPORTED_EXT) != len(expected) {
		t.Errorf("Expected %d supported extensions, got %d", len(expected), len(SUPPORTED_EXT))
	}

	for i, ext := range expected {
		if SUPPORTED_EXT[i] != ext {
			t.Errorf("Expected extension '%s' at index %d, got '%s'", ext, i, SUPPORTED_EXT[i])
		}
	}
}
