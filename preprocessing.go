package main

import (
	"fmt"
	"net/url"
	"regexp"
)

type preProcessor func(string) string

var imgUrlReg = regexp.MustCompile(`(?imU)^(https\:\/\/[^ ]+\.(svg|jpg|jpeg|gif|png|webp))$`)

func imgUrlPreprocessor(c string) string {
	return imgUrlReg.ReplaceAllString(c, `![]($1)`)
}

var tweetUrlReg = regexp.MustCompile(`(?imU)^(https\:\/\/twitter.com\/[^ ]+\/status\/[0-9]+)$`)

func tweetUrlPreprocessor(c string) string {
	return tweetUrlReg.ReplaceAllString(c, `
<blockquote class="twitter-tweet">
	<a href="$1"></a>
</blockquote><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>`)
}

var youtubeUrlReg = regexp.MustCompile(`(?imU)^https\:\/\/www\.youtube\.com\/watch\?v=([^ ]+)$`)

func youtubeUrlPreprocessor(c string) string {
	return youtubeUrlReg.ReplaceAllString(c, `
<figure class="image is-16by9">
	<iframe class="has-ratio" width="560" height="315" src="https://www.youtube-nocookie.com/embed/$1" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</figure>`)
}

var fbUrlReg = regexp.MustCompile(`(?imU)^(https\:\/\/www\.facebook\.com\/[^ \/]+/posts/[0-9]+)$`)

func fbUrlPreprocessor(c string) string {
	return fbUrlReg.ReplaceAllStringFunc(c, func(l string) string {
		return fmt.Sprintf(`
<iframe src="https://www.facebook.com/plugins/post.php?show_text=true&width=500&href=%s" width="500" height="271" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowfullscreen="true" allow="autoplay; clipboard-write; encrypted-media; picture-in-picture; web-share"></iframe>`, url.QueryEscape(l))
	})
}

var giphyUrlReg = regexp.MustCompile(`(?imU)^https\:\/\/giphy.com\/gifs\/[^ ]+\-([^ \-]+)$`)

func giphyUrlPreprocessor(c string) string {
	return giphyUrlReg.ReplaceAllString(c, `![](https://media.giphy.com/media/$1/giphy.gif)`)
}

var preProcessors = []preProcessor{
	imgUrlPreprocessor,
	tweetUrlPreprocessor,
	youtubeUrlPreprocessor,
	fbUrlPreprocessor,
	giphyUrlPreprocessor,
}

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
	fillInShortcodes()
}

func preProcess(content string) string {
	for _, v := range preProcessors {
		content = v(content)
	}

	return content
}

func fillInShortcodes() {
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

		preProcessors = append(preProcessors, preprocessor)
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

		preProcessors = append(preProcessors, multilinePreprocessor)
	}
}
