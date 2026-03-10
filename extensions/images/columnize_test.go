package images

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestContainsOnlyImages(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *ast.Paragraph
		expected bool
	}{
		{
			name: "single image returns false",
			setup: func() *ast.Paragraph {
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				return p
			},
			expected: false,
		},
		{
			name: "two images returns true",
			setup: func() *ast.Paragraph {
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				return p
			},
			expected: true,
		},
		{
			name: "images with soft line breaks returns true",
			setup: func() *ast.Paragraph {
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				softBreak := ast.NewText()
				softBreak.SetSoftLineBreak(true)
				p.AppendChild(p, softBreak)
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				return p
			},
			expected: true,
		},
		{
			name: "images with text content returns false",
			setup: func() *ast.Paragraph {
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				textNode := ast.NewText()
				textNode.SetSoftLineBreak(false)
				p.AppendChild(p, textNode)
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				return p
			},
			expected: false,
		},
		{
			name: "image with link returns false",
			setup: func() *ast.Paragraph {
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				p.AppendChild(p, ast.NewLink())
				return p
			},
			expected: false,
		},
		{
			name: "empty paragraph returns false",
			setup: func() *ast.Paragraph {
				return ast.NewParagraph()
			},
			expected: false,
		},
		{
			name: "three images returns true",
			setup: func() *ast.Paragraph {
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				return p
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.setup()
			result := containsOnlyImages(p)
			if result != tt.expected {
				t.Errorf("containsOnlyImages() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRemoveBreaks(t *testing.T) {
	tests := []struct {
		name           string
		setup          func() *ast.Paragraph
		expectedChildren int
	}{
		{
			name: "removes all text nodes",
			setup: func() *ast.Paragraph {
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				p.AppendChild(p, ast.NewText())
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				p.AppendChild(p, ast.NewText())
				return p
			},
			expectedChildren: 2, // only images remain
		},
		{
			name: "handles paragraph with no breaks",
			setup: func() *ast.Paragraph {
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				return p
			},
			expectedChildren: 2,
		},
		{
			name: "handles paragraph with only breaks",
			setup: func() *ast.Paragraph {
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewText())
				p.AppendChild(p, ast.NewText())
				return p
			},
			expectedChildren: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.setup()
			removeBreaks(p)
			if p.ChildCount() != tt.expectedChildren {
				t.Errorf("removeBreaks() resulted in %d children, want %d", p.ChildCount(), tt.expectedChildren)
			}
		})
	}
}

func TestReplaceWithColumns(t *testing.T) {
	// Create a document structure to test replacement
	doc := ast.NewDocument()
	section := ast.NewHeading(1)
	p := ast.NewParagraph()
	p.AppendChild(p, ast.NewImage(ast.NewLink()))
	p.AppendChild(p, ast.NewImage(ast.NewLink()))
	
	section.AppendChild(section, p)
	doc.AppendChild(doc, section)

	// Verify initial state
	if p.Parent() != section {
		t.Fatal("paragraph parent should be section")
	}

	// Replace paragraph with columns
	replaceWithColumns(p)

	// Verify the replacement
	replaced := section.FirstChild()
	if replaced == nil {
		t.Fatal("section should have a child after replacement")
	}

	// Check that it's now an imagesColumns node
	if _, ok := replaced.(*imagesColumns); !ok {
		t.Errorf("replacement should be imagesColumns type, got %T", replaced)
	}

	// Verify it has the same children (images)
	if replaced.ChildCount() != 2 {
		t.Errorf("replaced node should have 2 children, got %d", replaced.ChildCount())
	}
}

func TestColumnizeImagesParagraph_Transform(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *ast.Document
		validate func(*testing.T, *ast.Document)
	}{
		{
			name: "transforms paragraph with multiple images",
			setup: func() *ast.Document {
				doc := ast.NewDocument()
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				doc.AppendChild(doc, p)
				return doc
			},
			validate: func(t *testing.T, doc *ast.Document) {
				child := doc.FirstChild()
				if child == nil {
					t.Fatal("document should have a child")
				}
				if _, ok := child.(*imagesColumns); !ok {
					t.Errorf("child should be imagesColumns, got %T", child)
				}
			},
		},
		{
			name: "leaves single image paragraph unchanged",
			setup: func() *ast.Document {
				doc := ast.NewDocument()
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				doc.AppendChild(doc, p)
				return doc
			},
			validate: func(t *testing.T, doc *ast.Document) {
				child := doc.FirstChild()
				if child == nil {
					t.Fatal("document should have a child")
				}
				if _, ok := child.(*ast.Paragraph); !ok {
					t.Errorf("child should still be Paragraph, got %T", child)
				}
			},
		},
		{
			name: "leaves mixed content paragraph unchanged",
			setup: func() *ast.Document {
				doc := ast.NewDocument()
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				textNode := ast.NewText()
				textNode.SetSoftLineBreak(false)
				p.AppendChild(p, textNode)
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				doc.AppendChild(doc, p)
				return doc
			},
			validate: func(t *testing.T, doc *ast.Document) {
				child := doc.FirstChild()
				if child == nil {
					t.Fatal("document should have a child")
				}
				if _, ok := child.(*ast.Paragraph); !ok {
					t.Errorf("child should still be Paragraph, got %T", child)
				}
			},
		},
		{
			name: "handles nested structure",
			setup: func() *ast.Document {
				doc := ast.NewDocument()
				section := ast.NewHeading(1)
				p := ast.NewParagraph()
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				p.AppendChild(p, ast.NewImage(ast.NewLink()))
				section.AppendChild(section, p)
				doc.AppendChild(doc, section)
				return doc
			},
			validate: func(t *testing.T, doc *ast.Document) {
				section := doc.FirstChild()
				if section == nil {
					t.Fatal("document should have section")
				}
				child := section.FirstChild()
				if child == nil {
					t.Fatal("section should have a child")
				}
				if _, ok := child.(*imagesColumns); !ok {
					t.Errorf("nested paragraph should be transformed to imagesColumns, got %T", child)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := tt.setup()
			transformer := columnizeImagesParagraph{}
			transformer.Transform(doc, text.NewReader([]byte{}), nil)
			tt.validate(t, doc)
		})
	}
}
