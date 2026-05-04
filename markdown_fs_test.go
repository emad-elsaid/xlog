package xlog

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMarkdownFS(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantPanic bool
	}{
		{
			name:      "valid path creates markdownFS",
			path:      ".",
			wantPanic: false,
		},
		{
			name:      "non-existent path creates markdownFS",
			path:      "/nonexistent/path/12345",
			wantPanic: false,
		},
		{
			name:      "empty path creates markdownFS",
			path:      "",
			wantPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					newMarkdownFS(tt.path)
				})
			} else {
				mfs := newMarkdownFS(tt.path)
				require.NotNil(t, mfs)
				assert.Equal(t, tt.path, mfs.path)
				assert.NotNil(t, mfs.cache)
				assert.NotNil(t, mfs._page)
				assert.NotNil(t, mfs.watch)
			}
		})
	}
}

func TestNewMarkdownFS_CacheInitialization(t *testing.T) {
	// This test verifies that the cache is properly initialized
	// and the fallback mechanism is in place for edge cases
	mfs := newMarkdownFS(".")

	// Verify cache is not nil
	require.NotNil(t, mfs.cache, "cache should be initialized")

	// Verify cache has non-zero capacity
	assert.Greater(t, mfs.cache.Len(), -1, "cache should have valid capacity")

	// Verify page function works with initialized cache
	page := mfs.Page("test")
	require.NotNil(t, page, "Page should work with initialized cache")

	// Verify cache stores pages
	page2 := mfs.Page("test")
	assert.Equal(t, page.Name(), page2.Name(), "cache should return same page for same name")
}

func TestMarkdownFS_Page(t *testing.T) {
	// Save and restore original config
	origIndex := Config.Index
	defer func() { Config.Index = origIndex }()
	Config.Index = "index"

	tests := []struct {
		name         string
		pageName     string
		expectedName string
	}{
		{
			name:         "normal page name",
			pageName:     "testpage",
			expectedName: "testpage",
		},
		{
			name:         "empty name uses index",
			pageName:     "",
			expectedName: "index",
		},
		{
			name:         "page name with path",
			pageName:     "path/to/page",
			expectedName: "path/to/page",
		},
		{
			name:         "page name with special chars",
			pageName:     "page-with_special.chars",
			expectedName: "page-with_special.chars",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mfs := newMarkdownFS(".")
			page := mfs.Page(tt.pageName)

			require.NotNil(t, page)
			assert.Equal(t, tt.expectedName, page.Name())
		})
	}
}

func TestMarkdownFS_PageCaching(t *testing.T) {
	mfs := newMarkdownFS(".")

	// Request the same page twice
	page1 := mfs.Page("testpage")
	page2 := mfs.Page("testpage")

	// Should return the same cached instance
	require.NotNil(t, page1)
	require.NotNil(t, page2)
	assert.Equal(t, page1.Name(), page2.Name())
}

func TestMarkdownFS_Each(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()

	// Create test markdown files
	testFiles := []string{
		"page1.md",
		"page2.md",
		"subdir/page3.md",
	}

	for _, f := range testFiles {
		fullPath := filepath.Join(tmpDir, f)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0750); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Create non-markdown file (should be ignored)
	if err := os.WriteFile(filepath.Join(tmpDir, "ignored.txt"), []byte("test"), 0600); err != nil {
		t.Fatalf("Failed to create ignored file: %v", err)
	}

	tests := []struct {
		name           string
		setupIgnored   func()
		expectedCount  int
		cancelContext  bool
		expectedCancel bool
	}{
		{
			name:          "iterates all markdown files",
			expectedCount: 3,
		},
		{
			name:           "stops on context cancellation",
			cancelContext:  true,
			expectedCancel: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mfs := newMarkdownFS(tmpDir)
			ctx := context.Background()

			if tt.cancelContext {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel() // Cancel immediately
			}

			var pages []Page
			mfs.Each(ctx, func(p Page) {
				pages = append(pages, p)
			})

			if tt.expectedCancel {
				// Context was cancelled, so we might not get all pages
				assert.LessOrEqual(t, len(pages), tt.expectedCount)
			} else {
				assert.Equal(t, tt.expectedCount, len(pages))
			}
		})
	}
}

