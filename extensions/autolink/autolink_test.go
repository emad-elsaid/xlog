package autolink

import (
	"bytes"
	"testing"

	"github.com/emad-elsaid/xlog/markdown"
	"github.com/emad-elsaid/xlog/markdown/parser"
	"github.com/emad-elsaid/xlog/markdown/renderer"
	"github.com/emad-elsaid/xlog/markdown/util"
)

func TestAutoLinkName(t *testing.T) {
	al := AutoLink{}
	if got := al.Name(); got != "autolink" {
		t.Errorf("Name() = %q, want %q", got, "autolink")
	}
}

func TestAutoLinkRendering(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple HTTPS URL",
			input:    "Visit <https://example.com> for more info",
			expected: `<a href="https://example.com">https://example.com</a>`,
		},
		{
			name:     "simple HTTP URL",
			input:    "Visit <http://example.com> for more info",
			expected: `<a href="http://example.com">http://example.com</a>`,
		},
		{
			name:     "email address",
			input:    "Contact me at <user@example.com>",
			expected: `<a href="mailto:user@example.com">user@example.com</a>`,
		},
		{
			name:     "long URL with truncation",
			input:    "<https://example.com/very/long/path/that/exceeds/the/character/limit/for/display>",
			expected: `https://example.com/very/long/…</a>`,
		},
		{
			name:     "URL with query parameters",
			input:    "Check <https://example.com/path?query=value&foo=bar>",
			expected: `https://example.com/path?query…`,
		},
		{
			name:     "multiple URLs in text",
			input:    "Visit <https://example.com> and <https://test.org>",
			expected: `<a href="https://example.com">https://example.com</a>`,
		},
		{
			name:     "GitHub URL",
			input:    "<https://github.com/emad-elsaid/xlog>",
			expected: `https://github.com/emad-elsaid…`,
		},
		{
			name:     "short URL no truncation",
			input:    "<https://git.io/abc>",
			expected: `<a href="https://git.io/abc">https://git.io/abc</a>`,
		},
		{
			name:     "URL with anchor fragment",
			input:    "<https://example.com/page#section>",
			expected: `https://example.com/page#secti…`,
		},
		{
			name:     "localhost URL",
			input:    "<http://localhost:8080>",
			expected: `<a href="http://localhost:8080">http://localhost:8080</a>`,
		},
		{
			name:     "IP address URL",
			input:    "<http://192.168.1.1:3000>",
			expected: `<a href="http://192.168.1.1:3000">http://192.168.1.1:3000</a>`,
		},
		{
			name:     "FTP protocol URL",
			input:    "<ftp://files.example.com>",
			expected: `<a href="ftp://files.example.com">ftp://files.example.com</a>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := markdown.New()
			
			// Register autolink parser
			md.Parser().AddOptions(parser.WithInlineParsers(
				util.Prioritized(parser.NewAutoLinkParser(), 999),
			))
			
			// Register autolink renderer
			ext := &extension{}
			md.Renderer().AddOptions(renderer.WithNodeRenderers(
				util.Prioritized(ext, -1),
			))

			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Convert() error = %v", err)
			}

			got := buf.String()
			if !bytes.Contains([]byte(got), []byte(tt.expected)) {
				t.Errorf("AutoLink rendering failed\nInput: %q\nExpected substring: %q\nGot: %q", tt.input, tt.expected, got)
			}
		})
	}
}

func TestAutoLinkHTMLEscaping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		mustHave string
		mustNot  string
	}{
		{
			name:     "ampersand escaping in query",
			input:    "<https://example.com?a=1&b=2>",
			mustHave: "&amp;",
			mustNot:  "&b=2\"", // raw ampersand shouldn't be in href
		},
		{
			name:     "proper href attribute format",
			input:    "<https://example.com/path>",
			mustHave: `href="https://example.com/path"`,
			mustNot:  "<<",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := markdown.New()
			md.Parser().AddOptions(parser.WithInlineParsers(
				util.Prioritized(parser.NewAutoLinkParser(), 999),
			))
			ext := &extension{}
			md.Renderer().AddOptions(renderer.WithNodeRenderers(
				util.Prioritized(ext, -1),
			))

			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Convert() error = %v", err)
			}

			got := buf.String()
			if !bytes.Contains([]byte(got), []byte(tt.mustHave)) {
				t.Errorf("Expected to contain %q, got: %q", tt.mustHave, got)
			}
			if tt.mustNot != "" && bytes.Contains([]byte(got), []byte(tt.mustNot)) {
				t.Errorf("Expected NOT to contain %q, got: %q", tt.mustNot, got)
			}
		})
	}
}

func TestAutoLinkLabelTruncation(t *testing.T) {
	// Test the 30-character truncation limit with ellipsis
	longURL := "<https://example.com/this/is/a/very/long/path/that/definitely/exceeds/thirty/characters>"
	
	md := markdown.New()
	md.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(parser.NewAutoLinkParser(), 999),
	))
	ext := &extension{}
	md.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(ext, -1),
	))

	var buf bytes.Buffer
	err := md.Convert([]byte(longURL), &buf)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	got := buf.String()
	// The label should be truncated with ellipsis (…)
	if !bytes.Contains([]byte(got), []byte("…")) {
		t.Errorf("Expected truncated label with ellipsis, got: %q", got)
	}
	// But the href should still have the full URL
	if !bytes.Contains([]byte(got), []byte("href=\"https://example.com/this/is/a/very/long/path/that/definitely/exceeds/thirty/characters\"")) {
		t.Errorf("Expected full URL in href, got: %q", got)
	}
}

func TestAutoLinkEmailMailtoPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check string
	}{
		{
			name:  "email gets mailto prefix",
			input: "<user@example.com>",
			check: `href="mailto:user@example.com"`,
		},
		{
			name:  "email with dots",
			input: "<first.last@example.com>",
			check: `href="mailto:first.last@example.com"`,
		},
		{
			name:  "email with plus",
			input: "<user+tag@example.com>",
			check: `href="mailto:user+tag@example.com"`,
		},
		{
			name:  "subdomain email",
			input: "<admin@mail.example.com>",
			check: `href="mailto:admin@mail.example.com"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := markdown.New()
			md.Parser().AddOptions(parser.WithInlineParsers(
				util.Prioritized(parser.NewAutoLinkParser(), 999),
			))
			ext := &extension{}
			md.Renderer().AddOptions(renderer.WithNodeRenderers(
				util.Prioritized(ext, -1),
			))

			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)
			if err != nil {
				t.Fatalf("Convert() error = %v", err)
			}

			got := buf.String()
			if !bytes.Contains([]byte(got), []byte(tt.check)) {
				t.Errorf("Expected to contain %q, got: %q", tt.check, got)
			}
		})
	}
}

func TestAutoLinkMultipleInSingleLine(t *testing.T) {
	input := "Visit <https://example.com> and email <user@example.com> today"
	
	md := markdown.New()
	md.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(parser.NewAutoLinkParser(), 999),
	))
	ext := &extension{}
	md.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(ext, -1),
	))

	var buf bytes.Buffer
	err := md.Convert([]byte(input), &buf)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	got := buf.String()
	
	// Should contain both links
	if !bytes.Contains([]byte(got), []byte(`href="https://example.com"`)) {
		t.Errorf("Expected HTTPS link, got: %q", got)
	}
	if !bytes.Contains([]byte(got), []byte(`href="mailto:user@example.com"`)) {
		t.Errorf("Expected email link, got: %q", got)
	}
}

func BenchmarkAutoLinkRendering(b *testing.B) {
	md := markdown.New()
	md.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(parser.NewAutoLinkParser(), 999),
	))
	ext := &extension{}
	md.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(ext, -1),
	))

	input := []byte("Visit <https://example.com> and contact <user@example.com> for details")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_ = md.Convert(input, &buf)
	}
}
