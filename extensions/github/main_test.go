package github

import (
	"flag"
	"html/template"
	"strings"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
)

const testEditURL = "https://github.com/emad-elsaid/xlog/edit/master/docs"

// mockPage implements xlog.Page interface for testing.
type mockPage struct {
	name     string
	fileName string
	exists   bool
}

func (m mockPage) Name() string             { return m.name }
func (m mockPage) FileName() string         { return m.fileName }
func (m mockPage) Exists() bool             { return m.exists }
func (m mockPage) Render() template.HTML    { return "" }
func (m mockPage) Content() xlog.Markdown   { return xlog.Markdown("") }
func (m mockPage) Delete() bool             { return false }
func (m mockPage) Write(xlog.Markdown) bool { return false }
func (m mockPage) ModTime() time.Time       { return time.Now() }
func (m mockPage) AST() ([]byte, ast.Node)  { return []byte{}, nil }

func TestGithubExtensionName(t *testing.T) {
	ext := Github{}
	expected := "github"
	if got := ext.Name(); got != expected {
		t.Errorf("Name() = %q, want %q", got, expected)
	}
}

func TestGithubInit_EmptyEditUrl(t *testing.T) {
	// Save original editUrl
	originalEditUrl := editUrl
	defer func() { editUrl = originalEditUrl }()

	editUrl = ""

	// Init should not panic with empty URL
	ext := Github{}
	ext.Init()
	// No assertions needed - just verify it doesn't panic
}

func TestGithubInit_WithEditUrl(t *testing.T) {
	// Save original editUrl
	originalEditUrl := editUrl
	defer func() { editUrl = originalEditUrl }()

	editUrl = testEditURL

	ext := Github{}
	ext.Init()
	// Init should register quick command - no panic expected
}

func TestQuickCommands_EmptyFileName(t *testing.T) {
	// Save original editUrl
	originalEditUrl := editUrl
	defer func() { editUrl = originalEditUrl }()

	editUrl = testEditURL

	page := mockPage{
		name:     "test",
		fileName: "",
		exists:   true,
	}

	commands := quickCommands(page)
	if commands != nil {
		t.Errorf("quickCommands() with empty filename should return nil, got %d commands", len(commands))
	}
}

func TestQuickCommands_WithFileName(t *testing.T) {
	// Save original editUrl
	originalEditUrl := editUrl
	defer func() { editUrl = originalEditUrl }()

	editUrl = testEditURL

	page := mockPage{
		name:     "test-page",
		fileName: "test-page.md",
		exists:   true,
	}

	commands := quickCommands(page)
	if commands == nil {
		t.Fatal("quickCommands() should return commands for page with filename")
	}

	if len(commands) != 1 {
		t.Errorf("quickCommands() should return 1 command, got %d", len(commands))
	}
}

func TestEditOnGithub_Icon(t *testing.T) {
	page := mockPage{
		name:     "test",
		fileName: "test.md",
	}

	cmd := editOnGithub{page: page}
	expected := "fa-brands fa-github"
	if got := cmd.Icon(); got != expected {
		t.Errorf("Icon() = %q, want %q", got, expected)
	}
}

func TestEditOnGithub_Name(t *testing.T) {
	page := mockPage{
		name:     "test",
		fileName: "test.md",
	}

	cmd := editOnGithub{page: page}
	expected := "Edit on Github"
	if got := cmd.Name(); got != expected {
		t.Errorf("Name() = %q, want %q", got, expected)
	}
}

func TestEditOnGithub_Attrs(t *testing.T) {
	// Save original editUrl
	originalEditUrl := editUrl
	defer func() { editUrl = originalEditUrl }()

	tests := []struct {
		name     string
		editUrl  string
		fileName string
		wantHref string
	}{
		{
			name:     "basic path",
			editUrl:  "https://github.com/user/repo/edit/master",
			fileName: "docs/test.md",
			wantHref: "https://github.com/user/repo/edit/master/docs/test.md",
		},
		{
			name:     "path with special characters",
			fileName: "docs/my file.md",
			editUrl:  "https://github.com/user/repo/edit/main",
			wantHref: "https://github.com/user/repo/edit/main/docs/my file.md",
		},
		{
			name:     "root file",
			fileName: "README.md",
			editUrl:  "https://github.com/user/repo/edit/develop",
			wantHref: "https://github.com/user/repo/edit/develop/README.md",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			editUrl = tc.editUrl
			page := mockPage{fileName: tc.fileName}
			cmd := editOnGithub{page: page}

			attrs := cmd.Attrs()
			if attrs == nil {
				t.Fatal("Attrs() returned nil")
			}

			hrefAttr := template.HTMLAttr("href")
			href, ok := attrs[hrefAttr]
			if !ok {
				t.Fatal("Attrs() missing 'href' attribute")
			}

			hrefStr, ok := href.(string)
			if !ok {
				t.Fatalf("href attribute is not a string: %T", href)
			}

			if hrefStr != tc.wantHref {
				t.Errorf("href = %q, want %q", hrefStr, tc.wantHref)
			}
		})
	}
}

func TestGithubFlagRegistration(t *testing.T) {
	// Verify the flag was registered
	f := flag.Lookup("github.url")
	if f == nil {
		t.Fatal("github.url flag should be registered")
		return
	}

	expectedUsage := "Repository url for 'edit on Github' quick action e.g https://github.com/emad-elsaid/xlog/edit/master/docs"
	if f.Usage != expectedUsage {
		t.Errorf("flag usage = %q, want %q", f.Usage, expectedUsage)
	}
}

func TestEditOnGithub_AttrsFormat(t *testing.T) {
	// Save original editUrl
	originalEditUrl := editUrl
	defer func() { editUrl = originalEditUrl }()

	editUrl = "https://github.com/test/repo/edit/master"
	page := mockPage{fileName: "test.md"}
	cmd := editOnGithub{page: page}

	attrs := cmd.Attrs()

	// Verify it returns a non-empty map
	if len(attrs) == 0 {
		t.Error("Attrs() should return non-empty map")
	}

	// Verify href key exists
	hrefKey := template.HTMLAttr("href")
	if _, ok := attrs[hrefKey]; !ok {
		t.Error("Attrs() should contain 'href' key")
	}
}

func TestQuickCommands_Integration(t *testing.T) {
	// Save original editUrl
	originalEditUrl := editUrl
	defer func() { editUrl = originalEditUrl }()

	editUrl = "https://github.com/owner/repo/edit/main/content"

	page := mockPage{
		name:     "my-page",
		fileName: "my-page.md",
		exists:   true,
	}

	commands := quickCommands(page)
	if len(commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(commands))
	}

	cmd := commands[0]

	// Verify command properties
	if !strings.Contains(cmd.Icon(), "github") {
		t.Errorf("expected github icon, got %q", cmd.Icon())
	}

	if !strings.Contains(cmd.Name(), "Github") {
		t.Errorf("expected Github in name, got %q", cmd.Name())
	}

	attrs := cmd.Attrs()
	href, ok := attrs[template.HTMLAttr("href")]
	if !ok {
		t.Fatal("command should have href attribute")
	}

	hrefStr := href.(string)
	expectedHref := "https://github.com/owner/repo/edit/main/content/my-page.md"
	if hrefStr != expectedHref {
		t.Errorf("href = %q, want %q", hrefStr, expectedHref)
	}
}
