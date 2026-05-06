package xlog

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/emad-elsaid/xlog/markdown/ast"
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
	checkThemeValue(&warnings)
	checkBrokenLinks(&warnings)
	checkOrphanPages(&warnings)
	checkDuplicateContent(&warnings)

	// Configuration validation
	checkConfigurationFlags(&warnings)

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

const (
	themeDark  = "dark"
	themeLight = "light"
)

func checkThemeValue(warnings *[]string) {
	if Config.Theme == "" {
		slog.Info("✓ Theme not set (will use system preference)")
		return
	}

	validThemes := []string{themeLight, themeDark}
	for _, valid := range validThemes {
		if Config.Theme == valid {
			slog.Info("✓ Theme is valid", "theme", Config.Theme)
			return
		}
	}

	*warnings = append(*warnings, fmt.Sprintf("⚠ Invalid theme '%s' (valid options: light, dark). Will fall back to system preference.", Config.Theme))
}

func checkBrokenLinks(warnings *[]string) {
	broken := FindBrokenLinks(context.Background())

	if len(broken) == 0 {
		slog.Info("✓ No broken internal links found")
		return
	}

	// Group by source page to create a concise warning message
	bySource := make(map[string]int)
	for _, bl := range broken {
		bySource[bl.SourcePage]++
	}

	totalBroken := len(broken)
	affectedPages := len(bySource)

	if totalBroken == 1 {
		*warnings = append(*warnings, fmt.Sprintf("⚠ Found 1 broken internal link in %d page(s). Run with -list flag to see details.", affectedPages))
	} else {
		*warnings = append(*warnings, fmt.Sprintf("⚠ Found %d broken internal link(s) in %d page(s). Run with -list flag to see details.", totalBroken, affectedPages))
	}

	slog.Warn("Broken internal links detected", "total", totalBroken, "affected_pages", affectedPages)
}

// findOrphanedPages returns a list of page names that have no incoming links.
func findOrphanedPages(ctx context.Context) []string {
	allPages := Pages(ctx)
	incomingLinks := make(map[string]int)

	// Build incoming link counts
	for _, p := range allPages {
		select {
		case <-ctx.Done():
			return nil
		default:
			content := p.Content()

			// Extract implicit [[page]] links
			pageLinkPattern := `\[\[([^\]]+)\]\]`
			matches := findPageLinks(string(content), pageLinkPattern)
			for _, pageName := range matches {
				incomingLinks[pageName]++
			}

			// Extract explicit markdown links from AST
			_, tree := p.AST()
			markdownLinks := FindAllInAST[*ast.Link](tree)
			for _, link := range markdownLinks {
				dest := string(link.Destination)
				if dest == "" || isExternalLink(dest) || strings.HasPrefix(dest, "#") {
					continue
				}
				targetPageName := linkToPageName(dest)
				incomingLinks[targetPageName]++
			}
		}
	}

	// Find pages with no incoming links
	var orphans []string
	for _, p := range allPages {
		if incomingLinks[p.Name()] == 0 {
			orphans = append(orphans, p.Name())
		}
	}

	return orphans
}

func checkOrphanPages(warnings *[]string) {
	orphans := findOrphanedPages(context.Background())

	if len(orphans) == 0 {
		slog.Info("✓ No orphaned pages (all pages have incoming links)")
		return
	}

	if len(orphans) == 1 {
		*warnings = append(*warnings, fmt.Sprintf("⚠ Found 1 orphaned page with no incoming links: %s", orphans[0]))
	} else {
		pageList := strings.Join(orphans, ", ")
		*warnings = append(*warnings, fmt.Sprintf("⚠ Found %d orphaned page(s) with no incoming links: %s", len(orphans), pageList))
	}

	slog.Warn("Orphaned pages detected", "count", len(orphans), "pages", orphans)
}

