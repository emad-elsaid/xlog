package xlog

import (
	"html/template"
	"io/fs"
	"net/http"
	"regexp"
	"runtime"
	"testing"
	"time"

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

// TestCommandsBehavior tests that Commands function behavior is preserved
func TestCommandsBehavior(t *testing.T) {
	app := newTestApp()
	app.commands = []func(Page) []Command{}

	page := &page{name: "test"}
	commands := app.Commands(page)
	require.Len(t, commands, 0, "Expected 0 commands")

	testCommand := func(p Page) []Command {
		return []Command{
			&testCommandImpl{
				icon:  "fa-test",
				name:  "Test Command",
				attrs: map[template.HTMLAttr]any{"href": "/test"},
			},
		}
	}

	app.RegisterCommand(testCommand)
	commands = app.Commands(page)
	require.Len(t, commands, 1, "Expected 1 command")
	require.Equal(t, "Test Command", commands[0].Name(), "Expected command name 'Test Command'")
}

// TestQuickCommandsBehavior tests that QuickCommands function behavior is preserved
func TestQuickCommandsBehavior(t *testing.T) {
	app := newTestApp()
	app.quickCommands = []func(Page) []Command{}

	page := &page{name: "test"}
	commands := app.QuickCommands(page)
	require.Len(t, commands, 0, "Expected 0 quick commands")

	testCommand := func(p Page) []Command {
		return []Command{
			&testCommandImpl{
				icon:  "fa-quick",
				name:  "Quick Test",
				attrs: map[template.HTMLAttr]any{"href": "/quick"},
			},
		}
	}

	app.RegisterQuickCommand(testCommand)
	commands = app.QuickCommands(page)
	require.Len(t, commands, 1, "Expected 1 quick command")
	require.Equal(t, "Quick Test", commands[0].Name(), "Expected quick command name 'Quick Test'")
}

// TestLinksBehavior tests that Links function behavior is preserved
func TestLinksBehavior(t *testing.T) {
	app := newTestApp()
	app.links = []func(Page) []Command{}

	page := &page{name: "test"}
	links := app.Links(page)
	require.Len(t, links, 0, "Expected 0 links")

	testLink := func(p Page) []Command {
		return []Command{
			&testCommandImpl{
				icon:  "fa-link",
				name:  "Test Link",
				attrs: map[template.HTMLAttr]any{"href": "/link"},
			},
		}
	}

	app.RegisterLink(testLink)
	links = app.Links(page)
	require.Len(t, links, 1, "Expected 1 link")
	require.Equal(t, "Test Link", links[0].Name(), "Expected link name 'Test Link'")
}

// TestIsIgnoredPathBehavior tests that IsIgnoredPath function behavior is preserved
func TestIsIgnoredPathBehavior(t *testing.T) {
	app := newTestApp()

	require.True(t, app.IsIgnoredPath(".git"), "Expected .git to be ignored")
	require.True(t, app.IsIgnoredPath(".hidden"), "Expected .hidden to be ignored")
	require.False(t, app.IsIgnoredPath("normal"), "Expected 'normal' to not be ignored")
	require.False(t, app.IsIgnoredPath("normal/path"), "Expected 'normal/path' to not be ignored")

	customPattern := regexp.MustCompile(`^temp`)
	app.IgnorePath(customPattern)

	require.True(t, app.IsIgnoredPath("temp"), "Expected 'temp' to be ignored after registering custom pattern")
	require.True(t, app.IsIgnoredPath("tempfile.txt"), "Expected 'tempfile.txt' to be ignored after registering custom pattern")
	require.False(t, app.IsIgnoredPath("permanent"), "Expected 'permanent' to not be ignored")
}

// TestAgoBehavior tests that ago function behavior is preserved
func TestAgoBehavior(t *testing.T) {
	app := newTestApp()

	app.config.Readonly = true
	now := time.Now()
	result := app.ago(now)
	require.Equal(t, now.Format("Monday 2 January 2006"), result)

	app.config.Readonly = false

	recent := time.Now().Add(-500 * time.Millisecond)
	result = app.ago(recent)
	require.Contains(t, result, "Less than a second")

	oneMinuteAgo := time.Now().Add(-1 * time.Minute)
	result = app.ago(oneMinuteAgo)
	require.Contains(t, result, "1 minute")

	oneHourAgo := time.Now().Add(-1 * time.Hour)
	result = app.ago(oneHourAgo)
	require.Contains(t, result, "1 hour")

	oneDayAgo := time.Now().Add(-24 * time.Hour)
	result = app.ago(oneDayAgo)
	require.Contains(t, result, "1 day")
}

// TestPropertiesBehavior tests that Properties function behavior is preserved
func TestPropertiesBehavior(t *testing.T) {
	app := newTestApp()

	page := &page{name: "test"}
	props := app.Properties(page)
	require.Len(t, props, 0, "Expected 0 properties for page with zero mod time")

	// Note: Testing with actual file system operations would require more complex setup
	// This is a basic test to ensure the function doesn't panic
}

// TestRenderWidgetBehavior tests that RenderWidget function behavior is preserved
func TestRenderWidgetBehavior(t *testing.T) {
	app := newTestApp()
	app.widgets = make(map[WidgetSpace]*priorityList[WidgetFunc])

	page := &page{name: "test"}
	result := app.RenderWidget(WidgetAfterView, page)
	require.Equal(t, template.HTML(""), result, "Expected empty result for no widgets")

	testWidget := func(p Page) template.HTML {
		return template.HTML("<div>Test Widget</div>")
	}

	app.RegisterWidget(WidgetAfterView, 1.0, testWidget)
	result = app.RenderWidget(WidgetAfterView, page)
	require.Equal(t, template.HTML("<div>Test Widget</div>"), result)

	secondWidget := func(p Page) template.HTML {
		return template.HTML("<div>Second Widget</div>")
	}

	app.RegisterWidget(WidgetAfterView, 0.5, secondWidget)
	result = app.RenderWidget(WidgetAfterView, page)
	require.Equal(t, template.HTML("<div>Second Widget</div><div>Test Widget</div>"), result)
}

// TestPreProcessBehavior tests that PreProcess function behavior is preserved
func TestPreProcessBehavior(t *testing.T) {
	app := newTestApp()
	app.preprocessors = []Preprocessor{}

	content := Markdown("test content")
	result := app.PreProcess(content)
	require.Equal(t, content, result, "Expected unchanged content")

	preprocessor := func(content Markdown) Markdown {
		return Markdown("processed: " + string(content))
	}

	app.RegisterPreprocessor(preprocessor)
	result = app.PreProcess(content)
	require.Equal(t, Markdown("processed: test content"), result)

	secondPreprocessor := func(content Markdown) Markdown {
		return Markdown(string(content) + " (final)")
	}

	app.RegisterPreprocessor(secondPreprocessor)
	result = app.PreProcess(content)
	require.Equal(t, Markdown("processed: test content (final)"), result)
}

// TestEventHandling tests that event handling behavior is preserved
func TestEventHandling(t *testing.T) {
	app := newTestApp()
	app.pageEvents = make(map[PageEvent][]PageEventHandler)

	triggered := false
	handler := func(p Page) error {
		triggered = true
		return nil
	}

	app.Listen(PageChanged, handler)
	handlers, exists := app.pageEvents[PageChanged]
	require.True(t, exists, "Expected PageChanged event to be registered")
	require.Len(t, handlers, 1, "Expected 1 handler")

	page := &page{name: "test"}
	app.Trigger(PageChanged, page)
	require.True(t, triggered, "Expected handler to be triggered")
}

// TestHelperRegistration tests that helper registration behavior is preserved
func TestHelperRegistration(t *testing.T) {
	app := newTestApp()
	testHelper := func(s string) string {
		return "test: " + s
	}

	err := app.RegisterHelper("testHelper", testHelper)
	require.NoError(t, err, "Expected no error")
	require.NotNil(t, app.helpers["testHelper"], "Expected helper to be registered")

	err = app.RegisterHelper("testHelper", testHelper)
	require.ErrorIs(t, err, ErrHelperRegistered, "Expected ErrHelperRegistered")
}

// TestJavaScriptHandling tests that JavaScript handling behavior is preserved
func TestJavaScriptHandling(t *testing.T) {
	app := newTestApp()
	app.js = []string{}

	result := app.includeJS("/test.js")
	require.Equal(t, template.HTML(""), result, "Expected empty result")
	require.Len(t, app.js, 1, "Expected 1 JS file")
	require.Equal(t, "/test.js", app.js[0], "Expected '/test.js'")

	app.includeJS("/test.js")
	require.Len(t, app.js, 1, "Expected 1 JS file after duplicate")

	app.includeJS("/another.js")
	scripts := app.scripts()
	require.Equal(t, `<script src="/test.js" defer></script><script src="/another.js" defer></script>`, string(scripts))
}

// TestIsFontAwesome tests that IsFontAwesome function behavior is preserved
func TestIsFontAwesome(t *testing.T) {
	app := newTestApp()

	require.True(t, app.IsFontAwesome("fa-solid"), "Expected 'fa-solid' to be FontAwesome")
	require.True(t, app.IsFontAwesome("fa-regular"), "Expected 'fa-regular' to be FontAwesome")
	require.True(t, app.IsFontAwesome("fa-brands"), "Expected 'fa-brands' to be FontAwesome")
	require.False(t, app.IsFontAwesome("not-fa"), "Expected 'not-fa' to not be FontAwesome")
	require.False(t, app.IsFontAwesome(""), "Expected empty string to not be FontAwesome")
}

// TestBannerApp tests that Banner function behavior is preserved
func TestBannerApp(t *testing.T) {
	app := newTestApp()
	testPage := &page{name: "test"}
	result := app.Banner(testPage)
	require.Equal(t, "", result, "Expected empty banner for page with no AST")
	// Note: Testing with actual AST would require more complex setup
	// This is a basic test to ensure the function doesn't panic
}

// TestEmoji tests that Emoji function behavior is preserved
func TestEmoji(t *testing.T) {
	app := newTestApp()
	testPage := &page{name: "test"}
	result := app.Emoji(testPage)
	require.Equal(t, "", result, "Expected empty emoji for page with no AST")
	// Note: Testing with actual AST would require more complex setup
	// This is a basic test to ensure the function doesn't panic
}

// TestDir tests that dir function behavior is preserved
func TestDir(t *testing.T) {
	app := newTestApp()
	require.Equal(t, "", app.dir(""), "Expected empty string for empty path")
	require.Equal(t, "", app.dir("."), "Expected empty string for '.'")
	require.Equal(t, "", app.dir("file.txt"), "Expected empty string for 'file.txt'")
	require.Equal(t, "dir", app.dir("dir/file.txt"), "Expected 'dir' for 'dir/file.txt'")
	require.Equal(t, "a/b", app.dir("a/b/c.txt"), "Expected 'a/b' for 'a/b/c.txt'")
}

// TestRaw tests that raw function behavior is preserved
func TestRaw(t *testing.T) {
	app := newTestApp()
	input := "<div>test</div>"
	result := app.raw(input)
	require.Equal(t, template.HTML(input), result, "Expected HTML to match input")
}

// Test command implementation for testing
type testCommandImpl struct {
	icon  string
	name  string
	attrs map[template.HTMLAttr]any
}

func (c *testCommandImpl) Icon() string {
	return c.icon
}

func (c *testCommandImpl) Name() string {
	return c.name
}

func (c *testCommandImpl) Attrs() map[template.HTMLAttr]any {
	return c.attrs
}
