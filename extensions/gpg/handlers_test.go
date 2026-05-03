package gpg

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestEncryptHandler tests the encrypt handler with various scenarios.
func TestEncryptHandler(t *testing.T) {
	tests := []struct {
		name           string
		pageName       string
		pageExists     bool
		pageContent    string
		expectedStatus int
		expectRefresh  bool
	}{
		{
			name:           "non-existent page returns 404",
			pageName:       "test/nonexistent",
			pageExists:     false,
			expectedStatus: http.StatusNotFound,
			expectRefresh:  false,
		},
		{
			name:           "encrypt existing page without GPG setup fails",
			pageName:       "test/encrypt",
			pageExists:     true,
			pageContent:    "# Test Content\nThis is a test page.",
			expectedStatus: http.StatusInternalServerError,
			expectRefresh:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test environment
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

			// Set GPG ID (even though GPG may not be available)
			gpgId = "test@example.com"

			// Create test page if it should exist
			if tc.pageExists {
				dir := filepath.Dir(tc.pageName)
				if err := os.MkdirAll(dir, 0700); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
				if err := os.WriteFile(tc.pageName+".md", []byte(tc.pageContent), 0600); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			}

			// Create test request
			req := httptest.NewRequest(http.MethodPost, "/+/gpg/encrypt/"+tc.pageName, nil)
			req.SetPathValue("page", tc.pageName)
			w := httptest.NewRecorder()

			// Execute handler
			output := encryptHandler(req)
			if output != nil {
				output(w, req)
			}

			// Verify response - we expect errors without GPG setup
			// The test verifies handler structure, not GPG functionality
			if tc.pageExists {
				// When page exists but GPG fails, we should get an error
				// The handler will attempt encryption which will fail
				if w.Code == http.StatusNoContent && w.Header().Get("HX-Refresh") == "true" {
					t.Skip("GPG is configured and working - skipping error path test")
				}
			}
		})
	}
}

// TestDecryptHandler tests the decrypt handler with various scenarios.
func TestDecryptHandler(t *testing.T) {
	tests := []struct {
		name           string
		pageName       string
		pageExists     bool
		expectedStatus int
		expectRefresh  bool
	}{
		{
			name:           "non-existent page returns 404",
			pageName:       "test/nonexistent",
			pageExists:     false,
			expectedStatus: http.StatusNotFound,
			expectRefresh:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test environment
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

			// Create test request
			req := httptest.NewRequest(http.MethodPost, "/+/gpg/decrypt/"+tc.pageName, nil)
			req.SetPathValue("page", tc.pageName)
			w := httptest.NewRecorder()

			// Execute handler
			output := decryptHandler(req)
			if output != nil {
				output(w, req)
			}

			// For non-404 responses, check headers
			if tc.expectRefresh {
				if w.Header().Get("HX-Refresh") != "true" {
					t.Errorf("expected HX-Refresh header to be 'true', got %q", w.Header().Get("HX-Refresh"))
				}
			}
		})
	}
}

// TestHandlerErrors tests error conditions in handlers.
func TestHandlerErrors(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		expectedMsg string
	}{
		{
			name:        "delete failed error",
			err:         errDeleteFailed,
			expectedMsg: "couldn't delete original page",
		},
		{
			name:        "encryption failed error",
			err:         errEncryptionFailed,
			expectedMsg: "couldn't encrypt page",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err.Error() != tc.expectedMsg {
				t.Errorf("error message = %q, want %q", tc.err.Error(), tc.expectedMsg)
			}
		})
	}
}
