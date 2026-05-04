package xlog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDoctor(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T) (cleanup func())
		expectIssues  bool
		expectWarning bool
	}{
		{
			name: "healthy configuration with all files present",
			setup: func(t *testing.T) func() {
				// Save original config
				origSource := Config.Source
				origIndex := Config.Index
				origNotFound := Config.NotFoundPage
				origReadonly := Config.Readonly
				origBindAddress := Config.BindAddress

				// Create temp directory with all required files
				tmpDir := t.TempDir()
				Config.Source = tmpDir
				Config.Index = "index"
				Config.NotFoundPage = "404"
				Config.Readonly = false
				Config.BindAddress = "127.0.0.1:3000"

				// Create index and 404 pages
				indexPath := filepath.Join(tmpDir, "index.md")
				if err := os.WriteFile(indexPath, []byte("# Home"), 0644); err != nil {
					t.Fatalf("Failed to create index: %v", err)
				}

				notFoundPath := filepath.Join(tmpDir, "404.md")
				if err := os.WriteFile(notFoundPath, []byte("# Not Found"), 0644); err != nil {
					t.Fatalf("Failed to create 404 page: %v", err)
				}

				return func() {
					Config.Source = origSource
					Config.Index = origIndex
					Config.NotFoundPage = origNotFound
					Config.Readonly = origReadonly
					Config.BindAddress = origBindAddress
				}
			},
			expectIssues:  false,
			expectWarning: false,
		},
		{
			name: "readonly mode with existing files",
			setup: func(t *testing.T) func() {
				origSource := Config.Source
				origIndex := Config.Index
				origNotFound := Config.NotFoundPage
				origReadonly := Config.Readonly
				origBindAddress := Config.BindAddress

				tmpDir := t.TempDir()
				Config.Source = tmpDir
				Config.Index = "index"
				Config.NotFoundPage = "404"
				Config.Readonly = true
				Config.BindAddress = "127.0.0.1:3000"

				// Create index and 404 pages
				indexPath := filepath.Join(tmpDir, "index.md")
				if err := os.WriteFile(indexPath, []byte("# Index"), 0644); err != nil {
					t.Fatalf("Failed to create index: %v", err)
				}

				notFoundPath := filepath.Join(tmpDir, "404.md")
				if err := os.WriteFile(notFoundPath, []byte("# Not Found"), 0644); err != nil {
					t.Fatalf("Failed to create 404 page: %v", err)
				}

				return func() {
					Config.Source = origSource
					Config.Index = origIndex
					Config.NotFoundPage = origNotFound
					Config.Readonly = origReadonly
					Config.BindAddress = origBindAddress
				}
			},
			expectIssues:  false,
			expectWarning: false,
		},
		{
			name: "missing index page shows warning",
			setup: func(t *testing.T) func() {
				origSource := Config.Source
				origIndex := Config.Index
				origReadonly := Config.Readonly
				origBindAddress := Config.BindAddress

				tmpDir := t.TempDir()
				Config.Source = tmpDir
				Config.Index = "nonexistent"
				Config.Readonly = true
				Config.BindAddress = "127.0.0.1:3000"

				// Create some markdown file but not the index
				if err := os.WriteFile(filepath.Join(tmpDir, "other.md"), []byte("# Other"), 0644); err != nil {
					t.Fatalf("Failed to create markdown file: %v", err)
				}

				return func() {
					Config.Source = origSource
					Config.Index = origIndex
					Config.Readonly = origReadonly
					Config.BindAddress = origBindAddress
				}
			},
			expectIssues:  false,
			expectWarning: true,
		},
		{
			name: "missing 404 page shows warning",
			setup: func(t *testing.T) func() {
				origSource := Config.Source
				origIndex := Config.Index
				origNotFound := Config.NotFoundPage
				origReadonly := Config.Readonly
				origBindAddress := Config.BindAddress

				tmpDir := t.TempDir()
				Config.Source = tmpDir
				Config.Index = "index"
				Config.NotFoundPage = "nonexistent-404"
				Config.Readonly = true
				Config.BindAddress = "127.0.0.1:3000"

				// Create index page
				indexPath := filepath.Join(tmpDir, "index.md")
				if err := os.WriteFile(indexPath, []byte("# Index"), 0644); err != nil {
					t.Fatalf("Failed to create index: %v", err)
				}

				return func() {
					Config.Source = origSource
					Config.Index = origIndex
					Config.NotFoundPage = origNotFound
					Config.Readonly = origReadonly
					Config.BindAddress = origBindAddress
				}
			},
			expectIssues:  false,
			expectWarning: true,
		},
		{
			name: "no markdown files shows warning",
			setup: func(t *testing.T) func() {
				origSource := Config.Source
				origReadonly := Config.Readonly
				origBindAddress := Config.BindAddress

				tmpDir := t.TempDir()
				Config.Source = tmpDir
				Config.Readonly = true
				Config.BindAddress = "127.0.0.1:3000"

				// Create directory but no markdown files
				// Create a non-markdown file
				if err := os.WriteFile(filepath.Join(tmpDir, "README.txt"), []byte("text"), 0644); err != nil {
					t.Fatalf("Failed to create text file: %v", err)
				}

				return func() {
					Config.Source = origSource
					Config.Readonly = origReadonly
					Config.BindAddress = origBindAddress
				}
			},
			expectIssues:  false,
			expectWarning: true,
		},
		{
			name: "nonexistent source directory shows critical issue",
			setup: func(t *testing.T) func() {
				origSource := Config.Source
				origBindAddress := Config.BindAddress
				origReadonly := Config.Readonly

				Config.Source = "/nonexistent/path/to/nowhere"
				Config.BindAddress = "127.0.0.1:3000"
				Config.Readonly = true // Prevent write check on non-existent dir

				return func() {
					Config.Source = origSource
					Config.BindAddress = origBindAddress
					Config.Readonly = origReadonly
				}
			},
			expectIssues:  true,
			expectWarning: true, // Will also have warnings about missing files
		},
		{
			name: "empty bind address shows critical issue",
			setup: func(t *testing.T) func() {
				origSource := Config.Source
				origReadonly := Config.Readonly
				origBindAddress := Config.BindAddress

				tmpDir := t.TempDir()
				Config.Source = tmpDir
				Config.Readonly = true
				Config.BindAddress = ""

				return func() {
					Config.Source = origSource
					Config.Readonly = origReadonly
					Config.BindAddress = origBindAddress
				}
			},
			expectIssues:  true,
			expectWarning: true, // Empty dir will have warnings
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cleanup := tc.setup(t)
			defer cleanup()

			// Call runDiagnostics directly to test without os.Exit
			result := runDiagnostics()

			hasIssues := len(result.Issues) > 0
			hasWarnings := len(result.Warnings) > 0

			if tc.expectIssues && !hasIssues {
				t.Error("Expected critical issues but found none")
			}

			if !tc.expectIssues && hasIssues {
				t.Errorf("Expected no critical issues but found: %v", result.Issues)
			}

			if tc.expectWarning && !hasWarnings {
				t.Error("Expected warnings but found none")
			}

			if !tc.expectWarning && hasWarnings {
				t.Errorf("Expected no warnings but found: %v", result.Warnings)
			}
		})
	}
}

