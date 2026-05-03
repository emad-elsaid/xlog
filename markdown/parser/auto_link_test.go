package parser

import (
	"testing"

	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestAutoLinkParser_Trigger(t *testing.T) {
	parser := &autoLinkParser{}
	triggers := parser.Trigger()

	expected := []byte{'<'}
	if len(triggers) != len(expected) {
		t.Fatalf("Trigger() returned %d characters, want %d", len(triggers), len(expected))
	}

	if triggers[0] != expected[0] {
		t.Errorf("Trigger()[0] = %c, want %c", triggers[0], expected[0])
	}
}

func TestNewAutoLinkParser(t *testing.T) {
	parser := NewAutoLinkParser()
	if parser == nil {
		t.Fatal("NewAutoLinkParser() returned nil")
	}

	_, ok := parser.(*autoLinkParser)
	if !ok {
		t.Fatalf("NewAutoLinkParser() returned %T, want *autoLinkParser", parser)
	}
}

func TestAutoLinkParser_Parse(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType ast.AutoLinkType
		expectedText string
		shouldParse  bool
	}{
		{
			name:         "simple URL",
			input:        "<http://example.com>",
			expectedType: ast.AutoLinkURL,
			expectedText: "http://example.com",
			shouldParse:  true,
		},
		{
			name:         "simple email",
			input:        "<user@example.com>",
			expectedType: ast.AutoLinkEmail,
			expectedText: "user@example.com",
			shouldParse:  true,
		},
		{
			name:         "HTTPS URL",
			input:        "<https://secure.example.com>",
			expectedType: ast.AutoLinkURL,
			expectedText: "https://secure.example.com",
			shouldParse:  true,
		},
		{
			name:         "email with subdomain",
			input:        "<alice@mail.example.com>",
			expectedType: ast.AutoLinkEmail,
			expectedText: "alice@mail.example.com",
			shouldParse:  true,
		},
		{
			name:         "URL without closing bracket",
			input:        "<http://example.com",
			expectedType: ast.AutoLinkURL,
			expectedText: "",
			shouldParse:  false,
		},
		{
			name:         "empty brackets",
			input:        "<>",
			expectedType: ast.AutoLinkURL,
			expectedText: "",
			shouldParse:  false,
		},
		{
			name:         "URL with path",
			input:        "<http://example.com/path/to/page>",
			expectedType: ast.AutoLinkURL,
			expectedText: "http://example.com/path/to/page",
			shouldParse:  true,
		},
		{
			name:         "URL with query string",
			input:        "<http://example.com?key=value>",
			expectedType: ast.AutoLinkURL,
			expectedText: "http://example.com?key=value",
			shouldParse:  true,
		},
		{
			name:         "FTP URL",
			input:        "<ftp://files.example.com>",
			expectedType: ast.AutoLinkURL,
			expectedText: "ftp://files.example.com",
			shouldParse:  true,
		},
		{
			name:         "email with plus sign",
			input:        "<user+tag@example.com>",
			expectedType: ast.AutoLinkEmail,
			expectedText: "user+tag@example.com",
			shouldParse:  true,
		},
		{
			name:         "email with hyphen",
			input:        "<user-name@example-domain.com>",
			expectedType: ast.AutoLinkEmail,
			expectedText: "user-name@example-domain.com",
			shouldParse:  true,
		},
		{
			name:         "email with dots",
			input:        "<first.last@example.com>",
			expectedType: ast.AutoLinkEmail,
			expectedText: "first.last@example.com",
			shouldParse:  true,
		},
		{
			name:         "not a link - plain text",
			input:        "<not a link>",
			expectedType: ast.AutoLinkURL,
			expectedText: "",
			shouldParse:  false,
		},
		{
			name:         "URL with fragment",
			input:        "<http://example.com#section>",
			expectedType: ast.AutoLinkURL,
			expectedText: "http://example.com#section",
			shouldParse:  true,
		},
		{
			name:         "URL with port",
			input:        "<http://example.com:8080>",
			expectedType: ast.AutoLinkURL,
			expectedText: "http://example.com:8080",
			shouldParse:  true,
		},
		{
			name:         "localhost URL",
			input:        "<http://localhost:3000>",
			expectedType: ast.AutoLinkURL,
			expectedText: "http://localhost:3000",
			shouldParse:  true,
		},
		{
			name:         "IPv4 URL",
			input:        "<http://192.168.1.1>",
			expectedType: ast.AutoLinkURL,
			expectedText: "http://192.168.1.1",
			shouldParse:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &autoLinkParser{}
			parent := ast.NewDocument()
			reader := text.NewReader([]byte(tc.input))
			context := NewContext()

			result := parser.Parse(parent, reader, context)

			if !tc.shouldParse {
				if result != nil {
					t.Errorf("Parse() should return nil for invalid input, got %T", result)
				}
				return
			}

			if result == nil {
				t.Fatal("Parse() returned nil for valid autolink")
			}

			autoLink, ok := result.(*ast.AutoLink)
			if !ok {
				t.Fatalf("Parse() returned %T, want *ast.AutoLink", result)
			}

			if autoLink.AutoLinkType != tc.expectedType {
				t.Errorf("AutoLinkType = %v, want %v", autoLink.AutoLinkType, tc.expectedType)
			}

			source := []byte(tc.input)
			actualText := extractAutoLinkText(t, autoLink, source)

			if actualText != tc.expectedText {
				t.Errorf("AutoLink text = %q, want %q", actualText, tc.expectedText)
			}
		})
	}
}

