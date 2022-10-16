XLog
=========

<p align="center"><img width="256" src="assets/logo.png" /></p>

Local-first personal knowledge management application with focus on enriching markdown files and surfacing implicit links between pages.

## Core Features

- Uses any directory of markdown files
- Supports Github flavor markdown (GFM)
- One statically compiled binary
- `template.md` content will be used for new pages
- Support for nested directories (although not favored)
- Has a list of tools defined by extensions
- Add image to the start of a page to make it a cover photo
- A web server with a very small footprint thanks to Go.
- Minimal third party dependencies
- The first Emoji used in the page will be considered the icon of the page and displayed beside the title
- Shows task list (Done/Total tasks) beside page link (auto links)

## Extensions

Extensions are defined under `/extensions` sub package. each extension is a subpackage. importing the package in `cmd/xlog/xlog.go` has the side effect of registering the extension hooks. removing the extension from the list of imports will removing the features it provides.

* autolink:
  -Shorten a link string so it wouldn't take unnecessary space
* autolink_pages:
  - Convert a page name mentions in the middle of text to a link
  - List pages that links to the current one in a section at the end of the page.
* emoji:
  - Emoji autocomplete while editing
* file_operations
  - Add a tool item to delete current page
  - Add a tool item to rename current page
* hashtags
  - Support Hashtags `#hashtag`.
  - Convert any `#hashtag` to a link to list all pages the uses the hashtag
  - Adds an item in the sidebar to list all hashtags
  - Adds a section after the page to list all pages that used the same hashtags
* link_preview
  - Preview tweets, Facebook posts, youtube videos, Giphy links
* opengraph
  - Adds Opengraph meta tags for title, type, image
* recent
  - Adds an item to sidebar to list all pages ordered by last modified page file.
* search
  - Full text search
  - Adds a searchbox to the top of the sidebar to search pages and make it easier to create a page from selected text.
* shortcode
  - adds a way for short codes (one line and block)
  - Defines functions that can be used to add more shortcodes
  - '/' in editor autocompletes from the list of defined shortcodes
* star
  - Star pages to pin them to sidebar
* upload_file
  - Drop a file or use the tool to upload the file and include/append it to the current page
  - Record screen/window/tab
  - Screenshot
  - Record Camera + Audio
  - Record Audio only
* versions
  - Keeps list of pages older versions


## Installation

```
go install github.com/emad-elsaid/xlog/cmd/xlog@latest
```

## Usage

```
Usage of xlog:
  -bind string
        IP and port to bind the web server to (default "127.0.0.1:3000")
  -build string
        Build all pages as static site in this directory
  -index string
        Index file name used as home page (default "index")
  -readonly
        Should xlog hide write operations, read-only means all write operations will be disabled
  -sidebar
        Should render sidebar. (default true)
  -sitename string
        Site name is the name that appears on the header beside the logo and in the title tag (default "XLOG")
  -source string
        Directory that will act as a storage (default current directory)
```

Now you can access xlog with `localhost:3000`

## Generating static site

I used Xlog to generate [my personal blog](https://www.emadelsaid.com/). it uses github workflow to do that [here is an example](https://github.com/emad-elsaid/emad-elsaid.github.io/blob/master/.github/workflows/xlog.yml).


## License

Xlog is released under [MIT license](LICENSE)

## Logo

[Cassette tape icons created by Good Ware - Flaticon](https://www.flaticon.com/free-icons/cassette-tape)

## Screenshots

![](/public/f583bb0cbdf12641666e6f10b26171f61caee3330a68eb825ecbf77eab0227bd.png)

![](/public/e070d8b44a069a6b7336d9be02da4be9020d381aabeb073aa3399df71ca0492b.png)

![](/public/5a182dba45298d4bbc837a6d719c5c194c00f0fce8be363e33f85e9f7f849903.png)

![](/public/7fb69b1749f666f57bcc4044496efd943a84e0f2330cdd724177b4a225baef38.png)

![](/public/bbf6f8be4d374a33338938aeefe70a2f77eb841af8e77b6b79e2659b872e3933.png)

![](/public/9e732a37a60ea4e75c66e5acee8eb493e58d69f52c7121970f4ca24a8f69c8bd.png)

![](/public/8428b0390a330ebc3a815a9efac7b8b6a7b9891f772e6bec5ab73c07d0b9bda4.png)

![](/public/47ef8360209b86693a93794abafb5cccde08499ced4bd5250b77af1e0fce70cd.png)
