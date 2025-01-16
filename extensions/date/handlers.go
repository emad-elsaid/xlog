package date

import (
	"embed"
	"slices"
	"time"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func dateHandler(r Request) Output {
	dateV := r.PathValue("date")
	date, err := time.Parse("2-1-2006", dateV)
	if err != nil {
		return BadRequest(err.Error())
	}

	pages := MapPage(r.Context(), func(p Page) Page {
		_, tree := p.AST()
		allDates := FindAllInAST[*DateNode](tree)
		for _, d := range allDates {
			if d.time.Equal(date) {
				return p
			}
		}

		return nil
	})

	return Render("date", Locals{
		"page":  DynamicPage{NameVal: date.Format("2 January 2006")},
		"pages": pages,
	})
}

func calendarHandler(r Request) Output {
	calendar := []pair{}

	EachPage(r.Context(), func(p Page) {
		_, ast := p.AST()
		if ast == nil {
			return
		}

		for _, v := range FindAllInAST[*DateNode](ast) {
			calendar = append(calendar, pair{Time: v.time, Page: p})
		}
	})

	cal := organizeCalendar(calendar)

	slices.SortFunc(cal, func(a, b Year) int {
		return int(b.Year) - int(a.Year)
	})

	return Render("calendar", Locals{
		"page":     DynamicPage{NameVal: "Calendar"},
		"calendar": cal,
	})
}

type pair struct {
	Time time.Time
	Page Page
}

type Day struct {
	Date  time.Time
	Pages []Page
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
				for _, p := range pairs {
					if p.Time.Day() == currentDay.Day() {
						days[dayOffset][weekday].Pages = append(days[dayOffset][weekday].Pages, p.Page)
					}
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
