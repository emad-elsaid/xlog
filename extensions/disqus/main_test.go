package disqus

import (
	"flag"
	"html/template"
	"strings"
	"testing"
	"time"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
)

type mockPage struct {
	name string
}

func (m mockPage) Name() string                       { return m.name }
func (m mockPage) FileName() string                   { return m.name + ".md" }
func (m mockPage) Exists() bool                       { return true }
func (m mockPage) Render() template.HTML              { return "" }
func (m mockPage) Content() Markdown                  { return Markdown("") }
func (m mockPage) Delete() bool                       { return false }
func (m mockPage) Write(Markdown) bool                { return false }
func (m mockPage) ModTime() time.Time                 { return time.Now() }
func (m mockPage) AST() ([]byte, ast.Node)            { return []byte{}, nil }

func TestDisqusExtensionName(t *testing.T) {
	ext := Disqus{}
	if ext.Name() != "disqus" {
		t.Errorf("expected extension name 'disqus', got '%s'", ext.Name())
	}
}

func TestDisqusWidget_EmptyDomain(t *testing.T) {
	// Save original domain
	originalDomain := domain
	defer func() { domain = originalDomain }()

	domain = ""
	page := mockPage{name: "test-page"}

	result := widget(page)
	if result != "" {
		t.Errorf("expected empty widget when domain is empty, got: %s", result)
	}
}

func TestDisqusWidget_WithDomain(t *testing.T) {
	// Save original domain
	originalDomain := domain
	defer func() { domain = originalDomain }()

	domain = "xlog-test.disqus.com"
	page := mockPage{name: "test-page"}

	result := string(widget(page))

	// Check that result contains expected elements
	if !strings.Contains(result, "disqus_thread") {
		t.Error("widget output should contain 'disqus_thread' div")
	}

	if !strings.Contains(result, domain) {
		t.Errorf("widget output should contain domain '%s'", domain)
	}

	if !strings.Contains(result, "test-page") {
		t.Error("widget output should contain page identifier")
	}

	if !strings.Contains(result, "embed.js") {
		t.Error("widget output should contain embed.js script")
	}
}

func TestDisqusWidget_EscapesPageName(t *testing.T) {
	// Save original domain
	originalDomain := domain
	defer func() { domain = originalDomain }()

	domain = "xlog-test.disqus.com"
	page := mockPage{name: "test<script>alert('xss')</script>"}

	result := string(widget(page))

	// Page name should be JS-escaped
	if strings.Contains(result, "<script>") && !strings.Contains(result, "\\u003Cscript\\u003E") {
		t.Error("widget output should escape page name to prevent XSS")
	}

	// Should contain escaped version (\u003C instead of <)
	if !strings.Contains(result, "\\u003C") {
		t.Error("widget output should contain JS-escaped page name")
	}
}

func TestDisqusFlagRegistration(t *testing.T) {
	// Verify the flag was registered
	f := flag.Lookup("disqus")
	if f == nil {
		t.Fatal("disqus flag should be registered")
	}

	if f.Usage != "Disqus domain name for example: xlog-emadelsaid.disqus.com" {
		t.Errorf("unexpected flag usage: %s", f.Usage)
	}
}
