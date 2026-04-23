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
		name      string
		pairs     []pair
		wantYears int
	}{
		{
			name:      "empty calendar",
			pairs:     []pair{},
			wantYears: 0,
		},
		{
			name: "single year",
			pairs: []pair{
				{Time: time.Date(2026, time.March, 26, 0, 0, 0, 0, time.UTC)},
			},
			wantYears: 1,
		},
		{
			name: "multiple years",
			pairs: []pair{
				{Time: time.Date(2026, time.March, 26, 0, 0, 0, 0, time.UTC)},
				{Time: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)},
			},
			wantYears: 2,
		},
		{
			name: "same day multiple pages",
			pairs: []pair{
				{Time: time.Date(2026, time.March, 26, 0, 0, 0, 0, time.UTC)},
				{Time: time.Date(2026, time.March, 26, 0, 0, 0, 0, time.UTC)},
			},
			wantYears: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			years := organizeCalendar(tt.pairs)
			if len(years) != tt.wantYears {
				t.Errorf("organizeCalendar() years = %v, want %v", len(years), tt.wantYears)
			}

			// Verify year data structure
			for _, year := range years {
				if year.Year == 0 {
					t.Error("organizeCalendar() year = 0, expected valid year")
				}
				if len(year.Months) == 0 {
					t.Error("organizeCalendar() no months in year")
				}
				for _, month := range year.Months {
					if month.Name == "" {
						t.Error("organizeCalendar() month name is empty")
					}
					// Verify days array structure (6 weeks × 7 days)
					if len(month.Days) != 6 {
						t.Errorf("organizeCalendar() weeks = %d, want 6", len(month.Days))
					}
					for _, week := range month.Days {
						if len(week) != 7 {
							t.Errorf("organizeCalendar() days in week = %d, want 7", len(week))
						}
					}
				}
			}
		})
	}
}

func TestOrganizeCalendar_MonthLayout(t *testing.T) {
	// Test that days are correctly placed in the calendar grid
	// March 2026 starts on a Sunday (weekday 0)
	pairs := []pair{
		{Time: time.Date(2026, time.March, 1, 0, 0, 0, 0, time.UTC)},
		{Time: time.Date(2026, time.March, 15, 0, 0, 0, 0, time.UTC)},
		{Time: time.Date(2026, time.March, 31, 0, 0, 0, 0, time.UTC)},
	}

	years := organizeCalendar(pairs)
	if len(years) != 1 {
		t.Fatalf("expected 1 year, got %d", len(years))
	}

	year := years[0]
	if len(year.Months) != 1 {
		t.Fatalf("expected 1 month, got %d", len(year.Months))
	}

	month := year.Months[0]
	
	// March 1, 2026 is a Sunday, should be in first week, first day (index [0][0])
	if month.Days[0][0] == nil {
		t.Error("March 1st should be in [0][0]")
	} else if month.Days[0][0].Date.Day() != 1 {
		t.Errorf("Expected day 1 in [0][0], got %d", month.Days[0][0].Date.Day())
	}

	// Check that we have some days filled
	dayCount := 0
	for _, week := range month.Days {
		for _, day := range week {
			if day != nil {
				dayCount++
			}
		}
	}

	if dayCount < 3 {
		t.Errorf("Expected at least 3 days filled, got %d", dayCount)
	}
}

func TestOrganizeCalendar_PageAssociation(t *testing.T) {
	testPage := &mockPage{name: "test-page", date: time.Date(2026, time.March, 26, 0, 0, 0, 0, time.UTC)}
	
	pairs := []pair{
		{Time: time.Date(2026, time.March, 26, 0, 0, 0, 0, time.UTC), Page: testPage},
	}

	years := organizeCalendar(pairs)
	if len(years) != 1 {
		t.Fatalf("expected 1 year, got %d", len(years))
	}

	// Find the day with March 26
	month := years[0].Months[0]
	foundPage := false
	for _, week := range month.Days {
		for _, day := range week {
			if day != nil && day.Date.Day() == 26 {
				if len(day.Pages) != 1 {
					t.Errorf("Expected 1 page for March 26, got %d", len(day.Pages))
				}
				if len(day.Pages) > 0 && day.Pages[0].Name() != testPage.Name() {
					t.Error("Expected test page to be associated with March 26")
				}
				foundPage = true
				break
			}
		}
	}

	if !foundPage {
		t.Error("Could not find March 26 in calendar")
	}
}

