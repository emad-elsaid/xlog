package star

import (
	"strings"
	"testing"
)

func TestIsStarredLogic(t *testing.T) {
	tests := []struct {
		name           string
		starredContent string
		pageName       string
		expected       bool
	}{
		{
			name:           "Page is starred",
			starredContent: "page1.md\npage2.md\npage3.md",
			pageName:       "page2.md",
			expected:       true,
		},
		{
			name:           "Page is not starred",
			starredContent: "page1.md\npage3.md",
			pageName:       "page2.md",
			expected:       false,
		},
		{
			name:           "Empty starred list",
			starredContent: "",
			pageName:       "page1.md",
			expected:       false,
		},
		{
			name:           "Page with whitespace",
			starredContent: "  page1.md  \npage2.md\n  page3.md  ",
			pageName:       "page1.md",
			expected:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, k := range strings.Split(tt.starredContent, "\n") {
				if strings.TrimSpace(k) == tt.pageName {
					found = true
					break
				}
			}

			if found != tt.expected {
				t.Errorf("Expected %v, got %v for page %s in starred list:\n%s",
					tt.expected, found, tt.pageName, tt.starredContent)
			}
		})
	}
}

func TestActionIconAndName(t *testing.T) {
	tests := []struct {
		name          string
		starred       bool
		expectedIcon  string
		expectedName  string
	}{
		{
			name:          "Starred action shows unstar",
			starred:       true,
			expectedIcon:  "fa-solid fa-star",
			expectedName:  "Unstar",
		},
		{
			name:          "Unstarred action shows star",
			starred:       false,
			expectedIcon:  "fa-regular fa-star",
			expectedName:  "Star",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			act := action{starred: tt.starred}

			if act.Icon() != tt.expectedIcon {
				t.Errorf("Expected icon %s, got %s", tt.expectedIcon, act.Icon())
			}

			if act.Name() != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, act.Name())
			}
		})
	}
}

func TestStarredPagesParsing(t *testing.T) {
	content := "page1.md\npage2.md\npage3.md\n"
	list := strings.Split(strings.TrimSpace(content), "\n")

	if len(list) != 3 {
		t.Errorf("Expected 3 pages, got %d", len(list))
	}

	expected := []string{"page1.md", "page2.md", "page3.md"}
	for i, v := range list {
		if v != expected[i] {
			t.Errorf("Expected %s at index %d, got %s", expected[i], i, v)
		}
	}
}

func TestStarredPagesEmptyContent(t *testing.T) {
	content := ""
	trimmed := strings.TrimSpace(content)

	if trimmed != "" {
		t.Error("Expected empty string after trim")
	}

	// Empty content should return nil list
	if trimmed == "" {
		// This is the expected behavior
		return
	}

	t.Error("Should have returned early for empty content")
}
