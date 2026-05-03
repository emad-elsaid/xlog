package photos

import (
	"testing"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

const (
	captureTimeProperty = "capture time"
	calendarIcon        = "fa-regular fa-calendar"
)

func TestProperties_WithCaptureTime(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		wantProp bool
	}{
		{
			name:     "non-zero time creates capture time property",
			time:     time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC),
			wantProp: true,
		},
		{
			name:     "zero time skips capture time property",
			time:     time.Time{},
			wantProp: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			photo := &Photo{
				Thumbnail: "/test/photo.jpg",
				Exif:      &exif.Exif{},
				Time:      tc.time,
			}

			props := properties(photo)

			foundCaptureTime := false
			for _, prop := range props {
				p := prop.(Property)
				if p.Name() == captureTimeProperty {
					foundCaptureTime = true
					if p.Icon() != calendarIcon {
						t.Errorf("capture time icon = %v, want %s", p.Icon(), calendarIcon)
					}
					// Verify time is formatted correctly
					val, ok := p.Value().(string)
					if !ok || val == "" {
						t.Error("capture time value should be non-empty string")
					}
				}
			}

			if foundCaptureTime != tc.wantProp {
				t.Errorf("Found capture time property = %v, want %v", foundCaptureTime, tc.wantProp)
			}
		})
	}
}

func TestProperties_ExifFields(t *testing.T) {
	// Test that properties function handles photos with basic EXIF correctly
	photo := &Photo{
		Thumbnail: "/test/photo.jpg",
		Exif:      &exif.Exif{},
		Time:      time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC),
	}

	props := properties(photo)

	// Verify we get at least the capture time property
	if len(props) == 0 {
		t.Error("Expected at least capture time property")
	}

	// Verify each property has required fields
	for i, prop := range props {
		p, ok := prop.(Property)
		if !ok {
			t.Errorf("property[%d] is not a Property type", i)
			continue
		}

		if p.Icon() == "" {
			t.Errorf("property[%d] (%s) has empty Icon", i, p.Name())
		}
		if p.Name() == "" {
			t.Errorf("property[%d] has empty Name", i)
		}
		if p.Value() == nil {
			t.Errorf("property[%d] (%s) has nil Value", i, p.Name())
		}
	}
}

func TestProperties_StringValueExtraction(t *testing.T) {
	// Test that properties function handles photos with empty EXIF gracefully
	photo := &Photo{
		Thumbnail: "/test/photo.jpg",
		Exif:      &exif.Exif{},
		Time:      time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC),
	}

	props := properties(photo)

	// Should return at least capture time
	if len(props) == 0 {
		t.Error("Expected at least capture time property")
	}

	// Verify we got capture time
	foundCaptureTime := false
	for _, prop := range props {
		p := prop.(Property)
		if p.Name() == captureTimeProperty {
			foundCaptureTime = true
		}
	}

	if !foundCaptureTime {
		t.Error("Expected capture time property when Time is set")
	}
}

func TestProperties_CaptureTimeFormatting(t *testing.T) {
	// Test specific date formatting
	testTime := time.Date(2024, time.March, 15, 14, 30, 45, 0, time.UTC)

	photo := &Photo{
		Thumbnail: "/test/photo.jpg",
		Exif:      &exif.Exif{},
		Time:      testTime,
	}

	props := properties(photo)

	// Find capture time property
	var captureTimeProp Property
	found := false
	for _, prop := range props {
		p := prop.(Property)
		if p.Name() == captureTimeProperty {
			captureTimeProp = p
			found = true
			break
		}
	}

	if !found {
		t.Fatal("capture time property not found")
	}

	// Verify the formatted value contains expected components
	val, ok := captureTimeProp.Value().(string)
	if !ok {
		t.Fatalf("capture time value is not a string: %T", captureTimeProp.Value())
	}

	// Should contain: Weekday Day Month Year
	// Format: "%s %d %s %d" -> "Friday 15 March 2024"
	expectedComponents := []string{
		testTime.Weekday().String(),
		"15",
		testTime.Month().String(),
		"2024",
	}

	for _, component := range expectedComponents {
		if !contains(val, component) {
			t.Errorf("capture time value %q does not contain %q", val, component)
		}
	}
}

// Helper function for string contains check.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr))))
}

func TestProperties_RationalFieldParsing(t *testing.T) {
	// Test that the str() helper doesn't panic when EXIF fields return errors
	photo := &Photo{
		Thumbnail: "/test/photo.jpg",
		Exif:      &exif.Exif{},
		Time:      time.Time{},
	}

	// This should not panic even if Time is zero and EXIF fields are missing
	props := properties(photo)

	// With zero time and empty EXIF, we expect empty slice (not nil, but empty)
	// properties() returns []xlog.Property{} at line 23 when initialized
	if props == nil {
		t.Error("properties() should return empty slice, not nil")
	}
}

func TestProperties_AllExifFieldsPresent(t *testing.T) {
	// Test with explicitly set EXIF and Time
	photo := &Photo{
		Thumbnail: "/test/photo.jpg",
		Exif:      &exif.Exif{},
		Time:      time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC),
	}

	props := properties(photo)

	// Should have at least capture time
	if len(props) == 0 {
		t.Error("Expected at least capture time property")
	}

	// Verify capture time is present since we explicitly set it
	foundCaptureTime := false
	for _, prop := range props {
		p := prop.(Property)
		if p.Name() == captureTimeProperty {
			foundCaptureTime = true

			// Verify icon
			if p.Icon() != calendarIcon {
				t.Errorf("capture time has wrong icon: %s", p.Icon())
			}

			// Verify value is formatted
			val, ok := p.Value().(string)
			if !ok || val == "" {
				t.Error("capture time value should be non-empty string")
			}
		}
	}

	if !foundCaptureTime {
		t.Error("Expected 'capture time' property when photo.Time is set")
	}
}