func checkDuplicateContent(warnings *[]string) {
	duplicates := findDuplicateContent(context.Background())

	if len(duplicates) == 0 {
		slog.Info("✓ No duplicate or similar content found")
		return
	}

	// Count total pairs
	totalPairs := len(duplicates)

	if totalPairs == 1 {
		*warnings = append(*warnings, "⚠ Found 1 pair of pages with duplicate or similar content. Review and consolidate if appropriate.")
	} else {
		*warnings = append(*warnings, fmt.Sprintf("⚠ Found %d pair(s) of pages with duplicate or similar content. Review and consolidate if appropriate.", totalPairs))
	}

	slog.Warn("Duplicate content detected", "pairs", totalPairs)
}

// DuplicatePair represents two pages with duplicate or highly similar content.
type DuplicatePair struct {
	Page1      string
	Page2      string
	Similarity float64
}

// findDuplicateContent scans all pages and identifies pairs with duplicate or highly similar content.
// Similarity threshold is 0.9 (90% match).
func findDuplicateContent(ctx context.Context) []DuplicatePair {
	const similarityThreshold = 0.9
	duplicates := []DuplicatePair{}
	seen := make(map[string]bool)

	// Collect all pages with their normalized content
	type pageContent struct {
		name    string
		content string
	}
	var pages []pageContent
	var mu sync.Mutex

	MapPage(ctx, func(p Page) Page {
		// Skip index page
		if p.Name() == Config.Index {
			return nil
		}

		normalized := normalizeContent(string(p.Content()))
		if normalized != "" {
			mu.Lock()
			pages = append(pages, pageContent{
				name:    p.Name(),
				content: normalized,
			})
			mu.Unlock()
		}
		return nil
	})

	// Compare each pair of pages
	for i := 0; i < len(pages); i++ {
		for j := i + 1; j < len(pages); j++ {
			p1, p2 := pages[i], pages[j]

			// Create unique key for this pair
			pairKey := p1.name + "|" + p2.name
			if seen[pairKey] {
				continue
			}

			similarity := calculateSimilarity(p1.content, p2.content)
			if similarity >= similarityThreshold {
				duplicates = append(duplicates, DuplicatePair{
					Page1:      p1.name,
					Page2:      p2.name,
					Similarity: similarity,
				})
				seen[pairKey] = true
			}
		}
	}

	return duplicates
}

// normalizeContent normalizes page content for comparison by:
//   - Converting to lowercase
//   - Normalizing whitespace (multiple spaces/newlines to single space)
//   - Trimming leading/trailing whitespace.
func normalizeContent(content string) string {
	// Convert to lowercase
	normalized := strings.ToLower(content)

	// Replace multiple whitespace with single space
	normalized = strings.Join(strings.Fields(normalized), " ")

	return strings.TrimSpace(normalized)
}