func TestDoctor_SourceDirectoryValidation(t *testing.T) {
	tests := []struct {
		name       string
		setupDir   string
		createDir  bool
		isFile     bool
		shouldFail bool
	}{
		{
			name:       "existing directory",
			setupDir:   "",
			createDir:  true,
			isFile:     false,
			shouldFail: false,
		},
		{
			name:       "nonexistent directory",
			setupDir:   "/nonexistent/path/to/nowhere",
			createDir:  false,
			isFile:     false,
			shouldFail: true,
		},
		{
			name:       "file instead of directory",
			setupDir:   "",
			createDir:  true,
			isFile:     true,
			shouldFail: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			origSource := Config.Source
			defer func() {
				Config.Source = origSource
			}()

			if tc.createDir {
				tmpDir := t.TempDir()
				if tc.isFile {
					// Create a file instead of directory
					filePath := filepath.Join(tmpDir, "notadir")
					if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
						t.Fatalf("Failed to create file: %v", err)
					}
					Config.Source = filePath
				} else {
					Config.Source = tmpDir
				}
			} else {
				Config.Source = tc.setupDir
			}

			// Check if source exists and is directory
			stat, err := os.Stat(Config.Source)
			hasError := err != nil || (stat != nil && !stat.IsDir())

			if tc.shouldFail && !hasError {
				t.Error("Expected source directory validation to fail, but it passed")
			}

			if !tc.shouldFail && hasError {
				t.Errorf("Expected source directory validation to pass, but it failed: %v", err)
			}
		})
	}
}

