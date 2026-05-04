package xlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiagnosticResult_HasCriticalIssues(t *testing.T) {
	tests := []struct {
		name   string
		result DiagnosticResult
		want   bool
	}{
		{
			name:   "no issues",
			result: DiagnosticResult{Issues: []string{}, Warnings: []string{}},
			want:   false,
		},
		{
			name:   "has issues",
			result: DiagnosticResult{Issues: []string{"critical error"}, Warnings: []string{}},
			want:   true,
		},
		{
			name:   "only warnings",
			result: DiagnosticResult{Issues: []string{}, Warnings: []string{"warning"}},
			want:   false,
		},
		{
			name:   "both issues and warnings",
			result: DiagnosticResult{Issues: []string{"error"}, Warnings: []string{"warning"}},
			want:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.result.HasCriticalIssues()
			if got != tc.want {
				t.Errorf("HasCriticalIssues() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDiagnosticResult_IsHealthy(t *testing.T) {
	tests := []struct {
		name   string
		result DiagnosticResult
		want   bool
	}{
		{
			name:   "completely healthy",
			result: DiagnosticResult{Issues: []string{}, Warnings: []string{}},
			want:   true,
		},
		{
			name:   "has issues",
			result: DiagnosticResult{Issues: []string{"error"}, Warnings: []string{}},
			want:   false,
		},
		{
			name:   "has warnings",
			result: DiagnosticResult{Issues: []string{}, Warnings: []string{"warning"}},
			want:   false,
		},
		{
			name:   "has both",
			result: DiagnosticResult{Issues: []string{"error"}, Warnings: []string{"warning"}},
			want:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.result.IsHealthy()
			if got != tc.want {
				t.Errorf("IsHealthy() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRunDiagnostics(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T) string
		cleanup       func(t *testing.T, dir string)
		configSetup   func(t *testing.T, dir string)
		wantIssues    int
		wantWarnings  int
		checkIssues   func(t *testing.T, issues []string)
		checkWarnings func(t *testing.T, warnings []string)
	}{
		{
			name: "healthy configuration",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				// Create index page
				if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte("# Home"), 0644); err != nil {
					t.Fatal(err)
				}
				// Create 404 page
				if err := os.WriteFile(filepath.Join(dir, "404.md"), []byte("# Not Found"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			configSetup: func(t *testing.T, dir string) {
				Config.Source = dir
				Config.Index = "index"
				Config.NotFoundPage = "404"
				Config.BindAddress = "127.0.0.1:3000"
				Config.Readonly = false
			},
			wantIssues:   0,
			wantWarnings: 0,
		},
		{
			name: "missing source directory",
			setup: func(t *testing.T) string {
				return "/nonexistent/directory"
			},
			configSetup: func(t *testing.T, dir string) {
				Config.Source = dir
				Config.Index = "index"
				Config.NotFoundPage = "404"
				Config.BindAddress = "127.0.0.1:3000"
				Config.Readonly = false
			},
			wantIssues:   2,  // source doesn't exist + not writable (cascading failure)
			wantWarnings: -1, // don't check warnings (will vary)
			checkIssues: func(t *testing.T, issues []string) {
				// Should have at least one source directory related issue
				if len(issues) < 2 {
					t.Errorf("expected 2 issues, got %d: %v", len(issues), issues)
				}
			},
		},
		{
			name: "missing index page",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				// Create a markdown file but not index
				if err := os.WriteFile(filepath.Join(dir, "other.md"), []byte("# Other"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			configSetup: func(t *testing.T, dir string) {
				Config.Source = dir
				Config.Index = "index"
				Config.NotFoundPage = "404"
				Config.BindAddress = "127.0.0.1:3000"
				Config.Readonly = false
			},
			wantIssues:   0,
			wantWarnings: 2, // missing index and 404
		},
		{
			name: "no markdown files",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				// Create a non-markdown file
				if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("text"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			configSetup: func(t *testing.T, dir string) {
				Config.Source = dir
				Config.Index = "index"
				Config.NotFoundPage = "404"
				Config.BindAddress = "127.0.0.1:3000"
				Config.Readonly = false
			},
			wantIssues:   0,
			wantWarnings: 3, // no markdown files, missing index, missing 404
		},
		{
			name: "readonly mode skips write check",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte("# Home"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			configSetup: func(t *testing.T, dir string) {
				Config.Source = dir
				Config.Index = "index"
				Config.NotFoundPage = "404"
				Config.BindAddress = "127.0.0.1:3000"
				Config.Readonly = true
			},
			wantIssues:   0,
			wantWarnings: 1, // missing 404
		},
		{
			name: "empty bind address",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte("# Home"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			configSetup: func(t *testing.T, dir string) {
				Config.Source = dir
				Config.Index = "index"
				Config.NotFoundPage = "404"
				Config.BindAddress = ""
				Config.Readonly = false
			},
			wantIssues:   1,
			wantWarnings: 1, // missing 404
			checkIssues: func(t *testing.T, issues []string) {
				found := false
				for _, issue := range issues {
					if issue == "✗ Bind address is empty" {
						found = true
						break
					}
				}
				if !found {
					t.Error("expected bind address issue")
				}
			},
		},
		{
			name: "source is file not directory",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				filePath := filepath.Join(dir, "notadir.txt")
				if err := os.WriteFile(filePath, []byte("text"), 0644); err != nil {
					t.Fatal(err)
				}
				return filePath
			},
			configSetup: func(t *testing.T, path string) {
				Config.Source = path
				Config.Index = "index"
				Config.NotFoundPage = "404"
				Config.BindAddress = "127.0.0.1:3000"
				Config.Readonly = false
			},
			wantIssues:   2,  // not a directory + not writable (cascading failure)
			wantWarnings: -1, // don't check warnings (will vary)
			checkIssues: func(t *testing.T, issues []string) {
				// Should have at least 2 issues (not a dir + not writable)
				if len(issues) < 2 {
					t.Errorf("expected at least 2 issues, got %d: %v", len(issues), issues)
				}
			},
		},
		{
			name: "invalid theme produces warning",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte("# Home"), 0644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, "404.md"), []byte("# Not Found"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			configSetup: func(t *testing.T, dir string) {
				Config.Source = dir
				Config.Index = "index"
				Config.NotFoundPage = "404"
				Config.BindAddress = "127.0.0.1:3000"
				Config.Readonly = false
				Config.Theme = "invalid-theme"
			},
			wantIssues:   0,
			wantWarnings: 1,
			checkWarnings: func(t *testing.T, warnings []string) {
				found := false
				for _, warning := range warnings {
					if strings.Contains(warning, "Invalid theme 'invalid-theme'") {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected invalid theme warning, got: %v", warnings)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Save original config
			origConfig := Config

			// Setup test environment
			path := tc.setup(t)
			tc.configSetup(t, path)

			// Run diagnostics
			result := runDiagnostics()

			// Verify results
			if tc.wantIssues >= 0 && len(result.Issues) != tc.wantIssues {
				t.Errorf("got %d issues, want %d. Issues: %v", len(result.Issues), tc.wantIssues, result.Issues)
			}

			if tc.wantWarnings >= 0 && len(result.Warnings) != tc.wantWarnings {
				t.Errorf("got %d warnings, want %d. Warnings: %v", len(result.Warnings), tc.wantWarnings, result.Warnings)
			}

			// Additional checks if provided
			if tc.checkIssues != nil {
				tc.checkIssues(t, result.Issues)
			}
			if tc.checkWarnings != nil {
				tc.checkWarnings(t, result.Warnings)
			}

			// Restore original config
			Config = origConfig

			// Cleanup if provided
			if tc.cleanup != nil {
				tc.cleanup(t, path)
			}
		})
	}
}

func TestCheckSourceDirectory(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		setupFunc  func(t *testing.T) string
		wantIssues int
	}{
		{
			name: "valid directory",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			wantIssues: 0,
		},
		{
			name:       "nonexistent directory",
			source:     "/this/does/not/exist",
			wantIssues: 1,
		},
		{
			name: "file instead of directory",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				file := filepath.Join(dir, "file.txt")
				if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
					t.Fatal(err)
				}
				return file
			},
			wantIssues: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var source string
			if tc.setupFunc != nil {
				source = tc.setupFunc(t)
			} else {
				source = tc.source
			}

			origSource := Config.Source
			Config.Source = source

			issues := []string{}
			checkSourceDirectory(&issues)

			if len(issues) != tc.wantIssues {
				t.Errorf("got %d issues, want %d. Issues: %v", len(issues), tc.wantIssues, issues)
			}

			Config.Source = origSource
		})
	}
}

func TestCheckIndexPage(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(t *testing.T) (dir, indexName string)
		wantWarnings int
	}{
		{
			name: "index page exists",
			setupFunc: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte("# Home"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir, "index"
			},
			wantWarnings: 0,
		},
		{
			name: "index page missing",
			setupFunc: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				return dir, "index"
			},
			wantWarnings: 1,
		},
		{
			name: "custom index page exists",
			setupFunc: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "home.md"), []byte("# Home"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir, "home"
			},
			wantWarnings: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir, indexName := tc.setupFunc(t)

			origSource := Config.Source
			origIndex := Config.Index
			Config.Source = dir
			Config.Index = indexName

			warnings := []string{}
			checkIndexPage(&warnings)

			if len(warnings) != tc.wantWarnings {
				t.Errorf("got %d warnings, want %d. Warnings: %v", len(warnings), tc.wantWarnings, warnings)
			}

			Config.Source = origSource
			Config.Index = origIndex
		})
	}
}

func TestCheckNotFoundPage(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(t *testing.T) (dir, notFoundName string)
		wantWarnings int
	}{
		{
			name: "404 page exists",
			setupFunc: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "404.md"), []byte("# Not Found"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir, "404"
			},
			wantWarnings: 0,
		},
		{
			name: "404 page missing",
			setupFunc: func(t *testing.T) (string, string) {
				dir := t.TempDir()
				return dir, "404"
			},
			wantWarnings: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir, notFoundName := tc.setupFunc(t)

			origSource := Config.Source
			origNotFoundPage := Config.NotFoundPage
			Config.Source = dir
			Config.NotFoundPage = notFoundName

			warnings := []string{}
			checkNotFoundPage(&warnings)

			if len(warnings) != tc.wantWarnings {
				t.Errorf("got %d warnings, want %d. Warnings: %v", len(warnings), tc.wantWarnings, warnings)
			}

			Config.Source = origSource
			Config.NotFoundPage = origNotFoundPage
		})
	}
}

func TestCheckMarkdownFiles(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(t *testing.T) string
		wantWarnings int
	}{
		{
			name: "has markdown files",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "file.md"), []byte("# Test"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			wantWarnings: 0,
		},
		{
			name: "no markdown files",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("test"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			wantWarnings: 1,
		},
		{
			name: "empty directory",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			wantWarnings: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := tc.setupFunc(t)

			origSource := Config.Source
			Config.Source = dir

			warnings := []string{}
			checkMarkdownFiles(&warnings)

			if len(warnings) != tc.wantWarnings {
				t.Errorf("got %d warnings, want %d. Warnings: %v", len(warnings), tc.wantWarnings, warnings)
			}

			Config.Source = origSource
		})
	}
}

func TestCheckWritePermissions(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func(t *testing.T) string
		readonly   bool
		wantIssues int
	}{
		{
			name: "writable directory",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			readonly:   false,
			wantIssues: 0,
		},
		{
			name: "readonly mode skips check",
			setupFunc: func(t *testing.T) string {
				return t.TempDir()
			},
			readonly:   true,
			wantIssues: 0,
		},
		{
			name: "non-writable directory",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				// Make directory read-only
				if err := os.Chmod(dir, 0444); err != nil {
					t.Fatal(err)
				}
				// Restore permissions after test
				t.Cleanup(func() {
					_ = os.Chmod(dir, 0755)
				})
				return dir
			},
			readonly:   false,
			wantIssues: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := tc.setupFunc(t)

			origSource := Config.Source
			origReadonly := Config.Readonly
			Config.Source = dir
			Config.Readonly = tc.readonly

			issues := []string{}
			checkWritePermissions(&issues)

			if len(issues) != tc.wantIssues {
				t.Errorf("got %d issues, want %d. Issues: %v", len(issues), tc.wantIssues, issues)
			}

			Config.Source = origSource
			Config.Readonly = origReadonly
		})
	}
}

