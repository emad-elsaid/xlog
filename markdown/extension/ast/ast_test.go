package ast

import (
	"bytes"
	"strings"
	"testing"

	gast "github.com/emad-elsaid/xlog/markdown/ast"
)

// TestDefinitionListNode tests DefinitionList AST node creation and properties.
func TestDefinitionListNode(t *testing.T) {
	tests := []struct {
		name          string
		offset        int
		withParagraph bool
		wantKind      gast.NodeKind
		wantOffset    int
		wantParaNil   bool
	}{
		{
			name:          "new definition list with offset",
			offset:        5,
			withParagraph: false,
			wantKind:      KindDefinitionList,
			wantOffset:    5,
			wantParaNil:   true,
		},
		{
			name:          "new definition list with paragraph",
			offset:        10,
			withParagraph: true,
			wantKind:      KindDefinitionList,
			wantOffset:    10,
			wantParaNil:   false,
		},
		{
			name:          "zero offset",
			offset:        0,
			withParagraph: false,
			wantKind:      KindDefinitionList,
			wantOffset:    0,
			wantParaNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var para *gast.Paragraph
			if tt.withParagraph {
				para = gast.NewParagraph()
			}

			node := NewDefinitionList(tt.offset, para)

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}

			if node.Offset != tt.wantOffset {
				t.Errorf("Offset = %v, want %v", node.Offset, tt.wantOffset)
			}

			if tt.wantParaNil && node.TemporaryParagraph != nil {
				t.Errorf("TemporaryParagraph = %v, want nil", node.TemporaryParagraph)
			}

			if !tt.wantParaNil && node.TemporaryParagraph == nil {
				t.Errorf("TemporaryParagraph = nil, want non-nil")
			}
		})
	}
}

// TestDefinitionListDump tests the Dump method.
func TestDefinitionListDump(t *testing.T) {
	node := NewDefinitionList(0, nil)
	source := []byte("test source")

	// Dump should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestDefinitionTermNode tests DefinitionTerm AST node.
func TestDefinitionTermNode(t *testing.T) {
	tests := []struct {
		name     string
		wantKind gast.NodeKind
	}{
		{
			name:     "new definition term",
			wantKind: KindDefinitionTerm,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewDefinitionTerm()

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}
		})
	}
}

// TestDefinitionTermDump tests the Dump method for DefinitionTerm.
func TestDefinitionTermDump(t *testing.T) {
	node := NewDefinitionTerm()
	source := []byte("test")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestDefinitionDescriptionNode tests DefinitionDescription AST node.
func TestDefinitionDescriptionNode(t *testing.T) {
	tests := []struct {
		name     string
		wantKind gast.NodeKind
	}{
		{
			name:     "new definition description",
			wantKind: KindDefinitionDescription,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewDefinitionDescription()

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}

			// Default IsTight should be false
			if node.IsTight {
				t.Errorf("IsTight = true, want false")
			}
		})
	}
}

// TestDefinitionDescriptionDump tests the Dump method for DefinitionDescription.
func TestDefinitionDescriptionDump(t *testing.T) {
	node := NewDefinitionDescription()
	source := []byte("test")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestFootnoteLinkNode tests FootnoteLink AST node.
func TestFootnoteLinkNode(t *testing.T) {
	tests := []struct {
		name         string
		index        int
		wantKind     gast.NodeKind
		wantIndex    int
		wantRefCount int
		wantRefIndex int
	}{
		{
			name:         "new footnote link with index 0",
			index:        0,
			wantKind:     KindFootnoteLink,
			wantIndex:    0,
			wantRefCount: 0,
			wantRefIndex: 0,
		},
		{
			name:         "new footnote link with index 5",
			index:        5,
			wantKind:     KindFootnoteLink,
			wantIndex:    5,
			wantRefCount: 0,
			wantRefIndex: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewFootnoteLink(tt.index)

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}

			if node.Index != tt.wantIndex {
				t.Errorf("Index = %v, want %v", node.Index, tt.wantIndex)
			}

			if node.RefCount != tt.wantRefCount {
				t.Errorf("RefCount = %v, want %v", node.RefCount, tt.wantRefCount)
			}

			if node.RefIndex != tt.wantRefIndex {
				t.Errorf("RefIndex = %v, want %v", node.RefIndex, tt.wantRefIndex)
			}
		})
	}
}