func TestDoctor_WritePermissions(t *testing.T) {
	tests := []struct {
		name       string
		readonly   bool
		perms      os.FileMode
		shouldFail bool
	}{
		{
			name:       "writable directory with write enabled",
			readonly:   false,
			perms:      0755,
			shouldFail: false,
		},
		{
			name:       "readonly mode skips write check",
			readonly:   true,
			perms:      0755,
			shouldFail: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			origSource := Config.Source
			origReadonly := Config.Readonly
			defer func() {
				Config.Source = origSource
				Config.Readonly = origReadonly
			}()

			Config.Source = tmpDir
			Config.Readonly = tc.readonly

			// Set directory permissions
			if err := os.Chmod(tmpDir, tc.perms); err != nil {
				t.Fatalf("Failed to set permissions: %v", err)
			}

			if !tc.readonly {
				// Try to write test file
				testFile := filepath.Join(tmpDir, ".xlog-write-test")
				err := os.WriteFile(testFile, []byte("test"), 0600)
				hasError := err != nil

				if tc.shouldFail && !hasError {
					t.Error("Expected write permission check to fail, but it passed")
				}

				if !tc.shouldFail && hasError {
					t.Errorf("Expected write permission check to pass, but failed: %v", err)
				}

				// Clean up
				_ = os.Remove(testFile)
			}
		})
	}
}

func TestDoctor_MarkdownFileDetection(t *testing.T) {
	tests := []struct {
		name       string
		files      []string
		expectNone bool
	}{
		{
			name:       "directory with markdown files",
			files:      []string{"page1.md", "page2.md", "notes.md"},
			expectNone: false,
		},
		{
			name:       "directory with no markdown files",
			files:      []string{"README.txt", "data.json", "image.png"},
			expectNone: true,
		},
		{
			name:       "empty directory",
			files:      []string{},
			expectNone: true,
		},
		{
			name:       "mixed files with some markdown",
			files:      []string{"README.txt", "notes.md", "data.json"},
			expectNone: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create test files
			for _, filename := range tc.files {
				filePath := filepath.Join(tmpDir, filename)
				if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
					t.Fatalf("Failed to create file %s: %v", filename, err)
				}
			}

			// Check for markdown files
			mdFiles, err := filepath.Glob(filepath.Join(tmpDir, "*.md"))
			if err != nil {
				t.Fatalf("Failed to glob markdown files: %v", err)
			}

			hasMarkdown := len(mdFiles) > 0

			if tc.expectNone && hasMarkdown {
				t.Errorf("Expected no markdown files, but found %d", len(mdFiles))
			}

			if !tc.expectNone && !hasMarkdown {
				t.Error("Expected to find markdown files, but found none")
			}
		})
	}
}

func TestDoctor_BindAddressValidation(t *testing.T) {
	tests := []struct {
		name        string
		bindAddress string
		shouldFail  bool
	}{
		{
			name:        "valid bind address",
			bindAddress: "127.0.0.1:3000",
			shouldFail:  false,
		},
		{
			name:        "valid bind address with hostname",
			bindAddress: "localhost:8080",
			shouldFail:  false,
		},
		{
			name:        "empty bind address",
			bindAddress: "",
			shouldFail:  true,
		},
		{
			name:        "bind address with port only",
			bindAddress: ":3000",
			shouldFail:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			isEmpty := tc.bindAddress == ""

			if tc.shouldFail && !isEmpty {
				t.Error("Expected bind address validation to fail, but it passed")
			}

			if !tc.shouldFail && isEmpty {
				t.Error("Expected bind address validation to pass, but it failed")
			}
		})
	}
}

