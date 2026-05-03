package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestCodeSpanParser_Trigger(t *testing.T) {
	parser := &codeSpanParser{}
	triggers := parser.Trigger()

	expected := []byte{'`'}
	if len(triggers) != len(expected) {
		t.Fatalf("Trigger() returned %d characters, want %d", len(triggers), len(expected))
	}

	if triggers[0] != expected[0] {
		t.Errorf("Trigger()[0] = %c, want %c", triggers[0], expected[0])
	}
}

func TestNewCodeSpanParser(t *testing.T) {
	parser := NewCodeSpanParser()
	if parser == nil {
		t.Fatal("NewCodeSpanParser() returned nil")
	}

	_, ok := parser.(*codeSpanParser)
	if !ok {
		t.Fatalf("NewCodeSpanParser() returned %T, want *codeSpanParser", parser)
	}
}

func TestCodeSpanParser_Parse(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedType   string
		expectedText   string
		shouldHaveCode bool
	}{
		{
			name:           "simple code span",
			input:          "`code`",
			expectedType:   "CodeSpan",
			expectedText:   "code",
			shouldHaveCode: true,
		},
		{
			name:           "code span with space padding",
			input:          "` code `",
			expectedType:   "CodeSpan",
			expectedText:   "code",
			shouldHaveCode: true,
		},
		{
			name:           "code span without closing",
			input:          "`code",
			expectedType:   "TextSegment",
			expectedText:   "`",
			shouldHaveCode: false,
		},
		{
			name:           "double backtick code span",
			input:          "``code``",
			expectedType:   "CodeSpan",
			expectedText:   "code",
			shouldHaveCode: true,
		},
		{
			name:           "code span with single backtick inside",
			input:          "`` code ` inside ``",
			expectedType:   "CodeSpan",
			expectedText:   "code ` inside",
			shouldHaveCode: true,
		},
		{
			name:           "mismatched backtick count",
			input:          "``code`",
			expectedType:   "TextSegment",
			expectedText:   "``",
			shouldHaveCode: false,
		},
		{
			name:           "code span with special chars",
			input:          "`<html>&nbsp;</html>`",
			expectedType:   "CodeSpan",
			expectedText:   "<html>&nbsp;</html>",
			shouldHaveCode: true,
		},
		{
			name:           "empty code span",
			input:          "``",
			expectedType:   "TextSegment",
			expectedText:   "``",
			shouldHaveCode: false,
		},
		{
			name:           "code span with newline",
			input:          "`code\nline2`",
			expectedType:   "CodeSpan",
			expectedText:   "code\nline2",
			shouldHaveCode: true,
		},
		{
			name:           "triple backtick with double close",
			input:          "```code``",
			expectedType:   "TextSegment",
			expectedText:   "```",
			shouldHaveCode: false,
		},
		{
			name:           "code span at end of line",
			input:          "`code`\n",
			expectedType:   "CodeSpan",
			expectedText:   "code",
			shouldHaveCode: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &codeSpanParser{}
			parent := ast.NewDocument()
			reader := text.NewReader([]byte(tc.input))
			context := NewContext()

			result := parser.Parse(parent, reader, context)

			if result == nil {
				t.Fatal("Parse() returned nil")
			}

			switch tc.expectedType {
			case "CodeSpan":
				codeSpan, ok := result.(*ast.CodeSpan)
				if !ok {
					t.Fatalf("Parse() returned %T, want *ast.CodeSpan", result)
				}

				if !tc.shouldHaveCode {
					t.Error("Expected no code span, but got one")
				}

				// Extract text from code span
				source := []byte(tc.input)
				actualText := extractCodeSpanText(t, codeSpan, source)

				if actualText != tc.expectedText {
					t.Errorf("CodeSpan text = %q, want %q", actualText, tc.expectedText)
				}

			case "TextSegment":
				textSeg, ok := result.(*ast.Text)
				if !ok {
					t.Fatalf("Parse() returned %T, want *ast.Text", result)
				}

				if tc.shouldHaveCode {
					t.Error("Expected code span, but got text segment")
				}

				source := []byte(tc.input)
				actualText := string(textSeg.Segment.Value(source))

				if actualText != tc.expectedText {
					t.Errorf("TextSegment text = %q, want %q", actualText, tc.expectedText)
				}
			}
		})
	}
}

func TestCodeSpanParser_Parse_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType string
	}{
		{
			name:         "multiple consecutive backticks",
			input:        "````",
			expectedType: "TextSegment",
		},
		{
			name:         "code span with only spaces",
			input:        "`   `",
			expectedType: "CodeSpan",
		},
		{
			name:         "single space trimming",
			input:        "` a `",
			expectedType: "CodeSpan",
		},
		{
			name:         "newline trimming",
			input:        "`\ncode\n`",
			expectedType: "CodeSpan",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &codeSpanParser{}
			parent := ast.NewDocument()
			reader := text.NewReader([]byte(tc.input))
			context := NewContext()

			result := parser.Parse(parent, reader, context)

			if result == nil {
				t.Fatal("Parse() returned nil")
			}

			switch tc.expectedType {
			case "CodeSpan":
				if _, ok := result.(*ast.CodeSpan); !ok {
					t.Errorf("Parse() returned %T, want *ast.CodeSpan", result)
				}
			case "TextSegment":
				if _, ok := result.(*ast.Text); !ok {
					t.Errorf("Parse() returned %T, want *ast.Text", result)
				}
			}
		})
	}
}

func TestIsSpaceOrNewline(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{
			name:     "space is space or newline",
			input:    ' ',
			expected: true,
		},
		{
			name:     "newline is space or newline",
			input:    '\n',
			expected: true,
		},
		{
			name:     "tab is not space or newline",
			input:    '\t',
			expected: false,
		},
		{
			name:     "letter is not space or newline",
			input:    'a',
			expected: false,
		},
		{
			name:     "backtick is not space or newline",
			input:    '`',
			expected: false,
		},
		{
			name:     "carriage return is not space or newline",
			input:    '\r',
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isSpaceOrNewline(tc.input)
			if result != tc.expected {
				t.Errorf("isSpaceOrNewline(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

// Helper function to extract text from CodeSpan node.
func extractCodeSpanText(t *testing.T, node *ast.CodeSpan, source []byte) string {
	t.Helper()

	if node.IsBlank(source) {
		return ""
	}

	var result []byte
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		if textNode, ok := child.(*ast.Text); ok {
			result = append(result, textNode.Segment.Value(source)...)
		}
	}
	return string(result)
}
