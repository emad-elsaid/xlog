package xlog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			r := httptest.NewRequest(http.MethodDelete, "/", nil)

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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
	r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			expectedBody:   "\n",
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
			expectedBody:   "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			expectedBody:   "\n",
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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			r := httptest.NewRequest(tc.expectedMethod, tc.path, nil)
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
			r2 := httptest.NewRequest(wrongMethod, tc.path, nil)
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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
			r := httptest.NewRequest(http.MethodGet, "/", nil)

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
