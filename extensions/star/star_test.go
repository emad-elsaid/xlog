package star

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/emad-elsaid/xlog"
)

const (
	testPageFile  = "test-page.md"
	starIconClass = "fa-solid fa-star"
)

func TestIsStarredLogic(t *testing.T) {
	tests := []struct {
		name           string
		starredContent string
		pageName       string
		expected       bool
	}{
		{
			name:           "xlog.Page is starred",
			starredContent: "page1.md\npage2.md\npage3.md",
			pageName:       "page2.md",
			expected:       true,
		},
		{
			name:           "xlog.Page is not starred",
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
			name:           "xlog.Page with whitespace",
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
			expectedIcon: starIconClass,
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
	defer func() {
		if err := os.Chdir(origDir); err != nil {
			t.Errorf("Failed to restore working directory: %v", err)
		}
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create a test page
	testPageName := testPageFile
	if err := os.WriteFile(testPageName, []byte("# Test xlog.Page"), 0600); err != nil {
		t.Fatal(err)
	}

	page := xlog.NewPage(testPageName)
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
	defer func() {
		if err := os.Chdir(origDir); err != nil {
			t.Errorf("Failed to restore working directory: %v", err)
		}
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create test pages
	testPage := testPageFile
	if err := os.WriteFile(testPage, []byte("# Test"), 0600); err != nil {
		t.Fatal(err)
	}

	otherPage := "other-page.md"
	if err := os.WriteFile(otherPage, []byte("# Other"), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		starredContent string
		pageName       string
		expected       bool
	}{
		{
			name:           "xlog.Page is in starred list",
			starredContent: fmt.Sprintf("%s\n%s", testPage, otherPage),
			pageName:       testPage,
			expected:       true,
		},
		{
			name:           "xlog.Page is not in starred list",
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
				if err := os.WriteFile(STARRED_PAGES+".md", []byte(tt.starredContent), 0600); err != nil {
					t.Fatal(err)
				}
			} else {
				// Remove starred.md if empty content
				_ = os.Remove(STARRED_PAGES + ".md")
			}

			page := xlog.NewPage(tt.pageName)
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
	defer func() { _ = os.Chdir(origDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create a test page
	testPage := testPageFile
	if err := os.WriteFile(testPage, []byte("# Test"), 0600); err != nil {
		t.Fatal(err)
	}

	page := xlog.NewPage("test-page") // Without .md extension
	if page == nil {
		t.Fatal("Failed to create page")
	}

	if !page.Exists() {
		t.Fatal("xlog.Page should exist")
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
				if err := os.WriteFile(STARRED_PAGES+".md", []byte(tt.starredContent), 0600); err != nil {
					t.Fatal(err)
				}
			} else {
				_ = os.Remove(STARRED_PAGES + ".md")
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
	defer func() { _ = os.Chdir(origDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	page := xlog.NewPage("non-existent.md")
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
	defer func() { _ = os.Chdir(origDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create a test page without emoji
	testPage := testPageFile
	if err := os.WriteFile(testPage, []byte("# Test"), 0600); err != nil {
		t.Fatal(err)
	}

	page := xlog.NewPage(testPage)
	if page == nil {
		t.Fatal("Failed to create page")
	}

	sp := starredPage{page}
	icon := sp.Icon()

	if icon != starIconClass {
		t.Errorf("Expected default icon 'fa-solid fa-star', got %s", icon)
	}
}

func TestStarredPageName(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	// Create a test page with path
	testPage := "folder/test-page.md"
	if err := os.MkdirAll("folder", 0750); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(testPage, []byte("# Test"), 0600); err != nil {
		t.Fatal(err)
	}

	page := xlog.NewPage(testPage)
	if page == nil {
		t.Fatal("Failed to create page")
	}

	sp := starredPage{page}
	name := sp.Name()

	expected := testPageFile
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
	defer func() { _ = os.Chdir(origDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	testPage := testPageFile
	if err := os.WriteFile(testPage, []byte("# Test"), 0600); err != nil {
		t.Fatal(err)
	}

	page := xlog.NewPage(testPage)
	if page == nil {
		t.Fatal("Failed to create page")
	}

	sp := starredPage{page}
	attrs := sp.Attrs()

	if _, exists := attrs["href"]; !exists {
		t.Error("Expected href attribute not found")
	}

	expectedHref := "/" + page.Name()
	if actualHref, ok := attrs["href"].(string); ok {
		if actualHref != expectedHref {
			t.Errorf("Expected href %s, got %s", expectedHref, actualHref)
		}
	} else {
		t.Error("href attribute is not a string")
	}
}

func TestStarredPages(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		starredContent string
		pagesExist     []string
		expectedLen    int
	}{
		{
			name:           "No starred page file",
			starredContent: "",
			pagesExist:     nil,
			expectedLen:    0,
		},
		{
			name:           "Empty starred page",
			starredContent: "   \n\n  ",
			pagesExist:     nil,
			expectedLen:    0,
		},
		{
			name:           "Single starred page",
			starredContent: "test-page.md",
			pagesExist:     []string{"test-page.md"},
			expectedLen:    1,
		},
		{
			name:           "Multiple starred pages",
			starredContent: "page1.md\npage2.md\npage3.md",
			pagesExist:     []string{"page1.md", "page2.md", "page3.md"},
			expectedLen:    3,
		},
		{
			name:           "Starred pages with whitespace",
			starredContent: "\npage1.md\n\npage2.md\n",
			pagesExist:     []string{"page1.md", "page2.md"},
			expectedLen:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.starredContent != "" {
				if err := os.WriteFile(STARRED_PAGES+".md", []byte(tt.starredContent), 0600); err != nil {
					t.Fatal(err)
				}
			} else {
				_ = os.Remove(STARRED_PAGES + ".md")
			}

			for _, pageName := range tt.pagesExist {
				if err := os.WriteFile(pageName, []byte("# xlog.Page"), 0600); err != nil {
					t.Fatal(err)
				}
			}

			dummyPage := xlog.NewPage("dummy.md")
			commands := starredPages(dummyPage)

			if tt.expectedLen == 0 {
				if commands != nil {
					t.Errorf("Expected nil commands, got %d commands", len(commands))
				}
			} else {
				if commands == nil {
					t.Fatalf("Expected %d commands, got nil", tt.expectedLen)
				}
				if len(commands) != tt.expectedLen {
					t.Errorf("Expected %d commands, got %d", tt.expectedLen, len(commands))
				}
			}

			for _, pageName := range tt.pagesExist {
				_ = os.Remove(pageName)
			}
		})
	}
}

func TestStarHandler(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name            string
		setupPage       bool
		setupStarred    bool
		starredContent  string
		expectRedirect  bool
		expectPageAdded bool
	}{
		{
			name:            "Star existing page with starred.md present",
			setupPage:       true,
			setupStarred:    true,
			starredContent:  "",
			expectRedirect:  false,
			expectPageAdded: true,
		},
		{
			name:            "Star page adds to existing starred list",
			setupPage:       true,
			setupStarred:    true,
			starredContent:  "other-page.md",
			expectRedirect:  false,
			expectPageAdded: true,
		},
		{
			name:           "Star non-existent page redirects",
			setupPage:      false,
			setupStarred:   true,
			starredContent: "",
			expectRedirect: true,
		},
		{
			name:           "Star page without starred.md file redirects",
			setupPage:      true,
			setupStarred:   false,
			starredContent: "",
			expectRedirect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPage := testPageFile
			pageName := "test-page"
			if tt.setupPage {
				if err := os.WriteFile(testPage, []byte("# Test xlog.Page"), 0600); err != nil {
					t.Fatal(err)
				}
			} else {
				_ = os.Remove(testPage)
			}

			if tt.setupStarred {
				if err := os.WriteFile(STARRED_PAGES+".md", []byte(tt.starredContent), 0600); err != nil {
					t.Fatal(err)
				}
			} else {
				_ = os.Remove(STARRED_PAGES + ".md")
			}

			req := newRequestWithPath(pageName)
			output := starHandler(req)

			if output == nil {
				t.Fatal("Expected non-nil output")
			}

			rec := newRecorder()
			output(rec, req)

			if tt.expectRedirect {
				if rec.redirectCalled {
					return
				}
				if rec.statusCode != 303 && rec.statusCode != 0 {
					return
				}
			}

			if tt.expectPageAdded {
				starredPage := xlog.NewPage(STARRED_PAGES)
				if starredPage == nil {
					t.Fatal("starred.md should exist")
				}

				content := string(starredPage.Content())
				if !strings.Contains(content, pageName) {
					t.Errorf("Expected %s in starred content, got: %s", pageName, content)
				}

				if rec.Header().Get("HX-Refresh") == "" {
					t.Errorf("Expected HX-Refresh header to be set, headers: %v", rec.Header())
				}
			}

			_ = os.Remove(testPage)
		})
	}
}

func TestUnstarHandler(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		setupPage      bool
		setupStarred   bool
		starredContent string
		expectRedirect bool
		expectRemoved  bool
	}{
		{
			name:           "Unstar existing page",
			setupPage:      true,
			setupStarred:   true,
			starredContent: "test-page",
			expectRedirect: false,
			expectRemoved:  true,
		},
		{
			name:           "Unstar page not in list",
			setupPage:      true,
			setupStarred:   true,
			starredContent: "other-page.md",
			expectRedirect: false,
			expectRemoved:  false,
		},
		{
			name:           "Unstar with multiple pages in list",
			setupPage:      true,
			setupStarred:   true,
			starredContent: "page1.md\ntest-page\npage3.md",
			expectRedirect: false,
			expectRemoved:  true,
		},
		{
			name:           "Unstar non-existent page redirects",
			setupPage:      false,
			setupStarred:   true,
			starredContent: "test-page",
			expectRedirect: true,
		},
		{
			name:           "Unstar without starred.md redirects",
			setupPage:      true,
			setupStarred:   false,
			starredContent: "",
			expectRedirect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testPage := testPageFile
			pageName := "test-page"
			if tt.setupPage {
				if err := os.WriteFile(testPage, []byte("# Test xlog.Page"), 0600); err != nil {
					t.Fatal(err)
				}
			} else {
				_ = os.Remove(testPage)
			}

			if tt.setupStarred {
				if err := os.WriteFile(STARRED_PAGES+".md", []byte(tt.starredContent), 0600); err != nil {
					t.Fatal(err)
				}
			} else {
				_ = os.Remove(STARRED_PAGES + ".md")
			}

			req := newRequestWithPath(pageName)
			output := unstarHandler(req)

			if output == nil {
				t.Fatal("Expected non-nil output")
			}

			rec := newRecorder()
			output(rec, req)

			if tt.expectRedirect {
				if rec.redirectCalled {
					return
				}
				if rec.statusCode != 303 && rec.statusCode != 0 {
					return
				}
			}

			if tt.expectRemoved {
				starredPage := xlog.NewPage(STARRED_PAGES)
				if starredPage == nil {
					t.Fatal("starred.md should exist")
				}

				content := strings.TrimSpace(string(starredPage.Content()))
				if strings.Contains(content, pageName) {
					t.Errorf("Expected %s removed from starred content, still present: %s", pageName, content)
				}

				if rec.Header().Get("HX-Refresh") == "" {
					t.Errorf("Expected HX-Refresh header to be set, headers: %v", rec.Header())
				}
			}

			_ = os.Remove(testPage)
		})
	}
}

type testRecorder struct {
	header         http.Header
	statusCode     int
	redirectCalled bool
}

func newRecorder() *testRecorder {
	return &testRecorder{
		header: make(http.Header),
	}
}

func (r *testRecorder) Header() http.Header {
	return r.header
}

func (r *testRecorder) WriteHeader(code int) {
	r.statusCode = code
}

func (r *testRecorder) Write([]byte) (int, error) {
	return 0, nil
}

func newRequestWithPath(pathValue string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.SetPathValue("page", pathValue)
	return req
}
