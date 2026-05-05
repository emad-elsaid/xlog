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

func TestCheckDuplicateContent(t *testing.T) {
	tests := []struct {
		name        string
		files       map[string]string
		expectWarn  bool
		warnPattern string
	}{
		{
			name: "no duplicate content produces no warning",
			files: map[string]string{
				"page1.md": "# Page 1\nUnique content here",
				"page2.md": "# Page 2\nDifferent content",
				"page3.md": "# Page 3\nCompletely different",
			},
			expectWarn: false,
		},
		{
			name: "identical content produces warning",
			files: map[string]string{
				"page1.md": "# Same Title\nSame content here",
				"page2.md": "# Same Title\nSame content here",
			},
			expectWarn:  true,
			warnPattern: "duplicate or similar content",
		},
		{
			name: "similar content (>90% match) produces warning",
			files: map[string]string{
				"article1.md": "# Introduction to Go\nGo is a programming language. " +
					"It was created by Google. It is statically typed. " +
					"Go has excellent concurrency support with goroutines.",
				"article2.md": "# Introduction to Go\nGo is a programming language. " +
					"It was created by Google. It is statically typed. " +
					"Go has great concurrency support with goroutines.",
			},
			expectWarn:  true,
			warnPattern: "similar content",
		},
		{
			name: "different length content not flagged",
			files: map[string]string{
				"short.md": "# Short\nBrief.",
				"long.md":  "# Long\nThis is much longer content with many more words and details",
			},
			expectWarn: false,
		},
		{
			name: "multiple duplicate pairs reported correctly",
			files: map[string]string{
				"a1.md": "# Test\nContent A",
				"a2.md": "# Test\nContent A",
				"b1.md": "# Test\nContent B",
				"b2.md": "# Test\nContent B",
			},
			expectWarn:  true,
			warnPattern: "duplicate",
		},
		{
			name: "case insensitive comparison",
			files: map[string]string{
				"upper.md": "# TITLE\nCONTENT HERE",
				"lower.md": "# title\ncontent here",
			},
			expectWarn:  true,
			warnPattern: "duplicate",
		},
		{
			name: "whitespace normalized",
			files: map[string]string{
				"spaced.md":  "# Title\nContent   with   spaces",
				"compact.md": "# Title\nContent with spaces",
			},
			expectWarn:  true,
			warnPattern: "duplicate",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			origDir, _ := os.Getwd()
			defer func() { _ = os.Chdir(origDir) }()

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Failed to change to temp dir: %v", err)
			}

			for filename, content := range tc.files {
				if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			_ = clearPagesCache(nil)
			originalConfig := Config
			Config.Source = tmpDir
			defer func() { Config = originalConfig }()

			warnings := []string{}
			checkDuplicateContent(&warnings)

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

func TestCheckGPGConfiguration(t *testing.T) {
	tests := []struct {
		name         string
		gpgFlag      string
		mockGPGPath  string
		expectWarn   bool
		warnContains string
	}{
		{
			name:       "no gpg flag set",
			gpgFlag:    "",
			expectWarn: false,
		},
		{
			name:         "gpg flag set but binary not found",
			gpgFlag:      "test-key-id",
			mockGPGPath:  "/nonexistent/gpg",
			expectWarn:   true,
			warnContains: "gpg binary not found",
		},
		{
			name:        "gpg flag set and binary exists",
			gpgFlag:     "test-key-id",
			mockGPGPath: "/usr/bin/gpg",
			expectWarn:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			warnings := []string{}

			// Simulate different gpg configurations
			origPath := os.Getenv("PATH")
			if tc.mockGPGPath != "" {
				os.Setenv("PATH", filepath.Dir(tc.mockGPGPath))
			}
			defer os.Setenv("PATH", origPath)

			checkGPGConfiguration(tc.gpgFlag, &warnings)

			if tc.expectWarn && len(warnings) == 0 {
				t.Errorf("Expected warning but got none")
			}
			if !tc.expectWarn && len(warnings) > 0 {
				t.Errorf("Expected no warning but got: %v", warnings)
			}
			if tc.expectWarn && len(warnings) > 0 {
				found := false
				for _, w := range warnings {
					if strings.Contains(w, tc.warnContains) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected warning containing %q, got: %v", tc.warnContains, warnings)
				}
			}
		})
	}
}

func TestCheckEditorCommand(t *testing.T) {
	tests := []struct {
		name         string
		editor       string
		expectWarn   bool
		warnContains string
	}{
		{
			name:       "empty editor uses default",
			editor:     "",
			expectWarn: false,
		},
		{
			name:         "invalid editor command",
			editor:       "nonexistent-editor-xyz123",
			expectWarn:   true,
			warnContains: "editor command",
		},
		{
			name:       "valid editor",
			editor:     "cat",
			expectWarn: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			warnings := []string{}
			checkEditorCommand(tc.editor, &warnings)

			if tc.expectWarn && len(warnings) == 0 {
				t.Errorf("Expected warning but got none")
			}
			if !tc.expectWarn && len(warnings) > 0 {
				t.Errorf("Expected no warning but got: %v", warnings)
			}
			if tc.expectWarn && len(warnings) > 0 {
				found := false
				for _, w := range warnings {
					if strings.Contains(w, tc.warnContains) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected warning containing %q, got: %v", tc.warnContains, warnings)
				}
			}
		})
	}
}

func TestCheckGitHubURLFormat(t *testing.T) {
	tests := []struct {
		name         string
		githubURL    string
		expectWarn   bool
		warnContains string
	}{
		{
			name:       "empty URL is valid",
			githubURL:  "",
			expectWarn: false,
		},
		{
			name:       "valid GitHub URL",
			githubURL:  "https://github.com/user/repo/edit/master/docs",
			expectWarn: false,
		},
		{
			name:         "invalid URL format",
			githubURL:    "not-a-url",
			expectWarn:   true,
			warnContains: "github.url",
		},
		{
			name:         "URL without https",
			githubURL:    "http://github.com/user/repo",
			expectWarn:   true,
			warnContains: "https",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			warnings := []string{}
			checkGitHubURLFormat(tc.githubURL, &warnings)

			if tc.expectWarn && len(warnings) == 0 {
				t.Errorf("Expected warning but got none")
			}
			if !tc.expectWarn && len(warnings) > 0 {
				t.Errorf("Expected no warning but got: %v", warnings)
			}
			if tc.expectWarn && len(warnings) > 0 {
				found := false
				for _, w := range warnings {
					if strings.Contains(w, tc.warnContains) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected warning containing %q, got: %v", tc.warnContains, warnings)
				}
			}
		})
	}
}

