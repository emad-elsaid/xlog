package link_preview

import (
	"crypto/sha256"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	. "github.com/emad-elsaid/xlog"
)

//go:embed templates
var templates embed.FS

func init() {
	RegisterExtension(LinkPreview{})
}

type LinkPreview struct{}

func (LinkPreview) Name() string { return "link-preview" }
func (LinkPreview) Init() {
	RegisterPreprocessor(imgUrlPreprocessor)
	RegisterPreprocessor(tweetUrlPreprocessor)
	RegisterPreprocessor(youtubeUrlPreprocessor)
	RegisterPreprocessor(fbUrlPreprocessor)
	RegisterPreprocessor(giphyUrlPreprocessor)
	RegisterPreprocessor(fallbackURLPreprocessor)
	RegisterTemplate(templates, "templates")
}

var imgUrlReg = regexp.MustCompile(`(?imU)^(https\:\/\/\S+\.(svg|jpg|jpeg|gif|png|webp))$`)

func imgUrlPreprocessor(c Markdown) Markdown {
	return Markdown(imgUrlReg.ReplaceAllString(string(c), `![]($1)`))
}

var tweetUrlReg = regexp.MustCompile(`(?imU)^https\:\/\/(?:twitter|x)\.com(\/\S+\/status\/[0-9]+)$`)

func tweetUrlPreprocessor(c Markdown) Markdown {
	return Markdown(
		tweetUrlReg.ReplaceAllString(string(c), `
<blockquote class="twitter-tweet">
	<a href="https://twitter.com$1"></a>
</blockquote><script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>`),
	)
}

var youtubeUrlReg = regexp.MustCompile(`(?imU)^https\:\/\/www\.youtube\.com\/watch\?v=(\S+)$`)
var youtubeShortReg = regexp.MustCompile(`(?imU)^https\:\/\/youtu\.be\/(\S+)$`)

func youtubeUrlPreprocessor(c Markdown) Markdown {
	tmplt := `
<figure class="image is-16by9 mx-0">
	<iframe class="has-ratio" width="560" height="315" src="https://www.youtube-nocookie.com/embed/$1" style="border-radius:0.5em;" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</figure>`

	longUrlReplaced := youtubeUrlReg.ReplaceAllString(string(c), tmplt)
	shortUrlReplaced := youtubeShortReg.ReplaceAllString(longUrlReplaced, tmplt)
	return Markdown(shortUrlReplaced)
}

var fbUrlReg = regexp.MustCompile(`(?imU)^(https\:\/\/www\.facebook\.com\/[^\t\n\f\r \/]+/posts/[0-9a-zA-Z]+)$`)

func fbUrlPreprocessor(c Markdown) Markdown {
	return Markdown(
		fbUrlReg.ReplaceAllStringFunc(string(c), func(l string) string {
			return fmt.Sprintf(`
<iframe src="https://www.facebook.com/plugins/post.php?show_text=true&width=500&href=%s" width="500" height="271" style="border:none;overflow:hidden" scrolling="no" frameborder="0" allowfullscreen="true" allow="autoplay; clipboard-write; encrypted-media; picture-in-picture; web-share"></iframe>`, url.QueryEscape(l))
		}),
	)
}

var giphyUrlReg = regexp.MustCompile(`(?imU)^https\:\/\/giphy.com\/gifs\/\S+\-([^\t\n\f\r \-]+)$`)

func giphyUrlPreprocessor(c Markdown) Markdown {
	return Markdown(giphyUrlReg.ReplaceAllString(string(c), `![](https://media.giphy.com/media/$1/giphy.gif)`))
}

var (
	fallbackUrlReg = regexp.MustCompile(`(?imU)^(https?\:\/\/\S+)$`)
	titleReg       = regexp.MustCompile(`(?imU)<title>(.*)</title>`)
	metaReg        = regexp.MustCompile(`(?imU)<meta.*\>`)
	metaNameReg    = regexp.MustCompile(`(?imU)(?:name|property)\s*=\s*"(.*)"`)
	metaContentReg = regexp.MustCompile(`(?imU)content\s*=\s*"(.*)"`)
)

func fallbackURLPreprocessor(c Markdown) Markdown {
	output := fallbackUrlReg.ReplaceAllStringFunc(string(c), func(m string) string {
		m, _ = url.PathUnescape(m)

		meta, err := getUrlMeta(m)
		if err != nil {
			return m
		}

		var title string
		if len(meta.Title) > 0 {
			title = meta.Title
		} else {
			title = m
		}

		u, _ := url.Parse(meta.URL)

		image := meta.Image
		if len(image) > 0 && image[0] == '/' {
			image = u.Scheme + "://" + u.Hostname() + image
		}

		var view string = string(
			Partial("link-preview", Locals{
				"url":         m,
				"title":       title,
				"description": meta.Description,
				"image":       image,
			}),
		)

		return strings.ReplaceAll(view, "\n", "")
	})

	return Markdown(output)
}

type Meta struct {
	URL         string
	Title       string
	Description string
	Image       string
}

func getUrlMeta(url string) (*Meta, error) {
	const cacheDir = ".cache"
	os.Mkdir(cacheDir, 0700)

	cacheFile := path.Join(cacheDir, fmt.Sprintf("%x.json", sha256.Sum256([]byte(url))))
	cache, err := os.ReadFile(cacheFile)
	var meta Meta
	if err == nil {
		if err := json.Unmarshal(cache, &meta); err == nil {
			return &meta, nil
		}
	}

	resp, err := http.Get(url)
	if resp == nil || err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cont, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	html := string(cont)

	titleMatches := titleReg.FindStringSubmatch(html)
	title := url
	if len(titleMatches) >= 1 {
		title = titleMatches[1]
	}

	meta = Meta{
		URL:   url,
		Title: title,
	}

	metaMatches := metaReg.FindAllString(html, -1)
	for _, v := range metaMatches {
		n := metaNameReg.FindStringSubmatch(v)
		if len(n) < 2 {
			continue
		}

		v := metaContentReg.FindStringSubmatch(v)
		if len(v) < 2 {
			continue
		}

		name := strings.ToLower(n[1])
		value := v[1]

		if name == "description" || name == "og:description" {
			meta.Description = value
		} else if name == "og:image" {
			meta.Image = value
		}
	}

	js, _ := json.Marshal(meta)
	os.WriteFile(cacheFile, js, 0644)

	return &meta, nil
}
