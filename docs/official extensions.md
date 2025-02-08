Defined under `/extensions` sub package. each extension is a subpackage. **All extensions are imported** by default in `cmd/xlog/xlog.go`.


| Extension       | Description                                                                             |
|-----------------|-----------------------------------------------------------------------------------------|
| ActiviyPub      | Implements webfinger and activityPub actor and exposing pages as activitypub outbox     |
| Autolink        | Shorten a link string so it wouldn't take unnecessary space                             |
| Autolink pages  | Convert a page name mentions in the middle of text to a link                            |
| blocks          | Allows the user to define custom blocks that uses YAML block of codes as input          |
| Custom Widget   | Allow specifying content that is added in <head> tag, before or after the content       |
| Date            | Detects dates and converts them to link to a page which lists all pages mentions it     |
| Disqus          | Add Disqus comments after the view page if -disqus flag is passed                       |
| Editor          | Open the current page in your editor. it uses $EDITOR env variable                      |
| Embed           | Adds a shortcode to embed one page in another page                                      |
| File operations | Add a tool item to delete and rename current page                                       |
| Frontmatter     | Allow YAML frontmatter. displayed as properties and can override page title             |
| Github          | Adds "Edit on github" quick action                                                      |
| PGP             | PGP key ID to decrypt and edit .md.pgp files using gpg. if empty encryption will be off |
| Hashtags        | Add support for hashtags #hashtag syntax                                                |
| Heading         | Render heading as a link                                                                |
| hotreload       | Reload current page when modified on disk                                               |
| HTML            | Considers HTML files as pages. supports (html, htm, xhtml)                              |
| Images          | Display consecutive images in columns beside each other instead of under each other     |
| Link preview    | Preview tweets, Facebook posts, youtube videos, Giphy links                             |
| Manifest        | adds manifest.json to head tag and output proper JSON value.                            |
| MathJax         | Support MathJax syntax inline using $ and blocks using $$                               |
| Mermaid         | Support for MermaidJS graphing library                                                  |
| Opengraph       | Adds Opengraph meta tags for title, type, image                                         |
| Pandoc          | Use pandoc to render documents in other formats as pages like Org-mode files            |
| Photos          | lists images in a directory similar to instagram                                        |
| Recent          | Adds an item to footer to list all pages ordered by last modified page file.            |
| RSS             | Provides RSS feed served under /+/feed.rss and added to the header of pages             |
| RTL             | Fixes text direction for RTL languages in the view page                                 |
| Search          | Full text search                                                                        |
| Shortcode       | adds a way for short codes (one line and block)                                         |
| Sitemap         | adds support for sitemap.xml for search engine crawling                                 |
| sql_table       | For long tables adds SQL query form                                                     |
| Star            | Star pages to pin them to footer                                                        |
| TOC             | Adds table of contents                                                                  |
| Todo            | allow toggle checkboxes while viewing the page without going to edit mode               |
| Upload file     | Add support for upload files, screenshots, audio and camera recording                   |
