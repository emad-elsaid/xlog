package rtl

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestRTLExtensionName(t *testing.T) {
	ext := RTL{}
	expected := "rtl"
	if ext.Name() != expected {
		t.Errorf("Expected name %q, got %q", expected, ext.Name())
	}
}

func TestRTLInit(t *testing.T) {
	// Test that Init doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Init() panicked: %v", r)
		}
	}()

	ext := RTL{}
	ext.Init()
}

func TestAddDirAutoTransform(t *testing.T) {
	tests := []struct {
		name          string
		markdown      string
		expectedNodes []ast.NodeKind
		nodeCount     int
	}{
		{
			name:          "Paragraph with dir=auto",
			markdown:      "This is a simple paragraph.",
			expectedNodes: []ast.NodeKind{ast.KindParagraph},
			nodeCount:     1,
		},
		{
			name:          "Heading with dir=auto",
			markdown:      "# Heading Level 1",
			expectedNodes: []ast.NodeKind{ast.KindHeading},
			nodeCount:     1,
		},
		{
			name:          "Multiple headings",
			markdown:      "# Heading 1\n\n## Heading 2\n\n### Heading 3",
			expectedNodes: []ast.NodeKind{ast.KindHeading},
			nodeCount:     3,
		},
		{
			name:          "List with dir=auto",
			markdown:      "- Item 1\n- Item 2\n- Item 3",
			expectedNodes: []ast.NodeKind{ast.KindList},
			nodeCount:     1,
		},
		{
			name:          "Blockquote with dir=auto",
			markdown:      "> This is a quote\n> on multiple lines",
			expectedNodes: []ast.NodeKind{ast.KindBlockquote, ast.KindParagraph},
			nodeCount:     2, // Blockquote + paragraph inside
		},
		{
			name: "Mixed content",
			markdown: `# Heading

This is a paragraph.

> A blockquote

- List item 1
- List item 2`,
			expectedNodes: []ast.NodeKind{ast.KindHeading, ast.KindParagraph, ast.KindBlockquote, ast.KindList},
			nodeCount:     5, // Heading + 2 paragraphs (one standalone, one in blockquote) + blockquote + list
		},
		{
			name: "Nested blockquote with paragraphs",
			markdown: `> Outer quote
> 
> > Nested quote
> 
> Back to outer`,
			expectedNodes: []ast.NodeKind{ast.KindBlockquote, ast.KindParagraph},
			nodeCount:     5, // 2 blockquotes + 3 paragraphs
		},
		{
			name:          "Ordered list",
			markdown:      "1. First\n2. Second\n3. Third",
			expectedNodes: []ast.NodeKind{ast.KindList},
			nodeCount:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse markdown to create AST
			source := []byte(tt.markdown)
			reader := text.NewReader(source)
			md := markdown.New()
			doc := md.Parser().Parse(reader)

			// Apply the transform - need type assertion
			docNode, ok := doc.(*ast.Document)
			if !ok {
				t.Fatalf("Expected *ast.Document, got %T", doc)
			}

			transformer := addDirAuto{}
			transformer.Transform(docNode, reader, parser.NewContext())

		// Collect nodes that should have dir="auto"
		var nodesWithDir []ast.Node
		_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
			if !entering {
				return ast.WalkContinue, nil
			}

				kind := node.Kind()
				if kind == ast.KindParagraph ||
					kind == ast.KindHeading ||
					kind == ast.KindList ||
					kind == ast.KindBlockquote {
					nodesWithDir = append(nodesWithDir, node)
				}

				return ast.WalkContinue, nil
			})

			// Verify count
			if len(nodesWithDir) != tt.nodeCount {
				t.Errorf("Expected %d nodes with dir attribute, got %d", tt.nodeCount, len(nodesWithDir))
			}

			// Verify each node has dir="auto"
			for _, node := range nodesWithDir {
				dir, ok := node.AttributeString("dir")
				if !ok {
					t.Errorf("Node %v missing 'dir' attribute", node.Kind())
					continue
				}

				dirStr, ok := dir.(string)
				if !ok {
					dirBytes, ok := dir.([]byte)
					if !ok {
						t.Errorf("dir attribute is unexpected type: %T", dir)
						continue
					}
					dirStr = string(dirBytes)
				}

				if dirStr != "auto" {
					t.Errorf("Expected dir='auto', got dir='%s'", dirStr)
				}
			}
		})
	}
}

