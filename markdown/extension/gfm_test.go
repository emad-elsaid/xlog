package extension

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
	"github.com/emad-elsaid/xlog/markdown/testutil"
)

func TestGFM(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			GFM,
		),
	)

	// Test case combining all GFM features:
	// - Linkify: auto-linking URLs
	// - Table: pipe table syntax
	// - Strikethrough: ~~text~~
	// - TaskList: - [ ] and - [x]
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:       1,
			Markdown: "Visit https://github.com\n\n~~deleted~~ text\n\n- [x] Done\n- [ ] Todo\n\n| A | B |\n|---|---|\n| 1 | 2 |\n",
			Expected: `<p>Visit <a href="https://github.com">https://github.com</a></p>
<p><del>deleted</del> text</p>
<ul>
<li><input checked="" disabled="" type="checkbox"> Done</li>
<li><input disabled="" type="checkbox"> Todo</li>
</ul>
<table>
<thead>
<tr>
<th>A</th>
<th>B</th>
</tr>
</thead>
<tbody>
<tr>
<td>1</td>
<td>2</td>
</tr>
</tbody>
</table>
`,
		},
		t,
	)
}

func TestGFM_Linkify(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			GFM,
		),
	)

	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:       1,
			Markdown: "Check http://example.com and www.github.com",
			Expected: `<p>Check <a href="http://example.com">http://example.com</a> and <a href="http://www.github.com">www.github.com</a></p>
`,
		},
		t,
	)
}

func TestGFM_Table(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			GFM,
		),
	)

	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:       1,
			Markdown: "| Header 1 | Header 2 |\n|----------|----------|\n| Cell 1   | Cell 2   |",
			Expected: `<table>
<thead>
<tr>
<th>Header 1</th>
<th>Header 2</th>
</tr>
</thead>
<tbody>
<tr>
<td>Cell 1</td>
<td>Cell 2</td>
</tr>
</tbody>
</table>
`,
		},
		t,
	)
}

func TestGFM_Strikethrough(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			GFM,
		),
	)

	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:       1,
			Markdown: "This is ~~wrong~~ correct.",
			Expected: `<p>This is <del>wrong</del> correct.</p>
`,
		},
		t,
	)
}

func TestGFM_TaskList(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			GFM,
		),
	)

	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:       1,
			Markdown: "- [x] Completed task\n- [ ] Incomplete task",
			Expected: `<ul>
<li><input checked="" disabled="" type="checkbox"> Completed task</li>
<li><input disabled="" type="checkbox"> Incomplete task</li>
</ul>
`,
		},
		t,
	)
}
