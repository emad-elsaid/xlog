package xlog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPreProcessBehavior tests that PreProcess function behavior is preserved
func TestPreProcessBehavior(t *testing.T) {
	app := newTestApp()
	app.preprocessors = []Preprocessor{}

	content := Markdown("test content")
	result := app.PreProcess(content)
	require.Equal(t, content, result, "Expected unchanged content")

	preprocessor := func(content Markdown) Markdown {
		return Markdown("processed: " + string(content))
	}

	app.RegisterPreprocessor(preprocessor)
	result = app.PreProcess(content)
	require.Equal(t, Markdown("processed: test content"), result)

	secondPreprocessor := func(content Markdown) Markdown {
		return Markdown(string(content) + " (final)")
	}

	app.RegisterPreprocessor(secondPreprocessor)
	result = app.PreProcess(content)
	require.Equal(t, Markdown("processed: test content (final)"), result)
}
