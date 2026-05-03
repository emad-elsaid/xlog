package ast

import (
	"bytes"
	"testing"

	gast "github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/emoji/definition"
)

func TestNewEmoji(t *testing.T) {
	tests := []struct {
		name      string
		shortName []byte
		value     *definition.Emoji
	}{
		{
			name:      "simple emoji",
			shortName: []byte("smile"),
			value:     &definition.Emoji{Name: "smile", Unicode: []rune("😊"), ShortNames: []string{"smile"}},
		},
		{
			name:      "empty short name",
			shortName: []byte(""),
			value:     &definition.Emoji{Name: "party", Unicode: []rune("🎉"), ShortNames: []string{"party"}},
		},
		{
			name:      "nil value",
			shortName: []byte("test"),
			value:     nil,
		},
		{
			name:      "complex emoji",
			shortName: []byte("thumbsup"),
			value:     &definition.Emoji{Name: "thumbs up", Unicode: []rune("👍"), ShortNames: []string{"thumbsup", "+1"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			emoji := NewEmoji(tc.shortName, tc.value)

			if emoji == nil {
				t.Fatal("NewEmoji returned nil")
			}

			if !bytes.Equal(emoji.ShortName, tc.shortName) {
				t.Errorf("ShortName = %q, want %q", emoji.ShortName, tc.shortName)
			}

			if emoji.Value != tc.value {
				t.Errorf("Value = %v, want %v", emoji.Value, tc.value)
			}
		})
	}
}

func TestEmoji_Kind(t *testing.T) {
	emoji := NewEmoji([]byte("test"), nil)

	kind := emoji.Kind()

	if kind != KindEmoji {
		t.Errorf("Kind() = %v, want %v", kind, KindEmoji)
	}
}

func TestKindEmoji_Identity(t *testing.T) {
	// Verify KindEmoji is properly initialized.
	if KindEmoji == gast.NodeKind(0) {
		t.Error("KindEmoji should not be zero value")
	}

	// Create two emojis and verify they have the same kind.
	emoji1 := NewEmoji([]byte("one"), nil)
	emoji2 := NewEmoji([]byte("two"), nil)

	if emoji1.Kind() != emoji2.Kind() {
		t.Error("Different emoji instances should have the same Kind")
	}
}

func TestEmoji_Dump(t *testing.T) {
	tests := []struct {
		name      string
		shortName []byte
		value     *definition.Emoji
		source    []byte
		level     int
	}{
		{
			name:      "basic dump",
			shortName: []byte("smile"),
			value:     &definition.Emoji{Name: "smile", Unicode: []rune("😊"), ShortNames: []string{"smile"}},
			source:    []byte(":smile:"),
			level:     0,
		},
		{
			name:      "nested level",
			shortName: []byte("heart"),
			value:     &definition.Emoji{Name: "heart", Unicode: []rune("❤️"), ShortNames: []string{"heart"}},
			source:    []byte(":heart:"),
			level:     2,
		},
		{
			name:      "nil value dump",
			shortName: []byte("unknown"),
			value:     nil,
			source:    []byte(":unknown:"),
			level:     1,
		},
		{
			name:      "empty source",
			shortName: []byte("test"),
			value:     &definition.Emoji{Name: "target", Unicode: []rune("🎯"), ShortNames: []string{"target"}},
			source:    []byte(""),
			level:     0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			emoji := NewEmoji(tc.shortName, tc.value)

			// Dump should not panic.
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Dump() panicked: %v", r)
				}
			}()

			emoji.Dump(tc.source, tc.level)
		})
	}
}

func TestEmoji_BaseInlineEmbedding(t *testing.T) {
	// Verify that Emoji properly embeds BaseInline.
	emoji := NewEmoji([]byte("test"), nil)

	// Should be able to call BaseInline methods without panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("BaseInline method call panicked: %v", r)
		}
	}()

	// Test that it's a valid inline node.
	var _ gast.Node = emoji
}