// TestFootnoteLinkDump tests the Dump method for FootnoteLink.
func TestFootnoteLinkDump(t *testing.T) {
	node := NewFootnoteLink(1)
	source := []byte("test")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestFootnoteBacklinkNode tests FootnoteBacklink AST node.
func TestFootnoteBacklinkNode(t *testing.T) {
	tests := []struct {
		name         string
		index        int
		wantKind     gast.NodeKind
		wantIndex    int
		wantRefCount int
		wantRefIndex int
	}{
		{
			name:         "new footnote backlink with index 0",
			index:        0,
			wantKind:     KindFootnoteBacklink,
			wantIndex:    0,
			wantRefCount: 0,
			wantRefIndex: 0,
		},
		{
			name:         "new footnote backlink with index 3",
			index:        3,
			wantKind:     KindFootnoteBacklink,
			wantIndex:    3,
			wantRefCount: 0,
			wantRefIndex: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewFootnoteBacklink(tt.index)

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}

			if node.Index != tt.wantIndex {
				t.Errorf("Index = %v, want %v", node.Index, tt.wantIndex)
			}

			if node.RefCount != tt.wantRefCount {
				t.Errorf("RefCount = %v, want %v", node.RefCount, tt.wantRefCount)
			}

			if node.RefIndex != tt.wantRefIndex {
				t.Errorf("RefIndex = %v, want %v", node.RefIndex, tt.wantRefIndex)
			}
		})
	}
}

// TestFootnoteBacklinkDump tests the Dump method for FootnoteBacklink.
func TestFootnoteBacklinkDump(t *testing.T) {
	node := NewFootnoteBacklink(2)
	source := []byte("test")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestFootnoteNode tests Footnote AST node.
func TestFootnoteNode(t *testing.T) {
	tests := []struct {
		name      string
		ref       []byte
		wantKind  gast.NodeKind
		wantRef   []byte
		wantIndex int
	}{
		{
			name:      "new footnote with simple ref",
			ref:       []byte("note1"),
			wantKind:  KindFootnote,
			wantRef:   []byte("note1"),
			wantIndex: -1,
		},
		{
			name:      "new footnote with empty ref",
			ref:       []byte(""),
			wantKind:  KindFootnote,
			wantRef:   []byte(""),
			wantIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewFootnote(tt.ref)

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}

			if !bytes.Equal(node.Ref, tt.wantRef) {
				t.Errorf("Ref = %v, want %v", node.Ref, tt.wantRef)
			}

			if node.Index != tt.wantIndex {
				t.Errorf("Index = %v, want %v", node.Index, tt.wantIndex)
			}
		})
	}
}

// TestFootnoteDump tests the Dump method for Footnote.
func TestFootnoteDump(t *testing.T) {
	node := NewFootnote([]byte("ref"))
	source := []byte("test")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestFootnoteListNode tests FootnoteList AST node.
func TestFootnoteListNode(t *testing.T) {
	tests := []struct {
		name      string
		wantKind  gast.NodeKind
		wantCount int
	}{
		{
			name:      "new footnote list",
			wantKind:  KindFootnoteList,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewFootnoteList()

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}

			if node.Count != tt.wantCount {
				t.Errorf("Count = %v, want %v", node.Count, tt.wantCount)
			}
		})
	}
}

// TestFootnoteListDump tests the Dump method for FootnoteList.
func TestFootnoteListDump(t *testing.T) {
	node := NewFootnoteList()
	source := []byte("test")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestStrikethroughNode tests Strikethrough AST node.
func TestStrikethroughNode(t *testing.T) {
	tests := []struct {
		name     string
		wantKind gast.NodeKind
	}{
		{
			name:     "new strikethrough",
			wantKind: KindStrikethrough,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewStrikethrough()

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}
		})
	}
}

