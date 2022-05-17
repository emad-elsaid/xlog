package main

import (
	"fmt"
	"regexp"
)

var shortcodes = map[string]preProcessor{
	"info": func(c string) string {
		return fmt.Sprintf(`<pre class="notification is-info">%s</pre>`, c)
	},

	"success": func(c string) string {
		return fmt.Sprintf(`<pre class="notification is-success">%s</pre>`, c)
	},

	"warning": func(c string) string {
		return fmt.Sprintf(`<pre class="notification is-warning">%s</pre>`, c)
	},

	"alert": func(c string) string {
		return fmt.Sprintf(`<pre class="notification is-danger">%s</pre>`, c)
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
