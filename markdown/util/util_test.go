package util

import (
	"bytes"
	"testing"
)

func TestCopyOnWriteBuffer_Write(t *testing.T) {
	tests := []struct {
		name           string
		initialBuffer  []byte
		writeData      [][]byte
		expectedResult []byte
		expectedCopied bool
	}{
		{
			name:           "write to empty buffer",
			initialBuffer:  []byte{},
			writeData:      [][]byte{[]byte("hello")},
			expectedResult: []byte("hello"),
			expectedCopied: true,
		},
		{
			name:           "write multiple times",
			initialBuffer:  []byte("initial"),
			writeData:      [][]byte{[]byte("hello"), []byte(" world")},
			expectedResult: []byte("hello world"),
			expectedCopied: true,
		},
		{
			name:           "write empty data",
			initialBuffer:  []byte("test"),
			writeData:      [][]byte{[]byte("")},
			expectedResult: []byte(""),
			expectedCopied: true,
		},
		{
			name:           "write nil",
			initialBuffer:  []byte("test"),
			writeData:      [][]byte{nil},
			expectedResult: []byte(""),
			expectedCopied: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := NewCopyOnWriteBuffer(tc.initialBuffer)

			for _, data := range tc.writeData {
				buf.Write(data)
			}

			if !bytes.Equal(buf.Bytes(), tc.expectedResult) {
				t.Errorf("expected %q, got %q", tc.expectedResult, buf.Bytes())
			}

			if buf.IsCopied() != tc.expectedCopied {
				t.Errorf("expected copied=%v, got %v", tc.expectedCopied, buf.IsCopied())
			}
		})
	}
}

func TestCopyOnWriteBuffer_WriteString(t *testing.T) {
	tests := []struct {
		name           string
		initialBuffer  []byte
		writeData      []string
		expectedResult []byte
		expectedCopied bool
	}{
		{
			name:           "write string to empty buffer",
			initialBuffer:  []byte{},
			writeData:      []string{"hello"},
			expectedResult: []byte("hello"),
			expectedCopied: true,
		},
		{
			name:           "write multiple strings",
			initialBuffer:  []byte("initial"),
			writeData:      []string{"hello", " world"},
			expectedResult: []byte("hello world"),
			expectedCopied: true,
		},
		{
			name:           "write empty string",
			initialBuffer:  []byte("test"),
			writeData:      []string{""},
			expectedResult: []byte(""),
			expectedCopied: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := NewCopyOnWriteBuffer(tc.initialBuffer)

			for _, data := range tc.writeData {
				buf.WriteString(data)
			}

			if !bytes.Equal(buf.Bytes(), tc.expectedResult) {
				t.Errorf("expected %q, got %q", tc.expectedResult, buf.Bytes())
			}

			if buf.IsCopied() != tc.expectedCopied {
				t.Errorf("expected copied=%v, got %v", tc.expectedCopied, buf.IsCopied())
			}
		})
	}
}

func TestCopyOnWriteBuffer_Append(t *testing.T) {
	tests := []struct {
		name           string
		initialBuffer  []byte
		appendData     [][]byte
		expectedResult []byte
		expectedCopied bool
	}{
		{
			name:           "append to empty buffer",
			initialBuffer:  []byte{},
			appendData:     [][]byte{[]byte("hello")},
			expectedResult: []byte("hello"),
			expectedCopied: true,
		},
		{
			name:           "append preserves initial content",
			initialBuffer:  []byte("initial"),
			appendData:     [][]byte{[]byte(" data")},
			expectedResult: []byte("initial data"),
			expectedCopied: true,
		},
		{
			name:           "multiple appends",
			initialBuffer:  []byte("start"),
			appendData:     [][]byte{[]byte(" middle"), []byte(" end")},
			expectedResult: []byte("start middle end"),
			expectedCopied: true,
		},
		{
			name:           "append empty data",
			initialBuffer:  []byte("test"),
			appendData:     [][]byte{[]byte("")},
			expectedResult: []byte("test"),
			expectedCopied: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := NewCopyOnWriteBuffer(tc.initialBuffer)

			for _, data := range tc.appendData {
				buf.Append(data)
			}

			if !bytes.Equal(buf.Bytes(), tc.expectedResult) {
				t.Errorf("expected %q, got %q", tc.expectedResult, buf.Bytes())
			}

			if buf.IsCopied() != tc.expectedCopied {
				t.Errorf("expected copied=%v, got %v", tc.expectedCopied, buf.IsCopied())
			}
		})
	}
}

