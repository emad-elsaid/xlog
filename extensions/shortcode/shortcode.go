package shortcode

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/emad-elsaid/xlog"
)

var shortcodes = map[string]Preprocessor{
	"info": func(c Markdown) Markdown {
		return Markdown(fmt.Sprintf(`<p class="notification is-info">%s</p>`, strings.ReplaceAll(string(c), "\n", "<br/>")))
	},

	"success": func(c Markdown) Markdown {
		return Markdown(fmt.Sprintf(`<p class="notification is-success">%s</p>`, strings.ReplaceAll(string(c), "\n", "<br/>")))
	},

	"warning": func(c Markdown) Markdown {
		return Markdown(fmt.Sprintf(`<p class="notification is-warning">%s</p>`, strings.ReplaceAll(string(c), "\n", "<br/>")))
	},

	"alert": func(c Markdown) Markdown {
		return Markdown(fmt.Sprintf(`<p class="notification is-danger">%s</p>`, strings.ReplaceAll(string(c), "\n", "<br/>")))
	},
}

func init() {
	for k, v := range shortcodes {
		ShortCode(k, v)
	}

	RegisterAutocomplete(autocomplete(0))
}

func ShortCode(name string, shortcode Preprocessor) {
	shortcodes[name] = shortcode

	// single line
	reg := regexp.MustCompile(`(?imU)^\/` + regexp.QuoteMeta(name) + `\s+(.*)$`)
	skip := len("/" + name + " ")

	preprocessor := func(r *regexp.Regexp, skip int, v Preprocessor) Preprocessor {
		return func(c Markdown) Markdown {
			output := reg.ReplaceAllStringFunc(string(c), func(i string) string {
				return string(v(Markdown(i[skip:])))
			})

			return Markdown(output)
		}
	}(reg, skip, shortcode)

	RegisterPreprocessor(preprocessor)

	// multi line
	headerSkip := len("```" + name + "\n")
	multireg := regexp.MustCompile("(?imUs)^```" + regexp.QuoteMeta(name) + "$(.*)^```$")
	multilinePreprocessor := func(r *regexp.Regexp, skip int, v Preprocessor) Preprocessor {
		return func(c Markdown) Markdown {
			output := multireg.ReplaceAllStringFunc(string(c), func(i string) string {
				input := i[skip : len(i)-4]
				return string(v(Markdown(input)))
			})
			return Markdown(output)
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
			DisplayText: "/" + k,
		})
	}

	return suggestions
}