func TestAddDirAutoIgnoresOtherNodes(t *testing.T) {
	markdownText := "# Heading\n\nParagraph with `code` and **bold** text.\n\n---"
	source := []byte(markdownText)
	reader := text.NewReader(source)
	md := markdown.New()
	doc := md.Parser().Parse(reader)

	docNode, ok := doc.(*ast.Document)
	if !ok {
		t.Fatalf("Expected *ast.Document, got %T", doc)
	}

	transformer := addDirAuto{}
	transformer.Transform(docNode, reader, parser.NewContext())

	// Count nodes that should NOT have dir="auto"
	var codeNodes, emphasisNodes, thematicBreaks int
	_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		kind := node.Kind()
		switch kind {
		case ast.KindCodeSpan:
			codeNodes++
			// Code spans should NOT have dir attribute
			if _, ok := node.AttributeString("dir"); ok {
				t.Errorf("CodeSpan should not have dir attribute")
			}
		case ast.KindEmphasis:
			emphasisNodes++
			// Emphasis should NOT have dir attribute
			if _, ok := node.AttributeString("dir"); ok {
				t.Errorf("Emphasis should not have dir attribute")
			}
		case ast.KindThematicBreak:
			thematicBreaks++
			// ThematicBreak should NOT have dir attribute
			if _, ok := node.AttributeString("dir"); ok {
				t.Errorf("ThematicBreak should not have dir attribute")
			}
		}

		return ast.WalkContinue, nil
	})

	// Verify we actually tested these node types
	if codeNodes == 0 {
		t.Error("Expected to find code nodes in test markdown")
	}
	if emphasisNodes == 0 {
		t.Error("Expected to find emphasis nodes in test markdown")
	}
	if thematicBreaks == 0 {
		t.Error("Expected to find thematic break in test markdown")
	}
}

func TestAddDirAutoWithEmptyDocument(t *testing.T) {
	markdownText := ""
	source := []byte(markdownText)
	reader := text.NewReader(source)
	md := markdown.New()
	doc := md.Parser().Parse(reader)

	docNode, ok := doc.(*ast.Document)
	if !ok {
		t.Fatalf("Expected *ast.Document, got %T", doc)
	}

	transformer := addDirAuto{}

	// Should not panic with empty document
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Transform panicked with empty document: %v", r)
		}
	}()

	transformer.Transform(docNode, reader, parser.NewContext())

	// Count nodes with dir attribute (should be 0)
	var count int
	_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if entering {
			if _, ok := node.AttributeString("dir"); ok {
				count++
			}
		}
		return ast.WalkContinue, nil
	})

	if count != 0 {
		t.Errorf("Expected 0 nodes with dir attribute in empty document, got %d", count)
	}
}

func TestAddDirAutoWithRTLContent(t *testing.T) {
	// Test with actual RTL (Arabic) content
	markdownText := "# مرحبا\n\nهذا نص عربي."
	source := []byte(markdownText)
	reader := text.NewReader(source)
	md := markdown.New()
	doc := md.Parser().Parse(reader)

	docNode, ok := doc.(*ast.Document)
	if !ok {
		t.Fatalf("Expected *ast.Document, got %T", doc)
	}

	transformer := addDirAuto{}
	transformer.Transform(docNode, reader, parser.NewContext())

	// Verify heading and paragraph have dir="auto"
	var headingCount, paragraphCount int
	_ = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		kind := node.Kind()
		if kind == ast.KindHeading {
			headingCount++
			dir, ok := node.AttributeString("dir")
			if !ok {
				t.Error("Heading missing dir attribute")
			} else {
				dirStr, ok := dir.(string)
				if !ok {
					dirBytes, ok := dir.([]byte)
					if !ok {
						t.Errorf("dir attribute is unexpected type: %T", dir)
					} else {
						dirStr = string(dirBytes)
					}
				}
				if dirStr != "auto" {
					t.Errorf("Expected dir='auto', got '%s'", dirStr)
				}
			}
		} else if kind == ast.KindParagraph {
			paragraphCount++
			dir, ok := node.AttributeString("dir")
			if !ok {
				t.Error("Paragraph missing dir attribute")
			} else {
				dirStr, ok := dir.(string)
				if !ok {
					dirBytes, ok := dir.([]byte)
					if !ok {
						t.Errorf("dir attribute is unexpected type: %T", dir)
					} else {
						dirStr = string(dirBytes)
					}
				}
				if dirStr != "auto" {
					t.Errorf("Expected dir='auto', got '%s'", dirStr)
				}
			}
		}

		return ast.WalkContinue, nil
	})

	if headingCount != 1 {
		t.Errorf("Expected 1 heading, got %d", headingCount)
	}
	if paragraphCount != 1 {
		t.Errorf("Expected 1 paragraph, got %d", paragraphCount)
	}
}
