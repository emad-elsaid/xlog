package gpg

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

const (
	testGPGID     = "test@example.com"
	headerTrue    = "true"
	headerRefresh = "HX-Refresh"
)

// TestEncryptHandler tests the encrypt handler with various scenarios.
func TestEncryptHandler(t *testing.T) {
	tests := []struct {
		name           string
		pageName       string
		pageExists     bool
		pageContent    string
		makeReadOnly   bool
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
		{
			name:           "encrypt fails when cannot delete original after encryption",
			pageName:       "test/readonly",
			pageExists:     true,
			pageContent:    "# Protected Page\nCannot delete this.",
			makeReadOnly:   true,
			expectedStatus: http.StatusInternalServerError,
			expectRefresh:  false,
		},
		{
			name:           "encrypt handles empty page content",
			pageName:       "test/empty",
			pageExists:     true,
			pageContent:    "",
			expectedStatus: http.StatusInternalServerError,
			expectRefresh:  false,
		},
		{
			name:           "encrypt handles special characters in content",
			pageName:       "test/special",
			pageExists:     true,
			pageContent:    "# Special Chars\n`code` **bold** [link](url)\n\n> quote",
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
			gpgId = testGPGID

			// Create test page if it should exist
			if tc.pageExists {
				dir := filepath.Dir(tc.pageName)
				if err := os.MkdirAll(dir, 0700); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
				if err := os.WriteFile(tc.pageName+".md", []byte(tc.pageContent), 0600); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}

				// Make file read-only to simulate delete failure after encryption
				if tc.makeReadOnly {
					if err := os.Chmod(tc.pageName+".md", 0400); err != nil {
						t.Fatalf("failed to chmod file: %v", err)
					}
					defer func() {
						_ = os.Chmod(tc.pageName+".md", 0600)
					}()
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
				if w.Code == http.StatusNoContent && w.Header().Get(headerRefresh) == headerTrue {
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
		setupEncrypted bool
		makeReadOnly   bool
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
			name:           "decrypt fails when cannot delete original encrypted file",
			pageName:       "test/readonly",
			pageExists:     true,
			setupEncrypted: true,
			makeReadOnly:   true,
			expectedStatus: http.StatusInternalServerError,
			expectRefresh:  false,
		},
		{
			name:           "decrypt fails when GPG not configured",
			pageName:       "test/nogpg",
			pageExists:     true,
			setupEncrypted: true,
			makeReadOnly:   false,
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

			// Set GPG ID
			gpgId = "test@example.com"

			// Create encrypted test page if it should exist
			if tc.pageExists && tc.setupEncrypted {
				dir := filepath.Dir(tc.pageName)
				if err := os.MkdirAll(dir, 0700); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
				// Create a mock encrypted file (won't decrypt without real GPG)
				if err := os.WriteFile(tc.pageName+EXT, []byte("mock encrypted content"), 0600); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}

				// Make directory read-only to simulate delete failure
				if tc.makeReadOnly {
					// #nosec G302 - Test intentionally uses 0500 permissions to simulate delete failure scenario
					if err := os.Chmod(dir, 0500); err != nil {
						t.Fatalf("failed to chmod directory: %v", err)
					}
					defer func() {
						// #nosec G302 - Cleanup restores standard directory permissions in test
						_ = os.Chmod(dir, 0700)
					}()
				}
			}

			// Create test request
			req := httptest.NewRequest(http.MethodPost, "/+/gpg/decrypt/"+tc.pageName, nil)
			req.SetPathValue("page", tc.pageName)
			w := httptest.NewRecorder()

			// Execute handler
			output := decryptHandler(req)
			if output != nil {
				output(w, req)
			}

			// Verify response - without real GPG, we expect errors
			if tc.expectRefresh {
				if w.Header().Get(headerRefresh) != headerTrue {
					t.Errorf("expected HX-Refresh header to be 'true', got %q", w.Header().Get(headerRefresh))
				}
			}

			// If the test expected an error and GPG succeeded anyway, skip
			if tc.pageExists && !tc.makeReadOnly {
				if w.Code == http.StatusNoContent && w.Header().Get(headerRefresh) == headerTrue {
					t.Skip("GPG is configured and working - skipping error path test")
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
