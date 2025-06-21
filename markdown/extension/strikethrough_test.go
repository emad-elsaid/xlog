package extension

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
	"github.com/emad-elsaid/xlog/markdown/testutil"
)

func TestStrikethrough(t *testing.T) {
	markdown := markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			Strikethrough,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/strikethrough.txt", t, testutil.ParseCliCaseArg()...)
}
