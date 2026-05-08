package livecodes_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/emad-elsaid/xlog/extensions/livecodes"
	"github.com/emad-elsaid/xlog/markdown"
)

func newMarkdownWithLiveCodes() markdown.Markdown {
	return markdown.New(
		markdown.WithExtensions(
			&livecodes.LiveCodes{},
		),
	)
}

func TestLiveCodesExtensionName(t *testing.T) {
	ext := livecodes.LiveCodes{}
	expected := "livecodes"
	if ext.Name() != expected {
		t.Errorf("Expected extension name %q, got %q", expected, ext.Name())
	}
}

func TestTransformLiveCodesBlocks(t *testing.T) {
	tests := []struct {
		name          string
		markdown      string
		shouldConvert bool
		expectedLang  string
	}{
		{
			name:          "live-js block",
			markdown:      "```live-js\nconsole.log('hello');\n```",
			shouldConvert: true,
			expectedLang:  "js",
		},
		{
			name:          "live-python block",
			markdown:      "```live-python\nprint('hello')\n```",
			shouldConvert: true,
			expectedLang:  "python",
		},
		{
			name:          "live.html block",
			markdown:      "```live.html\n<h1>Hello</h1>\n```",
			shouldConvert: true,
			expectedLang:  "html",
		},
		{
			name:          "live.css block",
			markdown:      "```live.css\nbody { color: red; }\n```",
			shouldConvert: true,
			expectedLang:  "css",
		},
		{
			name:          "regular js block should not convert",
			markdown:      "```js\nconsole.log('hello');\n```",
			shouldConvert: false,
			expectedLang:  "",
		},
		{
			name:          "regular python block should not convert",
			markdown:      "```python\nprint('hello')\n```",
			shouldConvert: false,
			expectedLang:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := newMarkdownWithLiveCodes()
			var buf bytes.Buffer
			if err := md.Convert([]byte(tt.markdown), &buf); err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			result := buf.String()

			if tt.shouldConvert {
				if !strings.Contains(result, "livecodes-playground") {
					t.Error("Expected LiveCodes playground to be created")
				}
				if !strings.Contains(result, `data-lang="`+tt.expectedLang+`"`) {
					t.Errorf("Expected data-lang attribute with value %q", tt.expectedLang)
				}
			} else if strings.Contains(result, "livecodes-playground") {
				t.Error("Expected LiveCodes playground NOT to be created")
			}
		})
	}
}

func TestLiveCodesRendering(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		contains []string
	}{
		{
			name:     "live-js renders playground",
			markdown: "```live-js\nconsole.log('hello');\n```",
			contains: []string{
				`class="livecodes-playground"`,
				`data-lang="js"`,
				`console.log(&#39;hello&#39;);`,
				`<script`,
				`livecodes`,
			},
		},
		{
			name:     "live-python renders playground",
			markdown: "```live-python\nprint('hello')\n```",
			contains: []string{
				`class="livecodes-playground"`,
				`data-lang="python"`,
				`print(&#39;hello&#39;)`,
			},
		},
		{
			name:     "live.html renders playground",
			markdown: "```live.html\n<h1>Test</h1>\n```",
			contains: []string{
				`class="livecodes-playground"`,
				`data-lang="html"`,
				`&lt;h1&gt;Test&lt;/h1&gt;`,
			},
		},
		{
			name:     "special characters are escaped",
			markdown: "```live-js\nlet x = \"<>&'\";\n```",
			contains: []string{
				`&lt;`,
				`&gt;`,
				`&amp;`,
				`&quot;`,
				`&#39;`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := newMarkdownWithLiveCodes()
			var buf bytes.Buffer
			if err := md.Convert([]byte(tt.markdown), &buf); err != nil {
				t.Fatalf("Failed to convert markdown: %v", err)
			}

			result := buf.String()

			for _, substring := range tt.contains {
				if !strings.Contains(result, substring) {
					t.Errorf("Expected output to contain %q\nGot: %s", substring, result)
				}
			}
		})
	}
}

func TestLiveCodesDoesNotAffectRegularCodeBlocks(t *testing.T) {
	input := "```js\nconsole.log('regular');\n```"

	md := newMarkdownWithLiveCodes()
	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		t.Fatalf("Failed to convert markdown: %v", err)
	}

	result := buf.String()

	// Regular code blocks should NOT contain livecodes-playground class
	if strings.Contains(result, "livecodes-playground") {
		t.Error("Regular code block should not be converted to LiveCodes playground")
	}

	// It should still contain code block elements
	if !strings.Contains(result, "<code") {
		t.Error("Expected regular code block to be rendered")
	}
}

func TestLiveCodesScriptEmbedded(t *testing.T) {
	// Test that the script is embedded and contains expected content
	md := newMarkdownWithLiveCodes()
	input := "```live-js\ntest\n```"

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		t.Fatalf("Failed to convert markdown: %v", err)
	}

	result := buf.String()

	// Verify it contains livecodes-related content
	if !strings.Contains(result, "livecodes") {
		t.Error("Expected output to contain 'livecodes' reference")
	}

	// Verify it contains createPlayground function
	if !strings.Contains(result, "createPlayground") {
		t.Error("Expected output to contain 'createPlayground' function")
	}
}

func TestLiveCodesInit(t *testing.T) {
	// Verify that the extension works when initialized
	md := newMarkdownWithLiveCodes()
	input := "```live-js\ntest\n```"

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		t.Fatalf("Failed to convert markdown: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Expected non-empty output")
	}

	// Check that playground was created
	if !strings.Contains(buf.String(), "livecodes-playground") {
		t.Error("Expected LiveCodes playground to be created")
	}
}

func TestTransformWithMultipleBlocks(t *testing.T) {
	input := `
# Test Document

Regular code:
` + "```js\nconsole.log('regular');\n```" + `

Live code:
` + "```live-js\nconsole.log('live');\n```" + `

Another regular:
` + "```python\nprint('regular')\n```" + `

Another live:
` + "```live-python\nprint('live')\n```" + `
`

	md := newMarkdownWithLiveCodes()
	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		t.Fatalf("Failed to convert markdown: %v", err)
	}

	result := buf.String()

	// Should have exactly 2 LiveCodes playground divs (the word appears multiple times due to script)
	liveCodesCount := strings.Count(result, `class="livecodes-playground"`)
	if liveCodesCount != 2 {
		t.Errorf("Expected 2 LiveCodes playgrounds, got %d", liveCodesCount)
	}

	// Check both languages are present
	if !strings.Contains(result, `data-lang="js"`) {
		t.Error("Expected JS playground to be created")
	}
	if !strings.Contains(result, `data-lang="python"`) {
		t.Error("Expected Python playground to be created")
	}

	// Should still have regular code blocks (they get syntax highlighting)
	regularJSCount := strings.Count(result, `class="language-js"`)
	if regularJSCount == 0 {
		t.Error("Expected regular JS code block to be preserved")
	}
}
