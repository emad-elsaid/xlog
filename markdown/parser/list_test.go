package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestParseListItem(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantTyp listItemType
		wantRet [6]int
	}{
		{
			name:    "bullet list with dash",
			input:   []byte("- item"),
			wantTyp: bulletList,
			wantRet: [6]int{0, 0, 0, 1, 1, 6},
		},
		{
			name:    "bullet list with asterisk",
			input:   []byte("* item"),
			wantTyp: bulletList,
			wantRet: [6]int{0, 0, 0, 1, 1, 6},
		},
		{
			name:    "bullet list with plus",
			input:   []byte("+ item"),
			wantTyp: bulletList,
			wantRet: [6]int{0, 0, 0, 1, 1, 6},
		},
		{
			name:    "ordered list with dot",
			input:   []byte("1. item"),
			wantTyp: orderedList,
			wantRet: [6]int{0, 0, 0, 2, 2, 7},
		},
		{
			name:    "ordered list with paren",
			input:   []byte("1) item"),
			wantTyp: orderedList,
			wantRet: [6]int{0, 0, 0, 2, 2, 7},
		},
		{
			name:    "ordered list with large number",
			input:   []byte("123456789. item"),
			wantTyp: orderedList,
			wantRet: [6]int{0, 0, 0, 10, 10, 15},
		},
		{
			name:    "ordered list number too large",
			input:   []byte("1234567890. item"),
			wantTyp: notList,
			wantRet: [6]int{},
		},
		{
			name:    "indented bullet list 1 space",
			input:   []byte(" - item"),
			wantTyp: bulletList,
			wantRet: [6]int{0, 1, 1, 2, 2, 7},
		},
		{
			name:    "indented bullet list 2 spaces",
			input:   []byte("  - item"),
			wantTyp: bulletList,
			wantRet: [6]int{0, 2, 2, 3, 3, 8},
		},
		{
			name:    "indented bullet list 3 spaces",
			input:   []byte("   - item"),
			wantTyp: bulletList,
			wantRet: [6]int{0, 3, 3, 4, 4, 9},
		},
		{
			name:    "indented too much (4 spaces)",
			input:   []byte("    - item"),
			wantTyp: notList,
			wantRet: [6]int{},
		},
		{
			name:    "bullet list with tab",
			input:   []byte("\t- item"),
			wantTyp: notList,
			wantRet: [6]int{},
		},
		{
			name:    "empty list item",
			input:   []byte("-"),
			wantTyp: bulletList,
			wantRet: [6]int{0, 0, 0, 1, -1, -1},
		},
		{
			name:    "empty ordered list item",
			input:   []byte("1."),
			wantTyp: orderedList,
			wantRet: [6]int{0, 0, 0, 2, -1, -1},
		},
		{
			name:    "bullet list with newline",
			input:   []byte("- item\n"),
			wantTyp: bulletList,
			wantRet: [6]int{0, 0, 0, 1, 1, 6},
		},
		{
			name:    "not a list - no space after marker",
			input:   []byte("-item"),
			wantTyp: notList,
			wantRet: [6]int{},
		},
		{
			name:    "not a list - plain text",
			input:   []byte("plain text"),
			wantTyp: notList,
			wantRet: [6]int{},
		},
		{
			name:    "ordered list starting with zero",
			input:   []byte("0. item"),
			wantTyp: orderedList,
			wantRet: [6]int{0, 0, 0, 2, 2, 7},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotRet, gotTyp := parseListItem(tc.input)
			if gotTyp != tc.wantTyp {
				t.Errorf("parseListItem() type = %v, want %v", gotTyp, tc.wantTyp)
			}
			if tc.wantTyp != notList && gotRet != tc.wantRet {
				t.Errorf("parseListItem() ret = %v, want %v", gotRet, tc.wantRet)
			}
		})
	}
}

