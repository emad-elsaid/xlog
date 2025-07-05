package fuzz

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/extension"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func fuzz(f *testing.F) {
	f.Fuzz(func(t *testing.T, orig string) {
		markdown := markdown.New(
			markdown.WithParserOptions(
				parser.WithAutoHeadingID(),
				parser.WithAttribute(),
			),
			markdown.WithRendererOptions(
				html.WithUnsafe(),
				html.WithXHTML(),
			),
			markdown.WithExtensions(
				extension.DefinitionList,
				extension.Footnote,
				extension.GFM,
				extension.Typographer,
				extension.Linkify,
				extension.Table,
				extension.TaskList,
			),
		)
		var b bytes.Buffer
		if err := markdown.Convert(util.StringToReadOnlyBytes(orig), &b); err != nil {
			panic(err)
		}
	})
}

func FuzzDefault(f *testing.F) {
	bs, err := os.ReadFile("../_test/spec.json")
	if err != nil {
		panic(err)
	}
	var testCases []map[string]any
	if err := json.Unmarshal(bs, &testCases); err != nil {
		panic(err)
	}
	for _, c := range testCases {
		f.Add(c["markdown"])
	}
	fuzz(f)
}
