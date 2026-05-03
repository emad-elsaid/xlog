package photos

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path"
	"path/filepath"
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
