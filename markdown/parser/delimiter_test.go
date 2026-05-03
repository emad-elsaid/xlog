package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

// mockProcessor implements DelimiterProcessor for testing.
type mockProcessor struct {
	delimiter byte
	canOpen   func(opener, closer *Delimiter) bool
	onMatch   func(consumes int) ast.Node
}

func (m *mockProcessor) IsDelimiter(c byte) bool {
	return c == m.delimiter
}

func (m *mockProcessor) CanOpenCloser(opener, closer *Delimiter) bool {
	if m.canOpen != nil {
		return m.canOpen(opener, closer)
	}
	return opener.Char == closer.Char
}

func (m *mockProcessor) OnMatch(consumes int) ast.Node {
	if m.onMatch != nil {
		return m.onMatch(consumes)
	}
	return ast.NewString([]byte("matched"))
}

func TestDelimiter_Inline(t *testing.T) {
	proc := &mockProcessor{delimiter: '*'}
	d := NewDelimiter(true, true, 2, '*', proc)

	// Should not panic
	d.Inline()
}

func TestDelimiter_Dump(t *testing.T) {
	proc := &mockProcessor{delimiter: '*'}
	d := NewDelimiter(true, true, 2, '*', proc)
	d.Segment = text.NewSegment(0, 2)

	source := []byte("**bold**")

	// Should not panic, just verify it works
	d.Dump(source, 0)
}

func TestDelimiter_Kind(t *testing.T) {
	proc := &mockProcessor{delimiter: '*'}
	d := NewDelimiter(true, true, 2, '*', proc)

	kind := d.Kind()
	if kind.String() != "Delimiter" {
		t.Errorf("Kind() = %v, want Delimiter", kind)
	}
}

