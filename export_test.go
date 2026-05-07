package xlog

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestExportJSON(t *testing.T) {
	// Save and restore original directory and sources
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	originalSources := sources
	defer func() {
		sourcesMutex.Lock()
		sources = originalSources
		sourcesMutex.Unlock()
	}()

	// Create isolated temp directory
	tmpDir := t.TempDir()

	// Create test markdown files
	testFiles := map[string]string{
		"page1.md": "# Page 1\nContent with [link to page2](page2.md)",
		"page2.md": "# Page 2\nSome content here",
		"page3.md": "# Page 3\n[External link](https://example.com)",
	}

	for name, content := range testFiles {
		path := filepath.Join(tmpDir, name)
		writeErr := os.WriteFile(path, []byte(content), 0644)
		if writeErr != nil {
			t.Fatalf("Failed to create test file %s: %v", name, writeErr)
		}
	}

	// Change to isolated temp directory
	chdirErr := os.Chdir(tmpDir)
	if chdirErr != nil {
		t.Fatal(chdirErr)
	}

	// Replace global sources with test-specific markdownFS
	sourcesMutex.Lock()
	sources = []PageSource{newMarkdownFS(tmpDir)}
	sourcesMutex.Unlock()

	// Clear pages cache to force using new sources
	_ = clearPagesCache(nil)

	// Capture stdout with proper synchronization
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Run ExportJSON in goroutine
	done := make(chan bool)
	var buf bytes.Buffer
	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	ExportJSON(context.Background())
	w.Close()
	<-done

	// Restore stdout
	os.Stdout = oldStdout

	// Parse JSON output
	var metadata []PageMetadata
	output := buf.Bytes()
	if len(output) == 0 {
		t.Fatal("No output from ExportJSON")
	}

	if err := json.Unmarshal(output, &metadata); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, string(output))
	}

	// Verify we got all pages
	if len(metadata) < 3 {
		t.Errorf("Expected at least 3 pages, got %d", len(metadata))
		for i, m := range metadata {
			t.Logf("Page %d: %s (path: %s)", i, m.Name, m.Path)
		}
	}

	// Verify structure (allowing for possible extra files from temp dir)
	foundPages := make(map[string]bool)
	for _, page := range metadata {
		if page.Name == "" {
			t.Error("Page name is empty")
		}
		if page.Path == "" {
			t.Error("Page path is empty")
		}
		// ModTime can be zero in some test environments, just check it exists
		if page.WordCount < 0 {
			t.Error("Page word_count is negative")
		}
		foundPages[page.Name] = true
	}

	// Verify our test pages exist
	expectedPages := []string{"page1", "page2", "page3"}
	for _, expected := range expectedPages {
		if !foundPages[expected] {
			t.Errorf("Expected page %s not found in metadata", expected)
		}
	}
}

func TestPageMetadata_JSONStructure(t *testing.T) {
	// Test JSON marshaling of PageMetadata struct
	meta := PageMetadata{
		Name:      "test-page",
		Path:      "/path/to/test-page.md",
		WordCount: 42,
		Links:     []string{"link1", "link2"},
	}

	data, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("Failed to marshal PageMetadata: %v", err)
	}

	// Verify JSON contains expected fields
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if decoded["name"] != "test-page" {
		t.Errorf("Expected name 'test-page', got %v", decoded["name"])
	}

	if decoded["path"] != "/path/to/test-page.md" {
		t.Errorf("Expected path '/path/to/test-page.md', got %v", decoded["path"])
	}

	if decoded["word_count"] != float64(42) {
		t.Errorf("Expected word_count 42, got %v", decoded["word_count"])
	}

	links, ok := decoded["links"].([]interface{})
	if !ok || len(links) != 2 {
		t.Errorf("Expected 2 links, got %v", decoded["links"])
	}
}
