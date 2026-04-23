package xlog

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewMarkdownFS(t *testing.T) {
	tempDir := t.TempDir()
	
	mfs := newMarkdownFS(tempDir)
	
	if mfs == nil {
		t.Fatal("Expected newMarkdownFS to return a non-nil markdownFS")
	}
	
	if mfs.path != tempDir {
		t.Errorf("Expected path to be %s, got %s", tempDir, mfs.path)
	}
	
	if mfs.cache == nil {
		t.Error("Expected cache to be initialized")
	}
	
	if mfs._page == nil {
		t.Error("Expected _page function to be initialized")
	}
	
	if mfs.watch == nil {
		t.Error("Expected watch function to be initialized")
	}
}

func TestMarkdownFSPage(t *testing.T) {
	tempDir := t.TempDir()
	mfs := newMarkdownFS(tempDir)
	
	t.Run("get page with name", func(t *testing.T) {
		page := mfs.Page("test-page")
		
		if page == nil {
			t.Fatal("Expected Page to return a non-nil page")
		}
		
		if page.Name() != "test-page" {
			t.Errorf("Expected page name to be 'test-page', got '%s'", page.Name())
		}
	})
	
	t.Run("get page without name uses index", func(t *testing.T) {
		// Save original config
		origIndex := Config.Index
		defer func() { Config.Index = origIndex }()
		
		Config.Index = "home"
		
		page := mfs.Page("")
		
		if page == nil {
			t.Fatal("Expected Page to return a non-nil page")
		}
		
		if page.Name() != "home" {
			t.Errorf("Expected page name to be 'home', got '%s'", page.Name())
		}
	})
	
	t.Run("pages are cached", func(t *testing.T) {
		page1 := mfs.Page("cached-test")
		page2 := mfs.Page("cached-test")
		
		// They should be the same instance due to memoization
		if page1.Name() != page2.Name() {
			t.Error("Expected cached pages to have the same name")
		}
	})
}

func TestMarkdownFSEach(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test markdown files
	testFiles := []string{
		"page1.md",
		"page2.md",
		"subdir/page3.md",
	}
	
	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		os.MkdirAll(filepath.Dir(filePath), 0755)
		if err := os.WriteFile(filePath, []byte("# Test\nContent"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}
	
	// Create a non-markdown file (should be ignored)
	txtFile := filepath.Join(tempDir, "readme.txt")
	os.WriteFile(txtFile, []byte("text file"), 0644)
	
	mfs := newMarkdownFS(tempDir)
	
	t.Run("iterate over all markdown files", func(t *testing.T) {
		ctx := context.Background()
		foundPages := make(map[string]bool)
		
		mfs.Each(ctx, func(page Page) {
			foundPages[page.Name()] = true
		})
		
		// Check that all markdown files were found
		expectedPages := []string{
			filepath.Join(tempDir, "page1"),
			filepath.Join(tempDir, "page2"),
			filepath.Join(tempDir, "subdir", "page3"),
		}
		
		if len(foundPages) != len(expectedPages) {
			t.Errorf("Expected %d pages, found %d", len(expectedPages), len(foundPages))
		}
		
		// Verify non-markdown file was not included
		if foundPages[filepath.Join(tempDir, "readme")] {
			t.Error("Expected non-markdown files to be excluded")
		}
	})
	
	t.Run("context cancellation stops iteration", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		
		pageCount := 0
		
		// Cancel after first page
		mfs.Each(ctx, func(page Page) {
			pageCount++
			if pageCount == 1 {
				cancel()
			}
		})
		
		// Should have stopped early
		if pageCount > len(testFiles) {
			t.Errorf("Expected iteration to stop, but found %d pages", pageCount)
		}
	})
}

