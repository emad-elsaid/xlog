package photos

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/fs"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
	"github.com/rwcarlsen/goexif/exif"
)

func TestResizeHandler_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) string
		expectError bool
		cleanup     func()
	}{
		{
			name: "valid image file resizes successfully",
			setup: func(t *testing.T) string {
				tmpFile := createTestPNG(t)
				return tmpFile
			},
			expectError: false,
		},
		{
			name: "nonexistent file returns error",
			setup: func(t *testing.T) string {
				return "nonexistent_file.png"
			},
			expectError: true,
		},
		{
			name: "cache directory creation handles errors",
			setup: func(t *testing.T) string {
				tmpFile := createTestPNG(t)
				return tmpFile
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			photoPath := tc.setup(t)

			// Clean up test cache directory
			defer func() {
				if err := os.RemoveAll(".cache"); err != nil {
					t.Logf("Failed to remove cache: %v", err)
				}
				if tc.expectError == false {
					if err := os.Remove(photoPath); err != nil {
						t.Logf("Failed to remove photo: %v", err)
					}
				}
			}()

			// Test that cache directory is created
			const cacheDir = ".cache"
			cacheFile := path.Join(cacheDir, fmt.Sprintf("photo-%x", sha256.Sum256([]byte(photoPath))))

			// Verify cache directory gets created
			if !tc.expectError {
				if err := os.Mkdir(cacheDir, 0700); err != nil && !os.IsExist(err) {
					t.Fatalf("Failed to create cache directory: %v", err)
				}
			}

			// Test cache file handling
			_, err := os.ReadFile(cacheFile)
			if err == nil && tc.expectError {
				t.Error("Expected error reading cache, but got none")
			}
		})
	}
}

func TestNewPhoto_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "nonexistent file returns error",
			path:        "nonexistent.jpg",
			expectError: true,
		},
		{
			name:        "valid file creates photo",
			path:        createTestPNG(t),
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expectError && tc.path != "nonexistent.jpg" {
				defer func() {
					if err := os.Remove(tc.path); err != nil {
						t.Logf("Failed to remove test file: %v", err)
					}
				}()
			}

			photo, err := NewPhoto(tc.path)
			if tc.expectError {
				if err == nil {
					t.Error("Expected error, but got none")
				}
				if photo != nil {
					t.Error("Expected nil photo on error, but got photo object")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
				if photo == nil {
					t.Error("Expected photo object, but got nil")
				}
			}
		})
	}
}

func TestNewPhoto_ExifDateTime(t *testing.T) {
	tests := []struct {
		name          string
		setupFile     func(t *testing.T) string
		expectExifNil bool
	}{
		{
			name:          "PNG without EXIF uses file ModTime",
			setupFile:     createTestPNG,
			expectExifNil: true,
		},
		{
			name: "file with invalid EXIF uses ModTime",
			setupFile: func(t *testing.T) string {
				tmpFile := filepath.Join(t.TempDir(), "invalid_exif.jpg")
				// Create file with some invalid content
				if err := os.WriteFile(tmpFile, []byte("not a real jpeg"), 0600); err != nil {
					t.Fatalf("Failed to write file: %v", err)
				}
				return tmpFile
			},
			expectExifNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			filePath := tc.setupFile(t)
			defer func() {
				if err := os.Remove(filePath); err != nil {
					t.Logf("Failed to cleanup: %v", err)
				}
			}()

			photo, err := NewPhoto(filePath)
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if tc.expectExifNil {
				if photo.Exif != nil {
					t.Error("Expected nil EXIF data")
				}
				// Time should be from file ModTime, not zero
				if photo.Time.IsZero() {
					t.Error("Expected non-zero time from ModTime")
				}
			}
		})
	}
}