func TestCopyOnWriteBuffer_AppendString(t *testing.T) {
	tests := []struct {
		name           string
		initialBuffer  []byte
		appendData     []string
		expectedResult []byte
		expectedCopied bool
	}{
		{
			name:           "append string preserves initial content",
			initialBuffer:  []byte("initial"),
			appendData:     []string{" data"},
			expectedResult: []byte("initial data"),
			expectedCopied: true,
		},
		{
			name:           "multiple string appends",
			initialBuffer:  []byte("start"),
			appendData:     []string{" middle", " end"},
			expectedResult: []byte("start middle end"),
			expectedCopied: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := NewCopyOnWriteBuffer(tc.initialBuffer)

			for _, data := range tc.appendData {
				buf.AppendString(data)
			}

			if !bytes.Equal(buf.Bytes(), tc.expectedResult) {
				t.Errorf("expected %q, got %q", tc.expectedResult, buf.Bytes())
			}

			if buf.IsCopied() != tc.expectedCopied {
				t.Errorf("expected copied=%v, got %v", tc.expectedCopied, buf.IsCopied())
			}
		})
	}
}

func TestCopyOnWriteBuffer_WriteByte(t *testing.T) {
	tests := []struct {
		name           string
		initialBuffer  []byte
		writeBytes     []byte
		expectedResult []byte
		expectedCopied bool
	}{
		{
			name:           "write single byte",
			initialBuffer:  []byte{},
			writeBytes:     []byte{'a'},
			expectedResult: []byte("a"),
			expectedCopied: true,
		},
		{
			name:           "write multiple bytes",
			initialBuffer:  []byte("start"),
			writeBytes:     []byte{'a', 'b', 'c'},
			expectedResult: []byte("abc"),
			expectedCopied: true,
		},
		{
			name:           "write null byte",
			initialBuffer:  []byte{},
			writeBytes:     []byte{0},
			expectedResult: []byte{0},
			expectedCopied: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := NewCopyOnWriteBuffer(tc.initialBuffer)

			for _, b := range tc.writeBytes {
				if err := buf.WriteByte(b); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if !bytes.Equal(buf.Bytes(), tc.expectedResult) {
				t.Errorf("expected %q, got %q", tc.expectedResult, buf.Bytes())
			}

			if buf.IsCopied() != tc.expectedCopied {
				t.Errorf("expected copied=%v, got %v", tc.expectedCopied, buf.IsCopied())
			}
		})
	}
}

func TestCopyOnWriteBuffer_AppendByte(t *testing.T) {
	tests := []struct {
		name           string
		initialBuffer  []byte
		appendBytes    []byte
		expectedResult []byte
		expectedCopied bool
	}{
		{
			name:           "append single byte",
			initialBuffer:  []byte("test"),
			appendBytes:    []byte{'!'},
			expectedResult: []byte("test!"),
			expectedCopied: true,
		},
		{
			name:           "append multiple bytes",
			initialBuffer:  []byte("hello"),
			appendBytes:    []byte{' ', 'w', 'o', 'r', 'l', 'd'},
			expectedResult: []byte("hello world"),
			expectedCopied: true,
		},
		{
			name:           "append to empty buffer",
			initialBuffer:  []byte{},
			appendBytes:    []byte{'x'},
			expectedResult: []byte("x"),
			expectedCopied: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buf := NewCopyOnWriteBuffer(tc.initialBuffer)

			for _, b := range tc.appendBytes {
				buf.AppendByte(b)
			}

			if !bytes.Equal(buf.Bytes(), tc.expectedResult) {
				t.Errorf("expected %q, got %q", tc.expectedResult, buf.Bytes())
			}

			if buf.IsCopied() != tc.expectedCopied {
				t.Errorf("expected copied=%v, got %v", tc.expectedCopied, buf.IsCopied())
			}
		})
	}
}

