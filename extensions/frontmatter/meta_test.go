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
)

func TestMeta(t *testing.T) {
	markdown := markdown.New(
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
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
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
	markdown := markdown.New(
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
	if err := markdown.Convert([]byte(source), &buf); err != nil {
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
	markdown := markdown.New(
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
	if err := markdown.Convert([]byte(source), &buf, parser.WithContext(context)); err != nil {
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
	markdown := markdown.New(
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
	if err := markdown.Convert([]byte(source), &buf); err != nil {
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
	markdown := markdown.New(
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

	document := markdown.Parser().Parse(text.NewReader([]byte(source)))
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
