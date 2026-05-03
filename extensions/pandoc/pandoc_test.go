package pandoc

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
)

func TestPandocName(t *testing.T) {
	p := pandoc{}
	if got := p.Name(); got != "pandoc" {
		t.Errorf("Name() = %q, want %q", got, "pandoc")
	}
}

func TestPandocInit(t *testing.T) {
	tests := []struct {
		name           string
		pandocSupport  bool
		expectRegister bool
	}{
		{
			name:           "pandoc support enabled",
			pandocSupport:  true,
			expectRegister: true,
		},
		{
			name:           "pandoc support disabled",
			pandocSupport:  false,
			expectRegister: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			oldVal := pandoc_support
			defer func() { pandoc_support = oldVal }()

			pandoc_support = tc.pandocSupport
			p := pandoc{}
			p.Init()
		})
	}
}

func TestPandocPage(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd) //nolint:errcheck
	_ = os.Chdir(tmpDir)

	tests := []struct {
		name        string
		pageName    string
		createFiles map[string]string
		expectNil   bool
		expectedExt string
	}{
		{
			name:     "page with .org extension",
			pageName: "test",
			createFiles: map[string]string{
				"test.org": "* Heading\nContent",
			},
			expectNil:   false,
			expectedExt: ".org",
		},
		{
			name:     "page with .rst extension",
			pageName: "document",
			createFiles: map[string]string{
				"document.rst": "Title\n=====\n\nContent",
			},
			expectNil:   false,
			expectedExt: ".rst",
		},
		{
			name:     "page with .rtf extension",
			pageName: "notes",
			createFiles: map[string]string{
				"notes.rtf": "{\\rtf1 Content}",
			},
			expectNil:   false,
			expectedExt: ".rtf",
		},
		{
			name:     "page with .odt extension",
			pageName: "report",
			createFiles: map[string]string{
				"report.odt": "ODT content",
			},
			expectNil:   false,
			expectedExt: ".odt",
		},
		{
			name:        "page does not exist",
			pageName:    "nonexistent",
			createFiles: map[string]string{},
			expectNil:   true,
		},
		{
			name:     "page with unsupported extension",
			pageName: "test",
			createFiles: map[string]string{
				"test.md": "# Markdown",
			},
			expectNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testDir := t.TempDir()
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd) //nolint:errcheck
			_ = os.Chdir(testDir)

			for filename, content := range tc.createFiles {
				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file %s: %v", filename, err)
				}
			}

			p := &pandoc{}
			result := p.Page(tc.pageName)

			if tc.expectNil {
				if result != nil {
					t.Errorf("Expected nil page, got %v", result)
				}
			} else {
				if result == nil {
					t.Fatalf("Expected non-nil page")
				}
				if pg, ok := result.(*page); ok {
					if pg.ext != tc.expectedExt {
						t.Errorf("Expected extension %q, got %q", tc.expectedExt, pg.ext)
					}
					if pg.name != tc.pageName {
						t.Errorf("Expected name %q, got %q", tc.pageName, pg.name)
					}
				} else {
					t.Errorf("Expected *page type, got %T", result)
				}
			}
		})
	}
}

func TestPandocEach(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd) //nolint:errcheck
	_ = os.Chdir(tmpDir)

	testFiles := map[string]string{
		"page1.org":        "* Heading",
		"page2.rst":        "Title\n=====",
		"subdir/page3.rtf": "{\\rtf1}",
		"page4.odt":        "ODT",
		"ignored.md":       "# Markdown",
		"ignored.txt":      "Text",
	}

	for filename, content := range testFiles {
		dir := filepath.Dir(filename)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("Failed to create directory %s: %v", dir, err)
			}
		}
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	tests := []struct {
		name          string
		ctx           context.Context
		expectedPages []string
	}{
		{
			name: "iterate all supported pages",
			ctx:  context.Background(),
			expectedPages: []string{
				"page1",
				"page2",
				"page4",
				filepath.Join("subdir", "page3"),
			},
		},
		{
			name:          "context cancellation stops iteration",
			ctx:           func() context.Context { ctx, cancel := context.WithCancel(context.Background()); cancel(); return ctx }(),
			expectedPages: []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := &pandoc{}
			var foundPages []string

			p.Each(tc.ctx, func(page xlog.Page) {
				foundPages = append(foundPages, page.Name())
			})

			if tc.ctx.Err() != nil {
				return
			}

			if len(foundPages) != len(tc.expectedPages) {
				t.Errorf("Expected %d pages, got %d: %v", len(tc.expectedPages), len(foundPages), foundPages)
			}

			for _, expected := range tc.expectedPages {
				found := false
				for _, actual := range foundPages {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected to find page %q, but it was not in %v", expected, foundPages)
				}
			}
		})
	}
}

