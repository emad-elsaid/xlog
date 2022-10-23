package shortcode

import (
	"fmt"
	"regexp"
	"strings"

	. "github.com/emad-elsaid/xlog"
)

var shortcodes = map[string]PreProcessor{
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
		SHORTCODE(k, v)
	}

	AUTOCOMPLETE(autocompleter)
}

func SHORTCODE(name string, shortcode func(string) string) {
	shortcodes[name] = shortcode

	// single line
	reg := regexp.MustCompile(`(?imU)^\/` + regexp.QuoteMeta(name) + `\s+(.*)$`)
	skip := len("/" + name + " ")

	preprocessor := func(r *regexp.Regexp, skip int, v PreProcessor) PreProcessor {
		return func(c string) string {
			return reg.ReplaceAllStringFunc(c, func(i string) string {
				return v(i[skip:])
			})
		}
	}(reg, skip, shortcode)

	PREPROCESSOR(preprocessor)

	// multi line
	headerSkip := len("```" + name + "\n")
	multireg := regexp.MustCompile("(?imUs)^```" + regexp.QuoteMeta(name) + "$(.*)^```$")
	multilinePreprocessor := func(r *regexp.Regexp, skip int, v PreProcessor) PreProcessor {
		return func(c string) string {
			return multireg.ReplaceAllStringFunc(c, func(i string) string {
				return v(i[skip : len(i)-4])
			})
		}
	}(reg, headerSkip, shortcode)

	PREPROCESSOR(multilinePreprocessor)
}

func autocompleter() *Autocomplete {
	a := &Autocomplete{
		StartChar:   "/",
		Suggestions: []*Suggestion{},
	}

	for k := range shortcodes {
		a.Suggestions = append(a.Suggestions, &Suggestion{
			Text:        "/" + k,
			DisplayText: k,
		})
	}

	return a
}