func TestCopyOnWriteBuffer_CopyOnWriteBehavior(t *testing.T) {
	original := []byte("original")
	buf := NewCopyOnWriteBuffer(original)

	// Before modification, should not be copied
	if buf.IsCopied() {
		t.Error("buffer should not be copied initially")
	}

	// After Write, should be copied
	buf.Write([]byte("new"))

	if !buf.IsCopied() {
		t.Error("buffer should be copied after Write")
	}

	// Original buffer should remain unchanged
	if !bytes.Equal(original, []byte("original")) {
		t.Errorf("original buffer was modified: %q", original)
	}
}

func TestCopyOnWriteBuffer_AppendVsWrite(t *testing.T) {
	// Test that Append preserves initial content while Write clears it
	t.Run("Append preserves content", func(t *testing.T) {
		buf := NewCopyOnWriteBuffer([]byte("initial"))
		buf.Append([]byte(" appended"))

		expected := []byte("initial appended")
		if !bytes.Equal(buf.Bytes(), expected) {
			t.Errorf("Append: expected %q, got %q", expected, buf.Bytes())
		}
	})

	t.Run("Write clears content", func(t *testing.T) {
		buf := NewCopyOnWriteBuffer([]byte("initial"))
		buf.Write([]byte("written"))

		expected := []byte("written")
		if !bytes.Equal(buf.Bytes(), expected) {
			t.Errorf("Write: expected %q, got %q", expected, buf.Bytes())
		}
	})
}

