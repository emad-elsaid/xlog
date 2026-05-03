package gpg

import (
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
)

const testGPGKeyID = "test-key-id"

// mockPage implements xlog.Page interface for testing.
type mockPage struct {
	name     string
	fileName string
	exists   bool
}

func (m *mockPage) Name() string                     { return m.name }
func (m *mockPage) FileName() string                 { return m.fileName }
func (m *mockPage) Exists() bool                     { return m.exists }
func (m *mockPage) Render() template.HTML            { return "" }
func (m *mockPage) Content() xlog.Markdown           { return "" }
func (m *mockPage) ModTime() time.Time               { return time.Time{} }
func (m *mockPage) Delete() bool                     { return true }
func (m *mockPage) Write(content xlog.Markdown) bool { return true }
func (m *mockPage) AST() ([]byte, ast.Node)          { return nil, nil }

// TestCommands_NonExistentPage verifies commands returns nil for non-existent pages.
func TestCommands_NonExistentPage(t *testing.T) {
	p := &mockPage{
		name:   "test",
		exists: false,
	}

	cmds := commands(p)
	if cmds != nil {
		t.Errorf("commands() = %v, want nil for non-existent page", cmds)
	}
}

// TestCommands_NoGPGId verifies commands returns nil when gpgId is empty.
func TestCommands_NoGPGId(t *testing.T) {
	// Save and restore original gpgId
	origGpgId := gpgId
	defer func() { gpgId = origGpgId }()

	gpgId = ""

	p := &mockPage{
		name:   "test",
		exists: true,
	}

	cmds := commands(p)
	if cmds != nil {
		t.Errorf("commands() = %v, want nil when gpgId is empty", cmds)
	}
}

// TestCommands_EncryptablePages verifies encrypt command for non-encrypted pages.
func TestCommands_EncryptablePages(t *testing.T) {
	// Save and restore original gpgId
	origGpgId := gpgId
	defer func() { gpgId = origGpgId }()

	gpgId = testGPGKeyID

	tests := []struct {
		name     string
		fileName string
	}{
		{"markdown page", "test.md"},
		{"text page", "test.txt"},
		{"no extension", "test"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := &mockPage{
				name:     "test",
				fileName: tc.fileName,
				exists:   true,
			}

			cmds := commands(p)

			if len(cmds) != 1 {
				t.Fatalf("commands() returned %d commands, want 1", len(cmds))
			}

			encCmd, ok := cmds[0].(*encryptCommand)
			if !ok {
				t.Fatalf("commands()[0] type = %T, want *encryptCommand", cmds[0])
			}

			if encCmd.Icon() != "fa-solid fa-lock" {
				t.Errorf("Icon() = %q, want %q", encCmd.Icon(), "fa-solid fa-lock")
			}

			if encCmd.Name() != "Make private" {
				t.Errorf("Name() = %q, want %q", encCmd.Name(), "Make private")
			}

			attrs := encCmd.Attrs()
			expectedURL := "/+/gpg/encrypt/" + url.PathEscape(p.Name())
			if attrs["hx-post"] != expectedURL {
				t.Errorf("Attrs[hx-post] = %q, want %q", attrs["hx-post"], expectedURL)
			}
		})
	}
}

// TestCommands_DecryptablePages verifies decrypt command for encrypted pages.
func TestCommands_DecryptablePages(t *testing.T) {
	// Save and restore original gpgId
	origGpgId := gpgId
	defer func() { gpgId = origGpgId }()

	gpgId = testGPGKeyID

	p := &mockPage{
		name:     "test",
		fileName: "test.pgp",
		exists:   true,
	}

	cmds := commands(p)

	if len(cmds) != 1 {
		t.Fatalf("commands() returned %d commands, want 1", len(cmds))
	}

	decCmd, ok := cmds[0].(*decryptCommand)
	if !ok {
		t.Fatalf("commands()[0] type = %T, want *decryptCommand", cmds[0])
	}

	if decCmd.Icon() != "fa-solid fa-lock-open has-text-danger" {
		t.Errorf("Icon() = %q, want %q", decCmd.Icon(), "fa-solid fa-lock-open has-text-danger")
	}

	if decCmd.Name() != "Make public" {
		t.Errorf("Name() = %q, want %q", decCmd.Name(), "Make public")
	}

	attrs := decCmd.Attrs()
	expectedURL := "/+/gpg/decrypt/" + url.PathEscape(p.Name())
	if attrs["hx-post"] != expectedURL {
		t.Errorf("Attrs[hx-post] = %q, want %q", attrs["hx-post"], expectedURL)
	}
}

// TestEncryptCommand_AttrsWithSpecialCharacters verifies URL encoding in command attributes.
func TestEncryptCommand_AttrsWithSpecialCharacters(t *testing.T) {
	tests := []struct {
		name         string
		pageName     string
		expectedPath string
	}{
		{
			name:         "page with spaces",
			pageName:     "my test page",
			expectedPath: "/+/gpg/encrypt/my%20test%20page",
		},
		{
			name:         "page with special chars",
			pageName:     "test/page?name=value",
			expectedPath: "/+/gpg/encrypt/test%2Fpage%3Fname=value",
		},
		{
			name:         "page with unicode",
			pageName:     "测试页面",
			expectedPath: "/+/gpg/encrypt/%E6%B5%8B%E8%AF%95%E9%A1%B5%E9%9D%A2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := &mockPage{name: tc.pageName}
			cmd := &encryptCommand{page: p}

			attrs := cmd.Attrs()
			if attrs["hx-post"] != tc.expectedPath {
				t.Errorf("Attrs[hx-post] = %q, want %q", attrs["hx-post"], tc.expectedPath)
			}
		})
	}
}

// TestPageDelete verifies the Delete method.
func TestPageDelete(t *testing.T) {
	tests := []struct {
		name       string
		setupFile  bool
		expectTrue bool
	}{
		{
			name:       "non-existent file returns true",
			setupFile:  false,
			expectTrue: true,
		},
		{
			name:       "existing file deleted successfully",
			setupFile:  true,
			expectTrue: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			origWd, _ := os.Getwd()
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("failed to change directory: %v", err)
			}
			defer func() {
				if err := os.Chdir(origWd); err != nil {
					t.Errorf("failed to restore directory: %v", err)
				}
			}()

			p := &page{name: "test/delete"}

			if tc.setupFile {
				if err := os.MkdirAll(filepath.Dir(p.FileName()), 0700); err != nil {
					t.Fatalf("failed to create directory: %v", err)
				}
				if err := os.WriteFile(p.FileName(), []byte("test"), 0600); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			}

			result := p.Delete()

			if result != tc.expectTrue {
				t.Errorf("Delete() = %v, want %v", result, tc.expectTrue)
			}

			if tc.setupFile && p.Exists() {
				t.Errorf("File still exists after Delete()")
			}
		})
	}
}