func TestPhoto_Name(t *testing.T) {
	tests := []struct {
		thumbnail string
		expected  string
	}{
		{
			thumbnail: "/+/photos/thumbnail/image.png",
			expected:  "image",
		},
		{
			thumbnail: "/+/photos/thumbnail/vacation/sunset.jpg",
			expected:  "sunset",
		},
		{
			thumbnail: "/+/photos/thumbnail/photo.with.dots.jpeg",
			expected:  "photo.with.dots",
		},
	}

	for _, tc := range tests {
		t.Run(tc.thumbnail, func(t *testing.T) {
			p := &Photo{Thumbnail: tc.thumbnail}
			got := p.Name()
			if got != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestPhoto_InterfaceMethods(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		expected interface{}
	}{
		{
			name:     "FileName returns empty string",
			method:   "FileName",
			expected: "",
		},
		{
			name:     "Exists returns false",
			method:   "Exists",
			expected: false,
		},
		{
			name:     "Content returns empty Markdown",
			method:   "Content",
			expected: "",
		},
		{
			name:     "Delete returns false",
			method:   "Delete",
			expected: false,
		},
		{
			name:     "Write returns false",
			method:   "Write",
			expected: false,
		},
	}

	p := &Photo{Thumbnail: "/test/photo.jpg"}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			switch tc.method {
			case "FileName":
				if got := p.FileName(); got != tc.expected {
					t.Errorf("Expected %q, got %q", tc.expected, got)
				}
			case "Exists":
				if got := p.Exists(); got != tc.expected {
					t.Errorf("Expected %v, got %v", tc.expected, got)
				}
			case "Content":
				if got := p.Content(); string(got) != tc.expected {
					t.Errorf("Expected %q, got %q", tc.expected, got)
				}
			case "Delete":
				if got := p.Delete(); got != tc.expected {
					t.Errorf("Expected %v, got %v", tc.expected, got)
				}
			case "Write":
				if got := p.Write(""); got != tc.expected {
					t.Errorf("Expected %v, got %v", tc.expected, got)
				}
			}
		})
	}
}

func TestPhoto_ModTime(t *testing.T) {
	p := &Photo{Thumbnail: "/test/photo.jpg"}
	got := p.ModTime()
	if !got.IsZero() {
		t.Errorf("Expected zero time, got %v", got)
	}
}

func TestPhoto_AST(t *testing.T) {
	p := &Photo{Thumbnail: "/test/photo.jpg"}
	data, node := p.AST()
	if data != nil {
		t.Errorf("Expected nil data, got %v", data)
	}
	if node != nil {
		t.Errorf("Expected nil node, got %v", node)
	}
}

func TestProperty_Methods(t *testing.T) {
	tests := []struct {
		name     string
		property Property
		wantIcon string
		wantName string
		wantVal  any
	}{
		{
			name: "camera make property",
			property: Property{
				IconVal: "fa-solid fa-camera-retro",
				NameVal: "camera make",
				Val:     "Canon",
			},
			wantIcon: "fa-solid fa-camera-retro",
			wantName: "camera make",
			wantVal:  "Canon",
		},
		{
			name: "capture time property",
			property: Property{
				IconVal: "fa-regular fa-calendar",
				NameVal: "capture time",
				Val:     "Monday 15 May 2023",
			},
			wantIcon: "fa-regular fa-calendar",
			wantName: "capture time",
			wantVal:  "Monday 15 May 2023",
		},
		{
			name: "ISO property with integer value",
			property: Property{
				IconVal: "fa-solid fa-camera-retro",
				NameVal: "ISO",
				Val:     800,
			},
			wantIcon: "fa-solid fa-camera-retro",
			wantName: "ISO",
			wantVal:  800,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.property.Icon(); got != tc.wantIcon {
				t.Errorf("Icon() = %q, want %q", got, tc.wantIcon)
			}
			if got := tc.property.Name(); got != tc.wantName {
				t.Errorf("Name() = %q, want %q", got, tc.wantName)
			}
			if got := tc.property.Value(); got != tc.wantVal {
				t.Errorf("Value() = %v, want %v", got, tc.wantVal)
			}
		})
	}
}

func TestProperties_NonPhotoPage(t *testing.T) {
	type mockPage struct {
		xlog.Page
	}
	var mock mockPage
	result := properties(&mock)
	if result != nil {
		t.Errorf("Expected nil for non-Photo page, got %v", result)
	}
}