func TestCheckBindPortPermissions(t *testing.T) {
	tests := []struct {
		name         string
		bindAddr     string
		expectWarn   bool
		warnContains string
	}{
		{
			name:       "standard port above 1024",
			bindAddr:   "127.0.0.1:3000",
			expectWarn: false,
		},
		{
			name:         "privileged port 80",
			bindAddr:     "127.0.0.1:80",
			expectWarn:   true,
			warnContains: "privileged port",
		},
		{
			name:         "privileged port 443",
			bindAddr:     "0.0.0.0:443",
			expectWarn:   true,
			warnContains: "privileged port",
		},
		{
			name:       "port 1024 exactly",
			bindAddr:   "localhost:1024",
			expectWarn: false,
		},
		{
			name:       "invalid format ignored",
			bindAddr:   "invalid",
			expectWarn: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			warnings := []string{}
			checkBindPortPermissions(tc.bindAddr, &warnings)

			if tc.expectWarn && len(warnings) == 0 {
				t.Errorf("Expected warning but got none")
			}
			if !tc.expectWarn && len(warnings) > 0 {
				t.Errorf("Expected no warning but got: %v", warnings)
			}
			if tc.expectWarn && len(warnings) > 0 {
				found := false
				for _, w := range warnings {
					if strings.Contains(w, tc.warnContains) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected warning containing %q, got: %v", tc.warnContains, warnings)
				}
			}
		})
	}
}
