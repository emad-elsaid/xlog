package xlog

import (
	"html/template"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/stretchr/testify/require"
)

func TestRenderWidgetBehavior(t *testing.T) {
	app := newTestApp()
	app.widgets = make(map[WidgetSpace]*priorityList[WidgetFunc])

	page := &page{name: "test"}
	result := app.RenderWidget(WidgetAfterView, page)
	require.Equal(t, template.HTML(""), result, "Expected empty result for no widgets")

	testWidget := func(p Page) template.HTML {
		return template.HTML("<div>Test Widget</div>")
	}

	app.RegisterWidget(WidgetAfterView, 1.0, testWidget)
	result = app.RenderWidget(WidgetAfterView, page)
	require.Equal(t, template.HTML("<div>Test Widget</div>"), result)

	secondWidget := func(p Page) template.HTML {
		return template.HTML("<div>Second Widget</div>")
	}

	app.RegisterWidget(WidgetAfterView, 0.5, secondWidget)
	result = app.RenderWidget(WidgetAfterView, page)
	require.Equal(t, template.HTML("<div>Second Widget</div><div>Test Widget</div>"), result)
}

func TestBannerApp(t *testing.T) {
	app := newTestApp()
	tcs := []struct {
		name     string
		path     string
		content  string
		expected string
	}{
		{
			name:     "page in root and image is relative implicitly",
			path:     "home",
			content:  "![](image.jpg)",
			expected: "/image.jpg",
		},
		{
			name:     "page in root and image is relative explicitly",
			path:     "home",
			content:  "![](./image.jpg)",
			expected: "/image.jpg",
		},
		{
			name:     "page in root and image is relative explicitly in subdir",
			path:     "home",
			content:  "![](./images/image.jpg)",
			expected: "/images/image.jpg",
		},
		{
			name:     "page in subdir and image is relative implicitly",
			path:     "posts/home",
			content:  "![](image.jpg)",
			expected: "/posts/image.jpg",
		},
		{
			name:     "page in subdir and image is relative explicitly",
			path:     "posts/home",
			content:  "![](./image.jpg)",
			expected: "/posts/image.jpg",
		},
		{
			name:     "page in subdir and image is relative explicitly in subdir",
			path:     "posts/home",
			content:  "![](./images/image.jpg)",
			expected: "/posts/images/image.jpg",
		},
		{
			name:     "page in subdir and image is relative explicitly in parent",
			path:     "posts/home",
			content:  "![](../images/image.jpg)",
			expected: "/images/image.jpg",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			reader := text.NewReader([]byte(tc.content))
			p := page{
				name:       tc.path,
				lastUpdate: time.Time{},
				ast:        MarkdownConverter().Parser().Parse(reader),
				content:    (*Markdown)(&tc.content),
			}

			require.Equal(t, tc.expected, app.Banner(&p))
		})
	}
}

func TestEmoji(t *testing.T) {
	app := newTestApp()
	testCases := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Empty content",
			content:  "",
			expected: "",
		},
		{
			name:     "No emoji in content",
			content:  "Some text without an emoji",
			expected: "",
		},
		{
			name:     "Emoji in content",
			content:  ":smile:",
			expected: "😄",
		},
		{
			name:     "Emoji with other text",
			content:  "Hello :wave:",
			expected: "👋",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := &page{content: (*Markdown)(&tc.content)}
			result := app.Emoji(p)
			require.Equal(t, tc.expected, result)
		})
	}
}
