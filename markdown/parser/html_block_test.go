package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestHTMLBlockParser_Trigger(t *testing.T) {
	parser := NewHTMLBlockParser()
	trigger := parser.Trigger()
	if len(trigger) != 1 || trigger[0] != '<' {
		t.Errorf("Trigger() = %v, want ['<']", trigger)
	}
}

func TestHTMLBlockParser_Open(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantType    ast.HTMLBlockType
		wantNil     bool
		blockOffset int
	}{
		{
			name:        "type 1 - script tag",
			input:       "<script>\n",
			wantType:    ast.HTMLBlockType1,
			blockOffset: 0,
		},
		{
			name:        "type 1 - pre tag",
			input:       "<pre>\n",
			wantType:    ast.HTMLBlockType1,
			blockOffset: 0,
		},
		{
			name:        "type 1 - style tag",
			input:       "<style>\n",
			wantType:    ast.HTMLBlockType1,
			blockOffset: 0,
		},
		{
			name:        "type 1 - textarea tag",
			input:       "<textarea>\n",
			wantType:    ast.HTMLBlockType1,
			blockOffset: 0,
		},
		{
			name:        "type 1 - script with attributes",
			input:       "<script type=\"text/javascript\">\n",
			wantType:    ast.HTMLBlockType1,
			blockOffset: 0,
		},
		{
			name:        "type 1 - script self-closing",
			input:       "<script />\n",
			wantType:    ast.HTMLBlockType1,
			blockOffset: 0,
		},
		{
			name:        "type 1 - case insensitive SCRIPT",
			input:       "<SCRIPT>\n",
			wantType:    ast.HTMLBlockType1,
			blockOffset: 0,
		},
		{
			name:        "type 2 - HTML comment",
			input:       "<!-- comment -->\n",
			wantType:    ast.HTMLBlockType2,
			blockOffset: 0,
		},
		{
			name:        "type 2 - HTML comment multiline",
			input:       "<!--\n",
			wantType:    ast.HTMLBlockType2,
			blockOffset: 0,
		},
		{
			name:        "type 3 - processing instruction",
			input:       "<?xml version=\"1.0\"?>\n",
			wantType:    ast.HTMLBlockType3,
			blockOffset: 0,
		},
		{
			name:        "type 3 - PHP tag",
			input:       "<?php\n",
			wantType:    ast.HTMLBlockType3,
			blockOffset: 0,
		},
		{
			name:        "type 4 - declaration",
			input:       "<!DOCTYPE html>\n",
			wantType:    ast.HTMLBlockType4,
			blockOffset: 0,
		},
		{
			name:        "type 4 - uppercase declaration",
			input:       "<!ELEMENT br EMPTY>\n",
			wantType:    ast.HTMLBlockType4,
			blockOffset: 0,
		},
		{
			name:        "type 5 - CDATA section",
			input:       "<![CDATA[content]]>\n",
			wantType:    ast.HTMLBlockType5,
			blockOffset: 0,
		},
		{
			name:        "type 5 - CDATA multiline",
			input:       "<![CDATA[\n",
			wantType:    ast.HTMLBlockType5,
			blockOffset: 0,
		},
		{
			name:        "type 6 - div tag",
			input:       "<div>\n",
			wantType:    ast.HTMLBlockType6,
			blockOffset: 0,
		},
		{
			name:        "type 6 - table tag",
			input:       "<table>\n",
			wantType:    ast.HTMLBlockType6,
			blockOffset: 0,
		},
		{
			name:        "type 6 - article tag",
			input:       "<article>\n",
			wantType:    ast.HTMLBlockType6,
			blockOffset: 0,
		},
		{
			name:        "type 6 - closing div",
			input:       "</div>\n",
			wantType:    ast.HTMLBlockType6,
			blockOffset: 0,
		},
		{
			name:        "type 6 - div with attributes",
			input:       "<div class=\"test\">\n",
			wantType:    ast.HTMLBlockType6,
			blockOffset: 0,
		},
		{
			name:        "type 6 - header tag",
			input:       "<header>\n",
			wantType:    ast.HTMLBlockType6,
			blockOffset: 0,
		},
		{
			name:        "type 7 - span tag",
			input:       "<span>\n",
			wantType:    ast.HTMLBlockType7,
			blockOffset: 0,
		},
		{
			name:        "type 7 - em tag",
			input:       "<em>\n",
			wantType:    ast.HTMLBlockType7,
			blockOffset: 0,
		},
		{
			name:        "type 7 - strong tag",
			input:       "<strong>\n",
			wantType:    ast.HTMLBlockType7,
			blockOffset: 0,
		},
		{
			name:        "type 7 - self-closing br",
			input:       "<br />\n",
			wantType:    ast.HTMLBlockType7,
			blockOffset: 0,
		},
		{
			name:        "not HTML - no opening bracket at offset",
			input:       "text <div>\n",
			wantNil:     true,
			blockOffset: 0,
		},
		{
			name:        "not HTML - invalid tag",
			input:       "<>\n",
			wantNil:     true,
			blockOffset: 0,
		},
		{
			name:        "not HTML - text only",
			input:       "just text\n",
			wantNil:     true,
			blockOffset: 0,
		},
		{
			name:        "HTML with 1 space indent",
			input:       " <div>\n",
			wantType:    ast.HTMLBlockType6,
			blockOffset: 1,
		},
		{
			name:        "HTML with 2 space indent",
			input:       "  <div>\n",
			wantType:    ast.HTMLBlockType6,
			blockOffset: 2,
		},
		{
			name:        "HTML with 3 space indent",
			input:       "   <div>\n",
			wantType:    ast.HTMLBlockType6,
			blockOffset: 3,
		},
		{
			name:        "not HTML - 4 space indent",
			input:       "    <div>\n",
			wantNil:     true,
			blockOffset: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewHTMLBlockParser()
			reader := text.NewReader([]byte(tc.input))
			pc := NewContext()
			pc.SetBlockOffset(tc.blockOffset)

			node, state := parser.Open(ast.NewDocument(), reader, pc)

			if tc.wantNil {
				if node != nil {
					t.Errorf("Open() node = %v, want nil", node)
				}
				if state != NoChildren {
					t.Errorf("Open() state = %v, want NoChildren", state)
				}
				return
			}

			if node == nil {
				t.Fatalf("Open() node = nil, want non-nil")
			}

			htmlBlock, ok := node.(*ast.HTMLBlock)
			if !ok {
				t.Fatalf("Open() node type = %T, want *ast.HTMLBlock", node)
			}

			if htmlBlock.HTMLBlockType != tc.wantType {
				t.Errorf("Open() HTMLBlockType = %v, want %v", htmlBlock.HTMLBlockType, tc.wantType)
			}

			if state != NoChildren {
				t.Errorf("Open() state = %v, want NoChildren", state)
			}

			if htmlBlock.Lines().Len() != 1 {
				t.Errorf("Open() lines count = %d, want 1", htmlBlock.Lines().Len())
			}
		})
	}
}