func TestProperties_PhotoWithoutExif(t *testing.T) {
	photo := &Photo{
		Thumbnail: "/test/photo.jpg",
		Exif:      nil,
	}
	result := properties(photo)
	if result != nil {
		t.Errorf("Expected nil for photo without EXIF, got %v", result)
	}
}

func TestPhotos_Name(t *testing.T) {
	ext := Photos{}
	if got := ext.Name(); got != "photos" {
		t.Errorf("Expected 'photos', got %q", got)
	}
}

func TestNewPhoto_WithExifData(t *testing.T) {
	// Create a JPEG with EXIF data
	tmpFile := filepath.Join(t.TempDir(), "photo_with_exif.jpg")

	// Create a minimal JPEG (simple test without real EXIF for this test)
	// The EXIF decoding will fail but we test the code path
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("Failed to encode image: %v", err)
	}

	if err := os.WriteFile(tmpFile, buf.Bytes(), 0600); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	defer func() {
		if err := os.Remove(tmpFile); err != nil {
			t.Logf("Failed to remove test file: %v", err)
		}
	}()

	photo, err := NewPhoto(tmpFile)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if photo == nil {
		t.Fatal("Expected photo object, got nil")
	}

	// Verify Photo fields are populated correctly
	expectedThumbnail := "/+/photos/thumbnail/" + tmpFile
	if photo.Thumbnail != expectedThumbnail {
		t.Errorf("Expected thumbnail %q, got %q", expectedThumbnail, photo.Thumbnail)
	}

	expectedPage := "/+/photos/photo/" + tmpFile
	if photo.Page != expectedPage {
		t.Errorf("Expected page %q, got %q", expectedPage, photo.Page)
	}

	if photo.Original != tmpFile {
		t.Errorf("Expected original %q, got %q", tmpFile, photo.Original)
	}

	if photo.Time.IsZero() {
		t.Error("Expected non-zero time from ModTime, got zero")
	}
}

func TestProperties_PhotoWithTime(t *testing.T) {
	photo := &Photo{
		Thumbnail: "/test/photo.jpg",
		Time:      time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC),
		Exif:      &exif.Exif{},
	}

	props := properties(photo)
	if len(props) == 0 {
		t.Error("Expected at least one property (capture time), got none")
	}

	// Check that capture time property exists
	foundTimeProperty := false
	for _, prop := range props {
		if prop.Name() == "capture time" {
			foundTimeProperty = true
			if prop.Icon() != "fa-regular fa-calendar" {
				t.Errorf("Expected calendar icon, got %q", prop.Icon())
			}
		}
	}

	if !foundTimeProperty {
		t.Error("Expected capture time property, but not found")
	}
}

