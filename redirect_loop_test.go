package xlog

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestIndexRedirectLoop verifies that accessing root (/) doesn't cause an infinite redirect loop
// when a custom index page is configured. This is a regression test to ensure proper redirect
// behavior when using custom index pages (e.g., -index=README).
//
// Expected behavior:
// - "/" redirects to "/README" (one redirect).
// - "/README" serves content (200 OK).
// - No infinite redirect loop.
func TestIndexRedirectLoop(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	// Create a README.md file
	if err := os.WriteFile("README.md", []byte("# README Content\n\nThis is the index page."), 0600); err != nil {
		t.Fatalf("Failed to create README.md: %v", err)
	}

	// Save original config and router state
	oldIndex := Config.Index
	oldReadonly := Config.Readonly
	defer func() {
		Config.Index = oldIndex
		Config.Readonly = oldReadonly
	}()

	// Set custom index page
	Config.Index = "README"
	Config.Readonly = true // Prevent file watching issues in tests

	// Test 1: Accessing root (/) should redirect to /README
	t.Run("root_redirects_to_index", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", http.NoBody)
		w := httptest.NewRecorder()

		output := rootHandler(req)
		output.ServeHTTP(w, req)

		if w.Code != http.StatusFound && w.Code != http.StatusMovedPermanently {
			t.Errorf("Expected redirect from /, got status %d", w.Code)
		}

		location := w.Header().Get("Location")
		expectedLocation := "/README"
		if location != expectedLocation {
			t.Errorf("Expected redirect to %s, got %s", expectedLocation, location)
		}
	})

	// Test 2: Accessing /README directly should NOT redirect (this is the bug test)
	t.Run("readme_no_redirect_loop", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/README", http.NoBody)
		req.SetPathValue("page", "README")
		w := httptest.NewRecorder()

		output := getPageHandler(req)
		output.ServeHTTP(w, req)

		// This should return 200 OK, not another redirect
		if w.Code == http.StatusFound || w.Code == http.StatusMovedPermanently {
			location := w.Header().Get("Location")
			t.Errorf("REDIRECT LOOP BUG: /README redirects to %s (should serve content)", location)
		}

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200 OK for /README, got %d", w.Code)
		}

		// Verify we got actual content
		body := w.Body.String()
		if len(body) == 0 {
			t.Error("Expected non-empty response body for /README")
		}
	})

	// Test 3: Following redirect chain should not exceed 2 hops
	t.Run("no_infinite_redirect_chain", func(t *testing.T) {
		maxRedirects := 10
		redirectCount := 0

		// Start with root
		req := httptest.NewRequest("GET", "/", http.NoBody)
		w := httptest.NewRecorder()
		output := rootHandler(req)
		output.ServeHTTP(w, req)

		for i := 0; i < maxRedirects; i++ {
			if w.Code == http.StatusOK {
				// Success - we got content
				t.Logf("Reached content after %d redirect(s)", redirectCount)
				return
			}

			if w.Code == http.StatusFound || w.Code == http.StatusMovedPermanently {
				redirectCount++
				location := w.Header().Get("Location")
				if location == "" {
					t.Fatal("Redirect with no Location header")
				}

				t.Logf("Redirect %d: %s", redirectCount, location)

				// Follow the redirect - extract page name from location
				pageName := location
				if len(pageName) > 0 && pageName[0] == '/' {
					pageName = pageName[1:]
				}

				req2 := httptest.NewRequest("GET", location, http.NoBody)
				req2.SetPathValue("page", pageName)
				w = httptest.NewRecorder()
				output = getPageHandler(req2)
				output.ServeHTTP(w, req2)
			} else {
				t.Fatalf("Unexpected status code %d", w.Code)
			}
		}

		t.Fatalf("INFINITE REDIRECT LOOP: Exceeded %d redirects", maxRedirects)
	})
}

// TestDefaultIndexNoRedirectLoop verifies the default index.md doesn't have redirect issues.
func TestDefaultIndexNoRedirectLoop(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}()

	// Create default index.md file
	if err := os.WriteFile("index.md", []byte("# Index\n\nDefault index page."), 0600); err != nil {
		t.Fatalf("Failed to create index.md: %v", err)
	}

	// Save original config
	oldIndex := Config.Index
	oldReadonly := Config.Readonly
	defer func() {
		Config.Index = oldIndex
		Config.Readonly = oldReadonly
	}()

	// Use default index
	Config.Index = "index"
	Config.Readonly = true

	// Test: Following redirects should lead to content
	redirectCount := 0
	maxRedirects := 10

	// Start with root
	req := httptest.NewRequest("GET", "/", http.NoBody)
	w := httptest.NewRecorder()
	output := rootHandler(req)
	output.ServeHTTP(w, req)

	for i := 0; i < maxRedirects; i++ {
		if w.Code == http.StatusOK {
			t.Logf("Default index OK after %d redirect(s)", redirectCount)
			return
		}

		if w.Code == http.StatusFound || w.Code == http.StatusMovedPermanently {
			redirectCount++
			location := w.Header().Get("Location")

			// Extract page name from location
			pageName := location
			if len(pageName) > 0 && pageName[0] == '/' {
				pageName = pageName[1:]
			}

			req2 := httptest.NewRequest("GET", location, http.NoBody)
			req2.SetPathValue("page", pageName)
			w = httptest.NewRecorder()
			output = getPageHandler(req2)
			output.ServeHTTP(w, req2)
		} else {
			t.Fatalf("Unexpected status %d", w.Code)
		}
	}

	t.Fatalf("Too many redirects (%d) for default index", redirectCount)
}
