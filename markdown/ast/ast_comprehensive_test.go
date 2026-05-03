package ast

import (
	"bytes"
	"testing"

	textm "github.com/emad-elsaid/xlog/markdown/text"
)

// TestNodeKind tests NodeKind type and registration.
func TestNodeKind(t *testing.T) {
	tests := []struct {
		name         string
		kindName     string
		wantNonEmpty bool
	}{
		{"creates new kind", "TestKind", true},
		{"creates another kind", "AnotherKind", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kind := NewNodeKind(tt.kindName)
			if kind == 0 {
				t.Error("NewNodeKind() returned zero value")
			}
			if got := kind.String(); got != tt.kindName {
				t.Errorf("Kind.String() = %v, want %v", got, tt.kindName)
			}
		})
	}
}

// TestDocument tests Document node operations.
func TestDocument(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			"new document has correct type and kind",
			func(t *testing.T) {
				doc := NewDocument()
				if doc.Type() != TypeDocument {
					t.Errorf("Type() = %v, want %v", doc.Type(), TypeDocument)
				}
				if doc.Kind() != KindDocument {
					t.Errorf("Kind() = %v, want %v", doc.Kind(), KindDocument)
				}
			},
		},
		{
			"owner document returns self",
			func(t *testing.T) {
				doc := NewDocument()
				if doc.OwnerDocument() != doc {
					t.Error("OwnerDocument() did not return self")
				}
			},
		},
		{
			"meta operations work correctly",
			func(t *testing.T) {
				doc := NewDocument()
				meta := doc.Meta()
				if meta == nil {
					t.Error("Meta() returned nil")
				}
				doc.AddMeta("key1", "value1")
				if doc.Meta()["key1"] != "value1" {
					t.Errorf("AddMeta failed, got %v", doc.Meta()["key1"])
				}
				doc.SetMeta(map[string]any{"key2": "value2"})
				if doc.Meta()["key2"] != "value2" {
					t.Errorf("SetMeta failed, got %v", doc.Meta()["key2"])
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestNodeHierarchy tests parent/child relationships.
func TestNodeHierarchy(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			"append child establishes parent-child relationship",
			func(t *testing.T) {
				doc := NewDocument()
				para := NewParagraph()
				doc.AppendChild(doc, para)
				if para.Parent() != doc {
					t.Error("Parent() not set correctly")
				}
				if doc.FirstChild() != para {
					t.Error("FirstChild() not set correctly")
				}
				if doc.HasChildren() != true {
					t.Error("HasChildren() should return true")
				}
				if doc.ChildCount() != 1 {
					t.Errorf("ChildCount() = %d, want 1", doc.ChildCount())
				}
			},
		},
		{
			"multiple children maintain correct order",
			func(t *testing.T) {
				doc := NewDocument()
				para1 := NewParagraph()
				para2 := NewParagraph()
				doc.AppendChild(doc, para1)
				doc.AppendChild(doc, para2)
				if doc.FirstChild() != para1 {
					t.Error("FirstChild() incorrect")
				}
				if doc.LastChild() != para2 {
					t.Error("LastChild() incorrect")
				}
				if para1.NextSibling() != para2 {
					t.Error("NextSibling() incorrect")
				}
				if para2.PreviousSibling() != para1 {
					t.Error("PreviousSibling() incorrect")
				}
				if doc.ChildCount() != 2 {
					t.Errorf("ChildCount() = %d, want 2", doc.ChildCount())
				}
			},
		},
		{
			"remove child works correctly",
			func(t *testing.T) {
				doc := NewDocument()
				para := NewParagraph()
				doc.AppendChild(doc, para)
				doc.RemoveChild(doc, para)
				if doc.HasChildren() {
					t.Error("HasChildren() should return false after removal")
				}
				if para.Parent() != nil {
					t.Error("Parent() should be nil after removal")
				}
			},
		},
		{
			"remove children clears all children",
			func(t *testing.T) {
				doc := NewDocument()
				doc.AppendChild(doc, NewParagraph())
				doc.AppendChild(doc, NewParagraph())
				doc.RemoveChildren(doc)
				if doc.HasChildren() {
					t.Error("HasChildren() should return false")
				}
				if doc.ChildCount() != 0 {
					t.Error("ChildCount() should be 0")
				}
			},
		},
		{
			"replace child works correctly",
			func(t *testing.T) {
				doc := NewDocument()
				para1 := NewParagraph()
				para2 := NewParagraph()
				doc.AppendChild(doc, para1)
				doc.ReplaceChild(doc, para1, para2)
				if doc.FirstChild() != para2 {
					t.Error("ReplaceChild() did not replace correctly")
				}
				if para1.Parent() != nil {
					t.Error("Old child should have nil parent")
				}
			},
		},
		{
			"insert before works correctly",
			func(t *testing.T) {
				doc := NewDocument()
				para1 := NewParagraph()
				para2 := NewParagraph()
				doc.AppendChild(doc, para1)
				doc.InsertBefore(doc, para1, para2)
				if doc.FirstChild() != para2 {
					t.Error("InsertBefore() incorrect")
				}
				if para2.NextSibling() != para1 {
					t.Error("Sibling relationship incorrect")
				}
			},
		},
		{
			"insert after works correctly",
			func(t *testing.T) {
				doc := NewDocument()
				para1 := NewParagraph()
				para2 := NewParagraph()
				doc.AppendChild(doc, para1)
				doc.InsertAfter(doc, para1, para2)
				if para1.NextSibling() != para2 {
					t.Error("InsertAfter() incorrect")
				}
				if doc.LastChild() != para2 {
					t.Error("LastChild() incorrect")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestBlockNodes tests various block node types.
func TestBlockNodes(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		wantType NodeType
		wantKind NodeKind
	}{
		{"TextBlock", NewTextBlock(), TypeBlock, KindTextBlock},
		{"Paragraph", NewParagraph(), TypeBlock, KindParagraph},
		{"Heading level 1", NewHeading(1), TypeBlock, KindHeading},
		{"Heading level 6", NewHeading(6), TypeBlock, KindHeading},
		{"ThematicBreak", NewThematicBreak(), TypeBlock, KindThematicBreak},
		{"CodeBlock", NewCodeBlock(), TypeBlock, KindCodeBlock},
		{"FencedCodeBlock", NewFencedCodeBlock(nil), TypeBlock, KindFencedCodeBlock},
		{"Blockquote", NewBlockquote(), TypeBlock, KindBlockquote},
		{"List ordered", NewList('+'), TypeBlock, KindList},
		{"ListItem", NewListItem(0), TypeBlock, KindListItem},
		{"HTMLBlock type 1", NewHTMLBlock(HTMLBlockType1), TypeBlock, KindHTMLBlock},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.Type(); got != tt.wantType {
				t.Errorf("Type() = %v, want %v", got, tt.wantType)
			}
			if got := tt.node.Kind(); got != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", got, tt.wantKind)
			}
		})
	}
}

// TestInlineNodes tests various inline node types.
func TestInlineNodes(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		wantType NodeType
		wantKind NodeKind
	}{
		{"Text", NewText(), TypeInline, KindText},
		{"String", NewString([]byte("test")), TypeInline, KindString},
		{"CodeSpan", NewCodeSpan(), TypeInline, KindCodeSpan},
		{"Emphasis level 1", NewEmphasis(1), TypeInline, KindEmphasis},
		{"Emphasis level 2", NewEmphasis(2), TypeInline, KindEmphasis},
		{"Link", NewLink(), TypeInline, KindLink},
		{"Image", NewImage(NewLink()), TypeInline, KindImage},
		{"AutoLink email", NewAutoLink(AutoLinkEmail, nil), TypeInline, KindAutoLink},
		{"AutoLink URL", NewAutoLink(AutoLinkURL, nil), TypeInline, KindAutoLink},
		{"RawHTML", NewRawHTML(), TypeInline, KindRawHTML},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.Type(); got != tt.wantType {
				t.Errorf("Type() = %v, want %v", got, tt.wantType)
			}
			if got := tt.node.Kind(); got != tt.wantKind {
				t.Errorf("Kind() = %v, want %v", got, tt.wantKind)
			}
		})
	}
}

