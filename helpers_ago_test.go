package xlog

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAgo(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "less than a second",
			time:     now.Add(-500 * time.Millisecond),
			expected: "Less than a second ago",
		},
		{
			name:     "30 seconds ago",
			time:     now.Add(-30 * time.Second),
			expected: "30 seconds ago",
		},
		{
			name:     "1 minute ago",
			time:     now.Add(-1 * time.Minute),
			expected: "1 minutes ago",
		},
		{
			name:     "5 minutes 30 seconds ago",
			time:     now.Add(-5*time.Minute - 30*time.Second),
			expected: "5 minutes 30 seconds ago",
		},
		{
			name:     "1 hour ago",
			time:     now.Add(-1 * time.Hour),
			expected: "1 hours ago",
		},
		{
			name:     "2 hours 45 minutes ago",
			time:     now.Add(-2*time.Hour - 45*time.Minute),
			expected: "2 hours 45 minutes ago",
		},
		{
			name:     "1 day ago",
			time:     now.Add(-24 * time.Hour),
			expected: "1 days ago",
		},
		{
			name:     "3 days 5 hours ago",
			time:     now.Add(-3*24*time.Hour - 5*time.Hour),
			expected: "3 days 5 hours ago",
		},
		{
			name:     "1 week ago",
			time:     now.Add(-7 * 24 * time.Hour),
			expected: "1 weeks ago",
		},
		{
			name:     "2 weeks 3 days ago",
			time:     now.Add(-2*7*24*time.Hour - 3*24*time.Hour),
			expected: "2 weeks 3 days ago",
		},
		{
			name:     "1 month ago",
			time:     now.Add(-30 * 24 * time.Hour),
			expected: "1 months ago",
		},
		{
			name:     "2 months 1 week ago",
			time:     now.Add(-2*30*24*time.Hour - 7*24*time.Hour),
			expected: "2 months 1 weeks ago",
		},
		{
			name:     "1 year ago",
			time:     now.Add(-365 * 24 * time.Hour),
			expected: "1 years ago",
		},
		{
			name:     "2 years 3 months ago",
			time:     now.Add(-2*365*24*time.Hour - 3*30*24*time.Hour),
			expected: "2 years 3 months ago",
		},
		{
			name:     "5 years 2 months ago",
			time:     now.Add(-5*365*24*time.Hour - 2*30*24*time.Hour),
			expected: "5 years 2 months ago",
		},
	}

	// Save original readonly state and restore after test
	originalReadonly := Config.Readonly
	defer func() { Config.Readonly = originalReadonly }()

	// Test in non-readonly mode (relative time)
	Config.Readonly = false

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ago(tc.time)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestAgoReadonlyMode(t *testing.T) {
	// Save original readonly state and restore after test
	originalReadonly := Config.Readonly
	defer func() { Config.Readonly = originalReadonly }()

	// Test in readonly mode (absolute date)
	Config.Readonly = true

	testTime := time.Date(2023, time.March, 15, 14, 30, 0, 0, time.UTC)
	result := ago(testTime)

	// In readonly mode, should return formatted date
	require.Equal(t, "Wednesday 15 March 2023", result)
}

func TestAgoPrecision(t *testing.T) {
	// Save original readonly state and restore after test
	originalReadonly := Config.Readonly
	defer func() { Config.Readonly = originalReadonly }()
	Config.Readonly = false

	now := time.Now()

	// Test that precision is limited to 2 units
	// A very complex duration should still only show 2 time units
	complexTime := now.Add(-3*365*24*time.Hour - 7*30*24*time.Hour - 2*7*24*time.Hour - 4*24*time.Hour - 6*time.Hour - 30*time.Minute - 45*time.Second)
	result := ago(complexTime)

	// Should show exactly 2 units (years and months), with no more than 2 unit words
	// Count the number of time units in the output
	unitCount := 0
	for _, unit := range []string{"years", "months", "weeks", "days", "hours", "minutes", "seconds"} {
		if strings.Contains(result, unit) {
			unitCount++
		}
	}

	// Should have exactly 2 time units
	require.Equal(t, 2, unitCount, "Expected exactly 2 time units in: %s", result)

	// Should end with "ago"
	require.True(t, strings.HasSuffix(result, "ago"), "Expected result to end with 'ago': %s", result)
}

func TestAgoEdgeCases(t *testing.T) {
	// Save original readonly state and restore after test
	originalReadonly := Config.Readonly
	defer func() { Config.Readonly = originalReadonly }()
	Config.Readonly = false

	now := time.Now()

	t.Run("exact second boundary", func(t *testing.T) {
		result := ago(now.Add(-1 * time.Second))
		require.Equal(t, "1 seconds ago", result)
	})

	t.Run("exact minute boundary", func(t *testing.T) {
		result := ago(now.Add(-60 * time.Second))
		require.Equal(t, "1 minutes ago", result)
	})

	t.Run("exact hour boundary", func(t *testing.T) {
		result := ago(now.Add(-60 * time.Minute))
		require.Equal(t, "1 hours ago", result)
	})

	t.Run("just under a second", func(t *testing.T) {
		result := ago(now.Add(-999 * time.Millisecond))
		require.Equal(t, "Less than a second ago", result)
	})
}
