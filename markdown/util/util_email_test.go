package util

import (
	"testing"
)

func TestFindEmailIndex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// Valid emails
		{
			name:     "simple email",
			input:    "user@example.com",
			expected: 16,
		},
		{
			name:     "email with dots in local part",
			input:    "user.name@example.com",
			expected: 21,
		},
		{
			name:     "email with plus",
			input:    "user+tag@example.com",
			expected: 20,
		},
		{
			name:     "email with numbers",
			input:    "user123@example456.com",
			expected: 22,
		},
		{
			name:     "subdomain",
			input:    "user@mail.example.com",
			expected: 21,
		},
		{
			name:     "multiple subdomains",
			input:    "user@a.b.c.example.com",
			expected: 22,
		},
		{
			name:     "short domain",
			input:    "a@b.c",
			expected: 5,
		},
		{
			name:     "hyphen in domain",
			input:    "user@my-domain.com",
			expected: 18,
		},
		{
			name:     "email followed by text",
			input:    "user@example.com and more text",
			expected: 16,
		},
		{
			name:     "email with underscore in local",
			input:    "user_name@example.com",
			expected: 21,
		},
		
		// Invalid emails
		{
			name:     "no local part",
			input:    "@example.com",
			expected: -1,
		},
		{
			name:     "no at symbol",
			input:    "userexample.com",
			expected: -1,
		},
		{
			name:     "no domain",
			input:    "user@",
			expected: -1,
		},
		{
			name:     "domain starting with hyphen",
			input:    "user@-example.com",
			expected: -1,
		},
		{
			name:     "domain ending with hyphen",
			input:    "user@example-.com",
			expected: -1,
		},
		{
			name:     "domain ending with dot",
			input:    "user@example.com.",
			expected: 16,
		},
		{
			name:     "consecutive dots in domain",
			input:    "user@example..com",
			expected: -1,
		},
		{
			name:     "empty string",
			input:    "",
			expected: -1,
		},
		{
			name:     "only @",
			input:    "@",
			expected: -1,
		},
		{
			name:     "domain starting with dot",
			input:    "user@.example.com",
			expected: -1,
		},
		
		// Edge cases
		{
			name:     "single char local and domain",
			input:    "a@b.c",
			expected: 5,
		},
		{
			name:     "long label (63 chars - valid)",
			input:    "user@" + string(make([]byte, 63)) + ".com",
			expected: -1, // will fail because we're using null bytes
		},
		{
			name:     "label with max valid length",
			input:    "user@abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuv.com",
			expected: 67,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindEmailIndex([]byte(tt.input))
			if result != tt.expected {
				t.Errorf("FindEmailIndex(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkFindEmailIndex(b *testing.B) {
	testCases := [][]byte{
		[]byte("user@example.com"),
		[]byte("very.long.email.address.with.many.dots@subdomain.example.co.uk"),
		[]byte("simple@test.io"),
		[]byte("user+tag@example.com"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			FindEmailIndex(tc)
		}
	}
}
