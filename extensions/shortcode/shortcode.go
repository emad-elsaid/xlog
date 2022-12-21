package shortcode

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"

	. "github.com/emad-elsaid/xlog"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

type ShortCodeFunc func(Markdown) template.HTML

var shortcodes = map[string]ShortCodeFunc{
	"info": func(c Markdown) template.HTML {
		return template.HTML(fmt.Sprintf(`<p class="notification is-info">%s</p>`, strings.ReplaceAll(string(c), "\n", "<br/>")))
	},

	"success": func(c Markdown) template.HTML {
		return template.HTML(fmt.Sprintf(`<p class="notification is-success">%s</p>`, strings.ReplaceAll(string(c), "\n", "<br/>")))
	},

	"warning": func(c Markdown) template.HTML {
		return template.HTML(fmt.Sprintf(`<p class="notification is-warning">%s</p>`, strings.ReplaceAll(string(c), "\n", "<br/>")))
	},

	"alert": func(c Markdown) template.HTML {
		return template.HTML(fmt.Sprintf(`<p class="notification is-danger">%s</p>`, strings.ReplaceAll(string(c), "\n", "<br/>")))
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
