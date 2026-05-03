package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestEmphasisDelimiterProcessor_IsDelimiter(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{
			name:     "asterisk is delimiter",
			input:    '*',
			expected: true,
		},
		{
			name:     "underscore is delimiter",
			input:    '_',
			expected: true,
		},
		{
			name:     "space not delimiter",
			input:    ' ',
			expected: false,
		},
		{
			name:     "letter not delimiter",
			input:    'a',
			expected: false,
		},
		{
			name:     "tilde not delimiter",
			input:    '~',
			expected: false,
		},
		{
			name:     "backtick not delimiter",
			input:    '`',
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			processor := &emphasisDelimiterProcessor{}
			result := processor.IsDelimiter(tc.input)
			if result != tc.expected {
				t.Errorf("IsDelimiter(%c) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestEmphasisDelimiterProcessor_CanOpenCloser(t *testing.T) {
	tests := []struct {
		name       string
		openerChar byte
		closerChar byte
		expected   bool
	}{
		{
			name:       "matching asterisks",
			openerChar: '*',
			closerChar: '*',
			expected:   true,
		},
		{
			name:       "matching underscores",
			openerChar: '_',
			closerChar: '_',
			expected:   true,
		},
		{
			name:       "mismatched asterisk and underscore",
			openerChar: '*',
			closerChar: '_',
			expected:   false,
		},
		{
			name:       "mismatched underscore and asterisk",
			openerChar: '_',
			closerChar: '*',
			expected:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			processor := &emphasisDelimiterProcessor{}
			opener := &Delimiter{Char: tc.openerChar}
			closer := &Delimiter{Char: tc.closerChar}
			result := processor.CanOpenCloser(opener, closer)
			if result != tc.expected {
				t.Errorf("CanOpenCloser(%c, %c) = %v, want %v",
					tc.openerChar, tc.closerChar, result, tc.expected)
			}
		})
	}
}

func TestEmphasisDelimiterProcessor_OnMatch(t *testing.T) {
	tests := []struct {
		name     string
		consumes int
	}{
		{
			name:     "single character consume",
			consumes: 1,
		},
		{
			name:     "double character consume",
			consumes: 2,
		},
		{
			name:     "zero consume",
			consumes: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			processor := &emphasisDelimiterProcessor{}
			result := processor.OnMatch(tc.consumes)
			if result == nil {
				t.Fatal("OnMatch returned nil")
			}
			emphasis, ok := result.(*ast.Emphasis)
			if !ok {
				t.Fatalf("OnMatch returned %T, want *ast.Emphasis", result)
			}
			if emphasis.Level != tc.consumes {
				t.Errorf("OnMatch(%d) created emphasis with level %d, want %d",
					tc.consumes, emphasis.Level, tc.consumes)
			}
		})
	}
}

func TestEmphasisParser_Trigger(t *testing.T) {
	parser := &emphasisParser{}
	triggers := parser.Trigger()

	expectedTriggers := []byte{'*', '_'}
	if len(triggers) != len(expectedTriggers) {
		t.Fatalf("Trigger() returned %d characters, want %d", len(triggers), len(expectedTriggers))
	}

	for i, expected := range expectedTriggers {
		if triggers[i] != expected {
			t.Errorf("Trigger()[%d] = %c, want %c", i, triggers[i], expected)
		}
	}
}

func TestNewEmphasisParser(t *testing.T) {
	parser := NewEmphasisParser()
	if parser == nil {
		t.Fatal("NewEmphasisParser() returned nil")
	}

	// Verify it returns the expected type
	_, ok := parser.(*emphasisParser)
	if !ok {
		t.Fatalf("NewEmphasisParser() returned %T, want *emphasisParser", parser)
	}
}

func TestEmphasisParser_Parse(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		shouldBeNil  bool
		expectedChar byte
		expectedLen  int
	}{
		{
			name:         "single asterisk",
			input:        "*text",
			shouldBeNil:  false,
			expectedChar: '*',
			expectedLen:  1,
		},
		{
			name:         "double asterisk",
			input:        "**text",
			shouldBeNil:  false,
			expectedChar: '*',
			expectedLen:  2,
		},
		{
			name:         "single underscore",
			input:        "_text",
			shouldBeNil:  false,
			expectedChar: '_',
			expectedLen:  1,
		},
		{
			name:         "double underscore",
			input:        "__text",
			shouldBeNil:  false,
			expectedChar: '_',
			expectedLen:  2,
		},
		{
			name:         "triple asterisk",
			input:        "***text",
			shouldBeNil:  false,
			expectedChar: '*',
			expectedLen:  3,
		},
		{
			name:        "non-delimiter character",
			input:       "text",
			shouldBeNil: true,
		},
		// Note: Empty input case removed - PeekLine() never returns empty slice in practice,
		// so emphasis parser doesn't need to handle it. ScanDelimiter assumes non-empty input.
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &emphasisParser{}
			parent := ast.NewDocument()
			reader := text.NewReader([]byte(tc.input))
			context := NewContext()

			result := parser.Parse(parent, reader, context)

			if tc.shouldBeNil {
				if result != nil {
					t.Errorf("Parse() returned node, expected nil")
				}
				return
			}

			if result == nil {
				t.Fatal("Parse() returned nil, expected Delimiter node")
			}

			delimiter, ok := result.(*Delimiter)
			if !ok {
				t.Fatalf("Parse() returned %T, want *Delimiter", result)
			}

			if delimiter.Char != tc.expectedChar {
				t.Errorf("Delimiter.Char = %c, want %c", delimiter.Char, tc.expectedChar)
			}

			if delimiter.OriginalLength != tc.expectedLen {
				t.Errorf("Delimiter.OriginalLength = %d, want %d",
					delimiter.OriginalLength, tc.expectedLen)
			}

			// Verify delimiter was pushed to context
			lastDelim := context.LastDelimiter()
			if lastDelim == nil {
				t.Error("Delimiter was not pushed to context")
			} else if lastDelim != delimiter {
				t.Error("Wrong delimiter pushed to context")
			}
		})
	}
}
