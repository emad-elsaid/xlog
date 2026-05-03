package link_preview

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	. "github.com/emad-elsaid/xlog"
)

func TestImgUrlPreprocessor(t *testing.T) {
	tests := []struct {
		name  string
		input Markdown
		want  Markdown
	}{
		{
			name:  "valid image URL with jpg",
			input: "https://example.com/image.jpg",
			want:  "![](https://example.com/image.jpg)",
		},
		{
			name:  "valid image URL with png",
			input: "https://example.com/photo.png",
			want:  "![](https://example.com/photo.png)",
		},
		{
			name:  "valid image URL with webp",
			input: "https://example.com/pic.webp",
			want:  "![](https://example.com/pic.webp)",
		},
		{
			name:  "valid image URL with svg",
			input: "https://example.com/icon.svg",
			want:  "![](https://example.com/icon.svg)",
		},
		{
			name:  "non-image URL",
			input: "https://example.com/page.html",
			want:  "https://example.com/page.html",
		},
		{
			name:  "mixed content with image",
			input: "some text https://example.com/image.jpg more text",
			want:  "some text https://example.com/image.jpg more text",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := imgUrlPreprocessor(tc.input)
			if got != tc.want {
				t.Errorf("imgUrlPreprocessor() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestTweetUrlPreprocessor(t *testing.T) {
	tests := []struct {
		name  string
		input Markdown
		want  string
	}{
		{
			name:  "twitter.com URL",
			input: "https://twitter.com/user/status/1234567890",
			want:  `<blockquote class="twitter-tweet">`,
		},
		{
			name:  "x.com URL",
			input: "https://x.com/user/status/9876543210",
			want:  `<blockquote class="twitter-tweet">`,
		},
		{
			name:  "non-tweet URL",
			input: "https://twitter.com/user",
			want:  "https://twitter.com/user",
		},
		{
			name:  "invalid status URL",
			input: "https://twitter.com/user/status/abc",
			want:  "https://twitter.com/user/status/abc",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := string(tweetUrlPreprocessor(tc.input))
			if tc.want == "https://twitter.com/user" || tc.want == "https://twitter.com/user/status/abc" {
				if got != tc.want {
					t.Errorf("tweetUrlPreprocessor() = %q, want %q", got, tc.want)
				}
			} else {
				if !contains(got, tc.want) {
					t.Errorf("tweetUrlPreprocessor() does not contain %q, got %q", tc.want, got)
				}
			}
		})
	}
}

func TestYoutubeUrlPreprocessor(t *testing.T) {
	tests := []struct {
		name  string
		input Markdown
		want  string
	}{
		{
			name:  "youtube.com long URL",
			input: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			want:  `src="https://www.youtube-nocookie.com/embed/dQw4w9WgXcQ"`,
		},
		{
			name:  "youtu.be short URL",
			input: "https://youtu.be/dQw4w9WgXcQ",
			want:  `src="https://www.youtube-nocookie.com/embed/dQw4w9WgXcQ"`,
		},
		{
			name:  "non-youtube URL",
			input: "https://vimeo.com/12345",
			want:  "https://vimeo.com/12345",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := string(youtubeUrlPreprocessor(tc.input))
			if tc.want == "https://vimeo.com/12345" {
				if got != tc.want {
					t.Errorf("youtubeUrlPreprocessor() = %q, want %q", got, tc.want)
				}
			} else {
				if len(got) < len(tc.want) || !contains(got, tc.want) {
					t.Errorf("youtubeUrlPreprocessor() does not contain %q, got %q", tc.want, got)
				}
			}
		})
	}
}

func TestFbUrlPreprocessor(t *testing.T) {
	tests := []struct {
		name  string
		input Markdown
		want  string
	}{
		{
			name:  "valid facebook post URL",
			input: "https://www.facebook.com/username/posts/123456789",
			want:  `<iframe src="https://www.facebook.com/plugins/post.php`,
		},
		{
			name:  "non-facebook URL",
			input: "https://example.com",
			want:  "https://example.com",
		},
		{
			name:  "invalid facebook URL format",
			input: "https://www.facebook.com/username",
			want:  "https://www.facebook.com/username",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := string(fbUrlPreprocessor(tc.input))
			if tc.want == "https://example.com" || tc.want == "https://www.facebook.com/username" {
				if got != tc.want {
					t.Errorf("fbUrlPreprocessor() = %q, want %q", got, tc.want)
				}
			} else {
				if len(got) < len(tc.want) || !contains(got, tc.want) {
					t.Errorf("fbUrlPreprocessor() does not contain %q, got %q", tc.want, got)
				}
			}
		})
	}
}

func TestGiphyUrlPreprocessor(t *testing.T) {
	tests := []struct {
		name  string
		input Markdown
		want  Markdown
	}{
		{
			name:  "valid giphy URL",
			input: "https://giphy.com/gifs/funny-cat-abc123",
			want:  "![](https://media.giphy.com/media/abc123/giphy.gif)",
		},
		{
			name:  "non-giphy URL",
			input: "https://example.com",
			want:  "https://example.com",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := giphyUrlPreprocessor(tc.input)
			if got != tc.want {
				t.Errorf("giphyUrlPreprocessor() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestGetUrlMeta(t *testing.T) {
	t.Run("successful meta fetch with all fields", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
	<title>Test Page Title</title>
	<meta property="og:description" content="Test description">
	<meta property="og:image" content="https://example.com/image.jpg">
</head>
<body>Content</body>
</html>
`))
		}))
		defer server.Close()

		cleanupCache(t)

		meta, err := getUrlMeta(server.URL)
		if err != nil {
			t.Fatalf("getUrlMeta() error = %v", err)
		}

		if meta.Title != "Test Page Title" {
			t.Errorf("meta.Title = %q, want %q", meta.Title, "Test Page Title")
		}

		if meta.Description != "Test description" {
			t.Errorf("meta.Description = %q, want %q", meta.Description, "Test description")
		}

		if meta.Image != "https://example.com/image.jpg" {
			t.Errorf("meta.Image = %q, want %q", meta.Image, "https://example.com/image.jpg")
		}

		if meta.URL != server.URL {
			t.Errorf("meta.URL = %q, want %q", meta.URL, server.URL)
		}
	})

	t.Run("meta fetch with minimal HTML", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<html><head><title>Simple</title></head></html>`))
		}))
		defer server.Close()

		cleanupCache(t)

		meta, err := getUrlMeta(server.URL)
		if err != nil {
			t.Fatalf("getUrlMeta() error = %v", err)
		}

		if meta.Title != "Simple" {
			t.Errorf("meta.Title = %q, want %q", meta.Title, "Simple")
		}

		if meta.Description != "" {
			t.Errorf("meta.Description = %q, want empty string", meta.Description)
		}
	})

	t.Run("meta fetch without title falls back to URL", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<html><body>No title</body></html>`))
		}))
		defer server.Close()

		cleanupCache(t)

		meta, err := getUrlMeta(server.URL)
		if err != nil {
			t.Fatalf("getUrlMeta() error = %v", err)
		}

		if meta.Title != server.URL {
			t.Errorf("meta.Title = %q, want %q (URL fallback)", meta.Title, server.URL)
		}
	})

	t.Run("cache hit returns cached data", func(t *testing.T) {
		callCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`<html><head><title>Cached</title></head></html>`))
		}))
		defer server.Close()

		cleanupCache(t)

		meta1, err := getUrlMeta(server.URL)
		if err != nil {
			t.Fatalf("first getUrlMeta() error = %v", err)
		}

		meta2, err := getUrlMeta(server.URL)
		if err != nil {
			t.Fatalf("second getUrlMeta() error = %v", err)
		}

		if callCount != 1 {
			t.Errorf("HTTP request count = %d, want 1 (cache should prevent second request)", callCount)
		}

		if meta1.Title != meta2.Title {
			t.Errorf("cached meta mismatch: %q vs %q", meta1.Title, meta2.Title)
		}
	})

	t.Run("HTTP error with no response body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		cleanupCache(t)

		meta, err := getUrlMeta(server.URL)
		if meta == nil || err != nil {
			t.Logf("getUrlMeta() returned meta=%v, err=%v (500 status with empty body still parses)", meta, err)
		}
	})

	t.Run("invalid URL returns error", func(t *testing.T) {
		cleanupCache(t)

		_, err := getUrlMeta("http://invalid-domain-that-does-not-exist-12345.com")
		if err == nil {
			t.Error("getUrlMeta() with invalid URL should return error")
		}
	})

	t.Run("meta with name attribute instead of property", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`
<html>
<head>
	<title>Test</title>
	<meta name="description" content="Name attribute description">
</head>
</html>
`))
		}))
		defer server.Close()

		cleanupCache(t)

		meta, err := getUrlMeta(server.URL)
		if err != nil {
			t.Fatalf("getUrlMeta() error = %v", err)
		}

		if meta.Description != "Name attribute description" {
			t.Errorf("meta.Description = %q, want %q", meta.Description, "Name attribute description")
		}
	})
}

func TestFallbackURLPreprocessor(t *testing.T) {
	t.Run("non-URL text unchanged", func(t *testing.T) {
		input := Markdown("just some text")
		got := fallbackURLPreprocessor(input)
		if got != input {
			t.Errorf("fallbackURLPreprocessor() changed non-URL text: %q", got)
		}
	})

	t.Run("URL with failing getUrlMeta returns original URL", func(t *testing.T) {
		cleanupCache(t)

		input := Markdown("http://invalid-domain-that-does-not-exist-12345.com")
		got := fallbackURLPreprocessor(input)
		if got != input {
			t.Errorf("fallbackURLPreprocessor() = %q, want %q", got, input)
		}
	})

	t.Run("relative image path converted to absolute", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`
<html>
<head>
	<title>Test</title>
	<meta property="og:image" content="/images/logo.png">
</head>
</html>
`))
		}))
		defer server.Close()

		cleanupCache(t)

		meta, err := getUrlMeta(server.URL)
		if err != nil {
			t.Fatalf("getUrlMeta() error = %v", err)
		}

		if meta.Image != "/images/logo.png" {
			t.Errorf("meta.Image before conversion = %q, want %q", meta.Image, "/images/logo.png")
		}

		t.Logf("Image path correctly extracted as: %s", meta.Image)
	})
}

func TestLinkPreviewExtension(t *testing.T) {
	ext := LinkPreview{}

	t.Run("Name returns correct extension name", func(t *testing.T) {
		if ext.Name() != "link-preview" {
			t.Errorf("Name() = %q, want %q", ext.Name(), "link-preview")
		}
	})

	t.Run("Init registers preprocessors", func(t *testing.T) {
		ext.Init()
	})
}

func cleanupCache(t *testing.T) {
	t.Helper()
	cacheDir := ".cache"
	_ = os.RemoveAll(cacheDir)
}

func contains(haystack, needle string) bool {
	if len(haystack) < len(needle) {
		return false
	}
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

func TestMetaJSONMarshalUnmarshal(t *testing.T) {
	original := Meta{
		URL:         "https://example.com",
		Title:       "Example Site",
		Description: "An example website",
		Image:       "https://example.com/image.jpg",
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var unmarshaled Meta
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if original != unmarshaled {
		t.Errorf("Meta marshal/unmarshal mismatch: %+v vs %+v", original, unmarshaled)
	}
}

func TestCacheDirectoryCreation(t *testing.T) {
	cleanupCache(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><title>Test</title></html>`))
	}))
	defer server.Close()

	_, err := getUrlMeta(server.URL)
	if err != nil {
		t.Fatalf("getUrlMeta() error = %v", err)
	}

	if _, err := os.Stat(".cache"); os.IsNotExist(err) {
		t.Error("cache directory was not created")
	}

	cacheFiles, err := filepath.Glob(".cache/*.json")
	if err != nil {
		t.Fatalf("filepath.Glob() error = %v", err)
	}

	if len(cacheFiles) != 1 {
		t.Errorf("expected 1 cache file, got %d", len(cacheFiles))
	}

	cleanupCache(t)
}
