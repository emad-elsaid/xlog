Defined under `/extensions` sub package. each extension is a subpackage. **All extensions are imported** by default in `cmd/xlog/xlog.go`.


| Extension       | Description                                                                             |
|-----------------|-----------------------------------------------------------------------------------------|
| ActiviyPub      | Implements webfinger and activityPub actor and exposing pages as activitypub outbox     |
| Autolink        | Shorten a link string so it wouldn't take unnecessary space                             |
| Autolink pages  | Convert a page name mentions in the middle of text to a link                            |
| Custom CSS      | Allow to add custom CSS file to the head of the page                                    |
| Custom Widget   | Allow specifying content that is added in <head> tag, before or after the content       |
| Date            | Detects dates and converts them to link to a page which lists all pages mentions it     |
| Disqus          | Add Disqus comments after the view page if -disqus flag is passed                       |
| Emoji           | Emoji autocomplete while editing                                                        |
| File operations | Add a tool item to delete and rename current page                                       |
| Github          | Adds "Edit on github" quick action                                                      |
| Hashtags        | Add support for hashtags #hashtag syntax                                                |
| HTML            | Considers HTML files as pages. supports (html, htm, xhtml)                              |
| Images          | Display consecutive images in columns beside each other instead of under each other     |
| Embed           | Adds a shortcode to embed one page in another page                                      |
| Link preview    | Preview tweets, Facebook posts, youtube videos, Giphy links                             |
| Manifest        | adds manifest.json to head tag and output proper JSON value.                            |
| MathJax         | Support MathJax syntax inline using $ and blocks using $$                               |
| Mermaid         | Support for MermaidJS graphing library                                                  |
| Opengraph       | Adds Opengraph meta tags for title, type, image                                         |
| Pandoc          | Use pandoc to render documents in other formats as pages like Org-mode files            |
| Photos          | lists images in a directory similar to instagram                                        |
| PGP             | PGP key ID to decrypt and edit .md.pgp files using gpg. if empty encryption will be off |
| RSS             | Provides RSS feed served under /+/feed.rss and added to the header of pages             |
| RTL             | Fixes text direction for RTL languages in the view page                                 |
| Recent          | Adds an item to footer to list all pages ordered by last modified page file.            |
| Search          | Full text search                                                                        |
| Shortcode       | adds a way for short codes (one line and block)                                         |
| Sitemap         | adds support for sitemap.xml for search engine crawling                                 |
| Star            | Star pages to pin them to footer                                                        |
| Todo            | allow toggle checkboxes while viewing the page without going to edit mode               |
| Upload file     | Add support for upload files, screenshots, audio and camera recording                   |
| Versions        | Keeps list of pages older versions                                                      |
| Editor          | Open the current page in your editor. it uses $EDITOR env variable                      |
