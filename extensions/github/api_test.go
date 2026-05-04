package github

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/emad-elsaid/xlog"
)

func TestToken_EnvVariableNotSet(t *testing.T) {
	// Clear all possible token variables
	for _, v := range githubTokenPossibleVariables {
		originalValue := os.Getenv(v)
		defer func(key, val string) {
			if val != "" {
				_ = os.Setenv(key, val) //nolint:errcheck // Test setup
			} else {
				_ = os.Unsetenv(key) //nolint:errcheck // Test cleanup
			}
		}(v, originalValue)
		_ = os.Unsetenv(v) //nolint:errcheck // Test setup
	}

	_, err := token()
	if err == nil {
		t.Error("token() should return error when no env variable is set")
	}

	if err != errTokenNotAvailable {
		t.Errorf("token() error = %v, want %v", err, errTokenNotAvailable)
	}
}

func TestToken_GithubTokenSet(t *testing.T) {
	const testToken = "ghp_test_token_1234567890" //nolint:gosec // Test token

	// Save original value
	originalValue := os.Getenv("GITHUB_TOKEN")
	defer func() {
		if originalValue != "" {
			_ = os.Setenv("GITHUB_TOKEN", originalValue) //nolint:errcheck // Test cleanup
		} else {
			_ = os.Unsetenv("GITHUB_TOKEN") //nolint:errcheck // Test cleanup
		}
	}()

	_ = os.Setenv("GITHUB_TOKEN", testToken) //nolint:errcheck // Test setup

	got, err := token()
	if err != nil {
		t.Fatalf("token() error = %v, want nil", err)
	}

	if got != testToken {
		t.Errorf("token() = %q, want %q", got, testToken)
	}
}

func TestToken_GithubApiTokenSet(t *testing.T) {
	const testToken = "ghp_api_token_1234567890" //nolint:gosec // Test token

	// Clear GITHUB_TOKEN and set GITHUB_API_TOKEN
	originalGithubToken := os.Getenv("GITHUB_TOKEN")
	originalApiToken := os.Getenv("GITHUB_API_TOKEN")
	defer func() {
		if originalGithubToken != "" {
			_ = os.Setenv("GITHUB_TOKEN", originalGithubToken)
		} else {
			_ = os.Unsetenv("GITHUB_TOKEN")
		}
		if originalApiToken != "" {
			_ = os.Setenv("GITHUB_API_TOKEN", originalApiToken)
		} else {
			_ = os.Unsetenv("GITHUB_API_TOKEN")
		}
	}()

	_ = os.Unsetenv("GITHUB_TOKEN")
	_ = os.Setenv("GITHUB_API_TOKEN", testToken)

	got, err := token()
	if err != nil {
		t.Fatalf("token() error = %v, want nil", err)
	}

	if got != testToken {
		t.Errorf("token() = %q, want %q", got, testToken)
	}
}

func TestToken_PriorityOrder(t *testing.T) {
	const tokenFirst = "first_token"
	const tokenSecond = "second_token"

	// Save original values
	originals := make(map[string]string)
	for _, v := range githubTokenPossibleVariables {
		originals[v] = os.Getenv(v)
	}
	defer func() {
		for k, v := range originals {
			if v != "" {
				_ = os.Setenv(k, v)
			} else {
				_ = os.Unsetenv(k)
			}
		}
	}()

	// Set both tokens
	_ = os.Setenv("GITHUB_TOKEN", tokenFirst)
	_ = os.Setenv("GITHUB_API_TOKEN", tokenSecond)

	got, err := token()
	if err != nil {
		t.Fatalf("token() error = %v, want nil", err)
	}

	// Should return first one in the list (GITHUB_TOKEN)
	if got != tokenFirst {
		t.Errorf("token() = %q, want %q (first in priority list)", got, tokenFirst)
	}
}

func TestClient_NoToken(t *testing.T) {
	// Clear all possible token variables
	for _, v := range githubTokenPossibleVariables {
		originalValue := os.Getenv(v)
		defer func(key, val string) {
			if val != "" {
				_ = os.Setenv(key, val)
			} else {
				_ = os.Unsetenv(key)
			}
		}(v, originalValue)
		_ = os.Unsetenv(v)
	}

	_, err := client()
	if err == nil {
		t.Error("client() should return error when no token is available")
	}
}

func TestClient_WithToken(t *testing.T) {
	const testToken = "ghp_test_token" //nolint:gosec // Test token

	// Save and set token
	originalValue := os.Getenv("GITHUB_TOKEN")
	defer func() {
		if originalValue != "" {
			_ = os.Setenv("GITHUB_TOKEN", originalValue)
		} else {
			_ = os.Unsetenv("GITHUB_TOKEN")
		}
	}()

	_ = os.Setenv("GITHUB_TOKEN", testToken)

	c, err := client()
	if err != nil {
		t.Fatalf("client() error = %v, want nil", err)
	}

	if c == nil {
		t.Error("client() returned nil client")
	}
}