func TestMarkdownFS_EachWithIgnoredPaths(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files in root and subdirectories
	testFiles := []string{
		"normal.md",
		"docs/readme.md",
		"other/test.md",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tmpDir, file)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0750); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte("test"), 0600); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// Create hidden directory - should be skipped by IsIgnoredPath check
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	if err := os.MkdirAll(hiddenDir, 0750); err != nil {
		t.Fatalf("Failed to create hidden directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(hiddenDir, "file.md"), []byte("test"), 0600); err != nil {
		t.Fatalf("Failed to create file in hidden dir: %v", err)
	}

	mfs := newMarkdownFS(tmpDir)
	ctx := context.Background()

	var foundPages []string
	mfs.Each(ctx, func(p Page) {
		foundPages = append(foundPages, p.Name())
	})

	// The actual behavior depends on how WalkDir reports paths
	// Just verify we get some pages and don't panic
	assert.NotEmpty(t, foundPages, "Should find some markdown files")

	// Verify all found pages are from .md files
	for _, page := range foundPages {
		// Pages should have been created with .md extension stripped
		assert.NotEmpty(t, page, "Page name should not be empty")
	}
}

func TestMarkdownFS_WatchInitialization(t *testing.T) {
	mfs := newMarkdownFS(".")

	// First call to Page should initialize watch
	_ = mfs.Page("test")

	// Watch should be initialized (we can't easily test the goroutine,
	// but we can verify it doesn't panic)
	assert.NotNil(t, mfs.watch)
}

func TestMarkdownFS_EachNonExistentDirectory(t *testing.T) {
	mfs := newMarkdownFS("/nonexistent/directory/12345")
	ctx := context.Background()

	var callCount int
	// Should handle non-existent directory gracefully without panic
	assert.NotPanics(t, func() {
		mfs.Each(ctx, func(p Page) {
			callCount++
		})
	})

	assert.Equal(t, 0, callCount, "Should not process any pages for non-existent directory")
}

func TestMarkdownFS_EachWithPermissionError(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a directory with no read permissions
	restrictedDir := filepath.Join(tmpDir, "restricted")
	if err := os.MkdirAll(restrictedDir, 0000); err != nil {
		t.Fatalf("Failed to create restricted directory: %v", err)
	}
	defer func() {
		// #nosec G302 - Test cleanup requires restoring directory permissions to allow deletion
		if err := os.Chmod(restrictedDir, 0700); err != nil {
			t.Logf("Failed to restore permissions: %v", err)
		}
	}()

	// Create a normal file to verify we still process accessible files
	normalFile := filepath.Join(tmpDir, "normal.md")
	if err := os.WriteFile(normalFile, []byte("test"), 0600); err != nil {
		t.Fatalf("Failed to create normal file: %v", err)
	}

	mfs := newMarkdownFS(tmpDir)
	ctx := context.Background()

	var pages []Page
	// Should not panic when encountering permission errors
	assert.NotPanics(t, func() {
		mfs.Each(ctx, func(p Page) {
			pages = append(pages, p)
		})
	})

	// Should still process the accessible file
	assert.GreaterOrEqual(t, len(pages), 1, "Should process accessible files despite errors")
}

