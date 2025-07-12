package xlog

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsIgnoredPathBehavior(t *testing.T) {
	app := newTestApp()

	require.True(t, app.IsIgnoredPath(".git"), "Expected .git to be ignored")
	require.True(t, app.IsIgnoredPath(".hidden"), "Expected .hidden to be ignored")
	require.False(t, app.IsIgnoredPath("normal"), "Expected 'normal' to not be ignored")
	require.False(t, app.IsIgnoredPath("normal/path"), "Expected 'normal/path' to not be ignored")

	customPattern := regexp.MustCompile(`^temp`)
	app.IgnorePath(customPattern)

	require.True(t, app.IsIgnoredPath("temp"), "Expected 'temp' to be ignored after registering custom pattern")
	require.True(t, app.IsIgnoredPath("tempfile.txt"), "Expected 'tempfile.txt' to be ignored after registering custom pattern")
	require.False(t, app.IsIgnoredPath("permanent"), "Expected 'permanent' to not be ignored")
}

func TestIsNil(t *testing.T) {
	require.True(t, isNil[Page](nil))
	require.True(t, isNil[*Page](nil))
}
