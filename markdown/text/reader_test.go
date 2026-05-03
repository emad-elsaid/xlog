package text

import (
	"regexp"
	"testing"
)

func TestFindSubMatchReader(t *testing.T) {
	s := "微笑"
	r := NewReader([]byte(":" + s + ":"))
	reg := regexp.MustCompile(`:(\p{L}+):`)
	match := r.FindSubMatch(reg)
	if len(match) != 2 || string(match[1]) != s {
		t.Fatal("no match cjk")
	}
}

func TestSegment_Value(t *testing.T) {
	tests := []struct {
		name     string
		segment  Segment
		buffer   []byte
		expected []byte
	}{
		{
			name:     "simple segment without padding",
			segment:  Segment{Start: 0, Stop: 5, Padding: 0},
			buffer:   []byte("hello world"),
			expected: []byte("hello"),
		},
		{
			name:     "segment with padding",
			segment:  Segment{Start: 0, Stop: 5, Padding: 3},
			buffer:   []byte("hello"),
			expected: []byte("   hello"),
		},
		{
			name:     "segment with ForceNewline set",
			segment:  Segment{Start: 0, Stop: 5, Padding: 0, ForceNewline: true},
			buffer:   []byte("hello"),
			expected: []byte("hello\n"),
		},
		{
			name:     "segment with ForceNewline already has newline",
			segment:  Segment{Start: 0, Stop: 6, Padding: 0, ForceNewline: true},
			buffer:   []byte("hello\n"),
			expected: []byte("hello\n"),
		},
		{
			name:     "middle segment",
			segment:  Segment{Start: 6, Stop: 11, Padding: 0},
			buffer:   []byte("hello world"),
			expected: []byte("world"),
		},
		{
			name:     "empty segment",
			segment:  Segment{Start: 5, Stop: 5, Padding: 0},
			buffer:   []byte("hello"),
			expected: []byte(""),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.segment.Value(tc.buffer)
			if string(result) != string(tc.expected) {
				t.Errorf("got %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestSegment_Len(t *testing.T) {
	tests := []struct {
		name     string
		segment  Segment
		expected int
	}{
		{
			name:     "simple segment",
			segment:  Segment{Start: 0, Stop: 5, Padding: 0},
			expected: 5,
		},
		{
			name:     "segment with padding",
			segment:  Segment{Start: 0, Stop: 5, Padding: 3},
			expected: 8,
		},
		{
			name:     "empty segment",
			segment:  Segment{Start: 5, Stop: 5, Padding: 0},
			expected: 0,
		},
		{
			name:     "only padding",
			segment:  Segment{Start: 0, Stop: 0, Padding: 4},
			expected: 4,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.segment.Len()
			if result != tc.expected {
				t.Errorf("got %d, want %d", result, tc.expected)
			}
		})
	}
}

func TestSegment_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		segment  Segment
		expected bool
	}{
		{
			name:     "empty segment",
			segment:  Segment{Start: 5, Stop: 5, Padding: 0},
			expected: true,
		},
		{
			name:     "start greater than stop",
			segment:  Segment{Start: 5, Stop: 3, Padding: 0},
			expected: true,
		},
		{
			name:     "non-empty segment",
			segment:  Segment{Start: 0, Stop: 5, Padding: 0},
			expected: false,
		},
		{
			name:     "segment with only padding",
			segment:  Segment{Start: 5, Stop: 5, Padding: 3},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.segment.IsEmpty()
			if result != tc.expected {
				t.Errorf("got %v, want %v", result, tc.expected)
			}
		})
	}
}

