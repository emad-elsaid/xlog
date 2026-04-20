package date

import (
	"html/template"
	"testing"
	"time"

	. "github.com/emad-elsaid/xlog"
	"github.com/emad-elsaid/xlog/markdown/ast"
)

// testPage implements the Page interface for testing
type testPage struct {
	name string
}

func (t *testPage) Name() string                   { return t.name }
func (t *testPage) FileName() string               { return t.name + ".md" }
func (t *testPage) Exists() bool                   { return true }
func (t *testPage) Render() template.HTML          { return "" }
func (t *testPage) Content() Markdown              { return "" }
func (t *testPage) Delete() bool                   { return false }
func (t *testPage) Write(Markdown) bool            { return false }
func (t *testPage) ModTime() time.Time             { return time.Now() }
func (t *testPage) AST() ([]byte, ast.Node)        { return nil, nil }

func TestOrganizeCalendar(t *testing.T) {
	tests := []struct {
		name          string
		pairs         []pair
		expectedYears int
		checkYear     int
		checkMonths   int
	}{
		{
			name:          "Empty calendar",
			pairs:         []pair{},
			expectedYears: 0,
		},
		{
			name: "Single date single page",
			pairs: []pair{
				{Time: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "test"}},
			},
			expectedYears: 1,
			checkYear:     2026,
			checkMonths:   1,
		},
		{
			name: "Multiple dates same month",
			pairs: []pair{
				{Time: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "day1"}},
				{Time: time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "day15"}},
				{Time: time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "day30"}},
			},
			expectedYears: 1,
			checkYear:     2026,
			checkMonths:   1,
		},
		{
			name: "Multiple months same year",
			pairs: []pair{
				{Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "jan"}},
				{Time: time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "jun"}},
				{Time: time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "dec"}},
			},
			expectedYears: 1,
			checkYear:     2026,
			checkMonths:   3,
		},
		{
			name: "Multiple years",
			pairs: []pair{
				{Time: time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "2024"}},
				{Time: time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "2025"}},
				{Time: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "2026"}},
			},
			expectedYears: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := organizeCalendar(tt.pairs)

			if len(result) != tt.expectedYears {
				t.Errorf("organizeCalendar() returned %d years, want %d", len(result), tt.expectedYears)
			}

			if tt.checkYear > 0 {
				var foundYear *Year
				for i := range result {
					if result[i].Year == tt.checkYear {
						foundYear = &result[i]
						break
					}
				}

				if foundYear == nil {
					t.Fatalf("organizeCalendar() did not contain year %d", tt.checkYear)
				}

				if tt.checkMonths > 0 && len(foundYear.Months) != tt.checkMonths {
					t.Errorf("Year %d has %d months, want %d", tt.checkYear, len(foundYear.Months), tt.checkMonths)
				}
			}
		})
	}
}

func TestOrganizeCalendar_DayPlacement(t *testing.T) {
	// Test that days are placed in the correct week and weekday slots
	// April 2026 starts on Wednesday (weekday 3)
	pairs := []pair{
		{Time: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "april1"}},
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

	// April 1, 2026 is a Wednesday (weekday 3)
	// It should be in the first week (index 0), on Wednesday (weekday 3)
	firstDay := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	expectedWeekday := int(firstDay.Weekday())

	if month.Days[0][expectedWeekday] == nil {
		t.Errorf("Expected day at week 0, weekday %d, got nil", expectedWeekday)
	} else {
		if month.Days[0][expectedWeekday].Date.Day() != 1 {
			t.Errorf("Expected day 1, got day %d", month.Days[0][expectedWeekday].Date.Day())
		}
		if len(month.Days[0][expectedWeekday].Pages) != 1 {
			t.Errorf("Expected 1 page for April 1, got %d", len(month.Days[0][expectedWeekday].Pages))
		}
	}
}

func TestOrganizeCalendar_MultiplePagesSameDay(t *testing.T) {
	// Test that multiple pages on the same date are collected together
	pairs := []pair{
		{Time: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "page1"}},
		{Time: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "page2"}},
		{Time: time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "page3"}},
	}

	result := organizeCalendar(pairs)

	year := result[0]
	month := year.Months[0]

	// Find April 20 in the calendar grid
	found := false
	for week := range month.Days {
		for weekday := range month.Days[week] {
			day := month.Days[week][weekday]
			if day != nil && day.Date.Day() == 20 {
				if len(day.Pages) != 3 {
					t.Errorf("Expected 3 pages for April 20, got %d", len(day.Pages))
				}
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		t.Error("Did not find April 20 in calendar")
	}
}

func TestOrganizeCalendar_MonthBoundaries(t *testing.T) {
	// Test that dates don't bleed across month boundaries
	pairs := []pair{
		{Time: time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "march31"}},
		{Time: time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC), Page: &testPage{name: "april1"}},
	}

	result := organizeCalendar(pairs)

	if len(result) != 1 {
		t.Fatalf("Expected 1 year, got %d", len(result))
	}

	year := result[0]
	if len(year.Months) != 2 {
		t.Fatalf("Expected 2 months (March and April), got %d", len(year.Months))
	}

	// Verify that pages are only attached to the correct dates
	// (Note: organizeCalendar fills entire month calendar with all days,
	// but pages should only be on the dates we specified)
	for _, month := range year.Months {
		for week := range month.Days {
			for weekday := range month.Days[week] {
				day := month.Days[week][weekday]
				if day != nil && len(day.Pages) > 0 {
					if month.Name == "March" && day.Date.Day() != 31 {
						t.Errorf("March calendar has pages on day %d, expected only day 31", day.Date.Day())
					}
					if month.Name == "April" && day.Date.Day() != 1 {
						t.Errorf("April calendar has pages on day %d, expected only day 1", day.Date.Day())
					}
				}
			}
		}
	}
}