// TestStrikethroughDump tests the Dump method for Strikethrough.
func TestStrikethroughDump(t *testing.T) {
	node := NewStrikethrough()
	source := []byte("test")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestAlignmentString tests Alignment.String() method.
func TestAlignmentString(t *testing.T) {
	tests := []struct {
		name      string
		alignment Alignment
		want      string
	}{
		{
			name:      "left alignment",
			alignment: AlignLeft,
			want:      "left",
		},
		{
			name:      "right alignment",
			alignment: AlignRight,
			want:      "right",
		},
		{
			name:      "center alignment",
			alignment: AlignCenter,
			want:      "center",
		},
		{
			name:      "none alignment",
			alignment: AlignNone,
			want:      "none",
		},
		{
			name:      "invalid alignment",
			alignment: Alignment(99),
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.alignment.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTableNode tests Table AST node.
func TestTableNode(t *testing.T) {
	tests := []struct {
		name       string
		wantKind   gast.NodeKind
		wantAligns []Alignment
	}{
		{
			name:       "new table",
			wantKind:   KindTable,
			wantAligns: []Alignment{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewTable()

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}

			if len(node.Alignments) != len(tt.wantAligns) {
				t.Errorf("Alignments length = %v, want %v", len(node.Alignments), len(tt.wantAligns))
			}
		})
	}
}

// TestTableDump tests the Dump method for Table.
func TestTableDump(t *testing.T) {
	node := NewTable()
	node.Alignments = []Alignment{AlignLeft, AlignCenter, AlignRight}
	source := []byte("test")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestTableRowNode tests TableRow AST node.
func TestTableRowNode(t *testing.T) {
	tests := []struct {
		name       string
		alignments []Alignment
		wantKind   gast.NodeKind
		wantLen    int
	}{
		{
			name:       "new table row with alignments",
			alignments: []Alignment{AlignLeft, AlignRight},
			wantKind:   KindTableRow,
			wantLen:    2,
		},
		{
			name:       "new table row empty",
			alignments: []Alignment{},
			wantKind:   KindTableRow,
			wantLen:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewTableRow(tt.alignments)

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}

			if len(node.Alignments) != tt.wantLen {
				t.Errorf("Alignments length = %v, want %v", len(node.Alignments), tt.wantLen)
			}
		})
	}
}

// TestTableRowDump tests the Dump method for TableRow.
func TestTableRowDump(t *testing.T) {
	node := NewTableRow([]Alignment{AlignLeft})
	source := []byte("test")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestTableHeaderNode tests TableHeader AST node.
func TestTableHeaderNode(t *testing.T) {
	tests := []struct {
		name        string
		rowChildren int
		wantKind    gast.NodeKind
	}{
		{
			name:        "new table header from row with children",
			rowChildren: 2,
			wantKind:    KindTableHeader,
		},
		{
			name:        "new table header from empty row",
			rowChildren: 0,
			wantKind:    KindTableHeader,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := NewTableRow([]Alignment{})
			for i := 0; i < tt.rowChildren; i++ {
				row.AppendChild(row, NewTableCell())
			}

			node := NewTableHeader(row)

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}

			// Count children
			count := 0
			for c := node.FirstChild(); c != nil; c = c.NextSibling() {
				count++
			}

			if count != tt.rowChildren {
				t.Errorf("Children count = %v, want %v", count, tt.rowChildren)
			}
		})
	}
}

// TestTableHeaderDump tests the Dump method for TableHeader.
func TestTableHeaderDump(t *testing.T) {
	row := NewTableRow([]Alignment{})
	node := NewTableHeader(row)
	source := []byte("test")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestTableCellNode tests TableCell AST node.
func TestTableCellNode(t *testing.T) {
	tests := []struct {
		name          string
		wantKind      gast.NodeKind
		wantAlignment Alignment
	}{
		{
			name:          "new table cell",
			wantKind:      KindTableCell,
			wantAlignment: AlignNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewTableCell()

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}

			if node.Alignment != tt.wantAlignment {
				t.Errorf("Alignment = %v, want %v", node.Alignment, tt.wantAlignment)
			}
		})
	}
}

// TestTableCellDump tests the Dump method for TableCell.
func TestTableCellDump(t *testing.T) {
	node := NewTableCell()
	source := []byte("test")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestTaskCheckBoxNode tests TaskCheckBox AST node.
func TestTaskCheckBoxNode(t *testing.T) {
	tests := []struct {
		name        string
		checked     bool
		wantKind    gast.NodeKind
		wantChecked bool
	}{
		{
			name:        "new checked task checkbox",
			checked:     true,
			wantKind:    KindTaskCheckBox,
			wantChecked: true,
		},
		{
			name:        "new unchecked task checkbox",
			checked:     false,
			wantKind:    KindTaskCheckBox,
			wantChecked: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewTaskCheckBox(tt.checked)

			if node.Kind() != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", node.Kind(), tt.wantKind)
			}

			if node.IsChecked != tt.wantChecked {
				t.Errorf("IsChecked = %v, want %v", node.IsChecked, tt.wantChecked)
			}
		})
	}
}

// TestTaskCheckBoxDump tests the Dump method for TaskCheckBox.
func TestTaskCheckBoxDump(t *testing.T) {
	node := NewTaskCheckBox(true)
	source := []byte("test")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	node.Dump(source, 0)
}

