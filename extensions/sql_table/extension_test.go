package sql_table

import (
	"flag"
	"html/template"
	"strings"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
	"github.com/emad-elsaid/xlog/markdown/text"
)

func TestExtensionName(t *testing.T) {
	ext := Extension{}
	got := ext.Name()
	want := "sql_table"

	if got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestScript_NilPage(t *testing.T) {
	result := script(nil)

	if result != "" {
		t.Errorf("script(nil) = %q, want empty string", result)
	}
}

func TestScript_NoAST(t *testing.T) {
	// Create a mock page with nil AST.
	mockPage := &mockPageNoAST{}

	result := script(mockPage)

	if result != "" {
		t.Errorf("script with nil AST = %q, want empty string", result)
	}
}

func TestScript_NoTables(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "Empty document",
			content: "",
		},
		{
			name:    "Only paragraphs",
			content: "Just some text without any tables.",
		},
		{
			name: "Headings and lists",
			content: `# Heading

- List item 1
- List item 2`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := xlog.MarkdownConverter()
			doc := md.Parser().Parse(text.NewReader([]byte(tt.content)))

			mockPage := &mockPageWithAST{astNode: doc}
			result := script(mockPage)

			if result != "" {
				t.Errorf("script with no tables = %q, want empty string", result)
			}
		})
	}
}

func TestScript_SmallTable(t *testing.T) {
	// Create table with fewer rows than threshold.
	content := `| Header 1 | Header 2 |
|----------|----------|
| Row 1    | Data 1   |
| Row 2    | Data 2   |
| Row 3    | Data 3   |`

	// Set threshold higher than the table rows.
	oldThreshold := sqlTableThreshold
	sqlTableThreshold = 10
	defer func() { sqlTableThreshold = oldThreshold }()

	md := xlog.MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(content)))

	mockPage := &mockPageWithAST{astNode: doc}
	result := script(mockPage)

	if result != "" {
		t.Errorf("script with small table = %q, want empty string", result)
	}
}

func TestScript_LargeTable(t *testing.T) {
	// Generate a table with many rows (exceeding threshold).
	var sb strings.Builder
	sb.WriteString("| Col1 | Col2 |\n")
	sb.WriteString("|------|------|\n")
	for i := 1; i <= 150; i++ {
		sb.WriteString("| Data | More |\n")
	}

	// Set threshold lower than the table rows.
	oldThreshold := sqlTableThreshold
	sqlTableThreshold = 100
	defer func() { sqlTableThreshold = oldThreshold }()

	md := xlog.MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(sb.String())))

	mockPage := &mockPageWithAST{astNode: doc}
	result := script(mockPage)

	// Should return script content.
	if result == "" {
		t.Error("script with large table returned empty, want script content")
	}

	// Verify it contains expected elements.
	resultStr := string(result)
	if !strings.Contains(resultStr, "sql.js") {
		t.Error("script result missing sql.js reference")
	}
	if !strings.Contains(resultStr, "sqlTableThreshold") {
		t.Error("script result missing sqlTableThreshold variable")
	}
}

func TestScript_MixedTables(t *testing.T) {
	// Create content with both small and large tables.
	var sb strings.Builder
	// Small table.
	sb.WriteString("| A | B |\n|---|---|\n| 1 | 2 |\n\n")
	// Large table.
	sb.WriteString("| Col1 | Col2 |\n|------|------|\n")
	for i := 1; i <= 150; i++ {
		sb.WriteString("| Data | More |\n")
	}

	oldThreshold := sqlTableThreshold
	sqlTableThreshold = 100
	defer func() { sqlTableThreshold = oldThreshold }()

	md := xlog.MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(sb.String())))

	mockPage := &mockPageWithAST{astNode: doc}
	result := script(mockPage)

	// Should return script because at least one table exceeds threshold.
	if result == "" {
		t.Error("script with mixed tables returned empty, want script content")
	}
}

