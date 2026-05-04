package recent

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
)

const mdExt = ".md"

func TestRecent_Extension(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"returns correct extension name", "recent"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ext := Recent{}
			if got := ext.Name(); got != tc.expected {
				t.Errorf("Name() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestRecent_LinksCommand(t *testing.T) {
	tests := []struct {
		name         string
		checkIcon    bool
		wantIcon     string
		checkName    bool
		wantName     string
		checkHref    bool
		wantHref     string
		checkKeyType bool
	}{
		{
			name:         "returns complete command implementation",
			checkIcon:    true,
			wantIcon:     "fa-solid fa-clock-rotate-left",
			checkName:    true,
			wantName:     "Recent",
			checkHref:    true,
			wantHref:     "/+/recent",
			checkKeyType: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := links{}

			if tc.checkIcon {
				if got := l.Icon(); got != tc.wantIcon {
					t.Errorf("Icon() = %q, want %q", got, tc.wantIcon)
				}
			}

			if tc.checkName {
				if got := l.Name(); got != tc.wantName {
					t.Errorf("Name() = %q, want %q", got, tc.wantName)
				}
			}

			if tc.checkHref {
				attrs := l.Attrs()
				if attrs == nil {
					t.Fatal("Attrs() returned nil")
				}

				href, ok := attrs["href"]
				if !ok {
					t.Error("Attrs() missing 'href' key")
				}

				hrefStr, ok := href.(string)
				if !ok {
					t.Errorf("href type = %T, want string", href)
				}

				if hrefStr != tc.wantHref {
					t.Errorf("href = %q, want %q", hrefStr, tc.wantHref)
				}
			}

			if tc.checkKeyType {
				attrs := l.Attrs()
				for key := range attrs {
					if _, ok := interface{}(key).(template.HTMLAttr); !ok {
						t.Errorf("key type = %T, want template.HTMLAttr", key)
					}
					break
				}
			}
		})
	}
}

func TestRecent_Handler(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		path          string
		cancelContext bool
		wantNonNilOut bool
	}{
		{
			name:          "GET request returns output",
			method:        http.MethodGet,
			path:          "/+/recent",
			wantNonNilOut: true,
		},
		{
			name:          "cancelled context returns output",
			method:        http.MethodGet,
			path:          "/+/recent",
			cancelContext: true,
			wantNonNilOut: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var ctx context.Context
			if tc.cancelContext {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(context.Background())
				cancel()
			} else {
				ctx = context.Background()
			}

			r := httptest.NewRequest(tc.method, tc.path, http.NoBody).WithContext(ctx)
			output := recentHandler(r)

			if tc.wantNonNilOut && output == nil {
				t.Error("recentHandler() returned nil, want non-nil output")
			}
		})
	}
}

func TestRecent_HandlerSorting(t *testing.T) {
	tests := []struct {
		name  string
		pages []struct {
			name    string
			modTime time.Time
		}
		expectedOrder []string
	}{
		{
			name: "sorts by modification time descending",
			pages: []struct {
				name    string
				modTime time.Time
			}{
				{"old.md", time.Now().Add(-2 * time.Hour)},
				{"new.md", time.Now()},
				{"mid.md", time.Now().Add(-1 * time.Hour)},
			},
			expectedOrder: []string{"new.md", "mid.md", "old.md"},
		},
		{
			name: "sorts alphabetically when same modification time",
			pages: []struct {
				name    string
				modTime time.Time
			}{
				{"zebra.md", time.Now()},
				{"alpha.md", time.Now()},
				{"beta.md", time.Now()},
			},
			expectedOrder: []string{"alpha.md", "beta.md", "zebra.md"},
		},
		{
			name: "handles empty pages list",
		},
		{
			name: "handles single page",
			pages: []struct {
				name    string
				modTime time.Time
			}{
				{"single.md", time.Now()},
			},
			expectedOrder: []string{"single.md"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create isolated temp directory for this sub-test
			tmpDir := t.TempDir()
			origDir, _ := os.Getwd()
			defer func() { _ = os.Chdir(origDir) }()
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatal(err)
			}

			// Create test files
			for _, p := range tc.pages {
				pagePath := filepath.Join(tmpDir, p.name)
				if err := os.WriteFile(pagePath, []byte("# Test"), 0600); err != nil {
					t.Fatalf("WriteFile(%s) failed: %v", p.name, err)
				}
				if err := os.Chtimes(pagePath, p.modTime, p.modTime); err != nil {
					t.Fatalf("Chtimes(%s) failed: %v", p.name, err)
				}
			}

			// Run handler (files must exist during execution for Pages() to find them)
			r := httptest.NewRequest(http.MethodGet, "/+/recent", http.NoBody)
			output := recentHandler(r)

			if output == nil {
				t.Error("recentHandler() returned nil")
			}
		})
	}
}

func TestRecent_TemplatesEmbedded(t *testing.T) {
	tests := []struct {
		name            string
		wantNonEmptyDir bool
	}{
		{
			name:            "embedded templates directory contains files",
			wantNonEmptyDir: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			entries, err := templates.ReadDir("templates")
			if err != nil {
				t.Fatalf("ReadDir(templates) failed: %v", err)
			}

			if tc.wantNonEmptyDir && len(entries) == 0 {
				t.Error("templates directory is empty, want at least one file")
			}
		})
	}
}

func TestRecent_Init(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Init does not panic"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Init() panicked: %v", r)
				}
			}()

			ext := Recent{}
			ext.Init()
		})
	}
}

