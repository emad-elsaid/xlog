package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestCodeBlockParser_Trigger(t *testing.T) {
	parser := NewCodeBlockParser()
	trigger := parser.Trigger()
	if trigger != nil {
		t.Errorf("Trigger() = %v, want nil", trigger)
	}
}

func TestCodeBlockParser_Open(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantState State
		wantNil   bool
	}{
		{
			name:      "simple indented code block",
			input:     "    code line\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "code block with tab",
			input:     "\tcode line\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "insufficient indentation (3 spaces)",
			input:     "   not code\n",
			wantNil:   true,
			wantState: NoChildren,
		},
		{
			name:      "blank line not code block",
			input:     "    \n",
			wantNil:   true,
			wantState: NoChildren,
		},
		{
			name:      "no indentation",
			input:     "not code\n",
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
			name:      "exactly 4 spaces",
			input:     "    func main() {}\n",
			wantState: NoChildren,
			wantNil:   false,
		},
		{
			name:      "more than 4 spaces",
			input:     "      deep indent\n",
			wantState: NoChildren,
			wantNil:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewCodeBlockParser()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()

			node, state := parser.Open(ast.NewDocument(), reader, pc)

			if tc.wantNil {
				if node != nil {
					t.Errorf("Open() returned node = %v, want nil", node)
				}
			} else {
				if node == nil {
					t.Errorf("Open() returned nil node, want code block")
				} else if node.Kind() != ast.KindCodeBlock {
					t.Errorf("Open() node.Kind() = %v, want KindCodeBlock", node.Kind())
				}
				if node.Lines().Len() != 1 {
					t.Errorf("Open() node has %d lines, want 1", node.Lines().Len())
				}
			}

			if state != tc.wantState {
				t.Errorf("Open() state = %v, want %v", state, tc.wantState)
			}
		})
	}
}

func TestCodeBlockParser_Continue(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantState State
	}{
		{
			name:      "continue with indented line",
			input:     "    more code\n",
			wantState: Continue | NoChildren,
		},
		{
			name:      "continue with blank line",
			input:     "    \n",
			wantState: Continue | NoChildren,
		},
		{
			name:      "continue with tab",
			input:     "\tcode\n",
			wantState: Continue | NoChildren,
		},
		{
			name:      "close on non-indented line",
			input:     "not indented\n",
			wantState: Close,
		},
		{
			name:      "close on insufficient indent",
			input:     "   three spaces\n",
			wantState: Close,
		},
		{
			name:      "continue with deep indent",
			input:     "        deep\n",
			wantState: Continue | NoChildren,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewCodeBlockParser()
			node := ast.NewCodeBlock()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()

			state := parser.Continue(node, reader, pc)

			if state != tc.wantState {
				t.Errorf("Continue() = %v, want %v", state, tc.wantState)
			}
		})
	}
}

func TestCodeBlockParser_Close(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedLines int
	}{
		{
			name:          "trim trailing blank lines",
			input:         "    code\n    \n    \n",
			expectedLines: 1,
		},
		{
			name:          "no trailing blank lines",
			input:         "    code1\n    code2\n",
			expectedLines: 2,
		},
		{
			name:          "mixed content and blanks",
			input:         "    code\n    \n    more\n    \n",
			expectedLines: 3,
		},
		{
			name:          "single line no blank",
			input:         "    single\n",
			expectedLines: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewCodeBlockParser()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()

		// Open and continue to build the node naturally
		node, _ := parser.Open(ast.NewDocument(), reader, pc)
		if node == nil {
			t.Fatal("Open() returned nil")
		}
		reader.AdvanceLine()

		// Continue until close or EOF
		for reader.Peek() != text.EOF {
			state := parser.Continue(node, reader, pc)
			if state == Close {
				break
			}
			reader.AdvanceLine()
		}

			// Now test Close trims trailing blanks
			parser.Close(node, reader, pc)

			if node.Lines().Len() != tc.expectedLines {
				t.Errorf("Close() left %d lines, want %d", node.Lines().Len(), tc.expectedLines)
			}
		})
	}
}

func TestCodeBlockParser_CanInterruptParagraph(t *testing.T) {
	parser := NewCodeBlockParser()
	if parser.CanInterruptParagraph() {
		t.Errorf("CanInterruptParagraph() = true, want false")
	}
}

func TestCodeBlockParser_CanAcceptIndentedLine(t *testing.T) {
	parser := NewCodeBlockParser()
	if !parser.CanAcceptIndentedLine() {
		t.Errorf("CanAcceptIndentedLine() = false, want true")
	}
}

func TestCodeBlockParser_PreserveLeadingTab(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLines int
	}{
		{
			name:      "tab-indented code preserves tab",
			input:     "\tfunc main() {}\n",
			wantLines: 1,
		},
		{
			name:      "space-indented code",
			input:     "    func main() {}\n",
			wantLines: 1,
		},
		{
			name: "mixed tab and space continuation",
			input: `	first
	second
    third
`,
			wantLines: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewCodeBlockParser()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()

		node, _ := parser.Open(ast.NewDocument(), reader, pc)
		if node == nil {
			t.Fatal("Open() returned nil")
		}
		reader.AdvanceLine()

		// Continue reading all lines
		for reader.Peek() != text.EOF {
			state := parser.Continue(node, reader, pc)
			if state == Close {
				break
			}
			reader.AdvanceLine()
		}

			if node.Lines().Len() != tc.wantLines {
				t.Errorf("Got %d lines, want %d", node.Lines().Len(), tc.wantLines)
			}
		})
	}
}
