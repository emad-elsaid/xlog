package sitemap

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSitemapExtensionName(t *testing.T) {
	ext := Sitemap{}
	if ext.Name() != "sitemap" {
		t.Errorf("Expected extension name 'sitemap', got '%s'", ext.Name())
	}
}

func TestSitemapHandler(t *testing.T) {
	// Create temporary directory for test pages
	tmpDir := t.TempDir()
	
	// Create test pages
	testPages := []string{"index.md", "about.md", "blog/post1.md"}
	for _, page := range testPages {
		pageDir := filepath.Dir(filepath.Join(tmpDir, page))
		if err := os.MkdirAll(pageDir, 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(filepath.Join(tmpDir, page), []byte("# Test"), 0644); err != nil {
			t.Fatalf("Failed to create test page: %v", err)
		}
	}

	// Store original dir and change to tmpDir
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Set sitemap domain
	originalDomain := SITEMAP_DOMAIN
	SITEMAP_DOMAIN = "example.com"
	defer func() { SITEMAP_DOMAIN = originalDomain }()

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	w := httptest.NewRecorder()

	// Call handler
	result := handler(req)
	result(w, req)

	// Get response body
	body := w.Body.String()

	// Verify XML declaration
	if !strings.Contains(body, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("Expected XML declaration in output")
	}

	// Verify urlset tag
	if !strings.Contains(body, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`) {
		t.Error("Expected urlset opening tag")
	}

	if !strings.Contains(body, `</urlset>`) {
		t.Error("Expected urlset closing tag")
	}

	// Verify domain is used
	if !strings.Contains(body, "https://example.com/") {
		t.Errorf("Expected domain 'example.com' in URLs, got: %s", body)
	}

	// Verify at least one URL entry exists
	if !strings.Contains(body, "<url>") || !strings.Contains(body, "</url>") {
		t.Error("Expected at least one URL entry")
	}

	// Verify loc tags exist
	if !strings.Contains(body, "<loc>") || !strings.Contains(body, "</loc>") {
		t.Error("Expected loc tags in URL entries")
	}
}

func TestSitemapHandlerEmptyDirectory(t *testing.T) {
	// Create empty temporary directory
	tmpDir := t.TempDir()

	// Store original dir and change to tmpDir
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Set sitemap domain
	originalDomain := SITEMAP_DOMAIN
	SITEMAP_DOMAIN = "test.com"
	defer func() { SITEMAP_DOMAIN = originalDomain }()

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	w := httptest.NewRecorder()

	// Call handler
	result := handler(req)
	result(w, req)
	body := w.Body.String()

	// Should still have valid XML structure
	if !strings.Contains(body, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("Expected XML declaration even with no pages")
	}

	if !strings.Contains(body, `<urlset`) {
		t.Error("Expected urlset tag even with no pages")
	}
}

func TestSitemapHandlerURLEncoding(t *testing.T) {
	// Create temporary directory for test pages
	tmpDir := t.TempDir()
	
	// Store original dir and change to tmpDir
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Create test page in a subdirectory (will have / that needs encoding)
	if err := os.MkdirAll("blog", 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	if err := os.WriteFile("blog/post1.md", []byte("# Test Post"), 0644); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}

	// Set sitemap domain
	originalDomain := SITEMAP_DOMAIN
	SITEMAP_DOMAIN = "example.com"
	defer func() { SITEMAP_DOMAIN = originalDomain }()

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	w := httptest.NewRecorder()

	// Call handler
	result := handler(req)
	result(w, req)
	body := w.Body.String()

	// Verify URL encoding - forward slash should be encoded as %2F
	if !strings.Contains(body, "blog%2Fpost1") {
		t.Logf("Output: %s", body)
		t.Error("Expected URL-encoded path with / as %2F")
	}

	// Verify raw forward slash doesn't appear in the path part of URL
	// (it's OK in the domain https://example.com/)
	if strings.Contains(body, "example.com/blog/post1") {
		t.Error("Expected forward slash to be URL encoded in page path")
	}
}
