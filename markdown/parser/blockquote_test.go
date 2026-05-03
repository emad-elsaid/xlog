package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestBlockquoteParser_Trigger(t *testing.T) {
	parser := NewBlockquoteParser()
	trigger := parser.Trigger()
	if len(trigger) != 1 || trigger[0] != '>' {
		t.Errorf("Trigger() = %v, want ['>']", trigger)
	}
}

func TestBlockquoteParser_Open(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantState   State
		wantNil     bool
		blockOffset int
	}{
		{
			name:        "simple blockquote",
			input:       "> Hello\n",
			wantState:   HasChildren,
			blockOffset: 0,
		},
		{
			name:        "blockquote with space after marker",
			input:       ">  Hello world\n",
			wantState:   HasChildren,
			blockOffset: 0,
		},
		{
			name:        "blockquote with tab after marker",
			input:       ">\tHello tab\n",
			wantState:   HasChildren,
			blockOffset: 0,
		},
		{
			name:        "blockquote marker only",
			input:       ">\n",
			wantState:   HasChildren,
			blockOffset: 0,
		},
		{
			name:        "blockquote without space",
			input:       ">NoSpace\n",
			wantState:   HasChildren,
			blockOffset: 0,
		},
		{
			name:        "not a blockquote - no marker",
			input:       "Just text\n",
			wantNil:     true,
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "not a blockquote - indent too large",
			input:       "    > Indented 4 spaces\n",
			wantNil:     true,
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "blockquote with 1 space indent",
			input:       " > One space\n",
			wantState:   HasChildren,
			blockOffset: 0,
		},
		{
			name:        "blockquote with 2 space indent",
			input:       "  > Two spaces\n",
			wantState:   HasChildren,
			blockOffset: 0,
		},
		{
			name:        "blockquote with 3 space indent",
			input:       "   > Three spaces\n",
			wantState:   HasChildren,
			blockOffset: 0,
		},
		{
			name:        "empty input",
			input:       "",
			wantNil:     true,
			wantState:   NoChildren,
			blockOffset: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewBlockquoteParser()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()
			pc.SetBlockOffset(tc.blockOffset)

			node, state := parser.Open(ast.NewDocument(), reader, pc)

			if tc.wantNil {
				if node != nil {
					t.Errorf("Open() returned node = %v, want nil", node)
				}
			} else {
				if node == nil {
					t.Errorf("Open() returned nil node, want blockquote")
				} else if node.Kind() != ast.KindBlockquote {
					t.Errorf("Open() node.Kind() = %v, want KindBlockquote", node.Kind())
				}
			}

			if state != tc.wantState {
				t.Errorf("Open() state = %v, want %v", state, tc.wantState)
			}
		})
	}
}

func TestBlockquoteParser_Continue(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantState State
	}{
		{
			name:      "continue with blockquote marker",
			input:     "> Continued line\n",
			wantState: Continue | HasChildren,
		},
		{
			name:      "continue with marker only",
			input:     ">\n",
			wantState: Continue | HasChildren,
		},
		{
			name:      "close on non-blockquote line",
			input:     "Regular text\n",
			wantState: Close,
		},
		{
			name:      "close on empty line",
			input:     "\n",
			wantState: Close,
		},
		{
			name:      "continue with indented marker",
			input:     " > Still blockquote\n",
			wantState: Continue | HasChildren,
		},
		{
			name:      "close on over-indented line",
			input:     "    > Too indented\n",
			wantState: Close,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewBlockquoteParser()
			node := ast.NewBlockquote()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()

			state := parser.Continue(node, reader, pc)

			if state != tc.wantState {
				t.Errorf("Continue() = %v, want %v", state, tc.wantState)
			}
		})
	}
}

func TestBlockquoteParser_Close(t *testing.T) {
	parser := NewBlockquoteParser()
	node := ast.NewBlockquote()
	reader := text.NewReader([]byte("> test\n"))
	pc := NewContext()

	// Close should not panic and should do nothing
	parser.Close(node, reader, pc)
}

func TestBlockquoteParser_CanInterruptParagraph(t *testing.T) {
	parser := NewBlockquoteParser()
	if !parser.CanInterruptParagraph() {
		t.Errorf("CanInterruptParagraph() = false, want true")
	}
}

func TestBlockquoteParser_CanAcceptIndentedLine(t *testing.T) {
	parser := NewBlockquoteParser()
	if parser.CanAcceptIndentedLine() {
		t.Errorf("CanAcceptIndentedLine() = true, want false")
	}
}

func TestBlockquoteParser_process(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantProcessed bool
		wantAdvanced  int
	}{
		{
			name:          "simple marker with content",
			input:         "> content\n",
			wantProcessed: true,
			wantAdvanced:  2, // '>' + ' '
		},
		{
			name:          "marker only",
			input:         ">\n",
			wantProcessed: true,
			wantAdvanced:  1, // '>'
		},
		{
			name:          "marker with tab",
			input:         ">\tcontent\n",
			wantProcessed: true,
			wantAdvanced:  2, // '>' + '\t'
		},
		{
			name:          "no marker",
			input:         "content\n",
			wantProcessed: false,
			wantAdvanced:  0,
		},
		{
			name:          "over-indented",
			input:         "    > content\n",
			wantProcessed: false,
			wantAdvanced:  0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewBlockquoteParser().(*blockquoteParser)
			reader := text.NewReader([]byte(tc.input))
			initialPos := reader.LineOffset()

			processed := parser.process(reader)

			if processed != tc.wantProcessed {
				t.Errorf("process() = %v, want %v", processed, tc.wantProcessed)
			}

			if tc.wantProcessed {
				advanced := reader.LineOffset() - initialPos
				if advanced != tc.wantAdvanced {
					t.Errorf("Advanced %d positions, want %d", advanced, tc.wantAdvanced)
				}
			}
		})
	}
}
