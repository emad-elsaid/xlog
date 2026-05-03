package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestFencedCodeBlockParser_Trigger(t *testing.T) {
	parser := NewFencedCodeBlockParser()
	trigger := parser.Trigger()
	expected := []byte{'~', '`'}

	if len(trigger) != len(expected) {
		t.Fatalf("Trigger() length = %d, want %d", len(trigger), len(expected))
	}

	for i, c := range expected {
		if trigger[i] != c {
			t.Errorf("Trigger()[%d] = %c, want %c", i, trigger[i], c)
		}
	}
}

func TestFencedCodeBlockParser_Open(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantState   State
		wantNil     bool
		blockOffset int
		wantInfo    bool
	}{
		{
			name:        "simple backtick fence",
			input:       "```\n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantInfo:    false,
		},
		{
			name:        "backtick fence with language",
			input:       "```go\n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantInfo:    true,
		},
		{
			name:        "backtick fence with language and spaces",
			input:       "```  go  \n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantInfo:    true,
		},
		{
			name:        "tilde fence",
			input:       "~~~\n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantInfo:    false,
		},
		{
			name:        "tilde fence with language",
			input:       "~~~python\n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantInfo:    true,
		},
		{
			name:        "four backticks",
			input:       "````\n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantInfo:    false,
		},
		{
			name:        "five tildes",
			input:       "~~~~~\n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantInfo:    false,
		},
		{
			name:        "fence with complex info string",
			input:       "```javascript {highlight: true}\n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantInfo:    true,
		},
		{
			name:        "only two backticks - should fail",
			input:       "``\n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantNil:     true,
		},
		{
			name:        "only two tildes - should fail",
			input:       "~~\n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantNil:     true,
		},
		{
			name:        "backtick fence with backtick in info - should fail",
			input:       "```go`test\n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantNil:     true,
		},
		{
			name:        "indented fence (1 space)",
			input:       " ```\n",
			wantState:   NoChildren,
			blockOffset: 1,
			wantInfo:    false,
		},
		{
			name:        "indented fence (2 spaces)",
			input:       "  ```\n",
			wantState:   NoChildren,
			blockOffset: 2,
			wantInfo:    false,
		},
		{
			name:        "indented fence (3 spaces)",
			input:       "   ```\n",
			wantState:   NoChildren,
			blockOffset: 3,
			wantInfo:    false,
		},
		{
			name:        "empty line - should fail",
			input:       "\n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantNil:     true,
		},
		{
			name:        "regular text - should fail",
			input:       "hello world\n",
			wantState:   NoChildren,
			blockOffset: 0,
			wantNil:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewFencedCodeBlockParser()
			reader := text.NewReader([]byte(tc.input))
			ctx := NewContext()
			ctx.SetBlockOffset(tc.blockOffset)

			node, state := parser.Open(ast.NewDocument(), reader, ctx)

			if tc.wantNil {
				if node != nil {
					t.Errorf("Open() node = %v, want nil", node)
				}
				return
			}

			if node == nil {
				t.Fatal("Open() node = nil, want non-nil")
			}

			if state != tc.wantState {
				t.Errorf("Open() state = %v, want %v", state, tc.wantState)
			}

			fcb, ok := node.(*ast.FencedCodeBlock)
			if !ok {
				t.Fatalf("Open() node type = %T, want *ast.FencedCodeBlock", node)
			}

			if tc.wantInfo {
				if fcb.Info == nil {
					t.Error("Open() node.Info = nil, want non-nil")
				}
			} else {
				if fcb.Info != nil {
					t.Errorf("Open() node.Info = %v, want nil", fcb.Info)
				}
			}
		})
	}
}

func TestFencedCodeBlockParser_Continue(t *testing.T) {
	tests := []struct {
		name      string
		setup     string
		nextLine  string
		wantState State
		wantLines int
	}{
		{
			name:      "regular content line",
			setup:     "```\n",
			nextLine:  "code here\n",
			wantState: Continue | NoChildren,
			wantLines: 1,
		},
		{
			name:      "empty line in code block",
			setup:     "```\n",
			nextLine:  "\n",
			wantState: Continue | NoChildren,
			wantLines: 1,
		},
		{
			name:      "closing fence - same length",
			setup:     "```\n",
			nextLine:  "```\n",
			wantState: Close,
			wantLines: 0,
		},
		{
			name:      "closing fence - longer",
			setup:     "```\n",
			nextLine:  "````\n",
			wantState: Close,
			wantLines: 0,
		},
		{
			name:      "closing fence - tilde",
			setup:     "~~~\n",
			nextLine:  "~~~\n",
			wantState: Close,
			wantLines: 0,
		},
		{
			name:      "closing fence - shorter should not close",
			setup:     "````\n",
			nextLine:  "```\n",
			wantState: Continue | NoChildren,
			wantLines: 1,
		},
		{
			name:      "mismatched fence chars - backtick vs tilde",
			setup:     "```\n",
			nextLine:  "~~~\n",
			wantState: Continue | NoChildren,
			wantLines: 1,
		},
		{
			name:      "closing fence with trailing spaces",
			setup:     "```\n",
			nextLine:  "```   \n",
			wantState: Close,
			wantLines: 0,
		},
		{
			name:      "closing fence with trailing content - should not close",
			setup:     "```\n",
			nextLine:  "``` text\n",
			wantState: Continue | NoChildren,
			wantLines: 1,
		},
		{
			name:      "indented content",
			setup:     "```\n",
			nextLine:  "  code\n",
			wantState: Continue | NoChildren,
			wantLines: 1,
		},
		{
			name:      "tab-indented content",
			setup:     "```\n",
			nextLine:  "\tcode\n",
			wantState: Continue | NoChildren,
			wantLines: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewFencedCodeBlockParser()

			// Setup: Open the fence
			reader := text.NewReader([]byte(tc.setup))
			ctx := NewContext()
			ctx.SetBlockOffset(0)
			node, _ := parser.Open(ast.NewDocument(), reader, ctx)
			if node == nil {
				t.Fatal("Setup failed: Open() returned nil")
			}

			// Test: Continue with next line
			reader = text.NewReader([]byte(tc.nextLine))
			state := parser.Continue(node, reader, ctx)

			if state != tc.wantState {
				t.Errorf("Continue() state = %v, want %v", state, tc.wantState)
			}

			if node.Lines().Len() != tc.wantLines {
				t.Errorf("Continue() lines count = %d, want %d", node.Lines().Len(), tc.wantLines)
			}
		})
	}
}

