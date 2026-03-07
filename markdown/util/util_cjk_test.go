package util

import (
	"testing"
)

func TestIsEastAsianWideRune(t *testing.T) {
	tests := []struct {
		name     string
		r        rune
		expected bool
	}{
		// Hiragana
		{
			name:     "hiragana あ",
			r:        'あ',
			expected: true,
		},
		{
			name:     "hiragana ん",
			r:        'ん',
			expected: true,
		},
		
		// Katakana
		{
			name:     "katakana ア",
			r:        'ア',
			expected: true,
		},
		{
			name:     "katakana ン",
			r:        'ン',
			expected: true,
		},
		
		// Han (Chinese characters)
		{
			name:     "Han 中",
			r:        '中',
			expected: true,
		},
		{
			name:     "Han 文",
			r:        '文',
			expected: true,
		},
		
		// Hangul (Korean)
		{
			name:     "Hangul 한",
			r:        '한',
			expected: true,
		},
		{
			name:     "Hangul 글",
			r:        '글',
			expected: true,
		},
		
		// CJK symbols and punctuation
		{
			name:     "CJK ideographic comma 、",
			r:        '、',
			expected: true,
		},
		{
			name:     "CJK fullwidth period 。",
			r:        '。',
			expected: true,
		},
		
		// Latin characters (not wide)
		{
			name:     "ASCII letter A",
			r:        'A',
			expected: false,
		},
		{
			name:     "ASCII letter z",
			r:        'z',
			expected: false,
		},
		
		// ASCII symbols (not wide)
		{
			name:     "ASCII space",
			r:        ' ',
			expected: false,
		},
		{
			name:     "ASCII exclamation !",
			r:        '!',
			expected: false,
		},
		
		// Digits (not wide)
		{
			name:     "ASCII digit 0",
			r:        '0',
			expected: false,
		},
		{
			name:     "ASCII digit 9",
			r:        '9',
			expected: false,
		},
		
		// Greek (not wide)
		{
			name:     "Greek alpha α",
			r:        'α',
			expected: false,
		},
		
		// Cyrillic (not wide)
		{
			name:     "Cyrillic А",
			r:        'А',
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEastAsianWideRune(tt.r)
			if result != tt.expected {
				t.Errorf("IsEastAsianWideRune(%q) = %v, want %v", tt.r, result, tt.expected)
			}
		})
	}
}

func TestIsSpaceDiscardingUnicodeRune(t *testing.T) {
	tests := []struct {
		name     string
		r        rune
		expected bool
	}{
		// CJK Radicals Supplement (2E80–2EFF)
		{
			name:     "CJK Radical Supplement start",
			r:        0x2E80,
			expected: true,
		},
		{
			name:     "CJK Radical Supplement end",
			r:        0x2EFF,
			expected: true,
		},
		
		// Kangxi Radicals (2F00–2FDF)
		{
			name:     "Kangxi Radical start",
			r:        0x2F00,
			expected: true,
		},
		
		// Hiragana (3040–309F)
		{
			name:     "Hiragana あ",
			r:        'あ',
			expected: true,
		},
		{
			name:     "Hiragana ん",
			r:        'ん',
			expected: true,
		},
		
		// Katakana (30A0–30FF)
		{
			name:     "Katakana ア",
			r:        'ア',
			expected: true,
		},
		{
			name:     "Katakana ン",
			r:        'ン',
			expected: true,
		},
		
		// CJK Unified Ideographs (4E00–9FFF)
		{
			name:     "CJK Unified Ideograph 中",
			r:        '中',
			expected: true,
		},
		{
			name:     "CJK Unified Ideograph 文",
			r:        '文',
			expected: true,
		},
		
		// CJK Unified Ideographs Extension A (3400–4DBF)
		{
			name:     "CJK Extension A start",
			r:        0x3400,
			expected: true,
		},
		{
			name:     "CJK Extension A end",
			r:        0x4DBF,
			expected: true,
		},
		
		// CJK Unified Ideographs Extension B (20000–2A6DF)
		{
			name:     "CJK Extension B start",
			r:        0x20000,
			expected: true,
		},
		{
			name:     "CJK Extension B sample",
			r:        0x20001,
			expected: true,
		},
		
		// CJK Unified Ideographs Extension G (30000–3134F)
		{
			name:     "CJK Extension G start",
			r:        0x30000,
			expected: true,
		},
		
		// Halfwidth and Fullwidth Forms (FF00–FFEF)
		{
			name:     "Fullwidth Latin A",
			r:        0xFF21,
			expected: true,
		},
		{
			name:     "Fullwidth digit 0",
			r:        0xFF10,
			expected: true,
		},
		
		// Yi Syllables (A000–A48F)
		{
			name:     "Yi Syllable",
			r:        0xA000,
			expected: true,
		},
		
		// ASCII (not space-discarding)
		{
			name:     "ASCII letter A",
			r:        'A',
			expected: false,
		},
		{
			name:     "ASCII space",
			r:        ' ',
			expected: false,
		},
		{
			name:     "ASCII digit 5",
			r:        '5',
			expected: false,
		},
		
		// Latin Extended (not space-discarding)
		{
			name:     "Latin Extended é",
			r:        'é',
			expected: false,
		},
		
		// Greek (not space-discarding)
		{
			name:     "Greek alpha α",
			r:        'α',
			expected: false,
		},
		
		// Cyrillic (not space-discarding)
		{
			name:     "Cyrillic А",
			r:        'А',
			expected: false,
		},
		
		// Arabic (not space-discarding)
		{
			name:     "Arabic alef ا",
			r:        'ا',
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSpaceDiscardingUnicodeRune(tt.r)
			if result != tt.expected {
				t.Errorf("IsSpaceDiscardingUnicodeRune(%q/U+%04X) = %v, want %v", tt.r, tt.r, result, tt.expected)
			}
		})
	}
}