// calculateSimilarity computes a similarity score between two strings.
// Returns a value between 0.0 (completely different) and 1.0 (identical).
// Uses a simple Jaccard similarity on word sets.
func calculateSimilarity(s1, s2 string) float64 {
	// Handle edge cases
	if s1 == s2 {
		return 1.0
	}
	if s1 == "" || s2 == "" {
		return 0.0
	}

	// Split into word sets
	words1 := strings.Fields(s1)
	words2 := strings.Fields(s2)

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	// Create word frequency maps
	set1 := make(map[string]int)
	set2 := make(map[string]int)

	for _, w := range words1 {
		set1[w]++
	}
	for _, w := range words2 {
		set2[w]++
	}

	// Calculate intersection and union
	intersection := 0
	for word, count1 := range set1 {
		if count2, exists := set2[word]; exists {
			if count1 < count2 {
				intersection += count1
			} else {
				intersection += count2
			}
		}
	}

	union := len(words1) + len(words2) - intersection

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

func printDiagnosticSummary(issues, warnings []string) {
	exitCode := formatDiagnosticSummary(os.Stdout, issues, warnings)
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

// formatDiagnosticSummary formats and writes the diagnostic summary to the given writer.
// It returns the exit code that should be used (0 for success, 1 for critical issues).
// This function is testable as it doesn't call os.Exit.
func formatDiagnosticSummary(w io.Writer, issues, warnings []string) int {
	_, _ = fmt.Fprintln(w)
	if len(issues) == 0 && len(warnings) == 0 {
		_, _ = fmt.Fprintln(w, "✓ All checks passed! Your xlog configuration looks good.")
		return 0
	}

	if len(issues) > 0 {
		_, _ = fmt.Fprintln(w, "CRITICAL ISSUES:")
		for _, issue := range issues {
			_, _ = fmt.Fprintln(w, "  "+issue)
		}
		_, _ = fmt.Fprintln(w)
	}

	if len(warnings) > 0 {
		_, _ = fmt.Fprintln(w, "WARNINGS:")
		for _, warning := range warnings {
			_, _ = fmt.Fprintln(w, "  "+warning)
		}
		_, _ = fmt.Fprintln(w)
	}

	if len(issues) > 0 {
		_, _ = fmt.Fprintln(w, "Please fix critical issues before running xlog.")
		return 1
	}

	_, _ = fmt.Fprintln(w, "Warnings noted. Xlog should run, but you may want to address these.")
	return 0
}

// checkGPGConfiguration validates the GPG configuration.
func checkGPGConfiguration(gpgFlag string, warnings *[]string) {
	if gpgFlag == "" {
		return
	}

	// Check if gpg binary exists
	if _, err := exec.LookPath("gpg"); err != nil {
		*warnings = append(*warnings, "⚠ gpg flag is set but gpg binary not found in PATH. Install gpg or remove -gpg flag.")
	}
}

// checkEditorCommand validates the editor command configuration.
func checkEditorCommand(editor string, warnings *[]string) {
	if editor == "" {
		return
	}

	// Extract binary name from command (handle cases like "emacs -nw")
	fields := strings.Fields(editor)
	if len(fields) == 0 {
		return
	}

	// Check if editor binary exists
	if _, err := exec.LookPath(fields[0]); err != nil {
		*warnings = append(*warnings, fmt.Sprintf("⚠ editor command '%s' not found in PATH. Pages may fail to open for editing.", fields[0]))
	}
}

// checkGitHubURLFormat validates the GitHub URL format.
func checkGitHubURLFormat(githubURL string, warnings *[]string) {
	if githubURL == "" {
		return
	}

	// Basic URL validation
	if !strings.HasPrefix(githubURL, "http://") && !strings.HasPrefix(githubURL, "https://") {
		*warnings = append(*warnings, fmt.Sprintf("⚠ github.url should start with http:// or https://: %s", githubURL))
		return
	}

	// Recommend HTTPS for GitHub
	if strings.HasPrefix(githubURL, "http://") {
		*warnings = append(*warnings, "⚠ github.url should use https:// instead of http:// for security")
	}
}

// checkBindPortPermissions validates bind port permissions.
func checkBindPortPermissions(bindAddr string, warnings *[]string) {
	// Extract port from address (format: "host:port")
	parts := strings.Split(bindAddr, ":")
	if len(parts) < 2 {
		return
	}

	portStr := parts[len(parts)-1]
	var port int
	if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
		return
	}

	// Warn about privileged ports (< 1024)
	if port < 1024 && port > 0 {
		*warnings = append(*warnings, fmt.Sprintf("⚠ Binding to privileged port %d requires root/sudo. Use port >= 1024 or run with elevated privileges.", port))
	}
}

// checkConfigurationFlags validates command-line flag configurations.
func checkConfigurationFlags(warnings *[]string) {
	// Check GPG configuration
	if gpgFlag := flag.Lookup("gpg"); gpgFlag != nil && gpgFlag.Value.String() != "" {
		checkGPGConfiguration(gpgFlag.Value.String(), warnings)
	}

	// Check editor command
	if editorFlag := flag.Lookup("editor"); editorFlag != nil && editorFlag.Value.String() != "" {
		checkEditorCommand(editorFlag.Value.String(), warnings)
	}

	// Check GitHub URL
	if githubURLFlag := flag.Lookup("github.url"); githubURLFlag != nil && githubURLFlag.Value.String() != "" {
		checkGitHubURLFormat(githubURLFlag.Value.String(), warnings)
	}

	// Check bind port
	checkBindPortPermissions(Config.BindAddress, warnings)
}
