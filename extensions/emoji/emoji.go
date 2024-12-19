package emoji

import (
	_ "embed"
	"encoding/json"
	"sync"

	. "github.com/emad-elsaid/xlog"
)

//go:embed emoji.json
var emojiFile []byte

func init() {
	RegisterExtension(Emoji{})
}

type Emoji struct{}

func (Emoji) Name() string { return "emoji" }
func (Emoji) Init()        { RegisterAutocomplete(autocomplete{}) }

type autocomplete struct{}

func (a autocomplete) StartChar() string          { return ":" }
func (a autocomplete) Suggestions() []*Suggestion { return suggestions() }

var suggestions = sync.OnceValue(func() []*Suggestion {
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
})
