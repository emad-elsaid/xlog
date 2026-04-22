package sitemap

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestSitemapExtensionName(t *testing.T) {
	ext := Sitemap{}
	if ext.Name() != "sitemap" {
		t.Errorf("Expected name 'sitemap', got '%s'", ext.Name())
	}
}

func TestSitemapHandler(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Create test pages
	testPages := []string{"home", "about", "contact"}
	for _, pageName := range testPages {
		filename := pageName + ".md"
		if err := os.WriteFile(filename, []byte("# "+pageName), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Set domain for sitemap
	oldDomain := SITEMAP_DOMAIN
	SITEMAP_DOMAIN = "example.com"
	defer func() { SITEMAP_DOMAIN = oldDomain }()

	// Create request and response
	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	w := httptest.NewRecorder()

	// Call handler
	output := handler(req)
	output(w, req)

	// Verify status
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	// Verify output content
	body := w.Body.String()

	// Check XML declaration
	if !strings.Contains(body, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("Expected XML declaration in sitemap")
	}

	// Check urlset opening tag
	if !strings.Contains(body, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`) {
		t.Error("Expected urlset opening tag in sitemap")
	}

	// Check urlset closing tag
	if !strings.Contains(body, `</urlset>`) {
		t.Error("Expected urlset closing tag in sitemap")
	}

	// Check each page URL is present
	for _, pageName := range testPages {
		expectedURL := "<loc>https://example.com/" + pageName + "</loc>"
		if !strings.Contains(body, expectedURL) {
			t.Errorf("Expected URL for page '%s' in sitemap", pageName)
		}
	}

	// Check URL format
	if !strings.Contains(body, "<url>") || !strings.Contains(body, "</url>") {
		t.Error("Expected URL tags in sitemap")
	}
}

func TestSitemapHandlerEmptyDirectory(t *testing.T) {
	// Setup empty test environment
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Set domain
	oldDomain := SITEMAP_DOMAIN
	SITEMAP_DOMAIN = "example.com"
	defer func() { SITEMAP_DOMAIN = oldDomain }()

	// Create request and response
	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	w := httptest.NewRecorder()

	// Call handler - should not crash with empty directory
	output := handler(req)
	output(w, req)

	body := w.Body.String()

	// Should still have valid XML structure
	if !strings.Contains(body, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("Expected XML declaration even with no pages")
	}

	if !strings.Contains(body, `<urlset`) {
		t.Error("Expected urlset tag even with no pages")
	}

	if !strings.Contains(body, `</urlset>`) {
		t.Error("Expected closing urlset tag even with no pages")
	}
}

func TestSitemapHandlerDifferentDomain(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Create test page
	if err := os.WriteFile("test.md", []byte("# test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Set different domain
	oldDomain := SITEMAP_DOMAIN
	SITEMAP_DOMAIN = "different.org"
	defer func() { SITEMAP_DOMAIN = oldDomain }()

	// Create request and response
	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	w := httptest.NewRecorder()

	// Call handler
	output := handler(req)
	output(w, req)

	body := w.Body.String()

	// Should use the configured domain
	if !strings.Contains(body, "https://different.org/") {
		t.Error("Expected configured domain in sitemap URLs")
	}

	// Should not contain the default example.com
	if strings.Contains(body, "example.com") {
		t.Error("Should not contain default domain when custom domain is set")
	}
}

func TestSitemapHandlerHTTPS(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}

	// Create a test page
	if err := os.WriteFile("secure.md", []byte("# secure page"), 0644); err != nil {
		t.Fatal(err)
	}

	// Set domain
	oldDomain := SITEMAP_DOMAIN
	SITEMAP_DOMAIN = "secure.com"
	defer func() { SITEMAP_DOMAIN = oldDomain }()

	// Create request and response
	req := httptest.NewRequest(http.MethodGet, "/sitemap.xml", nil)
	w := httptest.NewRecorder()

	// Call handler
	output := handler(req)
	output(w, req)

	body := w.Body.String()

	// Verify HTTPS is used
	if !strings.Contains(body, "https://") {
		t.Error("Expected HTTPS protocol in sitemap URLs")
	}

	// Verify no HTTP (non-secure)
	if strings.Contains(body, "http://") && !strings.Contains(body, "https://") {
		t.Error("Should not use HTTP protocol, only HTTPS")
	}
}
