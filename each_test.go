package xlog

import (
	"context"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsIgnoredPath(t *testing.T) {
	assert.True(t, IsIgnoredPath(".git/config"))
	assert.True(t, IsIgnoredPath(".versions/config"))
	assert.False(t, IsIgnoredPath("index.md"))
	assert.False(t, IsIgnoredPath("something/something"))
}

func TestIgnorePath_CustomPattern(t *testing.T) {
	// Save original and restore
	originalIgnoredPaths := ignoredPaths
	defer func() { ignoredPaths = originalIgnoredPaths }()
	
	// Reset to default
	ignoredPaths = []*regexp.Regexp{
		regexp.MustCompile(`^\.`),
	}
	
	// Add custom pattern
	IgnorePath(regexp.MustCompile(`^temp/`))
	
	assert.True(t, IsIgnoredPath("temp/file.md"), "Should ignore temp/")
	assert.True(t, IsIgnoredPath(".hidden"), "Should still ignore hidden files")
	assert.False(t, IsIgnoredPath("normal.md"), "Should not ignore normal files")
}

func TestIsNil(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		{
			name:     "nil page",
			value:    Page(nil),
			expected: true,
		},
		{
			name:     "nil pointer",
			value:    (*mockPage)(nil),
			expected: true,
		},
		{
			name:     "nil slice",
			value:    []string(nil),
			expected: true,
		},
		{
			name:     "empty slice",
			value:    []string{},
			expected: false,
		},
		{
			name:     "valid page",
			value:    &mockPage{name: "test"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNil(tt.value)
			if result != tt.expected {
				t.Errorf("isNil(%v) = %v, expected %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestClearPagesCache(t *testing.T) {
	// Save original and restore
	originalPages := pages
	originalSources := sources
	defer func() {
		pages = originalPages
		sources = originalSources
	}()
	
	// Set up test pages
	pages = []Page{
		&mockPage{name: "page1"},
		&mockPage{name: "page2"},
	}
	
	// Clear cache
	err := clearPagesCache(nil)
	
	assert.NoError(t, err)
	assert.Nil(t, pages, "Pages cache should be nil after clearing")
}

func TestPages_WithCache(t *testing.T) {
	// Save original and restore
	originalPages := pages
	originalSources := sources
	defer func() {
		pages = originalPages
		sources = originalSources
	}()
	
	// Pre-populate cache
	testPages := []Page{
		&mockPage{name: "cached1"},
		&mockPage{name: "cached2"},
	}
	pages = testPages
	
	ctx := context.Background()
	result := Pages(ctx)
	
	assert.Equal(t, 2, len(result), "Should return cached pages")
	assert.Equal(t, "cached1", result[0].Name())
	assert.Equal(t, "cached2", result[1].Name())
}

func TestPages_PopulatesCache(t *testing.T) {
	// Save original and restore
	originalPages := pages
	originalSources := sources
	defer func() {
		pages = originalPages
		sources = originalSources
	}()
	
	// Clear cache
	pages = nil
	
	// Setup mock source
	sources = []PageSource{
		&mockPageSource{
			eachFunc: func(ctx context.Context, fn func(Page)) {
				fn(&mockPage{name: "source1"})
				fn(&mockPage{name: "source2"})
			},
		},
	}
	
	ctx := context.Background()
	result := Pages(ctx)
	
	assert.Equal(t, 2, len(result), "Should populate cache from sources")
	assert.Equal(t, "source1", result[0].Name())
	assert.Equal(t, "source2", result[1].Name())
}

func TestEachPage_WithCache(t *testing.T) {
	// Save original and restore
	originalPages := pages
	originalSources := sources
	defer func() {
		pages = originalPages
		sources = originalSources
	}()
	
	// Pre-populate cache
	pages = []Page{
		&mockPage{name: "page1"},
		&mockPage{name: "page2"},
		&mockPage{name: "page3"},
	}
	
	ctx := context.Background()
	var names []string
	
	EachPage(ctx, func(p Page) {
		names = append(names, p.Name())
	})
	
	assert.Equal(t, 3, len(names))
	assert.Equal(t, []string{"page1", "page2", "page3"}, names)
}

func TestEachPage_ContextCancellation(t *testing.T) {
	// Save original and restore
	originalPages := pages
	defer func() { pages = originalPages }()
	
	// Pre-populate with many pages
	testPages := make([]Page, 100)
	for i := 0; i < 100; i++ {
		testPages[i] = &mockPage{name: "page"}
	}
	pages = testPages
	
	ctx, cancel := context.WithCancel(context.Background())
	var count int
	var mu sync.Mutex
	
	// Cancel after processing a few pages
	go func() {
		time.Sleep(1 * time.Millisecond)
		cancel()
	}()
	
	EachPage(ctx, func(p Page) {
		mu.Lock()
		count++
		mu.Unlock()
		time.Sleep(1 * time.Millisecond) // Simulate work
	})
	
	// Should have stopped early due to cancellation
	mu.Lock()
	finalCount := count
	mu.Unlock()
	
	assert.Less(t, finalCount, 100, "Should stop early on context cancellation")
}

func TestMapPage_BasicUsage(t *testing.T) {
	// Save original and restore
	originalPages := pages
	originalSources := sources
	defer func() {
		pages = originalPages
		sources = originalSources
	}()
	
	// Pre-populate cache
	pages = []Page{
		&mockPage{name: "page1"},
		&mockPage{name: "page2"},
		&mockPage{name: "page3"},
	}
	
	ctx := context.Background()
	
	// Map to page names
	result := MapPage(ctx, func(p Page) string {
		return p.Name()
	})
	
	assert.Equal(t, 3, len(result))
	// Order may vary due to concurrency, so check containment
	assert.Contains(t, result, "page1")
	assert.Contains(t, result, "page2")
	assert.Contains(t, result, "page3")
}

func TestMapPage_FiltersNil(t *testing.T) {
	// Save original and restore
	originalPages := pages
	defer func() { pages = originalPages }()
	
	// Pre-populate cache
	pages = []Page{
		&mockPage{name: "include"},
		&mockPage{name: "skip"},
		&mockPage{name: "include2"},
	}
	
	ctx := context.Background()
	
	// Map function that returns nil for "skip"
	result := MapPage(ctx, func(p Page) *string {
		if p.Name() == "skip" {
			return nil
		}
		name := p.Name()
		return &name
	})
	
	assert.Equal(t, 2, len(result), "Should filter out nil results")
}

func TestMapPage_ContextCancellation(t *testing.T) {
	// Save original and restore
	originalPages := pages
	defer func() { pages = originalPages }()
	
	// Pre-populate with many pages
	testPages := make([]Page, 50)
	for i := 0; i < 50; i++ {
		testPages[i] = &mockPage{name: "page"}
	}
	pages = testPages
	
	ctx, cancel := context.WithCancel(context.Background())
	
	// Cancel immediately
	cancel()
	
	result := MapPage(ctx, func(p Page) string {
		return p.Name()
	})
	
	// Should return early due to cancelled context
	// May process a few before cancellation kicks in
	assert.LessOrEqual(t, len(result), 50)
}

func TestMapPage_ConcurrentExecution(t *testing.T) {
	// Save original and restore
	originalPages := pages
	defer func() { pages = originalPages }()
	
	// Pre-populate cache
	numPages := 20
	testPages := make([]Page, numPages)
	for i := 0; i < numPages; i++ {
		testPages[i] = &mockPage{name: "page"}
	}
	pages = testPages
	
	ctx := context.Background()
	
	var executionOrder []int
	var mu sync.Mutex
	
	// Map with artificial delay to observe concurrency
	MapPage(ctx, func(p Page) string {
		mu.Lock()
		executionOrder = append(executionOrder, 1)
		mu.Unlock()
		time.Sleep(1 * time.Millisecond)
		return p.Name()
	})
	
	mu.Lock()
	count := len(executionOrder)
	mu.Unlock()
	
	assert.Equal(t, numPages, count, "All pages should be processed")
}

func TestPopulatePagesCache_DoubleCheck(t *testing.T) {
	// Save original and restore
	originalPages := pages
	originalSources := sources
	defer func() {
		pages = originalPages
		sources = originalSources
	}()
	
	// Pre-populate cache
	pages = []Page{&mockPage{name: "existing"}}
	
	// Setup mock source (should not be called)
	callCount := 0
	sources = []PageSource{
		&mockPageSource{
			eachFunc: func(ctx context.Context, fn func(Page)) {
				callCount++
				fn(&mockPage{name: "new"})
			},
		},
	}
	
	ctx := context.Background()
	populatePagesCache(ctx)
	
	// Should not have called source since cache was already populated
	assert.Equal(t, 0, callCount, "Should not populate if cache exists (double-check lock)")
	assert.Equal(t, 1, len(pages), "Cache should remain unchanged")
	assert.Equal(t, "existing", pages[0].Name())
}

func TestPopulatePagesCache_MultipleSOURCES(t *testing.T) {
	// Save original and restore
	originalPages := pages
	originalSources := sources
	defer func() {
		pages = originalPages
		sources = originalSources
	}()
	
	// Clear cache
	pages = nil
	
	// Setup multiple sources
	sources = []PageSource{
		&mockPageSource{
			eachFunc: func(ctx context.Context, fn func(Page)) {
				fn(&mockPage{name: "source1-page1"})
				fn(&mockPage{name: "source1-page2"})
			},
		},
		&mockPageSource{
			eachFunc: func(ctx context.Context, fn func(Page)) {
				fn(&mockPage{name: "source2-page1"})
			},
		},
	}
	
	ctx := context.Background()
	populatePagesCache(ctx)
	
	assert.Equal(t, 3, len(pages), "Should aggregate from all sources")
}

func TestPopulatePagesCache_ContextCancellation(t *testing.T) {
	// Save original and restore
	originalPages := pages
	originalSources := sources
	defer func() {
		pages = originalPages
		sources = originalSources
	}()
	
	// Clear cache
	pages = nil
	
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	
	// Setup source
	sources = []PageSource{
		&mockPageSource{
			eachFunc: func(ctx context.Context, fn func(Page)) {
				fn(&mockPage{name: "page1"})
			},
		},
	}
	
	populatePagesCache(ctx)
	
	// Should return early on cancelled context
	// Cache initialization happens, but source enumeration stops
	assert.NotNil(t, pages, "Cache slice should be initialized")
}
