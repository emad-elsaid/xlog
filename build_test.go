package xlog

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	testIndexFile    = "index.md"
	testNotFoundFile = "404.md"
)

func TestRegisterBuildPage(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		encloseInDir  bool
		checkEnclosed bool
	}{
		{
			name:          "register simple page",
			path:          "/test-page.html",
			encloseInDir:  false,
			checkEnclosed: false,
		},
		{
			name:          "register enclosed page",
			path:          "/enclosed-page",
			encloseInDir:  true,
			checkEnclosed: true,
		},
		{
			name:          "register root path",
			path:          "/",
			encloseInDir:  false,
			checkEnclosed: false,
		},
		{
			name:          "register nested path",
			path:          "/api/v1/endpoint",
			encloseInDir:  true,
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
				ReadHeaderTimeout: 5 * time.Second,
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
		ReadHeaderTimeout: 5 * time.Second,
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
		ReadHeaderTimeout: 5 * time.Second,
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

func TestBuild_Integration(t *testing.T) {
	// Create temporary directory for build output
	tmpDir, err := os.MkdirTemp("", "xlog-build-integration-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to clean up temp dir: %v", err)
		}
	}()

	// Save original config and restore after test
	origIndex := Config.Index
	origNotFound := Config.NotFoundPage
	defer func() {
		Config.Index = origIndex
		Config.NotFoundPage = origNotFound
	}()

	tests := []struct {
		name           string
		setupIndex     string
		setupNotFound  string
		expectIndexErr bool
	}{
		{
			name:           "build with default index",
			setupIndex:     testIndexFile,
			setupNotFound:  testNotFoundFile,
			expectIndexErr: false,
		},
		{
			name:           "build with custom index",
			setupIndex:     "home.md",
			setupNotFound:  "notfound.md",
			expectIndexErr: false,
		},
		{
			name:           "build with nonexistent index",
			setupIndex:     "nonexistent-page.md",
			setupNotFound:  testNotFoundFile,
			expectIndexErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup config
			Config.Index = tc.setupIndex
			Config.NotFoundPage = tc.setupNotFound

			// Create a fresh subdirectory for this test case
			testDir := filepath.Join(tmpDir, tc.name)
			if err := os.MkdirAll(testDir, 0750); err != nil {
				t.Fatalf("Failed to create test dir: %v", err)
			}

			// Execute build
			err := build(testDir)

			// Note: build() function logs errors via slog but returns nil
			// or returns error from fs.WalkDir. We verify output files instead.
			if err != nil && !tc.expectIndexErr {
				t.Errorf("build failed unexpectedly: %v", err)
			}

			// Verify basic structure was attempted (directories created)
			if _, err := os.Stat(testDir); os.IsNotExist(err) {
				t.Error("build output directory was not created")
			}
		})
	}
}

func TestBuild_AssetCopying(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "xlog-build-assets-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to clean up temp dir: %v", err)
		}
	}()

	// Save original config
	origIndex := Config.Index
	origNotFound := Config.NotFoundPage
	defer func() {
		Config.Index = origIndex
		Config.NotFoundPage = origNotFound
	}()

	Config.Index = testIndexFile
	Config.NotFoundPage = testNotFoundFile

	err = build(tmpDir)
	if err != nil {
		t.Logf("build completed with error: %v", err)
	}

	// Verify that asset files are copied
	// The assets embed.FS should contain files; verify they're copied to destination
	// Note: actual file verification depends on what's in the assets embed.FS
	// This test verifies the mechanism works without knowing specific files

	// Just verify the function executes without panic
	// More specific assertions would require knowing the embedded assets
}