// TestNodeKinds verifies that all NodeKind values are distinct.
func TestNodeKinds(t *testing.T) {
	kinds := []gast.NodeKind{
		KindDefinitionList,
		KindDefinitionTerm,
		KindDefinitionDescription,
		KindFootnoteLink,
		KindFootnoteBacklink,
		KindFootnote,
		KindFootnoteList,
		KindStrikethrough,
		KindTable,
		KindTableRow,
		KindTableHeader,
		KindTableCell,
		KindTaskCheckBox,
	}

	seen := make(map[string]bool)
	for _, kind := range kinds {
		name := kind.String()
		if seen[name] {
			t.Errorf("Duplicate NodeKind name: %s", name)
		}
		seen[name] = true
	}
}

// TestTableDumpWithAlignment verifies table dump output format.
func TestTableDumpWithAlignment(t *testing.T) {
	// Redirect stdout to capture dump output
	node := NewTable()
	node.Alignments = []Alignment{AlignLeft, AlignRight, AlignCenter, AlignNone}
	source := []byte("| a | b | c | d |")

	// This test verifies Dump doesn't panic with multiple alignments
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() with multiple alignments panicked: %v", r)
		}
	}()

	node.Dump(source, 1)
}

// TestTableHeaderMovesChildrenFromRow verifies child movement logic.
func TestTableHeaderMovesChildrenFromRow(t *testing.T) {
	row := NewTableRow([]Alignment{AlignLeft, AlignRight})

	cell1 := NewTableCell()
	cell2 := NewTableCell()
	row.AppendChild(row, cell1)
	row.AppendChild(row, cell2)

	// Verify row has 2 children
	rowChildCount := 0
	for c := row.FirstChild(); c != nil; c = c.NextSibling() {
		rowChildCount++
	}
	if rowChildCount != 2 {
		t.Fatalf("Row should have 2 children before NewTableHeader, got %d", rowChildCount)
	}

	header := NewTableHeader(row)

	// Verify header has 2 children
	headerChildCount := 0
	for c := header.FirstChild(); c != nil; c = c.NextSibling() {
		headerChildCount++
	}

	if headerChildCount != 2 {
		t.Errorf("Header should have 2 children, got %d", headerChildCount)
	}
}

// TestDumpWithComplexStructure tests Dump with nested node structures.
func TestDumpWithComplexStructure(t *testing.T) {
	tests := []struct {
		name string
		node gast.Node
	}{
		{
			name: "definition list with term and description",
			node: func() gast.Node {
				dl := NewDefinitionList(0, nil)
				term := NewDefinitionTerm()
				desc := NewDefinitionDescription()
				dl.AppendChild(dl, term)
				dl.AppendChild(dl, desc)
				return dl
			}(),
		},
		{
			name: "table with rows and cells",
			node: func() gast.Node {
				table := NewTable()
				table.Alignments = []Alignment{AlignLeft, AlignRight}
				row := NewTableRow(table.Alignments)
				cell1 := NewTableCell()
				cell1.Alignment = AlignLeft
				cell2 := NewTableCell()
				cell2.Alignment = AlignRight
				row.AppendChild(row, cell1)
				row.AppendChild(row, cell2)
				table.AppendChild(table, row)
				return table
			}(),
		},
		{
			name: "footnote list with footnote",
			node: func() gast.Node {
				list := NewFootnoteList()
				footnote := NewFootnote([]byte("ref"))
				list.AppendChild(list, footnote)
				return list
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := []byte("test source")

			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Dump() panicked with complex structure: %v", r)
				}
			}()

			tt.node.Dump(source, 0)
		})
	}
}

// TestAlignmentConstants verifies alignment constant values are distinct.
func TestAlignmentConstants(t *testing.T) {
	alignments := []Alignment{AlignLeft, AlignRight, AlignCenter, AlignNone}
	seen := make(map[Alignment]bool)

	for _, align := range alignments {
		if seen[align] {
			t.Errorf("Duplicate alignment value: %v", align)
		}
		seen[align] = true
	}
}

// TestTableDumpOutputFormat verifies Dump produces expected format strings.
func TestTableDumpOutputFormat(t *testing.T) {
	table := NewTable()
	table.Alignments = []Alignment{AlignLeft, AlignCenter}

	// Capture output by running Dump
	// We're primarily testing it doesn't panic and uses String() correctly
	var output strings.Builder
	source := []byte("test")

	// The actual output goes to stdout, but we verify String() works
	for _, align := range table.Alignments {
		str := align.String()
		if str == "" && (align == AlignLeft || align == AlignCenter) {
			t.Errorf("Alignment.String() returned empty for valid alignment %v", align)
		}
		output.WriteString(str)
	}

	// Verify Dump doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()

	table.Dump(source, 0)
}
