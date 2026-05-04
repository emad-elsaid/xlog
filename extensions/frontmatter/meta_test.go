package frontmatter

import (
	"bytes"
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/extension"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
	"gopkg.in/yaml.v2"
)

func TestMeta(t *testing.T) {
	md := markdown.New(
		markdown.WithExtensions(
			Meta,
		),
	)
	source := `---
Title: goldmark-meta
Summary: Add YAML metadata to the document
Tags:
    - markdown
    - goldmark
---

# Hello goldmark-meta
`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	metaData := Get(context)
	title := metaData["Title"]
	s, ok := title.(string)
	if !ok {
		t.Error("Title not found in meta data or is not a string")
	}
	if s != "goldmark-meta" {
		t.Errorf("Title must be %s, but got %v", "goldmark-meta", s)
	}
	if buf.String() != "<h1>Hello goldmark-meta</h1>\n" {
		t.Errorf("should render '<h1>Hello goldmark-meta</h1>', but '%s'", buf.String())
	}
	tags, ok := metaData["Tags"].([]any)
	if !ok {
		t.Error("Tags not found in meta data or is not a slice")
	}
	if len(tags) != 2 {
		t.Error("Tags must be a slice that has 2 elements")
	}
	if tags[0] != "markdown" {
		t.Errorf("Tag#1 must be 'markdown', but got %s", tags[0])
	}
	if tags[1] != "goldmark" {
		t.Errorf("Tag#2 must be 'goldmark', but got %s", tags[1])
	}
}

func TestMetaTable(t *testing.T) {
	md := markdown.New(
		markdown.WithExtensions(
			New(WithTable()),
		),
		markdown.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(extension.NewTableHTMLRenderer(), 500),
			),
		),
	)
	source := `---
Title: goldmark-meta
Summary: Add YAML metadata to the document
Tags:
    - markdown
    - goldmark
---

# Hello goldmark-meta
`

	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		panic(err)
	}
	if buf.String() != `<table>
<thead>
<tr>
<th>Title</th>
<th>Summary</th>
<th>Tags</th>
</tr>
</thead>
<tbody>
<tr>
<td>goldmark-meta</td>
<td>Add YAML metadata to the document</td>
<td>[markdown goldmark]</td>
</tr>
</tbody>
</table>
<h1>Hello goldmark-meta</h1>
` {
		t.Error("invalid table output")
	}
}

func TestMetaError(t *testing.T) {
	md := markdown.New(
		markdown.WithExtensions(
			New(WithTable()),
		),
	)
	source := `---
Title: goldmark-meta
Summary: Add YAML metadata to the document
Tags:
  - : {
  }
    - markdown
    - goldmark
---

# Hello goldmark-meta
`

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}
	if buf.String() != `Title: goldmark-meta
Summary: Add YAML metadata to the document
Tags:
  - : {
  }
    - markdown
    - goldmark
<!-- yaml: line 3: did not find expected key -->
<h1>Hello goldmark-meta</h1>
` {
		t.Error("invalid error output")
	}

	v, err := TryGet(context)
	if err == nil {
		t.Error("error should not be nil")
	}
	if v != nil {
		t.Error("data should be nil when there are errors")
	}
}

func TestMetaTableWithBlankline(t *testing.T) {
	md := markdown.New(
		markdown.WithExtensions(
			New(WithTable()),
		),
		markdown.WithRendererOptions(
			renderer.WithNodeRenderers(
				util.Prioritized(extension.NewTableHTMLRenderer(), 500),
			),
		),
	)
	source := `---
Title: goldmark-meta
Summary: Add YAML metadata to the document

# comments
Tags:
    - markdown
    - goldmark
---

# Hello goldmark-meta
`

	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		panic(err)
	}
	if buf.String() != `<table>
<thead>
<tr>
<th>Title</th>
<th>Summary</th>
<th>Tags</th>
</tr>
</thead>
<tbody>
<tr>
<td>goldmark-meta</td>
<td>Add YAML metadata to the document</td>
<td>[markdown goldmark]</td>
</tr>
</tbody>
</table>
<h1>Hello goldmark-meta</h1>
` {
		t.Error("invalid table output")
	}
}

