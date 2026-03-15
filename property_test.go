package xlog

import (
	"html/template"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog/markdown/ast"
)

// testPropertyPage implements the Page interface for testing properties
type testPropertyPage struct {
	name    string
	modTime time.Time
}

func (t *testPropertyPage) Name() string                 { return t.name }
func (t *testPropertyPage) FileName() string             { return t.name + ".md" }
func (t *testPropertyPage) Exists() bool                 { return true }
func (t *testPropertyPage) Render() template.HTML        { return "" }
func (t *testPropertyPage) Content() Markdown            { return "" }
func (t *testPropertyPage) Delete() bool                 { return false }
func (t *testPropertyPage) Write(Markdown) bool          { return false }
func (t *testPropertyPage) ModTime() time.Time           { return t.modTime }
func (t *testPropertyPage) AST() ([]byte, ast.Node)      { return nil, nil }

// testProperty is a simple property implementation for testing
type testProperty struct {
	icon  string
	name  string
	value any
}

func (tp testProperty) Icon() string  { return tp.icon }
func (tp testProperty) Name() string  { return tp.name }
func (tp testProperty) Value() any    { return tp.value }

func TestRegisterProperty(t *testing.T) {
	// Save original state
	originalPropsSources := propsSources
	defer func() { propsSources = originalPropsSources }()

	// Reset to default
	propsSources = []func(Page) []Property{defaultProps}

	// Register a new property source
	customProp := func(p Page) []Property {
		return []Property{
			testProperty{icon: "fa-custom", name: "custom", value: "test-value"},
		}
	}

	RegisterProperty(customProp)

	if len(propsSources) != 2 {
		t.Errorf("Expected 2 property sources, got %d", len(propsSources))
	}
}

func TestPropertiesDefaultProps(t *testing.T) {
	// Save original state
	originalPropsSources := propsSources
	defer func() { propsSources = originalPropsSources }()

	// Reset to default
	propsSources = []func(Page) []Property{defaultProps}

	testTime := time.Now().Add(-2 * time.Hour)
	testP := &testPropertyPage{
		name:    "test-page",
		modTime: testTime,
	}

	props := Properties(testP)

	if len(props) != 1 {
		t.Errorf("Expected 1 default property, got %d", len(props))
	}

	modifiedProp, exists := props["modified"]
	if !exists {
		t.Fatal("Expected 'modified' property to exist")
	}

	if modifiedProp.Name() != "modified" {
		t.Errorf("Expected property name 'modified', got %s", modifiedProp.Name())
	}

	if modifiedProp.Icon() != "fa-solid fa-clock" {
		t.Errorf("Expected icon 'fa-solid fa-clock', got %s", modifiedProp.Icon())
	}

	// Value should be string from ago() function
	if _, ok := modifiedProp.Value().(string); !ok {
		t.Errorf("Expected Value() to return string, got %T", modifiedProp.Value())
	}
}

func TestPropertiesZeroModTime(t *testing.T) {
	// Save original state
	originalPropsSources := propsSources
	defer func() { propsSources = originalPropsSources }()

	// Reset to default
	propsSources = []func(Page) []Property{defaultProps}

	testP := &testPropertyPage{
		name:    "test-page",
		modTime: time.Time{}, // Zero time
	}

	props := Properties(testP)

	if len(props) != 0 {
		t.Errorf("Expected 0 properties for zero ModTime, got %d", len(props))
	}
}

func TestPropertiesMultipleSources(t *testing.T) {
	// Save original state
	originalPropsSources := propsSources
	defer func() { propsSources = originalPropsSources }()

	// Reset to default
	propsSources = []func(Page) []Property{defaultProps}

	// Register multiple property sources
	RegisterProperty(func(p Page) []Property {
		return []Property{
			testProperty{icon: "fa-words", name: "words", value: 500},
		}
	})

	RegisterProperty(func(p Page) []Property {
		return []Property{
			testProperty{icon: "fa-time", name: "reading-time", value: "5 min"},
		}
	})

	testTime := time.Now().Add(-1 * time.Hour)
	testP := &testPropertyPage{
		name:    "test-page",
		modTime: testTime,
	}

	props := Properties(testP)

	// Should have default 'modified' + 2 custom properties
	if len(props) != 3 {
		t.Errorf("Expected 3 properties, got %d", len(props))
	}

	// Check each property exists
	expectedProps := []string{"modified", "words", "reading-time"}
	for _, name := range expectedProps {
		if _, exists := props[name]; !exists {
			t.Errorf("Expected property '%s' to exist", name)
		}
	}
}