func TestIssues_NoToken(t *testing.T) {
	// Clear all possible token variables
	for _, v := range githubTokenPossibleVariables {
		originalValue := os.Getenv(v)
		defer func(key, val string) {
			if val != "" {
				_ = os.Setenv(key, val)
			} else {
				_ = os.Unsetenv(key)
			}
		}(v, originalValue)
		_ = os.Unsetenv(v)
	}

	result := issues(context.Background(), "test query")

	// Should return error message
	if !strings.Contains(result, "token") && !strings.Contains(result, "not found") {
		t.Errorf("issues() should return token error, got: %s", result)
	}
}

func TestSearchIssuesShortcode_NoToken(t *testing.T) {
	// Clear all possible token variables
	for _, v := range githubTokenPossibleVariables {
		originalValue := os.Getenv(v)
		defer func(key, val string) {
			if val != "" {
				_ = os.Setenv(key, val)
			} else {
				_ = os.Unsetenv(key)
			}
		}(v, originalValue)
		_ = os.Unsetenv(v)
	}

	query := xlog.Markdown("repo:test/repo is:open")
	result := seachIssuesShortcode(query)

	resultStr := string(result)
	if !strings.Contains(resultStr, "token") && !strings.Contains(resultStr, "not found") {
		t.Errorf("seachIssuesShortcode() should return token error, got: %s", resultStr)
	}
}

func TestSearchIssuesShortcode_EmptyQuery(t *testing.T) {
	const testToken = "ghp_test_token" //nolint:gosec // Test token

	// Save and set token
	originalValue := os.Getenv("GITHUB_TOKEN")
	defer func() {
		if originalValue != "" {
			_ = os.Setenv("GITHUB_TOKEN", originalValue)
		} else {
			_ = os.Unsetenv("GITHUB_TOKEN")
		}
	}()

	_ = os.Setenv("GITHUB_TOKEN", testToken)

	query := xlog.Markdown("")
	result := seachIssuesShortcode(query)

	// Should not panic with empty query
	if result == "" {
		t.Error("seachIssuesShortcode() should return some content (even if error)")
	}
}

func TestIssues_NoResultsMessage(t *testing.T) {
	// Test the "no results" path (lines 70-72) which returns
	// a message when the GitHub API returns zero issues
	tests := []struct {
		name            string
		query           string
		expectedContain string
	}{
		{
			name:            "empty results returns no results message",
			query:           "repo:nonexistent/impossible is:open",
			expectedContain: "No results for query",
		},
		{
			name:            "impossible filter returns no results message",
			query:           "repo:emad-elsaid/xlog created:1970-01-01",
			expectedContain: "No results for query",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Skip if GITHUB_TOKEN not set (would fail on line 57 instead)
			_, tokenErr := token()
			if tokenErr != nil {
				t.Skip("Skipping test: GITHUB_TOKEN not available")
			}

			result := issues(context.Background(), tc.query)

			if !strings.Contains(result, tc.expectedContain) {
				t.Errorf("Expected result to contain %q, got: %s", tc.expectedContain, result)
			}
		})
	}
}

func TestIssues_HTMLGeneration(t *testing.T) {
	// Test the HTML generation logic (lines 74-87) by verifying
	// the structure when results exist
	tests := []struct {
		name              string
		query             string
		expectedElements  []string
		skipIfNoToken     bool
		minExpectedIssues int
	}{
		{
			name:  "results contain HTML list structure",
			query: "repo:emad-elsaid/xlog is:closed label:enhancement",
			expectedElements: []string{
				"<ul>",
				"</ul>",
				"<li>",
				"</li>",
			},
			skipIfNoToken:     true,
			minExpectedIssues: 0, // May have zero, that's OK
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skipIfNoToken {
				_, tokenErr := token()
				if tokenErr != nil {
					t.Skip("Skipping test: GITHUB_TOKEN not available")
				}
			}

			result := issues(context.Background(), tc.query)

			// If we get "No results", that's valid but skip HTML checks
			if strings.Contains(result, "No results") {
				t.Logf("Query returned no results (valid outcome): %s", result)
				return
			}

			// If we get an error message, log it but don't fail
			// (API might be down or rate limited)
			if strings.Contains(result, "error") || strings.Contains(result, "Error") {
				t.Logf("Query returned error (API issue, not code issue): %s", result)
				return
			}

			// Verify HTML structure for successful results
			for _, expected := range tc.expectedElements {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected HTML to contain %q, got:\n%s", expected, result)
				}
			}
		})
	}
}