func TestScript_ThresholdEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		rowCount    int
		threshold   int
		wantScript  bool
		description string
	}{
		{
			name:        "Exactly at threshold",
			rowCount:    100,
			threshold:   100,
			wantScript:  true,
			description: "Table with exactly threshold rows should trigger script",
		},
		{
			name:        "One below threshold",
			rowCount:    99,
			threshold:   100,
			wantScript:  false,
			description: "Table with one less than threshold should not trigger",
		},
		{
			name:        "One above threshold",
			rowCount:    101,
			threshold:   100,
			wantScript:  true,
			description: "Table with one more than threshold should trigger script",
		},
		{
			name:        "Zero threshold",
			rowCount:    1,
			threshold:   0,
			wantScript:  true,
			description: "Any table should trigger with zero threshold",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			sb.WriteString("| Col1 | Col2 |\n|------|------|\n")
			for i := 0; i < tt.rowCount; i++ {
				sb.WriteString("| Data | More |\n")
			}

			oldThreshold := sqlTableThreshold
			sqlTableThreshold = tt.threshold
			defer func() { sqlTableThreshold = oldThreshold }()

			md := xlog.MarkdownConverter()
			doc := md.Parser().Parse(text.NewReader([]byte(sb.String())))

			mockPage := &mockPageWithAST{astNode: doc}
			result := script(mockPage)

			gotScript := result != ""
			if gotScript != tt.wantScript {
				t.Errorf("%s: got script=%v, want %v", tt.description, gotScript, tt.wantScript)
			}
		})
	}
}

func TestScript_EmbeddedContentIntegrity(t *testing.T) {
	// Generate large table.
	var sb strings.Builder
	sb.WriteString("| Col1 | Col2 |\n|------|------|\n")
	for i := 1; i <= 150; i++ {
		sb.WriteString("| Data | More |\n")
	}

	oldThreshold := sqlTableThreshold
	sqlTableThreshold = 100
	defer func() { sqlTableThreshold = oldThreshold }()

	md := xlog.MarkdownConverter()
	doc := md.Parser().Parse(text.NewReader([]byte(sb.String())))

	mockPage := &mockPageWithAST{astNode: doc}
	result := script(mockPage)

	resultStr := string(result)

	// Verify critical components of embedded JavaScript.
	requiredElements := []string{
		"sql-wasm.js",
		"query",
		"tableToJson",
		"createResultTable",
		"initializeTableQueries",
		"const sqlTableThreshold",
	}

	for _, elem := range requiredElements {
		if !strings.Contains(resultStr, elem) {
			t.Errorf("Embedded script missing expected element: %q", elem)
		}
	}
}

func TestFlagRegistration(t *testing.T) {
	// Verify the flag is registered.
	f := flag.Lookup("sql-table.threshold")
	if f == nil {
		t.Fatal("Flag 'sql-table.threshold' not registered")
	}

	if f.DefValue != "100" {
		t.Errorf("Flag default value = %q, want %q", f.DefValue, "100")
	}
}

// Mock implementations for testing.
type mockPageNoAST struct{}

func (m *mockPageNoAST) Name() string             { return "mock" }
func (m *mockPageNoAST) FileName() string         { return "mock.md" }
func (m *mockPageNoAST) Exists() bool             { return false }
func (m *mockPageNoAST) Render() template.HTML    { return "" }
func (m *mockPageNoAST) Content() xlog.Markdown   { return "" }
func (m *mockPageNoAST) Delete() bool             { return false }
func (m *mockPageNoAST) Write(xlog.Markdown) bool { return false }
func (m *mockPageNoAST) ModTime() time.Time       { return time.Time{} }
func (m *mockPageNoAST) AST() ([]byte, ast.Node)  { return nil, nil }

type mockPageWithAST struct {
	astNode ast.Node
}

func (m *mockPageWithAST) Name() string             { return "mock" }
func (m *mockPageWithAST) FileName() string         { return "mock.md" }
func (m *mockPageWithAST) Exists() bool             { return false }
func (m *mockPageWithAST) Render() template.HTML    { return "" }
func (m *mockPageWithAST) Content() xlog.Markdown   { return "" }
func (m *mockPageWithAST) Delete() bool             { return false }
func (m *mockPageWithAST) Write(xlog.Markdown) bool { return false }
func (m *mockPageWithAST) ModTime() time.Time       { return time.Time{} }
func (m *mockPageWithAST) AST() ([]byte, ast.Node)  { return []byte{}, m.astNode }
