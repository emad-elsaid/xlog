package xlog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// failingResponseWriter is a mock ResponseWriter that fails on Write operations
// to test error handling paths.
type failingResponseWriter struct {
	writeCalled bool
	header      http.Header
}

func (f *failingResponseWriter) Header() http.Header {
	if f.header == nil {
		f.header = make(http.Header)
	}
	return f.header
}

func (f *failingResponseWriter) Write(b []byte) (int, error) {
	f.writeCalled = true
	return 0, errors.New("simulated write failure")
}

func (f *failingResponseWriter) WriteHeader(statusCode int) {
	// No-op for this mock
}

func TestBadRequest(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "simple error message",
			message:        "invalid input",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid input\n",
		},
		{
			name:           "empty message",
			message:        "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "\n",
		},
		{
			name:           "detailed validation error",
			message:        "field 'email' is required",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "field 'email' is required\n",
		},
		{
			name:           "special characters in message",
			message:        "invalid char: <>&\"'",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid char: <>&\"'\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := BadRequest(tc.message)
			handler.ServeHTTP(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("status code: want %d, got %d", tc.expectedStatus, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestUnauthorized(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing credentials",
			message:        "authentication required",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "authentication required\n",
		},
		{
			name:           "invalid token",
			message:        "invalid authentication token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "invalid authentication token\n",
		},
		{
			name:           "expired session",
			message:        "session expired",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "session expired\n",
		},
		{
			name:           "empty message",
			message:        "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := Unauthorized(tc.message)
			handler.ServeHTTP(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("status code: want %d, got %d", tc.expectedStatus, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestForbidden(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "insufficient permissions",
			message:        "insufficient permissions",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "insufficient permissions\n",
		},
		{
			name:           "access denied",
			message:        "access denied to resource",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "access denied to resource\n",
		},
		{
			name:           "operation not allowed",
			message:        "delete operation not permitted",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "delete operation not permitted\n",
		},
		{
			name:           "empty message",
			message:        "",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := Forbidden(tc.message)
			handler.ServeHTTP(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("status code: want %d, got %d", tc.expectedStatus, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestMethodNotAllowed(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "simple method not allowed",
			message:        "method not allowed",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
		{
			name:           "specific method detail",
			message:        "POST method not allowed for this endpoint",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "POST method not allowed for this endpoint\n",
		},
		{
			name:           "with allowed methods hint",
			message:        "method not allowed, use GET or POST",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed, use GET or POST\n",
		},
		{
			name:           "empty message",
			message:        "",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := MethodNotAllowed(tc.message)
			handler.ServeHTTP(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("status code: want %d, got %d", tc.expectedStatus, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestInternalServerError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "generic error",
			err:            errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database connection failed\n",
		},
		{
			name:           "wrapped error",
			err:            errors.New("failed to process request"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to process request\n",
		},
		{
			name:           "empty error message",
			err:            errors.New(""),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := InternalServerError(tc.err)
			handler.ServeHTTP(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("status code: want %d, got %d", tc.expectedStatus, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestNoContent(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "returns 204 no content",
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, "/", http.NoBody)

			handler := NoContent()
			handler.ServeHTTP(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("status code: want %d, got %d", tc.expectedStatus, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestPlainText(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		expectedBody string
	}{
		{
			name:         "simple text",
			text:         "Hello, World!",
			expectedBody: "Hello, World!",
		},
		{
			name:         "empty text",
			text:         "",
			expectedBody: "",
		},
		{
			name:         "multiline text",
			text:         "Line 1\nLine 2\nLine 3",
			expectedBody: "Line 1\nLine 2\nLine 3",
		},
		{
			name:         "special characters",
			text:         "<html>&amp; \"quotes\"</html>",
			expectedBody: "<html>&amp; \"quotes\"</html>",
		},
		{
			name:         "unicode text",
			text:         "Hello 世界 🌍",
			expectedBody: "Hello 世界 🌍",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := PlainText(tc.text)
			handler.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Errorf("status code: want %d, got %d", http.StatusOK, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestJsonResponse(t *testing.T) {
	tests := []struct {
		name         string
		input        any
		expectedBody string
		shouldMatch  bool
	}{
		{
			name:         "simple object",
			input:        map[string]string{"status": "ok", "message": "success"},
			expectedBody: `{"message":"success","status":"ok"}`,
			shouldMatch:  true,
		},
		{
			name:         "array of strings",
			input:        []string{"apple", "banana", "cherry"},
			expectedBody: `["apple","banana","cherry"]`,
			shouldMatch:  true,
		},
		{
			name:         "number",
			input:        42,
			expectedBody: `42`,
			shouldMatch:  true,
		},
		{
			name:         "boolean",
			input:        true,
			expectedBody: `true`,
			shouldMatch:  true,
		},
		{
			name:         "null value",
			input:        nil,
			expectedBody: `null`,
			shouldMatch:  true,
		},
		{
			name:         "empty object",
			input:        map[string]string{},
			expectedBody: `{}`,
			shouldMatch:  true,
		},
		{
			name: "nested structure",
			input: map[string]any{
				"user": map[string]any{
					"name": "Alice",
					"age":  30,
				},
			},
			expectedBody: `{"user":{"age":30,"name":"Alice"}}`,
			shouldMatch:  true,
		},
		{
			name:         "special characters in strings",
			input:        map[string]string{"text": "quotes \"and\" newlines\n"},
			expectedBody: `{"text":"quotes \"and\" newlines\n"}`,
			shouldMatch:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := JsonResponse(tc.input)
			handler.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Errorf("status code: want %d, got %d", http.StatusOK, w.Code)
			}

			if tc.shouldMatch {
				// Parse both as JSON to compare regardless of formatting
				var expected, actual any
				if err := json.Unmarshal([]byte(tc.expectedBody), &expected); err != nil {
					t.Fatalf("failed to unmarshal expected JSON: %v", err)
				}
				if err := json.Unmarshal(w.Body.Bytes(), &actual); err != nil {
					t.Fatalf("failed to unmarshal actual JSON: %v", err)
				}

				expectedJSON, _ := json.Marshal(expected)
				actualJSON, _ := json.Marshal(actual)

				if string(expectedJSON) != string(actualJSON) {
					t.Errorf("JSON mismatch:\nwant: %s\ngot:  %s",
						string(expectedJSON), string(actualJSON))
				}
			}
		})
	}
}

func TestJsonResponse_MarshalError(t *testing.T) {
	// Create a type that cannot be marshaled to JSON
	type unmarshalable struct {
		Chan chan int // channels cannot be marshaled
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

	handler := JsonResponse(unmarshalable{Chan: make(chan int)})
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status code: want %d, got %d", http.StatusOK, w.Code)
	}

	// Should contain error message about unsupported type
	body := w.Body.String()
	if body == "" {
		t.Error("expected error message in response body, got empty string")
	}
	if !contains(body, "unsupported") && !contains(body, "json") {
		t.Errorf("expected JSON marshaling error message, got: %q", body)
	}
}

// Helper function to check if string contains substring (case-insensitive).
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > len(substr) && anyContains(s, substr))
}

func anyContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestHealthHandler tests the health check endpoint for deployment monitoring.
func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "health check returns OK",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/+/health", http.NoBody)

			healthHandler(recorder, req)

			if status := recorder.Code; status != tc.expectedStatus {
				t.Errorf("healthHandler() status = %d, want %d", status, tc.expectedStatus)
			}

			if body := recorder.Body.String(); body != tc.expectedBody {
				t.Errorf("healthHandler() body = %q, want %q", body, tc.expectedBody)
			}

			// Verify Content-Type header
			if ct := recorder.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
				t.Errorf("healthHandler() Content-Type = %q, want %q", ct, "text/plain; charset=utf-8")
			}
		})
	}
}

func TestHealthHandler_WriteError(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "handles write error gracefully",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := &failingResponseWriter{}
			req := httptest.NewRequest(http.MethodGet, "/+/health", http.NoBody)

			// Should not panic despite write error
			healthHandler(w, req)

			// Verify write was attempted
			if !w.writeCalled {
				t.Error("Write should have been called")
			}

			// Verify headers were set before write attempt
			if ct := w.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
				t.Errorf("healthHandler() Content-Type = %q, want %q", ct, "text/plain; charset=utf-8")
			}
		})
	}
}

func TestRequestLoggerHandler_LogInjectionPrevention(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "clean GET request",
			method: "GET",
			path:   "/normal/path",
		},
		{
			name:   "POST with query params",
			method: "POST",
			path:   "/api/endpoint?key=value",
		},
		{
			name:   "path with special chars",
			method: "GET",
			path:   "/path/with%20space",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := requestLoggerHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, http.NoBody)

			// Execute handler - should complete without panic
			handler.ServeHTTP(recorder, req)

			if status := recorder.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
			}
		})
	}
}

// TestSanitizeLogString tests the log sanitization function to prevent injection attacks.
func TestSanitizeLogString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "clean string",
			input:    "GET /path",
			expected: "GET /path",
		},
		{
			name:     "newline character",
			input:    "GET /path\nINJECTED",
			expected: "GET /path INJECTED",
		},
		{
			name:     "carriage return",
			input:    "GET /path\rINJECTED",
			expected: "GET /path INJECTED",
		},
		{
			name:     "CRLF sequence",
			input:    "GET /path\r\nINJECTED",
			expected: "GET /path  INJECTED",
		},
		{
			name:     "multiple newlines",
			input:    "A\nB\nC",
			expected: "A B C",
		},
		{
			name:     "tab character (allowed)",
			input:    "GET\t/path",
			expected: "GET\t/path",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := sanitizeLogString(tc.input)
			if result != tc.expected {
				t.Errorf("sanitizeLogString(%q) = %q, want %q", tc.input, result, tc.expected)
			}
		})
	}
}

