package date

import (
	"html/template"
	"testing"
	"time"

	. "github.com/emad-elsaid/xlog"
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

// mockPage implements the Page interface for testing.
type mockPage struct {
	name    string
	astTree ast.Node
}

func (m *mockPage) Name() string                    { return m.name }
func (m *mockPage) FileName() string                { return m.name + ".md" }
func (m *mockPage) Exists() bool                    { return true }
func (m *mockPage) Render() template.HTML           { return template.HTML("test") }
func (m *mockPage) Content() Markdown               { return Markdown("test content") }
func (m *mockPage) Delete() bool                    { return true }
func (m *mockPage) Write(md Markdown) bool          { return true }
func (m *mockPage) ModTime() time.Time              { return time.Now() }
func (m *mockPage) AST() ([]byte, ast.Node)         { return []byte("test"), m.astTree }
func (m *mockPage) SetAST(b []byte, n ast.Node)     { m.astTree = n }
func (m *mockPage) ClearCache()                     {}
func (m *mockPage) URL() string                     { return "/" + m.name }
func (m *mockPage) EditURL() string                 { return "/" + m.name + "/edit" }
func (m *mockPage) HistoryURL() string              { return "/" + m.name + "/history" }
func (m *mockPage) RelativePath() string            { return m.name }
func (m *mockPage) AbsolutePath() string            { return "/tmp/" + m.name }
func (m *mockPage) ChangeExtension(ext string) Page { return m }
func (m *mockPage) SetExtension(ext string) Page    { return m }
func (m *mockPage) Extension() string               { return ".md" }
func (m *mockPage) ChangeDirectory(dir string) Page { return m }
func (m *mockPage) Directory() string               { return "" }
func (m *mockPage) Directories() []string           { return []string{} }
func (m *mockPage) Rebase(base string) string       { return m.name }
func (m *mockPage) Equal(other Page) bool           { return m.name == other.Name() }
func (m *mockPage) IsSubPageOf(parent Page) bool    { return false }
func (m *mockPage) DirectSubPages() []Page          { return []Page{} }
func (m *mockPage) AllSubPages() []Page             { return []Page{} }
func (m *mockPage) IsDirIndex() bool                { return false }
