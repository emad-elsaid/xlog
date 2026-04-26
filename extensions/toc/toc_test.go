package toc

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtensionName(t *testing.T) {
	ext := Extension{}
	assert.Equal(t, "toc", ext.Name())
}

func TestInspectBasic(t *testing.T) {
	src := []byte(`# Section 1
## Subsection 1.1
## Subsection 1.2
# Section 2
## Subsection 2.1
# Section 3`)

	doc := parseMarkdown(src)
	toc, err := Inspect(doc, src)
	require.NoError(t, err)
	require.NotNil(t, toc)
	require.Len(t, toc.Items, 3)

	// Section 1
	assert.Equal(t, "Section 1", string(toc.Items[0].Title))
	assert.Equal(t, "section-1", string(toc.Items[0].ID))
	require.Len(t, toc.Items[0].Items, 2)
	assert.Equal(t, "Subsection 1.1", string(toc.Items[0].Items[0].Title))
	assert.Equal(t, "subsection-11", string(toc.Items[0].Items[0].ID))
	assert.Equal(t, "Subsection 1.2", string(toc.Items[0].Items[1].Title))
	assert.Equal(t, "subsection-12", string(toc.Items[0].Items[1].ID))

	// Section 2
	assert.Equal(t, "Section 2", string(toc.Items[1].Title))
	assert.Equal(t, "section-2", string(toc.Items[1].ID))
	require.Len(t, toc.Items[1].Items, 1)
	assert.Equal(t, "Subsection 2.1", string(toc.Items[1].Items[0].Title))
	assert.Equal(t, "subsection-21", string(toc.Items[1].Items[0].ID))

	// Section 3
	assert.Equal(t, "Section 3", string(toc.Items[2].Title))
	assert.Equal(t, "section-3", string(toc.Items[2].ID))
	assert.Len(t, toc.Items[2].Items, 0)
}

func TestInspectEmpty(t *testing.T) {
	src := []byte("No headings here, just text.")
	doc := parseMarkdown(src)
	toc, err := Inspect(doc, src)
	require.NoError(t, err)
	require.NotNil(t, toc)
	assert.Len(t, toc.Items, 0)
}

func TestInspectMaxDepth(t *testing.T) {
	src := []byte(`# Foo
## Bar
### Baz
# Quux
## Qux`)

	doc := parseMarkdown(src)
	toc, err := Inspect(doc, src, MaxDepth(1))
	require.NoError(t, err)
	require.Len(t, toc.Items, 2)
	assert.Equal(t, "Foo", string(toc.Items[0].Title))
	assert.Equal(t, "Quux", string(toc.Items[1].Title))

	// MaxDepth(2) should include level 2 headings
	toc, err = Inspect(doc, src, MaxDepth(2))
	require.NoError(t, err)
	require.Len(t, toc.Items, 2)
	assert.Equal(t, "Foo", string(toc.Items[0].Title))
	require.Len(t, toc.Items[0].Items, 1)
	assert.Equal(t, "Bar", string(toc.Items[0].Items[0].Title))
	assert.Equal(t, "Quux", string(toc.Items[1].Title))
	require.Len(t, toc.Items[1].Items, 1)
	assert.Equal(t, "Qux", string(toc.Items[1].Items[0].Title))
}

func TestInspectMinDepth(t *testing.T) {
	src := []byte(`# Foo
## Bar
### Baz
# Quux
## Qux`)

	doc := parseMarkdown(src)
	// MinDepth(3) with Compact shows only level 3+ headings
	toc, err := Inspect(doc, src, MinDepth(3), Compact(true))
	require.NoError(t, err)
	require.Len(t, toc.Items, 1)
	assert.Equal(t, "Baz", string(toc.Items[0].Title))

	// MinDepth(2) with Compact should only include level 2+ headings
	toc, err = Inspect(doc, src, MinDepth(2), Compact(true))
	require.NoError(t, err)
	require.Len(t, toc.Items, 2)
	assert.Equal(t, "Bar", string(toc.Items[0].Title))
	require.Len(t, toc.Items[0].Items, 1)
	assert.Equal(t, "Baz", string(toc.Items[0].Items[0].Title))
	assert.Equal(t, "Qux", string(toc.Items[1].Title))
}

func TestInspectCompact(t *testing.T) {
	src := []byte(`# A
### B
#### C
# D
#### E`)

	doc := parseMarkdown(src)

	// Without compact, should have empty items
	toc, err := Inspect(doc, src)
	require.NoError(t, err)
	require.Len(t, toc.Items, 2)
	assert.Equal(t, "A", string(toc.Items[0].Title))
	require.Len(t, toc.Items[0].Items, 1)
	assert.Equal(t, "", string(toc.Items[0].Items[0].Title)) // Empty item at level 2
	require.Len(t, toc.Items[0].Items[0].Items, 1)
	assert.Equal(t, "B", string(toc.Items[0].Items[0].Items[0].Title))

	// With compact, empty items should be removed
	toc, err = Inspect(doc, src, Compact(true))
	require.NoError(t, err)
	require.Len(t, toc.Items, 2)
	assert.Equal(t, "A", string(toc.Items[0].Title))
	require.Len(t, toc.Items[0].Items, 1)
	assert.Equal(t, "B", string(toc.Items[0].Items[0].Title)) // B promoted
	require.Len(t, toc.Items[0].Items[0].Items, 1)
	assert.Equal(t, "C", string(toc.Items[0].Items[0].Items[0].Title))

	assert.Equal(t, "D", string(toc.Items[1].Title))
	require.Len(t, toc.Items[1].Items, 1)
	assert.Equal(t, "E", string(toc.Items[1].Items[0].Title)) // E promoted
}

