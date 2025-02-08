```
Usage of xlog:
  -activitypub.domain string
    	domain used for activitypub stream absolute URLs
  -activitypub.icon string
    	the path to the activitypub profile icon. mastodon use it as profile picture for example. (default "/public/logo.png")
  -activitypub.image string
    	the path to the activitypub profile image. mastodon use it as profile cover for example. (default "/public/logo.png")
  -activitypub.summary string
    	summary of the user for activitypub actor
  -activitypub.username string
    	username for activitypub actor
  -bind string
    	IP and port to bind the web server to (default "127.0.0.1:3000")
  -build string
    	Build all pages as static site in this directory
  -codestyle string
    	code highlighting style name from the list supported by https://pkg.go.dev/github.com/alecthomas/chroma/v2/styles (default "dracula")
  -csrf-cookie string
    	CSRF cookie name (default "xlog_csrf")
  -custom.after_view string
    	path to a file it's content will be included in every page AFTER the content of the page
  -custom.before_view string
    	path to a file it's content will be included in every page BEFORE the content of the page
  -custom.head string
    	path to a file it's content will be included in every page <head> tag
  -disabled-extensions string
    	disable list of extensions by name, comma separated
  -disqus string
    	Disqus domain name for example: xlog-emadelsaid.disqus.com
  -editor string
    	command to use to open pages for editing (default "emacsclient -n -a emacs")
  -github.url string
    	Repository url for 'edit on Github' quick action e.g https://github.com/emad-elsaid/xlog/edit/master/docs
  -gpg string
    	PGP key ID to decrypt and edit .md.pgp files using gpg. if empty encryption will be off
  -html
    	Consider HTML files as pages
  -index string
    	Index file name used as home page (default "index")
  -notfoundpage string
    	Custom not found page (default "404")
  -og.domain string
    	opengraph domain name to be used for meta tags of og:* and twitter:*
  -pandoc
    	Use pandoc to render .org, .rst, .rtf, .odt
  -readonly
    	Should xlog hide write operations, read-only means all write operations will be disabled
  -rss.description string
    	RSS feed description
  -rss.domain string
    	RSS domain name to be used for RSS feed. without HTTPS://
  -rss.limit int
    	Limit the number of items in the RSS feed to this amount (default 30)
  -serve-insecure
    	Accept http connections and forward crsf cookie over non secure connections
  -sitemap.domain string
    	domain name without protocol or trailing / to use for sitemap loc
  -sitename string
    	Site name is the name that appears on the header beside the logo and in the title tag (default "XLOG")
  -source string
    	Directory that will act as a storage (default "/home/emad/code/xlog")
  -sql-table.threshold int
    	If a table rows is more than this threshold it'll allow users to query it with SQL (default 100)
  -theme string
    	bulma theme to use. (light, dark). empty value means system preference is used
  -twitter.username string
    	user twitter account @handle. including the @
```