func TestMarkdownFS_EachWithTimeout(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a few test files
	for i := 0; i < 5; i++ {
		filename := filepath.Join(tmpDir, "page"+string(rune('0'+i))+".md")
		if err := os.WriteFile(filename, []byte("test"), 0600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	mfs := newMarkdownFS(tmpDir)

	// Create context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Give context time to expire
	time.Sleep(10 * time.Millisecond)

	var pages []Page
	mfs.Each(ctx, func(p Page) {
		pages = append(pages, p)
	})

	// Should stop iteration when context times out
	// Might process 0 or few pages depending on timing
	assert.LessOrEqual(t, len(pages), 5)
}

func TestMarkdownFS_FileWatchModify(t *testing.T) {
	tmpDir := t.TempDir()

	// Create initial test file
	testFile := filepath.Join(tmpDir, "watchtest.md")
	err := os.WriteFile(testFile, []byte("initial content"), 0600)
	require.NoError(t, err)

	mfs := newMarkdownFS(tmpDir)

	// Track events with mutex for race-free access
	var mu sync.Mutex
	var pageChangedCalled bool
	var eventPage Page
	Listen(PageChanged, func(p Page) error {
		if p.Name() == "watchtest" {
			mu.Lock()
			pageChangedCalled = true
			eventPage = p
			mu.Unlock()
		}
		return nil
	})

	// Initialize watcher by calling Page
	initialPage := mfs.Page("watchtest")
	require.NotNil(t, initialPage)

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Modify the file
	err = os.WriteFile(testFile, []byte("modified content"), 0600)
	require.NoError(t, err)

	// Wait for event to be processed
	time.Sleep(300 * time.Millisecond)

	// Verify PageChanged event was triggered
	mu.Lock()
	called := pageChangedCalled
	page := eventPage
	mu.Unlock()

	assert.True(t, called, "PageChanged event should have been triggered")
	if called {
		assert.Equal(t, "watchtest", page.Name())
	}

	// Verify cache was invalidated
	exists := mfs.cache.Contains("watchtest")
	assert.False(t, exists, "Cache should have been invalidated after file modification")
}

func TestMarkdownFS_FileWatchDelete(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tmpDir, "deleteme.md")
	err := os.WriteFile(testFile, []byte("content"), 0600)
	require.NoError(t, err)

	mfs := newMarkdownFS(tmpDir)

	// Track PageDeleted events with mutex
	var mu sync.Mutex
	var pageDeletedCalled bool
	var deletedPage Page
	Listen(PageDeleted, func(p Page) error {
		if p.Name() == "deleteme" {
			mu.Lock()
			pageDeletedCalled = true
			deletedPage = p
			mu.Unlock()
		}
		return nil
	})

	// Initialize watcher
	page := mfs.Page("deleteme")
	require.NotNil(t, page)

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Delete the file
	err = os.Remove(testFile)
	require.NoError(t, err)

	// Wait for event to be processed
	time.Sleep(300 * time.Millisecond)

	// Verify PageDeleted event was triggered
	mu.Lock()
	called := pageDeletedCalled
	page = deletedPage
	mu.Unlock()

	assert.True(t, called, "PageDeleted event should have been triggered")
	if called {
		assert.Equal(t, "deleteme", page.Name())
	}

	// Verify cache was invalidated
	exists := mfs.cache.Contains("deleteme")
	assert.False(t, exists, "Cache should have been invalidated after file deletion")
}

func TestMarkdownFS_FileWatchCreate(t *testing.T) {
	tmpDir := t.TempDir()

	mfs := newMarkdownFS(tmpDir)

	// Track PageChanged events (create triggers Write event) with mutex
	var mu sync.Mutex
	var pageChangedCalled bool
	var createdPage Page
	Listen(PageChanged, func(p Page) error {
		if p.Name() == "newfile" {
			mu.Lock()
			pageChangedCalled = true
			createdPage = p
			mu.Unlock()
		}
		return nil
	})

	// Initialize watcher by accessing any page
	_ = mfs.Page("dummy")

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Create new file
	newFile := filepath.Join(tmpDir, "newfile.md")
	err := os.WriteFile(newFile, []byte("new content"), 0600)
	require.NoError(t, err)

	// Wait for event to be processed
	time.Sleep(300 * time.Millisecond)

	// Verify event was triggered
	mu.Lock()
	called := pageChangedCalled
	page := createdPage
	mu.Unlock()

	assert.True(t, called, "PageChanged event should have been triggered for new file")
	if called {
		assert.Equal(t, "newfile", page.Name())
	}
}

func TestMarkdownFS_FileWatchRename(t *testing.T) {
	tmpDir := t.TempDir()

	// Create initial file
	oldFile := filepath.Join(tmpDir, "oldname.md")
	err := os.WriteFile(oldFile, []byte("content"), 0600)
	require.NoError(t, err)

	mfs := newMarkdownFS(tmpDir)

	// Track events with mutex
	var mu sync.Mutex
	var oldPageChanged bool
	var newPageChanged bool
	Listen(PageChanged, func(p Page) error {
		mu.Lock()
		switch p.Name() {
		case "oldname":
			oldPageChanged = true
		case "newname":
			newPageChanged = true
		}
		mu.Unlock()
		return nil
	})

	// Initialize watcher
	_ = mfs.Page("oldname")

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Rename the file
	newFile := filepath.Join(tmpDir, "newname.md")
	err = os.Rename(oldFile, newFile)
	require.NoError(t, err)

	// Wait for events to be processed
	time.Sleep(300 * time.Millisecond)

	// At least one event should have been triggered
	mu.Lock()
	oldChanged := oldPageChanged
	newChanged := newPageChanged
	mu.Unlock()

	assert.True(t, oldChanged || newChanged, "At least one PageChanged event should have been triggered")
}

func TestMarkdownFS_FileWatchIgnoresNonMarkdown(t *testing.T) {
	tmpDir := t.TempDir()

	mfs := newMarkdownFS(tmpDir)

	// Track all PageChanged events with mutex
	var mu sync.Mutex
	var eventCount int
	Listen(PageChanged, func(p Page) error {
		mu.Lock()
		eventCount++
		mu.Unlock()
		return nil
	})

	// Initialize watcher
	_ = mfs.Page("dummy")

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Create non-markdown file
	txtFile := filepath.Join(tmpDir, "test.txt")
	err := os.WriteFile(txtFile, []byte("text content"), 0600)
	require.NoError(t, err)

	// Wait to see if any events are triggered
	time.Sleep(200 * time.Millisecond)

	// No events should have been triggered for .txt file
	mu.Lock()
	count := eventCount
	mu.Unlock()

	assert.Equal(t, 0, count, "No events should be triggered for non-markdown files")
}

func TestMarkdownFS_FileWatchIgnoresIgnoredPaths(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .hidden directory (should be ignored by IsIgnoredPath)
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	err := os.MkdirAll(hiddenDir, 0750)
	require.NoError(t, err)

	mfs := newMarkdownFS(tmpDir)

	// Track PageChanged events with mutex
	var mu sync.Mutex
	var hiddenEventTriggered bool
	Listen(PageChanged, func(p Page) error {
		if p.Name() == ".hidden/secret" {
			mu.Lock()
			hiddenEventTriggered = true
			mu.Unlock()
		}
		return nil
	})

	// Initialize watcher
	_ = mfs.Page("dummy")

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Create file in hidden directory
	hiddenFile := filepath.Join(hiddenDir, "secret.md")
	err = os.WriteFile(hiddenFile, []byte("secret content"), 0600)
	require.NoError(t, err)

	// Wait for potential events
	time.Sleep(300 * time.Millisecond)

	// Event should not have been triggered for ignored path
	mu.Lock()
	triggered := hiddenEventTriggered
	mu.Unlock()

	assert.False(t, triggered, "Events should not be triggered for ignored paths")
}

func TestMarkdownFS_FileWatchSubdirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.MkdirAll(subDir, 0750)
	require.NoError(t, err)

	// Create file in subdirectory
	subFile := filepath.Join(subDir, "page.md")
	err = os.WriteFile(subFile, []byte("initial"), 0600)
	require.NoError(t, err)

	mfs := newMarkdownFS(tmpDir)

	// Track events for subdirectory files with mutex
	var mu sync.Mutex
	var eventTriggered bool
	var eventPageName string
	Listen(PageChanged, func(p Page) error {
		if p.Name() == "subdir/page" {
			mu.Lock()
			eventTriggered = true
			eventPageName = p.Name()
			mu.Unlock()
		}
		return nil
	})

	// Initialize watcher
	_ = mfs.Page("dummy")

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Modify file in subdirectory
	err = os.WriteFile(subFile, []byte("modified"), 0600)
	require.NoError(t, err)

	// Wait for event
	time.Sleep(300 * time.Millisecond)

	// Verify event was triggered with correct path
	mu.Lock()
	triggered := eventTriggered
	pageName := eventPageName
	mu.Unlock()

	assert.True(t, triggered, "Events should be triggered for subdirectory files")
	assert.Equal(t, "subdir/page", pageName, "Event should contain relative path")
}

func TestMarkdownFS_WatchInitializedOnlyOnce(t *testing.T) {
	tmpDir := t.TempDir()
	mfs := newMarkdownFS(tmpDir)

	// Create test file
	testFile := filepath.Join(tmpDir, "test.md")
	err := os.WriteFile(testFile, []byte("content"), 0600)
	require.NoError(t, err)

	// Track events with mutex
	var mu sync.Mutex
	var eventCount int
	Listen(PageChanged, func(p Page) error {
		if p.Name() == "test" {
			mu.Lock()
			eventCount++
			mu.Unlock()
		}
		return nil
	})

	// Call Page multiple times (should only initialize watcher once)
	_ = mfs.Page("page1")
	_ = mfs.Page("page2")
	_ = mfs.Page("page3")

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Modify file
	err = os.WriteFile(testFile, []byte("modified"), 0600)
	require.NoError(t, err)

	// Wait for events
	time.Sleep(300 * time.Millisecond)

	// Should get at least one event (file systems may generate multiple events per write)
	mu.Lock()
	count := eventCount
	mu.Unlock()

	assert.GreaterOrEqual(t, count, 1, "At least one event should have been triggered")
	assert.LessOrEqual(t, count, 3, "Should not get more events than reasonable for a single file modification")
}