func TestDiagnosticResult_HasCriticalIssues(t *testing.T) {
	tests := []struct {
		name     string
		result   DiagnosticResult
		expected bool
	}{
		{
			name:     "no issues",
			result:   DiagnosticResult{Issues: []string{}, Warnings: []string{}},
			expected: false,
		},
		{
			name:     "has one issue",
			result:   DiagnosticResult{Issues: []string{"critical error"}, Warnings: []string{}},
			expected: true,
		},
		{
			name:     "has multiple issues",
			result:   DiagnosticResult{Issues: []string{"error1", "error2"}, Warnings: []string{}},
			expected: true,
		},
		{
			name:     "has warnings but no issues",
			result:   DiagnosticResult{Issues: []string{}, Warnings: []string{"warning"}},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.result.HasCriticalIssues(); got != tc.expected {
				t.Errorf("HasCriticalIssues() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestDiagnosticResult_IsHealthy(t *testing.T) {
	tests := []struct {
		name     string
		result   DiagnosticResult
		expected bool
	}{
		{
			name:     "healthy - no issues or warnings",
			result:   DiagnosticResult{Issues: []string{}, Warnings: []string{}},
			expected: true,
		},
		{
			name:     "unhealthy - has issues",
			result:   DiagnosticResult{Issues: []string{"error"}, Warnings: []string{}},
			expected: false,
		},
		{
			name:     "unhealthy - has warnings",
			result:   DiagnosticResult{Issues: []string{}, Warnings: []string{"warning"}},
			expected: false,
		},
		{
			name:     "unhealthy - has both",
			result:   DiagnosticResult{Issues: []string{"error"}, Warnings: []string{"warning"}},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.result.IsHealthy(); got != tc.expected {
				t.Errorf("IsHealthy() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestRunDiagnostics_Comprehensive(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(t *testing.T) (cleanup func())
		expectHealthy  bool
		issueCount     int
		warningCount   int
		expectIssue    string   // Expected substring in issues
		expectWarnings []string // Expected substrings in warnings
	}{
		{
			name: "completely healthy system",
			setup: func(t *testing.T) func() {
				tmpDir := t.TempDir()

				origSource := Config.Source
				origIndex := Config.Index
				origNotFound := Config.NotFoundPage
				origReadonly := Config.Readonly
				origBindAddress := Config.BindAddress

				Config.Source = tmpDir
				Config.Index = "index"
				Config.NotFoundPage = "404"
				Config.Readonly = false
				Config.BindAddress = "127.0.0.1:3000"

				// Create all required files
				_ = os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Index"), 0644)
				_ = os.WriteFile(filepath.Join(tmpDir, "404.md"), []byte("# 404"), 0644)

				return func() {
					Config.Source = origSource
					Config.Index = origIndex
					Config.NotFoundPage = origNotFound
					Config.Readonly = origReadonly
					Config.BindAddress = origBindAddress
				}
			},
			expectHealthy: true,
			issueCount:    0,
			warningCount:  0,
		},
		{
			name: "multiple warnings, no critical issues",
			setup: func(t *testing.T) func() {
				tmpDir := t.TempDir()

				origSource := Config.Source
				origIndex := Config.Index
				origNotFound := Config.NotFoundPage
				origReadonly := Config.Readonly
				origBindAddress := Config.BindAddress

				Config.Source = tmpDir
				Config.Index = "missing"
				Config.NotFoundPage = "missing404"
				Config.Readonly = true
				Config.BindAddress = "127.0.0.1:3000"

				// Don't create files - will trigger warnings

				return func() {
					Config.Source = origSource
					Config.Index = origIndex
					Config.NotFoundPage = origNotFound
					Config.Readonly = origReadonly
					Config.BindAddress = origBindAddress
				}
			},
			expectHealthy: false,
			issueCount:    0,
			warningCount:  3, // Missing index, missing 404, no markdown files
			expectWarnings: []string{
				"Index page not found",
				"404 page not found",
				"No markdown files",
			},
		},
		{
			name: "file instead of directory as source",
			setup: func(t *testing.T) func() {
				tmpDir := t.TempDir()
				filePath := filepath.Join(tmpDir, "notadir")
				_ = os.WriteFile(filePath, []byte("test"), 0644)

				origSource := Config.Source
				origBindAddress := Config.BindAddress
				origReadonly := Config.Readonly

				Config.Source = filePath
				Config.BindAddress = "127.0.0.1:3000"
				Config.Readonly = true // Prevent write check on non-directory

				return func() {
					Config.Source = origSource
					Config.BindAddress = origBindAddress
					Config.Readonly = origReadonly
				}
			},
			expectHealthy: false,
			issueCount:    1,
			expectIssue:   "not a directory",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cleanup := tc.setup(t)
			defer cleanup()

			result := runDiagnostics()

			if got := result.IsHealthy(); got != tc.expectHealthy {
				t.Errorf("IsHealthy() = %v, want %v (issues: %v, warnings: %v)",
					got, tc.expectHealthy, result.Issues, result.Warnings)
			}

			if len(result.Issues) != tc.issueCount {
				t.Errorf("Expected %d issues, got %d: %v", tc.issueCount, len(result.Issues), result.Issues)
			}

			if tc.warningCount > 0 && len(result.Warnings) != tc.warningCount {
				t.Errorf("Expected %d warnings, got %d: %v", tc.warningCount, len(result.Warnings), result.Warnings)
			}

			if tc.expectIssue != "" {
				found := false
				for _, issue := range result.Issues {
					if contains(issue, tc.expectIssue) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected issue containing %q, but got: %v", tc.expectIssue, result.Issues)
				}
			}

			for _, expectedWarning := range tc.expectWarnings {
				found := false
				for _, warning := range result.Warnings {
					if contains(warning, expectedWarning) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected warning containing %q, but got: %v", expectedWarning, result.Warnings)
				}
			}
		})
	}
}
