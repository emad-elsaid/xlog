package extension

import (
	"github.com/emad-elsaid/xlog/markdown"
)

type gfm struct {
}

// GFM is an extension that provides Github Flavored markdown functionalities.
var GFM = &gfm{}

func (e *gfm) Extend(m markdown.Markdown) {
	Linkify.Extend(m)
	Table.Extend(m)
	Strikethrough.Extend(m)
	TaskList.Extend(m)
}
