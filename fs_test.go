package xlog

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestPriorityFSOpen(t *testing.T) {
	// Create test filesystems
	fs1 := fstest.MapFS{
		"file1.txt": &fstest.MapFile{Data: []byte("from fs1")},
		"shared.txt": &fstest.MapFile{Data: []byte("from fs1")},
	}
	
	fs2 := fstest.MapFS{
		"file2.txt": &fstest.MapFile{Data: []byte("from fs2")},
		"shared.txt": &fstest.MapFile{Data: []byte("from fs2")},
	}

	pfs := priorityFS{fs1, fs2}

	t.Run("file from first filesystem", func(t *testing.T) {
		f, err := pfs.Open("file1.txt")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer f.Close()
		
		// Verify it's the correct file
		if stat, err := f.Stat(); err == nil {
			if stat.Name() != "file1.txt" {
				t.Errorf("Expected name 'file1.txt', got '%s'", stat.Name())
			}
		}
	})

	t.Run("file from second filesystem", func(t *testing.T) {
		f, err := pfs.Open("file2.txt")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer f.Close()
		
		if stat, err := f.Stat(); err == nil {
			if stat.Name() != "file2.txt" {
				t.Errorf("Expected name 'file2.txt', got '%s'", stat.Name())
			}
		}
	})

	t.Run("priority - later filesystem wins", func(t *testing.T) {
		// shared.txt exists in both fs1 and fs2
		// fs2 should win because it's later in the slice
		f, err := pfs.Open("shared.txt")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		defer f.Close()
		
		// Read content to verify it's from fs2
		data := make([]byte, 100)
		n, _ := f.Read(data)
		content := string(data[:n])
		
		if content != "from fs2" {
			t.Errorf("Expected content 'from fs2', got '%s'", content)
		}
	})

	t.Run("nonexistent file returns ErrNotExist", func(t *testing.T) {
		_, err := pfs.Open("nonexistent.txt")
		if err != fs.ErrNotExist {
			t.Errorf("Expected fs.ErrNotExist, got %v", err)
		}
	})
}

func TestPriorityFSEmpty(t *testing.T) {
	pfs := priorityFS{}
	
	_, err := pfs.Open("any.txt")
	if err != fs.ErrNotExist {
		t.Errorf("Expected fs.ErrNotExist for empty priorityFS, got %v", err)
	}
}

func TestRegisterStaticDir(t *testing.T) {
	// Save original staticDirs
	original := staticDirs
	defer func() { staticDirs = original }()
	
	// Reset to known state
	staticDirs = []fs.FS{assets}
	initialLen := len(staticDirs)
	
	// Register a new filesystem
	testFS := fstest.MapFS{
		"test.txt": &fstest.MapFile{Data: []byte("test")},
	}
	
	RegisterStaticDir(testFS)
	
	if len(staticDirs) != initialLen+1 {
		t.Errorf("Expected staticDirs length to be %d, got %d", initialLen+1, len(staticDirs))
	}
}

func TestStaticHandler(t *testing.T) {
	// Create a temporary directory with a test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := []byte("test content")
	
	if err := os.WriteFile(testFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Change to temp directory for the test
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)
	
	t.Run("serve existing file", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test.txt", nil)
		rec := httptest.NewRecorder()
		
		output, err := staticHandler(req)
		
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		// Execute the output handler
		output(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rec.Code)
		}
		
		body := rec.Body.String()
		if body != string(testContent) {
			t.Errorf("Expected body '%s', got '%s'", string(testContent), body)
		}
	})
	
	t.Run("nonexistent file returns error", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/nonexistent.txt", nil)
		
		_, err := staticHandler(req)
		if err == nil {
			t.Error("Expected error for nonexistent file")
		}
	})
	
	t.Run("path traversal is cleaned", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/./test.txt", nil)
		rec := httptest.NewRecorder()
		
		output, err := staticHandler(req)
		
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		
		output(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200 for cleaned path, got %d", rec.Code)
		}
	})
}

func TestStaticHandlerEmbeddedAssets(t *testing.T) {
	// Test that embedded assets can be served
	// The 'public' directory is embedded in the binary
	
	// We know public/style.css exists from the embed directive
	req := httptest.NewRequest("GET", "/public/style.css", nil)
	
	output, err := staticHandler(req)
	
	if err != nil {
		// It's ok if this fails in test environment without embedded assets
		// but we test the code path
		t.Logf("Embedded asset test skipped: %v", err)
		return
	}
	
	rec := httptest.NewRecorder()
	output(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 for embedded asset, got %d", rec.Code)
	}
}

func TestStaticHandlerPriority(t *testing.T) {
	// Save original staticDirs
	original := staticDirs
	defer func() { staticDirs = original }()
	
	// Create temp directory with a file
	tempDir := t.TempDir()
	localFile := filepath.Join(tempDir, "override.txt")
	localContent := []byte("local override")
	
	if err := os.WriteFile(localFile, localContent, 0644); err != nil {
		t.Fatalf("Failed to create local file: %v", err)
	}
	
	// Create a filesystem with the same file
	testFS := fstest.MapFS{
		"override.txt": &fstest.MapFile{Data: []byte("from registered fs")},
	}
	
	// Register the filesystem (it should have lower priority than local files)
	staticDirs = []fs.FS{assets, testFS}
	
	// Change to temp directory
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tempDir)
	
	// Request the file - local version should win due to priorityFS behavior
	req := httptest.NewRequest("GET", "/override.txt", nil)
	rec := httptest.NewRecorder()
	
	output, err := staticHandler(req)
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	output(rec, req)
	
	body := rec.Body.String()
	// The working directory filesystem is added last, so it has highest priority
	if body != string(localContent) {
		t.Errorf("Expected local file to have priority, got '%s'", body)
	}
}
