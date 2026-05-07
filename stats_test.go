package xlog

import (
	"bytes"
	"context"
	"fmt"
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
				"Total Links:",
				"Orphaned Pages:",
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
	
	originalSources := sources
	defer func() {
		sourcesMutex.Lock()
		sources = originalSources
		sourcesMutex.Unlock()
	}()

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
	
	// Replace global sources with test-specific markdownFS
	sourcesMutex.Lock()
	sources = []PageSource{newMarkdownFS(tmpDir)}
	sourcesMutex.Unlock()
	
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

	// Verify link statistics are calculated
	if stats.TotalLinks < 0 {
		t.Errorf("TotalLinks should not be negative: got %d", stats.TotalLinks)
	}

	if stats.OrphanedPages < 0 {
		t.Errorf("OrphanedPages should not be negative: got %d", stats.OrphanedPages)
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

	if stats.TotalLinks != 0 {
		t.Errorf("TotalLinks = %d, want 0", stats.TotalLinks)
	}

	if stats.OrphanedPages != 0 {
		t.Errorf("OrphanedPages = %d, want 0", stats.OrphanedPages)
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

func TestCalculateStatsWithLinks(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	
	originalSources := sources
	defer func() {
		sourcesMutex.Lock()
		sources = originalSources
		sourcesMutex.Unlock()
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Create a network of linked pages:
	// hub.md <- page1.md, page2.md, page3.md (hub page)
	// page1.md -> hub.md
	// page2.md -> hub.md, page1.md
	// page3.md -> hub.md
	// orphan.md (no incoming links)
	files := map[string]string{
		"hub.md":    "# Hub Page\nCentral page",
		"page1.md":  "# Page 1\nLinks to [[hub]]",
		"page2.md":  "# Page 2\nLinks to [[hub]] and [[page1]]",
		"page3.md":  "# Page 3\nLinks to [[hub]]",
		"orphan.md": "# Orphan\nNo one links to me",
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	Config.Source = tmpDir
	
	// Replace global sources with test-specific markdownFS
	sourcesMutex.Lock()
	sources = []PageSource{newMarkdownFS(tmpDir)}
	sourcesMutex.Unlock()
	
	_ = clearPagesCache(nil)

	stats := calculateStats(context.Background())

	if stats.TotalPages != 5 {
		t.Errorf("TotalPages = %d, want 5", stats.TotalPages)
	}

	if stats.TotalLinks != 4 {
		t.Errorf("TotalLinks = %d, want 4 (page1->hub + page2->hub + page2->page1 + page3->hub)", stats.TotalLinks)
	}

	// orphan.md should have no incoming links
	if stats.OrphanedPages < 1 {
		t.Errorf("OrphanedPages = %d, want at least 1 (orphan.md)", stats.OrphanedPages)
	}

	// hub.md should be identified as a hub page (most incoming links)
	foundHub := false
	for _, hubPage := range stats.HubPages {
		if hubPage == "hub" {
			foundHub = true
			break
		}
	}

	if !foundHub {
		t.Errorf("Expected 'hub' to be in HubPages %v", stats.HubPages)
	}
}

func TestFindHubPages(t *testing.T) {
	tests := []struct {
		name          string
		incomingLinks map[string]int
		topN          int
		want          []string
	}{
		{
			name: "single hub page",
			incomingLinks: map[string]int{
				"hub":   10,
				"page1": 2,
				"page2": 1,
			},
			topN: 3,
			want: []string{"hub", "page1", "page2"},
		},
		{
			name: "limit to top 2",
			incomingLinks: map[string]int{
				"hub":   10,
				"page1": 5,
				"page2": 3,
				"page3": 1,
			},
			topN: 2,
			want: []string{"hub", "page1"},
		},
		{
			name:          "empty links",
			incomingLinks: map[string]int{},
			topN:          3,
			want:          []string{},
		},
		{
			name: "pages with zero links excluded",
			incomingLinks: map[string]int{
				"hub":    5,
				"orphan": 0,
				"page1":  2,
			},
			topN: 3,
			want: []string{"hub", "page1"},
		},
		{
			name: "fewer pages than topN",
			incomingLinks: map[string]int{
				"page1": 3,
			},
			topN: 5,
			want: []string{"page1"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := findHubPages(tc.incomingLinks, tc.topN)

			if len(got) != len(tc.want) {
				t.Errorf("findHubPages() returned %d pages, want %d\nGot: %v\nWant: %v",
					len(got), len(tc.want), got, tc.want)
				return
			}

			// Verify the expected pages are present (order matters for top hubs)
			for i, wantPage := range tc.want {
				if got[i] != wantPage {
					t.Errorf("Position %d: got %q, want %q", i, got[i], wantPage)
				}
			}
		})
	}
}

// BenchmarkFindHubPages_Small benchmarks hub page finding with small dataset.
func BenchmarkFindHubPages_Small(b *testing.B) {
	links := make(map[string]int)
	for i := 0; i < 10; i++ {
		links[fmt.Sprintf("page%d", i)] = i * 3
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = findHubPages(links, 3)
	}
}

// BenchmarkFindHubPages_Medium benchmarks with medium-sized garden (100 pages).
func BenchmarkFindHubPages_Medium(b *testing.B) {
	links := make(map[string]int)
	for i := 0; i < 100; i++ {
		links[fmt.Sprintf("page%d", i)] = (i * 17) % 50 // Varied link counts
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = findHubPages(links, 3)
	}
}

// BenchmarkFindHubPages_Large benchmarks with large garden (1000 pages).
func BenchmarkFindHubPages_Large(b *testing.B) {
	links := make(map[string]int)
	for i := 0; i < 1000; i++ {
		links[fmt.Sprintf("page%d", i)] = (i * 37) % 100 // Varied link counts
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = findHubPages(links, 3)
	}
}
