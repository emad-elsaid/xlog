package xlog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterPreprocessor(t *testing.T) {
	// Save original preprocessors
	original := preprocessors
	defer func() { preprocessors = original }()

	// Reset preprocessors
	preprocessors = []Preprocessor{}

	// Register a simple preprocessor
	upperCasePreprocessor := func(m Markdown) Markdown {
		return Markdown(strings.ToUpper(string(m)))
	}
	RegisterPreprocessor(upperCasePreprocessor)

	assert.Len(t, preprocessors, 1)
}

func TestPreProcess_EmptyInput(t *testing.T) {
	// Save original preprocessors
	original := preprocessors
	defer func() { preprocessors = original }()

	// Reset preprocessors
	preprocessors = []Preprocessor{}

	result := PreProcess("")
	assert.Equal(t, Markdown(""), result)
}

func TestPreProcess_NoPreprocessors(t *testing.T) {
	// Save original preprocessors
	original := preprocessors
	defer func() { preprocessors = original }()

	// Reset preprocessors
	preprocessors = []Preprocessor{}

	input := Markdown("# Hello World")
	result := PreProcess(input)
	assert.Equal(t, input, result)
}

func TestPreProcess_SinglePreprocessor(t *testing.T) {
	// Save original preprocessors
	original := preprocessors
	defer func() { preprocessors = original }()

	// Reset preprocessors
	preprocessors = []Preprocessor{}

	// Register a preprocessor that converts to uppercase
	RegisterPreprocessor(func(m Markdown) Markdown {
		return Markdown(strings.ToUpper(string(m)))
	})

	input := Markdown("# Hello World")
	result := PreProcess(input)
	assert.Equal(t, Markdown("# HELLO WORLD"), result)
}

func TestPreProcess_MultiplePreprocessors(t *testing.T) {
	// Save original preprocessors
	original := preprocessors
	defer func() { preprocessors = original }()

	// Reset preprocessors
	preprocessors = []Preprocessor{}

	// Register multiple preprocessors
	// First: add prefix
	RegisterPreprocessor(func(m Markdown) Markdown {
		return Markdown("PREFIX: " + string(m))
	})

	// Second: convert to uppercase
	RegisterPreprocessor(func(m Markdown) Markdown {
		return Markdown(strings.ToUpper(string(m)))
	})

	// Third: add suffix
	RegisterPreprocessor(func(m Markdown) Markdown {
		return Markdown(string(m) + " :SUFFIX")
	})

	input := Markdown("hello")
	result := PreProcess(input)
	
	// Should be: "hello" -> "PREFIX: hello" -> "PREFIX: HELLO" -> "PREFIX: HELLO :SUFFIX"
	assert.Equal(t, Markdown("PREFIX: HELLO :SUFFIX"), result)
}

func TestPreProcess_PreprocessorOrder(t *testing.T) {
	// Save original preprocessors
	original := preprocessors
	defer func() { preprocessors = original }()

	// Reset preprocessors
	preprocessors = []Preprocessor{}

	// Register preprocessors in specific order to test pipeline
	// First: replace "A" with "B"
	RegisterPreprocessor(func(m Markdown) Markdown {
		return Markdown(strings.ReplaceAll(string(m), "A", "B"))
	})

	// Second: replace "B" with "C"
	RegisterPreprocessor(func(m Markdown) Markdown {
		return Markdown(strings.ReplaceAll(string(m), "B", "C"))
	})

	input := Markdown("AAA")
	result := PreProcess(input)
	
	// Should be: "AAA" -> "BBB" -> "CCC"
	// This proves preprocessors run in order as a pipeline
	assert.Equal(t, Markdown("CCC"), result)
}

func TestPreProcess_PreservesNewlines(t *testing.T) {
	// Save original preprocessors
	original := preprocessors
	defer func() { preprocessors = original }()

	// Reset preprocessors
	preprocessors = []Preprocessor{}

	// Register a preprocessor that adds prefix to each line
	RegisterPreprocessor(func(m Markdown) Markdown {
		lines := strings.Split(string(m), "\n")
		for i, line := range lines {
			lines[i] = "> " + line
		}
		return Markdown(strings.Join(lines, "\n"))
	})

	input := Markdown("line1\nline2\nline3")
	result := PreProcess(input)
	
	expected := Markdown("> line1\n> line2\n> line3")
	assert.Equal(t, expected, result)
}

func TestPreProcess_RealWorldScenario(t *testing.T) {
	// Save original preprocessors
	original := preprocessors
	defer func() { preprocessors = original }()

	// Reset preprocessors
	preprocessors = []Preprocessor{}

	// Simulate real-world preprocessors
	// 1. Expand wiki links [[Page Name]] -> [Page Name](/Page_Name)
	RegisterPreprocessor(func(m Markdown) Markdown {
		content := string(m)
		// Simple wiki link expansion
		content = strings.ReplaceAll(content, "[[Home]]", "[Home](/Home)")
		content = strings.ReplaceAll(content, "[[About]]", "[About](/About)")
		return Markdown(content)
	})

	// 2. Process shortcodes {{date}} -> actual date placeholder
	RegisterPreprocessor(func(m Markdown) Markdown {
		content := string(m)
		content = strings.ReplaceAll(content, "{{date}}", "2026-03-13")
		return Markdown(content)
	})

	input := Markdown("Visit [[Home]] and [[About]] on {{date}}")
	result := PreProcess(input)
	
	expected := Markdown("Visit [Home](/Home) and [About](/About) on 2026-03-13")
	assert.Equal(t, expected, result)
}

func TestPreProcess_IdentityPreprocessor(t *testing.T) {
	// Save original preprocessors
	original := preprocessors
	defer func() { preprocessors = original }()

	// Reset preprocessors
	preprocessors = []Preprocessor{}

	// Register an identity preprocessor (does nothing)
	RegisterPreprocessor(func(m Markdown) Markdown {
		return m
	})

	input := Markdown("# Unchanged Content\n\nThis should remain exactly the same.")
	result := PreProcess(input)
	assert.Equal(t, input, result)
}

func TestPreProcess_EmptyMarkdownType(t *testing.T) {
	// Save original preprocessors
	original := preprocessors
	defer func() { preprocessors = original }()

	// Reset preprocessors
	preprocessors = []Preprocessor{}

	// Register a preprocessor
	RegisterPreprocessor(func(m Markdown) Markdown {
		return Markdown(strings.ToUpper(string(m)))
	})

	// Test with empty Markdown type
	var input Markdown
	result := PreProcess(input)
	assert.Equal(t, Markdown(""), result)
}