func TestAutoLinkParser_Parse_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldParse bool
	}{
		{
			name:        "closing bracket immediately after opening",
			input:       "<>",
			shouldParse: false,
		},
		{
			name:        "text before autolink",
			input:       "text <http://example.com>",
			shouldParse: false,
		},
		{
			name:        "autolink at end of longer string",
			input:       "<http://example.com> more text",
			shouldParse: true,
		},
		{
			name:        "multiple brackets",
			input:       "<<http://example.com>>",
			shouldParse: false,
		},
		{
			name:        "URL with spaces",
			input:       "<http://example.com/path with spaces>",
			shouldParse: false,
		},
		{
			name:        "newline in autolink",
			input:       "<http://example.com\n>",
			shouldParse: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &autoLinkParser{}
			parent := ast.NewDocument()
			reader := text.NewReader([]byte(tc.input))
			context := NewContext()

			result := parser.Parse(parent, reader, context)

			if tc.shouldParse {
				if result == nil {
					t.Error("Parse() returned nil, expected valid autolink")
				} else if _, ok := result.(*ast.AutoLink); !ok {
					t.Errorf("Parse() returned %T, want *ast.AutoLink", result)
				}
			} else {
				if result != nil {
					t.Errorf("Parse() should return nil, got %T", result)
				}
			}
		})
	}
}

func TestAutoLinkParser_Parse_TypePrecedence(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType ast.AutoLinkType
	}{
		{
			name:         "email takes precedence over URL pattern",
			input:        "<user@example.com>",
			expectedType: ast.AutoLinkEmail,
		},
		{
			name:         "URL when no email found",
			input:        "<http://example.com>",
			expectedType: ast.AutoLinkURL,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			parser := &autoLinkParser{}
			parent := ast.NewDocument()
			reader := text.NewReader([]byte(tc.input))
			context := NewContext()

			result := parser.Parse(parent, reader, context)

			if result == nil {
				t.Fatal("Parse() returned nil")
			}

			autoLink, ok := result.(*ast.AutoLink)
			if !ok {
				t.Fatalf("Parse() returned %T, want *ast.AutoLink", result)
			}

			if autoLink.AutoLinkType != tc.expectedType {
				t.Errorf("AutoLinkType = %v, want %v", autoLink.AutoLinkType, tc.expectedType)
			}
		})
	}
}

// Helper function to extract text from AutoLink node.
func extractAutoLinkText(t *testing.T, node *ast.AutoLink, source []byte) string {
	t.Helper()
	return string(node.Label(source))
}
