package xlog

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStats(t *testing.T) {
	// Setup test environment with markdown files
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create test markdown files
	files := map[string]string{
		"page1.md": "# Page 1\nThis is page 1 with some content. It links to [[page2]].",
		"page2.md": "# Page 2\nAnother page with content. Links to [[page1]] and [[page3]].",
		"page3.md": "# Page 3\nShort page.",
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	Config.Source = tmpDir
	_ = clearPagesCache(nil)

	tests := []struct {
		name         string
		wantContains []string
	}{
		{
			name: "displays statistics",
			wantContains: []string{
				"Digital Garden Statistics",
				"Total Pages:",
				"Total Words:",
				"Average Words per Page:",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			printStats(context.Background(), &buf)

			output := buf.String()
			for _, want := range tc.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("Output missing %q\nGot: %s", want, output)
				}
			}
		})
	}
}

func TestCalculateStats(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create test files
	files := map[string]string{
		"test1.md": "# Title\nOne two three four five",
		"test2.md": "# Another\nSix seven eight",
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	Config.Source = tmpDir
	_ = clearPagesCache(nil)

	stats := calculateStats(context.Background())

	if stats.TotalPages != 2 {
		t.Errorf("TotalPages = %d, want 2", stats.TotalPages)
	}

	if stats.TotalWords < 8 { // Should have at least 8 words
		t.Errorf("TotalWords = %d, want >= 8", stats.TotalWords)
	}

	if stats.AvgWordsPerPage == 0 {
		t.Errorf("AvgWordsPerPage should not be 0")
	}
}

func TestStatsWithNoPages(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	Config.Source = tmpDir
	_ = clearPagesCache(nil)

	stats := calculateStats(context.Background())

	if stats.TotalPages != 0 {
		t.Errorf("TotalPages = %d, want 0", stats.TotalPages)
	}

	if stats.TotalWords != 0 {
		t.Errorf("TotalWords = %d, want 0", stats.TotalWords)
	}

	if stats.AvgWordsPerPage != 0 {
		t.Errorf("AvgWordsPerPage = %d, want 0", stats.AvgWordsPerPage)
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{
			name:  "simple sentence",
			input: "Hello world this is a test",
			want:  6,
		},
		{
			name:  "empty string",
			input: "",
			want:  0,
		},
		{
			name:  "multiple spaces",
			input: "one  two   three",
			want:  3,
		},
		{
			name:  "newlines and spaces",
			input: "one\ntwo\n\nthree",
			want:  3,
		},
		{
			name:  "markdown content",
			input: "# Title\n\nThis is **bold** and *italic* text.",
			want:  8, // Fields includes markdown symbols as separate words
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := countWords(tc.input)
			if got != tc.want {
				t.Errorf("countWords(%q) = %d, want %d", tc.input, got, tc.want)
			}
		})
	}
}
