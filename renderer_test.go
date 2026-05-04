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

func BenchmarkFindInAST_SmallDocument(b *testing.B) {
	content := `# Heading

Simple paragraph with [a link](https://example.com).`

	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindInAST[*ast.Link](doc)
	}
}

func BenchmarkFindInAST_MediumDocument(b *testing.B) {
	content := `# Main Heading

## Section 1

Paragraph with some **bold** and *italic* text.

[Link 1](https://example.com)

## Section 2

Another paragraph with [link 2](https://example.org).

- List item 1
- List item 2
- List item 3

Final paragraph.`

	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindInAST[*ast.Link](doc)
	}
}

func BenchmarkFindInAST_LargeDocument(b *testing.B) {
	// Generate large document with many elements
	content := `# Main Title

## Introduction

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Here's [link 1](https://example.com).

### Background

More text with [link 2](https://example.org) and some **bold content**.

## Methods

1. First method
2. Second method
3. Third method

Paragraph explaining the methods with [link 3](https://example.net).

### Detailed Approach

- Point one with *emphasis*
- Point two with **strong emphasis**
- Point three with [link 4](https://example.co.uk)

Code example here.

## Results

Multiple paragraphs here describing results. [Link 5](https://example.io).

Another paragraph with more details and [link 6](https://example.dev).

## Conclusion

Final thoughts with [link 7](https://example.com/conclusion).`

	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FindInAST[*ast.Link](doc)
	}
}

func BenchmarkFindAllInAST_SmallDocument(b *testing.B) {
	content := `# Heading

Paragraph with [link 1](https://example.com) and [link 2](https://example.org).`

	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindAllInAST[*ast.Link](doc)
	}
}

func BenchmarkFindAllInAST_MediumDocument(b *testing.B) {
	content := `# Main Heading

## Section 1

Paragraph with [link 1](https://example.com) and [link 2](https://example.org).

## Section 2

Another paragraph with [link 3](https://example.net) and [link 4](https://example.io).

- List with [link 5](https://example.dev)
- Another item with [link 6](https://example.co.uk)

Final paragraph with [link 7](https://example.com/final).`

	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindAllInAST[*ast.Link](doc)
	}
}

func BenchmarkFindAllInAST_LargeDocument(b *testing.B) {
	// Generate document with many links
	content := `# Main Title

## Introduction

Lorem ipsum [link 1](https://example.com) dolor sit amet.

### Background

More text with [link 2](https://example.org) and [link 3](https://example.net).

## Methods

1. First [link 4](https://example.io)
2. Second [link 5](https://example.dev)
3. Third [link 6](https://example.co.uk)

Paragraph with [link 7](https://example.com/methods).

### Detailed Approach

- Point [link 8](https://example.com/p1)
- Point [link 9](https://example.com/p2)
- Point [link 10](https://example.com/p3)

## Results

Text with [link 11](https://example.com/r1).
More [link 12](https://example.com/r2).
And [link 13](https://example.com/r3).

## Discussion

Analysis with [link 14](https://example.com/d1).
Further [link 15](https://example.com/d2).

## Conclusion

Final [link 16](https://example.com/conclusion).`

	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FindAllInAST[*ast.Link](doc)
	}
}

func BenchmarkFindAllInAST_MultipleTypes(b *testing.B) {
	content := `# Main Title

## Section 1

Paragraph with [link](https://example.com) and **bold** text.

### Subsection

More text with *italic* here.

## Section 2

- List item 1
- List item 2
- List item 3

Final paragraph.`

	md := MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))

	b.Run("FindHeadings", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FindAllInAST[*ast.Heading](doc)
		}
	})

	b.Run("FindParagraphs", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FindAllInAST[*ast.Paragraph](doc)
		}
	})

	b.Run("FindLinks", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FindAllInAST[*ast.Link](doc)
		}
	})

	b.Run("FindEmphasis", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FindAllInAST[*ast.Emphasis](doc)
		}
	})

	b.Run("FindLists", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FindAllInAST[*ast.List](doc)
		}
	})
}
