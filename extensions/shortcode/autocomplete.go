package shortcode

import . "github.com/emad-elsaid/xlog"

func init() {
	RegisterAutocomplete(autocomplete(0))
}

type autocomplete int

func (a autocomplete) StartChar() string {
	return string(trigger)
}

func (a autocomplete) Suggestions() []*Suggestion {
	suggestions := []*Suggestion{}

	for k := range shortcodes {
		suggestions = append(suggestions, &Suggestion{
			Text:        "/" + k,
			DisplayText: "/" + k,
		})
	}

	return suggestions
}
