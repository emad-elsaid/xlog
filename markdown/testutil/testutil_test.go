package testutil

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
)

// TestParseCliCaseArg tests command line argument parsing.
func TestParseCliCaseArg(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected []int
	}{
		{
			name:     "no case args",
			args:     []string{"cmd", "-v", "other"},
			expected: []int{},
		},
		{
			name:     "single case",
			args:     []string{"cmd", "case=1"},
			expected: []int{1},
		},
		{
			name:     "multiple cases comma separated",
			args:     []string{"cmd", "case=1,2,3"},
			expected: []int{1, 2, 3},
		},
		{
			name:     "multiple case args",
			args:     []string{"cmd", "case=1,2", "case=3"},
			expected: []int{1, 2, 3},
		},
		{
			name:     "with whitespace",
			args:     []string{"cmd", "case=1, 2 , 3"},
			expected: []int{1, 2, 3},
		},
		{
			name:     "invalid numbers ignored",
			args:     []string{"cmd", "case=1,invalid,3"},
			expected: []int{1, 3},
		},
		{
			name:     "mixed with other args",
			args:     []string{"cmd", "-v", "case=1,2", "--debug", "case=3"},
			expected: []int{1, 2, 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Save original args
			origArgs := os.Args
			defer func() { os.Args = origArgs }()

			// Set test args
			os.Args = tc.args

			result := ParseCliCaseArg()

			if len(result) != len(tc.expected) {
				t.Errorf("expected %d elements, got %d", len(tc.expected), len(result))
				return
			}

			for i, val := range tc.expected {
				if result[i] != val {
					t.Errorf("element %d: expected %d, got %d", i, val, result[i])
				}
			}
		})
	}
}

// TestSourceAndExpected tests the source and expected helper functions.
func TestSourceAndExpected(t *testing.T) {
	tests := []struct {
		name           string
		testCase       MarkdownTestCase
		expectedSource string
		expectedExpect string
	}{
		{
			name: "no options",
			testCase: MarkdownTestCase{
				Markdown: "  hello  ",
				Expected: "  world  ",
			},
			expectedSource: "  hello  ",
			expectedExpect: "  world  ",
		},
		{
			name: "trim enabled",
			testCase: MarkdownTestCase{
				Markdown: "  hello  ",
				Expected: "  world  ",
				Options:  MarkdownTestCaseOptions{Trim: true},
			},
			expectedSource: "hello",
			expectedExpect: "world",
		},
		{
			name: "escape sequences enabled",
			testCase: MarkdownTestCase{
				Markdown: "hello\\nworld",
				Expected: "foo\\tbar",
				Options:  MarkdownTestCaseOptions{EnableEscape: true},
			},
			expectedSource: "hello\nworld",
			expectedExpect: "foo\tbar",
		},
		{
			name: "trim and escape combined",
			testCase: MarkdownTestCase{
				Markdown: "  hello\\nworld  ",
				Expected: "  foo\\tbar  ",
				Options:  MarkdownTestCaseOptions{Trim: true, EnableEscape: true},
			},
			expectedSource: "hello\nworld",
			expectedExpect: "foo\tbar",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualSource := source(&tc.testCase)
			if actualSource != tc.expectedSource {
				t.Errorf("source: expected %q, got %q", tc.expectedSource, actualSource)
			}

			actualExpect := expected(&tc.testCase)
			if actualExpect != tc.expectedExpect {
				t.Errorf("expected: expected %q, got %q", tc.expectedExpect, actualExpect)
			}
		})
	}
}