func TestCheckBindAddress(t *testing.T) {
	tests := []struct {
		name        string
		bindAddress string
		wantIssues  int
	}{
		{
			name:        "valid bind address",
			bindAddress: "127.0.0.1:3000",
			wantIssues:  0,
		},
		{
			name:        "empty bind address",
			bindAddress: "",
			wantIssues:  1,
		},
		{
			name:        "localhost with port",
			bindAddress: "localhost:8080",
			wantIssues:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			origBindAddress := Config.BindAddress
			Config.BindAddress = tc.bindAddress

			issues := []string{}
			checkBindAddress(&issues)

			if len(issues) != tc.wantIssues {
				t.Errorf("got %d issues, want %d. Issues: %v", len(issues), tc.wantIssues, issues)
			}

			Config.BindAddress = origBindAddress
		})
	}
}

func TestPrintDiagnosticSummary(t *testing.T) {
	tests := []struct {
		name            string
		issues          []string
		warnings        []string
		wantContains    []string
		wantNotContains []string
		wantExitCode    int
	}{
		{
			name:     "no issues or warnings",
			issues:   []string{},
			warnings: []string{},
			wantContains: []string{
				"✓ All checks passed!",
				"Your xlog configuration looks good.",
			},
			wantNotContains: []string{
				"CRITICAL ISSUES:",
				"WARNINGS:",
				"Please fix critical issues",
			},
			wantExitCode: 0,
		},
		{
			name:     "only warnings",
			issues:   []string{},
			warnings: []string{"⚠ Index page not found", "⚠ 404 page not found"},
			wantContains: []string{
				"WARNINGS:",
				"⚠ Index page not found",
				"⚠ 404 page not found",
				"Warnings noted. Xlog should run",
			},
			wantNotContains: []string{
				"CRITICAL ISSUES:",
				"Please fix critical issues",
			},
			wantExitCode: 0,
		},
		{
			name:     "only critical issues",
			issues:   []string{"✗ Source directory does not exist", "✗ Bind address is empty"},
			warnings: []string{},
			wantContains: []string{
				"CRITICAL ISSUES:",
				"✗ Source directory does not exist",
				"✗ Bind address is empty",
				"Please fix critical issues before running xlog.",
			},
			wantNotContains: []string{
				"WARNINGS:",
				"✓ All checks passed!",
			},
			wantExitCode: 1,
		},
		{
			name:     "both issues and warnings",
			issues:   []string{"✗ Source directory not writable"},
			warnings: []string{"⚠ Index page not found"},
			wantContains: []string{
				"CRITICAL ISSUES:",
				"✗ Source directory not writable",
				"WARNINGS:",
				"⚠ Index page not found",
				"Please fix critical issues before running xlog.",
			},
			wantExitCode: 1,
		},
		{
			name:     "multiple issues and warnings",
			issues:   []string{"✗ Issue 1", "✗ Issue 2", "✗ Issue 3"},
			warnings: []string{"⚠ Warning 1", "⚠ Warning 2"},
			wantContains: []string{
				"CRITICAL ISSUES:",
				"✗ Issue 1",
				"✗ Issue 2",
				"✗ Issue 3",
				"WARNINGS:",
				"⚠ Warning 1",
				"⚠ Warning 2",
			},
			wantExitCode: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Capture output
			var buf strings.Builder
			exitCode := formatDiagnosticSummary(&buf, tc.issues, tc.warnings)

			output := buf.String()

			// Check that expected strings are present
			for _, want := range tc.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output missing expected string %q\nGot:\n%s", want, output)
				}
			}

			// Check that unexpected strings are not present
			for _, notWant := range tc.wantNotContains {
				if strings.Contains(output, notWant) {
					t.Errorf("output contains unexpected string %q\nGot:\n%s", notWant, output)
				}
			}

			// Check exit code
			if exitCode != tc.wantExitCode {
				t.Errorf("got exit code %d, want %d", exitCode, tc.wantExitCode)
			}
		})
	}
}

