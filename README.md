XLog
=========

<p align="center"><img width="256" src="assets/logo.png" /></p>

Local-first personal knowledge management application with focus on enriching markdown files and surfacing implicit links between pages.

## Features
- Uses any directory of markdown files
- Supports Github flavor markdown (GFM)
- One statically compiled binary
- Converts a page name to link automatically on render time
- Lists pages that links to the current one.
- `template.md` content will be used for new pages
- Support Hashtags `#hashtag`
- Support for nested directories (although not favored)
- Lists pages that uses same hashtags in a `See Also` section
- Full text search
- Keeps list of pages older versions
- Supports editing pages with source code
- Has a list of tools for:
  - Drop a file or use the tool to upload the file and include/append it to the current page
  - Record screen/window/tab
  - Screenshot
  - Record Camera + Audio
  - Record Audio only
  - Tools works in both edit and view modes
- Preview tweets, Facebook posts, youtube videos, Giphy links
- Has a system for short codes (one line and block)
- Shows task list (Done/Total tasks) beside page link (auto links)
- Star pages to pin them to sidebar
- Add image to the start of a page to make it a cover photo
- A web server with a very small footprint thanks to Go.
- Minimal third party dependencies
- The first Emoji used in the page will be considered the icon of the page and displayed beside the title
- Uploading the same file multiple times it will be saved once
- Checkout [index.md](index) for additional features

## Installation

```
go install github.com/emad-elsaid/xlog@latest
```

## Usage

```
Usage of xlog:
  -bind string
        IP and port to bind the web server to (default "127.0.0.1:3000")
  -source string
        Directory that will act as a storage (default ".")
```

Now you can access xlog with `localhost:3000`

## Generating static site

I used Xlog to generate [my personal blog](https://www.emadelsaid.com/). it uses github workflow to do that [here is an example](https://github.com/emad-elsaid/emad-elsaid.github.io/blob/master/.github/workflows/xlog.yml).



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
