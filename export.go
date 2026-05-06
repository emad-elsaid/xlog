package xlog

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/emad-elsaid/xlog/markdown/ast"
)

// PageMetadata represents exported page information in JSON format.
type PageMetadata struct {
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	ModTime   time.Time `json:"mod_time"`
	WordCount int       `json:"word_count"`
	Links     []string  `json:"links"`
}

// ExportJSON exports all pages' metadata as JSON to stdout.
func ExportJSON(ctx context.Context) {
	allPages := Pages(ctx)
	metadata := make([]PageMetadata, 0, len(allPages))

	for _, p := range allPages {
		// Extract links from page content
		content := string(p.Content())
		links := extractLinks(p, content)

		// Count words in content
		wordCount := countWords(content)

		meta := PageMetadata{
			Name:      p.Name(),
			Path:      p.FileName(),
			ModTime:   p.ModTime(),
			WordCount: wordCount,
			Links:     links,
		}
		metadata = append(metadata, meta)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metadata); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

// extractLinks extracts link targets from page content and AST.
func extractLinks(p Page, content string) []string {
	linkSet := make(map[string]bool)

	// Extract implicit [[page]] links
	matches := findPageLinks(content, pageLinkPattern)
	for _, match := range matches {
		linkSet[match] = true
	}

	// Extract markdown links from AST
	_, tree := p.AST()
	astLinks := FindAllInAST[*ast.Link](tree)
	for _, link := range astLinks {
		dest := string(link.Destination)
		if dest != "" {
			linkSet[dest] = true
		}
	}

	// Convert set to sorted slice
	links := make([]string, 0, len(linkSet))
	for link := range linkSet {
		links = append(links, link)
	}
	return links
}
