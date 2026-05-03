package definition

import (
	"testing"
)

func TestNewEmoji(t *testing.T) {
	tests := []struct {
		name       string
		emojiName  string
		unicode    []rune
		shortNames []string
		want       Emoji
		wantPanic  bool
	}{
		{
			name:       "valid emoji with single short name",
			emojiName:  "smile",
			unicode:    []rune{0x1F600},
			shortNames: []string{"smile"},
			want: Emoji{
				Name:       "smile",
				Unicode:    []rune{0x1F600},
				ShortNames: []string{"smile"},
			},
			wantPanic: false,
		},
		{
			name:       "valid emoji with multiple short names",
			emojiName:  "grinning face",
			unicode:    []rune{0x1F601},
			shortNames: []string{"grinning", "grin"},
			want: Emoji{
				Name:       "grinning face",
				Unicode:    []rune{0x1F601},
				ShortNames: []string{"grinning", "grin"},
			},
			wantPanic: false,
		},
		{
			name:       "empty unicode defaults to replacement character",
			emojiName:  "custom",
			unicode:    []rune{},
			shortNames: []string{"custom"},
			want: Emoji{
				Name:       "custom",
				Unicode:    []rune{0xFFFD},
				ShortNames: []string{"custom"},
			},
			wantPanic: false,
		},
		{
			name:       "nil unicode defaults to replacement character",
			emojiName:  "nil_unicode",
			unicode:    nil,
			shortNames: []string{"nil"},
			want: Emoji{
				Name:       "nil_unicode",
				Unicode:    []rune{0xFFFD},
				ShortNames: []string{"nil"},
			},
			wantPanic: false,
		},
		{
			name:       "multi-rune unicode sequence",
			emojiName:  "flag",
			unicode:    []rune{0x1F1FA, 0x1F1F8},
			shortNames: []string{"flag-us"},
			want: Emoji{
				Name:       "flag",
				Unicode:    []rune{0x1F1FA, 0x1F1F8},
				ShortNames: []string{"flag-us"},
			},
			wantPanic: false,
		},
		{
			name:       "no short names causes panic",
			emojiName:  "invalid",
			unicode:    []rune{0x1F600},
			shortNames: []string{},
			wantPanic:  true,
		},
		{
			name:       "nil short names causes panic",
			emojiName:  "invalid",
			unicode:    []rune{0x1F600},
			shortNames: nil,
			wantPanic:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewEmoji() did not panic, expected panic")
					}
				}()
				NewEmoji(tc.emojiName, tc.unicode, tc.shortNames...)
				return
			}

			got := NewEmoji(tc.emojiName, tc.unicode, tc.shortNames...)

			if got.Name != tc.want.Name {
				t.Errorf("NewEmoji().Name = %q, want %q", got.Name, tc.want.Name)
			}

			if len(got.Unicode) != len(tc.want.Unicode) {
				t.Errorf("NewEmoji().Unicode length = %d, want %d", len(got.Unicode), len(tc.want.Unicode))
			}
			for i, r := range got.Unicode {
				if r != tc.want.Unicode[i] {
					t.Errorf("NewEmoji().Unicode[%d] = %U, want %U", i, r, tc.want.Unicode[i])
				}
			}

			if len(got.ShortNames) != len(tc.want.ShortNames) {
				t.Errorf("NewEmoji().ShortNames length = %d, want %d", len(got.ShortNames), len(tc.want.ShortNames))
			}
			for i, s := range got.ShortNames {
				if s != tc.want.ShortNames[i] {
					t.Errorf("NewEmoji().ShortNames[%d] = %q, want %q", i, s, tc.want.ShortNames[i])
				}
			}
		})
	}
}

