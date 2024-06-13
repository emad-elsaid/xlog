package metadata

import (
	"bufio"
	"strings"

	. "github.com/emad-elsaid/xlog"
)

func init() {
	RegisterPreprocessor(stripYAML)
}

// TODO move to an extension. see link_preview
func stripYAML(content Markdown) Markdown {
	reader := strings.NewReader(string(content))

	scanner := bufio.NewScanner(reader)
	var body strings.Builder
	inFrontMatter := false

	for scanner.Scan() {
		line := scanner.Text()

		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "---" {
			if inFrontMatter {
				// End of front matter
				inFrontMatter = false
			} else {
				// Start of front matter
				inFrontMatter = true
			}
		} else if !inFrontMatter {
			// Append lines not part of the front matter
			body.WriteString(line)
			body.WriteString("\n")
		}
	}

	if err := scanner.Err(); err != nil {
		// ignore error, and just return original content
		return content
	}
	return Markdown(body.String())
}
