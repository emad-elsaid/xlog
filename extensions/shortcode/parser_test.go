package shortcode_test

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/extensions/shortcode"
)

func TestShortCode(t *testing.T) {
	tcs := []struct {
		name          string
		input         string
		handlerOutput string
		output        string
	}{
		{
			name:          "page with one line",
			input:         "/test",
			handlerOutput: "output",
			output:        "output",
		},
		{
			name:          "short code with new line before it",
			input:         "\n/test",
			handlerOutput: "output",
			output:        "output",
		},
		{
			name:          "short code with new line after it",
			input:         "/test\n",
			handlerOutput: "output",
			output:        "output",
		},
		{
			name:          "short code with new line before and after it",
			input:         "\n/test\n",
			handlerOutput: "output",
			output:        "output",
		},
		{
			name:          "short code with space after",
			input:         "/test ",
			handlerOutput: "output",
			output:        "output",
		},
		{
			name:          "two short codes",
			input:         "/test\n\n/test",
			handlerOutput: "output",
			output:        "outputoutput",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			handler := func(xlog.Markdown) template.HTML { return template.HTML(tc.handlerOutput) }
			shortcode.ShortCode("test", handler)

			output := bytes.NewBufferString("")
			xlog.MarkDownRenderer.Convert([]byte(tc.input), output)
			if output.String() != tc.output {
				t.Errorf("input: %s\nexpected: %s\noutput: %s", tc.input, tc.output, output.String())
			}
		})
	}
}