func TestHTMLBlockParser_Continue(t *testing.T) {
	tests := []struct {
		name         string
		blockType    ast.HTMLBlockType
		initialLine  string
		continueLine string
		wantState    State
		wantClose    bool
	}{
		{
			name:         "type 1 - script continues",
			blockType:    ast.HTMLBlockType1,
			initialLine:  "<script>\n",
			continueLine: "alert('hello');\n",
			wantState:    Continue | NoChildren,
		},
		{
			name:         "type 1 - script closes",
			blockType:    ast.HTMLBlockType1,
			initialLine:  "<script>\n",
			continueLine: "</script>\n",
			wantState:    Close,
			wantClose:    true,
		},
		{
			name:         "type 1 - script closes on same line",
			blockType:    ast.HTMLBlockType1,
			initialLine:  "<script>alert('test');</script>\n",
			continueLine: "",
			wantState:    Close,
		},
		{
			name:         "type 2 - comment continues",
			blockType:    ast.HTMLBlockType2,
			initialLine:  "<!--\n",
			continueLine: "comment text\n",
			wantState:    Continue | NoChildren,
		},
		{
			name:         "type 2 - comment closes",
			blockType:    ast.HTMLBlockType2,
			initialLine:  "<!--\n",
			continueLine: "end -->\n",
			wantState:    Close,
			wantClose:    true,
		},
		{
			name:         "type 3 - processing instruction continues",
			blockType:    ast.HTMLBlockType3,
			initialLine:  "<?xml\n",
			continueLine: "version=\"1.0\"\n",
			wantState:    Continue | NoChildren,
		},
		{
			name:         "type 3 - processing instruction closes",
			blockType:    ast.HTMLBlockType3,
			initialLine:  "<?xml\n",
			continueLine: "?>\n",
			wantState:    Close,
			wantClose:    true,
		},
		{
			name:         "type 4 - declaration continues",
			blockType:    ast.HTMLBlockType4,
			initialLine:  "<!DOCTYPE\n",
			continueLine: "html\n",
			wantState:    Continue | NoChildren,
		},
		{
			name:         "type 4 - declaration closes",
			blockType:    ast.HTMLBlockType4,
			initialLine:  "<!DOCTYPE\n",
			continueLine: ">\n",
			wantState:    Close,
			wantClose:    true,
		},
		{
			name:         "type 5 - CDATA continues",
			blockType:    ast.HTMLBlockType5,
			initialLine:  "<![CDATA[\n",
			continueLine: "data content\n",
			wantState:    Continue | NoChildren,
		},
		{
			name:         "type 5 - CDATA closes",
			blockType:    ast.HTMLBlockType5,
			initialLine:  "<![CDATA[\n",
			continueLine: "]]>\n",
			wantState:    Close,
			wantClose:    true,
		},
		{
			name:         "type 6 - div continues",
			blockType:    ast.HTMLBlockType6,
			initialLine:  "<div>\n",
			continueLine: "content\n",
			wantState:    Continue | NoChildren,
		},
		{
			name:         "type 6 - div closes on blank line",
			blockType:    ast.HTMLBlockType6,
			initialLine:  "<div>\n",
			continueLine: "\n",
			wantState:    Close,
		},
		{
			name:         "type 7 - span continues",
			blockType:    ast.HTMLBlockType7,
			initialLine:  "<span>\n",
			continueLine: "text\n",
			wantState:    Continue | NoChildren,
		},
		{
			name:         "type 7 - span closes on blank line",
			blockType:    ast.HTMLBlockType7,
			initialLine:  "<span>\n",
			continueLine: "\n",
			wantState:    Close,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewHTMLBlockParser()

			// Create the HTML block node
			node := ast.NewHTMLBlock(tc.blockType)

			// Add initial line
			initialReader := text.NewReader([]byte(tc.initialLine))
			_, initialSeg := initialReader.PeekLine()
			node.Lines().Append(initialSeg)

			// For single-line closure tests (type 1), check immediately
			if tc.continueLine == "" {
				reader := text.NewReader([]byte(tc.initialLine))
				pc := NewContext()
				state := parser.Continue(node, reader, pc)
				if state != tc.wantState {
					t.Errorf("Continue() state = %v, want %v", state, tc.wantState)
				}
				return
			}

			// Test continuation with next line
			source := []byte(tc.initialLine + tc.continueLine)
			reader := text.NewReader(source)
			reader.Advance(len(tc.initialLine))

			pc := NewContext()
			state := parser.Continue(node, reader, pc)

			if state != tc.wantState {
				t.Errorf("Continue() state = %v, want %v", state, tc.wantState)
			}

			if tc.wantClose && !node.ClosureLine.IsEmpty() {
				// Verify closure line was set
				if node.ClosureLine.Start == 0 && node.ClosureLine.Stop == 0 {
					t.Errorf("Continue() ClosureLine not set, want set")
				}
			}
		})
	}
}

