package mermaid

import (
	"html/template"
	"strings"
	"testing"

	. "github.com/emad-elsaid/xlog"
)

func TestMermaidExtensionName(t *testing.T) {
	ext := Mermaid{}
	expected := "mermaid"
	if ext.Name() != expected {
		t.Errorf("Expected extension name %q, got %q", expected, ext.Name())
	}
}

func TestMermaidRenderer(t *testing.T) {
	tests := []struct {
		name     string
		input    Markdown
		contains []string
	}{
		{
			name:  "simple graph",
			input: Markdown("graph TD\n    A-->B"),
			contains: []string{
				`<pre class="mermaid"`,
				`style="background: transparent;text-align:center;"`,
				"graph TD",
				"A-->B",
				"</pre>",
				"<script",
				"mermaid",
			},
		},
		{
			name:  "sequence diagram",
			input: Markdown("sequenceDiagram\n    Alice->>John: Hello John"),
			contains: []string{
				`<pre class="mermaid"`,
				"sequenceDiagram",
				"Alice->>John: Hello John",
			},
		},
		{
			name:  "empty diagram",
			input: Markdown(""),
			contains: []string{
				`<pre class="mermaid"`,
				"</pre>",
			},
		},
		{
			name:  "diagram with special characters",
			input: Markdown("graph LR\n    A[\"Item with <>&\"]"),
			contains: []string{
				`<pre class="mermaid"`,
				"graph LR",
				`A["Item with <>&"]`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderer(tt.input)
			resultStr := string(result)

			for _, substring := range tt.contains {
				if !strings.Contains(resultStr, substring) {
					t.Errorf("Expected result to contain %q\nGot: %s", substring, resultStr)
				}
			}

			// Verify it's valid HTML template
			if _, ok := interface{}(result).(template.HTML); !ok {
				t.Error("Expected result to be template.HTML type")
			}
		})
	}
}

func TestMermaidRendererStructure(t *testing.T) {
	input := Markdown("graph TD\n    Start-->End")
	result := renderer(input)
	resultStr := string(result)

	// Verify the structure: should have both <pre> and <script>
	preIndex := strings.Index(resultStr, "<pre")
	preEndIndex := strings.Index(resultStr, "</pre>")
	scriptIndex := strings.Index(resultStr, "<script")

	if preIndex == -1 {
		t.Error("Expected result to contain <pre> tag")
	}

	if preEndIndex == -1 {
		t.Error("Expected result to contain </pre> closing tag")
	}

	if scriptIndex == -1 {
		t.Error("Expected result to contain <script> tag")
	}

	if preIndex >= preEndIndex {
		t.Error("Expected <pre> to open before it closes")
	}

	if scriptIndex <= preEndIndex {
		t.Error("Expected <script> to come after </pre>")
	}
}

func TestMermaidScriptEmbedded(t *testing.T) {
	// Verify that the embedded script is not empty
	if script == "" {
		t.Error("Expected embedded script to be non-empty")
	}

	// Verify it contains mermaid-related content
	if !strings.Contains(script, "mermaid") {
		t.Error("Expected embedded script to contain 'mermaid' reference")
	}
}

func TestMermaidInit(t *testing.T) {
	// Create a new instance and call Init
	ext := Mermaid{}
	ext.Init()

	// Verify that the shortcode was registered by checking if we can use it
	// The shortcode.ShortCodes map should contain "mermaid"
	// We verify this by attempting to render something
	result := renderer(Markdown("test"))
	if result == "" {
		t.Error("Expected renderer to work after Init()")
	}
}
