package xlog

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/emad-elsaid/xlog/markdown/ast"
)

// BrokenLink represents a broken internal link found in a page.
type BrokenLink struct {
	SourcePage      string // Page containing the broken link
	TargetPage      string // Referenced page that doesn't exist
	LinkDestination string // Original link text
}

// FindBrokenLinks scans all pages and returns broken internal links.
// An internal link is considered broken if it references a page that doesn't exist.
// Only checks links that appear to be internal (not http/https URLs).
func FindBrokenLinks(ctx context.Context) []BrokenLink {
	results := MapPage(ctx, func(p Page) []BrokenLink {
		var broken []BrokenLink

		_, tree := p.AST()
		links := FindAllInAST[*ast.Link](tree)

		for _, link := range links {
			dest := string(link.Destination)

			// Skip empty links
			if dest == "" {
				continue
			}

			// Skip external links (http, https, mailto, etc.)
			if isExternalLink(dest) {
				continue
			}

			// Skip anchors and fragments
			if strings.HasPrefix(dest, "#") {
				continue
			}

			// Extract the page name from the link
			targetPageName := linkToPageName(dest)

			// Check if the target page exists
			targetPage := NewPage(targetPageName)
			if !targetPage.Exists() {
				broken = append(broken, BrokenLink{
					SourcePage:      p.Name(),
					TargetPage:      targetPageName,
					LinkDestination: dest,
				})
			}
		}

		return broken
	})

	// Flatten results from all pages
	var allBroken []BrokenLink
	for _, pageResults := range results {
		allBroken = append(allBroken, pageResults...)
	}

	return allBroken
}

// isExternalLink checks if a link destination is external (http/https/mailto/etc.)
func isExternalLink(dest string) bool {
	// Quick check for common schemes
	if strings.HasPrefix(dest, "http://") ||
		strings.HasPrefix(dest, "https://") ||
		strings.HasPrefix(dest, "mailto:") ||
		strings.HasPrefix(dest, "ftp://") ||
		strings.HasPrefix(dest, "//") {
		return true
	}

	// Parse as URL to catch other schemes
	if u, err := url.Parse(dest); err == nil && u.Scheme != "" {
		return true
	}

	return false
}

// linkToPageName converts a link destination to a page name.
// Removes leading slashes, removes fragment/query, and removes .md extension.
func linkToPageName(dest string) string {
	// Remove fragment (e.g., "page#section" -> "page")
	if idx := strings.IndexByte(dest, '#'); idx >= 0 {
		dest = dest[:idx]
	}

	// Remove query string (e.g., "page?query=foo" -> "page")
	if idx := strings.IndexByte(dest, '?'); idx >= 0 {
		dest = dest[:idx]
	}

	// Clean the path (remove redundant slashes, resolve .. etc.)
	dest = path.Clean(dest)

	// Remove leading slash after cleaning
	dest = strings.TrimPrefix(dest, "/")

	// Remove .md extension if present
	dest = strings.TrimSuffix(dest, ".md")

	return dest
}

// PrintBrokenLinks formats and prints broken links to stdout.
func PrintBrokenLinks(broken []BrokenLink) {
	if len(broken) == 0 {
		fmt.Println("✓ No broken internal links found")
		return
	}

	fmt.Printf("Found %d broken internal link(s):\n\n", len(broken))

	// Group by source page for better readability
	bySource := make(map[string][]BrokenLink)
	for _, bl := range broken {
		bySource[bl.SourcePage] = append(bySource[bl.SourcePage], bl)
	}

	for sourcePage, links := range bySource {
		fmt.Printf("In page '%s':\n", sourcePage)
		for _, bl := range links {
			fmt.Printf("  → Broken link to '%s' (link: %s)\n",
				bl.TargetPage, bl.LinkDestination)
		}
		fmt.Println()
	}
}