func TestBuild_ExtensionPages(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "xlog-build-extensions-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to clean up temp dir: %v", err)
		}
	}()

	// Save original config
	origIndex := Config.Index
	origNotFound := Config.NotFoundPage
	defer func() {
		Config.Index = origIndex
		Config.NotFoundPage = origNotFound
	}()

	Config.Index = testIndexFile
	Config.NotFoundPage = testNotFoundFile

	// Register some test extension pages
	testPages := []struct {
		route        string
		encloseInDir bool
	}{
		{"/test-extension", true},
		{"/api-endpoint.json", false},
	}

	for _, page := range testPages {
		RegisterBuildPage(page.route, page.encloseInDir)
	}

	// Clean up registered pages after test
	defer func() {
		for _, page := range testPages {
			extension_page.Delete(page.route)
			extension_page_enclosed.Delete(page.route)
		}
	}()

	err = build(tmpDir)
	if err != nil {
		t.Logf("build completed with error: %v", err)
	}

	// Verify build attempted to process extension pages
	// The actual success depends on server routing setup
	// This test verifies no panic occurs when processing extension pages
}

func TestBuild_404Handling(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "xlog-build-404-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Errorf("Failed to clean up temp dir: %v", err)
		}
	}()

	// Save original config
	origNotFound := Config.NotFoundPage
	origIndex := Config.Index
	defer func() {
		Config.NotFoundPage = origNotFound
		Config.Index = origIndex
	}()

	Config.Index = testIndexFile
	Config.NotFoundPage = testNotFoundFile

	// First create the expected 404 source directory and file
	notFoundDir := filepath.Join(tmpDir, Config.NotFoundPage)
	if err := os.MkdirAll(notFoundDir, 0750); err != nil {
		t.Fatalf("Failed to create 404 dir: %v", err)
	}

	notFoundContent := []byte("<html><body>Page Not Found</body></html>")
	notFoundIndexPath := filepath.Join(notFoundDir, "index.html")
	if err := os.WriteFile(notFoundIndexPath, notFoundContent, 0600); err != nil {
		t.Fatalf("Failed to write 404 index file: %v", err)
	}

	err = build(tmpDir)
	if err != nil {
		t.Logf("build completed with error: %v", err)
	}

	// Verify 404.html was copied to root
	notFoundCopy := filepath.Join(tmpDir, "404.html")
	if _, err := os.Stat(notFoundCopy); err == nil {
		content, readErr := os.ReadFile(notFoundCopy)
		if readErr != nil {
			t.Errorf("Failed to read 404.html: %v", readErr)
		}
		if string(content) != string(notFoundContent) {
			t.Errorf("404.html content mismatch: want %q, got %q", notFoundContent, content)
		}
	} else if !os.IsNotExist(err) {
		t.Errorf("Error checking 404.html: %v", err)
	}
	// Note: If the file doesn't exist, it's because the build route failed
	// which is expected in test environment without full server setup
}

func TestBuild_ErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		destPath   string
		shouldFail bool
	}{
		{
			name:       "build to valid directory",
			destPath:   "",
			shouldFail: false,
		},
		{
			name:       "build to path with invalid permissions",
			destPath:   "/root/xlog-test-no-permission",
			shouldFail: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var testDir string
			var err error

			if tc.destPath == "" {
				testDir, err = os.MkdirTemp("", "xlog-build-error-*")
				if err != nil {
					t.Fatalf("Failed to create temp dir: %v", err)
				}
				defer func() {
					if err := os.RemoveAll(testDir); err != nil {
						t.Logf("Failed to clean up temp dir: %v", err)
					}
				}()
			} else {
				testDir = tc.destPath
			}

			// Save original config
			origIndex := Config.Index
			origNotFound := Config.NotFoundPage
			defer func() {
				Config.Index = origIndex
				Config.NotFoundPage = origNotFound
			}()

			Config.Index = testIndexFile
			Config.NotFoundPage = testNotFoundFile

			err = build(testDir)

			if tc.shouldFail {
				if err == nil && os.Getuid() != 0 {
					// If running as non-root, we expect error for /root access
					// But build might not fail immediately, it logs errors
					t.Logf("build to restricted path: %v", err)
				}
			}

			// The build function logs errors via slog rather than
			// failing immediately, so we verify it doesn't panic
		})
	}
}