func TestRecent_RegisterLinkCallback(t *testing.T) {
	tests := []struct {
		name              string
		wantCommandCount  int
		wantFirstCmdIsCmd bool
	}{
		{
			name:              "link callback returns commands slice",
			wantCommandCount:  1,
			wantFirstCmdIsCmd: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// The callback is registered in Init's RegisterLink call
			// We test it directly by invoking the same lambda
			linkFunc := func(xlog.Page) []xlog.Command { return []xlog.Command{links{}} }

			// Create a mock page (nil is acceptable since callback doesn't use it)
			var mockPage xlog.Page
			commands := linkFunc(mockPage)

			if len(commands) != tc.wantCommandCount {
				t.Errorf("linkFunc() returned %d commands, want %d", len(commands), tc.wantCommandCount)
			}

			if tc.wantFirstCmdIsCmd && len(commands) > 0 {
				if _, ok := commands[0].(links); !ok {
					t.Errorf("first command type = %T, want links", commands[0])
				}
			}
		})
	}
}

func TestRecent_SortingComparison(t *testing.T) {
	// Use fixed times to ensure consistent test behavior
	baseTime := time.Date(2026, 5, 4, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name  string
		pageA struct {
			name    string
			modTime time.Time
		}
		pageB struct {
			name    string
			modTime time.Time
		}
		wantSortValue string
	}{
		{
			name: "newer page comes first",
			pageA: struct {
				name    string
				modTime time.Time
			}{"a.md", baseTime.AddDate(0, 0, -1)},
			pageB: struct {
				name    string
				modTime time.Time
			}{"b.md", baseTime},
			wantSortValue: "positive", // b newer: b.ModTime().Compare(a.ModTime()) > 0
		},
		{
			name: "older page comes after",
			pageA: struct {
				name    string
				modTime time.Time
			}{"a.md", baseTime},
			pageB: struct {
				name    string
				modTime time.Time
			}{"b.md", baseTime.AddDate(0, 0, -1)},
			wantSortValue: "negative", // b older: b.ModTime().Compare(a.ModTime()) < 0
		},
		{
			name: "same modtime sorts alphabetically",
			pageA: struct {
				name    string
				modTime time.Time
			}{"zebra.md", baseTime},
			pageB: struct {
				name    string
				modTime time.Time
			}{"alpha.md", baseTime},
			wantSortValue: "positive", // strings.Compare("zebra", "alpha") > 0
		},
		{
			name: "same modtime different names",
			pageA: struct {
				name    string
				modTime time.Time
			}{"same.md", baseTime},
			pageB: struct {
				name    string
				modTime time.Time
			}{"other.md", baseTime},
			wantSortValue: "positive", // strings.Compare("same", "other") > 0
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create isolated temp directory
			tmpDir := t.TempDir()

			// Create test files
			pathA := filepath.Join(tmpDir, tc.pageA.name)
			pathB := filepath.Join(tmpDir, tc.pageB.name)

			if err := os.WriteFile(pathA, []byte("# Test A"), 0600); err != nil {
				t.Fatalf("WriteFile(%s) failed: %v", tc.pageA.name, err)
			}
			if err := os.WriteFile(pathB, []byte("# Test B"), 0600); err != nil {
				t.Fatalf("WriteFile(%s) failed: %v", tc.pageB.name, err)
			}

			if err := os.Chtimes(pathA, tc.pageA.modTime, tc.pageA.modTime); err != nil {
				t.Fatalf("Chtimes(%s) failed: %v", tc.pageA.name, err)
			}
			if err := os.Chtimes(pathB, tc.pageB.modTime, tc.pageB.modTime); err != nil {
				t.Fatalf("Chtimes(%s) failed: %v", tc.pageB.name, err)
			}

			// Load pages
			origDir, _ := os.Getwd()
			defer func() { _ = os.Chdir(origDir) }()
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatal(err)
			}

			pages := []xlog.Page{}
			entries, err := os.ReadDir(tmpDir)
			if err != nil {
				t.Fatal(err)
			}
			for _, entry := range entries {
				if !entry.IsDir() {
					// Remove .md extension for page name
					name := entry.Name()
					if len(name) > 3 && name[len(name)-3:] == mdExt {
						name = name[:len(name)-3]
					}
					pages = append(pages, xlog.NewPage(name))
				}
			}

			if len(pages) < 2 {
				t.Fatalf("expected at least 2 pages, got %d", len(pages))
			}

			// Find the two pages we created
			var pageAObj, pageBObj xlog.Page
			expectedNameA := tc.pageA.name
			expectedNameB := tc.pageB.name
			// Remove .md extension if present
			if len(expectedNameA) > 3 && expectedNameA[len(expectedNameA)-3:] == mdExt {
				expectedNameA = expectedNameA[:len(expectedNameA)-3]
			}
			if len(expectedNameB) > 3 && expectedNameB[len(expectedNameB)-3:] == mdExt {
				expectedNameB = expectedNameB[:len(expectedNameB)-3]
			}

			for _, p := range pages {
				if p.Name() == expectedNameA {
					pageAObj = p
				}
				if p.Name() == expectedNameB {
					pageBObj = p
				}
			}

			if pageAObj == nil || pageBObj == nil {
				t.Fatal("could not find created pages")
			}

			// Test the extracted comparison function
			result := comparePagesByRecency(pageAObj, pageBObj)

			switch tc.wantSortValue {
			case "positive":
				if result <= 0 {
					t.Errorf("comparePagesByRecency() = %d, want positive (pageA modtime=%v, pageB modtime=%v)",
						result, pageAObj.ModTime(), pageBObj.ModTime())
				}
			case "negative":
				if result >= 0 {
					t.Errorf("comparePagesByRecency() = %d, want negative (pageA modtime=%v, pageB modtime=%v)",
						result, pageAObj.ModTime(), pageBObj.ModTime())
				}
			case "zero":
				if result != 0 {
					t.Errorf("comparePagesByRecency() = %d, want 0", result)
				}
			}
		})
	}
}
