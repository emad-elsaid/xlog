package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestMatchesSetextHeadingBar(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantChar  byte
		wantMatch bool
	}{
		{
			name:      "level 1 equals",
			line:      "===",
			wantChar:  '=',
			wantMatch: true,
		},
		{
			name:      "level 2 dashes",
			line:      "---",
			wantChar:  '-',
			wantMatch: true,
		},
		{
			name:      "level 1 with trailing spaces",
			line:      "=== ",
			wantChar:  '=',
			wantMatch: true,
		},
		{
			name:      "level 1 with leading space (1)",
			line:      " ===",
			wantChar:  '=',
			wantMatch: true,
		},
		{
			name:      "level 1 with leading spaces (3 max)",
			line:      "   ===",
			wantChar:  '=',
			wantMatch: true,
		},
		{
			name:      "too many leading spaces (>3)",
			line:      "    ===",
			wantChar:  0,
			wantMatch: false,
		},
		{
			name:      "level 2 with spaces",
			line:      "  ---  ",
			wantChar:  '-',
			wantMatch: true,
		},
		{
			name:      "mixed markers invalid",
			line:      "=-=",
			wantChar:  0,
			wantMatch: false,
		},
		{
			name:      "single equals",
			line:      "=",
			wantChar:  '=',
			wantMatch: true,
		},
		{
			name:      "single dash",
			line:      "-",
			wantChar:  '-',
			wantMatch: true,
		},
		{
			name:      "empty line",
			line:      "",
			wantChar:  0,
			wantMatch: false,
		},
		{
			name:      "only spaces",
			line:      "   ",
			wantChar:  0,
			wantMatch: false,
		},
		{
			name:      "text not markers",
			line:      "text",
			wantChar:  0,
			wantMatch: false,
		},
		{
			name:      "partial match with text",
			line:      "==text",
			wantChar:  0,
			wantMatch: false,
		},
		{
			name:      "long equals line",
			line:      "================",
			wantChar:  '=',
			wantMatch: true,
		},
		{
			name:      "long dashes line",
			line:      "----------------",
			wantChar:  '-',
			wantMatch: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotChar, gotMatch := matchesSetextHeadingBar([]byte(tc.line))
			if gotChar != tc.wantChar {
				t.Errorf("matchesSetextHeadingBar(%q) char = %c, want %c", tc.line, gotChar, tc.wantChar)
			}
			if gotMatch != tc.wantMatch {
				t.Errorf("matchesSetextHeadingBar(%q) match = %v, want %v", tc.line, gotMatch, tc.wantMatch)
			}
		})
	}
}

func TestSetextHeadingParser_Trigger(t *testing.T) {
	parser := NewSetextHeadingParser()
	trigger := parser.Trigger()

	expectedTriggers := []byte{'-', '='}
	if len(trigger) != len(expectedTriggers) {
		t.Fatalf("Trigger() length = %d, want %d", len(trigger), len(expectedTriggers))
	}

	for i, expected := range expectedTriggers {
		if trigger[i] != expected {
			t.Errorf("Trigger()[%d] = %c, want %c", i, trigger[i], expected)
		}
	}
}

