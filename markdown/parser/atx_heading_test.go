package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

const testHeadingID = "test"

func TestATXHeadingParser_Trigger(t *testing.T) {
	parser := NewATXHeadingParser()
	trigger := parser.Trigger()
	if len(trigger) != 1 || trigger[0] != '#' {
		t.Errorf("Trigger() = %v, want ['#']", trigger)
	}
}

func TestATXHeadingParser_Open(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantLevel   int
		wantText    string
		wantState   State
		wantNil     bool
		blockOffset int
	}{
		{
			name:        "level 1 heading",
			input:       "# Heading 1\n",
			wantLevel:   1,
			wantText:    "Heading 1",
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "level 2 heading",
			input:       "## Heading 2\n",
			wantLevel:   2,
			wantText:    "Heading 2",
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "level 3 heading",
			input:       "### Heading 3\n",
			wantLevel:   3,
			wantText:    "Heading 3",
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "level 6 heading max",
			input:       "###### Heading 6\n",
			wantLevel:   6,
			wantText:    "Heading 6",
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "level 7 exceeds max - should fail",
			input:       "####### Not a heading\n",
			wantNil:     true,
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "no space after hash - should fail",
			input:       "#NoSpace\n",
			wantNil:     true,
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "just hash alone",
			input:       "#",
			wantLevel:   1,
			wantText:    "",
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "heading with trailing hashes",
			input:       "## Heading ##\n",
			wantLevel:   2,
			wantText:    "Heading ",
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "heading with trailing hashes no space",
			input:       "## Heading##\n",
			wantLevel:   2,
			wantText:    "Heading##",
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "empty heading with just hashes",
			input:       "### ###\n",
			wantLevel:   3,
			wantText:    "",
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "heading with multiple spaces",
			input:       "#    Multiple   Spaces\n",
			wantLevel:   1,
			wantText:    "Multiple   Spaces",
			wantState:   NoChildren,
			blockOffset: 0,
		},
		{
			name:        "negative block offset",
			input:       "# Heading\n",
			wantNil:     true,
			wantState:   NoChildren,
			blockOffset: -1,
		},
		{
			name:        "only spaces after hash",
			input:       "#   \n",
			wantLevel:   1,
			wantText:    "",
			wantState:   NoChildren,
			blockOffset: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewATXHeadingParser()
			reader := text.NewReader([]byte(tc.input))
			ctx := NewContext()
			ctx.SetBlockOffset(tc.blockOffset)

			node, state := parser.Open(nil, reader, ctx)

			if tc.wantNil {
				if node != nil {
					t.Errorf("Open() node = %v, want nil", node)
				}
			} else {
				if node == nil {
					t.Fatal("Open() node = nil, want non-nil")
				}
				heading, ok := node.(*ast.Heading)
				if !ok {
					t.Fatalf("Open() node type = %T, want *ast.Heading", node)
				}
				if heading.Level != tc.wantLevel {
					t.Errorf("Open() level = %d, want %d", heading.Level, tc.wantLevel)
				}

				// Check text content if expected
				if tc.wantText != "" && heading.Lines().Len() > 0 {
					line := heading.Lines().At(0)
					text := string(line.Value(reader.Source()))
					if text != tc.wantText {
						t.Errorf("Open() text = %q, want %q", text, tc.wantText)
					}
				}
			}

			if state != tc.wantState {
				t.Errorf("Open() state = %v, want %v", state, tc.wantState)
			}
		})
	}
}

func TestATXHeadingParser_Continue(t *testing.T) {
	parser := NewATXHeadingParser()
	reader := text.NewReader([]byte("# Heading\n"))
	ctx := NewContext()
	heading := ast.NewHeading(1)

	state := parser.Continue(heading, reader, ctx)
	if state != Close {
		t.Errorf("Continue() = %v, want Close", state)
	}
}

func TestATXHeadingParser_CanInterruptParagraph(t *testing.T) {
	parser := NewATXHeadingParser()
	result := parser.CanInterruptParagraph()
	if !result {
		t.Error("CanInterruptParagraph() = false, want true")
	}
}

func TestATXHeadingParser_CanAcceptIndentedLine(t *testing.T) {
	parser := NewATXHeadingParser()
	result := parser.CanAcceptIndentedLine()
	if result {
		t.Error("CanAcceptIndentedLine() = true, want false")
	}
}

