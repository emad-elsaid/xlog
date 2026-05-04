package extension

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
	"github.com/emad-elsaid/xlog/markdown/testutil"
)

func TestTypographer(t *testing.T) {
	markdown := markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			Typographer,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/typographer.txt", t, testutil.ParseCliCaseArg()...)
}

func TestUnclosedCounter_Reset(t *testing.T) {
	tests := []struct {
		name          string
		initialSingle int
		initialDouble int
	}{
		{
			name:          "reset zero values",
			initialSingle: 0,
			initialDouble: 0,
		},
		{
			name:          "reset non-zero single quote count",
			initialSingle: 3,
			initialDouble: 0,
		},
		{
			name:          "reset non-zero double quote count",
			initialSingle: 0,
			initialDouble: 5,
		},
		{
			name:          "reset both non-zero counts",
			initialSingle: 7,
			initialDouble: 4,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			counter := &unclosedCounter{
				Single: tc.initialSingle,
				Double: tc.initialDouble,
			}

			counter.Reset()

			if counter.Single != 0 {
				t.Errorf("After Reset(), Single = %d; want 0", counter.Single)
			}
			if counter.Double != 0 {
				t.Errorf("After Reset(), Double = %d; want 0", counter.Double)
			}
		})
	}
}
