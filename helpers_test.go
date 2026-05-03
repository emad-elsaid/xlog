package xlog

import (
	"testing"
	"time"

	"github.com/emad-elsaid/xlog/markdown/text"
	"github.com/stretchr/testify/require"
)

func TestBanner(t *testing.T) {
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

			require.Equal(t, tc.expected, Banner(&p))
		})
	}
}

func TestBanner_EdgeCases(t *testing.T) {
	tcs := []struct {
		name     string
		path     string
		content  string
		expected string
	}{
		{
			name:     "no content returns empty",
			path:     "test",
			content:  "",
			expected: "",
		},
		{
			name:     "no image returns empty",
			path:     "test",
			content:  "Just text without image",
			expected: "",
		},
		{
			name:     "image not first element returns empty",
			path:     "test",
			content:  "Text first\n\n![](image.jpg)",
			expected: "",
		},
		{
			name:     "empty image destination returns empty",
			path:     "test",
			content:  "![]()",
			expected: "",
		},
		{
			name:     "hash image destination returns empty",
			path:     "test",
			content:  "![](#)",
			expected: "",
		},
		{
			name:     "absolute path preserved",
			path:     "posts/home",
			content:  "![](/images/banner.jpg)",
			expected: "/images/banner.jpg",
		},
		{
			name:     "http URL preserved",
			path:     "posts/home",
			content:  "![](http://example.com/image.jpg)",
			expected: "http://example.com/image.jpg",
		},
		{
			name:     "https URL preserved",
			path:     "posts/home",
			content:  "![](https://example.com/image.jpg)",
			expected: "https://example.com/image.jpg",
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

			require.Equal(t, tc.expected, Banner(&p))
		})
	}
}

func TestEmoji(t *testing.T) {
	tcs := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "page with emoji returns unicode",
			content:  ":smile:",
			expected: "😄",
		},
		{
			name:     "page without emoji returns empty",
			content:  "Just plain text",
			expected: "",
		},
		{
			name:     "empty page returns empty",
			content:  "",
			expected: "",
		},
		{
			name:     "multiple emojis returns first",
			content:  ":smile: :heart:",
			expected: "😄",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			reader := text.NewReader([]byte(tc.content))
			p := page{
				name:       "test",
				lastUpdate: time.Time{},
				ast:        MarkdownConverter().Parser().Parse(reader),
				content:    (*Markdown)(&tc.content),
			}

			require.Equal(t, tc.expected, Emoji(&p))
		})
	}
}

func TestDir(t *testing.T) {
	tcs := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "file in directory returns directory",
			input:    "posts/article.md",
			expected: "posts",
		},
		{
			name:     "nested path returns parent",
			input:    "a/b/c/file.md",
			expected: "a/b/c",
		},
		{
			name:     "file in root returns empty",
			input:    "index.md",
			expected: "",
		},
		{
			name:     "single element returns empty",
			input:    "file",
			expected: "",
		},
		{
			name:     "slash returns slash",
			input:    "/",
			expected: "/",
		},
		{
			name:     "absolute path returns parent",
			input:    "/posts/article.md",
			expected: "/posts",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := dir(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func BenchmarkBanner(b *testing.B) {
	tests := []struct {
		name    string
		path    string
		content string
	}{
		{
			name:    "simple banner",
			path:    "post.md",
			content: "![](banner.jpg)",
		},
		{
			name:    "banner in subdirectory",
			path:    "posts/article.md",
			content: "![](./images/banner.png)",
		},
		{
			name:    "absolute URL banner",
			path:    "post.md",
			content: "![](https://example.com/banner.jpg)",
		},
		{
			name:    "no banner",
			path:    "post.md",
			content: "Just plain text without image",
		},
		{
			name:    "large document with banner",
			path:    "docs/guide.md",
			content: "![](hero.jpg)\n\n" + string(make([]byte, 10000)),
		},
	}

	for _, tc := range tests {
		b.Run(tc.name, func(b *testing.B) {
			reader := text.NewReader([]byte(tc.content))
			p := page{
				name:       tc.path,
				lastUpdate: time.Time{},
				ast:        MarkdownConverter().Parser().Parse(reader),
				content:    (*Markdown)(&tc.content),
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = Banner(&p)
			}
		})
	}
}
