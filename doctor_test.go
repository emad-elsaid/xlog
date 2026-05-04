package xlog

import (
	"bytes"
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
