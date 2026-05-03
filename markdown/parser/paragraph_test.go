package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestParagraphParser_Trigger(t *testing.T) {
	parser := NewParagraphParser()
	trigger := parser.Trigger()
	if trigger != nil {
		t.Errorf("Trigger() = %v, want nil (paragraph has no specific trigger)", trigger)
	}
}

func TestParagraphParser_Open(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantState State
		wantNil   bool
	}{
		{
			name:      "simple paragraph",
			input:     "Hello world\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "paragraph with leading spaces",
			input:     "   Content with spaces\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "paragraph with tabs",
			input:     "\tTabbed content\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "blank line returns nil",
			input:     "\n",
			wantState: NoChildren,
			wantNil:   true,
		},
		{
			name:      "spaces only returns nil",
			input:     "   \n",
			wantState: NoChildren,
			wantNil:   true,
		},
		{
			name:      "tabs only returns nil",
			input:     "\t\t\n",
			wantState: NoChildren,
			wantNil:   true,
		},
		{
			name:      "empty input returns nil",
			input:     "",
			wantState: NoChildren,
			wantNil:   true,
		},
		{
			name:      "paragraph with mixed whitespace",
			input:     " \t Content\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "paragraph with special characters",
			input:     "Hello *world* [link](url)\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "paragraph with unicode",
			input:     "Café résumé 日本語\n",
			wantState: NoChildren,
			wantNil:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParagraphParser()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()

			node, state := parser.Open(ast.NewDocument(), reader, pc)

			if tc.wantNil {
				if node != nil {
					t.Errorf("Open() returned node = %v, want nil", node)
				}
			} else {
				if node == nil {
					t.Errorf("Open() returned nil node, want paragraph")
				} else if node.Kind() != ast.KindParagraph {
					t.Errorf("Open() node.Kind() = %v, want KindParagraph", node.Kind())
				}

				// Verify the paragraph has content
				para := node.(*ast.Paragraph)
				if para.Lines().Len() == 0 {
					t.Errorf("Paragraph has no lines")
				}

				// Note: Open calls AdvanceToEOL which positions reader AFTER newline
				// Remaining line may be empty or the next line - both are acceptable
			}

			if state != tc.wantState {
				t.Errorf("Open() state = %v, want %v", state, tc.wantState)
			}
		})
	}
}

func TestParagraphParser_Continue(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantState State
	}{
		{
			name:      "continue with text",
			input:     "More content\n",
			wantState: Continue | NoChildren,
		},
		{
			name:      "continue with leading spaces",
			input:     "  Indented content\n",
			wantState: Continue | NoChildren,
		},
		{
			name:      "close on blank line",
			input:     "\n",
			wantState: Close,
		},
		{
			name:      "close on spaces only",
			input:     "   \n",
			wantState: Close,
		},
		{
			name:      "close on tabs only",
			input:     "\t\t\n",
			wantState: Close,
		},
		{
			name:      "continue with special characters",
			input:     "**bold** and *italic*\n",
			wantState: Continue | NoChildren,
		},
		{
			name:      "continue with unicode",
			input:     "国際化テスト\n",
			wantState: Continue | NoChildren,
		},
		{
			name:      "continue with mixed whitespace prefix",
			input:     " \t Text\n",
			wantState: Continue | NoChildren,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParagraphParser()
			node := ast.NewParagraph()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()

			initialLineCount := node.Lines().Len()
			state := parser.Continue(node, reader, pc)

			if state != tc.wantState {
				t.Errorf("Continue() = %v, want %v", state, tc.wantState)
			}

			// If continuing, verify line was added
			if state == (Continue | NoChildren) {
				if node.Lines().Len() != initialLineCount+1 {
					t.Errorf("Continue() did not add line: before=%d, after=%d",
						initialLineCount, node.Lines().Len())
				}
			}
		})
	}
}