func TestIsEscapedPunctuation(t *testing.T) {
	tests := []struct {
		name     string
		source   []byte
		index    int
		expected bool
	}{
		{
			name:     "escaped period",
			source:   []byte(`\.`),
			index:    0,
			expected: true,
		},
		{
			name:     "escaped asterisk",
			source:   []byte(`\*`),
			index:    0,
			expected: true,
		},
		{
			name:     "not escaped",
			source:   []byte(`a*`),
			index:    0,
			expected: false,
		},
		{
			name:     "backslash but not punctuation",
			source:   []byte(`\a`),
			index:    0,
			expected: false,
		},
		{
			name:     "middle of string",
			source:   []byte(`hello\*world`),
			index:    5,
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsEscapedPunctuation(tc.source, tc.index)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestReadWhile(t *testing.T) {
	tests := []struct {
		name        string
		source      []byte
		index       [2]int
		pred        func(byte) bool
		expectedPos int
		expectedOk  bool
	}{
		{
			name:        "read spaces",
			source:      []byte("   abc"),
			index:       [2]int{0, 6},
			pred:        func(b byte) bool { return b == ' ' },
			expectedPos: 3,
			expectedOk:  true,
		},
		{
			name:        "read digits",
			source:      []byte("12345abc"),
			index:       [2]int{0, 8},
			pred:        func(b byte) bool { return b >= '0' && b <= '9' },
			expectedPos: 5,
			expectedOk:  true,
		},
		{
			name:        "no match",
			source:      []byte("abc123"),
			index:       [2]int{0, 6},
			pred:        func(b byte) bool { return b >= '0' && b <= '9' },
			expectedPos: 0,
			expectedOk:  false,
		},
		{
			name:        "partial range",
			source:      []byte("  abc  "),
			index:       [2]int{2, 5},
			pred:        func(b byte) bool { return b == ' ' },
			expectedPos: 2,
			expectedOk:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pos, ok := ReadWhile(tc.source, tc.index, tc.pred)

			if pos != tc.expectedPos {
				t.Errorf("expected pos %d, got %d", tc.expectedPos, pos)
			}

			if ok != tc.expectedOk {
				t.Errorf("expected ok %v, got %v", tc.expectedOk, ok)
			}
		})
	}
}

func TestIsBlank(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "all spaces",
			input:    []byte("   "),
			expected: true,
		},
		{
			name:     "spaces and tabs",
			input:    []byte(" \t \t "),
			expected: true,
		},
		{
			name:     "empty",
			input:    []byte(""),
			expected: true,
		},
		{
			name:     "contains text",
			input:    []byte("  a  "),
			expected: false,
		},
		{
			name:     "newline is space",
			input:    []byte(" \n "),
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsBlank(tc.input)
			if result != tc.expected {
				t.Errorf("IsBlank(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestVisualizeSpaces(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "spaces",
			input:    []byte("hello world"),
			expected: []byte("hello[SPACE]world"),
		},
		{
			name:     "tabs",
			input:    []byte("hello\tworld"),
			expected: []byte("hello[TAB]world"),
		},
		{
			name:     "newlines",
			input:    []byte("hello\nworld"),
			expected: []byte("hello[NEWLINE]\nworld"),
		},
		{
			name:     "mixed whitespace",
			input:    []byte(" \t\n\r"),
			expected: []byte("[SPACE][TAB][NEWLINE]\n[CR]"),
		},
		{
			name:     "null byte",
			input:    []byte("a\x00b"),
			expected: []byte("a[NUL]b"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := VisualizeSpaces(tc.input)
			if !bytes.Equal(result, tc.expected) {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestTabWidth(t *testing.T) {
	tests := []struct {
		name          string
		currentPos    int
		expectedWidth int
	}{
		{
			name:          "position 0",
			currentPos:    0,
			expectedWidth: 4,
		},
		{
			name:          "position 1",
			currentPos:    1,
			expectedWidth: 3,
		},
		{
			name:          "position 2",
			currentPos:    2,
			expectedWidth: 2,
		},
		{
			name:          "position 3",
			currentPos:    3,
			expectedWidth: 1,
		},
		{
			name:          "position 4 (new tab stop)",
			currentPos:    4,
			expectedWidth: 4,
		},
		{
			name:          "position 5",
			currentPos:    5,
			expectedWidth: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := TabWidth(tc.currentPos)
			if result != tc.expectedWidth {
				t.Errorf("TabWidth(%d) = %d, expected %d", tc.currentPos, result, tc.expectedWidth)
			}
		})
	}
}

func TestIndentPosition(t *testing.T) {
	tests := []struct {
		name            string
		input           []byte
		currentPos      int
		width           int
		expectedPos     int
		expectedPadding int
	}{
		{
			name:            "spaces only",
			input:           []byte("    text"),
			currentPos:      0,
			width:           2,
			expectedPos:     2,
			expectedPadding: 0,
		},
		{
			name:            "tab at start",
			input:           []byte("\ttext"),
			currentPos:      0,
			width:           2,
			expectedPos:     1,
			expectedPadding: 2,
		},
		{
			name:            "zero width",
			input:           []byte("  text"),
			currentPos:      0,
			width:           0,
			expectedPos:     0,
			expectedPadding: 0,
		},
		{
			name:            "width exceeds spaces",
			input:           []byte("  text"),
			currentPos:      0,
			width:           10,
			expectedPos:     -1,
			expectedPadding: -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pos, padding := IndentPosition(tc.input, tc.currentPos, tc.width)

			if pos != tc.expectedPos {
				t.Errorf("expected pos %d, got %d", tc.expectedPos, pos)
			}

			if padding != tc.expectedPadding {
				t.Errorf("expected padding %d, got %d", tc.expectedPadding, padding)
			}
		})
	}
}

func TestTrimLeftSpace(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "leading spaces",
			input:    []byte("   hello"),
			expected: []byte("hello"),
		},
		{
			name:     "leading tabs",
			input:    []byte("\t\thello"),
			expected: []byte("hello"),
		},
		{
			name:     "mixed leading whitespace",
			input:    []byte(" \t \nhello"),
			expected: []byte("hello"),
		},
		{
			name:     "no leading whitespace",
			input:    []byte("hello world"),
			expected: []byte("hello world"),
		},
		{
			name:     "only whitespace",
			input:    []byte("   "),
			expected: []byte(""),
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: []byte(""),
		},
		{
			name:     "preserves trailing spaces",
			input:    []byte("  hello  "),
			expected: []byte("hello  "),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := TrimLeftSpace(tc.input)
			if !bytes.Equal(result, tc.expected) {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestTrimRightSpace(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "trailing spaces",
			input:    []byte("hello   "),
			expected: []byte("hello"),
		},
		{
			name:     "trailing tabs",
			input:    []byte("hello\t\t"),
			expected: []byte("hello"),
		},
		{
			name:     "mixed trailing whitespace",
			input:    []byte("hello \t \n"),
			expected: []byte("hello"),
		},
		{
			name:     "no trailing whitespace",
			input:    []byte("hello world"),
			expected: []byte("hello world"),
		},
		{
			name:     "only whitespace",
			input:    []byte("   "),
			expected: []byte(""),
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: []byte(""),
		},
		{
			name:     "preserves leading spaces",
			input:    []byte("  hello  "),
			expected: []byte("  hello"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := TrimRightSpace(tc.input)
			if !bytes.Equal(result, tc.expected) {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "no special characters",
			input:    []byte("hello world"),
			expected: []byte("hello world"),
		},
		{
			name:     "ampersand",
			input:    []byte("Tom & Jerry"),
			expected: []byte("Tom &amp; Jerry"),
		},
		{
			name:     "less than",
			input:    []byte("5 < 10"),
			expected: []byte("5 &lt; 10"),
		},
		{
			name:     "greater than",
			input:    []byte("10 > 5"),
			expected: []byte("10 &gt; 5"),
		},
		{
			name:     "quote",
			input:    []byte(`He said "hello"`),
			expected: []byte("He said &quot;hello&quot;"),
		},
		{
			name:     "all special characters",
			input:    []byte(`<div>"text" & 'stuff'</div>`),
			expected: []byte("&lt;div&gt;&quot;text&quot; &amp; 'stuff'&lt;/div&gt;"),
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: []byte(""),
		},
		{
			name:     "multiple ampersands",
			input:    []byte("&&&"),
			expected: []byte("&amp;&amp;&amp;"),
		},
		{
			name:     "html tag",
			input:    []byte("<script>alert('xss')</script>"),
			expected: []byte("&lt;script&gt;alert('xss')&lt;/script&gt;"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := EscapeHTML(tc.input)
			if !bytes.Equal(result, tc.expected) {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestToLinkReference(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "simple text",
			input:    []byte("hello"),
			expected: "hello",
		},
		{
			name:     "with leading spaces",
			input:    []byte("  hello"),
			expected: "hello",
		},
		{
			name:     "with trailing spaces",
			input:    []byte("hello  "),
			expected: "hello",
		},
		{
			name:     "with multiple spaces",
			input:    []byte("hello   world"),
			expected: "hello world",
		},
		{
			name:     "uppercase to lowercase",
			input:    []byte("HELLO WORLD"),
			expected: "hello world",
		},
		{
			name:     "mixed case",
			input:    []byte("Hello World"),
			expected: "hello world",
		},
		{
			name:     "with tabs",
			input:    []byte("hello\t\tworld"),
			expected: "hello world",
		},
		{
			name:     "complex whitespace",
			input:    []byte("  Hello   World  "),
			expected: "hello world",
		},
		{
			name:     "unicode characters",
			input:    []byte("Héllo Wörld"),
			expected: "héllo wörld",
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: "",
		},
		{
			name:     "only spaces",
			input:    []byte("   "),
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ToLinkReference(tc.input)
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{
			name:     "digit 0",
			input:    '0',
			expected: true,
		},
		{
			name:     "digit 5",
			input:    '5',
			expected: true,
		},
		{
			name:     "digit 9",
			input:    '9',
			expected: true,
		},
		{
			name:     "lowercase letter",
			input:    'a',
			expected: false,
		},
		{
			name:     "uppercase letter",
			input:    'Z',
			expected: false,
		},
		{
			name:     "space",
			input:    ' ',
			expected: false,
		},
		{
			name:     "special character",
			input:    '!',
			expected: false,
		},
		{
			name:     "before 0",
			input:    '/',
			expected: false,
		},
		{
			name:     "after 9",
			input:    ':',
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsNumeric(tc.input)
			if result != tc.expected {
				t.Errorf("IsNumeric(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsHexDecimal(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected bool
	}{
		{
			name:     "digit 0",
			input:    '0',
			expected: true,
		},
		{
			name:     "digit 9",
			input:    '9',
			expected: true,
		},
		{
			name:     "lowercase a",
			input:    'a',
			expected: true,
		},
		{
			name:     "lowercase f",
			input:    'f',
			expected: true,
		},
		{
			name:     "uppercase A",
			input:    'A',
			expected: true,
		},
		{
			name:     "uppercase F",
			input:    'F',
			expected: true,
		},
		{
			name:     "lowercase g (not hex)",
			input:    'g',
			expected: false,
		},
		{
			name:     "uppercase G (not hex)",
			input:    'G',
			expected: false,
		},
		{
			name:     "lowercase z",
			input:    'z',
			expected: false,
		},
		{
			name:     "space",
			input:    ' ',
			expected: false,
		},
		{
			name:     "special character",
			input:    '@',
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsHexDecimal(tc.input)
			if result != tc.expected {
				t.Errorf("IsHexDecimal(%q) = %v, expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestDedentPosition(t *testing.T) {
	tests := []struct {
		name            string
		input           []byte
		currentPos      int
		width           int
		expectedPos     int
		expectedPadding int
	}{
		{
			name:            "zero width",
			input:           []byte("    text"),
			currentPos:      0,
			width:           0,
			expectedPos:     0,
			expectedPadding: 0,
		},
		{
			name:            "dedent spaces",
			input:           []byte("    text"),
			currentPos:      0,
			width:           2,
			expectedPos:     4,
			expectedPadding: 2,
		},
		{
			name:            "dedent exact match",
			input:           []byte("    text"),
			currentPos:      0,
			width:           4,
			expectedPos:     4,
			expectedPadding: 0,
		},
		{
			name:            "dedent exceeds spaces",
			input:           []byte("  text"),
			currentPos:      0,
			width:           5,
			expectedPos:     2,
			expectedPadding: 0,
		},
		{
			name:            "dedent with tab",
			input:           []byte("\ttext"),
			currentPos:      0,
			width:           2,
			expectedPos:     1,
			expectedPadding: 2,
		},
		{
			name:            "dedent mixed spaces and tabs",
			input:           []byte("  \t  text"),
			currentPos:      0,
			width:           4,
			expectedPos:     5,
			expectedPadding: 2,
		},
		{
			name:            "empty input",
			input:           []byte(""),
			currentPos:      0,
			width:           2,
			expectedPos:     0,
			expectedPadding: 0,
		},
		{
			name:            "no whitespace",
			input:           []byte("text"),
			currentPos:      0,
			width:           2,
			expectedPos:     0,
			expectedPadding: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pos, padding := DedentPosition(tc.input, tc.currentPos, tc.width)
			if pos != tc.expectedPos {
				t.Errorf("expected pos %d, got %d", tc.expectedPos, pos)
			}
			if padding != tc.expectedPadding {
				t.Errorf("expected padding %d, got %d", tc.expectedPadding, padding)
			}
		})
	}
}

func TestIndentWidth(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		currentPos    int
		expectedWidth int
		expectedPos   int
	}{
		{
			name:          "no indentation",
			input:         []byte("text"),
			currentPos:    0,
			expectedWidth: 0,
			expectedPos:   0,
		},
		{
			name:          "spaces only",
			input:         []byte("    text"),
			currentPos:    0,
			expectedWidth: 4,
			expectedPos:   4,
		},
		{
			name:          "single tab",
			input:         []byte("\ttext"),
			currentPos:    0,
			expectedWidth: 4,
			expectedPos:   1,
		},
		{
			name:          "tab at position 1",
			input:         []byte("\ttext"),
			currentPos:    1,
			expectedWidth: 3,
			expectedPos:   1,
		},
		{
			name:          "tab at position 2",
			input:         []byte("\ttext"),
			currentPos:    2,
			expectedWidth: 2,
			expectedPos:   1,
		},
		{
			name:          "mixed spaces and tabs",
			input:         []byte("  \t  text"),
			currentPos:    0,
			expectedWidth: 6,
			expectedPos:   5,
		},
		{
			name:          "empty input",
			input:         []byte(""),
			currentPos:    0,
			expectedWidth: 0,
			expectedPos:   0,
		},
		{
			name:          "all whitespace",
			input:         []byte("    \t  "),
			currentPos:    0,
			expectedWidth: 10,
			expectedPos:   7,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			width, pos := IndentWidth(tc.input, tc.currentPos)
			if width != tc.expectedWidth {
				t.Errorf("expected width %d, got %d", tc.expectedWidth, width)
			}
			if pos != tc.expectedPos {
				t.Errorf("expected pos %d, got %d", tc.expectedPos, pos)
			}
		})
	}
}

func TestFirstNonSpacePosition(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected int
	}{
		{
			name:     "no leading whitespace",
			input:    []byte("text"),
			expected: 0,
		},
		{
			name:     "spaces before text",
			input:    []byte("    text"),
			expected: 4,
		},
		{
			name:     "tabs before text",
			input:    []byte("\t\ttext"),
			expected: 2,
		},
		{
			name:     "mixed whitespace",
			input:    []byte("  \t  text"),
			expected: 5,
		},
		{
			name:     "newline only",
			input:    []byte("\n"),
			expected: -1,
		},
		{
			name:     "spaces then newline",
			input:    []byte("  \n"),
			expected: -1,
		},
		{
			name:     "empty input",
			input:    []byte(""),
			expected: -1,
		},
		{
			name:     "only spaces",
			input:    []byte("    "),
			expected: -1,
		},
		{
			name:     "only tabs",
			input:    []byte("\t\t"),
			expected: -1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FirstNonSpacePosition(tc.input)
			if result != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, result)
			}
		})
	}
}

func TestFindClosure(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		opener       byte
		closure      byte
		codeSpan     bool
		allowNesting bool
		expected     int
	}{
		{
			name:         "simple balanced brackets",
			input:        []byte("text]"),
			opener:       '[',
			closure:      ']',
			codeSpan:     false,
			allowNesting: false,
			expected:     4,
		},
		{
			name:         "simple balanced parentheses",
			input:        []byte("hello)"),
			opener:       '(',
			closure:      ')',
			codeSpan:     false,
			allowNesting: false,
			expected:     5,
		},
		{
			name:         "no closure found",
			input:        []byte("text without closure"),
			opener:       '[',
			closure:      ']',
			codeSpan:     false,
			allowNesting: false,
			expected:     -1,
		},
		{
			name:         "empty input",
			input:        []byte(""),
			opener:       '[',
			closure:      ']',
			codeSpan:     false,
			allowNesting: false,
			expected:     -1,
		},
		{
			name:         "nested brackets with nesting allowed",
			input:        []byte("outer [inner] outer]"),
			opener:       '[',
			closure:      ']',
			codeSpan:     false,
			allowNesting: true,
			expected:     19,
		},
		{
			name:         "nested brackets without nesting allowed",
			input:        []byte("outer [inner] outer]"),
			opener:       '[',
			closure:      ']',
			codeSpan:     false,
			allowNesting: false,
			expected:     -1,
		},
		{
			name:         "escaped closure character",
			input:        []byte(`text \] more]`),
			opener:       '[',
			closure:      ']',
			codeSpan:     false,
			allowNesting: false,
			expected:     12,
		},
		{
			name:         "escaped opener character",
			input:        []byte(`text \[ more]`),
			opener:       '[',
			closure:      ']',
			codeSpan:     false,
			allowNesting: false,
			expected:     12,
		},
		{
			name:         "code span with backticks",
			input:        []byte("text `code]` more]"),
			opener:       '[',
			closure:      ']',
			codeSpan:     true,
			allowNesting: false,
			expected:     17,
		},
		{
			name:         "closure inside code span ignored",
			input:        []byte("`]`]"),
			opener:       '[',
			closure:      ']',
			codeSpan:     true,
			allowNesting: false,
			expected:     3,
		},
		{
			name:         "multiple backtick code span",
			input:        []byte("``code]``]"),
			opener:       '[',
			closure:      ']',
			codeSpan:     true,
			allowNesting: false,
			expected:     9,
		},
		{
			name:         "mismatched backtick counts leaves code span open",
			input:        []byte("``code`]"),
			opener:       '[',
			closure:      ']',
			codeSpan:     true,
			allowNesting: false,
			expected:     -1,
		},
		{
			name:         "closure at start",
			input:        []byte("]"),
			opener:       '[',
			closure:      ']',
			codeSpan:     false,
			allowNesting: false,
			expected:     0,
		},
		{
			name:         "multiple nested levels",
			input:        []byte("a[b[c]d]e]"),
			opener:       '[',
			closure:      ']',
			codeSpan:     false,
			allowNesting: true,
			expected:     9,
		},
		{
			name:         "escaped punctuation before closure",
			input:        []byte(`\*\.]`),
			opener:       '[',
			closure:      ']',
			codeSpan:     false,
			allowNesting: false,
			expected:     4,
		},
		{
			name:         "code span disabled finds first closure",
			input:        []byte("`backticks]ignored]"),
			opener:       '[',
			closure:      ']',
			codeSpan:     false,
			allowNesting: false,
			expected:     10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FindClosure(tc.input, tc.opener, tc.closure, tc.codeSpan, tc.allowNesting)
			if result != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, result)
			}
		})
	}
}
