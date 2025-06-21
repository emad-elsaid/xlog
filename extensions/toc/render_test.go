package toc

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/ast"
)

type list []*listItem

func (tt list) Match(t *testing.T, got ast.Node) {
	t.Helper()

	if len(tt) == 0 {
		assert.Nil(t, got, "should not have a list")
		return
	}

	assert.Equal(t, len(tt), got.ChildCount(), "child count mismatch")

	child := got.FirstChild()
	for _, want := range tt {
		li, ok := child.(*ast.ListItem)
		if assert.True(t, ok, "child must be ListItem, got %T", child) {
			want.Match(t, li)
		}
		child = child.NextSibling()
	}
}

type listItem struct {
	Text string

	// If non-empty, Text should be inside a link.
	Href string

	List list
}

func (tt *listItem) Match(t *testing.T, got *ast.ListItem) {
	t.Helper()

	childCount := 0
	if len(tt.Text) > 0 {
		childCount++
	}
	if len(tt.List) > 0 {
		childCount++
	}

	assert.Equal(t, childCount, got.ChildCount(), "child count mismatch")

	child := got.FirstChild()
	if want := tt.Text; len(want) > 0 {
		if href := tt.Href; len(href) > 0 {
			a, ok := child.(*ast.Link)
			if assert.True(t, ok, "expected link, got %T", child) {
				assert.Equal(t, href, string(a.Destination), "destination mismatch")
			}
		}

		assert.Equal(t, want, string(nodeText(nil /* src */, child)))
		child = child.NextSibling()
	}

	if want := tt.List; len(want) > 0 {
		ul, ok := child.(*ast.List)
		if assert.True(t, ok, "child must be List, got %T", child) {
			want.Match(t, ul)
		}
	}
}

func TestRenderList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		give Items
		want list
	}{
		{
			desc: "empty",
			want: list{},
		},
		{
			desc: "plain",
			give: Items{
				item("foo", ""),
			},
			want: list{
				{Text: "foo"},
			},
		},
		{
			desc: "id",
			give: Items{
				item("foo", "foo"),
			},
			want: list{
				{Text: "foo", Href: "#foo"},
			},
		},
		{
			desc: "siblings",
			give: Items{
				item("foo", "foo"),
				item("bar", ""),
				item("baz", "baz"),
				item("qux", ""),
			},
			want: list{
				{Text: "foo", Href: "#foo"},
				{Text: "bar"},
				{Text: "baz", Href: "#baz"},
				{Text: "qux"},
			},
		},
		{
			desc: "subitems",
			give: Items{
				item("Foo", "foo",
					item("Bar", "bar"),
					item("Baz", "baz"),
				),
			},
			want: list{
				{
					Text: "Foo",
					Href: "#foo",
					List: list{
						{Text: "Bar", Href: "#bar"},
						{Text: "Baz", Href: "#baz"},
					},
				},
			},
		},
		{
			desc: "decrease level",
			give: Items{
				item("Foo", "foo",
					item("Bar", "bar"),
				),
				item("Baz", "baz"),
				item("Qux", "qux"),
			},
			want: list{
				{
					Text: "Foo",
					Href: "#foo",
					List: list{
						{Text: "Bar", Href: "#bar"},
					},
				},
				{Text: "Baz", Href: "#baz"},
				{Text: "Qux", Href: "#qux"},
			},
		},
		{
			desc: "several levels offset",
			// 	1	2	3	4	5	6
			give: Items{
				item("A", "a",
					item("", "",
						item("", "",
							item("", "",
								item("", "",
									item("B", "b"),
								),
							),
						),
						item("C", "c",
							item("", "",
								item("D", "d"),
							),
						),
					),
					item("E", "e"),
				),
				item("F", "f"),
				item("G", "g"),
			},
			//	1		2		3		4		5		6
			want: list{
				{
					Text: "A",
					Href: "#a",
					List: list{
						{
							List: list{
								{
									List: list{
										{
											List: list{
												{
													List: list{
														{
															Text: "B",
															Href: "#b",
														},
													},
												},
											},
										},
									},
								},
								{
									Text: "C",
									Href: "#c",
									List: list{
										{
											List: list{
												{Text: "D", Href: "#d"},
											},
										},
									},
								},
							},
						},
						{Text: "E", Href: "#e"},
					},
				},
				{Text: "F", Href: "#f"},
				{Text: "G", Href: "#g"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			got := RenderList(&TOC{Items: tt.give})
			tt.want.Match(t, got)
		})
	}
}

func TestRenderList_id(t *testing.T) {
	t.Parallel()

	node := RenderList(&TOC{
		Items: Items{
			item("Foo", "foo",
				item("Bar", "bar"),
				item("Baz", "baz"),
			),
		},
	})
	node.SetAttribute([]byte("id"), []byte("toc"))

	var buf bytes.Buffer
	err := markdown.DefaultRenderer().Render(&buf, nil, node)
	require.NoError(t, err)

	assert.Contains(t, buf.String(), `<ul id="toc">`)
}

func TestRenderList_nil(t *testing.T) {
	t.Parallel()

	assert.Nil(t, RenderList(nil))
}

func TestOrderedList(t *testing.T) {
	t.Parallel()

	node := RenderOrderedList(&TOC{
		Items: Items{
			item("Foo", "foo",
				item("Bar", "bar"),
				item("Baz", "baz"),
			),
		},
	})

	var buf bytes.Buffer
	err := markdown.DefaultRenderer().Render(&buf, nil, node)
	require.NoError(t, err)

	assert.Contains(t, buf.String(), `<ol>`)
	assert.NotContains(t, buf.String(), "start=")
}
