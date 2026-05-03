package xlog

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestRegisterBuildPage(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		encloseInDir bool
		checkEnclosed bool
	}{
		{
			name:         "register simple page",
			path:         "/test-page.html",
			encloseInDir: false,
			checkEnclosed: false,
		},
		{
			name:         "register enclosed page",
			path:         "/enclosed-page",
			encloseInDir: true,
			checkEnclosed: true,
		},
		{
			name:         "register root path",
			path:         "/",
			encloseInDir: false,
			checkEnclosed: false,
		},
		{
			name:         "register nested path",
			path:         "/api/v1/endpoint",
			encloseInDir: true,
			checkEnclosed: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clear the maps by deleting the test path if it exists
			extension_page.Delete(tc.path)
			extension_page_enclosed.Delete(tc.path)

			RegisterBuildPage(tc.path, tc.encloseInDir)

			if tc.checkEnclosed {
				if val, exists := extension_page_enclosed.Load(tc.path); !exists || !val {
					t.Errorf("expected path %q to be in extension_page_enclosed", tc.path)
				}
				if _, exists := extension_page.Load(tc.path); exists {
					t.Errorf("path %q should not be in extension_page when enclosed", tc.path)
				}
			} else {
				if val, exists := extension_page.Load(tc.path); !exists || !val {
					t.Errorf("expected path %q to be in extension_page", tc.path)
				}
				if _, exists := extension_page_enclosed.Load(tc.path); exists {
					t.Errorf("path %q should not be in extension_page_enclosed", tc.path)
				}
			}
		})
	}
}

func TestBuildRoute(t *testing.T) {
	// Create temporary directory for test output
	tmpDir, err := os.MkdirTemp("", "xlog-build-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to clean up temp dir: %v", err)
		}
	}()

	tests := []struct {
		name           string
		route          string
		dir            string
		file           string
		expectedStatus int
		shouldFail     bool
	}{
		{
			name:           "build valid route",
			route:          "/",
			dir:            filepath.Join(tmpDir, "test1"),
			file:           filepath.Join(tmpDir, "test1", "index.html"),
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name:           "build creates nested directories",
			route:          "/",
			dir:            filepath.Join(tmpDir, "nested", "deep", "path"),
			file:           filepath.Join(tmpDir, "nested", "deep", "path", "page.html"),
			expectedStatus: http.StatusOK,
			shouldFail:     false,
		},
		{
			name:           "build non-existent route fails",
			route:          "/nonexistent-page-404-error",
			dir:            filepath.Join(tmpDir, "fail"),
			file:           filepath.Join(tmpDir, "fail", "index.html"),
			expectedStatus: http.StatusNotFound,
			shouldFail:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup minimal server for testing
			srv := &http.Server{
				Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/nonexistent-page-404-error" {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					w.WriteHeader(http.StatusOK)
					if _, err := w.Write([]byte("<html><body>Test Content</body></html>")); err != nil {
						t.Errorf("Failed to write response: %v", err)
					}
				}),
			}

			err := buildRoute(srv, tc.route, tc.dir, tc.file)

			if tc.shouldFail {
				if err == nil {
					t.Errorf("expected buildRoute to fail for route %q, but it succeeded", tc.route)
				}
				return
			}

			if err != nil {
				t.Fatalf("buildRoute failed: %v", err)
			}

			// Verify directory was created
			if _, err := os.Stat(tc.dir); os.IsNotExist(err) {
				t.Errorf("expected directory %q to exist", tc.dir)
			}

			// Verify file was created and contains content
			content, err := os.ReadFile(tc.file)
			if err != nil {
				t.Fatalf("failed to read output file: %v", err)
			}

			if len(content) == 0 {
				t.Error("output file is empty")
			}

			expectedContent := "<html><body>Test Content</body></html>"
			if string(content) != expectedContent {
				t.Errorf("content mismatch:\nwant: %s\ngot:  %s", expectedContent, string(content))
			}
		})
	}
}

func TestBuildRoute_InvalidRequest(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "xlog-build-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to clean up temp dir: %v", err)
		}
	}()

	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	// Test with invalid route containing control characters
	invalidRoute := "/test\x00invalid"
	err = buildRoute(srv, invalidRoute, tmpDir, filepath.Join(tmpDir, "output.html"))
	
	if err == nil {
		t.Error("expected error for invalid route with control characters")
	}
}

func TestBuildRoute_FilePermissions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "xlog-build-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to clean up temp dir: %v", err)
		}
	}()

	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("content")); err != nil {
				t.Errorf("Failed to write response: %v", err)
			}
		}),
	}

	outputFile := filepath.Join(tmpDir, "output.html")
	err = buildRoute(srv, "/", tmpDir, outputFile)
	if err != nil {
		t.Fatalf("buildRoute failed: %v", err)
	}

	// Check file permissions
	info, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("failed to stat output file: %v", err)
	}

	expectedPerms := build_perms
	actualPerms := info.Mode().Perm()
	if actualPerms != expectedPerms {
		t.Errorf("file permissions mismatch: want %v, got %v", expectedPerms, actualPerms)
	}
}