func TestPrintDiagnosticSummary_OutputFormat(t *testing.T) {
	tests := []struct {
		name     string
		issues   []string
		warnings []string
		validate func(t *testing.T, output string)
	}{
		{
			name:     "issues appear before warnings",
			issues:   []string{"✗ Critical error"},
			warnings: []string{"⚠ Minor warning"},
			validate: func(t *testing.T, output string) {
				issuesIdx := strings.Index(output, "CRITICAL ISSUES:")
				warningsIdx := strings.Index(output, "WARNINGS:")
				if issuesIdx == -1 {
					t.Error("missing CRITICAL ISSUES section")
				}
				if warningsIdx == -1 {
					t.Error("missing WARNINGS section")
				}
				if issuesIdx >= warningsIdx {
					t.Error("CRITICAL ISSUES should appear before WARNINGS")
				}
			},
		},
		{
			name:     "issues are indented",
			issues:   []string{"✗ Error message"},
			warnings: []string{},
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "  ✗ Error message") {
					t.Error("issue should be indented with two spaces")
				}
			},
		},
		{
			name:     "warnings are indented",
			issues:   []string{},
			warnings: []string{"⚠ Warning message"},
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "  ⚠ Warning message") {
					t.Error("warning should be indented with two spaces")
				}
			},
		},
		{
			name:     "output starts with blank line",
			issues:   []string{"✗ Error"},
			warnings: []string{},
			validate: func(t *testing.T, output string) {
				if !strings.HasPrefix(output, "\n") {
					t.Error("output should start with a blank line")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf strings.Builder
			formatDiagnosticSummary(&buf, tc.issues, tc.warnings)
			tc.validate(t, buf.String())
		})
	}
}

func TestFormatDiagnosticSummaryEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		issues       []string
		warnings     []string
		wantContains []string
		wantExitCode int
	}{
		{
			name:     "handles very long issue messages",
			issues:   []string{"✗ Source directory path is extremely long and contains many nested subdirectories which might cause issues on some filesystems: /very/long/path/that/goes/on/and/on"},
			warnings: []string{},
			wantContains: []string{
				"CRITICAL ISSUES:",
				"/very/long/path/",
			},
			wantExitCode: 1,
		},
		{
			name:     "handles multiple lines in single issue",
			issues:   []string{"✗ Critical:\n  Nested detail line 1\n  Nested detail line 2"},
			warnings: []string{},
			wantContains: []string{
				"CRITICAL ISSUES:",
				"Nested detail line 1",
			},
			wantExitCode: 1,
		},
		{
			name:     "handles empty string issue",
			issues:   []string{""},
			warnings: []string{},
			wantContains: []string{
				"CRITICAL ISSUES:",
			},
			wantExitCode: 1,
		},
		{
			name:     "handles unicode in messages",
			issues:   []string{"✗ ファイルが見つかりません (file not found)"},
			warnings: []string{"⚠ 警告: Configuration parameter使用してください"},
			wantContains: []string{
				"CRITICAL ISSUES:",
				"ファイルが見つかりません",
				"WARNINGS:",
				"警告",
			},
			wantExitCode: 1,
		},
		{
			name:     "handles large number of issues",
			issues:   make([]string, 50),
			warnings: []string{},
			wantContains: []string{
				"CRITICAL ISSUES:",
				"Please fix critical issues",
			},
			wantExitCode: 1,
		},
		{
			name:     "handles large number of warnings without issues",
			issues:   []string{},
			warnings: make([]string, 30),
			wantContains: []string{
				"WARNINGS:",
				"Warnings noted",
			},
			wantExitCode: 0,
		},
		{
			name:     "handles messages with special characters",
			issues:   []string{"✗ Path contains special chars: ~!@#$%^&*()_+-={}[]|\\:\";<>?,./"},
			warnings: []string{},
			wantContains: []string{
				"special chars",
				"~!@#$%^&*",
			},
			wantExitCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Fill arrays with content if they were created with make()
			for i := range tt.issues {
				if tt.issues[i] == "" {
					tt.issues[i] = fmt.Sprintf("✗ Issue #%d", i+1)
				}
			}
			for i := range tt.warnings {
				if tt.warnings[i] == "" {
					tt.warnings[i] = fmt.Sprintf("⚠ Warning #%d", i+1)
				}
			}

			var buf strings.Builder
			exitCode := formatDiagnosticSummary(&buf, tt.issues, tt.warnings)

			output := buf.String()

			// Check expected strings
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output missing expected string %q", want)
				}
			}

			// Check exit code
			if exitCode != tt.wantExitCode {
				t.Errorf("got exit code %d, want %d", exitCode, tt.wantExitCode)
			}
		})
	}
}