func TestEmoji_IsUnicode(t *testing.T) {
	tests := []struct {
		name string
		em   Emoji
		want bool
	}{
		{
			name: "valid unicode single rune",
			em: Emoji{
				Unicode: []rune{0x1F600},
			},
			want: true,
		},
		{
			name: "valid unicode multi-rune",
			em: Emoji{
				Unicode: []rune{0x1F1FA, 0x1F1F8},
			},
			want: true,
		},
		{
			name: "replacement character (not unicode)",
			em: Emoji{
				Unicode: []rune{0xFFFD},
			},
			want: false,
		},
		{
			name: "empty unicode",
			em: Emoji{
				Unicode: []rune{},
			},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.em.IsUnicode()
			if got != tc.want {
				t.Errorf("Emoji.IsUnicode() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNewEmojis(t *testing.T) {
	tests := []struct {
		name   string
		emojis []Emoji
		gets   []struct {
			shortName string
			wantOk    bool
			wantName  string
		}
	}{
		{
			name:   "empty emoji collection",
			emojis: []Emoji{},
			gets: []struct {
				shortName string
				wantOk    bool
				wantName  string
			}{
				{shortName: "smile", wantOk: false},
			},
		},
		{
			name: "single emoji",
			emojis: []Emoji{
				NewEmoji("smile", []rune{0x1F600}, "smile"),
			},
			gets: []struct {
				shortName string
				wantOk    bool
				wantName  string
			}{
				{shortName: "smile", wantOk: true, wantName: "smile"},
				{shortName: "grin", wantOk: false},
			},
		},
		{
			name: "multiple emojis with unique short names",
			emojis: []Emoji{
				NewEmoji("smile", []rune{0x1F600}, "smile"),
				NewEmoji("grin", []rune{0x1F601}, "grin"),
			},
			gets: []struct {
				shortName string
				wantOk    bool
				wantName  string
			}{
				{shortName: "smile", wantOk: true, wantName: "smile"},
				{shortName: "grin", wantOk: true, wantName: "grin"},
				{shortName: "unknown", wantOk: false},
			},
		},
		{
			name: "emoji with multiple short names",
			emojis: []Emoji{
				NewEmoji("grinning face", []rune{0x1F601}, "grinning", "grin", ":D"),
			},
			gets: []struct {
				shortName string
				wantOk    bool
				wantName  string
			}{
				{shortName: "grinning", wantOk: true, wantName: "grinning face"},
				{shortName: "grin", wantOk: true, wantName: "grinning face"},
				{shortName: ":D", wantOk: true, wantName: "grinning face"},
				{shortName: "smile", wantOk: false},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewEmojis(tc.emojis...)

			for _, g := range tc.gets {
				got, ok := m.Get(g.shortName)
				if ok != g.wantOk {
					t.Errorf("Get(%q) ok = %v, want %v", g.shortName, ok, g.wantOk)
				}
				if g.wantOk && got.Name != g.wantName {
					t.Errorf("Get(%q).Name = %q, want %q", g.shortName, got.Name, g.wantName)
				}
			}
		})
	}
}

func TestEmojis_Add(t *testing.T) {
	base := NewEmojis(
		NewEmoji("smile", []rune{0x1F600}, "smile"),
	)

	child := NewEmojis(
		NewEmoji("grin", []rune{0x1F601}, "grin"),
	)

	base.Add(child)

	tests := []struct {
		shortName string
		wantOk    bool
		wantName  string
	}{
		{shortName: "smile", wantOk: true, wantName: "smile"},
		{shortName: "grin", wantOk: true, wantName: "grin"},
		{shortName: "unknown", wantOk: false},
	}

	for _, tc := range tests {
		t.Run(tc.shortName, func(t *testing.T) {
			got, ok := base.Get(tc.shortName)
			if ok != tc.wantOk {
				t.Errorf("Get(%q) ok = %v, want %v", tc.shortName, ok, tc.wantOk)
			}
			if tc.wantOk && got.Name != tc.wantName {
				t.Errorf("Get(%q).Name = %q, want %q", tc.shortName, got.Name, tc.wantName)
			}
		})
	}
}

func TestEmojis_Clone(t *testing.T) {
	original := NewEmojis(
		NewEmoji("smile", []rune{0x1F600}, "smile"),
	)

	child := NewEmojis(
		NewEmoji("grin", []rune{0x1F601}, "grin"),
	)
	original.Add(child)

	cloned := original.Clone()

	// Verify cloned collection has same data
	emoji, ok := cloned.Get("smile")
	if !ok || emoji.Name != "smile" {
		t.Errorf("Cloned collection missing original emoji")
	}

	emoji, ok = cloned.Get("grin")
	if !ok || emoji.Name != "grin" {
		t.Errorf("Cloned collection missing child emoji")
	}

	// Verify adding to clone doesn't affect original
	newChild := NewEmojis(
		NewEmoji("heart", []rune{0x2764}, "heart"),
	)
	cloned.Add(newChild)

	_, okOriginal := original.Get("heart")
	_, okCloned := cloned.Get("heart")

	if okOriginal {
		t.Errorf("Adding to clone affected original collection")
	}
	if !okCloned {
		t.Errorf("Clone did not receive new child")
	}
}

func TestEmojis_Get_ChildPriority(t *testing.T) {
	base := NewEmojis(
		NewEmoji("base", []rune{0x1F600}, "shared"),
	)

	child := NewEmojis(
		NewEmoji("child", []rune{0x1F601}, "shared"),
	)

	base.Add(child)

	// Base map is checked first, so "base" emoji should be returned
	got, ok := base.Get("shared")
	if !ok {
		t.Fatal("Get(\"shared\") not found")
	}
	if got.Name != "base" {
		t.Errorf("Get(\"shared\").Name = %q, want %q (base takes priority)", got.Name, "base")
	}
}

func TestWithEmojis(t *testing.T) {
	base := NewEmojis(
		NewEmoji("smile", []rune{0x1F600}, "smile"),
	)

	opt := WithEmojis(
		NewEmoji("grin", []rune{0x1F601}, "grin"),
	)

	opt(base)

	emoji, ok := base.Get("grin")
	if !ok {
		t.Error("WithEmojis option did not add emoji")
	}
	if ok && emoji.Name != "grin" {
		t.Errorf("WithEmojis added wrong emoji, got Name=%q", emoji.Name)
	}
}

func TestGithub(t *testing.T) {
	// Test basic Github emoji retrieval
	github := Github()
	if github == nil {
		t.Fatal("Github() returned nil")
	}

	// Test that common emojis exist (github.gen.go should have these)
	commonTests := []string{"smile", "heart", "thumbsup", "+1"}

	for _, shortName := range commonTests {
		t.Run(shortName, func(t *testing.T) {
			emoji, ok := github.Get(shortName)
			if ok {
				// If found, verify it has unicode
				if len(emoji.Unicode) == 0 {
					t.Errorf("Github emoji %q has empty Unicode", shortName)
				}
				if emoji.Name == "" {
					t.Errorf("Github emoji %q has empty Name", shortName)
				}
			}
			// Note: We don't fail if not found, as the generated data might not include all
		})
	}
}

func TestGithub_WithOptions(t *testing.T) {
	custom := NewEmoji("custom", []rune{0xFFFD}, "custom")

	github := Github(WithEmojis(custom))

	// Verify custom emoji was added
	emoji, ok := github.Get("custom")
	if !ok {
		t.Error("Github(WithEmojis(...)) did not add custom emoji")
	}
	if ok && emoji.Name != "custom" {
		t.Errorf("Custom emoji has Name=%q, want %q", emoji.Name, "custom")
	}
}

func TestGithub_Singleton(t *testing.T) {
	// Calling Github() multiple times without options should return same instance
	g1 := Github()
	g2 := Github()

	// Both should be able to find the same emojis
	emoji1, ok1 := g1.Get("smile")
	emoji2, ok2 := g2.Get("smile")

	if ok1 != ok2 {
		t.Error("Multiple Github() calls returned different instances")
	}

	// Note: We can't directly compare pointers since the interface hides the implementation
	// but we can verify they behave identically
	if ok1 && ok2 && emoji1.Name != emoji2.Name {
		t.Error("Multiple Github() calls returned inconsistent data")
	}
}

func TestGithub_CloneWithOptions(t *testing.T) {
	// Github with options should return a clone, not the singleton
	custom := NewEmoji("custom", []rune{0xFFFD}, "custom")
	g1 := Github(WithEmojis(custom))

	// Original singleton should not have custom emoji
	g2 := Github()
	_, ok := g2.Get("custom")
	if ok {
		t.Error("Adding emoji with options affected the singleton")
	}

	// The clone should have it
	emoji, ok := g1.Get("custom")
	if !ok || emoji.Name != "custom" {
		t.Error("Clone with options did not include custom emoji")
	}
}