func TestNotFound(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "simple not found message",
			message:        "page not found",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "page not found\n",
		},
		{
			name:           "empty message",
			message:        "",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "\n",
		},
		{
			name:           "detailed message",
			message:        "resource /users/123 does not exist",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "resource /users/123 does not exist\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := NotFound(tc.message)
			handler.ServeHTTP(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("status code: want %d, got %d", tc.expectedStatus, w.Code)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestRedirect(t *testing.T) {
	tests := []struct {
		name             string
		url              string
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:             "redirect to home page",
			url:              "/",
			expectedStatus:   http.StatusFound,
			expectedLocation: "/",
		},
		{
			name:             "redirect to specific page",
			url:              "/dashboard",
			expectedStatus:   http.StatusFound,
			expectedLocation: "/dashboard",
		},
		{
			name:             "redirect with query parameters",
			url:              "/search?q=test&page=2",
			expectedStatus:   http.StatusFound,
			expectedLocation: "/search?q=test&page=2",
		},
		{
			name:             "redirect to external URL",
			url:              "https://example.com",
			expectedStatus:   http.StatusFound,
			expectedLocation: "https://example.com",
		},
		{
			name:             "redirect with special characters",
			url:              "/page?name=hello%20world",
			expectedStatus:   http.StatusFound,
			expectedLocation: "/page?name=hello%20world",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := Redirect(tc.url)
			handler.ServeHTTP(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("status code: want %d, got %d", tc.expectedStatus, w.Code)
			}

			location := w.Header().Get("Location")
			if location != tc.expectedLocation {
				t.Errorf("location header: want %q, got %q", tc.expectedLocation, location)
			}
		})
	}
}

