Defined under `/extensions` sub package. each extension is a subpackage. **All extensions are imported** by default in `cmd/xlog/xlog.go`.


| Extension       | Description                                                                         |
|-----------------|-------------------------------------------------------------------------------------|
| ActiviyPub      | Implements webfinger and activityPub actor and exposing pages as activitypub outbox |
| Autolink        | Shorten a link string so it wouldn't take unnecessary space                         |
| Autolink pages  | Convert a page name mentions in the middle of text to a link                        |
| Date            | Detects dates and converts them to link to a page which lists all pages mentions it |
| Disqus          | Add Disqus comments after the view page if -disqus flag is passed                   |
| Emoji           | Emoji autocomplete while editing                                                    |
| File operations | Add a tool item to delete and rename current page                                   |
| Github          | Adds "Edit on github" quick action                                                  |
| Hashtags        | Add support for hashtags #hashtag syntax                                            |
| Link preview    | Preview tweets, Facebook posts, youtube videos, Giphy links                         |
| Manifest        | adds manifest.json to head tag and output proper JSON value.                        |
| Mermaid         | Support for MermaidJS graphing library                                              |
| Opengraph       | Adds Opengraph meta tags for title, type, image                                     |
| RSS             | Provides RSS feed served under /+/feed.rss and added to the header of pages         |
| RTL             | Fixes text direction for RTL languages in the view page                             |
| Recent          | Adds an item to footer to list all pages ordered by last modified page file.        |
| Search          | Full text search                                                                    |
| Shortcode       | adds a way for short codes (one line and block)                                     |
| Sitemap         | adds support for sitemap.xml for search engine crawling                             |
| Star            | Star pages to pin them to footer                                                    |
| Upload file     | Add support for upload files, screenshots, audio and camera recording               |
| Versions        | Keeps list of pages older versions                                                  |
