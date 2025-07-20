package xlog

import (
	"html/template"
	"io/fs"
	"net/http"
	"regexp"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// newTestApp creates a new App instance for testing
func newTestApp() *App {
	return &App{
		config:                &Config,
		router:                http.NewServeMux(),
		widgets:               make(map[WidgetSpace]*priorityList[WidgetFunc]),
		pageEvents:            make(map[PageEvent][]PageEventHandler),
		ignoredPaths:          []*regexp.Regexp{regexp.MustCompile(`^\.`)},
		concurrency:           runtime.NumCPU() * 4,
		propsSources:          []func(Page) []Property{DefaultProps},
		sources:               []PageSource{newMarkdownFS(".")},
		preprocessors:         []Preprocessor{},
		helpers:               template.FuncMap{},
		js:                    []string{},
		extensionPage:         make(map[string]bool),
		extensionPageEnclosed: make(map[string]bool),
		buildPerms:            0744,
		staticDirs:            []fs.FS{assets},
	}
}

// TestAppInitialization tests that the global app is properly initialized
func TestAppInitialization(t *testing.T) {
	app := GetApp()
	require.NotNil(t, app, "GetApp() returned nil")

	// Test that default values are set correctly
	require.NotNil(t, app.config, "config should not be nil")
	require.NotNil(t, app.router, "router should not be nil")
	require.NotNil(t, app.widgets, "widgets should not be nil")
	require.NotNil(t, app.pageEvents, "pageEvents should not be nil")
	require.NotNil(t, app.ignoredPaths, "ignoredPaths should not be nil")
	require.NotZero(t, app.concurrency, "concurrency should not be zero")
	require.NotNil(t, app.propsSources, "propsSources should not be nil")
	require.NotNil(t, app.sources, "sources should not be nil")
	require.NotNil(t, app.preprocessors, "preprocessors should not be nil")
	require.NotNil(t, app.helpers, "helpers should not be nil")
	require.NotNil(t, app.js, "js should not be nil")
	require.NotNil(t, app.extensionPage, "extensionPage should not be nil")
	require.NotNil(t, app.extensionPageEnclosed, "extensionPageEnclosed should not be nil")
	require.NotZero(t, app.buildPerms, "buildPerms should not be zero")
	require.NotNil(t, app.staticDirs, "staticDirs should not be nil")
}