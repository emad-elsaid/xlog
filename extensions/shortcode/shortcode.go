package shortcode

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/emad-elsaid/xlog"
)

var shortcodes = map[string]Preprocessor{
	"info": func(c string) string {
		return fmt.Sprintf(`<p class="notification is-info">%s</p>`, strings.ReplaceAll(c, "\n", "<br/>"))
	},

	"success": func(c string) string {
		return fmt.Sprintf(`<p class="notification is-success">%s</p>`, strings.ReplaceAll(c, "\n", "<br/>"))
	},

	"warning": func(c string) string {
		return fmt.Sprintf(`<p class="notification is-warning">%s</p>`, strings.ReplaceAll(c, "\n", "<br/>"))
	},

	"alert": func(c string) string {
		return fmt.Sprintf(`<p class="notification is-danger">%s</p>`, strings.ReplaceAll(c, "\n", "<br/>"))
	},
}

func init() {
	for k, v := range shortcodes {
		ShortCode(k, v)
	}

	RegisterAutocomplete(autocomplete(0))
}

func ShortCode(name string, shortcode func(string) string) {
	shortcodes[name] = shortcode

	// single line
	reg := regexp.MustCompile(`(?imU)^\/` + regexp.QuoteMeta(name) + `\s+(.*)$`)
	skip := len("/" + name + " ")

	preprocessor := func(r *regexp.Regexp, skip int, v Preprocessor) Preprocessor {
		return func(c string) string {
			return reg.ReplaceAllStringFunc(c, func(i string) string {
				return v(i[skip:])
			})
		}
	}(reg, skip, shortcode)

	RegisterPreprocessor(preprocessor)

	// multi line
	headerSkip := len("```" + name + "\n")
	multireg := regexp.MustCompile("(?imUs)^```" + regexp.QuoteMeta(name) + "$(.*)^```$")
	multilinePreprocessor := func(r *regexp.Regexp, skip int, v Preprocessor) Preprocessor {
		return func(c string) string {
			return multireg.ReplaceAllStringFunc(c, func(i string) string {
				return v(i[skip : len(i)-4])
			})
		}
	}(reg, headerSkip, shortcode)

	RegisterPreprocessor(multilinePreprocessor)
}

type autocomplete int

func (a autocomplete) StartChar() string {
	return "/"
}

func (a autocomplete) Suggestions() []*Suggestion {
	suggestions := []*Suggestion{}

	for k := range shortcodes {
		suggestions = append(suggestions, &Suggestion{
			Text:        "/" + k,
			DisplayText: k,
		})
	}

	return suggestions
}
