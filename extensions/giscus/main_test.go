package giscus

import (
	"flag"
	"html/template"
	"strings"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
)

type mockPage struct {
	name string
}

func (m mockPage) Name() string             { return m.name }
func (m mockPage) FileName() string         { return m.name + ".md" }
func (m mockPage) Exists() bool             { return true }
func (m mockPage) Render() template.HTML    { return "" }
func (m mockPage) Content() xlog.Markdown   { return xlog.Markdown("") }
func (m mockPage) Delete() bool             { return false }
func (m mockPage) Write(xlog.Markdown) bool { return false }
func (m mockPage) ModTime() time.Time       { return time.Now() }
func (m mockPage) AST() ([]byte, ast.Node)  { return []byte{}, nil }

func TestGiscusExtensionName(t *testing.T) {
	ext := Giscus{}
	if ext.Name() != "giscus" {
		t.Errorf("expected extension name 'giscus', got '%s'", ext.Name())
	}
}

func TestGiscusWidget_MissingConfiguration(t *testing.T) {
	tests := []struct {
		name          string
		repoVal       string
		repoIDVal     string
		categoryVal   string
		categoryIDVal string
		shouldBeEmpty bool
		description   string
	}{
		{
			name:          "all empty",
			repoVal:       "",
			repoIDVal:     "",
			categoryVal:   "",
			categoryIDVal: "",
			shouldBeEmpty: true,
			description:   "widget should be empty when all values are empty",
		},
		{
			name:          "only repo set",
			repoVal:       "owner/repo",
			repoIDVal:     "",
			categoryVal:   "",
			categoryIDVal: "",
			shouldBeEmpty: true,
			description:   "widget should be empty when only repo is set",
		},
		{
			name:          "missing category ID",
			repoVal:       "owner/repo",
			repoIDVal:     "R_123",
			categoryVal:   "Announcements",
			categoryIDVal: "",
			shouldBeEmpty: true,
			description:   "widget should be empty when category ID is missing",
		},
		{
			name:          "all required fields set",
			repoVal:       "owner/repo",
			repoIDVal:     "R_123",
			categoryVal:   "Announcements",
			categoryIDVal: "DIC_456",
			shouldBeEmpty: false,
			description:   "widget should render when all required fields are set",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Save original values
			origRepo := repo
			origRepoID := repoID
			origCategory := category
			origCategoryID := categoryID
			defer func() {
				repo = origRepo
				repoID = origRepoID
				category = origCategory
				categoryID = origCategoryID
			}()

			repo = tc.repoVal
			repoID = tc.repoIDVal
			category = tc.categoryVal
			categoryID = tc.categoryIDVal

			page := mockPage{name: "test-page"}
			result := widget(page)

			if tc.shouldBeEmpty && result != "" {
				t.Errorf("%s, got: %s", tc.description, result)
			}

			if !tc.shouldBeEmpty && result == "" {
				t.Errorf("%s, but got empty", tc.description)
			}
		})
	}
}

func TestGiscusWidget_WithValidConfiguration(t *testing.T) {
	// Save original values
	origRepo := repo
	origRepoID := repoID
	origCategory := category
	origCategoryID := categoryID
	origMapping := mapping
	origTheme := theme
	origLang := lang
	defer func() {
		repo = origRepo
		repoID = origRepoID
		category = origCategory
		categoryID = origCategoryID
		mapping = origMapping
		theme = origTheme
		lang = origLang
	}()

	repo = "emad-elsaid/xlog"
	repoID = "R_kgDOG1234"
	category = "Announcements"
	categoryID = "DIC_kwDOG5678"
	mapping = "pathname"
	theme = "preferred_color_scheme"
	lang = "en"

	page := mockPage{name: "test-page"}
	result := string(widget(page))

	// Check that result contains expected elements
	expectedElements := []string{
		"giscus.app/client.js",
		"data-repo=\"emad-elsaid/xlog\"",
		"data-repo-id=\"R_kgDOG1234\"",
		"data-category=\"Announcements\"",
		"data-category-id=\"DIC_kwDOG5678\"",
		"data-mapping=\"pathname\"",
		"data-theme=\"preferred_color_scheme\"",
		"data-lang=\"en\"",
		"crossorigin=\"anonymous\"",
		"async",
	}

	for _, elem := range expectedElements {
		if !strings.Contains(result, elem) {
			t.Errorf("widget output should contain '%s'", elem)
		}
	}
}

func TestGiscusWidget_EscapesValues(t *testing.T) {
	// Save original values
	origRepo := repo
	origRepoID := repoID
	origCategory := category
	origCategoryID := categoryID
	defer func() {
		repo = origRepo
		repoID = origRepoID
		category = origCategory
		categoryID = origCategoryID
	}()

	repo = "owner/repo<script>alert('xss')</script>"
	repoID = "R_123"
	category = "Cat<script>"
	categoryID = "DIC_456"

	page := mockPage{name: "test-page"}
	result := string(widget(page))

	// Should not contain raw script tags
	if strings.Contains(result, "<script>alert") {
		t.Error("widget output should escape HTML in configuration values")
	}

	// Should contain escaped version
	if !strings.Contains(result, "&lt;script&gt;") {
		t.Error("widget output should contain HTML-escaped configuration values")
	}
}

func TestGiscusFlagRegistration(t *testing.T) {
	tests := []struct {
		flagName      string
		expectedUsage string
	}{
		{
			flagName:      "giscus-repo",
			expectedUsage: "GitHub repository in format 'owner/repo' (e.g., 'emad-elsaid/xlog')",
		},
		{
			flagName:      "giscus-repo-id",
			expectedUsage: "GitHub repository ID (get from giscus.app)",
		},
		{
			flagName:      "giscus-category",
			expectedUsage: "GitHub Discussions category name",
		},
		{
			flagName:      "giscus-category-id",
			expectedUsage: "GitHub Discussions category ID (get from giscus.app)",
		},
		{
			flagName:      "giscus-mapping",
			expectedUsage: "Page-discussion mapping method (pathname, url, title, og:title)",
		},
		{
			flagName:      "giscus-theme",
			expectedUsage: "Giscus theme (e.g., light, dark, preferred_color_scheme)",
		},
		{
			flagName:      "giscus-lang",
			expectedUsage: "Language code for Giscus interface",
		},
	}

	for _, tc := range tests {
		t.Run(tc.flagName, func(t *testing.T) {
			f := flag.Lookup(tc.flagName)
			if f == nil {
				t.Fatalf("%s flag should be registered", tc.flagName)
			}

			if f.Usage != tc.expectedUsage {
				t.Errorf("unexpected flag usage for %s: %s", tc.flagName, f.Usage)
			}
		})
	}
}

func TestGiscusWidget_DifferentMappings(t *testing.T) {
	// Save original values
	origRepo := repo
	origRepoID := repoID
	origCategory := category
	origCategoryID := categoryID
	origMapping := mapping
	defer func() {
		repo = origRepo
		repoID = origRepoID
		category = origCategory
		categoryID = origCategoryID
		mapping = origMapping
	}()

	repo = "owner/repo"
	repoID = "R_123"
	category = "General"
	categoryID = "DIC_456"

	mappings := []string{"pathname", "url", "title", "og:title"}

	for _, m := range mappings {
		t.Run(m, func(t *testing.T) {
			mapping = m
			page := mockPage{name: "test-page"}
			result := string(widget(page))

			expected := "data-mapping=\"" + m + "\""
			if !strings.Contains(result, expected) {
				t.Errorf("widget should contain mapping '%s'", m)
			}
		})
	}
}