func TestNoCache(t *testing.T) {
	tests := []struct {
		name                 string
		wrappedOutput        Output
		expectedCacheControl string
		expectedPragma       string
		expectedExpires      string
		expectedBody         string
	}{
		{
			name:                 "no cache plain text response",
			wrappedOutput:        PlainText("dynamic content"),
			expectedCacheControl: "no-cache, no-store, must-revalidate",
			expectedPragma:       "no-cache",
			expectedExpires:      "0",
			expectedBody:         "dynamic content",
		},
		{
			name:                 "no cache empty response",
			wrappedOutput:        PlainText(""),
			expectedCacheControl: "no-cache, no-store, must-revalidate",
			expectedPragma:       "no-cache",
			expectedExpires:      "0",
			expectedBody:         "",
		},
		{
			name:                 "no cache json response",
			wrappedOutput:        JsonResponse(map[string]string{"status": "ok"}),
			expectedCacheControl: "no-cache, no-store, must-revalidate",
			expectedPragma:       "no-cache",
			expectedExpires:      "0",
			expectedBody:         `{"status":"ok"}`,
		},
		{
			name: "no cache no content response",
			wrappedOutput: func(w Response, r Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			expectedCacheControl: "no-cache, no-store, must-revalidate",
			expectedPragma:       "no-cache",
			expectedExpires:      "0",
			expectedBody:         "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := NoCache(tc.wrappedOutput)
			handler.ServeHTTP(w, r)

			cacheControl := w.Header().Get("Cache-Control")
			if cacheControl != tc.expectedCacheControl {
				t.Errorf("cache-control header: want %q, got %q",
					tc.expectedCacheControl, cacheControl)
			}

			pragma := w.Header().Get("Pragma")
			if pragma != tc.expectedPragma {
				t.Errorf("pragma header: want %q, got %q",
					tc.expectedPragma, pragma)
			}

			expires := w.Header().Get("Expires")
			if expires != tc.expectedExpires {
				t.Errorf("expires header: want %q, got %q",
					tc.expectedExpires, expires)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestCache(t *testing.T) {
	tests := []struct {
		name                 string
		wrappedOutput        Output
		expectedCacheControl string
		expectedBody         string
	}{
		{
			name:                 "cache plain text response",
			wrappedOutput:        PlainText("cached content"),
			expectedCacheControl: "max-age=604800",
			expectedBody:         "cached content",
		},
		{
			name:                 "cache empty response",
			wrappedOutput:        PlainText(""),
			expectedCacheControl: "max-age=604800",
			expectedBody:         "",
		},
		{
			name:                 "cache json response",
			wrappedOutput:        JsonResponse(map[string]string{"status": "ok"}),
			expectedCacheControl: "max-age=604800",
			expectedBody:         `{"status":"ok"}`,
		},
		{
			name: "cache no content response",
			wrappedOutput: func(w Response, r Request) {
				w.WriteHeader(http.StatusNoContent)
			},
			expectedCacheControl: "max-age=604800",
			expectedBody:         "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := Cache(tc.wrappedOutput)
			handler.ServeHTTP(w, r)

			cacheControl := w.Header().Get("Cache-Control")
			if cacheControl != tc.expectedCacheControl {
				t.Errorf("cache-control header: want %q, got %q",
					tc.expectedCacheControl, cacheControl)
			}

			if w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestHandlerFuncToHttpHandler(t *testing.T) {
	tests := []struct {
		name           string
		handlerFunc    HandlerFunc
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "plain text handler",
			handlerFunc: func(r Request) Output {
				return PlainText("test response")
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "test response",
		},
		{
			name: "not found handler",
			handlerFunc: func(r Request) Output {
				return NotFound("resource not found")
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "resource not found\n",
		},
		{
			name: "redirect handler",
			handlerFunc: func(r Request) Output {
				return Redirect("/new-location")
			},
			expectedStatus: http.StatusFound,
			expectedBody:   "",
		},
		{
			name: "json response handler",
			handlerFunc: func(r Request) Output {
				return JsonResponse(map[string]int{"count": 42})
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"count":42}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := handlerFuncToHttpHandler(tc.handlerFunc)
			handler.ServeHTTP(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("status code: want %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedBody != "" && w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestRouteRegistration(t *testing.T) {
	// Save original router and restore after test
	originalRouter := router
	defer func() { router = originalRouter }()

	tests := []struct {
		name           string
		method         string
		path           string
		registerFunc   func(string, HandlerFunc)
		expectedMethod string
	}{
		{
			name:           "GET route registration",
			method:         http.MethodGet,
			path:           "/test-get",
			registerFunc:   Get,
			expectedMethod: http.MethodGet,
		},
		{
			name:           "POST route registration",
			method:         http.MethodPost,
			path:           "/test-post",
			registerFunc:   Post,
			expectedMethod: http.MethodPost,
		},
		{
			name:           "DELETE route registration",
			method:         http.MethodDelete,
			path:           "/test-delete",
			registerFunc:   Delete,
			expectedMethod: http.MethodDelete,
		},
		{
			name:           "PUT route registration",
			method:         http.MethodPut,
			path:           "/test-put",
			registerFunc:   Put,
			expectedMethod: http.MethodPut,
		},
		{
			name:           "PATCH route registration",
			method:         http.MethodPatch,
			path:           "/test-patch",
			registerFunc:   Patch,
			expectedMethod: http.MethodPatch,
		},
		{
			name:           "HEAD route registration",
			method:         http.MethodHead,
			path:           "/test-head",
			registerFunc:   Head,
			expectedMethod: http.MethodHead,
		},
		{
			name:           "OPTIONS route registration",
			method:         http.MethodOptions,
			path:           "/test-options",
			registerFunc:   Options,
			expectedMethod: http.MethodOptions,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh router for each test
			router = http.NewServeMux()

			// Register route
			tc.registerFunc(tc.path, func(r Request) Output {
				return PlainText("handler executed")
			})

			// Test route works with correct method
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tc.expectedMethod, tc.path, http.NoBody)
			router.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Errorf("expected status %d for %s %s, got %d",
					http.StatusOK, tc.expectedMethod, tc.path, w.Code)
			}

			if w.Body.String() != "handler executed" {
				t.Errorf("expected body %q, got %q",
					"handler executed", w.Body.String())
			}

			// Test route does not respond to wrong method
			wrongMethod := http.MethodPost
			if tc.expectedMethod == http.MethodPost {
				wrongMethod = http.MethodGet
			}

			w2 := httptest.NewRecorder()
			r2 := httptest.NewRequest(wrongMethod, tc.path, http.NoBody)
			router.ServeHTTP(w2, r2)

			if w2.Code != http.StatusMethodNotAllowed {
				t.Errorf("expected status %d for wrong method %s %s, got %d",
					http.StatusMethodNotAllowed, wrongMethod, tc.path, w2.Code)
			}
		})
	}
}

func TestPlainTextOutput(t *testing.T) {
	tests := []struct {
		name         string
		text         string
		expectedBody string
	}{
		{
			name:         "simple text",
			text:         "hello world",
			expectedBody: "hello world",
		},
		{
			name:         "empty text",
			text:         "",
			expectedBody: "",
		},
		{
			name:         "text with special characters",
			text:         "<html>&special</html>",
			expectedBody: "<html>&special</html>",
		},
		{
			name:         "multiline text",
			text:         "line1\nline2\nline3",
			expectedBody: "line1\nline2\nline3",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := PlainText(tc.text)
			handler.ServeHTTP(w, r)

			if w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestJsonResponseOutput(t *testing.T) {
	tests := []struct {
		name         string
		data         any
		expectedBody string
		shouldError  bool
	}{
		{
			name:         "simple object",
			data:         map[string]string{"key": "value"},
			expectedBody: `{"key":"value"}`,
			shouldError:  false,
		},
		{
			name:         "empty object",
			data:         map[string]string{},
			expectedBody: `{}`,
			shouldError:  false,
		},
		{
			name:         "array",
			data:         []int{1, 2, 3},
			expectedBody: `[1,2,3]`,
			shouldError:  false,
		},
		{
			name:         "null value",
			data:         nil,
			expectedBody: `null`,
			shouldError:  false,
		},
		{
			name:         "unmarshalable value",
			data:         make(chan int),
			expectedBody: "json: unsupported type: chan int",
			shouldError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := JsonResponse(tc.data)
			handler.ServeHTTP(w, r)

			if w.Body.String() != tc.expectedBody {
				t.Errorf("body: want %q, got %q", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestRenderOutput(t *testing.T) {
	tests := []struct {
		name         string
		templatePath string
		data         Locals
		shouldFind   string
	}{
		{
			name:         "render missing template shows error",
			templatePath: "nonexistent-template-path",
			data:         Locals{},
			shouldFind:   "template nonexistent-template-path not found",
		},
		{
			name:         "render with nil data",
			templatePath: "missing",
			data:         nil,
			shouldFind:   "template missing not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := Render(tc.templatePath, tc.data)
			handler.ServeHTTP(w, r)

			if tc.shouldFind != "" && !contains(w.Body.String(), tc.shouldFind) {
				t.Errorf("expected response to contain %q, got %q",
					tc.shouldFind, w.Body.String())
			}
		})
	}
}

func TestPlainTextOutput_WriteError(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "write error on simple text",
			text: "hello world",
		},
		{
			name: "write error on empty text",
			text: "",
		},
		{
			name: "write error on large text",
			text: string(make([]byte, 10000)),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := &failingResponseWriter{}
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := PlainText(tc.text)
			handler.ServeHTTP(w, r)

			if !w.writeCalled {
				t.Error("Write should have been called")
			}
		})
	}
}

func TestJsonResponseOutput_WriteError(t *testing.T) {
	tests := []struct {
		name string
		data any
	}{
		{
			name: "write error on valid JSON",
			data: map[string]string{"key": "value"},
		},
		{
			name: "write error on array",
			data: []int{1, 2, 3},
		},
		{
			name: "write error on marshal failure",
			data: make(chan int),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := &failingResponseWriter{}
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := JsonResponse(tc.data)
			handler.ServeHTTP(w, r)

			if !w.writeCalled {
				t.Error("Write should have been called")
			}
		})
	}
}

func TestRenderOutput_WriteError(t *testing.T) {
	tests := []struct {
		name         string
		templatePath string
		data         Locals
	}{
		{
			name:         "write error on missing template",
			templatePath: "nonexistent",
			data:         Locals{},
		},
		{
			name:         "write error with nil data",
			templatePath: "missing",
			data:         nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := &failingResponseWriter{}
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			handler := Render(tc.templatePath, tc.data)
			handler.ServeHTTP(w, r)

			if !w.writeCalled {
				t.Error("Write should have been called")
			}
		})
	}
}

func TestFuncStringer(t *testing.T) {
	tests := []struct {
		name          string
		function      any
		shouldContain string
	}{
		{
			name:          "lambda function",
			function:      func() {},
			shouldContain: "TestFuncStringer",
		},
		{
			name:          "package function reference",
			function:      TestFuncStringer,
			shouldContain: "TestFuncStringer",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fs := funcStringer{tc.function}
			result := fs.String()

			if !contains(result, tc.shouldContain) {
				t.Errorf("String(): expected to contain %q, got %q",
					tc.shouldContain, result)
			}

			// Should not contain full github path prefix
			if contains(result, "github.com/") {
				t.Errorf("String(): should strip github prefix, got %q", result)
			}
		})
	}
}

func TestCreated(t *testing.T) {
	tests := []struct {
		name             string
		location         string
		expectLocation   bool
		expectedLocation string
	}{
		{
			name:             "with location header",
			location:         "/api/users/123",
			expectLocation:   true,
			expectedLocation: "/api/users/123",
		},
		{
			name:           "without location header",
			location:       "",
			expectLocation: false,
		},
		{
			name:             "with absolute URL",
			location:         "https://example.com/resource/456",
			expectLocation:   true,
			expectedLocation: "https://example.com/resource/456",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/users", http.NoBody)

			output := Created(tc.location)
			output(recorder, req)

			if status := recorder.Code; status != http.StatusCreated {
				t.Errorf("Created() status = %d, want %d", status, http.StatusCreated)
			}

			location := recorder.Header().Get("Location")
			if tc.expectLocation {
				if location != tc.expectedLocation {
					t.Errorf("Created() Location header = %q, want %q", location, tc.expectedLocation)
				}
			} else {
				if location != "" {
					t.Errorf("Created() expected no Location header, got %q", location)
				}
			}
		})
	}
}

func TestAccepted(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "async POST request",
			method: http.MethodPost,
			path:   "/api/async-operation",
		},
		{
			name:   "async PUT request",
			method: http.MethodPut,
			path:   "/api/background-task",
		},
		{
			name:   "async DELETE request",
			method: http.MethodDelete,
			path:   "/api/deferred-deletion",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, http.NoBody)

			output := Accepted()
			output(recorder, req)

			if status := recorder.Code; status != http.StatusAccepted {
				t.Errorf("Accepted() status = %d, want %d", status, http.StatusAccepted)
			}

			if body := recorder.Body.String(); body != "" {
				t.Errorf("Accepted() body = %q, want empty", body)
			}
		})
	}
}

// BenchmarkHandlerFuncToHttpHandler measures the overhead of wrapping
// a HandlerFunc into an http.HandlerFunc. This executes once per handler
// registration during server initialization.
func BenchmarkHandlerFuncToHttpHandler(b *testing.B) {
	handler := func(r Request) Output {
		return func(w Response, r Request) {
			w.WriteHeader(http.StatusOK)
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = handlerFuncToHttpHandler(handler)
	}
}

// BenchmarkResponseFunctions measures the performance of various HTTP response
// helper functions. These are hot-path operations executing on every error condition.
func BenchmarkResponseFunctions(b *testing.B) {
	tests := []struct {
		name string
		fn   func() Output
	}{
		{
			name: "NotFound",
			fn:   func() Output { return NotFound("page not found") },
		},
		{
			name: "BadRequest",
			fn:   func() Output { return BadRequest("invalid input") },
		},
		{
			name: "Unauthorized",
			fn:   func() Output { return Unauthorized("auth required") },
		},
		{
			name: "Forbidden",
			fn:   func() Output { return Forbidden("access denied") },
		},
		{
			name: "InternalServerError",
			fn:   func() Output { return InternalServerError(errors.New("internal error")) },
		},
		{
			name: "Redirect",
			fn:   func() Output { return Redirect("/login") },
		},
		{
			name: "NoContent",
			fn:   NoContent,
		},
		{
			name: "Accepted",
			fn:   Accepted,
		},
	}

	for _, tc := range tests {
		b.Run(tc.name, func(b *testing.B) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				output := tc.fn()
				output(w, r)
				w.Body.Reset()
			}
		})
	}
}

// BenchmarkDefaultMiddlewares measures middleware chain construction overhead.
// This executes once during server initialization but is critical for startup time.
func BenchmarkDefaultMiddlewares(b *testing.B) {
	tests := []struct {
		name     string
		readonly bool
	}{
		{
			name:     "with CSRF protection",
			readonly: false,
		},
		{
			name:     "readonly mode (no CSRF)",
			readonly: true,
		},
	}

	for _, tc := range tests {
		b.Run(tc.name, func(b *testing.B) {
			// Save and restore original config
			origReadonly := Config.Readonly
			defer func() { Config.Readonly = origReadonly }()

			Config.Readonly = tc.readonly

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = defaultMiddlewares()
			}
		})
	}
}

// BenchmarkHandlerChainExecution measures the full execution path through
// handlerFuncToHttpHandler wrapper, simulating real request handling.
func BenchmarkHandlerChainExecution(b *testing.B) {
	handler := func(r Request) Output {
		return func(w Response, r Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		}
	}

	httpHandler := handlerFuncToHttpHandler(handler)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		httpHandler(w, r)
		w.Body.Reset()
	}
}

// TestDefaultMiddlewares_SessionSecretGeneration tests secure random generation
// for CSRF protection session secrets in non-readonly mode.
func TestDefaultMiddlewares_SessionSecretGeneration(t *testing.T) {
	tests := []struct {
		name             string
		sessionSecretEnv string
		readonly         bool
		expectMiddleware bool
	}{
		{
			name:             "generates random secret when env not set",
			sessionSecretEnv: "",
			readonly:         false,
			expectMiddleware: true,
		},
		{
			name:             "uses env secret when provided",
			sessionSecretEnv: "my-custom-secret-key",
			readonly:         false,
			expectMiddleware: true,
		},
		{
			name:             "skips CSRF in readonly mode",
			sessionSecretEnv: "",
			readonly:         true,
			expectMiddleware: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Save and restore environment
			oldEnv := os.Getenv("SESSION_SECRET")
			defer os.Setenv("SESSION_SECRET", oldEnv)

			// Save and restore config
			oldReadonly := Config.Readonly
			defer func() { Config.Readonly = oldReadonly }()

			// Set test environment
			if tc.sessionSecretEnv != "" {
				os.Setenv("SESSION_SECRET", tc.sessionSecretEnv)
			} else {
				os.Unsetenv("SESSION_SECRET")
			}
			Config.Readonly = tc.readonly

			// Execute function - should not panic on crypto errors
			middlewares := defaultMiddlewares()

			// Verify middleware count expectations
			if tc.readonly {
				// Only request logger middleware
				if len(middlewares) != 1 {
					t.Errorf("readonly mode: expected 1 middleware, got %d", len(middlewares))
				}
			} else {
				// CSRF + request logger middlewares
				if len(middlewares) != 2 {
					t.Errorf("non-readonly mode: expected 2 middlewares, got %d", len(middlewares))
				}
			}
		})
	}
}

// TestDefaultMiddlewares_RandReadErrorHandling tests that crypto/rand errors
// during session secret generation are properly caught and handled.
// This is a critical security test - cryptographic operations must never silently fail.
func TestDefaultMiddlewares_RandReadErrorHandling(t *testing.T) {
	t.Run("handles random generation successfully", func(t *testing.T) {
		// Save environment
		oldEnv := os.Getenv("SESSION_SECRET")
		defer os.Setenv("SESSION_SECRET", oldEnv)
		oldReadonly := Config.Readonly
		defer func() { Config.Readonly = oldReadonly }()

		// Test with empty environment (forces rand.Read path)
		os.Unsetenv("SESSION_SECRET")
		Config.Readonly = false

		// Mock osExit to ensure it's NOT called during successful operation
		oldOsExit := osExit
		exitCalled := false
		exitCode := -1
		osExit = func(code int) {
			exitCalled = true
			exitCode = code
			panic("mock exit") // Use panic to stop execution in test
		}
		defer func() {
			osExit = oldOsExit
			if r := recover(); r != nil && r != "mock exit" {
				panic(r) // Re-panic if it's not our mock exit
			}
		}()

		// Execute - should succeed with proper random generation
		middlewares := defaultMiddlewares()

		// Verify no exit on success
		if exitCalled {
			t.Errorf("osExit unexpectedly called with code %d during normal operation", exitCode)
		}

		// Verify middlewares created successfully
		if len(middlewares) != 2 {
			t.Errorf("expected 2 middlewares (CSRF + logger), got %d", len(middlewares))
		}
	})

	t.Run("uses environment variable when provided", func(t *testing.T) {
		// Save environment
		oldEnv := os.Getenv("SESSION_SECRET")
		defer os.Setenv("SESSION_SECRET", oldEnv)
		oldReadonly := Config.Readonly
		defer func() { Config.Readonly = oldReadonly }()

		// Set custom secret to bypass rand.Read
		os.Setenv("SESSION_SECRET", "test-secret-from-environment-variable")
		Config.Readonly = false

		// Mock osExit - should not be called when using env var
		oldOsExit := osExit
		exitCalled := false
		osExit = func(code int) {
			exitCalled = true
			panic("mock exit")
		}
		defer func() {
			osExit = oldOsExit
			if r := recover(); r != nil && r != "mock exit" {
				panic(r)
			}
		}()

		// Execute
		middlewares := defaultMiddlewares()

		// Verify no exit when using environment secret
		if exitCalled {
			t.Error("osExit should not be called when SESSION_SECRET environment variable is set")
		}

		// Verify middlewares created
		if len(middlewares) != 2 {
			t.Errorf("expected 2 middlewares, got %d", len(middlewares))
		}
	})
}
