package xlog

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestHelperRegistration tests that helper registration behavior is preserved
func TestHelperRegistration(t *testing.T) {
	app := newTestApp()
	testHelper := func(s string) string {
		return "test: " + s
	}

	err := app.RegisterHelper("testHelper", testHelper)
	require.NoError(t, err, "Expected no error")
	require.NotNil(t, app.helpers["testHelper"], "Expected helper to be registered")

	err = app.RegisterHelper("testHelper", testHelper)
	require.ErrorIs(t, err, ErrHelperRegistered, "Expected ErrHelperRegistered")
}

// TestJavaScriptHandling tests that JavaScript handling behavior is preserved
func TestJavaScriptHandling(t *testing.T) {
	app := newTestApp()
	app.js = []string{}

	result := app.includeJS("/test.js")
	require.Equal(t, template.HTML(""), result, "Expected empty result")
	require.Len(t, app.js, 1, "Expected 1 JS file")
	require.Equal(t, "/test.js", app.js[0], "Expected '/test.js'")

	app.includeJS("/test.js")
	require.Len(t, app.js, 1, "Expected 1 JS file after duplicate")

	app.includeJS("/another.js")
	scripts := app.scripts()
	require.Equal(t, `<script src="/test.js" defer></script><script src="/another.js" defer></script>`, string(scripts))
}