func TestOrganizeCalendar_MultipleMonths(t *testing.T) {
	pairs := []pair{
		{Time: time.Date(2026, time.January, 15, 0, 0, 0, 0, time.UTC)},
		{Time: time.Date(2026, time.February, 20, 0, 0, 0, 0, time.UTC)},
		{Time: time.Date(2026, time.March, 26, 0, 0, 0, 0, time.UTC)},
	}

	years := organizeCalendar(pairs)
	if len(years) != 1 {
		t.Fatalf("expected 1 year, got %d", len(years))
	}

	if len(years[0].Months) != 3 {
		t.Errorf("expected 3 months, got %d", len(years[0].Months))
	}

	// Verify month names
	monthNames := make(map[string]bool)
	for _, month := range years[0].Months {
		monthNames[month.Name] = true
	}

	expectedMonths := []string{"January", "February", "March"}
	for _, expected := range expectedMonths {
		if !monthNames[expected] {
			t.Errorf("expected month %s not found", expected)
		}
	}
}

func TestOrganizeCalendar_MultiplePagesPerDay(t *testing.T) {
	page1 := &mockPage{name: "page1"}
	page2 := &mockPage{name: "page2"}
	page3 := &mockPage{name: "page3"}

	pairs := []pair{
		{Time: time.Date(2026, time.March, 26, 0, 0, 0, 0, time.UTC), Page: page1},
		{Time: time.Date(2026, time.March, 26, 0, 0, 0, 0, time.UTC), Page: page2},
		{Time: time.Date(2026, time.March, 26, 0, 0, 0, 0, time.UTC), Page: page3},
	}

	years := organizeCalendar(pairs)
	month := years[0].Months[0]
	
	// Find March 26
	foundDay := false
	for _, week := range month.Days {
		for _, day := range week {
			if day != nil && day.Date.Day() == 26 {
				if len(day.Pages) != 3 {
					t.Errorf("Expected 3 pages for March 26, got %d", len(day.Pages))
				}
				foundDay = true
				break
			}
		}
	}

	if !foundDay {
		t.Error("Could not find March 26 in calendar")
	}
}

func TestOrganizeCalendar_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		pairs     []pair
		checkFunc func(*testing.T, []Year)
	}{
		{
			name: "leap year February",
			pairs: []pair{
				{Time: time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)},
			},
			checkFunc: func(t *testing.T, years []Year) {
				if len(years) != 1 {
					t.Errorf("expected 1 year, got %d", len(years))
				}
				if years[0].Year != 2024 {
					t.Errorf("expected year 2024, got %d", years[0].Year)
				}
			},
		},
		{
			name: "year boundary",
			pairs: []pair{
				{Time: time.Date(2025, time.December, 31, 0, 0, 0, 0, time.UTC)},
				{Time: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)},
			},
			checkFunc: func(t *testing.T, years []Year) {
				if len(years) != 2 {
					t.Errorf("expected 2 years, got %d", len(years))
				}
			},
		},
		{
			name: "distant years",
			pairs: []pair{
				{Time: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)},
				{Time: time.Date(2050, time.December, 31, 0, 0, 0, 0, time.UTC)},
			},
			checkFunc: func(t *testing.T, years []Year) {
				if len(years) != 2 {
					t.Errorf("expected 2 years, got %d", len(years))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			years := organizeCalendar(tt.pairs)
			tt.checkFunc(t, years)
		})
	}
}

// Mock Page implementation for testing
type mockPage struct {
	name string
	date time.Time
}

func (m *mockPage) Name() string                   { return m.name }
func (m *mockPage) FileName() string               { return m.name + ".md" }
func (m *mockPage) Exists() bool                   { return true }
func (m *mockPage) Render() template.HTML          { return template.HTML(m.name) }
func (m *mockPage) Content() Markdown              { return Markdown("# " + m.name) }
func (m *mockPage) Delete() bool                   { return true }
func (m *mockPage) Write(Markdown) bool            { return true }
func (m *mockPage) ModTime() time.Time             { return m.date }
func (m *mockPage) AST() ([]byte, ast.Node) {
	doc := ast.NewDocument()
	if !m.date.IsZero() {
		doc.AppendChild(doc, &DateNode{time: m.date})
	}
	return []byte("# " + m.name), doc
}