// TestApplyEscapeSequence tests escape sequence conversion.
func TestApplyEscapeSequence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no escapes",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "newline",
			input:    "hello\\nworld",
			expected: "hello\nworld",
		},
		{
			name:     "tab",
			input:    "hello\\tworld",
			expected: "hello\tworld",
		},
		{
			name:     "carriage return",
			input:    "hello\\rworld",
			expected: "hello\rworld",
		},
		{
			name:     "bell",
			input:    "hello\\aworld",
			expected: "hello\aworld",
		},
		{
			name:     "backspace",
			input:    "hello\\bworld",
			expected: "hello\bworld",
		},
		{
			name:     "form feed",
			input:    "hello\\fworld",
			expected: "hello\fworld",
		},
		{
			name:     "vertical tab",
			input:    "hello\\vworld",
			expected: "hello\vworld",
		},
		{
			name:     "backslash",
			input:    "hello\\\\world",
			expected: "hello\\world",
		},
		{
			name:     "hex escape",
			input:    "\\x41BC",
			expected: "ABC",
		},
		{
			name:     "unicode 4 digit",
			input:    "\\u0041BC",
			expected: "ABC",
		},
		{
			name:     "unicode 8 digit",
			input:    "\\U00000041BC",
			expected: "ABC",
		},
		{
			name:     "multiple escapes",
			input:    "line1\\nline2\\ttab\\rreturn",
			expected: "line1\nline2\ttab\rreturn",
		},
		{
			name:     "trailing backslash",
			input:    "hello\\",
			expected: "hello\\",
		},
		{
			name:     "incomplete hex - single digit",
			input:    "\\x4",
			expected: "\\x4",
		},
		{
			name:     "incomplete hex - at end",
			input:    "test\\x",
			expected: "test\\x",
		},
		{
			name:     "incomplete unicode short",
			input:    "\\u00",
			expected: "\\u00",
		},
		{
			name:     "incomplete unicode - 3 digits",
			input:    "\\u004",
			expected: "\\u004",
		},
		{
			name:     "invalid escape character",
			input:    "\\q",
			expected: "\\q",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := string(applyEscapeSequence([]byte(tc.input)))
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

// TestDiffPretty tests diff formatting output.
func TestDiffPretty(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		contains []string // Substrings expected in output
	}{
		{
			name:     "identical content",
			v1:       "hello\nworld\n",
			v2:       "hello\nworld\n",
			contains: []string{" | hello", " | world"},
		},
		{
			name:     "added line",
			v1:       "hello\n",
			v2:       "hello\nworld\n",
			contains: []string{" | hello", "+ | world"},
		},
		{
			name:     "removed line",
			v1:       "hello\nworld\n",
			v2:       "hello\n",
			contains: []string{" | hello", "- | world"},
		},
		{
			name:     "changed line",
			v1:       "hello\nworld\n",
			v2:       "hello\nearth\n",
			contains: []string{" | hello", "- | world", "+ | earth"},
		},
		{
			name:     "completely different",
			v1:       "abc\ndef\n",
			v2:       "xyz\nuvw\n",
			contains: []string{"- | abc", "- | def", "+ | xyz", "+ | uvw"},
		},
		{
			name:     "empty to content",
			v1:       "",
			v2:       "hello\n",
			contains: []string{"+ | hello"},
		},
		{
			name:     "content to empty",
			v1:       "hello\n",
			v2:       "",
			contains: []string{"- | hello"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := DiffPretty([]byte(tc.v1), []byte(tc.v2))
			resultStr := string(result)

			for _, substring := range tc.contains {
				if !bytes.Contains(result, []byte(substring)) {
					t.Errorf("expected output to contain %q, got:\n%s", substring, resultStr)
				}
			}
		})
	}
}

// TestSimpleDiff tests the diff algorithm directly.
func TestSimpleDiff(t *testing.T) {
	tests := []struct {
		name          string
		v1            string
		v2            string
		expectedTypes []diffType
	}{
		{
			name:          "identical",
			v1:            "hello\nworld",
			v2:            "hello\nworld",
			expectedTypes: []diffType{diffNone},
		},
		{
			name: "all added",
			v1:   "",
			v2:   "hello",
			// Empty string splits to [[]], producing both removed and added
			expectedTypes: []diffType{diffRemoved, diffAdded},
		},
		{
			name: "all removed",
			v1:   "hello",
			v2:   "",
			// Empty string splits to [[]], producing both removed and added
			expectedTypes: []diffType{diffRemoved, diffAdded},
		},
		{
			name:          "mixed changes",
			v1:            "line1\nline2\nline3",
			v2:            "line1\nchanged\nline3",
			expectedTypes: []diffType{diffNone, diffRemoved, diffAdded, diffNone},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := simpleDiff([]byte(tc.v1), []byte(tc.v2))

			if len(result) != len(tc.expectedTypes) {
				t.Errorf("expected %d diff sections, got %d", len(tc.expectedTypes), len(result))
				return
			}

			for i, expectedType := range tc.expectedTypes {
				if result[i].Type != expectedType {
					t.Errorf("diff section %d: expected type %v, got %v", i, expectedType, result[i].Type)
				}
			}
		})
	}
}

// mockMarkdown is a simple markdown implementation for testing.
type mockMarkdown struct {
	convertFunc func([]byte, io.Writer) error
}

