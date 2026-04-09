package embed

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/emad-elsaid/xlog"
)

func TestEmbedExtensionName(t *testing.T) {
	ext := Embed{}
	expected := "embed"
	if ext.Name() != expected {
		t.Errorf("Expected name %q, got %q", expected, ext.Name())
	}
}

func TestEmbedInit(t *testing.T) {
	// Test that Init doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Init() panicked: %v", r)
		}
	}()
	
	ext := Embed{}
	ext.Init()
}

func TestEmbedShortcode_NonExistentPage(t *testing.T) {
	// Test embedding a non-existent page
	// This test uses a name that's highly unlikely to exist
	input := Markdown("non-existent-page-9876543210")
	result := embedShortcode(input)

	resultStr := string(result)
	// Should contain error message about page not existing
	if !strings.Contains(resultStr, "doesn't exist") {
		t.Errorf("Expected error message for non-existent page, got: %s", resultStr)
	}
}

func TestEmbedShortcode_EmptyInput(t *testing.T) {
	input := Markdown("")
	result := embedShortcode(input)

	expected := template.HTML("Page:  doesn't exist")
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestEmbedShortcode_WhitespaceTrimming(t *testing.T) {
	// Test that whitespace is trimmed for page lookup
	// The trimming happens with strings.TrimSpace in the embedShortcode function
	tests := []struct {
		name  string
		input string
	}{
		{"leading space", " non-existent-test-page-xyz"},
		{"trailing space", "non-existent-test-page-xyz "},
		{"both spaces", " non-existent-test-page-xyz "},
		{"tab", "\tnon-existent-test-page-xyz"},
		{"newline", "non-existent-test-page-xyz\n"},
		{"multiple whitespace", "  non-existent-test-page-xyz  \n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := embedShortcode(Markdown(tt.input))
			resultStr := string(result)

			// Should still report page doesn't exist (trimming happens for lookup)
			if !strings.Contains(resultStr, "doesn't exist") {
				t.Errorf("Expected 'doesn't exist' message, got: %s", resultStr)
			}
			// The page name in the message reflects original input
			if !strings.Contains(resultStr, "non-existent-test-page-xyz") {
				t.Errorf("Expected page name in error message, got: %s", resultStr)
			}
		})
	}
}

func TestEmbedShortcode_Integration(t *testing.T) {
	// Create a temporary test page file for integration testing
	// This tests the full flow with actual filesystem
	tmpDir := t.TempDir()
	
	// Save current directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)
	
	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	
	// Create a test markdown page
	testPageName := "test-embed-page"
	testContent := "# Test Embed Page\n\nThis is test content for embedding."
	testFilePath := filepath.Join(tmpDir, testPageName+".md")
	
	if err := os.WriteFile(testFilePath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}
	
	// Test embedding the page
	input := Markdown(testPageName)
	result := embedShortcode(input)
	
	resultStr := string(result)
	
	// Verify the result is not empty
	if resultStr == "" {
		t.Error("Expected non-empty result for existing page")
	}
	
	// Verify it's not an error message
	if strings.Contains(resultStr, "doesn't exist") {
		t.Errorf("Expected page to be found, got error: %s", resultStr)
	}
	
	// Verify the content is rendered (should contain heading)
	if !strings.Contains(resultStr, "Test Embed Page") {
		t.Errorf("Expected rendered content to contain page heading, got: %s", resultStr)
	}
}

func TestEmbedShortcode_Integration_WithWhitespace(t *testing.T) {
	// Test that pages with whitespace in input are found correctly
	tmpDir := t.TempDir()
	
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)
	
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	
	testPageName := "whitespace-page"
	testContent := "# Whitespace Test"
	testFilePath := filepath.Join(tmpDir, testPageName+".md")
	
	if err := os.WriteFile(testFilePath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test page: %v", err)
	}
	
	// Test with various whitespace - trimming should make it work
	inputs := []string{
		"whitespace-page",
		" whitespace-page",
		"whitespace-page ",
		" whitespace-page ",
		"\twhitespace-page",
		"whitespace-page\n",
	}
	
	for _, input := range inputs {
		t.Run(fmt.Sprintf("input_%q", input), func(t *testing.T) {
			result := embedShortcode(Markdown(input))
			resultStr := string(result)
			
			if strings.Contains(resultStr, "doesn't exist") {
				t.Errorf("Expected page to be found with input %q, got error: %s", input, resultStr)
			}
			
			if !strings.Contains(resultStr, "Whitespace Test") {
				t.Errorf("Expected rendered content with input %q, got: %s", input, resultStr)
			}
		})
	}
}

func TestEmbedShortcode_OutputFormat(t *testing.T) {
	// Test that output is valid HTML (template.HTML type)
	input := Markdown("some-page")
	result := embedShortcode(input)
	
	// Verify the result is of type template.HTML
	_, ok := interface{}(result).(template.HTML)
	if !ok {
		t.Errorf("Expected result to be template.HTML, got %T", result)
	}
}

func TestEmbedShortcode_ErrorMessageFormat(t *testing.T) {
	// Test that error messages follow expected format
	testCases := []struct {
		input    string
		expected string
	}{
		{"test-page-xyz123", "Page: test-page-xyz123 doesn't exist"},
		{"another-page-abc456", "Page: another-page-abc456 doesn't exist"},
		{"", "Page:  doesn't exist"},
	}
	
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("input_%s", tc.input), func(t *testing.T) {
			result := embedShortcode(Markdown(tc.input))
			resultStr := string(result)
			
			// For pages that don't exist, verify error format
			if strings.Contains(resultStr, "doesn't exist") {
				if resultStr != tc.expected {
					t.Errorf("Expected error %q, got %q", tc.expected, resultStr)
				}
			}
		})
	}
}

func TestEmbedShortcode_ComplexPageName(t *testing.T) {
	// Test various page name formats
	tmpDir := t.TempDir()
	
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldWd)
	
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	
	// Test with page name containing various characters
	testCases := []string{
		"simple-page",
		"page_with_underscores",
		"page-with-dashes",
		"page123",
		"PageWithCaps",
	}
	
	for _, pageName := range testCases {
		t.Run(pageName, func(t *testing.T) {
			// Create the page
			content := fmt.Sprintf("# %s", pageName)
			filePath := filepath.Join(tmpDir, pageName+".md")
			
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to create test page: %v", err)
			}
			
			// Test embedding
			result := embedShortcode(Markdown(pageName))
			resultStr := string(result)
			
			if strings.Contains(resultStr, "doesn't exist") {
				t.Errorf("Expected page %q to be found, got error", pageName)
			}
			
			if !strings.Contains(resultStr, pageName) {
				t.Errorf("Expected rendered content to contain page name %q, got: %s", pageName, resultStr)
			}
		})
	}
}
