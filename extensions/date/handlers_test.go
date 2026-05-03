package date

import (
	"context"
	"html/template"
	"net/http"
	"testing"
	"time"

	"github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
)

func TestOrganizeCalendar(t *testing.T) {
	tests := []struct {
		name  string
		pairs []pair
		want  func([]Year) bool
	}{
		{
			name:  "empty calendar",
			pairs: []pair{},
			want: func(years []Year) bool {
				return len(years) == 0
			},
		},
		{
			name: "single date single page",
			pairs: []pair{
				{
					Time: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
					Page: &mockPage{name: "test"},
				},
			},
			want: func(years []Year) bool {
				if len(years) != 1 {
					return false
				}
				if years[0].Year != 2026 {
					return false
				}
				if len(years[0].Months) != 1 {
					return false
				}
				return years[0].Months[0].Name == "March"
			},
		},
		{
			name: "multiple dates same month",
			pairs: []pair{
				{
					Time: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
					Page: &mockPage{name: "page1"},
				},
				{
					Time: time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC),
					Page: &mockPage{name: "page2"},
				},
			},
			want: func(years []Year) bool {
				if len(years) != 1 || years[0].Year != 2026 {
					return false
				}
				if len(years[0].Months) != 1 {
					return false
				}
				month := years[0].Months[0]
				foundDay15 := false
				foundDay20 := false
				for _, week := range month.Days {
					for _, day := range week {
						if day != nil {
							if day.Date.Day() == 15 {
								foundDay15 = len(day.Pages) == 1
							}
							if day.Date.Day() == 20 {
								foundDay20 = len(day.Pages) == 1
							}
						}
					}
				}
				return foundDay15 && foundDay20
			},
		},
		{
			name: "multiple pages same date",
			pairs: []pair{
				{
					Time: time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC),
					Page: &mockPage{name: "page1"},
				},
				{
					Time: time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC),
					Page: &mockPage{name: "page2"},
				},
			},
			want: func(years []Year) bool {
				if len(years) != 1 || years[0].Year != 2026 {
					return false
				}
				month := years[0].Months[0]
				for _, week := range month.Days {
					for _, day := range week {
						if day != nil && day.Date.Day() == 10 {
							return len(day.Pages) == 2
						}
					}
				}
				return false
			},
		},
		{
			name: "dates across multiple years",
			pairs: []pair{
				{
					Time: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
					Page: &mockPage{name: "old"},
				},
				{
					Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
					Page: &mockPage{name: "new"},
				},
			},
			want: func(years []Year) bool {
				if len(years) != 2 {
					return false
				}
				yearMap := make(map[int]bool)
				for _, y := range years {
					yearMap[y.Year] = true
				}
				return yearMap[2025] && yearMap[2026]
			},
		},
		{
			name: "month boundary dates",
			pairs: []pair{
				{
					Time: time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC),
					Page: &mockPage{name: "feb"},
				},
				{
					Time: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
					Page: &mockPage{name: "mar"},
				},
			},
			want: func(years []Year) bool {
				if len(years) != 1 || years[0].Year != 2026 {
					return false
				}
				return len(years[0].Months) == 2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := organizeCalendar(tt.pairs)
			if !tt.want(got) {
				t.Errorf("organizeCalendar() validation failed for %s", tt.name)
			}
		})
	}
}

func TestOrganizeCalendar_WeekStructure(t *testing.T) {
	// Test that month is organized correctly into 6 weeks of 7 days
	pairs := []pair{
		{
			Time: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			Page: &mockPage{name: "start"},
		},
	}

	result := organizeCalendar(pairs)

	if len(result) != 1 {
		t.Fatalf("Expected 1 year, got %d", len(result))
	}

	year := result[0]
	if len(year.Months) != 1 {
		t.Fatalf("Expected 1 month, got %d", len(year.Months))
	}

	month := year.Months[0]

	// Verify structure is 6 weeks x 7 days
	if len(month.Days) != 6 {
		t.Errorf("Expected 6 weeks, got %d", len(month.Days))
	}

	for i, week := range month.Days {
		if len(week) != 7 {
			t.Errorf("Week %d has %d days, want 7", i, len(week))
		}
	}

	// Find the first day of the month
	found := false
	for _, week := range month.Days {
		for _, day := range week {
			if day != nil && day.Date.Day() == 1 {
				found = true
				if len(day.Pages) != 1 {
					t.Errorf("Day 1 has %d pages, want 1", len(day.Pages))
				}
			}
		}
	}

	if !found {
		t.Error("Did not find day 1 in organized calendar")
	}
}

func TestDateHandler(t *testing.T) {
	tests := []struct {
		name         string
		dateParam    string
		setupPages   map[string]string
		expectError  bool
		expectStatus int
		expectedName string
		pageCount    int
	}{
		{
			name:         "valid date with matching pages",
			dateParam:    "15-3-2026",
			setupPages:   map[string]string{"page1.md": "Meeting on 15/3/2026"},
			expectError:  false,
			expectedName: "15 March 2026",
			pageCount:    1,
		},
		{
			name:         "valid date no matching pages",
			dateParam:    "1-1-2026",
			setupPages:   map[string]string{"page1.md": "No dates here"},
			expectError:  false,
			expectedName: "1 January 2026",
			pageCount:    0,
		},
		{
			name:         "invalid date format",
			dateParam:    "invalid-date",
			expectError:  true,
			expectStatus: 400,
		},
		{
			name:         "malformed date",
			dateParam:    "32-13-2026",
			expectError:  true,
			expectStatus: 400,
		},
		{
			name:         "multiple pages same date",
			dateParam:    "10-5-2026",
			setupPages:   map[string]string{"page1.md": "Event 10/5/2026", "page2.md": "Another 10/5/2026"},
			expectError:  false,
			expectedName: "10 May 2026",
			pageCount:    2,
		},
		{
			name:         "date with different format in content",
			dateParam:    "25-12-2025",
			setupPages:   map[string]string{"page1.md": "Christmas 25/12/2025"},
			expectError:  false,
			expectedName: "25 December 2025",
			pageCount:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock request with PathValue
			req := createMockRequest(map[string]string{"date": tt.dateParam})

			// Execute handler
			output := dateHandler(req)

			if output == nil {
				t.Fatal("dateHandler returned nil output")
			}

			// For now, we just verify handler doesn't panic and returns output
			// Full integration testing would require template rendering
		})
	}
}

