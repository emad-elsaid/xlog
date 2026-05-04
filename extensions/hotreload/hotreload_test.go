package hotreload

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/gorilla/websocket"
)

func TestHotreload_Name(t *testing.T) {
	h := Hotreload{}
	if got := h.Name(); got != "hotreload" {
		t.Errorf("Name() = %q, want %q", got, "hotreload")
	}
}

func TestNotifyPageChange(t *testing.T) {
	tests := []struct {
		name     string
		pageName string
		exists   bool
		wantErr  bool
	}{
		{
			name:     "notify existing page",
			pageName: "test-page",
			exists:   true,
			wantErr:  false,
		},
		{
			name:     "skip non-existing page",
			pageName: "missing",
			exists:   false,
			wantErr:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Clear clients map before test
			clientsMutex.Lock()
			oldClients := clients
			clients = make(map[*websocket.Conn]bool)
			clientsMutex.Unlock()
			defer func() {
				clientsMutex.Lock()
				clients = oldClients
				clientsMutex.Unlock()
			}()

			// Create mock page
			page := mockPage{name: tc.pageName, exists: tc.exists}

			err := NotifyPageChange(page)
			if (err != nil) != tc.wantErr {
				t.Errorf("NotifyPageChange() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestNotifyPageChange_WithMockConnection(t *testing.T) {
	// Clear and setup
	clientsMutex.Lock()
	oldClients := clients
	clients = make(map[*websocket.Conn]bool)
	clientsMutex.Unlock()
	defer func() {
		clientsMutex.Lock()
		clients = oldClients
		clientsMutex.Unlock()
	}()

	// Create mock connection using httptest
	msgReceived := make(chan map[string]string, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer func() {
			if err := conn.Close(); err != nil {
				t.Logf("Failed to close connection: %v", err)
			}
		}()

		// Read one message
		var msg map[string]string
		if err := conn.ReadJSON(&msg); err == nil {
			msgReceived <- msg
		}
	}))
	defer server.Close()

	// Connect as websocket client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Logf("Failed to close client connection: %v", err)
		}
	}()

	// Add to clients manually
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	// Notify
	page := mockPage{name: "test-page", exists: true}
	if err := NotifyPageChange(page); err != nil {
		t.Errorf("NotifyPageChange() unexpected error: %v", err)
	}

	// Verify message received
	select {
	case msg := <-msgReceived:
		expectedURL := "/test-page"
		if msg["url"] != expectedURL {
			t.Errorf("Got URL %q, want %q", msg["url"], expectedURL)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timeout waiting for WebSocket message")
	}
}

func TestClientWidget(t *testing.T) {
	result := clientWidget(mockPage{name: "test", exists: true})

	// Should return embedded script
	if result == "" {
		t.Error("clientWidget() returned empty string")
	}

	// Should be HTML template type
	resultStr := string(result)
	if !strings.Contains(resultStr, "script") && !strings.Contains(resultStr, "websocket") {
		t.Error("clientWidget() does not appear to contain script/websocket content")
	}
}

func TestHandleWebSocket(t *testing.T) {
	// Clear clients
	clientsMutex.Lock()
	clients = make(map[*websocket.Conn]bool)
	clientsMutex.Unlock()

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := handleWebSocket(nil)
		handler(w, r)
	}))
	defer server.Close()

	// Connect as websocket client
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Logf("Failed to close connection: %v", err)
		}
	}()

	// Give time for connection to be registered
	time.Sleep(50 * time.Millisecond)

	// Verify client was added
	clientsMutex.Lock()
	count := len(clients)
	clientsMutex.Unlock()

	if count != 1 {
		t.Errorf("Expected 1 client, got %d", count)
	}

	// Close connection
	if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		t.Logf("Failed to send close message: %v", err)
	}
	time.Sleep(50 * time.Millisecond)

	// Verify client was removed
	clientsMutex.Lock()
	finalCount := len(clients)
	clientsMutex.Unlock()

	if finalCount != 0 {
		t.Errorf("Expected client to be removed after close, got %d", finalCount)
	}
}

func TestNotifyPageChange_WriteErrorHandling(t *testing.T) {
	tests := []struct {
		name             string
		closeConnBefore  bool
		expectClientGone bool
	}{
		{
			name:             "removes client when write fails on closed connection",
			closeConnBefore:  true,
			expectClientGone: true,
		},
		{
			name:             "keeps client when write succeeds",
			closeConnBefore:  false,
			expectClientGone: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			clientsMutex.Lock()
			oldClients := clients
			clients = make(map[*websocket.Conn]bool)
			clientsMutex.Unlock()
			defer func() {
				clientsMutex.Lock()
				clients = oldClients
				clientsMutex.Unlock()
			}()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				conn, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					return
				}
				defer func() {
					if err := conn.Close(); err != nil {
						t.Logf("Server close connection: %v", err)
					}
				}()

				var msg map[string]string
				_ = conn.ReadJSON(&msg)
			}))
			defer server.Close()

			wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			if err != nil {
				t.Fatal(err)
			}

			clientsMutex.Lock()
			clients[conn] = true
			clientsMutex.Unlock()

			if tc.closeConnBefore {
				if err := conn.Close(); err != nil {
					t.Logf("Failed to close connection: %v", err)
				}
				time.Sleep(50 * time.Millisecond)
			} else {
				defer func() {
					if err := conn.Close(); err != nil {
						t.Logf("Failed to close connection: %v", err)
					}
				}()
			}

			page := mockPage{name: "test-page", exists: true}
			err = NotifyPageChange(page)
			if err != nil {
				t.Errorf("NotifyPageChange() unexpected error: %v", err)
			}

			time.Sleep(50 * time.Millisecond)

			clientsMutex.Lock()
			_, stillPresent := clients[conn]
			clientsMutex.Unlock()

			if tc.expectClientGone && stillPresent {
				t.Error("Expected client to be removed after write error, but still present")
			}
			if !tc.expectClientGone && !stillPresent {
				t.Error("Expected client to remain after successful write, but was removed")
			}
		})
	}
}

// mockPage implements xlog.Page interface for testing.
type mockPage struct {
	name   string
	exists bool
}

func (m mockPage) Name() string             { return m.name }
func (m mockPage) Exists() bool             { return m.exists }
func (m mockPage) FileName() string         { return m.name + ".md" }
func (m mockPage) Render() template.HTML    { return "" }
func (m mockPage) Content() xlog.Markdown   { return xlog.Markdown("") }
func (m mockPage) Delete() bool             { return true }
func (m mockPage) Write(xlog.Markdown) bool { return true }
func (m mockPage) ModTime() time.Time       { return time.Time{} }
func (m mockPage) AST() ([]byte, ast.Node)  { return nil, nil }
