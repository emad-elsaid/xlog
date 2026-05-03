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
	}

	// Verify it's a proper GitHub client
	if c.Issues == nil {
		t.Error("client should have Issues service initialized")
	}

	if c.Search == nil {
		t.Error("client should have Search service initialized")
	}
}