func TestATXHeadingParser_WithAutoHeadingID(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantID    string
		wantLevel int
	}{
		{
			name:      "simple heading generates id",
			input:     "# Hello World\n",
			wantID:    "hello-world",
			wantLevel: 1,
		},
		{
			name:      "heading with special chars",
			input:     "## Hello, World!\n",
			wantID:    "hello-world",
			wantLevel: 2,
		},
		{
			name:      "heading with hyphens",
			input:     "### Hello-World\n",
			wantID:    "hello-world",
			wantLevel: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewATXHeadingParser(WithAutoHeadingID())
			reader := text.NewReader([]byte(tc.input))
			ctx := NewContext()
			ctx.SetBlockOffset(0)

			node, _ := parser.Open(nil, reader, ctx)
			if node == nil {
				t.Fatal("Open() returned nil node")
			}

			// Close is where auto ID generation happens
			parser.Close(node, reader, ctx)

			heading := node.(*ast.Heading)
			if heading.Level != tc.wantLevel {
				t.Errorf("Level = %d, want %d", heading.Level, tc.wantLevel)
			}

			id, ok := heading.AttributeString("id")
			if !ok {
				t.Fatal("Expected id attribute, got none")
			}

			if string(id.([]byte)) != tc.wantID {
				t.Errorf("ID = %q, want %q", id, tc.wantID)
			}
		})
	}
}

func TestATXHeadingParser_WithAttribute(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantID    string
		wantClass string
		wantLevel int
		hasID     bool
		hasClass  bool
	}{
		{
			name:      "heading with id attribute",
			input:     "# Heading {#custom-id}\n",
			wantID:    "custom-id",
			wantLevel: 1,
			hasID:     true,
		},
		{
			name:      "heading with class attribute",
			input:     "## Heading {.my-class}\n",
			wantClass: "my-class",
			wantLevel: 2,
			hasClass:  true,
		},
		{
			name:      "heading with trailing hashes and attributes",
			input:     "### Heading ### {#custom}\n",
			wantID:    "custom",
			wantLevel: 3,
			hasID:     true,
		},
		{
			name:      "heading without attributes",
			input:     "# Plain Heading\n",
			wantLevel: 1,
			hasID:     false,
			hasClass:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewATXHeadingParser(WithHeadingAttribute())
			reader := text.NewReader([]byte(tc.input))
			ctx := NewContext()
			ctx.SetBlockOffset(0)

			node, _ := parser.Open(nil, reader, ctx)
			if node == nil {
				t.Fatal("Open() returned nil node")
			}

			parser.Close(node, reader, ctx)

			heading := node.(*ast.Heading)
			if heading.Level != tc.wantLevel {
				t.Errorf("Level = %d, want %d", heading.Level, tc.wantLevel)
			}

			if tc.hasID {
				id, ok := heading.AttributeString("id")
				if !ok {
					t.Fatal("Expected id attribute, got none")
				}
				if string(id.([]byte)) != tc.wantID {
					t.Errorf("ID = %q, want %q", id, tc.wantID)
				}
			}

			if tc.hasClass {
				class, ok := heading.AttributeString("class")
				if !ok {
					t.Fatal("Expected class attribute, got none")
				}
				if string(class.([]byte)) != tc.wantClass {
					t.Errorf("Class = %q, want %q", class, tc.wantClass)
				}
			}
		})
	}
}

func TestHeadingConfig_SetOption(t *testing.T) {
	tests := []struct {
		name       string
		optionName OptionName
		wantAutoID bool
		wantAttr   bool
	}{
		{
			name:       "set auto heading id",
			optionName: optAutoHeadingID,
			wantAutoID: true,
			wantAttr:   false,
		},
		{
			name:       "set attribute",
			optionName: optAttribute,
			wantAutoID: false,
			wantAttr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config := &HeadingConfig{}
			config.SetOption(tc.optionName, nil)

			if config.AutoHeadingID != tc.wantAutoID {
				t.Errorf("AutoHeadingID = %v, want %v", config.AutoHeadingID, tc.wantAutoID)
			}
			if config.Attribute != tc.wantAttr {
				t.Errorf("Attribute = %v, want %v", config.Attribute, tc.wantAttr)
			}
		})
	}
}

func TestWithAutoHeadingID_Options(t *testing.T) {
	opt := WithAutoHeadingID()

	// Test SetParserOption
	cfg := NewConfig()
	opt.SetParserOption(cfg)
	if val, ok := cfg.Options[optAutoHeadingID]; !ok || val != true {
		t.Error("SetParserOption should set optAutoHeadingID to true")
	}

	// Test SetHeadingOption
	headingCfg := &HeadingConfig{}
	opt.SetHeadingOption(headingCfg)
	if !headingCfg.AutoHeadingID {
		t.Error("SetHeadingOption should set AutoHeadingID to true")
	}
}

