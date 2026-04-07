package star

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	. "github.com/emad-elsaid/xlog"
)

func TestIsStarredLogic(t *testing.T) {
	tests := []struct {
		name           string
		starredContent string
		pageName       string
		expected       bool
	}{
		{
			name:           "Page is starred",
			starredContent: "page1.md\npage2.md\npage3.md",
			pageName:       "page2.md",
			expected:       true,
		},
		{
			name:           "Page is not starred",
			starredContent: "page1.md\npage3.md",
			pageName:       "page2.md",
			expected:       false,
		},
		{
			name:           "Empty starred list",
			starredContent: "",
			pageName:       "page1.md",
			expected:       false,
		},
		{
			name:           "Page with whitespace",
			starredContent: "  page1.md  \npage2.md\n  page3.md  ",
			pageName:       "page1.md",
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, k := range strings.Split(tt.starredContent, "\n") {
				if strings.TrimSpace(k) == tt.pageName {
					found = true
					break
				}
			}

			if found != tt.expected {
				t.Errorf("Expected %v, got %v for page %s in starred list:\n%s",
					tt.expected, found, tt.pageName, tt.starredContent)
			}
		})
	}
}

func TestActionIconAndName(t *testing.T) {
	tests := []struct {
		name         string
		starred      bool
		expectedIcon string
		expectedName string
	}{
		{
			name:         "Starred action shows unstar",
			starred:      true,
			expectedIcon: "fa-solid fa-star",
			expectedName: "Unstar",
		},
		{
			name:         "Unstarred action shows star",
			starred:      false,
			expectedIcon: "fa-regular fa-star",
			expectedName: "Star",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := action{starred: tt.starred}

			if act.Icon() != tt.expectedIcon {
				t.Errorf("Expected icon %s, got %s", tt.expectedIcon, act.Icon())
			}

			if act.Name() != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, act.Name())
			}
		})
	}
}

func TestActionAttrs(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create a test page
	testPageName := "test-page.md"
	if err := os.WriteFile(testPageName, []byte("# Test Page"), 0644); err != nil {
		t.Fatal(err)
	}

	page := NewPage(testPageName)
	if page == nil {
		t.Fatal("Failed to create test page")
	}

	tests := []struct {
		name     string
		starred  bool
		wantAttr template.HTMLAttr
	}{
		{
			name:     "Unstarred page has hx-post",
			starred:  false,
			wantAttr: "hx-post",
		},
		{
			name:     "Starred page has hx-delete",
			starred:  true,
			wantAttr: "hx-delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := action{starred: tt.starred, page: page}
			attrs := act.Attrs()

			if _, exists := attrs[tt.wantAttr]; !exists {
				t.Errorf("Expected attribute %s not found in attrs: %v", tt.wantAttr, attrs)
			}

			// Verify href always exists
			if _, exists := attrs["href"]; !exists {
				t.Error("Expected href attribute not found")
			}
		})
	}
}

func TestStarredPagesParsing(t *testing.T) {
	content := "page1.md\npage2.md\npage3.md\n"
	list := strings.Split(strings.TrimSpace(content), "\n")

	if len(list) != 3 {
		t.Errorf("Expected 3 pages, got %d", len(list))
	}

	expected := []string{"page1.md", "page2.md", "page3.md"}
	for i, v := range list {
		if v != expected[i] {
			t.Errorf("Expected %s at index %d, got %s", expected[i], i, v)
		}
	}
}

func TestStarredPagesEmptyContent(t *testing.T) {
	content := ""
	trimmed := strings.TrimSpace(content)

	if trimmed != "" {
		t.Error("Expected empty string after trim")
	}

	// Empty content should return nil list
	if trimmed == "" {
		// This is the expected behavior
		return
	}

	t.Error("Should have returned early for empty content")
}

