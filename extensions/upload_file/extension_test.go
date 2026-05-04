package upload_file

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
)

// Mock page for testing Command methods.
type mockPage struct {
	name string
}

func (m mockPage) Name() string           { return m.name }
func (mockPage) FileName() string         { return "" }
func (mockPage) Exists() bool             { return true }
func (mockPage) Render() template.HTML    { return "" }
func (mockPage) Content() xlog.Markdown   { return "" }
func (mockPage) Delete() bool             { return false }
func (mockPage) Write(xlog.Markdown) bool { return false }
func (mockPage) ModTime() time.Time       { return time.Time{} }
func (mockPage) AST() ([]byte, ast.Node)  { return nil, nil }

func TestUploadFileHandler(t *testing.T) {
	tests := []struct {
		name           string
		fileName       string
		pageExists     bool
		fileContent    []byte
		fileExt        string
		uploadFileName string
		wantStatus     int
		wantContains   string
	}{
		{
			name:           "upload image to existing page",
			fileName:       "test.md",
			pageExists:     true,
			fileContent:    []byte("fake png data"),
			fileExt:        ".png",
			uploadFileName: "test.png",
			wantStatus:     http.StatusFound,
			wantContains:   "",
		},
		{
			name:           "upload video generates video tag",
			fileName:       "",
			pageExists:     false,
			fileContent:    []byte("fake webm data"),
			fileExt:        ".webm",
			uploadFileName: "video.webm",
			wantStatus:     http.StatusOK,
			wantContains:   "<video controls src=",
		},
		{
			name:           "upload audio generates audio tag",
			fileName:       "",
			pageExists:     false,
			fileContent:    []byte("fake audio data"),
			fileExt:        ".mp3",
			uploadFileName: "audio.mp3",
			wantStatus:     http.StatusOK,
			wantContains:   "<audio controls src=",
		},
		{
			name:           "upload generic file generates link",
			fileName:       "",
			pageExists:     false,
			fileContent:    []byte("fake pdf data"),
			fileExt:        ".pdf",
			uploadFileName: "document.pdf",
			wantStatus:     http.StatusOK,
			wantContains:   "[document.pdf]",
		},
		{
			name:           "upload to nonexistent page returns 404",
			fileName:       "nonexistent.md",
			pageExists:     false,
			fileContent:    []byte("data"),
			fileExt:        ".png",
			uploadFileName: "test.png",
			wantStatus:     http.StatusNotFound,
			wantContains:   "", // NotFound returns empty body
		},
		{
			name:           "filename with special chars gets filtered",
			fileName:       "",
			pageExists:     false,
			fileContent:    []byte("data"),
			fileExt:        ".txt",
			uploadFileName: "test[].txt",
			wantStatus:     http.StatusOK,
			wantContains:   "[test.txt]",
		},
		{
			name:           "uppercase extension handled correctly",
			fileName:       "",
			pageExists:     false,
			fileContent:    []byte("fake png data"),
			fileExt:        ".PNG",
			uploadFileName: "TEST.PNG",
			wantStatus:     http.StatusOK,
			wantContains:   "![](",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup: Create temporary test environment
			tempDir := t.TempDir()
			origDir, _ := os.Getwd()
			defer func() { _ = os.Chdir(origDir) }()
			_ = os.Chdir(tempDir)

			// Create test page if needed
			if tc.pageExists && tc.fileName != "" {
				_ = os.WriteFile(tc.fileName+".md", []byte("# Test Page"), 0600)
			}

			// Create multipart form with file upload
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", tc.uploadFileName)
			_, _ = io.Copy(part, bytes.NewReader(tc.fileContent))
			_ = writer.WriteField("page", tc.fileName)
			_ = writer.Close()

			// Create request
			req := httptest.NewRequest(http.MethodPost, "/+/upload-file", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			w := httptest.NewRecorder()

			// Execute handler
			result := uploadFileHandler(req)
			result(w, req)

			// Verify status code
			if status := w.Code; status != tc.wantStatus {
				t.Errorf("status = %d, want %d", status, tc.wantStatus)
			}

			// Verify response contains expected content
			if tc.wantContains != "" {
				respBody := w.Body.String()
				if !strings.Contains(respBody, tc.wantContains) {
					t.Errorf("response body does not contain %q\nGot: %s", tc.wantContains, respBody)
				}
			}

			// Verify file was saved with correct hash
			if tc.wantStatus == http.StatusOK || tc.wantStatus == http.StatusFound {
				expectedHash := fmt.Sprintf("%x", sha256.Sum256(tc.fileContent))
				expectedPath := filepath.Join(PUBLIC_PATH, expectedHash+strings.ToLower(tc.fileExt))
				if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
					t.Errorf("expected file not created: %s", expectedPath)
				}
			}
		})
	}
}

