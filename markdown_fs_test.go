package xlog

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMarkdownFS_Page(t *testing.T) {
	tmpDir := t.TempDir()
	fs := newMarkdownFS(tmpDir)

	// Test getting a page
	page := fs.Page("test-page")
	if page == nil {
		t.Fatal("Expected page to be non-nil")
	}

	if page.Name() != "test-page" {
		t.Errorf("Expected page name to be 'test-page', got '%s'", page.Name())
	}
}

func TestMarkdownFS_Page_DefaultIndex(t *testing.T) {
	tmpDir := t.TempDir()
	fs := newMarkdownFS(tmpDir)

	// Test empty name defaults to Config.Index
	originalIndex := Config.Index
	Config.Index = "index"
	defer func() { Config.Index = originalIndex }()

	page := fs.Page("")
	if page == nil {
		t.Fatal("Expected page to be non-nil")
	}

	if page.Name() != "index" {
		t.Errorf("Expected page name to be 'index', got '%s'", page.Name())
	}
}

func TestMarkdownFS_Each(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test markdown files
	testFiles := []string{"page1.md", "page2.md", "subdir/page3.md"}
	for _, file := range testFiles {
		fullPath := filepath.Join(tmpDir, file)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte("# Test"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create a non-markdown file that should be ignored
	if err := os.WriteFile(filepath.Join(tmpDir, "ignore.txt"), []byte("ignore"), 0644); err != nil {
		t.Fatalf("Failed to create ignore file: %v", err)
	}

	fs := newMarkdownFS(tmpDir)

	// Collect pages
	var pages []string
	ctx := context.Background()
	fs.Each(ctx, func(p Page) {
		pages = append(pages, p.Name())
	})

	// Should find 3 markdown files
	if len(pages) != 3 {
		t.Errorf("Expected 3 pages, got %d: %v", len(pages), pages)
	}

	// Verify .txt file was not included
	for _, name := range pages {
		if filepath.Ext(name) == ".txt" {
			t.Errorf("Non-markdown file should not be included: %s", name)
		}
	}
}

func TestMarkdownFS_Each_ContextCancellation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create many test files
	for i := 0; i < 100; i++ {
		file := filepath.Join(tmpDir, filepath.Join("test", string(rune('a'+i%26)), "page.md"))
		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(file, []byte("# Test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	fs := newMarkdownFS(tmpDir)

	// Create context with immediate cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Try to iterate (should stop early due to cancelled context)
	var count int
	fs.Each(ctx, func(p Page) {
		count++
	})

	// Should process very few or zero pages due to immediate cancellation
	if count > 10 {
		t.Errorf("Expected very few pages processed due to cancelled context, got %d", count)
	}
}

func TestMarkdownFS_Cache(t *testing.T) {
	tmpDir := t.TempDir()
	fs := newMarkdownFS(tmpDir)

	// Get the same page twice
	page1 := fs.Page("test-page")
	page2 := fs.Page("test-page")

	// Should return the same cached instance
	if page1 != page2 {
		t.Error("Expected cached page instances to be identical")
	}
}

func TestMarkdownFS_Watch_FileChange(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping file watch test in short mode")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")

	// Create initial file
	if err := os.WriteFile(testFile, []byte("# Original"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	fs := newMarkdownFS(tmpDir)

	// Set up event listener
	eventReceived := make(chan bool, 1)
	Listen(PageChanged, func(p Page) error {
		if p.Name() == "test" {
			eventReceived <- true
		}
		return nil
	})

	// Trigger watch by accessing a page
	_ = fs.Page("test")

	// Give watch goroutine time to start
	time.Sleep(100 * time.Millisecond)

	// Modify the file
	if err := os.WriteFile(testFile, []byte("# Modified"), 0644); err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}

	// Wait for event with timeout
	select {
	case <-eventReceived:
		// Success - event was received
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for PageChanged event")
	}
}

func TestMarkdownFS_Watch_FileDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping file watch test in short mode")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")

	// Create initial file
	if err := os.WriteFile(testFile, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	fs := newMarkdownFS(tmpDir)

	// Set up event listener
	eventReceived := make(chan bool, 1)
	Listen(PageDeleted, func(p Page) error {
		if p.Name() == "test" {
			eventReceived <- true
		}
		return nil
	})

	// Trigger watch by accessing a page
	_ = fs.Page("test")

	// Give watch goroutine time to start
	time.Sleep(100 * time.Millisecond)

	// Delete the file
	if err := os.Remove(testFile); err != nil {
		t.Fatalf("Failed to delete test file: %v", err)
	}

	// Wait for event with timeout
	select {
	case <-eventReceived:
		// Success - event was received
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for PageDeleted event")
	}
}
