package toc

import (
	"bytes"
	"os"
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	testsdata, err := os.ReadFile("testdata/tests.yaml")
	require.NoError(t, err)

	var tests []struct {
		Desc       string `yaml:"desc"`
		Give       string `yaml:"give"`
		Want       string `yaml:"want"`
		Title      string `yaml:"title"`
		TitleDepth int    `yaml:"titleDepth"`
		ListID     string `yaml:"listID"`
		TitleID    string `yaml:"titleID"`

		MinDepth int  `yaml:"minDepth"`
		MaxDepth int  `yaml:"maxDepth"`
		Compact  bool `yaml:"compact"`
	}
	require.NoError(t, yaml.Unmarshal(testsdata, &tests))

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Desc, func(t *testing.T) {
			t.Parallel()

			md := markdown.New(
				markdown.WithExtensions(&Extender{
					Title:      tt.Title,
					TitleDepth: tt.TitleDepth,
					MinDepth:   tt.MinDepth,
					MaxDepth:   tt.MaxDepth,
					Compact:    tt.Compact,
					ListID:     tt.ListID,
					TitleID:    tt.TitleID,
				}),
				markdown.WithParserOptions(parser.WithAutoHeadingID()),
			)

			var buf bytes.Buffer
			require.NoError(t, md.Convert([]byte(tt.Give), &buf))
			require.Equal(t, tt.Want, buf.String())
		})
	}
}