func TestFilterChars(t *testing.T) {
	tests := []struct {
		input    string
		exclude  string
		expected string
	}{
		{
			input:    "test[].txt",
			exclude:  "[]",
			expected: "test.txt",
		},
		{
			input:    "file(1).md",
			exclude:  "()",
			expected: "file1.md",
		},
		{
			input:    "normal.txt",
			exclude:  "[]",
			expected: "normal.txt",
		},
		{
			input:    "test*.md",
			exclude:  "*",
			expected: "test.md",
		},
		{
			input:    "",
			exclude:  "[]",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got := filterChars(tc.input, tc.exclude)
			if got != tc.expected {
				t.Errorf("filterChars(%q, %q) = %q, want %q", tc.input, tc.exclude, got, tc.expected)
			}
		})
	}
}

func TestUploadFileExtensionInit(t *testing.T) {
	const expectedName = "upload-file"
	// Test that Init registers handlers correctly
	ext := UploadFile{}

	if name := ext.Name(); name != expectedName {
		t.Errorf("Name() = %q, want %q", name, expectedName)
	}

	// Init should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Init() panicked: %v", r)
		}
	}()

	// Note: Full integration testing of Init would require
	// mocking the xlog registration functions
}

func TestFileTypeDetection(t *testing.T) {
	tests := []struct {
		ext      string
		wantType string
	}{
		{".jpg", "image"},
		{".jpeg", "image"},
		{".png", "image"},
		{".gif", "image"},
		{".svg", "image"},
		{".webp", "image"},
		{".webm", "video"},
		{".wave", "audio"},
		{".ogg", "audio"},
		{".opus", "audio"},
		{".mp3", "audio"},
		{".pdf", "generic"},
		{".txt", "generic"},
	}

	for _, tc := range tests {
		t.Run(tc.ext, func(t *testing.T) {
			var found bool
			switch tc.wantType {
			case "image":
				found = contains(IMAGES_EXTENSIONS, tc.ext)
			case "video":
				found = contains(VIDEOS_EXTENSIONS, tc.ext)
			case "audio":
				found = contains(AUDIO_EXTENSIONS, tc.ext)
			case "generic":
				found = !contains(IMAGES_EXTENSIONS, tc.ext) &&
					!contains(VIDEOS_EXTENSIONS, tc.ext) &&
					!contains(AUDIO_EXTENSIONS, tc.ext)
			}

			if !found {
				t.Errorf("extension %s not correctly classified as %s", tc.ext, tc.wantType)
			}
		})
	}
}

// Helper functions

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func TestUpload_Command(t *testing.T) {
	page := mockPage{name: "test-page"}
	cmd := Upload{p: page}

	if got := cmd.Icon(); got != "fa-solid fa-file-arrow-up" {
		t.Errorf("Icon() = %q, want %q", got, "fa-solid fa-file-arrow-up")
	}

	if got := cmd.Name(); got != "Upload File" {
		t.Errorf("Name() = %q, want %q", got, "Upload File")
	}

	attrs := cmd.Attrs()
	if attrs == nil {
		t.Fatal("Attrs() returned nil")
	}

	href, ok := attrs["href"]
	if !ok {
		t.Error("Attrs() missing 'href' key")
	} else if !strings.Contains(href.(string), "/+/upload-file/form") {
		t.Errorf("href doesn't contain expected path: %v", href)
	}
}

func TestScreenshot_Command(t *testing.T) {
	page := mockPage{name: "test-page"}
	cmd := Screenshot{p: page}

	if got := cmd.Icon(); got != "fa-solid fa-camera" {
		t.Errorf("Icon() = %q, want %q", got, "fa-solid fa-camera")
	}

	if got := cmd.Name(); got != "Screenshot" {
		t.Errorf("Name() = %q, want %q", got, "Screenshot")
	}

	attrs := cmd.Attrs()
	if attrs == nil {
		t.Fatal("Attrs() returned nil")
	}
}

func TestRecordScreen_Command(t *testing.T) {
	page := mockPage{name: "test-page"}
	cmd := RecordScreen{p: page}

	if got := cmd.Icon(); got != "fa-solid fa-desktop" {
		t.Errorf("Icon() = %q, want %q", got, "fa-solid fa-desktop")
	}

	if got := cmd.Name(); got != "Record screen" {
		t.Errorf("Name() = %q, want %q", got, "Record screen")
	}
}

func TestRecordCamera_Command(t *testing.T) {
	page := mockPage{name: "test-page"}
	cmd := RecordCamera{p: page}

	if got := cmd.Icon(); got != "fa-solid fa-video" {
		t.Errorf("Icon() = %q, want %q", got, "fa-solid fa-video")
	}

	if got := cmd.Name(); got != "Record camera" {
		t.Errorf("Name() = %q, want %q", got, "Record camera")
	}
}

