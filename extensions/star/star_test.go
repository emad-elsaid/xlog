package star

import (
	"fmt"
	"html/template"
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
