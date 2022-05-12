package main

import (
	"fmt"
	"regexp"
)

type preProcessor func(string) string

var (
	imgUrlReg     = regexp.MustCompile(`(?imU)^(https\:\/\/[^ ]+\.(svg|jpg|jpeg|gif|png|webp))$`)
	tweetUrlReg   = regexp.MustCompile(`(?imU)^(https\:\/\/twitter.com\/[^ ]+\/status\/[0-9]+)$`)
	youtubeUrlReg = regexp.MustCompile(`(?imU)^https\:\/\/www\.youtube\.com\/watch\?v=([^ ]+)$`)

	preProcessors = []preProcessor{
		// image
		func(c string) string { return imgUrlReg.ReplaceAllString(c, `<img src="$1"/>`) },

		// twitter
		func(c string) string {
			return tweetUrlReg.ReplaceAllString(c, `
<blockquote class="twitter-tweet">
	<a href="$1"></a>
</blockquote><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>`)
		},

		// youtube
		func(c string) string {
			return youtubeUrlReg.ReplaceAllString(c, `
<figure class="image is-16by9">
	<iframe class="has-ratio" width="560" height="315" src="https://www.youtube-nocookie.com/embed/$1" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</figure>`)
		},
	}
)

var shortcodes = map[string]preProcessor{
	"info": func(c string) string {
		return fmt.Sprintf(`<div class="notification is-info">%s</div>`, c)
	},

	"success": func(c string) string {
		return fmt.Sprintf(`<div class="notification is-success">%s</div>`, c)
	},

	"warning": func(c string) string {
		return fmt.Sprintf(`<div class="notification is-warning">%s</div>`, c)
	},

	"alert": func(c string) string {
		return fmt.Sprintf(`<div class="notification is-danger">%s</div>`, c)
	},
}

func init() {
	for k, v := range shortcodes {
		reg := regexp.MustCompile(`(?imU)^\/` + regexp.QuoteMeta(k) + `\s+(.*)$`)
		skip := len("/" + k + " ")

		preprocessor := func(r *regexp.Regexp, skip int, v preProcessor) preProcessor {
			return func(c string) string {
				return reg.ReplaceAllStringFunc(c, func(i string) string {
					return v(i[skip:])
				})
			}
		}(reg, skip, v)

		preProcessors = append(preProcessors, preprocessor)
	}
}

func preProcess(content string) string {
	for _, v := range preProcessors {
		content = v(content)
	}

	return content
}
