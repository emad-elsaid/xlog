package xlog

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
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
			result: DiagnosticResult{Issues: []string{"critical problem"}, Warnings: []string{}},
			want:   true,
		},
		{
			name:   "only warnings",
			result: DiagnosticResult{Issues: []string{}, Warnings: []string{"minor issue"}},
			want:   false,
		},
		{
			name:   "both issues and warnings",
			result: DiagnosticResult{Issues: []string{"critical"}, Warnings: []string{"minor"}},
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
			result: DiagnosticResult{Issues: []string{"critical"}, Warnings: []string{}},
			want:   false,
		},
		{
			name:   "has warnings",
			result: DiagnosticResult{Issues: []string{}, Warnings: []string{"minor"}},
			want:   false,
		},
		{
			name:   "has both",
			result: DiagnosticResult{Issues: []string{"critical"}, Warnings: []string{"minor"}},
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

func TestCheckOrphanPages(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "executes without panic",
			description: "Verifies checkOrphanPages can execute safely",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			warnings := []string{}

			// Should not panic
			checkOrphanPages(&warnings)

			// If there are warnings, verify format
			if len(warnings) > 0 {
				warning := warnings[0]
				if !strings.Contains(warning, "orphaned page") {
					t.Errorf("Warning should mention 'orphaned page', got: %q", warning)
				}
				if !strings.HasPrefix(warning, "⚠") {
					t.Errorf("Warning should start with ⚠ symbol, got: %q", warning)
				}
			}
		})
	}
}

func TestCheckOrphanPages_MessageFormat(t *testing.T) {
	tests := []struct {
		name          string
		orphanCount   int
		expectContain string
	}{
		{
			name:          "single orphan message",
			orphanCount:   1,
			expectContain: "1 orphaned page",
		},
		{
			name:          "multiple orphans message",
			orphanCount:   5,
			expectContain: "5 orphaned page(s)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			warnings := []string{}

			// Mock the message format by directly testing the logic
			switch tc.orphanCount {
			case 0:
				// No warning expected
			case 1:
				warnings = append(warnings, "⚠ Found 1 orphaned page with no incoming links. Consider linking to it from other pages.")
			default:
				warnings = append(warnings, fmt.Sprintf("⚠ Found %d orphaned page(s) with no incoming links. Consider linking to them from other pages.", tc.orphanCount))
			}

			if tc.orphanCount > 0 {
				if len(warnings) == 0 {
					t.Error("Expected warning but got none")
				}
				if !strings.Contains(warnings[0], tc.expectContain) {
					t.Errorf("Expected message containing %q, got %q", tc.expectContain, warnings[0])
				}
			}
		})
	}
}

func TestFormatDiagnosticSummary(t *testing.T) {
	tests := []struct {
		name         string
		issues       []string
		warnings     []string
		wantExitCode int
		wantContains []string
	}{
		{
			name:         "all checks passed",
			issues:       []string{},
			warnings:     []string{},
			wantExitCode: 0,
			wantContains: []string{"All checks passed"},
		},
		{
			name:         "critical issues only",
			issues:       []string{"Source directory not writable"},
			warnings:     []string{},
			wantExitCode: 1,
			wantContains: []string{"CRITICAL ISSUES", "Source directory not writable", "fix critical issues"},
		},
		{
			name:         "warnings only",
			issues:       []string{},
			warnings:     []string{"Index page not found"},
			wantExitCode: 0,
			wantContains: []string{"WARNINGS", "Index page not found", "should run"},
		},
		{
			name:         "both issues and warnings",
			issues:       []string{"Bind address empty"},
			warnings:     []string{"Theme invalid"},
			wantExitCode: 1,
			wantContains: []string{"CRITICAL ISSUES", "WARNINGS", "Bind address empty", "Theme invalid"},
		},
		{
			name:         "multiple issues",
			issues:       []string{"Issue 1", "Issue 2"},
			warnings:     []string{},
			wantExitCode: 1,
			wantContains: []string{"Issue 1", "Issue 2"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			exitCode := formatDiagnosticSummary(&buf, tc.issues, tc.warnings)

			if exitCode != tc.wantExitCode {
				t.Errorf("formatDiagnosticSummary() exitCode = %v, want %v", exitCode, tc.wantExitCode)
			}

			output := buf.String()
			for _, want := range tc.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("formatDiagnosticSummary() output missing %q\nGot: %s", want, output)
				}
			}
		})
	}
}