func TestInspectCombinedOptions(t *testing.T) {
	src := []byte(`# Level 1
## Level 2
### Level 3
#### Level 4
## Another Level 2`)

	doc := parseMarkdown(src)
	toc, err := Inspect(doc, src, MinDepth(2), MaxDepth(3), Compact(true))
	require.NoError(t, err)
	require.Len(t, toc.Items, 2)
	assert.Equal(t, "Level 2", string(toc.Items[0].Title))
	require.Len(t, toc.Items[0].Items, 1)
	assert.Equal(t, "Level 3", string(toc.Items[0].Items[0].Title))
	assert.Equal(t, "Another Level 2", string(toc.Items[1].Title))
}

func TestInspectSpecialCharacters(t *testing.T) {
	src := []byte(`# Hello & World
## Code: ` + "`example`" + `
### Link [text](url)`)

	doc := parseMarkdown(src)
	toc, err := Inspect(doc, src)
	require.NoError(t, err)
	require.Len(t, toc.Items, 1)
	assert.Equal(t, "Hello & World", string(toc.Items[0].Title))
	assert.Equal(t, "hello--world", string(toc.Items[0].ID))
}

func TestInspectNestedDeep(t *testing.T) {
	src := []byte(`# H1
## H2
### H3
#### H4
##### H5
###### H6`)

	doc := parseMarkdown(src)
	toc, err := Inspect(doc, src)
	require.NoError(t, err)
	require.Len(t, toc.Items, 1)

	// Walk the nested structure
	current := toc.Items[0]
	assert.Equal(t, "H1", string(current.Title))
	require.Len(t, current.Items, 1)

	current = current.Items[0]
	assert.Equal(t, "H2", string(current.Title))
	require.Len(t, current.Items, 1)

	current = current.Items[0]
	assert.Equal(t, "H3", string(current.Title))
	require.Len(t, current.Items, 1)

	current = current.Items[0]
	assert.Equal(t, "H4", string(current.Title))
	require.Len(t, current.Items, 1)

	current = current.Items[0]
	assert.Equal(t, "H5", string(current.Title))
	require.Len(t, current.Items, 1)

	current = current.Items[0]
	assert.Equal(t, "H6", string(current.Title))
	assert.Len(t, current.Items, 0)
}

func TestInspectOptionsString(t *testing.T) {
	assert.Equal(t, "MinDepth(3)", MinDepth(3).(minDepthOption).String())
	assert.Equal(t, "MaxDepth(2)", MaxDepth(2).(maxDepthOption).String())
	assert.Equal(t, "Compact(true)", Compact(true).(compactOption).String())
	assert.Equal(t, "Compact(false)", Compact(false).(compactOption).String())
}

func TestInspectMultipleSectionsAtSameLevel(t *testing.T) {
	// Level 1 headings at same level
	src := []byte(`# First
# Second
# Third`)

	doc := parseMarkdown(src)
	toc, err := Inspect(doc, src)
	require.NoError(t, err)
	require.Len(t, toc.Items, 3)
	assert.Equal(t, "First", string(toc.Items[0].Title))
	assert.Equal(t, "Second", string(toc.Items[1].Title))
	assert.Equal(t, "Third", string(toc.Items[2].Title))
}

func TestInspectSkippedLevels(t *testing.T) {
	// H1 -> H3 (skipping H2)
	src := []byte(`# Title
### Subsection`)

	doc := parseMarkdown(src)
	toc, err := Inspect(doc, src)
	require.NoError(t, err)
	require.Len(t, toc.Items, 1)
	assert.Equal(t, "Title", string(toc.Items[0].Title))
	require.Len(t, toc.Items[0].Items, 1)
	// Empty level 2 item
	assert.Equal(t, "", string(toc.Items[0].Items[0].Title))
	require.Len(t, toc.Items[0].Items[0].Items, 1)
	assert.Equal(t, "Subsection", string(toc.Items[0].Items[0].Items[0].Title))
}

func TestInspectWithZeroDepthOptions(t *testing.T) {
	src := []byte(`# H1
## H2
### H3`)

	doc := parseMarkdown(src)

	// MinDepth(0) should be no limit
	toc, err := Inspect(doc, src, MinDepth(0))
	require.NoError(t, err)
	require.Len(t, toc.Items, 1)
	assert.Equal(t, "H1", string(toc.Items[0].Title))

	// MaxDepth(0) should be no limit
	toc, err = Inspect(doc, src, MaxDepth(0))
	require.NoError(t, err)
	require.Len(t, toc.Items, 1)
	assert.Equal(t, "H1", string(toc.Items[0].Title))
	require.Len(t, toc.Items[0].Items, 1)
	assert.Equal(t, "H2", string(toc.Items[0].Items[0].Title))
	require.Len(t, toc.Items[0].Items[0].Items, 1)
	assert.Equal(t, "H3", string(toc.Items[0].Items[0].Items[0].Title))
}

func TestInspectNegativeDepthOptions(t *testing.T) {
	src := []byte(`# H1
## H2`)

	doc := parseMarkdown(src)

	// Negative values should be no limit
	toc, err := Inspect(doc, src, MinDepth(-1), MaxDepth(-1))
	require.NoError(t, err)
	require.Len(t, toc.Items, 1)
	assert.Equal(t, "H1", string(toc.Items[0].Title))
	require.Len(t, toc.Items[0].Items, 1)
	assert.Equal(t, "H2", string(toc.Items[0].Items[0].Title))
}

// parseMarkdown is a helper to parse markdown with auto heading IDs enabled
func parseMarkdown(src []byte) ast.Node {
	return parser.NewParser(
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithAutoHeadingID(),
	).Parse(text.NewReader(src))
}