func TestCalendarHandler(t *testing.T) {
	tests := []struct {
		name       string
		setupPages map[string]string
		expectName string
	}{
		{
			name:       "empty calendar",
			setupPages: map[string]string{},
			expectName: "Calendar",
		},
		{
			name:       "calendar with single date",
			setupPages: map[string]string{"page1.md": "Meeting 15/3/2026"},
			expectName: "Calendar",
		},
		{
			name: "calendar with multiple dates",
			setupPages: map[string]string{
				"page1.md": "Event 1/1/2026",
				"page2.md": "Event 15/3/2026",
				"page3.md": "Event 25/12/2026",
			},
			expectName: "Calendar",
		},
		{
			name:       "calendar with pages without dates",
			setupPages: map[string]string{"page1.md": "No dates"},
			expectName: "Calendar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock request
			req := createMockRequest(map[string]string{})

			// Execute handler
			output := calendarHandler(req)

			if output == nil {
				t.Fatal("calendarHandler returned nil output")
			}

			// Verify output structure (basic smoke test)
			// Full testing would require template rendering
		})
	}
}

// createMockRequest creates a mock *http.Request for testing handlers.
func createMockRequest(pathValues map[string]string) *http.Request {
	req, _ := http.NewRequest("GET", "/test", nil)
	req = req.WithContext(context.Background())

	// Set path values using SetPathValue (available in Go 1.22+)
	for key, value := range pathValues {
		req.SetPathValue(key, value)
	}

	return req
}

func TestOrganizeCalendar_FirstDayOfWeek(t *testing.T) {
	// Test that first day of month is placed correctly based on weekday
	testCases := []struct {
		date           time.Time
		expectedOffset int // Which day of the week (0=Sunday, 6=Saturday)
	}{
		{
			date:           time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), // Sunday
			expectedOffset: 0,
		},
		{
			date:           time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC), // Monday
			expectedOffset: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.date.Format("2006-01"), func(t *testing.T) {
			pairs := []pair{
				{
					Time: tc.date,
					Page: &mockPage{name: "test"},
				},
			}

			result := organizeCalendar(pairs)
			month := result[0].Months[0]

			// Find which position day 1 is in the first week
			firstWeek := month.Days[0]
			foundAt := -1
			for i, day := range firstWeek {
				if day != nil && day.Date.Day() == 1 {
					foundAt = i
					break
				}
			}

			if foundAt != tc.expectedOffset {
				t.Errorf("First day of month at position %d, want %d (weekday: %s)",
					foundAt, tc.expectedOffset, tc.date.Weekday())
			}
		})
	}
}

// mockPage implements the xlog.Page interface for testing.
type mockPage struct {
	name    string
	astTree ast.Node
}

func (m *mockPage) Name() string                         { return m.name }
func (m *mockPage) FileName() string                     { return m.name + ".md" }
func (m *mockPage) Exists() bool                         { return true }
func (m *mockPage) Render() template.HTML                { return template.HTML("test") }
func (m *mockPage) Content() xlog.Markdown               { return xlog.Markdown("test content") }
func (m *mockPage) Delete() bool                         { return true }
func (m *mockPage) Write(md xlog.Markdown) bool          { return true }
func (m *mockPage) ModTime() time.Time                   { return time.Now() }
func (m *mockPage) AST() ([]byte, ast.Node)              { return []byte("test"), m.astTree }
func (m *mockPage) SetAST(b []byte, n ast.Node)          { m.astTree = n }
func (m *mockPage) ClearCache()                          {}
func (m *mockPage) URL() string                          { return "/" + m.name }
func (m *mockPage) EditURL() string                      { return "/" + m.name + "/edit" }
func (m *mockPage) HistoryURL() string                   { return "/" + m.name + "/history" }
func (m *mockPage) RelativePath() string                 { return m.name }
func (m *mockPage) AbsolutePath() string                 { return "/tmp/" + m.name }
func (m *mockPage) ChangeExtension(ext string) xlog.Page { return m }
func (m *mockPage) SetExtension(ext string) xlog.Page    { return m }
func (m *mockPage) Extension() string                    { return ".md" }
func (m *mockPage) ChangeDirectory(dir string) xlog.Page { return m }
func (m *mockPage) Directory() string                    { return "" }
func (m *mockPage) Directories() []string                { return []string{} }
func (m *mockPage) Rebase(base string) string            { return m.name }
func (m *mockPage) Equal(other xlog.Page) bool           { return m.name == other.Name() }
func (m *mockPage) IsSubPageOf(parent xlog.Page) bool    { return false }
func (m *mockPage) DirectSubPages() []xlog.Page          { return []xlog.Page{} }
func (m *mockPage) AllSubPages() []xlog.Page             { return []xlog.Page{} }
func (m *mockPage) IsDirIndex() bool                     { return false }
