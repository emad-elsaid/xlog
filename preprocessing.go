package main

import "regexp"

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
			return tweetUrlReg.ReplaceAllString(c, `<blockquote class="twitter-tweet"><a href="$1"></a></blockquote><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>`)
		},

		// youtube
		func(c string) string {
			return youtubeUrlReg.ReplaceAllString(c, `<figure class="image is-16by9"><iframe class="has-ratio" width="560" height="315" src="https://www.youtube-nocookie.com/embed/$1" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe></figure>`)
		},
	}
)

func preProcess(content string) string {
	for _, v := range preProcessors {
		content = v(content)
	}

	return content
}
