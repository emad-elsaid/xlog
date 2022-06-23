package main

import (
	_ "embed"
	"encoding/json"
)

//go:embed emoji.json
var emojiFile []byte
var emojiAC = Autocomplete{
	StartChar:   ":",
	Suggestions: []*Suggestion{},
}

func init() {
	AUTOCOMPLETE(emojiAutocomplete)

	emojis := []struct {
		Emoji   string   `json:"emoji"`
		Aliases []string `json:"aliases"`
	}{}

	json.Unmarshal(emojiFile, &emojis)

	for _, v := range emojis {
		for _, alias := range v.Aliases {
			emojiAC.Suggestions = append(emojiAC.Suggestions, &Suggestion{
				Text:        ":" + alias + ":",
				DisplayText: v.Emoji + " " + alias,
			})
		}
	}
}

func emojiAutocomplete() *Autocomplete {
	return &emojiAC
}
