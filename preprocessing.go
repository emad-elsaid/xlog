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

func PREPROCESSOR(f preProcessor) {
	preProcessors = append(preProcessors, f)
}

func preProcess(content string) string {
	for _, v := range preProcessors {
		content = v(content)
	}

	return content
}
