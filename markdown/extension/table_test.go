package extension

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/ast"
	east "github.com/emad-elsaid/xlog/markdown/extension/ast"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer/html"
	"github.com/emad-elsaid/xlog/markdown/testutil"
	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func TestTable(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
			html.WithXHTML(),
		),
		markdown.WithExtensions(
			Table,
		),
	)
	testutil.DoTestCaseFile(md, "_test/table.txt", t, testutil.ParseCliCaseArg()...)
}

func TestTableWithAlignDefault(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			NewTable(
				WithTableCellAlignMethod(TableCellAlignDefault),
			),
		),
	)
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          1,
			Description: "Cell with TableCellAlignDefault and XHTML should be rendered as an align attribute",
			Markdown: `
| abc | defghi |
:-: | -----------:
bar | baz
`,
			Expected: `<table>
<thead>
<tr>
<th align="center">abc</th>
<th align="right">defghi</th>
</tr>
</thead>
<tbody>
<tr>
<td align="center">bar</td>
<td align="right">baz</td>
</tr>
</tbody>
</table>`,
		},
		t,
	)

	md = markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			NewTable(
				WithTableCellAlignMethod(TableCellAlignDefault),
			),
		),
	)
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          2,
			Description: "Cell with TableCellAlignDefault and HTML5 should be rendered as a style attribute",
			Markdown: `
| abc | defghi |
:-: | -----------:
bar | baz
`,
			Expected: `<table>
<thead>
<tr>
<th style="text-align:center">abc</th>
<th style="text-align:right">defghi</th>
</tr>
</thead>
<tbody>
<tr>
<td style="text-align:center">bar</td>
<td style="text-align:right">baz</td>
</tr>
</tbody>
</table>`,
		},
		t,
	)
}

func TestTableWithAlignAttribute(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			NewTable(
				WithTableCellAlignMethod(TableCellAlignAttribute),
			),
		),
	)
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          1,
			Description: "Cell with TableCellAlignAttribute and XHTML should be rendered as an align attribute",
			Markdown: `
| abc | defghi |
:-: | -----------:
bar | baz
`,
			Expected: `<table>
<thead>
<tr>
<th align="center">abc</th>
<th align="right">defghi</th>
</tr>
</thead>
<tbody>
<tr>
<td align="center">bar</td>
<td align="right">baz</td>
</tr>
</tbody>
</table>`,
		},
		t,
	)

	md = markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			NewTable(
				WithTableCellAlignMethod(TableCellAlignAttribute),
			),
		),
	)
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          2,
			Description: "Cell with TableCellAlignAttribute and HTML5 should be rendered as an align attribute",
			Markdown: `
| abc | defghi |
:-: | -----------:
bar | baz
`,
			Expected: `<table>
<thead>
<tr>
<th align="center">abc</th>
<th align="right">defghi</th>
</tr>
</thead>
<tbody>
<tr>
<td align="center">bar</td>
<td align="right">baz</td>
</tr>
</tbody>
</table>`,
		},
		t,
	)
}

type tableStyleTransformer struct {
}

func (a *tableStyleTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	cell := node.FirstChild().FirstChild().FirstChild().(*east.TableCell)
	cell.SetAttributeString("style", []byte("font-size:1em"))
}

func TestTableWithAlignStyle(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			NewTable(
				WithTableCellAlignMethod(TableCellAlignStyle),
			),
		),
	)
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          1,
			Description: "Cell with TableCellAlignStyle and XHTML should be rendered as a style attribute",
			Markdown: `
| abc | defghi |
:-: | -----------:
bar | baz
`,
			Expected: `<table>
<thead>
<tr>
<th style="text-align:center">abc</th>
<th style="text-align:right">defghi</th>
</tr>
</thead>
<tbody>
<tr>
<td style="text-align:center">bar</td>
<td style="text-align:right">baz</td>
</tr>
</tbody>
</table>`,
		},
		t,
	)

	md = markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			NewTable(
				WithTableCellAlignMethod(TableCellAlignStyle),
			),
		),
	)
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          2,
			Description: "Cell with TableCellAlignStyle and HTML5 should be rendered as a style attribute",
			Markdown: `
| abc | defghi |
:-: | -----------:
bar | baz
`,
			Expected: `<table>
<thead>
<tr>
<th style="text-align:center">abc</th>
<th style="text-align:right">defghi</th>
</tr>
</thead>
<tbody>
<tr>
<td style="text-align:center">bar</td>
<td style="text-align:right">baz</td>
</tr>
</tbody>
</table>`,
		},
		t,
	)

	md = markdown.New(
		markdown.WithParserOptions(
			parser.WithASTTransformers(
				util.Prioritized(&tableStyleTransformer{}, 100),
			),
		),
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			NewTable(
				WithTableCellAlignMethod(TableCellAlignStyle),
			),
		),
	)

	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          3,
			Description: "Styled cell should not be broken the style by the alignments",
			Markdown: `
| abc | defghi |
:-: | -----------:
bar | baz
`,
			Expected: `<table>
<thead>
<tr>
<th style="font-size:1em;text-align:center">abc</th>
<th style="text-align:right">defghi</th>
</tr>
</thead>
<tbody>
<tr>
<td style="text-align:center">bar</td>
<td style="text-align:right">baz</td>
</tr>
</tbody>
</table>`,
		},
		t,
	)
}

func TestTableWithAlignNone(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			NewTable(
				WithTableCellAlignMethod(TableCellAlignNone),
			),
		),
	)
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          1,
			Description: "Cell with TableCellAlignNone and XHTML should not be rendered",
			Markdown: `
| abc | defghi |
:-: | -----------:
bar | baz
`,
			Expected: `<table>
<thead>
<tr>
<th>abc</th>
<th>defghi</th>
</tr>
</thead>
<tbody>
<tr>
<td>bar</td>
<td>baz</td>
</tr>
</tbody>
</table>`,
		},
		t,
	)

	md = markdown.New(
		markdown.WithRendererOptions(
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			NewTable(
				WithTableCellAlignMethod(TableCellAlignNone),
			),
		),
	)
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          2,
			Description: "Cell with TableCellAlignNone and HTML5 should not be rendered",
			Markdown: `
| abc | defghi |
:-: | -----------:
bar | baz
`,
			Expected: `<table>
<thead>
<tr>
<th>abc</th>
<th>defghi</th>
</tr>
</thead>
<tbody>
<tr>
<td>bar</td>
<td>baz</td>
</tr>
</tbody>
</table>`,
		},
		t,
	)
}

func TestTableFuzzedPanics(t *testing.T) {
	md := markdown.New(
		markdown.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
		markdown.WithExtensions(
			NewTable(),
		),
	)
	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          1,
			Description: "This should not panic",
			Markdown:    "* 0\n-|\n\t0",
			Expected: `<ul>
<li>
<table>
<thead>
<tr>
<th>0</th>
</tr>
</thead>
<tbody>
<tr>
<td>0</td>
</tr>
</tbody>
</table>
</li>
</ul>`,
		},
		t,
	)

	testutil.DoTestCase(
		md,
		testutil.MarkdownTestCase{
			No:          2,
			Description: "This should not panic",
			Markdown:    "* 0\n-|\n\t0",
			Expected: `<ul>
<li>
<table>
<thead>
<tr>
<th>0</th>
</tr>
</thead>
<tbody>
<tr>
<td>0</td>
</tr>
</tbody>
</table>
</li>
</ul>`,
		},
		t,
	)
}
