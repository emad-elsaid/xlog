package metadata

import (
	"bufio"
	"strings"

	"github.com/emad-elsaid/xlog"
	. "github.com/emad-elsaid/xlog"
	"gopkg.in/yaml.v2"
)

func init() {
	RegisterPreprocessor(stripYAML)
}

func stripYAML(content Markdown) Markdown {
	reader := strings.NewReader(string(content))

	scanner := bufio.NewScanner(reader)
	var body strings.Builder
	var y strings.Builder
	inFrontMatter := false
	frontMatterFound := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "---" {
			if inFrontMatter {
				// End of front matter
				inFrontMatter = false
				frontMatterFound = true
			} else if !frontMatterFound {
				// Start of front matter
				inFrontMatter = true
			} else {
				body.WriteString(line)
				body.WriteString("\n")
			}
		} else {
			if inFrontMatter {
				y.WriteString(line)
				y.WriteString("\n")
			} else if !inFrontMatter {
				body.WriteString(line)
				body.WriteString("\n")
			}

		}
	}

	if err := scanner.Err(); err != nil {
		// ignore error, and just return original content
		return content
	}

	var meta xlog.Metadata
	// only strip if valid yaml
	err := yaml.Unmarshal([]byte(y.String()), &meta)
	if err != nil {
		return content
	}
	return Markdown(body.String())
}
