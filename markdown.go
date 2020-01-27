package xlog

import (
	"bytes"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
)

func renderMarkdown(content string) string {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			highlighting.Highlighting,
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return err.Error()
	}

	return buf.String()
}