func TestCheckThemeValue(t *testing.T) {
	tests := []struct {
		name         string
		theme        string
		wantWarnings int
		checkWarning func(t *testing.T, warnings []string)
	}{
		{
			name:         "empty theme is valid (uses system preference)",
			theme:        "",
			wantWarnings: 0,
		},
		{
			name:         "light theme is valid",
			theme:        "light",
			wantWarnings: 0,
		},
		{
			name:         "dark theme is valid",
			theme:        "dark",
			wantWarnings: 0,
		},
		{
			name:         "invalid theme produces warning",
			theme:        "blue",
			wantWarnings: 1,
			checkWarning: func(t *testing.T, warnings []string) {
				if len(warnings) != 1 {
					t.Fatalf("expected 1 warning, got %d", len(warnings))
				}
				want := "⚠ Invalid theme 'blue' (valid options: light, dark). Will fall back to system preference."
				if warnings[0] != want {
					t.Errorf("got warning %q, want %q", warnings[0], want)
				}
			},
		},
		{
			name:         "uppercase theme produces warning",
			theme:        "Light",
			wantWarnings: 1,
			checkWarning: func(t *testing.T, warnings []string) {
				if !strings.Contains(warnings[0], "Invalid theme 'Light'") {
					t.Errorf("warning should mention invalid theme 'Light', got: %q", warnings[0])
				}
			},
		},
		{
			name:         "random value produces warning",
			theme:        "rainbow",
			wantWarnings: 1,
		},
		{
			name:         "numeric value produces warning",
			theme:        "123",
			wantWarnings: 1,
		},
		{
			name:         "special characters produce warning",
			theme:        "light-dark",
			wantWarnings: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			origTheme := Config.Theme
			Config.Theme = tc.theme

			warnings := []string{}
			checkThemeValue(&warnings)

			if len(warnings) != tc.wantWarnings {
				t.Errorf("got %d warnings, want %d. Warnings: %v", len(warnings), tc.wantWarnings, warnings)
			}

			if tc.checkWarning != nil {
				tc.checkWarning(t, warnings)
			}

			Config.Theme = origTheme
		})
	}
}

