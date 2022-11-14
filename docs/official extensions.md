Defined under `/extensions` sub package. each extension is a subpackage. **All extensions are imported** by default in `cmd/xlog/xlog.go`.


| Extension       | Description                                                                   |
|-----------------|-------------------------------------------------------------------------------|
| Autolink        | Shorten a link string so it wouldn't take unnecessary space                   |
| Autolink pages  | Convert a page name mentions in the middle of text to a link                  |
| Emoji           | Emoji autocomplete while editing                                              |
| File operations | Add a tool item to delete and rename current page                             |
| Hashtags        | Add support for hashtags #hashtag syntax                                      |
| Link preview    | Preview tweets, Facebook posts, youtube videos, Giphy links                   |
| Opengraph       | Adds Opengraph meta tags for title, type, image                               |
| Recent          | Adds an item to sidebar to list all pages ordered by last modified page file. |
| Search          | Full text search                                                              |
| Shortcode       | adds a way for short codes (one line and block)                               |
| Star            | Star pages to pin them to sidebar                                             |
| Upload file     | Add support for upload files, screenshots, audio and camera recording         |
| Versions        | Keeps list of pages older versions                                            |
| Manifest        | adds manifest.json to head tag and output proper JSON value.                  |
| Sitemap         | adds support for sitemap.xml for search engine crawling                       |
