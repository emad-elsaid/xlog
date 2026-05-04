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

// Benchmark FindBrokenLinks with varying garden sizes to measure scalability.
// Digital gardens can grow to thousands of pages, so performance at scale matters.
func BenchmarkFindBrokenLinks(b *testing.B) {
	tests := []struct {
		name      string
		pageCount int
		linksPer  int // Average links per page
	}{
		{"small garden (10 pages)", 10, 3},
		{"medium garden (100 pages)", 100, 5},
		{"large garden (500 pages)", 500, 8},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			// Setup test environment once
			tmpDir := b.TempDir()
			origDir, _ := os.Getwd()
			defer func() { _ = os.Chdir(origDir) }()

			if err := os.Chdir(tmpDir); err != nil {
				b.Fatalf("Failed to change directory: %v", err)
			}

			// Create test pages with realistic content
			for i := 0; i < tt.pageCount; i++ {
				filename := filepath.Join(tmpDir, "page"+filepath.FromSlash(string(rune('0'+i%10))), "content.md")
				dir := filepath.Dir(filename)
				if err := os.MkdirAll(dir, 0755); err != nil {
					b.Fatalf("Failed to create directory: %v", err)
				}

				// Generate page content with internal links
				content := "# Page " + string(rune('0'+i)) + "\n\nSome content here.\n\n"
				for j := 0; j < tt.linksPer; j++ {
					targetIdx := (i + j + 1) % tt.pageCount
					content += "Link to [page" + string(rune('0'+targetIdx%10)) + "](page" +
						string(rune('0'+targetIdx%10)) + "/content)\n"
				}

				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
					b.Fatalf("Failed to write file: %v", err)
				}
			}

			Config.Source = tmpDir
			_ = clearPagesCache(nil)

			// Reset timer after setup
			b.ResetTimer()

			// Run benchmark
			for i := 0; i < b.N; i++ {
				_ = FindBrokenLinks(context.Background())
			}
		})
	}
}

// BenchmarkIsExternalLink measures performance of external link detection.
// This function is called for every link in every page, so it's a hot path.
func BenchmarkIsExternalLink(b *testing.B) {
	tests := []struct {
		name string
		dest string
	}{
		{"http URL", "http://example.com/path/to/page"},
		{"https URL", "https://example.com/very/long/path/to/resource"},
		{"mailto", "mailto:user@example.com"},
		{"internal relative", "page-name"},
		{"internal absolute", "/folder/page-name"},
		{"protocol relative", "//cdn.example.com/resource"},
		{"tel link", "tel:+1234567890"},
		{"ftp link", "ftp://files.example.com/file.zip"},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = isExternalLink(tt.dest)
			}
		})
	}
}

// BenchmarkLinkToPageName measures performance of link-to-page-name conversion.
// Like isExternalLink, this is called for every internal link.
func BenchmarkLinkToPageName(b *testing.B) {
	tests := []struct {
		name string
		dest string
	}{
		{"simple", "page-name"},
		{"with slash", "/page-name"},
		{"with extension", "page-name.md"},
		{"with fragment", "page-name#section"},
		{"with query", "page-name?foo=bar"},
		{"folder path", "folder/subfolder/page"},
		{"complex", "/folder/page.md#anchor?query=value"},
		{"redundant slashes", "//folder//subfolder//page"},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = linkToPageName(tt.dest)
			}
		})
	}
}

// BenchmarkFindBrokenLinks_BrokenRatio tests performance with varying ratios of broken links.
// Pages with many broken links might perform differently due to allocation patterns.
func BenchmarkFindBrokenLinks_BrokenRatio(b *testing.B) {
	tests := []struct {
		name        string
		pageCount   int
		brokenRatio float64 // Ratio of links that are broken (0.0 to 1.0)
	}{
		{"no broken links", 50, 0.0},
		{"10% broken", 50, 0.1},
		{"50% broken", 50, 0.5},
		{"all broken", 50, 1.0},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			tmpDir := b.TempDir()
			origDir, _ := os.Getwd()
			defer func() { _ = os.Chdir(origDir) }()

			if err := os.Chdir(tmpDir); err != nil {
				b.Fatalf("Failed to change directory: %v", err)
			}

			linksPerPage := 10
			brokenCount := int(float64(linksPerPage) * tt.brokenRatio)
			validCount := linksPerPage - brokenCount

			// Create pages
			for i := 0; i < tt.pageCount; i++ {
				filename := "page" + string(rune('0'+(i%10))) + ".md"
				content := "# Page " + string(rune('0'+i)) + "\n\n"

				// Add valid links
				for j := 0; j < validCount; j++ {
					targetIdx := (i + j + 1) % tt.pageCount
					content += "[link](page" + string(rune('0'+(targetIdx%10))) + ")\n"
				}

				// Add broken links
				for j := 0; j < brokenCount; j++ {
					content += "[broken](nonexistent-page-" + string(rune('0'+j)) + ")\n"
				}

				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
					b.Fatalf("Failed to write file: %v", err)
				}
			}

			Config.Source = tmpDir
			_ = clearPagesCache(nil)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = FindBrokenLinks(context.Background())
			}
		})
	}
}
