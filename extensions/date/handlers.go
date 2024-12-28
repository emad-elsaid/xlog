package date

import (
	"embed"
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
