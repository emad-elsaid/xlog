package date

import (
	"html/template"

	"github.com/emad-elsaid/xlog"
)

func links(xlog.Page) []xlog.Command {
	return []xlog.Command{
		Calendar{},
	}
}

type Calendar struct{}

func (Calendar) Icon() string { return "fa-regular fa-calendar-days" }
func (Calendar) Name() string { return "Calendar" }
func (Calendar) Attrs() map[template.HTMLAttr]any {
	return map[template.HTMLAttr]any{
		"href": "/+/calendar",
	}
}
