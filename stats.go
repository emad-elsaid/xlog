package xlog

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/emad-elsaid/xlog/markdown/ast"
)

// GardenStats holds statistics about the digital garden.
type GardenStats struct {
	TotalPages      int
	TotalWords      int
	AvgWordsPerPage int
	TotalLinks      int
	OrphanedPages   int
	HubPages        []string
}

// Stats displays statistics about the digital garden.
// This is the CLI entry point.
func Stats(ctx context.Context) {
	printStats(ctx, os.Stdout)
}

// calculateStats computes statistics across all pages.
// This function is testable as it doesn't perform I/O.
func calculateStats(ctx context.Context) GardenStats {
	allPages := Pages(ctx)

	totalPages := len(allPages)
	if totalPages == 0 {
		return GardenStats{}
	}

	totalWords := 0
	linkCounts := make(map[string]int)
	incomingLinks := make(map[string]int)

	for _, p := range allPages {
		select {
		case <-ctx.Done():
			return GardenStats{}
		default:
			content := p.Content()
			totalWords += countWords(string(content))

			// Extract implicit [[page]] links from markdown content using regex
			// This approach avoids circular dependencies with extension packages
			pageLinkPattern := `\[\[([^\]]+)\]\]`
			matches := findPageLinks(string(content), pageLinkPattern)

			// Extract explicit markdown links from AST
			_, tree := p.AST()
			markdownLinks := FindAllInAST[*ast.Link](tree)

			internalLinkCount := len(matches)

			// Count page link references
			for _, pageName := range matches {
				incomingLinks[pageName]++
			}

			// Count explicit markdown links to internal pages
			for _, link := range markdownLinks {
				dest := string(link.Destination)

				// Skip empty, external, and fragment-only links
				if dest == "" || isExternalLink(dest) || strings.HasPrefix(dest, "#") {
					continue
				}

				targetPageName := linkToPageName(dest)
				incomingLinks[targetPageName]++
				internalLinkCount++
			}

			// Store total outgoing links for this page
			linkCounts[p.Name()] = internalLinkCount
		}
	}

	avgWords := 0
	if totalPages > 0 {
		avgWords = totalWords / totalPages
	}

	// Calculate total links
	totalLinks := 0
	for _, count := range linkCounts {
		totalLinks += count
	}

	// Find orphaned pages (no incoming links)
	orphanedCount := 0
	for _, p := range allPages {
		if incomingLinks[p.Name()] == 0 {
			orphanedCount++
		}
	}

	// Find hub pages (top 3 pages with most incoming links)
	hubPages := findHubPages(incomingLinks, 3)

	return GardenStats{
		TotalPages:      totalPages,
		TotalWords:      totalWords,
		AvgWordsPerPage: avgWords,
		TotalLinks:      totalLinks,
		OrphanedPages:   orphanedCount,
		HubPages:        hubPages,
	}
}

// findHubPages returns the top N pages with the most incoming links.
func findHubPages(incomingLinks map[string]int, topN int) []string {
	type pageLink struct {
		name  string
		count int
	}

	var pages []pageLink
	for name, count := range incomingLinks {
		if count > 0 {
			pages = append(pages, pageLink{name, count})
		}
	}

	// Sort by count descending using stdlib sort for O(n log n) performance
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].count > pages[j].count
	})

	// Take top N
	result := []string{}
	limit := topN
	if len(pages) < limit {
		limit = len(pages)
	}

	for i := 0; i < limit; i++ {
		result = append(result, pages[i].name)
	}

	return result
}

// findPageLinks extracts page names from [[page]] syntax in markdown content.
func findPageLinks(content, pattern string) []string {
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(content, -1)

	var pageNames []string
	for _, match := range matches {
		if len(match) > 1 {
			pageNames = append(pageNames, match[1])
		}
	}

	return pageNames
}

// countWords counts the number of words in text.
// Words are sequences of non-whitespace characters.
func countWords(text string) int {
	if text == "" {
		return 0
	}

	fields := strings.Fields(text)
	return len(fields)
}

// printStats formats and writes statistics to the given writer.
func printStats(ctx context.Context, w io.Writer) {
	stats := calculateStats(ctx)

	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, "═══════════════════════════════════════")
	_, _ = fmt.Fprintln(w, "    Digital Garden Statistics")
	_, _ = fmt.Fprintln(w, "═══════════════════════════════════════")
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintf(w, "  Total Pages:            %d\n", stats.TotalPages)
	_, _ = fmt.Fprintf(w, "  Total Words:            %d\n", stats.TotalWords)
	_, _ = fmt.Fprintf(w, "  Average Words per Page: %d\n", stats.AvgWordsPerPage)
	_, _ = fmt.Fprintf(w, "  Total Links:            %d\n", stats.TotalLinks)
	_, _ = fmt.Fprintf(w, "  Orphaned Pages:         %d\n", stats.OrphanedPages)
	_, _ = fmt.Fprintln(w)

	if len(stats.HubPages) > 0 {
		_, _ = fmt.Fprintln(w, "  Hub Pages (most connected):")
		for i, page := range stats.HubPages {
			_, _ = fmt.Fprintf(w, "    %d. %s\n", i+1, page)
		}
		_, _ = fmt.Fprintln(w)
	}

	_, _ = fmt.Fprintln(w, "═══════════════════════════════════════")
	_, _ = fmt.Fprintln(w)
}