func TestDoctor_Integration(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(t *testing.T) string
		wantExitCode int
		wantOutput   []string
	}{
		{
			name: "healthy configuration exits with code 0",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte("# Home"), 0644); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, "404.md"), []byte("# Not Found"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			wantExitCode: 0,
			wantOutput: []string{
				"✓ All checks passed!",
			},
		},
		{
			name: "missing source exits with code 1",
			setupFunc: func(t *testing.T) string {
				return "/nonexistent/directory/that/does/not/exist"
			},
			wantExitCode: 1,
			wantOutput: []string{
				"CRITICAL ISSUES:",
				"Please fix critical issues",
			},
		},
		{
			name: "warnings only exits with code 0",
			setupFunc: func(t *testing.T) string {
				dir := t.TempDir()
				if err := os.WriteFile(filepath.Join(dir, "other.md"), []byte("# Other"), 0644); err != nil {
					t.Fatal(err)
				}
				return dir
			},
			wantExitCode: 0,
			wantOutput: []string{
				"WARNINGS:",
				"Warnings noted",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := tc.setupFunc(t)

			// Save and restore original config
			origConfig := Config
			t.Cleanup(func() { Config = origConfig })

			// Configure for test
			Config.Source = dir
			Config.Index = "index"
			Config.NotFoundPage = "404"
			Config.BindAddress = "127.0.0.1:3000"
			Config.Readonly = false

			// Capture output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run Doctor and expect it to exit
			exitCode := -1
			done := make(chan struct{})
			go func() {
				defer func() {
					if r := recover(); r != nil {
						// Doctor calls os.Exit which we can't catch directly in tests
						// So we use the testable formatDiagnosticSummary pathway
						close(done)
					}
				}()
				result := runDiagnostics()
				var buf strings.Builder
				exitCode = formatDiagnosticSummary(&buf, result.Issues, result.Warnings)
				_, _ = w.Write([]byte(buf.String()))
				close(done)
			}()

			<-done
			w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var output strings.Builder
			_, _ = io.Copy(&output, r)
			outputStr := output.String()

			// Verify exit code
			if exitCode != tc.wantExitCode {
				t.Errorf("exit code = %d, want %d", exitCode, tc.wantExitCode)
			}

			// Verify output contains expected strings
			for _, want := range tc.wantOutput {
				if !strings.Contains(outputStr, want) {
					t.Errorf("output missing expected string %q\nGot:\n%s", want, outputStr)
				}
			}
		})
	}
}