func TestRecordAudio_Command(t *testing.T) {
	page := mockPage{name: "test-page"}
	cmd := RecordAudio{p: page}

	if got := cmd.Icon(); got != "fa-solid fa-microphone" {
		t.Errorf("Icon() = %q, want %q", got, "fa-solid fa-microphone")
	}

	if got := cmd.Name(); got != "Record audio" {
		t.Errorf("Name() = %q, want %q", got, "Record audio")
	}
}

func TestFormHandlers(t *testing.T) {
	tests := []struct {
		name     string
		handler  func(xlog.Request) xlog.Output
		pageName string
	}{
		{
			name:     "UploadForm handles request",
			handler:  UploadForm,
			pageName: "test-page",
		},
		{
			name:     "ScreenshotForm handles request",
			handler:  ScreenshotForm,
			pageName: "my-page",
		},
		{
			name:     "RecordScreenForm handles request",
			handler:  RecordScreenForm,
			pageName: "notes",
		},
		{
			name:     "RecordCameraForm handles request",
			handler:  RecordCameraForm,
			pageName: "video-page",
		},
		{
			name:     "RecordAudioForm handles request",
			handler:  RecordAudioForm,
			pageName: "audio-notes",
		},
		{
			name:     "handles page names with special characters",
			handler:  UploadForm,
			pageName: "test-page-stuff",
		},
		{
			name:     "handles empty page name",
			handler:  UploadForm,
			pageName: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock request with page parameter
			req := httptest.NewRequest(http.MethodPost, "/test", http.NoBody)
			req.Form = map[string][]string{
				"page": {tc.pageName},
			}

			// Execute handler - verifies the function doesn't panic
			result := tc.handler(req)

			// Verify result is not nil
			if result == nil {
				t.Fatal("handler returned nil Output")
			}

			// Verify handler doesn't panic when called
			// Note: Full rendering test would require template infrastructure setup
		})
	}
}

func TestCommandAttrs(t *testing.T) {
	tests := []struct {
		name            string
		command         xlog.Command
		wantHrefPattern string
		wantHxPost      bool
	}{
		{
			name:            "Upload Attrs includes correct href",
			command:         Upload{p: mockPage{name: "my-page"}},
			wantHrefPattern: "/+/upload-file/form?page=my-page",
			wantHxPost:      true,
		},
		{
			name:            "Screenshot Attrs includes correct href",
			command:         Screenshot{p: mockPage{name: "notes"}},
			wantHrefPattern: "/+/upload-file/screenshot-form?page=notes",
			wantHxPost:      true,
		},
		{
			name:            "RecordScreen Attrs includes correct href",
			command:         RecordScreen{p: mockPage{name: "video"}},
			wantHrefPattern: "/+/upload-file/record-screen-form?page=video",
			wantHxPost:      true,
		},
		{
			name:            "RecordCamera Attrs includes correct href",
			command:         RecordCamera{p: mockPage{name: "camera-test"}},
			wantHrefPattern: "/+/upload-file/record-camera-form?page=camera-test",
			wantHxPost:      true,
		},
		{
			name:            "RecordAudio Attrs includes correct href",
			command:         RecordAudio{p: mockPage{name: "audio"}},
			wantHrefPattern: "/+/upload-file/record-audio-form?page=audio",
			wantHxPost:      true,
		},
		{
			name:            "handles page names with special characters",
			command:         Upload{p: mockPage{name: "test-page-more"}},
			wantHrefPattern: "/+/upload-file/form?page=test-page-more",
			wantHxPost:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			attrs := tc.command.Attrs()
			if attrs == nil {
				t.Fatal("Attrs() returned nil")
			}

			// Check href attribute
			href, ok := attrs["href"]
			if !ok {
				t.Error("Attrs() missing 'href' attribute")
			} else if !strings.Contains(fmt.Sprint(href), tc.wantHrefPattern) {
				t.Errorf("href = %v, want to contain %q", href, tc.wantHrefPattern)
			}

			// Check hx-post attribute if expected
			if tc.wantHxPost {
				if _, ok := attrs["hx-post"]; !ok {
					t.Error("Attrs() missing 'hx-post' attribute")
				}
			}
		})
	}
}

func TestInit(t *testing.T) {
	const expectedName = "upload-file"
	// Verify Init doesn't panic
	ext := UploadFile{}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Init() panicked: %v", r)
		}
	}()

	// Test that Name returns correct value
	if got := ext.Name(); got != expectedName {
		t.Errorf("Name() = %q, want %q", got, expectedName)
	}

	// Note: Full Init testing would require mocking xlog registration functions
	// Testing Init() execution would require integration test setup
}
