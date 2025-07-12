package xlog

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPropertiesBehavior tests that Properties function behavior is preserved
func TestPropertiesBehavior(t *testing.T) {
	app := newTestApp()

	// Create a temporary file to test with
	tempFile, err := os.CreateTemp("", "test-*.md")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Get the file's ModTime
	stat, err := tempFile.Stat()
	require.NoError(t, err)
	modTime := stat.ModTime()

	// Create a page with the temporary file's name (without the .md extension)
	pageName := strings.TrimSuffix(tempFile.Name(), ".md")
	page := &page{name: pageName}

	// Get the page's properties
	props := app.Properties(page)

	// Check that the ModTime property is correct
	require.Len(t, props, 1, "Expected 1 property")
	require.Contains(t, props, "modified", "Expected properties to contain 'modified'")
	require.Equal(t, app.ago(modTime), props["modified"].Value(), "Expected property value to be the file's modification time")
}
