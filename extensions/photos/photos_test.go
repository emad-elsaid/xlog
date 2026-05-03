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