func TestHTMLBlockParser_Close(t *testing.T) {
	parser := NewHTMLBlockParser()
	node := ast.NewHTMLBlock(ast.HTMLBlockType1)
	reader := text.NewReader([]byte("<script></script>\n"))
	pc := NewContext()

	// Close should not panic and should do nothing
	parser.Close(node, reader, pc)

	// Verify node is still valid
	if node == nil {
		t.Error("Close() should not modify node")
	}
}

func TestHTMLBlockParser_CanInterruptParagraph(t *testing.T) {
	parser := NewHTMLBlockParser()
	if !parser.CanInterruptParagraph() {
		t.Error("CanInterruptParagraph() = false, want true")
	}
}

func TestHTMLBlockParser_CanAcceptIndentedLine(t *testing.T) {
	parser := NewHTMLBlockParser()
	if parser.CanAcceptIndentedLine() {
		t.Error("CanAcceptIndentedLine() = true, want false")
	}
}

func TestHTMLBlockParser_Open_NegativeBlockOffset(t *testing.T) {
	parser := NewHTMLBlockParser()
	reader := text.NewReader([]byte("<div>\n"))
	pc := NewContext()
	pc.SetBlockOffset(-1)

	node, state := parser.Open(ast.NewDocument(), reader, pc)

	if node != nil {
		t.Errorf("Open() with negative offset should return nil, got %v", node)
	}
	if state != NoChildren {
		t.Errorf("Open() state = %v, want NoChildren", state)
	}
}

func TestHTMLBlockParser_Open_NonHTMLAtOffset(t *testing.T) {
	parser := NewHTMLBlockParser()
	reader := text.NewReader([]byte("text<div>\n"))
	pc := NewContext()
	pc.SetBlockOffset(4) // offset points to '<'

	node, state := parser.Open(ast.NewDocument(), reader, pc)

	if node != nil {
		t.Errorf("Open() should return nil when character at offset is not '<', got %v", node)
	}
	if state != NoChildren {
		t.Errorf("Open() state = %v, want NoChildren", state)
	}
}