// TestHeading tests Heading-specific functionality.
func TestHeading(t *testing.T) {
	tests := []struct {
		name  string
		level int
		want  int
	}{
		{"level 1", 1, 1},
		{"level 3", 3, 3},
		{"level 6", 6, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHeading(tt.level)
			if h.Level != tt.want {
				t.Errorf("Level = %d, want %d", h.Level, tt.want)
			}
		})
	}
}

// TestList tests List-specific functionality.
func TestList(t *testing.T) {
	tests := []struct {
		name       string
		marker     byte
		wantMarker byte
	}{
		{"unordered dash", '-', '-'},
		{"unordered plus", '+', '+'},
		{"unordered star", '*', '*'},
		{"ordered", '.', '.'},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewList(tt.marker)
			if list.Marker != tt.wantMarker {
				t.Errorf("Marker = %c, want %c", list.Marker, tt.wantMarker)
			}
			isOrdered := tt.marker == '.' || tt.marker == ')'
			if list.IsOrdered() != isOrdered {
				t.Errorf("IsOrdered() = %v, want %v", list.IsOrdered(), isOrdered)
			}
		})
	}
}

// TestFencedCodeBlock tests FencedCodeBlock-specific functionality.
func TestFencedCodeBlock(t *testing.T) {
	tests := []struct {
		name     string
		info     *Text
		wantLang []byte
	}{
		{
			"with language info",
			&Text{Segment: textm.NewSegment(0, 2)},
			[]byte("go"),
		},
		{
			"without language info",
			nil,
			nil,
		},
	}
	source := []byte("go\npython\nrust")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fcb := NewFencedCodeBlock(tt.info)
			lang := fcb.Language(source)
			if !bytes.Equal(lang, tt.wantLang) {
				t.Errorf("Language() = %s, want %s", lang, tt.wantLang)
			}
		})
	}
}

