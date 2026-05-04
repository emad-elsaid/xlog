package images

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
)

func TestImagesColumns_Kind(t *testing.T) {
	tests := []struct {
		name     string
		expected ast.NodeKind
	}{
		{
			name:     "returns KindColumns",
			expected: KindColumns,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &imagesColumns{}
			result := node.Kind()
			if result != tt.expected {
				t.Errorf("Kind() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestKindColumns_Uniqueness(t *testing.T) {
	// Verify KindColumns is properly initialized as a unique kind
	if KindColumns == 0 {
		t.Error("KindColumns should be non-zero")
	}

	// Verify it's distinct from standard paragraph kind
	if KindColumns == ast.KindParagraph {
		t.Error("KindColumns should be distinct from KindParagraph")
	}
}

func TestImagesColumns_InheritsFromParagraph(t *testing.T) {
	// Create an imagesColumns node
	node := &imagesColumns{}

	// Verify it can be treated as BaseBlock (inherited from Paragraph)
	// This tests the inheritance relationship
	if node.FirstChild() != nil && node.ChildCount() != 0 {
		// If there are children (shouldn't be for new node), this is unexpected
		t.Error("new imagesColumns should have no children")
	}

	// Test that we can append children (inherited Paragraph behavior)
	img1 := ast.NewImage(ast.NewLink())
	node.AppendChild(node, img1)

	if node.ChildCount() != 1 {
		t.Errorf("ChildCount() = %d, want 1", node.ChildCount())
	}

	if node.FirstChild() != img1 {
		t.Error("FirstChild() should return the appended image")
	}
}
