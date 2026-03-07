package html

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
)

func TestRenderAttributes_NumericValues(t *testing.T) {
	tests := []struct {
		name     string
		attrName string
		value    any
		expected string
	}{
		{
			name:     "int value",
			attrName: "data-count",
			value:    42,
			expected: ` data-count="42"`,
		},
		{
			name:     "int8 value",
			attrName: "data-small",
			value:    int8(8),
			expected: ` data-small="8"`,
		},
		{
			name:     "int16 value",
			attrName: "data-medium",
			value:    int16(16),
			expected: ` data-medium="16"`,
		},
		{
			name:     "int32 value",
			attrName: "data-large",
			value:    int32(32),
			expected: ` data-large="32"`,
		},
		{
			name:     "int64 value",
			attrName: "data-xlarge",
			value:    int64(64),
			expected: ` data-xlarge="64"`,
		},
		{
			name:     "uint value",
			attrName: "data-unsigned",
			value:    uint(100),
			expected: ` data-unsigned="100"`,
		},
		{
			name:     "uint8 value",
			attrName: "data-byte",
			value:    uint8(255),
			expected: ` data-byte="255"`,
		},
		{
			name:     "uint16 value",
			attrName: "data-ushort",
			value:    uint16(65535),
			expected: ` data-ushort="65535"`,
		},
		{
			name:     "uint32 value",
			attrName: "data-ulong",
			value:    uint32(4294967295),
			expected: ` data-ulong="4294967295"`,
		},
		{
			name:     "uint64 value",
			attrName: "data-uxlarge",
			value:    uint64(18446744073709551615),
			expected: ` data-uxlarge="18446744073709551615"`,
		},
		{
			name:     "float32 value",
			attrName: "data-ratio",
			value:    float32(3.14),
			expected: ` data-ratio="3.14"`,
		},
		{
			name:     "float64 value",
			attrName: "data-pi",
			value:    float64(3.141592653589793),
			expected: ` data-pi="3.141592653589793"`,
		},
		{
			name:     "bool true value",
			attrName: "data-enabled",
			value:    true,
			expected: ` data-enabled="true"`,
		},
		{
			name:     "bool false value",
			attrName: "data-disabled",
			value:    false,
			expected: ` data-disabled="false"`,
		},
		{
			name:     "string value (unchanged)",
			attrName: "class",
			value:    "foo bar",
			expected: ` class="foo bar"`,
		},
		{
			name:     "[]byte value (unchanged)",
			attrName: "id",
			value:    []byte("test-id"),
			expected: ` id="test-id"`,
		},
		{
			name:     "negative int",
			attrName: "data-negative",
			value:    -42,
			expected: ` data-negative="-42"`,
		},
		{
			name:     "zero value",
			attrName: "data-zero",
			value:    0,
			expected: ` data-zero="0"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := ast.NewTextBlock()
			node.SetAttributeString(tt.attrName, tt.value)

			var buf bytes.Buffer
			writer := bufio.NewWriter(&buf)
			RenderAttributes(writer, node, nil)
			writer.Flush()

			got := buf.String()
			if got != tt.expected {
				t.Errorf("RenderAttributes() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestRenderAttributes_MixedTypes(t *testing.T) {
	node := ast.NewTextBlock()
	node.SetAttributeString("class", "container")
	node.SetAttributeString("data-count", 42)
	node.SetAttributeString("data-enabled", true)
	node.SetAttributeString("id", []byte("test"))

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	RenderAttributes(writer, node, nil)
	writer.Flush()

	result := buf.String()

	// Check that all attributes are present (order may vary)
	expectedParts := []string{
		`class="container"`,
		`data-count="42"`,
		`data-enabled="true"`,
		`id="test"`,
	}

	for _, part := range expectedParts {
		if !bytes.Contains([]byte(result), []byte(part)) {
			t.Errorf("RenderAttributes() missing expected part %q in %q", part, result)
		}
	}
}

func BenchmarkRenderAttributes_NumericConversion(b *testing.B) {
	node := ast.NewTextBlock()
	node.SetAttributeString("data-int", 42)
	node.SetAttributeString("data-float", 3.14)
	node.SetAttributeString("data-bool", true)
	node.SetAttributeString("class", "test")

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		writer.Reset(&buf)
		RenderAttributes(writer, node, nil)
		writer.Flush()
	}
}

// Helper benchmark to show old behavior would panic or produce incorrect output
func BenchmarkFormatConversions(b *testing.B) {
	values := []any{
		42,
		int8(8),
		int16(16),
		int32(32),
		int64(64),
		uint(100),
		uint8(255),
		uint16(65535),
		uint32(4294967295),
		uint64(18446744073709551615),
		float32(3.14),
		float64(3.141592653589793),
		true,
		false,
	}

	b.Run("fmt.Sprint", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range values {
				_ = fmt.Sprint(v)
			}
		}
	})

	b.Run("strconv.Format", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, v := range values {
				switch typed := v.(type) {
				case int:
					_ = strconv.Itoa(typed)
				case int64:
					_ = strconv.FormatInt(typed, 10)
				case uint64:
					_ = strconv.FormatUint(typed, 10)
				case float64:
					_ = strconv.FormatFloat(typed, 'f', -1, 64)
				case bool:
					_ = strconv.FormatBool(typed)
				}
			}
		}
	})
}