func TestMatchesListItem(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		strict  bool
		wantTyp listItemType
	}{
		{
			name:    "strict mode - valid bullet",
			input:   []byte("- item"),
			strict:  true,
			wantTyp: bulletList,
		},
		{
			name:    "strict mode - 4 spaces indented",
			input:   []byte("    - item"),
			strict:  true,
			wantTyp: notList,
		},
		{
			name:    "non-strict mode - 4 spaces indented",
			input:   []byte("    - item"),
			strict:  false,
			wantTyp: notList,
		},
		{
			name:    "strict mode - 3 spaces ok",
			input:   []byte("   - item"),
			strict:  true,
			wantTyp: bulletList,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, gotTyp := matchesListItem(tc.input, tc.strict)
			if gotTyp != tc.wantTyp {
				t.Errorf("matchesListItem() = %v, want %v", gotTyp, tc.wantTyp)
			}
		})
	}
}

func TestCalcListOffset(t *testing.T) {
	tests := []struct {
		name       string
		source     []byte
		match      [6]int
		wantOffset int
	}{
		{
			name:       "blank line after marker",
			source:     []byte("-   "),
			match:      [6]int{0, 0, 0, 1, 1, 4},
			wantOffset: 1,
		},
		{
			name:       "normal item with content",
			source:     []byte("- item"),
			match:      [6]int{0, 0, 0, 1, 1, 6},
			wantOffset: 1,
		},
		{
			name:       "marker followed by large indent (codeblock)",
			source:     []byte("-     code"),
			match:      [6]int{0, 0, 0, 1, 1, 10},
			wantOffset: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotOffset := calcListOffset(tc.source, tc.match)
			if gotOffset != tc.wantOffset {
				t.Errorf("calcListOffset() = %v, want %v", gotOffset, tc.wantOffset)
			}
		})
	}
}

func TestLastOffset(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() ast.Node
		wantOffset int
	}{
		{
			name: "list with one item",
			setup: func() ast.Node {
				list := ast.NewList('-')
				item := ast.NewListItem(2)
				list.AppendChild(list, item)
				return list
			},
			wantOffset: 2,
		},
		{
			name: "empty list",
			setup: func() ast.Node {
				return ast.NewList('-')
			},
			wantOffset: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node := tc.setup()
			gotOffset := lastOffset(node)
			if gotOffset != tc.wantOffset {
				t.Errorf("lastOffset() = %v, want %v", gotOffset, tc.wantOffset)
			}
		})
	}
}

func TestListParser_Trigger(t *testing.T) {
	parser := NewListParser()
	trigger := parser.Trigger()
	expected := []byte{'-', '+', '*', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	if len(trigger) != len(expected) {
		t.Errorf("Trigger() length = %v, want %v", len(trigger), len(expected))
		return
	}
	for i, b := range expected {
		if trigger[i] != b {
			t.Errorf("Trigger()[%d] = %v, want %v", i, trigger[i], b)
		}
	}
}

func TestListParser_Open(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantState State
		wantNil   bool
		checkNode func(*testing.T, ast.Node)
	}{
		{
			name:      "simple bullet list",
			input:     "- item\n",
			wantState: HasChildren,
			wantNil:   false,
			checkNode: func(t *testing.T, node ast.Node) {
				list, ok := node.(*ast.List)
				if !ok {
					t.Errorf("expected *ast.List, got %T", node)
					return
				}
				if list.Marker != '-' {
					t.Errorf("Marker = %c, want -", list.Marker)
				}
			},
		},
		{
			name:      "simple ordered list",
			input:     "1. item\n",
			wantState: HasChildren,
			wantNil:   false,
			checkNode: func(t *testing.T, node ast.Node) {
				list, ok := node.(*ast.List)
				if !ok {
					t.Errorf("expected *ast.List, got %T", node)
					return
				}
				if list.Start != 1 {
					t.Errorf("Start = %d, want 1", list.Start)
				}
			},
		},
		{
			name:      "ordered list with asterisk marker",
			input:     "* item\n",
			wantState: HasChildren,
			wantNil:   false,
			checkNode: func(t *testing.T, node ast.Node) {
				list, ok := node.(*ast.List)
				if !ok {
					t.Errorf("expected *ast.List, got %T", node)
					return
				}
				if list.Marker != '*' {
					t.Errorf("Marker = %c, want *", list.Marker)
				}
			},
		},
		{
			name:      "not a list",
			input:     "plain text\n",
			wantState: NoChildren,
			wantNil:   true,
		},
		{
			name:      "indented too much",
			input:     "    - item\n",
			wantState: NoChildren,
			wantNil:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewListParser()
			parent := ast.NewDocument()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()

			node, state := parser.Open(parent, reader, pc)

			if tc.wantNil && node != nil {
				t.Errorf("Open() node = %v, want nil", node)
			}
			if !tc.wantNil && node == nil {
				t.Errorf("Open() node = nil, want non-nil")
			}
			if state != tc.wantState {
				t.Errorf("Open() state = %v, want %v", state, tc.wantState)
			}
			if !tc.wantNil && tc.checkNode != nil {
				tc.checkNode(t, node)
			}
		})
	}
}