func TestPageName(t *testing.T) {
	tests := []struct {
		name     string
		page     page
		expected string
	}{
		{
			name:     "simple page name",
			page:     page{name: "test", ext: ".org"},
			expected: "test",
		},
		{
			name:     "page with path",
			page:     page{name: "path/to/page", ext: ".rst"},
			expected: "path/to/page",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.page.Name(); got != tc.expected {
				t.Errorf("Name() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestPageFileName(t *testing.T) {
	tests := []struct {
		name     string
		page     page
		expected string
	}{
		{
			name:     "simple filename",
			page:     page{name: "test", ext: ".org"},
			expected: "test.org",
		},
		{
			name:     "filename with path",
			page:     page{name: "path/to/page", ext: ".rst"},
			expected: filepath.FromSlash("path/to/page") + ".rst",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.page.FileName(); got != tc.expected {
				t.Errorf("FileName() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestPageExists(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd) //nolint:errcheck
	_ = os.Chdir(tmpDir)

	if err := os.WriteFile("existing.org", []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		page     page
		expected bool
	}{
		{
			name:     "file exists",
			page:     page{name: "existing", ext: ".org"},
			expected: true,
		},
		{
			name:     "file does not exist",
			page:     page{name: "nonexistent", ext: ".org"},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.page.Exists(); got != tc.expected {
				t.Errorf("Exists() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestPageContent(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd) //nolint:errcheck
	_ = os.Chdir(tmpDir)

	testContent := "* Test Heading\nTest content"
	if err := os.WriteFile("test.org", []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		page     page
		expected xlog.Markdown
	}{
		{
			name:     "read existing file",
			page:     page{name: "test", ext: ".org"},
			expected: xlog.Markdown(testContent),
		},
		{
			name:     "read nonexistent file",
			page:     page{name: "nonexistent", ext: ".org"},
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.page.Content(); got != tc.expected {
				t.Errorf("Content() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestPageModTime(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd) //nolint:errcheck
	_ = os.Chdir(tmpDir)

	now := time.Now()
	if err := os.WriteFile("test.org", []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name         string
		page         page
		expectZero   bool
		expectRecent bool
	}{
		{
			name:         "existing file has recent modtime",
			page:         page{name: "test", ext: ".org"},
			expectZero:   false,
			expectRecent: true,
		},
		{
			name:         "nonexistent file has zero modtime",
			page:         page{name: "nonexistent", ext: ".org"},
			expectZero:   true,
			expectRecent: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.page.ModTime()

			if tc.expectZero {
				if !got.IsZero() {
					t.Errorf("Expected zero time, got %v", got)
				}
			} else {
				if got.IsZero() {
					t.Errorf("Expected non-zero time, got zero")
				}
			}

			if tc.expectRecent {
				if got.Before(now.Add(-1 * time.Minute)) {
					t.Errorf("Expected recent modtime (within 1 minute of %v), got %v", now, got)
				}
			}
		})
	}
}

func TestPageDelete(t *testing.T) {
	tests := []struct {
		name           string
		createFile     bool
		expectedResult bool
		expectFileGone bool
	}{
		{
			name:           "delete existing file",
			createFile:     true,
			expectedResult: true,
			expectFileGone: true,
		},
		{
			name:           "delete nonexistent file",
			createFile:     false,
			expectedResult: true,
			expectFileGone: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd) //nolint:errcheck
			_ = os.Chdir(tmpDir)

			p := &page{name: "test", ext: ".org"}

			if tc.createFile {
				if err := os.WriteFile(p.FileName(), []byte("content"), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			result := p.Delete()

			if result != tc.expectedResult {
				t.Errorf("Delete() = %v, want %v", result, tc.expectedResult)
			}

			if tc.expectFileGone {
				if _, err := os.Stat(p.FileName()); !os.IsNotExist(err) {
					t.Errorf("Expected file to be deleted, but it still exists")
				}
			}
		})
	}
}

func TestPageWrite(t *testing.T) {
	tests := []struct {
		name             string
		content          xlog.Markdown
		readonlyMode     bool
		expectedResult   bool
		expectFileExists bool
		createNestedDir  bool
	}{
		{
			name:             "write new file",
			content:          "* New Content",
			readonlyMode:     false,
			expectedResult:   true,
			expectFileExists: true,
		},
		{
			name:             "write with CRLF normalization",
			content:          "Line1\r\nLine2\r\nLine3",
			readonlyMode:     false,
			expectedResult:   true,
			expectFileExists: true,
		},
		{
			name:             "write in readonly mode",
			content:          "Content",
			readonlyMode:     true,
			expectedResult:   false,
			expectFileExists: false,
		},
		{
			name:             "write to nested directory",
			content:          "Nested content",
			readonlyMode:     false,
			expectedResult:   true,
			expectFileExists: true,
			createNestedDir:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd) //nolint:errcheck
			_ = os.Chdir(tmpDir)

			oldReadonly := xlog.Config.Readonly
			defer func() { xlog.Config.Readonly = oldReadonly }()
			xlog.Config.Readonly = tc.readonlyMode

			pageName := "test"
			if tc.createNestedDir {
				pageName = "nested/dir/test"
			}

			p := &page{name: pageName, ext: ".org"}
			result := p.Write(tc.content)

			if result != tc.expectedResult {
				t.Errorf("Write() = %v, want %v", result, tc.expectedResult)
			}

			fileExists := false
			if data, err := os.ReadFile(p.FileName()); err == nil {
				fileExists = true

				if tc.expectedResult {
					expectedContent := strings.ReplaceAll(string(tc.content), "\r\n", "\n")
					if string(data) != expectedContent {
						t.Errorf("File content = %q, want %q", string(data), expectedContent)
					}
				}
			}

			if fileExists != tc.expectFileExists {
				t.Errorf("File exists = %v, want %v", fileExists, tc.expectFileExists)
			}
		})
	}
}

func TestPageAST(t *testing.T) {
	p := &page{name: "test", ext: ".org"}
	data, node := p.AST()

	if len(data) != 0 {
		t.Errorf("Expected empty byte slice, got %v", data)
	}

	if node == nil {
		t.Errorf("Expected non-nil AST node, got nil")
	}
}

func TestSupportedExtensions(t *testing.T) {
	expected := []string{".org", ".rst", ".rtf", ".odt"}

	if len(SUPPORTED_EXT) != len(expected) {
		t.Errorf("Expected %d supported extensions, got %d", len(expected), len(SUPPORTED_EXT))
	}

	for i, ext := range expected {
		if SUPPORTED_EXT[i] != ext {
			t.Errorf("Expected extension %q at index %d, got %q", ext, i, SUPPORTED_EXT[i])
		}
	}
}
