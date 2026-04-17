package rss

import (
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/emad-elsaid/xlog"
)

func TestRSSExtensionName(t *testing.T) {
	ext := RSS{}
	if ext.Name() != "rss" {
		t.Errorf("Expected extension name 'rss', got '%s'", ext.Name())
	}
}

func setupTestEnv(t *testing.T) func() {
	t.Helper()

	tmpDir := t.TempDir()
	wd, _ := os.Getwd()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Create test pages
	os.WriteFile("page1.md", []byte("# Page 1\nContent 1"), 0644)
	os.WriteFile("page2.md", []byte("# Page 2\nContent 2"), 0644)

	// Save original values
	originalDomain := domain
	originalDescription := description
	originalLimit := limit
	originalSitename := Config.Sitename

	// Set test values
	domain = "example.com"
	description = "Test RSS Feed"
	limit = 30
	Config.Sitename = "Test Site"

	cleanup := func() {
		os.Chdir(wd)
		domain = originalDomain
		description = originalDescription
		limit = originalLimit
		Config.Sitename = originalSitename
	}

	return cleanup
}

func TestFeedXMLStructure(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/+/feed.rss", nil)
	w := httptest.NewRecorder()

	result := feed(req)
	result(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()

	// Parse XML
	var rssData rss
	if err := xml.Unmarshal([]byte(body), &rssData); err != nil {
		t.Fatalf("Failed to parse RSS XML: %v", err)
	}

	// Verify RSS structure
	if rssData.Version != "2.0" {
		t.Errorf("Expected RSS version 2.0, got %s", rssData.Version)
	}

	if rssData.Channel.Title != "Test Site" {
		t.Errorf("Expected title 'Test Site', got '%s'", rssData.Channel.Title)
	}

	if rssData.Channel.Description != "Test RSS Feed" {
		t.Errorf("Expected description 'Test RSS Feed', got '%s'", rssData.Channel.Description)
	}

	if rssData.Channel.Language != "en-US" {
		t.Errorf("Expected language 'en-US', got '%s'", rssData.Channel.Language)
	}

	if !strings.Contains(rssData.Channel.Link, "example.com") {
		t.Errorf("Expected link to contain 'example.com', got '%s'", rssData.Channel.Link)
	}
}

func TestFeedItemsCount(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Set limit to 2
	limit = 2

	// Create a third page
	os.WriteFile("page3.md", []byte("# Page 3\nContent 3"), 0644)

	req := httptest.NewRequest(http.MethodGet, "/+/feed.rss", nil)
	w := httptest.NewRecorder()

	result := feed(req)
	result(w, req)

	body := w.Body.String()

	// Parse XML
	var rssData rss
	if err := xml.Unmarshal([]byte(body), &rssData); err != nil {
		t.Fatalf("Failed to parse RSS XML: %v", err)
	}

	// Should respect the limit of 2
	if len(rssData.Channel.Items) > 2 {
		t.Errorf("Expected at most 2 items (respecting limit), got %d", len(rssData.Channel.Items))
	}
}

func TestFeedItemsSortedByModTime(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create old page first
	os.WriteFile("old.md", []byte("Old content"), 0644)

	// Sleep to ensure different modtimes
	time.Sleep(10 * time.Millisecond)

	// Create new page
	os.WriteFile("new.md", []byte("New content"), 0644)

	req := httptest.NewRequest(http.MethodGet, "/+/feed.rss", nil)
	w := httptest.NewRecorder()

	result := feed(req)
	result(w, req)

	body := w.Body.String()

	// Parse XML
	var rssData rss
	if err := xml.Unmarshal([]byte(body), &rssData); err != nil {
		t.Fatalf("Failed to parse RSS XML: %v", err)
	}

	if len(rssData.Channel.Items) < 2 {
		t.Fatal("Expected at least 2 items")
	}

	// First item should be the newest (most recently modified)
	// The feed is sorted by ModTime descending
	firstPubDate := rssData.Channel.Items[0].PubDate
	secondPubDate := rssData.Channel.Items[1].PubDate

	if !firstPubDate.After(secondPubDate) && !firstPubDate.Equal(secondPubDate) {
		t.Errorf("Expected items to be sorted by modification time (newest first)")
	}
}

func TestFeedItemContent(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create a specific test page
	os.WriteFile("test-page.md", []byte("# Test Content\n\nSome content here"), 0644)

	req := httptest.NewRequest(http.MethodGet, "/+/feed.rss", nil)
	w := httptest.NewRecorder()

	result := feed(req)
	result(w, req)

	body := w.Body.String()

	// Parse XML
	var rssData rss
	if err := xml.Unmarshal([]byte(body), &rssData); err != nil {
		t.Fatalf("Failed to parse RSS XML: %v", err)
	}

	if len(rssData.Channel.Items) < 1 {
		t.Fatal("Expected at least 1 item")
	}

	// Verify items exist and have basic structure
	for _, item := range rssData.Channel.Items {
		if item.Title == "" {
			t.Error("Expected item to have a title")
		}

		if item.GUID == "" {
			t.Error("Expected item to have a GUID")
		}

		if !strings.Contains(item.Link, "example.com") {
			t.Errorf("Expected link to contain domain, got '%s'", item.Link)
		}

		if item.Description == "" {
			t.Error("Expected item to have description")
		}
	}
}

func TestFeedContentType(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/+/feed.rss", nil)
	w := httptest.NewRecorder()

	result := feed(req)
	result(w, req)

	// Check that XML header is present
	body := w.Body.String()
	if !strings.HasPrefix(body, xml.Header) {
		t.Error("Expected XML to start with XML header")
	}

	// Verify it's valid XML
	var rssData rss
	if err := xml.Unmarshal([]byte(body), &rssData); err != nil {
		t.Fatalf("Generated invalid XML: %v", err)
	}
}