func TestEastAsianWidth(t *testing.T) {
	tests := []struct {
		name     string
		r        rune
		expected string
	}{
		// Fullwidth (F) - characters that take up full width
		{
			name:     "Ideographic space",
			r:        0x3000,
			expected: "F",
		},
		{
			name:     "Fullwidth exclamation",
			r:        0xFF01,
			expected: "F",
		},
		{
			name:     "Fullwidth tilde",
			r:        0xFF5E,
			expected: "F",
		},
		{
			name:     "Fullwidth won sign",
			r:        0xFFE0,
			expected: "F",
		},
		
		// Halfwidth (H) - characters that take up half width
		{
			name:     "Won sign",
			r:        0x20A9,
			expected: "H",
		},
		{
			name:     "Halfwidth katakana ｱ",
			r:        0xFF71,
			expected: "H",
		},
		{
			name:     "Halfwidth katakana ﾝ",
			r:        0xFFDC, // Use 0xFFDC which is in the halfwidth range
			expected: "H",
		},
		
		// Wide (W) - East Asian wide characters
		{
			name:     "Hiragana あ",
			r:        'あ',
			expected: "W",
		},
		{
			name:     "Katakana ア",
			r:        'ア',
			expected: "W",
		},
		{
			name:     "Han 中",
			r:        '中',
			expected: "W",
		},
		{
			name:     "Hangul 한",
			r:        '한',
			expected: "W",
		},
		{
			name:     "CJK left corner bracket",
			r:        0x300C,
			expected: "W",
		},
		{
			name:     "CJK Extension B",
			r:        0x20000,
			expected: "W",
		},
		
		// Narrow (Na) - narrow ASCII
		{
			name:     "ASCII space",
			r:        ' ',
			expected: "Na",
		},
		{
			name:     "ASCII letter A",
			r:        'A',
			expected: "Na",
		},
		{
			name:     "ASCII letter z",
			r:        'z',
			expected: "Na",
		},
		{
			name:     "ASCII tilde ~",
			r:        '~',
			expected: "Na",
		},
		{
			name:     "Cent sign ¢",
			r:        0x00A2,
			expected: "Na",
		},
		{
			name:     "Pound sign £",
			r:        0x00A3,
			expected: "Na",
		},
		{
			name:     "Yen sign ¥",
			r:        0x00A5,
			expected: "Na",
		},
		
		// Ambiguous (A) - characters with ambiguous width
		{
			name:     "Inverted exclamation ¡",
			r:        0x00A1,
			expected: "A",
		},
		{
			name:     "Currency sign ¤",
			r:        0x00A4,
			expected: "A",
		},
		{
			name:     "Section sign §",
			r:        0x00A7,
			expected: "A",
		},
		{
			name:     "Diaeresis ¨",
			r:        0x00A8,
			expected: "A",
		},
		{
			name:     "Degree sign °",
			r:        0x00B0,
			expected: "A",
		},
		{
			name:     "Greek capital alpha Α",
			r:        0x0391,
			expected: "A",
		},
		{
			name:     "Greek capital omega Ω",
			r:        0x03A9,
			expected: "A",
		},
		{
			name:     "Cyrillic capital A А",
			r:        0x0410,
			expected: "A",
		},
		{
			name:     "En dash –",
			r:        0x2013,
			expected: "A",
		},
		{
			name:     "Em dash —",
			r:        0x2014,
			expected: "A",
		},
		{
			name:     "Left single quote '",
			r:        0x2018,
			expected: "A",
		},
		{
			name:     "Left double quote",
			r:        0x201C,
			expected: "A",
		},
		{
			name:     "Bullet •",
			r:        0x2022,
			expected: "A",
		},
		{
			name:     "Euro sign €",
			r:        0x20AC,
			expected: "A",
		},
		{
			name:     "Degree Celsius ℃",
			r:        0x2103,
			expected: "A",
		},
		{
			name:     "Number sign №",
			r:        0x2116,
			expected: "A",
		},
		{
			name:     "Trademark ™",
			r:        0x2122,
			expected: "A",
		},
		{
			name:     "One third ⅓",
			r:        0x2153,
			expected: "A",
		},
		{
			name:     "Roman numeral one Ⅰ",
			r:        0x2160,
			expected: "A",
		},
		{
			name:     "Leftward arrow ←",
			r:        0x2190,
			expected: "A",
		},
		{
			name:     "For all ∀",
			r:        0x2200,
			expected: "A",
		},
		{
			name:     "Partial differential ∂",
			r:        0x2202,
			expected: "A",
		},
		{
			name:     "Not equal ≠",
			r:        0x2260,
			expected: "A",
		},
		{
			name:     "Black square ■",
			r:        0x25A0,
			expected: "A",
		},
		{
			name:     "White circle ○",
			r:        0x25CB,
			expected: "A",
		},
		{
			name:     "Black star ★",
			r:        0x2605,
			expected: "A",
		},
		
		// Neutral (N) - all other characters
		{
			name:     "Copyright sign ©",
			r:        0x00A9,
			expected: "N",
		},
		{
			name:     "Latin Extended é",
			r:        'é',
			expected: "A", // Actually ambiguous width
		},
		{
			name:     "Arabic alef ا",
			r:        'ا',
			expected: "N",
		},
		{
			name:     "Hebrew alef א",
			r:        'א',
			expected: "N",
		},
		{
			name:     "Thai character ก",
			r:        'ก',
			expected: "N",
		},
		{
			name:     "Devanagari क",
			r:        'क',
			expected: "N",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EastAsianWidth(tt.r)
			if result != tt.expected {
				t.Errorf("EastAsianWidth(%q/U+%04X) = %s, want %s", tt.r, tt.r, result, tt.expected)
			}
		})
	}
}

func BenchmarkIsEastAsianWideRune(b *testing.B) {
	testRunes := []rune{
		'A',    // ASCII (false)
		'あ',   // Hiragana (true)
		'ア',   // Katakana (true)
		'中',   // Han (true)
		'한',   // Hangul (true)
		'α',    // Greek (false)
		0x3000, // CJK space (true)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range testRunes {
			IsEastAsianWideRune(r)
		}
	}
}

func BenchmarkIsSpaceDiscardingUnicodeRune(b *testing.B) {
	testRunes := []rune{
		'A',     // ASCII (false)
		'あ',    // Hiragana (true)
		'中',    // CJK Unified Ideograph (true)
		0x2E80,  // CJK Radicals Supplement (true)
		0x20000, // CJK Extension B (true)
		'α',     // Greek (false)
		0xFF21,  // Fullwidth Latin (true)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range testRunes {
			IsSpaceDiscardingUnicodeRune(r)
		}
	}
}

func BenchmarkEastAsianWidth(b *testing.B) {
	testRunes := []rune{
		'A',     // Na
		'あ',    // W
		0x3000,  // F
		0xFF71,  // H
		0x00A1,  // A
		'é',     // N
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range testRunes {
			EastAsianWidth(r)
		}
	}
}
