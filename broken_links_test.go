package xlog

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestIsExternalLink(t *testing.T) {
	tests := []struct {
		name     string
		dest     string
		expected bool
	}{
		{
			name:     "http URL",
			dest:     "http://example.com",
			expected: true,
		},
		{
			name:     "https URL",
			dest:     "https://example.com",
			expected: true,
		},
		{
			name:     "mailto link",
			dest:     "mailto:user@example.com",
			expected: true,
		},
		{
			name:     "ftp link",
			dest:     "ftp://example.com/file",
			expected: true,
		},
		{
			name:     "protocol-relative URL",
			dest:     "//example.com/path",
			expected: true,
		},
		{
			name:     "internal relative link",
			dest:     "page-name",
			expected: false,
		},
		{
			name:     "internal absolute link",
			dest:     "/page-name",
			expected: false,
		},
		{
			name:     "internal link with .md",
			dest:     "folder/page.md",
			expected: false,
		},
		{
			name:     "anchor only",
			dest:     "#section",
			expected: false,
		},
		{
			name:     "tel link",
			dest:     "tel:+1234567890",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isExternalLink(tt.dest)
			if result != tt.expected {
				t.Errorf("isExternalLink(%q) = %v, expected %v",
					tt.dest, result, tt.expected)
			}
		})
	}
}

func TestLinkToPageName(t *testing.T) {
	tests := []struct {
		name     string
		dest     string
		expected string
	}{
		{
			name:     "simple page name",
			dest:     "page-name",
			expected: "page-name",
		},
		{
			name:     "removes leading slash",
			dest:     "/page-name",
			expected: "page-name",
		},
		{
			name:     "removes .md extension",
			dest:     "page-name.md",
			expected: "page-name",
		},
		{
			name:     "removes both slash and extension",
			dest:     "/page-name.md",
			expected: "page-name",
		},
		{
			name:     "removes fragment",
			dest:     "page-name#section",
			expected: "page-name",
		},
		{
			name:     "removes query string",
			dest:     "page-name?foo=bar",
			expected: "page-name",
		},
		{
			name:     "handles folder paths",
			dest:     "folder/page-name",
			expected: "folder/page-name",
		},
		{
			name:     "handles folder with slash and extension",
			dest:     "/folder/page-name.md",
			expected: "folder/page-name",
		},
		{
			name:     "handles complex path",
			dest:     "/folder/subfolder/page.md#anchor",
			expected: "folder/subfolder/page",
		},
		{
			name:     "cleans redundant slashes",
			dest:     "//folder//page",
			expected: "folder/page",
		},
		{
			name:     "handles fragment with extension",
			dest:     "page.md#section",
			expected: "page",
		},
		{
			name:     "handles query with extension",
			dest:     "page.md?query",
			expected: "page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := linkToPageName(tt.dest)
			if result != tt.expected {
				t.Errorf("linkToPageName(%q) = %q, expected %q",
					tt.dest, result, tt.expected)
			}
		})
	}
}

