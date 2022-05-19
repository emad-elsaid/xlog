package main

import (
	"fmt"
	"regexp"
	"strings"
)

var shortcodes = map[string]preProcessor{
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
		// single line
		reg := regexp.MustCompile(`(?imU)^\/` + regexp.QuoteMeta(k) + `\s+(.*)$`)
		skip := len("/" + k + " ")

		preprocessor := func(r *regexp.Regexp, skip int, v preProcessor) preProcessor {
			return func(c string) string {
				return reg.ReplaceAllStringFunc(c, func(i string) string {
					return v(i[skip:])
				})
			}
		}(reg, skip, v)

		PREPROCESSOR(preprocessor)
		headerSkip := len("```" + k + "\n")

		// multi line
		multireg := regexp.MustCompile("(?imUs)^```" + regexp.QuoteMeta(k) + "$(.*)^```$")
		multilinePreprocessor := func(r *regexp.Regexp, skip int, v preProcessor) preProcessor {
			return func(c string) string {
				return multireg.ReplaceAllStringFunc(c, func(i string) string {
					return v(i[skip : len(i)-4])
				})
			}
		}(reg, headerSkip, v)

		PREPROCESSOR(multilinePreprocessor)
	}
}
