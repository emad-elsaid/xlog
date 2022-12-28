package date

import (
	"embed"
	"time"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	RegisterTemplate(templates, "templates")
	Get(`/\+/date/{date}`, dateHandler)
}

func dateHandler(w Response, r Request) Output {
	vars := Vars(r)
	dateV := vars["date"]
	date, err := time.Parse("2-1-2006", dateV)
	if err != nil {
		return BadRequest(err.Error())
	}

	pages := []Page{}
	EachPage(r.Context(), func(p Page) {
		allDates := FindAllInAST[*DateNode](p.AST(), KindDate)
		for _, d := range allDates {
			if d.time.Equal(date) {
				pages = append(pages, p)
				break
			}
		}
	})

	return Render("date", Locals{
		"title": date.Format("2 January 2006"),
		"pages": pages,
	})
}