func TestFindBrokenLinks(t *testing.T) {
	tests := []struct {
		name          string
		files         map[string]string
		expectedCount int
		checkBroken   func(t *testing.T, broken []BrokenLink)
	}{
		{
			name: "no broken links",
			files: map[string]string{
				"page1.md": "Content linking to [page2](page2)",
				"page2.md": "Content linking to [page1](/page1)",
			},
			expectedCount: 0,
		},
		{
			name: "one broken link",
			files: map[string]string{
				"page1.md": "Content linking to [nonexistent](nonexistent)",
			},
			expectedCount: 1,
			checkBroken: func(t *testing.T, broken []BrokenLink) {
				if broken[0].SourcePage != "page1" {
					t.Errorf("Expected source page 'page1', got '%s'", broken[0].SourcePage)
				}
				if broken[0].TargetPage != "nonexistent" {
					t.Errorf("Expected target page 'nonexistent', got '%s'", broken[0].TargetPage)
				}
			},
		},
		{
			name: "multiple broken links in same page",
			files: map[string]string{
				"page1.md": "Links to [missing1](missing1) and [missing2](missing2)",
			},
			expectedCount: 2,
		},
		{
			name: "ignores external links",
			files: map[string]string{
				"page1.md": "External [link](https://example.com) and broken [internal](missing)",
			},
			expectedCount: 1,
			checkBroken: func(t *testing.T, broken []BrokenLink) {
				if broken[0].TargetPage != "missing" {
					t.Errorf("Should only detect internal broken link, got '%s'", broken[0].TargetPage)
				}
			},
		},
		{
			name: "ignores anchors",
			files: map[string]string{
				"page1.md": "Anchor [link](#section) and broken [page](missing)",
			},
			expectedCount: 1,
		},
		{
			name: "handles .md extension in links",
			files: map[string]string{
				"page1.md": "Link to [existing](page2.md) and [missing](missing.md)",
				"page2.md": "Content here",
			},
			expectedCount: 1,
		},
		{
			name: "handles links with fragments",
			files: map[string]string{
				"page1.md": "Link to [existing section](page2#section) and [missing](missing#section)",
				"page2.md": "Content here",
			},
			expectedCount: 1,
			checkBroken: func(t *testing.T, broken []BrokenLink) {
				if broken[0].TargetPage != "missing" {
					t.Errorf("Fragment should be stripped, got '%s'", broken[0].TargetPage)
				}
			},
		},
		{
			name: "handles folder structure",
			files: map[string]string{
				"folder/page1.md": "Link to [page2](page2) in same folder",
				"folder/page2.md": "Content here",
			},
			expectedCount: 0,
		},
		{
			name: "mixed working and broken links",
			files: map[string]string{
				"page1.md": `
Links: 
- [existing1](page2)
- [broken1](missing1)
- [existing2](https://example.com)
- [broken2](missing2)
- [anchor](#test)
				`,
				"page2.md": "Content",
			},
			expectedCount: 2,
		},
		{
			name: "empty page",
			files: map[string]string{
				"empty.md": "",
			},
			expectedCount: 0,
		},
		{
			name: "page with only external links",
			files: map[string]string{
				"external.md": "[Link1](https://example.com) and [Link2](http://test.org)",
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			tmpDir := t.TempDir()
			origDir, _ := os.Getwd()
			defer func() {
				_ = os.Chdir(origDir)
			}()

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Failed to change to temp dir: %v", err)
			}

			// Create test files
			for filename, content := range tt.files {
				dir := filepath.Dir(filename)
				if dir != "." {
					if err := os.MkdirAll(dir, 0755); err != nil {
						t.Fatalf("Failed to create directory %s: %v", dir, err)
					}
				}

				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write file %s: %v", filename, err)
				}
			}

			// Run test
			broken := FindBrokenLinks(context.Background())

			// Check count
			if len(broken) != tt.expectedCount {
				t.Errorf("Expected %d broken links, got %d", tt.expectedCount, len(broken))
				for _, bl := range broken {
					t.Logf("  Found: %s -> %s (%s)", bl.SourcePage, bl.TargetPage, bl.LinkDestination)
				}
			}

			// Run additional checks if provided
			if tt.checkBroken != nil && len(broken) > 0 {
				tt.checkBroken(t, broken)
			}
		})
	}
}

func TestPrintBrokenLinks(t *testing.T) {
	tests := []struct {
		name   string
		broken []BrokenLink
	}{
		{
			name:   "empty list",
			broken: []BrokenLink{},
		},
		{
			name: "single broken link",
			broken: []BrokenLink{
				{
					SourcePage:      "page1",
					TargetPage:      "missing",
					LinkDestination: "missing",
				},
			},
		},
		{
			name: "multiple broken links",
			broken: []BrokenLink{
				{
					SourcePage:      "page1",
					TargetPage:      "missing1",
					LinkDestination: "missing1",
				},
				{
					SourcePage:      "page1",
					TargetPage:      "missing2",
					LinkDestination: "/missing2",
				},
				{
					SourcePage:      "page2",
					TargetPage:      "missing3",
					LinkDestination: "missing3.md",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just ensure it doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("PrintBrokenLinks panicked: %v", r)
				}
			}()

			PrintBrokenLinks(tt.broken)
		})
	}
}

func TestFindBrokenLinks_CaseInsensitive(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() {
		_ = os.Chdir(origDir)
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Create files
	files := map[string]string{
		"Page1.md": "Link to [page2](Page2)", // Different case
		"page2.md": "Link to [page1](page1)", // Different case
	}

	for filename, content := range files {
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}

	broken := FindBrokenLinks(context.Background())

	// On case-sensitive filesystems, these might be broken
	// On case-insensitive filesystems (macOS, Windows), they shouldn't be
	// This test just ensures the function doesn't crash with different cases
	if len(broken) > 0 {
		t.Logf("Note: Found %d broken links (may vary by filesystem case sensitivity)", len(broken))
	}
}

func TestBrokenLink_Struct(t *testing.T) {
	// Test that BrokenLink struct can be created and accessed
	bl := BrokenLink{
		SourcePage:      "source",
		TargetPage:      "target",
		LinkDestination: "/target",
	}

	if bl.SourcePage != "source" {
		t.Errorf("SourcePage = %q, expected 'source'", bl.SourcePage)
	}
	if bl.TargetPage != "target" {
		t.Errorf("TargetPage = %q, expected 'target'", bl.TargetPage)
	}
	if bl.LinkDestination != "/target" {
		t.Errorf("LinkDestination = %q, expected '/target'", bl.LinkDestination)
	}
}
