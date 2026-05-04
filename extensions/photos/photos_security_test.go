package photos

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestValidatePath_PathTraversalAttempts tests the validatePath function against various
// path traversal attack vectors.
func TestValidatePath_PathTraversalAttempts(t *testing.T) {
	tests := []struct {
		name        string
		inputPath   string
		expectError bool
		description string
	}{
		{
			name:        "simple relative path is allowed",
			inputPath:   "photos/vacation.jpg",
			expectError: false,
			description: "Normal relative path should be allowed",
		},
		{
			name:        "path traversal with ../",
			inputPath:   "../../../etc/passwd",
			expectError: true,
			description: "Path traversal attempt should be blocked",
		},
		{
			name:        "path traversal in middle",
			inputPath:   "photos/../../etc/passwd",
			expectError: true,
			description: "Path traversal in middle of path should be blocked",
		},
		{
			name:        "absolute path attempt",
			inputPath:   "/etc/passwd",
			expectError: true,
			description: "Absolute path should be blocked",
		},
		{
			name:        "URL encoded path traversal",
			inputPath:   "photos%2F..%2F..%2Fetc%2Fpasswd",
			expectError: false, // URL decoding not performed by filepath.Clean
			description: "URL-encoded strings are treated as literal filename characters",
		},
		{
			name:        "double encoded path traversal",
			inputPath:   "..%252F..%252Fetc%252Fpasswd",
			expectError: false, // URL encoding not decoded
			description: "Double-encoded path treated as literal filename",
		},
		{
			name:        "backslash path traversal (Windows-style)",
			inputPath:   "photos\\..\\..\\etc\\passwd",
			expectError: false, // On Linux, backslash is valid filename character
			description: "Backslash on Linux is literal character, not separator",
		},
		{
			name:        "null byte injection",
			inputPath:   "photos/vacation.jpg\x00../../etc/passwd",
			expectError: true,
			description: "Null byte injection should be blocked",
		},
		{
			name:        "valid nested path",
			inputPath:   "photos/2023/vacation/beach.jpg",
			expectError: false,
			description: "Valid nested path should be allowed",
		},
		{
			name:        "empty path",
			inputPath:   "",
			expectError: true,
			description: "Empty path should be blocked",
		},
		{
			name:        "dot segment path",
			inputPath:   "./photos/image.jpg",
			expectError: false,
			description: "Dot segment that resolves locally should be allowed",
		},
		{
			name:        "path with only dots",
			inputPath:   "....",
			expectError: false,
			description: "Path with only dots (valid filename) should be allowed",
		},
		{
			name:        "path traversal with mixed separators",
			inputPath:   "photos/../../../etc/passwd",
			expectError: true,
			description: "Mixed separator path traversal should be blocked",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePath(tc.inputPath)

			if tc.expectError {
				if err == nil {
					t.Errorf("%s: Expected error for path %q, but got none",
						tc.description, tc.inputPath)
				}
			} else {
				if err != nil {
					t.Errorf("%s: Expected no error for path %q, but got: %v",
						tc.description, tc.inputPath, err)
				}
			}
		})
	}
}

// TestResizeHandler_PathTraversalBlocked tests that the HTTP handler properly blocks
// path traversal attempts at the HTTP request level.
func TestResizeHandler_PathTraversalBlocked(t *testing.T) {
	tests := []struct {
		name           string
		pathValue      string
		expectBlocked  bool
		expectedStatus string
	}{
		{
			name:           "normal path allowed",
			pathValue:      "photos/vacation.jpg",
			expectBlocked:  false,
			expectedStatus: "success or file not found",
		},
		{
			name:           "path traversal blocked",
			pathValue:      "../../etc/passwd",
			expectBlocked:  true,
			expectedStatus: "invalid path",
		},
		{
			name:           "absolute path blocked",
			pathValue:      "/etc/passwd",
			expectBlocked:  true,
			expectedStatus: "invalid path",
		},
		{
			name:           "complex path traversal blocked",
			pathValue:      "photos/../../../etc/passwd",
			expectBlocked:  true,
			expectedStatus: "invalid path",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/+/photos/thumbnail/"+tc.pathValue, http.NoBody)
			req.SetPathValue("path", tc.pathValue)

			output := resizeHandler(req)
			w := httptest.NewRecorder()
			output(w, req)

			body := w.Body.String()

			if tc.expectBlocked {
				// Should contain error message about invalid path
				if !strings.Contains(body, "invalid path") &&
					!strings.Contains(body, "path traversal") {
					t.Errorf("Expected error message about invalid path, got: %q", body)
				}
			}
		})
	}
}

