package emoji

import (
	_ "embed"
	"encoding/json"

	. "github.com/emad-elsaid/xlog"
)

//go:embed emoji.json
var emojiFile []byte
var autocomplete = Autocomplete{
	StartChar:   ":",
	Suggestions: []*Suggestion{},
}

func init() {
	AUTOCOMPLETE(autocompleter)

	emojis := []struct {
		Emoji   string   `json:"emoji"`
		Aliases []string `json:"aliases"`
	}{}

	json.Unmarshal(emojiFile, &emojis)

	for _, v := range emojis {
		for _, alias := range v.Aliases {
			autocomplete.Suggestions = append(autocomplete.Suggestions, &Suggestion{
				Text:        ":" + alias + ":",
				DisplayText: v.Emoji + " " + alias,
			})
		}
	}
}

func autocompleter() *Autocomplete {
	return &autocomplete
}
