package date

import (
	"embed"
	"slices"
	"time"

	"github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

const calendarPageName = "Calendar"

func dateHandler(r xlog.Request) xlog.Output {
	dateV := r.PathValue("date")
	date, err := time.Parse("2-1-2006", dateV)
	if err != nil {
		return xlog.BadRequest(err.Error())
	}

	pages := xlog.MapPage(r.Context(), func(p xlog.Page) xlog.Page {
		_, tree := p.AST()
		allDates := xlog.FindAllInAST[*DateNode](tree)
		for _, d := range allDates {
			if d.time.Equal(date) {
				return p
			}
		}

		return nil
	})

	return xlog.Render("date", xlog.Locals{
		"page":  xlog.DynamicPage{NameVal: date.Format("2 January 2006")},
		"pages": pages,
	})
}

func calendarHandler(r xlog.Request) xlog.Output {
	calendar := []pair{}

	xlog.EachPage(r.Context(), func(p xlog.Page) {
		_, ast := p.AST()
		if ast == nil {
			return
		}

		// Track dates already seen for this page to avoid duplicates
		seenDates := make(map[time.Time]bool)
		for _, v := range xlog.FindAllInAST[*DateNode](ast) {
			// Normalize time to day precision (ignore hour/minute/second)
			normalized := time.Date(v.time.Year(), v.time.Month(), v.time.Day(), 0, 0, 0, 0, v.time.Location())
			if !seenDates[normalized] {
				calendar = append(calendar, pair{Time: normalized, Page: p})
				seenDates[normalized] = true
			}
		}
	})

	cal := organizeCalendar(calendar)

	slices.SortFunc(cal, func(a, b Year) int {
		return int(b.Year) - int(a.Year)
	})

	return xlog.Render("calendar", xlog.Locals{
		"page":     xlog.DynamicPage{NameVal: calendarPageName},
		"calendar": cal,
	})
}

type pair struct {
	Time time.Time
	Page xlog.Page
}

type Day struct {
	Date  time.Time
	Pages []xlog.Page
}

type Month struct {
	Name string
	Days [6][7]*Day // 6 weeks, 7 days per week
}

type Year struct {
	Year   int
	Months []Month
}

func organizeCalendar(pairs []pair) []Year {
	grouped := make(map[int]map[time.Month][]pair)

	for _, p := range pairs {
		year, month := p.Time.Year(), p.Time.Month()
		if grouped[year] == nil {
			grouped[year] = make(map[time.Month][]pair)
		}
		grouped[year][month] = append(grouped[year][month], p)
	}

	var years []Year
	for year, months := range grouped {
		var yearData Year
		yearData.Year = year

		for month, pairs := range months {
			var monthData Month
			monthData.Name = month.String()

			firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
			startOffset := int(firstDay.Weekday())

			var days [6][7]*Day
			for i := range days {
				for j := range days[i] {
					days[i][j] = nil
				}
			}

			currentDay := firstDay
			for day := 1; day <= 31; day++ {
				if currentDay.Month() != month {
					break
				}
				dayOffset := (day + startOffset - 1) / 7
				weekday := (day + startOffset - 1) % 7
				if days[dayOffset][weekday] == nil {
					days[dayOffset][weekday] = &Day{Date: currentDay}
				}
				// Use a map to deduplicate pages for this day
				pageMap := make(map[string]xlog.Page)
				for _, p := range pairs {
					if p.Time.Day() == currentDay.Day() {
						pageMap[p.Page.Name()] = p.Page
					}
				}
				// Convert map back to slice
				for _, page := range pageMap {
					days[dayOffset][weekday].Pages = append(days[dayOffset][weekday].Pages, page)
				}
				currentDay = currentDay.AddDate(0, 0, 1)
			}

			monthData.Days = days
			yearData.Months = append(yearData.Months, monthData)
		}

		years = append(years, yearData)
	}

	return years
}
