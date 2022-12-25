package shortcode

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"

	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

type ShortCodeFunc func(Markdown) template.HTML

func render(i Markdown) string {
	var b bytes.Buffer
	MarkDownRenderer.Convert([]byte(i), &b)
	return b.String()
}

var shortcodes = map[string]ShortCodeFunc{
	"info": func(c Markdown) template.HTML {
		return template.HTML(fmt.Sprintf(`<div class="notification is-info is-light">%s</div>`, render(c)))
	},

	"success": func(c Markdown) template.HTML {
		return template.HTML(fmt.Sprintf(`<div class="notification is-success is-light">%s</div>`, render(c)))
	},

	"warning": func(c Markdown) template.HTML {
		return template.HTML(fmt.Sprintf(`<div class="notification is-warning is-light">%s</div>`, render(c)))
	},

	"alert": func(c Markdown) template.HTML {
		return template.HTML(fmt.Sprintf(`<div class="notification is-danger is-light">%s</div>`, render(c)))
	},
}

func init() {
	for k, v := range shortcodes {
		ShortCode(k, v)
	}

	MarkDownRenderer.Parser().AddOptions(parser.WithBlockParsers(
		util.Prioritized(&shortCodeParser{}, 0),
	))
}

func ShortCode(name string, shortcode ShortCodeFunc) {
	shortcodes[name] = shortcode

	headerSkip := len("```" + name + "\n")
	multireg := regexp.MustCompile("(?imUs)^```" + regexp.QuoteMeta(name) + "$(.*)^```$")
	multilinePreprocessor := func(skip int, v ShortCodeFunc) Preprocessor {
		return func(c Markdown) Markdown {
			output := multireg.ReplaceAllStringFunc(string(c), func(i string) string {
				input := i[skip : len(i)-4]
				return string(v(Markdown(input)))
			})
			return Markdown(output)
		}
	}(headerSkip, shortcode)

	RegisterPreprocessor(multilinePreprocessor)
}
