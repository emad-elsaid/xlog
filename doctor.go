package xlog

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// DiagnosticResult holds the results of running diagnostics.
type DiagnosticResult struct {
	Issues   []string
	Warnings []string
}

// HasCriticalIssues returns true if there are any critical issues.
func (d DiagnosticResult) HasCriticalIssues() bool {
	return len(d.Issues) > 0
}

// IsHealthy returns true if there are no issues or warnings.
func (d DiagnosticResult) IsHealthy() bool {
	return len(d.Issues) == 0 && len(d.Warnings) == 0
}

// Doctor runs diagnostic checks on the xlog configuration and environment.
// It reports potential issues that could affect xlog's operation.
// This is the CLI entry point that calls os.Exit for critical issues.
func Doctor() {
	slog.Info("Running xlog diagnostics...")
	result := runDiagnostics()
	printDiagnosticSummary(result.Issues, result.Warnings)
}

// runDiagnostics performs all diagnostic checks and returns the results.
// This function is testable as it doesn't call os.Exit.
func runDiagnostics() DiagnosticResult {
	issues := []string{}
	warnings := []string{}

	checkSourceDirectory(&issues)
	checkIndexPage(&warnings)
	checkNotFoundPage(&warnings)
	checkMarkdownFiles(&warnings)
	checkWritePermissions(&issues)
	checkBindAddress(&issues)

	return DiagnosticResult{
		Issues:   issues,
		Warnings: warnings,
	}
}

func checkSourceDirectory(issues *[]string) {
	stat, err := os.Stat(Config.Source)
	switch {
	case err != nil && os.IsNotExist(err):
		*issues = append(*issues, fmt.Sprintf("✗ Source directory does not exist: %s", Config.Source))
	case err != nil:
		*issues = append(*issues, fmt.Sprintf("✗ Cannot access source directory: %v", err))
	case !stat.IsDir():
		*issues = append(*issues, fmt.Sprintf("✗ Source path is not a directory: %s", Config.Source))
	default:
		slog.Info("✓ Source directory accessible", "path", Config.Source)
	}
}

func checkIndexPage(warnings *[]string) {
	indexPath := filepath.Join(Config.Source, Config.Index+".md")
	_, err := os.Stat(indexPath)
	switch {
	case err != nil && os.IsNotExist(err):
		*warnings = append(*warnings, fmt.Sprintf("⚠ Index page not found: %s (xlog will show error on homepage)", indexPath))
	case err != nil:
		*warnings = append(*warnings, fmt.Sprintf("⚠ Cannot check index page: %v", err))
	default:
		slog.Info("✓ Index page exists", "path", indexPath)
	}
}

func checkNotFoundPage(warnings *[]string) {
	notFoundPath := filepath.Join(Config.Source, Config.NotFoundPage+".md")
	if _, err := os.Stat(notFoundPath); err != nil {
		if os.IsNotExist(err) {
			*warnings = append(*warnings, fmt.Sprintf("⚠ 404 page not found: %s (xlog will use default error page)", notFoundPath))
		}
	} else {
		slog.Info("✓ 404 page exists", "path", notFoundPath)
	}
}

func checkMarkdownFiles(warnings *[]string) {
	mdFiles, err := filepath.Glob(filepath.Join(Config.Source, "*.md"))
	switch {
	case err != nil:
		*warnings = append(*warnings, fmt.Sprintf("⚠ Cannot scan for markdown files: %v", err))
	case len(mdFiles) == 0:
		*warnings = append(*warnings, fmt.Sprintf("⚠ No markdown files (*.md) found in source directory: %s", Config.Source))
	default:
		slog.Info("✓ Found markdown files", "count", len(mdFiles))
	}
}

func checkWritePermissions(issues *[]string) {
	if Config.Readonly {
		slog.Info("✓ Running in readonly mode (write operations disabled)")
		return
	}

	testFile := filepath.Join(Config.Source, ".xlog-write-test")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		*issues = append(*issues, fmt.Sprintf("✗ Source directory not writable (required when not in readonly mode): %v", err))
	} else {
		_ = os.Remove(testFile)
		slog.Info("✓ Source directory is writable")
	}
}

func checkBindAddress(issues *[]string) {
	if Config.BindAddress == "" {
		*issues = append(*issues, "✗ Bind address is empty")
	} else {
		slog.Info("✓ Bind address configured", "address", Config.BindAddress)
	}
}

func printDiagnosticSummary(issues, warnings []string) {
	fmt.Println()
	if len(issues) == 0 && len(warnings) == 0 {
		fmt.Println("✓ All checks passed! Your xlog configuration looks good.")
		return
	}

	if len(issues) > 0 {
		fmt.Println("CRITICAL ISSUES:")
		for _, issue := range issues {
			fmt.Println("  " + issue)
		}
		fmt.Println()
	}

	if len(warnings) > 0 {
		fmt.Println("WARNINGS:")
		for _, warning := range warnings {
			fmt.Println("  " + warning)
		}
		fmt.Println()
	}

	if len(issues) > 0 {
		fmt.Println("Please fix critical issues before running xlog.")
		os.Exit(1)
	} else {
		fmt.Println("Warnings noted. Xlog should run, but you may want to address these.")
	}
}
