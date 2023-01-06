Xlog ships with a CLI that includes the core and all official extensions. There are cases where you need custom set of extensions:

* Maybe you need only the core features without any extension
* Maybe there is an extension that you don't need or want or misbehaving
* Maybe you developed a set of extensions and you want to include them in your installations

Here is how you can build your own custom xlog with features you select.

# Creating a Go module

Create a directory for your custom installation and initialize a go module in it.

```shell
mkdir custom_xlog
cd custom_xlog
go mod init github.com/yourusername/custom_xlog
```

# Main file

Then create a file `xlog.go` for example with the following content

```go
package main

import (
	// Core
	"github.com/emad-elsaid/xlog"

	// Extensions
	_ "github.com/emad-elsaid/xlog/extensions/activitypub"
	_ "github.com/emad-elsaid/xlog/extensions/autolink"
	_ "github.com/emad-elsaid/xlog/extensions/autolink_pages"
	_ "github.com/emad-elsaid/xlog/extensions/date"
	_ "github.com/emad-elsaid/xlog/extensions/disqus"
	_ "github.com/emad-elsaid/xlog/extensions/emoji"
	_ "github.com/emad-elsaid/xlog/extensions/file_operations"
	_ "github.com/emad-elsaid/xlog/extensions/github"
	_ "github.com/emad-elsaid/xlog/extensions/hashtags"
	_ "github.com/emad-elsaid/xlog/extensions/images"
	_ "github.com/emad-elsaid/xlog/extensions/link_preview"
	_ "github.com/emad-elsaid/xlog/extensions/manifest"
	_ "github.com/emad-elsaid/xlog/extensions/mermaid"
	_ "github.com/emad-elsaid/xlog/extensions/opengraph"
	_ "github.com/emad-elsaid/xlog/extensions/recent"
	_ "github.com/emad-elsaid/xlog/extensions/rss"
	_ "github.com/emad-elsaid/xlog/extensions/rtl"
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
```

# Selecting extensions

The previous file is what xlog ships in `cmd/xlog/xlog.go` if you missed up at any point feel free to go back to it and copy it from there. 

You can now select specific extensions from the list of official extensions. or add custom extensions you developed.

# Running your custom xlog

Now you can use Go to run your custom installation 

```shell
go get github.com/emad-elsaid/xlog
go run xlog.go
```