// TestText tests Text node functionality.
func TestText(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			"soft line break flag",
			func(t *testing.T) {
				txt := NewText()
				if txt.SoftLineBreak() {
					t.Error("SoftLineBreak() should be false initially")
				}
				txt.SetSoftLineBreak(true)
				if !txt.SoftLineBreak() {
					t.Error("SoftLineBreak() should be true after setting")
				}
				txt.SetSoftLineBreak(false)
				if txt.SoftLineBreak() {
					t.Error("SoftLineBreak() should be false after unsetting")
				}
			},
		},
		{
			"hard line break flag",
			func(t *testing.T) {
				txt := NewText()
				if txt.HardLineBreak() {
					t.Error("HardLineBreak() should be false initially")
				}
				txt.SetHardLineBreak(true)
				if !txt.HardLineBreak() {
					t.Error("HardLineBreak() should be true after setting")
				}
				txt.SetHardLineBreak(false)
				if txt.HardLineBreak() {
					t.Error("HardLineBreak() should be false after unsetting")
				}
			},
		},
		{
			"raw flag",
			func(t *testing.T) {
				txt := NewText()
				if txt.IsRaw() {
					t.Error("IsRaw() should be false initially")
				}
				txt.SetRaw(true)
				if !txt.IsRaw() {
					t.Error("IsRaw() should be true after setting")
				}
				txt.SetRaw(false)
				if txt.IsRaw() {
					t.Error("IsRaw() should be false after unsetting")
				}
			},
		},
		{
			"merge compatible segments",
			func(t *testing.T) {
				source := []byte("hello world")
				txt1 := NewText()
				txt1.Segment = textm.NewSegment(0, 5)
				txt2 := NewText()
				txt2.Segment = textm.NewSegment(5, 11)
				if !txt1.Merge(txt2, source) {
					t.Error("Merge() should succeed for adjacent segments")
				}
				if txt1.Segment.Stop != 11 {
					t.Errorf("Merged segment stop = %d, want 11", txt1.Segment.Stop)
				}
			},
		},
		{
			"merge incompatible types",
			func(t *testing.T) {
				source := []byte("hello world")
				txt := NewText()
				para := NewParagraph()
				if txt.Merge(para, source) {
					t.Error("Merge() should fail for different node types")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestString tests String node functionality.
func TestString(t *testing.T) {
	tests := []struct {
		name  string
		value []byte
		want  []byte
	}{
		{"simple string", []byte("hello"), []byte("hello")},
		{"empty string", []byte(""), []byte("")},
		{"unicode string", []byte("世界"), []byte("世界")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewString(tt.value)
			if !bytes.Equal(s.Value, tt.want) {
				t.Errorf("Value = %s, want %s", s.Value, tt.want)
			}
		})
	}
}

// TestEmphasis tests Emphasis node functionality.
func TestEmphasis(t *testing.T) {
	tests := []struct {
		name  string
		level int
		want  int
	}{
		{"single emphasis", 1, 1},
		{"double emphasis", 2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			em := NewEmphasis(tt.level)
			if em.Level != tt.want {
				t.Errorf("Level = %d, want %d", em.Level, tt.want)
			}
		})
	}
}

// TestLink tests Link node functionality.
func TestLink(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			"new link has empty destination and title",
			func(t *testing.T) {
				link := NewLink()
				if len(link.Destination) != 0 {
					t.Error("Destination should be empty")
				}
				if len(link.Title) != 0 {
					t.Error("Title should be empty")
				}
			},
		},
		{
			"set destination and title",
			func(t *testing.T) {
				link := NewLink()
				link.Destination = []byte("https://example.com")
				link.Title = []byte("Example")
				if !bytes.Equal(link.Destination, []byte("https://example.com")) {
					t.Errorf("Destination = %s", link.Destination)
				}
				if !bytes.Equal(link.Title, []byte("Example")) {
					t.Errorf("Title = %s", link.Title)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestImage tests Image node functionality.
func TestImage(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			"image wraps link",
			func(t *testing.T) {
				link := NewLink()
				link.Destination = []byte("image.png")
				img := NewImage(link)
				if img.Destination == nil {
					t.Error("Destination should not be nil")
				}
				if !bytes.Equal(img.Destination, []byte("image.png")) {
					t.Errorf("Destination = %s", img.Destination)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestAutoLink tests AutoLink node functionality.
func TestAutoLink(t *testing.T) {
	tests := []struct {
		name     string
		autoType AutoLinkType
		wantType AutoLinkType
	}{
		{"email autolink", AutoLinkEmail, AutoLinkEmail},
		{"URL autolink", AutoLinkURL, AutoLinkURL},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := []byte("test@example.com")
			txt := NewText()
			txt.Segment = textm.NewSegment(0, len(source))
			al := NewAutoLink(tt.autoType, txt)
			if al.AutoLinkType != tt.wantType {
				t.Errorf("AutoLinkType = %v, want %v", al.AutoLinkType, tt.wantType)
			}
			label := al.Label(source)
			if !bytes.Equal(label, source) {
				t.Errorf("Label() = %s, want %s", label, source)
			}
		})
	}
}

// TestAttributes tests node attribute operations.
func TestAttributes(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			"set and get attribute",
			func(t *testing.T) {
				doc := NewDocument()
				doc.SetAttribute([]byte("class"), []byte("highlight"))
				val, ok := doc.Attribute([]byte("class"))
				if !ok {
					t.Error("Attribute() returned false")
				}
				if !bytes.Equal(val.([]byte), []byte("highlight")) {
					t.Errorf("Attribute value = %s, want highlight", val)
				}
			},
		},
		{
			"set attribute string",
			func(t *testing.T) {
				doc := NewDocument()
				doc.SetAttributeString("id", "main")
				val, ok := doc.AttributeString("id")
				if !ok {
					t.Error("AttributeString() returned false")
				}
				if val != "main" {
					t.Errorf("AttributeString() = %s, want main", val)
				}
			},
		},
		{
			"get all attributes",
			func(t *testing.T) {
				doc := NewDocument()
				doc.SetAttribute([]byte("class"), []byte("test"))
				doc.SetAttribute([]byte("id"), []byte("main"))
				attrs := doc.Attributes()
				if len(attrs) != 2 {
					t.Errorf("len(Attributes()) = %d, want 2", len(attrs))
				}
			},
		},
		{
			"remove attributes",
			func(t *testing.T) {
				doc := NewDocument()
				doc.SetAttribute([]byte("class"), []byte("test"))
				doc.RemoveAttributes()
				attrs := doc.Attributes()
				if len(attrs) != 0 {
					t.Errorf("len(Attributes()) = %d, want 0 after removal", len(attrs))
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestBlockLineOperations tests block node line operations.
func TestBlockLineOperations(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			"set and get lines",
			func(t *testing.T) {
				para := NewParagraph()
				lines := textm.NewSegments()
				lines.Append(textm.NewSegment(0, 5))
				para.SetLines(lines)
				retrieved := para.Lines()
				if retrieved.Len() != 1 {
					t.Errorf("Lines().Len() = %d, want 1", retrieved.Len())
				}
			},
		},
		{
			"blank previous lines flag",
			func(t *testing.T) {
				para := NewParagraph()
				if para.HasBlankPreviousLines() {
					t.Error("HasBlankPreviousLines() should be false initially")
				}
				para.SetBlankPreviousLines(true)
				if !para.HasBlankPreviousLines() {
					t.Error("HasBlankPreviousLines() should be true after setting")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestIsParagraph tests the IsParagraph utility function.
func TestIsParagraph(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want bool
	}{
		{"paragraph is paragraph", NewParagraph(), true},
		{"document is not paragraph", NewDocument(), false},
		{"heading is not paragraph", NewHeading(1), false},
		{"text is not paragraph", NewText(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsParagraph(tt.node); got != tt.want {
				t.Errorf("IsParagraph() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestHTMLBlock tests HTMLBlock-specific functionality.
func TestHTMLBlock(t *testing.T) {
	tests := []struct {
		name     string
		htmlType HTMLBlockType
		want     HTMLBlockType
	}{
		{"type 1", HTMLBlockType1, HTMLBlockType1},
		{"type 2", HTMLBlockType2, HTMLBlockType2},
		{"type 3", HTMLBlockType3, HTMLBlockType3},
		{"type 4", HTMLBlockType4, HTMLBlockType4},
		{"type 5", HTMLBlockType5, HTMLBlockType5},
		{"type 6", HTMLBlockType6, HTMLBlockType6},
		{"type 7", HTMLBlockType7, HTMLBlockType7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hb := NewHTMLBlock(tt.htmlType)
			if hb.HTMLBlockType != tt.want {
				t.Errorf("HTMLBlockType = %v, want %v", hb.HTMLBlockType, tt.want)
			}
			if !hb.IsRaw() {
				t.Error("HTMLBlock should have IsRaw() = true")
			}
			if hb.HasClosure() {
				t.Error("HTMLBlock should have HasClosure() = false initially")
			}
		})
	}
}

// TestListItem tests ListItem-specific functionality.
func TestListItem(t *testing.T) {
	tests := []struct {
		name   string
		offset int
		want   int
	}{
		{"zero offset", 0, 0},
		{"positive offset", 4, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			li := NewListItem(tt.offset)
			if li.Offset != tt.want {
				t.Errorf("Offset = %d, want %d", li.Offset, tt.want)
			}
		})
	}
}

// TestSortChildren tests SortChildren functionality.
func TestSortChildren(t *testing.T) {
	doc := NewDocument()
	h1 := NewHeading(1)
	h2 := NewHeading(2)
	h3 := NewHeading(3)
	doc.AppendChild(doc, h3)
	doc.AppendChild(doc, h1)
	doc.AppendChild(doc, h2)

	doc.SortChildren(func(n1, n2 Node) int {
		l1 := n1.(*Heading).Level
		l2 := n2.(*Heading).Level
		return l1 - l2
	})

	if doc.FirstChild().(*Heading).Level != 1 {
		t.Errorf("First child after sort should be level 1")
	}
	if doc.LastChild().(*Heading).Level != 3 {
		t.Errorf("Last child after sort should be level 3")
	}
}

// TestOwnerDocument tests OwnerDocument navigation.
func TestOwnerDocument(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			"deeply nested node finds owner document",
			func(t *testing.T) {
				doc := NewDocument()
				para := NewParagraph()
				txt := NewText()
				doc.AppendChild(doc, para)
				para.AppendChild(para, txt)
				if txt.OwnerDocument() != doc {
					t.Error("OwnerDocument() should return root document")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestDeprecatedTextMethods tests deprecated Text() methods.
func TestDeprecatedTextMethods(t *testing.T) {
	source := []byte("hello world\ntest content")
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			"TextBlock.Text returns lines value",
			func(t *testing.T) {
				tb := NewTextBlock()
				lines := textm.NewSegments()
				lines.Append(textm.NewSegment(0, 5))
				tb.SetLines(lines)
				text := tb.Text(source)
				if !bytes.Equal(text, []byte("hello")) {
					t.Errorf("Text() = %s, want hello", text)
				}
			},
		},
		{
			"Paragraph.Text returns lines value",
			func(t *testing.T) {
				p := NewParagraph()
				lines := textm.NewSegments()
				lines.Append(textm.NewSegment(0, 11))
				p.SetLines(lines)
				text := p.Text(source)
				if !bytes.Equal(text, []byte("hello world")) {
					t.Errorf("Text() = %s, want 'hello world'", text)
				}
			},
		},
		{
			"CodeBlock.Text returns lines value",
			func(t *testing.T) {
				cb := NewCodeBlock()
				lines := textm.NewSegments()
				lines.Append(textm.NewSegment(12, 24))
				cb.SetLines(lines)
				text := cb.Text(source)
				if !bytes.Equal(text, []byte("test content")) {
					t.Errorf("Text() = %s", text)
				}
			},
		},
		{
			"FencedCodeBlock.Text returns lines value",
			func(t *testing.T) {
				fcb := NewFencedCodeBlock(nil)
				lines := textm.NewSegments()
				lines.Append(textm.NewSegment(0, 5))
				fcb.SetLines(lines)
				text := fcb.Text(source)
				if !bytes.Equal(text, []byte("hello")) {
					t.Errorf("Text() = %s", text)
				}
			},
		},
		{
			"HTMLBlock.Text returns lines value",
			func(t *testing.T) {
				hb := NewHTMLBlock(HTMLBlockType1)
				lines := textm.NewSegments()
				lines.Append(textm.NewSegment(0, 5))
				hb.SetLines(lines)
				text := hb.Text(source)
				if !bytes.Equal(text, []byte("hello")) {
					t.Errorf("Text() = %s", text)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

// TestIsRawMethods tests IsRaw() for various node types.
func TestIsRawMethods(t *testing.T) {
	tests := []struct {
		name    string
		node    Node
		wantRaw bool
	}{
		{"BaseBlock is not raw", NewTextBlock(), false},
		{"CodeBlock is raw", NewCodeBlock(), true},
		{"FencedCodeBlock is raw", NewFencedCodeBlock(nil), true},
		{"HTMLBlock is raw", NewHTMLBlock(HTMLBlockType1), true},
		{"BaseInline is not raw", NewText(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.IsRaw(); got != tt.wantRaw {
				t.Errorf("IsRaw() = %v, want %v", got, tt.wantRaw)
			}
		})
	}
}

// TestDumpMethods tests Dump() methods for various nodes.
func TestDumpMethods(t *testing.T) {
	source := []byte("test content")
	tests := []struct {
		name string
		node Node
	}{
		{"Document", NewDocument()},
		{"TextBlock", NewTextBlock()},
		{"Paragraph", NewParagraph()},
		{"Heading", NewHeading(1)},
		{"ThematicBreak", NewThematicBreak()},
		{"CodeBlock", NewCodeBlock()},
		{"FencedCodeBlock", NewFencedCodeBlock(nil)},
		{"Blockquote", NewBlockquote()},
		{"List", NewList('-')},
		{"ListItem", NewListItem(0)},
		{"HTMLBlock", NewHTMLBlock(HTMLBlockType1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just ensure Dump doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Dump() panicked: %v", r)
				}
			}()
			tt.node.Dump(source, 0)
		})
	}
}

// TestListCanContinue tests List.CanContinue.
func TestListCanContinue(t *testing.T) {
	tests := []struct {
		name        string
		marker      byte
		isOrdered   bool
		testMarker  byte
		testOrdered bool
		want        bool
	}{
		{"same marker unordered", '-', false, '-', false, true},
		{"different marker unordered", '-', false, '*', false, false},
		{"same marker ordered", '.', true, '.', true, true},
		{"ordered vs unordered mismatch", '.', true, '.', false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewList(tt.marker)
			if got := list.CanContinue(tt.testMarker, tt.testOrdered); got != tt.want {
				t.Errorf("CanContinue(%c, %v) = %v, want %v", tt.testMarker, tt.testOrdered, got, tt.want)
			}
		})
	}
}

// TestInlinePanicMethods tests that inline nodes panic on block-only methods.
func TestInlinePanicMethods(t *testing.T) {
	txt := NewText()
	tests := []struct {
		name string
		fn   func()
	}{
		{
			"HasBlankPreviousLines panics",
			func() { txt.HasBlankPreviousLines() },
		},
		{
			"SetBlankPreviousLines panics",
			func() { txt.SetBlankPreviousLines(true) },
		},
		{
			"Lines panics",
			func() { txt.Lines() },
		},
		{
			"SetLines panics",
			func() { txt.SetLines(textm.NewSegments()) },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("Expected panic but didn't get one")
				}
			}()
			tt.fn()
		})
	}
}

// TestCodeSpan tests CodeSpan functionality.
func TestCodeSpan(t *testing.T) {
	cs := NewCodeSpan()
	if cs.Kind() != KindCodeSpan {
		t.Errorf("Kind() = %v, want %v", cs.Kind(), KindCodeSpan)
	}
	if cs.Type() != TypeInline {
		t.Errorf("Type() = %v, want %v", cs.Type(), TypeInline)
	}
}

// TestRawHTML tests RawHTML functionality.
func TestRawHTML(t *testing.T) {
	rh := NewRawHTML()
	if rh.Kind() != KindRawHTML {
		t.Errorf("Kind() = %v, want %v", rh.Kind(), KindRawHTML)
	}
	if rh.Type() != TypeInline {
		t.Errorf("Type() = %v, want %v", rh.Type(), TypeInline)
	}
	if rh.Segments == nil {
		t.Error("Segments should not be nil")
	}
}

// TestTextSegmentOperations tests TextSegment creation and operations.
func TestTextSegmentOperations(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			"NewTextSegment creates text node with segment",
			func(t *testing.T) {
				txt := NewTextSegment(textm.NewSegment(0, 5))
				if txt.Segment.Start != 0 || txt.Segment.Stop != 5 {
					t.Errorf("Segment = %v, want [0:5]", txt.Segment)
				}
			},
		},
		{
			"NewRawTextSegment creates raw text",
			func(t *testing.T) {
				txt := NewRawTextSegment(textm.NewSegment(0, 5))
				if !txt.IsRaw() {
					t.Error("NewRawTextSegment should create raw text")
				}
			},
		},
		{
			"MergeOrAppendTextSegment merges compatible segments",
			func(t *testing.T) {
				parent := NewParagraph()
				txt1 := NewTextSegment(textm.NewSegment(0, 5))
				parent.AppendChild(parent, txt1)
				MergeOrAppendTextSegment(parent, textm.NewSegment(5, 11))
				if parent.ChildCount() != 1 {
					t.Errorf("Should have merged into 1 child, got %d", parent.ChildCount())
				}
			},
		},
		{
			"MergeOrReplaceTextSegment replaces non-mergeable",
			func(t *testing.T) {
				parent := NewParagraph()
				txt1 := NewTextSegment(textm.NewSegment(0, 5))
				parent.AppendChild(parent, txt1)
				MergeOrReplaceTextSegment(parent, txt1, textm.NewSegment(12, 16))
				if parent.ChildCount() != 1 {
					t.Error("Should have 1 child after replace")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestCodeFlag tests String.IsCode and SetCode.
func TestCodeFlag(t *testing.T) {
	str := NewString([]byte("code"))
	if str.IsCode() {
		t.Error("IsCode() should be false initially")
	}
	str.SetCode(true)
	if !str.IsCode() {
		t.Error("IsCode() should be true after setting")
	}
	str.SetCode(false)
	if str.IsCode() {
		t.Error("IsCode() should be false after unsetting")
	}
}

// TestCodeSpanIsBlank tests CodeSpan.IsBlank.
func TestCodeSpanIsBlank(t *testing.T) {
	tests := []struct {
		name   string
		source []byte
		setup  func(*CodeSpan)
		want   bool
	}{
		{
			"empty code span is blank",
			[]byte("   "),
			func(cs *CodeSpan) {
				txt := NewTextSegment(textm.NewSegment(0, 3))
				cs.AppendChild(cs, txt)
			},
			true,
		},
		{
			"code span with text is not blank",
			[]byte("hello"),
			func(cs *CodeSpan) {
				txt := NewTextSegment(textm.NewSegment(0, 5))
				cs.AppendChild(cs, txt)
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := NewCodeSpan()
			tt.setup(cs)
			if got := cs.IsBlank(tt.source); got != tt.want {
				t.Errorf("IsBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestInlineMethods tests Inline() marker method.
func TestInlineMethods(t *testing.T) {
	nodes := []Node{
		NewText(),
		NewString([]byte("test")),
		NewCodeSpan(),
		NewEmphasis(1),
		NewLink(),
		NewImage(NewLink()),
		NewAutoLink(AutoLinkURL, NewText()),
		NewRawHTML(),
	}
	for _, node := range nodes {
		// Inline() is just a marker method, calling it shouldn't panic
		if inl, ok := node.(interface{ Inline() }); ok {
			inl.Inline()
		}
	}
}

// TestStringIsRawAndSetRaw tests String raw flag.
func TestStringIsRawAndSetRaw(t *testing.T) {
	str := NewString([]byte("test"))
	if str.IsRaw() {
		t.Error("IsRaw() should be false initially")
	}
	str.SetRaw(true)
	if !str.IsRaw() {
		t.Error("IsRaw() should be true after setting")
	}
	str.SetRaw(false)
	if str.IsRaw() {
		t.Error("IsRaw() should be false after unsetting")
	}
}

// TestTextDumpMethod tests Text.Dump.
func TestTextDumpMethod(t *testing.T) {
	source := []byte("hello world")
	txt := NewTextSegment(textm.NewSegment(0, 5))
	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()
	txt.Dump(source, 0)
}

// TestStringDumpMethod tests String.Dump.
func TestStringDumpMethod(t *testing.T) {
	source := []byte("test")
	str := NewString([]byte("test"))
	// Should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()
	str.Dump(source, 0)
}

// TestCodeSpanDump tests CodeSpan.Dump.
func TestCodeSpanDump(t *testing.T) {
	source := []byte("code")
	cs := NewCodeSpan()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()
	cs.Dump(source, 0)
}

// TestEmphasisDump tests Emphasis.Dump.
func TestEmphasisDump(t *testing.T) {
	source := []byte("text")
	em := NewEmphasis(1)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()
	em.Dump(source, 0)
}

// TestLinkDump tests Link.Dump.
func TestLinkDump(t *testing.T) {
	source := []byte("link")
	link := NewLink()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()
	link.Dump(source, 0)
}

// TestImageDump tests Image.Dump.
func TestImageDump(t *testing.T) {
	source := []byte("image")
	img := NewImage(NewLink())
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()
	img.Dump(source, 0)
}

// TestAutoLinkDump tests AutoLink.Dump.
func TestAutoLinkDump(t *testing.T) {
	source := []byte("link")
	al := NewAutoLink(AutoLinkURL, NewText())
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Dump() panicked: %v", r)
		}
	}()
	al.Dump(source, 0)
}

// TestAutoLinkURL tests AutoLink.URL with protocol.
func TestAutoLinkURL(t *testing.T) {
	tests := []struct {
		name     string
		protocol []byte
		source   []byte
		segment  textm.Segment
		want     []byte
	}{
		{
			"with protocol",
			[]byte("http"),
			[]byte("example.com"),
			textm.NewSegment(0, 11),
			[]byte("http://example.com"),
		},
		{
			"without protocol",
			nil,
			[]byte("test@example.com"),
			textm.NewSegment(0, 16),
			[]byte("test@example.com"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txt := NewText()
			txt.Segment = tt.segment
			al := NewAutoLink(AutoLinkURL, txt)
			al.Protocol = tt.protocol
			url := al.URL(tt.source)
			if !bytes.Equal(url, tt.want) {
				t.Errorf("URL() = %s, want %s", url, tt.want)
			}
		})
	}
}

// TestDeprecatedNodeTextMethods tests deprecated Text() methods on inline nodes.
func TestDeprecatedNodeTextMethods(t *testing.T) {
	source := []byte("hello world")
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			"BaseNode.Text",
			func(t *testing.T) {
				doc := NewDocument()
				txt := doc.Text(source)
				if txt != nil {
					t.Logf("BaseNode.Text() = %s", txt)
				}
			},
		},
		{
			"Text.Text returns segment value",
			func(t *testing.T) {
				txt := NewTextSegment(textm.NewSegment(0, 5))
				val := txt.Text(source)
				if !bytes.Equal(val, []byte("hello")) {
					t.Errorf("Text() = %s, want hello", val)
				}
			},
		},
		{
			"String.Text returns value",
			func(t *testing.T) {
				str := NewString([]byte("test"))
				val := str.Text(source)
				if !bytes.Equal(val, []byte("test")) {
					t.Errorf("Text() = %s, want test", val)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
