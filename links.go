package xlog

import (
	"net/url"
	"path"
	"strings"
)

// isExternalLink checks if a link destination points to an external resource.
// Returns true for URLs with http://, https://, mailto:, ftp://, protocol-relative (//),
// or any other URL scheme.
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
// This is used to normalize links for page lookups.
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
