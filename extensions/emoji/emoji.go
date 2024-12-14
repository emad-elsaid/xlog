package emoji

import (
	_ "embed"
	"encoding/json"

	. "github.com/emad-elsaid/xlog"
)

//go:embed emoji.json
var emojiFile []byte

func init() {
	RegisterExtension(Emoji{})
}

type Emoji struct{}

func (Emoji) Name() string { return "emoji" }
func (Emoji) Init() {
	RegisterAutocomplete(autocomplete(0))
}

type autocomplete int

func (a autocomplete) StartChar() string {
	return ":"
}

// TODO this is a bit inefficient as it parses the emoji json everytime
func (a autocomplete) Suggestions() []*Suggestion {
	emojis := []struct {
		Emoji   string   `json:"emoji"`
		Aliases []string `json:"aliases"`
	}{}

	json.Unmarshal(emojiFile, &emojis)

	suggestions := []*Suggestion{}

	for _, v := range emojis {
		for _, alias := range v.Aliases {
			suggestions = append(suggestions, &Suggestion{
				Text:        ":" + alias + ":",
				DisplayText: ":" + v.Emoji + " " + alias,
			})
		}
	}

	return suggestions
}
