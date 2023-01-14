package autolink_pages

import (
	"context"

	. "github.com/emad-elsaid/xlog"
)

func init() {
	RegisterAutocomplete(autocomplete(0))
}

type autocomplete int

func (a autocomplete) StartChar() string {
	return "@"
}

func (a autocomplete) Suggestions() []*Suggestion {
	suggestions := []*Suggestion{}

	EachPage(context.Background(), func(p Page) {
		suggestions = append(suggestions, &Suggestion{
			Text:        p.Name(),
			DisplayText: "@" + p.Name(),
		})
	})

	return suggestions
}