func TestFencedCodeBlockParser_Close(t *testing.T) {
	parser := NewFencedCodeBlockParser()

	// Setup: Open a fence
	input := "```\n"
	reader := text.NewReader([]byte(input))
	ctx := NewContext()
	ctx.SetBlockOffset(0)
	node, _ := parser.Open(ast.NewDocument(), reader, ctx)

	if node == nil {
		t.Fatal("Setup failed: Open() returned nil")
	}

	// Verify context has fence data before close
	fdata := ctx.Get(fencedCodeBlockInfoKey)
	if fdata == nil {
		t.Fatal("Context missing fence data before Close()")
	}

	// Test: Close the block
	parser.Close(node, reader, ctx)

	// Verify context cleared fence data
	fdata = ctx.Get(fencedCodeBlockInfoKey)
	if fdata != nil {
		t.Error("Close() should clear fence data from context")
	}
}

func TestFencedCodeBlockParser_CanInterruptParagraph(t *testing.T) {
	parser := NewFencedCodeBlockParser()
	if !parser.CanInterruptParagraph() {
		t.Error("CanInterruptParagraph() = false, want true")
	}
}

func TestFencedCodeBlockParser_CanAcceptIndentedLine(t *testing.T) {
	parser := NewFencedCodeBlockParser()
	if parser.CanAcceptIndentedLine() {
		t.Error("CanAcceptIndentedLine() = true, want false")
	}
}

func TestFencedCodeBlockParser_MultilineCode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLines int
	}{
		{
			name:      "simple multiline code",
			input:     "```\nline1\nline2\nline3\n```\n",
			wantLines: 3,
		},
		{
			name:      "code with empty lines",
			input:     "```\nline1\n\nline2\n```\n",
			wantLines: 3,
		},
		{
			name:      "code with indentation preserved",
			input:     "```\n  func main() {\n    println()\n  }\n```\n",
			wantLines: 3,
		},
		{
			name:      "empty code block",
			input:     "```\n```\n",
			wantLines: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewFencedCodeBlockParser()
			lines := []byte(tc.input)
			reader := text.NewReader(lines)
			ctx := NewContext()
			ctx.SetBlockOffset(0)

			// Open the fence
			node, state := parser.Open(ast.NewDocument(), reader, ctx)
			if node == nil {
				t.Fatal("Open() returned nil")
			}
			if state != NoChildren {
				t.Fatalf("Open() state = %v, want NoChildren", state)
			}

			// Process all lines until close
			_, seg := reader.PeekLine()
			reader.Advance(seg.Stop - seg.Start)
			for {
				line, seg := reader.PeekLine()
				if len(line) == 0 {
					break
				}

				state = parser.Continue(node, reader, ctx)
				if state == Close {
					break
				}
				reader.Advance(seg.Stop - seg.Start)
			}

			if node.Lines().Len() != tc.wantLines {
				t.Errorf("Lines count = %d, want %d", node.Lines().Len(), tc.wantLines)
			}
		})
	}
}

func TestFencedCodeBlockParser_InfoString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantInfo string
	}{
		{
			name:     "language only",
			input:    "```go\n",
			wantInfo: "go",
		},
		{
			name:     "language with spaces",
			input:    "```  python  \n",
			wantInfo: "python",
		},
		{
			name:     "complex info string",
			input:    "```javascript {highlight: [1,2,3]}\n",
			wantInfo: "javascript {highlight: [1,2,3]}",
		},
		{
			name:     "no info string",
			input:    "```\n",
			wantInfo: "",
		},
		{
			name:     "tilde with language",
			input:    "~~~rust\n",
			wantInfo: "rust",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewFencedCodeBlockParser()
			reader := text.NewReader([]byte(tc.input))
			ctx := NewContext()
			ctx.SetBlockOffset(0)

			node, _ := parser.Open(ast.NewDocument(), reader, ctx)
			if node == nil {
				t.Fatal("Open() returned nil")
			}

			fcb, ok := node.(*ast.FencedCodeBlock)
			if !ok {
				t.Fatalf("node type = %T, want *ast.FencedCodeBlock", node)
			}

			if tc.wantInfo == "" {
				if fcb.Info != nil {
					t.Errorf("Info = %v, want nil", fcb.Info)
				}
			} else {
				if fcb.Info == nil {
					t.Fatal("Info = nil, want non-nil")
				}

				// Extract the actual info text from the segment
				infoSegment := fcb.Info.Segment
				actualInfo := string([]byte(tc.input)[infoSegment.Start:infoSegment.Stop])

				if actualInfo != tc.wantInfo {
					t.Errorf("Info text = %q, want %q", actualInfo, tc.wantInfo)
				}
			}
		})
	}
}
