package extension

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
	"github.com/emad-elsaid/xlog/markdown/testutil"
)

func TestDefinitionList(t *testing.T) {
	markdown := markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			DefinitionList,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/definition_list.txt", t, testutil.ParseCliCaseArg()...)
}