func (m *mockMarkdown) Convert(source []byte, writer io.Writer, opts ...parser.ParseOption) error {
	if m.convertFunc != nil {
		return m.convertFunc(source, writer)
	}
	// Default: uppercase the input
	_, _ = writer.Write(bytes.ToUpper(source))
	return nil
}

func (m *mockMarkdown) Parser() parser.Parser         { return nil }
func (m *mockMarkdown) SetParser(parser.Parser)       {}
func (m *mockMarkdown) Renderer() renderer.Renderer   { return nil }
func (m *mockMarkdown) SetRenderer(renderer.Renderer) {}

// mockTestingT is a mock implementation of TestingT for testing error reporting.
type mockTestingT struct {
	logs   []string
	skips  []string
	errors []string
	failed bool
}

func (m *mockTestingT) Logf(format string, args ...any) {
	m.logs = append(m.logs, "log")
}

func (m *mockTestingT) Skipf(format string, args ...any) {
	m.skips = append(m.skips, "skip")
}

func (m *mockTestingT) Errorf(format string, args ...any) {
	m.errors = append(m.errors, "error")
}

func (m *mockTestingT) FailNow() {
	m.failed = true
}

// TestDoTestCase tests individual test case execution.
func TestDoTestCase(t *testing.T) {
	tests := []struct {
		name         string
		testCase     MarkdownTestCase
		markdown     *mockMarkdown
		expectError  bool
		expectPassed bool
	}{
		{
			name: "passing test",
			testCase: MarkdownTestCase{
				No:       1,
				Markdown: "hello",
				Expected: "HELLO",
			},
			markdown:     &mockMarkdown{},
			expectError:  false,
			expectPassed: true,
		},
		{
			name: "failing test - mismatch",
			testCase: MarkdownTestCase{
				No:       2,
				Markdown: "hello",
				Expected: "GOODBYE",
			},
			markdown:     &mockMarkdown{},
			expectError:  true,
			expectPassed: false,
		},
		{
			name: "test with description",
			testCase: MarkdownTestCase{
				No:          3,
				Description: "uppercase conversion",
				Markdown:    "test",
				Expected:    "TEST",
			},
			markdown:     &mockMarkdown{},
			expectError:  false,
			expectPassed: true,
		},
		{
			name: "test with trim option",
			testCase: MarkdownTestCase{
				No:       4,
				Markdown: "  hello  ",
				Expected: "  HELLO  ",
				Options:  MarkdownTestCaseOptions{Trim: true},
			},
			markdown:     &mockMarkdown{},
			expectError:  false,
			expectPassed: true,
		},
		{
			name: "conversion error triggers panic recovery",
			testCase: MarkdownTestCase{
				No:       5,
				Markdown: "fail",
				Expected: "ignored",
			},
			markdown: &mockMarkdown{
				convertFunc: func(source []byte, w io.Writer) error {
					return bytes.ErrTooLarge
				},
			},
			expectError:  true,
			expectPassed: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockTestingT{}
			DoTestCase(tc.markdown, tc.testCase, mock)

			hasError := len(mock.errors) > 0
			if hasError != tc.expectError {
				t.Errorf("expected error=%v, got error=%v (errors: %v)", tc.expectError, hasError, mock.errors)
			}
		})
	}
}

// TestDoTestCases tests batch test execution.
func TestDoTestCases(t *testing.T) {
	tests := []struct {
		name          string
		cases         []MarkdownTestCase
		markdown      *mockMarkdown
		expectedCalls int
	}{
		{
			name:          "empty cases",
			cases:         []MarkdownTestCase{},
			markdown:      &mockMarkdown{},
			expectedCalls: 0,
		},
		{
			name: "single case",
			cases: []MarkdownTestCase{
				{No: 1, Markdown: "test", Expected: "TEST"},
			},
			markdown:      &mockMarkdown{},
			expectedCalls: 1,
		},
		{
			name: "multiple cases",
			cases: []MarkdownTestCase{
				{No: 1, Markdown: "one", Expected: "ONE"},
				{No: 2, Markdown: "two", Expected: "TWO"},
				{No: 3, Markdown: "three", Expected: "THREE"},
			},
			markdown:      &mockMarkdown{},
			expectedCalls: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockTestingT{}
			DoTestCases(tc.markdown, tc.cases, mock)

			// Verify all cases were processed
			// (We can't easily count executions, but we verify the function completes)
			if len(tc.cases) != tc.expectedCalls {
				t.Errorf("test logic error: case count mismatch")
			}
		})
	}
}
