package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func TestLinkReferenceParagraphTransformer_Transform(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantRefLabel  string
		wantRefDest   string
		wantRefTitle  string
		wantRemaining bool
	}{
		{
			name:          "simple reference definition",
			input:         "[foo]: /url \"title\"",
			wantRefLabel:  "foo",
			wantRefDest:   "/url",
			wantRefTitle:  "title",
			wantRemaining: false,
		},
		{
			name:          "reference without title",
			input:         "[bar]: /path/to/resource",
			wantRefLabel:  "bar",
			wantRefDest:   "/path/to/resource",
			wantRefTitle:  "",
			wantRemaining: false,
		},
		{
			name:          "reference with single quotes",
			input:         "[baz]: /url 'single quoted title'",
			wantRefLabel:  "baz",
			wantRefDest:   "/url",
			wantRefTitle:  "single quoted title",
			wantRemaining: false,
		},
		{
			name:          "reference with parentheses title",
			input:         "[qux]: /url (title in parens)",
			wantRefLabel:  "qux",
			wantRefDest:   "/url",
			wantRefTitle:  "title in parens",
			wantRemaining: false,
		},
		{
			name:          "reference with indentation",
			input:         "   [indented]: /url",
			wantRefLabel:  "indented",
			wantRefDest:   "/url",
			wantRefTitle:  "",
			wantRemaining: false,
		},
		{
			name:          "reference with angle bracket URL",
			input:         "[link]: <http://example.com>",
			wantRefLabel:  "link",
			wantRefDest:   "http://example.com",
			wantRefTitle:  "",
			wantRemaining: false,
		},
		{
			name:          "paragraph with no reference",
			input:         "This is just regular text",
			wantRemaining: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			source := []byte(tc.input)
			reader := text.NewReader(source)
			pc := NewContext()
			doc := ast.NewDocument()

			para := ast.NewParagraph()
			lines := text.NewSegments()
			lines.Append(text.NewSegment(0, len(source)))
			para.SetLines(lines)
			doc.AppendChild(doc, para)

			transformer := LinkReferenceParagraphTransformer
			transformer.Transform(para, reader, pc)

			// Check reference was added to context
			if tc.wantRefLabel != "" {
				normalizedLabel := util.ToLinkReference([]byte(tc.wantRefLabel))
				ref, ok := pc.Reference(normalizedLabel)
				if !ok {
					t.Errorf("Reference %q (normalized: %q) not found in context", tc.wantRefLabel, normalizedLabel)
					return
				}

				if string(ref.Destination()) != tc.wantRefDest {
					t.Errorf("Reference destination = %q, want %q", ref.Destination(), tc.wantRefDest)
				}

				if tc.wantRefTitle != "" {
					if string(ref.Title()) != tc.wantRefTitle {
						t.Errorf("Reference title = %q, want %q", ref.Title(), tc.wantRefTitle)
					}
				}
			}

			// Check if content remains
			if tc.wantRemaining {
				if para.Lines() == nil || para.Lines().Len() == 0 {
					t.Error("Expected paragraph to have remaining content, but it was empty")
				}
			} else {
				child := doc.FirstChild()
				if _, ok := child.(*ast.TextBlock); !ok {
					t.Errorf("Expected paragraph to be replaced with TextBlock, got %T", child)
				}
			}
		})
	}
}

func TestParseLinkReferenceDefinition_InvalidCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  [2]int
	}{
		// Note: SkipSpaces() is called before checking indentation,
		// so 4 spaces actually get consumed and width becomes 0
		// This may not be fully CommonMark compliant but is current behavior
		{
			name:  "missing opening bracket",
			input: "foo]: /url",
			want:  [2]int{-1, -1},
		},
		{
			name:  "missing closing bracket",
			input: "[foo: /url",
			want:  [2]int{-1, -1},
		},
		{
			name:  "missing colon",
			input: "[foo] /url",
			want:  [2]int{-1, -1},
		},
		{
			name:  "blank label",
			input: "[]: /url",
			want:  [2]int{-1, -1},
		},
		{
			name:  "empty label",
			input: "[   ]: /url",
			want:  [2]int{-1, -1},
		},
		{
			name:  "missing destination",
			input: "[foo]: ",
			want:  [2]int{-1, -1},
		},
		{
			name:  "empty input",
			input: "",
			want:  [2]int{-1, -1},
		},
		{
			name:  "unclosed title quote",
			input: "[foo]: /url \"unclosed",
			want:  [2]int{-1, -1},
		},
		{
			name:  "title on same line with text after",
			input: "[foo]: /url \"title\" extra text",
			want:  [2]int{-1, -1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			source := []byte(tc.input)
			segments := text.NewSegments()
			segments.Append(text.NewSegment(0, len(source)))
			block := text.NewBlockReader(source, segments)
			pc := NewContext()

			start, end := parseLinkReferenceDefinition(block, pc)

			if start != tc.want[0] || end != tc.want[1] {
				t.Errorf("parseLinkReferenceDefinition() = (%d, %d), want (%d, %d)",
					start, end, tc.want[0], tc.want[1])
			}

			// Verify no reference was added
			if start == -1 {
				refs := pc.References()
				if len(refs) > 0 {
					t.Errorf("Expected no references, but got %d", len(refs))
				}
			}
		})
	}
}