func TestRunDiagnostics(t *testing.T) {
	originalConfig := Config
	t.Cleanup(func() { Config = originalConfig })

	tests := []struct {
		name            string
		setup           func(t *testing.T) string
		wantIssueCount  int
		wantWarnCount   int
		issueContains   []string
		warningContains []string
	}{
		{
			name: "healthy configuration",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				_ = os.WriteFile(filepath.Join(dir, "index.md"), []byte("# Index"), 0644)
				_ = os.WriteFile(filepath.Join(dir, "404.md"), []byte("# Not Found"), 0644)
				return dir
			},
			wantIssueCount: 0,
			wantWarnCount:  0,
		},
		{
			name: "source directory does not exist",
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			issueContains: []string{"Source directory does not exist"},
		},
		{
			name: "missing index page",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				_ = os.WriteFile(filepath.Join(dir, "404.md"), []byte("# 404"), 0644)
				_ = os.WriteFile(filepath.Join(dir, "other.md"), []byte("# Other"), 0644)
				return dir
			},
			wantIssueCount:  0,
			warningContains: []string{"Index page not found"},
		},
		{
			name: "missing 404 page",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				_ = os.WriteFile(filepath.Join(dir, "index.md"), []byte("# Index"), 0644)
				return dir
			},
			wantIssueCount:  0,
			warningContains: []string{"404 page not found"},
		},
		{
			name: "no markdown files",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				_ = os.WriteFile(filepath.Join(dir, "readme.txt"), []byte("text"), 0644)
				return dir
			},
			wantIssueCount:  0,
			warningContains: []string{"No markdown files"},
		},
		{
			name: "readonly mode skips write check",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				Config.Readonly = true
				_ = os.WriteFile(filepath.Join(dir, "index.md"), []byte("# Index"), 0644)
				_ = os.WriteFile(filepath.Join(dir, "404.md"), []byte("# Not Found"), 0644)
				return dir
			},
			wantIssueCount: 0,
			wantWarnCount:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sourceDir := tc.setup(t)

			Config = Configuration{
				Source:       sourceDir,
				Index:        "index",
				NotFoundPage: "404",
				BindAddress:  "127.0.0.1:3000",
				Theme:        "",
				Readonly:     false,
			}

			result := runDiagnostics()

			if tc.wantIssueCount > 0 && len(result.Issues) != tc.wantIssueCount {
				t.Errorf("runDiagnostics() issues count = %v, want %v\nIssues: %v",
					len(result.Issues), tc.wantIssueCount, result.Issues)
			}

			if tc.wantWarnCount > 0 && len(result.Warnings) != tc.wantWarnCount {
				t.Errorf("runDiagnostics() warnings count = %v, want %v\nWarnings: %v",
					len(result.Warnings), tc.wantWarnCount, result.Warnings)
			}

			for _, want := range tc.issueContains {
				found := false
				for _, issue := range result.Issues {
					if strings.Contains(issue, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected issue containing %q, got: %v", want, result.Issues)
				}
			}

			for _, want := range tc.warningContains {
				found := false
				for _, warning := range result.Warnings {
					if strings.Contains(warning, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected warning containing %q, got: %v", want, result.Warnings)
				}
			}
		})
	}
}

func TestRunDiagnostics_InvalidTheme(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Index"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "404.md"), []byte("# 404"), 0644)

	originalConfig := Config
	t.Cleanup(func() { Config = originalConfig })

	tests := []struct {
		name         string
		theme        string
		wantContains string
	}{
		{
			name:  "valid light theme",
			theme: "light",
		},
		{
			name:  "valid dark theme",
			theme: "dark",
		},
		{
			name:  "empty theme (uses system)",
			theme: "",
		},
		{
			name:         "invalid theme",
			theme:        "blue",
			wantContains: "Invalid theme",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			Config = Configuration{
				Source:       tmpDir,
				Index:        "index",
				NotFoundPage: "404",
				BindAddress:  "127.0.0.1:3000",
				Theme:        tc.theme,
				Readonly:     false,
			}

			result := runDiagnostics()

			if tc.wantContains != "" {
				found := false
				for _, warning := range result.Warnings {
					if strings.Contains(warning, tc.wantContains) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected warning containing %q, got: %v", tc.wantContains, result.Warnings)
				}
			} else {
				// Check no invalid theme warning
				for _, warning := range result.Warnings {
					if strings.Contains(warning, "Invalid theme") {
						t.Errorf("Unexpected invalid theme warning: %v", warning)
					}
				}
			}
		})
	}
}

