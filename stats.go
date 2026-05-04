package xlog

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

// GardenStats holds statistics about the digital garden.
type GardenStats struct {
	TotalPages      int
	TotalWords      int
	AvgWordsPerPage int
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
	for _, p := range allPages {
		select {
		case <-ctx.Done():
			return GardenStats{}
		default:
			content := p.Content()
			totalWords += countWords(string(content))
		}
	}

	avgWords := 0
	if totalPages > 0 {
		avgWords = totalWords / totalPages
	}

	return GardenStats{
		TotalPages:      totalPages,
		TotalWords:      totalWords,
		AvgWordsPerPage: avgWords,
	}
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
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, "═══════════════════════════════════════")
	_, _ = fmt.Fprintln(w)
}
