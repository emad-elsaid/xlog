package xlog

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFrontmatterTitleRendering tests that pages with frontmatter titles
// render correctly without type errors. This is a regression test for issue #134.
//
// The bug occurred when a page had frontmatter with a title property.
// The template tried to call emoji() on the MetaProperty instead of the Page,
// causing: "wrong type for value; expected xlog.Page; got frontmatter.MetaProperty"
func TestFrontmatterTitleRendering(t *testing.T) {
	// Store original state
	originalTemplates := templates
	originalConfig := Config
	defer func() {
		templates = originalTemplates
		Config = originalConfig
	}()

	// Create a temporary directory for test pages
	tmpDir := t.TempDir()
	
	// Create a test markdown file with frontmatter
	testFile := filepath.Join(tmpDir, "test.md")
	content := `---
title: My Custom Page Title
tags:
 - golang
 - notes
author: John
---

# Heading

This is the content of the page.
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Set up configuration
	Config.Source = tmpDir
	Config.Index = "index"
	Config.Sitename = "Test Site"

	// Load the page
	page := NewPage(testFile)

	// Compile templates
	compileTemplates()

	// Create template data matching the actual structure used in xlog
	data := map[string]interface{}{
		"page":   page,
		"config": Config,
	}

	// Try to render just the header template (contains the bug location at line 38)
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "header", data)
	
	if err != nil {
		// Check if this is the specific error we're testing for
		if strings.Contains(err.Error(), "wrong type for value; expected xlog.Page; got frontmatter.MetaProperty") {
			t.Fatalf("Regression detected (issue #134): %v\n\nThe template is incorrectly passing MetaProperty to emoji() instead of Page.\nLocation: templates/layout.html line 38", err)
		}
		if strings.Contains(err.Error(), "executing \"header\"") {
			t.Fatalf("Template execution failed (possible regression of #134): %v", err)
		}
		t.Fatalf("Failed to render template: %v", err)
	}

	// Verify no error messages in output
	output := buf.String()
	if strings.Contains(output, "rendering error") {
		t.Errorf("Rendered output contains 'rendering error': %s", output)
	}
}

// TestFrontmatterWithoutTitle tests that pages without frontmatter title
// still render correctly (baseline test).
func TestFrontmatterWithoutTitle(t *testing.T) {
	originalTemplates := templates
	originalConfig := Config
	defer func() {
		templates = originalTemplates
		Config = originalConfig
	}()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "no-title.md")
	content := `---
tags:
 - test
---

# Content

No title in frontmatter.
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	Config.Source = tmpDir
	Config.Index = "index"
	Config.Sitename = "Test Site"

	page := NewPage(testFile)
	compileTemplates()

	data := map[string]interface{}{
		"page":   page,
		"config": Config,
	}

	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "header", data)
	
	if err != nil {
		t.Fatalf("Failed to render template without frontmatter title: %v", err)
	}
}
