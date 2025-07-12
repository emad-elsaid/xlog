package xlog

import (
	"html/template"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAgoBehavior(t *testing.T) {
	app := newTestApp()

	app.config.Readonly = true
	now := time.Now()
	result := app.ago(now)
	require.Equal(t, now.Format("Monday 2 January 2006"), result)

	app.config.Readonly = false

	recent := time.Now().Add(-500 * time.Millisecond)
	result = app.ago(recent)
	require.Contains(t, result, "Less than a second")

	oneMinuteAgo := time.Now().Add(-1 * time.Minute)
	result = app.ago(oneMinuteAgo)
	require.Contains(t, result, "1 minute")

	oneHourAgo := time.Now().Add(-1 * time.Hour)
	result = app.ago(oneHourAgo)
	require.Contains(t, result, "1 hour")

	oneDayAgo := time.Now().Add(-24 * time.Hour)
	result = app.ago(oneDayAgo)
	require.Contains(t, result, "1 day")
}

func TestIsFontAwesome(t *testing.T) {
	require.True(t, IsFontAwesome("fa-solid"), "Expected 'fa-solid' to be FontAwesome")
	require.True(t, IsFontAwesome("fa-regular"), "Expected 'fa-regular' to be FontAwesome")
	require.True(t, IsFontAwesome("fa-brands"), "Expected 'fa-brands' to be FontAwesome")
	require.False(t, IsFontAwesome("not-fa"), "Expected 'not-fa' to not be FontAwesome")
	require.False(t, IsFontAwesome(""), "Expected empty string to not be FontAwesome")
}

func TestDir(t *testing.T) {
	require.Equal(t, "", dir(""), "Expected empty string for empty path")
	require.Equal(t, "", dir("."), "Expected empty string for '.'")
	require.Equal(t, "", dir("file.txt"), "Expected empty string for 'file.txt'")
	require.Equal(t, "dir", dir("dir/file.txt"), "Expected 'dir' for 'dir/file.txt'")
	require.Equal(t, "a/b", dir("a/b/c.txt"), "Expected 'a/b' for 'a/b/c.txt'")
}

func TestRaw(t *testing.T) {
	input := "<div>test</div>"
	result := raw(input)
	require.Equal(t, template.HTML(input), result, "Expected HTML to match input")
}