func TestDelimiter_Text(t *testing.T) {
	tests := []struct {
		name     string
		segment  text.Segment
		source   []byte
		expected string
	}{
		{
			name:     "asterisk delimiter",
			segment:  text.NewSegment(0, 2),
			source:   []byte("**bold"),
			expected: "**",
		},
		{
			name:     "underscore delimiter",
			segment:  text.NewSegment(0, 1),
			source:   []byte("_italic_"),
			expected: "_",
		},
		{
			name:     "middle segment",
			segment:  text.NewSegment(3, 5),
			source:   []byte("abc**def"),
			expected: "**",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			proc := &mockProcessor{delimiter: '*'}
			d := NewDelimiter(true, true, 2, '*', proc)
			d.Segment = tc.segment

			got := string(d.Text(tc.source))
			if got != tc.expected {
				t.Errorf("Text() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestDelimiter_ConsumeCharacters(t *testing.T) {
	tests := []struct {
		name           string
		initialLength  int
		consume        int
		expectedLength int
		segmentStart   int
	}{
		{
			name:           "consume 1 from 2",
			initialLength:  2,
			consume:        1,
			expectedLength: 1,
			segmentStart:   0,
		},
		{
			name:           "consume 2 from 3",
			initialLength:  3,
			consume:        2,
			expectedLength: 1,
			segmentStart:   0,
		},
		{
			name:           "consume all",
			initialLength:  2,
			consume:        2,
			expectedLength: 0,
			segmentStart:   0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			proc := &mockProcessor{delimiter: '*'}
			d := NewDelimiter(true, true, tc.initialLength, '*', proc)
			d.Segment = text.NewSegment(tc.segmentStart, tc.segmentStart+tc.initialLength)

			d.ConsumeCharacters(tc.consume)

			if d.Length != tc.expectedLength {
				t.Errorf("Length = %d, want %d", d.Length, tc.expectedLength)
			}

			expectedStop := tc.segmentStart + tc.expectedLength
			if d.Segment.Stop != expectedStop {
				t.Errorf("Segment.Stop = %d, want %d", d.Segment.Stop, expectedStop)
			}
		})
	}
}

func TestDelimiter_CalcComsumption(t *testing.T) {
	tests := []struct {
		name           string
		openerCanClose bool
		openerLength   int
		openerOrigLen  int
		closerCanOpen  bool
		closerLength   int
		closerOrigLen  int
		expected       int
	}{
		{
			name:           "both length >= 2, consume 2",
			openerCanClose: false,
			openerLength:   2,
			openerOrigLen:  2,
			closerCanOpen:  false,
			closerLength:   2,
			closerOrigLen:  2,
			expected:       2,
		},
		{
			name:           "both length >= 2, different lengths",
			openerCanClose: false,
			openerLength:   3,
			openerOrigLen:  3,
			closerCanOpen:  false,
			closerLength:   2,
			closerOrigLen:  2,
			expected:       2,
		},
		{
			name:           "one length < 2, consume 1",
			openerCanClose: false,
			openerLength:   1,
			openerOrigLen:  1,
			closerCanOpen:  false,
			closerLength:   2,
			closerOrigLen:  2,
			expected:       1,
		},
		{
			name:           "rule of 3 violation - sum divisible by 3, closer not divisible by 3",
			openerCanClose: false,
			openerLength:   3,
			openerOrigLen:  3,
			closerCanOpen:  true,
			closerLength:   3,
			closerOrigLen:  3,
			expected:       2,
		},
		{
			name:           "rule of 3 - opener CanClose, sum=6, closer%3!=0",
			openerCanClose: true,
			openerLength:   2,
			openerOrigLen:  2,
			closerCanOpen:  false,
			closerLength:   2,
			closerOrigLen:  4,
			expected:       0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			proc := &mockProcessor{delimiter: '*'}

			opener := NewDelimiter(true, tc.openerCanClose, tc.openerLength, '*', proc)
			opener.OriginalLength = tc.openerOrigLen

			closer := NewDelimiter(tc.closerCanOpen, true, tc.closerLength, '*', proc)
			closer.OriginalLength = tc.closerOrigLen

			got := opener.CalcComsumption(closer)
			if got != tc.expected {
				t.Errorf("CalcComsumption() = %d, want %d", got, tc.expected)
			}
		})
	}
}

func TestNewDelimiter(t *testing.T) {
	tests := []struct {
		name     string
		canOpen  bool
		canClose bool
		length   int
		char     byte
	}{
		{
			name:     "opener only",
			canOpen:  true,
			canClose: false,
			length:   2,
			char:     '*',
		},
		{
			name:     "closer only",
			canOpen:  false,
			canClose: true,
			length:   1,
			char:     '_',
		},
		{
			name:     "both open and close",
			canOpen:  true,
			canClose: true,
			length:   3,
			char:     '*',
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			proc := &mockProcessor{delimiter: tc.char}
			d := NewDelimiter(tc.canOpen, tc.canClose, tc.length, tc.char, proc)

			if d.CanOpen != tc.canOpen {
				t.Errorf("CanOpen = %v, want %v", d.CanOpen, tc.canOpen)
			}
			if d.CanClose != tc.canClose {
				t.Errorf("CanClose = %v, want %v", d.CanClose, tc.canClose)
			}
			if d.Length != tc.length {
				t.Errorf("Length = %d, want %d", d.Length, tc.length)
			}
			if d.OriginalLength != tc.length {
				t.Errorf("OriginalLength = %d, want %d", d.OriginalLength, tc.length)
			}
			if d.Char != tc.char {
				t.Errorf("Char = %c, want %c", d.Char, tc.char)
			}
			if d.Processor != proc {
				t.Error("Processor not set correctly")
			}
			if d.PreviousDelimiter != nil {
				t.Error("PreviousDelimiter should be nil")
			}
			if d.NextDelimiter != nil {
				t.Error("NextDelimiter should be nil")
			}
		})
	}
}

func TestScanDelimiter(t *testing.T) {
	tests := []struct {
		name             string
		line             []byte
		before           rune
		minimum          int
		delimiter        byte
		expectedNil      bool
		expectedCanOpen  bool
		expectedCanClose bool
		expectedLength   int
	}{
		{
			name:             "asterisk at start",
			line:             []byte("**bold"),
			before:           ' ',
			minimum:          1,
			delimiter:        '*',
			expectedNil:      false,
			expectedCanOpen:  true,
			expectedCanClose: false,
			expectedLength:   2,
		},
		{
			name:             "underscore at start",
			line:             []byte("_italic"),
			before:           ' ',
			minimum:          1,
			delimiter:        '_',
			expectedNil:      false,
			expectedCanOpen:  true,
			expectedCanClose: false,
			expectedLength:   1,
		},
		{
			name:        "not a delimiter",
			line:        []byte("abc"),
			before:      ' ',
			minimum:     1,
			delimiter:   '*',
			expectedNil: true,
		},
		{
			name:        "below minimum length",
			line:        []byte("*a"),
			before:      ' ',
			minimum:     2,
			delimiter:   '*',
			expectedNil: true,
		},
		{
			name:             "asterisk after text",
			line:             []byte("**"),
			before:           'a',
			minimum:          1,
			delimiter:        '*',
			expectedNil:      false,
			expectedCanOpen:  false,
			expectedCanClose: true,
			expectedLength:   2,
		},
		{
			name:             "asterisk between spaces (space before, at end)",
			line:             []byte("**"),
			before:           ' ',
			minimum:          1,
			delimiter:        '*',
			expectedNil:      false,
			expectedCanOpen:  false,
			expectedCanClose: false,
			expectedLength:   2,
		},
		{
			name:             "underscore after punctuation",
			line:             []byte("_italic"),
			before:           '.',
			minimum:          1,
			delimiter:        '_',
			expectedNil:      false,
			expectedCanOpen:  true,
			expectedCanClose: false,
			expectedLength:   1,
		},
		{
			name:             "underscore after letter (cannot open)",
			line:             []byte("_"),
			before:           'a',
			minimum:          1,
			delimiter:        '_',
			expectedNil:      false,
			expectedCanOpen:  false,
			expectedCanClose: true,
			expectedLength:   1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			proc := &mockProcessor{delimiter: tc.delimiter}
			d := ScanDelimiter(tc.line, tc.before, tc.minimum, proc)

			if tc.expectedNil {
				if d != nil {
					t.Errorf("ScanDelimiter() = %v, want nil", d)
				}
				return
			}

			if d == nil {
				t.Fatal("ScanDelimiter() = nil, want non-nil")
			}

			if d.CanOpen != tc.expectedCanOpen {
				t.Errorf("CanOpen = %v, want %v", d.CanOpen, tc.expectedCanOpen)
			}
			if d.CanClose != tc.expectedCanClose {
				t.Errorf("CanClose = %v, want %v", d.CanClose, tc.expectedCanClose)
			}
			if d.Length != tc.expectedLength {
				t.Errorf("Length = %d, want %d", d.Length, tc.expectedLength)
			}
			if d.Char != tc.delimiter {
				t.Errorf("Char = %c, want %c", d.Char, tc.delimiter)
			}
		})
	}
}