func TestSetextHeadingParser_Open(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLevel int
		wantState State
		wantNil   bool
	}{
		{
			name:      "level 1 heading with equals",
			input:     "Heading Text\n===\n",
			wantLevel: 1,
			wantState: NoChildren | RequireParagraph,
		},
		{
			name:      "level 2 heading with dashes",
			input:     "Heading Text\n---\n",
			wantLevel: 2,
			wantState: NoChildren | RequireParagraph,
		},
		{
			name:      "level 1 with spaces",
			input:     "Heading\n  ===  \n",
			wantLevel: 1,
			wantState: NoChildren | RequireParagraph,
		},
		{
			name:      "level 2 with leading spaces",
			input:     "Heading\n   ---\n",
			wantLevel: 2,
			wantState: NoChildren | RequireParagraph,
		},
		{
			name:      "multiline paragraph before heading",
			input:     "First line\nSecond line\n===\n",
			wantLevel: 1,
			wantState: NoChildren | RequireParagraph,
		},
		{
			name:      "single word heading level 1",
			input:     "Title\n=\n",
			wantLevel: 1,
			wantState: NoChildren | RequireParagraph,
		},
		{
			name:      "single word heading level 2",
			input:     "Title\n-\n",
			wantLevel: 2,
			wantState: NoChildren | RequireParagraph,
		},
		{
			name:    "no paragraph before marker (invalid)",
			input:   "===\n",
			wantNil: true,
		},
		{
			name:    "too many leading spaces (>3)",
			input:   "Heading\n    ===\n",
			wantNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewSetextHeadingParser()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()
			doc := ast.NewDocument()

			// Simulate parsing the paragraph first
			if !tc.wantNil {
				para := ast.NewParagraph()
				doc.AppendChild(doc, para)

				// Read first line(s) into paragraph
				lines := 0
				for {
					line, segment := reader.PeekLine()
					if len(line) == 0 {
						break
					}

					// Check if this line is a setext marker
					if _, ok := matchesSetextHeadingBar(line); ok {
						break
					}

					para.Lines().Append(segment)
					reader.Advance(segment.Len())
					lines++

					// For single-line tests, just read one line
					if lines > 0 && (tc.name == "single word heading level 1" ||
						tc.name == "single word heading level 2" ||
						tc.name == "level 1 heading with equals" ||
						tc.name == "level 2 heading with dashes" ||
						tc.name == "level 1 with spaces" ||
						tc.name == "level 2 with leading spaces") {
						break
					}

					// For multiline, read two lines
					if lines >= 2 {
						break
					}
				}

				// Set up parser context with the paragraph as last opened block
				pc.SetOpenedBlocks([]Block{{Node: para}})
			}

			node, state := parser.Open(doc, reader, pc)

			if tc.wantNil {
				if node != nil {
					t.Errorf("Open() = %v, want nil", node)
				}
				return
			}

			if node == nil {
				t.Fatal("Open() = nil, want heading node")
			}

			heading, ok := node.(*ast.Heading)
			if !ok {
				t.Fatalf("Open() node type = %T, want *ast.Heading", node)
			}

			if heading.Level != tc.wantLevel {
				t.Errorf("heading.Level = %d, want %d", heading.Level, tc.wantLevel)
			}

			if state != tc.wantState {
				t.Errorf("Open() state = %v, want %v", state, tc.wantState)
			}
		})
	}
}

func TestSetextHeadingParser_Continue(t *testing.T) {
	parser := NewSetextHeadingParser()
	reader := text.NewReader([]byte("===\n"))
	pc := NewContext()
	heading := ast.NewHeading(1)

	state := parser.Continue(heading, reader, pc)

	if state != Close {
		t.Errorf("Continue() = %v, want Close", state)
	}
}

