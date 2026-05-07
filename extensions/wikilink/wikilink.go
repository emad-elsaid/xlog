package wikilink

import (
	"regexp"
	"strings"

	"github.com/emad-elsaid/xlog"
)

const extensionName = "wikilink"

// Pattern to match wiki links: [[Page Name]] or [[Subdirectory/Page Name]].
// Also matches escaped wiki links: \[[Page Name]].
var wikiLinkRegex = regexp.MustCompile(`\\?\[\[([^\]]+)\]\]`)

func init() {
	xlog.RegisterExtension(WikiLink{})
}

type WikiLink struct{}

func (WikiLink) Name() string { return extensionName }
func (WikiLink) Init() {
	xlog.RegisterPreprocessor(preprocessor)
}

// preprocessor converts [[Page Name]] syntax to markdown links [Page Name](/Page_Name).
// Escaped syntax \[[Page Name]] is left as literal [[Page Name]] (backslash removed).
func preprocessor(md xlog.Markdown) xlog.Markdown {
	content := string(md)

	// Replace all wiki links with markdown links
	result := wikiLinkRegex.ReplaceAllStringFunc(content, func(match string) string {
		// Check if this is an escaped wiki link (starts with backslash)
		if strings.HasPrefix(match, `\`) {
			// Remove the escape backslash, leave the rest as-is
			return match[1:]
		}

		// Extract the page name from [[Page Name]]
		pageName := strings.TrimSpace(match[2 : len(match)-2])

		// Convert page name to URL path
		// Replace spaces with underscores for URL compatibility
		urlPath := strings.ReplaceAll(pageName, " ", "_")

		// Create markdown link: [Display Name](/URL_Path)
		return "[" + pageName + "](/" + urlPath + ")"
	})

	return xlog.Markdown(result)
}