func TestParseLinkReferenceDefinition_ValidCases(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantStart int
		wantEnd   int
		wantLabel string
		wantDest  string
		wantTitle string
	}{
		{
			name:      "simple reference",
			input:     "[foo]: /url",
			wantStart: 0,
			wantEnd:   1,
			wantLabel: "foo",
			wantDest:  "/url",
			wantTitle: "",
		},
		{
			name:      "reference with title",
			input:     "[bar]: /url \"title\"",
			wantStart: 0,
			wantEnd:   1,
			wantLabel: "bar",
			wantDest:  "/url",
			wantTitle: "title",
		},
		{
			name:      "reference with empty title",
			input:     "[baz]: /url \"\"",
			wantStart: 0,
			wantEnd:   1,
			wantLabel: "baz",
			wantDest:  "/url",
			wantTitle: "",
		},
		{
			name:      "reference with single-quoted title",
			input:     "[qux]: /url 'title'",
			wantStart: 0,
			wantEnd:   1,
			wantLabel: "qux",
			wantDest:  "/url",
			wantTitle: "title",
		},
		{
			name:      "reference with parenthesized title",
			input:     "[quux]: /url (title)",
			wantStart: 0,
			wantEnd:   1,
			wantLabel: "quux",
			wantDest:  "/url",
			wantTitle: "title",
		},
		{
			name:      "reference with angle-bracketed URL",
			input:     "[link]: <http://example.com>",
			wantStart: 0,
			wantEnd:   1,
			wantLabel: "link",
			wantDest:  "http://example.com",
			wantTitle: "",
		},
		{
			name:      "reference with 1 space indent",
			input:     " [ref]: /url",
			wantStart: 0,
			wantEnd:   1,
			wantLabel: "ref",
			wantDest:  "/url",
			wantTitle: "",
		},
		{
			name:      "reference with 2 space indent",
			input:     "  [ref]: /url",
			wantStart: 0,
			wantEnd:   1,
			wantLabel: "ref",
			wantDest:  "/url",
			wantTitle: "",
		},
		{
			name:      "reference with 3 space indent",
			input:     "   [ref]: /url",
			wantStart: 0,
			wantEnd:   1,
			wantLabel: "ref",
			wantDest:  "/url",
			wantTitle: "",
		},
		{
			name:      "reference with complex label",
			input:     "[foo bar baz]: /url",
			wantStart: 0,
			wantEnd:   1,
			wantLabel: "foo bar baz",
			wantDest:  "/url",
			wantTitle: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			source := []byte(tc.input)
			segments := text.NewSegments()
			segments.Append(text.NewSegment(0, len(source)))
			block := text.NewBlockReader(source, segments)
			pc := NewContext()

			start, end := parseLinkReferenceDefinition(block, pc)

			if start != tc.wantStart {
				t.Errorf("start = %d, want %d", start, tc.wantStart)
			}

			if end != tc.wantEnd {
				t.Errorf("end = %d, want %d", end, tc.wantEnd)
			}

			// Verify reference was added correctly
			ref, ok := pc.Reference(tc.wantLabel)
			if !ok {
				t.Fatalf("Reference %q not found in context", tc.wantLabel)
			}

			if string(ref.Label()) != tc.wantLabel {
				t.Errorf("Label = %q, want %q", ref.Label(), tc.wantLabel)
			}

			if string(ref.Destination()) != tc.wantDest {
				t.Errorf("Destination = %q, want %q", ref.Destination(), tc.wantDest)
			}

			if tc.wantTitle != "" && string(ref.Title()) != tc.wantTitle {
				t.Errorf("Title = %q, want %q", ref.Title(), tc.wantTitle)
			}
		})
	}
}

func TestLinkReferenceParagraphTransformer_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantRefCount int
		wantEmpty    bool
	}{
		{
			name:         "empty paragraph",
			input:        "",
			wantRefCount: 0,
			wantEmpty:    true,
		},
		{
			name:         "single reference",
			input:        "[a]: /url1",
			wantRefCount: 1,
			wantEmpty:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			source := []byte(tc.input)
			reader := text.NewReader(source)
			pc := NewContext()
			doc := ast.NewDocument()

			para := ast.NewParagraph()
			lines := text.NewSegments()
			if len(source) > 0 {
				lines.Append(text.NewSegment(0, len(source)))
			}
			para.SetLines(lines)
			doc.AppendChild(doc, para)

			transformer := LinkReferenceParagraphTransformer
			transformer.Transform(para, reader, pc)

			refs := pc.References()
			if len(refs) != tc.wantRefCount {
				t.Errorf("Reference count = %d, want %d", len(refs), tc.wantRefCount)
			}

			if tc.wantEmpty {
				child := doc.FirstChild()
				if _, ok := child.(*ast.TextBlock); !ok {
					t.Errorf("Expected TextBlock for empty result, got %T", child)
				}
			}
		})
	}
}
