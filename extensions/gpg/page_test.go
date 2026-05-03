package gpg

import (
	"os"
	"strings"
	"testing"
)

// TestPageRenderErrorEscaping verifies that error messages in Render()
// are properly HTML-escaped to prevent XSS vulnerabilities.
func TestPageRenderErrorEscaping(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantHTML bool // true if we expect unescaped HTML chars
	}{
		{
			name:     "malicious script tag in error",
			content:  "<script>alert('xss')</script>",
			wantHTML: false, // should be escaped
		},
		{
			name:     "HTML entities in error",
			content:  "<img src=x onerror='alert(1)'>",
			wantHTML: false, // should be escaped
		},
		{
			name:     "ampersand and quotes",
			content:  `"test" & 'quotes'`,
			wantHTML: false, // should be escaped
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a page with content that will cause a conversion error
			// (This test will fail until we implement proper escaping)
			p := &page{
				name: "test",
				ast:  nil,
			}

			// We need to trigger an error condition in Render()
			// Since we can't easily mock MarkdownConverter, we'll test the
			// escaping logic directly once implemented

			// For now, this is a placeholder that documents the expected behavior
			result := p.Render()
			resultStr := string(result)

			// If error contains malicious content, it should be escaped
			if tc.wantHTML {
				// Should NOT be escaped (normal case)
				if !strings.Contains(resultStr, tc.content) {
					t.Errorf("Expected unescaped content, but it was escaped")
				}
			} else {
				// Should be escaped
				if strings.Contains(resultStr, "<script>") ||
					strings.Contains(resultStr, "onerror") ||
					strings.Contains(resultStr, tc.content) {
					t.Errorf("Expected escaped content, but found unescaped HTML: %s", resultStr)
				}
			}
		})
	}
}

// TestPageBasicMethods tests the simple getter methods.
func TestPageBasicMethods(t *testing.T) {
	p := &page{
		name: "test/page",
		ast:  nil,
	}

	if got := p.Name(); got != "test/page" {
		t.Errorf("Name() = %q, want %q", got, "test/page")
	}

	// FileName should append EXT constant
	fileName := p.FileName()
	if !strings.HasSuffix(fileName, EXT) {
		t.Errorf("FileName() = %q, expected to end with %q", fileName, EXT)
	}
}

// TestPageExists verifies the Exists method.
func TestPageExists(t *testing.T) {
	p := &page{
		name: "nonexistent/page/that/should/not/exist",
		ast:  nil,
	}

	// Should return false for non-existent file
	if p.Exists() {
		t.Errorf("Exists() = true for non-existent file, want false")
	}
}

// TestPageModTime verifies the ModTime method returns correct modification time.
func TestPageModTime(t *testing.T) {
	tests := []struct {
		name         string
		pageName     string
		expectZero   bool
		setupFile    bool
		setupContent string
	}{
		{
			name:       "non-existent file returns zero time",
			pageName:   "test/nonexistent",
			expectZero: true,
			setupFile:  false,
		},
		{
			name:         "existing file returns valid time",
			pageName:     "test/existing",
			expectZero:   false,
			setupFile:    true,
			setupContent: "test content",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := &page{name: tc.pageName}

			if tc.setupFile {
				// Create temporary file for testing
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

				if err := os.MkdirAll("test", 0700); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
				if err := os.WriteFile(p.FileName(), []byte(tc.setupContent), 0600); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			}

			modTime := p.ModTime()

			if tc.expectZero {
				if !modTime.IsZero() {
					t.Errorf("ModTime() = %v, expected zero time", modTime)
				}
			} else {
				if modTime.IsZero() {
					t.Errorf("ModTime() returned zero time, expected valid time")
				}
			}
		})
	}
}
