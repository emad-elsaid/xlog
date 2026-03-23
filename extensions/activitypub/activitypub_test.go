package activitypub

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/emad-elsaid/xlog"
)

func setupTestEnv(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	wd, _ := os.Getwd()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	// Create test pages
	os.WriteFile("test-page.md", []byte("# Test Page\n\nContent"), 0644)
	os.WriteFile("another-page.md", []byte("# Another Page\n\nMore content"), 0644)

	// Set ActivityPub flags
	domain = "example.com"
	username = "testuser"
	summary = "Test user summary"
	icon = "/public/icon.png"
	image = "/public/image.png"

	cleanup := func() {
		os.Chdir(wd)
		domain = ""
		username = ""
		summary = ""
		icon = "/public/logo.png"
		image = "/public/logo.png"
	}

	return tmpDir, cleanup
}

func TestActivityPubName(t *testing.T) {
	ap := ActivityPub{}
	if got := ap.Name(); got != "activitypub" {
		t.Errorf("Name() = %q, want %q", got, "activitypub")
	}
}

func TestWebfingerWithoutConfig(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	// Reset config to test without domain/username
	domain = ""
	username = ""

	req := httptest.NewRequest(http.MethodGet, "/.well-known/webfinger?resource=acct:test@example.com", nil)
	w := httptest.NewRecorder()

	result := webfinger(req)
	result(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("webfinger without config: got status %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestWebfingerResponse(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/.well-known/webfinger?resource=acct:testuser@example.com", nil)
	w := httptest.NewRecorder()

	result := webfinger(req)
	result(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("webfinger: got status %d, want %d", w.Code, http.StatusOK)
	}

	var response webfingerResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	expectedSubject := "acct:testuser@example.com"
	if response.Subject != expectedSubject {
		t.Errorf("Subject = %q, want %q", response.Subject, expectedSubject)
	}

	if len(response.Aliases) != 2 {
		t.Errorf("got %d aliases, want 2", len(response.Aliases))
	}

	if len(response.Links) != 3 {
		t.Errorf("got %d links, want 3", len(response.Links))
	}

	// Verify self link exists
	foundSelfLink := false
	for _, link := range response.Links {
		if link["rel"] == "self" && link["type"] == "application/activity+json" {
			foundSelfLink = true
			expectedHref := "https://example.com/+/activitypub/@testuser"
			if link["href"] != expectedHref {
				t.Errorf("self link href = %q, want %q", link["href"], expectedHref)
			}
		}
	}
	if !foundSelfLink {
		t.Error("webfinger response missing self link")
	}
}

func TestProfileWithoutConfig(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	domain = ""
	username = ""

	req := httptest.NewRequest(http.MethodGet, "/+/activitypub/@testuser", nil)
	req.SetPathValue("user", "@testuser")
	w := httptest.NewRecorder()

	result := profile(req)
	result(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("profile without config: got status %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestProfileWrongUser(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/+/activitypub/@wronguser", nil)
	req.SetPathValue("user", "@wronguser")
	w := httptest.NewRecorder()

	result := profile(req)
	result(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("profile wrong user: got status %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestProfileResponse(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/+/activitypub/@testuser", nil)
	req.SetPathValue("user", "@testuser")
	w := httptest.NewRecorder()

	result := profile(req)
	result(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("profile: got status %d, want %d", w.Code, http.StatusOK)
	}

	var response profileResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if response.Type != "Person" {
		t.Errorf("Type = %q, want %q", response.Type, "Person")
	}

	if response.PreferredUsername != "testuser" {
		t.Errorf("PreferredUsername = %q, want %q", response.PreferredUsername, "testuser")
	}

	if response.Summary != "Test user summary" {
		t.Errorf("Summary = %q, want %q", response.Summary, "Test user summary")
	}

	expectedOutbox := "https://example.com/+/activitypub/@testuser/outbox"
	if response.Outbox != expectedOutbox {
		t.Errorf("Outbox = %q, want %q", response.Outbox, expectedOutbox)
	}
}

func TestOutboxWithoutConfig(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	domain = ""
	username = ""

	req := httptest.NewRequest(http.MethodGet, "/+/activitypub/@testuser/outbox", nil)
	req.SetPathValue("user", "@testuser")
	w := httptest.NewRecorder()

	result := outbox(req)
	result(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("outbox without config: got status %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestOutboxWrongUser(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/+/activitypub/@wronguser/outbox", nil)
	req.SetPathValue("user", "@wronguser")
	w := httptest.NewRecorder()

	result := outbox(req)
	result(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("outbox wrong user: got status %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestOutboxResponse(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/+/activitypub/@testuser/outbox", nil)
	req.SetPathValue("user", "@testuser")
	w := httptest.NewRecorder()

	result := outbox(req)
	result(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("outbox: got status %d, want %d (dir: %s)", w.Code, http.StatusOK, tmpDir)
	}

	var response outboxResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if response.Type != "OrderedCollection" {
		t.Errorf("Type = %q, want %q", response.Type, "OrderedCollection")
	}

	if response.TotalItems < 2 {
		t.Errorf("TotalItems = %d, want at least 2", response.TotalItems)
	}
}

func TestOutboxPageWithoutConfig(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	domain = ""
	username = ""

	req := httptest.NewRequest(http.MethodGet, "/+/activitypub/@testuser/outbox/1", nil)
	req.SetPathValue("user", "@testuser")
	req.SetPathValue("page", "1")
	w := httptest.NewRecorder()

	result := outboxPage(req)
	result(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("outboxPage without config: got status %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestOutboxPageWrongUser(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/+/activitypub/@wronguser/outbox/1", nil)
	req.SetPathValue("user", "@wronguser")
	req.SetPathValue("page", "1")
	w := httptest.NewRecorder()

	result := outboxPage(req)
	result(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("outboxPage wrong user: got status %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestOutboxPageInvalidIndex(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/+/activitypub/@testuser/outbox/999", nil)
	req.SetPathValue("user", "@testuser")
	req.SetPathValue("page", "999")
	w := httptest.NewRecorder()

	result := outboxPage(req)
	result(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("outboxPage invalid index: got status %d, want %d (dir: %s)", w.Code, http.StatusNotFound, tmpDir)
	}
}

func TestOutboxPageResponse(t *testing.T) {
	tmpDir, cleanup := setupTestEnv(t)
	defer cleanup()

	// Set a known modtime for predictable sorting
	modTime := time.Now().Add(-1 * time.Hour)
	os.Chtimes(filepath.Join(tmpDir, "test-page.md"), modTime, modTime)

	req := httptest.NewRequest(http.MethodGet, "/+/activitypub/@testuser/outbox/1", nil)
	req.SetPathValue("user", "@testuser")
	req.SetPathValue("page", "1")
	w := httptest.NewRecorder()

	result := outboxPage(req)
	result(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("outboxPage: got status %d, want %d (dir: %s)", w.Code, http.StatusOK, tmpDir)
	}

	var response outboxPageResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode JSON response: %v", err)
	}

	if response.Type != "OrderedCollectionPage" {
		t.Errorf("Type = %q, want %q", response.Type, "OrderedCollectionPage")
	}

	if len(response.OrderedItems) != 1 {
		t.Errorf("got %d items, want 1", len(response.OrderedItems))
	}

	if len(response.OrderedItems) > 0 {
		item := response.OrderedItems[0]
		if item.Type != "Create" {
			t.Errorf("item Type = %q, want %q", item.Type, "Create")
		}

		expectedActor := "https://example.com/+/activitypub/@testuser"
		if item.Actor != expectedActor {
			t.Errorf("Actor = %q, want %q", item.Actor, expectedActor)
		}

		if item.Object.Type != "Note" {
			t.Errorf("Object Type = %q, want %q", item.Object.Type, "Note")
		}
	}
}

func TestMetaWithoutConfig(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	domain = ""
	username = ""

	// meta() expects a Page interface, but we can pass nil since it only checks config
	var p Page = nil
	result := meta(p)

	if result != "" {
		t.Errorf("meta without config: got %q, want empty string", result)
	}
}

func TestMetaOutput(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	var p Page = nil
	result := meta(p)

	expected := `<link href='https://example.com/+/activitypub/@testuser' rel='alternate' type='application/activity+json'>`
	if string(result) != expected {
		t.Errorf("meta output = %q, want %q", result, expected)
	}
}
