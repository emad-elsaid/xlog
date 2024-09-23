package shortcode

import (
	"strings"

	. "github.com/emad-elsaid/xlog"
)

func init() {
	RegisterAutocomplete(autocomplete(0))
}

type autocomplete int

func (a autocomplete) StartChar() string {
	return string(trigger)
}

func (a autocomplete) Suggestions() []*Suggestion {
	suggestions := []*Suggestion{}

	for k, s := range shortcodes {
		if strings.ContainsRune(s.Default, '\n') {
			suggestions = append(suggestions, &Suggestion{
				Text:        "```" + k + "\n" + s.Default + "\n```",
				DisplayText: "/" + k,
			})
		} else if s.Default != "" {
			suggestions = append(suggestions, &Suggestion{
				Text:        "/" + k + " " + s.Default,
				DisplayText: "/" + k,
			})
		} else {
			suggestions = append(suggestions, &Suggestion{
				Text:        "/" + k,
				DisplayText: "/" + k,
			})
		}
	}

	return suggestions
}