func TestMetaStoreInDocument(t *testing.T) {
	md := markdown.New(
		markdown.WithExtensions(
			New(
				WithStoresInDocument(),
			),
		),
	)
	source := `---
Title: goldmark-meta
Summary: Add YAML metadata to the document
Tags:
    - markdown
    - goldmark
---
`

	document := md.Parser().Parse(text.NewReader([]byte(source)))
	metaData := document.OwnerDocument().Meta()
	title := metaData["Title"]
	s, ok := title.(string)
	if !ok {
		t.Error("Title not found in meta data or is not a string")
	}
	if s != "goldmark-meta" {
		t.Errorf("Title must be %s, but got %v", "goldmark-meta", s)
	}
	tags, ok := metaData["Tags"].([]any)
	if !ok {
		t.Error("Tags not found in meta data or is not a slice")
	}
	if len(tags) != 2 {
		t.Error("Tags must be a slice that has 2 elements")
	}
	if tags[0] != "markdown" {
		t.Errorf("Tag#1 must be 'markdown', but got %s", tags[0])
	}
	if tags[1] != "goldmark" {
		t.Errorf("Tag#2 must be 'goldmark', but got %s", tags[1])
	}
}

func TestTryGetItems(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectError bool
		expectNil   bool
		validate    func(t *testing.T, items yaml.MapSlice)
	}{
		{
			name: "valid frontmatter with preserved order",
			source: `---
Title: Test Page
Author: John Doe
Tags:
  - golang
  - testing
---

# Content`,
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, items yaml.MapSlice) {
				if len(items) != 3 {
					t.Errorf("expected 3 items, got %d", len(items))
				}
				if items[0].Key != "Title" {
					t.Errorf("expected first key to be Title, got %v", items[0].Key)
				}
				if items[1].Key != "Author" {
					t.Errorf("expected second key to be Author, got %v", items[1].Key)
				}
				if items[2].Key != "Tags" {
					t.Errorf("expected third key to be Tags, got %v", items[2].Key)
				}
			},
		},
		{
			name: "no frontmatter",
			source: `# Just a heading

Regular content without frontmatter.`,
			expectError: false,
			expectNil:   true,
			validate:    nil,
		},
		{
			name: "invalid YAML syntax",
			source: `---
Title: Test
Invalid: { unclosed bracket
Tags:
  - test
---

# Content`,
			expectError: true,
			expectNil:   true,
			validate:    nil,
		},
		{
			name: "empty frontmatter",
			source: `---
---

# Content`,
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, items yaml.MapSlice) {
				if len(items) != 0 {
					t.Errorf("expected 0 items for empty frontmatter, got %d", len(items))
				}
			},
		},
		{
			name: "complex nested structure",
			source: `---
Title: Complex
Metadata:
  nested:
    value: test
  list:
    - item1
    - item2
---

# Content`,
			expectError: false,
			expectNil:   false,
			validate: func(t *testing.T, items yaml.MapSlice) {
				if len(items) != 2 {
					t.Errorf("expected 2 items, got %d", len(items))
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			md := markdown.New(
				markdown.WithExtensions(Meta),
			)

			var buf bytes.Buffer
			context := parser.NewContext()
			if err := md.Convert([]byte(tc.source), &buf, parser.WithContext(context)); err != nil {
				t.Fatalf("failed to convert markdown: %v", err)
			}

			items, err := TryGetItems(context)

			if tc.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}

			if tc.expectNil && items != nil {
				t.Error("expected nil items but got non-nil")
			}
			if !tc.expectNil && !tc.expectError && items == nil {
				t.Error("expected non-nil items but got nil")
			}

			if tc.validate != nil && items != nil {
				tc.validate(t, items)
			}
		})
	}
}
