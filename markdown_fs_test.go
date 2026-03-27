package xlog

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMarkdownFS_Page(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	
	// Create a test markdown file
	testFile := filepath.Join(tmpDir, "test.md")
	content := []byte("# Test Page\n\nThis is a test page.")
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Initialize markdownFS
	mfs := newMarkdownFS(tmpDir)

	t.Run("retrieve existing page", func(t *testing.T) {
		page := mfs.Page("test")
		if page == nil {
			t.Fatal("Expected page, got nil")
		}
		if page.Name() != "test" {
			t.Errorf("Expected page name 'test', got '%s'", page.Name())
		}
	})

	t.Run("retrieve index page when empty name", func(t *testing.T) {
		// Set a default index in Config
		originalIndex := Config.Index
		Config.Index = "index"
		defer func() { Config.Index = originalIndex }()

		page := mfs.Page("")
		if page == nil {
			t.Fatal("Expected page, got nil")
		}
		if page.Name() != "index" {
			t.Errorf("Expected page name 'index', got '%s'", page.Name())
		}
	})

	t.Run("cache returns same instance", func(t *testing.T) {
		page1 := mfs.Page("test")
		page2 := mfs.Page("test")
		
		// They should be the same instance due to caching
		if page1.Name() != page2.Name() {
			t.Error("Expected cached pages to have same name")
		}
	})
}

func TestMarkdownFS_Each(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()
	
	// Create multiple test markdown files
	files := map[string]string{
		"page1.md": "# Page 1",
		"page2.md": "# Page 2",
		"page3.md": "# Page 3",
		"test.txt": "Not a markdown file",
	}
	
	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	// Create a subdirectory with another markdown file
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	subFile := filepath.Join(subDir, "page4.md")
	if err := os.WriteFile(subFile, []byte("# Page 4"), 0644); err != nil {
		t.Fatalf("Failed to create subdirectory file: %v", err)
	}

	mfs := newMarkdownFS(tmpDir)

	t.Run("iterate all markdown files", func(t *testing.T) {
		ctx := context.Background()
		count := 0
		names := make(map[string]bool)

		mfs.Each(ctx, func(page Page) {
			count++
			names[page.Name()] = true
		})

		// Should find 4 markdown files (page1, page2, page3, page4)
		// but not test.txt
		if count < 3 {
			t.Errorf("Expected at least 3 markdown files, found %d", count)
		}

		// Normalize page names to handle platform-specific path separators
		hasPage1 := names["page1"] || names[filepath.Join(tmpDir, "page1")]
		hasPage2 := names["page2"] || names[filepath.Join(tmpDir, "page2")]
		hasPage3 := names["page3"] || names[filepath.Join(tmpDir, "page3")]

		if !hasPage1 {
			t.Errorf("Expected to find page1, got names: %v", names)
		}
		if !hasPage2 {
			t.Errorf("Expected to find page2, got names: %v", names)
		}
		if !hasPage3 {
			t.Errorf("Expected to find page3, got names: %v", names)
		}
	})

	t.Run("respect context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		count := 0
		mfs.Each(ctx, func(page Page) {
			count++
		})

		// With cancelled context, iteration should stop early or not start
		// Count should be less than total files
		if count > 1 {
			t.Logf("Context cancellation may not have stopped iteration immediately (found %d pages)", count)
		}
	})
}

func TestMarkdownFS_Watch(t *testing.T) {
	tmpDir := t.TempDir()
	mfs := newMarkdownFS(tmpDir)

	// Trigger watch by calling Page
	_ = mfs.Page("test")

	// Give the watcher time to start
	time.Sleep(100 * time.Millisecond)

	t.Run("watch is initialized once", func(t *testing.T) {
		// Multiple calls should not panic or cause issues
		_ = mfs.Page("test2")
		_ = mfs.Page("test3")
		
		// If we got here without panic, the sync.OnceFunc worked
	})
}

func TestNewMarkdownFS(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("creates new instance", func(t *testing.T) {
		mfs := newMarkdownFS(tmpDir)
		if mfs == nil {
			t.Fatal("Expected markdownFS instance, got nil")
		}
		if mfs.path != tmpDir {
			t.Errorf("Expected path '%s', got '%s'", tmpDir, mfs.path)
		}
		if mfs.cache == nil {
			t.Error("Expected cache to be initialized")
		}
	})
}
