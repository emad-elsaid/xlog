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
	Get(`/+/date/{date}`, dateHandler)
}

func dateHandler(w Response, r Request) Output {
	dateV := r.PathValue("date")
	date, err := time.Parse("2-1-2006", dateV)
	if err != nil {
		return BadRequest(err.Error())
	}

	pages := MapPageCon(r.Context(), func(p Page) *Page {
		allDates := FindAllInAST[*DateNode](p.AST())
		for _, d := range allDates {
			if d.time.Equal(date) {
				return &p
			}
		}

		return nil
	})

	return Render("date", Locals{
		"title": date.Format("2 January 2006"),
		"pages": pages,
	})
}