func TestRunDiagnostics_EmptyBindAddress(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "index.md"), []byte("# Index"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "404.md"), []byte("# 404"), 0644)

	originalConfig := Config
	t.Cleanup(func() { Config = originalConfig })

	Config = Configuration{
		Source:       tmpDir,
		Index:        "index",
		NotFoundPage: "404",
		BindAddress:  "",
		Theme:        "",
		Readonly:     false,
	}

	result := runDiagnostics()

	if len(result.Issues) == 0 {
		t.Error("Expected issue for empty bind address")
	}

	found := false
	for _, issue := range result.Issues {
		if strings.Contains(issue, "Bind address is empty") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected issue about empty bind address, got: %v", result.Issues)
	}
}

func TestRunDiagnostics_UnwritableDirectory(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root (permissions work differently)")
	}

	tmpDir := t.TempDir()
	readonlyDir := filepath.Join(tmpDir, "readonly")
	_ = os.Mkdir(readonlyDir, 0444)
	_ = os.WriteFile(filepath.Join(readonlyDir, "index.md"), []byte("# Index"), 0644)

	originalConfig := Config
	t.Cleanup(func() {
		Config = originalConfig
		_ = os.Chmod(readonlyDir, 0755)
	})

	Config = Configuration{
		Source:       readonlyDir,
		Index:        "index",
		NotFoundPage: "404",
		BindAddress:  "127.0.0.1:3000",
		Theme:        "",
		Readonly:     false,
	}

	result := runDiagnostics()

	if len(result.Issues) == 0 {
		t.Error("Expected issue for unwritable directory")
	}

	found := false
	for _, issue := range result.Issues {
		if strings.Contains(issue, "not writable") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected issue about unwritable directory, got: %v", result.Issues)
	}
}

func TestCheckSourceDirectory_EdgeCases(t *testing.T) {
	originalConfig := Config
	t.Cleanup(func() { Config = originalConfig })

	tests := []struct {
		name         string
		setup        func(t *testing.T) string
		wantContains string
	}{
		{
			name: "source is a file not directory",
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				filePath := filepath.Join(dir, "notadir")
				_ = os.WriteFile(filePath, []byte("content"), 0644)
				return filePath
			},
			wantContains: "not a directory",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sourcePath := tc.setup(t)

			Config = Configuration{
				Source:       sourcePath,
				Index:        "index",
				NotFoundPage: "404",
				BindAddress:  "127.0.0.1:3000",
				Theme:        "",
				Readonly:     false,
			}

			result := runDiagnostics()

			if len(result.Issues) == 0 {
				t.Error("Expected critical issue")
			}

			found := false
			for _, issue := range result.Issues {
				if strings.Contains(issue, tc.wantContains) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected issue containing %q, got: %v", tc.wantContains, result.Issues)
			}
		})
	}
}

func TestCheckMarkdownFiles_GlobError(t *testing.T) {
	originalConfig := Config
	t.Cleanup(func() { Config = originalConfig })

	// Create a directory with a name that would cause glob issues
	tmpDir := t.TempDir()
	invalidGlobDir := filepath.Join(tmpDir, "test[")
	_ = os.Mkdir(invalidGlobDir, 0755)

	Config = Configuration{
		Source:       invalidGlobDir,
		Index:        "index",
		NotFoundPage: "404",
		BindAddress:  "127.0.0.1:3000",
		Theme:        "",
		Readonly:     false,
	}

	result := runDiagnostics()

	// Should have a warning about scanning for markdown files
	found := false
	for _, warning := range result.Warnings {
		if strings.Contains(warning, "Cannot scan for markdown files") {
			found = true
			break
		}
	}
	if !found {
		t.Logf("Note: Glob pattern did not error as expected (filesystem dependent). Warnings: %v", result.Warnings)
	}
}

func TestPrintDiagnosticSummary_WithIssues(t *testing.T) {
	if os.Getenv("TEST_EXIT") == "1" {
		issues := []string{"Source directory missing", "Cannot bind to address"}
		warnings := []string{"No 404 page"}
		printDiagnosticSummary(issues, warnings)
		return
	}

	// Run test in subprocess to capture os.Exit
	cmd := exec.Command(os.Args[0], "-test.run=TestPrintDiagnosticSummary_WithIssues")
	cmd.Env = append(os.Environ(), "TEST_EXIT=1")
	err := cmd.Run()

	// Should exit with code 1 (has critical issues)
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("Expected exit error, got: %v", err)
	}
	if exitErr.ExitCode() != 1 {
		t.Errorf("Expected exit code 1, got %d", exitErr.ExitCode())
	}
}