func TestSegment_WithStart(t *testing.T) {
	tests := []struct {
		name     string
		segment  Segment
		newStart int
		expected Segment
	}{
		{
			name:     "update start position",
			segment:  Segment{Start: 0, Stop: 5, Padding: 2},
			newStart: 3,
			expected: Segment{Start: 3, Stop: 5, Padding: 2},
		},
		{
			name:     "update to zero",
			segment:  Segment{Start: 10, Stop: 15, Padding: 0},
			newStart: 0,
			expected: Segment{Start: 0, Stop: 15, Padding: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.segment.WithStart(tc.newStart)
			if result != tc.expected {
				t.Errorf("got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func TestSegment_WithStop(t *testing.T) {
	tests := []struct {
		name     string
		segment  Segment
		newStop  int
		expected Segment
	}{
		{
			name:     "update stop position",
			segment:  Segment{Start: 0, Stop: 5, Padding: 2},
			newStop:  10,
			expected: Segment{Start: 0, Stop: 10, Padding: 2},
		},
		{
			name:     "shrink segment",
			segment:  Segment{Start: 0, Stop: 15, Padding: 0},
			newStop:  5,
			expected: Segment{Start: 0, Stop: 5, Padding: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.segment.WithStop(tc.newStop)
			if result != tc.expected {
				t.Errorf("got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func TestSegment_ConcatPadding(t *testing.T) {
	tests := []struct {
		name     string
		segment  Segment
		input    []byte
		expected []byte
	}{
		{
			name:     "no padding",
			segment:  Segment{Start: 0, Stop: 5, Padding: 0},
			input:    []byte("hello"),
			expected: []byte("hello"),
		},
		{
			name:     "with padding",
			segment:  Segment{Start: 0, Stop: 5, Padding: 3},
			input:    []byte("hello"),
			expected: []byte("hello   "),
		},
		{
			name:     "empty input with padding",
			segment:  Segment{Start: 0, Stop: 0, Padding: 2},
			input:    []byte(""),
			expected: []byte("  "),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.segment.ConcatPadding(tc.input)
			if string(result) != string(tc.expected) {
				t.Errorf("got %q, want %q", result, tc.expected)
			}
		})
	}
}

func TestSegment_TrimRightSpace(t *testing.T) {
	tests := []struct {
		name     string
		segment  Segment
		buffer   []byte
		expected Segment
	}{
		{
			name:     "trim trailing spaces",
			segment:  Segment{Start: 0, Stop: 8},
			buffer:   []byte("hello   "),
			expected: Segment{Start: 0, Stop: 5, Padding: 0},
		},
		{
			name:     "no trailing spaces",
			segment:  Segment{Start: 0, Stop: 5},
			buffer:   []byte("hello"),
			expected: Segment{Start: 0, Stop: 5, Padding: 0},
		},
		{
			name:     "all spaces",
			segment:  Segment{Start: 0, Stop: 3},
			buffer:   []byte("   "),
			expected: Segment{Start: 0, Stop: 0, Padding: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.segment.TrimRightSpace(tc.buffer)
			if result != tc.expected {
				t.Errorf("got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func TestSegment_TrimLeftSpace(t *testing.T) {
	tests := []struct {
		name     string
		segment  Segment
		buffer   []byte
		expected Segment
	}{
		{
			name:     "trim leading spaces",
			segment:  Segment{Start: 0, Stop: 8},
			buffer:   []byte("   hello"),
			expected: Segment{Start: 3, Stop: 8},
		},
		{
			name:     "no leading spaces",
			segment:  Segment{Start: 0, Stop: 5},
			buffer:   []byte("hello"),
			expected: Segment{Start: 0, Stop: 5},
		},
		{
			name:     "all spaces",
			segment:  Segment{Start: 0, Stop: 3},
			buffer:   []byte("   "),
			expected: Segment{Start: 3, Stop: 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.segment.TrimLeftSpace(tc.buffer)
			if result != tc.expected {
				t.Errorf("got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func TestSegment_TrimLeftSpaceWidth(t *testing.T) {
	tests := []struct {
		name     string
		segment  Segment
		width    int
		buffer   []byte
		expected Segment
	}{
		{
			name:     "trim 3 spaces",
			segment:  Segment{Start: 0, Stop: 8, Padding: 0},
			width:    3,
			buffer:   []byte("   hello"),
			expected: Segment{Start: 3, Stop: 8, Padding: 0},
		},
		{
			name:     "trim from padding only",
			segment:  Segment{Start: 0, Stop: 5, Padding: 3},
			width:    2,
			buffer:   []byte("hello"),
			expected: Segment{Start: 0, Stop: 5, Padding: 1},
		},
		{
			name:     "width exceeds available",
			segment:  Segment{Start: 0, Stop: 3, Padding: 2},
			width:    5,
			buffer:   []byte("  a"),
			expected: Segment{Start: 2, Stop: 3, Padding: 0},
		},
		{
			name:     "tab character handling",
			segment:  Segment{Start: 0, Stop: 4, Padding: 0},
			width:    4,
			buffer:   []byte("\thel"),
			expected: Segment{Start: 1, Stop: 4, Padding: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.segment.TrimLeftSpaceWidth(tc.width, tc.buffer)
			if result != tc.expected {
				t.Errorf("got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func TestNewSegment(t *testing.T) {
	tests := []struct {
		name     string
		start    int
		stop     int
		expected Segment
	}{
		{
			name:     "basic segment",
			start:    0,
			stop:     5,
			expected: Segment{Start: 0, Stop: 5, Padding: 0},
		},
		{
			name:     "empty segment",
			start:    5,
			stop:     5,
			expected: Segment{Start: 5, Stop: 5, Padding: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NewSegment(tc.start, tc.stop)
			if result != tc.expected {
				t.Errorf("got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func TestNewSegmentPadding(t *testing.T) {
	tests := []struct {
		name     string
		start    int
		stop     int
		padding  int
		expected Segment
	}{
		{
			name:     "segment with padding",
			start:    0,
			stop:     5,
			padding:  3,
			expected: Segment{Start: 0, Stop: 5, Padding: 3},
		},
		{
			name:     "segment without padding",
			start:    0,
			stop:     5,
			padding:  0,
			expected: Segment{Start: 0, Stop: 5, Padding: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := NewSegmentPadding(tc.start, tc.stop, tc.padding)
			if result != tc.expected {
				t.Errorf("got %+v, want %+v", result, tc.expected)
			}
		})
	}
}

func TestSegments_Operations(t *testing.T) {
	t.Run("Append", func(t *testing.T) {
		s := NewSegments()
		s.Append(Segment{Start: 0, Stop: 5})
		s.Append(Segment{Start: 5, Stop: 10})

		if s.Len() != 2 {
			t.Errorf("got length %d, want 2", s.Len())
		}
		if s.At(0).Start != 0 || s.At(0).Stop != 5 {
			t.Errorf("got first segment %+v, want {Start:0 Stop:5}", s.At(0))
		}
	})

	t.Run("AppendAll", func(t *testing.T) {
		s := NewSegments()
		segments := []Segment{
			{Start: 0, Stop: 5},
			{Start: 5, Stop: 10},
			{Start: 10, Stop: 15},
		}
		s.AppendAll(segments)

		if s.Len() != 3 {
			t.Errorf("got length %d, want 3", s.Len())
		}
	})

	t.Run("Set", func(t *testing.T) {
		s := NewSegments()
		s.Append(Segment{Start: 0, Stop: 5})
		s.Set(0, Segment{Start: 0, Stop: 10})

		if s.At(0).Stop != 10 {
			t.Errorf("got stop %d, want 10", s.At(0).Stop)
		}
	})

	t.Run("SetSliced", func(t *testing.T) {
		s := NewSegments()
		s.Append(Segment{Start: 0, Stop: 5})
		s.Append(Segment{Start: 5, Stop: 10})
		s.Append(Segment{Start: 10, Stop: 15})
		s.SetSliced(1, 3)

		if s.Len() != 2 {
			t.Errorf("got length %d, want 2", s.Len())
		}
		if s.At(0).Start != 5 {
			t.Errorf("got first segment start %d, want 5", s.At(0).Start)
		}
	})

	t.Run("Sliced", func(t *testing.T) {
		s := NewSegments()
		s.Append(Segment{Start: 0, Stop: 5})
		s.Append(Segment{Start: 5, Stop: 10})
		s.Append(Segment{Start: 10, Stop: 15})

		sliced := s.Sliced(1, 3)
		if len(sliced) != 2 {
			t.Errorf("got length %d, want 2", len(sliced))
		}
		if sliced[0].Start != 5 {
			t.Errorf("got first segment start %d, want 5", sliced[0].Start)
		}
	})

	t.Run("Clear", func(t *testing.T) {
		s := NewSegments()
		s.Append(Segment{Start: 0, Stop: 5})
		s.Clear()

		if s.Len() != 0 {
			t.Errorf("got length %d, want 0", s.Len())
		}
	})

	t.Run("Unshift", func(t *testing.T) {
		s := NewSegments()
		s.Append(Segment{Start: 5, Stop: 10})
		s.Unshift(Segment{Start: 0, Stop: 5})

		if s.Len() != 2 {
			t.Errorf("got length %d, want 2", s.Len())
		}
		if s.At(0).Start != 0 {
			t.Errorf("got first segment start %d, want 0", s.At(0).Start)
		}
	})

	t.Run("Value", func(t *testing.T) {
		s := NewSegments()
		s.Append(Segment{Start: 0, Stop: 5})
		s.Append(Segment{Start: 6, Stop: 11})

		buffer := []byte("hello world")
		result := s.Value(buffer)

		if string(result) != "helloworld" {
			t.Errorf("got %q, want %q", result, "helloworld")
		}
	})

	t.Run("Len on nil", func(t *testing.T) {
		s := NewSegments()
		if s.Len() != 0 {
			t.Errorf("got length %d, want 0", s.Len())
		}
	})
}

func TestReader_BasicOperations(t *testing.T) {
	t.Run("Source", func(t *testing.T) {
		source := []byte("hello world")
		r := NewReader(source)
		if string(r.Source()) != string(source) {
			t.Errorf("got %q, want %q", r.Source(), source)
		}
	})

	t.Run("Peek", func(t *testing.T) {
		r := NewReader([]byte("hello"))
		if r.Peek() != 'h' {
			t.Errorf("got %c, want 'h'", r.Peek())
		}
		r.Advance(1)
		if r.Peek() != 'e' {
			t.Errorf("got %c, want 'e'", r.Peek())
		}
	})

	t.Run("Peek EOF", func(t *testing.T) {
		r := NewReader([]byte(""))
		if r.Peek() != EOF {
			t.Errorf("got %v, want EOF", r.Peek())
		}
	})

	t.Run("Value", func(t *testing.T) {
		r := NewReader([]byte("hello world"))
		seg := Segment{Start: 0, Stop: 5}
		if string(r.Value(seg)) != "hello" {
			t.Errorf("got %q, want %q", r.Value(seg), "hello")
		}
	})

	t.Run("PeekLine", func(t *testing.T) {
		r := NewReader([]byte("hello\nworld"))
		line, seg := r.PeekLine()
		if string(line) != "hello\n" {
			t.Errorf("got %q, want %q", line, "hello\n")
		}
		if seg.Start != 0 || seg.Stop != 6 {
			t.Errorf("got segment %+v, want {Start:0 Stop:6}", seg)
		}
	})

	t.Run("PrecendingCharacter", func(t *testing.T) {
		r := NewReader([]byte("abc"))
		r.Advance(1)
		if r.PrecendingCharacter() != 'a' {
			t.Errorf("got %c, want 'a'", r.PrecendingCharacter())
		}
	})

	t.Run("PrecendingCharacter at start", func(t *testing.T) {
		r := NewReader([]byte("abc"))
		if r.PrecendingCharacter() != '\n' {
			t.Errorf("got %c, want newline", r.PrecendingCharacter())
		}
	})

	t.Run("LineOffset", func(t *testing.T) {
		r := NewReader([]byte("hello world"))
		r.Advance(6)
		offset := r.LineOffset()
		if offset != 6 {
			t.Errorf("got offset %d, want 6", offset)
		}
	})

	t.Run("Position and SetPosition", func(t *testing.T) {
		r := NewReader([]byte("hello\nworld"))
		line, seg := r.Position()
		if line != 0 {
			t.Errorf("got line %d, want 0", line)
		}

		r.SetPosition(1, Segment{Start: 6, Stop: 11})
		line, seg = r.Position()
		if line != 1 || seg.Start != 6 {
			t.Errorf("got line=%d seg=%+v, want line=1 seg.Start=6", line, seg)
		}
	})

	t.Run("SetPadding", func(t *testing.T) {
		r := NewReader([]byte("hello"))
		r.SetPadding(3)
		_, seg := r.Position()
		if seg.Padding != 3 {
			t.Errorf("got padding %d, want 3", seg.Padding)
		}
	})
}

func TestReader_Advance(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		advances      []int
		expectedPeeks []byte
		expectedLines []int
	}{
		{
			name:          "simple advance",
			source:        "hello",
			advances:      []int{1, 1, 1},
			expectedPeeks: []byte{'e', 'l', 'l'},
			expectedLines: []int{0, 0, 0},
		},
		{
			name:          "advance across newline",
			source:        "hello\nworld",
			advances:      []int{6},
			expectedPeeks: []byte{'w'},
			expectedLines: []int{1},
		},
		{
			name:          "advance to EOF",
			source:        "hi",
			advances:      []int{2},
			expectedPeeks: []byte{EOF},
			expectedLines: []int{0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewReader([]byte(tc.source))
			for i, adv := range tc.advances {
				r.Advance(adv)
				if r.Peek() != tc.expectedPeeks[i] {
					t.Errorf("after advance %d: got peek %c, want %c",
						i, r.Peek(), tc.expectedPeeks[i])
				}
				line, _ := r.Position()
				if line != tc.expectedLines[i] {
					t.Errorf("after advance %d: got line %d, want %d",
						i, line, tc.expectedLines[i])
				}
			}
		})
	}
}

func TestReader_AdvanceAndSetPadding(t *testing.T) {
	r := NewReader([]byte("hello"))
	r.AdvanceAndSetPadding(2, 3)

	_, seg := r.Position()
	if seg.Start != 2 {
		t.Errorf("got start %d, want 2", seg.Start)
	}
	if seg.Padding != 3 {
		t.Errorf("got padding %d, want 3", seg.Padding)
	}
}

func TestReader_AdvanceToEOL(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected byte
	}{
		{
			name:     "advance to newline",
			source:   "hello\nworld",
			expected: '\n',
		},
		{
			name:     "advance to EOF",
			source:   "hello",
			expected: EOF,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewReader([]byte(tc.source))
			r.AdvanceToEOL()
			if r.Peek() != tc.expected {
				t.Errorf("got %c, want %c", r.Peek(), tc.expected)
			}
		})
	}
}

func TestReader_AdvanceLine(t *testing.T) {
	r := NewReader([]byte("hello\nworld\ntest"))
	line, _ := r.Position()
	if line != 0 {
		t.Errorf("initial line: got %d, want 0", line)
	}

	r.AdvanceLine()
	line, _ = r.Position()
	if line != 1 {
		t.Errorf("after first advance: got %d, want 1", line)
	}
	if r.Peek() != 'w' {
		t.Errorf("after first advance: got %c, want 'w'", r.Peek())
	}

	r.AdvanceLine()
	line, _ = r.Position()
	if line != 2 {
		t.Errorf("after second advance: got %d, want 2", line)
	}
}

func TestReader_SkipSpaces(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expectedChars int
		expectedOk    bool
		expectedPeek  byte
	}{
		{
			name:          "skip leading spaces",
			source:        "   hello",
			expectedChars: 3,
			expectedOk:    true,
			expectedPeek:  'h',
		},
		{
			name:          "no spaces to skip",
			source:        "hello",
			expectedChars: 0,
			expectedOk:    true,
			expectedPeek:  'h',
		},
		{
			name:          "only spaces then EOF",
			source:        "   ",
			expectedChars: 3,
			expectedOk:    false,
			expectedPeek:  EOF,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewReader([]byte(tc.source))
			_, chars, ok := r.SkipSpaces()
			if chars != tc.expectedChars {
				t.Errorf("got chars %d, want %d", chars, tc.expectedChars)
			}
			if ok != tc.expectedOk {
				t.Errorf("got ok %v, want %v", ok, tc.expectedOk)
			}
			if r.Peek() != tc.expectedPeek {
				t.Errorf("got peek %c, want %c", r.Peek(), tc.expectedPeek)
			}
		})
	}
}

func TestReader_SkipBlankLines(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expectedLines int
		expectedOk    bool
		expectedPeek  byte
	}{
		{
			name:          "skip blank lines",
			source:        "\n\n\nhello",
			expectedLines: 3,
			expectedOk:    true,
			expectedPeek:  'h',
		},
		{
			name:          "no blank lines",
			source:        "hello\nworld",
			expectedLines: 0,
			expectedOk:    true,
			expectedPeek:  'h',
		},
		{
			name:          "only blank lines",
			source:        "\n\n",
			expectedLines: 2,
			expectedOk:    false,
			expectedPeek:  EOF,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewReader([]byte(tc.source))
			_, lines, ok := r.SkipBlankLines()
			if lines != tc.expectedLines {
				t.Errorf("got lines %d, want %d", lines, tc.expectedLines)
			}
			if ok != tc.expectedOk {
				t.Errorf("got ok %v, want %v", ok, tc.expectedOk)
			}
			if r.Peek() != tc.expectedPeek {
				t.Errorf("got peek %c, want %c", r.Peek(), tc.expectedPeek)
			}
		})
	}
}

func TestReader_Match(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		pattern  string
		expected bool
		peekAt   byte
	}{
		{
			name:     "match found",
			source:   "hello world",
			pattern:  `hello`,
			expected: true,
			peekAt:   ' ',
		},
		{
			name:     "match not found",
			source:   "hello world",
			pattern:  `goodbye`,
			expected: false,
			peekAt:   'h',
		},
		{
			name:     "partial match",
			source:   "hello world",
			pattern:  `hel`,
			expected: true,
			peekAt:   'l',
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewReader([]byte(tc.source))
			reg := regexp.MustCompile(tc.pattern)
			result := r.Match(reg)
			if result != tc.expected {
				t.Errorf("got %v, want %v", result, tc.expected)
			}
			if r.Peek() != tc.peekAt {
				t.Errorf("got peek %c, want %c", r.Peek(), tc.peekAt)
			}
		})
	}
}

func TestReader_FindClosure(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		opener  byte
		closer  byte
		options FindClosureOptions
		wantOk  bool
	}{
		{
			name:    "simple closure",
			source:  "(hello)",
			opener:  '(',
			closer:  ')',
			options: FindClosureOptions{Advance: true},
			wantOk:  true,
		},
		{
			name:    "nested closure",
			source:  "((hello))",
			opener:  '(',
			closer:  ')',
			options: FindClosureOptions{Nesting: true, Advance: true},
			wantOk:  true,
		},
		{
			name:    "multiline closure",
			source:  "(hello\nworld)",
			opener:  '(',
			closer:  ')',
			options: FindClosureOptions{Newline: true, Advance: true},
			wantOk:  true,
		},
		{
			name:    "no closure found",
			source:  "(hello",
			opener:  '(',
			closer:  ')',
			options: FindClosureOptions{Advance: false},
			wantOk:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := NewReader([]byte(tc.source))
			r.Advance(1)
			_, ok := r.FindClosure(tc.opener, tc.closer, tc.options)
			if ok != tc.wantOk {
				t.Errorf("got ok=%v, want %v", ok, tc.wantOk)
			}
		})
	}
}