func TestSetextHeadingParser_Close(t *testing.T) {
	tests := []struct {
		name          string
		paragraphText string
		markerLine    string
		wantLevel     int
		skipOpen      bool // For empty paragraph case that doesn't open
		withAttribute bool
		withAutoID    bool
	}{
		{
			name:          "simple level 1",
			paragraphText: "Heading Text",
			markerLine:    "===",
			wantLevel:     1,
		},
		{
			name:          "simple level 2",
			paragraphText: "Heading Text",
			markerLine:    "---",
			wantLevel:     2,
		},
		{
			name:          "multiline paragraph",
			paragraphText: "First Line\nSecond Line",
			markerLine:    "===",
			wantLevel:     1,
		},
		{
			name:          "empty paragraph should not open heading",
			paragraphText: "",
			markerLine:    "===",
			wantLevel:     1,
			skipOpen:      true, // Empty paragraph won't trigger Open
		},
		{
			name:          "with attributes",
			paragraphText: "Heading {#custom-id}",
			markerLine:    "===",
			wantLevel:     1,
			withAttribute: true,
		},
		{
			name:          "with auto id",
			paragraphText: "My Heading",
			markerLine:    "===",
			wantLevel:     1,
			withAutoID:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Build parser with options
			var opts []HeadingOption
			if tc.withAttribute {
				opts = append(opts, WithHeadingAttribute())
			}
			if tc.withAutoID {
				opts = append(opts, WithAutoHeadingID())
			}

			parser := NewSetextHeadingParser(opts...)

			// Construct input
			input := tc.paragraphText + "\n" + tc.markerLine + "\n"
			reader := text.NewReader([]byte(input))
			pc := NewContext()
			doc := ast.NewDocument()

			// Create paragraph
			para := ast.NewParagraph()
			doc.AppendChild(doc, para)

			// Read paragraph lines
			if tc.paragraphText != "" {
				for {
					line, segment := reader.PeekLine()
					if len(line) == 0 {
						break
					}
					if _, ok := matchesSetextHeadingBar(line); ok {
						break
					}
					para.Lines().Append(segment)
					reader.Advance(segment.Len())
				}
			}

			// Set up context
			pc.SetOpenedBlocks([]Block{{Node: para}})

			// Open the heading
			heading, _ := parser.Open(doc, reader, pc)

			// Skip rest of test if we expect Open to return nil
			if tc.skipOpen {
				if heading != nil {
					t.Errorf("Open() = %v, want nil for empty paragraph", heading)
				}
				return
			}

			if heading == nil {
				t.Fatal("Open() returned nil")
			}

			// Close the heading
			parser.Close(heading, reader, pc)

			// Verify the heading was processed correctly
			h, ok := heading.(*ast.Heading)
			if !ok {
				t.Fatalf("heading type = %T, want *ast.Heading", heading)
			}

			if h.Level != tc.wantLevel {
				t.Errorf("heading.Level = %d, want %d", h.Level, tc.wantLevel)
			}

			// Check if attributes were parsed
			if tc.withAttribute && tc.paragraphText == "Heading {#custom-id}" {
				id, ok := h.AttributeString("id")
				if !ok {
					t.Error("expected id attribute to be set")
				} else if string(id.([]byte)) != "custom-id" {
					t.Errorf("id = %s, want 'custom-id'", id)
				}
			}

			// Check if auto-ID was generated
			if tc.withAutoID && !tc.withAttribute {
				_, ok := h.AttributeString("id")
				if !ok {
					t.Error("expected auto-generated id attribute")
				}
			}
		})
	}
}

func TestSetextHeadingParser_CanInterruptParagraph(t *testing.T) {
	parser := NewSetextHeadingParser()
	if !parser.CanInterruptParagraph() {
		t.Error("CanInterruptParagraph() = false, want true")
	}
}

func TestSetextHeadingParser_CanAcceptIndentedLine(t *testing.T) {
	parser := NewSetextHeadingParser()
	if parser.CanAcceptIndentedLine() {
		t.Error("CanAcceptIndentedLine() = true, want false")
	}
}

func TestSetextHeadingParser_WithOptions(t *testing.T) {
	t.Run("with attribute option", func(t *testing.T) {
		parser := NewSetextHeadingParser(WithHeadingAttribute())
		p, ok := parser.(*setextHeadingParser)
		if !ok {
			t.Fatal("parser type assertion failed")
		}
		if !p.Attribute {
			t.Error("Attribute = false, want true")
		}
	})

	t.Run("with auto heading id option", func(t *testing.T) {
		parser := NewSetextHeadingParser(WithAutoHeadingID())
		p, ok := parser.(*setextHeadingParser)
		if !ok {
			t.Fatal("parser type assertion failed")
		}
		if !p.AutoHeadingID {
			t.Error("AutoHeadingID = false, want true")
		}
	})

	t.Run("with both options", func(t *testing.T) {
		parser := NewSetextHeadingParser(WithHeadingAttribute(), WithAutoHeadingID())
		p, ok := parser.(*setextHeadingParser)
		if !ok {
			t.Fatal("parser type assertion failed")
		}
		if !p.Attribute {
			t.Error("Attribute = false, want true")
		}
		if !p.AutoHeadingID {
			t.Error("AutoHeadingID = false, want true")
		}
	})
}
