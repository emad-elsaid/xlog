package main

import (
	// Core
	"github.com/emad-elsaid/xlog"

	// Extensions
	_ "github.com/emad-elsaid/xlog/extensions/activitypub"
	_ "github.com/emad-elsaid/xlog/extensions/autolink"
	_ "github.com/emad-elsaid/xlog/extensions/autolink_pages"
	_ "github.com/emad-elsaid/xlog/extensions/emoji"
	_ "github.com/emad-elsaid/xlog/extensions/file_operations"
	_ "github.com/emad-elsaid/xlog/extensions/github"
	_ "github.com/emad-elsaid/xlog/extensions/hashtags"
	_ "github.com/emad-elsaid/xlog/extensions/link_preview"
	_ "github.com/emad-elsaid/xlog/extensions/manifest"
	_ "github.com/emad-elsaid/xlog/extensions/opengraph"
	_ "github.com/emad-elsaid/xlog/extensions/recent"
	_ "github.com/emad-elsaid/xlog/extensions/rss"
	_ "github.com/emad-elsaid/xlog/extensions/search"
	_ "github.com/emad-elsaid/xlog/extensions/shortcode"
	_ "github.com/emad-elsaid/xlog/extensions/sitemap"
	_ "github.com/emad-elsaid/xlog/extensions/star"
	_ "github.com/emad-elsaid/xlog/extensions/upload_file"
	_ "github.com/emad-elsaid/xlog/extensions/versions"
)

func main() {
	xlog.Start()
}