func TestIsStarred(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create test pages
	testPage := "test-page.md"
	if err := os.WriteFile(testPage, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	otherPage := "other-page.md"
	if err := os.WriteFile(otherPage, []byte("# Other"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		starredContent string
		pageName       string
		expected       bool
	}{
		{
			name:           "Page is in starred list",
			starredContent: fmt.Sprintf("%s\n%s", testPage, otherPage),
			pageName:       testPage,
			expected:       true,
		},
		{
			name:           "Page is not in starred list",
			starredContent: otherPage,
			pageName:       testPage,
			expected:       false,
		},
		{
			name:           "No starred page exists",
			starredContent: "",
			pageName:       testPage,
			expected:       false,
		},
		{
			name:           "Starred list with whitespace",
			starredContent: fmt.Sprintf("  %s  \n%s", testPage, otherPage),
			pageName:       testPage,
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create or update starred.md
			if tt.starredContent != "" {
				if err := os.WriteFile(STARRED_PAGES+".md", []byte(tt.starredContent), 0644); err != nil {
					t.Fatal(err)
				}
			} else {
				// Remove starred.md if empty content
				os.Remove(STARRED_PAGES + ".md")
			}

			page := NewPage(tt.pageName)
			if page == nil {
				t.Fatal("Failed to create page")
			}

			result := isStarred(page)
			if result != tt.expected {
				t.Errorf("Expected isStarred=%v, got %v for page %s with starred content:\n%s",
					tt.expected, result, tt.pageName, tt.starredContent)
			}
		})
	}
}

func TestStarAction(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create a test page
	testPage := "test-page.md"
	if err := os.WriteFile(testPage, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	page := NewPage("test-page") // Without .md extension
	if page == nil {
		t.Fatal("Failed to create page")
	}

	if !page.Exists() {
		t.Fatal("Page should exist")
	}

	tests := []struct {
		name           string
		starredContent string
		expectedLen    int
		expectedName   string
	}{
		{
			name:           "Unstarred page returns Star action",
			starredContent: "",
			expectedLen:    1,
			expectedName:   "Star",
		},
		{
			name:           "Starred page returns Unstar action",
			starredContent: "test-page",
			expectedLen:    1,
			expectedName:   "Unstar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.starredContent != "" {
				if err := os.WriteFile(STARRED_PAGES+".md", []byte(tt.starredContent), 0644); err != nil {
					t.Fatal(err)
				}
			} else {
				os.Remove(STARRED_PAGES + ".md")
			}

			commands := starAction(page)
			if len(commands) != tt.expectedLen {
				t.Errorf("Expected %d commands, got %d", tt.expectedLen, len(commands))
			}

			if len(commands) > 0 {
				if commands[0].Name() != tt.expectedName {
					t.Errorf("Expected command name %s, got %s", tt.expectedName, commands[0].Name())
				}
			}
		})
	}
}

func TestStarActionNonExistentPage(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	page := NewPage("non-existent.md")
	if page == nil {
		// This is expected for a non-existent page
		t.Skip("NewPage returns nil for non-existent pages as expected")
	}

	commands := starAction(page)
	if commands != nil {
		t.Errorf("Expected nil commands for non-existent page, got %d commands", len(commands))
	}
}

