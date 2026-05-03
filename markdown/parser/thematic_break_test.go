package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestThematicBreakParser_Trigger(t *testing.T) {
	parser := NewThematicBreakParser()
	trigger := parser.Trigger()
	expected := []byte{'-', '*', '_'}

	if len(trigger) != len(expected) {
		t.Errorf("Trigger() length = %d, want %d", len(trigger), len(expected))
		return
	}

	for i, c := range expected {
		if trigger[i] != c {
			t.Errorf("Trigger()[%d] = %c, want %c", i, trigger[i], c)
		}
	}
}

func TestThematicBreakParser_Open(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantState State
		wantNil   bool
	}{
		{
			name:      "simple asterisk thematic break",
			input:     "***\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "simple dash thematic break",
			input:     "---\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "simple underscore thematic break",
			input:     "___\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "asterisk with spaces",
			input:     "* * *\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "dash with spaces",
			input:     "- - -\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "underscore with spaces",
			input:     "_ _ _\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "many asterisks",
			input:     "*****\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "many dashes",
			input:     "-----\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "mixed spaces and characters",
			input:     "  * * *  \n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "with leading spaces (up to 3)",
			input:     "   ***\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "not thematic break - only 2 characters",
			input:     "**\n",
			wantNil:   true,
			wantState: NoChildren,
		},
		{
			name:      "not thematic break - mixed characters",
			input:     "***___\n",
			wantNil:   true,
			wantState: NoChildren,
		},
		{
			name:      "not thematic break - mixed dash and asterisk",
			input:     "-*-\n",
			wantNil:   true,
			wantState: NoChildren,
		},
		{
			name:      "not thematic break - too much indentation",
			input:     "    ***\n",
			wantNil:   true,
			wantState: NoChildren,
		},
		{
			name:      "not thematic break - invalid character",
			input:     "###\n",
			wantNil:   true,
			wantState: NoChildren,
		},
		{
			name:      "not thematic break - text before",
			input:     "text***\n",
			wantNil:   true,
			wantState: NoChildren,
		},
		{
			name:      "empty input",
			input:     "",
			wantNil:   true,
			wantState: NoChildren,
		},
		{
			name:      "just spaces",
			input:     "   \n",
			wantNil:   true,
			wantState: NoChildren,
		},
		{
			name:      "single asterisk",
			input:     "*\n",
			wantNil:   true,
			wantState: NoChildren,
		},
		{
			name:      "trailing spaces allowed",
			input:     "***   \n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "tabs as spaces - tab counts as 4 width (too much)",
			input:     "\t***\n",
			wantNil:   true,
			wantState: NoChildren,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewThematicBreakParser()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()

			node, state := parser.Open(ast.NewDocument(), reader, pc)

			if tc.wantNil {
				if node != nil {
					t.Errorf("Open() returned node = %v, want nil", node)
				}
			} else {
				if node == nil {
					t.Errorf("Open() returned nil node, want ThematicBreak")
				} else if node.Kind() != ast.KindThematicBreak {
					t.Errorf("Open() node.Kind() = %v, want KindThematicBreak", node.Kind())
				}
			}

			if state != tc.wantState {
				t.Errorf("Open() state = %v, want %v", state, tc.wantState)
			}
		})
	}
}

func TestThematicBreakParser_Continue(t *testing.T) {
	parser := NewThematicBreakParser()
	node := ast.NewThematicBreak()
	reader := text.NewReader([]byte("any input\n"))
	pc := NewContext()

	state := parser.Continue(node, reader, pc)

	if state != Close {
		t.Errorf("Continue() = %v, want Close", state)
	}
}

func TestThematicBreakParser_Close(t *testing.T) {
	parser := NewThematicBreakParser()
	node := ast.NewThematicBreak()
	reader := text.NewReader([]byte("***\n"))
	pc := NewContext()

	// Close should not panic and should do nothing
	parser.Close(node, reader, pc)
}

func TestThematicBreakParser_CanInterruptParagraph(t *testing.T) {
	parser := NewThematicBreakParser()
	if !parser.CanInterruptParagraph() {
		t.Errorf("CanInterruptParagraph() = false, want true")
	}
}

func TestThematicBreakParser_CanAcceptIndentedLine(t *testing.T) {
	parser := NewThematicBreakParser()
	if parser.CanAcceptIndentedLine() {
		t.Errorf("CanAcceptIndentedLine() = true, want false")
	}
}

func TestIsThematicBreak(t *testing.T) {
	tests := []struct {
		name   string
		line   string
		offset int
		want   bool
	}{
		{
			name:   "three asterisks",
			line:   "***",
			offset: 0,
			want:   true,
		},
		{
			name:   "three dashes",
			line:   "---",
			offset: 0,
			want:   true,
		},
		{
			name:   "three underscores",
			line:   "___",
			offset: 0,
			want:   true,
		},
		{
			name:   "more than three",
			line:   "*****",
			offset: 0,
			want:   true,
		},
		{
			name:   "with spaces between",
			line:   "* * *",
			offset: 0,
			want:   true,
		},
		{
			name:   "with multiple spaces",
			line:   "*  *  *",
			offset: 0,
			want:   true,
		},
		{
			name:   "with leading spaces",
			line:   "  ***",
			offset: 0,
			want:   true,
		},
		{
			name:   "with three spaces leading",
			line:   "   ***",
			offset: 0,
			want:   true,
		},
		{
			name:   "not thematic - four spaces leading",
			line:   "    ***",
			offset: 0,
			want:   false,
		},
		{
			name:   "not thematic - only two",
			line:   "**",
			offset: 0,
			want:   false,
		},
		{
			name:   "not thematic - mixed chars",
			line:   "*_*",
			offset: 0,
			want:   false,
		},
		{
			name:   "not thematic - invalid char",
			line:   "###",
			offset: 0,
			want:   false,
		},
		{
			name:   "empty line",
			line:   "",
			offset: 0,
			want:   false,
		},
		{
			name:   "only spaces",
			line:   "   ",
			offset: 0,
			want:   false,
		},
		{
			name:   "no offset from start",
			line:   "***",
			offset: 0,
			want:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isThematicBreak([]byte(tc.line), tc.offset)
			if got != tc.want {
				t.Errorf("isThematicBreak(%q, %d) = %v, want %v", tc.line, tc.offset, got, tc.want)
			}
		})
	}
}
