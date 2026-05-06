package xlog

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

// TestMapPageDeadlock reproduces the deadlock bug when MapPage is called
// with more pages than the concurrency limit.
func TestMapPageDeadlock(t *testing.T) {
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

	// Create MORE pages than concurrency limit (which is NumCPU * 4)
	// With 16 CPUs, concurrency = 64, so let's create 258 pages like the real repo
	numPages := 258
	t.Logf("Creating %d pages (concurrency limit is %d)", numPages, concurrency)

	for i := 0; i < numPages; i++ {
		filename := fmt.Sprintf("page%d.md", i)
		content := fmt.Sprintf("# Page %d\nContent\n", i)
		if err := os.WriteFile(filename, []byte(content), 0600); err != nil {
			t.Fatalf("Failed to create %s: %v", filename, err)
		}
	}

	// Register page source
	RegisterPageSource(newMarkdownFS("."))
	_ = clearPagesCache(nil)

	// This should complete, not deadlock
	done := make(chan bool)
	timeout := time.After(10 * time.Second)

	go func() {
		t.Logf("Starting MapPage with %d pages...", numPages)
		results := MapPage(context.Background(), func(p Page) Page {
			// Return all pages (no filtering)
			return p
		})
		t.Logf("MapPage completed, got %d results", len(results))
		done <- true
	}()

	select {
	case <-done:
		t.Logf("SUCCESS: MapPage completed without deadlock")
	case <-timeout:
		t.Fatal("DEADLOCK BUG: MapPage() hung for 10 seconds with 258 pages!")
	}
}
