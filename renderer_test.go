package xlog

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestFindAllInAST_EmptyDocument(t *testing.T) {
	// Parse empty markdown
	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte("")))
	
	// Try to find paragraphs (should be none)
	paragraphs := FindAllInAST[*ast.Paragraph](doc)
	
	if len(paragraphs) != 0 {
		t.Errorf("Expected 0 paragraphs in empty document, got %d", len(paragraphs))
	}
}

func TestFindAllInAST_SimpleParagraphs(t *testing.T) {
	// Parse markdown with multiple paragraphs
	content := `First paragraph.

Second paragraph.

Third paragraph.`
	
	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))
	
	paragraphs := FindAllInAST[*ast.Paragraph](doc)
	
	if len(paragraphs) != 3 {
		t.Errorf("Expected 3 paragraphs, got %d", len(paragraphs))
	}
}

func TestFindAllInAST_Headings(t *testing.T) {
	content := `# Heading 1

Some text.

## Heading 2

More text.

### Heading 3`
	
	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))
	
	headings := FindAllInAST[*ast.Heading](doc)
	
	if len(headings) != 3 {
		t.Errorf("Expected 3 headings, got %d", len(headings))
	}
	
	// Verify heading levels
	if headings[0].Level != 1 {
		t.Errorf("Expected level 1, got %d", headings[0].Level)
	}
	if headings[1].Level != 2 {
		t.Errorf("Expected level 2, got %d", headings[1].Level)
	}
	if headings[2].Level != 3 {
		t.Errorf("Expected level 3, got %d", headings[2].Level)
	}
}

func TestFindAllInAST_Links(t *testing.T) {
	content := `Here is [link one](https://example.com).

And [link two](https://example.org).

And [link three](https://example.net).`
	
	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))
	
	links := FindAllInAST[*ast.Link](doc)
	
	if len(links) != 3 {
		t.Errorf("Expected 3 links, got %d", len(links))
	}
	
	// Verify destinations
	expectedDests := []string{
		"https://example.com",
		"https://example.org",
		"https://example.net",
	}
	
	for i, link := range links {
		dest := string(link.Destination)
		if dest != expectedDests[i] {
			t.Errorf("Expected destination '%s', got '%s'", expectedDests[i], dest)
		}
	}
}

func TestFindAllInAST_CodeBlocks(t *testing.T) {
	content := "```go\nfunc main() {}\n```\n\nSome text.\n\n```python\nprint('hello')\n```"
	
	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))
	
	codeBlocks := FindAllInAST[*ast.FencedCodeBlock](doc)
	
	if len(codeBlocks) != 2 {
		t.Errorf("Expected 2 code blocks, got %d", len(codeBlocks))
	}
	
	// Verify languages
	if string(codeBlocks[0].Language([]byte(content))) != "go" {
		t.Errorf("Expected language 'go', got '%s'", codeBlocks[0].Language([]byte(content)))
	}
	if string(codeBlocks[1].Language([]byte(content))) != "python" {
		t.Errorf("Expected language 'python', got '%s'", codeBlocks[1].Language([]byte(content)))
	}
}

func TestFindAllInAST_Lists(t *testing.T) {
	content := `- Item 1
- Item 2
- Item 3

Some text.

1. Numbered 1
2. Numbered 2`
	
	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))
	
	lists := FindAllInAST[*ast.List](doc)
	
	if len(lists) != 2 {
		t.Errorf("Expected 2 lists, got %d", len(lists))
	}
	
	// First should be unordered
	if lists[0].IsOrdered() {
		t.Error("Expected first list to be unordered")
	}
	
	// Second should be ordered
	if !lists[1].IsOrdered() {
		t.Error("Expected second list to be ordered")
	}
}

func TestFindAllInAST_NilNode(t *testing.T) {
	// Test with nil node
	paragraphs := FindAllInAST[*ast.Paragraph](nil)
	
	if paragraphs != nil {
		t.Errorf("Expected nil result for nil node, got %v", paragraphs)
	}
}

func TestFindAllInAST_NoMatches(t *testing.T) {
	// Parse markdown with no images
	content := "Just some plain text without any images."
	
	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))
	
	images := FindAllInAST[*ast.Image](doc)
	
	if len(images) != 0 {
		t.Errorf("Expected 0 images, got %d", len(images))
	}
}

func TestFindAllInAST_MixedContent(t *testing.T) {
	content := `# Title

First paragraph with [a link](https://example.com).

Second paragraph.

- List item 1
- List item 2

Final paragraph.`
	
	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))
	
	// Count all text nodes
	textNodes := FindAllInAST[*ast.Text](doc)
	
	// Should have multiple text nodes from heading, paragraphs, link, list items
	if len(textNodes) < 5 {
		t.Errorf("Expected at least 5 text nodes, got %d", len(textNodes))
	}
}

func TestFindAllInAST_NestedStructures(t *testing.T) {
	content := `- Item with **bold** text
- Item with *italic* text
- Item with [link](https://example.com)`
	
	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))
	
	// Find all emphasis nodes (bold and italic)
	emphasis := FindAllInAST[*ast.Emphasis](doc)
	
	// Should find at least 2 (bold and italic)
	if len(emphasis) < 2 {
		t.Errorf("Expected at least 2 emphasis nodes, got %d", len(emphasis))
	}
}

func TestFindAllInAST_ComplexDocument(t *testing.T) {
	content := `# Main Title

## Section 1

Paragraph with [link1](https://example.com) and [link2](https://example.org).

### Subsection

More text here.

## Section 2

Final paragraph.`
	
	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))
	
	// Test multiple types
	headings := FindAllInAST[*ast.Heading](doc)
	paragraphs := FindAllInAST[*ast.Paragraph](doc)
	links := FindAllInAST[*ast.Link](doc)
	
	if len(headings) != 4 {
		t.Errorf("Expected 4 headings, got %d", len(headings))
	}
	
	if len(paragraphs) != 3 {
		t.Errorf("Expected 3 paragraphs, got %d", len(paragraphs))
	}
	
	if len(links) != 2 {
		t.Errorf("Expected 2 links, got %d", len(links))
	}
}