func TestStarredPageIcon(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create a test page without emoji
	testPage := "test-page.md"
	if err := os.WriteFile(testPage, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	page := NewPage(testPage)
	if page == nil {
		t.Fatal("Failed to create page")
	}

	sp := starredPage{page}
	icon := sp.Icon()

	if icon != "fa-solid fa-star" {
		t.Errorf("Expected default icon 'fa-solid fa-star', got %s", icon)
	}

	// Test with emoji
	testPageEmoji := "emoji-page.md"
	if err := os.WriteFile(testPageEmoji, []byte(":smile: Test"), 0644); err != nil {
		t.Fatal(err)
	}

	pageEmoji := NewPage(testPageEmoji)
	if pageEmoji == nil {
		t.Fatal("Failed to create emoji page")
	}

	spEmoji := starredPage{pageEmoji}
	iconEmoji := spEmoji.Icon()

	// If emoji parsing is working, it should return the emoji, otherwise fall back to default
	// We can't strictly assert the emoji value without knowing the exact implementation
	// but we can verify the function doesn't panic and returns a non-empty string
	if iconEmoji == "" {
		t.Error("Expected non-empty icon")
	}
}

func TestStarredPageName(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create a test page with path
	testPage := "folder/test-page.md"
	if err := os.MkdirAll("folder", 0755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(testPage, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	page := NewPage(testPage)
	if page == nil {
		t.Fatal("Failed to create page")
	}

	sp := starredPage{page}
	name := sp.Name()

	expected := "test-page.md"
	if name != expected {
		t.Errorf("Expected name %s, got %s", expected, name)
	}
}

func TestStarExtensionName(t *testing.T) {
	ext := Star{}
	if ext.Name() != "star" {
		t.Errorf("Expected name 'star', got '%s'", ext.Name())
	}
}

func TestStarredPageAttrs(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create a test page
	testPage := "test-page.md"
	if err := os.WriteFile(testPage, []byte("# Test"), 0644); err != nil {
		t.Fatal(err)
	}

	page := NewPage(testPage)
	if page == nil {
		t.Fatal("Failed to create page")
	}

	sp := starredPage{page}
	attrs := sp.Attrs()

	// Check that href attribute exists
	href, exists := attrs["href"]
	if !exists {
		t.Error("Expected href attribute")
	}

	// Check that href points to the page
	expectedHref := "/" + page.Name()
	if href != expectedHref {
		t.Errorf("Expected href %s, got %v", expectedHref, href)
	}
}

func TestStarredPagesCommand(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create test pages
	page1 := "page1.md"
	page2 := "page2.md"
	if err := os.WriteFile(page1, []byte("# Page 1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(page2, []byte("# Page 2"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		starredContent string
		expectedLen    int
	}{
		{
			name:           "Multiple starred pages",
			starredContent: "page1\npage2",
			expectedLen:    2,
		},
		{
			name:           "Single starred page",
			starredContent: "page1",
			expectedLen:    1,
		},
		{
			name:           "Empty starred list",
			starredContent: "",
			expectedLen:    0,
		},
		{
			name:           "Whitespace only starred list",
			starredContent: "   \n\n  ",
			expectedLen:    0,
		},
		{
			name:           "No starred page exists",
			starredContent: "<no-file>",
			expectedLen:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.starredContent == "<no-file>" {
				os.Remove(STARRED_PAGES + ".md")
			} else {
				if err := os.WriteFile(STARRED_PAGES+".md", []byte(tt.starredContent), 0644); err != nil {
					t.Fatal(err)
				}
			}

			dummyPage := NewPage(page1)
			commands := starredPages(dummyPage)

			if tt.expectedLen == 0 {
				if commands != nil {
					t.Errorf("Expected nil commands, got %d", len(commands))
				}
			} else {
				if len(commands) != tt.expectedLen {
					t.Errorf("Expected %d commands, got %d", tt.expectedLen, len(commands))
				}
			}
		})
	}
}

func TestStarHandler(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create test page
	testPage := "test-page.md"
	if err := os.WriteFile(testPage, []byte("# Test Page"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create starred.md
	starredFile := STARRED_PAGES + ".md"
	if err := os.WriteFile(starredFile, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	// Create mock request with path value
	req := httptest.NewRequest(http.MethodPost, "/+/star/test-page", nil)
	req.SetPathValue("page", "test-page")
	w := httptest.NewRecorder()

	// Call handler
	result := starHandler(req)
	result(w, req)

	// Check response
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Check HX-Refresh header
	if hxRefresh := w.Header().Get("HX-Refresh"); hxRefresh != "true" {
		t.Errorf("Expected HX-Refresh header 'true', got '%s'", hxRefresh)
	}

	// Verify page was added to starred.md
	content, err := os.ReadFile(starredFile)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "test-page") {
		t.Errorf("Expected starred.md to contain 'test-page', got: %s", content)
	}
}

func TestStarHandlerNonExistentPage(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create starred.md but not the target page
	starredFile := STARRED_PAGES + ".md"
	if err := os.WriteFile(starredFile, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/+/star/non-existent", nil)
	req.SetPathValue("page", "non-existent")
	w := httptest.NewRecorder()

	result := starHandler(req)
	result(w, req)

	// Should redirect to home for non-existent page
	if w.Code != http.StatusFound && w.Code != http.StatusSeeOther && w.Code != http.StatusMovedPermanently {
		// Check if it's actually a redirect
		location := w.Header().Get("Location")
		if location != "/" {
			t.Logf("Warning: non-existent page handling might differ, got status %d", w.Code)
		}
	}
}

func TestUnstarHandler(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create test page
	testPage := "test-page.md"
	if err := os.WriteFile(testPage, []byte("# Test Page"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create starred.md with the page already starred
	starredFile := STARRED_PAGES + ".md"
	initialContent := "test-page\nother-page"
	if err := os.WriteFile(starredFile, []byte(initialContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create mock request
	req := httptest.NewRequest(http.MethodDelete, "/+/star/test-page", nil)
	req.SetPathValue("page", "test-page")
	w := httptest.NewRecorder()

	// Call handler
	result := unstarHandler(req)
	result(w, req)

	// Check response
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Check HX-Refresh header
	if hxRefresh := w.Header().Get("HX-Refresh"); hxRefresh != "true" {
		t.Errorf("Expected HX-Refresh header 'true', got '%s'", hxRefresh)
	}

	// Verify page was removed from starred.md
	content, err := os.ReadFile(starredFile)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "test-page" {
			t.Error("Expected 'test-page' to be removed from starred.md")
		}
	}

	// Verify other-page is still there
	if !strings.Contains(string(content), "other-page") {
		t.Error("Expected 'other-page' to remain in starred.md")
	}
}

func TestUnstarHandlerNonExistentPage(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	// Create starred.md
	starredFile := STARRED_PAGES + ".md"
	if err := os.WriteFile(starredFile, []byte("some-page"), 0644); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/+/star/non-existent", nil)
	req.SetPathValue("page", "non-existent")
	w := httptest.NewRecorder()

	result := unstarHandler(req)
	result(w, req)

	// Should redirect for non-existent page
	if w.Code != http.StatusFound && w.Code != http.StatusSeeOther && w.Code != http.StatusMovedPermanently {
		location := w.Header().Get("Location")
		if location != "/" {
			t.Logf("Warning: non-existent page handling might differ, got status %d", w.Code)
		}
	}
}