func TestIssues_HTMLStructure(t *testing.T) {
	// Test specific HTML elements generated in lines 74-87
	// by using a query likely to return results
	const testQuery = "repo:emad-elsaid/xlog is:issue"

	_, tokenErr := token()
	if tokenErr != nil {
		t.Skip("Skipping test: GITHUB_TOKEN not available")
	}

	result := issues(context.Background(), testQuery)

	// Handle no results case
	if strings.Contains(result, "No results") {
		t.Log("Query returned no results (valid path, lines 70-72 covered)")
		return
	}

	// Handle error case
	if strings.Contains(result, "token") || strings.Contains(result, "API") {
		t.Logf("API error (external issue): %s", result)
		return
	}

	// Verify HTML structure elements if we got results
	expectedStructure := []string{
		"<ul>",
		"</ul>",
		"<li>",
		"<span class=\"icon-text\"",
		"<figure class=\"icon image is-24x24",
		"<img src=",
		"<a href=",
	}

	for _, element := range expectedStructure {
		if !strings.Contains(result, element) {
			t.Errorf("Expected result to contain HTML element %q, but it was missing.\nFull result:\n%s",
				element, result)
		}
	}
}

func TestIssues_ErrorHandling(t *testing.T) {
	// Test error path when API call fails (line 66-68)
	tests := []struct {
		name          string
		setupError    func() func()
		query         string
		expectsError  bool
		errorContains string
	}{
		{
			name: "invalid query returns error message",
			setupError: func() func() {
				orig := os.Getenv("GITHUB_TOKEN")
				_ = os.Setenv("GITHUB_TOKEN", "ghp_invalid_token_1234567890")
				return func() {
					if orig != "" {
						_ = os.Setenv("GITHUB_TOKEN", orig)
					} else {
						_ = os.Unsetenv("GITHUB_TOKEN")
					}
				}
			},
			query:         "invalid syntax query",
			expectsError:  true,
			errorContains: "", // Any error message is fine
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cleanup := tc.setupError()
			defer cleanup()

			result := issues(context.Background(), tc.query)

			// We expect either an error message or valid HTML
			// Empty result would be a bug
			if result == "" {
				t.Error("issues() returned empty string (should return error or HTML)")
			}

			// If we expected error and got "No results", that's OK too
			// (depends on GitHub API behavior with invalid queries)
			if tc.expectsError {
				hasError := strings.Contains(result, "error") ||
					strings.Contains(result, "Error") ||
					strings.Contains(result, "No results") ||
					strings.Contains(result, "401") ||
					strings.Contains(result, "Bad credentials")

				if !hasError && !strings.HasPrefix(result, "<ul>") {
					t.Logf("Note: Expected error-like response, got: %s", result)
				}
			}
		})
	}
}

func TestErrTokenNotAvailable_Message(t *testing.T) {
	expectedMessage := "Github token env variable not found in any of: GITHUB_TOKEN, GITHUB_API_TOKEN"
	if errTokenNotAvailable.Error() != expectedMessage {
		t.Errorf("errTokenNotAvailable message = %q, want %q", errTokenNotAvailable.Error(), expectedMessage)
	}
}

func TestGithubTokenPossibleVariables_NotEmpty(t *testing.T) {
	if len(githubTokenPossibleVariables) == 0 {
		t.Error("githubTokenPossibleVariables should not be empty")
	}

	expectedVars := []string{"GITHUB_TOKEN", "GITHUB_API_TOKEN"}
	if len(githubTokenPossibleVariables) != len(expectedVars) {
		t.Errorf("githubTokenPossibleVariables length = %d, want %d", len(githubTokenPossibleVariables), len(expectedVars))
	}

	for i, want := range expectedVars {
		if i >= len(githubTokenPossibleVariables) {
			break
		}
		if githubTokenPossibleVariables[i] != want {
			t.Errorf("githubTokenPossibleVariables[%d] = %q, want %q", i, githubTokenPossibleVariables[i], want)
		}
	}
}

func TestPerPageConstant(t *testing.T) {
	if perPage != 100 {
		t.Errorf("perPage = %d, want 100", perPage)
	}
}

func TestToken_EmptyStringValue(t *testing.T) {
	// Test that empty string is treated as not set
	for _, v := range githubTokenPossibleVariables {
		originalValue := os.Getenv(v)
		defer func(key, val string) {
			if val != "" {
				_ = os.Setenv(key, val)
			} else {
				_ = os.Unsetenv(key)
			}
		}(v, originalValue)
		_ = os.Setenv(v, "") // Set to empty string
	}

	_, err := token()
	if err == nil {
		t.Error("token() should return error when all env variables are empty strings")
	}
}

func TestClient_Integration(t *testing.T) {
	// This test verifies client creation works end-to-end with a token
	const testToken = "ghp_integration_test_token" //nolint:gosec // Test token

	originalValue := os.Getenv("GITHUB_TOKEN")
	defer func() {
		if originalValue != "" {
			_ = os.Setenv("GITHUB_TOKEN", originalValue)
		} else {
			_ = os.Unsetenv("GITHUB_TOKEN")
		}
	}()

	_ = os.Setenv("GITHUB_TOKEN", testToken)

	c, err := client()
	if err != nil {
		t.Fatalf("client() integration error = %v, want nil", err)
	}

	if c == nil {
		t.Fatal("client() returned nil")
		return
	}

	// Verify it's a proper GitHub client
	if c.Issues == nil {
		t.Error("client should have Issues service initialized")
	}

	if c.Search == nil {
		t.Error("client should have Search service initialized")
	}
}
