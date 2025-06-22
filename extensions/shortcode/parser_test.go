package shortcode_test

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/extensions/shortcode"
	"github.com/emad-elsaid/xlog/markdown"
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

	md := markdown.New(
		markdown.WithExtensions(
			&shortcode.ShortCodeEx{},
		),
	)

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			handler := func(xlog.Markdown) template.HTML { return template.HTML(tc.handlerOutput) }
			shortcode.RegisterShortCode("test", shortcode.ShortCode{Render: handler, Default: ""})

			output := bytes.NewBufferString("")
			md.Convert([]byte(tc.input), output)
			if output.String() != tc.output {
				t.Errorf("input: %s\nexpected: %s\noutput: %s", tc.input, tc.output, output.String())
			}
		})
	}
}
