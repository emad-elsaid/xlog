Defined under `/extensions` sub package. each extension is a subpackage. **All extensions are imported** by default in `cmd/xlog/xlog.go`.


| Extension       | Description                                                                         |
|-----------------|-------------------------------------------------------------------------------------|
| ActiviyPub      | Implements webfinger and activityPub actor and exposing pages as activitypub outbox |
| Autolink        | Shorten a link string so it wouldn't take unnecessary space                         |
| Autolink pages  | Convert a page name mentions in the middle of text to a link                        |
| Emoji           | Emoji autocomplete while editing                                                    |
| File operations | Add a tool item to delete and rename current page                                   |
| Github          | Adds "Edit on github" quick action                                                  |
| Hashtags        | Add support for hashtags #hashtag syntax                                            |
| Link preview    | Preview tweets, Facebook posts, youtube videos, Giphy links                         |
| Manifest        | adds manifest.json to head tag and output proper JSON value.                        |
| Opengraph       | Adds Opengraph meta tags for title, type, image                                     |
| Recent          | Adds an item to sidebar to list all pages ordered by last modified page file.       |
| Search          | Full text search                                                                    |
| Shortcode       | adds a way for short codes (one line and block)                                     |
| Sitemap         | adds support for sitemap.xml for search engine crawling                             |
| Star            | Star pages to pin them to sidebar                                                   |
| Upload file     | Add support for upload files, screenshots, audio and camera recording               |
| Versions        | Keeps list of pages older versions                                                  |