func TestPropertiesOverwriteDuplicateNames(t *testing.T) {
	// Save original state
	originalPropsSources := propsSources
	defer func() { propsSources = originalPropsSources }()

	// Reset to default
	propsSources = []func(Page) []Property{defaultProps}

	// Register two sources with same property name
	RegisterProperty(func(p Page) []Property {
		return []Property{
			testProperty{icon: "fa-first", name: "duplicate", value: "first"},
		}
	})

	RegisterProperty(func(p Page) []Property {
		return []Property{
			testProperty{icon: "fa-second", name: "duplicate", value: "second"},
		}
	})

	testTime := time.Now().Add(-1 * time.Hour)
	testP := &testPropertyPage{
		name:    "test-page",
		modTime: testTime,
	}

	props := Properties(testP)

	duplicateProp, exists := props["duplicate"]
	if !exists {
		t.Fatal("Expected 'duplicate' property to exist")
	}

	// Second registration should overwrite first
	if duplicateProp.Value() != "second" {
		t.Errorf("Expected last registered property to win, got value: %v", duplicateProp.Value())
	}

	if duplicateProp.Icon() != "fa-second" {
		t.Errorf("Expected icon from last registration, got: %s", duplicateProp.Icon())
	}
}

func TestPropertiesEmptySource(t *testing.T) {
	// Save original state
	originalPropsSources := propsSources
	defer func() { propsSources = originalPropsSources }()

	// Reset to default
	propsSources = []func(Page) []Property{defaultProps}

	// Register a source that returns empty slice
	RegisterProperty(func(p Page) []Property {
		return []Property{}
	})

	// Register a source that returns nil
	RegisterProperty(func(p Page) []Property {
		return nil
	})

	testTime := time.Now().Add(-1 * time.Hour)
	testP := &testPropertyPage{
		name:    "test-page",
		modTime: testTime,
	}

	props := Properties(testP)

	// Should only have the default 'modified' property
	if len(props) != 1 {
		t.Errorf("Expected 1 property (only default), got %d", len(props))
	}

	if _, exists := props["modified"]; !exists {
		t.Error("Expected default 'modified' property to exist")
	}
}

func TestLastUpdatePropInterface(t *testing.T) {
	testTime := time.Now().Add(-3 * time.Hour)
	testP := &testPropertyPage{
		name:    "test-page",
		modTime: testTime,
	}

	prop := lastUpdateProp{page: testP}

	if prop.Icon() != "fa-solid fa-clock" {
		t.Errorf("Expected icon 'fa-solid fa-clock', got %s", prop.Icon())
	}

	if prop.Name() != "modified" {
		t.Errorf("Expected name 'modified', got %s", prop.Name())
	}

	value := prop.Value()
	if value == nil {
		t.Error("Expected non-nil value")
	}

	// Value should be string from ago() function
	if _, ok := value.(string); !ok {
		t.Errorf("Expected Value() to return string, got %T", value)
	}
}

func TestPropertiesReturnsMapNotSlice(t *testing.T) {
	// Save original state
	originalPropsSources := propsSources
	defer func() { propsSources = originalPropsSources }()

	// Reset to default
	propsSources = []func(Page) []Property{defaultProps}

	testTime := time.Now().Add(-1 * time.Hour)
	testP := &testPropertyPage{
		name:    "test-page",
		modTime: testTime,
	}

	props := Properties(testP)

	// Verify it's actually a map
	if props == nil {
		t.Fatal("Expected non-nil map")
	}

	// Maps allow direct key access
	_, exists := props["modified"]
	if !exists {
		t.Error("Expected map to support key-based access")
	}
}
