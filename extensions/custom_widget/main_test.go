package custom_widget

import (
	"html/template"
	"os"
	"path/filepath"
	"testing"
)

func TestCustomWidgetExtensionName(t *testing.T) {
	ext := CustomWidget{}
	expected := "custom-widget"
	if ext.Name() != expected {
		t.Errorf("Expected extension name to be %q, got %q", expected, ext.Name())
	}
}

func TestCustomWidgetInit_NoFiles(t *testing.T) {
	// Reset flags for test
	oldHead := head_file
	oldBefore := before_view_file
	oldAfter := after_view_file
	defer func() {
		head_file = oldHead
		before_view_file = oldBefore
		after_view_file = oldAfter
	}()

	head_file = ""
	before_view_file = ""
	after_view_file = ""

	// Initialize extension (should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Init() panicked with no files: %v", r)
		}
	}()

	ext := CustomWidget{}
	ext.Init()
}

func TestCustomWidgetInit_WithHeadFile(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	headFile := filepath.Join(tmpDir, "head.html")
	headContent := "<meta name=\"test\" content=\"value\">"
	if err := os.WriteFile(headFile, []byte(headContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Save old values and restore after test
	oldHead := head_file
	defer func() {
		head_file = oldHead
	}()

	head_file = headFile
	before_view_file = ""
	after_view_file = ""

	// Initialize extension
	ext := CustomWidget{}
	ext.Init()

	// Test that readFile works correctly
	output := readFile(headFile)
	if output != template.HTML(headContent) {
		t.Errorf("Expected readFile output to be %q, got %q", headContent, output)
	}
}

func TestCustomWidgetInit_WithBeforeViewFile(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	beforeFile := filepath.Join(tmpDir, "before.html")
	beforeContent := "<div class=\"before-content\">Before</div>"
	if err := os.WriteFile(beforeFile, []byte(beforeContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Save old values and restore after test
	oldBefore := before_view_file
	defer func() {
		before_view_file = oldBefore
	}()

	head_file = ""
	before_view_file = beforeFile
	after_view_file = ""

	// Initialize extension
	ext := CustomWidget{}
	ext.Init()

	// Test that readFile works correctly
	output := readFile(beforeFile)
	if output != template.HTML(beforeContent) {
		t.Errorf("Expected readFile output to be %q, got %q", beforeContent, output)
	}
}

func TestCustomWidgetInit_WithAfterViewFile(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	afterFile := filepath.Join(tmpDir, "after.html")
	afterContent := "<div class=\"after-content\">After</div>"
	if err := os.WriteFile(afterFile, []byte(afterContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Save old values and restore after test
	oldAfter := after_view_file
	defer func() {
		after_view_file = oldAfter
	}()

	head_file = ""
	before_view_file = ""
	after_view_file = afterFile

	// Initialize extension
	ext := CustomWidget{}
	ext.Init()

	// Test that readFile works correctly
	output := readFile(afterFile)
	if output != template.HTML(afterContent) {
		t.Errorf("Expected readFile output to be %q, got %q", afterContent, output)
	}
}

func TestCustomWidgetInit_WithMultipleFiles(t *testing.T) {
	// Create temporary files
	tmpDir := t.TempDir()
	headFile := filepath.Join(tmpDir, "head.html")
	beforeFile := filepath.Join(tmpDir, "before.html")
	afterFile := filepath.Join(tmpDir, "after.html")

	headContent := "<meta name=\"test\" content=\"value\">"
	beforeContent := "<div>Before</div>"
	afterContent := "<div>After</div>"

	if err := os.WriteFile(headFile, []byte(headContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(beforeFile, []byte(beforeContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(afterFile, []byte(afterContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Save old values and restore after test
	oldHead := head_file
	oldBefore := before_view_file
	oldAfter := after_view_file
	defer func() {
		head_file = oldHead
		before_view_file = oldBefore
		after_view_file = oldAfter
	}()

	head_file = headFile
	before_view_file = beforeFile
	after_view_file = afterFile

	// Initialize extension (should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Init() panicked with multiple files: %v", r)
		}
	}()

	ext := CustomWidget{}
	ext.Init()

	// Verify all files can be read
	headOutput := readFile(headFile)
	beforeOutput := readFile(beforeFile)
	afterOutput := readFile(afterFile)

	if headOutput != template.HTML(headContent) {
		t.Errorf("Expected head output %q, got %q", headContent, headOutput)
	}
	if beforeOutput != template.HTML(beforeContent) {
		t.Errorf("Expected before output %q, got %q", beforeContent, beforeOutput)
	}
	if afterOutput != template.HTML(afterContent) {
		t.Errorf("Expected after output %q, got %q", afterContent, afterOutput)
	}
}

func TestReadFile_ValidFile(t *testing.T) {
	// Create temporary file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.html")
	content := "<div>Test Content</div>"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Read file
	result := readFile(testFile)
	if result != template.HTML(content) {
		t.Errorf("Expected %q, got %q", content, result)
	}
}

func TestReadFile_InvalidFile(t *testing.T) {
	// Test with non-existent file
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected readFile to panic with non-existent file")
		}
	}()

	readFile("/nonexistent/file.html")
}

func TestReadFile_EmptyFile(t *testing.T) {
	// Create temporary empty file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "empty.html")
	if err := os.WriteFile(testFile, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	// Read file
	result := readFile(testFile)
	if result != template.HTML("") {
		t.Errorf("Expected empty string, got %q", result)
	}
}