func TestWithHeadingAttribute_Options(t *testing.T) {
	opt := WithHeadingAttribute()

	// Test SetHeadingOption
	headingCfg := &HeadingConfig{}
	opt.SetHeadingOption(headingCfg)
	if !headingCfg.Attribute {
		t.Error("SetHeadingOption should set Attribute to true")
	}
}

func TestATXHeadingParser_Close_WithoutOptions(t *testing.T) {
	parser := NewATXHeadingParser()
	reader := text.NewReader([]byte("# Heading\n"))
	ctx := NewContext()
	ctx.SetBlockOffset(0)

	node, _ := parser.Open(nil, reader, ctx)
	if node == nil {
		t.Fatal("Open() returned nil node")
	}

	parser.Close(node, reader, ctx)

	heading := node.(*ast.Heading)
	_, hasID := heading.AttributeString("id")
	if hasID {
		t.Error("Heading should not have auto-generated id when WithAutoHeadingID not set")
	}
}

func TestATXHeadingParser_EmptyHeading(t *testing.T) {
	parser := NewATXHeadingParser(WithAutoHeadingID())
	reader := text.NewReader([]byte("# \n"))
	ctx := NewContext()
	ctx.SetBlockOffset(0)

	node, _ := parser.Open(nil, reader, ctx)
	if node == nil {
		t.Fatal("Open() returned nil for empty heading - actual behavior creates empty heading node")
	}

	heading := node.(*ast.Heading)
	if heading.Level != 1 {
		t.Errorf("Level = %d, want 1", heading.Level)
	}
}

func TestATXHeadingParser_AttributeWithEscapedBrace(t *testing.T) {
	parser := NewATXHeadingParser(WithHeadingAttribute())
	reader := text.NewReader([]byte("# Heading \\{not-attribute}\n"))
	ctx := NewContext()
	ctx.SetBlockOffset(0)

	node, _ := parser.Open(nil, reader, ctx)
	if node == nil {
		t.Fatal("Open() returned nil node")
	}

	parser.Close(node, reader, ctx)

	heading := node.(*ast.Heading)
	_, hasID := heading.AttributeString("id")
	if hasID {
		t.Error("Escaped brace should not be parsed as attribute")
	}
}

func TestATXHeadingParser_DuplicateIDs(t *testing.T) {
	parser := NewATXHeadingParser(WithAutoHeadingID())
	ctx := NewContext()
	ctx.SetBlockOffset(0)

	// First heading with "test"
	reader1 := text.NewReader([]byte("# Test\n"))
	node1, _ := parser.Open(nil, reader1, ctx)
	if node1 == nil {
		t.Fatal("First Open() returned nil node")
	}
	parser.Close(node1, reader1, ctx)

	id1, _ := node1.(*ast.Heading).AttributeString("id")

	// Second heading with "test" should get "test-1"
	reader2 := text.NewReader([]byte("# Test\n"))
	node2, _ := parser.Open(nil, reader2, ctx)
	if node2 == nil {
		t.Fatal("Second Open() returned nil node")
	}
	parser.Close(node2, reader2, ctx)

	id2, _ := node2.(*ast.Heading).AttributeString("id")

	if string(id1.([]byte)) != testHeadingID {
		t.Errorf("First ID = %q, want %q", id1, testHeadingID)
	}
	if string(id2.([]byte)) != "test-1" {
		t.Errorf("Second ID = %q, want 'test-1'", id2)
	}
}

func TestATXHeadingParser_CustomIDPreventsDuplication(t *testing.T) {
	parser := NewATXHeadingParser(WithHeadingAttribute(), WithAutoHeadingID())
	ctx := NewContext()
	ctx.SetBlockOffset(0)

	// Heading with custom id
	reader1 := text.NewReader([]byte("# Test {#custom}\n"))
	node1, _ := parser.Open(nil, reader1, ctx)
	if node1 == nil {
		t.Fatal("First Open() returned nil node")
	}
	parser.Close(node1, reader1, ctx)

	id1, _ := node1.(*ast.Heading).AttributeString("id")
	if string(id1.([]byte)) != "custom" {
		t.Errorf("Custom ID = %q, want 'custom'", id1)
	}

	// Second heading should not regenerate id if already has custom id
	reader2 := text.NewReader([]byte("# Other {#other}\n"))
	node2, _ := parser.Open(nil, reader2, ctx)
	if node2 == nil {
		t.Fatal("Second Open() returned nil node")
	}
	parser.Close(node2, reader2, ctx)

	id2, _ := node2.(*ast.Heading).AttributeString("id")
	if string(id2.([]byte)) != "other" {
		t.Errorf("Custom ID = %q, want 'other'", id2)
	}
}