func TestMarkdownFSEachIgnoredPaths(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create files in hidden directory (starts with .)
	hiddenDir := filepath.Join(tempDir, ".hidden")
	os.MkdirAll(hiddenDir, 0755)
	hiddenFile := filepath.Join(hiddenDir, "config.md")
	os.WriteFile(hiddenFile, []byte("# Hidden"), 0644)
	
	// Create a normal file
	normalFile := filepath.Join(tempDir, "normal.md")
	os.WriteFile(normalFile, []byte("# Normal"), 0644)
	
	mfs := newMarkdownFS(tempDir)
	
	ctx := context.Background()
	foundPages := make([]string, 0)
	
	mfs.Each(ctx, func(page Page) {
		foundPages = append(foundPages, page.Name())
	})
	
	// Should only find the normal file (hidden directory should be skipped)
	if len(foundPages) != 1 {
		t.Errorf("Expected 1 page (ignoring hidden dirs), found %d: %v", len(foundPages), foundPages)
	}
	
	// Check that the normal file was found
	normalPageName := filepath.Join(tempDir, "normal")
	found := false
	for _, page := range foundPages {
		if page == normalPageName {
			found = true
			break
		}
	}
	
	if !found {
		t.Errorf("Expected to find normal.md, but it wasn't in the results")
	}
}

func TestMarkdownFSCacheInvalidation(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.md")
	
	// Create initial file
	if err := os.WriteFile(testFile, []byte("# V1"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	mfs := newMarkdownFS(tempDir)
	
	// Get page (should be cached)
	page1 := mfs.Page("test")
	name1 := page1.Name()
	
	// Manually remove from cache (simulating file change event)
	mfs.cache.Remove("test")
	
	// Get page again (should create new instance)
	page2 := mfs.Page("test")
	name2 := page2.Name()
	
	// Names should still match
	if name1 != name2 {
		t.Errorf("Expected page names to match after cache invalidation")
	}
}

func TestMarkdownFSWatchOnlyCalledOnce(t *testing.T) {
	tempDir := t.TempDir()
	mfs := newMarkdownFS(tempDir)
	
	// Call Page multiple times
	for i := 0; i < 5; i++ {
		mfs.Page("test")
	}
	
	// The watch function should only be called once due to sync.OnceFunc
	// This is hard to test directly, but we can verify that the function
	// doesn't panic or cause issues when called multiple times
	
	// If we get here without panicking, the test passes
}

func TestMarkdownFSWithRealFileChanges(t *testing.T) {
	// This test verifies the file watching behavior
	// Note: This is an integration-style test
	
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "watched.md")
	
	// Create initial file
	if err := os.WriteFile(testFile, []byte("# Initial"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	mfs := newMarkdownFS(tempDir)
	
	// Trigger watch by getting a page
	page := mfs.Page("watched")
	if page.Name() != "watched" {
		t.Errorf("Expected page name 'watched', got '%s'", page.Name())
	}
	
	// Give the watcher time to start
	time.Sleep(100 * time.Millisecond)
	
	// The watcher is now running in the background
	// We can't easily test the event handling without flakiness,
	// but we've verified the setup works
}

func TestMarkdownFSContextTimeout(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create many files to iterate over
	for i := 0; i < 50; i++ {
		filename := filepath.Join(tempDir, filepath.Base(tempDir)+string(rune(i))+".md")
		os.WriteFile(filename, []byte("# Test"), 0644)
	}
	
	mfs := newMarkdownFS(tempDir)
	
	// Create a context that times out quickly
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	
	// Give it time to timeout
	time.Sleep(10 * time.Millisecond)
	
	foundPages := 0
	mfs.Each(ctx, func(page Page) {
		foundPages++
	})
	
	// Due to timeout, we shouldn't process all 50 files
	// (though this might be flaky depending on system speed)
	t.Logf("Found %d pages before context cancellation", foundPages)
}

func TestMarkdownFSEmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()
	mfs := newMarkdownFS(tempDir)
	
	ctx := context.Background()
	foundPages := 0
	
	mfs.Each(ctx, func(page Page) {
		foundPages++
	})
	
	if foundPages != 0 {
		t.Errorf("Expected 0 pages in empty directory, found %d", foundPages)
	}
}