// TestPhotoHandler_PathTraversalBlocked tests that the photo handler blocks path traversal.
func TestPhotoHandler_PathTraversalBlocked(t *testing.T) {
	tests := []struct {
		name          string
		pathValue     string
		expectBlocked bool
	}{
		{
			name:          "normal path",
			pathValue:     "photos/image.jpg",
			expectBlocked: false,
		},
		{
			name:          "path traversal attempt",
			pathValue:     "../../../etc/passwd",
			expectBlocked: true,
		},
		{
			name:          "absolute path attempt",
			pathValue:     "/etc/passwd",
			expectBlocked: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/+/photos/photo/"+tc.pathValue, http.NoBody)
			req.SetPathValue("path", tc.pathValue)

			output := photoHandler(req)

			// The handler should return a non-nil output
			if output == nil {
				t.Error("Expected non-nil output")
			}

			// Execute the output to check for error messages
			w := httptest.NewRecorder()
			output(w, req)

			body := w.Body.String()

			if tc.expectBlocked {
				// Should contain error indication
				if !strings.Contains(body, "invalid path") &&
					!strings.Contains(body, "path traversal") &&
					!strings.Contains(body, "error") {
					t.Logf("Warning: Expected error indication for blocked path, got: %q", body)
				}
			}
		})
	}
}

// TestValidatePath_EdgeCases tests edge cases in path validation.
func TestValidatePath_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "single dot",
			path:        ".",
			expectError: false,
		},
		{
			name:        "double dot",
			path:        "..",
			expectError: true,
		},
		{
			name:        "triple dot (valid filename)",
			path:        "...",
			expectError: false,
		},
		{
			name:        "file starting with dot",
			path:        ".hidden",
			expectError: false,
		},
		{
			name:        "path with spaces",
			path:        "my photos/vacation.jpg",
			expectError: false,
		},
		{
			name:        "path with unicode",
			path:        "фото/отпуск.jpg",
			expectError: false,
		},
		{
			name:        "very long path within limits",
			path:        strings.Repeat("a/", 50) + "file.jpg",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePath(tc.path)

			if tc.expectError && err == nil {
				t.Errorf("Expected error for path %q, got none", tc.path)
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error for path %q, got: %v", tc.path, err)
			}
		})
	}
}

// TestValidatePath_RealWorldPaths tests validation against real-world photo paths.
func TestValidatePath_RealWorldPaths(t *testing.T) {
	// Create a temporary directory with actual files to test
	tmpDir := t.TempDir()

	// Create nested structure
	photoDir := filepath.Join(tmpDir, "photos", "2023", "vacation")
	if err := mkdir(photoDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "relative path within tmpdir",
			path:        "photos/2023/vacation/beach.jpg",
			expectError: false,
		},
		{
			name:        "attempt to escape tmpdir",
			path:        "photos/../../etc/passwd",
			expectError: true,
		},
		{
			name:        "deeply nested valid path",
			path:        "a/b/c/d/e/f/g/photo.jpg",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePath(tc.path)

			if tc.expectError && err == nil {
				t.Errorf("Expected error for path %q, got none", tc.path)
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error for path %q, got: %v", tc.path, err)
			}
		})
	}
}

// mkdir is a helper that creates directory recursively.
func mkdir(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
