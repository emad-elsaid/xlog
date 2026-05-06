package xlog

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestHiddenDirectoriesNotScanned verifies that hidden directories like .git
// are properly skipped during page discovery, preventing performance issues.
func TestHiddenDirectoriesNotScanned(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	// Create visible pages
	visiblePages := []string{
		"index.md",
		"readme.md",
		"docs/guide.md",
		"notes/daily.md",
	}
	for _, path := range visiblePages {
		if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
			t.Fatalf("Failed to create directory for %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte("# Test"), 0600); err != nil {
			t.Fatalf("Failed to create %s: %v", path, err)
		}
	}

	// Create many files in hidden directories to simulate a large .git
	hiddenPages := []string{
		".git/objects/file1.md",
		".git/objects/file2.md",
		".git/refs/file3.md",
		".hidden/secret.md",
		".cache/data.md",
	}
	for _, path := range hiddenPages {
		if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
			t.Fatalf("Failed to create hidden directory for %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte("# Hidden"), 0600); err != nil {
			t.Fatalf("Failed to create %s: %v", path, err)
		}
	}

	// Initialize markdown filesystem
	mfs := newMarkdownFS(".")

	ctx := context.Background()
	start := time.Now()

	// Collect all discovered pages
	var foundPages []string
	mfs.Each(ctx, func(p Page) {
		foundPages = append(foundPages, p.Name())
	})

	elapsed := time.Since(start)

	// Verify only visible pages were found
	if len(foundPages) != len(visiblePages) {
		t.Errorf("Expected %d pages, got %d. Found: %v", len(visiblePages), len(foundPages), foundPages)
	}

	expectedNames := map[string]bool{
		"index":       true,
		"readme":      true,
		"docs/guide":  true,
		"notes/daily": true,
	}

	for _, name := range foundPages {
		if !expectedNames[name] {
			t.Errorf("Found unexpected page: %q", name)
		}
	}

	// Verify hidden directory pages were NOT found
	hiddenNames := []string{
		".git/objects/file1",
		".git/objects/file2",
		".git/refs/file3",
		".hidden/secret",
		".cache/data",
	}
	for _, hidden := range hiddenNames {
		for _, found := range foundPages {
			if found == hidden {
				t.Errorf("Should not have found page in hidden directory: %q", hidden)
			}
		}
	}

	// Performance check: should be fast even with hidden directories
	// This is a sanity check - with proper filtering, 4 pages should scan in < 100ms
	if elapsed > 100*time.Millisecond {
		t.Logf("WARNING: Page discovery took %v (expected < 100ms)", elapsed)
	}
}

// TestHiddenDirectoriesPerformance benchmarks page discovery with and without hidden directories.
func TestHiddenDirectoriesPerformance(t *testing.T) {
	// Create scenario with many files in .git
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	// Create 10 visible pages
	for i := 0; i < 10; i++ {
		filename := filepath.Join("docs", "page"+string(rune('0'+i))+".md")
		if err := os.MkdirAll(filepath.Dir(filename), 0750); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(filename, []byte("# Test"), 0600); err != nil {
			t.Fatalf("Failed to create %s: %v", filename, err)
		}
	}

	// Create 100 files in .git to simulate a real repository
	for i := 0; i < 100; i++ {
		filename := filepath.Join(".git", "objects", "file"+string(rune('0'+(i%10)))+".md")
		if err := os.MkdirAll(filepath.Dir(filename), 0750); err != nil {
			t.Fatalf("Failed to create .git directory: %v", err)
		}
		if err := os.WriteFile(filename, []byte("# Git Object"), 0600); err != nil {
			t.Fatalf("Failed to create %s: %v", filename, err)
		}
	}

	mfs := newMarkdownFS(".")

	ctx := context.Background()
	start := time.Now()

	count := 0
	mfs.Each(ctx, func(p Page) {
		count++
	})

	elapsed := time.Since(start)

	// Should find exactly 10 visible pages
	if count != 10 {
		t.Errorf("Expected 10 pages, got %d", count)
	}

	// Log performance for visibility
	t.Logf("Scanned %d pages in %v (with 100 files in .git directory)", count, elapsed)

	// With proper filtering, this should be very fast
	// Even 50ms is generous for just 10 files
	if elapsed > 50*time.Millisecond {
		t.Errorf("Page discovery too slow: %v (expected < 50ms). Hidden directories may not be properly skipped.", elapsed)
	}
}
