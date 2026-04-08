package recent

import (
	"html/template"
	"testing"
)

func TestRecentExtensionName(t *testing.T) {
	ext := Recent{}

	if ext.Name() != "recent" {
		t.Errorf("Expected name 'recent', got '%s'", ext.Name())
	}
}

func TestLinksIcon(t *testing.T) {
	l := links{}

	expectedIcon := "fa-solid fa-clock-rotate-left"
	if l.Icon() != expectedIcon {
		t.Errorf("Expected icon '%s', got '%s'", expectedIcon, l.Icon())
	}
}

func TestLinksName(t *testing.T) {
	l := links{}

	expectedName := "Recent"
	if l.Name() != expectedName {
		t.Errorf("Expected name '%s', got '%s'", expectedName, l.Name())
	}
}

func TestLinksAttrs(t *testing.T) {
	l := links{}

	attrs := l.Attrs()
	if attrs == nil {
		t.Fatal("Expected attrs map, got nil")
	}

	href, ok := attrs["href"]
	if !ok {
		t.Error("Expected 'href' attribute in attrs map")
	}

	expectedHref := "/+/recent"
	if href != expectedHref {
		t.Errorf("Expected href '%s', got '%v'", expectedHref, href)
	}
}

func TestLinksAttrsType(t *testing.T) {
	l := links{}

	attrs := l.Attrs()

	// Verify the map key type is template.HTMLAttr
	for key := range attrs {
		_, ok := interface{}(key).(template.HTMLAttr)
		if !ok {
			t.Errorf("Expected key type template.HTMLAttr, got %T", key)
		}
		break // Just check first key
	}
}

func TestLinksAttrsHrefValue(t *testing.T) {
	l := links{}

	attrs := l.Attrs()
	
	// Check that href points to the recent page endpoint
	href := attrs["href"]
	if href == nil {
		t.Fatal("href attribute is nil")
	}

	hrefStr, ok := href.(string)
	if !ok {
		t.Fatalf("Expected href to be string, got %T", href)
	}

	if hrefStr != "/+/recent" {
		t.Errorf("Expected href '/+/recent', got '%s'", hrefStr)
	}
}