func TestListParser_Continue(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() (ast.Node, text.Reader, Context)
		wantState State
	}{
		{
			name: "blank line in list",
			setup: func() (ast.Node, text.Reader, Context) {
				list := ast.NewList('-')
				item := ast.NewListItem(2)
				list.AppendChild(list, item)
				reader := text.NewReader([]byte("\n"))
				pc := NewContext()
				return list, reader, pc
			},
			wantState: Continue | HasChildren,
		},
		{
			name: "continued list item",
			setup: func() (ast.Node, text.Reader, Context) {
				list := ast.NewList('-')
				item := ast.NewListItem(2)
				list.AppendChild(list, item)
				reader := text.NewReader([]byte("  continued text\n"))
				pc := NewContext()
				return list, reader, pc
			},
			wantState: Continue | HasChildren,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewListParser()
			node, reader, pc := tc.setup()
			state := parser.Continue(node, reader, pc)
			if state != tc.wantState {
				t.Errorf("Continue() = %v, want %v", state, tc.wantState)
			}
		})
	}
}

func TestListParser_Close(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() ast.Node
		wantTight bool
	}{
		{
			name: "tight list",
			setup: func() ast.Node {
				list := ast.NewList('-')
				item1 := ast.NewListItem(2)
				para1 := ast.NewParagraph()
				item1.AppendChild(item1, para1)
				list.AppendChild(list, item1)

				item2 := ast.NewListItem(2)
				para2 := ast.NewParagraph()
				item2.AppendChild(item2, para2)
				list.AppendChild(list, item2)
				return list
			},
			wantTight: true,
		},
		{
			name: "loose list with blank lines",
			setup: func() ast.Node {
				list := ast.NewList('-')
				item1 := ast.NewListItem(2)
				para1 := ast.NewParagraph()
				item1.AppendChild(item1, para1)
				list.AppendChild(list, item1)

				item2 := ast.NewListItem(2)
				item2.SetBlankPreviousLines(true)
				para2 := ast.NewParagraph()
				item2.AppendChild(item2, para2)
				list.AppendChild(list, item2)
				return list
			},
			wantTight: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewListParser()
			node := tc.setup()
			reader := text.NewReader([]byte(""))
			pc := NewContext()

			parser.Close(node, reader, pc)

			list := node.(*ast.List)
			if list.IsTight != tc.wantTight {
				t.Errorf("Close() IsTight = %v, want %v", list.IsTight, tc.wantTight)
			}
		})
	}
}

func TestListParser_CanInterruptParagraph(t *testing.T) {
	parser := NewListParser()
	if !parser.CanInterruptParagraph() {
		t.Error("CanInterruptParagraph() = false, want true")
	}
}

func TestListParser_CanAcceptIndentedLine(t *testing.T) {
	parser := NewListParser()
	if parser.CanAcceptIndentedLine() {
		t.Error("CanAcceptIndentedLine() = true, want false")
	}
}
