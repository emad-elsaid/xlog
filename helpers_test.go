package xlog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/emad-elsaid/xlog/markdown/text"
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