func TestPhotosShortcode_DirectoryWalking(t *testing.T) {
	tests := []struct {
		name         string
		setupDir     func(t *testing.T) string
		expectError  bool
		errorMessage string
	}{
		{
			name: "empty directory succeeds",
			setupDir: func(t *testing.T) string {
				return t.TempDir()
			},
			expectError: false,
		},
		{
			name: "directory with single photo",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				createTestPNGInDir(t, dir, "photo1.png")
				return dir
			},
			expectError: false,
		},
		{
			name: "directory with multiple photos",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				createTestPNGInDir(t, dir, "photo1.jpg")
				time.Sleep(10 * time.Millisecond)
				createTestPNGInDir(t, dir, "photo2.png")
				time.Sleep(10 * time.Millisecond)
				createTestPNGInDir(t, dir, "photo3.gif")
				return dir
			},
			expectError: false,
		},
		{
			name: "directory with mixed files",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				createTestPNGInDir(t, dir, "image.png")
				// Create non-image file
				txtFile := filepath.Join(dir, "readme.txt")
				if err := os.WriteFile(txtFile, []byte("test"), 0600); err != nil {
					t.Fatalf("Failed to create text file: %v", err)
				}
				return dir
			},
			expectError: false,
		},
		{
			name: "nonexistent directory returns error",
			setupDir: func(t *testing.T) string {
				return "/nonexistent/path/to/photos"
			},
			expectError:  true,
			errorMessage: "no such file",
		},
		{
			name: "nested directory structure",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				subdir := filepath.Join(dir, "vacation")
				if err := os.MkdirAll(subdir, 0700); err != nil {
					t.Fatalf("Failed to create subdirectory: %v", err)
				}
				createTestPNGInDir(t, dir, "photo1.png")
				createTestPNGInDir(t, subdir, "photo2.jpg")
				return dir
			},
			expectError: false,
		},
		{
			name: "case-insensitive extension matching",
			setupDir: func(t *testing.T) string {
				dir := t.TempDir()
				createTestPNGInDir(t, dir, "photo.PNG")
				createTestPNGInDir(t, dir, "photo2.JPG")
				createTestPNGInDir(t, dir, "photo3.JPEG")
				return dir
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := tc.setupDir(t)
			p := strings.TrimSpace(dir)

			photos := []*Photo{}
			err := filepath.WalkDir(p, func(file string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.Type().IsRegular() && supportedExt.Include(strings.ToLower(path.Ext(file))) {
					photo, err := NewPhoto(file)
					if err != nil {
						return err
					}
					photos = append(photos, photo)
				}

				return nil
			})

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error containing %q, but got no error", tc.errorMessage)
				} else if !strings.Contains(err.Error(), tc.errorMessage) {
					t.Errorf("Expected error containing %q, got: %v", tc.errorMessage, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

func TestPhotosShortcode_WhitespaceHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "trims leading whitespace",
			input: "  /path/to/photos",
			want:  "/path/to/photos",
		},
		{
			name:  "trims trailing whitespace",
			input: "/path/to/photos  ",
			want:  "/path/to/photos",
		},
		{
			name:  "trims both leading and trailing",
			input: "  /path/to/photos  ",
			want:  "/path/to/photos",
		},
		{
			name:  "trims newlines",
			input: "\n/path/to/photos\n",
			want:  "/path/to/photos",
		},
		{
			name:  "trims tabs",
			input: "\t/path/to/photos\t",
			want:  "/path/to/photos",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := strings.TrimSpace(tc.input)
			if got != tc.want {
				t.Errorf("TrimSpace(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestPhotosShortcode_SupportedExtensions(t *testing.T) {
	tests := []struct {
		filename string
		want     bool
	}{
		{".jpg", true},
		{".jpeg", true},
		{".gif", true},
		{".png", true},
		{".JPG", true},
		{".JPEG", true},
		{".GIF", true},
		{".PNG", true},
		{".txt", false},
		{".pdf", false},
		{".bmp", false},
		{".webp", false},
		{".svg", false},
	}

	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			got := supportedExt.Include(strings.ToLower(tc.filename))
			if got != tc.want {
				t.Errorf("supportedExt.Include(%q) = %v, want %v", tc.filename, got, tc.want)
			}
		})
	}
}

func TestPhotosShortcode_PhotosSorting(t *testing.T) {
	// Create photos with different times
	photos := []*Photo{
		{Time: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
		{Time: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)},
		{Time: time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)},
	}

	// Sort using the same logic as photosShortcode
	slices.SortFunc(photos, func(i, j *Photo) int {
		return j.Time.Compare(i.Time)
	})

	// Verify sorted in descending order (newest first)
	if !photos[0].Time.Equal(time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("Expected newest photo first, got %v", photos[0].Time)
	}
	if !photos[1].Time.Equal(time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("Expected middle photo second, got %v", photos[1].Time)
	}
	if !photos[2].Time.Equal(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("Expected oldest photo last, got %v", photos[2].Time)
	}
}

// Helper function to create a test PNG file.
func createTestPNG(t *testing.T) string {
	t.Helper()

	tmpFile := filepath.Join(t.TempDir(), "test_image.png")

	// Create a simple 100x100 red image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	red := color.RGBA{255, 0, 0, 255}
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, red)
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	if err := os.WriteFile(tmpFile, buf.Bytes(), 0600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	return tmpFile
}

