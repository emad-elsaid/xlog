package wikilink

import (
	"testing"

	"github.com/emad-elsaid/xlog"
	"github.com/stretchr/testify/assert"
)

func TestPreprocessor(t *testing.T) {
	tests := []struct {
		name     string
		input    xlog.Markdown
		expected xlog.Markdown
	}{
		{
			name:     "simple wiki link",
			input:    "Check out [[Home]]",
			expected: "Check out [Home](/Home)",
		},
		{
			name:     "wiki link with spaces",
			input:    "Learn about [[Machine Learning]]",
			expected: "Learn about [Machine Learning](/Machine_Learning)",
		},
		{
			name:     "wiki link with subdirectory",
			input:    "See [[docs/Installation]]",
			expected: "See [docs/Installation](/docs/Installation)",
		},
		{
			name:     "multiple wiki links",
			input:    "[[Home]] and [[About]] pages",
			expected: "[Home](/Home) and [About](/About) pages",
		},
		{
			name:     "wiki link in sentence",
			input:    "I'm learning about [[Neural Networks]] and [[Machine Learning]].",
			expected: "I'm learning about [Neural Networks](/Neural_Networks) and [Machine Learning](/Machine_Learning).",
		},
		{
			name:     "no wiki links",
			input:    "This is plain text with [regular](link)",
			expected: "This is plain text with [regular](link)",
		},
		{
			name:     "wiki link with special characters",
			input:    "[[C++ Programming]]",
			expected: "[C++ Programming](/C++_Programming)",
		},
		{
			name:     "empty content",
			input:    "",
			expected: "",
		},
		{
			name:     "wiki link at start and end",
			input:    "[[Start]] middle [[End]]",
			expected: "[Start](/Start) middle [End](/End)",
		},
		{
			name:     "escaped wiki link with backslash",
			input:    `\[[Page Name]]`,
			expected: "[[Page Name]]",
		},
		{
			name:     "mixed escaped and unescaped",
			input:    `\[[Escaped]] and [[Not Escaped]]`,
			expected: "[[Escaped]] and [Not Escaped](/Not_Escaped)",
		},
		{
			name:     "escaped in code description",
			input:    `Use \[[Page Name]] to link`,
			expected: "Use [[Page Name]] to link",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := preprocessor(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestWikiLinkRegex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single match",
			input:    "[[Page]]",
			expected: []string{"Page"},
		},
		{
			name:     "multiple matches",
			input:    "[[First]] and [[Second]]",
			expected: []string{"First", "Second"},
		},
		{
			name:     "no matches",
			input:    "No wiki links here",
			expected: []string{},
		},
		{
			name:     "nested brackets not matched",
			input:    "[[[Invalid]]]",
			expected: []string{"[Invalid"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matches := wikiLinkRegex.FindAllStringSubmatch(tc.input, -1)
			var results []string
			for _, match := range matches {
				if len(match) > 1 {
					results = append(results, match[1])
				}
			}
			if tc.expected == nil {
				tc.expected = []string{}
			}
			if results == nil {
				results = []string{}
			}
			assert.Equal(t, tc.expected, results)
		})
	}
}