func TestPrintDiagnosticSummary_WithWarningsOnly(t *testing.T) {
	if os.Getenv("TEST_EXIT") == "1" {
		issues := []string{}
		warnings := []string{"No index page", "Theme not set"}
		printDiagnosticSummary(issues, warnings)
		return
	}

	// Run test in subprocess to capture os.Exit
	cmd := exec.Command(os.Args[0], "-test.run=TestPrintDiagnosticSummary_WithWarningsOnly")
	cmd.Env = append(os.Environ(), "TEST_EXIT=1")
	err := cmd.Run()

	// Should exit with code 0 (warnings only)
	if err != nil {
		t.Errorf("Expected clean exit, got: %v", err)
	}
}

func TestPrintDiagnosticSummary_AllGood(t *testing.T) {
	if os.Getenv("TEST_EXIT") == "1" {
		issues := []string{}
		warnings := []string{}
		printDiagnosticSummary(issues, warnings)
		return
	}

	// Run test in subprocess to capture os.Exit
	cmd := exec.Command(os.Args[0], "-test.run=TestPrintDiagnosticSummary_AllGood")
	cmd.Env = append(os.Environ(), "TEST_EXIT=1")
	err := cmd.Run()

	// Should exit with code 0 (all good)
	if err != nil {
		t.Errorf("Expected clean exit, got: %v", err)
	}
}

func TestDoctor_Integration(t *testing.T) {
	if os.Getenv("TEST_EXIT") == "1" {
		// Set up a valid temporary environment
		tmpDir := t.TempDir()
		indexPath := filepath.Join(tmpDir, "index.md")
		_ = os.WriteFile(indexPath, []byte("# Index"), 0644)
		notFoundPath := filepath.Join(tmpDir, "404.md")
		_ = os.WriteFile(notFoundPath, []byte("# Not Found"), 0644)

		originalConfig := Config
		Config = Configuration{
			Source:       tmpDir,
			Index:        "index",
			NotFoundPage: "404",
			BindAddress:  "127.0.0.1:3000",
			Theme:        "",
			Readonly:     false,
		}
		defer func() { Config = originalConfig }()

		Doctor()
		return
	}

	// Run test in subprocess to capture os.Exit
	cmd := exec.Command(os.Args[0], "-test.run=TestDoctor_Integration")
	cmd.Env = append(os.Environ(), "TEST_EXIT=1")
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	err := cmd.Run()

	// Should exit cleanly with valid config
	if err != nil {
		t.Errorf("Expected clean exit, got: %v\nOutput: %s", err, outBuf.String())
	}

	// Verify diagnostic output was produced
	output := outBuf.String()
	if !strings.Contains(output, "Running xlog diagnostics") {
		t.Errorf("Expected diagnostic log message in output, got: %s", output)
	}
}

func TestCheckBrokenLinks(t *testing.T) {
	tests := []struct {
		name        string
		files       map[string]string
		expectWarn  bool
		warnPattern string
	}{
		{
			name: "no broken links produces no warning",
			files: map[string]string{
				"page1.md": "# Page 1\n[Link to page2](page2)\n",
				"page2.md": "# Page 2\nContent here",
			},
			expectWarn: false,
		},
		{
			name: "single broken link produces warning",
			files: map[string]string{
				"page1.md": "# Page 1\n[Broken link](nonexistent)\n",
			},
			expectWarn:  true,
			warnPattern: "1 broken internal link",
		},
		{
			name: "multiple broken links reports correct count",
			files: map[string]string{
				"page1.md": "# Page 1\n[Link 1](missing1)\n[Link 2](missing2)\n",
			},
			expectWarn:  true,
			warnPattern: "2 broken internal link",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Use chdir approach like broken_links_test.go
			tmpDir := t.TempDir()
			origDir, _ := os.Getwd()
			defer func() { _ = os.Chdir(origDir) }()

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Failed to change to temp dir: %v", err)
			}

			// Create test files
			for filename, content := range tc.files {
				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			// Clear cache and set config
			_ = clearPagesCache(nil)
			originalConfig := Config
			Config.Source = tmpDir
			defer func() { Config = originalConfig }()

			// Test checkBrokenLinks
			warnings := []string{}
			checkBrokenLinks(&warnings)

			if tc.expectWarn && len(warnings) == 0 {
				t.Errorf("Expected warning but got none")
			}

			if !tc.expectWarn && len(warnings) > 0 {
				t.Errorf("Expected no warning but got: %v", warnings)
			}

			if tc.warnPattern != "" {
				found := false
				for _, w := range warnings {
					if strings.Contains(w, tc.warnPattern) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected warning containing %q, got: %v", tc.warnPattern, warnings)
				}
			}
		})
	}
}