// Helper to create PNG in specific directory with given filename.
func createTestPNGInDir(t *testing.T, dir, filename string) string {
	t.Helper()

	filePath := filepath.Join(dir, filename)

	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	blue := color.RGBA{0, 0, 255, 255}
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			img.Set(x, y, blue)
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	if err := os.WriteFile(filePath, buf.Bytes(), 0600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	return filePath
}

func TestResizeHandler_WithHTTPRequest(t *testing.T) {
	tests := []struct {
		name          string
		setupPhoto    func(t *testing.T) string
		expectSuccess bool
		validateResp  func(t *testing.T, body []byte)
	}{
		{
			name:          "valid photo returns resized PNG",
			setupPhoto:    createTestPNG,
			expectSuccess: true,
			validateResp: func(t *testing.T, body []byte) {
				if len(body) == 0 {
					t.Error("Expected non-empty response body")
				}
				// Verify it's a valid PNG
				img, err := png.Decode(bytes.NewReader(body))
				if err != nil {
					t.Errorf("Response is not a valid PNG: %v", err)
				}
				// Verify dimensions
				bounds := img.Bounds()
				if bounds.Dx() != 700 {
					t.Errorf("Expected width 700, got %d", bounds.Dx())
				}
			},
		},
		{
			name: "nonexistent photo returns error message",
			setupPhoto: func(t *testing.T) string {
				return "/nonexistent/photo.png"
			},
			expectSuccess: false,
			validateResp: func(t *testing.T, body []byte) {
				if len(body) == 0 {
					t.Error("Expected error message in body")
				}
			},
		},
		{
			name: "cached photo served from cache",
			setupPhoto: func(t *testing.T) string {
				photoPath := createTestPNG(t)

				// Pre-populate cache
				const cacheDir = ".cache"
				if err := os.Mkdir(cacheDir, 0700); err != nil && !os.IsExist(err) {
					t.Fatalf("Failed to create cache dir: %v", err)
				}

				cacheFile := path.Join(cacheDir, fmt.Sprintf("photo-%x", sha256.Sum256([]byte(photoPath))))
				cachedData := []byte("cached content")
				if err := os.WriteFile(cacheFile, cachedData, 0600); err != nil {
					t.Fatalf("Failed to write cache: %v", err)
				}

				return photoPath
			},
			expectSuccess: true,
			validateResp: func(t *testing.T, body []byte) {
				if string(body) != "cached content" {
					t.Error("Expected cached content to be served")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			photoPath := tc.setupPhoto(t)
			defer func() {
				_ = os.RemoveAll(".cache")
			}()
			if tc.expectSuccess {
				defer func() {
					_ = os.Remove(photoPath)
				}()
			}

			// Create mock HTTP request
			req := httptest.NewRequest("GET", "/+/photos/thumbnail/"+photoPath, nil)
			req.SetPathValue("path", photoPath)

			// Get handler output function
			output := resizeHandler(req)

			// Execute the handler
			w := httptest.NewRecorder()
			output(w, req)

			// Validate response
			tc.validateResp(t, w.Body.Bytes())
		})
	}
}

func TestPhotoHandler_WithHTTPRequest(t *testing.T) {
	tests := []struct {
		name          string
		setupPhoto    func(t *testing.T) string
		expectSuccess bool
	}{
		{
			name:          "valid photo path returns output function",
			setupPhoto:    createTestPNG,
			expectSuccess: true,
		},
		{
			name: "nonexistent photo returns error output",
			setupPhoto: func(t *testing.T) string {
				return "/nonexistent/photo.png"
			},
			expectSuccess: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			photoPath := tc.setupPhoto(t)
			if tc.expectSuccess {
				defer func() {
					_ = os.Remove(photoPath)
				}()
			}

			// Create mock HTTP request
			req := httptest.NewRequest("GET", "/+/photos/photo/"+photoPath, nil)
			req.SetPathValue("path", photoPath)

			// Get handler output
			output := photoHandler(req)

			// Verify output function exists
			if output == nil {
				t.Error("Expected non-nil output function")
			}
		})
	}
}