func TestParagraphParser_Close(t *testing.T) {
	tests := []struct {
		name            string
		lines           []string
		wantLineCount   int
		wantFirstLine   string
		wantLastLine    string
		shouldBeRemoved bool
	}{
		{
			name:          "single line paragraph",
			lines:         []string{"Hello world"},
			wantLineCount: 1,
			wantFirstLine: "Hello world",
			wantLastLine:  "Hello world",
		},
		{
			name:          "multiple lines",
			lines:         []string{"First line", "Second line", "Third line"},
			wantLineCount: 3,
			wantFirstLine: "First line",
			wantLastLine:  "Third line",
		},
		{
			name:          "lines with leading spaces trimmed",
			lines:         []string{"  First", "  Second", "  Third"},
			wantLineCount: 3,
			wantFirstLine: "First",
			wantLastLine:  "Third",
		},
		{
			name:          "lines with trailing spaces trimmed on last only",
			lines:         []string{"First  ", "Second  ", "Third  "},
			wantLineCount: 3,
			wantFirstLine: "First  ", // First line keeps trailing spaces
			wantLastLine:  "Third",   // Last line has trailing spaces trimmed
		},
		{
			name:          "mixed leading/trailing spaces",
			lines:         []string{"  First  ", "  Second  ", "  Third  "},
			wantLineCount: 3,
			wantFirstLine: "First  ", // Leading trimmed, trailing kept (not last line)
			wantLastLine:  "Third",   // Both leading and trailing trimmed (last line)
		},
		{
			name:            "empty paragraph removed",
			lines:           []string{},
			shouldBeRemoved: true,
		},
		{
			name:          "single line with spaces",
			lines:         []string{"   Content   "},
			wantLineCount: 1,
			wantFirstLine: "Content",
			wantLastLine:  "Content",
		},
		{
			name:          "unicode content preserved with spacing",
			lines:         []string{"  日本語  ", "  中文  "},
			wantLineCount: 2,
			wantFirstLine: "日本語  ", // Leading trimmed, trailing kept
			wantLastLine:  "中文",    // Both trimmed (last line)
		},
		{
			name:          "tabs in content trimmed appropriately",
			lines:         []string{"\tFirst\t", "\tSecond\t"},
			wantLineCount: 2,
			wantFirstLine: "First\t", // Leading trimmed, trailing kept
			wantLastLine:  "Second",  // Both trimmed (last line)
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParagraphParser()
			parent := ast.NewDocument()
			node := ast.NewParagraph()
			parent.AppendChild(parent, node)

			// Construct full input source for proper segment handling
			inputLines := append([]string(nil), tc.lines...)
			var source []byte
			for _, line := range inputLines {
				source = append(source, []byte(line+"\n")...)
			}

			reader := text.NewReader(source)

			// Add lines to paragraph with proper segments
			offset := 0
			for _, line := range tc.lines {
				lineBytes := []byte(line + "\n")
				segment := text.NewSegment(offset, offset+len(lineBytes)-1)
				node.Lines().Append(segment)
				offset += len(lineBytes)
			}

			pc := NewContext()
			parser.Close(node, reader, pc)

			if tc.shouldBeRemoved {
				// Verify node was removed from parent
				if parent.ChildCount() != 0 {
					t.Errorf("Close() did not remove empty paragraph, child count = %d",
						parent.ChildCount())
				}
				return
			}

			// Verify line count
			if node.Lines().Len() != tc.wantLineCount {
				t.Errorf("Close() resulted in %d lines, want %d",
					node.Lines().Len(), tc.wantLineCount)
			}

			if tc.wantLineCount > 0 {
				// Verify first line content (may have trailing newline)
				firstSeg := node.Lines().At(0)
				firstLine := string(firstSeg.Value(source))
				firstLineTrimmed := firstLine
				if len(firstLine) > 0 && firstLine[len(firstLine)-1] == '\n' {
					firstLineTrimmed = firstLine[:len(firstLine)-1]
				}
				if firstLineTrimmed != tc.wantFirstLine {
					t.Errorf("First line = %q, want %q", firstLineTrimmed, tc.wantFirstLine)
				}

				// Verify last line content (may have trailing newline)
				lastSeg := node.Lines().At(tc.wantLineCount - 1)
				lastLine := string(lastSeg.Value(source))
				lastLineTrimmed := lastLine
				if len(lastLine) > 0 && lastLine[len(lastLine)-1] == '\n' {
					lastLineTrimmed = lastLine[:len(lastLine)-1]
				}
				if lastLineTrimmed != tc.wantLastLine {
					t.Errorf("Last line = %q, want %q", lastLineTrimmed, tc.wantLastLine)
				}
			}
		})
	}
}

func TestParagraphParser_CanInterruptParagraph(t *testing.T) {
	parser := NewParagraphParser()
	if parser.CanInterruptParagraph() {
		t.Errorf("CanInterruptParagraph() = true, want false (paragraph cannot interrupt paragraph)")
	}
}

func TestParagraphParser_CanAcceptIndentedLine(t *testing.T) {
	parser := NewParagraphParser()
	if parser.CanAcceptIndentedLine() {
		t.Errorf("CanAcceptIndentedLine() = true, want false")
	}
}

func TestParagraphParser_DefaultInstance(t *testing.T) {
	// Test that NewParagraphParser returns the shared default instance
	p1 := NewParagraphParser()
	p2 := NewParagraphParser()

	if p1 != p2 {
		t.Errorf("NewParagraphParser() returns different instances, want same default instance")
	}
}

func TestParagraphParser_Integration(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantLineCount int
	}{
		{
			name: "multi-line paragraph",
			input: `This is the first line.
This is the second line.
This is the third line.
`,
			wantLineCount: 3,
		},
		{
			name: "paragraph with blank terminator",
			input: `Line one.
Line two.

`,
			wantLineCount: 2,
		},
		{
			name: "single line paragraph",
			input: `Just one line.
`,
			wantLineCount: 1,
		},
		{
			name: "paragraph with various formatting",
			input: `   Leading spaces.
Trailing spaces.   
	Tabs here.
Normal line.
`,
			wantLineCount: 4,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewParagraphParser()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()
			parent := ast.NewDocument()

			// Open the paragraph
			node, state := parser.Open(parent, reader, pc)
			if node == nil {
				t.Fatal("Open() returned nil")
			}
			if state != NoChildren {
				t.Errorf("Open() state = %v, want NoChildren", state)
			}

			// Continue reading lines until we close
			for {
				reader.AdvanceLine()
				line, _ := reader.PeekLine()
				if len(line) == 0 {
					break
				}

				state := parser.Continue(node, reader, pc)
				if state == Close {
					break
				}
				if state != (Continue | NoChildren) {
					t.Errorf("Continue() state = %v, want Continue|NoChildren", state)
				}
			}

			// Close the paragraph
			parent.AppendChild(parent, node)
			parser.Close(node, reader, pc)

			// Verify line count
			if node.Lines().Len() != tc.wantLineCount {
				t.Errorf("Final paragraph has %d lines, want %d",
					node.Lines().Len(), tc.wantLineCount)
			}
		})
	}
}
