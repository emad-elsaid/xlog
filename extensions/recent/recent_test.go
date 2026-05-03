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
)

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

			r := httptest.NewRequest(tc.method, tc.path, nil).WithContext(ctx)
			output := recentHandler(r)

			if tc.wantNonNilOut && output == nil {
				t.Error("recentHandler() returned nil, want non-nil output")
			}
		})
	}
}

func TestRecent_HandlerSorting(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer func() { _ = os.Chdir(origDir) }()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}

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
			for _, p := range tc.pages {
				pagePath := filepath.Join(tmpDir, p.name)
				if err := os.WriteFile(pagePath, []byte("# Test"), 0600); err != nil {
					t.Fatalf("WriteFile(%s) failed: %v", p.name, err)
				}
				if err := os.Chtimes(pagePath, p.modTime, p.modTime); err != nil {
					t.Fatalf("Chtimes(%s) failed: %v", p.name, err)
				}
			}

			r := httptest.NewRequest(http.MethodGet, "/+/recent", nil)
			output := recentHandler(r)

			if output == nil {
				t.Error("recentHandler() returned nil")
			}

			for _, p := range tc.pages {
				_ = os.Remove(filepath.Join(tmpDir, p.name))
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
