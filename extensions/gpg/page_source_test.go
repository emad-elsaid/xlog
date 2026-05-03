package gpg

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/emad-elsaid/xlog"
)

// TestEncryptedPages_Page_NoGPGId verifies Page returns nil when gpgId is empty.
func TestEncryptedPages_Page_NoGPGId(t *testing.T) {
	origGpgId := gpgId
	defer func() { gpgId = origGpgId }()

	gpgId = ""

	ps := &encryptedPages{}
	p := ps.Page("test")

	if p != nil {
		t.Errorf("Page() = %v, want nil when gpgId is empty", p)
	}
}

// TestEncryptedPages_Page_NonExistent verifies Page returns nil for non-existent pages.
func TestEncryptedPages_Page_NonExistent(t *testing.T) {
	origGpgId := gpgId
	defer func() { gpgId = origGpgId }()

	gpgId = testGPGKeyID

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(origWd); err != nil {
			t.Errorf("failed to restore directory: %v", err)
		}
	}()

	ps := &encryptedPages{}
	p := ps.Page("nonexistent/page")

	if p != nil {
		t.Errorf("Page() = %v, want nil for non-existent page", p)
	}
}

// TestEncryptedPages_Page_Exists verifies Page returns valid page for existing encrypted files.
func TestEncryptedPages_Page_Exists(t *testing.T) {
	origGpgId := gpgId
	defer func() { gpgId = origGpgId }()

	gpgId = testGPGKeyID

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(origWd); err != nil {
			t.Errorf("failed to restore directory: %v", err)
		}
	}()

	// Create encrypted page file
	pageName := "test/encrypted"
	pg := &page{name: pageName}
	if err := os.MkdirAll(filepath.Dir(pg.FileName()), 0700); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(pg.FileName(), []byte("encrypted content"), 0600); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	ps := &encryptedPages{}
	result := ps.Page(pageName)

	if result == nil {
		t.Fatal("Page() = nil, want non-nil for existing encrypted page")
	}

	if result.Name() != pageName {
		t.Errorf("Page().Name() = %q, want %q", result.Name(), pageName)
	}
}

// TestEncryptedPages_Each_NoGPGId verifies Each does nothing when gpgId is empty.
func TestEncryptedPages_Each_NoGPGId(t *testing.T) {
	origGpgId := gpgId
	defer func() { gpgId = origGpgId }()

	gpgId = ""

	ps := &encryptedPages{}
	called := false

	ps.Each(context.Background(), func(p xlog.Page) {
		called = true
	})

	if called {
		t.Error("Each() called function when gpgId is empty")
	}
}

// TestEncryptedPages_Each_FindsEncryptedPages verifies Each iterates over encrypted pages.
func TestEncryptedPages_Each_FindsEncryptedPages(t *testing.T) {
	origGpgId := gpgId
	defer func() { gpgId = origGpgId }()

	gpgId = testGPGKeyID

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(origWd); err != nil {
			t.Errorf("failed to restore directory: %v", err)
		}
	}()

	// Create test structure
	testPages := []string{
		"page1.md.pgp",
		"dir/page2.md.pgp",
		"dir/subdir/page3.md.pgp",
	}

	for _, pagePath := range testPages {
		if err := os.MkdirAll(filepath.Dir(pagePath), 0700); err != nil {
			t.Fatalf("failed to create directory for %s: %v", pagePath, err)
		}
		if err := os.WriteFile(pagePath, []byte("encrypted"), 0600); err != nil {
			t.Fatalf("failed to write file %s: %v", pagePath, err)
		}
	}

	// Create non-encrypted files that should be ignored
	if err := os.WriteFile("normal.md", []byte("normal"), 0600); err != nil {
		t.Fatalf("failed to write normal.md: %v", err)
	}
	if err := os.WriteFile("other.txt", []byte("text"), 0600); err != nil {
		t.Fatalf("failed to write other.txt: %v", err)
	}

	ps := &encryptedPages{}
	foundPages := make(map[string]bool)

	ps.Each(context.Background(), func(p xlog.Page) {
		foundPages[p.Name()] = true
	})

	expectedPages := map[string]bool{
		"page1":            true,
		"dir/page2":        true,
		"dir/subdir/page3": true,
	}

	if len(foundPages) != len(expectedPages) {
		t.Errorf("Each() found %d pages, want %d", len(foundPages), len(expectedPages))
	}

	for name := range expectedPages {
		if !foundPages[name] {
			t.Errorf("Each() did not find page %q", name)
		}
	}
}

// TestEncryptedPages_Each_ContextCancellation verifies Each respects context cancellation.
func TestEncryptedPages_Each_ContextCancellation(t *testing.T) {
	origGpgId := gpgId
	defer func() { gpgId = origGpgId }()

	gpgId = testGPGKeyID

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(origWd); err != nil {
			t.Errorf("failed to restore directory: %v", err)
		}
	}()

	// Create many encrypted pages
	for i := 0; i < 100; i++ {
		pagePath := filepath.Join("dir", string(rune('a'+i))+".md.pgp")
		if err := os.MkdirAll(filepath.Dir(pagePath), 0700); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}
		if err := os.WriteFile(pagePath, []byte("encrypted"), 0600); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	ps := &encryptedPages{}
	callCount := 0

	// Cancel after first call
	ps.Each(ctx, func(p xlog.Page) {
		callCount++
		if callCount == 1 {
			cancel()
		}
	})

	// Should have stopped after cancellation
	// Note: exact count depends on filesystem iteration order, but should be small
	if callCount > 50 {
		t.Errorf("Each() processed %d pages after cancellation, expected early termination", callCount)
	}
}

// TestEncryptedPages_Each_IgnoresNonEncryptedFiles verifies Each only processes .md.pgp files.
func TestEncryptedPages_Each_IgnoresNonEncryptedFiles(t *testing.T) {
	origGpgId := gpgId
	defer func() { gpgId = origGpgId }()

	gpgId = testGPGKeyID

	tmpDir := t.TempDir()
	origWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(origWd); err != nil {
			t.Errorf("failed to restore directory: %v", err)
		}
	}()

	// Create various file types
	filesToCreate := map[string]bool{
		"valid.md.pgp":      true,  // should be found
		"normal.md":         false, // should be ignored
		"doc.pgp":           false, // should be ignored (not .md.pgp)
		"text.txt":          false, // should be ignored
		"dir/valid2.md.pgp": true,  // should be found
	}

	for path := range filesToCreate {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatalf("failed to create directory for %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte("content"), 0600); err != nil {
			t.Fatalf("failed to write file %s: %v", path, err)
		}
	}

	ps := &encryptedPages{}
	foundPages := make(map[string]bool)

	ps.Each(context.Background(), func(p xlog.Page) {
		foundPages[p.Name()] = true
	})

	// Only "valid" and "dir/valid2" should be found
	expectedCount := 2
	if len(foundPages) != expectedCount {
		t.Errorf("Each() found %d pages, want %d", len(foundPages), expectedCount)
	}

	expectedNames := map[string]bool{
		"valid":      true,
		"dir/valid2": true,
	}

	for name := range foundPages {
		if !expectedNames[name] {
			t.Errorf("Each() found unexpected page %q", name)
		}
	}
}

// TestPGPExtension_Name verifies the extension name.
func TestPGPExtension_Name(t *testing.T) {
	ext := PGP{}
	if name := ext.Name(); name != "pgp" {
		t.Errorf("Name() = %q, want %q", name, "pgp")
	}
}
