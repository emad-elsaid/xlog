```
xlog --help

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
  -github.branch string
        Github repository branch to use for 'edit on Github' quick action (default "master")
  -github.repo string
        Github repository to use for 'edit on Github' quick action
  -index string
        Index file name used as home page (default "index")
  -og.domain string
        opengraph domain name to be used for meta tags of og:* and twitter:*
  -readonly
        Should xlog hide write operations, read-only means all write operations will be disabled
  -rss.description string
        RSS feed description
  -rss.domain string
        RSS domain name to be used for RSS feed. without HTTPS://
  -rss.limit int
        Limit the number of items in the RSS feed to this amount (default 30)
  -sidebar
        Should render sidebar. (default true)
  -sitemap.domain string
        domain name without protocol or trailing / to use for sitemap loc
  -sitename string
        Site name is the name that appears on the header beside the logo and in the title tag (default "XLOG")
  -source string
        Directory that will act as a storage (default "/path/to/source")
  -twitter.username string
        user twitter account @handle. including the @
```
